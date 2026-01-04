# 小窝项目开发规范

## 📋 文档信息
- **项目名称**: 小窝同步观影平台
- **文档版本**: v1.0
- **适用范围**: 全体开发团队成员
- **最后更新**: 2025-12-30

---

## 🔄 Git 分支规范

### 分支命名规则

#### 主要分支
- `main`: 主分支，生产环境代码
- `develop`: 开发分支，日常开发集成
- `feature/*`: 功能分支，新功能开发
- `bugfix/*`: 缺陷修复分支
- `hotfix/*`: 紧急修复分支
- `release/*`: 发布准备分支

#### 命名格式
```
分支类型/分支描述/版本号或日期

examples:
- feature/websocket-sync-v1.0
- bugfix/room-cleanup-issue
- hotfix/security-vulnerability
- release/v1.0.0
```

### 提交信息规范

#### 提交信息格式
```
类型(范围): 简短描述

详细描述（可选）

相关Issue: #123
```

#### 提交类型
- `feat`: 新功能
- `fix`: 缺陷修复
- `docs`: 文档更新
- `style`: 代码格式调整
- `refactor`: 代码重构
- `test`: 测试相关
- `chore`: 构建过程或辅助工具的变动
- `perf`: 性能优化
- `ci`: 持续集成相关

#### 示例
```bash
# 功能开发
feat(room): 添加房间成员管理功能
- 实现成员列表展示
- 添加成员权限控制
- 支持成员踢出功能
相关Issue: #15

# 缺陷修复
fix(websocket): 修复WebSocket连接断开后重连失败
- 添加自动重连机制
- 优化连接状态管理
- 增加重连次数限制
修复Issue: #28

# 文档更新
docs(api): 更新用户注册API文档
- 添加请求参数示例
- 更新响应格式说明
相关Issue: #12
```

### 代码审查流程

#### Pull Request 要求
1. **必须填写 PR 描述**
   - 功能描述
   - 测试情况
   - 相关Issue链接

2. **代码审查标准**
   - 代码功能正确性
   - 代码可读性和维护性
   - 性能影响评估
   - 安全性检查

3. **审查通过条件**
   - 至少1人Code Review通过
   - 所有CI/CD检查通过
   - 单元测试覆盖率不低于80%

---

## 🚨 错误码规范

### 错误码结构

#### 格式定义
```
XXXYYYZZZ

- XXX: 模块代码 (3位)
- YYY: 业务代码 (3位)  
- ZZZ: 错误类型 (3位)
```

#### 模块代码定义
```
100: 用户模块 (User)
200: 房间模块 (Room)
300: 消息模块 (Message)
400: WebSocket模块 (WebSocket)
500: 数据库模块 (Database)
600: 系统模块 (System)
```

#### 业务代码定义
```
001: 业务逻辑错误
002: 参数验证错误
003: 权限错误
004: 资源不存在
005: 状态错误
006: 业务规则违反
```

#### 错误类型定义
```
001: 未知错误
002: 参数错误
003: 权限不足
004: 资源不存在
005: 状态冲突
006: 业务逻辑错误
007: 系统异常
008: 网络错误
009: 超时错误
010: 资源冲突
```

### 常见错误码示例

#### 用户模块错误 (100)
```
100001001: 用户模块 - 未知错误
100002002: 用户模块 - 参数错误 (无效的用户名格式)
100003003: 用户模块 - 权限不足 (未登录访问)
100004004: 用户模块 - 资源不存在 (用户不存在)
100006006: 用户模块 - 业务规则违反 (用户名已存在)
```

#### 房间模块错误 (200)
```
200001001: 房间模块 - 未知错误
200002002: 房间模块 - 参数错误 (房间名过长)
200003003: 房间模块 - 权限不足 (非房主操作)
200004004: 房间模块 - 资源不存在 (房间不存在)
200005005: 房间模块 - 状态冲突 (房间已满员)
200006006: 房间模块 - 业务规则违反 (房间名重复)
```

#### WebSocket模块错误 (400)
```
400001001: WebSocket模块 - 未知错误
400002002: WebSocket模块 - 参数错误 (无效的消息格式)
400003003: WebSocket模块 - 权限不足 (未授权连接)
400004004: WebSocket模块 - 资源不存在 (连接不存在)
400008008: WebSocket模块 - 网络错误 (连接断开)
400009009: WebSocket模块 - 超时错误 (消息发送超时)
```

### 错误处理标准

#### 后端错误处理
```go
// 统一错误响应结构
type ErrorResponse struct {
    Code    string `json:"code"`    // 错误码
    Message string `json:"message"` // 错误信息
    Detail  string `json:"detail"`  // 错误详情
}

// 错误创建函数
func NewError(module, business, errorType, message string) error {
    return &AppError{
        Code:    fmt.Sprintf("%s%s%s", module, business, errorType),
        Message: message,
    }
}

// 使用示例
func (s *UserService) CreateUser(req *CreateUserRequest) error {
    // 参数验证
    if len(req.Username) < 3 {
        return NewError("100", "002", "002", "用户名长度不能少于3个字符")
    }
    
    // 业务逻辑
    if s.repo.UserExists(req.Username) {
        return NewError("100", "006", "006", "用户名已存在")
    }
    
    return nil
}
```

