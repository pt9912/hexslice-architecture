// Package cancelorder implements the cancel-order use case as a vertical slice:
// command, validator, handler, result, and its use-case-local ports.
package cancelorder

// Command is the input of the cancel-order use case.
type Command struct {
	OrderID string
	Reason  string
}
