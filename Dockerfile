FROM golang:1 AS builder
COPY . /app
ARG CGO_ENABLED=0
ARG GOPROXY=https://goproxy.cn,direct
RUN cd /app && \
    go build -o /medalhelper /app

FROM gcr.io/distroless/static
COPY --from=builder /medalhelper /medalhelper
WORKDIR /config
ENTRYPOINT ["/medalhelper"]
