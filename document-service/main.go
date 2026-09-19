// Command document-service exposes documents.extract, completed
// asynchronously: the invocation is accepted immediately and the result is
// delivered when the work finishes.
package main

import (
	"context"
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
	svc.Handle("documents.extract", func(ctx context.Context, inv *sdk.Invocation) (any, error) {
		inv.Accept() // accept now; complete when the work is done
		go func() {
			var in ExtractInput
			if err := inv.Decode(&in); err != nil {
				log.Printf("decode: %v", err)
				return
			}
			text, err := extract(in.URI)
			if err != nil {
				log.Printf("extract: %v", err)
				return
			}
			inv.Complete(ctx, ExtractOutput{Text: text})
		}()
		return nil, nil
	})
	svc.Serve(ctx)
}

// extract stands in for real document processing.
func extract(uri string) (string, error) {
	return "extracted text from " + uri, nil
}
