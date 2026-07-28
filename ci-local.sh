#!/usr/bin/env bash
# 本地跑一遍 CI 的全部檢查，與 .github/workflows/ci.yml 同步。
# runner 還在調整期間，用這支確認提交是乾淨的。
set -euo pipefail
cd "$(dirname "$0")"

echo "== 格式檢查 =="
unformatted=$(gofmt -l .)
[ -z "$unformatted" ] || { echo "未格式化：$unformatted"; exit 1; }
echo "ok"

echo "== go vet =="
go vet ./...
echo "ok"

echo "== 單元測試 =="
go test ./...

echo "== 零依賴檢查 =="
deps=$(go list -deps -test ./...)
ext=$(printf '%s\n' "$deps" \
  | grep -E '^[a-zA-Z0-9-]+(\.[a-zA-Z0-9-]+)+/' \
  | grep -v '^github.com/LukeLogix/destiny-core' \
  || true)
[ -z "$ext" ] || { echo "主 module 出現外部依賴："; echo "$ext"; exit 1; }
echo "✅ 零外部依賴（已含測試檔的 import）"

echo "== 交叉驗證（獨立 module）=="
(cd test && go test ./...)

echo
echo "全部通過。"
