-- +goose Up
-- +goose StatementBegin
CREATE TABLE public.Orders (
    login TEXT PRIMARY KEY,
    finalizedOrders JSONB,
    notFinalizedOrders JSONB
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE public.Orders;
-- +goose StatementEnd
