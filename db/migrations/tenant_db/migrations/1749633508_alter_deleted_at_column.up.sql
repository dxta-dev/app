ALTER TABLE organizations ALTER COLUMN deleted_at TO deleted_at DATETIME DEFAULT NULL;
ALTER TABLE github_organizations ALTER COLUMN deleted_at TO deleted_at DATETIME DEFAULT NULL;