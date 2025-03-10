FROM golang:1.22.2 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Defina as variáveis de build usando ARG
# ARG JWT_SECRET
# ARG ENTERPRISES_DB_URL
# ARG ULLER_DB_URL
# ARG ULLER_DB_NAME

# # Defina as variáveis de ambiente dentro da imagem
# ENV JWT_SECRET=$JWT_SECRET
# ENV ENTERPRISES_DB_URL=$ENTERPRISES_DB_URL
# ENV ULLER_DB_URL=$ULLER_DB_URL
# ENV ULLER_DB_NAME=$ULLER_DB_NAME

RUN CGO_ENABLED=0 GOOS=linux go build -o uller-backend

FROM scratch
COPY --from=builder /app/uller-backend /
CMD ["/uller-backend"]