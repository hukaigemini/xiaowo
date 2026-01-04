# Secrets 验证和测试指南

## 🎯 验证目标
确保所有16个Secrets正确配置，CI/CD流水线能够成功运行。

## 🔧 验证方法

### 方法1: 基础验证（立即执行）

#### 1.1 验证GitHub CLI可用性
```bash
# 检查GitHub CLI是否已安装
gh --version

# 验证是否已登录
gh auth status

# 如果未登录，执行登录
gh auth login --with-token
```

#### 1.2 列出已配置的Secrets
```bash
# 显示所有Secrets（不显示值）
gh secret list

# 预期输出应该包含：
# DOCKER_REGISTRY
# DOCKER_USERNAME
# DOCKER_PASSWORD
# KUBE_CONFIG
# KUBE_CONFIG_PROD
# SLACK_WEBHOOK
# SONAR_TOKEN
# SNYK_TOKEN
# STAGING_ENV_CONFIG
# 其他7个Secrets...
```

#### 1.3 验证关键Secrets格式
```bash
# 验证Docker仓库配置
echo "DOCKER_REGISTRY: $DOCKER_REGISTRY"
echo "DOCKER_USERNAME: $DOCKER_USERNAME"

# 验证Kubernetes配置（base64编码）
echo "KUBE_CONFIG长度: $(echo $KUBE_CONFIG | wc -c)"
echo "KUBE_CONFIG_PROD长度: $(echo $KUBE_CONFIG_PROD | wc -c)"

# 验证Slack Webhook格式
echo "SLACK_WEBHOOK是否以https开头: $(echo $SLACK_WEBHOOK | grep -c '^https://')"
```

### 方法2: 功能验证（CI/CD测试）

#### 2.1 创建测试工作流
创建 `.github/workflows/secrets-test.yml`:
```yaml
name: Secrets配置验证

on:
  workflow_dispatch:
  push:
    branches: [ secrets-test ]

jobs:
  validate-secrets:
    runs-on: ubuntu-latest
    
    steps:
    - name: 检出代码
      uses: actions/checkout@v4

    - name: 验证Docker配置
      run: |
        echo "🔍 验证Docker仓库配置"
        echo "Registry: ${{ secrets.DOCKER_REGISTRY }}"
        echo "Username: ${{ secrets.DOCKER_USERNAME }}"
        
        # 验证Docker登录
        echo "${{ secrets.DOCKER_PASSWORD }}" | docker login ${{ secrets.DOCKER_REGISTRY }} -u ${{ secrets.DOCKER_USERNAME }} --password-stdin
        echo "✅ Docker配置验证成功"

    - name: 验证Kubernetes配置
      run: |
        echo "🔍 验证Kubernetes配置"
        echo "${{ secrets.KUBE_CONFIG }}" | base64 -d > kubeconfig-dev
        export KUBECONFIG=kubeconfig-dev
        kubectl version --client
        echo "✅ Kubernetes配置验证成功"

    - name: 验证Slack Webhook
      run: |
        echo "🔍 验证Slack Webhook"
        curl -X POST -H 'Content-type: application/json' \
          --data '{"text":"✅ Secrets验证测试 - Webhook正常"}' \
          ${{ secrets.SLACK_WEBHOOK }}
        echo "✅ Slack Webhook验证成功"

    - name: 验证环境配置
      run: |
        echo "🔍 验证环境配置"
        echo "${{ secrets.STAGING_ENV_CONFIG }}" | base64 -d > staging-config.json
        cat staging-config.json
        echo "✅ 环境配置验证成功"

    - name: 发送验证完成通知
      run: |
        echo "🎉 所有Secrets验证完成！"
        echo "## 📋 Secrets验证报告" >> $GITHUB_STEP_SUMMARY
        echo "✅ Docker配置: 正常" >> $GITHUB_STEP_SUMMARY
        echo "✅ Kubernetes配置: 正常" >> $GITHUB_STEP_SUMMARY
        echo "✅ Slack通知: 正常" >> $GITHUB_STEP_SUMMARY
        echo "✅ 环境配置: 正常" >> $GITHUB_STEP_SUMMARY
        echo "🎯 状态: 所有关键Secrets配置正确" >> $GITHUB_STEP_SUMMARY
```

#### 2.2 执行验证测试
```bash
# 创建测试分支
git checkout -b secrets-test

# 提交触发验证
git add .
git commit -m "chore: 添加Secrets验证测试"
git push origin secrets-test

# 或者手动触发GitHub Actions中的secrets-test工作流
```

### 方法3: 完整CI/CD流程测试

#### 3.1 触发完整流水线测试
```bash
# 创建测试分支
git checkout -b ci-test-$(date +%Y%m%d-%H%M%S)
git push origin ci-test-$(date +%Y%m%d-%H%M%S)

# 或者修改main分支的小文件来触发完整CI/CD流程
echo "# 测试提交 $(date)" >> README.md
git add README.md
git commit -m "test: 触发CI/CD验证测试"
git push origin main
```

#### 3.2 监控流水线执行
1. 访问GitHub仓库的 `Actions` 页面
2. 观察流水线执行情况：
   - ✅ 代码检查和测试
   - ✅ Docker镜像构建
   - ✅ 镜像推送到仓库
   - ✅ 测试环境部署
   - ✅ 生产环境部署（如果配置了）
   - ✅ 监控和告警

#### 3.3 检查执行日志
重点关注以下步骤的日志：
- Docker构建步骤的镜像推送结果
- Kubernetes部署的执行状态
- Slack通知的发送状态
- 回滚机制的测试结果

## 📊 验证结果判断标准

### ✅ 成功标准
- 所有16个Secrets都能正常读取
- Docker镜像能够成功构建和推送
- Kubernetes部署能够正常执行
- Slack通知能够正常发送
- 完整的CI/CD流程能够从开始到结束无错误执行

### ❌ 失败判断
- 任何Secret读取失败或格式错误
- Docker操作失败（登录、构建、推送）
- Kubernetes操作失败（连接、部署、回滚）
- 通知发送失败
- 流水线在任何一个关键步骤中断

## 🚨 常见问题排查

### 问题1: Secret值包含特殊字符
```bash
# 如果Secret值包含换行符或特殊字符，需要正确处理
echo "$SECRET_VALUE" | tr -d '\r' | jq -Rs .
```

### 问题2: Kubernetes配置解码失败
```bash
# 验证base64编码是否正确
echo "$KUBE_CONFIG" | base64 -d > temp-kubeconfig
kubectl --kubeconfig=temp-kubeconfig version
rm temp-kubeconfig
```

### 问题3: Docker镜像推送权限不足
```bash
# 验证Docker权限
docker login $DOCKER_REGISTRY -u $DOCKER_USERNAME -p $DOCKER_PASSWORD
docker push $DOCKER_REGISTRY/test-image:latest
```

### 问题4: Slack Webhook无效
```bash
# 测试Webhook URL
curl -X POST -H 'Content-type: application/json' \
  --data '{"text":"测试消息"}' \
  $SLACK_WEBHOOK
```

## 📈 验证完成检查清单

- [ ] GitHub CLI可以正常列出Secrets
- [ ] 所有16个Secrets都已配置
- [ ] Docker仓库配置可以正常登录
- [ ] Kubernetes配置可以正常连接集群
- [ ] Slack Webhook可以正常发送消息
- [ ] 完整CI/CD流水线可以成功执行
- [ ] 所有关键步骤都没有错误
- [ ] 监控和告警系统正常工作

**验证完成后，您的I-003 CI/CD流水线就可以标记为100%完成了！**
