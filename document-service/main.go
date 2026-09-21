// Command document-service exposes documents.extract, completed
// asynchronously: the invocation is accepted immediately and the result is
// delivered when the work finishes.
package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/pagnet-code/pagnet/sdk"
)

type ExtractInput struct {
	URI string `json:"uri"`
}

type ExtractOutput struct {
	Text string `json:"text"`
}

func main() {
	ctx := context.Background()
	client, err := sdk.Connect(ctx, sdk.ConfigFromEnv())
	if err != nil {
		log.Fatal(err)
	}
	svc := client.Service("document-service")
	// Async: return (nil, sdk.ErrAsync) to defer completion; the work runs
	// in the background and completes the invocation when it is done.
	svc.Handle("documents.extract", func(ctx context.Context, inv *sdk.Invocation) (any, error) {
		go func() {
			_ = inv.Accept(ctx) // confirmation (the SDK already accepted on dispatch)
			var in ExtractInput
			if err := json.Unmarshal(inv.InputRaw, &in); err != nil {
				_ = inv.CompleteError(ctx, err)
				return
			}
			text, err := extract(in.URI)
			if err != nil {
				_ = inv.CompleteError(ctx, err)
				return
			}
			_ = inv.Complete(ctx, ExtractOutput{Text: text})
		}()
		return nil, sdk.ErrAsync
	})
	svc.Serve(ctx)
}

// extract stands in for real document processing.
func extract(uri string) (string, error) {
	return "extracted text from " + uri, nil
}
