# AJIRASCAN Deep Technical Audit

Audit date: 2026-09-16  
Scope: all 56 files tracked at repository HEAD: Go source, templates, CSS, tests, module manifests, sample files and generated PDFs. Environment-managed .git and .kilo content was not treated as product source. No application code/configuration was changed.

Evidence standard: a feature is described as present only if current source makes it reachable. README roadmap statements are not implementation evidence.

## 1. Executive Summary

AJIRASCAN is not yet the ATS/HR/CV-writing, interview, payment platform described in the request. The current implementation is a small, stateless Go proof of concept:

* a CLI reads two local files, performs deterministic text matching and can write a local PDF;
* a Go HTTP server exposes a landing page and an anonymous form accepting pasted CV and job-description text;
* a synchronous in-memory heuristic generates a score and a few simple findings.

There are no accounts, organizations, HR users, candidates, jobs, applications, database, object storage, uploads, AI/LLM provider, payment provider, email/SMS, calendar, worker, public API, deployment definition, CI/CD pipeline or monitoring integration in the tracked source.

The implementation cannot currently leak stored candidate records, enable a paid subscription, or make an automated hiring decision because it has no such records, capability or workflow. It is nevertheless **not production-ready for real CV data**. The web process accepts unrestricted text anonymously, has no TLS configuration, request limits, server timeouts, rate limiting, logging, security headers, privacy controls or operational controls.

The product must not be represented as AI-powered in its current form. It has no model/API call. It uses static word lists, literal phrase checks, simple token matching, fixed weights and canned rewrites. It can produce invalid scores above 100, and some displayed report values are fixed strings rather than results of checks.

Production decision: **NO-GO for public production or personal-data processing.** It is appropriate only as a local/demo prototype with a clear deterministic-heuristic disclosure.

### Top findings

| Severity | Finding | Evidence | Risk, impact and recommended fix |
|---|---|---|---|
| **CRITICAL** | Direct HTTP handler has no TLS, input-size limits, rate limits or timeout policy | cmd/web/main.go; internal/web/handlers.go | CV text may be exposed if directly deployed; large/slow/repeated requests can exhaust service resources. Put it behind managed TLS/HSTS/WAF, configure HTTP timeouts, cap request/form/text bytes and add quotas/queues. |
| **CRITICAL** | Required identity, authorization, privacy and operational foundations are absent | Entire tracked source; only three routes in cmd/web/main.go | The intended multi-user HR/payments platform has no way to establish identity, tenant ownership or permitted access. Build tenant-scoped identity/RBAC, persistence, auditing and privacy foundations first. |
| **HIGH** | CV-like sample/generated files are committed | sample_cv.txt, output.pdf, output_cv.pdf | A named career profile may be replicated to clones/backups/forks without evidence of consent. Verify provenance, use synthetic fixtures, remediate history only under an approved process, and scan for PII. |
| **HIGH** | POST analysis accepts arbitrary free text without bounds | internal/web/handlers.go calls FormValue; internal/ats/engine.go processes it synchronously | Availability abuse and memory/CPU pressure. Reject unsupported methods/types, bound input and concurrency, validate encoding/length and rate-limit. |
| **HIGH** | Privacy governance is absent for planned employment profiling | no notice, consent, DPIA, retention, export/deletion, processor or audit materials exist | A future launch would lack demonstrable lawful basis/transparency/rights support. Complete DPIA/legal/DPO review and implement controls before real CVs or rankings. |
| **MEDIUM** | Score/report quality is misleading | internal/ats/weighted_score.go; internal/ats/report_formatter.go; cmd/cli/main.go | Uncapped boost can exceed 100; grammar/contact/action-verb values are printed as fixed values. Cap and validate scores; show only measured evidence and limitations. |

## 2. Current Architecture

### What currently exists

cmd/web/main.go uses Go standard-library net/http. It registers:

* / and /analyze, both pointing to web.HomeHandler;
* /static/, a FileServer for the local static directory.

