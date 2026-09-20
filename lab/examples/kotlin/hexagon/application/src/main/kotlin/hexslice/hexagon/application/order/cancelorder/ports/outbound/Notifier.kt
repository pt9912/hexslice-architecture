// The outbound (driven) port that only the cancel-order slice needs.
//
// Announcing a cancellation is a need specific to this slice, so the port lives
// inside it.
package hexslice.hexagon.application.order.cancelorder.ports.outbound

import hexslice.hexagon.domain.order.OrderId

/** Announces that an order has been cancelled. */
interface Notifier {
    fun orderCancelled(id: OrderId, reason: String)
}
