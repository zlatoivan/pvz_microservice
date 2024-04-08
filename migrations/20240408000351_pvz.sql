-- +goose Up
-- +goose StatementBegin
ALTER TABLE orders ADD packaging_type TEXT NOT NULL DEFAULT '';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
