# Release process

GitHub Actions are used as a backbone to get releases going.

## Rules

1. Releases are only created from `master`.
1. `master` is meant to be stable, so before tagging and pushing a tag, make sure that the CI checks pass.
1. Releases are GitHub releases.
1. Releases are following *semantic versioning*.
1. Releases are to be named in pattern of `vX.Y.Z`. The produced binary artifacts contain the `vX.Y.Z` in their names.
1. Changelog must up-to-date with what's going to be released. Check [CHANGELOG](./CHANGELOG.md).

## Flow

1. Create a new GitHub a new tag from `master`
1. Push it to the remote git repository.
1. Wait for GitHub action workflow to finish

