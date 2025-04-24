-- +goose Up
-- +goose StatementBegin
ALTER TABLE patrons
ADD COLUMN password VARCHAR(255)  NOT NULL,
ADD COLUMN email VARCHAR(60) UNIQUE NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE patrons
DROP COLUMN password,
DROP COLUMN email;
-- +goose StatementEnd
