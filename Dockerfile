# ビルドステージ（共通）
FROM golang:1.23-alpine AS builder

WORKDIR /app

# CA証明書をインストール（prdステージ用）
RUN apk --no-cache add ca-certificates

# 依存関係をコピー
COPY go.mod go.sum ./
RUN go mod download

# ソースコードをコピー
COPY . .

# バイナリをビルド（基本）- vendorディレクトリを無視
RUN CGO_ENABLED=0 GOOS=linux go build -mod=mod -a -installsuffix cgo -o server ./cmd/server/main.go

# バイナリをビルド（最適化版 - dev/prd用）- vendorディレクトリを無視
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=mod \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o server-optimized ./cmd/server/main.go

# ============================================
# local ステージ: ローカル開発環境用
# ============================================
FROM golang:1.23-alpine AS local

WORKDIR /app

# 開発ツールをインストール
RUN apk --no-cache add \
    ca-certificates \
    tzdata \
    git \
    curl \
    vim

# Airをインストール（ホットリロード用 - Go 1.23互換バージョン）
RUN go install github.com/cosmtrek/air@v1.49.0

# ソースコードをコピー（ホットリロード用）
COPY . .

# 依存関係をダウンロード
COPY go.mod go.sum ./
RUN go mod download

# ビルド済みバイナリをコピー（オプション）
COPY --from=builder /app/server /app/server

# ポートを公開
EXPOSE 8080

# Airでホットリロード実行（.air.tomlが存在しない場合はgo runにフォールバック）
CMD ["sh", "-c", "if [ -f .air.toml ]; then air -c .air.toml; else go run ./cmd/server/main.go; fi"]

# ============================================
# dev ステージ: 開発環境（クラウド）用
# ============================================
FROM alpine:latest AS dev

RUN apk --no-cache add \
    ca-certificates \
    tzdata \
    curl

WORKDIR /root/

# ビルドしたバイナリをコピー（最適化版）
COPY --from=builder /app/server-optimized ./server

# ポートを公開
EXPOSE 8080

# ヘルスチェック用のcurlを利用可能
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

# アプリケーションを実行
CMD ["./server"]

# ============================================
# prd ステージ: 本番環境用（最小イメージ）
# ============================================
FROM scratch AS prd

# CA証明書をコピー（HTTPS接続用）
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# タイムゾーンデータをコピー
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo

WORKDIR /root/

# ビルドしたバイナリをコピー（最適化版）
COPY --from=builder /app/server-optimized ./server

# ポートを公開
EXPOSE 8080

# アプリケーションを実行
CMD ["./server"]
