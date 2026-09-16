# Database

Migration 001_initial.up.sql creates users, sessions, verification/reset tokens, organizations, organization_members, applicant_profiles, jobs, applications, application_status_history, cv_documents, cv_versions, cover_letters, analysis runs/results/explanations, interviews/participants/feedback, plans/subscriptions/payment transactions/events, notifications, consents, audit/security events and deletion requests.

The design uses UUIDs, foreign keys, unique constraints, timestamps, indexes, tenant organization_id columns and RLS on key tables. Run: go run ./cmd/migrate. Production must use a least-privilege role, encrypted managed backups/PITR and tested restores.
