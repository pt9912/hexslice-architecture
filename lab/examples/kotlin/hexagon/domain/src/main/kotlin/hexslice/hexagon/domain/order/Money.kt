package hexslice.hexagon.domain.order

/**
 * A value object representing an amount in minor units (e.g. cents) together
 * with a 3-letter ISO currency code.
 *
 * Immutable: operations return new values instead of mutating the receiver, and
 * the constructor is private so a Money can only come into existence validated —
 * [of] is the only way in.
 */
@ConsistentCopyVisibility
data class Money private constructor(val amountMinor: Long, val currency: String) {
    /** Returns the sum of two Money values of the same currency. */
    fun add(other: Money): Money {
        if (currency != other.currency) throw OrderError.CurrencyMismatch()
        return Money(amountMinor + other.amountMinor, currency)
    }

    /** Renders the amount in major units, e.g. `19.99 EUR`. */
    override fun toString(): String =
        "${amountMinor / 100}.${(amountMinor % 100).toString().padStart(2, '0')} $currency"

    companion object {
        /** Creates a validated Money value. */
        fun of(amountMinor: Long, currency: String): Money {
            if (amountMinor < 0) throw OrderError.NegativeAmount()
            if (currency.length != 3) throw OrderError.InvalidCurrency()
            return Money(amountMinor, currency)
        }
    }
}
