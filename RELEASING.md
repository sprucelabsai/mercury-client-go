# Releasing Mercury Client Go

Follow this checklist whenever you cut a new version so the Go module stays consumable from other projects.

## 1. Prep Your Environment
- Ensure the Mercury Theatre is running locally and reachable at the URL you plan to use (`http://127.0.0.1:8081` by default).
- Copy `.env` (or create one) so it contains `TEST_HOST=http://127.0.0.1:8081`.
- Confirm you can authenticate against the Theatre with the demo credentials used in the integration tests.

## 2. Verify the Client
- Run the full test suite: `go test ./...`.  
  The tests hit the live Theatre, so they will fail if `TEST_HOST` is wrong or the server is offline.
- Optionally run targeted checks (e.g. `go test -run TestFactory ./...`) while iterating.
- Review and update docs or examples that need to change for this release.

## 3. Pick a Version
- Use semantic versioning: bump `MAJOR` for breaking changes, `MINOR` for new features, `PATCH` for fixes.
- For the first public release, start at `v0.1.0` (pre-1.0 indicates the API can still change rapidly).
- Record highlights in `CHANGELOG.md` (create it if it does not exist yet).

## 4. Tag and Publish
```sh
git checkout main
git pull origin main
git tag -a vX.Y.Z -m "Release vX.Y.Z"
git push origin main
git push origin vX.Y.Z
```
- Replace `vX.Y.Z` with your chosen version.
- After pushing the tag, create a GitHub release (optional but recommended) so consumers see the notes.

## 5. Consume from Other Projects
- In the project that depends on this client, run:  
  `go get github.com/sprucelabsai/mercury-client-go@vX.Y.Z`
- Commit the updated `go.mod`/`go.sum` files in that consumer repo.
- When updating later, repeat the release process and bump the version in the consuming project.

## 6. Verify the Published Module
- Run `go list -m github.com/sprucelabsai/mercury-client-go@vX.Y.Z` to confirm the proxy can resolve the new tag.
- If the command fails, double-check that the tag is reachable (public repo or authenticated GOPRIVATE setup) and that CI finished indexing the tag on GitHub.

Following these steps keeps the Mercury Go client reproducible for downstream projects and makes it easy to automate releases later.
