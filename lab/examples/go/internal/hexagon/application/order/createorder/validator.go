package createorder

import (
	"errors"

	"hexslice/example/internal/hexagon/application/order/createorder/ports/inbound"
)

// Validation errors describe malformed create-order requests. They guard the
// use case before any domain object is constructed.
var (
	ErrNoCustomer = errors.New("createorder: customer id is required")
	ErrNoLines    = errors.New("createorder: at least one line is required")
)

// Validate checks the request shape before the domain is involved.
func Validate(req inbound.Request) error {
	if req.CustomerID == "" {
		return ErrNoCustomer
	}
	if len(req.Lines) == 0 {
		return ErrNoLines
	}
	return nil
}
