# Go Web Sample

Echoフレームワークを使用したGo Webアプリケーションのサンプルプロジェクトです。

## 機能

- REST API
- データベース連携（GORM + SQLite）
- Redisキャッシュ
- JWT認証・認可
- CORS設定
- グレースフルシャットダウン（SIGINT/SIGTERM対応）
- Dockerコンテナ対応
- AWS ECS Fargateデプロイ対応

## 技術スタック

- **フレームワーク**: Echo v4
- **ORM**: GORM
- **データベース**: SQLite（開発環境）、PostgreSQL/MySQL対応可能
- **キャッシュ**: Redis (go-redis/v9)
- **認証**: JWT (golang-jwt/jwt/v5)
- **設定管理**: godotenv
- **コンテナ**: Docker
- **デプロイ**: AWS ECS Fargate

## プロジェクト構造

```
go_web_sample/
├── cmd/
│   └── server/
│       ├── main.go
│       └── t/
│           └── internal/
│               └── handler/
│                   └── hello_test.go
├── internal/
│   ├── config/              # 設定管理
│   ├── handler/             # HTTPハンドラー
│   ├── model/               # データモデル
│   ├── repository/          # データアクセス層
│   ├── service/             # ビジネスロジック層
│   └── middleware/          # カスタムミドルウェア
├── pkg/
│   └── redis/               # Redis接続管理とキャッシュユーティリティ
├── migrations/              # データベースマイグレーション
├── deploy/                   # AWSデプロイ関連
│   ├── task-definition.json
│   ├── deploy.sh
│   └── README.md
├── Dockerfile                # マルチステージDockerfile（local/dev/prd）
├── docker-compose.yml        # 開発環境用
├── .dockerignore             # Docker除外設定
├── .env.example              # 環境変数テンプレート
└── README.md
```

## セットアップ

### 1. 前提条件

- Go 1.23以上
- Git
- Docker（オプション）

### 2. 環境変数の設定

`.env.example`をコピーして`.env`ファイルを作成し、必要に応じて値を変更してください。

```bash
cp .env.example .env
```

### 3. 依存関係のインストール

```bash
go mod download
```

### 4. サーバーの起動

#### ローカル実行（通常）

```bash
go run ./cmd/server
```

#### ローカル実行（ホットリロード付き）

```bash
# Airをインストール（初回のみ、Go 1.23互換バージョン）
go install github.com/cosmtrek/air@v1.49.0

# ホットリロードで起動
air
```

ファイルを変更すると、自動的に再ビルド・再起動されます。

サーバーはデフォルトで `http://localhost:8080` で起動します。

**ホットリロードツールの詳細は [`HOT_RELOAD.md`](HOT_RELOAD.md) を参照してください。**

#### Docker Composeで実行

```bash
# ビルドと起動（localステージ）
docker-compose up -d

# ログ確認
docker-compose logs -f app

# 停止
docker-compose down
```

アプリケーションとRedisがコンテナで起動します。

## API エンドポイント

### ヘルスチェック

```
GET /health
```

### HelloWorld

```
GET /hello
```

### API v1

```
GET /api/v1/
```

## 開発

### ビルド

```bash
go build -o bin/server ./cmd/server
```

### テスト

```bash
go test ./...
```

### グレースフルシャットダウン

サーバーは `signal.NotifyContext` を使用してシグナル（SIGINT、SIGTERM）を受信すると、グレースフルシャットダウンを実行します。

- **SIGINT**（Ctrl+C）: ローカル開発時に手動で停止する場合
- **SIGTERM**: DockerコンテナやECSなどのコンテナ環境で停止シグナルとして送信される場合

シャットダウン時のタイムアウトは `SHUTDOWN_TIMEOUT_SECONDS` 環境変数で設定できます（デフォルト: 10秒）。

## Docker コマンド一覧

### ステージの選択

Dockerfileは以下の3つのステージを提供しています：

- **local**: ローカル開発環境用（デバッグツール、ホットリロード対応）
- **dev**: 開発環境（クラウド）用（軽量、デバッグ可能）
- **prd**: 本番環境用（最小イメージ、最適化）

### イメージのビルド

```bash
# ローカル環境用イメージのビルド
docker build --target local -t go-web-sample:local .

# 開発環境用イメージのビルド
docker build --target dev -t go-web-sample:dev .

# 本番環境用イメージのビルド
docker build --target prd -t go-web-sample:prd .
```

### コンテナの起動（docker run）

