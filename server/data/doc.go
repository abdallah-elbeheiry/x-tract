// Package data owns database access for the server project.
//
// The package keeps a small split:
//   - config.go opens the database and runs migrations.
//   - store.go owns the shared sql.DB wrapper.
//   - admins.go, customers.go, salesmen.go, and guest_employees.go expose
//     typed CRUD stores.
//   - sql_helpers.go contains shared SQL helpers used by the typed stores.
package data
