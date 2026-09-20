package hexslice.hexagon.domain.order

import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith

class OrderTest {
    private fun line(): Line = Line.of("SKU-1", 1, Money.of(1000, "EUR"))

    @Test
    fun `an order requires at least one line`() {
        assertFailsWith<OrderError.NoLines> {
            Order.new(OrderId("ORD-1"), CustomerId("cust"), emptyList())
        }
    }

    @Test
    fun `cancelling twice is rejected`() {
        val order = Order.new(OrderId("ORD-1"), CustomerId("cust"), listOf(line()))

        order.cancel()
        assertEquals(Status.CANCELLED, order.status)
        assertFailsWith<OrderError.AlreadyCancelled> { order.cancel() }
    }

    @Test
    fun `total sums the line subtotals`() {
        val order = Order.new(
            OrderId("ORD-1"),
            CustomerId("cust"),
            listOf(Line.of("BOOK-1", 2, Money.of(1999, "EUR"))),
        )

        assertEquals("39.98 EUR", order.total().toString())
    }

    @Test
    fun `an order rejects mixed currencies`() {
        val eur = Line.of("A", 1, Money.of(100, "EUR"))
        val usd = Line.of("B", 1, Money.of(100, "USD"))

        assertFailsWith<OrderError.CurrencyMismatch> {
            Order.new(OrderId("ORD-1"), CustomerId("cust"), listOf(eur, usd))
        }
    }
}
