package hexslice.hexagon.application.order.cancelorder

import hexslice.hexagon.application.order.cancelorder.ports.inbound.Request
import hexslice.hexagon.application.order.cancelorder.ports.outbound.Notifier
import hexslice.hexagon.application.order.ports.outbound.OrderRepository
import hexslice.hexagon.domain.order.CustomerId
import hexslice.hexagon.domain.order.Line
import hexslice.hexagon.domain.order.Money
import hexslice.hexagon.domain.order.Order
import hexslice.hexagon.domain.order.OrderError
import hexslice.hexagon.domain.order.OrderId
import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith
import kotlin.test.assertTrue

private class FakeOrderRepository(private val stored: Order? = null) : OrderRepository {
    override fun save(order: Order) = Unit

    override fun findById(id: OrderId): Order =
        stored?.takeIf { it.id == id } ?: throw OrderError.NotFound()
}

private class RecordingNotifier : Notifier {
    var called = false

    override fun orderCancelled(id: OrderId, reason: String) {
        called = true
    }
}

private fun seededOrder(): Order = Order.new(
    OrderId("ORD-1"),
    CustomerId("cust"),
    listOf(Line.of("SKU-1", 1, Money.of(1000, "EUR"))),
)

class CancelOrderHandlerTest {
    @Test
    fun `cancels the order and notifies`() {
        val notifier = RecordingNotifier()
        val handler = CancelOrderHandler(FakeOrderRepository(seededOrder()), notifier)

        val result = handler.cancel(Request(orderId = "ORD-1", reason = "test"))

        assertEquals("CANCELLED", result.status)
        assertTrue(notifier.called, "expected the notifier port to be invoked")
    }

    @Test
    fun `reports a missing order`() {
        val handler = CancelOrderHandler(FakeOrderRepository(), RecordingNotifier())

        assertFailsWith<OrderError.NotFound> {
            handler.cancel(Request(orderId = "missing", reason = ""))
        }
    }
}
