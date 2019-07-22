FROM golang:1.10 as builder
RUN go get -d github.com/newrelic/nri-consul/... && \
    cd /go/src/github.com/newrelic/nri-consul && \
    make && \
    strip ./bin/nr-consul

FROM newrelic/infrastructure:latest
ENV NRIA_IS_FORWARD_ONLY true
ENV NRIA_K8S_INTEGRATION true
COPY --from=builder /go/src/github.com/newrelic/nri-consul/bin/nr-consul /var/db/newrelic-infra/newrelic-integrations/bin/nr-consul
COPY --from=builder /go/src/github.com/newrelic/nri-consul/consul-definition.yml /var/db/newrelic-infra/newrelic-integrations/definition.yml
