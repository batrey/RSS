CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE SEQUENCE articles_id_seq;

-- updated_at column trigger
CREATE OR REPLACE FUNCTION set_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$$ language 'plpgsql';





CREATE TABLE IF NOT EXISTS "articles"(
    "id" integer NOT NULL DEFAULT nextval('articles_id_seq'),
    "uuid" UUID NOT NULL  DEFAULT uuid_generate_v4(),
    "category" varchar(255)  NOT NULL,
    "article"  jsonb NOT NULL DEFAULT '{}'::jsonb,
    "updated_at" TIMESTAMP,
    "created_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);



create trigger articles_updated_at 
  before update on articles
  for each row
  execute procedure set_updated_at_column();
