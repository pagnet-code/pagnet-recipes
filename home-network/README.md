# home-network

A home starter: a Home Assistant service (exposes `home.status`, publishes
`energy.*` and `lights.*` events) plus a family-representative agent
template.

    pagnet recipe apply home-network

After applying, launch the representative and give it the suggested
subscriptions (wake mode): `energy.*` and `lights.*` — it then reacts to
home events on the family's behalf.

Interesting property: the same model that runs a coding team runs a
household — services publish, the representative reacts, nothing is
home-specific.
