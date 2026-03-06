# review CLI

Git の最新コミット差分を取得し、Copilot (gpt-5-mini) でレビュー結果をストリーミング表示する CLI です。

## 使い方

```bash
go run ./cmd/review
```

## 仕様

- diff 取得: `git diff -W -U3 HEAD~1 HEAD`
- diff 行数が 300 行を超えた場合は `Diff too large (max 300 lines)` で終了
- diff が空の場合は `No diff found` を表示して終了

## 依存

- Git
- GitHub Copilot CLI (`copilot` コマンド)
