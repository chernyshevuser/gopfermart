-- +goose Up
-- +goose StatementBegin
CREATE TABLE public.Orders (
    number TEXT PRIMARY KEY,
    login TEXT,
    status TEXT,
    accrual FLOAT8 DEFAULT 0,
    uploaded_at TIMESTAMP WITH TIME ZONE NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE public.Orders;
-- +goose StatementEnd
