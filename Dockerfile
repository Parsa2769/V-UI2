# Multi-stage Dockerfile for 3X-UI Modern

# Stage 1: Backend Builder
FROM golang:1.21-alpine AS backend-builder

WORKDIR /build

# Install build dependencies
RUN apk add --no-cache git make gcc musl-dev

# Copy go mod files
COPY backend/go.mod backend/go.sum ./
RUN go mod download && go mod tidy

# Copy source code (including the new go.sum)
COPY backend/ ./
# Ensure go.sum is up to date with all dependencies
RUN go mod tidy && go mod download

# Build backend
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /app/server \
    ./cmd/server

# Stage 2: Frontend Builder
FROM node:18-alpine AS frontend-builder

WORKDIR /build

# Copy package files
COPY frontend/package*.json ./
# Install dependencies (works with or without package-lock.json)
RUN npm install

# Copy source code
COPY frontend/ ./

# Build frontend
RUN npm run build

# Stage 3: Backend Runtime
FROM alpine:latest AS backend

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    wget \
    unzip \
    && rm -rf /var/cache/apk/*

# Install Xray
RUN wget -O /tmp/Xray-linux-64.zip https://github.com/XTLS/Xray-core/releases/latest/download/Xray-linux-64.zip \
    && unzip /tmp/Xray-linux-64.zip -d /usr/local/bin/ \
    && chmod +x /usr/local/bin/xray \
    && rm /tmp/Xray-linux-64.zip

# Copy backend binary
COPY --from=backend-builder /app/server /app/server

# Copy migrations
COPY backend/migrations /app/migrations

# Create necessary directories
RUN mkdir -p /app/data /app/config /etc/xray

# Set up user
RUN addgroup -g 1000 app && \
    adduser -D -u 1000 -G app app && \
    chown -R app:app /app /etc/xray

USER app

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 \
    CMD wget --quiet --tries=1 --spider http://localhost:8080/health || exit 1

CMD ["/app/server"]

# Stage 4: Frontend Runtime (nginx serves static files)
FROM nginx:alpine AS frontend

COPY --from=frontend-builder /build/dist /usr/share/nginx/html
COPY nginx/default.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
