# event-subscriber

A service that subscribes to `test.*` and logs every delivery — the
publish/subscribe path in its smallest form. Publish with
`pagnet event publish test.hello --payload file` and watch it land.

    pagnet service create event-subscriber   # prints the one-time credential once
    export PAGNET_SERVER=... PAGNET_CREDENTIAL=...
    go run .

Interesting property: the handler's nil return acks the delivery; a non-nil
error makes the SDK retry, and duplicate deliveries are deduplicated by
delivery id.
