
# Build the manager binary
FROM golang:1.14 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY main.go main.go
COPY mutate.go mutate.go

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o pod-mutating-webhook main.go mutate.go


FROM alpine:3.10
COPY --from=builder /workspace/pod-mutating-webhook /
ENTRYPOINT ["/pod-mutating-webhook"]