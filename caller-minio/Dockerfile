FROM golang:1.24-bookworm
WORKDIR /app
COPY . .
RUN go build -o caller-minio
CMD ["./caller-minio"]
