# review CLI

Git の最新コミット差分を取得し、GitHub Copilot SDK for Go 経由でレビュー結果をストリーミング表示する CLI です。

## ビルド

```bash
go build -o review ./cmd/review
```

Windows の場合:

```powershell
go build -o review.exe .\cmd\review
```

## 使い方

```bash
go run ./cmd/review
```

プロンプトも確認する場合:

```bash
go run ./cmd/review --verbose
```

300 行制限を 3000 行まで拡張する場合:

```bash
go run ./cmd/review --limit-break
```

ビルド済みバイナリを使う場合:

```bash
./review
```

Windows の場合:

```powershell
.\review.exe
```

Windows でプロンプトも表示する場合:

```powershell
.\review.exe --verbose
```

Windows で 300 行制限を 3000 行まで拡張する場合:

```powershell
.\review.exe --limit-break
```

## 仕様

- diff 取得: `git diff -W -U3 HEAD~1 HEAD`
- diff 行数が 300 行を超えた場合は `Diff too large (<current> lines, max 300 lines)` で終了
- `--limit-break` 指定時は diff 行数上限を 3000 行に拡張
- diff が空の場合は `No diff found` を表示して終了
- `--verbose` 指定時は Copilot に送信するプロンプトを標準出力へ表示

## 依存

- Git
- Go 1.24 以上
- GitHub Copilot CLI (`copilot` コマンド)
