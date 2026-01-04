# 开发设计指导文档

## 📋 文档信息

- **文档版本**: v1.0
- **创建日期**: 2025-12-31
- **创建者**: 优优 (UI/UX设计师)
- **目标受众**: 前端开发工程师
- **适用范围**: 小窝项目设计系统实施指导
- **技术栈**: Vue 3 + TypeScript + Tailwind CSS

## 🎯 开发目标与原则

### 开发目标
1. **快速开发**: 基于设计系统的组件化开发
2. **一致体验**: 确保跨页面、跨设备的一致性
3. **性能优先**: 优化加载速度和交互响应
4. **可维护性**: 清晰的代码结构和文档

### 开发原则
- **移动优先**: 从移动端开始设计和开发
- **渐进增强**: 在基础功能上逐步添加特性
- **组件化**: 构建可复用的UI组件
- **类型安全**: 充分利用TypeScript的类型系统

## 🛠️ 技术栈配置

### Tailwind CSS 配置
```javascript
// tailwind.config.js
module.exports = {
  content: [
    "./index.html",
    "./src/**/*.{vue,js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      // 自定义颜色
      colors: {
        primary: {
          50: '#eff6ff',
          100: '#dbeafe',
          200: '#bfdbfe',
          300: '#93c5fd',
          400: '#60a5fa',
          500: '#3b82f6',
          600: '#2563eb',
          700: '#1d4ed8',
          800: '#1e40af',
          900: '#1e3a8a',
        },
        success: {
          50: '#ecfdf5',
          500: '#10b981',
          600: '#059669',
          700: '#047857',
        },
        warning: {
          50: '#fffbeb',
          500: '#f59e0b',
          600: '#d97706',
          700: '#b45309',
        },
        error: {
          50: '#fef2f2',
          500: '#ef4444',
          600: '#dc2626',
          700: '#b91c1c',
        }
      },
      
      // 自定义间距
      spacing: {
        '18': '4.5rem',
        '88': '22rem',
      },
      
      // 自定义字体大小
      fontSize: {
        '2xs': ['0.625rem', { lineHeight: '1rem' }],
      },
      
      // 动画
      animation: {
        'fade-in': 'fadeIn 0.5s ease-in-out',
        'slide-up': 'slideUp 0.3s ease-out',
        'scale-in': 'scaleIn 0.2s ease-out',
      },
      
      // 关键帧
      keyframes: {
        fadeIn: {
          '0%': { opacity: '0' },
          '100%': { opacity: '1' },
        },
        slideUp: {
          '0%': { transform: 'translateY(10px)', opacity: '0' },
          '100%': { transform: 'translateY(0)', opacity: '1' },
        },
        scaleIn: {
          '0%': { transform: 'scale(0.95)', opacity: '0' },
          '100%': { transform: 'scale(1)', opacity: '1' },
        },
      }
    },
  },
  plugins: [
    require('@tailwindcss/forms'),
    require('@tailwindcss/typography'),
    require('@tailwindcss/aspect-ratio'),
  ],
}
```

### PostCSS 配置
```javascript
// postcss.config.js
module.exports = {
  plugins: {
    tailwindcss: {},
    autoprefixer: {},
    // CSS 优化插件
    'postcss-preset-env': {
      stage: 3,
      features: {
        'nesting-rules': true,
      },
    },
  },
}
```

