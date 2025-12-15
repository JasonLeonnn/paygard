-- transactions lookups
CREATE INDEX IF NOT EXISTS idx_transactions_category_date
ON transactions (category, transaction_date);

-- baselines lookup
CREATE UNIQUE INDEX IF NOT EXISTS idx_baselines_category
ON baselines (category);

-- alerts sorting
CREATE INDEX IF NOT EXISTS idx_alerts_created_at
ON alerts (created_at DESC);

CREATE INDEX IF NOT EXISTS idx_alerts_transaction_id
ON alerts(transaction_id);