CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    amount NUMERIC(12, 2) NOT NULL,
    category TEXT NOT NULL,
    merchant TEXT,
    transaction_date TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE baselines (
    category TEXT PRIMARY KEY,
    average_amount NUMERIC(12, 2) NOT NULL,
    stddev_amount NUMERIC(12, 2) NOT NULL,
    updated_at TIMESTAMP DEFAULT now()
);

CREATE TABLE alerts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    transaction_id UUID REFERENCES transactions(id) ON DELETE CASCADE,
    alert_message TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT now()
);