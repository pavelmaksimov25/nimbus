# Security

## Known Vulnerabilities

The following vulnerabilities have been identified by `govulncheck` and are tracked for resolution. All are currently muted in CI.

| ID | Severity | Component | Summary | Fix |
|----|----------|-----------|---------|-----|
| GO-2026-4341 | Medium | `net/url` (stdlib) | Memory exhaustion in query param parsing | Upgrade Go to 1.25.7+ |
| GO-2026-4340 | Medium | `crypto/tls` (stdlib) | Handshake at incorrect encryption level | Upgrade Go to 1.25.7+ |
| GO-2026-4337 | Medium | `crypto/tls` (stdlib) | Unexpected session resumption | Upgrade Go to 1.25.7+ |
| GO-2025-4175 | Medium | `crypto/x509` (stdlib) | Improper DNS wildcard constraint handling | Upgrade Go to 1.25.7+ |
| GO-2025-4155 | Medium | `crypto/x509` (stdlib) | Excessive resource consumption in cert validation | Upgrade Go to 1.25.7+ |
| GO-2025-4233 | Medium | `quic-go` v0.55.0 (indirect) | HTTP/3 QPACK Header Expansion DoS | Upgrade to v0.57.0+ |
