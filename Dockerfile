# Stage 1
FROM golang:1.23-alpine AS builder
LABEL maintainer="savioruz <jakueenak@gmail.com>"

# Install dependencies
RUN apk update && apk --no-cache add ca-certificates gcc musl-dev

# Move to working directory (/build).
WORKDIR /build

# Copy and download dependency using go mod.
COPY go.mod go.sum ./
RUN go mod download

# Copy the code into the container.
COPY . .

# Set necessary environment variables needed for our image and build the API server.
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -ldflags="-s -w" -o todos .

# Stage 2
FROM scratch

# Copy CA certificates from the builder stage to enable SSL verification
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy binary and config files from /build to root folder of scratch container.
COPY --from=builder ["/build/todos", "/"]

# Command to run when starting the container.
ENTRYPOINT ["/todos"]

# Expose port 3000 to the outside world.
EXPOSE 3000