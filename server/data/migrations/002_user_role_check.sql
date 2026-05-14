-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
    ADD CONSTRAINT users_role_check
    CHECK (role IN ('ADMIN', 'CUSTOMER', 'SALES_MAN', 'GUEST_EMPLOYEE'));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
    DROP CONSTRAINT IF EXISTS users_role_check;
-- +goose StatementEnd
