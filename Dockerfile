FROM golang:1.24-alpine AS builder

WORKDIR /build-dir

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Вызов необходимых утилит для 

RUN CGO_ENABLE=0 go build -ldflags="-w -s" -o /build-dir/app ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /build-dir/app ./start

# -v ./config:/app/config
# --env-file .env


RUN addgroup -S appgroup && adduser -S appuser -G appgroup
USER appuser

EXPOSE 80

CMD [ "/app/start" ]
