# 📚 GitHub Secrets 手把手配置指南

## 🎯 目标
为CI/CD流水线配置必需的GitHub Secrets，确保自动化部署正常运行。

## 📋 配置清单

### ✅ 必需Secrets（必须配置）

| Secret名称 | 用途 | 如何获取 |
|------------|------|----------|
| `DOCKER_REGISTRY` | Docker镜像仓库地址 | 您的镜像仓库URL |
| `DOCKER_USERNAME` | Docker仓库用户名 | 您的仓库登录名 |
| `DOCKER_PASSWORD` | Docker仓库密码/令牌 | 您的密码或访问令牌 |
| `KUBE_CONFIG` | Kubernetes集群配置 | base64编码的kubeconfig文件 |

### ⚠️ 可选Secrets（建议配置）

| Secret名称 | 用途 | 如何获取 |
|------------|------|----------|
| `KUBE_CONFIG_PROD` | 生产环境K8s配置 | 生产环境kubeconfig |
| `SLACK_WEBHOOK` | Slack通知URL | Slack webhook地址 |
| `SONAR_TOKEN` | SonarQube分析令牌 | SonarQube项目令牌 |
| `SNYK_TOKEN` | 安全扫描令牌 | Snyk项目令牌 |

---

## 🚀 详细配置步骤

### 第一步：进入GitHub Secrets设置

1. **打开GitHub仓库**
   - 登录GitHub
   - 进入您的`xiaowo`项目仓库

2. **导航到Actions Secrets**
   - 点击顶部的 `Settings` 标签
   - 在左侧菜单中点击 `Secrets and variables`
   - 选择 `Actions`

### 第二步：配置Docker相关Secrets

#### 2.1 配置DOCKER_REGISTRY

1. 点击 `New repository secret` 按钮
2. **名称**: `DOCKER_REGISTRY`
3. **Secret值**: 
   - Docker Hub: `docker.io`
   - GitHub Container Registry: `ghcr.io`
   - 阿里云容器镜像: `registry.cn-hangzhou.aliyuncs.com`
   - 腾讯云容器镜像: `ccr.ccs.tencentyun.com`
4. 点击 `Add secret` 保存

#### 2.2 配置DOCKER_USERNAME

1. 点击 `New repository secret` 按钮
2. **名称**: `DOCKER_USERNAME`
3. **Secret值**: 您的Docker仓库用户名
4. 点击 `Add secret` 保存

#### 2.3 配置DOCKER_PASSWORD

1. 点击 `New repository secret` 按钮
2. **名称**: `DOCKER_PASSWORD`
3. **Secret值**: 
   - **推荐使用访问令牌而不是密码**
   - Docker Hub访问令牌获取：
     - 登录Docker Hub
     - Account Settings → Security → Access Tokens
     - 创建新的访问令牌
4. 点击 `Add secret` 保存

### 第三步：配置Kubernetes相关Secrets

#### 3.1 获取kubeconfig文件

在您的本地机器上执行：

```bash
# 查看当前kubeconfig
kubectl config view --minify

# 导出并Base64编码
kubectl config view --minify --raw > kubeconfig-temp.yaml
base64 -w 0 kubeconfig-temp.yaml
```

#### 3.2 配置KUBE_CONFIG

1. 点击 `New repository secret` 按钮
2. **名称**: `KUBE_CONFIG`
3. **Secret值**: 粘贴上一步的Base64编码结果
4. 点击 `Add secret` 保存

#### 3.3 配置KUBE_CONFIG_PROD (可选)

如果有生产环境kubeconfig：

```bash
# 导出生产环境配置
kubectl config view --minify --raw --context production > kubeconfig-prod.yaml
base64 -w 0 kubeconfig-prod.yaml
```

然后创建对应的Secret。

### 第四步：配置通知相关Secrets

#### 4.1 获取Slack Webhook URL (可选)

1. 在Slack中创建Incoming Webhook
2. 复制WebHook URL
3. 创建Secret:
   - **名称**: `SLACK_WEBHOOK`
   - **值**: 您的Slack Webhook URL

### 第五步：验证配置

配置完成后，使用我们创建的测试工作流来验证：

1. 进入GitHub Actions页面
2. 选择 `🔧 Secrets配置测试` 工作流
3. 点击 `Run workflow` 按钮
4. 查看运行结果

---

## 🔍 常见问题解决

### Docker登录失败

**问题**: `docker login failed`
**解决**:
- 检查DOCKER_REGISTRY是否正确
- 确认DOCKER_USERNAME和DOCKER_PASSWORD有效
- 尝试使用访问令牌而不是密码

### Kubernetes连接失败

**问题**: `kubectl cluster-info failed`
**解决**:
- 确认kubeconfig文件格式正确
- 检查Base64编码是否正确
- 验证集群访问权限

### Secrets未生效

**问题**: Secrets显示为空
**解决**:
- 确认Secrets名称完全匹配
- 等待几分钟让配置生效
- 重新运行测试工作流

---

## 📞 需要帮助？

如果您在配置过程中遇到任何问题，请告诉我具体的错误信息，我会帮您解决！

## ✅ 配置完成后

配置完成后，请运行我们创建的测试工作流来验证一切配置正确。如果所有测试通过，您就可以开始使用完整的CI/CD流水线了！