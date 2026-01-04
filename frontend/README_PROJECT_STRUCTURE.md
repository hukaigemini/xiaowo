# 小窝项目 - Vue3前端目录结构

## 📁 项目结构树

```
xiaowo-frontend/
├── public/                       # 静态资源
│   ├── favicon.ico              # 网站图标
│   ├── robots.txt               # 搜索引擎爬虫配置
│   └── index.html               # HTML模板
│
├── src/                         # 源代码
│   ├── assets/                  # 资源文件
│   │   ├── images/              # 图片资源
│   │   ├── icons/               # 图标资源
│   │   └── styles/              # 样式文件
│   │       ├── base.css         # 基础样式
│   │       ├── main.css         # 主样式文件
│   │       └── variables.css    # CSS变量
│   │
│   ├── components/              # 公共组件
│   │   ├── common/              # 通用组件
│   │   │   ├── AppHeader.vue    # 应用头部
│   │   │   ├── AppFooter.vue    # 应用底部
│   │   │   ├── LoadingSpinner.vue  # 加载动画
│   │   │   └── ErrorMessage.vue   # 错误提示
│   │   │
│   │   ├── chat/                # 聊天相关组件
│   │   │   ├── MessageList.vue  # 消息列表
│   │   │   ├── MessageInput.vue # 消息输入框
│   │   │   └── MemberList.vue   # 成员列表
│   │   │
│   │   └── video/               # 视频相关组件
│   │       ├── VideoPlayer.vue  # 视频播放器
│   │       ├── VideoControls.vue # 播放控制
│   │       ├── VideoProgress.vue # 播放进度
│   │       └── URLInput.vue     # URL输入框
│   │
│   ├── views/                   # 页面组件
│   │   ├── HomePage.vue         # 首页
│   │   ├── RoomPage.vue         # 房间页面
│   │   ├── NotFound.vue         # 404页面
│   │   └── components/          # 页面私有组件
│   │       └── RoomHeader.vue   # 房间头部
│   │
│   ├── router/                  # 路由配置
│   │   ├── index.ts             # 路由主文件
│   │   └── guards.ts            # 路由守卫
│   │
│   ├── store/                   # 状态管理
│   │   ├── index.ts             # Store主文件
│   │   ├── modules/             # 模块化状态
│   │   │   ├── user.ts          # 用户状态
│   │   │   ├── room.ts          # 房间状态
│   │   │   ├── websocket.ts     # WebSocket状态
│   │   │   └── video.ts         # 视频状态
│   │   └── types/               # Store类型定义
│   │       └── index.ts         # Store接口定义
│   │
│   ├── services/                # API服务层
│   │   ├── api/                 # API接口
│   │   │   ├── user.ts          # 用户API
│   │   │   ├── room.ts          # 房间API
│   │   │   └── message.ts       # 消息API
│   │   ├── websocket.ts         # WebSocket服务
│   │   └── types.ts             # API类型定义
│   │
│   ├── composables/             # Composition API复用逻辑
│   │   ├── useWebSocket.ts      # WebSocket Hook
│   │   ├── useVideoPlayer.ts    # 视频播放器 Hook
│   │   ├── useRoom.ts           # 房间管理 Hook
│   │   └── useResponsive.ts     # 响应式设计 Hook
│   │
│   ├── utils/                   # 工具函数
│   │   ├── api.ts               # API请求封装
│   │   ├── validation.ts        # 表单验证
│   │   ├── format.ts            # 格式化工具
│   │   └── constants.ts         # 常量定义
│   │
│   ├── types/                   # TypeScript类型定义
│   │   ├── api.ts               # API相关类型
│   │   ├── components.ts        # 组件类型
│   │   └── global.ts            # 全局类型
│   │
│   ├── App.vue                  # 根组件
│   └── main.ts                  # 应用入口
│
├── .env.example                 # 环境变量示例
├── .env.development             # 开发环境配置
├── .env.production              # 生产环境配置
├── .env.local                   # 本地环境配置（Git忽略）
│
├── vite.config.ts               # Vite配置文件
├── tsconfig.json                # TypeScript配置
├── tailwind.config.js           # Tailwind CSS配置
├── postcss.config.js            # PostCSS配置
├── eslint.config.js             # ESLint配置
├── prettier.config.js           # Prettier配置
│
├── index.html                   # HTML入口文件
├── package.json                 # 项目依赖配置
├── package-lock.json            # 依赖锁定文件
├── yarn.lock                    # Yarn锁定文件
├── pnpm-lock.yaml               # PNPM锁定文件
├── .gitignore                   # Git忽略文件
├── .editorconfig                # 编辑器配置
└── README.md                    # 项目说明文档
```

