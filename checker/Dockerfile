FROM golang:1.16 AS builder
WORKDIR /src/checker
COPY ./checker/ /src/checker
COPY ./service/ /src/service
RUN CGO_ENABLED=0 go build -ldflags="-extldflags '-static' -s -w" ./cmd/checker

FROM scratch
WORKDIR /app
COPY --from=builder /src/checker .
EXPOSE 8080
CMD ["./checker"]