For POST, internal/web/handlers.go obtains form fields named cv and job. It rejects only empty fields by redirecting to /. It calls ats.Analyze and renders templates/index.html using html/template. Contextual escaping from html/template is a positive XSS control for fields rendered from the result.

The CLI in cmd/cli/main.go takes -cv and -job paths, reads them using os.ReadFile, runs the same analysis, prints a report and invokes ats.ExportCVToPDF only when an -out path is supplied.

internal/ats is a pure synchronous algorithm: normalize and tokenize; compare CV/job tokens; filter a static ignore list; check 20 fixed phrases; find headings by substring; infer job type from static indicators; apply static keyword weights; count numeric strings as achievements; and return canned suggestions.

~~~text
Browser -- GET /, POST /analyze --> Go net/http :8080
                                      |
                                      +--> web.HomeHandler
                                      |       |
                                      |       +--> ats.Analyze (in memory)
                                      |       +--> html/template response
                                      |
                                      +--> /static/ local CSS FileServer

Local operator --> CLI --> os.ReadFile --> ats.Analyze --> stdout
                                           |
                                           +--> gofpdf --> local output path

Database, object storage, auth, AI provider, payments, email, calendar,
queue, deployment and third-party services: NOT IMPLEMENTED/CONNECTED.
~~~

### Partially implemented

internal/fileparser/parser.go contains TXT, DOCX and PDF extraction helpers. A repository-wide reference search found no caller outside that package: it is not reachable from web or CLI. PDF generation is implemented only for local CLI output; there is no server download/storage flow.

### Missing and not verifiable

No service/API layer beyond Go packages, database repository, background worker, third-party integration or hosting code exists. A TLS proxy, WAF, cloud account, backups, environment secrets, branch protections and runtime operations are outside this repository and cannot be verified.

## 3. Technology Stack

| Technology | Version | Purpose/location | Production suitability and risk |
|---|---:|---|---|
| Go | module says 1.24.3; audit host ran 1.25.4 | go.mod and all Go source | Strong base, but no pinned build image or patch/update policy. |
| net/http | Go runtime | cmd/web/main.go, internal/web/handlers.go | Suitable only after explicit TLS/timeouts/recovery/hardening. |
| html/template | Go runtime | internal/web/handlers.go, templates/index.html | Contextual escaping helps XSS; no CSP or browser headers. |
| Server HTML/CSS/vanilla JS | no package version | templates/index.html, static/style.css | Appropriate for demo; no frontend test/accessibility/security build process. |
| github.com/jung-kurt/gofpdf | v1.16.2 | internal/ats/export_pdf.go | Needs SCA review; exports raw CV content to a local file. |
| github.com/nguyenthenguyen/docx | 2023-06-21 pseudo-version | unused internal/fileparser/parser.go | Do not expose to untrusted upload without isolation. |
| github.com/ledongthuc/pdf | 2025-05-11 pseudo-version | unused internal/fileparser/parser.go | Do not expose to untrusted upload without isolation. |
| Transitive modules | go.sum | barcode, gofpdi, x/image, x/text and others | No SBOM, scanner, dependency update policy or vulnerability evidence. |

No frontend framework, ORM, SQL/NoSQL driver, cloud SDK, auth/JWT library, payment SDK, mail SDK or AI SDK is present.

## 4. Database Architecture

### Actual implementation

There is **no database**. The audit found no driver, schema, migration, entity/model, repository, connection configuration, transaction, index, foreign key, constraint, cache or backup code. Form input and ats.Result exist only in request memory. They disappear after the response except for any external infrastructure logging, which cannot be verified.

~~~text
Current data relationships

Browser fields (cv, job) --> in-memory Go strings --> ats.Result --> HTML
CLI local files ---------> in-memory Go strings --> stdout / local PDF

Candidate, HR user, Organization, Job, Application, Subscription,
Payment and AuditEvent: [not modelled]
~~~

### Required future relationship boundary (not a current schema)

