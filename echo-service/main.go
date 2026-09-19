// Command echo-service is the Pagnet Hello World: one service, one
// capability, no state.
package main

import (
	"context"
	"log"

	"github.com/pagnet-code/pagnet/sdk"
)

type SayInput struct {
	Name string `json:"name"`
}

type SayOutput struct {
	Message string `json:"message"`
}

func main() {
	ctx := context.Background()
	client, err := sdk.Connect(ctx, sdk.ConfigFromEnv())
	if err != nil {
		log.Fatal(err)
	}
	svc := client.Service("echo-service")
	svc.Handle("echo.say", func(ctx context.Context, in SayInput) (SayOutput, error) {
		return SayOutput{Message: "hello " + in.Name}, nil
	})
	svc.Serve(ctx)
}
