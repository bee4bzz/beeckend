-- reverse: create "user_cheptels" table
DROP TABLE "public"."user_cheptels";
-- reverse: create index "idx_users_deleted_at" to table: "users"
DROP INDEX "public"."idx_users_deleted_at";
-- reverse: create "users" table
DROP TABLE "public"."users";
-- reverse: create index "idx_hive_note_photos_deleted_at" to table: "hive_note_photos"
DROP INDEX "public"."idx_hive_note_photos_deleted_at";
-- reverse: create "hive_note_photos" table
DROP TABLE "public"."hive_note_photos";
-- reverse: create index "idx_hive_note_albums_key" to table: "hive_note_albums"
DROP INDEX "public"."idx_hive_note_albums_key";
-- reverse: create index "idx_hive_note_albums_deleted_at" to table: "hive_note_albums"
DROP INDEX "public"."idx_hive_note_albums_deleted_at";
-- reverse: create "hive_note_albums" table
DROP TABLE "public"."hive_note_albums";
-- reverse: create index "idx_name_hive_id" to table: "hive_notes"
DROP INDEX "public"."idx_name_hive_id";
-- reverse: create index "idx_hive_notes_deleted_at" to table: "hive_notes"
DROP INDEX "public"."idx_hive_notes_deleted_at";
-- reverse: create "hive_notes" table
DROP TABLE "public"."hive_notes";
-- reverse: create index "idx_name_cheptel_id" to table: "hives"
DROP INDEX "public"."idx_name_cheptel_id";
-- reverse: create index "idx_hives_deleted_at" to table: "hives"
DROP INDEX "public"."idx_hives_deleted_at";
-- reverse: create "hives" table
DROP TABLE "public"."hives";
-- reverse: create index "idx_cheptel_photos_deleted_at" to table: "cheptel_photos"
DROP INDEX "public"."idx_cheptel_photos_deleted_at";
-- reverse: create "cheptel_photos" table
DROP TABLE "public"."cheptel_photos";
-- reverse: create index "idx_cheptel_notes_deleted_at" to table: "cheptel_notes"
DROP INDEX "public"."idx_cheptel_notes_deleted_at";
-- reverse: create "cheptel_notes" table
DROP TABLE "public"."cheptel_notes";
-- reverse: create index "idx_cheptel_albums_key" to table: "cheptel_albums"
DROP INDEX "public"."idx_cheptel_albums_key";
-- reverse: create index "idx_cheptel_albums_deleted_at" to table: "cheptel_albums"
DROP INDEX "public"."idx_cheptel_albums_deleted_at";
-- reverse: create "cheptel_albums" table
DROP TABLE "public"."cheptel_albums";
-- reverse: create index "idx_cheptels_deleted_at" to table: "cheptels"
DROP INDEX "public"."idx_cheptels_deleted_at";
-- reverse: create "cheptels" table
DROP TABLE "public"."cheptels";
-- reverse: create "tokens" table
DROP TABLE "public"."tokens";
-- reverse: create index "idx_albums_key" to table: "albums"
DROP INDEX "public"."idx_albums_key";
-- reverse: create "albums" table
DROP TABLE "public"."albums";
-- reverse: create index "idx_refresh_tokens_deleted_at" to table: "refresh_tokens"
DROP INDEX "public"."idx_refresh_tokens_deleted_at";
-- reverse: create "refresh_tokens" table
DROP TABLE "public"."refresh_tokens";
