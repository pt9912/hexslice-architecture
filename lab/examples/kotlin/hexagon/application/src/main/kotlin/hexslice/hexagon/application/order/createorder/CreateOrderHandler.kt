// The create-order use case as a vertical slice: its handler and validator, plus
// its use-case-local ports. The slice publishes itself through the inbound port
// and states its needs through the outbound port; the request and result types it
// works with are the inbound port's.
package hexslice.hexagon.application.order.createorder

import hexslice.hexagon.application.order.createorder.ports.inbound.CreateOrder
import hexslice.hexagon.application.order.createorder.ports.inbound.Line as RequestLine
import hexslice.hexagon.application.order.createorder.ports.inbound.Request
import hexslice.hexagon.application.order.createorder.ports.inbound.Result
import hexslice.hexagon.application.order.createorder.ports.outbound.IdGenerator
import hexslice.hexagon.application.order.ports.outbound.OrderRepository
import hexslice.hexagon.domain.order.CustomerId
import hexslice.hexagon.domain.order.Line
import hexslice.hexagon.domain.order.Money
import hexslice.hexagon.domain.order.Order

/**
 * Implements the create-order inbound port. It depends only on ports, never on
 * concrete adapters.
 *
 * One difference to the Go twin is the type system's, not the architecture's:
 * Kotlin is nominal, so implementing the port means naming it in the class
 * header — the slice imports its own inbound port. Go satisfies it structurally
 * and needs no import.
 *
 * The alias above is readability, not structure: the port's line and the domain
 * line are two different things, and the request carries the first.
 */
class CreateOrderHandler(
    private val orders: OrderRepository,
    private val ids: IdGenerator,
) : CreateOrder {
    /**
     * Validates the request, builds the order aggregate, persists it, and returns
     * a result describing the created order.
     */
    override fun create(request: Request): Result {
        validate(request)

        val lines = request.lines.map { line -> toDomain(line) }

        val order = Order.new(ids.newOrderId(), CustomerId(request.customerId), lines)
        orders.save(order)

        return Result(
            orderId = order.id.value,
            total = order.total().toString(),
            status = order.status.name,
        )
    }

    private fun toDomain(line: RequestLine): Line =
        Line.of(line.sku, line.quantity, Money.of(line.unitAmount, line.currency))
}
