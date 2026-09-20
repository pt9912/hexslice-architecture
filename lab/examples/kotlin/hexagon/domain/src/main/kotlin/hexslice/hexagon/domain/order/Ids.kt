// Identifiers of the order business area. Value classes, so an OrderId can never
// be passed where a CustomerId is expected — the type-level counterpart of the
// named string types in the Go twin.
package hexslice.hexagon.domain.order

/** Uniquely identifies an Order. */
@JvmInline
value class OrderId(val value: String)

/** Identifies the customer that owns an order. */
@JvmInline
value class CustomerId(val value: String)
