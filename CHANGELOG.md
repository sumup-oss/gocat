# Version History
 
## Intro
 
The version history is motivated by https://semver.org/ and https://keepachangelog.com/en/1.0.0/ .
 
## Structure
 
Types of changes that can be seen in the changelog
 
```
Added: for new features/functionality.
Changed: for changes in existing features/functionality.
Deprecated: for soon-to-be removed features. Removed in the
Removed: for now removed features.
Fixed: for any bug fixes.
Security: in case of vulnerabilities.
```
 
## How deprecation of functionality is handled?
 
tl;dr 1 minor release stating that the functionality is going to be deprecated. Then in the next major - removed.
 
```
Deprecating existing functionality is a normal part of software development and
is often required to make forward progress.
 
When you deprecate part of your public API, you should do two things:
 
(1) update your documentation to let users know about the change,
(2) issue a new minor release with the deprecation in place.
Before you completely remove the functionality in a new major
release there should be at least one minor release
that contains the deprecation so that users can smoothly transition to the new API
```
 
As per https://semver.org/ .
 
As per rule-of-thumb, moving the project forward is very important,
  but providing stability is the most important thing to anyone using `gocat`.
 
Introducing breaking changes under a feature flag can be ok in some cases where new functionality needs user feedback before being introduced in next major release.
 
## Changelog
 
Change line format:
 
```
* <Change title/PR title/content> ; Ref: <pr link>
```

## Unreleased (master)

### Added

* Added `tcp-to-tcp` cmd that relays TCP4/TCP6 to TCP4/TCP6 ; Ref: https://github.com/sumup-oss/gocat/pull/5

## v0.2.0

### Fixed

* Fixed context cancellation, noticeably for health checking ; Ref: https://github.com/sumup-oss/gocat/pull/3

### Changed

* Improved overall performance of `gocat` by building it with Golang 1.14.0. (https://golang.org/doc/go1.14#runtime) ; Ref: https://github.com/sumup-oss/gocat/pull/3

## v0.1.0

### Added

* Open-sourced the project
* Supports TCP-to-Unix relay
* Supports Unix-to-TCP relay
