-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS pvz (
    id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    name TEXT NOT NULL DEFAULT '',
    address TEXT NOT NULL DEFAULT '',
    contacts TEXT NOT NULL DEFAULT ''
);

COMMENT ON TABLE pvz IS 'Таблица ПВЗ';
COMMENT ON COLUMN pvz.id IS 'Уникальный идентификатор ПВЗ';
COMMENT ON COLUMN pvz.name IS 'Название ПВЗ';
COMMENT ON COLUMN pvz.address IS 'Адрес ПВЗ';
COMMENT ON COLUMN pvz.contacts IS 'Контакты ПВЗ';

CREATE TABLE IF NOT EXISTS "order" (
    id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    client_id UUID NOT NULL DEFAULT gen_random_uuid(),
    stores_till TIMESTAMP NOT NULL DEFAULT now(),
    is_deleted BOOL NOT NULL DEFAULT FALSE,
    give_out_time TIMESTAMP NOT NULL DEFAULT now(),
    is_returned BOOL NOT NULL DEFAULT FALSE
);

COMMENT ON TABLE "order" IS 'Таблица заказов';
COMMENT ON COLUMN "order".id IS 'Уникальный идентификатор заказа';
COMMENT ON COLUMN "order".client_id IS 'Идентификатор клиента';
COMMENT ON COLUMN "order".stores_till IS 'Конечная дата хранения заказа';
COMMENT ON COLUMN "order".is_deleted IS 'Флаг, удален ли заказ';
COMMENT ON COLUMN "order".give_out_time IS 'Дата выдачи заказа';
COMMENT ON COLUMN "order".is_returned IS 'Флаг, возвращен ли заказ';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS pvz;

DROP TABLE IF EXISTS "order";
-- +goose StatementEnd
