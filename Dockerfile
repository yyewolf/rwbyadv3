### Build the main web assets
###
FROM node:22.5-alpine AS web-builder-main
# Install make 
RUN apk add --no-cache make
WORKDIR /app
# Install NodeJS dependencies
COPY package.json package-lock.json ./
RUN npm install
COPY . .
RUN make assets

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
FROM golang:1.22-alpine AS builder
WORKDIR /app
# Disable CGO
ENV CGO_ENABLED=0
COPY go.mod go.sum ./
RUN --mount=type=ssh go mod download && go mod verify
COPY . .
RUN go install github.com/a-h/templ/cmd/templ@latest
RUN go generate templ.go
COPY --from=web-builder-main /app/static /app/static
COPY --from=web-builder-dungeons /app/dungeons/dist /app/dungeons/dist
RUN go build -o /app/rwbyadv3 /app/cmd/bot/main.go
# Install CA certificates for scratch image
RUN apk --no-cache add ca-certificates && update-ca-certificates

### Final image
###
FROM scratch
COPY --from=builder /app/rwbyadv3 .
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY cards /cards
COPY sql /sql
USER 1000
ENTRYPOINT ["/rwbyadv3"]