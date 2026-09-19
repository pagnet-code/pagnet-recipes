# echo-service

The Pagnet Hello World: a deterministic service with one capability,
`echo.say`. Connect, handle, serve — nothing else.

    pagnet service create echo-service   # prints the one-time credential once
    export PAGNET_SERVER=... PAGNET_CREDENTIAL=...
    go run .

Add the service to a network (network → Add service).

Interesting property: the input/output JSON schemas in the manifest are
enforced by the SDK on both ends — before encryption out, after decryption in.
