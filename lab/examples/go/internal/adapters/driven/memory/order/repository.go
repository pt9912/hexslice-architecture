// Package order is a driven adapter providing an in-memory implementation of
// the order-area outbound OrderRepository port. It depends on the application
// core; the core never depends on it.
package order

import (
	"context"
	"sync"

	domain "hexslice/example/internal/hexagon/domain/order"
)

// Repository is a thread-safe in-memory OrderRepository. It stores immutable
// snapshots so callers cannot mutate persisted state by accident.
type Repository struct {
	mu     sync.RWMutex
	orders map[domain.ID]*domain.Order
}

// NewRepository creates an empty in-memory repository.
func NewRepository() *Repository {
	return &Repository{orders: make(map[domain.ID]*domain.Order)}
}

// Save stores a snapshot of the order.
func (r *Repository) Save(_ context.Context, o *domain.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.orders[o.ID()] = domain.Restore(o.ID(), o.CustomerID(), o.Lines(), o.Status())
	return nil
}

// FindByID returns a snapshot of the stored order or domain.ErrNotFound.
func (r *Repository) FindByID(_ context.Context, id domain.ID) (*domain.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	o, ok := r.orders[id]
	if !ok {
		return nil, domain.ErrNotFound
	}
	return domain.Restore(o.ID(), o.CustomerID(), o.Lines(), o.Status()), nil
}
