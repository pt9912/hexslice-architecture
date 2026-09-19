package cancelorder

import (
	"errors"

	"hexslice/example/internal/hexagon/application/order/cancelorder/ports/inbound"
)

// ErrNoOrderID guards the cancel-order use case against empty input.
var ErrNoOrderID = errors.New("cancelorder: order id is required")

// Validate checks the request shape before the domain is involved.
func Validate(req inbound.Request) error {
	if req.OrderID == "" {
		return ErrNoOrderID
	}
	return nil
}
