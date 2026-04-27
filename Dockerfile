FROM postgres:17-alpine AS db

COPY packages/db/init/*.sql /docker-entrypoint-initdb.d/


FROM node:22-alpine AS frontend-builder

WORKDIR /src

RUN corepack enable

COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
COPY apps/omni-pixel/package.json apps/omni-pixel/package.json
RUN pnpm install --frozen-lockfile --filter ./apps/omni-pixel...

COPY apps/omni-pixel apps/omni-pixel
RUN pnpm --filter ./apps/omni-pixel build


FROM nginx:1.27-alpine AS frontend

COPY docker/frontend/nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=frontend-builder /src/apps/omni-pixel/dist /usr/share/nginx/html

EXPOSE 80


FROM golang:1.25-alpine AS backend-builder

WORKDIR /src/apps/server

COPY apps/server/go.mod apps/server/go.sum ./
RUN go mod download

COPY apps/server ./
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/server ./cmd


FROM alpine:3.22 AS backend

RUN addgroup -S app && adduser -S app -G app

WORKDIR /app

COPY --from=backend-builder /out/server /app/server

ENV APP_ENV=production
ENV SERVER_ADDRESS=:8080

EXPOSE 8080

USER app

ENTRYPOINT ["/app/server"]
