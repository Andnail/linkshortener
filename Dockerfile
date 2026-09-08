FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /linkshortener ./cmd/linkshortener

FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata

RUN adduser -D -g '' appuser

WORKDIR /home/appuser

COPY --from=builder /linkshortener .

RUN chown -R appuser:appuser .

USER appuser

EXPOSE 8080

CMD ["./linkshortener"]