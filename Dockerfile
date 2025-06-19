# --- Build Stage ---
FROM golang:1.23.3-alpine AS builder

WORKDIR /app

# ติดตั้ง git (สำหรับ go mod) และ timezone
RUN apk add --no-cache git tzdata

# คัดลอก go.mod และ go.sum ก่อน เพื่อ cache dependency
COPY go.mod go.sum ./
RUN go mod download

# คัดลอก source code ทั้งหมด
COPY . .

# Build binary (ชื่อ main)
RUN go build -o main ./cmd/main.go

# --- Run Stage ---
FROM debian:bookworm-slim

WORKDIR /app

RUN apt-get update && apt-get install -y tzdata

COPY --from=builder /app/main .
COPY .env .env

ENV PORT=8080

EXPOSE 8080

CMD ["./main"] 