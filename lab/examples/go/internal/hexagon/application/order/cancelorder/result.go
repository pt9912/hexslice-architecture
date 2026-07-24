package cancelorder

// Result is the output of the cancel-order use case.
type Result struct {
	OrderID string
	Status  string
}
