# TODO Management Application

TODO管理Webアプリケーション - Next.js + Go + PostgreSQL

## 技術スタック

- **フロントエンド**: Next.js 14+ (App Router), TypeScript, Tailwind CSS
- **バックエンド**: Go, Gin, GORM
- **データベース**: PostgreSQL

## 機能

- ユーザー登録・ログイン（JWT認証）
- TODO作成・編集・削除・完了マーク
- 管理者によるユーザー管理

## 環境構築方法

### 必要なもの

- Docker
- Docker Compose

### ファイル構成

```
/
├── docker/
│   ├── backend.Dockerfile       # Go本番ビルド
│   ├── frontend.Dockerfile      # Next.js本番ビルド
│   ├── backend.dev.Dockerfile   # Go開発用（ホットリロード）
│   ├── frontend.dev.Dockerfile  # Next.js開発用（ホットリロード）
│   └── README.md                # Dockerセットアップ詳細
├── docker-compose.yml           # 本番用構成
├── docker-compose.dev.yml       # 開発用構成
├── .env.example                 # 環境変数テンプレート
├── backend/                     # Goバックエンド
└── frontend/                    # Next.jsフロントエンド
```

### クイックスタート

#### 1. リポジトリをクローン

```bash
git clone <repository-url>
cd TODO-App
```

#### 2. 環境変数の設定（任意）

```bash
cp .env.example .env
# .envファイルを編集してSESSION_SECRETを安全な値に変更
```

#### 3. 起動

**本番モード:**
```bash
docker-compose up --build
```

**開発モード（ホットリロード対応）:**
```bash
docker-compose -f docker-compose.dev.yml up --build
```

#### 4. アクセス

| サービス | URL |
|----------|-----|
| フロントエンド | http://localhost:3000 |
| バックエンドAPI | http://localhost:8080 |

### 停止方法

```bash
docker-compose down
```

### データを含めて完全削除

```bash
docker-compose down -v
```

## 初回セットアップ

1. http://localhost:3000/register にアクセス
2. 最初に登録したユーザーが自動的に管理者になります
3. 管理者はナビゲーションの「Admin」からユーザー管理が可能

## 環境変数

| 変数名 | 説明 | デフォルト値 |
|--------|------|-------------|
| `POSTGRES_USER` | DBユーザー名 | postgres |
| `POSTGRES_PASSWORD` | DBパスワード | postgres |
| `POSTGRES_DB` | DB名 | todo_app |
| `SESSION_SECRET` | JWT署名キー | (設定必須) |

## トラブルシューティング

### ポートが使用中の場合

`docker-compose.yml`のポート設定を変更:
```yaml
ports:
  - "3001:3000"  # フロントエンド
  - "8081:8080"  # バックエンド
```

### データベース接続エラー

PostgreSQLコンテナの起動を待ってください。healthcheckにより自動で処理されます。

## ライセンス

MIT
