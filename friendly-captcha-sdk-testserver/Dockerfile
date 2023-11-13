FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download
RUN go mod verify
RUN CGO_ENABLED=0 go build -o ./friendly-captcha-sdk-testserver
RUN chmod +x ./friendly-captcha-sdk-testserver

FROM scratch
WORKDIR /opt
COPY --from=builder /app/friendly-captcha-sdk-testserver .

CMD ["/opt/friendly-captcha-sdk-testserver", "serve"]
