-- +goose Up
-- +goose StatementBegin
ALTER TABLE patrons
ADD COLUMN isAdmin BOOLEAN DEFAULT true;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE patrons
DROP COLUMN isAdmin;
-- +goose StatementEnd
