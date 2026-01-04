#!/bin/bash

# 替换所有文件中的导入路径
for file in $(grep -r "xiaowo/internal" . --include="*.go" --exclude-dir=internal/test | cut -d: -f1 | sort | uniq); do
  echo "修改文件: $file"
  sed -i '' 's/xiaowo\/internal/xiaowo\/backend\/internal/g' "$file"
done

echo "所有导入路径已更新"