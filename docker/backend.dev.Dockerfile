FROM golang:1.25-alpine

RUN apk add --no-cache git

# Install Air for hot-reload
RUN go install github.com/air-verse/air@latest

# Install gqlgen for GraphQL code generation
RUN go install github.com/99designs/gqlgen@latest

WORKDIR /app

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ .

# Generate GraphQL code on initial build
RUN go generate ./...

EXPOSE 8080

CMD ["air", "-c", ".air.toml"]
