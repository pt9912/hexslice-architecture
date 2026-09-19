// Package cancelorder implements the cancel-order use case as a vertical slice:
// its handler and validator, plus its use-case-local ports. The slice publishes
// itself through the inbound port and states its needs through the outbound
// port; the request and result types it works with are the inbound port's.
package cancelorder

import (
	"context"

	"hexslice/example/internal/hexagon/application/order/cancelorder/ports/inbound"
	"hexslice/example/internal/hexagon/application/order/cancelorder/ports/outbound"
	areaoutbound "hexslice/example/internal/hexagon/application/order/ports/outbound"
	"hexslice/example/internal/hexagon/domain/order"
)

// Handler implements the cancel-order inbound port. It depends only on ports.
type Handler struct {
	orders   areaoutbound.OrderRepository
	notifier outbound.Notifier
}

// NewHandler wires the cancel-order handler to its ports.
func NewHandler(orders areaoutbound.OrderRepository, notifier outbound.Notifier) *Handler {
	return &Handler{orders: orders, notifier: notifier}
}

// Cancel implements inbound.CancelOrder: it loads the order, applies the domain
// cancellation rule, persists the change, and notifies interested parties.
func (h *Handler) Cancel(ctx context.Context, req inbound.Request) (inbound.Result, error) {
	if err := Validate(req); err != nil {
		return inbound.Result{}, err
	}

	o, err := h.orders.FindByID(ctx, order.ID(req.OrderID))
	if err != nil {
		return inbound.Result{}, err
	}

	if err := o.Cancel(); err != nil {
		return inbound.Result{}, err
	}

	if err := h.orders.Save(ctx, o); err != nil {
		return inbound.Result{}, err
	}

	// Notification is a secondary side effect. The cancellation is already
	// committed, so a failing notifier must not turn a successful cancel into a
	// reported failure (which would also make a retry hit ErrAlreadyCancelled).
	// Best-effort: the error is deliberately ignored.
	_ = h.notifier.OrderCancelled(ctx, o.ID(), req.Reason)

	return inbound.Result{
		OrderID: string(o.ID()),
		Status:  string(o.Status()),
	}, nil
}
