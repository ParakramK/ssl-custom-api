-- Create "modules" table
CREATE TABLE "public"."modules" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "name" character varying(100) NOT NULL,
  "code" character varying(100) NOT NULL,
  "description" text NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_modules_code" to table: "modules"
CREATE UNIQUE INDEX "idx_modules_code" ON "public"."modules" ("code");
-- Create index "idx_modules_name" to table: "modules"
CREATE UNIQUE INDEX "idx_modules_name" ON "public"."modules" ("name");