```bash
# ローカル環境（バックグラウンド実行）
docker run -d -p 8080:8080 --name go-web-sample-local go-web-sample:local

# ローカル環境（環境変数ファイルを使用）
docker run -d -p 8080:8080 --env-file .env --name go-web-sample-local go-web-sample:local

# 開発環境（バックグラウンド実行）
docker run -d -p 8080:8080 --name go-web-sample-dev go-web-sample:dev

# 本番環境（バックグラウンド実行）
docker run -d -p 8080:8080 --name go-web-sample-prd go-web-sample:prd

# フォアグラウンド実行（ログを直接確認）
docker run -p 8080:8080 go-web-sample:local
```

### コンテナの管理

```bash
# 実行中のコンテナ一覧を表示
docker ps

# すべてのコンテナ（停止中も含む）を表示
docker ps -a

# コンテナの停止
docker stop go-web-sample-local

# コンテナの起動（再起動）
docker start go-web-sample-local

# コンテナの削除
docker rm go-web-sample-local

# コンテナの停止と削除を同時に実行
docker rm -f go-web-sample-local
```

### ログの確認

```bash
# ログを表示
docker logs go-web-sample-local

# ログをリアルタイムで表示（フォロー）
docker logs -f go-web-sample-local

# 最新の100行を表示
docker logs --tail 100 go-web-sample-local

# タイムスタンプ付きで表示
docker logs -t go-web-sample-local
```

### Docker Compose を使用した起動

```bash
# ビルドと起動（バックグラウンド）
docker-compose up -d

# ビルドと起動（フォアグラウンド、ログ表示）
docker-compose up

# ログ確認（リアルタイム）
docker-compose logs -f app

# ログ確認（すべてのサービス）
docker-compose logs -f

# 停止
docker-compose down

# 停止とボリュームの削除
docker-compose down -v

# 再ビルドして起動
docker-compose up -d --build
```

### イメージの管理

```bash
# イメージ一覧を表示
docker images

# イメージの削除
docker rmi go-web-sample:local

# 未使用のイメージを一括削除
docker image prune -a

# イメージの詳細情報を表示
docker inspect go-web-sample:local
```

### APIテスト用curlコマンド

```bash
# ヘルスチェック
curl http://localhost:8080/health

# HelloWorld API
curl http://localhost:8080/hello

# リクエストIDとユーザーIDをヘッダーで指定
curl -H "X-Request-ID: my-request-123" -H "X-User-ID: user-456" http://localhost:8080/hello

# API v1
curl http://localhost:8080/api/v1/

# 詳細なレスポンス情報を表示
curl -v http://localhost:8080/hello

# JSONを整形して表示（jqがインストールされている場合）
curl http://localhost:8080/hello | jq
```

## 環境変数

| 変数名 | 説明 | デフォルト値 |
|--------|------|-------------|
| PORT | サーバーポート | 8080 |
| ENV | 環境（development/production） | development |
| DB_TYPE | データベースタイプ | sqlite |
| DB_PATH | データベースパス | ./data/app.db |
| REDIS_HOST | Redisホスト | localhost |
| REDIS_PORT | Redisポート | 6379 |
| REDIS_PASSWORD | Redisパスワード | （空） |
| REDIS_DB | Redisデータベース番号 | 0 |
| JWT_SECRET | JWT秘密鍵 | your-secret-key-change-in-production |
| JWT_EXPIRATION_HOURS | JWT有効期限（時間） | 24 |
| CORS_ALLOWED_ORIGINS | CORS許可オリジン（カンマ区切り） | <http://localhost:3000,http://localhost:8080> |
| SHUTDOWN_TIMEOUT_SECONDS | サーバーシャットダウンのタイムアウト（秒） | 10 |
| REQUEST_TIMEOUT_SECONDS | リクエストごとのタイムアウト（秒） | 30 |

## AWS ECS Fargate デプロイ

詳細なデプロイ手順は [`deploy/README.md`](deploy/README.md) を参照してください。

### クイックスタート

1. ECRリポジトリとECSクラスターを作成
2. `deploy/task-definition.json`を編集（プレースホルダーを実際の値に置き換え）
3. デプロイスクリプトを実行：

```bash
# 開発環境にデプロイ
ENV_STAGE=dev ./deploy/deploy.sh

# 本番環境にデプロイ
ENV_STAGE=prd ./deploy/deploy.sh
```

### 必要なAWSリソース

- ECRリポジトリ
- ECSクラスター
- VPCとサブネット
- セキュリティグループ
- IAMロール（タスク実行ロール、タスクロール）
- CloudWatch Logsグループ
- Secrets Manager（シークレット管理用）
- ElastiCache（Redis、オプション）
- RDS（データベース、オプション）

## ライセンス

MIT
