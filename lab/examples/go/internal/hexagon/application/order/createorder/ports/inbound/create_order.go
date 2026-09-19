// Package inbound holds the create-order slice's inbound (driving) port: the
// contract the outside world calls the use case with.
//
// The port owns the request and result types as well as the interface, so a
// driving adapter depends on the port alone and never on the slice. The slice
// implements the interface and imports its own port for the contract types;
// that import points inward (app -> ports), which is what keeps the direction
// of the dependency intact.
package inbound

import "context"

// Line is one requested order line, as handed in from the outside.
type Line struct {
	SKU        string
	Quantity   int
	UnitAmount int64 // amount in minor units, e.g. cents
	Currency   string
}

// Request is the input of the create-order use case. It carries raw,
// adapter-neutral data; the slice's handler turns it into domain objects.
type Request struct {
	CustomerID string
	Lines      []Line
}

// Result is the output of the create-order use case.
type Result struct {
	OrderID string
	Total   string
	Status  string
}

// CreateOrder is the inbound port of the create-order slice. A driving adapter
// (CLI, API controller, message consumer) calls it; the slice's handler
// implements it.
type CreateOrder interface {
	Create(ctx context.Context, req Request) (Result, error)
}
