# review CLI

Git の最新コミット差分を取得し、Copilot (gpt-5-mini) でレビュー結果をストリーミング表示する CLI です。

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

## 仕様

- diff 取得: `git diff -W -U3 HEAD~1 HEAD`
- diff 行数が 300 行を超えた場合は `Diff too large (<current> lines, max 300 lines)` で終了
- diff が空の場合は `No diff found` を表示して終了
- `--verbose` 指定時は Copilot に送信するプロンプトを標準出力へ表示

## 依存

- Git
- GitHub Copilot CLI (`copilot` コマンド)
