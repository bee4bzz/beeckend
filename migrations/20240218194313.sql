-- Create "cheptels" table
CREATE TABLE "public"."cheptels" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_cheptels_deleted_at" to table: "cheptels"
CREATE INDEX "idx_cheptels_deleted_at" ON "public"."cheptels" ("deleted_at");
-- Create "users" table
CREATE TABLE "public"."users" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "name" text NULL,
  "email" text NULL,
  "password" text NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_users_deleted_at" to table: "users"
CREATE INDEX "idx_users_deleted_at" ON "public"."users" ("deleted_at");
-- Create "user_cheptels" table
CREATE TABLE "public"."user_cheptels" (
  "user_id" bigint NOT NULL,
  "cheptel_id" bigint NOT NULL,
  PRIMARY KEY ("user_id", "cheptel_id"),
  CONSTRAINT "fk_user_cheptels_cheptel" FOREIGN KEY ("cheptel_id") REFERENCES "public"."cheptels" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_user_cheptels_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
