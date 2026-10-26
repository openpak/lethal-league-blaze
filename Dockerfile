# syntax=docker/dockerfile:1
# Lethal League Blaze front (lethal-league-blazed): /healthz + /readyz over the
# title's NEX listeners. Plain HTTP on 23172 (requested port base 23170 + 2).
FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS build
ARG TARGETOS TARGETARCH
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH GOWORK=off go build -trimpath -buildvcs=false -ldflags="-s -w" -o /out/lethal-league-blazed ./cmd/lethal-league-blazed

FROM gcr.io/distroless/static
COPY --from=build /out/lethal-league-blazed /lethal-league-blazed
USER nonroot:nonroot
EXPOSE 23172
ENTRYPOINT ["/lethal-league-blazed"]
