# API

Current public endpoints:

| Method | Path | Behavior |
|---|---|---|
| GET | / | CV/job form |
| POST | / or /analyze | bounded URL-encoded comparison |
| GET | /healthz | liveness JSON |
| GET | /readyz | PostgreSQL readiness JSON |
| GET | /static/* | static assets |

Persisted HR/applicant/upload/billing/interview HTTP workflows are not yet exposed. No unauthenticated data-management endpoint exists.
