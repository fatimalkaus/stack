FROM alpine:latest

COPY bin/stack-server /stack-server
COPY config.yaml /etc/stack/config.yaml


CMD ["/stack-server", "--config", "/etc/stack/config.yaml"]