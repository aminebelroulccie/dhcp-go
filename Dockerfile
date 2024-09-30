FROM golang:1.23 as builder
# Define build env
ENV GOOS linux
ENV CGO_ENABLED 0
# Add a work directory
WORKDIR /app
# Cache and install dependencies
COPY go.mod go.sum ./
RUN go mod download
# Copy app files
COPY . .
# Build app
RUN go build  -o dhcp svc/agent/main.go

# FROM alpine:3.14 as production
# # Add certificates
# RUN apk add --no-cache ca-certificates
# # Create Configuration Folders
# RUN mkdir -p /opt/ncnf
# # Copy built binary from builder
# COPY --from=builder app/ncnf-controller /bin/usr/
# # Exec built binary
# CMD  /bin/usr/ncnf-controller