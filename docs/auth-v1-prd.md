# Auth v1 PRD

## Problem Statement

SchoolCRM needs real authentication across the web admin, the primary SchoolCRM API, and the RAG service. The current backend has a strong foundation for JWT generation, token validation, password verification, role checks, and auth middleware, but the product experience is incomplete: the web admin still uses mock authentication, the token endpoint is not shaped like a normal login flow, and the RAG service accepts Bearer tokens without validating them.

This creates three user-facing problems:

1. Users cannot log in to the web admin using real SchoolCRM accounts.
2. The RAG service cannot reliably know which user made a request.
3. School-specific authorization cannot be expressed with the current generic `ADMIN` and `USER` roles.

## Solution

Auth v1 will make the standalone auth service the single identity provider for SchoolCRM. The web admin will log in against the auth service, receive a JWT access token and user profile, and send the same token to both the SchoolCRM API and the RAG service. The SchoolCRM API and RAG service will both validate tokens through the auth service.

Auth v1 will also migrate the role model from generic roles to school roles:

- `SUPER_ADMIN`
- `SCHOOL_ADMIN`
- `TEACHER`
- `STUDENT`
- `PARENT`

Existing `ADMIN` users will migrate to `SCHOOL_ADMIN`. Existing `USER` users will migrate to `STUDENT`. The migration should be hard, not alias-based; new code and seed/test data should use only school roles.

The first version will use 8-hour access tokens, no refresh tokens, localStorage-backed frontend token storage, and role-based RAG authorization without school-scoped document filtering.

## User Stories

1. As a SchoolCRM user, I want to log in with my email and password, so that I can access the web admin with my real account.
2. As a SchoolCRM user, I want my login session to work across SchoolCRM API features and RAG features, so that I do not need separate credentials.
3. As a SchoolCRM user, I want my session to expire after a school-day-length window, so that stale sessions do not remain valid indefinitely.
4. As a SchoolCRM user, I want to be denied access when my token is missing or invalid, so that protected resources are not exposed.
5. As a SchoolCRM user, I want the frontend to know my profile and roles after login, so that the UI can route me to the right areas.
6. As a SchoolCRM user, I want logout to remove my local session, so that another person using the browser cannot continue as me.
7. As a school administrator, I want school-specific roles, so that authorization reflects how schools actually operate.
8. As a school administrator, I want old admin accounts migrated to `SCHOOL_ADMIN`, so that existing privileged users keep appropriate access.
9. As a school administrator, I want old generic user accounts migrated to `STUDENT`, so that existing low-privilege users remain constrained.
10. As a `SUPER_ADMIN`, I want full platform access, so that I can operate and support the SchoolCRM installation.
11. As a `SCHOOL_ADMIN`, I want to manage school-level resources, so that I can administer the school’s data.
12. As a `TEACHER`, I want access to teaching and learning workflows, so that I can support students.
13. As a `STUDENT`, I want access to allowed learning and school resources, so that I can use the system safely.
14. As a `PARENT`, I want access appropriate to a parent role, so that I can view or query permitted school information.
15. As a RAG user, I want RAG requests to use my authenticated identity, so that document and query audit trails show the real requester.
16. As a RAG user, I want invalid tokens rejected, so that only authenticated users can use RAG in normal environments.
17. As a RAG user in local development, I want anonymous mode to remain available by explicit configuration, so that local testing remains easy.
18. As a `SUPER_ADMIN`, `SCHOOL_ADMIN`, or `TEACHER`, I want to ingest RAG documents, so that authorized staff can add knowledge sources.
19. As a `SUPER_ADMIN` or `SCHOOL_ADMIN`, I want to delete RAG documents, so that administrative users can manage the knowledge base.
20. As any authenticated school role, I want to query RAG, so that I can retrieve relevant school information.
21. As an engineer, I want the auth service to own login and token validation, so that identity rules are centralized.
22. As an engineer, I want RAG to validate tokens through the auth service, so that Python does not duplicate Go JWT and OPA policy behavior.
23. As an engineer, I want the frontend auth service to hide token storage details, so that token storage can be hardened later without rewriting the app.
24. As an engineer, I want auth failures to return clear `401` or `403` responses, so that clients can respond correctly.
25. As an engineer, I want the first slice to avoid refresh tokens, revocation, and session storage, so that real login can ship without overbuilding auth infrastructure.

## Implementation Decisions

