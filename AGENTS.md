# Repository Guidelines

## プロジェクト構成とモジュール配置
このリポジトリは、最新の Git コミット差分を GitHub Copilot CLI でレビューする小規模な Go 製 CLI です。

- `cmd/review/`: エントリーポイント (`main.go`)
- `pkg/git/`: Git 実行まわりの補助処理
- `pkg/diff/`: diff の行数確認と解析
- `pkg/copilot/`: Copilot CLI 呼び出しラッパー
- `pkg/review/`: プロンプト生成とレビュー実行
- `README.md`: ビルド方法と使い方

新しい業務ロジックは `pkg/<name>` 配下に追加し、`cmd/` には実行ファイルの配線だけを置いてください。

## ビルド・テスト・開発コマンド
- `go run ./cmd/review`: ビルドせずに CLI を実行
- `go build -o review ./cmd/review`: ローカル用バイナリを作成
- `go test ./...`: 全パッケージのテストを実行
- `gofmt -w cmd pkg`: ソースコードを整形

実行には `git` と `copilot` コマンドが `PATH` にある必要があります。対象 diff は `HEAD~1..HEAD` で、差分が空、または 300 行超の場合は早期終了します。

## コーディング規約と命名
Go の標準規約に従ってください。インデントはタブ、整形は `gofmt`、パッケージ名は短く明確にし、公開識別子は他パッケージから必要なものに限定します。

コンストラクタは `NewRunner` のような説明的な名前を使い、API はできるだけ狭く保ってください。ファイル名は `diff.go`、`review.go` のように小文字で責務に合わせます。

## 設定とツール利用の注意
認証情報やモデル用トークンをコードに埋め込まないでください。`copilot` CLI と、その既存の認証状態を前提に実装します。
