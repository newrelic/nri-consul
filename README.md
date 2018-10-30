# New Relic Infrastructure Integration for HashiCorp Consul

The New Relic Infrastructure Integration for HashiCorp Consul captures critical performance metrics and inventory reported by Consul clusters. Data on Agents and the Raft Cluster as a whole is collected.

All data is obtained via the REST API.

## Requirements

No requirements at this time.

## Installation

- download an archive file for the `Consul` Integration
- extract `consul-definition.yml` and `/bin` directory into `/var/db/newrelic-infra/newrelic-integrations`
- add execute permissions for the binary file `nr-consul` (if required)
- extract `consul-config.yml.sample` into `/etc/newrelic-infra/integrations.d`

## Usage

This is the description about how to run the Consul Integration with New Relic Infrastructure agent, so it is required to have the agent installed (see [agent installation](https://docs.newrelic.com/docs/infrastructure/new-relic-infrastructure/installation/install-infrastructure-linux)).

In order to use the Consul Integration it is required to configure `consul-config.yml.sample` file. Firstly, rename the file to `consul-config.yml`. Then, depending on your needs, specify all instances that you want to monitor. Once this is done, restart the Infrastructure agent.

You can view your data in Insights by creating your own custom NRQL queries. To do so use the **ConsulClusterSample** and **ConsulAgentSample** event type.

## Compatibility

* Supported OS: No limitations
* Consul versions: 1.0+

## Integration Development usage

Assuming that you have source code you can build and run the Consul Integration locally.

* Go to directory of the Consul Integration and build it
```bash
$ make
```
* The command above will execute tests for the Consul Integration and build an executable file called `nr-consul` in `bin` directory.
```bash
$ ./bin/nr-consul
```
* If you want to know more about usage of `./nr-consul` check
```bash
$ ./bin/nr-consul -help
```

For managing external dependencies [govendor tool](https://github.com/kardianos/govendor) is used. It is required to lock all external dependencies to specific version (if possible) into vendor directory.
