# Build stage
FROM golang:1.23-alpine AS build

# Set the Current Working Directory inside the container
WORKDIR /app

ARG BUILD_TYPE
ARG APP_VERSION
ARG APP_NAME
ARG HOST_NAME

# Set environment variables
ENV BUILD_TYPE=$BUILD_TYPE
ENV APP_VERSION=$APP_VERSION
ENV APP_NAME=$APP_NAME
ENV HOST_NAME=$HOST_NAME

RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source code into the container
COPY . .

RUN swag init

# Build the Go app
RUN go build -o main main.go

# Final stage
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=build /app/main .
COPY --from=build /app/docs ./docs
COPY version.info ./

ARG BUILD_TYPE
ARG APP_VERSION
ARG APP_NAME
ARG HOST_NAME

# Set environment variables
ENV BUILD_TYPE=$BUILD_TYPE
ENV APP_VERSION=$APP_VERSION
ENV APP_NAME=$APP_NAME
ENV HOST_NAME=$HOST_NAME

# Expose port 3000 to the outside world
EXPOSE 3000

# Command to run the executable
CMD ["./main"]