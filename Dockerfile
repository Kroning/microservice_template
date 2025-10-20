FROM golang:1.23 AS builder

ARG OUTPUT_BINARY
ARG APP_VERSION
ARG BUILD_DIR
ARG GITHUB_CREDS

ADD . /app/
WORKDIR /app

ENV CGO_ENABLED=0
RUN go mod tidy -v

RUN go build -o ./main ./cmd/main.go

FROM alpine:latest

ARG OUTPUT_BINARY
USER 0

COPY --from=builder /app/config config
COPY --from=builder /app/main /app/bin

RUN chown -R -f 888:888 /app \
 && chmod +x /app
USER 888

ENTRYPOINT ["./app/bin"]
