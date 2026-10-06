# Receptionist

A general-purpose inbound receptionist for routing callers, answering simple front-desk requests, and taking messages.

## Behavior

The agent should identify the caller's intent before taking action, verify the intended employee or department before transferring, and fall back to message taking when a transfer cannot be completed.

## Tools

`lookup_employee` resolves people and departments. `transfer_call` uses Leamout's built-in call control. `take_message` persists a callback request or message through a tenant-owned webhook.

## Deployment

Choose either the composable or realtime engine in Leamout, configure the required provider bindings, create the tools from `tools.json`, and replace the webhook URL placeholders with tenant endpoints.
