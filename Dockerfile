FROM golang:1.22 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./cmd/accounts_storage

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/app .
COPY --from=builder /app/cmd/accounts_storage/configs/apiserver.toml ./cmd/accounts_storage/configs/apiserver.toml 

EXPOSE 8080

CMD ["./app"]
