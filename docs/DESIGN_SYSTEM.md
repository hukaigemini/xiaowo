# 设计系统文档

## 📋 文档信息

- **文档版本**: v1.0
- **创建日期**: 2025-12-31
- **创建者**: 优优 (UI/UX设计师)
- **适用范围**: 小窝项目完整设计系统
- **基于实践**: RoomPage.vue + 响应式优化经验
- **技术栈**: Vue 3 + TypeScript + Tailwind CSS

## 🎯 设计价值观与原则

### 核心价值观
1. **用户至上**: 以用户需求为中心的设计决策
2. **一致性优先**: 统一的设计语言和交互模式
3. **性能驱动**: 快速、流畅的用户体验
4. **可访问性**: 包容性设计，支持所有用户
5. **简洁高效**: 最小化认知负担，最大化效率

### 设计原则
#### 1. 清晰性 (Clarity)
- 信息层级明确，关键内容突出
- 视觉引导用户注意力
- 减少用户认知负担

#### 2. 一致性 (Consistency)
- 统一的设计语言和交互模式
- 可预测的用户界面行为
- 品牌识别度保持

#### 3. 效率性 (Efficiency)
- 快速完成任务路径
- 减少操作步骤
- 智能默认设置

#### 4. 反馈性 (Feedback)
- 即时的操作反馈
- 清晰的状态提示
- 错误处理友好

## 🎨 颜色系统 (Color System)

### 主品牌色 (Primary Colors)
```css
/* 主品牌色 - 蓝色系 */
--primary-50:  #eff6ff   /* 最浅色 - 背景色 */
--primary-100: #dbeafe   /* 浅色 - 悬停背景 */
--primary-200: #bfdbfe   /* 更浅色 */
--primary-300: #93c5fd   /* 浅色 */
--primary-400: #60a5fa   /* 次要强调色 */
--primary-500: #3b82f6   /* 主要强调色 - 品牌色 */
--primary-600: #2563eb   /* 深色 - 悬停状态 */
--primary-700: #1d4ed8   /* 更深色 - 激活状态 */
--primary-800: #1e40af   /* 深色 - 文字色 */
--primary-900: #1e3a8a   /* 最深色 - 强调文字 */
```

### 功能色彩 (Functional Colors)
```css
/* 成功色 - 绿色系 */
--success-50:  #ecfdf5
--success-500: #10b981
--success-600: #059669
--success-700: #047857

/* 警告色 - 橙色系 */
--warning-50:  #fffbeb
--warning-500: #f59e0b
--warning-600: #d97706
--warning-700: #b45309

/* 错误色 - 红色系 */
--error-50:    #fef2f2
--error-500:   #ef4444
--error-600:   #dc2626
--error-700:   #b91c1c

/* 信息色 - 蓝色系 */
--info-50:     #eff6ff
--info-500:    #3b82f6
--info-600:    #2563eb
--info-700:    #1d4ed8
```

### 中性色 (Neutral Colors)
```css
/* 灰度系统 */
--gray-50:  #f9fafb    /* 背景色 */
--gray-100: #f3f4f6    /* 卡片背景 */
--gray-200: #e5e7eb    /* 边框色 */
--gray-300: #d1d5db    /* 分割线 */
--gray-400: #9ca3af    /* 占位符 */
--gray-500: #6b7280    /* 次要文字 */
--gray-600: #4b5563    /* 弱化文字 */
--gray-700: #374151    /* 主要文字 */
--gray-800: #1f2937    /* 深色文字 */
--gray-900: #111827    /* 最深文字 */
```

### 语义化颜色使用
```vue
<!-- 主要按钮 -->
<button class="bg-primary-500 hover:bg-primary-600 text-white">
  主要操作
</button>

<!-- 成功状态 -->
<div class="bg-success-50 border border-success-200 text-success-700">
  操作成功
</div>

<!-- 警告状态 -->
<div class="bg-warning-50 border border-warning-200 text-warning-700">
  请注意
</div>

<!-- 错误状态 -->
<div class="bg-error-50 border border-error-200 text-error-700">
  操作失败
</div>
```

## 📝 字体系统 (Typography)

### 字体族 (Font Families)
```css
/* 主字体 - 系统字体栈 */
--font-sans: 'Inter', 'SF Pro Display', -apple-system, BlinkMacSystemFont, 'Segoe UI', 'PingFang SC', 'Hiragino Sans GB', 'Microsoft YaHei', sans-serif;

/* 等宽字体 - 代码和数字 */
--font-mono: 'SF Mono', Monaco, 'Cascadia Code', 'Roboto Mono', Consolas, 'Courier New', monospace;
```

### 字体大小层级 (Font Sizes)
```css
/* 响应式字体大小 */
--text-xs:   0.75rem   /* 12px  - 辅助信息 */
--text-sm:   0.875rem  /* 14px  - 标签文字 */
--text-base: 1rem      /* 16px  - 正文文字 */
--text-lg:   1.125rem  /* 18px  - 强调文字 */
--text-xl:   1.25rem   /* 20px  - 小标题 */
--text-2xl:  1.5rem    /* 24px  - 标题 */
--text-3xl:  1.875rem  /* 30px  - 大标题 */
--text-4xl:  2.25rem   /* 36px  - 特大标题 */

/* 响应式字体 */
@media (min-width: 640px) {
  --text-sm: 0.875rem;  /* 14px */
  --text-base: 1rem;    /* 16px */
  --text-lg: 1.125rem;  /* 18px */
  --text-xl: 1.25rem;   /* 20px */
}

@media (min-width: 768px) {
  --text-base: 1rem;    /* 16px */
  --text-lg: 1.125rem;  /* 18px */
  --text-xl: 1.25rem;   /* 20px */
  --text-2xl: 1.5rem;   /* 24px */
}
```

### 字体权重 (Font Weights)
```css
--font-light:    300    /* 细体 - 弱化文字 */
--font-normal:   400    /* 常规 - 正文 */
--font-medium:   500    /* 中等 - 标签 */
--font-semibold: 600    /* 半粗 - 小标题 */
--font-bold:     700    /* 粗体 - 标题 */
```

### 行高系统 (Line Heights)
```css
--leading-none:    1      /* 无行高 - 标题 */
--leading-tight:   1.25   /* 紧凑 - 标题 */
--leading-snug:    1.375  /* 略紧 - 小标题 */
--leading-normal:  1.5    /* 常规 - 正文 */
--leading-relaxed: 1.625  /* 放松 - 长文 */
--leading-loose:   2      /* 宽松 - 诗歌 */
```

