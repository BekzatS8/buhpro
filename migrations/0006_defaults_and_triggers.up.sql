BEGIN;

-- helper function to update updated_at automatically
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = now();
RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- attach triggers
DO $$
BEGIN
  IF NOT EXISTS (
    SELECT 1 FROM pg_trigger WHERE tgname = 'set_timestamp_users'
  ) THEN
CREATE TRIGGER set_timestamp_users
    BEFORE UPDATE ON users FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_timestamp();
END IF;

  IF NOT EXISTS (
    SELECT 1 FROM pg_trigger WHERE tgname = 'set_timestamp_organizations'
  ) THEN
CREATE TRIGGER set_timestamp_organizations
    BEFORE UPDATE ON organizations FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_timestamp();
END IF;

  IF NOT EXISTS (
    SELECT 1 FROM pg_trigger WHERE tgname = 'set_timestamp_orders'
  ) THEN
CREATE TRIGGER set_timestamp_orders
    BEFORE UPDATE ON orders FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_timestamp();
END IF;

  IF NOT EXISTS (
    SELECT 1 FROM pg_trigger WHERE tgname = 'set_timestamp_bids'
  ) THEN
CREATE TRIGGER set_timestamp_bids
    BEFORE UPDATE ON bids FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_timestamp();
END IF;

  IF NOT EXISTS (
    SELECT 1 FROM pg_trigger WHERE tgname = 'set_timestamp_payments'
  ) THEN
CREATE TRIGGER set_timestamp_payments
    BEFORE UPDATE ON payments FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_timestamp();
END IF;
END;
$$ LANGUAGE plpgsql;

COMMIT;
