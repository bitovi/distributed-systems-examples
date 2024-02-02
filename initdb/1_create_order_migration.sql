CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE SEQUENCE order_number_seq START 1;

CREATE TABLE IF NOT EXISTS orders (
  order_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  order_number INT NOT NULL DEFAULT nextval('order_number_seq'),
  customer JSONB NOT NULL,
  status TEXT NOT NULL,
  products JSONB NOT NULL,
  total INT NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);