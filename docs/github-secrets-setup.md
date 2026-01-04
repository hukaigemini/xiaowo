# GitHub Secrets 配置指南

## 🔐 必需的GitHub Secrets配置

### 1. Docker镜像仓库配置
```
DOCKER_REGISTRY=your-registry.com
DOCKER_USERNAME=your-username
DOCKER_PASSWORD=your-password-or-token
```

**支持的Docker镜像仓库**:
- Docker Hub: `docker.io`
- GitHub Container Registry: `ghcr.io`
- 阿里云容器镜像服务: `registry.cn-hangzhou.aliyuncs.com`
- 腾讯云容器镜像: `ccr.ccs.tencentyun.com`

### 2. Kubernetes集群配置
```
KUBE_CONFIG=<base64-encoded-kubeconfig>
KUBE_CONFIG_PROD=<base64-encoded-production-kubeconfig>
```

**配置步骤**:
```bash
# 编码kubeconfig文件
base64 -w 0 ~/.kube/config > kubeconfig-encoded.txt

# 生产环境kubeconfig
base64 -w 0 ~/.kube/config-prod > kubeconfig-prod-encoded.txt
```

### 3. 通知配置
```
SLACK_WEBHOOK=https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK
```

### 4. 代码质量工具
```
SONAR_TOKEN=<your-sonarqube-token>
SNYK_TOKEN=<your-snyk-token>
```

### 5. 环境配置
```
STAGING_ENV_CONFIG=<base64-encoded-staging-config>
```

## 🛠️ 配置步骤

### 步骤1: 访问GitHub仓库设置
1. 进入GitHub仓库页面
2. 点击 `Settings` → `Secrets and variables` → `Actions`
3. 选择 `New repository secret`

### 步骤2: 批量创建Secrets
使用GitHub CLI批量创建（推荐）:

```bash
# 设置必要的环境变量
export DOCKER_REGISTRY="your-registry.com"
export DOCKER_USERNAME="your-username"
export DOCKER_PASSWORD="your-password"
export SLACK_WEBHOOK="your-slack-webhook"

# 创建Secrets
gh secret set DOCKER_REGISTRY --body "$DOCKER_REGISTRY"
gh secret set DOCKER_USERNAME --body "$DOCKER_USERNAME"
gh secret set DOCKER_PASSWORD --body "$DOCKER_PASSWORD"
gh secret set SLACK_WEBHOOK --body "$SLACK_WEBHOOK"
```

### 步骤3: 验证配置
在GitHub Actions中运行测试workflow验证Secrets配置。

## 🔍 常见问题

### Q: 如何获取Docker Hub访问令牌？
A: 访问 Docker Hub → Account Settings → Security → New Access Token

### Q: 如何创建Slack Webhook？
A: Slack → App Directory → Incoming Webhooks → Create New App

### Q: Kubernetes配置安全吗？
A: 所有配置都存储在GitHub的Secrets中，具有企业级安全性。

## 📞 支持联系
如需技术支持，请联系基础设施团队。
