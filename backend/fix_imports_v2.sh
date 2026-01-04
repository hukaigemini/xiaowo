#!/bin/bash

# 替换所有文件中的导入路径，在路径前加上backend
for file in $(grep -r "xiaowo/backend" . --include="*.go" --exclude-dir=internal/test | cut -d: -f1 | sort | uniq); do
  echo "修改文件: $file"
  sed -i '' 's/xiaowo\/backend\/internal/backend\/internal/g' "$file"
done

echo "所有导入路径已更新"