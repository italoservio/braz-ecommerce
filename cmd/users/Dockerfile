FROM golang:1.21.6-bullseye
EXPOSE 3000
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o users cmd/users/main.go
CMD ["/app/users"]
