# pagnet-recipes

Starter bundles for Pagnet networks: small, shareable sets of participant
manifests — and tiny Go examples — showing what a Pagnet network can be.
Nothing here is special-cased in the protocol; a recipe is just data that
`pagnet recipe apply` turns into normal API calls.

## Layout

| Directory | What it is |
|---|---|
| `coding-team/` | Recipe: coder, reviewer, coordinator agent templates |
| `custom-go-agent/` | Minimal external agent in Go (SDK) |
| `echo-service/` | Minimal deterministic service in Go (SDK) |
| `event-subscriber/` | Service subscribing to `test.*` events |
| `document-service/` | Service with a schema'd, asynchronously completed capability |
| `connector-example/` | Service adapting an external HTTP API |
| `home-network/` | Recipe: Home Assistant service + family representative template |
| `office-network/` | Recipe: sales + support templates, CRM/document service references |

## Applying a recipe

    pagnet recipe apply <path-or-url>

A Recipe manifest lists the sub-manifests to apply (`contains:`).
Interactively `pagnet recipe apply` shows a summary first;
`--non-interactive` applies deterministically.

## Go examples

Each example is its own module requiring the published SDK
(`github.com/pagnet-code/pagnet v0.3.0`) and reads `PAGNET_SERVER` +
`PAGNET_CREDENTIAL` from the environment: create the participant in the web
UI or CLI, copy the one-time credential once, export it, run.
