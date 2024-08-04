-- +goose Up
-- +goose StatementBegin
CREATE TABLE public.Withdrawals (
    "order" TEXT PRIMARY KEY,
    login TEXT,
    sum FLOAT8 DEFAULT 0,
    processed_at TIMESTAMP WITH TIME ZONE NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE public.Withdrawals;
-- +goose StatementEnd
