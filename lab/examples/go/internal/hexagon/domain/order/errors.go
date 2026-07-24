package order

import "errors"

// Domain errors express violations of the order business rules. They are part
// of the domain vocabulary and carry no technical/infrastructure meaning.
var (
	ErrNegativeAmount   = errors.New("order: money amount must not be negative")
	ErrInvalidCurrency  = errors.New("order: currency must be a 3-letter ISO code")
	ErrCurrencyMismatch = errors.New("order: cannot combine different currencies")
	ErrEmptySKU         = errors.New("order: line SKU must not be empty")
	ErrInvalidQuantity  = errors.New("order: line quantity must be greater than zero")
	ErrNoLines          = errors.New("order: an order must contain at least one line")
	ErrAlreadyCancelled = errors.New("order: order is already cancelled")
	ErrNotFound         = errors.New("order: order not found")
)
