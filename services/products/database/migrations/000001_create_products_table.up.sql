DO $$
BEGIN
  IF NOT EXISTS (SELECT table_name FROM information_schema.tables WHERE table_name='products') THEN
      CREATE TYPE currency AS ENUM ('GBP', 'EURO', 'USD');

      CREATE TABLE products(
        id int GENERATED ALWAYS AS IDENTITY,
        name VARCHAR(255),
        description VARCHAR(500),
        price float8,
        quantity INT,
        currency currency,
        image text,
        created_at TIMESTAMPTZ,
        updated_at TIMESTAMPTZ
      );

      ALTER TABLE products ADD CONSTRAINT products_pk PRIMARY KEY (id);
      ALTER TABLE products ADD CONSTRAINT products_name_unique UNIQUE (name);
      ALTER TABLE products ADD CONSTRAINT products_name_desc_image_price_unique UNIQUE (name, description, image, price); 
      ALTER TABLE products ADD CONSTRAINT products_price_constraint CHECK (price > 0);
      ALTER TABLE products ADD CONSTRAINT products_quantity_constraint CHECK (quantity >= 0);

      ALTER TABLE products ALTER COLUMN price SET NOT NULL;
      ALTER TABLE products ALTER COLUMN description SET NOT NULL;
      ALTER TABLE products ALTER COLUMN name SET NOT NULL;
      ALTER TABLE products ALTER COLUMN created_at SET NOT NULL;
      ALTER TABLE products ALTER COLUMN updated_at SET NOT NULL;
      ALTER TABLE products ALTER COLUMN image SET NOT NULL;

      ALTER TABLE products ALTER COLUMN quantity SET DEFAULT 0;
      ALTER TABLE products ALTER COLUMN created_at SET DEFAULT NOW();
      ALTER TABLE products ALTER COLUMN currency SET DEFAULT 'GBP';
  ELSE
    RAISE NOTICE 'Database is already filled with products table...';
  END IF;

END; $$;