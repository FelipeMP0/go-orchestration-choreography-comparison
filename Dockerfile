FROM golang:1.14.11-alpine3.11 AS builder

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o main .

WORKDIR /dist

RUN cp /build/main .

FROM golang:1.14.11-alpine3.11

COPY --from=builder /dist/main /

ENTRYPOINT ["/main"]
