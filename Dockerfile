# Go version argument
ARG GO_VERSION
ARG GO_IMAGE=golang

# Build image
FROM ${GO_IMAGE}:${GO_VERSION} AS builder

WORKDIR /var/tmp/build

# Copy go mod files first to avoid invalidating
# mod cache if they don't change
COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /kubectl-draincheck .

# Flatten into scratch image
FROM scratch

COPY --from=builder /kubectl-draincheck /

ENTRYPOINT ["/kubectl-draincheck"]
CMD []
