package cancelorder

import "errors"

// ErrNoOrderID guards the cancel-order use case against empty input.
var ErrNoOrderID = errors.New("cancelorder: order id is required")

// Validate checks the command shape before the domain is involved.
func Validate(cmd Command) error {
	if cmd.OrderID == "" {
		return ErrNoOrderID
	}
	return nil
}
