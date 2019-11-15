FROM golang:1.10 as builder
COPY . /go/src/github.com/newrelic/nri-consul/
RUN cd /go/src/github.com/newrelic/nri-consul && \
    make && \
    strip ./bin/nri-consul

FROM newrelic/infrastructure:latest
ENV NRIA_IS_FORWARD_ONLY true
ENV NRIA_K8S_INTEGRATION true
COPY --from=builder /go/src/github.com/newrelic/nri-consul/bin/nri-consul /nri-sidecar/newrelic-infra/newrelic-integrations/bin/nri-consul
COPY --from=builder /go/src/github.com/newrelic/nri-consul/consul-definition.yml /nri-sidecar/newrelic-infra/newrelic-integrations/definition.yml
USER 1000
