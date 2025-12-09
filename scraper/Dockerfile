FROM golang:alpine AS builder

WORKDIR /app

COPY go.* ./

RUN go mod download

COPY . .

RUN go build -o main .

FROM alpine AS runner

WORKDIR /root

COPY --from=builder /app/main .

EXPOSE 8080

CMD [ "./main" ]