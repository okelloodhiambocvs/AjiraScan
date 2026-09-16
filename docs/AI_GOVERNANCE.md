# AI Governance

The current browser comparison is deterministic, not AI. internal/ai provides provider-independent bounded input, timeout, prompt-injection pattern rejection, output validation and disabled-provider behavior.

A real provider adapter must use controlled instructions, treat documents as untrusted data, minimize PII, require structured/schema-validated output, track provider/model/version, cap retries/cost, redact logs and apply business rules. AI may never directly authorize, reject, hire or make employment decisions; human review, evidence, limitations and appeal are required.
