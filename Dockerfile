FROM golang:1.20 as builder

# Set destination for COPY commands
WORKDIR /app

# Copy the Go Mod and Sum files and download the dependencies separately
# This step will be cached if we won't change mod/sum files
COPY ../go.mod .
COPY ../go.sum .
RUN go mod download

COPY . .

# Install dependencies and build the applications
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /out/apiservice ./cmd/api/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /out/datarefresherservice ./cmd/datarefresher/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /out/storageservice ./cmd/storageservice/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /out/websocketservice ./cmd/websocket/main.go

# Run with Alpine image
FROM alpine:latest as api-runner
COPY --from=builder /out/apiservice /app/
COPY --from=builder /app/config.yaml /
COPY --from=builder /app/queries/ /queries/
CMD ["/app/apiservice"]

FROM alpine:latest as refresher-runner
COPY --from=builder /out/datarefresherservice /app/
COPY --from=builder /app/config.yaml /
COPY --from=builder /app/queries/ /queries/
CMD ["/app/datarefresherservice"]

FROM alpine:latest as storage-runner
COPY --from=builder /out/storageservice /app/
COPY --from=builder /app/config.yaml /
CMD ["/app/storageservice"]

FROM alpine:latest as websocket-runner
COPY --from=builder /out/websocketservice /app/
COPY --from=builder /app/config.yaml /
CMD ["/app/websocketservice"]
