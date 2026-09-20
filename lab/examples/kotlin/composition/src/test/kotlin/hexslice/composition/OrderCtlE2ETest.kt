package hexslice.composition

import hexslice.adapters.driving.cli.order.CliError
import hexslice.hexagon.domain.order.OrderError
import kotlin.test.Test
import kotlin.test.assertEquals
import kotlin.test.assertFailsWith
import kotlin.test.assertTrue

/**
 * Drives the real thing end to end: the composition root's own wiring, the real
 * adapters, the real slices and the real CLI. Only the process boundary — main's
 * exit code — is out of scope.
 *
 * It lives in this module, and that is forced rather than chosen: the driving
 * adapter's module does not have the driven adapters on its classpath, so a test
 * over there would not even compile. The composition root is the only place that
 * sees both sides — the same reason it exists.
 */
class OrderCtlE2ETest {
    @Test
    fun `the demo runs the whole lifecycle through the real wiring`() {
        val out = StringBuilder()
        val app = application(out, StringBuilder())

        app.run(listOf("demo"))

        assertEquals(
            listOf(
                "== HexSlice order demo ==",
                "created order ORD-000001 | total 47.48 EUR | status CREATED",
                "notification: order ORD-000001 cancelled (reason: customer changed their mind)",
                "cancelled order ORD-000001 | status CANCELLED",
                "second cancel rejected by domain: order: order is already cancelled",
            ),
            out.nonBlankLines(),
        )
    }

    @Test
    fun `create parses the flag surface and persists through the adapter`() {
        val out = StringBuilder()
        val app = application(out, StringBuilder())

        app.run(listOf("create", "-customer", "cust-1", "-line", "BOOK-1:2:1999:EUR"))

        assertEquals(
            listOf("created order ORD-000001 | total 39.98 EUR | status CREATED"),
            out.nonBlankLines(),
        )
    }

    @Test
    fun `cancelling an unknown order surfaces the domain error`() {
        val app = application(StringBuilder(), StringBuilder())

        assertFailsWith<OrderError.NotFound> {
            app.run(listOf("cancel", "-id", "does-not-exist"))
        }
    }

    @Test
    fun `a malformed flag is reported on stderr with the usage line`() {
        val errOut = StringBuilder()
        val app = application(StringBuilder(), errOut)

        assertFailsWith<CliError> {
            app.run(listOf("create", "-customer"))
        }
        assertTrue(
            errOut.toString().contains("usage: orderctl"),
            "expected the usage line on stderr, got \"$errOut\"",
        )
    }

    /** The non-empty output lines, so assertions are about content and order. */
    private fun StringBuilder.nonBlankLines(): List<String> = toString().lines().filter { it.isNotBlank() }
}
