# Discord Bot メッセージスケジューラ

Discord の Bot からのメッセージ送信をスケジューリングできるアプリのリポジトリです。  
特定の日時に Bot からメッセージを送信できます。

## 開発要件

以下決まっているものから記載

| 要件対象       | 詳細                                |
| -------------- | ----------------------------------- |
| バックエンド   | golang v1.17                        |
| フレームワーク | fibar v2.2                          |
| フロントエンド | テンプレートレンダリング + AdminLTE |
| DB             | MySQL v8（CloudSQL）                |
| インフラ       | サーバレスコンテナ想定              |
| ジョブサービス | Cron or CloudScheduler              |
| UML            | plantUML                            |

## 開発プロセス

分析と設計の成果物については、plantUML で記述を行いリポジトリで管理していく。

### 1. 分析

以下成果物

- 機能一覧
- 用語集
- 概念モデル図（ドメインモデル）
- ユースケース図（ユースケース記述は図中にメモとして記載）
- ロバストネス図

※ 今回 UI は AdminLTE で済ませるので画面・UI 設計は無し

### 2. 設計

以下成果物

- クラス図（主にドメインクラス・オブジェクト）
- ER 図（データモデル）

### 3. 実装

以下のプロセスで実装

1. バックエンド実装
2. フロント実装
3. テスト実装
