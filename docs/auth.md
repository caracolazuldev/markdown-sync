# Authentication — markdown-sync

Supported flows
----------------
- OAuth2 (user consent / interactive): opens a browser to obtain user tokens. Tokens are cached locally under the user's config directory.
- Service account (headless): accepts a path to a JSON key; suitable for CI or server use.

Credentials & storage
---------------------
- Local tokens: stored at `$XDG_CONFIG_HOME/markdown-sync/tokens.json` with restricted file permissions.
- Service account: recommended to use `GOOGLE_APPLICATION_CREDENTIALS` or `--credentials` cli flag.

Scopes
------
- `https://www.googleapis.com/auth/documents.readonly` for export/preview
- `https://www.googleapis.com/auth/documents` for import/update
- Additional scopes may be documented later for Drive access if asset uploads are needed.

Security notes
--------------
- Do not check credentials into source control.
- For CI, use project-scoped service accounts with limited roles and rotate keys regularly.
