# Architecture

AJIRASCAN is a modular monolith. Go modules separate web, security, auth, tenancy, documents, ATS, AI, billing, interviews, privacy and audit concerns. Business modules use interfaces for external providers.

~~~text
Browser -> CDN/WAF/TLS -> Go HTTP -> security middleware -> modules
                                                     |-> PostgreSQL
                                                     |-> private object storage + scanner
                                                     |-> worker queue
                                                     |-> AI/payment/email/calendar adapters
~~~

Project tree: cmd/{web,cli,migrate}; internal/{ai,ats,audit,auth,billing,config,database,documents,interviews,privacy,security,tenancy,web}; migrations; templates; static; docs; .github/workflows.

Use tenancy.Begin to set transaction-local organization and user context for tenant queries. PostgreSQL RLS is defense in depth; server-side authorization remains mandatory.
