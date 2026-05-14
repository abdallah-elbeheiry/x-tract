# Server API Guide (**last updated: 14/05/2026 at 17:05 EET**)

This README is written for people who want to use the API from the frontend.

You do **not** need to understand the Go code to use this server.

This file explains:

- what the API is for
- what URLs exist
- what request body to send
- what response body comes back
- what error responses look like
- how the different user types are shaped

## What This API Is

This backend manages **people in the system**.

There are **4 public user types**:

1. `admin`
2. `customer`
3. `salesman`
4. `guest employee`

All of them share the same basic user information:

- `id`
- `username`
- `email`
- `number` (optional)
- `role`

Some types also have extra fields:

- `customer` has `other_statistics`
- `salesman` has `other_statistics`
- `guest employee` has `group_id`
- `admin` has no extra fields beyond `user`

Important:

- There is **no public `/users` endpoint**
- The backend stores users internally in a shared table, but the frontend should only use:
  - `/admins`
  - `/customers`
  - `/salesmen`
  - `/guest-employees`

## Base URL

By default, the server runs on:

```text
http://localhost:8080
```

If the backend is deployed somewhere else, replace that host with the deployed URL.

Example:

```text
http://localhost:8080/customers
```

## Quick Start

### Health check

Use this to confirm the backend is running:

```http
GET /health
```

Response:

```json
{
  "status": "ok"
}
```

## Response Style

Most successful responses come back in this shape:

```json
{
  "data": { ... }
}
```

For list endpoints:

```json
{
  "data": [ ... ]
}
```

For delete endpoints:

- success response has **no body**
- status code is `204 No Content`

## Common User Object

Every person type (admin, salesman, guest employee, customer) contains a nested `user` object.

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

| Field      | Type          | Meaning                                                   |
|------------|---------------|-----------------------------------------------------------|
| `id`       | string (UUID) | Unique ID for this user                                   |
| `username` | string        | Person name                                               |
| `email`    | string        | Email address                                             |
| `number`   | string        | Optional phone number                                     |
| `role`     | string        | One of `ADMIN`, `CUSTOMER`, `SALES_MAN`, `GUEST_EMPLOYEE` |

## Error Responses

Errors come back like this:

```json
{
  "error": "some message"
}
```

### Common status codes

| Status | Meaning                                                    |
|--------|------------------------------------------------------------|
| `200`  | Success                                                    |
| `201`  | Created successfully                                       |
| `204`  | Deleted successfully                                       |
| `400`  | Bad request, invalid input, invalid UUID, validation issue |
| `404`  | Resource not found                                         |
| `500`  | Internal server error                                      |

### Common error messages

Examples you may see:

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

`resource conflict` usually means something like:

- email already exists

`invalid reference` usually means:

- you sent a `group_id` that does not exist

## Request Rules

### Validation rules

When creating or updating users:

- `email` must be a valid email
- `password` must be at least 8 characters
- `number` should be in international phone format
- `group_id` must be a valid UUID if provided

### PATCH behavior

All update endpoints use `PATCH`.

That means:

- send only the fields you want to change
- leave fields out if you do not want to update them

Example:

```json
{
  "email": "new-email@example.com"
}
```

## API Endpoints

---

## Admins

Admins only contain the nested `user` object, this is subject to change as more admin data is needed.

### Create admin

```http
POST /admins
Content-Type: application/json
```

Request body:

```json
{
  "name": "Ali Hassan",
  "email": "ali@example.com",
  "number": "+201001234567",
  "password": "password123"
}
```

Success response:

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

Success response:

```json
{
  "data": [
    {
      "user": {
        "id": "11111111-1111-1111-1111-111111111111",
        "username": "Ali Hassan",
        "email": "ali@example.com",
        "number": "+201001234567",
        "role": "ADMIN"
      }
    }
  ]
}
```

### Get one admin

```http
GET /admins/:id
```

Example:

```http
GET /admins/11111111-1111-1111-1111-111111111111
```

### Update admin

```http
PATCH /admins/:id
Content-Type: application/json
```

Request body example:

```json
{
  "name": "Ali Updated",
  "number": "+201009999999"
}
```

### Delete admin

```http
DELETE /admins/:id
```

Success:

```http
204 No Content
```

---

## Customers

Customers contain:

- `user`
- `other_statistics`

Right now `other_statistics` is just a placeholder field. You can still send it and read it.

### Create customer

```http
POST /customers
Content-Type: application/json
```

Request body:

```json
{
  "name": "Aya Hassan",
  "email": "aya@example.com",
  "number": "+201001234567",
  "password": "password123",
  "stats": ""
}
```

