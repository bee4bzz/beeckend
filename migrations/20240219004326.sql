-- Modify "albums" table
ALTER TABLE "public"."albums" ALTER COLUMN "name" SET NOT NULL;
-- Modify "cheptel_notes" table
ALTER TABLE "public"."cheptel_notes" ALTER COLUMN "name" SET NOT NULL, ALTER COLUMN "flora" SET NOT NULL;
-- Modify "hive_notes" table
ALTER TABLE "public"."hive_notes" ALTER COLUMN "name" SET NOT NULL, ALTER COLUMN "nb_risers" SET NOT NULL, ALTER COLUMN "operation" SET NOT NULL;
-- Modify "hives" table
ALTER TABLE "public"."hives" ALTER COLUMN "name" SET NOT NULL;
-- Modify "users" table
ALTER TABLE "public"."users" ALTER COLUMN "name" SET NOT NULL, ALTER COLUMN "email" SET NOT NULL, ALTER COLUMN "password" SET NOT NULL;
-- Modify "user_cheptels" table
ALTER TABLE "public"."user_cheptels" DROP CONSTRAINT "fk_user_cheptels_cheptel", DROP CONSTRAINT "fk_user_cheptels_user", ADD
 CONSTRAINT "fk_user_cheptels_cheptel" FOREIGN KEY ("cheptel_id") REFERENCES "public"."cheptels" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, ADD
 CONSTRAINT "fk_user_cheptels_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
