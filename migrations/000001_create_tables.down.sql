DROP INDEX IF EXISTS idx_orders_status;
DROP INDEX IF EXISTS idx_orders_uploaded_at_user_id_order_number;
DROP INDEX IF EXISTS idx_balance_processed_at_user_id;
DROP INDEX IF EXISTS idx_user_balance_user_id;

DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS user_balance_flow;
DROP TABLE IF EXISTS user_balance;

DROP TYPE IF EXISTS order_status;