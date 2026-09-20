// Command orderctl is the composition root of the Go HexSlice example. It wires
// driven adapters to the application slices and drives them from a driving CLI
// adapter. This is the only place where core and infrastructure meet.
package main

import (
	"context"
	"fmt"
	"io"
	"os"

	idgen "hexslice/example/internal/adapters/driven/id"
	memorder "hexslice/example/internal/adapters/driven/memory/order"
	"hexslice/example/internal/adapters/driven/notify"

	cliorder "hexslice/example/internal/adapters/driving/cli/order"

	"hexslice/example/internal/hexagon/application/order/cancelorder"
	"hexslice/example/internal/hexagon/application/order/createorder"
)

func main() {
	app := newCLI(os.Stdout, os.Stderr)

	if err := app.Run(context.Background(), os.Args[1:]); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

// newCLI is the wiring itself: driven adapters implement the outbound ports owned
// by the application core, the slices depend only on ports and implement their
// inbound port, and the driving adapter drives the use cases through those ports.
// That the handlers satisfy the ports is checked here, at the only place that
// knows both sides.
//
// The sinks are parameters so the end-to-end test drives exactly this wiring
// (see main_test.go) instead of a lookalike that could drift from it.
func newCLI(out, errOut io.Writer) *cliorder.CLI {
	orders := memorder.NewRepository()
	ids := idgen.NewGenerator()
	notifier := notify.NewWriter(out)

	createHandler := createorder.NewHandler(orders, ids)
	cancelHandler := cancelorder.NewHandler(orders, notifier)

	return cliorder.NewCLI(createHandler, cancelHandler, out, errOut)
}
