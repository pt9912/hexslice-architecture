package createorder_test

import (
	"context"
	"testing"

	"hexslice/example/internal/hexagon/application/order/createorder"
	"hexslice/example/internal/hexagon/domain/order"
)

// fakeRepo is a hand-written stub for the OrderRepository port. It shows that
// the use case is testable without any real infrastructure.
type fakeRepo struct {
	saved *order.Order
}

func (f *fakeRepo) Save(_ context.Context, o *order.Order) error {
	f.saved = o
	return nil
}

func (f *fakeRepo) FindByID(_ context.Context, id order.ID) (*order.Order, error) {
	if f.saved != nil && f.saved.ID() == id {
		return f.saved, nil
	}
	return nil, order.ErrNotFound
}

// fakeIDs is a deterministic stub for the IDGenerator port.
type fakeIDs struct{}

func (fakeIDs) NewOrderID() order.ID { return "ORD-TEST" }

func TestCreateOrderSucceeds(t *testing.T) {
	repo := &fakeRepo{}
	handler := createorder.NewHandler(repo, fakeIDs{})

	res, err := handler.Handle(context.Background(), createorder.Command{
		CustomerID: "cust-1",
		Lines: []createorder.LineInput{
			{SKU: "A", Quantity: 2, UnitAmount: 500, Currency: "EUR"},
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.OrderID != "ORD-TEST" {
		t.Fatalf("expected generated id ORD-TEST, got %s", res.OrderID)
	}
	if res.Total != "10.00 EUR" {
		t.Fatalf("expected total 10.00 EUR, got %s", res.Total)
	}
	if res.Status != "CREATED" {
		t.Fatalf("expected status CREATED, got %s", res.Status)
	}
	if repo.saved == nil {
		t.Fatal("expected the order to be persisted through the port")
	}
}

func TestCreateOrderValidationFails(t *testing.T) {
	handler := createorder.NewHandler(&fakeRepo{}, fakeIDs{})
	if _, err := handler.Handle(context.Background(), createorder.Command{}); err == nil {
		t.Fatal("expected a validation error for an empty command")
	}
}
