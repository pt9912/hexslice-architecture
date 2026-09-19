// Package outbound holds the contracts the cancel-order slice needs from the
// outside world.
//
// Notifier is a use-case-local port: announcing a cancellation is a need
// specific to this slice, so the port lives inside it.
package outbound

import (
	"context"

	"hexslice/example/internal/hexagon/domain/order"
)

// Notifier announces that an order has been cancelled.
type Notifier interface {
	OrderCancelled(ctx context.Context, id order.ID, reason string) error
}
