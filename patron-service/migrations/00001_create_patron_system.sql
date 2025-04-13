-- +goose Up
-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Create ENUM types
CREATE TYPE membership_level AS ENUM ('Bronze', 'Silver', 'Gold');
CREATE TYPE status AS ENUM ('Good', 'Warned', 'Banned', 'Pending');
CREATE TYPE reservation_status AS ENUM ('Pending', 'Fulfilled', 'Canceled');
CREATE TYPE violation_type AS ENUM ('Late Return', 'Unpaid Fees', 'Damaged Book');

-- Create Patrons Table
CREATE TABLE patrons (
    patron_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    phone_number VARCHAR(15) CHECK (phone_number ~ '^[0-9]{10,15}$') NOT NULL
);

-- Create Memberships Table
CREATE TABLE memberships (
    membership_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    patron_id UUID REFERENCES patrons (patron_id) ON DELETE CASCADE,
    level membership_level DEFAULT 'Bronze' NOT NULL
);

-- Create Patron Status Table
CREATE TABLE patron_status (
    patron_id UUID PRIMARY KEY REFERENCES patrons (patron_id) ON DELETE CASCADE,
    warning_count INTEGER DEFAULT 0 NOT NULL,
    patron_status status DEFAULT 'Good' NOT NULL,
    unpaid_fees DECIMAL(10,2) DEFAULT 0 CHECK (unpaid_fees >= 0)
);

-- Create Violation Records Table
CREATE TABLE violation_records (
    violation_record_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    patron_id UUID REFERENCES patrons (patron_id) ON DELETE CASCADE,
    violation_type violation_type NOT NULL,
    violation_info TEXT NOT NULL
);

-- Create Reservations Table
CREATE TABLE reservations (
    reservation_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    patron_id UUID REFERENCES patrons (patron_id) ON DELETE CASCADE,
    book_id UUID NOT NULL, -- Reference to book service (external DB)
    reservation_date TIMESTAMP DEFAULT NOW(),
    status reservation_status DEFAULT 'Pending'
);

-- Create Borrowed Books Table
CREATE TABLE borrowed_books (
    borrow_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    patron_id UUID REFERENCES patrons (patron_id) ON DELETE CASCADE,
    book_id UUID NOT NULL, -- Reference to book service (external DB)
    borrow_date TIMESTAMP DEFAULT NOW(),
    due_date TIMESTAMP NOT NULL,
    return_date TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS borrowed_books;
DROP TABLE IF EXISTS reservations;
DROP TABLE IF EXISTS violation_records;
DROP TABLE IF EXISTS patron_status;
DROP TABLE IF EXISTS memberships;
DROP TABLE IF EXISTS patrons;
DROP TYPE IF EXISTS violation_type;
DROP TYPE IF EXISTS reservation_status;
DROP TYPE IF EXISTS status;
DROP TYPE IF EXISTS membership_level;
DROP EXTENSION IF EXISTS "pgcrypto";
