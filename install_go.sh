#!/bin/bash

# Go安装脚本
echo "开始安装Go环境..."

# 检查是否已安装Go
if command -v go &> /dev/null; then
    echo "Go已经安装"
    go version
    exit 0
fi

# 安装Homebrew (如果未安装)
/usr/bin/ruby -e "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install)"

# 安装Go
brew install go

# 验证安装
echo "验证Go安装..."
go version

# 设置环境变量
export GOROOT=/usr/local/go
export GOPATH=$HOME/go
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin

echo "Go安装完成!"
echo "请运行: source ~/.zshrc 或重启终端以应用环境变量"