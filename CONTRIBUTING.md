# Contributing

## Prerequisites

* **Signed and verified CLA**
* Golang 1.12.x
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