~~~text
Organization 1---* Membership *---1 User
Organization 1---* Job 1---* Application *---1 CandidateProfile ---* CVVersion
Application 1---* AnalysisRun ---* AnalysisExplanation
Organization 1---* Subscription 1---* PaymentEvent
User/Service 1---* AuditEvent
CVVersion 1---* StoredObject (private and tenant-scoped)
~~~

Future records need immutable IDs, organization_id on every tenant-owned row, database/query tenant enforcement, foreign keys, composite unique/index constraints, controlled state transitions and retention/deletion metadata. Use managed KMS encryption and keep raw documents in private object storage rather than database fields.

| Severity | Finding -> evidence -> risk -> impact -> fix |
|---|---|
| **INFORMATIONAL** | No database -> no schema/migration/import anywhere -> no current persistent-data exposure but planned records cannot exist -> design migrations, constraints, backup/restore and tenant rules before adding persistence. |
| **HIGH** | No tenant model -> anonymous form fields only -> future cross-company leakage is likely if scope is bolted on -> make organization scope an authorization and database invariant from first migration. |
| **HIGH** | No audit/security records -> no audit/log package/model -> future CV/payment access would be unattributable -> add append-only, PII-redacted audit events for reads/exports/admin/payment actions. |

## 5. API Architecture

There is no versioned, JSON, REST or GraphQL API. The complete current HTTP inventory is:

