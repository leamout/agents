# Customer Support

A deployable support package for account lookup, case status, issue triage, case creation, staffed-hours escalation, and after-hours case capture.

## Package

```text
customer-support/
├── agents/
│   ├── customer-support.json
│   └── after-hours.json
├── routing.json
├── tools.json
└── README.md
```

`routing.json` uses the Agent manifest filenames without `.json` as local aliases: `customer-support` and `after-hours`.

## Agents

`customer-support` handles staffed-hours support and may use `transfer_call` when human escalation is appropriate.

`after-hours` performs safe self-service lookup and case creation but does not expose live transfer.

## Routing

The included schedule treats Monday through Friday, 08:00–18:00, as staffed support hours in the example timezone. All other times route to `after-hours`.

Replace `America/New_York` with the deployment's IANA timezone before production use.

## Tools

`lookup_customer` identifies the customer record. `lookup_case` reads an existing case. `create_case` persists unresolved work. `transfer_call` hands the conversation to a human support destination during staffed hours.

## Deployment inputs

Connect the webhook tools to the tenant's CRM or support platform, configure the staffed-hours escalation destination, resolve provider credentials/integrations, and set the routing timezone.
