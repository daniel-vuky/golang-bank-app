# Build stage
FROM golang:1.22.2-alpine3.19 AS BUILDER
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run stage
FROM alpine:3.19
WORKDIR /app
COPY --from=BUILDER /app/main .
COPY app.env .

EXPOSE 8080
CMD ["/app/main"]