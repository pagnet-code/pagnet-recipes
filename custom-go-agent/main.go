// Command custom-go-agent is the minimal external agent: it connects with
// the Pagnet SDK, receives messages, and exposes one capability.
package main

import (
	"context"
	"log"

	"github.com/pagnet-code/pagnet/sdk"
)

type PingInput struct {
	Note string `json:"note"`
}

type PingOutput struct {
	Pong string `json:"pong"`
}

func main() {
	ctx := context.Background()
	client, err := sdk.Connect(ctx, sdk.ConfigFromEnv())
	if err != nil {
		log.Fatal(err)
	}
	agent := client.Agent("custom-go-agent")
	agent.OnMessage(func(ctx context.Context, m *sdk.Message) {
		log.Printf("message: %v", m)
	})
	agent.Handle("agent.ping", func(ctx context.Context, in PingInput) (PingOutput, error) {
		return PingOutput{Pong: "pong: " + in.Note}, nil
	})
	agent.Run(ctx)
}
