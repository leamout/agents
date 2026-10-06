# Appointment Booking

A deployable scheduling package for checking availability, creating appointments, rescheduling existing bookings, and cancellations during staffed and after-hours periods.

## Package

```text
appointment-booking/
├── agents/
│   ├── appointment-booking.json
│   └── after-hours.json
├── routing.json
├── tools.json
└── README.md
```

`routing.json` uses the Agent manifest filenames without `.json` as local aliases: `appointment-booking` and `after-hours`.

## Agents

`appointment-booking` is the primary booking Agent during staffed service-desk hours.

`after-hours` keeps automated scheduling available outside staffed hours and uses the same authoritative scheduling tools without implying that a human desk is open.

## Routing

The included schedule treats Monday through Friday, 08:00–18:00, as staffed service-desk hours in the example timezone. All other times route to `after-hours`.

Replace `America/New_York` with the deployment's IANA timezone before production use.

## Tools

`check_availability` returns valid slots. `create_appointment` commits a new booking. `reschedule_appointment` changes an existing booking. `cancel_appointment` removes a confirmed appointment.

## Deployment inputs

Connect the webhook tools to the tenant's scheduling system, resolve provider credentials/integrations, and set the routing timezone expected by the business.