### 字体使用示例
```vue
<!-- 标题层级 -->
<h1 class="text-3xl sm:text-4xl font-bold leading-tight text-gray-900">
  页面主标题
</h1>
<h2 class="text-2xl sm:text-3xl font-semibold leading-tight text-gray-800">
  区块标题
</h2>
<h3 class="text-xl sm:text-2xl font-medium leading-snug text-gray-700">
  子标题
</h3>

<!-- 正文文字 -->
<p class="text-base sm:text-lg leading-normal text-gray-600">
  正文内容，描述详细的信息和说明。
</p>

<!-- 标签和辅助信息 -->
<label class="text-sm font-medium text-gray-700">
  表单标签
</label>
<span class="text-xs text-gray-500">
  辅助信息
</span>
```

## 📐 间距系统 (Spacing)

### 基础间距单位
```css
/* 间距变量 - 基于4px网格 */
--space-0:  0px      /* 0 */
--space-1:  0.25rem  /* 4px  - 最小间距 */
--space-2:  0.5rem   /* 8px  - 紧密间距 */
--space-3:  0.75rem  /* 12px - 小间距 */
--space-4:  1rem     /* 16px - 基础间距 */
--space-5:  1.25rem  /* 20px - 中间距 */
--space-6:  1.5rem   /* 24px - 标准间距 */
--space-8:  2rem     /* 32px - 大间距 */
--space-10: 2.5rem   /* 40px - 特大间距 */
--space-12: 3rem     /* 48px - 页面间距 */
--space-16: 4rem     /* 64px - 区块间距 */
--space-20: 5rem     /* 80px - 大区块间距 */
--space-24: 6rem     /* 96px - 页面边距 */
```

### 响应式间距
```css
/* 移动端紧凑间距 */
@media (max-width: 767px) {
  --space-4: 1rem;   /* 16px */
  --space-6: 1.5rem; /* 24px */
  --space-8: 2rem;   /* 32px */
}

/* 平板间距 */
@media (min-width: 768px) and (max-width: 1023px) {
  --space-4: 1rem;   /* 16px */
  --space-6: 1.5rem; /* 24px */
  --space-8: 2rem;   /* 32px */
}

/* 桌面端宽松间距 */
@media (min-width: 1024px) {
  --space-4: 1rem;   /* 16px */
  --space-6: 1.5rem; /* 24px */
  --space-8: 2.5rem; /* 40px */
}
```

### 间距使用场景
```vue
<!-- 页面布局间距 -->
<div class="px-4 sm:px-6 lg:px-8 py-6 sm:py-8 lg:py-12">
  <div class="max-w-7xl mx-auto">
    <!-- 内容 -->
  </div>
</div>

<!-- 组件内部间距 -->
<div class="p-4 sm:p-6 space-y-4 sm:space-y-6">
  <div class="flex items-center gap-3 sm:gap-4">
    <!-- 元素间距 -->
  </div>
</div>

<!-- 按钮组间距 -->
<div class="flex flex-col sm:flex-row gap-3 sm:gap-4">
  <button class="min-h-[48px] px-6">按钮1</button>
  <button class="min-h-[48px] px-6">按钮2</button>
</div>
```

## 🔧 组件库 (Component Library)

### 按钮组件 (Buttons)

#### 主要按钮
```vue
<template>
  <button
    :class="buttonClasses"
    :disabled="disabled"
    @click="$emit('click')"
  >
    <svg v-if="loading" class="animate-spin -ml-1 mr-3 h-4 w-4" fill="none" viewBox="0 0 24 24">
      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
      <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
    </svg>
    <slot />
  </button>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  variant: {
    type: String,
    default: 'primary',
    validator: (value) => ['primary', 'secondary', 'ghost', 'danger'].includes(value)
  },
  size: {
    type: String,
    default: 'md',
    validator: (value) => ['sm', 'md', 'lg'].includes(value)
  },
  disabled: Boolean,
  loading: Boolean,
  fullWidth: Boolean
})

const buttonClasses = computed(() => [
  'inline-flex items-center justify-center font-medium rounded-lg transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2',
  'disabled:opacity-50 disabled:cursor-not-allowed',
  {
    'min-w-[44px] min-h-[44px] px-4 py-3 text-sm': props.size === 'sm',
    'min-w-[48px] min-h-[48px] px-6 py-3 text-base': props.size === 'md',
    'min-w-[52px] min-h-[52px] px-8 py-4 text-lg': props.size === 'lg'
  },
  {
    'w-full': props.fullWidth,
    'bg-primary-500 text-white hover:bg-primary-600 focus:ring-primary-500': props.variant === 'primary',
    'border border-gray-300 bg-white text-gray-700 hover:bg-gray-50 focus:ring-primary-500': props.variant === 'secondary',
    'text-gray-700 hover:bg-gray-100 focus:ring-primary-500': props.variant === 'ghost',
    'bg-red-500 text-white hover:bg-red-600 focus:ring-red-500': props.variant === 'danger'
  }
])
</script>
```

#### 图标按钮
```vue
<template>
  <button
    :class="iconButtonClasses"
    :title="title"
    :aria-label="ariaLabel"
    @click="$emit('click')"
  >
    <slot />
  </button>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  size: {
    type: String,
    default: 'md',
    validator: (value) => ['sm', 'md', 'lg'].includes(value)
  },
  variant: {
    type: String,
    default: 'ghost',
    validator: (value) => ['ghost', 'primary', 'danger'].includes(value)
  },
  title: String,
  ariaLabel: String
})

const iconButtonClasses = computed(() => [
  'inline-flex items-center justify-center rounded-full transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2',
  {
    'w-8 h-8 p-1.5': props.size === 'sm',
    'w-10 h-10 p-2': props.size === 'md',
    'w-12 h-12 p-3': props.size === 'lg'
  },
  {
    'hover:bg-gray-100 text-gray-600 focus:ring-gray-500': props.variant === 'ghost',
    'hover:bg-primary-100 text-primary-600 focus:ring-primary-500': props.variant === 'primary',
    'hover:bg-red-100 text-red-600 focus:ring-red-500': props.variant === 'danger'
  }
])
</script>
```

