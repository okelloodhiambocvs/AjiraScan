FROM golang:1.25.4-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/ajirascan ./cmd/web

FROM gcr.io/distroless/static-debian12:nonroot
WORKDIR /app
COPY --from=build /out/ajirascan /app/ajirascan
COPY templates /app/templates
COPY static /app/static
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/app/ajirascan"]
