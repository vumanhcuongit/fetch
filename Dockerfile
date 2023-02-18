FROM golang:1.19-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o fetch-go cmd/app/main.go

FROM alpine:3.15

COPY --from=build /app/fetch-go /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/fetch-go"]