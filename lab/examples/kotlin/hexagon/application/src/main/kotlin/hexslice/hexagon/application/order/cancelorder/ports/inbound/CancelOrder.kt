// The inbound (driving) port of the cancel-order slice: the contract the outside
// world calls the use case with.
//
// As in every inbound port, the request and result types belong to the port, so a
// driving adapter never has to know the slice's package.
package hexslice.hexagon.application.order.cancelorder.ports.inbound

/** The input of the cancel-order use case. */
data class Request(
    val orderId: String,
    val reason: String,
)

/** The output of the cancel-order use case. */
data class Result(
    val orderId: String,
    val status: String,
)

/**
 * The inbound port of the cancel-order slice. A driving adapter calls it; the
 * slice's handler implements it.
 */
interface CancelOrder {
    fun cancel(request: Request): Result
}
