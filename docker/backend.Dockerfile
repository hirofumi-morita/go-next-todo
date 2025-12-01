FROM golang:1.25-alpine AS builder

# Install gqlgen for GraphQL code generation
RUN go install github.com/99designs/gqlgen@latest

WORKDIR /app

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ .

# Generate GraphQL code and build
RUN go generate ./... && CGO_ENABLED=0 GOOS=linux go build -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]
