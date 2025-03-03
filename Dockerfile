FROM golang:1.23.6-alpine AS builder

WORKDIR /usr/local/src

RUN apk add --no-cache bash git make gcc gettext musl-dev

COPY ["./go.mod", "./go.sum", "."]

RUN go mod download

COPY ./ ./

RUN go build -o ./bin/app ./cmd/app/main.go


FROM alpine:3.21 AS runner

RUN apk add --no-cache ca-certificates postgresql-client

COPY --from=builder /usr/local/src/bin/app /app

CMD ["/app"]
