FROM mwader/static-ffmpeg AS ffmpeg

FROM golang:1.26.1-alpine3.23 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY docs ./docs
COPY internal ./internal

RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=0 go build -o /app/trophy ./cmd/trophy

RUN chmod +x /app/trophy
RUN mkdir -p /runtime/var/uploads && chown 65532:65532 /runtime/var/uploads && chmod 0775 /runtime/var/uploads

FROM gcr.io/distroless/static-debian13

USER nonroot
WORKDIR /app
COPY --from=mwader/static-ffmpeg:8.1 /ffmpeg /usr/local/bin/
COPY --from=mwader/static-ffmpeg:8.1 /ffprobe /usr/local/bin/
COPY --from=build /app/trophy /app/trophy
COPY --from=build --chown=65532:65532 /runtime/var/uploads /var/uploads

CMD [ "/app/trophy" ]