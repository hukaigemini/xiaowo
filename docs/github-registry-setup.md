# 🚀 GitHub Container Registry 配置指南

## 🎯 推荐方案：使用GitHub Container Registry

**为什么选择这个方案？**
- ✅ 稳定可靠，不需要额外平台
- ✅ 与GitHub完美集成
- ✅ 令牌创建简单直接
- ✅ 免费使用

## 📋 配置步骤

### 第一步：创建GitHub个人访问令牌

1. **打开GitHub令牌页面**
   - 访问：https://github.com/settings/tokens
   - 使用您的GitHub账户登录

2. **创建新令牌**
   - 点击 "Generate new token" 按钮
   - 选择 "Generate new token (classic)"

3. **配置令牌权限**
   - **Note**: `CI/CD for xiaowo project`
   - **Expiration**: 根据需要选择（建议90天）
   - **Select scopes** (勾选以下权限):
     - ✅ `write:packages` - 推送Docker镜像
     - ✅ `read:packages` - 拉取Docker镜像
     - ✅ `delete:packages` - 删除Docker镜像

4. **生成并保存令牌**
   - 点击 "Generate token"
   - **重要：立即复制生成的令牌**
   - 令牌格式类似：`ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx`

### 第二步：在GitHub仓库中配置Secrets

1. **进入您的xiaowo仓库**
2. **导航到Secrets设置**
   - 点击顶部的 `Settings` 标签
   - 左侧菜单：`Secrets and variables`
   - 选择 `Actions`

3. **添加第一个Secret**
   - 点击 `New repository secret`
   - **Name**: `DOCKER_REGISTRY`
   - **Secret**: `ghcr.io`
   - 点击 `Add secret`

4. **添加第二个Secret**
   - 点击 `New repository secret`
   - **Name**: `DOCKER_USERNAME`
   - **Secret**: 您的GitHub用户名（注意：是用户名，不是邮箱）
   - 点击 `Add secret`

5. **添加第三个Secret**
   - 点击 `New repository secret`
   - **Name**: `DOCKER_PASSWORD`
   - **Secret**: 刚才复制的GitHub令牌
   - 点击 `Add secret`

### 第三步：验证配置

配置完成后，您的Secrets列表应该显示：

```
DOCKER_REGISTRY    ✅ Recently updated
DOCKER_USERNAME    ✅ Recently updated
DOCKER_PASSWORD    ✅ Recently updated
```

## 🧪 测试配置

我们之前创建的工作流可以验证配置：

1. 进入GitHub Actions页面
2. 找到 `🔧 Secrets配置测试` 工作流
3. 点击 `Run workflow`
4. 等待运行完成

### 期望的成功输出：
```
✅ DOCKER_REGISTRY: ghcr.io
✅ DOCKER_USERNAME 已设置
✅ DOCKER_PASSWORD 已设置
✅ Docker登录成功
```

## 🔧 如果遇到问题

### 问题1：GitHub用户名vs邮箱
- **使用用户名，不是邮箱**
- 用户名在GitHub个人资料页面顶部显示

### 问题2：令牌权限不足
- 确保勾选了 `write:packages`, `read:packages`, `delete:packages`
- 重新生成令牌确保权限正确

### 问题3：令牌失效
- 令牌可能已过期
- 在GitHub设置中查看令牌状态

## 🎉 完成后的效果

配置成功后，您的CI/CD流水线可以：
- ✅ 构建Docker镜像
- ✅ 推送到GitHub Container Registry
- ✅ 使用最新的镜像进行部署
- ✅ 支持私有镜像仓库

## 📦 镜像访问

推送的镜像可以通过以下方式访问：
```
ghcr.io/您的用户名/xiaowo:latest
```

---

**您想开始配置吗？我可以一步步指导您完成每个步骤！**