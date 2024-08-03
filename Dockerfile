FROM golang:alpine3.18 AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin .


FROM cgr.dev/chainguard/go:latest AS production
COPY --from=builder /app/bin /bin/app
EXPOSE 9090
ENTRYPOINT ["/bin/app"]