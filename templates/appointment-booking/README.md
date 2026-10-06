# Appointment Booking

A scheduling template for checking availability, creating appointments, rescheduling existing bookings, and cancellations.

## Behavior

The agent must treat the scheduling system as authoritative. It should never promise a slot before an availability check or report a booking change before the write tool succeeds.

## Tools

`check_availability` returns valid slots. `create_appointment` commits a new booking. `reschedule_appointment` changes an existing booking. `cancel_appointment` removes a confirmed appointment.

## Deployment

Connect the webhook tools to the tenant's scheduling system, define the expected timezone behavior, then configure the AI engine and provider bindings in Leamout.
