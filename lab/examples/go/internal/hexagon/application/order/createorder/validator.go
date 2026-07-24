package createorder

import "errors"

// Validation errors describe malformed create-order commands. They guard the
// use case before any domain object is constructed.
var (
	ErrNoCustomer = errors.New("createorder: customer id is required")
	ErrNoLines    = errors.New("createorder: at least one line is required")
)

// Validate checks the command shape before the domain is involved.
func Validate(cmd Command) error {
	if cmd.CustomerID == "" {
		return ErrNoCustomer
	}
	if len(cmd.Lines) == 0 {
		return ErrNoLines
	}
	return nil
}
