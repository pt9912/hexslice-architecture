package hexslice.hexagon.application.order.createorder

import hexslice.hexagon.application.order.createorder.ports.inbound.Line
import hexslice.hexagon.application.order.createorder.ports.inbound.Request
import hexslice.hexagon.application.order.createorder.ports.outbound.IdGenerator
import hexslice.hexagon.application.order.ports.outbound.OrderRepository
import hexslice.hexagon.domain.order.Order
import hexslice.hexagon.domain.order.OrderError
import hexslice.hexagon.domain.order.OrderId
import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith
import kotlin.test.assertNotNull

/** A hand-written stub for the outbound OrderRepository port. */
private class FakeOrderRepository : OrderRepository {
    var saved: Order? = null

    override fun save(order: Order) {
        saved = order
    }

    override fun findById(id: OrderId): Order =
        saved?.takeIf { it.id == id } ?: throw OrderError.NotFound()
}

/** A deterministic stub for the outbound IdGenerator port. */
private class FakeIdGenerator : IdGenerator {
    override fun newOrderId(): OrderId = OrderId("ORD-TEST")
}

class CreateOrderHandlerTest {
    @Test
    fun `creates an order and persists it`() {
        val repository = FakeOrderRepository()
        val handler = CreateOrderHandler(repository, FakeIdGenerator())

        val result = handler.create(
            Request(
                customerId = "cust-1",
                lines = listOf(Line(sku = "A", quantity = 2, unitAmount = 500, currency = "EUR")),
            ),
        )

        assertEquals("ORD-TEST", result.orderId)
        assertEquals("10.00 EUR", result.total)
        assertEquals("CREATED", result.status)
        assertNotNull(repository.saved, "expected the order to be persisted through the port")
    }

    @Test
    fun `rejects an empty request`() {
        val handler = CreateOrderHandler(FakeOrderRepository(), FakeIdGenerator())

        assertFailsWith<CreateOrderError.NoCustomer> {
            handler.create(Request(customerId = "", lines = emptyList()))
        }
    }
}
