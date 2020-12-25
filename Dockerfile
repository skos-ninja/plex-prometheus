FROM golang:latest as builder

LABEL maintainer="Jake Oliver <docker@skos.ninja>"

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ENV CGO_ENABLED=0
ENV GOOS=linux

RUN go build -a -installsuffix cgo -o main *.go

# Create new image and import just the binary
FROM alpine:latest

# Alpine doesn't include cert auth certificates
RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/main .

ENV CGO_ENABLED=0

EXPOSE 8080

ENTRYPOINT ["./main"]