FROM golang:1.16 AS builder
WORKDIR /src

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-extldflags '-static' -s -w" ./cmd/checker

FROM scratch
WORKDIR /app
COPY --from=builder /src/checker .
EXPOSE 2002
CMD ["./checker"]
