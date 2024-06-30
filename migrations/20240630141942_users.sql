-- +goose Up
-- +goose StatementBegin
CREATE TABLE public.Users (
    login TEXT PRIMARY KEY,
    password TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE public.Users;
-- +goose StatementEnd
