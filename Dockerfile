### Build the www web assets
###
FROM node:22.5-alpine AS web-builder-www
# Install make 
RUN apk add --no-cache make
WORKDIR /app/www
# Install NodeJS dependencies
COPY ./www/package.json ./www/package-lock.json ./
RUN npm install
COPY ./www .
RUN npm run build

### Build the dungeons assets
###
FROM node:22.5-alpine AS web-builder-dungeons
# Install make 
RUN apk add --no-cache make
WORKDIR /app/dungeons
# Install NodeJS dependencies
COPY ./dungeons/package.json ./dungeons/package-lock.json ./
RUN npm install
COPY ./dungeons .
RUN npm run build

### Build the main bot
###
FROM golang:1.24-alpine AS builder
WORKDIR /app
# Disable CGO
ENV CGO_ENABLED=0
COPY go.mod go.sum ./
RUN --mount=type=ssh go mod download && go mod verify
COPY . .
COPY --from=web-builder-www /app/www/build /app/www/build
COPY --from=web-builder-dungeons /app/dungeons/dist /app/dungeons/dist
RUN go build -ldflags "-s -w" -o /app/rwbyadv3 /app/cmd/bot/main.go
# Install CA certificates for scratch image
RUN apk --no-cache add ca-certificates && update-ca-certificates

### Final image
###
FROM scratch
COPY --from=builder /app/rwbyadv3 .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY cards /cards
USER 1000
ENTRYPOINT ["/rwbyadv3"]