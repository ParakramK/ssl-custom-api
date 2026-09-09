-- Create "api_permissions" table
CREATE TABLE "public"."api_permissions" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "api_key_id" uuid NOT NULL,
  "resource_code" character varying(100) NOT NULL,
  "permission" character varying(20) NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_api_permissions_api_key" FOREIGN KEY ("api_key_id") REFERENCES "public"."api_keys" ("id") ON UPDATE CASCADE ON DELETE CASCADE
);
-- Create index "uq_api_key_perm" to table: "api_permissions"
CREATE UNIQUE INDEX "uq_api_key_perm" ON "public"."api_permissions" ("api_key_id", "resource_code", "permission");
-- Create "user_module_permissions" table
CREATE TABLE "public"."user_module_permissions" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "user_id" uuid NOT NULL,
  "module_id" uuid NOT NULL,
  "resource_code" character varying(100) NOT NULL DEFAULT '',
  "permission" character varying(20) NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_user_module_permissions_module" FOREIGN KEY ("module_id") REFERENCES "public"."modules" ("id") ON UPDATE CASCADE ON DELETE CASCADE,
  CONSTRAINT "fk_user_module_permissions_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE CASCADE ON DELETE CASCADE
);
-- Create index "uq_user_module_perm" to table: "user_module_permissions"
CREATE UNIQUE INDEX "uq_user_module_perm" ON "public"."user_module_permissions" ("user_id", "module_id", "resource_code", "permission");
