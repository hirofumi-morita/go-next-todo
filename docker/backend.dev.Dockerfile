FROM golang:1.21-alpine

RUN apk add --no-cache git

RUN go install github.com/air-verse/air@latest

WORKDIR /app

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ .

EXPOSE 8080

CMD ["air", "-c", ".air.toml"]
