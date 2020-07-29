FROM golang:1.12.0-alpine3.9
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN apk add git && \
    go get google.golang.org/grpc && \
    go get github.com/lib/pq && \
    go build -o server .
CMD ["/app/server"]