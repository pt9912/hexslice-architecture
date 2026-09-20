// A driven adapter providing an in-memory implementation of the order-area
// outbound OrderRepository port. It depends on the application core; the core
// never depends on it.
package hexslice.adapters.driven.memory.order

import hexslice.hexagon.application.order.ports.outbound.OrderRepository
import hexslice.hexagon.domain.order.Order
import hexslice.hexagon.domain.order.OrderError
import hexslice.hexagon.domain.order.OrderId
import java.util.concurrent.ConcurrentHashMap

/**
 * A thread-safe in-memory [OrderRepository]. It stores immutable snapshots so
 * callers cannot mutate persisted state by accident — the same reason the Go twin
 * restores a copy on both save and load.
 */
class InMemoryOrderRepository : OrderRepository {
    private val orders = ConcurrentHashMap<OrderId, Order>()

    override fun save(order: Order) {
        orders[order.id] = snapshotOf(order)
    }

    override fun findById(id: OrderId): Order {
        val stored = orders[id] ?: throw OrderError.NotFound()
        return snapshotOf(stored)
    }

    private fun snapshotOf(order: Order): Order =
        Order.restore(order.id, order.customerId, order.lines(), order.status)
}
