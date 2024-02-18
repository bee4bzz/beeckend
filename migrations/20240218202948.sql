-- Create "albums" table
CREATE TABLE "public"."albums" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "name" text NULL,
  "observation" text NULL,
  "owner_id" bigint NULL,
  "owner_type" text NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_albums_deleted_at" to table: "albums"
CREATE INDEX "idx_albums_deleted_at" ON "public"."albums" ("deleted_at");
-- Create "cheptel_notes" table
CREATE TABLE "public"."cheptel_notes" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "cheptel_id" bigint NULL,
  "name" text NULL,
  "temperature_day" numeric NULL,
  "temperature_night" numeric NULL,
  "weather" text NULL,
  "flora" text NULL,
  "state" text NULL,
  "observation" text NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_cheptels_notes" FOREIGN KEY ("cheptel_id") REFERENCES "public"."cheptels" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_cheptel_notes_deleted_at" to table: "cheptel_notes"
CREATE INDEX "idx_cheptel_notes_deleted_at" ON "public"."cheptel_notes" ("deleted_at");
-- Create "hives" table
CREATE TABLE "public"."hives" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "name" text NULL,
  "cheptel_id" bigint NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_cheptels_hives" FOREIGN KEY ("cheptel_id") REFERENCES "public"."cheptels" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_hives_deleted_at" to table: "hives"
CREATE INDEX "idx_hives_deleted_at" ON "public"."hives" ("deleted_at");
-- Create "hive_notes" table
CREATE TABLE "public"."hive_notes" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  "hive_id" bigint NULL,
  "name" text NULL,
  "nb_risers" bigint NULL,
  "operation" text NULL,
  "observation" text NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_hives_notes" FOREIGN KEY ("hive_id") REFERENCES "public"."hives" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- Create index "idx_hive_notes_deleted_at" to table: "hive_notes"
CREATE INDEX "idx_hive_notes_deleted_at" ON "public"."hive_notes" ("deleted_at");
