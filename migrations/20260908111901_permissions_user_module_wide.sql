-- Drop index "uq_user_module_perm" from table: "user_module_permissions"
DROP INDEX "public"."uq_user_module_perm";
-- Modify "user_module_permissions" table
ALTER TABLE "public"."user_module_permissions" DROP COLUMN "resource_code";
-- Create index "uq_user_module_perm" to table: "user_module_permissions"
CREATE UNIQUE INDEX "uq_user_module_perm" ON "public"."user_module_permissions" ("user_id", "module_id", "permission");
