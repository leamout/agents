# Sales Qualifier

A deployable inbound sales package for qualification, lead capture, follow-up scheduling, staffed-hours live handoff, and after-hours lead capture.

## Package

```text
sales-qualifier/
├── agents/
│   ├── sales-qualifier.json
│   └── after-hours.json
├── routing.json
├── tools.json
└── README.md
```

`routing.json` uses the Agent manifest filenames without `.json` as local aliases: `sales-qualifier` and `after-hours`.

## Agents

`sales-qualifier` handles staffed-hours qualification and may transfer a caller to a salesperson when appropriate.

`after-hours` captures qualification details and schedules follow-up but does not expose live transfer.

## Routing

The included schedule treats Monday through Friday, 09:00–17:00, as staffed sales hours in the example timezone. All other times route to `after-hours`.

Replace `America/New_York` with the deployment's IANA timezone before production use.

## Tools

`upsert_lead` records confirmed lead information. `schedule_followup` creates the next sales action. `transfer_call` performs a live handoff during staffed hours.

## Deployment inputs

Connect the webhook tools to the tenant CRM or scheduling service, configure the live sales transfer destination, resolve provider credentials/integrations, and set the routing timezone.
