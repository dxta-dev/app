CREATE TABLE IF NOT EXISTS "organizations" (
    "id" INTEGER PRIMARY KEY NOT NULL,
    "external_id" TEXT NOT NULL UNIQUE,
    "created_at" DATETIME NOT NULL DEFAULT (datetime('now')),
    "updated_at" DATETIME NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS "organizations_github_organizations" (
	"organization_id" TEXT NOT NULL,
	"github_app_installation_id" TEXT NOT NULL,
    "created_at" DATETIME NOT NULL DEFAULT (datetime('now')),
    "updated_at" DATETIME NOT NULL DEFAULT (datetime('now')),
	 PRIMARY KEY ("organization_id", "github_app_installation_id"),
     FOREIGN KEY ("organization_id") REFERENCES "organizations"("external_id") ON DELETE CASCADE,
     FOREIGN KEY ("github_app_installation_id") REFERENCES "github_organizations"("github_app_installation_id") ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS "github_organizations" (
    "id" INTEGER PRIMARY KEY NOT NULL,
    "name" TEXT,
    "github_app_installation_id" TEXT NOT NULL UNIQUE,
    "created_at" DATETIME NOT NULL DEFAULT (datetime('now')),
    "updated_at" DATETIME NOT NULL DEFAULT (datetime('now'))
    
);

CREATE TABLE IF NOT EXISTS "settings" (
    "id" TEXT PRIMARY KEY NOT NULL,
    "created_at" DATETIME NOT NULL DEFAULT (datetime('now')),
    "updated_at" DATETIME NOT NULL DEFAULT (datetime('now'))
);