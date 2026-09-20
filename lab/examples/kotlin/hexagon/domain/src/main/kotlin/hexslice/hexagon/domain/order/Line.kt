package hexslice.hexagon.domain.order

/**
 * A value object describing one ordered item. Created through [of], which
 * enforces the line invariants; the Go twin's `NewLine` returns the same two
 * violations as errors.
 */
@ConsistentCopyVisibility
data class Line private constructor(
    val sku: String,
    val quantity: Int,
    val unitPrice: Money,
) {
    /** Returns `unitPrice * quantity`. */
    fun subtotal(): Money = Money.of(unitPrice.amountMinor * quantity, unitPrice.currency)

    companion object {
        /** Creates a validated order line. */
        fun of(sku: String, quantity: Int, unitPrice: Money): Line {
            if (sku.isEmpty()) throw OrderError.EmptySku()
            if (quantity <= 0) throw OrderError.InvalidQuantity()
            return Line(sku, quantity, unitPrice)
        }
    }
}