### CSS 基础配置
```css
/* src/assets/styles/main.css */
@tailwind base;
@tailwind components;
@tailwind utilities;

/* 自定义基础样式 */
@layer base {
  html {
    scroll-behavior: smooth;
    -webkit-text-size-adjust: 100%;
  }
  
  body {
    @apply text-gray-900 antialiased;
    font-feature-settings: 'rlig' 1, 'calt' 1;
  }
  
  /* 滚动条样式 */
  ::-webkit-scrollbar {
    width: 8px;
    height: 8px;
  }
  
  ::-webkit-scrollbar-track {
    @apply bg-gray-100;
  }
  
  ::-webkit-scrollbar-thumb {
    @apply bg-gray-300 rounded-full;
  }
  
  ::-webkit-scrollbar-thumb:hover {
    @apply bg-gray-400;
  }
  
  /* 选中文本样式 */
  ::selection {
    @apply bg-primary-100 text-primary-900;
  }
  
  /* 焦点样式 */
  :focus {
    outline: 2px solid theme('colors.primary.500');
    outline-offset: 2px;
  }
  
  :focus:not(:focus-visible) {
    outline: none;
  }
}

/* 组件样式 */
@layer components {
  /* 按钮基础样式 */
  .btn {
    @apply inline-flex items-center justify-center font-medium rounded-lg transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2;
    @apply disabled:opacity-50 disabled:cursor-not-allowed;
    min-height: 44px; /* 触摸友好 */
  }
  
  .btn-primary {
    @apply btn bg-primary-500 text-white hover:bg-primary-600 focus:ring-primary-500;
  }
  
  .btn-secondary {
    @apply btn border border-gray-300 bg-white text-gray-700 hover:bg-gray-50 focus:ring-primary-500;
  }
  
  .btn-ghost {
    @apply btn text-gray-700 hover:bg-gray-100 focus:ring-primary-500;
  }
  
  /* 卡片样式 */
  .card {
    @apply bg-white rounded-xl shadow-sm border border-gray-100;
  }
  
  .card-body {
    @apply p-6;
  }
  
  /* 输入框样式 */
  .input {
    @apply block w-full border border-gray-300 rounded-lg px-4 py-3 text-base;
    @apply focus:ring-2 focus:ring-primary-500 focus:border-primary-500;
    @apply disabled:bg-gray-50 disabled:text-gray-500 disabled:cursor-not-allowed;
    min-height: 48px;
  }
  
  /* 文本样式 */
  .text-heading {
    @apply text-gray-900 font-semibold;
  }
  
  .text-body {
    @apply text-gray-600 leading-relaxed;
  }
  
  .text-caption {
    @apply text-sm text-gray-500;
  }
}

/* 工具样式 */
@layer utilities {
  /* 安全区域支持 */
  .safe-top {
    padding-top: env(safe-area-inset-top);
  }
  
  .safe-bottom {
    padding-bottom: env(safe-area-inset-bottom);
  }
  
  .safe-left {
    padding-left: env(safe-area-inset-left);
  }
  
  .safe-right {
    padding-right: env(safe-area-inset-right);
  }
  
  /* 触摸优化 */
  .touch-manipulation {
    touch-action: manipulation;
  }
  
  .no-tap-highlight {
    -webkit-tap-highlight-color: transparent;
  }
  
  /* 文本处理 */
  .text-balance {
    text-wrap: balance;
  }
  
  /* 硬件加速 */
  .gpu-accelerated {
    transform: translateZ(0);
    will-change: transform;
  }
  
  /* 隐藏滚动条 */
  .hide-scrollbar {
    -ms-overflow-style: none;
    scrollbar-width: none;
  }
  
  .hide-scrollbar::-webkit-scrollbar {
    display: none;
  }
}
```

## 🎨 组件开发指南

### 组件结构规范
```vue
<!-- 组件模板结构 -->
<template>
  <div :class="wrapperClasses">
    <!-- 组件内容 -->
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'

// Props 类型定义
interface Props {
  variant?: 'primary' | 'secondary' | 'ghost'
  size?: 'sm' | 'md' | 'lg'
  disabled?: boolean
  loading?: boolean
  fullWidth?: boolean
}

// 默认值
const props = withDefaults(defineProps<Props>(), {
  variant: 'primary',
  size: 'md',
  disabled: false,
  loading: false,
  fullWidth: false
})

// Emits 定义
const emit = defineEmits<{
  click: [event: MouseEvent]
  focus: [event: FocusEvent]
  blur: [event: FocusEvent]
}>()

// 响应式数据
const isProcessing = ref(false)

// 计算属性
const wrapperClasses = computed(() => [
  'base-component-class',
  {
    'variant-primary': props.variant === 'primary',
    'variant-secondary': props.variant === 'secondary',
    'size-sm': props.size === 'sm',
    'size-md': props.size === 'md',
    'size-lg': props.size === 'lg',
    'disabled': props.disabled,
    'loading': props.loading,
    'full-width': props.fullWidth
  }
])

// 方法
function handleClick(event: MouseEvent) {
  if (props.disabled || props.loading) return
  emit('click', event)
}

function handleFocus(event: FocusEvent) {
  emit('focus', event)
}

function handleBlur(event: FocusEvent) {
  emit('blur', event)
}
</script>

<style scoped lang="scss">
/* 组件样式 */
.base-component-class {
  /* 基础样式 */
}

.variant-primary {
  /* 主要样式变体 */
}

.size-sm {
  /* 小尺寸样式 */
}

/* 响应式样式 */
@media (max-width: 640px) {
  .base-component-class {
    /* 移动端样式 */
  }
}
</style>
```

