FROM golang:1.21 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o uller-backend

FROM scratch
COPY --from=builder /app/uller-backend /
CMD ["/uller-backend"]