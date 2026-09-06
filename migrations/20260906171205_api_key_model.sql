-- Create "api_keys" table
CREATE TABLE "public"."api_keys" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "user_id" uuid NOT NULL,
  "key" character varying(255) NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_api_keys_user" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE CASCADE ON DELETE CASCADE
);
-- Create index "idx_api_keys_key" to table: "api_keys"
CREATE UNIQUE INDEX "idx_api_keys_key" ON "public"."api_keys" ("key");
-- Create index "idx_api_keys_user_id" to table: "api_keys"
CREATE INDEX "idx_api_keys_user_id" ON "public"."api_keys" ("user_id");
