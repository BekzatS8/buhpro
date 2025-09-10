BEGIN;

DROP TRIGGER IF EXISTS set_timestamp_payments ON payments;
DROP TRIGGER IF EXISTS set_timestamp_bids ON bids;
DROP TRIGGER IF EXISTS set_timestamp_orders ON orders;
DROP TRIGGER IF EXISTS set_timestamp_organizations ON organizations;
DROP TRIGGER IF EXISTS set_timestamp_users ON users;

DROP FUNCTION IF EXISTS trigger_set_timestamp();

DROP TABLE IF EXISTS wallet_transactions;
DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS bids;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS organizations;
DROP TABLE IF EXISTS users;

COMMIT;