### 输入框组件 (Input)
```vue
<template>
  <div class="w-full">
    <label v-if="label" :for="inputId" class="block text-sm font-medium text-gray-700 mb-2">
      {{ label }}
      <span v-if="required" class="text-red-500">*</span>
    </label>
    
    <div class="relative">
      <input
        :id="inputId"
        v-model="inputValue"
        :type="type"
        :placeholder="placeholder"
        :disabled="disabled"
        :class="inputClasses"
        @input="handleInput"
        @blur="handleBlur"
        @focus="handleFocus"
      />
      
      <div v-if="$slots.icon" class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
        <slot name="icon" />
      </div>
      
      <div v-if="error" class="absolute inset-y-0 right-0 pr-3 flex items-center pointer-events-none">
        <svg class="w-5 h-5 text-red-400" fill="currentColor" viewBox="0 0 20 20">
          <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd"></path>
        </svg>
      </div>
    </div>
    
    <p v-if="error" class="mt-2 text-sm text-red-600">{{ error }}</p>
    <p v-else-if="hint" class="mt-2 text-sm text-gray-500">{{ hint }}</p>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'

const props = defineProps({
  modelValue: [String, Number],
  label: String,
  type: {
    type: String,
    default: 'text'
  },
  placeholder: String,
  disabled: Boolean,
  required: Boolean,
  error: String,
  hint: String,
  size: {
    type: String,
    default: 'md'
  }
})

const emit = defineEmits(['update:modelValue', 'blur', 'focus'])

const inputValue = ref(props.modelValue)
const inputId = `input-${Math.random().toString(36).substr(2, 9)}`

const inputClasses = computed(() => [
  'block w-full border rounded-lg transition-colors focus:outline-none focus:ring-2 focus:ring-offset-0',
  'disabled:bg-gray-50 disabled:text-gray-500 disabled:cursor-not-allowed',
  {
    'px-3 py-2 text-sm': props.size === 'sm',
    'px-4 py-3 text-base': props.size === 'md',
    'px-6 py-4 text-lg': props.size === 'lg'
  },
  {
    'pl-10': !!useSlots().icon,
    'pr-10': !!props.error,
    'border-gray-300 focus:border-primary-500 focus:ring-primary-500': !props.error,
    'border-red-300 focus:border-red-500 focus:ring-red-500': props.error
  }
])

function handleInput(event) {
  emit('update:modelValue', event.target.value)
}

function handleBlur(event) {
  emit('blur', event)
}

function handleFocus(event) {
  emit('focus', event)
}
</script>
```

### 卡片组件 (Card)
```vue
<template>
  <div :class="cardClasses">
    <div v-if="$slots.header" :class="headerClasses">
      <slot name="header" />
    </div>
    
    <div :class="bodyClasses">
      <slot />
    </div>
    
    <div v-if="$slots.footer" :class="footerClasses">
      <slot name="footer" />
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  variant: {
    type: String,
    default: 'default',
    validator: (value) => ['default', 'elevated', 'bordered'].includes(value)
  },
  padding: {
    type: String,
    default: 'md',
    validator: (value) => ['none', 'sm', 'md', 'lg'].includes(value)
  },
  rounded: {
    type: String,
    default: 'lg',
    validator: (value) => ['none', 'sm', 'md', 'lg', 'xl'].includes(value)
  }
})

const cardClasses = computed(() => [
  'bg-white',
  {
    'shadow-sm border border-gray-100': props.variant === 'default',
    'shadow-lg border border-gray-100': props.variant === 'elevated',
    'border border-gray-200': props.variant === 'bordered'
  },
  {
    'rounded-none': props.rounded === 'none',
    'rounded': props.rounded === 'sm',
    'rounded-lg': props.rounded === 'md',
    'rounded-xl': props.rounded === 'lg',
    'rounded-2xl': props.rounded === 'xl'
  }
])

const headerClasses = computed(() => [
  {
    'p-4 sm:p-6': props.padding === 'md',
    'p-3 sm:p-4': props.padding === 'sm',
    'p-6 sm:p-8': props.padding === 'lg',
    'p-0': props.padding === 'none'
  },
  'border-b border-gray-100'
])

const bodyClasses = computed(() => [
  {
    'p-4 sm:p-6': props.padding === 'md',
    'p-3 sm:p-4': props.padding === 'sm',
    'p-6 sm:p-8': props.padding === 'lg',
    'p-0': props.padding === 'none'
  }
])

const footerClasses = computed(() => [
  {
    'p-4 sm:p-6': props.padding === 'md',
    'p-3 sm:p-4': props.padding === 'sm',
    'p-6 sm:p-8': props.padding === 'lg',
    'p-0': props.padding === 'none'
  },
  'border-t border-gray-100'
])
</script>
```

### 模态框组件 (Modal)
```vue
<template>
  <Teleport to="body">
    <Transition
      enter-active-class="transition ease-out duration-300"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-active-class="transition ease-in duration-200"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div v-if="show" class="fixed inset-0 z-50 overflow-y-auto">
        <!-- 背景遮罩 -->
        <div 
          class="fixed inset-0 bg-black bg-opacity-50 transition-opacity"
          @click="handleBackdropClick"
        ></div>
        
        <!-- 模态框内容 -->
        <div class="flex min-h-full items-center justify-center p-4">
          <Transition
            enter-active-class="transition ease-out duration-300"
            enter-from-class="opacity-0 scale-95"
            enter-to-class="opacity-100 scale-100"
            leave-active-class="transition ease-in duration-200"
            leave-from-class="opacity-100 scale-100"
            leave-to-class="opacity-0 scale-95"
          >
            <div
              v-if="show"
              :class="modalClasses"
              @click.stop
            >
              <!-- 头部 -->
              <div v-if="$slots.header" class="flex items-center justify-between p-6 border-b border-gray-200">
                <slot name="header" />
                <button
                  @click="close"
                  class="text-gray-400 hover:text-gray-600 focus:outline-none focus:ring-2 focus:ring-primary-500 rounded-full p-1"
                >
                  <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
                  </svg>
                </button>
              </div>
              
              <!-- 内容 -->
              <div class="p-6">
                <slot />
              </div>
              
              <!-- 底部 -->
              <div v-if="$slots.footer" class="flex items-center justify-end gap-3 p-6 border-t border-gray-200">
                <slot name="footer" />
              </div>
            </div>
          </Transition>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup>
import { computed, watch, onMounted, onUnmounted } from 'vue'

const props = defineProps({
  show: Boolean,
  size: {
    type: String,
    default: 'md',
    validator: (value) => ['sm', 'md', 'lg', 'xl', 'full'].includes(value)
  },
  closeOnBackdrop: {
    type: Boolean,
    default: true
  }
})

const emit = defineEmits(['close', 'show'])

const modalClasses = computed(() => [
  'relative bg-white rounded-xl shadow-xl w-full mx-auto',
  {
    'max-w-sm': props.size === 'sm',
    'max-w-md': props.size === 'md',
    'max-w-lg': props.size === 'lg',
    'max-w-2xl': props.size === 'xl',
    'max-w-4xl': props.size === 'full'
  }
])

function close() {
  emit('close')
}

function handleBackdropClick() {
  if (props.closeOnBackdrop) {
    close()
  }
}

function handleEscapeKey(event) {
  if (event.key === 'Escape') {
    close()
  }
}

watch(() => props.show, (newShow) => {
  if (newShow) {
    document.addEventListener('keydown', handleEscapeKey)
    document.body.style.overflow = 'hidden'
  } else {
    document.removeEventListener('keydown', handleEscapeKey)
    document.body.style.overflow = ''
  }
})

onMounted(() => {
  if (props.show) {
    document.body.style.overflow = 'hidden'
  }
})

onUnmounted(() => {
  document.body.style.overflow = ''
  document.removeEventListener('keydown', handleEscapeKey)
})
</script>
```

