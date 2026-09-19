# coding-team

A software-development starter: three managed agent templates extracted out
of the product core — coder (implement/test, read-write workspace), reviewer
(analyze/review, read-only), coordinator (coordinate/delegate). Manifests
only, no code.

    pagnet recipe apply coding-team        # the whole team
    pagnet recipe apply coding-team/coder  # or a single template

Each template launches as a managed agent (runtime `qwen`) with its default
mission; add the launched agents to a network through the normal membership
UI.

Interesting property: these are plain data — the same manifests a user could
write by hand. Nothing about coding is special-cased in the protocol.
