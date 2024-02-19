-- Modify "user_cheptels" table
ALTER TABLE "public"."user_cheptels" DROP CONSTRAINT "fk_user_cheptels_cheptel", DROP CONSTRAINT "fk_user_cheptels_user", ADD
 CONSTRAINT "fk_user_cheptels_cheptel" FOREIGN KEY ("cheptel_id") REFERENCES "public"."cheptels" ("id") ON UPDATE CASCADE ON DELETE SET NULL, ADD
 CONSTRAINT "fk_user_cheptels_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE CASCADE ON DELETE SET NULL;
