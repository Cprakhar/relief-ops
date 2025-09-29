FROM golang:1.25.1-alpine AS build

RUN apk update --no-cache && \
    apk add --no-cache gcc musl-dev && \
    rm -rf /var/cache/apk/*

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY services/api-gateway/ ./services/api-gateway/
COPY shared ./shared/

RUN CGO_ENABLED=1 GOOS=linux go build -tags musl -a -installsuffix cgo -o api-gateway ./services/api-gateway

FROM alpine:latest

RUN apk update --no-cache && \
    apk add --no-cache ca-certificates curl && \
    rm -rf /var/cache/apk/*

WORKDIR /root/

COPY --from=build /app/api-gateway .

EXPOSE 8080

CMD ["./api-gateway"]