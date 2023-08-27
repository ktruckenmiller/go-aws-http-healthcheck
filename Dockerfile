FROM golang:alpine as builder
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates
WORKDIR /app
ENV GO111MODULE=on
COPY ["go.mod", "go.sum", "/app/"]
RUN go mod download
