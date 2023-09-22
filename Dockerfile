# Setp 2: install docker
FROM docker as docker

# Step 3: Modules caching
FROM golang:1.21.1-alpine3.17 as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Step 4: Builder
FROM golang:1.21.1-alpine3.17 as builder
RUN apk add --no-cache \
    # Important: required for go-sqlite3
    gcc \
    # Required for Alpine
    musl-dev

COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=1 GOOS=linux go build -o /bin/app -a -ldflags '-linkmode external -extldflags "-static"' .

# RUN GOOS=linux GOARCH=amd64 \
#     go build  -o /bin/app .

# Step 5: Final
FROM scratch

# use these 3 lines for debugging purposes
# FROM ubuntu 
# RUN apt update
# RUN apt install -y git curl iputils-ping
COPY --from=docker /usr/local/bin/docker /usr/local/bin/docker
COPY --from=docker /usr/local/libexec/docker/cli-plugins/docker-buildx /usr/local/libexec/docker/cli-plugins/docker-buildx
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/dockerfiles /dockerfiles
COPY --from=builder /app/build /build
COPY --from=builder /app/config.yml /
COPY --from=builder /bin/app /app
EXPOSE 8080
CMD ["/app"]

