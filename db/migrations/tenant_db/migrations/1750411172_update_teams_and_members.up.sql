CREATE TABLE "external_teams" (
    "id" INTEGER PRIMARY KEY NOT NULL,
    "external_id" INTEGER NOT NULL,
    "name" TEXT NOT NULL,
    "organization_id" TEXT NOT NULL,
    "created_at" DATETIME NOT NULL DEFAULT (datetime('now')),
    "updated_at" DATETIME NOT NULL DEFAULT (datetime('now')),
    "deleted_at" DATETIME DEFAULT NULL,
    FOREIGN KEY ("organization_id") references "organizations" ("external_id") ON DELETE CASCADE
);

CREATE TABLE "external_members" (
    "id" INTEGER PRIMARY KEY NOT NULL,
    "external_id" INTEGER NOT NULL,
    "username" TEXT NOT NULL,
    "email" TEXT DEFAULT NULL,
    "created_at" DATETIME NOT NULL DEFAULT (datetime('now')),
    "updated_at" DATETIME NOT NULL DEFAULT (datetime('now')),
    "deleted_at" DATETIME DEFAULT NULL
);

CREATE TABLE "external_teams_members" (
    "member_id" INTEGER NOT NULL,
    "team_id" INTEGER NOT NULL,
    "created_at" DATETIME NOT NULL DEFAULT (datetime('now')),
    "updated_at" DATETIME NOT NULL DEFAULT (datetime('now')),
	PRIMARY KEY ("member_id", "team_id"),
    FOREIGN KEY ("team_id") REFERENCES "external_teams" ("id") ON DELETE CASCADE,
    FOREIGN KEY ("member_id") REFERENCES "external_members" ("id") ON DELETE CASCADE
);

CREATE TABLE "custom_teams" (
    "id" INTEGER PRIMARY KEY NOT NULL,
    "name" TEXT NOT NULL,
    "organization_id" TEXT NOT NULL,
    "created_at" DATETIME NOT NULL DEFAULT (datetime('now')),
    "updated_at" DATETIME NOT NULL DEFAULT (datetime('now')),
    "deleted_at" DATETIME DEFAULT NULL
);

CREATE TABLE "custom_members" (
    "id" INTEGER PRIMARY KEY NOT NULL,
    "external_id" INTEGER NOT NULL,
    "username" TEXT NOT NULL,
    "email" TEXT DEFAULT NULL,
    "created_at" DATETIME NOT NULL DEFAULT (datetime('now')),
    "updated_at" DATETIME NOT NULL DEFAULT (datetime('now')),
    "deleted_at" DATETIME DEFAULT NULL
);

CREATE TABLE "custom_teams_members" (
    "member_id" INTEGER NOT NULL,
    "team_id" INTEGER NOT NULL,
    "created_at" DATETIME NOT NULL DEFAULT (datetime('now')),
    "updated_at" DATETIME NOT NULL DEFAULT (datetime('now')),
	PRIMARY KEY ("member_id", "team_id"),
    FOREIGN KEY ("team_id") REFERENCES "custom_teams" ("id") ON DELETE CASCADE,
    FOREIGN KEY ("member_id") REFERENCES "custom_members" ("id") ON DELETE CASCADE
);

CREATE TABLE "custom_external_teams" (
    "custom_team_id" INTEGER NOT NULL,
    "external_team_id" INTEGER NOT NULL,
    PRIMARY KEY ("custom_team_id","external_team_id"),
    FOREIGN KEY ("custom_team_id") REFERENCES "custom_teams" ("id") ON DELETE CASCADE,
    FOREIGN KEY ("external_team_id") REFERENCES "external_teams" ("id") ON DELETE CASCADE
);

CREATE TRIGGER "external_teams_set_updated_at"
AFTER UPDATE ON "external_teams"
FOR EACH ROW
BEGIN
 UPDATE "external_teams" 
 SET updated_at = datetime('now') 
 WHERE id = OLD.id;
END;

CREATE TRIGGER "external_members_set_updated_at"
AFTER UPDATE ON "external_members"
FOR EACH ROW
BEGIN
 UPDATE "external_members" 
 SET updated_at = datetime('now') 
 WHERE id = OLD.id;
END;

CREATE TRIGGER "external_teams_members_set_updated_at"
AFTER UPDATE ON "external_teams_members"
FOR EACH ROW
BEGIN
 UPDATE "external_teams_members" 
 SET updated_at = datetime('now') 
 WHERE id = OLD.id;
END;

CREATE TRIGGER "custom_teams_set_updated_at"
AFTER UPDATE ON "custom_teams"
FOR EACH ROW
BEGIN
 UPDATE "custom_teams" 
 SET updated_at = datetime('now') 
 WHERE id = OLD.id;
END;

CREATE TRIGGER "custom_members_set_updated_at"
AFTER UPDATE ON "custom_members"
FOR EACH ROW
BEGIN
 UPDATE "custom_members" 
 SET updated_at = datetime('now') 
 WHERE id = OLD.id;
END;

CREATE TRIGGER "custom_teams_members_set_updated_at"
AFTER UPDATE ON "custom_teams_members"
FOR EACH ROW
BEGIN
 UPDATE "custom_teams_members" 
 SET updated_at = datetime('now') 
 WHERE id = OLD.id;
END;