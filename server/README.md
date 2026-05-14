# Server

This server is organized in a small, explicit set of layers:

- `main.go`: application startup and dependency wiring.
- `endpoints`: route registration only.
- `controllers`: HTTP handlers and request/response concerns.
- `data`: database access and SQL.
- `models`: typed API and persistence shapes.

## Request flow

1. `main` builds the database and Gin router.
2. `endpoints.Register` attaches routes.
3. A controller validates HTTP input and calls its store interface.
4. A typed store in `data` performs the SQL work and returns models.
5. The controller serializes the result as JSON.

## Design notes

- The public API exposes concrete user types only: `admins`, `customers`, `salesmen`, and `guest-employees`.
- The `users` table remains internal and acts as the shared base record for all typed users.
- Controllers depend on interfaces, not concrete SQL details.
- Database operations use short context timeouts for reads and writes.
- Each typed resource has its own small store.
- Shared SQL helpers live in one place to keep CRUD methods small.
- `PATCH` endpoints use pointer fields so partial updates stay explicit.

## Endpoints

- `GET /health`
- `GET /admins`
- `POST /admins`
- `GET /admins/:id`
- `PATCH /admins/:id`
- `DELETE /admins/:id`
- `GET /customers`
- `POST /customers`
- `GET /customers/:id`
- `PATCH /customers/:id`
- `DELETE /customers/:id`
- `GET /salesmen`
- `POST /salesmen`
- `GET /salesmen/:id`
- `PATCH /salesmen/:id`
- `DELETE /salesmen/:id`
- `GET /guest-employees`
- `POST /guest-employees`
- `GET /guest-employees/:id`
- `PATCH /guest-employees/:id`
- `DELETE /guest-employees/:id`
