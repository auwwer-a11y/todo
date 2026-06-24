FROM golang:1.22-alpine AS builder
WORKDIR /app
ENV GOTOOLCHAIN=local
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o todo cmd/app/main.go

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/todo .
EXPOSE 8080
CMD ["./todo"]