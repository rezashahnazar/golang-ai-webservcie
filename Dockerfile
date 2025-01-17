FROM golang:1.23.5-alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o main .

FROM golang:1.23.5-alpine

WORKDIR /app
COPY --from=builder /app/main .

EXPOSE 8080
CMD ["./main"]