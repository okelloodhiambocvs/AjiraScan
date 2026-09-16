CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS citext;

CREATE TYPE user_role AS ENUM ('PLATFORM_ADMIN', 'ORGANIZATION_ADMIN', 'HR_RECRUITER', 'HIRING_MANAGER', 'APPLICANT');
CREATE TYPE job_status AS ENUM ('DRAFT', 'OPEN', 'CLOSED', 'ARCHIVED');
CREATE TYPE application_status AS ENUM ('SUBMITTED', 'UNDER_REVIEW', 'SHORTLISTED', 'INTERVIEW', 'OFFERED', 'REJECTED', 'WITHDRAWN');
CREATE TYPE document_status AS ENUM ('QUARANTINED', 'SCANNING', 'SAFE', 'REJECTED', 'DELETED');
CREATE TYPE subscription_status AS ENUM ('TRIALING', 'ACTIVE', 'PAST_DUE', 'CANCELLED', 'EXPIRED');

CREATE TABLE users (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(), email CITEXT UNIQUE NOT NULL, password_hash TEXT NOT NULL,
 email_verified_at TIMESTAMPTZ, failed_login_count INTEGER NOT NULL DEFAULT 0 CHECK (failed_login_count >= 0),
 locked_until TIMESTAMPTZ, mfa_enabled BOOLEAN NOT NULL DEFAULT FALSE, mfa_secret_encrypted BYTEA,
 deleted_at TIMESTAMPTZ, created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE sessions (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(), user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
 token_hash BYTEA UNIQUE NOT NULL, expires_at TIMESTAMPTZ NOT NULL, revoked_at TIMESTAMPTZ,
 created_at TIMESTAMPTZ NOT NULL DEFAULT now(), last_seen_at TIMESTAMPTZ
);
CREATE INDEX sessions_active_idx ON sessions(user_id, expires_at) WHERE revoked_at IS NULL;
CREATE TABLE email_verification_tokens (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(), user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
 token_hash BYTEA UNIQUE NOT NULL, expires_at TIMESTAMPTZ NOT NULL, consumed_at TIMESTAMPTZ, created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE password_reset_tokens (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(), user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
 token_hash BYTEA UNIQUE NOT NULL, expires_at TIMESTAMPTZ NOT NULL, consumed_at TIMESTAMPTZ, created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE organizations (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(), name TEXT NOT NULL CHECK (char_length(name) BETWEEN 2 AND 200),
 slug TEXT UNIQUE NOT NULL, created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now(), deleted_at TIMESTAMPTZ
);
CREATE TABLE organization_members (
 organization_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
 user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE, role user_role NOT NULL,
 created_at TIMESTAMPTZ NOT NULL DEFAULT now(), PRIMARY KEY (organization_id, user_id)
);
CREATE INDEX organization_members_user_idx ON organization_members(user_id, organization_id);
CREATE TABLE applicant_profiles (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(), user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
 full_name TEXT, phone TEXT, location TEXT, headline TEXT, linkedin_content TEXT,
 created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now(), deleted_at TIMESTAMPTZ
);
CREATE TABLE jobs (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(), organization_id UUID NOT NULL REFERENCES organizations(id),
 created_by UUID NOT NULL REFERENCES users(id), title TEXT NOT NULL, description TEXT NOT NULL,
 requirements TEXT NOT NULL DEFAULT '', status job_status NOT NULL DEFAULT 'DRAFT',
 closes_at TIMESTAMPTZ, created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now(), deleted_at TIMESTAMPTZ
);
CREATE INDEX jobs_tenant_status_idx ON jobs(organization_id, status, created_at DESC) WHERE deleted_at IS NULL;
CREATE TABLE cv_documents (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(), owner_user_id UUID NOT NULL REFERENCES users(id),
 organization_id UUID REFERENCES organizations(id), object_key TEXT UNIQUE NOT NULL, original_filename TEXT NOT NULL,
 media_type TEXT NOT NULL, size_bytes BIGINT NOT NULL CHECK (size_bytes > 0), sha256 BYTEA NOT NULL,
 status document_status NOT NULL DEFAULT 'QUARANTINED', retention_until TIMESTAMPTZ, deleted_at TIMESTAMPTZ,
 created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX cv_documents_owner_idx ON cv_documents(owner_user_id, created_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX cv_documents_tenant_idx ON cv_documents(organization_id, created_at DESC) WHERE organization_id IS NOT NULL AND deleted_at IS NULL;
CREATE TABLE cv_versions (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(), document_id UUID NOT NULL REFERENCES cv_documents(id),
 version_number INTEGER NOT NULL CHECK (version_number > 0), extracted_text TEXT,
 parser_version TEXT, created_at TIMESTAMPTZ NOT NULL DEFAULT now(), UNIQUE(document_id, version_number)
);
CREATE TABLE cover_letters (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(), user_id UUID NOT NULL REFERENCES users(id), job_id UUID REFERENCES jobs(id),
 content TEXT NOT NULL, model_provider TEXT, model_name TEXT, created_at TIMESTAMPTZ NOT NULL DEFAULT now(), deleted_at TIMESTAMPTZ
);
CREATE TABLE applications (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(), organization_id UUID NOT NULL REFERENCES organizations(id),
 job_id UUID NOT NULL REFERENCES jobs(id), applicant_user_id UUID NOT NULL REFERENCES users(id),
 cv_version_id UUID REFERENCES cv_versions(id), cover_letter_id UUID REFERENCES cover_letters(id),
 status application_status NOT NULL DEFAULT 'SUBMITTED', human_review_required BOOLEAN NOT NULL DEFAULT TRUE,
 created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now(), deleted_at TIMESTAMPTZ,
 UNIQUE(job_id, applicant_user_id)
);
CREATE INDEX applications_tenant_status_idx ON applications(organization_id, status, created_at DESC) WHERE deleted_at IS NULL;
CREATE INDEX applications_applicant_idx ON applications(applicant_user_id, created_at DESC) WHERE deleted_at IS NULL;
CREATE TABLE application_status_history (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(), application_id UUID NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
 organization_id UUID NOT NULL REFERENCES organizations(id), from_status application_status, to_status application_status NOT NULL,
 changed_by UUID NOT NULL REFERENCES users(id), reason TEXT, created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX application_history_tenant_idx ON application_status_history(organization_id, application_id, created_at);
CREATE TABLE analysis_runs (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(), organization_id UUID REFERENCES organizations(id),
 applicant_user_id UUID NOT NULL REFERENCES users(id), job_id UUID REFERENCES jobs(id), cv_version_id UUID REFERENCES cv_versions(id),
 algorithm_version TEXT NOT NULL, model_provider TEXT, model_name TEXT, status TEXT NOT NULL CHECK (status IN ('PENDING','COMPLETED','FAILED')),
 created_at TIMESTAMPTZ NOT NULL DEFAULT now(), completed_at TIMESTAMPTZ, error_code TEXT
);
CREATE TABLE analysis_results (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(), analysis_run_id UUID NOT NULL UNIQUE REFERENCES analysis_runs(id) ON DELETE CASCADE,
 score INTEGER CHECK (score BETWEEN 0 AND 100), matched_requirements JSONB NOT NULL DEFAULT '[]', missing_requirements JSONB NOT NULL DEFAULT '[]',
 relevant_experience JSONB NOT NULL DEFAULT '[]', skills JSONB NOT NULL DEFAULT '[]', limitations JSONB NOT NULL DEFAULT '[]',
 created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE analysis_explanations (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(), analysis_run_id UUID NOT NULL REFERENCES analysis_runs(id) ON DELETE CASCADE,
 explanation TEXT NOT NULL, evidence JSONB NOT NULL DEFAULT '[]', created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE interviews (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(), organization_id UUID NOT NULL REFERENCES organizations(id),
 application_id UUID NOT NULL REFERENCES applications(id), starts_at TIMESTAMPTZ NOT NULL, ends_at TIMESTAMPTZ NOT NULL,
 timezone TEXT NOT NULL, status TEXT NOT NULL CHECK (status IN ('SCHEDULED','COMPLETED','CANCELLED','RESCHEDULED')),
 meeting_reference TEXT, created_by UUID NOT NULL REFERENCES users(id), created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
 CHECK (ends_at > starts_at)
);
CREATE INDEX interviews_tenant_time_idx ON interviews(organization_id, starts_at);
CREATE TABLE interview_participants (
 interview_id UUID NOT NULL REFERENCES interviews(id) ON DELETE CASCADE, user_id UUID NOT NULL REFERENCES users(id),
 participant_role TEXT NOT NULL CHECK (participant_role IN ('APPLICANT','INTERVIEWER','OBSERVER')), PRIMARY KEY(interview_id, user_id)
);
CREATE TABLE interview_feedback (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(), interview_id UUID NOT NULL REFERENCES interviews(id) ON DELETE CASCADE,
 organization_id UUID NOT NULL REFERENCES organizations(id), author_user_id UUID NOT NULL REFERENCES users(id),
 rating INTEGER CHECK (rating BETWEEN 1 AND 5), feedback TEXT NOT NULL, created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE subscription_plans (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(), code TEXT UNIQUE NOT NULL, name TEXT NOT NULL, amount_minor BIGINT NOT NULL CHECK (amount_minor >= 0),
 currency CHAR(3) NOT NULL, interval_months INTEGER NOT NULL CHECK (interval_months IN (1,12)), trial_days INTEGER NOT NULL DEFAULT 0 CHECK (trial_days >= 0),
 active BOOLEAN NOT NULL DEFAULT TRUE, created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE subscriptions (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(), organization_id UUID NOT NULL REFERENCES organizations(id),
 plan_id UUID NOT NULL REFERENCES subscription_plans(id), status subscription_status NOT NULL,
 provider TEXT NOT NULL, provider_subscription_id TEXT, current_period_start TIMESTAMPTZ, current_period_end TIMESTAMPTZ,
 cancel_at_period_end BOOLEAN NOT NULL DEFAULT FALSE, created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
 UNIQUE(provider, provider_subscription_id)
);
CREATE INDEX subscriptions_tenant_idx ON subscriptions(organization_id, status);
CREATE TABLE payment_transactions (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(), organization_id UUID NOT NULL REFERENCES organizations(id),
 subscription_id UUID REFERENCES subscriptions(id), provider TEXT NOT NULL, provider_transaction_id TEXT NOT NULL,
 amount_minor BIGINT NOT NULL CHECK (amount_minor >= 0), currency CHAR(3) NOT NULL, status TEXT NOT NULL,
 idempotency_key TEXT UNIQUE NOT NULL, verified_at TIMESTAMPTZ, created_at TIMESTAMPTZ NOT NULL DEFAULT now(), updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
 UNIQUE(provider, provider_transaction_id)
);
CREATE TABLE payment_events (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(), provider TEXT NOT NULL, provider_event_id TEXT NOT NULL,
 payload_sha256 BYTEA NOT NULL, received_at TIMESTAMPTZ NOT NULL DEFAULT now(), processed_at TIMESTAMPTZ,
 processing_error TEXT, UNIQUE(provider, provider_event_id)
);
CREATE TABLE notifications (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(), user_id UUID NOT NULL REFERENCES users(id), organization_id UUID REFERENCES organizations(id),
 channel TEXT NOT NULL, template_key TEXT NOT NULL, payload JSONB NOT NULL DEFAULT '{}', status TEXT NOT NULL DEFAULT 'PENDING',
 created_at TIMESTAMPTZ NOT NULL DEFAULT now(), sent_at TIMESTAMPTZ
);
CREATE TABLE consents (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(), user_id UUID NOT NULL REFERENCES users(id), purpose TEXT NOT NULL,
 policy_version TEXT NOT NULL, granted BOOLEAN NOT NULL, captured_at TIMESTAMPTZ NOT NULL DEFAULT now(), withdrawn_at TIMESTAMPTZ,
 UNIQUE(user_id, purpose, policy_version)
);
CREATE TABLE audit_events (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(), organization_id UUID REFERENCES organizations(id), actor_user_id UUID REFERENCES users(id),
 action TEXT NOT NULL, entity_type TEXT NOT NULL, entity_id UUID, metadata JSONB NOT NULL DEFAULT '{}',
 request_id TEXT, created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX audit_events_tenant_time_idx ON audit_events(organization_id, created_at DESC);
CREATE TABLE security_events (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(), user_id UUID REFERENCES users(id), event_type TEXT NOT NULL,
 source_ip INET, metadata JSONB NOT NULL DEFAULT '{}', created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE TABLE data_deletion_requests (
 id UUID PRIMARY KEY DEFAULT gen_random_uuid(), user_id UUID NOT NULL REFERENCES users(id), status TEXT NOT NULL CHECK (status IN ('PENDING','APPROVED','COMPLETED','REJECTED')),
 requested_at TIMESTAMPTZ NOT NULL DEFAULT now(), completed_at TIMESTAMPTZ, notes TEXT
);

ALTER TABLE jobs ENABLE ROW LEVEL SECURITY;
ALTER TABLE applications ENABLE ROW LEVEL SECURITY;
ALTER TABLE cv_documents ENABLE ROW LEVEL SECURITY;
ALTER TABLE interviews ENABLE ROW LEVEL SECURITY;
ALTER TABLE audit_events ENABLE ROW LEVEL SECURITY;

CREATE POLICY jobs_tenant_policy ON jobs USING (organization_id = NULLIF(current_setting('app.organization_id', true), '')::uuid);
CREATE POLICY applications_tenant_policy ON applications USING (organization_id = NULLIF(current_setting('app.organization_id', true), '')::uuid);
CREATE POLICY documents_tenant_policy ON cv_documents USING (organization_id = NULLIF(current_setting('app.organization_id', true), '')::uuid OR owner_user_id = NULLIF(current_setting('app.user_id', true), '')::uuid);
CREATE POLICY interviews_tenant_policy ON interviews USING (organization_id = NULLIF(current_setting('app.organization_id', true), '')::uuid);
CREATE POLICY audit_events_tenant_policy ON audit_events USING (organization_id = NULLIF(current_setting('app.organization_id', true), '')::uuid);
