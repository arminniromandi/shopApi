FROM golang:1.24-alpine AS builder

WORKDIR /app

# نصب ابزارهای مورد نیاز برای دانلود پکیج‌ها
RUN apk add --no-cache git

# کپی فایل‌های وابستگی و دانلود آن‌ها
COPY go.mod go.sum ./
RUN go mod download

# کپی کدهای منبع
COPY . .

# کامپایل بهینه‌شده به صورت استاتیک
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main .

# ایمیج نهایی بسیار سبک (حدود ۲۰ مگابایت)
FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add ca-certificates tzdata

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
