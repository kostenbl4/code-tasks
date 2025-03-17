FROM golang:1.24-alpine AS build

WORKDIR /app

COPY go.mod go.sum* ./

COPY . .

RUN go mod download

WORKDIR /app/task-service
RUN go build -o main main.go

FROM alpine AS runner

WORKDIR /app

RUN apk add --no-cache curl

COPY --from=build /app/task-service/main .

COPY --from=build /app/task-service/config.yaml ./

CMD ["/app/main"]