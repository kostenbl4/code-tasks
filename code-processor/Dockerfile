FROM golang:1.24-alpine AS build

WORKDIR /app

COPY go.mod go.sum* ./

COPY . .

RUN go mod download

WORKDIR /app/code-processor
RUN go build -o main main.go

FROM alpine AS runner

WORKDIR /app

RUN apk add --no-cache curl

COPY --from=build /app/code-processor/main .

COPY --from=build /app/code-processor/config.yaml ./

CMD ["/app/main"]