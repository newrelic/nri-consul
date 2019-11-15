# New Relic Infrastructure Integration for HashiCorp Consul

The New Relic Infrastructure Integration for HashiCorp Consul captures critical performance metrics and inventory reported by Consul clusters. Data on Agents and the Datacenter as a whole is collected.

All data is obtained via the REST API.

*HashiCorp components* The on host integration contains open source software (in unmodified form) provided by HashiCorp, distributed under the Mozilla Public License 2.0 (https://www.mozilla.org/en-US/MPL/2.0/FAQ/.).  Please see [LICENSE](LICENSE) for additional information.  Do not remove the license file.  Mozilla Public License 2.0 is a copyleft light open source license, if you do not want to be subject the copyleft light open source provisions, do not modify the MPL 2.0 files.

All other components of this on host integration are licensed under an MIT license, Copyright 2019, Blue Medora Inc. Please see [LICENSE](LICENSE) for additional information.  

See our [documentation web site](https://docs.newrelic.com/docs/integrations/host-integrations/host-integrations-list/consul-monitoring-integration) for more details.

## Requirements

No requirements at this time.

## Installation

- download an archive file for the `Consul` Integration
- extract `consul-definition.yml` and `/bin` directory into `/var/db/newrelic-infra/newrelic-integrations`
- add execute permissions for the binary file `nri-consul` (if required)
- extract `consul-config.yml.sample` into `/etc/newrelic-infra/integrations.d`

## Usage

This is the description about how to run the Consul Integration with New Relic Infrastructure agent, so it is required to have the agent installed (see [agent installation](https://docs.newrelic.com/docs/infrastructure/new-relic-infrastructure/installation/install-infrastructure-linux)).

In order to use the Consul Integration it is required to configure `consul-config.yml.sample` file. Firstly, rename the file to `consul-config.yml`. Then, depending on your needs, specify all instances that you want to monitor. Once this is done, restart the Infrastructure agent.

You can view your data in Insights by creating your own custom NRQL queries. To do so use the **ConsulDatacenterSample** and **ConsulAgentSample** event type.

## Compatibility

* Supported OS: No limitations
* Consul versions: 1.0+

## Integration Development usage

Assuming that you have source code you can build and run the Consul Integration locally.

* Go to directory of the Consul Integration and build it
```bash
$ make
```
* The command above will execute tests for the Consul Integration and build an executable file called `nri-consul` in `bin` directory.
```bash
$ ./bin/nri-consul
```
* If you want to know more about usage of `./nri-consul` check
```bash
$ ./bin/nri-consul -help
```

For managing external dependencies [govendor tool](https://github.com/kardianos/govendor) is used. It is required to lock all external dependencies to specific version (if possible) into vendor directory.