Success response:

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

Success response:

```json
{
  "data": [
    {
      "user": {
        "id": "d0a26a24-0bd9-45c6-b74c-ee3aa8073d17",
        "username": "Aya Hassan",
        "email": "aya@example.com",
        "number": "+201001234567",
        "role": "CUSTOMER"
      },
      "other_statistics": ""
    }
  ]
}
```

### Get one customer

```http
GET /customers/:id
```

### Update customer

```http
PATCH /customers/:id
Content-Type: application/json
```

Request body example:

```json
{
  "name": "Aya Updated",
  "stats": "VIP customer"
}
```

Success response:

```json
{
  "data": {
    "user": {
      "id": "d0a26a24-0bd9-45c6-b74c-ee3aa8073d17",
      "username": "Aya Updated",
      "email": "aya@example.com",
      "number": "+201001234567",
      "role": "CUSTOMER"
    },
    "other_statistics": "VIP customer"
  }
}
```

### Delete customer

```http
DELETE /customers/:id
```

Success:

```http
204 No Content
```

---

## Salesmen

Salesmen have the same shape as customers, but with role `SALES_MAN`.

### Create salesman

```http
POST /salesmen
Content-Type: application/json
```

Request body:

```json
{
  "name": "Omar Adel",
  "email": "omar@example.com",
  "number": "+201002223334",
  "password": "password123",
  "stats": ""
}
```

Success response:

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

### Get one salesman

```http
GET /salesmen/:id
```

### Update salesman

```http
PATCH /salesmen/:id
Content-Type: application/json
```

Request body example:

```json
{
  "stats": "Top seller"
}
```

### Delete salesman

```http
DELETE /salesmen/:id
```

Success:

```http
204 No Content
```

---

## Guest Employees

Guest employees contain:

- `user`
- `group_id` (optional)

If no group is set, `group_id` may be missing from the response.

### Create guest employee

```http
POST /guest-employees
Content-Type: application/json
```

Request body without group:

```json
{
  "name": "Mona Sameh",
  "email": "mona@example.com",
  "number": "+201007778889",
  "password": "password123"
}
```

Request body with group:

```json
{
  "name": "Mona Sameh",
  "email": "mona@example.com",
  "number": "+201007778889",
  "password": "password123",
  "group_id": "33333333-3333-3333-3333-333333333333"
}
```

Success response:

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
    "group_id": "33333333-3333-3333-3333-333333333333"
  }
}
```

### List guest employees

```http
GET /guest-employees
```

### Get one guest employee

```http
GET /guest-employees/:id
```

### Update guest employee

```http
PATCH /guest-employees/:id
Content-Type: application/json
```

Request body example:

```json
{
  "group_id": "33333333-3333-3333-3333-333333333333"
}
```

### Delete guest employee

```http
DELETE /guest-employees/:id
```

Success:

```http
204 No Content
```

---

## Frontend Tips

### 1. Use the nested `user.id` as the main ID

There is no separate `user_id` field in the API response anymore.

Use:

```json
data.user.id
```

### 2. Use the role from the nested user

Use:

```json
data.user.role
```

Example values:

- `ADMIN`
- `CUSTOMER`
- `SALES_MAN`
- `GUEST_EMPLOYEE`

### 3. Keep create and update payloads separate

Create requests need `password`.

Update requests do **not** need `password` unless you are changing it.

### 4. For optional values

- `number` may be missing
- `group_id` may be missing
- `other_statistics` may be an empty string

Your UI should handle those cases safely.

## Example Frontend Fetch

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

### Update customer

```js
const response = await fetch("http://localhost:8080/customers/d0a26a24-0bd9-45c6-b74c-ee3aa8073d17", {
  method: "PATCH",
  headers: {
    "Content-Type": "application/json"
  },
  body: JSON.stringify({
    stats: "VIP customer"
  })
});

const json = await response.json();
console.log(json);
```

## Running the Server

The server needs a PostgreSQL connection string in the `DB_CONNECTION` environment variable.

Optional:

- `SERVER_ADDRESS` to change the port or listening address

### Default behavior

If `SERVER_ADDRESS` is not set, the server runs on:

```text
:8080
```

### Example environment variables

```powershell
$env:DB_CONNECTION="postgres://username:password@localhost:5432/xtract?sslmode=disable"
$env:SERVER_ADDRESS=":8080"
```

Then run the server normally with Go.

## Final Notes

- This API is CRUD-based.
- Every person type is a typed user.
- The backend automatically runs database migrations on startup.
- The frontend should treat this README as the contract for request and response shapes.

If the backend changes later, update this README at the same time so the frontend team never has to guess.
