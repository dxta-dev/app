CREATE TABLE users_waitlist (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_email TEXT NOT NULL,
  repository_url TEXT NOT NULL,
  __created_at INTEGER DEFAULT ((strftime('%s', 'now') || substr(printf('%.3f', julianday('now') - julianday('1970-01-01')), 4, 3)) * 1),
  __updated_at INTEGER DEFAULT ((strftime('%s', 'now') || substr(printf('%.3f', julianday('now') - julianday('1970-01-01')), 4, 3)) * 1),
  UNIQUE (user_email, repository_url)
)
