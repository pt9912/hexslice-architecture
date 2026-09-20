package hexslice.hexagon.domain.order

/** The lifecycle state of an Order. Its names are the wire values, as in the Go twin. */
enum class Status { CREATED, CANCELLED }

/**
 * The aggregate root of the order business area.
 *
 * Instances are only obtainable through [new] (which enforces the creation
 * invariants) or [restore] (which trusts persisted state) — an invalid order can
 * never be constructed, so it can never be persisted.
 */
class Order private constructor(
    val id: OrderId,
    val customerId: CustomerId,
    private val lines: List<Line>,
    status: Status,
) {
    /** The current lifecycle state; only [cancel] may change it. */
    var status: Status = status
        private set

    /** Returns a copy of the order lines. */
    fun lines(): List<Line> = lines.toList()

    /**
     * Transitions the order into [Status.CANCELLED]. It is a business rule that an
     * already cancelled order cannot be cancelled again.
     */
    fun cancel() {
        if (status == Status.CANCELLED) throw OrderError.AlreadyCancelled()
        status = Status.CANCELLED
    }

    /** Sums the subtotals of all lines. */
    fun total(): Money {
        if (lines.isEmpty()) throw OrderError.NoLines()
        return lines.map { it.subtotal() }.reduce { sum, subtotal -> sum.add(subtotal) }
    }

    companion object {
        /**
         * Creates a new order in the [Status.CREATED] state, enforcing two
         * invariants: an order must contain at least one line, and all lines must
         * share one currency.
         */
        fun new(id: OrderId, customerId: CustomerId, lines: List<Line>): Order {
            if (lines.isEmpty()) throw OrderError.NoLines()
            val currency = lines.first().unitPrice.currency
            if (lines.any { it.unitPrice.currency != currency }) throw OrderError.CurrencyMismatch()
            return Order(id, customerId, lines.toList(), Status.CREATED)
        }

        /**
         * Rebuilds an order from persisted state without re-running creation
         * invariants. Driven adapters use it when loading stored orders.
         */
        fun restore(id: OrderId, customerId: CustomerId, lines: List<Line>, status: Status): Order =
            Order(id, customerId, lines.toList(), status)
    }
}
