#!/bin/bash

set -e

# 設定変数（環境に応じて変更）
AWS_REGION="ap-northeast-1"
AWS_ACCOUNT_ID=$(aws sts get-caller-identity --query Account --output text)
ECR_REPOSITORY_NAME="go-web-sample"
ECS_CLUSTER_NAME="go-web-sample-cluster"
ECS_SERVICE_NAME="go-web-sample-service"
ECS_TASK_DEFINITION_FAMILY="go-web-sample"
IMAGE_TAG=${IMAGE_TAG:-latest}

# 環境ステージ（dev または prd）
ENV_STAGE=${ENV_STAGE:-prd}

# ステージの検証
if [ "$ENV_STAGE" != "dev" ] && [ "$ENV_STAGE" != "prd" ]; then
  echo "エラー: ENV_STAGEは 'dev' または 'prd' である必要があります"
  exit 1
fi

echo "=== ECS Fargate デプロイスクリプト ==="
echo "Region: $AWS_REGION"
echo "Account ID: $AWS_ACCOUNT_ID"
echo "Repository: $ECR_REPOSITORY_NAME"
echo "Image Tag: $IMAGE_TAG"
echo "Environment Stage: $ENV_STAGE"
echo ""

# ECRリポジトリURI
ECR_REPOSITORY_URI="${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/${ECR_REPOSITORY_NAME}"

# 1. ECRにログイン
echo "1. ECRにログイン中..."
aws ecr get-login-password --region $AWS_REGION | docker login --username AWS --password-stdin $ECR_REPOSITORY_URI

# 2. 環境ステージに応じたDockerイメージをビルド
echo "2. Dockerイメージをビルド中（ステージ: $ENV_STAGE）..."
docker build --target ${ENV_STAGE} -f Dockerfile -t ${ECR_REPOSITORY_NAME}:${IMAGE_TAG} .

# 3. イメージにタグを付与
echo "3. イメージにタグを付与中..."
docker tag ${ECR_REPOSITORY_NAME}:${IMAGE_TAG} ${ECR_REPOSITORY_URI}:${IMAGE_TAG}
docker tag ${ECR_REPOSITORY_NAME}:${IMAGE_TAG} ${ECR_REPOSITORY_URI}:latest

# 4. ECRにプッシュ
echo "4. ECRにイメージをプッシュ中..."
docker push ${ECR_REPOSITORY_URI}:${IMAGE_TAG}
docker push ${ECR_REPOSITORY_URI}:latest

# 5. タスク定義を更新
echo "5. タスク定義を更新中..."
# タスク定義ファイルのプレースホルダーを置換
sed -i.bak "s|YOUR_ACCOUNT_ID|${AWS_ACCOUNT_ID}|g" deploy/task-definition.json
sed -i.bak "s|YOUR_ECR_REPOSITORY_URI|${ECR_REPOSITORY_URI}|g" deploy/task-definition.json
sed -i.bak "s|REGION|${AWS_REGION}|g" deploy/task-definition.json

# タスク定義を登録
TASK_DEFINITION_ARN=$(aws ecs register-task-definition \
  --cli-input-json file://deploy/task-definition.json \
  --query 'taskDefinition.taskDefinitionArn' \
  --output text)

echo "タスク定義が登録されました: $TASK_DEFINITION_ARN"

# 6. ECSサービスを更新（サービスが存在する場合）
if aws ecs describe-services --cluster $ECS_CLUSTER_NAME --services $ECS_SERVICE_NAME --region $AWS_REGION > /dev/null 2>&1; then
  echo "6. ECSサービスを更新中..."
  aws ecs update-service \
    --cluster $ECS_CLUSTER_NAME \
    --service $ECS_SERVICE_NAME \
    --task-definition $TASK_DEFINITION_ARN \
    --force-new-deployment \
    --region $AWS_REGION > /dev/null
  
  echo "サービス更新が開始されました。デプロイの進行状況を確認してください:"
  echo "aws ecs describe-services --cluster $ECS_CLUSTER_NAME --services $ECS_SERVICE_NAME --region $AWS_REGION"
else
  echo "6. ECSサービスが見つかりません。手動で作成してください。"
fi

# バックアップファイルを削除
rm -f deploy/task-definition.json.bak

echo ""
echo "=== デプロイ完了 ==="