### 基础组件实现

#### Button 组件
```vue
<!-- components/ui/Button.vue -->
<template>
  <component
    :is="tag"
    :type="tag === 'button' ? type : undefined"
    :disabled="disabled || loading"
    :class="buttonClasses"
    @click="handleClick"
    @focus="handleFocus"
    @blur="handleBlur"
  >
    <!-- 加载指示器 -->
    <svg
      v-if="loading"
      class="animate-spin -ml-1 mr-3 h-4 w-4"
      :class="loadingIconColor"
      fill="none"
      viewBox="0 0 24 24"
    >
      <circle
        class="opacity-25"
        cx="12"
        cy="12"
        r="10"
        stroke="currentColor"
        stroke-width="4"
      />
      <path
        class="opacity-75"
        fill="currentColor"
        d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
      />
    </svg>
    
    <!-- 图标插槽 -->
    <slot v-if="!loading" name="icon-left" />
    
    <!-- 文字内容 -->
    <span v-if="$slots.default || loading" :class="{ 'sr-only': loading }">
      <slot>{{ text }}</slot>
    </span>
    
    <!-- 右侧图标 -->
    <slot v-if="!loading" name="icon-right" />
  </component>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Props {
  variant?: 'primary' | 'secondary' | 'ghost' | 'danger'
  size?: 'sm' | 'md' | 'lg'
  type?: 'button' | 'submit' | 'reset'
  disabled?: boolean
  loading?: boolean
  fullWidth?: boolean
  text?: string
  tag?: 'button' | 'a' | 'router-link'
  href?: string
  to?: string
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'primary',
  size: 'md',
  type: 'button',
  disabled: false,
  loading: false,
  fullWidth: false,
  tag: 'button'
})

const emit = defineEmits<{
  click: [event: MouseEvent]
  focus: [event: FocusEvent]
  blur: [event: FocusEvent]
}>()

const buttonClasses = computed(() => [
  'inline-flex items-center justify-center font-medium rounded-lg transition-all duration-200',
  'focus:outline-none focus:ring-2 focus:ring-offset-2',
  'disabled:opacity-50 disabled:cursor-not-allowed disabled:pointer-events-none',
  {
    // 尺寸变体
    'min-w-[44px] min-h-[44px] px-3 py-2 text-sm': props.size === 'sm',
    'min-w-[48px] min-h-[48px] px-6 py-3 text-base': props.size === 'md',
    'min-w-[52px] min-h-[52px] px-8 py-4 text-lg': props.size === 'lg',
    
    // 颜色变体
    'bg-primary-500 text-white hover:bg-primary-600 focus:ring-primary-500 shadow-sm hover:shadow-md':
      props.variant === 'primary',
    'border border-gray-300 bg-white text-gray-700 hover:bg-gray-50 focus:ring-primary-500 shadow-sm':
      props.variant === 'secondary',
    'text-gray-700 hover:bg-gray-100 focus:ring-primary-500':
      props.variant === 'ghost',
    'bg-red-500 text-white hover:bg-red-600 focus:ring-red-500 shadow-sm':
      props.variant === 'danger',
    
    // 布局变体
    'w-full': props.fullWidth,
    
    // 状态变体
    'loading': props.loading
  }
])

const loadingIconColor = computed(() => {
  switch (props.variant) {
    case 'primary':
    case 'danger':
      return 'text-white'
    case 'secondary':
    case 'ghost':
      return 'text-gray-500'
    default:
      return 'text-current'
  }
})

function handleClick(event: MouseEvent) {
  if (props.disabled || props.loading) {
    event.preventDefault()
    return
  }
  emit('click', event)
}

function handleFocus(event: FocusEvent) {
  emit('focus', event)
}

function handleBlur(event: FocusEvent) {
  emit('blur', event)
}
</script>

<!-- 使用示例 -->
<template>
  <!-- 主要按钮 -->
  <Button variant="primary" @click="handleSubmit">
    提交表单
  </Button>
  
  <!-- 带图标的按钮 -->
  <Button variant="secondary">
    <template #icon-left>
      <svg class="w-4 h-4"><!-- 图标 --></svg>
    </template>
    编辑
  </Button>
  
  <!-- 加载状态按钮 -->
  <Button :loading="isLoading" variant="primary">
    保存
  </Button>
  
  <!-- 全宽移动端按钮 -->
  <Button variant="primary" fullWidth size="lg" class="sm:hidden">
    移动端按钮
  </Button>
</template>
```

