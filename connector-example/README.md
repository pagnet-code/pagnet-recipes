# connector-example

A service (`rest-connector`) that adapts an external HTTP API: a Pagnet
invocation in, an HTTP call out, the result mapped back. The provider is
configured with environment variables:

    export CONNECTOR_BASE_URL=https://api.example.com
    export CONNECTOR_API_KEY=***          # optional, sent as a Bearer token
    export PAGNET_SERVER=... PAGNET_CREDENTIAL=...
    go run .

Interesting property: the provider's secret lives only in the connector's
process environment — it never appears in a Pagnet payload, and network
members cannot read it.
