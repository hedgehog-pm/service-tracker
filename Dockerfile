FROM golang:1.27-alpine AS build
WORKDIR /src

COPY go.mod go.sum ./
COPY cmd ./cmd
COPY internal ./internal

RUN CGO_ENABLED=0 go build -o /out/ ./cmd/service ./cmd/migrate

FROM alpine:3
WORKDIR /app

COPY --from=build /out/ ./
COPY migrations ./migrations

USER 65532:65532
CMD ["/app/service"]
