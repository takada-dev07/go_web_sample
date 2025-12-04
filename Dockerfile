# ビルドステージ
FROM golang:1.23-alpine AS builder

WORKDIR /app

# 依存関係をコピー
COPY go.mod go.sum ./
RUN go mod download

# ソースコードをコピー
COPY . .

# バイナリをビルド
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server/main.go

# 実行ステージ
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# ビルドしたバイナリをコピー
COPY --from=builder /app/server .

# ポートを公開
EXPOSE 8080

# アプリケーションを実行
CMD ["./server"]

