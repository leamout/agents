# Customer Support

A support template for account lookup, case status, issue triage, case creation, and human escalation.

## Behavior

The agent should resolve straightforward questions using authoritative tool results, avoid guessing customer-specific information, and escalate when the request is sensitive, ambiguous, unsupported, or explicitly requires a person.

## Tools

`lookup_customer` identifies the customer record. `lookup_case` reads an existing case. `create_case` persists unresolved work. `transfer_call` hands the conversation to a human support queue.

## Deployment

Connect the webhook tools to the tenant's CRM or support platform, configure the escalation destination, then choose and bind either a composable or realtime AI path in Leamout.
