# User Management System with LDAP Server - Design Document

## Overview

A user management system that provides both HTTP REST API and LDAP protocol access. Users can be organized into hierarchical groups. The system acts as an LDAP Server supporting both OpenLDAP and Microsoft AD conventions.

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Backend | Go + Gin + Ent + Zap |
| Frontend | React + Ant Design (in `web/`) |
| Database | PostgreSQL |
| LDAP | Embedded LDAP server (dual-port with HTTP) |

## Architecture

Single binary, dual-port listener:
- HTTP on port 8080 (Gin REST API + static frontend)
- LDAP on port 389/636 (LDAP protocol)

Both share the same Service and DAO layers (Clean Architecture).

```
cmd/server (Cobra CLI)
├── HTTP Handler (Gin, :8080)
├── LDAP Handler (:389)
└── Shared: Service → DAO (Ent) → PostgreSQL
```

## Directory Structure

```
├── cmd/server/             # Cobra CLI entry point
├── configs/                # Configuration files
├── docs/design/            # Design documents
├── internal/
│   ├── handler/
│   │   ├── http/           # Gin HTTP handlers
│   │   └── ldap/           # LDAP protocol handlers (Bind/Search)
│   ├── service/            # Business logic
│   ├── dao/                # Data access (Ent)
│   ├── domain/             # Domain models
│   ├── schema/             # Ent schema definitions
│   ├── ldap/
│   │   ├── filter/         # RFC 4515 filter parsing & evaluation
│   │   ├── dn/             # DN parsing & building
│   │   └── protocol/       # LDAP protocol codec (BER/ASN.1)
│   ├── middleware/         # HTTP middleware
│   └── config/             # Config structs
├── pkg/                    # Reusable utilities
├── web/                    # React + Ant Design frontend
└── Makefile
```

## Data Models

### User

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| username | string (unique) | Login name |
| display_name | string | Display name |
| email | string (unique) | Email |
| password_hash | string | bcrypt hash |
| phone | string (optional) | Phone number |
| status | enum | enabled / disabled |
| created_at | timestamp | |
| updated_at | timestamp | |

### Group

| Field | Type | Description |
|-------|------|-------------|
| id | UUID | Primary key |
| name | string (unique) | Group name |
| description | string | Description |
| parent_id | UUID (nullable) | Parent group (hierarchy) |
| created_at | timestamp | |
| updated_at | timestamp | |

### UserGroup (many-to-many)

| Field | Type | Description |
|-------|------|-------------|
| user_id | UUID | FK -> User |
| group_id | UUID | FK -> Group |

## LDAP Server

### Supported Operations

- **Bind**: Simple Bind (DN + password authentication)
- **Search**: Users and groups with full filter support
- **Unbind**: Disconnect

### Filter Implementation (RFC 4515)

Full filter grammar:

```
Filter     = AND / OR / NOT / Item
AND        = '(' '&' FilterList ')'
OR         = '(' '|' FilterList ')'
NOT        = '(' '!' Filter ')'
Item       = '(' attr filtertype value ')'
filtertype = '=' / '~=' / '>=' / '<=' / '=*' / substring
```

Supported operators: `=` (equal), `=*` (presence), `>=`, `<=`, `~=` (approx, case-insensitive), `*` (substring), `&` (AND), `|` (OR), `!` (NOT).

Filter-to-SQL conversion: simple filters become Ent predicates; compound filters build predicate trees recursively; substring becomes LIKE/ILIKE; presence becomes IS NOT NULL.

### OpenLDAP vs AD Compatibility

Configurable mode (`openldap` / `activedirectory`):

| Aspect | OpenLDAP | AD |
|--------|----------|-----|
| User objectClass | inetOrgPerson | user, person |
| Group objectClass | groupOfNames | group |
| Username attr | uid | sAMAccountName |
| Member attr | member (DN) | member (DN) |
| User's groups | search group members | memberOf |
| DN format | uid=john,ou=users,dc=... | cn=John,cn=Users,dc=... |

### Attribute Mapping

```
User field     → OpenLDAP attr        → AD attr
username       → uid                  → sAMAccountName
display_name   → cn, displayName      → cn, displayName
email          → mail                 → mail
phone          → telephoneNumber      → telephoneNumber
status         → (custom) status      → userAccountControl
```

## HTTP API

### Users
- `POST /api/v1/users` — Create
- `GET /api/v1/users` — List (paginated, searchable)
- `GET /api/v1/users/:id` — Detail
- `PUT /api/v1/users/:id` — Update
- `DELETE /api/v1/users/:id` — Delete
- `PUT /api/v1/users/:id/password` — Change password
- `PUT /api/v1/users/:id/status` — Enable/disable

### Groups
- `POST /api/v1/groups` — Create
- `GET /api/v1/groups` — List (tree structure)
- `GET /api/v1/groups/:id` — Detail
- `PUT /api/v1/groups/:id` — Update
- `DELETE /api/v1/groups/:id` — Delete

### Membership
- `POST /api/v1/groups/:id/members` — Add members
- `DELETE /api/v1/groups/:id/members/:uid` — Remove member
- `GET /api/v1/groups/:id/members` — Group members
- `GET /api/v1/users/:id/groups` — User's groups

### Auth
- `POST /api/v1/auth/login` — Admin login
- `POST /api/v1/auth/logout` — Logout

### LDAP Config
- `GET /api/v1/ldap/config` — Get config
- `PUT /api/v1/ldap/config` — Update config
- `GET /api/v1/ldap/status` — Service status

## Frontend Pages (React + Ant Design)

| Page | Features |
|------|----------|
| Login | Admin authentication |
| User List | Table, search, pagination, enable/disable |
| User Detail/Edit | User form, group membership |
| Group List | Tree structure display |
| Group Detail | Info, member list, add/remove members |
| LDAP Config | Base DN, mode selection, port settings |
