DROP TRIGGER IF EXISTS "settings_set_updated_at";
DROP TRIGGER IF EXISTS "organizations_set_updated_at";
DROP TRIGGER IF EXISTS "teams_set_updated_at";
DROP TRIGGER IF EXISTS "members_set_updated_at";
DROP TRIGGER IF EXISTS "github_organizations_set_updated_at";
DROP TRIGGER IF EXISTS "github_members_set_updated_at";
DROP TRIGGER IF EXISTS "github_teams_set_updated_at";

DROP TABLE IF EXISTS "settings";
DROP TABLE IF EXISTS "organizations";
DROP TABLE IF EXISTS "teams";
DROP TABLE IF EXISTS "members";
DROP TABLE IF EXISTS "teams_members";
DROP TABLE IF EXISTS "github_organizations";
DROP TABLE IF EXISTS "organizations__github_organizations";
DROP TABLE IF EXISTS "github_members";
DROP TABLE IF EXISTS "github_teams";
DROP TABLE IF EXISTS "github_teams__github_members";