// Command orderctl is the composition root of the Go HexSlice example. It wires
// outbound adapters to the application slices and drives them from an inbound
// CLI adapter. This is the only place where core and infrastructure meet.
package main

import (
	"context"
	"fmt"
	"os"

	idgen "hexslice/example/internal/adapters/outbound/id"
	memorder "hexslice/example/internal/adapters/outbound/memory/order"
	"hexslice/example/internal/adapters/outbound/notify"

	cliorder "hexslice/example/internal/adapters/inbound/cli/order"

	"hexslice/example/internal/hexagon/application/order/cancelorder"
	"hexslice/example/internal/hexagon/application/order/createorder"
)

func main() {
	// Outbound adapters implement the ports owned by the application core.
	orders := memorder.NewRepository()
	ids := idgen.NewGenerator()
	notifier := notify.NewWriter(os.Stdout)

	// Application slices depend only on ports.
	createHandler := createorder.NewHandler(orders, ids)
	cancelHandler := cancelorder.NewHandler(orders, notifier)

	// Inbound adapter drives the use cases.
	app := cliorder.NewCLI(createHandler, cancelHandler, os.Stdout, os.Stderr)

	if err := app.Run(context.Background(), os.Args[1:]); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
