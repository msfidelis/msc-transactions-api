FROM golang:1.23 AS builder

WORKDIR /app

COPY ./ .

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .


FROM scratch

COPY --from=builder /app/main ./
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080

CMD ["./main"]