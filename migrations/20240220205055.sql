-- Modify "albums" table
ALTER TABLE "public"."albums" ADD COLUMN "paths" text[] NOT NULL;
