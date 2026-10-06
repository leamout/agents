# Scheduled routing

This example shows the minimal Leamout `routing.json` shape for selecting an Agent when an inbound call begins.

```text
Inbound call
    │
    ├── after hours ──→ after_hours_agent
    │
    └── otherwise ────→ host_agent
```

`default_agent` is the fallback when no route matches. In this example the normal business-hours path is the default, so only the exceptional after-hours schedule needs an explicit route.

The `schedule` field references a named entry under `schedules`. Schedules use an IANA timezone and explicit day/time windows; the routing manifest intentionally does not define a string expression language.

Agent references such as `host_agent` and `after_hours_agent` are portable aliases. A deployment/import step resolves them to the actual Voice Agent resources for that tenant.

`routing.json` is only responsible for initial inbound-call selection. Once an Agent owns the call, conversation behavior, tools, transfers, hangup, and other call control remain the responsibility of the Agent Runtime and the Agent's configured tools.

The routing manifest therefore does not contain SIP destinations, intent transitions, workflow nodes, provider credentials, or tenant integration IDs.
