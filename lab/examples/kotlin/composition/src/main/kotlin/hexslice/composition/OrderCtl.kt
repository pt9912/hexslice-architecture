// The composition root of the Kotlin HexSlice example. It wires driven adapters to
// the application slices and drives them from a driving CLI adapter. This is the
// only module where core and infrastructure meet — and the only one that *can*
// wire them, because it is the only module that sees both sides.
package hexslice.composition

import hexslice.adapters.driven.id.SequentialIdGenerator
import hexslice.adapters.driven.memory.order.InMemoryOrderRepository
import hexslice.adapters.driven.notify.StdoutNotifier
import hexslice.adapters.driving.cli.order.OrderCli
import hexslice.hexagon.application.order.cancelorder.CancelOrderHandler
import hexslice.hexagon.application.order.createorder.CreateOrderHandler
import kotlin.system.exitProcess

fun main(args: Array<String>) {
    val app = application(System.out, System.err)

    try {
        app.run(args.toList())
    } catch (error: RuntimeException) {
        System.err.println("error: ${error.message}")
        exitProcess(1)
    }
}

/**
 * The wiring itself: driven adapters implement the outbound ports owned by the
 * application core, the slices depend only on ports and implement their inbound
 * port, and the driving adapter drives the use cases through those ports. That the
 * handlers satisfy the ports is checked here, by the compiler — this is the only
 * place that knows both sides.
 *
 * The sinks are parameters so the end-to-end test drives exactly this wiring (see
 * `OrderCtlE2ETest`) instead of a lookalike that could drift from it.
 */
fun application(out: Appendable, errOut: Appendable): OrderCli {
    val orders = InMemoryOrderRepository()
    val ids = SequentialIdGenerator()
    val notifier = StdoutNotifier(out)

    val createHandler = CreateOrderHandler(orders, ids)
    val cancelHandler = CancelOrderHandler(orders, notifier)

    return OrderCli(createHandler, cancelHandler, out, errOut)
}
