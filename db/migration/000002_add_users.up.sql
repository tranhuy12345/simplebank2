CREATE TABLE "users" (
  "username" varchar NOT null PRIMARY KEY,
  "hash_password" varchar NOT null,
  "full_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "password_changed_at" timestamp not null default '0001-01-01 00:00:00Z',
  "created_at" timestamp default (now())
);

alter table "accounts" add foreign key ("owner") REFERENCES "users"("username");
--create unique index on "accounts"("owner","currency");
alter table "accounts" add constraint "owner_currency_key" unique("owner","currency");