package order

// Line is a value object describing one ordered item.
type Line struct {
	sku       string
	quantity  int
	unitPrice Money
}

// NewLine creates a validated order line.
func NewLine(sku string, quantity int, unitPrice Money) (Line, error) {
	if sku == "" {
		return Line{}, ErrEmptySKU
	}
	if quantity <= 0 {
		return Line{}, ErrInvalidQuantity
	}
	return Line{sku: sku, quantity: quantity, unitPrice: unitPrice}, nil
}

// SKU returns the stock-keeping unit of the line.
func (l Line) SKU() string { return l.sku }

// Quantity returns the ordered quantity.
func (l Line) Quantity() int { return l.quantity }

// UnitPrice returns the price of a single unit.
func (l Line) UnitPrice() Money { return l.unitPrice }

// Subtotal returns unitPrice * quantity.
func (l Line) Subtotal() Money {
	return Money{
		amountMinor: l.unitPrice.amountMinor * int64(l.quantity),
		currency:    l.unitPrice.currency,
	}
}
