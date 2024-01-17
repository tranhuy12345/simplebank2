#Build stage
FROM golang:1.21.6-alpine3.18 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

#Build stage copy từ /app/main
FROM golang:alpine3.18
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .


EXPOSE 8080
CMD [ "/app/main" ]