#### 前端错误处理
```typescript
// 统一错误处理
interface ApiError {
  code: string
  message: string
  detail: string
}

// 错误处理工具
class ErrorHandler {
  static handle(error: ApiError | Error): string {
    if ('code' in error) {
      // API错误
      return this.getUserFriendlyMessage(error.code)
    }
    // 系统错误
    return '系统异常，请稍后重试'
  }
  
  private static getUserFriendlyMessage(code: string): string {
    const messages: Record<string, string> = {
      '100002002': '用户名格式不正确',
      '100006006': '用户名已存在',
      '200005005': '房间已满，无法加入',
      '400008008': '网络连接异常',
    }
    
    return messages[code] || '操作失败，请重试'
  }
}

// 组件中使用
try {
  await createRoom(params)
} catch (error) {
  const message = ErrorHandler.handle(error as ApiError)
  showToast(message)
}
```

---

## 🎨 代码风格规范

### Go 代码风格

#### 格式化工具
- **gofmt**: 代码格式化
- **goimports**: 导入包管理
- **golangci-lint**: 代码质量检查

#### 命名规范
```go
// 变量命名：驼峰命名，避免缩写
var userName string        // ✓
var usrName string         // ✗

// 常量命名：驼峰命名，全大写用于导出
const MaxRoomMembers = 7   // ✓
const MAX_ROOM_MEMBERS = 7 // ✗

// 函数命名：驼峰命名
func CreateRoom() {}       // ✓
func create_room() {}      // ✗

// 结构体命名：驼峰命名
type UserService struct {} // ✓
type user_service struct{} // ✗
```

#### 导入规范
```go
import (
    // 标准库
    "context"
    "encoding/json"
    "net/http"
    
    // 第三方库
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "gorm.io/gorm"
    
    // 内部包
    "github.com/xiaowo/internal/service"
    "github.com/xiaowo/pkg/types"
)
```

#### 注释规范
```go
// GetUserByID 根据用户ID获取用户信息
// 返回用户信息和可能的错误
func (s *UserService) GetUserByID(ctx context.Context, id string) (*User, error) {
    // 函数实现
}

// User 用户模型
type User struct {
    ID       uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
    Username string    `json:"username" gorm:"uniqueIndex;not null;size:50"`
    Avatar   string    `json:"avatar" gorm:"size:255"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### Vue3 + TypeScript 代码风格

#### 格式化工具
- **ESLint**: 代码质量检查
- **Prettier**: 代码格式化
- **TypeScript**: 类型检查

#### 命名规范
```typescript
// 变量命名：驼峰命名
const userName = ref('')           // ✓
const user_name = ref('')          // ✗

// 组件命名：PascalCase
const UserProfile = defineComponent({})  // ✓
const userProfile = defineComponent({})  // ✗

// Props命名：camelCase（模板中使用kebab-case）
interface Props {
  userName: string    // ✓
  user-name: string   // ✗
}

// 事件命名：kebab-case
const emit = defineEmits<{
  'user-updated': [user: User]
  'user-deleted': [id: string]
}>()

// emits: ['user-updated', 'user-deleted'] // ✓
// emits: ['userUpdated', 'userDeleted']   // ✗
```

#### 组件结构规范
```vue
<template>
  <div class="component-name">
    <!-- 模板内容 -->
  </div>
</template>

<script setup lang="ts">
// 1. 导入依赖
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import type { User } from '@/types'

// 2. Props定义
interface Props {
  user: User
  readonly?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  readonly: false
})

// 3. Emits定义
const emit = defineEmits<{
  'update:user': [user: User]
  'delete': [id: string]
}>()

// 4. 响应式数据
const isEditing = ref(false)
const formData = ref({ ...props.user })

// 5. 计算属性
const displayName = computed(() => 
  props.user.displayName || props.user.username
)

// 6. 方法
const handleUpdate = () => {
  emit('update:user', formData.value)
  isEditing.value = false
}

// 7. 生命周期
onMounted(() => {
  // 组件挂载逻辑
})
</script>

<style scoped>
.component-name {
  /* 组件样式 */
}
</style>
```

#### API调用规范
```typescript
// 服务层封装
import axios from 'axios'
import type { User, CreateUserRequest } from '@/types'

class UserService {
  private readonly baseURL = import.meta.env.VITE_API_BASE_URL
  
  async createUser(data: CreateUserRequest): Promise<User> {
    try {
      const response = await axios.post(`${this.baseURL}/api/v1/users`, data)
      return response.data
    } catch (error) {
      throw this.handleError(error)
    }
  }
  
  private handleError(error: any): Error {
    if (error.response?.data?.code) {
      return new ApiError(error.response.data.code, error.response.data.message)
    }
    return new Error('网络异常，请重试')
  }
}

// 使用示例
const userService = new UserService()
try {
  const user = await userService.createUser({
    username: 'newuser',
    avatar: 'avatar.jpg'
  })
} catch (error) {
  console.error('创建用户失败:', error)
}
```

