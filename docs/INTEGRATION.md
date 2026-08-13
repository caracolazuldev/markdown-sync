Integration: Google Docs service account setup
===========================================

This file describes a minimal setup to run integration tests and allow the
CLI to update a Google Doc using a service account.

1. Enable the Google Docs API

- Go to GCP Console → APIs & Services → Library
- Search for "Google Docs API" and click Enable for your project.

2. Create a service account and download a JSON key

- Go to IAM & Admin → Service accounts → Create service account
- Give it a name and optional description
- After creation, create a JSON key and download it (store securely)

3. Share the target Google Doc with the service account

- Open the target Google Doc in the browser.
- Click "Share" and add the service account email (e.g. `svc-name@PROJECT.iam.gserviceaccount.com`)
- Give the service account the "Editor" role on the document.

4. Run the integration test (local)

Set the following environment variables:

```
export GOOGLE_APPLICATION_CREDENTIALS=/path/to/sa-key.json
export INTEGRATION_DOC_ID=1aB2cD3EfGhiJkLmnopQRsTUvWXyz
```

Run the single integration test:

```
go test -tags=integration ./cmd/gdocs-markdown-sync -run TestApplyDocumentIntegration
```

5. Run the CLI against the doc (example)

Export a markdown file and apply it (single-tab / legacy body only; tabbed Docs fail `export`):

```
./cmd/gdocs-markdown-sync export -auth service -doc $INTEGRATION_DOC_ID -out /tmp/out.md
./cmd/gdocs-markdown-sync import -auth service -file /path/to/local.md -doc $INTEGRATION_DOC_ID
```

Optional: one-way `track` of a multi-tab fixture (not required for unit tests):

```
export INTEGRATION_TABBED_DOC_ID=1aB2cD3EfGhiJkLmnopQRsTUvWXyz
./bin/gdocs-markdown-sync track -auth service -doc $INTEGRATION_TABBED_DOC_ID -out /tmp/tracked
```

Share the tabbed Doc with the service account as in step 3. `track` only reads the Doc.

Security notes
--------------
- Keep the JSON key secret and avoid committing it to git.
- Prefer using short-lived credentials or Workload Identity where possible.
