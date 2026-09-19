// Package createorder implements the create-order use case as a vertical slice:
// its handler and validator, plus its use-case-local ports. The slice publishes
// itself through the inbound port and states its needs through the outbound
// port; the request and result types it works with are the inbound port's.
package createorder

import (
	"context"

	"hexslice/example/internal/hexagon/application/order/createorder/ports/inbound"
	"hexslice/example/internal/hexagon/application/order/createorder/ports/outbound"
	areaoutbound "hexslice/example/internal/hexagon/application/order/ports/outbound"
	"hexslice/example/internal/hexagon/domain/order"
)

// Handler implements the create-order inbound port. It depends only on ports,
// never on concrete adapters.
type Handler struct {
	orders areaoutbound.OrderRepository
	ids    outbound.IDGenerator
}

// NewHandler wires the create-order handler to its ports.
func NewHandler(orders areaoutbound.OrderRepository, ids outbound.IDGenerator) *Handler {
	return &Handler{orders: orders, ids: ids}
}

// Create implements inbound.CreateOrder: it validates the request, builds the
// order aggregate, persists it, and returns a result describing the created
// order.
func (h *Handler) Create(ctx context.Context, req inbound.Request) (inbound.Result, error) {
	if err := Validate(req); err != nil {
		return inbound.Result{}, err
	}

	lines := make([]order.Line, 0, len(req.Lines))
	for _, in := range req.Lines {
		price, err := order.NewMoney(in.UnitAmount, in.Currency)
		if err != nil {
			return inbound.Result{}, err
		}
		line, err := order.NewLine(in.SKU, in.Quantity, price)
		if err != nil {
			return inbound.Result{}, err
		}
		lines = append(lines, line)
	}

	o, err := order.New(h.ids.NewOrderID(), order.CustomerID(req.CustomerID), lines)
	if err != nil {
		return inbound.Result{}, err
	}

	if err := h.orders.Save(ctx, o); err != nil {
		return inbound.Result{}, err
	}

	total, err := o.Total()
	if err != nil {
		return inbound.Result{}, err
	}

	return inbound.Result{
		OrderID: string(o.ID()),
		Total:   total.String(),
		Status:  string(o.Status()),
	}, nil
}
