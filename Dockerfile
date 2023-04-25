FROM scratch
WORKDIR /app
COPY dist/go-template-ddd /app/
ENTRYPOINT ["/app/go-template-ddd"]
