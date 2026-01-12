CREATE TABLE IF NOT EXISTS users (
    id            INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    login         VARCHAR(255) UNIQUE NOT NULL,
    password      VARCHAR(255)NOT NULL,
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

DROP TYPE IF EXISTS order_status;
CREATE TYPE order_status AS ENUM ('NEW', 'PROCESSING', 'INVALID', 'PROCESSED');

CREATE TABLE IF NOT EXISTS orders (
    order_number      VARCHAR(255)  UNIQUE NOT NULL  PRIMARY KEY,
    user_id           INTEGER        NOT NULL,
    status            order_status,
    accrual           DECIMAL(10, 2),
    uploaded_at       TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    processed_at      TIMESTAMP WITH TIME ZONE
);
CREATE INDEX idx_orders_status ON orders (status);
CREATE INDEX idx_orders_uploaded_at_user_id_order_number ON orders (order_number, user_id, uploaded_at);
CREATE INDEX idx_orders_uploaded_at_user_id ON orders (user_id, uploaded_at);


CREATE TABLE IF NOT EXISTS user_balance_flow (
    user_id           INTEGER        NOT NULL,
    order_number      VARCHAR(255) PRIMARY KEY,
    amount            DECIMAL(10, 2) NOT NULL,
    processed_at      TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_balance_processed_at_user_id ON user_balance_flow (user_id, processed_at);

CREATE TABLE IF NOT EXISTS user_balance
(
    user_id         INTEGER PRIMARY KEY,
    current         DECIMAL(10, 2) DEFAULT 0,
    withdrawn       DECIMAL(10, 2) DEFAULT 0,
    updated_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX idx_user_balance_user_id ON user_balance (user_id);









