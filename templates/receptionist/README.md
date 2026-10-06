# Receptionist

A deployable inbound receptionist package for routing callers, answering front-desk requests, transferring during staffed hours, and taking messages after hours.

## Package

```text
receptionist/
├── agents/
│   ├── receptionist.json
│   └── after-hours.json
├── routing.json
├── tools.json
└── README.md
```

`routing.json` uses the Agent manifest filenames without `.json` as local aliases: `receptionist` and `after-hours`.

## Agents

`receptionist` is the staffed-hours Agent. It can look up employees, transfer a confirmed caller, or take a message.

`after-hours` is the fallback Agent outside staffed hours. It can identify the intended employee or department and take a message, but it does not expose the live-transfer tool.

## Routing

The included schedule treats Monday through Friday, 09:00–17:00, as staffed hours in the example timezone. All other times route to `after-hours`.

Replace `America/New_York` with the deployment's IANA timezone before production use.

## Tools

`lookup_employee` resolves people and departments. `transfer_call` uses Leamout's built-in call control. `take_message` persists a callback request or message through a tenant-owned webhook.

## Deployment inputs

Resolve provider credentials/integrations, replace the webhook URL placeholders in `tools.json`, configure the live transfer destination, and set the routing timezone for the deployment.
