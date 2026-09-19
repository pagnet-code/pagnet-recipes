# custom-go-agent

The minimal external agent in Go: connects with the SDK, receives messages,
and exposes one capability (`agent.ping`). The agent's body is whatever you
put in the handlers — an LLM, deterministic code, or neither.

Create the participant (web: network → Add to Network → Connect agent; the
one-time credential is printed once), then:

    export PAGNET_SERVER=... PAGNET_CREDENTIAL=...
    go run .

Add the agent to a network from the network's member list.

Interesting property: Pagnet does not care what the agent contains — the
protocol only sees its capability and its messages.
