-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS violation_records;
DROP TYPE IF EXISTS violation_type;
DROP TYPE IF EXISTS violation_status;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE TYPE violation_status AS ENUM ('Ongoing', 'Resolved');
CREATE TYPE violation_type AS ENUM ('Late Return', 'Unpaid Fees', 'Damaged Book'); 
CREATE TABLE violation_records (
    violation_record_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    patron_id UUID REFERENCES patrons (patron_id) ON DELETE CASCADE,
    violation_type violation_type NOT NULL,
    violation_info TEXT NOT NULL,
    violation_created TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    violation_status violation_status DEFAULT 'Ongoing' NOT NULL
);
-- +goose StatementEnd
