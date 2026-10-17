# syntax=docker/dockerfile:1

# Build stages run on the build host's platform and cross-compile, so
# multi-arch builds don't need emulation.
FROM --platform=$BUILDPLATFORM node:24-alpine AS frontend
WORKDIR /src
COPY package.json package-lock.json ./
RUN npm ci
COPY index.html vite.config.js ./
COPY src ./src
RUN npm run build

FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS backend
ARG TARGETOS
ARG TARGETARCH
ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_DATE=unknown
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend /src/internal/server/dist ./internal/server/dist
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build \
    -trimpath \
    -ldflags "-s -w -X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${BUILD_DATE}" \
    -o /out/oneoff .

FROM alpine:3.22
# bash, curl and jq are here for shell jobs, which run inside this container.
# tzdata provides zone info for job timezones.
RUN apk add --no-cache bash ca-certificates curl jq tzdata \
    && addgroup -S -g 1000 oneoff \
    && adduser -S -D -u 1000 -G oneoff -h /data oneoff \
    && mkdir -p /data \
    && chown oneoff:oneoff /data
COPY --from=backend /out/oneoff /usr/local/bin/oneoff

ENV HOST=0.0.0.0 \
    PORT=8080 \
    DB_PATH=/data/oneoff.db

WORKDIR /data
VOLUME /data
USER oneoff
EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget -q -O /dev/null "http://127.0.0.1:${PORT}/health" || exit 1

ENTRYPOINT ["oneoff"]
CMD ["serve"]
