FROM golang:1.22-alpine AS builder
WORKDIR /app
ENV CGO_ENABLED=0
COPY go.mod go.sum ./
RUN --mount=type=ssh go mod download && go mod verify
COPY . .
RUN go build -o /app/rwbyadv3 /app/cmd/bot/main.go

RUN apk --no-cache add ca-certificates && update-ca-certificates

FROM scratch
COPY --from=builder /app/rwbyadv3 .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY sql /sql
USER 1000
ENTRYPOINT ["/rwbyadv3"]