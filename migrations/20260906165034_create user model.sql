-- Create "users" table
CREATE TABLE "public"."users" (
  "id" uuid NOT NULL DEFAULT gen_random_uuid(),
  "username" character varying(100) NOT NULL,
  "email" character varying(100) NOT NULL,
  "password" character varying(255) NOT NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_users_email" UNIQUE ("email"),
  CONSTRAINT "uni_users_username" UNIQUE ("username")
);
