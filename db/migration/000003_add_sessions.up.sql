CREATE TABLE "sessions" (
  "id" uuid PRIMARY KEY,
  "username" varchar NOT NULL,
  "refresh_token" varchar NOT null,
  "user_agent" varchar NOT NULL,
  "client_ip" varchar NOT NULL,
  "is_blocked" boolean NOT NULL default false,
  "expires_at" timestamp NOT NULL,    
  "created_at" timestamp NOT NULL default (now())
);

alter table "sessions" add foreign key ("username") REFERENCES "users"("username");