- The standalone auth service is the identity provider for Auth v1.
- The auth service will expose a normal email/password login endpoint.
- Login uses email and password only in v1.
- Login returns an access token, token expiry, and a sanitized user profile.
- The login response includes user ID, name, email, and roles.
- Access tokens expire after 8 hours.
- Refresh tokens are out of scope for v1.
- Logout is frontend-only in v1 and clears local token/user state.
- The frontend stores the v1 access token in localStorage behind the Angular auth service abstraction.
- The frontend replaces mock users with real auth service API calls.
- The frontend sends the same Bearer token to SchoolCRM API and RAG API requests.
- SchoolCRM API continues to validate tokens through the auth service.
- RAG validates Bearer tokens by calling the auth service authentication endpoint.
- RAG should reject missing or invalid tokens with `401 Unauthorized` when anonymous mode is disabled.
- RAG anonymous mode remains available only through explicit local/dev configuration.
- RAG anonymous mode should be disabled by default.
- RAG should use the authenticated user ID from the auth service response for audit fields.
- RAG authorization is role-based in v1.
- RAG document access is not school-scoped in v1.
- RAG ingestion is allowed for `SUPER_ADMIN`, `SCHOOL_ADMIN`, and `TEACHER`.
- RAG deletion is allowed for `SUPER_ADMIN` and `SCHOOL_ADMIN`.
- RAG querying is allowed for `SUPER_ADMIN`, `SCHOOL_ADMIN`, `TEACHER`, `STUDENT`, and `PARENT`.
- The role model migrates from `ADMIN` and `USER` to school roles.
- `ADMIN` migrates to `SCHOOL_ADMIN`.
- `USER` migrates to `STUDENT`.
- Compatibility aliases for `ADMIN` and `USER` will not be kept.
- School roles are global for v1.
- Per-school role scoping is deferred.
- Seed data, tests, authorization policies, and frontend route guards should use school roles after migration.
- Auth service token validation remains the central point for key handling, user-enabled checks, issuer validation, and role claims.

## Testing Decisions

- Tests should cover external behavior and contracts, not internal implementation details.
- Auth service tests should cover successful login, invalid password, unknown email, disabled user, returned token expiry, and returned sanitized user profile.
- Role tests should cover the new school role parser/validator and reject old `ADMIN`/`USER` values after migration.
- Authorization policy tests should cover each school role against the protected rules used by SchoolCRM and RAG.
- RAG auth tests should mock the auth service response and verify missing token, invalid token, valid token, and anonymous-mode behavior.
- RAG controller tests should verify that `uploaded_by` and `requested_by` use the authenticated subject returned by the auth service.
- Frontend auth tests should verify that login calls the backend, stores the token/profile, exposes current user state, sends Bearer tokens through the interceptor, and clears state on logout.
- Integration-style backend compile checks should include the auth service, SchoolCRM service, app SDK, and RAG service after auth changes.
- Existing auth tests in the Go auth SDK should remain as prior art for token generation and validation behavior.
- Existing RAG controller tests should be extended rather than bypassing auth validation.

## Out of Scope

- Refresh tokens.
- Server-side session storage.
- Token revocation or blacklist.
- Password reset or forgot-password flows.
- Account lockout or rate limiting.
- Multi-factor authentication.
- Username or student-ID login.
- HttpOnly cookie auth and CSRF strategy.
- Per-school role assignments.
- School-scoped RAG document filtering.
- Role management UI.
- User registration/self-signup.
- Local JWT validation in RAG using public keys or JWKS.

## Further Notes

- Auth v1 intentionally prioritizes one real login flow shared across services.
- Calling the auth service from RAG on each request is acceptable for v1 because it centralizes correctness. A short validation cache can be considered later if RAG traffic requires it.
- The role migration is a breaking change for any fixtures, tests, or clients that still use `ADMIN` or `USER`.
- The frontend localStorage decision is a pragmatic v1 choice. The auth service and frontend should be structured so token storage can move to a hardened cookie/session model later.

## Default Local Users

The local seed data creates four default users for development and test workflows:

┌──────────────────────────────────────┬────────────────────┬────────────────────────┬──────────────┬──────────┐
│ User ID                              │ Name               │ Email                  │ Role         │ Password │
├──────────────────────────────────────┼────────────────────┼────────────────────────┼──────────────┼──────────┤
│ b7f1d86f-0f1f-4c7b-84d7-99a570f14b6f │ Super Admin Gopher │ superadmin@example.com │ SUPER_ADMIN  │ gophers  │
│ 5cf37266-3473-4006-984f-9325122678b7 │ Admin Gopher       │ admin@example.com      │ SCHOOL_ADMIN │ gophers  │
│ c41fa5d3-d61f-45f1-b054-d2c7a3704019 │ Teacher Gopher     │ teacher@example.com    │ TEACHER      │ gophers  │
│ 45b5fbd3-755f-4379-8f07-a58d4a30fa2f │ User Gopher        │ user@example.com       │ STUDENT      │ gophers  │
└──────────────────────────────────────┴────────────────────┴────────────────────────┴──────────────┴──────────┘

Use the auth login endpoint with email/password credentials:

```bash
curl -i -X POST http://localhost:6000/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"email":"superadmin@example.com","password":"gophers"}'
```
