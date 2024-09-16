CREATE TABLE repos (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  organization TEXT NOT NULL,
  repository TEXT NOT NULL,
  db_url TEXT NOT NULL,
  __created_at INTEGER DEFAULT ((strftime('%s', 'now') || substr(printf('%.3f', julianday('now') - julianday('1970-01-01')), 4, 3)) * 1),
  __updated_at INTEGER DEFAULT ((strftime('%s', 'now') || substr(printf('%.3f', julianday('now') - julianday('1970-01-01')), 4, 3)) * 1)
  UNIQUE(organization, repository)
);

CREATE INDEX idx_org_repo_repos ON repos (organization, repository);

CREATE TRIGGER repos_update_updated_at after update on repos BEGIN
  UPDATE repos SET __updated_at = ((strftime('%s', 'now') || substr(printf('%.3f', julianday('now') - julianday('1970-01-01')), 4, 3)) * 1) WHERE id = old.id;
END;
