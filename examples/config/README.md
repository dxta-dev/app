## .toml configuration examples

This directory contains example .toml setups for configuring the runtime behaviour of the web app.

### Config file loading
When the `/cmd/web` app starts, it loads the `./config.toml` by default if the `CONFIG_PATH` environment variable evaluates to `""` (not set).

### Important Notes:

**_Note 1:_** The use of the configuration file `./config.toml` is discouraged in development. Since it is version-controlled, it may lead to the unintentional exposure of confidential data.

**_Note 2:_** During development, follow these guidelines for test configuration files:

- Test configuration files should adhere to the naming convention `config*.toml`, as specified in `.gitignore`.
- This practice helps minimize the risk of leaking confidential data when testing different configuration strategies.
