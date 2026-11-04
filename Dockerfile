FROM golang:1.27.2-alpine AS build

WORKDIR /src/backend

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/laluna-server ./cmd/server

FROM alpine:3.23

RUN apk add --no-cache ca-certificates \
    && addgroup -S app \
    && adduser -S -G app app

COPY --from=build --chown=app:app /out/laluna-server /usr/local/bin/laluna-server

USER app
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/laluna-server"]
