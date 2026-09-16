# Security

Implemented: explicit HTTP timeouts, bounded body/form/text handling, methods/content-type checks, CSP and browser headers, configured-origin CORS, in-process IP rate limit, recovery, safe errors, Argon2id, hashed session primitives, tenant transaction scope, audit schema and file signature/type/size checks.

Production requires TLS ingress, trusted-proxy IP handling, managed secrets, strict TRUSTED_ORIGINS, private storage, a malware scanner, shared rate limiting, vulnerability/secret scans and least-privilege database roles. Never rely on UI state, client organization IDs, filenames/MIME headers, RLS alone or LLM output for authorization.
