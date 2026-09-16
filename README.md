# AJIRASCAN

AJIRASCAN is a Go modular-monolith foundation for a privacy-aware career and recruitment platform. The runnable UI is a deterministic CV-to-job comparison demo, not an AI hiring decision system. Do not use it to automatically reject, shortlist, hire, or make similarly significant employment decisions.

## Implemented foundation

* Deterministic ATS comparison and CLI analysis, with scores bounded to 0-100.
* PostgreSQL migration runner and initial schema for users, organizations, applicant profiles, jobs, applications, documents, analysis, interviews, subscriptions, payments, consents, audit/security events and deletion requests.
* Configuration, Argon2id password primitives, opaque hashed sessions, roles, tenant transaction scope and audit event writer.
* Security middleware for timeouts, bounded request bodies, response headers, configured-origin CORS, rate limiting, request IDs and panic recovery.
* Document validation/quarantine interfaces and disabled-by-default AI, billing, object-storage and malware-scanner adapters.

## Start locally

1. Copy `.env.example` to `.env`. Its values are development-only placeholders.
2. Start the local database: `docker compose up -d db`. The default host port is `5433` so it can coexist with an existing local PostgreSQL service.
3. Run migrations: `go run ./cmd/migrate`.
4. Run the server: `go run ./cmd/web`.
5. Open http://127.0.0.1:8080. Health: /healthz. Readiness: /readyz.

After signup or login, AjiraScan redirects to the protected `/dashboard`. The database must be running and migrated before account creation is available.

Checks:

    go test ./...
    go vet ./...

## Production boundary

Production requires TLS/WAF, PostgreSQL, managed secrets, private object storage, malware scanning, workers, email/calendar/AI/payment providers, monitoring, backups and external legal/privacy governance. No credentials are committed or invented. The sample CV is synthetic.

Read docs/ARCHITECTURE.md, docs/DATABASE.md, docs/API.md, docs/SECURITY.md, docs/PRIVACY.md, docs/AI_GOVERNANCE.md, docs/PAYMENTS.md, docs/DEPLOYMENT.md and docs/PRODUCTION_READINESS.md.
