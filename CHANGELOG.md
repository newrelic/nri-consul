# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

Unreleased section should follow [Release Toolkit](https://github.com/newrelic/release-toolkit#render-markdown-and-update-markdown)
## Unreleased

## v2.9.2 - 2025-06-30

### ⛓️ Dependencies
- Updated golang version to v1.24.4
- Updated github.com/hashicorp/consul/api to v1.32.1

## v2.9.1 - 2025-05-26

### ⛓️ Dependencies
- Updated github.com/hashicorp/consul/api to v1.32.0

## v2.9.1 - 2025-04-07

### ⛓️ Dependencies
- Updated github.com/hashicorp/consul/api to v1.32.0

## v2.9.0 - 2025-03-17

### 🚀 Enhancements
- Add FIPS compliant packages

## v2.8.4 - 2025-03-10

### ⛓️ Dependencies
- Updated golang patch version to v1.23.6
- Updated github.com/hashicorp/consul/api to v1.31.2

## v2.8.3 - 2025-01-20

### ⛓️ Dependencies
- Updated github.com/hashicorp/consul/api to v1.31.0
- Updated golang patch version to v1.23.5

## v2.8.2 - 2024-12-09

### ⛓️ Dependencies
- Updated golang patch version to v1.23.4

## v2.8.1 - 2024-10-21

### ⛓️ Dependencies
- Updated github.com/hashicorp/consul/api to v1.30.0

## v2.8.0 - 2024-10-14

### dependency
- Upgrade go to 1.23.2

### 🚀 Enhancements
- Upgrade integrations SDK so the interval is variable and allows intervals up to 5 minutes

## v2.7.15 - 2024-09-16

### ⛓️ Dependencies
- Updated golang version to v1.23.1

## v2.7.14 - 2024-09-09

### ⛓️ Dependencies
- Updated github.com/hashicorp/consul/api to v1.29.4

## v2.7.13 - 2024-07-08

### ⛓️ Dependencies
- Updated golang version to v1.22.5

## v2.7.12 - 2024-06-03

### ⛓️ Dependencies
- Updated github.com/hashicorp/consul/api to v1.29.1

## v2.7.11 - 2024-05-20

### ⛓️ Dependencies
- Updated github.com/hashicorp/consul/api to v1.28.3

## v2.7.10 - 2024-05-13

### ⛓️ Dependencies
- Updated golang version to v1.22.3

## v2.7.9 - 2024-04-15

### ⛓️ Dependencies
- Updated golang version to v1.22.2

## v2.7.8 - 2024-03-07

### 🐞 Bug fixes
- Updated golang to version v1.21.7 to fix a vulnerability

## v2.7.7 - 2024-03-04

### 🐞 Bug fixes
- Fix Windows packaging signtool path

### ⛓️ Dependencies
- Updated github.com/hashicorp/consul/api to v1.28.2

## v2.7.6 - 2024-02-26

### ⛓️ Dependencies
- Updated github.com/newrelic/infra-integrations-sdk to v3.8.2+incompatible

## v2.7.5 - 2024-02-12

### ⛓️ Dependencies
- Updated github.com/newrelic/infra-integrations-sdk to v3.8.0+incompatible

## v2.7.4 - 2024-01-22

### ⛓️ Dependencies
- Updated github.com/hashicorp/consul/api to v1.27.0

## v2.7.3 - 2023-11-06

### ⛓️ Dependencies
- Updated github.com/hashicorp/consul/api to v1.26.1

## v2.7.2 - 2023-09-25

### ⛓️ Dependencies
- Updated github.com/hashicorp/consul/api to v1.25.1

## v2.7.1 - 2023-08-07

### ⛓️ Dependencies
- Updated golang to v1.20.7
- Updated github.com/hashicorp/consul/api to v1.24.0

## v2.7.0 - 2023-07-24

### 🚀 Enhancements
- bumped golang version pining 1.20.6

### ⛓️ Dependencies
- Updated github.com/hashicorp/consul/api to v1.23.0

## 2.6.0 (2023-06-06)
### Changed
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
