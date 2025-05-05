CREATE TABLE IF NOT EXISTS fines (
  fine_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  patron_id UUID NOT NULL,
  book_id UUID NOT NULL,
  days_late INT NOT NULL,
  rate_per_day FLOAT NOT NULL,
  amount FLOAT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);



