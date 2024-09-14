#!/bin/bash

# 检查是否提供了版本号
if [ $# -eq 0 ]; then
    echo "请提供版本号，例如: ./release.sh v1.0.0"
    exit 1
fi

VERSION=$1

# 确保工作目录干净
if [[ -n $(git status -s) ]]; then
    echo "工作目录不干净，请提交或存储您的更改"
    exit 1
fi

# 创建并推送 tag
git tag -a $VERSION -m "Release $VERSION"
git push origin $VERSION

echo "已创建并推送 tag $VERSION。GitHub Actions 将开始构建和发布过程。"