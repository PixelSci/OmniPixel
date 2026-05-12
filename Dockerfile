# OmniPixel Multi-Stage Dockerfile
# Stage 1: Build Frontend
# Stage 2: Build Backend
# Stage 3: Final minimal image

ARG NODE_IMAGE=node:24-alpine
ARG GOLANG_IMAGE=golang:1.25.0-alpine
ARG ALPINE_IMAGE=alpine:3.21
ARG POSTGRES_IMAGE=postgres:18-alpine
ARG GOPROXY=https://goproxy.cn,direct
ARG GOSUMDB=sum.golang.google.cn

# Stage 1
# --------------------------------------
FROM ${NODE_IMAGE} AS frontend-builder
WORKDIR /build/omni-pixel
# Install pnpm
RUN corepack enable && corepack prepare pnpm@latest --activate
# Install dependencies first (better caching)
COPY apps/omni-pixel/package.json apps/omni-pixel/pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile
# Build frontend source
COPY apps/omni-pixel/ ./
RUN pnpm run -F omni-pixel build
# --------------------------------------

# stage 2
# --------------------------------------
FROM ${GOLANG_IMAGE} AS backend-builder
ENV GOPROXY=${GOPROXY}
ENV GOSUMDB=${GOSUMDB}
WORKDIR /build/server
# Copy go mod files first (better caching)
COPY apps/server/go.mod apps/server/go.sum ./
RUN go mod download
# Copy backend source first
COPY apps/server/ ./
RUN go build -o main ./cmd/
# --------------------------------------

# stage 3
# --------------------------------------
FROM ${POSTGRES_IMAGE} AS pg-client
