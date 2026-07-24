package cancelorder

import (
	"context"

	sliceports "hexslice/example/internal/hexagon/application/order/cancelorder/ports"
	areaports "hexslice/example/internal/hexagon/application/order/ports"
	"hexslice/example/internal/hexagon/domain/order"
)

// Handler executes the cancel-order use case. It depends only on ports.
type Handler struct {
	orders   areaports.OrderRepository
	notifier sliceports.Notifier
}

// NewHandler wires the cancel-order handler to its ports.
func NewHandler(orders areaports.OrderRepository, notifier sliceports.Notifier) *Handler {
	return &Handler{orders: orders, notifier: notifier}
}

// Handle loads the order, applies the domain cancellation rule, persists the
// change, and notifies interested parties.
func (h *Handler) Handle(ctx context.Context, cmd Command) (Result, error) {
	if err := Validate(cmd); err != nil {
		return Result{}, err
	}

	o, err := h.orders.FindByID(ctx, order.ID(cmd.OrderID))
	if err != nil {
		return Result{}, err
	}

	if err := o.Cancel(); err != nil {
		return Result{}, err
	}

	if err := h.orders.Save(ctx, o); err != nil {
		return Result{}, err
	}

	// Notification is a secondary side effect. The cancellation is already
	// committed, so a failing notifier must not turn a successful cancel into a
	// reported failure (which would also make a retry hit ErrAlreadyCancelled).
	// Best-effort: the error is deliberately ignored.
	_ = h.notifier.OrderCancelled(ctx, o.ID(), cmd.Reason)

	return Result{
		OrderID: string(o.ID()),
		Status:  string(o.Status()),
	}, nil
}
