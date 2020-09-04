FROM golang:1.15 AS builder
WORKDIR /go/src/github.com/form3tech-oss/cilium-etcd-watchdog
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /cilium-etcd-watchdog -v ./main.go

FROM gcr.io/distroless/base
USER nobody:nobody
WORKDIR /
COPY --from=builder /cilium-etcd-watchdog /cilium-etcd-watchdog
ENTRYPOINT ["/cilium-etcd-watchdog"]