## 🎯 布局网格 (Layout Grid)

### 响应式网格系统
```css
/* 网格容器 */
.container {
  width: 100%;
  max-width: 1280px;
  margin: 0 auto;
  padding: 0 1rem;
}

/* 响应式断点 */
@media (min-width: 640px) {
  .container {
    max-width: 640px;
    padding: 0 1.5rem;
  }
}

@media (min-width: 768px) {
  .container {
    max-width: 768px;
    padding: 0 2rem;
  }
}

@media (min-width: 1024px) {
  .container {
    max-width: 1024px;
  }
}

@media (min-width: 1280px) {
  .container {
    max-width: 1280px;
  }
}

/* 网格系统 */
.grid {
  display: grid;
  gap: 1.5rem;
}

.grid-cols-1 { grid-template-columns: repeat(1, minmax(0, 1fr)); }
.grid-cols-2 { grid-template-columns: repeat(2, minmax(0, 1fr)); }
.grid-cols-3 { grid-template-columns: repeat(3, minmax(0, 1fr)); }
.grid-cols-4 { grid-template-columns: repeat(4, minmax(0, 1fr)); }

/* 响应式网格 */
@media (min-width: 640px) {
  .sm\:grid-cols-2 { grid-template-columns: repeat(2, minmax(0, 1fr)); }
  .sm\:grid-cols-3 { grid-template-columns: repeat(3, minmax(0, 1fr)); }
}

@media (min-width: 768px) {
  .md\:grid-cols-3 { grid-template-columns: repeat(3, minmax(0, 1fr)); }
  .md\:grid-cols-4 { grid-template-columns: repeat(4, minmax(0, 1fr)); }
}

@media (min-width: 1024px) {
  .lg\:grid-cols-4 { grid-template-columns: repeat(4, minmax(0, 1fr)); }
  .lg\:grid-cols-5 { grid-template-columns: repeat(5, minmax(0, 1fr)); }
}
```

### 弹性布局模式
```vue
<!-- 水平布局 -->
<div class="flex items-center gap-4">
  <div class="flex-shrink-0">
    <!-- 固定宽度内容 -->
  </div>
  <div class="flex-1 min-w-0">
    <!-- 自适应内容 -->
  </div>
</div>

<!-- 垂直布局 -->
<div class="flex flex-col gap-6">
  <div class="flex-shrink-0">
    <!-- 固定高度内容 -->
  </div>
  <div class="flex-1 min-h-0">
    <!-- 自适应内容 -->
  </div>
</div>

<!-- 响应式布局 -->
<div class="flex flex-col lg:flex-row gap-4 lg:gap-8">
  <div class="flex-1 min-w-0">
    <!-- 主内容区域 -->
  </div>
  <div class="w-full lg:w-80 flex-shrink-0">
    <!-- 侧边栏 -->
  </div>
</div>
```

## 🔄 交互规范 (Interaction Patterns)

### 触摸优化 (Touch Optimization)
```css
/* 触摸目标最小尺寸 */
.touch-target {
  min-width: 44px;
  min-height: 44px;
  padding: 12px;
  margin: 4px;
}

/* 重要按钮更大 */
.primary-button {
  min-width: 48px;
  min-height: 48px;
  padding: 16px 24px;
}

/* 列表项触摸区域 */
.list-item {
  min-height: 56px;
  padding: 12px 16px;
}

/* 安全区域适配 */
.safe-area-top {
  padding-top: env(safe-area-inset-top);
}

.safe-area-bottom {
  padding-bottom: env(safe-area-inset-bottom);
}

.safe-area-left {
  padding-left: env(safe-area-inset-left);
}

.safe-area-right {
  padding-right: env(safe-area-inset-right);
}
```

### 动画系统 (Animation System)
```css
/* 缓动函数 */
--ease-linear: linear;
--ease-in: cubic-bezier(0.4, 0, 1, 1);
--ease-out: cubic-bezier(0, 0, 0.2, 1);
--ease-in-out: cubic-bezier(0.4, 0, 0.2, 1);

/* 动画时长 */
--duration-fast: 150ms;
--duration-normal: 300ms;
--duration-slow: 500ms;

/* 常用动画 */
@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}

@keyframes slideUp {
  from { transform: translateY(20px); opacity: 0; }
  to { transform: translateY(0); opacity: 1; }
}

@keyframes scaleIn {
  from { transform: scale(0.95); opacity: 0; }
  to { transform: scale(1); opacity: 1; }
}

/* 动画类 */
.animate-fade-in {
  animation: fadeIn var(--duration-normal) var(--ease-out);
}

.animate-slide-up {
  animation: slideUp var(--duration-normal) var(--ease-out);
}

.animate-scale-in {
  animation: scaleIn var(--duration-fast) var(--ease-out);
}

/* 过渡效果 */
.transition-base {
  transition: all var(--duration-normal) var(--ease-in-out);
}

.transition-colors {
  transition: color var(--duration-fast) var(--ease-in-out),
              background-color var(--duration-fast) var(--ease-in-out),
              border-color var(--duration-fast) var(--ease-in-out);
}

.transition-transform {
  transition: transform var(--duration-normal) var(--ease-in-out);
}
```

