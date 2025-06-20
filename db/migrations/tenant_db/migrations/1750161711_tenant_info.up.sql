CREATE TABLE "tenant_info" (
    "id" INTEGER PRIMARY KEY NOT NULL,
    "name" TEXT NOT NULL,
    "domain" TEXT NOT NULL,
    "created_at" DATETIME NOT NULL DEFAULT (datetime('now')),
    "updated_at" DATETIME NOT NULL DEFAULT (datetime('now')),
    "deleted_at" DATETIME DEFAULT NULL
);

CREATE TRIGGER "tenant_info_set_updated_at"
AFTER UPDATE ON "tenant_info"
FOR EACH ROW
BEGIN
  UPDATE "tenant_info" 
  SET updated_at = datetime('now')
  WHERE id = OLD.id;
END;