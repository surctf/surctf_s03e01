FROM golang:alpine AS builder

WORKDIR /build

COPY .. .

ENV CGO_ENABLED 0
ENV GOOS linux
RUN go build -ldflags="-s -w" -o main.bin ./cmd/main.go

FROM alpine

COPY --from=builder /build/main.bin /godra/main.bin

#COPY ./build/main.bin /db-interface/main.bin
RUN ["chmod", "+x", "/godra/main.bin"]

WORKDIR /godra/
ENTRYPOINT ["/godra/main.bin"]
