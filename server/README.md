# Server API Guide

This README is for frontend developers.

You do not need to read the Go code to use this API.

This file explains:

- what the API is for
- what endpoints exist
- what request body to send
- what response body comes back
- what errors look like
- how login works
- how JWT works
- how groups and guest employees relate to each other

## What This API Is

This backend manages people in the system.

There are 4 public person types:

1. `admin`
2. `customer`
3. `salesman`
4. `guest employee`

Each person has a nested `user` object with the shared fields:

- `id`
- `username`
- `email`
- `number` (optional)
- `role`

Some types also have extra fields:

- customer: `other_statistics`
- salesman: `other_statistics`
- guest employee: `group`
- admin: no extra fields beyond `user`

There is also a separate resource:

- `group`

## Important API Rule

There is no public `/users` endpoint.

Even though the backend stores all people in a shared internal users table, the frontend should only call:

- `/auth/login`
- `/admins`
- `/customers`
- `/salesmen`
- `/guest-employees`
- `/groups`

## Base URL

By default the server runs on:

```text
http://localhost:8080
```

Example:

```text
http://localhost:8080/customers
```

## Health Check

Use this to check if the backend is alive:

```http
GET /health
```

Response:

```json
{
  "status": "ok"
}
```

## Response Format

Most successful responses look like this:

```json
{
  "data": {}
}
```

List endpoints look like this:

```json
{
  "data": []
}
```

Delete endpoints return:

- HTTP status `204 No Content`
- no JSON body

## Common User Object

Every person type contains a nested `user` object.

Example:

```json
{
  "user": {
    "id": "d0a26a24-0bd9-45c6-b74c-ee3aa8073d17",
    "username": "Aya Hassan",
    "email": "aya@example.com",
    "number": "+201001234567",
    "role": "CUSTOMER"
  }
}
```

### User fields

| Field | Type | Meaning |
|---|---|---|
| `id` | string (UUID) | Unique ID of the person |
| `username` | string | Person name |
| `email` | string | Email address |
| `number` | string | Optional phone number |
| `role` | string | `ADMIN`, `CUSTOMER`, `SALES_MAN`, or `GUEST_EMPLOYEE` |

## Error Format

Errors come back like this:

```json
{
  "error": "some message"
}
```

### Common status codes

| Status | Meaning |
|---|---|
| `200` | Success |
| `201` | Created |
| `204` | Deleted |
| `400` | Bad request / validation error / invalid UUID / invalid reference |
| `401` | Invalid login or invalid token |
| `404` | Resource not found |
| `500` | Internal server error |

### Common error messages

```json
{ "error": "invalid uuid" }
```

```json
{ "error": "resource not found" }
```

```json
{ "error": "resource conflict" }
```

```json
{ "error": "invalid reference" }
```

```json
{ "error": "invalid credentials" }
```

`resource conflict` usually means:

- duplicate email

`invalid reference` usually means:

- you sent a `group_id` that does not exist

## Validation Rules

- `email` must be a valid email
- `password` must be at least 8 characters
- `number` should be in international phone format
- `group_id` must be a valid UUID if provided

## PATCH Rules

All update endpoints use `PATCH`.

That means:

- send only the fields you want to change
- leave other fields out

Example:

```json
{
  "email": "new-email@example.com"
}
```

## Authentication and JWT

The backend now supports login and JWT creation.

What exists now:

- `POST /auth/login`
- access tokens are returned as JWT strings
- JWT verification middleware exists in the backend code

What is not enforced yet:

- the CRUD endpoints are not protected yet
- you do not need to send a token for those endpoints right now

That means JWT is set up and ready, but route protection is still the next step.

### Login

```http
POST /auth/login
Content-Type: application/json
```

Request body:

```json
{
  "email": "aya@example.com",
  "password": "password123"
}
```

Success response:

