# AWS ECS Fargate デプロイ手順

このディレクトリには、AWS ECS Fargateにアプリケーションをデプロイするためのファイルが含まれています。

## 環境ステージについて

Dockerfileは以下のステージを提供しています：

- **dev**: 開発環境用（Alpine Linuxベース、デバッグツール付き、ヘルスチェック対応）
- **prd**: 本番環境用（scratchベース、最小イメージ、最適化済み）

デプロイスクリプトは `ENV_STAGE` 環境変数でステージを選択します。

## 前提条件

- AWS CLIがインストール・設定されていること
- Dockerがインストールされていること
- 以下のAWSリソースが作成されていること：
  - ECRリポジトリ
  - ECSクラスター
  - VPCとサブネット
  - セキュリティグループ
  - Application Load Balancer（オプション）
  - ElastiCache（Redis）クラスター（オプション）
  - RDSデータベース（オプション）
  - Secrets Manager（シークレット管理用）

## セットアップ手順

### 1. ECRリポジトリの作成

```bash
aws ecr create-repository \
  --repository-name go-web-sample \
  --region ap-northeast-1
```

### 2. ECSクラスターの作成

```bash
aws ecs create-cluster \
  --cluster-name go-web-sample-cluster \
  --region ap-northeast-1
```

### 3. IAMロールの作成

#### タスク実行ロール（ecsTaskExecutionRole）

ECSタスクがECRからイメージをプルし、CloudWatch Logsにログを送信するためのロール。

必要なポリシー：

- `AmazonEC2ContainerRegistryReadOnly`
- `CloudWatchLogsFullAccess`

#### タスクロール（ecsTaskRole）

アプリケーションがAWSサービス（Secrets Manager、ElastiCache等）にアクセスするためのロール。

必要なポリシー：

- `SecretsManagerReadWrite`（シークレット読み取り用）

### 4. Secrets Managerにシークレットを登録

```bash
# JWT秘密鍵
aws secretsmanager create-secret \
  --name go-web-sample/jwt-secret \
  --secret-string "your-production-jwt-secret" \
  --region ap-northeast-1

# データベースパスワード
aws secretsmanager create-secret \
  --name go-web-sample/db-password \
  --secret-string "your-db-password" \
  --region ap-northeast-1

# Redisパスワード
aws secretsmanager create-secret \
  --name go-web-sample/redis-password \
  --secret-string "your-redis-password" \
  --region ap-northeast-1
```

### 5. CloudWatch Logsグループの作成

```bash
aws logs create-log-group \
  --log-group-name /ecs/go-web-sample \
  --region ap-northeast-1
```

### 6. タスク定義ファイルの編集

`task-definition.json`を編集し、以下のプレースホルダーを実際の値に置き換えてください：

- `YOUR_ACCOUNT_ID`: AWSアカウントID
- `YOUR_ECR_REPOSITORY_URI`: ECRリポジトリURI
- `REGION`: AWSリージョン
- `YOUR_ELASTICACHE_ENDPOINT`: ElastiCacheエンドポイント（使用する場合）
- Secrets ManagerのARNも実際の値に置き換える

### 7. デプロイスクリプトの実行

```bash
# 開発環境にデプロイ
ENV_STAGE=dev IMAGE_TAG=v1.0.0-dev ./deploy/deploy.sh

# 本番環境にデプロイ
ENV_STAGE=prd IMAGE_TAG=v1.0.0 ./deploy/deploy.sh
```

**注意**: `ENV_STAGE` を指定しない場合、デフォルトで `prd` が使用されます。

## ECSサービスの作成（初回のみ）

初回デプロイ時は、ECSサービスを手動で作成する必要があります：

```bash
aws ecs create-service \
  --cluster go-web-sample-cluster \
  --service-name go-web-sample-service \
  --task-definition go-web-sample \
  --desired-count 1 \
  --launch-type FARGATE \
  --network-configuration "awsvpcConfiguration={subnets=[subnet-xxx,subnet-yyy],securityGroups=[sg-xxx],assignPublicIp=ENABLED}" \
  --region ap-northeast-1
```

## 注意事項

- 本番環境では、必ずSecrets Managerを使用してシークレットを管理してください
- ElastiCacheやRDSを使用する場合は、VPC内の適切なセキュリティグループ設定が必要です
- タスク定義のCPU/メモリ設定は、アプリケーションの要件に応じて調整してください
- ヘルスチェックの設定を確認し、適切な間隔とタイムアウトを設定してください

## トラブルシューティング

### タスクが起動しない

- CloudWatch Logsでエラーログを確認
- セキュリティグループの設定を確認
- タスク実行ロールの権限を確認

### イメージのプルに失敗する

- ECRへのアクセス権限を確認
- タスク実行ロールに`AmazonEC2ContainerRegistryReadOnly`ポリシーが付与されているか確認

### アプリケーションが外部サービスに接続できない

- セキュリティグループのアウトバウンドルールを確認
- VPCのルーティングテーブルを確認
- ElastiCache/RDSのセキュリティグループでECSタスクのセキュリティグループからのアクセスを許可
