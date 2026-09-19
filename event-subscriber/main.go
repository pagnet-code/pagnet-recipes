// Command event-subscriber is a service that subscribes to test.* events
// and logs every delivery.
package main

import (
	"context"
	"log"

	"github.com/pagnet-code/pagnet/sdk"
)

func main() {
	ctx := context.Background()
	client, err := sdk.Connect(ctx, sdk.ConfigFromEnv())
	if err != nil {
		log.Fatal(err)
	}
	svc := client.Service("event-subscriber")
	svc.OnEvent("test.*", func(ctx context.Context, e *sdk.Event) error {
		log.Printf("event: %v", e)
		return nil // nil acks the delivery
	})
	svc.Serve(ctx)
}
