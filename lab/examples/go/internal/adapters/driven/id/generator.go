// Package id is a driven adapter implementing the create-order outbound
// IDGenerator port with sequential identifiers of the form ORD-000001.
package id

import (
	"fmt"
	"sync/atomic"

	"hexslice/example/internal/hexagon/domain/order"
)

// Generator produces sequential, unique order ids.
type Generator struct {
	counter atomic.Uint64
}

// NewGenerator creates a fresh generator starting at ORD-000001.
func NewGenerator() *Generator { return &Generator{} }

// NewOrderID returns the next order id.
func (g *Generator) NewOrderID() order.ID {
	n := g.counter.Add(1)
	return order.ID(fmt.Sprintf("ORD-%06d", n))
}
