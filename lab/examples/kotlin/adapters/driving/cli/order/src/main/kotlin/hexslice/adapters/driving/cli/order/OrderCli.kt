// The driving CLI adapter for the order business area. It translates command-line
// input into the slices' inbound-port requests and prints the results. It depends
// on the inbound ports only — never on the slices themselves, and never the other
// way around.
package hexslice.adapters.driving.cli.order

import hexslice.hexagon.application.order.cancelorder.ports.inbound.CancelOrder
import hexslice.hexagon.application.order.cancelorder.ports.inbound.Request as CancelRequest
import hexslice.hexagon.application.order.createorder.ports.inbound.CreateOrder
import hexslice.hexagon.application.order.createorder.ports.inbound.Line
import hexslice.hexagon.application.order.createorder.ports.inbound.Request as CreateRequest

/** Raised for command-line input the adapter itself rejects. */
class CliError(message: String) : RuntimeException(message)

/**
 * Drives the create-order and cancel-order use cases from the command line. Both
 * are held as inbound ports, so the adapter knows the contract but not the
 * implementation behind it.
 *
 * The two inbound ports each declare their own `Request`/`Result`, hence the
 * import aliases above — in the Go twin the package name disambiguated them.
 */
class OrderCli(
    private val create: CreateOrder,
    private val cancel: CancelOrder,
    private val out: Appendable,
    private val errOut: Appendable,
) {
    /** Dispatches a subcommand. With no arguments it runs the demo scenario. */
    fun run(args: List<String>) {
        if (args.isEmpty()) {
            runDemo()
            return
        }

        val subcommand = args.first()
        val rest = args.drop(1)
        try {
            when (subcommand) {
                "demo" -> runDemo()
                "create" -> runCreate(rest)
                "cancel" -> runCancel(rest)
                else -> throw CliError(
                    "unknown subcommand \"$subcommand\" (expected demo | create | cancel)",
                )
            }
        } catch (error: CliError) {
            // Flag and usage errors are the adapter's own. They go to stderr, and
            // the caller decides the exit code — as in the Go twin, where the flag
            // package wrote there and returned the error.
            errOut.appendLine("usage: orderctl $USAGE")
            throw error
        }
    }

    /** Runs a self-contained create -> cancel -> cancel-again scenario. */
    private fun runDemo() {
        appendLine("== HexSlice order demo ==")

        val created = create.create(
            CreateRequest(
                customerId = "cust-42",
                lines = listOf(
                    Line(sku = "BOOK-1", quantity = 2, unitAmount = 1999, currency = "EUR"),
                    Line(sku = "PEN-7", quantity = 5, unitAmount = 150, currency = "EUR"),
                ),
            ),
        )
        printf("created order %s | total %s | status %s\n", created.orderId, created.total, created.status)

        val cancelled = cancel.cancel(
            CancelRequest(orderId = created.orderId, reason = "customer changed their mind"),
        )
        printf("cancelled order %s | status %s\n", cancelled.orderId, cancelled.status)

        // Business rule: cancelling an already cancelled order is rejected by the
        // domain, not by the adapter.
        try {
            cancel.cancel(CancelRequest(orderId = created.orderId, reason = "again"))
        } catch (error: RuntimeException) {
            printf("second cancel rejected by domain: %s\n", error.message)
        }
    }

    private fun runCreate(args: List<String>) {
        val flags = Flags.parse(args)
        val lines = flags.values("line").map { parseLine(it) }

        val result = create.create(
            CreateRequest(
                customerId = flags.value("customer"),
                lines = lines,
            ),
        )
        printf("created order %s | total %s | status %s\n", result.orderId, result.total, result.status)
    }

    private fun runCancel(args: List<String>) {
        val flags = Flags.parse(args)

        val result = cancel.cancel(
            CancelRequest(
                orderId = flags.value("id"),
                reason = flags.value("reason"),
            ),
        )
        printf("cancelled order %s | status %s\n", result.orderId, result.status)
    }

    /** Parses one `-line sku:qty:amount:currency` value. */
    private fun parseLine(raw: String): Line {
        val parts = raw.split(":")
        if (parts.size != 4) {
            throw CliError("line must be sku:qty:amount:currency, got \"$raw\"")
        }
        val quantity = parts[1].toIntOrNull()
            ?: throw CliError("invalid quantity \"${parts[1]}\"")
        val amount = parts[2].toLongOrNull()
            ?: throw CliError("invalid amount \"${parts[2]}\"")
        return Line(sku = parts[0], quantity = quantity, unitAmount = amount, currency = parts[3])
    }

    private fun printf(format: String, vararg args: Any?) {
        out.appendLine(format.format(*args))
    }

    private fun appendLine(line: String) {
        out.appendLine(line)
    }

    /**
     * Collects repeated flags. Recognizes `-name value` and `-name=value`; the Go
     * twin's `flag` package accepts both, and so does the CLI surface here.
     */
    private class Flags(private val parsed: Map<String, List<String>>) {
        fun value(name: String): String = parsed[name].orEmpty().lastOrNull() ?: ""

        fun values(name: String): List<String> = parsed[name].orEmpty()

        companion object {
            fun parse(args: List<String>): Flags {
                val parsed = LinkedHashMap<String, MutableList<String>>()
                var index = 0
                while (index < args.size) {
                    val arg = args[index]
                    if (!arg.startsWith("-")) {
                        throw CliError("unexpected argument \"$arg\"")
                    }
                    val body = arg.trimStart('-')
                    val separator = body.indexOf('=')
                    if (separator >= 0) {
                        parsed.getOrPut(body.substring(0, separator)) { mutableListOf() }
                            .add(body.substring(separator + 1))
                        index++
                        continue
                    }
                    val next = args.getOrNull(index + 1)
                        ?: throw CliError("flag needs an argument: $arg")
                    parsed.getOrPut(body) { mutableListOf() }.add(next)
                    index += 2
                }
                return Flags(parsed)
            }
        }
    }

    private companion object {
        const val USAGE =
            "[demo | create -customer id -line sku:qty:amount:currency ... | cancel -id id -reason text]"
    }
}