#### Input 组件
```vue
<!-- components/ui/Input.vue -->
<template>
  <div class="w-full">
    <!-- 标签 -->
    <label
      v-if="label"
      :for="inputId"
      class="block text-sm font-medium text-gray-700 mb-2"
    >
      {{ label }}
      <span v-if="required" class="text-red-500">*</span>
    </label>
    
    <!-- 输入框容器 -->
    <div class="relative">
      <!-- 左侧图标 -->
      <div
        v-if="$slots.icon"
        class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none"
      >
        <slot name="icon" />
      </div>
      
      <!-- 输入框 -->
      <input
        :id="inputId"
        ref="inputRef"
        v-model="inputValue"
        :type="type"
        :placeholder="placeholder"
        :disabled="disabled"
        :readonly="readonly"
        :required="required"
        :min="min"
        :max="max"
        :step="step"
        :pattern="pattern"
        :class="inputClasses"
        :aria-describedby="helpTextId"
        :aria-invalid="!!error"
        @input="handleInput"
        @focus="handleFocus"
        @blur="handleBlur"
        @keydown="handleKeydown"
      />
      
      <!-- 右侧操作 -->
      <div class="absolute inset-y-0 right-0 flex items-center pr-3">
        <!-- 清除按钮 -->
        <button
          v-if="clearable && inputValue && !disabled"
          type="button"
          class="text-gray-400 hover:text-gray-600 focus:outline-none focus:ring-2 focus:ring-primary-500 rounded-full p-1"
          @click="clearInput"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
        
        <!-- 错误图标 -->
        <svg
          v-if="error"
          class="w-5 h-5 text-red-400"
          fill="currentColor"
          viewBox="0 0 20 20"
        >
          <path
            fill-rule="evenodd"
            d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z"
            clip-rule="evenodd"
          />
        </svg>
      </div>
    </div>
    
    <!-- 帮助文本和错误信息 -->
    <p
      v-if="error"
      :id="helpTextId"
      class="mt-2 text-sm text-red-600"
    >
      {{ error }}
    </p>
    <p
      v-else-if="helpText"
      :id="helpTextId"
      class="mt-2 text-sm text-gray-500"
    >
      {{ helpText }}
    </p>
  </div>
</template>

<script setup lang="ts">
import { computed, ref, nextTick } from 'vue'

interface Props {
  modelValue?: string | number
  label?: string
  type?: 'text' | 'email' | 'password' | 'number' | 'tel' | 'url' | 'search'
  placeholder?: string
  disabled?: boolean
  readonly?: boolean
  required?: boolean
  clearable?: boolean
  error?: string
  helpText?: string
  size?: 'sm' | 'md' | 'lg'
  min?: number | string
  max?: number | string
  step?: number | string
  pattern?: string
}

const props = withDefaults(defineProps<Props>(), {
  type: 'text',
  disabled: false,
  readonly: false,
  required: false,
  clearable: false,
  size: 'md'
})

const emit = defineEmits<{
  'update:modelValue': [value: string | number]
  focus: [event: FocusEvent]
  blur: [event: FocusEvent]
  keydown: [event: KeyboardEvent]
  clear: []
}>()

const inputValue = ref(props.modelValue)
const inputId = `input-${Math.random().toString(36).substr(2, 9)}`
const inputRef = ref<HTMLInputElement>()

const inputClasses = computed(() => [
  'block w-full border rounded-lg transition-colors focus:outline-none focus:ring-2 focus:ring-offset-0',
  'disabled:bg-gray-50 disabled:text-gray-500 disabled:cursor-not-allowed',
  'placeholder-gray-400',
  {
    // 尺寸变体
    'px-3 py-2 text-sm': props.size === 'sm',
    'px-4 py-3 text-base': props.size === 'md',
    'px-6 py-4 text-lg': props.size === 'lg',
    
    // 状态变体
    'pl-10': !!useSlots().icon,
    'pr-10': props.clearable || !!props.error,
    'border-gray-300 focus:border-primary-500 focus:ring-primary-500': !props.error,
    'border-red-300 focus:border-red-500 focus:ring-red-500': props.error
  }
])

const helpTextId = computed(() => `${inputId}-help`)

function handleInput(event: Event) {
  const target = event.target as HTMLInputElement
  const value = props.type === 'number' ? Number(target.value) : target.value
  inputValue.value = value
  emit('update:modelValue', value)
}

function handleFocus(event: FocusEvent) {
  emit('focus', event)
}

function handleBlur(event: FocusEvent) {
  emit('blur', event)
}

function handleKeydown(event: KeyboardEvent) {
  emit('keydown', event)
}

function clearInput() {
  inputValue.value = ''
  emit('update:modelValue', '')
  emit('clear')
  nextTick(() => {
    inputRef.value?.focus()
  })
}

// 监听外部值变化
watch(() => props.modelValue, (newValue) => {
  inputValue.value = newValue
})
</script>

<!-- 使用示例 -->
<template>
  <!-- 基本输入框 -->
  <Input
    v-model="email"
    label="邮箱地址"
    type="email"
    placeholder="请输入邮箱"
    required
  />
  
  <!-- 带图标的输入框 -->
  <Input
    v-model="searchQuery"
    placeholder="搜索..."
    size="lg"
  >
    <template #icon>
      <svg class="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
      </svg>
    </template>
  </Input>
  
  <!-- 错误状态 -->
  <Input
    v-model="username"
    label="用户名"
    :error="usernameError"
    help-text="用户名只能包含字母和数字"
  />
</template>
```

