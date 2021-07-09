FROM golang:1.16 AS builder
WORKDIR /src

COPY go.* .
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-extldflags '-static' -s -w" ./cmd/gosship

FROM scratch
WORKDIR /app
COPY --from=builder /src/gosship .
EXPOSE 2222
CMD ["./gosship"]
