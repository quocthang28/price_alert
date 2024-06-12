FROM golang:1.22-alpine as builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o main .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/app.env .
COPY --from=builder /app/config.json .
COPY entrypoint.sh . 

RUN chmod +x entrypoint.sh

CMD ["./entrypoint.sh"]

#docker run --rm -d -v rmb-volume:/app/config rmb-vol