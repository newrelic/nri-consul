# Change Log

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

## 2.6.0 (2023-06-06)
# Changed
- Upgrade Go version to 1.20

## 2.5.1 (2022-06-27)
### Changed
- Bump dependencies

### Added
Added support for more distributions:
- RHEL(EL) 9
- Ubuntu 22.04

## 2.5.0 (2022-05-18)
### Added
Added TIMEOUT for the client.

## 2.4.1 (2021-10-20)
### Added
Added support for more distributions:
- Debian 11
- Ubuntu 20.10
- Ubuntu 21.04
- SUSE 12.15
- SUSE 15.1
- SUSE 15.2
- SUSE 15.3
- Oracle Linux 7
- Oracle Linux 8

## 2.4.0 (2021-08-30)
### Changed
- Moved default config.sample to [V4](https://docs.newrelic.com/docs/create-integrations/infrastructure-integrations-sdk/specifications/host-integrations-newer-configuration-format/), added a dependency for infra-agent version 1.20.0

Please notice that old [V3](https://docs.newrelic.com/docs/create-integrations/infrastructure-integrations-sdk/specifications/host-integrations-standard-configuration-format/) configuration format is deprecated, but still supported.

## 2.3.1 (2021-06-09)
### Changed
- Added ARM support.

## 2.3.0 (2021-05-03)
### Fixed
- Bump to integrations-sdk to v3.6.7 containing fix for multiple instance sharing storer.
### Changed
- Migrate to gomod and go v1.16.
- Update CI to support go mod and bump go version.
- Bump non core dependencies to last minor verison.

## 2.2.0 (2021-03-31)
### Added
- Arm and Arm64 packages for Linux
### Changed
- The CI pipeline has been migrated to Github Actions

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