```json
{
  "data": {
    "access_token": "jwt-token-here",
    "token_type": "Bearer",
    "expires_at": "2026-05-16T10:00:00Z",
    "user": {
      "id": "d0a26a24-0bd9-45c6-b74c-ee3aa8073d17",
      "username": "Aya Hassan",
      "email": "aya@example.com",
      "number": "+201001234567",
      "role": "CUSTOMER"
    }
  }
}
```

### How to send the token later

When route protection is enabled, send:

```http
Authorization: Bearer your-jwt-token
```

## Admin Endpoints

Admins only contain the nested `user`.

### Create admin

```http
POST /admins
Content-Type: application/json
```

Request:

```json
{
  "name": "Ali Hassan",
  "email": "ali@example.com",
  "number": "+201001234567",
  "password": "password123"
}
```

Response:

```json
{
  "data": {
    "user": {
      "id": "11111111-1111-1111-1111-111111111111",
      "username": "Ali Hassan",
      "email": "ali@example.com",
      "number": "+201001234567",
      "role": "ADMIN"
    }
  }
}
```

### List admins

```http
GET /admins
```

### Get admin

```http
GET /admins/:id
```

### Update admin

```http
PATCH /admins/:id
Content-Type: application/json
```

Example request:

```json
{
  "name": "Ali Updated"
}
```

### Delete admin

```http
DELETE /admins/:id
```

## Customer Endpoints

Customers contain:

- `user`
- `other_statistics`

`other_statistics` is currently just a placeholder text field.

### Create customer

```http
POST /customers
Content-Type: application/json
```

Request:

```json
{
  "name": "Aya Hassan",
  "email": "aya@example.com",
  "number": "+201001234567",
  "password": "password123",
  "stats": ""
}
```

Response:

```json
{
  "data": {
    "user": {
      "id": "d0a26a24-0bd9-45c6-b74c-ee3aa8073d17",
      "username": "Aya Hassan",
      "email": "aya@example.com",
      "number": "+201001234567",
      "role": "CUSTOMER"
    },
    "other_statistics": ""
  }
}
```

### List customers

```http
GET /customers
```

### Get customer

```http
GET /customers/:id
```

### Update customer

```http
PATCH /customers/:id
Content-Type: application/json
```

Example request:

```json
{
  "stats": "VIP customer"
}
```

### Delete customer

```http
DELETE /customers/:id
```

## Salesman Endpoints

Salesmen have the same shape as customers, but their role is `SALES_MAN`.

### Create salesman

```http
POST /salesmen
Content-Type: application/json
```

Request:

```json
{
  "name": "Omar Adel",
  "email": "omar@example.com",
  "number": "+201002223334",
  "password": "password123",
  "stats": ""
}
```

Response:

```json
{
  "data": {
    "user": {
      "id": "22222222-2222-2222-2222-222222222222",
      "username": "Omar Adel",
      "email": "omar@example.com",
      "number": "+201002223334",
      "role": "SALES_MAN"
    },
    "other_statistics": ""
  }
}
```

### List salesmen

```http
GET /salesmen
```

### Get salesman

```http
GET /salesmen/:id
```

### Update salesman

```http
PATCH /salesmen/:id
Content-Type: application/json
```

Example request:

```json
{
  "stats": "Top seller"
}
```

### Delete salesman

```http
DELETE /salesmen/:id
```

## Group Endpoints

Groups are mainly used to organize guest employees.

### Create group

```http
POST /groups
Content-Type: application/json
```

Request:

```json
{
  "name": "Factory A"
}
```

Response:

```json
{
  "data": {
    "id": "33333333-3333-3333-3333-333333333333",
    "name": "Factory A"
  }
}
```

### List groups

```http
GET /groups
```

Response:

```json
{
  "data": [
    {
      "id": "33333333-3333-3333-3333-333333333333",
      "name": "Factory A"
    }
  ]
}
```

### Get group

```http
GET /groups/:id
```

### Update group

```http
PATCH /groups/:id
Content-Type: application/json
```

Example request:

```json
{
  "name": "Factory B"
}
```

### Delete group

