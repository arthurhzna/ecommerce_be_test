FROM golang:1.24-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ecommerce-be main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

ENV TZ=Asia/Jakarta

WORKDIR /root/

COPY --from=builder /app/ecommerce-be .

EXPOSE 8080

CMD ["./ecommerce-be"]
