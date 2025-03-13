FROM golang:1.24-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main main.go

FROM alpine AS runner

WORKDIR /app

RUN apk add --no-cache curl

COPY --from=build /app/main .

CMD ["/app/main"]