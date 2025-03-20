FROM alpine

WORKDIR /root

COPY build/${PROJECT_NAME} /root/
COPY config/config.toml /root/config/
COPY migrations /root/migrations

CMD ["/root/go-template-ddd", "-c", "/root/config/config.toml"]
