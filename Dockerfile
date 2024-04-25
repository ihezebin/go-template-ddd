FROM alpine

WORKDIR /root

COPY build/${PROJECT_NAME} /root/
COPY config/config.toml /root/

CMD ["/root/go-template-ddd", "-c", "/root/config.toml"]
