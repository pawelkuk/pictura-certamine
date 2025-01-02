# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.23.4 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /pictura-certamine ./cmd/app/main.go


# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /pictura-certamine /pictura-certamine

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/pictura-certamine"]