

DO $$
BEGIN
  
  IF NOT EXISTS(SELECT table_name FROM information_schema.tables WHERE table_name='users') THEN
    RAISE NOTICE 'Running migration...';

    CREATE TABLE IF NOT EXISTS users (
      id int GENERATED ALWAYS AS IDENTITY,
      username VARCHAR(255),
      email VARCHAR(255),
      password VARCHAR(255),
      balance float8,
      created_at TIMESTAMPTZ,
      updated_at TIMESTAMPTZ
    );

    ALTER TABLE users ADD CONSTRAINT users_pk PRIMARY KEY (id);
    ALTER TABLE users ADD CONSTRAINT users_username_unique UNIQUE(username);
    ALTER TABLE users ADD CONSTRAINT users_email_unique UNIQUE(email);
    ALTER TABLE users ADD CONSTRAINT users_username_email_balance UNIQUE(username,email,balance);

    ALTER TABLE users ALTER COLUMN username SET NOT NULL;
    ALTER TABLE users ALTER COLUMN email SET NOT NULL;
    ALTER TABLE users ALTER COLUMN password SET NOT NULL;

    ALTER TABLE users ALTER COLUMN balance SET DEFAULT 0;
    ALTER TABLE users ALTER COLUMN created_at SET DEFAULT NOW();
    ALTER TABLE users ALTER COLUMN updated_at SET DEFAULT NOW();
    ALTER TABLE users ADD CONSTRAINT users_balance_check CHECK (balance > 0);

  ELSE
    RAISE NOTICE 'Database is already filled with data...';
  END IF;

END; $$;