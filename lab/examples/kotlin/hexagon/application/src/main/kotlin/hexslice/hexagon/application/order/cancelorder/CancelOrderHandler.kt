// The cancel-order use case as a vertical slice: its handler and validator, plus
// its use-case-local ports. The slice publishes itself through the inbound port
// and states its needs through the outbound port; the request and result types it
// works with are the inbound port's.
package hexslice.hexagon.application.order.cancelorder

import hexslice.hexagon.application.order.cancelorder.ports.inbound.CancelOrder
import hexslice.hexagon.application.order.cancelorder.ports.inbound.Request
import hexslice.hexagon.application.order.cancelorder.ports.inbound.Result
import hexslice.hexagon.application.order.cancelorder.ports.outbound.Notifier
import hexslice.hexagon.application.order.ports.outbound.OrderRepository
import hexslice.hexagon.domain.order.OrderId

/**
 * Implements the cancel-order inbound port. It depends only on ports.
 */
class CancelOrderHandler(
    private val orders: OrderRepository,
    private val notifier: Notifier,
) : CancelOrder {
    /**
     * Loads the order, applies the domain cancellation rule, persists the change,
     * and notifies interested parties.
     */
    override fun cancel(request: Request): Result {
        validate(request)

        val order = orders.findById(OrderId(request.orderId))
        order.cancel()
        orders.save(order)

        // Notification is a secondary side effect. The cancellation is already
        // committed, so a failing notifier must not turn a successful cancel into
        // a reported failure (which would also make a retry hit AlreadyCancelled).
        // Best-effort: the failure is deliberately swallowed.
        runCatching { notifier.orderCancelled(order.id, request.reason) }

        return Result(orderId = order.id.value, status = order.status.name)
    }
}
