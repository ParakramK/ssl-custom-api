-- Normalize legacy rows to one level per (user, module) before tightening
-- the unique constraint. The old grant endpoint stored '' (empty
-- permission), which authorized exactly like 'read', so it ranks as read
-- and any surviving '' row is mapped to 'read' below.
-- Collapse duplicates first (otherwise mapping '' to 'read' would collide
-- with the old (user, module, permission) unique index), keeping the
-- highest level (admin > write > delete > read); break ties by keeping
-- one row.
DELETE FROM "public"."user_module_permissions" dup
USING "public"."user_module_permissions" keep
WHERE dup."user_id" = keep."user_id"
  AND dup."module_id" = keep."module_id"
  AND dup."id" <> keep."id"
  AND (
    CASE keep."permission" WHEN 'admin' THEN 4 WHEN 'write' THEN 3 WHEN 'delete' THEN 2 ELSE 1 END
    > CASE dup."permission" WHEN 'admin' THEN 4 WHEN 'write' THEN 3 WHEN 'delete' THEN 2 ELSE 1 END
    OR (
      CASE keep."permission" WHEN 'admin' THEN 4 WHEN 'write' THEN 3 WHEN 'delete' THEN 2 ELSE 1 END
      = CASE dup."permission" WHEN 'admin' THEN 4 WHEN 'write' THEN 3 WHEN 'delete' THEN 2 ELSE 1 END
      AND dup."id" > keep."id"
    )
  );
UPDATE "public"."user_module_permissions" SET "permission" = 'read' WHERE "permission" = '';
-- Drop index "uq_user_module_perm" from table: "user_module_permissions"
DROP INDEX "public"."uq_user_module_perm";
-- Create index "uq_user_module" to table: "user_module_permissions"
CREATE UNIQUE INDEX "uq_user_module" ON "public"."user_module_permissions" ("user_id", "module_id");