## 📱 响应式开发指南

### 移动优先 CSS
```css
/* 移动端基础样式 (默认) */
.component {
  /* 移动端样式 */
  padding: 16px;
  font-size: 14px;
  flex-direction: column;
}

/* 平板及以上设备 */
@media (min-width: 768px) {
  .component {
    padding: 24px;
    font-size: 16px;
    flex-direction: row;
  }
}

/* 桌面设备 */
@media (min-width: 1024px) {
  .component {
    padding: 32px;
    font-size: 18px;
  }
}
```

### 响应式组件示例
```vue
<template>
  <div class="responsive-component">
    <!-- 移动端：单列布局 -->
    <div class="mobile-layout">
      <div class="content-block">
        <!-- 主内容 -->
      </div>
      <div class="action-block">
        <!-- 操作按钮 -->
      </div>
    </div>
    
    <!-- 桌面端：双列布局 -->
    <div class="desktop-layout hidden md:grid">
      <div class="content-area">
        <!-- 主内容 -->
      </div>
      <div class="sidebar-area">
        <!-- 侧边栏 -->
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
// 组件逻辑
</script>

<style scoped lang="scss">
.responsive-component {
  @apply w-full;
}

.mobile-layout {
  @apply flex flex-col gap-4 p-4;
  
  .content-block {
    @apply flex-1;
  }
  
  .action-block {
    @apply w-full;
    
    button {
      @apply w-full min-h-[48px];
    }
  }
}

.desktop-layout {
  @apply grid-cols-3 gap-6 p-6;
  
  .content-area {
    @apply col-span-2;
  }
  
  .sidebar-area {
    @apply col-span-1;
  }
}

/* 平板断点 */
@media (min-width: 768px) {
  .mobile-layout {
    @apply hidden;
  }
  
  .desktop-layout {
    @apply grid;
  }
}
</style>
```

