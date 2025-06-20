CREATE TABLE IF NOT EXISTS "tenants" (
    "id" INTEGER PRIMARY KEY NOT NULL,
    "organization_id" TEXT NOT NULL UNIQUE,
    "db_url" TEXT NOT NULL,
    "deleted_at" DATETIME,
    "created_at" DATETIME NOT NULL DEFAULT (datetime('now')),
    "updated_at" DATETIME NOT NULL DEFAULT (datetime('now'))
);

CREATE TRIGGER "tenants_updated_at"
AFTER UPDATE ON "tenants" 
FOR EACH ROW 
BEGIN
UPDATE "tenants"
SET updated_at = datetime('now')
WHERE id = OLD.id;
END;