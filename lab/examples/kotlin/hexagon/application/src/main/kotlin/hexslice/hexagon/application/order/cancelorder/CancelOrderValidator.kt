package hexslice.hexagon.application.order.cancelorder

import hexslice.hexagon.application.order.cancelorder.ports.inbound.Request

/**
 * Violations of the cancel-order request shape — the Go twin's `ErrNoOrderID`.
 */
sealed class CancelOrderError(message: String) : RuntimeException(message) {
    class NoOrderId : CancelOrderError("cancelorder: order id is required")
}

/** Checks the request shape before the domain is involved. */
internal fun validate(request: Request) {
    if (request.orderId.isEmpty()) throw CancelOrderError.NoOrderId()
}
