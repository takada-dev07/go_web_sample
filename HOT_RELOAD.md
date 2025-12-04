# ホットリロード設定ガイド

このプロジェクトでは、開発時のホットリロード機能をサポートしています。ファイルを変更すると、自動的にアプリケーションが再ビルド・再起動されます。

## ホットリロードツールの比較

### 1. Air（推奨）⭐

**特徴:**

- 最も一般的で人気のあるGo用ホットリロードツール
- 設定が簡単で柔軟
- ファイル変更の検知が高速
- カスタマイズ可能な設定ファイル（`.air.toml`）

**インストール:**

```bash
go install github.com/cosmtrek/air@latest
```

**使用方法:**

```bash
# ローカルで実行
air

# または設定ファイルを指定
air -c .air.toml
```

**メリット:**

- 設定が豊富
- コミュニティが大きい
- ドキュメントが充実

**デメリット:**

- やや重い（他の選択肢と比較）

---

### 2. Fresh

**特徴:**

- 軽量でシンプル
- 設定不要で即座に使用可能
- ファイル変更の検知が高速

**インストール:**

```bash
go install github.com/gravityblast/fresh@latest
```

**使用方法:**

```bash
fresh
```

**メリット:**

- 設定不要
- 軽量
- シンプル

**デメリット:**

- カスタマイズが限定的
- メンテナンスが活発でない

---

### 3. Realize

**特徴:**

- 機能が豊富（ビルド、テスト、実行を統合）
- Web UI付き
- 複数のプロジェクトを管理可能

**インストール:**

```bash
go install github.com/oxequa/realize@latest
```

**使用方法:**

```bash
realize start
```

**メリット:**

- 機能が豊富
- Web UIで管理可能
- 複数プロジェクト対応

**デメリット:**

- 設定が複雑
- やや重い

---

### 4. CompileDaemon

**特徴:**

- シンプルで軽量
- 設定が最小限

**インストール:**

```bash
go install github.com/githubnemo/CompileDaemon@latest
```

**使用方法:**

```bash
CompileDaemon -command="./tmp/main"
```

**メリット:**

- 軽量
- シンプル

**デメリット:**

- 機能が限定的
- 設定オプションが少ない

---

### 5. カスタムスクリプト（watcher + go run）

**特徴:**

- 完全にカスタマイズ可能
- 外部依存なし

**実装例:**

```bash
#!/bin/bash
while true; do
    go run ./cmd/server/main.go &
    PID=$!
    inotifywait -e modify -r .
    kill $PID
done
```

**メリット:**

- 完全な制御
- 外部ツール不要

**デメリット:**

- 実装が必要
- メンテナンスが必要

---

## 現在の設定（Air）

このプロジェクトでは**Air**を使用しています。

### 設定ファイル

`.air.toml`ファイルで設定を管理しています。主な設定項目：

- `build.cmd`: ビルドコマンド
- `build.bin`: ビルド後のバイナリパス
- `build.delay`: ファイル変更検知後の待機時間
- `build.exclude_dir`: 監視から除外するディレクトリ
- `build.include_ext`: 監視するファイル拡張子

### 使用方法

#### ローカル環境

```bash
# Airをインストール
go install github.com/cosmtrek/air@latest

# ホットリロードで起動
air
```

#### Docker環境

```bash
# docker-composeで起動（自動的にAirが使用される）
docker-compose up -d

# ログを確認
docker-compose logs -f app
```

### カスタマイズ

`.air.toml`を編集して、以下の設定を変更できます：

- ビルドコマンド
- 除外するディレクトリ/ファイル
- 監視するファイル拡張子
- ログの表示設定

詳細は[Airの公式ドキュメント](https://github.com/cosmtrek/air)を参照してください。

---

## 他のツールに切り替える場合

### Freshに切り替える場合

1. `.air.toml`を削除
2. `Dockerfile`の`local`ステージを修正：

   ```dockerfile
   RUN go install github.com/gravityblast/fresh@latest
   CMD ["fresh"]
   ```

### Realizeに切り替える場合

1. `.air.toml`を削除
2. `Dockerfile`の`local`ステージを修正：

   ```dockerfile
   RUN go install github.com/oxequa/realize@latest
   CMD ["realize", "start"]
   ```

---

## トラブルシューティング

### ファイル変更が検知されない

- `.air.toml`の`exclude_dir`を確認
- ボリュームマウントが正しく設定されているか確認（Dockerの場合）

### ビルドエラーが表示されない

- `.air.toml`の`log`設定を確認
- `build-errors.log`ファイルを確認

### 再起動が頻繁に発生する

- `build.delay`を増やす
- `build.exclude_dir`に不要なディレクトリを追加