### 状态反馈 (State Feedback)
```vue
<!-- 加载状态 -->
<div class="animate-pulse">
  <div class="h-4 bg-gray-200 rounded w-3/4 mb-2"></div>
  <div class="h-4 bg-gray-200 rounded w-1/2"></div>
</div>

<!-- 悬停状态 -->
<button class="hover:bg-primary-50 hover:text-primary-600 transition-colors">
  悬停效果
</button>

<!-- 激活状态 -->
<button class="active:scale-95 active:bg-primary-700 transition-transform">
  按下效果
</button>

<!-- 焦点状态 -->
<button class="focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2">
  焦点效果
</button>

<!-- 禁用状态 -->
<button class="opacity-50 cursor-not-allowed" disabled>
  禁用状态
</button>
```

## 📱 移动端优化 (Mobile Optimization)

### 移动端断点系统
```css
/* 移动端断点 */
@media (max-width: 767px) {
  /* 小屏手机 */
  .mobile-only {
    display: block;
  }
  
  .desktop-only {
    display: none;
  }
}

@media (min-width: 768px) {
  /* 平板及以上 */
  .mobile-only {
    display: none;
  }
  
  .desktop-only {
    display: block;
  }
}
```

### 移动端布局模式
```vue
<!-- 移动端优先布局 -->
<div class="flex flex-col gap-4 sm:flex-row sm:gap-6 lg:gap-8">
  <!-- 移动端垂直排列，桌面端水平排列 -->
  <div class="w-full sm:w-1/2 lg:w-1/3">
    <Card>内容1</Card>
  </div>
  <div class="w-full sm:w-1/2 lg:w-1/3">
    <Card>内容2</Card>
  </div>
  <div class="w-full sm:w-1/2 lg:w-1/3">
    <Card>内容3</Card>
  </div>
</div>

<!-- 移动端底部导航 -->
<div class="fixed bottom-0 left-0 right-0 bg-white border-t border-gray-200 safe-area-bottom sm:hidden">
  <div class="flex items-center justify-around py-2">
    <button class="flex flex-col items-center py-2 px-3">
      <svg class="w-6 h-6 mb-1">...</svg>
      <span class="text-xs">首页</span>
    </button>
    <button class="flex flex-col items-center py-2 px-3">
      <svg class="w-6 h-6 mb-1">...</svg>
      <span class="text-xs">房间</span>
    </button>
    <button class="flex flex-col items-center py-2 px-3">
      <svg class="w-6 h-6 mb-1">...</svg>
      <span class="text-xs">我的</span>
    </button>
  </div>
</div>
```

### 移动端交互模式
```vue
<!-- 滑动检测组件 -->
<template>
  <div 
    @touchstart="handleTouchStart"
    @touchend="handleTouchEnd"
    @touchmove="handleTouchMove"
  >
    <slot />
  </div>
</template>

<script setup>
const emit = defineEmits(['swipeLeft', 'swipeRight', 'swipeUp', 'swipeDown'])

let touchStartX = 0
let touchStartY = 0
let touchEndX = 0
let touchEndY = 0

function handleTouchStart(e) {
  touchStartX = e.touches[0].clientX
  touchStartY = e.touches[0].clientY
}

function handleTouchEnd(e) {
  touchEndX = e.changedTouches[0].clientX
  touchEndY = e.changedTouches[0].clientY
  handleSwipe()
}

function handleTouchMove(e) {
  e.preventDefault() // 防止页面滚动
}

function handleSwipe() {
  const deltaX = touchEndX - touchStartX
  const deltaY = touchEndY - touchStartY
  const minSwipeDistance = 50

  if (Math.abs(deltaX) > Math.abs(deltaY)) {
    // 水平滑动
    if (Math.abs(deltaX) > minSwipeDistance) {
      if (deltaX > 0) {
        emit('swipeRight')
      } else {
        emit('swipeLeft')
      }
    }
  } else {
    // 垂直滑动
    if (Math.abs(deltaY) > minSwipeDistance) {
      if (deltaY > 0) {
        emit('swipeDown')
      } else {
        emit('swipeUp')
      }
    }
  }
}
</script>
```

## ♿ 无障碍性 (Accessibility)

### ARIA 属性使用
```vue
<!-- 按钮 -->
<button 
  class="px-4 py-2 bg-primary text-white rounded"
  aria-label="关闭弹窗"
  aria-describedby="modal-description"
>
  <svg aria-hidden="true" class="w-4 h-4">...</svg>
</button>

<!-- 输入框 -->
<label for="username" class="block text-sm font-medium mb-1">
  用户名 <span aria-label="必填">*</span>
</label>
<input
  id="username"
  v-model="username"
  type="text"
  required
  aria-describedby="username-error"
  :aria-invalid="!!errors.username"
>
<div 
  v-if="errors.username" 
  id="username-error" 
  role="alert"
  class="text-red-600 text-sm mt-1"
>
  {{ errors.username }}
</div>

<!-- 模态框 -->
<div 
  v-if="showModal"
  role="dialog"
  aria-modal="true"
  aria-labelledby="modal-title"
  aria-describedby="modal-description"
  class="fixed inset-0 bg-black bg-opacity-50"
>
  <div class="flex items-center justify-center min-h-screen p-4">
    <div class="bg-white rounded-lg p-6 max-w-md w-full">
      <h2 id="modal-title" class="text-lg font-semibold mb-2">
        确认操作
      </h2>
      <p id="modal-description" class="text-gray-600 mb-4">
        此操作无法撤销，确定要继续吗？
      </p>
      <div class="flex gap-3 justify-end">
        <button @click="closeModal">取消</button>
        <button @click="confirm" class="bg-red-500 text-white px-4 py-2 rounded">
          确定
        </button>
      </div>
    </div>
  </div>
</div>
```

