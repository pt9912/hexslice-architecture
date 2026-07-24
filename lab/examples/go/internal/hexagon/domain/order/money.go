package order

import "fmt"

// Money is a value object representing an amount in minor units (e.g. cents)
// together with a 3-letter ISO currency code. It is immutable: operations
// return new values instead of mutating the receiver.
type Money struct {
	amountMinor int64
	currency    string
}

// NewMoney creates a validated Money value.
func NewMoney(amountMinor int64, currency string) (Money, error) {
	if amountMinor < 0 {
		return Money{}, ErrNegativeAmount
	}
	if len(currency) != 3 {
		return Money{}, ErrInvalidCurrency
	}
	return Money{amountMinor: amountMinor, currency: currency}, nil
}

// AmountMinor returns the amount in minor units.
func (m Money) AmountMinor() int64 { return m.amountMinor }

// Currency returns the ISO currency code.
func (m Money) Currency() string { return m.currency }

// Add returns the sum of two Money values of the same currency.
func (m Money) Add(other Money) (Money, error) {
	if m.currency != other.currency {
		return Money{}, ErrCurrencyMismatch
	}
	return Money{amountMinor: m.amountMinor + other.amountMinor, currency: m.currency}, nil
}

// String renders the amount in major units, e.g. "19.99 EUR".
func (m Money) String() string {
	return fmt.Sprintf("%d.%02d %s", m.amountMinor/100, m.amountMinor%100, m.currency)
}
