-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS pvz (
    id uuid PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    name TEXT NOT NULL DEFAULT '',
    address TEXT NOT NULL DEFAULT '',
    contacts TEXT NOT NULL DEFAULT ''
);

COMMENT ON TABLE pvz IS 'Таблица ПВЗ';
COMMENT ON COLUMN pvz.id IS 'Уникальный идентификатор ПВЗ';
COMMENT ON COLUMN pvz.name IS 'Название ПВЗ';
COMMENT ON COLUMN pvz.address IS 'Адрес ПВЗ';
COMMENT ON COLUMN pvz.contacts IS 'Контакты ПВЗ';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS pvz;
-- +goose StatementEnd
