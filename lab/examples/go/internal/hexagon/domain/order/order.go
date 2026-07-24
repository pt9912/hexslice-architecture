// Package order is the business core of the order area. It contains the Order
// aggregate and its value objects and business rules. It depends on nothing
// outside the domain: no frameworks, no databases, no adapters.
package order

// Status represents the lifecycle state of an Order.
type Status string

// Order lifecycle states.
const (
	StatusCreated   Status = "CREATED"
	StatusCancelled Status = "CANCELLED"
)

// Order is the aggregate root of the order business area.
type Order struct {
	id         ID
	customerID CustomerID
	lines      []Line
	status     Status
}

// New creates a new order in the CREATED state. It enforces two invariants: an
// order must contain at least one line, and all lines must share one currency.
// An invalid order is never constructed, so it can never be persisted.
func New(id ID, customerID CustomerID, lines []Line) (*Order, error) {
	if len(lines) == 0 {
		return nil, ErrNoLines
	}
	currency := lines[0].unitPrice.currency
	for _, l := range lines[1:] {
		if l.unitPrice.currency != currency {
			return nil, ErrCurrencyMismatch
		}
	}
	return &Order{
		id:         id,
		customerID: customerID,
		lines:      append([]Line(nil), lines...),
		status:     StatusCreated,
	}, nil
}

// Restore rebuilds an order from persisted state without re-running creation
// invariants. Outbound adapters use it when loading stored orders.
func Restore(id ID, customerID CustomerID, lines []Line, status Status) *Order {
	return &Order{
		id:         id,
		customerID: customerID,
		lines:      append([]Line(nil), lines...),
		status:     status,
	}
}

// ID returns the order identifier.
func (o *Order) ID() ID { return o.id }

// CustomerID returns the owning customer.
func (o *Order) CustomerID() CustomerID { return o.customerID }

// Status returns the current lifecycle state.
func (o *Order) Status() Status { return o.status }

// Lines returns a copy of the order lines.
func (o *Order) Lines() []Line { return append([]Line(nil), o.lines...) }

// Cancel transitions the order into the CANCELLED state. It is a business rule
// that an already cancelled order cannot be cancelled again.
func (o *Order) Cancel() error {
	if o.status == StatusCancelled {
		return ErrAlreadyCancelled
	}
	o.status = StatusCancelled
	return nil
}

// Total sums the subtotals of all lines.
func (o *Order) Total() (Money, error) {
	if len(o.lines) == 0 {
		return Money{}, ErrNoLines
	}
	total := o.lines[0].Subtotal()
	for _, l := range o.lines[1:] {
		var err error
		total, err = total.Add(l.Subtotal())
		if err != nil {
			return Money{}, err
		}
	}
	return total, nil
}
