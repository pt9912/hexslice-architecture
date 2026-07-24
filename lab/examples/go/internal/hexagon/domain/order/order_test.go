package order_test

import (
	"errors"
	"testing"

	"hexslice/example/internal/hexagon/domain/order"
)

func newLine(t *testing.T) order.Line {
	t.Helper()
	price, err := order.NewMoney(1000, "EUR")
	if err != nil {
		t.Fatal(err)
	}
	line, err := order.NewLine("SKU-1", 1, price)
	if err != nil {
		t.Fatal(err)
	}
	return line
}

func TestOrderRequiresAtLeastOneLine(t *testing.T) {
	if _, err := order.New("ORD-1", "cust", nil); !errors.Is(err, order.ErrNoLines) {
		t.Fatalf("expected ErrNoLines, got %v", err)
	}
}

func TestCancelTwiceIsRejected(t *testing.T) {
	o, err := order.New("ORD-1", "cust", []order.Line{newLine(t)})
	if err != nil {
		t.Fatal(err)
	}
	if err := o.Cancel(); err != nil {
		t.Fatalf("first cancel should succeed: %v", err)
	}
	if o.Status() != order.StatusCancelled {
		t.Fatalf("expected status CANCELLED, got %s", o.Status())
	}
	if err := o.Cancel(); !errors.Is(err, order.ErrAlreadyCancelled) {
		t.Fatalf("expected ErrAlreadyCancelled, got %v", err)
	}
}

func TestTotalSumsLineSubtotals(t *testing.T) {
	price, _ := order.NewMoney(1999, "EUR")
	line, _ := order.NewLine("BOOK-1", 2, price)
	o, err := order.New("ORD-1", "cust", []order.Line{line})
	if err != nil {
		t.Fatal(err)
	}
	total, err := o.Total()
	if err != nil {
		t.Fatal(err)
	}
	if total.String() != "39.98 EUR" {
		t.Fatalf("expected 39.98 EUR, got %s", total.String())
	}
}

func TestOrderRejectsMixedCurrencies(t *testing.T) {
	eur, _ := order.NewMoney(100, "EUR")
	usd, _ := order.NewMoney(100, "USD")
	eurLine, _ := order.NewLine("A", 1, eur)
	usdLine, _ := order.NewLine("B", 1, usd)

	if _, err := order.New("ORD-1", "cust", []order.Line{eurLine, usdLine}); !errors.Is(err, order.ErrCurrencyMismatch) {
		t.Fatalf("expected ErrCurrencyMismatch, got %v", err)
	}
}
