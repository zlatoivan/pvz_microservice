-- +goose Up
-- +goose StatementBegin
CREATE TABLE pvzs (
                      id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
                      name TEXT NOT NULL DEFAULT '',
                      address TEXT NOT NULL DEFAULT '',
                      contacts TEXT NOT NULL DEFAULT ''
);

COMMENT ON TABLE pvzs IS 'Таблица ПВЗ';
COMMENT ON COLUMN pvzs.id IS 'Уникальный идентификатор ПВЗ';
COMMENT ON COLUMN pvzs.name IS 'Название ПВЗ';
COMMENT ON COLUMN pvzs.address IS 'Адрес ПВЗ';
COMMENT ON COLUMN pvzs.contacts IS 'Контакты ПВЗ';

CREATE TABLE orders (
                        id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
                        client_id UUID NOT NULL DEFAULT gen_random_uuid(),
                        weight INT NOT NULL DEFAULT 0,
                        cost INT NOT NULL DEFAULT 0,
                        stores_till TIMESTAMP NOT NULL DEFAULT now(),
                        give_out_time TIMESTAMP NOT NULL DEFAULT now(),
                        is_returned BOOL NOT NULL DEFAULT FALSE,
                        is_deleted BOOL NOT NULL DEFAULT FALSE
);

COMMENT ON TABLE orders IS 'Таблица заказов';
COMMENT ON COLUMN orders.id IS 'Уникальный идентификатор заказа';
COMMENT ON COLUMN orders.client_id IS 'Идентификатор клиента';
COMMENT ON COLUMN orders.weight IS 'Вес заказа';
COMMENT ON COLUMN orders.cost IS 'Стоимость заказа';
COMMENT ON COLUMN orders.stores_till IS 'Конечная дата хранения заказа';
COMMENT ON COLUMN orders.give_out_time IS 'Дата выдачи заказа';
COMMENT ON COLUMN orders.is_returned IS 'Флаг, возвращен ли заказ';
COMMENT ON COLUMN orders.is_deleted IS 'Флаг, удален ли заказ';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE pvzs;

DROP TABLE orders;
-- +goose StatementEnd
