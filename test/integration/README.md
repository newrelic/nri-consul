# nri-consul integration test

This integration test is running a 3 node consul cluster with TLS and the nri-consul integration.

## Running the test

```bash
$ make integration-test
```


## certs

In order to create the certs used for this test cluster, we ran the following commands in a Consul node.

To create the consul-agent-ca.pem expiring in 10 years:
```bash
$ consul tls ca create -days=3650
```

To create the server cert and key allowing access through the hostnames and ips we set in the docker-compose:
```bash
$ consul tls cert create -server -dc dc1 -days=3650 \
        -additional-dnsname="consul-server1" \
        -additional-dnsname="consul-server2" \
        -additional-dnsname="consul-server3" \
        -additional-ipaddress 10.5.0.2 \
        -additional-ipaddress 10.5.0.3 \
        -additional-ipaddress 10.5.0.4 
```
