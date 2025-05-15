-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS memberships;
DROP TYPE IF EXISTS membership_level;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE TYPE membership_level AS ENUM ('Bronze', 'Silver', 'Gold');
CREATE TABLE memberships (
    membership_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    patron_id UUID REFERENCES patrons (patron_id) ON DELETE CASCADE,
    level membership_level DEFAULT 'Bronze' NOT NULL
);
-- +goose StatementEnd
