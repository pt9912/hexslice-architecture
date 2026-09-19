// Package outbound holds the contracts the create-order slice needs from the
// outside world.
//
// IDGenerator is a use-case-local port: it lives inside the slice because no
// other use case requires it. This is the "as local as possible" level of a
// HexSlice port.
package outbound

import "hexslice/example/internal/hexagon/domain/order"

// IDGenerator produces new, unique order identifiers.
type IDGenerator interface {
	NewOrderID() order.ID
}
