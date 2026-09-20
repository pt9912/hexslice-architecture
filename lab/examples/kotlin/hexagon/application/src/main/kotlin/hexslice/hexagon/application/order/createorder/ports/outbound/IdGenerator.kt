// The outbound (driven) port that only the create-order slice needs.
//
// It lives inside the slice because no other use case requires it: the "as local
// as possible" level of a HexSlice port.
package hexslice.hexagon.application.order.createorder.ports.outbound

import hexslice.hexagon.domain.order.OrderId

/** Produces new, unique order identifiers. */
interface IdGenerator {
    fun newOrderId(): OrderId
}
