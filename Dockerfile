FROM golang:alpine as builder
RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates
WORKDIR /app
ENV GO111MODULE=on
COPY ["go.mod", "go.sum", "/app/"]
RUN go mod download


FROM builder as binary
RUN GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o bootstrap

FROM alpine
ENV AWS_REGION=us-west-2
COPY --from=binary /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=binary /etc/passwd /etc/passwd
COPY --from=binary /app/bootstrap /bootstrap
ENTRYPOINT ["/bootstrap"]
