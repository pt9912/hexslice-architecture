// The outbound (driven) port of the order business area.
//
// OrderRepository lives here rather than inside a single slice because both the
// create-order and cancel-order use cases genuinely need order persistence.
// This is the "as shared as necessary" level of a HexSlice port.
package hexslice.hexagon.application.order.ports.outbound

import hexslice.hexagon.domain.order.Order
import hexslice.hexagon.domain.order.OrderId

/**
 * The persistence port for the order business area. It is owned by the
 * application core; driven adapters implement it.
 *
 * Where the Go twin returns `(Order, error)`, [findById] throws
 * [hexslice.hexagon.domain.order.OrderError.NotFound] — same contract, the
 * language's way of carrying the failure.
 */
interface OrderRepository {
    fun save(order: Order)

    fun findById(id: OrderId): Order
}
