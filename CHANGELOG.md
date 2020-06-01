# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.0.0] - 2020-06-01

A functional redesign: by default most methods that return a bitfield return a new
bitfield and leave the original intact. This however can be changed with calling
Mut(): It will make replace in-place as default.
### Added
- Mut

### Changed
- Copy - it is now a bit-copy from src to destination.

## [1.2.0] - 2020-05-31
### Changed
- clearEnd: speed improvement affects the overal speed of the package

## [1.1.0] - 2020-05-28
### Added
- Resize
- BitCopy
- Shift
- Rotate
- Left
- Right
- Mid
- Append
- SetMul
- ClearMul

## Changed
- New: in case len<=0 it no longer returns nil. It returns a Len(0) BitField.
- bumped up bitfield64 version number

## Deprecated
- Copy: this method will be renamed to Clone for added clarity.

## [1.0.0] - 2020-05-23

First release.
