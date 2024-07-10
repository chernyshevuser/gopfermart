-- +goose Up
-- +goose StatementBegin
CREATE TABLE public.Numbers (
    number NUMERIC PRIMARY KEY,
    login TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE public.Numbers;
-- +goose StatementEnd