### 键盘导航
```vue
<!-- 焦点管理 -->
<template>
  <div 
    ref="container"
    tabindex="0"
    @keydown="handleKeydown"
    class="focus:outline-none focus:ring-2 focus:ring-primary-500"
  >
    <!-- 内容 -->
  </div>
</template>

<script setup>
function handleKeydown(event) {
  switch (event.key) {
    case 'ArrowUp':
      focusPrevious()
      event.preventDefault()
      break
    case 'ArrowDown':
      focusNext()
      event.preventDefault()
      break
    case 'Enter':
    case ' ':
      activateCurrent()
      event.preventDefault()
      break
    case 'Escape':
      closeModal()
      event.preventDefault()
      break
  }
}
</script>

<!-- 跳转链接 -->
<nav aria-label="主导航">
  <ul class="flex gap-4">
    <li>
      <a 
        href="#main-content" 
        class="sr-only focus:not-sr-only focus:absolute focus:top-0 focus:left-0 bg-primary text-white px-4 py-2 z-50"
      >
        跳转到主内容
      </a>
    </li>
    <li><a href="/">首页</a></li>
    <li><a href="/about">关于</a></li>
    <li><a href="/contact">联系</a></li>
  </ul>
</nav>

<main id="main-content" tabindex="-1">
  <!-- 主内容 -->
</main>
```

## 🎯 项目特定规范 (Project-Specific Guidelines)

### RoomPage.vue 样式约定
```css
/* 基于实际项目的颜色系统 */
:root {
  /* 主背景色 */
  --bg-dark: #0a0a0a;
  --bg-surface: rgba(255, 255, 255, 0.1);
  
  /* 文字颜色 */
  --text-main: #ffffff;
  --text-muted: rgba(255, 255, 255, 0.7);
  --text-subtle: rgba(255, 255, 255, 0.5);
  
  /* 功能色 */
  --primary: #3b82f6;
  --primary-hover: #2563eb;
  --secondary: #a3c9a8;
  --success: #10b981;
  --warning: #f59e0b;
  --error: #ef4444;
}

/* 组件样式类 */
.room-page {
  @apply relative z-10 flex-1 flex flex-col h-screen min-h-screen overflow-hidden;
  @apply bg-dark-bg transition-colors duration-1000;
}

.surface {
  @apply bg-surface backdrop-blur-md;
  @apply border border-white/10 shadow-wood;
}

.text-main {
  @apply text-text-main;
}

.text-muted {
  @apply text-text-muted;
}

.touch-optimize {
  @apply min-w-[44px] min-h-[44px];
  @apply active:scale-95 transition-transform duration-150;
}
```

### 响应式组件示例
```vue
<!-- 响应式成员列表组件 -->
<template>
  <div class="member-list">
    <!-- 移动端紧凑显示 -->
    <div class="sm:hidden flex items-center gap-2">
      <div class="flex -space-x-1.5">
        <div 
          v-for="(member, index) in displayMembers" 
          :key="member.id"
          :class="['w-6 h-6 rounded-full border-2 border-surface flex items-center justify-center text-[10px]', memberClass(member)]"
        >
          {{ member.avatar }}
        </div>
        <div v-if="remainingCount > 0" class="w-6 h-6 rounded-full border-2 border-surface bg-surface flex items-center justify-center text-[10px] text-text-muted">
          +{{ remainingCount }}
        </div>
      </div>
      <span class="text-xs text-text-main font-medium">{{ totalCount }}人</span>
    </div>
    
    <!-- 桌面端详细显示 -->
    <div class="hidden sm:block space-y-2">
      <div 
        v-for="member in members" 
        :key="member.id"
        class="flex items-center justify-between p-3 rounded-lg hover:bg-white/5 transition-colors"
      >
        <div class="flex items-center gap-3">
          <div :class="['w-10 h-10 rounded-full flex items-center justify-center', memberClass(member)]">
            {{ member.avatar }}
          </div>
          <div>
            <div class="text-text-main font-medium">{{ member.name }}</div>
            <div class="text-text-muted text-sm">{{ member.role }}</div>
          </div>
        </div>
        <div class="flex items-center gap-2">
          <div :class="['w-2 h-2 rounded-full', networkStatusColors[member.status]]"></div>
          <span class="text-xs text-text-muted">{{ member.status }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  members: {
    type: Array,
    default: () => []
  },
  maxDisplay: {
    type: Number,
    default: 3
  }
})

const displayMembers = computed(() => props.members.slice(0, props.maxDisplay))
const remainingCount = computed(() => Math.max(0, props.members.length - props.maxDisplay))
const totalCount = computed(() => props.members.length)

function memberClass(member) {
  const classes = {
    online: 'bg-green-500/20 text-green-300',
    away: 'bg-yellow-500/20 text-yellow-300',
    offline: 'bg-gray-500/20 text-gray-400'
  }
  return classes[member.status] || classes.offline
}

const networkStatusColors = {
  good: '#10b981',
  poor: '#f59e0b',
  offline: '#6b7280'
}
</script>
```

## 📋 开发指南 (Developer Guidelines)

### 代码组织规范
```typescript
// 样式组织结构
src/
├── assets/
│   ├── styles/
│   │   ├── variables.css     # CSS 变量定义
│   │   ├── components.css    # 组件样式
│   │   ├── utilities.css     # 工具类
│   │   └── responsive.css    # 响应式样式
│   └── images/
├── components/
│   ├── ui/                   # 基础组件
│   ├── layout/               # 布局组件
│   └── features/             # 功能组件
└── views/
```

### 样式命名规范
```css
/* BEM 命名规范 */
.component-name {}
.component-name__element {}
.component-name--modifier {}

/* 示例 */
.btn {}
.btn__icon {}
.btn--primary {}
.btn--large {}

/* 状态类 */
.is-loading {}
.is-error {}
.is-success {}

/* 工具类 */
.u-hidden {}
.u-visually-hidden {}
.u-text-center {}
```

### 响应式开发最佳实践
```vue
<!-- 移动端优先 -->
<template>
  <div class="container">
    <!-- 移动端布局 -->
    <div class="stack gap-4 sm:flex sm:gap-6 lg:gap-8">
      <div class="w-full sm:w-1/2 lg:w-1/3">
        <!-- 内容 -->
      </div>
    </div>
    
    <!-- 响应式文本 -->
    <h1 class="text-2xl sm:text-3xl lg:text-4xl font-bold">
      响应式标题
    </h1>
    
    <!-- 响应式间距 -->
    <div class="p-4 sm:p-6 lg:p-8">
      <!-- 内容 -->
    </div>
  </div>
</template>

<!-- 性能优化 -->
<template>
  <!-- 使用 CSS containment 优化渲染 -->
  <div class="component" style="contain: layout style paint">
    <!-- 大列表使用虚拟滚动 -->
    <VirtualList :items="items" />
    
    <!-- 图片懒加载 -->
    <img 
      :src="imageSrc" 
      loading="lazy" 
      decoding="async"
      class="w-full h-auto"
    >
  </div>
</template>
```

