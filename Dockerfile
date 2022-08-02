# Build image
FROM golang:1.18 AS builder

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
