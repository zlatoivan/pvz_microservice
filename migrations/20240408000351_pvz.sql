-- +goose Up
-- +goose StatementBegin
ALTER TABLE orders ADD packaging_type TEXT NOT NULL DEFAULT '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS pvzs;

DROP TABLE IF EXISTS orders;
-- +goose StatementEnd
