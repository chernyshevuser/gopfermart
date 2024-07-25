-- +goose Up
-- +goose StatementBegin
CREATE TABLE public.Users (
    login TEXT PRIMARY KEY,
    password TEXT,
    balance FLOAT8 DEFAULT 0
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE public.Users;
-- +goose StatementEnd