### CSS 样式规范

#### Tailwind CSS 使用规范
```vue
<template>
  <!-- 基础样式 -->
  <div class="bg-white rounded-lg shadow-md p-6">
    <!-- 响应式设计 -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      <!-- 状态样式 -->
      <button 
        :class="[
          'px-4 py-2 rounded font-medium transition-colors',
          isActive 
            ? 'bg-blue-500 text-white hover:bg-blue-600' 
            : 'bg-gray-200 text-gray-700 hover:bg-gray-300'
        ]"
      >
        {{ buttonText }}
      </button>
    </div>
  </div>
</template>

<!-- 自定义样式补充 -->
<style scoped>
/* 复杂的CSS动画或特殊效果 */
.slide-enter-active,
.slide-leave-active {
  transition: transform 0.3s ease;
}

.slide-enter-from,
.slide-leave-to {
  transform: translateX(-100%);
}
</style>
```

---

## 🧪 测试规范

### Go 测试规范
```go
// 单元测试命名：TestFunctionName
func TestUserService_CreateUser(t *testing.T) {
    // Arrange：准备测试数据
    service := NewUserService(mockRepo)
    request := &CreateUserRequest{
        Username: "testuser",
        Avatar:   "avatar.jpg",
    }
    
    // Act：执行被测试函数
    user, err := service.CreateUser(context.Background(), request)
    
    // Assert：验证结果
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, "testuser", user.Username)
    assert.Equal(t, "avatar.jpg", user.Avatar)
}

// 表格驱动测试
func TestUserService_ValidateUsername(t *testing.T) {
    tests := []struct {
        name     string
        username string
        wantErr  bool
    }}{
        {"valid username", "testuser", false},
        {"too short", "ab", true},
        {"too long", strings.Repeat("a", 51), true},
        {"contains special chars", "user@123", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            service := NewUserService(nil)
            err := service.ValidateUsername(tt.username)
            
            if tt.wantErr && err == nil {
                t.Errorf("Expected error for username %s", tt.username)
            }
            if !tt.wantErr && err != nil {
                t.Errorf("Unexpected error for username %s: %v", tt.username, err)
            }
        })
    }
}
```

### Vue3 测试规范
```typescript
import { describe, it, expect, vi } from 'vitest'
import { mount } from '@vue/test-utils'
import UserProfile from '@/components/UserProfile.vue'

describe('UserProfile', () => {
  it('renders user information correctly', () => {
    const user = {
      id: '1',
      username: 'testuser',
      avatar: 'avatar.jpg'
    }
    
    const wrapper = mount(UserProfile, {
      props: { user }
    })
    
    expect(wrapper.find('[data-testid="username"]').text()).toBe('testuser')
    expect(wrapper.find('img').attributes('src')).toBe('avatar.jpg')
  })
  
  it('emits update event when save button is clicked', async () => {
    const user = { id: '1', username: 'testuser' }
    const wrapper = mount(UserProfile, {
      props: { user }
    })
    
    await wrapper.find('[data-testid="save-button"]').trigger('click')
    
    expect(wrapper.emitted('update:user')).toBeTruthy()
  })
})
```

---

## 📋 CI/CD 规范

### GitHub Actions 配置
```yaml
# .github/workflows/ci.yml
name: CI/CD Pipeline

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main, develop]

jobs:
  # 后端测试和构建
  backend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Go
        uses: actions/setup-go@v3
        with:
          go-version: '1.21'
          
      - name: Install dependencies
        run: go mod tidy
        
      - name: Run tests
        run: go test -v -coverprofile=coverage.out ./...
        
      - name: Run linter
        run: golangci-lint run
        
      - name: Build application
        run: go build -o bin/xiaowo-server cmd/server/main.go

  # 前端测试和构建
  frontend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '18'
          cache: 'npm'
          
      - name: Install dependencies
        run: npm ci
        
      - name: Run type check
        run: npm run type-check
        
      - name: Run tests
        run: npm run test:coverage
        
      - name: Run linter
        run: npm run lint
        
      - name: Build application
        run: npm run build
```

---

## 🎯 代码质量指标

### 覆盖率要求
- **单元测试覆盖率**: ≥ 80%
- **集成测试覆盖率**: ≥ 60%
- **关键路径覆盖率**: 100%

### 性能要求
- **API响应时间**: < 200ms (95%分位数)
- **页面加载时间**: < 2s (首次访问)
- **WebSocket消息延迟**: < 100ms

### 安全要求
- **SQL注入防护**: 使用参数化查询
- **XSS防护**: 输入输出过滤
- **CSRF防护**: Token验证
- **权限验证**: 所有敏感操作权限检查

---

## 📞 联系方式

- **技术问题**: @老架 (架构相关)
- **代码审查**: @后盾 (后端), @阿码 (前端)
- **环境配置**: @稳当 (基础设施)

**文档更新**: 遇到问题或有改进建议，请及时更新文档，保持规范与实践同步。