### 触摸友好设计
```vue
<template>
  <div class="touch-friendly-component">
    <!-- 触摸目标最小44x44px -->
    <button class="touch-target">
      <svg class="touch-icon"><!-- 图标 --></svg>
      <span class="touch-text">按钮文字</span>
    </button>
    
    <!-- 触摸间距 -->
    <div class="touch-spacing">
      <!-- 内容 -->
    </div>
  </div>
</template>

<style scoped lang="scss">
.touch-friendly-component {
  @apply touch-manipulation no-tap-highlight;
}

.touch-target {
  @apply flex items-center gap-3 min-w-[44px] min-h-[44px] px-4 py-3;
  @apply bg-white border border-gray-300 rounded-lg;
  @apply hover:bg-gray-50 active:scale-95 transition-all duration-150;
  @apply focus:outline-none focus:ring-2 focus:ring-primary-500;
}

.touch-icon {
  @apply w-5 h-5 flex-shrink-0;
}

.touch-text {
  @apply text-base font-medium;
}

.touch-spacing {
  @apply mt-4 p-4; // 最小8px间距
}

/* 防止用户选择文本 */
.touch-target {
  -webkit-user-select: none;
  -moz-user-select: none;
  -ms-user-select: none;
  user-select: none;
}
</style>
```

## 🚀 性能优化实践

### 图片优化
```vue
<!-- 响应式图片组件 -->
<template>
  <picture>
    <!-- 移动端图片 -->
    <source
      v-if="sources.mobile"
      media="(max-width: 640px)"
      :srcset="sources.mobile"
      :type="getImageType(sources.mobile)"
    />
    
    <!-- 平板图片 -->
    <source
      v-if="sources.tablet"
      media="(min-width: 641px) and (max-width: 1024px)"
      :srcset="sources.tablet"
      :type="getImageType(sources.tablet)"
    />
    
    <!-- 桌面图片 -->
    <source
      v-if="sources.desktop"
      media="(min-width: 1025px)"
      :srcset="sources.desktop"
      :type="getImageType(sources.desktop)"
    />
    
    <!-- 默认图片 -->
    <img
      :src="fallbackSrc"
      :alt="alt"
      :class="imageClasses"
      :loading="lazy ? 'lazy' : 'eager'"
      :decoding="lazy ? 'async' : 'sync'"
      @error="handleImageError"
    />
  </picture>
</template>

<script setup lang="ts">
interface ImageSources {
  mobile?: string
  tablet?: string
  desktop?: string
}

interface Props {
  sources: ImageSources
  fallbackSrc: string
  alt: string
  lazy?: boolean
  class?: string
  width?: number | string
  height?: number | string
}

const props = withDefaults(defineProps<Props>(), {
  lazy: true,
  class: ''
})

const emit = defineEmits<{
  error: [event: Event]
}>()

const imageClasses = computed(() => [
  'w-full h-auto object-cover',
  props.class
])

function getImageType(src: string): string {
  const ext = src.split('.').pop()?.toLowerCase()
  switch (ext) {
    case 'webp':
      return 'image/webp'
    case 'avif':
      return 'image/avif'
    case 'jpg':
    case 'jpeg':
      return 'image/jpeg'
    case 'png':
      return 'image/png'
    default:
      return 'image/jpeg'
  }
}

function handleImageError(event: Event) {
  emit('error', event)
}
</script>
```

### CSS 优化
```css
/* 硬件加速 */
.gpu-layer {
  transform: translateZ(0);
  will-change: transform;
}

/* 避免重排 */
.avoid-reflow {
  will-change: auto;
  backface-visibility: hidden;
}

/* 优化动画性能 */
.smooth-animation {
  animation-fill-mode: both;
  animation-duration: 0.3s;
  animation-timing-function: cubic-bezier(0.4, 0, 0.2, 1);
}

/* 图片懒加载优化 */
.lazy-image {
  opacity: 0;
  transition: opacity 0.3s;
}

.lazy-image.loaded {
  opacity: 1;
}

/* 减少重绘 */
.reduce-paints {
  contain: layout style paint;
}
```

