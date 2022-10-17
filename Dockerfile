FROM golang:1.17-alpine AS builder

LABEL maintainer="acornsoft"

# Move to working directory (/build).
WORKDIR /build

# Copy and download dependency using go mod.
COPY go.mod go.sum ./
RUN go mod download

# Copy the code into the container.
COPY . .

# Set necessary environment variables needed for our image and build the API server.
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -ldflags="-s -w" -o edgecraft-api .

FROM scratch

# Copy binary and config files from /build to root folder of scratch container.
COPY --from=builder ["/build/edgecraft-api", "/"]
COPY --from=builder ["/build/conf", "/conf"]
COPY --from=builder ["/build/templates", "/templates"]

# Command to run when starting the container.
ENTRYPOINT ["/edgecraft-api"]