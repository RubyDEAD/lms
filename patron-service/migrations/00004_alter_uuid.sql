-- +goose Up
-- +goose StatementBegin
ALTER TABLE patrons
ALTER COLUMN patron_id DROP DEFAULT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE patrons
ALTER COLUMN patron_id SET DEFAULT gen_random_uuid();
-- +goose StatementEnd
