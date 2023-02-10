FROM golang:1.20-alpine AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -ldflags '-w -s' -o medalhelper .

FROM istio/distroless
COPY --from=builder ["/build/medalhelper", "/"]"
WORKDIR /config
ENTRYPOINT ["/medalhelper"]
