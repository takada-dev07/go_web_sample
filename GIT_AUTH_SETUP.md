# Git認証設定ガイド

## 概要

このドキュメントでは、GitHubリポジトリへのプッシュ時に発生する認証エラーを解決し、プロジェクトごとに異なるGitHubアカウントを使用する設定方法を説明します。

## 問題の症状

以下のようなエラーが発生する場合：

```
remote: Permission to takada-dev07/go_web_sample.git denied to hideaki-takada07.
fatal: unable to access 'https://github.com/takada-dev07/go_web_sample.git/': The requested URL returned error: 403
```

### 原因

- リモートリポジトリの所有者: `takada-dev07`
- 現在認証中のアカウント: `hideaki-takada07`
- アカウントの不一致により、権限エラー（403）が発生

## 解決方法：プロジェクトごとにアカウントを設定

プロジェクトごとに異なるGitHubアカウントを使用するには、リモートURLにユーザー名を含める方法が推奨されます。

### 手順

#### 1. リモートURLにユーザー名を含める

```bash
git remote set-url origin https://takada-dev07@github.com/takada-dev07/go_web_sample.git
```

#### 2. 設定の確認

```bash
git remote -v
```

以下のように表示されればOK：

```
origin https://takada-dev07@github.com/takada-dev07/go_web_sample.git (fetch)
origin https://takada-dev07@github.com/takada-dev07/go_web_sample.git (push)
```

#### 3. 既存の認証情報を削除（必要に応じて）

macOSのKeychainに保存されている古い認証情報を削除：

```bash
security delete-internet-password -s github.com
```

#### 4. プッシュを実行

```bash
git push
```

初回は認証情報の入力を求められます：

- **Username**: `takada-dev07`（URLに含まれているため自動入力される場合があります）
- **Password**: `takada-dev07`のPersonal Access Token（PAT）を入力

### Personal Access Token（PAT）の作成方法

1. GitHubに`takada-dev07`でログイン
2. Settings → Developer settings → Personal access tokens → Tokens (classic)
3. Generate new token (classic) をクリック
4. 以下を設定：
   - **Note**: `go_web_sample repository access`（任意の名前）
   - **Expiration**: お好みの期間
   - **Scopes**: `repo` にチェック（リポジトリへのフルアクセス）
5. Generate token をクリック
6. 表示されたトークンをコピー（再表示不可のため、必ず保存）

## 他のプロジェクトで別アカウントを使う場合

他のプロジェクトでも同様に、リモートURLにユーザー名を含めます：

```bash
git remote set-url origin https://別のアカウント名@github.com/ユーザー名/リポジトリ名.git
```

これにより、プロジェクトごとに異なるアカウントの認証情報を保存でき、他のプロジェクトに影響を与えません。

## メリット

- ✅ プロジェクトごとに異なるアカウントを使用可能
- ✅ 他のプロジェクトの認証情報に影響しない
- ✅ 認証情報はKeychainに自動保存される（次回以降は自動認証）

## 注意事項

- PATは再表示できないため、作成時に必ずコピーして安全な場所に保存してください
- PATには適切な有効期限を設定し、定期的に更新することを推奨します
- リポジトリへのアクセス権限がない場合は、リポジトリの所有者に権限を付与してもらう必要があります

## 参考

- [GitHub Personal Access Tokens](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token)
- [Git Credential Storage](https://git-scm.com/book/en/v2/Git-Tools-Credential-Storage)
