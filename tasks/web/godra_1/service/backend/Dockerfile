FROM golang:alpine AS builder

WORKDIR /build

COPY . .

ENV CGO_ENABLED 0
ENV GOOS linux
RUN go build -ldflags="-s -w" -o main.bin ./cmd/main.go

FROM alpine

COPY --from=builder /build/main.bin /backend/main.bin

RUN ["chmod", "+x", "/backend/main.bin"]

WORKDIR /backend/
ENTRYPOINT ["/backend/main.bin"]
