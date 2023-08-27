FROM golang:1.17 AS builder

WORKDIR /app
COPY . .

RUN go build -o /app/bin/organize

FROM debian:buster-slim
COPY --from=builder /app/bin/organize /organize

EXPOSE 23234

CMD ["/organize"]