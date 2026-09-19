// Command connector-example adapts an external HTTP API: a Pagnet
// invocation in, an HTTP call out, the result mapped back. Provider secrets
// come from this process's environment, never from Pagnet payloads.
package main

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/pagnet-code/pagnet/sdk"
)

type GetInput struct {
	Path string `json:"path"`
}

type GetOutput struct {
	Status int    `json:"status"`
	Body   string `json:"body"`
}

func main() {
	ctx := context.Background()
	client, err := sdk.Connect(ctx, sdk.ConfigFromEnv())
	if err != nil {
		log.Fatal(err)
	}
	svc := client.Service("rest-connector")
	svc.Handle("rest.get", func(ctx context.Context, in GetInput) (GetOutput, error) {
		return doGet(ctx, in.Path)
	})
	svc.Serve(ctx)
}

func doGet(ctx context.Context, path string) (GetOutput, error) {
	base := os.Getenv("CONNECTOR_BASE_URL")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, base+path, nil)
	if err != nil {
		return GetOutput{}, err
	}
	if key := os.Getenv("CONNECTOR_API_KEY"); key != "" {
		req.Header.Set("Authorization", "Bearer "+key)
	}
	resp, err := (&http.Client{Timeout: 30 * time.Second}).Do(req)
	if err != nil {
		return GetOutput{}, err
	}
	defer resp.Body.Close()
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		return GetOutput{}, err
	}
	return GetOutput{Status: resp.StatusCode, Body: buf.String()}, nil
}
