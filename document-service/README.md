# document-service

A service with one schema'd capability, `documents.extract`, completed
asynchronously: the invocation is accepted immediately and the result is
delivered when the work finishes.

    pagnet service create document-service   # prints the one-time credential once
    export PAGNET_SERVER=... PAGNET_CREDENTIAL=...
    go run .

Interesting property: async completion — the caller's invocation stays
pending (idempotency-key safe) until `inv.Complete` delivers the result, so
long work never holds a handler goroutine.
