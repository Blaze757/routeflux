# Routeflux Development Guidelines

## Code Quality

### Go formatting
Before committing any Go code changes, always run:
  gofmt -l internal/ pkg/ cmd/ test/
Fix any files listed before committing. Never push unformatted code.

## CI Checklist

After pushing, verify CI via API:
  gh api repos/Blaze757/routeflux/actions/runs?per_page=1 --jq '.workflow_runs[0].conclusion'
Expect "success". If "failure", investigate before proceeding.

## CI Workflow Rules

The `.github/workflows/ci.yml` must be clean production state:
- `permissions: contents: read`
- No `github-script` steps
- No `contents: write`
- No artifact commits to repo
If you see debug artifacts, remove them before pushing.

## Project Structure

- `internal/cli/` — CLI commands (cobra)
- `internal/probe/` — node latency probing
- `internal/store/` — local state management
- `internal/platform/openwrt/` — OpenWRT-specific logic
- `pkg/sub/` — subscription parser
- `cmd/routeflux/` — entrypoint

## Security

Never commit secrets, keys, or tokens.
Validators in `internal/platform/openwrt/security.go` are critical — do not bypass.
