
DO $$
BEGIN 

  IF NOT EXISTS(SELECT table_name FROM information_schema.tables WHERE table_name='sessions') THEN
    RAISE NOTICE 'Creating sessions table...';

    CREATE TABLE IF NOT EXISTS sessions (
      id int GENERATED ALWAYS AS IDENTITY,
      user_id int,
      guid varchar(255),
      created_at TIMESTAMPTZ,
      expires_at TIMESTAMPTZ
    );

    ALTER TABLE sessions ADD CONSTRAINT sessions_pk PRIMARY KEY (id);
    ALTER TABLE sessions ADD CONSTRAINT sessions_users_fk FOREIGN KEY (user_id) REFERENCES users;
    ALTER TABLE sessions ADD CONSTRAINT sessions_guid_unique UNIQUE (guid);
    ALTER TABLE sessions ADD CONSTRAINT sessions_user_id_guid_created_at_expires_at UNIQUE(user_id,guid,created_at, expires_at);
    ALTER TABLE sessions ALTER COLUMN user_id SET NOT NULL;
    ALTER TABLE sessions ALTER COLUMN guid SET NOT NULL;
    ALTER TABLE sessions ALTER COLUMN expires_at SET NOT NULL;
    ALTER TABLE sessions ALTER COLUMN created_at SET NOT NULL;
    ALTER TABLE sessions ALTER COLUMN created_at SET DEFAULT NOW();

  ELSE
    RAISE NOTICE 'Database is already filled with data...';
  END IF;

END $$;