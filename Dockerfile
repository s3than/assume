# ---
# Go Builder Image
FROM golang:1.11-alpine AS builder

# Create and set working directory
WORKDIR /app

# copy sources
COPY . .

ENV CGO_ENABLED=0
ENV GO111MODULE=on
ENV GOOS=linux
ENV GOARCH=amd64

RUN apk --no-cache add git

RUN name=$(basename "$dir") \
    set -x && \
    go build -a \
    -installsuffix cgo \
    -ldflags "-w -s" \
    -o "$name" .

# ---
# Application Runtime Image
FROM alpine:3.8

# Copy from builder
COPY --from=builder /app/assume /usr/bin/assume

CMD ["assume", "--help"]
