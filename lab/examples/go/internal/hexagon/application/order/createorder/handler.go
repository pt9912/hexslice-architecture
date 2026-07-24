package createorder

import (
	"context"

	sliceports "hexslice/example/internal/hexagon/application/order/createorder/ports"
	areaports "hexslice/example/internal/hexagon/application/order/ports"
	"hexslice/example/internal/hexagon/domain/order"
)

// Handler executes the create-order use case. It depends only on ports, never
// on concrete adapters.
type Handler struct {
	orders areaports.OrderRepository
	ids    sliceports.IDGenerator
}

// NewHandler wires the create-order handler to its ports.
func NewHandler(orders areaports.OrderRepository, ids sliceports.IDGenerator) *Handler {
	return &Handler{orders: orders, ids: ids}
}

// Handle validates the command, builds the order aggregate, persists it, and
// returns a result describing the created order.
func (h *Handler) Handle(ctx context.Context, cmd Command) (Result, error) {
	if err := Validate(cmd); err != nil {
		return Result{}, err
	}

	lines := make([]order.Line, 0, len(cmd.Lines))
	for _, in := range cmd.Lines {
		price, err := order.NewMoney(in.UnitAmount, in.Currency)
		if err != nil {
			return Result{}, err
		}
		line, err := order.NewLine(in.SKU, in.Quantity, price)
		if err != nil {
			return Result{}, err
		}
		lines = append(lines, line)
	}

	o, err := order.New(h.ids.NewOrderID(), order.CustomerID(cmd.CustomerID), lines)
	if err != nil {
		return Result{}, err
	}

	if err := h.orders.Save(ctx, o); err != nil {
		return Result{}, err
	}

	total, err := o.Total()
	if err != nil {
		return Result{}, err
	}

	return Result{
		OrderID: string(o.ID()),
		Total:   total.String(),
		Status:  string(o.Status()),
	}, nil
}
