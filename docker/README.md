# Docker環境でのセットアップ

## 必要なもの
- Docker
- Docker Compose

## 起動方法

### 1. リポジトリをクローン
```bash
git clone <repository-url>
cd TODO-App
```

### 2. 環境変数の設定（オプション）
```bash
cp .env.example .env
# .envファイルを編集してSESSION_SECRETを安全な値に変更
```

### 3. Docker Composeで起動
```bash
docker-compose up --build
```

### 4. アクセス
- フロントエンド: http://localhost:3000
- バックエンドAPI: http://localhost:8080

## 停止方法
```bash
docker-compose down
```

## データを含めて完全に削除
```bash
docker-compose down -v
```

## 開発モード
開発時はホットリロードを有効にするため、以下のコマンドを使用：

```bash
docker-compose -f docker-compose.dev.yml up --build
```

### GraphQL開発

バックエンドはgqlgenを使用したGraphQLレイヤーを持っています。

**スキーマの変更時：**
開発モードでは、`backend/graph/schema.graphqls`を編集すると、Airが自動的にコード再生成とリビルドを行います。

**手動でコード生成を実行する場合：**
```bash
docker-compose -f docker-compose.dev.yml exec backend go generate ./...
```

**GraphQL構成ファイル：**
- `backend/graph/schema.graphqls` - GraphQLスキーマ定義
- `backend/gqlgen.yml` - gqlgen設定ファイル
- `backend/graph/schema.resolvers.go` - リゾルバー実装

## トラブルシューティング

### ポートが使用中の場合
docker-compose.ymlのポート設定を変更：
```yaml
ports:
  - "3001:3000"  # フロントエンド
  - "8081:8080"  # バックエンド
```

### データベース接続エラー
PostgreSQLコンテナが完全に起動するまで待ってください。
healthcheckが設定されているため、通常は自動的に処理されます。
