DO $$
BEGIN
  IF NOT EXISTS (SELECT table_name FROM information_schema.tables WHERE table_name='orders') THEN
      CREATE TYPE order_status AS ENUM ('pending', 'accepted', 'declined');

      CREATE TABLE orders(
        id int GENERATED ALWAYS AS IDENTITY,
        uid varchar(255),
        product_id int,
        quantity int,
        total float8,
        user_id int,
        status order_status,
        created_at TIMESTAMPTZ,
        updated_at TIMESTAMPTZ
      );

      ALTER TABLE orders ADD CONSTRAINT orders_pk PRIMARY KEY (id);
      ALTER TABLE orders ADD CONSTRAINT orders_product_id_fk FOREIGN KEY (product_id) REFERENCES products;
      ALTER TABLE orders ADD CONSTRAINT orders_uid_product_id_quantity_total_user_id_status_unique UNIQUE (uid,product_id,quantity,total,user_id,status); 
      ALTER TABLE orders ADD CONSTRAINT orders_total_constraint CHECK (total > 0.0);
      ALTER TABLE orders ADD CONSTRAINT orders_quantity_constraint CHECK (quantity >= 0);

      ALTER TABLE orders ALTER COLUMN uid SET NOT NULL;
      ALTER TABLE orders ALTER COLUMN product_id SET NOT NULL;
      ALTER TABLE orders ALTER COLUMN quantity SET NOT NULL;
      ALTER TABLE orders ALTER COLUMN total SET NOT NULL;
      ALTER TABLE orders ALTER COLUMN user_id SET NOT NULL;
      ALTER TABLE orders ALTER COLUMN created_at SET NOT NULL;
      ALTER TABLE orders ALTER COLUMN updated_at SET NOT NULL;

      ALTER TABLE orders ALTER COLUMN created_at SET DEFAULT NOW();
      ALTER TABLE orders ALTER COLUMN status SET DEFAULT 'pending';

  ELSE
    RAISE NOTICE 'Database is already filled with products table...';
  END IF;

END; $$;