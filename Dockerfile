FROM golang:1.25-alpine AS builder

COPY . /myapp
WORKDIR /myapp
ENV CGO_ENABLED==0
RUN go build -o app ./cmd/rest/main.go

FROM alpine:latest AS deployment
COPY --from=builder /myapp/app .
CMD ["./app"]