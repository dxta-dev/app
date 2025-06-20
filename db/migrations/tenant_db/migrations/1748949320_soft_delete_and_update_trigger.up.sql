ALTER TABLE "organizations" ADD COLUMN deleted_at DATETIME;
ALTER TABLE "github_organizations" ADD COLUMN deleted_at DATETIME;
DROP TABLE IF EXISTS "settings";

CREATE TRIGGER "organizations_set_updated_at"
AFTER UPDATE ON "organizations"
FOR EACH ROW
BEGIN
  UPDATE "organizations" 
  SET updated_at = datetime('now')
  WHERE id = OLD.id;
END;

CREATE TRIGGER "github_organizations_set_updated_at"
AFTER UPDATE ON "github_organizations"
FOR EACH ROW
BEGIN
  UPDATE "github_organizations" 
  SET updated_at = datetime('now')
  WHERE id = OLD.id;
END;

CREATE TRIGGER "organizations_github_organizations_set_updated_at"
AFTER UPDATE ON "organizations_github_organizations"
FOR EACH ROW
BEGIN
 UPDATE "organizations_github_organizations" 
 SET updated_at = datetime('now') 
 WHERE organization_id = OLD.organization_id
 AND github_app_installation_id = OLD.github_app_installation_id;
END;