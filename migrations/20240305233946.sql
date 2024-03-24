-- Create "albums" table
CREATE TABLE "public"."albums" (
    "id" bigserial NOT NULL,
    "created_at" timestamptz NULL,
    "updated_at" timestamptz NULL,
    "deleted_at" timestamptz NULL,
    "name" text NOT NULL,
    "paths" text [] NOT NULL,
    "observation" text NULL,
    "owner_id" bigint NOT NULL,
    "owner_type" text NOT NULL,
    PRIMARY KEY ("id")
);
-- Create index "idx_albums_deleted_at" to table: "albums"
CREATE INDEX "idx_albums_deleted_at" ON "public"."albums" ("deleted_at");
-- Create index "idx_name_owner_id" to table: "albums"
CREATE UNIQUE INDEX "idx_name_owner_id" ON "public"."albums" (
    "name", "owner_id"
);
-- Create "cheptels" table
CREATE TABLE "public"."cheptels" (
    "id" bigserial NOT NULL,
    "created_at" timestamptz NULL,
    "updated_at" timestamptz NULL,
    "deleted_at" timestamptz NULL,
    "name" text NOT NULL,
    PRIMARY KEY ("id")
);
-- Create index "idx_cheptels_deleted_at" to table: "cheptels"
CREATE INDEX "idx_cheptels_deleted_at" ON "public"."cheptels" ("deleted_at");
-- Create "cheptel_notes" table
CREATE TABLE "public"."cheptel_notes" (
    "id" bigserial NOT NULL,
    "created_at" timestamptz NULL,
    "updated_at" timestamptz NULL,
    "deleted_at" timestamptz NULL,
    "cheptel_id" bigint NULL,
    "name" text NOT NULL,
    "temperature_day" numeric NULL,
    "temperature_night" numeric NULL,
    "weather" text NULL DEFAULT 'UNKNOWN',
    "flora" text NOT NULL,
    "state" text NULL,
    "observation" text NULL,
    PRIMARY KEY ("id"),
    CONSTRAINT "fk_cheptels_notes" FOREIGN KEY (
        "cheptel_id"
    ) REFERENCES "public"."cheptels" (
        "id"
    ) ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_cheptel_notes_deleted_at" to table: "cheptel_notes"
CREATE INDEX "idx_cheptel_notes_deleted_at" ON "public"."cheptel_notes" (
    "deleted_at"
);
-- Create "hives" table
CREATE TABLE "public"."hives" (
    "id" bigserial NOT NULL,
    "created_at" timestamptz NULL,
    "updated_at" timestamptz NULL,
    "deleted_at" timestamptz NULL,
    "name" text NOT NULL,
    "cheptel_id" bigint NOT NULL,
    PRIMARY KEY ("id"),
    CONSTRAINT "fk_cheptels_hives" FOREIGN KEY (
        "cheptel_id"
    ) REFERENCES "public"."cheptels" (
        "id"
    ) ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_hives_deleted_at" to table: "hives"
CREATE INDEX "idx_hives_deleted_at" ON "public"."hives" ("deleted_at");
-- Create index "idx_name_cheptel_id" to table: "hives"
CREATE UNIQUE INDEX "idx_name_cheptel_id" ON "public"."hives" (
    "name", "cheptel_id"
);
-- Create "hive_notes" table
CREATE TABLE "public"."hive_notes" (
    "id" bigserial NOT NULL,
    "created_at" timestamptz NULL,
    "updated_at" timestamptz NULL,
    "deleted_at" timestamptz NULL,
    "hive_id" bigint NOT NULL,
    "name" text NOT NULL,
    "nb_risers" bigint NOT NULL,
    "operation" text NOT NULL,
    "observation" text NULL,
    PRIMARY KEY ("id"),
    CONSTRAINT "fk_hives_notes" FOREIGN KEY (
        "hive_id"
    ) REFERENCES "public"."hives" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- Create index "idx_hive_notes_deleted_at" to table: "hive_notes"
CREATE INDEX "idx_hive_notes_deleted_at" ON "public"."hive_notes" (
    "deleted_at"
);
-- Create index "idx_name_hive_id" to table: "hive_notes"
CREATE UNIQUE INDEX "idx_name_hive_id" ON "public"."hive_notes" (
    "hive_id", "name"
);
-- Create "users" table
CREATE TABLE "public"."users" (
    "id" bigserial NOT NULL,
    "created_at" timestamptz NULL,
    "updated_at" timestamptz NULL,
    "deleted_at" timestamptz NULL,
    "name" text NOT NULL,
    "email" text NOT NULL,
    PRIMARY KEY ("id")
);
-- Create index "idx_users_deleted_at" to table: "users"
CREATE INDEX "idx_users_deleted_at" ON "public"."users" ("deleted_at");
-- Create index "uni_users_email" to table: "users"
CREATE UNIQUE INDEX "uni_users_email" ON "public"."users" ("email");
-- Create "user_cheptels" table
CREATE TABLE "public"."user_cheptels" (
    "user_id" bigint NOT NULL,
    "cheptel_id" bigint NOT NULL,
    PRIMARY KEY ("user_id", "cheptel_id"),
    CONSTRAINT "fk_user_cheptels_cheptel" FOREIGN KEY (
        "cheptel_id"
    ) REFERENCES "public"."cheptels" (
        "id"
    ) ON UPDATE NO ACTION ON DELETE CASCADE,
    CONSTRAINT "fk_user_cheptels_user" FOREIGN KEY (
        "user_id"
    ) REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
