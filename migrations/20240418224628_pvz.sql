-- +goose Up
-- +goose StatementBegin
ALTER TABLE orders ADD packaging_type TEXT NOT NULL DEFAULT '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE orders DROP COLUMN packaging_type;
-- +goose StatementEnd
