ALTER TABLE fines
ADD COLUMN violation_record_id UUID;


ALTER TABLE fines
ADD CONSTRAINT fk_violation_record
FOREIGN KEY (violation_record_id)
REFERENCES violation_records(violation_record_id);
