FROM golang:1.21.1

WORKDIR /go/src/hash


COPY . .
RUN go mod init main
RUN go mod tidy
RUN go build -o hash main.go

CMD t1=$(date +%s%3N); ./hash; t2=$(date +%s%3N); echo "\n$((t2-t1))"