```http
DELETE /groups/:id
```

## Guest Employee Endpoints

Guest employees contain:

- `user`
- `group` (optional)

Request bodies still use `group_id` when you want to assign a group.

Responses return the nested group object instead of only the ID.

### Create guest employee

```http
POST /guest-employees
Content-Type: application/json
```

Request without group:

```json
{
  "name": "Mona Sameh",
  "email": "mona@example.com",
  "number": "+201007778889",
  "password": "password123"
}
```

Request with group:

```json
{
  "name": "Mona Sameh",
  "email": "mona@example.com",
  "number": "+201007778889",
  "password": "password123",
  "group_id": "33333333-3333-3333-3333-333333333333"
}
```

Response:

```json
{
  "data": {
    "user": {
      "id": "44444444-4444-4444-4444-444444444444",
      "username": "Mona Sameh",
      "email": "mona@example.com",
      "number": "+201007778889",
      "role": "GUEST_EMPLOYEE"
    },
    "group": {
      "id": "33333333-3333-3333-3333-333333333333",
      "name": "Factory A"
    }
  }
}
```

If the guest employee has no group, the `group` field may be missing.

### List guest employees

```http
GET /guest-employees
```

### Get guest employee

```http
GET /guest-employees/:id
```

### Update guest employee

```http
PATCH /guest-employees/:id
Content-Type: application/json
```

Example request:

```json
{
  "group_id": "33333333-3333-3333-3333-333333333333"
}
```

### Delete guest employee

```http
DELETE /guest-employees/:id
```

## Frontend Tips

### Use `user.id` as the main ID

There is no `user_id` in the response.

Use:

```json
data.user.id
```

### Use the nested role

Use:

```json
data.user.role
```

### Group assignment rule

For guest employees:

- request uses `group_id`
- response returns `group`

Request:

```json
{
  "group_id": "33333333-3333-3333-3333-333333333333"
}
```

Response:

```json
{
  "group": {
    "id": "33333333-3333-3333-3333-333333333333",
    "name": "Factory A"
  }
}
```

### Optional fields

These may be missing or empty:

- `number`
- `group`
- `other_statistics`

Your UI should handle that safely.

## Frontend Examples

### Login

```js
const loginResponse = await fetch("http://localhost:8080/auth/login", {
  method: "POST",
  headers: {
    "Content-Type": "application/json"
  },
  body: JSON.stringify({
    email: "aya@example.com",
    password: "password123"
  })
});

const loginJson = await loginResponse.json();
const token = loginJson.data.access_token;
```

### Create customer

```js
const response = await fetch("http://localhost:8080/customers", {
  method: "POST",
  headers: {
    "Content-Type": "application/json"
  },
  body: JSON.stringify({
    name: "Aya Hassan",
    email: "aya@example.com",
    number: "+201001234567",
    password: "password123",
    stats: ""
  })
});

const json = await response.json();
console.log(json);
```

### Example future authenticated request

```js
const response = await fetch("http://localhost:8080/customers", {
  method: "GET",
  headers: {
    Authorization: `Bearer ${token}`
  }
});
```

## Environment Variables

Required:

- `DB_CONNECTION`

Optional:

- `SERVER_ADDRESS`
- `JWT_SECRET`

### Example

```powershell
$env:DB_CONNECTION="postgres://username:password@localhost:5432/xtract?sslmode=disable"
$env:SERVER_ADDRESS=":8080"
$env:JWT_SECRET="change-this-in-real-environments"
```

If `SERVER_ADDRESS` is not set, the server listens on:

```text
:8080
```

If `JWT_SECRET` is not set, the backend falls back to a development secret. That is okay for local development, but not okay for production.

## Final Notes

- This API is CRUD-based.
- Every person type is a typed user.
- Groups are their own resource.
- Login now returns JWT access tokens.
- JWT protection middleware exists, but route protection is not switched on yet.
- The backend automatically runs migrations on startup.

If the backend changes, update this README at the same time so the frontend never has to guess.
