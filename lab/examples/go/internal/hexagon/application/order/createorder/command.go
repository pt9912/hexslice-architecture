// Package createorder implements the create-order use case as a vertical slice:
// command, validator, handler, result, and its use-case-local ports.
package createorder

// Command is the input of the create-order use case. It carries raw,
// adapter-neutral data; the handler turns it into domain objects.
type Command struct {
	CustomerID string
	Lines      []LineInput
}

// LineInput is the raw line data carried by the command.
type LineInput struct {
	SKU        string
	Quantity   int
	UnitAmount int64 // amount in minor units, e.g. cents
	Currency   string
}