## 📋 目录说明

### src/components/ - 组件层
- **common/**: 通用组件，可复用的基础组件
  - `AppHeader.vue`: 导航栏组件
  - `LoadingSpinner.vue`: 加载动画组件
  - `ErrorMessage.vue`: 错误提示组件

- **chat/**: 聊天相关组件
  - `MessageList.vue`: 消息列表显示
  - `MessageInput.vue`: 消息输入和发送
  - `MemberList.vue`: 房间成员列表

- **video/**: 视频相关组件
  - `VideoPlayer.vue`: HTML5视频播放器包装
  - `VideoControls.vue`: 播放控制按钮
  - `URLInput.vue`: 视频URL输入框

### src/views/ - 页面层
- **HomePage.vue**: 首页，创建/加入房间入口
- **RoomPage.vue**: 房间页面，主功能界面
- **components/**: 页面私有组件，只在该页面使用

### src/store/ - 状态管理
- **modules/**: 模块化状态管理
  - `user.ts`: 用户信息和会话状态
  - `room.ts`: 房间信息和成员状态
  - `websocket.ts`: WebSocket连接状态
  - `video.ts`: 视频播放状态

### src/services/ - 服务层
- **api/**: RESTful API接口
- **websocket.ts**: WebSocket服务封装
- **types.ts**: API请求/响应类型定义

### src/composables/ - 复用逻辑
- `useWebSocket.ts`: WebSocket连接管理
- `useVideoPlayer.ts`: 视频播放器控制
- `useRoom.ts`: 房间操作逻辑
- `useResponsive.ts`: 响应式设计逻辑

### src/utils/ - 工具函数
- `api.ts`: Axios拦截器和请求封装
- `validation.ts`: 表单验证规则
- `format.ts`: 数据格式化工具
- `constants.ts`: 应用常量定义

## 🎯 核心设计原则

### 1. 组件化设计
```
View → Components → Composables → Services
   ↓         ↓           ↓          ↓
 路由     组件库      复用逻辑    API服务
```

### 2. 状态管理规范
- **本地状态**: 使用 `ref/reactive`
- **跨组件状态**: 使用 `provide/inject`
- **全局状态**: 使用 Pinia Store
- **持久化状态**: 使用 localStorage

### 3. 响应式设计
- **移动端优先**: 从小屏幕设计开始
- **Tailwind CSS**: 原子化CSS框架
- **断点管理**: 使用 `useResponsive` Hook

## 🚀 快速开始

### 环境要求
```bash
Node.js >= 16.0.0
npm >= 8.0.0 或 yarn >= 1.22.0 或 pnpm >= 7.0.0
```

### 安装依赖
```bash
# 使用 npm
npm install

# 或使用 yarn
yarn install

# 或使用 pnpm
pnpm install
```

### 开发环境
```bash
# 启动开发服务器
npm run dev

# 类型检查
npm run type-check

# 代码检查
npm run lint

# 代码格式化
npm run format
```

### 构建部署
```bash
# 构建生产版本
npm run build

# 预览构建结果
npm run preview

# 部署到静态托管
npm run deploy
```

### 项目配置

#### Vite 配置要点
```typescript
// vite.config.ts
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'

export default defineConfig({
  plugins: [vue()],
  resolve: {
    alias: {
      '@': resolve(__dirname, 'src'),
      '@components': resolve(__dirname, 'src/components'),
      '@views': resolve(__dirname, 'src/views'),
      '@utils': resolve(__dirname, 'src/utils'),
    },
  },
  server: {
    port: 3000,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true,
      },
      '/ws': {
        target: 'ws://localhost:8080',
        ws: true,
      },
    },
  },
})
```

#### 环境变量配置
```env
# .env.development
VITE_API_BASE_URL=http://localhost:8080
VITE_WS_URL=ws://localhost:8080
VITE_APP_TITLE=小窝同步观影
VITE_APP_VERSION=1.0.0
```

这个目录结构遵循Vue3 + TypeScript最佳实践，确保代码类型安全、组件复用性强、维护性高。