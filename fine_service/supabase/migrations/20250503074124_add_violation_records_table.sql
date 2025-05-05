CREATE TABLE violation_records (
  violation_record_id UUID PRIMARY KEY,
  patron_id UUID NOT NULL,
  violation_type TEXT NOT NULL CHECK (violation_type IN ('Late_Return', 'Unpaid_Fees', 'Damaged_Book')),
  violation_info TEXT NOT NULL,
  violation_created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  violation_status TEXT NOT NULL CHECK (violation_status IN ('Ongoing', 'Resolved'))
);
