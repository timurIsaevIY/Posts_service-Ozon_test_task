FROM golang:1.23.1-alpine AS builder
COPY . /github.com/Ozon_Post_comment_system
WORKDIR /github.com/Ozon_Post_comment_system
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -mod=readonly -o ./.bin ./cmd/main.go
FROM scratch AS runner
WORKDIR /build
COPY --from=builder /github.com/Ozon_Post_comment_system/.bin .
COPY --from=builder /github.com/Ozon_Post_comment_system/configs configs/
COPY --from=builder /github.com/Ozon_Post_comment_system/.env .
EXPOSE 8080
ENTRYPOINT ["./.bin"]
