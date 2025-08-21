FROM golang:1.24.2
WORKDIR /app
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /api
EXPOSE 8080
CMD ["/api", "--mode=production"]


