// The inbound (driving) port of the create-order slice: the contract the outside
// world calls the use case with.
//
// The port owns the request and result types as well as the interface, so a
// driving adapter depends on the port alone and never on the slice. The slice
// implements the interface and imports its own port for the contract types; that
// import points inward (app -> ports), which is what keeps the direction of the
// dependency intact.
package hexslice.hexagon.application.order.createorder.ports.inbound

/** One requested order line, as handed in from the outside. */
data class Line(
    val sku: String,
    val quantity: Int,
    /** Amount in minor units, e.g. cents. */
    val unitAmount: Long,
    val currency: String,
)

/**
 * The input of the create-order use case. It carries raw, adapter-neutral data;
 * the slice's handler turns it into domain objects.
 */
data class Request(
    val customerId: String,
    val lines: List<Line>,
)

/** The output of the create-order use case. */
data class Result(
    val orderId: String,
    val total: String,
    val status: String,
)

/**
 * The inbound port of the create-order slice. A driving adapter (CLI, API
 * controller, message consumer) calls it; the slice's handler implements it.
 */
interface CreateOrder {
    fun create(request: Request): Result
}
