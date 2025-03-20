FROM alpine

WORKDIR /root

COPY build/${PROJECT_NAME} /root/
COPY config/config.toml /root/config/
COPY migration /root/migration

CMD ["/root/go-template-ddd", "-c", "/root/config/config.toml"]