### 组件懒加载
```typescript
// 路由懒加载
const HomePage = () => import('@/views/HomePage.vue')
const RoomPage = () => import('@/views/RoomPage.vue')

// 组件懒加载
const HeavyComponent = defineAsyncComponent({
  loader: () => import('@/components/HeavyComponent.vue'),
  loadingComponent: LoadingSpinner,
  errorComponent: ErrorMessage,
  delay: 200,
  timeout: 3000
})

// 图片懒加载
const LazyImage = {
  setup() {
    const imageRef = ref<HTMLElement>()
    const isLoaded = ref(false)
    
    const observer = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) {
            // 开始加载图片
            isLoaded.value = true
            observer.unobserve(entry.target)
          }
        })
      },
      {
        rootMargin: '50px'
      }
    )
    
    onMounted(() => {
      if (imageRef.value) {
        observer.observe(imageRef.value)
      }
    })
    
    onUnmounted(() => {
      observer.disconnect()
    })
    
    return {
      imageRef,
      isLoaded
    }
  }
}
```

## 🎯 状态管理

### 响应式状态管理
```typescript
// stores/ui.ts
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useUiStore = defineStore('ui', () => {
  // 响应式状态
  const isLoading = ref(false)
  const toast = ref<{
    id: string
    message: string
    type: 'success' | 'error' | 'warning' | 'info'
    duration?: number
  } | null>(null)
  
  const modal = ref<{
    isOpen: boolean
    component?: string
    props?: Record<string, any>
  }>({
    isOpen: false
  })
  
  const sidebar = ref({
    isOpen: false,
    mode: 'over' | 'side'
  })
  
  // 计算属性
  const isToastVisible = computed(() => !!toast.value)
  
  // 方法
  function setLoading(loading: boolean) {
    isLoading.value = loading
  }
  
  function showToast(message: string, type: 'success' | 'error' | 'warning' | 'info' = 'info', duration = 3000) {
    const id = Date.now().toString()
    toast.value = { id, message, type, duration }
    
    if (duration > 0) {
      setTimeout(() => {
        hideToast()
      }, duration)
    }
  }
  
  function hideToast() {
    toast.value = null
  }
  
  function openModal(component: string, props?: Record<string, any>) {
    modal.value = {
      isOpen: true,
      component,
      props
    }
  }
  
  function closeModal() {
    modal.value = {
      isOpen: false
    }
  }
  
  function toggleSidebar() {
    sidebar.value.isOpen = !sidebar.value.isOpen
  }
  
  function setSidebarState(isOpen: boolean) {
    sidebar.value.isOpen = isOpen
  }
  
  return {
    // state
    isLoading,
    toast,
    modal,
    sidebar,
    
    // getters
    isToastVisible,
    
    // actions
    setLoading,
    showToast,
    hideToast,
    openModal,
    closeModal,
    toggleSidebar,
    setSidebarState
  }
})
```

### 状态管理使用示例
```vue
<script setup lang="ts">
import { useUiStore } from '@/stores/ui'

const uiStore = useUiStore()

// 显示加载状态
function handleSubmit() {
  uiStore.setLoading(true)
  
  // 模拟异步操作
  setTimeout(() => {
    uiStore.setLoading(false)
    uiStore.showToast('提交成功！', 'success')
  }, 2000)
}

// 显示错误提示
function handleError() {
  uiStore.showToast('操作失败，请重试', 'error')
}

// 打开模态框
function openEditModal() {
  uiStore.openModal('EditModal', {
    title: '编辑用户',
    userId: '123'
  })
}
</script>

<template>
  <div>
    <!-- 全局加载状态 -->
    <div v-if="uiStore.isLoading" class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50">
      <div class="bg-white rounded-lg p-6 flex items-center gap-3">
        <svg class="animate-spin h-6 w-6 text-primary-500" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
        <span class="text-gray-900">加载中...</span>
      </div>
    </div>
    
    <!-- 全局 Toast 通知 -->
    <Transition
      enter-active-class="transition ease-out duration-300"
      enter-from-class="opacity-0 transform translate-y-2"
      enter-to-class="opacity-100 transform translate-y-0"
      leave-active-class="transition ease-in duration-200"
      leave-from-class="opacity-100 transform translate-y-0"
      leave-to-class="opacity-0 transform translate-y-2"
    >
      <div
        v-if="uiStore.toast"
        class="fixed top-4 right-4 z-50 max-w-sm w-full bg-white border border-gray-200 rounded-lg shadow-lg p-4"
      >
        <div class="flex items-start gap-3">
          <div :class="getToastIconClass(uiStore.toast.type)">
            <component :is="getToastIcon(uiStore.toast.type)" />
          </div>
          <div class="flex-1">
            <p class="text-sm font-medium text-gray-900">{{ uiStore.toast.message }}</p>
          </div>
          <button
            @click="uiStore.hideToast"
            class="text-gray-400 hover:text-gray-600 focus:outline-none"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
      </div>
    </Transition>
    
    <!-- 页面内容 -->
    <button @click="handleSubmit" class="btn-primary">
      提交表单
    </button>
  </div>
</template>

<script setup lang="ts">
function getToastIcon(type: string) {
  const icons = {
    success: 'CheckCircleIcon',
    error: 'XCircleIcon',
    warning: 'ExclamationTriangleIcon',
    info: 'InformationCircleIcon'
  }
  return icons[type as keyof typeof icons] || 'InformationCircleIcon'
}

function getToastIconClass(type: string) {
  const classes = {
    success: 'w-6 h-6 text-green-500',
    error: 'w-6 h-6 text-red-500',
    warning: 'w-6 h-6 text-yellow-500',
    info: 'w-6 h-6 text-blue-500'
  }
  return classes[type as keyof typeof classes] || classes.info
}
</script>
```