## 📊 性能指南 (Performance Guidelines)

### 关键性能指标
```css
/* 关键渲染路径优化 */
.critical {
  /* 内联关键 CSS */
  font-display: swap;
  contain: layout style paint;
}

/* 避免布局抖动 */
.stable-layout {
  contain: layout;
  content-visibility: auto;
}

/* 图片优化 */
.responsive-image {
  object-fit: cover;
  object-position: center;
  loading: lazy;
  decoding: async;
}
```

### 移动端性能优化
```css
/* 减少重绘和重排 */
.optimized-animation {
  will-change: transform;
  transform: translateZ(0); /* 硬件加速 */
  backface-visibility: hidden;
}

/* 触摸响应优化 */
.touch-responsive {
  touch-action: manipulation;
  -webkit-tap-highlight-color: transparent;
}

/* 滚动优化 */
.smooth-scroll {
  -webkit-overflow-scrolling: touch;
  scroll-behavior: smooth;
}
```

---

## 📝 文档维护

### 版本历史
- **v1.0** (2025-12-31): 初始版本，包含完整的设计系统基础
  - 颜色系统
  - 字体系统  
  - 间距系统
  - 组件库
  - 响应式规范
  - 移动端优化
  - 无障碍性指南
  - 项目特定规范

### 更新指南
1. **颜色系统变更**: 更新 CSS 变量并同步到所有组件
2. **新增组件**: 在组件库中添加并提供使用示例
3. **响应式调整**: 更新断点系统并测试兼容性
4. **无障碍性更新**: 遵循 WCAG 2.1 AA 标准

### 贡献指南
- 所有设计决策需要文档化
- 新组件需要包含无障碍性考虑
- 响应式设计必须测试多个设备
- 性能影响需要评估和记录

---

**版本**: v1.0.1
  <div class="flex-1 min-w-0">
    <!-- 自适应内容 -->
  </div>
  <div class="flex-shrink-0">
    <!-- 固定宽度内容 -->
  </div>
</div>

<!-- 垂直布局 -->
<div class="flex flex-col gap-4">
  <div class="flex-1">
    <!-- 自适应高度 -->
  </div>
  <div class="flex-shrink-0">
    <!-- 固定高度 -->
  </div>
</div>

<!-- 居中对齐 -->
<div class="flex items-center justify-center min-h-[200px]">
  <div class="text-center">
    <!-- 居中内容 -->
  </div>
</div>
```

## 🎨 图标系统 (Icon System)

### 图标规格
- **尺寸**: 16px, 20px, 24px, 32px
- **风格**: 线性图标，2px线宽
- **颜色**: 继承父元素颜色
- **库**: Heroicons v2

### 图标使用
```vue
<!-- 基本图标 -->
<svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
</svg>

<!-- 语义化图标组件 -->
<template>
  <component :is="iconComponent" :class="iconClasses" />
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  name: {
    type: String,
    required: true
  },
  size: {
    type: String,
    default: 'md',
    validator: (value) => ['sm', 'md', 'lg', 'xl'].includes(value)
  },
  color: {
    type: String,
    default: 'currentColor'
  }
})

const iconComponent = computed(() => {
  const icons = {
    'home': 'HomeIcon',
    'user': 'UserIcon',
    'settings': 'CogIcon',
    'search': 'MagnifyingGlassIcon',
    'close': 'XMarkIcon',
    'menu': 'Bars3Icon'
  }
  return icons[props.name] || 'QuestionMarkCircleIcon'
})

const iconClasses = computed(() => [
  'inline-flex',
  {
    'w-4 h-4': props.size === 'sm',
    'w-5 h-5': props.size === 'md',
    'w-6 h-6': props.size === 'lg',
    'w-8 h-8': props.size === 'xl'
  }
])
</script>
```

## 🔍 交互模式 (Interaction Patterns)

### 状态变化
```css
/* 悬停状态 */
.hover\:scale-105:hover {
  transform: scale(1.05);
}

.hover\:bg-gray-100:hover {
  background-color: rgb(243 244 246);
}

/* 焦点状态 */
.focus\:ring-2:focus {
  box-shadow: 0 0 0 2px rgb(59 130 246 / 0.5);
}

.focus\:outline-none:focus {
  outline: 2px solid transparent;
  outline-offset: 2px;
}

/* 激活状态 */
.active\:scale-95:active {
  transform: scale(0.95);
}

/* 禁用状态 */
.disabled\:opacity-50:disabled {
  opacity: 0.5;
}

.disabled\:cursor-not-allowed:disabled {
  cursor: not-allowed;
}
```

### 动画效果
```css
/* 过渡动画 */
.transition {
  transition-property: color, background-color, border-color, text-decoration-color, fill, stroke, opacity, box-shadow, transform, filter, backdrop-filter;
  transition-timing-function: cubic-bezier(0.4, 0, 0.2, 1);
  transition-duration: 150ms;
}

.transition-colors {
  transition-property: color, background-color, border-color, text-decoration-color, fill, stroke;
  transition-timing-function: cubic-bezier(0.4, 0, 0.2, 1);
  transition-duration: 150ms;
}

.transition-transform {
  transition-property: transform;
  transition-timing-function: cubic-bezier(0.4, 0, 0.2, 1);
  transition-duration: 150ms;
}

/* 加载动画 */
@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.animate-spin {
  animation: spin 1s linear infinite;
}

@keyframes pulse {
  50% {
    opacity: 0.5;
  }
}

.animate-pulse {
  animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
}
```

### 手势交互
```vue
<!-- 滑动手势支持 -->
<div
  @touchstart="handleTouchStart"
  @touchmove="handleTouchMove"
  @touchend="handleTouchEnd"
  class="touch-pan-x touch-pan-y"
>
  <!-- 可交互内容 -->
</div>

