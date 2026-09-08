-- Create "roles" table
CREATE TABLE "public"."roles" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" character varying(255) NOT NULL,
  "is_admin" boolean NOT NULL DEFAULT false,
  PRIMARY KEY ("id")
);
-- Create index "idx_roles_name" to table: "roles"
CREATE UNIQUE INDEX "idx_roles_name" ON "public"."roles" ("name");
-- Modify "users" table
ALTER TABLE "public"."users" ADD COLUMN "role_id" uuid NOT NULL, ADD CONSTRAINT "fk_users_role" FOREIGN KEY ("role_id") REFERENCES "public"."roles" ("id") ON UPDATE RESTRICT ON DELETE RESTRICT;
