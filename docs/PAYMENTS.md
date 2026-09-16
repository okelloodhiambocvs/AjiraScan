# Payments

internal/billing defines verified webhook and idempotent event-store contracts. The schema holds plans, subscriptions, transactions and unique provider event IDs. No payment provider is configured.

A production adapter must create checkout server-side; verify raw signatures/replay window and transaction state/amount/currency; apply entitlements transactionally; handle cancellation, expiry, failure, refund and chargeback. Never accept client-side entitlement activation or store CVV/raw card data.
