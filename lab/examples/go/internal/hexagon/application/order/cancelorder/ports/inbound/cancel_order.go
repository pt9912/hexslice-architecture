// Package inbound holds the cancel-order slice's inbound (driving) port: the
// contract the outside world calls the use case with.
//
// As in every inbound port, the request and result types belong to the port,
// so a driving adapter never has to know the slice's package.
package inbound

import "context"

// Request is the input of the cancel-order use case.
type Request struct {
	OrderID string
	Reason  string
}

// Result is the output of the cancel-order use case.
type Result struct {
	OrderID string
	Status  string
}

// CancelOrder is the inbound port of the cancel-order slice. A driving adapter
// calls it; the slice's handler implements it.
type CancelOrder interface {
	Cancel(ctx context.Context, req Request) (Result, error)
}
