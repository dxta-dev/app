## Local Single-tenant setup

```toml
[tenants.company]
  database_url = "sqlite://company.db"
```

## Local Multi-tenant setup

```toml
[tenants.acme]
  database_url = "sqlite://acme.db"

[tenants.]
  database_url = "sqlite://acme.db"
