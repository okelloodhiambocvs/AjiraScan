# Deployment

Build with Dockerfile, inject environment variables from a secret store, run migrations as a controlled release step, and expose only behind CDN/WAF/TLS ingress. Health is /healthz; readiness is /readyz.

Required external services: PostgreSQL, TLS/WAF, secret manager, private object storage, malware scanner, worker queue, email/calendar, AI/payment providers, monitoring/alerting, backups/PITR and restore drills. Set APP_ENV=production, DATABASE_URL, TRUSTED_ORIGINS and provider credentials only in deployment secrets.
