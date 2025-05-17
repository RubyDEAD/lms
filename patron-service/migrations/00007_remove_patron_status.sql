-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS patron_status;
DROP TYPE IF EXISTS status;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE TYPE status AS ENUM ('Good', 'Warned', 'Banned', 'Pending');
CREATE TABLE patron_status (
    patron_id UUID PRIMARY KEY REFERENCES patrons (patron_id) ON DELETE CASCADE,
    warning_count INTEGER DEFAULT 0 NOT NULL,
    patron_status status DEFAULT 'Good' NOT NULL,
    unpaid_fees DECIMAL(10,2) DEFAULT 0 CHECK (unpaid_fees >= 0)
);
-- +goose StatementEnd
