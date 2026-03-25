# UI Builder Stage
FROM node:22-alpine AS ui-builder

WORKDIR /app/web
COPY web/package.json web/package-lock.json ./
RUN npm ci
COPY web/ .
RUN npm run build

# Go Builder Stage
FROM golang:1.26-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ cmd/
COPY internal/ internal/
COPY utils/ utils/

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/goma cmd/main.go

# Final Stage
FROM alpine:3.21

ENV TZ=UTC
ENV GOMA_WEB_DIR=/app/web/dist

# Install runtime dependencies
RUN apk --update --no-cache add tzdata ca-certificates curl

# Copy built binary
COPY --from=build /app/goma /usr/local/bin/goma
RUN chmod a+x /usr/local/bin/goma && ln -s /usr/local/bin/goma /goma

# Copy frontend assets
COPY --from=ui-builder /app/web/dist /app/web/dist

ENV POSTA_WEB_DIR=/app/web/dist

WORKDIR /app
EXPOSE 9000

ENTRYPOINT ["/goma"]
