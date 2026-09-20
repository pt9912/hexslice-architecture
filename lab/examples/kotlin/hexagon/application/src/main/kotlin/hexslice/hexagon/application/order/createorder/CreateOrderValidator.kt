package hexslice.hexagon.application.order.createorder

import hexslice.hexagon.application.order.createorder.ports.inbound.Request

/**
 * Violations of the create-order request shape. They guard the use case before
 * any domain object is constructed — the Go twin's `ErrNoCustomer`/`ErrNoLines`.
 */
sealed class CreateOrderError(message: String) : RuntimeException(message) {
    class NoCustomer : CreateOrderError("createorder: customer id is required")

    class NoLines : CreateOrderError("createorder: at least one line is required")
}

/** Checks the request shape before the domain is involved. */
internal fun validate(request: Request) {
    if (request.customerId.isEmpty()) throw CreateOrderError.NoCustomer()
    if (request.lines.isEmpty()) throw CreateOrderError.NoLines()
}
