# Go Web Sample

Echoフレームワークを使用したGo Webアプリケーションのサンプルプロジェクトです。

## 機能

- REST API
- データベース連携（GORM + SQLite）
- JWT認証・認可
- CORS設定

## 技術スタック

- **フレームワーク**: Echo v4
- **ORM**: GORM
- **データベース**: SQLite（開発環境）、PostgreSQL/MySQL対応可能
- **認証**: JWT (golang-jwt/jwt/v5)
- **設定管理**: godotenv

## プロジェクト構造

```
go_web_sample/
├── cmd/
│   └── server/
│       └── main.go          # エントリーポイント
├── internal/
│   ├── config/              # 設定管理
│   ├── handler/             # HTTPハンドラー
│   ├── model/               # データモデル
│   ├── repository/          # データアクセス層
│   ├── service/             # ビジネスロジック層
│   └── middleware/          # カスタムミドルウェア
├── pkg/                     # 共通パッケージ
├── migrations/              # データベースマイグレーション
├── .env.example             # 環境変数テンプレート
└── README.md
```

## セットアップ

### 1. 前提条件

- Go 1.23以上
- Git

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

```bash
go run cmd/server/main.go
```

サーバーはデフォルトで `http://localhost:8080` で起動します。

## API エンドポイント

### ヘルスチェック

```
GET /health
```

### API v1

```
GET /api/v1/
```

## 開発

### ビルド

```bash
go build -o bin/server cmd/server/main.go
```

### テスト

```bash
go test ./...
```

## 環境変数

| 変数名 | 説明 | デフォルト値 |
|--------|------|-------------|
| PORT | サーバーポート | 8080 |
| ENV | 環境（development/production） | development |
| DB_TYPE | データベースタイプ | sqlite |
| DB_PATH | データベースパス | ./data/app.db |
| JWT_SECRET | JWT秘密鍵 | your-secret-key-change-in-production |
| JWT_EXPIRATION_HOURS | JWT有効期限（時間） | 24 |
| CORS_ALLOWED_ORIGINS | CORS許可オリジン（カンマ区切り） | <http://localhost:3000,http://localhost:8080> |

## ライセンス

MIT