<!-- 长按手势 -->
<div
  @touchstart="handleLongPressStart"
  @touchend="handleLongPressEnd"
  class="relative"
>
  <!-- 长按内容 -->
</div>
```

## ♿ 可访问性指南 (Accessibility)

### 键盘导航
```vue
<!-- Tab 导航顺序 -->
<button tabindex="0">第一个按钮</button>
<button tabindex="1">第二个按钮</button>
<button tabindex="-1">不可聚焦按钮</button>

<!-- 跳过链接 -->
<a href="#main-content" class="sr-only focus:not-sr-only focus:absolute focus:top-0 focus:left-0 bg-primary-500 text-white p-2 z-50">
  跳转到主内容
</a>

<!-- 键盘快捷键 -->
<div
  @keydown="handleKeyboard"
  tabindex="0"
  role="button"
  aria-label="可点击的卡片"
>
  <!-- 内容 -->
</div>
```

### 屏幕阅读器支持
```vue
<!-- 语义化HTML -->
<main role="main" aria-labelledby="page-title">
  <h1 id="page-title">页面标题</h1>
  
  <nav role="navigation" aria-label="主导航">
    <ul>
      <li><a href="/" aria-current="page">首页</a></li>
      <li><a href="/about">关于</a></li>
    </ul>
  </nav>
  
  <section aria-labelledby="section-title">
    <h2 id="section-title">区块标题</h2>
    <p>区块内容</p>
  </section>
</main>

<!-- ARIA 标签 -->
<button aria-label="关闭对话框">
  <svg aria-hidden="true"><!-- 图标 --></svg>
</button>

<div role="status" aria-live="polite">
  <!-- 状态信息 -->
</div>

<div role="alert" aria-live="assertive">
  <!-- 重要警告信息 -->
</div>
```

### 颜色对比度
```css
/* 确保足够的对比度 */
.text-gray-900 { color: #111827; }      /* 对比度: 21:1 */
.text-gray-700 { color: #374151; }      /* 对比度: 12.6:1 */
.text-gray-600 { color: #4b5563; }      /* 对比度: 7.3:1 */
.text-gray-500 { color: #6b7280; }      /* 对比度: 4.5:1 */
.text-gray-400 { color: #9ca3af; }      /* 对比度: 2.8:1 */

/* 背景色 */
.bg-white { background-color: #ffffff; }
.bg-gray-50 { background-color: #f9fafb; }
.bg-gray-100 { background-color: #f3f4f6; }
```

## 📱 移动端适配

### 触摸优化
```css
/* 触摸目标最小尺寸 */
.touch-target {
  min-width: 44px;
  min-height: 44px;
}

/* 触摸友好的间距 */
.touch-spacing {
  margin: 8px;
  padding: 12px;
}

/* 防止用户选择 */
.no-select {
  -webkit-user-select: none;
  -moz-user-select: none;
  -ms-user-select: none;
  user-select: none;
}

/* 平滑滚动 */
.smooth-scroll {
  -webkit-overflow-scrolling: touch;
  scroll-behavior: smooth;
}
```

### 响应式图片
```vue
<!-- 响应式图片 -->
<picture>
  <source media="(max-width: 767px)" srcset="image-mobile.webp">
  <source media="(min-width: 768px) and (max-width: 1023px)" srcset="image-tablet.webp">
  <source media="(min-width: 1024px)" srcset="image-desktop.webp">
  <img 
    src="image-mobile.webp" 
    alt="描述文字"
    class="w-full h-auto object-cover"
    loading="lazy"
  >
</picture>
```

## 🚀 开发规范

### CSS 组织结构
```css
/* 1. 基础样式 */
@tailwind base;
@tailwind components;
@tailwind utilities;

/* 2. 自定义基础样式 */
@layer base {
  html {
    scroll-behavior: smooth;
  }
  
  body {
    @apply text-gray-900 antialiased;
  }
}

/* 3. 组件样式 */
@layer components {
  .btn-primary {
    @apply bg-primary-500 text-white px-6 py-3 rounded-lg font-medium hover:bg-primary-600 focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 transition-colors;
  }
  
  .card {
    @apply bg-white rounded-xl shadow-sm border border-gray-100 p-6;
  }
}

/* 4. 工具样式 */
@layer utilities {
  .text-balance {
    text-wrap: balance;
  }
  
  .safe-area-bottom {
    padding-bottom: env(safe-area-inset-bottom);
  }
}
```

### Vue 组件规范
```vue
<template>
  <!-- 模板内容 -->
</template>

<script setup>
import { ref, computed } from 'vue'

// Props 定义
const props = defineProps({
  // props
})

// Emits 定义
const emit = defineEmits([
  'update:modelValue',
  'change',
  'focus',
  'blur'
])

// 响应式数据
const isLoading = ref(false)

// 计算属性
const classes = computed(() => [
  // class logic
])

// 方法
function handleClick() {
  // method logic
}
</script>

<style scoped>
/* 组件样式 */
.component {
  /* styles */
}
</style>
```

### 性能优化
```css
/* 硬件加速 */
.gpu-accelerated {
  transform: translateZ(0);
  will-change: transform;
}

/* 减少重绘 */
.reduce-paints {
  will-change: auto;
  transform: translateZ(0);
}

/* 图片优化 */
.lazy-image {
  loading: lazy;
  decoding: async;
}
```

## 📋 实施检查清单

### 设计一致性检查
- [ ] 颜色使用符合设计系统
- [ ] 字体大小层级正确
- [ ] 间距遵循8px网格
- [ ] 组件样式统一
- [ ] 交互状态完整

### 响应式设计检查
- [ ] 移动端布局适配
- [ ] 触摸目标尺寸合适
- [ ] 字体大小响应式
- [ ] 图片自适应
- [ ] 导航适配移动端

### 可访问性检查
- [ ] 键盘导航完整
- [ ] 屏幕阅读器兼容
- [ ] 颜色对比度充足
- [ ] ARIA标签正确
- [ ] 焦点状态清晰

### 性能优化检查
- [ ] CSS优化实现
- [ ] 图片懒加载
- [ ] 动画性能良好
- [ ] 组件代码优化
- [ ] 加载状态友好

---

**文档维护**: 请在设计新组件时参考此文档，确保设计一致性
**更新周期**: 建议每季度回顾更新设计系统
**最后更新**: 2025-12-31
**版本**: v1.0