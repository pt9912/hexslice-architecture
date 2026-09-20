// A driven adapter implementing the create-order outbound IdGenerator port with
// sequential identifiers of the form ORD-000001.
package hexslice.adapters.driven.id

import hexslice.hexagon.application.order.createorder.ports.outbound.IdGenerator
import hexslice.hexagon.domain.order.OrderId
import java.util.concurrent.atomic.AtomicLong

/** Produces sequential, unique order ids. */
class SequentialIdGenerator : IdGenerator {
    private val counter = AtomicLong()

    /** Returns the next order id, starting at ORD-000001. */
    override fun newOrderId(): OrderId = OrderId("ORD-%06d".format(counter.incrementAndGet()))
}
