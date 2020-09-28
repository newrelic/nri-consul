# Change Log

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

## 2.1.2 (2020-09-28)
### Fixed
- Added a fallback for leader detection on old versions

## 2.1.1 (2020-05-04)
### Changed
- Updated the Consul API library

## 2.1.0 (2019-11-18)
### Changed
- Renamed the integration executable from nr-consul to nri-consul in order to be consistent with the package naming. **Important Note:** if you have any security module rules (eg. SELinux), alerts or automation that depends on the name of this binary, these will have to be updated.

## 2.0.4 - 2019-11-13
### Fixed
- Use unique component GUIDs

## 2.0.3 - 2019-10-23
### Added
- Windows installer packaging

## 2.0.2 - 2019-07-17
### Fixed
- Use agent name for latency calculations

## 2.0.0 - 2019-04-25
### Changed
- Prefixed namespaces for uniqueness
- Updated SDK
- Added ID attributes

## 1.1.1 - 2019-04-17
### Added
- Use address rather than name to connect

## 1.1.0 - 2019-04-08
### Added
- Local-only collection option

## 1.0.1 - 2019-03-19
### Fixed
- Stop failing when the leader can't be contacted

## 0.1.2 - 2019-03-14
### Fixed
- Stop failing when the leader can't be contacted

## 0.1.1 - 2018-11-15
### Added
- Datacenter and IP attributes to all Agent samples

## 0.1.0 - 2018-11-14
### Added
- Initial version: Includes Metrics and Inventory data
