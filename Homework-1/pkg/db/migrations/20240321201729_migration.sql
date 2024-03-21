-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS pvz (
     id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY NOT NULL,
     name VARCHAR(100) NOT NULL DEFAULT '',
     address VARCHAR(100) NOT NULL DEFAULT '',
     contacts VARCHAR(50) NOT NULL DEFAULT ''
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS pvz;
-- +goose StatementEnd
