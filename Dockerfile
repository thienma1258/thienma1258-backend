ARG DOCKER_IMAGE_TAG="alpine"


# ===== Builder =====
FROM golang:${DOCKER_IMAGE_TAG} as builder

WORKDIR /app

RUN \
  apk add \
  --no-cache \
  --update \
  --virtual "shared-dependencies" \
  "bash" \
  "bc" \
  "curl" \
  "upx" \
  "git" \
  "jq"

# Core
COPY "./lib" "./lib"

# Install
COPY "./src" "./src"
COPY "./project.config.json" "./project.config.json"
COPY "./pipeline/install" "./pipeline/install"
COPY "./pipeline/build" "./pipeline/build"

RUN  dos2unix ./pipeline/install ./pipeline/build ./lib/bash/core.sh ./lib/bash/go.sh
RUN ./pipeline/install

# Build
RUN ./pipeline/build


# ===== Production =====
FROM alpine:latest

WORKDIR /app

# Built
COPY --from=builder "/app/build/appBuilt" "./appBuilt"

ENTRYPOINT ["./appBuilt"]