## 📋 开发检查清单

### 代码质量检查
- [ ] TypeScript 类型定义完整
- [ ] 组件 Props 和 Emits 明确定义
- [ ] 响应式数据使用 ref/computed
- [ ] 事件处理函数正确绑定
- [ ] 错误边界处理完善

### 响应式设计检查
- [ ] 移动端优先的 CSS 实现
- [ ] 触摸目标最小尺寸 44×44px
- [ ] 断点使用符合设计规范
- [ ] 图片响应式适配
- [ ] 字体大小响应式调整

### 性能优化检查
- [ ] 图片懒加载实现
- [ ] 组件懒加载使用
- [ ] CSS 动画性能优化
- [ ] 避免不必要的重排重绘
- [ ] 资源压缩和缓存

### 可访问性检查
- [ ] 语义化 HTML 标签
- [ ] ARIA 标签正确使用
- [ ] 键盘导航支持
- [ ] 屏幕阅读器兼容
- [ ] 颜色对比度充足

### 状态管理检查
- [ ] 全局状态正确使用 Pinia
- [ ] 本地状态合理管理
- [ ] 异步状态加载处理
- [ ] 错误状态友好提示
- [ ] 加载状态用户反馈

## 🔧 开发工具配置

### VS Code 配置
```json
// .vscode/settings.json
{
  "editor.formatOnSave": true,
  "editor.defaultFormatter": "esbenp.prettier-vscode",
  "editor.codeActionsOnSave": {
    "source.fixAll.eslint": true
  },
  "typescript.preferences.importModuleSpecifier": "relative",
  "emmet.includeLanguages": {
    "vue": "html"
  },
  "files.associations": {
    "*.vue": "vue"
  }
}
```

### ESLint 配置
```javascript
// .eslintrc.js
module.exports = {
  extends: [
    '@vue/typescript/recommended',
    'prettier'
  ],
  rules: {
    // Vue 特定规则
    'vue/multi-word-component-names': 'off',
    'vue/no-v-html': 'warn',
    'vue/require-default-prop': 'off',
    'vue/require-prop-types': 'off',
    
    // TypeScript 规则
    '@typescript-eslint/no-unused-vars': 'error',
    '@typescript-eslint/no-explicit-any': 'warn',
    '@typescript-eslint/explicit-function-return-type': 'off',
    
    // 通用规则
    'no-console': process.env.NODE_ENV === 'production' ? 'error' : 'warn',
    'no-debugger': process.env.NODE_ENV === 'production' ? 'error' : 'warn'
  }
}
```

### Prettier 配置
```javascript
// .prettierrc.js
module.exports = {
  semi: false,
  singleQuote: true,
  tabWidth: 2,
  useTabs: false,
  trailingComma: 'es5',
  printWidth: 100,
  bracketSpacing: true,
  arrowParens: 'avoid',
  endOfLine: 'lf'
}
```

---

**开发指南维护**: 请在开发新组件时参考此文档
**技术更新**: 建议每季度更新技术栈和最佳实践
**最后更新**: 2025-12-31
**版本**: v1.0