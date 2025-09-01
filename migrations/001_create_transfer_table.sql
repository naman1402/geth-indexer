-- Create transfer table for storing ERC-20 transfer events
CREATE TABLE IF NOT EXISTS transfer (
    id SERIAL PRIMARY KEY,
    "name" VARCHAR(50) NOT NULL,
    "blockNumber" BIGINT NOT NULL,
    "txnHash" VARCHAR(66) NOT NULL,
    "contract" VARCHAR(42) NOT NULL,
    "from" VARCHAR(42) NOT NULL,
    "to" VARCHAR(42) NOT NULL,
    "value" NUMERIC NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE("txnHash", "contract", "from", "to", "value") -- Prevent duplicates
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_transfer_contract ON transfer("contract");
CREATE INDEX IF NOT EXISTS idx_transfer_from ON transfer("from");
CREATE INDEX IF NOT EXISTS idx_transfer_to ON transfer("to");
CREATE INDEX IF NOT EXISTS idx_transfer_block ON transfer("blockNumber");
CREATE INDEX IF NOT EXISTS idx_transfer_txn ON transfer("txnHash");
