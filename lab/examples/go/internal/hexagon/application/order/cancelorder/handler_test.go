package cancelorder_test

import (
	"context"
	"errors"
	"testing"

	"hexslice/example/internal/hexagon/application/order/cancelorder"
	"hexslice/example/internal/hexagon/domain/order"
)

type fakeRepo struct {
	stored *order.Order
}

func (f *fakeRepo) Save(_ context.Context, o *order.Order) error {
	f.stored = o
	return nil
}

func (f *fakeRepo) FindByID(_ context.Context, id order.ID) (*order.Order, error) {
	if f.stored != nil && f.stored.ID() == id {
		return order.Restore(f.stored.ID(), f.stored.CustomerID(), f.stored.Lines(), f.stored.Status()), nil
	}
	return nil, order.ErrNotFound
}

type recordingNotifier struct {
	called bool
}

func (n *recordingNotifier) OrderCancelled(_ context.Context, _ order.ID, _ string) error {
	n.called = true
	return nil
}

func seededOrder(t *testing.T) *order.Order {
	t.Helper()
	price, err := order.NewMoney(1000, "EUR")
	if err != nil {
		t.Fatal(err)
	}
	line, err := order.NewLine("SKU-1", 1, price)
	if err != nil {
		t.Fatal(err)
	}
	o, err := order.New("ORD-1", "cust", []order.Line{line})
	if err != nil {
		t.Fatal(err)
	}
	return o
}

func TestCancelOrderSucceeds(t *testing.T) {
	repo := &fakeRepo{stored: seededOrder(t)}
	notifier := &recordingNotifier{}
	handler := cancelorder.NewHandler(repo, notifier)

	res, err := handler.Handle(context.Background(), cancelorder.Command{OrderID: "ORD-1", Reason: "test"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Status != "CANCELLED" {
		t.Fatalf("expected status CANCELLED, got %s", res.Status)
	}
	if !notifier.called {
		t.Fatal("expected the notifier port to be invoked")
	}
}

func TestCancelOrderNotFound(t *testing.T) {
	handler := cancelorder.NewHandler(&fakeRepo{}, &recordingNotifier{})
	_, err := handler.Handle(context.Background(), cancelorder.Command{OrderID: "missing"})
	if !errors.Is(err, order.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