| Method | Route | Purpose | Auth/role | Input -> output | Database | Validation, rate limit and concern |
|---|---|---|---|---|---|---|
| GET and other non-POST verbs | / | Render landing page | Public | none -> HTML | none | Unexpected verbs are not rejected; no headers. |
| POST | / | Analyze pasted form text | Public | cv, job -> HTML result | in memory only | Empty check only; no body/field/type/UTF-8 limit or rate limit. |
| GET and other non-POST verbs | /analyze | Also renders landing page | Public | none -> HTML | none | Ambiguous route/method semantics. |
| POST | /analyze | Form action analysis | Public | cv, job -> HTML result | in memory only | Same availability/input risks. |
| GET/HEAD normally | /static/* | Serve static assets | Public | path -> file | local static directory | No explicit cache/security header policy. |

IDOR/BOLA is not currently applicable: there are no IDs, accounts, organizations or stored objects. That is not a future control. New candidate/job/document endpoints require centralized object and tenant authorization before lookup.

The application sets no CORS headers, so ordinary browser same-origin defaults apply. It has no CSRF token/origin check. There is no authenticated state today, but a third-party form can consume analysis resources; CSRF controls become mandatory with cookies or side effects.

## 6. Authentication & Authorization

No registration, login, password, hash, JWT, session, refresh token, logout, reset, MFA, email verification, RBAC, applicant, HR, administrator or organization model exists. The public handler never examines identity.

Can one company, HR user or applicant access another user's information? **There are no such users or persisted information in the current implementation. Isolation cannot be verified or credited.**

For a future product use mature OIDC/session management (or securely implemented Argon2id/bcrypt storage if passwords are retained), short-lived rotating/revocable sessions, MFA for privileged roles, secure reset/email verification, least-privilege roles and server-side authorization on every action. Enforce organization scope before record lookup, not by UI filtering; create negative cross-tenant/BOLA tests.

## 7. Personal Data & Privacy

### Data processed and trace

The form intentionally accepts unrestricted CV and job-description text. CVs can include identifiers, contact details, address, employment/education history, skills, references, nationality, birth date, photographs and potentially sensitive data. Engine output derives score, matched/missing keywords and phrases, job type, section findings and rewrites.

sample_cv.txt contains a named career profile and employment history. output.pdf and output_cv.pdf are tracked generated artifacts. This establishes a repository privacy risk, not consent or lawful basis.

~~~text
User --> browser textareas --> POST / or /analyze --> Go request memory
     --> normalize/tokenize/heuristic score --> escaped HTML --> user browser

CLI operator --> local paths --> os.ReadFile --> memory --> analysis --> stdout
                                                --> optional local PDF

Verified external AI/vendor transfer: none
Verified application database/object-store retention: none
~~~

| Area | Exists | Missing / remediation |
|---|---|---|
| Transport encryption | No TLS configuration; server listens on :8080. | Require HTTPS, redirect HTTP and HSTS. Proxy protection is unverified. |
| At-rest encryption | No application store; CV-like artifacts are tracked repository files. | Review source-control access/provenance; KMS encrypted future storage/backups. |
| Notice/lawful basis/consent | None. | Versioned privacy notice, purpose/lawful basis/consent evidence where used. |
| Minimization | Entire unbounded text is accepted. | Defined fields/limits, avoid unneeded identity/sensitive data. |
| Retention/deletion/export | None. | Lifecycle policy and verified delete/export workflows across objects, indexes and backups. |
| Audit/cross-border/backups | None in source. | Processor/vendor/region/transfer inventory, DPA and audited access controls. |

This is not legal advice. The planned employment profiling/ranking may be high-risk processing. Kenya's official Data Protection Act provides for a DPIA before processing likely to create high risk; the General Regulations list automated decision-making/profiling with significant effect, large-scale/sensitive data and innovative solutions as high-risk examples. See the official [Data Protection Act, 2019](https://www.odpc.go.ke/wp-content/uploads/2024/02/TheDataProtectionAct__No24of2019.pdf) and [General Regulations, 2021](https://www.odpc.go.ke/wp-content/uploads/2024/03/THE-DATA-PROTECTION-GENERAL-REGULATIONS-2021-1.pdf). Before real CVs at scale, employer ranking, video/biometric features, AI-vendor transfer or international hosting, complete a DPIA with Kenyan counsel/DPO review and consult relevant [ODPC guidance](https://www.odpc.go.ke/guidelines-2/).

## 8. File/CV Security

There is **no web file upload**. templates/index.html uses textareas; there is no file input, FormFile or parser call. The CLI reads any local path that its local operator may read and writes an optional PDF to an operator-selected path.

The dormant parser is not a safe upload pipeline:

* ReadTXT copies an unbounded reader and ignores the copy error.
* ReadDOCX and ReadPDF receive local paths, not safe server object handles.
* No parser function has a caller or tests.
* There are no type/signature checks, sizes/page/text limits, filename controls, malware scanning, quarantine, private storage, signed URLs, authorization, sandboxing, retention or deletion controls.

No malicious upload is currently web-reachable, and no unauthorized stored-CV access exists because nothing is stored. Enabling these parsers without isolation would create file parser, path traversal and resource exhaustion risk. Use generated keys, private quarantine, MIME plus magic-byte verification, malware scanning, isolated workers with CPU/memory/time/page/text ceilings and tenant-authorized signed downloads.

## 9. AI & ATS Architecture

| Feature | Actual implementation | Input -> processing -> output | Status |
|---|---|---|---|
| CV parsing | Pasted text only; dormant DOCX/PDF helpers | text -> tokens | Partial utility; no structured CV. |
| ATS scoring | Static weighted keyword comparison | CV/job -> MatchKeywords/WeightedScore -> integer | Exists; deterministic, not AI. |
| Phrase match | 20 fixed phrases | strings -> contains checks -> matched/missing | Exists; no semantics. |
| Category/skill analysis | Static category maps | CV tokens -> report | Exists, not rendered in current UI. |
| CV improvement | Three literal patterns/canned text | CV lines -> fixed rewrite | Exists, mostly CLI/PDF output. |
| Job type | Static indicators | job tokens -> category | Exists, heuristic. |
| Ranking/rejection/shortlisting | None | N/A | Missing; no multi-candidate data. |
| Cover letter, LinkedIn, job matching, interview prep/evaluation, JD generation | None | N/A | Missing. |
| LLM/model/API/prompt/vector store | None | N/A | Missing. |

No prompt injection, hallucination, model-output validation, provider PII transfer, model key or token-cost risk exists in current code because there is no model call. These are future requirements, not current safeguards.

| Severity | Finding -> evidence -> risk/impact -> fix |
|---|---|
| **HIGH** | UI claims AI-powered ATS intelligence -> templates/index.html; no AI dependency/client -> misleading representation -> correct disclosure or implement governed AI later. |
| **MEDIUM** | Score can exceed 100 -> ApplyJobContextBoost in internal/ats/weighted_score.go adds without cap -> invalid percentage/verdict -> clamp 0..100 and add boundary/property tests. |
| **MEDIUM** | Report prints checks not performed -> fixed grammar/contact/title/action-verb values in cmd/cli/main.go and internal/ats/report_formatter.go -> false confidence -> show computed evidence only. |
| **MEDIUM** | Keyword stuffing distorts scoring -> repeated job tokens/static weights in matcher/weights -> inaccurate outcomes -> deduplicate/validate against labelled data and version algorithm. |
| **MEDIUM** | Intended HR use has no fairness/explanation/human-review policy -> English/static lists in weights.go, phrases.go, job_classifier.go -> risk of invalid employment effects -> prohibit automated decisions; test bias, explain output and retain human review/appeal. |

No current code automatically rejects, shortlists or materially ranks candidates against each other. It produces a single CV-to-job heuristic score only.

## 10. Payment Architecture

No payment/subscription function is implemented. Source search found no provider SDK, checkout, webhook, signature verification, transaction ID, subscription/trial/invoice/refund state, payment secret or record. The Pricing link in templates/index.html is a non-functional # link.

Frontend manipulation cannot currently activate a subscription because no entitlement exists. It is not payment-ready. A future design needs a PCI-compliant hosted/tokenized provider, server-side product/amount mapping, raw-body timestamped signature verification, replay protection, idempotent unique provider event/payment IDs, verified payment before entitlement, durable pending/failed/refunded/cancelled/chargeback transitions, authorization and audit records. Never store CVV or raw card credentials.

## 11. Security Assessment

| OWASP area | Evidence-based assessment | Severity |
|---|---|---|
| Access control/IDOR | No objects/accounts/auth; intended product boundary is entirely absent. | **CRITICAL** for planned product. |
| Injection | No SQL/NoSQL, shell execution or dynamic query. Form text stays Go strings. | **LOW** current exposure. |
| XSS | html/template contextual escaping is positive; no user-controlled DOM sink seen. CSP, nosniff, frame/referrer/permissions headers absent. | **MEDIUM** hardening gap. |
| CSRF/CORS | No CSRF defense; no auth state now. No explicit CORS headers, so default same-origin behavior is safer than wildcard. | **MEDIUM** availability now; **HIGH** before stateful routes. |
| SSRF | No URL fetcher/outbound HTTP client/webhook receiver. | **INFORMATIONAL**. |
| Abuse/validation | No limiter/quota/queue; only empty checks; synchronous analysis. | **HIGH**. |
| Upload | No reachable upload; dormant parser unsafe to expose as-is. | **HIGH** if enabled. |
| Transport | Plain http.ListenAndServe on :8080. | **CRITICAL** if direct/public. |
| Secrets | Current HEAD pattern scan found no application secret/API key/password/token. Historical/dangling-object absence is not verified. | **INFORMATIONAL**; add secret scan. |
| Dependencies | Direct third-party packages; no SCA/SBOM. govulncheck was unavailable in audit environment, so dependency vulnerability status is unverified. | **MEDIUM** process gap. |
| Errors/logging | Template errors use err.Error in HTTP response; startup/CLI failures panic; no logs/metrics/audit/recovery. | **MEDIUM/HIGH**. |
| Availability | No server timeouts/body limits; work is unbounded. | **HIGH**. |

## 12. DevOps & Deployment

Verified: go.mod/go.sum exist. The audit ran go test ./... and go vet ./... successfully with an isolated temporary Go build cache. The initial default-cache test attempt was blocked by sandbox filesystem access before compilation, not by an application failure.

Missing: Dockerfile/compose, Kubernetes/Helm/Terraform, reverse proxy, Makefile, CI, release/deploy pipeline, environment template, secret manager, config validation, health/readiness, graceful shutdown, logs, metrics, tracing, alerting, SLOs, backup/restore, disaster recovery and dev/stage/prod separation.

Although the startup message says localhost, the listener address :8080 binds all interfaces. Actual network exposure and compensating infrastructure cannot be verified.

## 13. Testing Assessment

Existing tests: 14 small unit-test files under internal/ats and one under internal/text. They cover normalization/tokenization, keyword/phrase matching, weighted/skills/section scoring, synonyms, suggestions, category/frequency analysis and achievement counts.

Critical missing tests:

* HTTP methods, body/field limits, errors, headers, rate limits and slow/large request behavior;
* XSS/CSRF/security regression tests;
* hostile DOCX/PDF upload/parser isolation tests;
* authentication/session/MFA/RBAC/BOLA/cross-tenant tests;
* database migrations/constraints/transactions/RLS/backup restore tests;
* payment signature/replay/idempotency/refund tests;
* AI PII/prompt/output/bias/explanation/human-review tests;
* integration, E2E, browser accessibility, load/performance and DR tests;
* CI-enforced coverage, static analysis, dependency/SBOM, secret and PII scanning.

## 14. Data Flow Diagram

~~~text
Applicant browser -- CV/JD text --> Go HTTP server --> in-memory analysis --> HTML
        |                               |                   |
        |                               |                   +--> no database/object store
        |                               +----------------------> no AI/payment/email/calendar
        |
        +-- confidentiality relies on external TLS/hosting, which is unverified

CLI operator -- local files --> os.ReadFile --> ATS heuristic --> terminal
                                             +--> optional local PDF
~~~

## 15. Critical Vulnerabilities

1. **CRITICAL — insecure, unbounded public data processing.** Finding: plain HTTP server with no configured timeouts, size limit or rate control. Evidence: cmd/web/main.go and internal/web/handlers.go. Risk: CV disclosure in transit and denial of service. Impact: confidentiality and availability loss. Fix: managed TLS/HSTS/WAF, explicit timeouts, maximum bytes and bounded concurrency before public deployment.

2. **CRITICAL — missing platform authorization and data protections.** Finding: no identity, roles, tenant data model, persistence or audit. Evidence: complete tracked source; anonymous handler. Risk: inevitable unsafe feature implementation if candidate/HR/payment capability is added prematurely. Impact: unauthorized cross-tenant access/noncompliance. Fix: build threat-modelled identity, RBAC, tenant persistence and audits before feature work.

3. **HIGH — possible personal data committed to source.** Finding: named CV profile and generated PDFs are tracked. Evidence: sample_cv.txt, output.pdf, output_cv.pdf. Risk: uncontrolled replication. Impact: privacy incident. Fix: verify authorization, replace with synthetic samples and perform approved remediation.

4. **HIGH — misleading employment-analysis claims.** Finding: UI says AI-powered while report includes non-computed checks. Evidence: templates/index.html, report_formatter.go, cmd/cli/main.go. Risk: reliance on unreliable score. Impact: unfair career/hiring choices and trust/legal risk. Fix: correct claims, validate method and introduce explanations/human review before employment impact.

## 16. Missing Production Requirements

Security/privacy: TLS/HSTS, WAF/rate limiting, security headers/CSP, method/body validation, recovery/safe errors, SCA/SBOM/secret/PII scans, threat model, DPIA, notice, lawful basis, retention/deletion/export, breach response and vendor governance.

Identity/data: account lifecycle, MFA, RBAC, tenant database/constraints/RLS, private object storage, malware quarantine, audit trail, encrypted backups and restore tests.

Payments/communications: hosted/tokenized checkout, verified idempotent webhooks, server-side entitlements, refunds/cancellations, transactional notification and safe calendar OAuth.

Operations: reproducible build/container, managed config/secrets, isolated environments, CI/CD gates, health checks, logs/metrics/traces/alerts, capacity tests, SLOs and DR runbooks.

Product/AI: durable job/application workflow, human approval, validated scoring/versioning/explanation/appeal, fairness tests and, if LLMs are introduced, redaction/provider agreements/schema validation/prompt-injection/cost controls.

## 17. Recommended Architecture Improvements

Use a modular monolith first:

~~~text
Browser/mobile --> CDN + WAF + TLS ingress --> Go API
                                                |
         +-------------------- tenant auth/validation/rate policy
         +--> PostgreSQL: tenant-scoped records, migrations, RLS
         +--> private object storage: KMS, quarantine, malware scan, signed URLs
         +--> worker queue: isolated parsing/analysis/email
         +--> audit/log/metric pipeline
         +--> identity, payment and calendar providers
         +--> approved AI gateway: minimization, schemas, budgets, human review
~~~

Separate raw documents from derived analysis. Enforce organization policy in every query before lookup. Make parsing/AI jobs isolated, retryable and idempotent. Create data classification and retention policy before selecting vendors/storage.

## 18. Production Readiness Assessment

| Domain | Assessment | Reason |
|---|---|---|
| Local demo | Limited ready | Basic CLI/web proof of concept; tests and vet pass. |
| Internet-facing analysis | Not ready | No TLS configuration, bounds, rate controls, headers, monitoring or incident controls. |
| Candidate/HR SaaS | Not implemented/not ready | No identity, database, jobs, applications or tenant protection. |
| Privacy | Not ready | No privacy program/rights controls; committed CV-like data needs review. |
| AI/ATS decision support | Not ready | Heuristic only; no validation, disclosure, fairness or review controls. |
| Payments | Not implemented/not ready | No provider/lifecycle/webhook. |
| DevOps/reliability | Not ready | No deployment/CI/observability/backups/DR evidence. |

## 19. Prioritized Remediation Roadmap

### Phase 1 — Must Fix Before Production

1. Correct AI/ATS marketing and disclose deterministic limitations.
2. Verify lawful basis/consent for tracked CV-like artifacts; use synthetic fixtures and approved remediation.
3. Complete data inventory, threat model and DPIA with privacy/legal/DPO review before real CVs, ranking or vendor AI.
4. Deploy HTTPS-only with HSTS/WAF/rate limits, server timeouts, strict request/text bounds, safe errors/recovery and PII-redacted logs.
5. Build identity, MFA for privileged roles, RBAC, tenant isolation, audit records and secure storage/database before HR/payment features.
6. Do not expose DOCX/PDF parsing until private quarantine, type validation, malware scanning and isolated workers exist.
7. Do not launch payments until verified signed/idempotent server-side payment state activates entitlements. Never store CVV/raw card data.

### Phase 2 — Production Hardening

1. Add CI for formatting, tests, vet, static analysis, dependency/SBOM/license, secret and PII scans; patch dependencies.
2. Add migration/constraint/index/tenant authorization and negative BOLA suites.
3. Add health/readiness, redacted logs, metrics/traces/alerts, backups, restore and DR exercises.
4. Add HTTP/integration/E2E/load/upload/payment/privacy-lifecycle tests and enforce coverage.
5. Isolate dev/stage/prod with managed secrets and least-privilege IAM.

### Phase 3 — Improvements

1. Validate scoring with consented labelled data; cap scores; remove fabricated indicators; publish explanations and algorithm versions.
2. Add applicant-controlled CV/job/application and recruiter-review workflows only after privacy/RBAC foundations.
3. If an LLM is added, implement minimization/redaction, provider DPA, output schemas, injection controls, cost budgets and mandatory human review for employment-impacting use.
4. Improve accessibility, mobile UX, localization, transparent retention controls and maintainability.
5. Split services only when measured scale/team boundaries warrant it.
