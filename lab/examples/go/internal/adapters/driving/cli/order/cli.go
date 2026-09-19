// Package order is the driving CLI adapter for the order business area. It
// translates command-line input into the slices' inbound-port requests and
// prints the results. It depends on the inbound ports only — never on the
// slices themselves, and never the other way around.
package order

import (
	"context"
	"flag"
	"fmt"
	"io"
	"strconv"
	"strings"

	cancelin "hexslice/example/internal/hexagon/application/order/cancelorder/ports/inbound"
	createin "hexslice/example/internal/hexagon/application/order/createorder/ports/inbound"
)

// CLI drives the create-order and cancel-order use cases from the command line.
// Both are held as inbound ports, so the adapter knows the contract but not the
// implementation behind it.
type CLI struct {
	create createin.CreateOrder
	cancel cancelin.CancelOrder
	out    io.Writer // results (stdout)
	errOut io.Writer // diagnostics and flag/usage errors (stderr)
}

// NewCLI wires the driving adapter to the inbound ports of the two use cases.
// Results go to out (stdout); flag parsing and usage errors go to errOut
// (stderr). The composition root passes the slices' handlers here; they satisfy
// the ports structurally.
func NewCLI(create createin.CreateOrder, cancel cancelin.CancelOrder, out, errOut io.Writer) *CLI {
	return &CLI{create: create, cancel: cancel, out: out, errOut: errOut}
}

// printf writes formatted output to the adapter's sink. Write errors are
// ignored: a CLI cannot meaningfully recover from a failed stdout write.
func (c *CLI) printf(format string, args ...any) {
	_, _ = fmt.Fprintf(c.out, format, args...)
}

// println writes a single line to the adapter's sink, ignoring write errors.
func (c *CLI) println(s string) {
	_, _ = fmt.Fprintln(c.out, s)
}

// Run dispatches a subcommand. With no arguments it runs the demo scenario.
func (c *CLI) Run(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return c.runDemo(ctx)
	}
	switch args[0] {
	case "demo":
		return c.runDemo(ctx)
	case "create":
		return c.runCreate(ctx, args[1:])
	case "cancel":
		return c.runCancel(ctx, args[1:])
	default:
		return fmt.Errorf("unknown subcommand %q (expected demo | create | cancel)", args[0])
	}
}

// lineFlag collects repeated -line flags of the form sku:qty:amount:currency.
type lineFlag []createin.Line

func (l *lineFlag) String() string { return "" }

func (l *lineFlag) Set(v string) error {
	parts := strings.Split(v, ":")
	if len(parts) != 4 {
		return fmt.Errorf("line must be sku:qty:amount:currency, got %q", v)
	}
	qty, err := strconv.Atoi(parts[1])
	if err != nil {
		return fmt.Errorf("invalid quantity %q: %w", parts[1], err)
	}
	amount, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid amount %q: %w", parts[2], err)
	}
	*l = append(*l, createin.Line{
		SKU:        parts[0],
		Quantity:   qty,
		UnitAmount: amount,
		Currency:   parts[3],
	})
	return nil
}

func (c *CLI) runCreate(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("create", flag.ContinueOnError)
	fs.SetOutput(c.errOut)
	customer := fs.String("customer", "", "customer id")
	var lines lineFlag
	fs.Var(&lines, "line", "order line as sku:qty:amount(minor):currency (repeatable)")
	if err := fs.Parse(args); err != nil {
		return err
	}

	res, err := c.create.Create(ctx, createin.Request{
		CustomerID: *customer,
		Lines:      lines,
	})
	if err != nil {
		return err
	}
	c.printf("created order %s | total %s | status %s\n", res.OrderID, res.Total, res.Status)
	return nil
}

func (c *CLI) runCancel(ctx context.Context, args []string) error {
	fs := flag.NewFlagSet("cancel", flag.ContinueOnError)
	fs.SetOutput(c.errOut)
	id := fs.String("id", "", "order id")
	reason := fs.String("reason", "", "cancellation reason")
	if err := fs.Parse(args); err != nil {
		return err
	}

	res, err := c.cancel.Cancel(ctx, cancelin.Request{
		OrderID: *id,
		Reason:  *reason,
	})
	if err != nil {
		return err
	}
	c.printf("cancelled order %s | status %s\n", res.OrderID, res.Status)
	return nil
}

// runDemo runs a self-contained create -> cancel -> cancel-again scenario so
// the whole order lifecycle is visible in a single process (the in-memory
// repository does not persist across separate invocations).
func (c *CLI) runDemo(ctx context.Context) error {
	c.println("== HexSlice order demo ==")

	created, err := c.create.Create(ctx, createin.Request{
		CustomerID: "cust-42",
		Lines: []createin.Line{
			{SKU: "BOOK-1", Quantity: 2, UnitAmount: 1999, Currency: "EUR"},
			{SKU: "PEN-7", Quantity: 5, UnitAmount: 150, Currency: "EUR"},
		},
	})
	if err != nil {
		return err
	}
	c.printf("created order %s | total %s | status %s\n", created.OrderID, created.Total, created.Status)

	cancelled, err := c.cancel.Cancel(ctx, cancelin.Request{
		OrderID: created.OrderID,
		Reason:  "customer changed their mind",
	})
	if err != nil {
		return err
	}
	c.printf("cancelled order %s | status %s\n", cancelled.OrderID, cancelled.Status)

	// Business rule: cancelling an already cancelled order is rejected by the
	// domain, not by the adapter.
	if _, err := c.cancel.Cancel(ctx, cancelin.Request{OrderID: created.OrderID, Reason: "again"}); err != nil {
		c.printf("second cancel rejected by domain: %v\n", err)
	}
	return nil
}
