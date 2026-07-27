# UI Builder Stage
FROM node:22-alpine AS ui-builder

WORKDIR /app/web
COPY web/package.json web/package-lock.json ./
RUN npm ci
COPY web/ .
RUN npm run build

# Go Builder Stage
FROM golang:1.26-alpine AS build

ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_DATE=2026-03-25


WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ cmd/
COPY internal/ internal/
COPY utils/ utils/

# Stage the built dashboard where //go:embed reads it (see internal/web), so the
# binary serves the UI itself.
COPY --from=ui-builder /app/web/dist internal/web/dist

# Build the binary
RUN CGO_ENABLED=0 go build \
    -ldflags "-s -w -X 'github.com/jkaninda/goma-admin/internal/config.Version=${VERSION}' -X 'github.com/jkaninda/goma-admin/internal/config.CommitID=${COMMIT}' -X 'github.com/jkaninda/goma-admin/internal/config.BuildDate=${BUILD_DATE}'" \
    -o /bin/goma ./cmd
# Final Stage
FROM alpine:3.21

ENV TZ=UTC

# Install runtime dependencies
RUN apk --update --no-cache add tzdata ca-certificates curl

# Copy built binary — the dashboard is embedded in it, so no static assets ship
# alongside and GOMA_WEB_DIR stays unset.
COPY --from=build /bin/goma /usr/local/bin/goma
RUN chmod a+x /usr/local/bin/goma && ln -s /usr/local/bin/goma /goma

WORKDIR /app
EXPOSE 9000

ENTRYPOINT ["/goma"]
