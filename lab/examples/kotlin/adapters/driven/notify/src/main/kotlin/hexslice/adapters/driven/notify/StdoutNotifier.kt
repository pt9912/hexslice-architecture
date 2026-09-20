// A driven adapter implementing the cancel-order outbound Notifier port by
// writing human-readable messages to an Appendable.
package hexslice.adapters.driven.notify

import hexslice.hexagon.application.order.cancelorder.ports.outbound.Notifier
import hexslice.hexagon.domain.order.OrderId

/** Implements the Notifier port over a sink such as stdout. */
class StdoutNotifier(private val out: Appendable) : Notifier {
    /** Writes a cancellation notice. */
    override fun orderCancelled(id: OrderId, reason: String) {
        out.appendLine("notification: order ${id.value} cancelled (reason: $reason)")
    }
}
