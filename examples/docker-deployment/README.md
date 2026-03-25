# Goma Admin - Docker Deployment Example

A production-ready Docker Compose deployment for Goma Admin with PostgreSQL.

## Quick Start

1. Copy the example environment file and update the values:

```bash
cp .env.example .env
```

2. Update `.env` with your production values — at minimum, change `GOMA_JWT_SECRET` and `GOMA_ADMIN_PASSWORD`.

3. Start the services:

```bash
docker compose up -d
```

4. Access Goma Admin at `http://localhost:9000`.

## Environment Variables

| Variable | Description | Default |
|---|---|---|
| `GOMA_DB_HOST` | PostgreSQL host | `localhost` |
| `GOMA_DB_USER` | Database user | `goma` |
| `GOMA_DB_PASSWORD` | Database password | `goma` |
| `GOMA_DB_NAME` | Database name | `goma` |
| `GOMA_DB_PORT` | Database port | `5432` |
| `GOMA_DB_SSL_MODE` | SSL mode (`disable`, `require`) | `disable` |
| `GOMA_DB_URL` | Full database URL (overrides individual DB vars) | — |
| `GOMA_PORT` | HTTP server port | `9000` |
| `GOMA_ENVIRONMENT` | Environment (`development`, `production`) | `development` |
| `GOMA_LOG_LEVEL` | Log level (`debug`, `info`, `warn`, `error`) | `info` |
| `GOMA_JWT_SECRET` | JWT signing secret | `default-secret-key` |
| `GOMA_JWT_ISSUER` | JWT issuer claim | `goma-admin` |
| `GOMA_CORS_ALLOWED_ORIGINS` | CORS origins (comma-separated) | `*` |
| `GOMA_ADMIN_EMAIL` | Default admin email | `admin@example.com` |
| `GOMA_ADMIN_PASSWORD` | Default admin password | `admin` |
| `GOMA_ENABLE_DOCS` | Enable OpenAPI documentation | `true` |
| `GOMA_WEB_DIR` | Frontend assets directory | `web/dist` |

## Stopping

```bash
docker compose down
```

To also remove the database volume:

```bash
docker compose down -v
```
