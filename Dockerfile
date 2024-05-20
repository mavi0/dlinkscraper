FROM golang:1.22-alpine3.19 AS builder

WORKDIR /root/dlinkscraper
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN go build -trimpath -o out/dlinkscraper ./cmd/dlinkscraper/*.go


FROM alpine:3.19

RUN apk add ca-certificates
COPY --from=builder /root/dlinkscraper/out/dlinkscraper /usr/local/bin/dlinkscraper

ENTRYPOINT ["dlinkscraper"]