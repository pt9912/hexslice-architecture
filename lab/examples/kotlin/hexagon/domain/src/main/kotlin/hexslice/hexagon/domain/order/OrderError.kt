// Domain errors of the order business area. They are part of the domain
// vocabulary and carry no technical or infrastructure meaning — the Go twin
// returns the same set as sentinel errors, so a caller's `errors.Is` reads as
// `assertFailsWith` here.
package hexslice.hexagon.domain.order

sealed class OrderError(message: String) : RuntimeException(message) {
    class NegativeAmount : OrderError("order: money amount must not be negative")

    class InvalidCurrency : OrderError("order: currency must be a 3-letter ISO code")

    class CurrencyMismatch : OrderError("order: cannot combine different currencies")

    class EmptySku : OrderError("order: line SKU must not be empty")

    class InvalidQuantity : OrderError("order: line quantity must be greater than zero")

    class NoLines : OrderError("order: an order must contain at least one line")

    class AlreadyCancelled : OrderError("order: order is already cancelled")

    class NotFound : OrderError("order: order not found")
}
