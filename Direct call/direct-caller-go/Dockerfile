FROM golang:1.24-bookworm
WORKDIR /app
COPY . .
RUN go build -o direct-caller-go
CMD ["./direct-caller-go"]