CREATE TABLE "teams" (
    "id" INTEGER PRIMARY KEY NOT NULL,
    "name" TEXT NOT NULL,
    "organization_id" TEXT NOT NULL,
    "created_at" DATETIME NOT NULL DEFAULT (datetime('now')),
    "updated_at" DATETIME NOT NULL DEFAULT (datetime('now')),
    "deleted_at" DATETIME DEFAULT NULL,
    FOREIGN KEY ("organization_id") references "organizations" ("external_id") ON DELETE CASCADE
);

CREATE TABLE "team_members" (
    "id" INTEGER PRIMARY KEY NOT NULL,
    "team_id" INTEGER NOT NULL,
    "external_id" INTEGER NOT NULL,
    "username" TEXT NOT NULL,
    "created_at" DATETIME NOT NULL DEFAULT (datetime('now')),
    "updated_at" DATETIME NOT NULL DEFAULT (datetime('now')),
    "deleted_at" DATETIME DEFAULT NULL,
    FOREIGN KEY ("team_id") REFERENCES teams ("id") ON DELETE CASCADE
);

CREATE TRIGGER "teams_set_updated_at"
AFTER UPDATE ON "teams"
FOR EACH ROW
BEGIN
 UPDATE "teams" 
 SET updated_at = datetime('now') 
 WHERE id = OLD.id;
END;

CREATE TRIGGER "team_members_set_updated_at"
AFTER UPDATE ON "team_members"
FOR EACH ROW
BEGIN
 UPDATE "team_members" 
 SET updated_at = datetime('now') 
 WHERE id = OLD.id;
END;