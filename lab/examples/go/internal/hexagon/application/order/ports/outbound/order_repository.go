// Package outbound holds the order business area's outbound (driven) port.
//
// OrderRepository lives here rather than inside a single slice because both the
// create-order and cancel-order use cases genuinely need order persistence.
// This is the "as shared as necessary" level of a HexSlice port.
package outbound

import (
	"context"

	"hexslice/example/internal/hexagon/domain/order"
)

// OrderRepository is the persistence port for the order business area. It is
// owned by the application core; driven adapters implement it.
type OrderRepository interface {
	Save(ctx context.Context, o *order.Order) error
	FindByID(ctx context.Context, id order.ID) (*order.Order, error)
}
