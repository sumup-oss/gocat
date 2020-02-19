# Contributing

## Prerequisites

* **Signed and verified CLA**
* Golang 1.13.x
* (To run linter) https://github.com/golangci/golangci-lint in `$PATH`

## Common commands

### Running the linter

```shell
> go run mage.go -v lint
```
### Running the tests

```shell
> go run mage.go -v test
```

### Running the benchmarks

```shell
> go run mage.go -v bench
```

### Running the benchmarks and generate graph

```shell
> go run mage.go -v benchandgraph
```

## Workflows

### Submitting an issue

1. Check existing issues and verify that your issue is not already submitted.
 If it is, it's highly recommended to add  to that issue with your reports.
 
2. Open issue

3. Be as detailed as possible - `go` version, what did you do, 
what did you expect to happen, what actually happened.

### Submitting a PR

1. Find an existing issue to work on or follow `Submitting an issue` to create one
 that you're also going to fix. 
 Make sure to notify that you're working on a fix for the issue you picked.
1. Branch out from latest `master`.
1. Code, add, commit and push your changes in your branch.
1. Make sure that tests and linter(s) pass locally for you.
1. Submit a PR.
1. Collaborate with the codeowners/reviewers to merge this in `master`.

### Releasing

#### Rules
    
1. Releases are only created from **tags**.
1. Tags are based on `master`.
1. `master` is meant to be stable, so before tagging a new release, make sure that the CI checks pass for `master`.
1. Releases and tags are following *semantic versioning*.
1. Releases are GitHub releases.
1. Releases and tags are to be named in pattern of `vX.Y.Z`. 
  The produced binary artifacts contain the `vX.Y.Z` in their names. 
  This is due how GitHub's CDNs work. They must be unique regardless if it's a different release.
1. Changelog must up-to-date with what's going to be released. Check [CHANGELOG](./CHANGELOG.md).
    
#### Flow

1. Push a new Git tag in pattern of `vX.Y.Z`.
1. Wait for GitHub action for `release` to workflow to kick in at https://github.com/sumup-oss/gocat/actions .
1. If all is good, a new GitHub release will be created and assets will be uploaded there.
