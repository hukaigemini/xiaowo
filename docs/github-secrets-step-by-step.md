# GitHub Secrets 配置操作指南

## 🎯 完整操作步骤

### 方法一：GitHub Web界面配置（推荐）

#### 步骤1: 访问仓库设置
1. 进入GitHub仓库页面
2. 点击右上角的 `Settings` 标签
3. 在左侧菜单中找到 `Secrets and variables`
4. 点击 `Actions`

#### 步骤2: 创建新的Secrets
1. 点击 `New repository secret` 按钮
2. 在 `Name` 字段输入密钥名称（如 `DOCKER_REGISTRY`）
3. 在 `Secret` 字段输入密钥值
4. 点击 `Add secret` 保存

#### 步骤3: 批量创建Secrets
按照上述步骤，逐一创建以下16个Secrets：

```
Name: DOCKER_REGISTRY
Secret: registry.cn-hangzhou.aliyuncs.com

Name: DOCKER_USERNAME  
Secret: your-docker-username

Name: DOCKER_PASSWORD
Secret: your-docker-password-or-token

Name: KUBE_CONFIG
Secret: <base64-encoded-kubeconfig>

Name: KUBE_CONFIG_PROD
Secret: <base64-encoded-prod-kubeconfig>

Name: SLACK_WEBHOOK
Secret: https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK

Name: SONAR_TOKEN
Secret: your-sonarqube-token

Name: SNYK_TOKEN
Secret: your-snyk-token

Name: STAGING_ENV_CONFIG
Secret: <base64-encoded-staging-config>

Name: ALIYUN_ACCESS_KEY_ID
Secret: your-access-key-id

Name: ALIYUN_ACCESS_KEY_SECRET
Secret: your-access-key-secret

Name: PROD_DB_PASSWORD
Secret: your-production-db-password

Name: JWT_SECRET
Secret: your-jwt-secret-key

Name: REDIS_PASSWORD
Secret: your-redis-password

Name: SMTP_PASSWORD
Secret: your-smtp-password

Name: PROMETHEUS_TOKEN
Secret: your-prometheus-token
```

### 方法二：GitHub CLI配置（高级用户）

#### 前置条件
```bash
# 安装GitHub CLI
# macOS: brew install gh
# Ubuntu: sudo apt install gh

# 登录GitHub
gh auth login
```

#### 批量配置脚本
```bash
#!/bin/bash
# secrets-setup.sh

echo "🔐 开始配置GitHub Secrets..."

# 设置环境变量（请根据实际情况修改）
export DOCKER_REGISTRY="registry.cn-hangzhou.aliyuncs.com"
export DOCKER_USERNAME="your-username"
export DOCKER_PASSWORD="your-password"
export SLACK_WEBHOOK="https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK"

# 创建Secrets
gh secret set DOCKER_REGISTRY --body "$DOCKER_REGISTRY"
gh secret set DOCKER_USERNAME --body "$DOCKER_USERNAME"
gh secret set DOCKER_PASSWORD --body "$DOCKER_PASSWORD"
gh secret set SLACK_WEBHOOK --body "$SLACK_WEBHOOK"

echo "✅ 基础Secrets配置完成"
echo "请手动配置其他需要特殊值的Secrets（如kubeconfig等）"
```

## 🔍 验证配置

### 检查Secrets是否正确设置
```bash
# 列出所有已配置的Secrets
gh secret list

# 验证特定Secret
gh secret get DOCKER_REGISTRY
```

### 测试配置
1. 创建一个简单的测试分支
2. 提交代码触发CI/CD流水线
3. 检查Actions页面中的流水线执行情况
4. 确认各个步骤是否成功

## 🚨 常见问题和解决方案

### Q1: 提示"Secrets not found"
**解决方案**: 
- 检查Secret名称是否完全匹配（大小写敏感）
- 确认在正确的仓库中配置
- 等待1-2分钟让配置生效

### Q2: 流水线执行失败，提示权限错误
**解决方案**:
- 检查Docker仓库权限
- 确认Kubernetes配置是否有效
- 验证Slack Webhook URL格式

### Q3: Kubernetes配置解码失败
**解决方案**:
```bash
# 验证base64编码
echo "your-base64-string" | base64 -d

# 重新编码kubeconfig
kubectl config view --raw | base64 -w 0
```

### Q4: 镜像推送失败
**解决方案**:
- 确认Docker仓库地址正确
- 检查用户名和密码是否有效
- 确认有推送权限

## 📊 配置完成检查清单

- [ ] DOCKER_REGISTRY 已配置
- [ ] DOCKER_USERNAME 已配置  
- [ ] DOCKER_PASSWORD 已配置
- [ ] KUBE_CONFIG 已配置（开发环境）
- [ ] KUBE_CONFIG_PROD 已配置（生产环境）
- [ ] SLACK_WEBHOOK 已配置
- [ ] SONAR_TOKEN 已配置（可选）
- [ ] SNYK_TOKEN 已配置（可选）
- [ ] STAGING_ENV_CONFIG 已配置
- [ ] 其他7个Secrets根据需要配置

**配置完成后，记得运行一次完整的CI/CD测试来验证所有配置是否正确！**
