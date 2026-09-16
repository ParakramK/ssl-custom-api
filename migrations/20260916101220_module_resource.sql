-- Create "resources" table
CREATE TABLE "public"."resources" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "module_id" uuid NOT NULL,
  "name" character varying(100) NOT NULL,
  "code" character varying(100) NOT NULL,
  "description" text NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_resources_module" FOREIGN KEY ("module_id") REFERENCES "public"."modules" ("id") ON UPDATE RESTRICT ON DELETE RESTRICT
);
-- Create index "idx_resources_module_id" to table: "resources"
CREATE INDEX "idx_resources_module_id" ON "public"."resources" ("module_id");
