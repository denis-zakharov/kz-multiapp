FROM golang:1.16-alpine as builder
ENV CGO_ENABLED=0
WORKDIR /go/src/
ADD go.mod go.sum main.go ./
RUN go mod download
ADD cmd cmd
RUN go build -ldflags '-w -s' -v -o /usr/local/bin/function ./

FROM alpine:latest
COPY --from=builder /usr/local/bin/function /usr/local/bin/function
ENTRYPOINT ["function"]
