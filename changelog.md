# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](http://semver.org/spec/v2.0.0.html).

## [unreleased]

- Added -min flag, uncover fails if below

## [0.3.0] 2020-03-14
### Changed

- Shorten multiline funcs to oneline
- Only use two colors, either it's covered or not.
- Special print uncovered funcs, hiding their body
- Using go/printer to print the entire func signature

### Fixed

- Signature panic when having small aligned funcs

## [0.2.0] 2020-01-04
### Changed

- Never show func main as it cannot be covered
- Use tabwriter for output to minimize width of result

## [0.1.1] 2019-03-26
### Changed

- Using go mod

## [0.1.0] 2018-11-22
### Added

- Show only coverage for named function
- Colorized output to stdout by func, only uncovered are shown
