# Sales Qualifier

An inbound sales template for understanding the caller's use case, capturing qualification details, scheduling follow-up, and optionally handing off to a salesperson.

## Behavior

The agent should collect useful context without turning the call into an interrogation. Product claims, pricing, discounts, and commercial commitments must come from authoritative business systems or a human salesperson rather than model inference.

## Tools

`upsert_lead` records confirmed lead information. `schedule_followup` creates the next sales action. `transfer_call` performs a live handoff when the caller or qualification flow calls for one.

## Deployment

Connect the webhook tools to the tenant CRM or scheduling service, configure a sales transfer destination if live handoff is enabled, then choose the Leamout execution engine and provider bindings.
