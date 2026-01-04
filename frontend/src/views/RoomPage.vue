<template>
  <div id="room-page" class="relative z-10 flex-1 flex flex-col h-screen min-h-screen overflow-hidden bg-dark-bg transition-colors duration-1000">
    <!-- Toast 通知 (默认隐藏) - 响应式优化 -->
    <div id="toast" v-if="toastVisible" class="absolute top-20 sm:top-24 left-1/2 -translate-x-1/2 z-[60] pointer-events-none opacity-100 transition-opacity duration-300">
      <div class="bg-black/60 backdrop-blur-md text-white px-4 sm:px-6 py-2.5 sm:py-2 rounded-full text-sm sm:text-sm font-medium shadow-lg flex items-center space-x-2 max-w-[90vw] sm:max-w-md">
        <span class="w-2 h-2 bg-primary rounded-full animate-pulse flex-shrink-0"></span>
        <span id="toast-message" class="truncate">{{ toastMessage }}</span>
      </div>
    </div>
    
    <!-- 新手引导遮罩 (首次进房显示) -->
    <div id="newbie-guide" v-if="showGuide" class="fixed inset-0 bg-black/80 backdrop-blur-md z-[100] flex items-center justify-center p-6 opacity-100 transition-opacity duration-500 pointer-events-auto">
      <div class="bg-surface rounded-3xl p-8 max-w-md text-center shadow-2xl">
        <div class="text-5xl mb-6">🌟</div>
        <h2 class="text-2xl font-bold text-text-main mb-4">欢迎来到小窝！</h2>
        <p class="text-text-muted mb-6">这是一个温馨的在线观影室，让你和远方的亲友一起看电影。</p>
        <div class="space-y-4 mb-8 text-left">
          <div class="flex items-start space-x-3">
            <span class="text-primary mt-1">1.</span>
            <p class="text-text-main">粘贴视频链接到输入框</p>
          </div>
          <div class="flex items-start space-x-3">
            <span class="text-primary mt-1">2.</span>
            <p class="text-text-main">点击「熄灯，开始放映」</p>
          </div>
          <div class="flex items-start space-x-3">
            <span class="text-primary mt-1">3.</span>
            <p class="text-text-main">复制链接邀请亲友加入</p>
          </div>
        </div>
        <button id="guide-close" @click="showGuide = false" class="w-full py-3 sm:py-3 bg-primary text-white font-bold rounded-2xl hover:bg-primary-hover transition-all duration-300 shadow-lg active:scale-95 min-h-[48px] sm:min-h-[44px] text-sm sm:text-base">
          我知道了，开始使用
        </button>
      </div>
    </div>

    <!-- 顶部木质栏 - 响应式优化 -->
    <div id="top-bar" class="absolute top-0 left-0 right-0 px-3 sm:px-6 pt-3 sm:pt-6 pb-2 flex flex-col sm:flex-row gap-2 sm:gap-0 sm:justify-between sm:items-center z-50 transition-all duration-500 transform translate-y-0 opacity-100">
      <!-- 第一行：退出按钮 + 成员列表 (移动端) / 退出按钮 (桌面端) -->
      <div class="flex justify-between items-center">
        <!-- 左侧：退出按钮 -->
        <button id="top-exit-btn" @click="exitRoom" class="bg-surface backdrop-blur-md hover:bg-red-50 hover:text-red-500 transition-all duration-300 p-3 sm:p-2.5 rounded-full shadow-wood border border-white/10 text-text-muted active:scale-95 min-w-[44px] min-h-[44px] sm:min-w-[40px] sm:min-h-[40px]" title="退出房间">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
          </svg>
        </button>

        <!-- 右侧：成员列表 (移动端显示) / 房间信息 (桌面端) -->
        <div class="flex items-center gap-2">
          <!-- 成员列表 (移动端优先显示) -->
          <div class="sm:hidden bg-surface backdrop-blur-md hover:bg-surface/90 transition-all duration-300 px-3 py-2 rounded-full flex items-center shadow-wood border border-white/10 min-w-[80px]">
            <div class="flex -space-x-1.5" id="member-avatars">
              <!-- 只显示前3个头像 -->
              <div v-for="(member, index) in members.slice(0, 3)" :key="member.id" 
                   :class="[memberClass(member), 'z-' + (3 - index)]"
                   class="w-6 h-6 rounded-full border-2 border-surface flex items-center justify-center text-[10px] shadow-sm relative flex-shrink-0">
                <span>{{ member.avatar }}</span>
                <!-- 网络状态指示灯 -->
                <div :style="{ backgroundColor: networkStatusColors[member.status] }"
                     :class="member.status === 'good' ? 'pulse-animation' : ''"
                     class="absolute -top-0.5 -right-0.5 w-2 h-2 rounded-full border border-surface shadow-sm"></div>
              </div>
              <!-- 如果成员超过3个，显示+数字 -->
              <div v-if="members.length > 3" class="w-6 h-6 rounded-full border-2 border-surface bg-surface flex items-center justify-center text-[10px] text-text-muted font-medium shadow-sm flex-shrink-0">
                +{{ members.length - 3 }}
              </div>
            </div>
            <span id="member-count" class="ml-2 text-xs text-text-main font-medium opacity-80 whitespace-nowrap">{{ members.length }}人</span>
          </div>

          <!-- 房间信息 (桌面端显示) -->
          <div id="top-room-info" class="hidden sm:block bg-surface backdrop-blur-md hover:bg-surface/90 transition-all duration-300 px-4 py-2.5 rounded-full flex items-center space-x-3 shadow-wood border border-white/10">
            <div class="flex items-center space-x-1.5">
              <span id="status-dot" :class="statusDotClass" class="w-2 h-2 bg-secondary rounded-full shadow-[0_0_10px_rgba(163,201,168,0.6)]"></span>
              <span id="status-text" class="text-xs text-text-main font-bold tracking-wider">{{ statusText }}</span>
            </div>
            <div class="h-3 w-0.5 bg-text-muted/20"></div>
            <span id="sync-status" class="text-xs text-text-muted font-medium hidden lg:inline">{{ syncStatusText }}</span>
            <div class="h-3 w-0.5 bg-text-muted/20 hidden lg:block"></div>
            <span class="text-xs sm:text-sm text-text-main font-mono font-medium tracking-widest">房间号 {{ roomId }}</span>
            
            <!-- 复制链接按钮 -->
            <button id="share-btn" @click="copyLink" class="ml-1.5 p-2 sm:p-1.5 hover:bg-black/5 rounded-full transition-colors text-text-muted hover:text-primary min-w-[40px] min-h-[40px] flex items-center justify-center" title="复制邀请链接">
              <svg v-if="!isShared" class="w-4 h-4 sm:w-3.5 sm:h-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
              </svg>
              <svg v-else class="w-4 h-4 sm:w-3.5 sm:h-3.5 text-primary" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
              </svg>
            </button>
          </div>
        </div>
      </div>

      <!-- 第二行：房间信息 (移动端) / 成员列表 (桌面端) -->
      <div class="flex justify-center sm:justify-end">
        <!-- 房间信息 (移动端显示) -->
        <div id="top-room-info-mobile" class="sm:hidden bg-surface backdrop-blur-md hover:bg-surface/90 transition-all duration-300 px-3 py-2 rounded-full flex items-center space-x-2 shadow-wood border border-white/10">
          <div class="flex items-center space-x-1.5">
            <span id="status-dot-mobile" :class="statusDotClass" class="w-2 h-2 bg-secondary rounded-full shadow-[0_0_10px_rgba(163,201,168,0.6)]"></span>
            <span id="status-text-mobile" class="text-xs text-text-main font-bold tracking-wider">{{ statusText }}</span>
          </div>
          <div class="h-3 w-0.5 bg-text-muted/20"></div>
          <span class="text-xs text-text-main font-mono font-medium tracking-widest">{{ roomId }}</span>
          
          <!-- 复制链接按钮 (移动端) -->
          <button id="share-btn-mobile" @click="copyLink" class="ml-1.5 p-2 hover:bg-black/5 rounded-full transition-colors text-text-muted hover:text-primary min-w-[40px] min-h-[40px] flex items-center justify-center" title="复制邀请链接">
            <svg v-if="!isShared" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
            </svg>
            <svg v-else class="w-4 h-4 text-primary" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
            </svg>
          </button>
        </div>

        <!-- 成员列表 (桌面端显示) -->
        <div id="top-member-list" class="hidden sm:flex bg-surface backdrop-blur-md hover:bg-surface/90 transition-all duration-300 px-3 py-2 rounded-full items-center shadow-wood border border-white/10 max-w-[200px]">
          <div class="flex -space-x-2" id="member-avatars">
            <!-- 动态生成的成员头像将显示在这里 -->
            <div v-for="(member, index) in members.slice(0, 5)" :key="member.id" 
                 :class="[memberClass(member), 'z-' + (members.length - index)]"
                 class="w-7 h-7 rounded-full border-2 border-surface flex items-center justify-center text-xs shadow-sm relative group cursor-help flex-shrink-0">
              <span>{{ member.avatar }}</span>
              <!-- 网络状态指示灯 -->
              <div :style="{ backgroundColor: networkStatusColors[member.status] }"
                   :class="member.status === 'good' ? 'pulse-animation' : ''"
                   class="absolute -top-1 -right-1 w-2.5 h-2.5 rounded-full border-2 border-surface shadow-sm"></div>
              
              <!-- tooltip -->
              <div class="absolute -bottom-8 left-1/2 -translate-x-1/2 bg-black/70 text-white text-[10px] px-2 py-0.5 rounded opacity-0 group-hover:opacity-100 transition-opacity whitespace-nowrap pointer-events-none">
                {{ member.name }}
              </div>
            </div>
            <!-- 如果成员超过5个，显示+数字 -->
            <div v-if="members.length > 5" class="w-7 h-7 rounded-full border-2 border-surface bg-surface flex items-center justify-center text-xs text-text-muted font-medium shadow-sm flex-shrink-0">
              +{{ members.length - 5 }}
            </div>
          </div>
          <span id="member-count" class="ml-2 text-xs text-text-main font-medium opacity-80 whitespace-nowrap">{{ members.length }}人在线</span>
        </div>
      </div>
    </div>

    <!-- 播放器主体 (全屏沉浸) -->
    <div class="flex-1 flex items-center justify-center relative" @click="toggleControls">
      <!-- 空状态/链接输入：窝窝·待映厅 - 响应式优化 -->
      <div v-if="!isPlaying" id="link-input-section" class="w-full max-w-sm sm:max-w-2xl text-center space-y-6 sm:space-y-10 fade-in px-4 sm:px-6 relative z-20" @click.stop>
        <!-- 氛围背景光晕 -->
        <div class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[140%] sm:w-[120%] h-[140%] sm:h-[120%] bg-primary/5 rounded-full blur-[80px] sm:blur-[100px] pointer-events-none animate-pulse"></div>

        <!-- 主体容器：磨砂票根 -->
        <div class="relative bg-surface/40 backdrop-blur-xl rounded-[24px] sm:rounded-[40px] p-6 sm:p-10 shadow-[0_8px_32px_rgba(0,0,0,0.1)] border border-white/20 overflow-hidden group hover:bg-surface/50 transition-all duration-500">
          <!-- 顶部装饰线 -->
          <div class="absolute top-0 left-0 w-full h-0.5 sm:h-1 bg-gradient-to-r from-transparent via-primary/30 to-transparent opacity-50"></div>

          <!-- 场景插画 -->
          <div class="flex justify-center mb-6 sm:mb-8 space-x-2 sm:space-x-4">
            <span class="text-3xl sm:text-5xl drop-shadow-lg transform -rotate-6 hover:rotate-0 transition-transform duration-300 cursor-default">🛋️</span>
            <span class="text-3xl sm:text-5xl drop-shadow-lg transform translate-y-1 sm:translate-y-2 hover:translate-y-0 transition-transform duration-300 cursor-default">🍿</span>
            <span class="text-3xl sm:text-5xl drop-shadow-lg transform rotate-6 hover:rotate-0 transition-transform duration-300 cursor-default">🕯️</span>
          </div>

          <!-- 文案 -->
          <div class="space-y-2 sm:space-y-3 mb-6 sm:mb-10">
            <h2 class="text-lg sm:text-2xl text-text-main font-bold tracking-wide">窝窝已搭好，等一部好戏。</h2>
          </div>
          
          <!-- 输入体验：聚光灯槽 -->
          <div class="relative group/input">
            <div class="absolute -inset-0.5 bg-gradient-to-r from-primary/20 to-secondary/20 rounded-xl sm:rounded-2xl blur opacity-0 group-hover/input:opacity-100 transition duration-500"></div>
            <div class="relative bg-surface rounded-xl sm:rounded-2xl p-1.5 sm:p-2 flex flex-col sm:flex-row items-center shadow-inner-warm border border-white/10 transition-all duration-300 focus-within:ring-2 focus-within:ring-primary/30 focus-within:bg-white/80">
              <!-- 移动端：按钮在上方，桌面端：按钮在右侧 -->
              <button 
                id="start-play-btn"
                @click="startPlay" 
                class="w-full sm:w-auto mb-2 sm:mb-0 sm:ml-3 px-4 sm:px-6 py-2.5 sm:py-3 bg-primary text-white rounded-lg sm:rounded-xl hover:bg-primary-hover transition-all duration-300 font-bold text-sm sm:text-sm shadow-lg active:scale-95 whitespace-nowrap flex items-center justify-center space-x-2 min-h-[44px] sm:min-h-[40px]">
                <span>✨ 熄灯，开始放映</span>
              </button>
              
              <div class="relative flex-1 w-full sm:w-auto">
                <input 
                  id="video-url-input"
                  v-model="videoUrl"
                  type="text" 
                  placeholder="把影片链接贴在这里，就像把光盘放入播放机" 
                  class="w-full bg-transparent border-none text-text-main placeholder-text-muted/50 px-3 sm:px-4 py-2.5 sm:py-3 focus:ring-0 focus:outline-none text-sm sm:text-base font-medium min-h-[44px] sm:min-h-[40px]"
                  @input="updateUrlStatus"
                  @blur="validateUrlOnBlur"
                >
                <div id="url-status-indicator" :class="urlStatusClass" class="absolute right-3 top-1/2 -translate-y-1/2 w-2.5 h-2.5 rounded-full"></div>
              </div>
            </div>
          </div>
          
          <!-- 支持格式提示 -->
          <p class="text-xs text-text-muted/50 mt-3 sm:mt-2">支持 MP4 / m3u8 链接</p>
          
          <!-- 底部装饰：印章 -->
          <div class="absolute -bottom-4 sm:-bottom-6 -right-4 sm:-right-6 w-16 sm:w-24 h-16 sm:h-24 bg-text-muted/5 rounded-full blur-xl"></div>
        </div>

        <!-- 快速测试源 (响应式布局) -->
        <div class="flex flex-col items-center space-y-3 sm:space-y-4 opacity-60 hover:opacity-100 transition-opacity duration-300">
          <div class="flex flex-wrap justify-center gap-2 sm:gap-3">
            <button 
              v-for="(source, index) in demoSources" 
              :key="index"
              @click="selectDemoSource(source.url, source.title)" 
              class="px-3 sm:px-3 py-2 sm:py-1.5 bg-white/10 hover:bg-white/20 border border-white/10 rounded-full text-xs text-text-muted/80 hover:text-primary transition-all cursor-pointer backdrop-blur-sm min-h-[40px] min-w-[48px] sm:min-h-[36px] sm:min-w-[44px] text-center">
              🎬 {{ source.title }}
            </button>
          </div>
        </div>
      </div>

      <!-- 视频容器 -->
      <div v-else id="player-section" class="absolute inset-0 bg-black z-0">
        <!-- Artplayer 挂载点 -->
        <div id="artplayer-container" class="w-full h-full"></div>
      </div>
    </div>

    <!-- Autoplay 引导覆盖层 (默认隐藏) - 响应式优化 -->
    <div id="autoplay-overlay" v-if="showAutoplayGuide" class="absolute inset-0 z-[70] bg-black/60 backdrop-blur-sm flex flex-col items-center justify-center cursor-pointer px-4" @click="handleAutoplayClick">
      <div class="bg-surface/20 p-4 sm:p-6 rounded-full mb-3 sm:mb-4 animate-pulse">
        <svg class="w-10 h-10 sm:w-12 sm:h-12 text-primary" fill="currentColor" viewBox="0 0 20 20">
          <path d="M6.3 2.841A1.5 1.5 0 004 4.11V15.89a1.5 1.5 0 002.3 1.269l9.344-5.89a1.5 1.5 0 000-2.538L6.3 2.84z" />
        </svg>
      </div>
      <p class="text-white text-base sm:text-lg font-medium tracking-wide text-center">点击屏幕，加入放映</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'

// 扩展window对象类型
declare global {
  interface Window {
    Artplayer?: any;
    Hls?: any;
  }
}

// 定义类型
interface Member {
  id: string;
  name: string;
  isOwner: boolean;
  isMe: boolean;
  avatar: string;
  status: 'good' | 'medium' | 'poor' | 'offline';
}

interface SyncLevel {
  text: string;
  color: string;
}

// 路由
const route = useRoute();
const router = useRouter();

// 响应式状态
const toastVisible = ref(false);
const toastMessage = ref('');
const showGuide = ref(true);
const controlsVisible = ref(true);
const isPlaying = ref(false);
const videoUrl = ref('');
const showAutoplayGuide = ref(false);
const isShared = ref(false);

// 房间信息
const roomId = computed(() => (route.params.roomId as string) || '888888');
const statusText = ref('等待中');
const syncStatusText = ref('同步精度：完美');
const statusDotClass = computed(() => {
  if (isPlaying.value) {
    return 'w-2.5 h-2.5 bg-secondary pulse-animation';
  }
  return 'w-2.5 h-2.5 bg-gray-400';
});

// 网络状态颜色
const networkStatusColors: Record<string, string> = {
  good: '#A3C9A8',
  medium: '#FFD93D',
  poor: '#FF6B6B',
  offline: '#8D7B68'
};

// URL状态类
const urlStatusClass = computed(() => {
  if (!videoUrl.value) {
    return 'bg-gray-300';
  }
  if (validateUrl(videoUrl.value) === 'valid') {
    return 'bg-secondary pulse-animation';
  }
  return 'bg-red-500';
});

// 随机头像生成器
const randomAvatars = ['🐱', '🐶', '🐰', '🦊', '🐻', '🐼', '🐨', '🐯', '🦁', '🐮', '🐷', '🐸', '🐵', '🐔', '🐧', '🐦', '🐤', '🐺', '🐗', '🐴'];

// 生成随机网络状态
function getRandomNetworkStatus(): 'good' | 'medium' | 'poor' | 'offline' {
  const statuses: Array<'good' | 'medium' | 'poor' | 'offline'> = ['good', 'medium', 'poor', 'offline'];
  return statuses[Math.floor(Math.random() * statuses.length)];
}

// 生成随机成员列表
function generateRandomMembers(): Member[] {
  const memberCount = Math.floor(Math.random() * 6) + 2;
  const members: Member[] = [];
  
  members.push({ id: 'owner', name: '房主', isOwner: true, isMe: false, avatar: '👑', status: 'good' });
  members.push({ id: 'me', name: '我', isOwner: false, isMe: true, avatar: randomAvatars[Math.floor(Math.random() * randomAvatars.length)], status: 'good' });
  
  for (let i = 2; i < memberCount; i++) {
    members.push({
      id: `member-${i}`,
      name: `访客${i}`,
      isOwner: false,
      isMe: false,
      avatar: randomAvatars[Math.floor(Math.random() * randomAvatars.length)],
      status: getRandomNetworkStatus()
    });
  }
  
  return members;
}

// 成员列表
const members = ref<Member[]>(generateRandomMembers());

// 成员样式类
function memberClass(member: Member) {
  if (member.isOwner) {
    return 'bg-primary text-white';
  } else if (member.isMe) {
    return 'bg-secondary text-white font-bold';
  } else {
    return 'bg-surface text-text-main font-medium';
  }
}

// 演示视频源
const demoSources = [
  { title: '大闹天宫', url: 'https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/BigBuckBunny.mp4' },
  { title: '奇幻旅程', url: 'https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ElephantsDream.mp4' },
  { title: '速度激情', url: 'https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerBlazes.mp4' },
  { title: '冒险之旅', url: 'https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/ForBiggerEscapes.mp4' },
  { title: '最后的片段', url: 'https://commondatastorage.googleapis.com/gtv-videos-bucket/sample/Sintel.mp4' }
];

// URL校验
const supportedFormats = ['.mp4', '.m3u8', '.webm', '.ogg', '.mov', '.avi'];

// 错误类型枚举
const ERROR_TYPES = {
  URL_EMPTY: '请先输入视频链接',
  URL_INVALID: '请输入有效的视频链接，支持 MP4 / m3u8 等格式',
  ROOM_FULL: '房间已满，请稍后再试或创建新房间',
  ROOM_EXPIRED: '房间已过期，请创建新房间',
  NETWORK_ERROR: '网络连接异常，请检查网络后重试',
  PLAY_ERROR: '视频播放失败，请更换链接',
  COPY_FAILED: '复制失败，请手动复制链接',
  SHARE_FAILED: '分享失败，请尝试手动分享',
  UNKNOWN_ERROR: '操作失败，请稍后重试'
};

function validateUrl(url: string): 'valid' | 'invalid' | 'empty' {
  if (!url.trim()) {
    return 'empty';
  }
  
  try {
    new URL(url);
  } catch (error) {
    return 'invalid';
  }
  
  const lowerUrl = url.toLowerCase();
  for (const format of supportedFormats) {
    if (lowerUrl.includes(format)) {
      return 'valid';
    }
  }
  
  return 'invalid';
}

// 更新URL状态
function updateUrlStatus() {
  // 这会在UI中自动更新，通过urlStatusClass计算属性
}

// 失焦时校验URL
function validateUrlOnBlur() {
  if (validateUrl(videoUrl.value) === 'invalid' && videoUrl.value) {
    showToast(ERROR_TYPES.URL_INVALID, 'error');
  }
}

// 显示Toast
function showToast(message: string, type: 'success' | 'warning' | 'error' = 'success') {
  toastMessage.value = message;
  toastVisible.value = true;
  
  setTimeout(() => {
    toastVisible.value = false;
  }, 3000);
}

// 选择演示源
function selectDemoSource(url: string, title: string) {
  videoUrl.value = url;
  showToast(`已选择视频：${title}`, 'success');
}

// 复制链接功能 (支持移动端原生分享)
async function copyLink() {
  // 获取房间号（当前是固定的，实际应从URL或后端获取）
  const roomIdValue = roomId.value;
  const baseUrl = window.location.origin;
  const shareUrl = `${baseUrl}/room/${roomIdValue}`;
  
  const shareData = {
    title: '小窝',
    text: '我在[小窝]建了个温馨小窝，快来一起看电影吧！',
    url: shareUrl
  };

  // 优先尝试原生分享 (移动端)
  if (navigator.share && /Mobi|Android/i.test(navigator.userAgent)) {
    try {
      await navigator.share(shareData);
      return; // 分享成功后无需 Toast，系统会有反馈
    } catch (err) {
      console.log('Share failed:', err);
      // 降级到复制链接
    }
  }

  // 降级方案：复制到剪贴板
  try {
    await navigator.clipboard.writeText(`${shareData.text} 房间号：${roomIdValue}，点击加入：${shareData.url}`);
    isShared.value = true;
    showToast('✨ 邀请函已复制，快去发给 Ta 吧！', 'success');
    
    setTimeout(() => {
      isShared.value = false;
    }, 2000);
  } catch (err) {
    showToast(ERROR_TYPES.COPY_FAILED, 'error');
  }
}

// 开始播放
function startPlay() {
  const url = videoUrl.value.trim();
  if (!url) {
    showToast(ERROR_TYPES.URL_EMPTY, 'warning');
    return;
  }
  
  if (validateUrl(url) === 'invalid') {
    showToast(ERROR_TYPES.URL_INVALID, 'error');
    return;
  }
  
  isPlaying.value = true;
  
  // 初始化Artplayer
  initArtplayer(url);
  
  showToast('准备就绪，开始放映', 'success');
}

// 初始化Artplayer
let art: any = null;
function initArtplayer(url: string) {
  if (art) {
    art.destroy();
  }
  
  // 动态加载Artplayer和Hls.js
  const loadScripts = async () => {
    // 加载Artplayer
    if (!window.Artplayer) {
      const artplayerScript = document.createElement('script');
      artplayerScript.src = 'https://cdn.jsdelivr.net/npm/artplayer/dist/artplayer.js';
      artplayerScript.async = true;
      document.head.appendChild(artplayerScript);
      await new Promise((resolve) => artplayerScript.onload = resolve);
    }
    
    // 加载Hls.js
    if (!window.Hls) {
      const hlsScript = document.createElement('script');
      hlsScript.src = 'https://cdn.jsdelivr.net/npm/hls.js@latest';
      hlsScript.async = true;
      document.head.appendChild(hlsScript);
      await new Promise((resolve) => hlsScript.onload = resolve);
    }
    
    // 初始化Artplayer
    const Artplayer = window.Artplayer;
    const Hls = window.Hls;
    
    if (Artplayer && Hls) {
      art = new Artplayer({
        container: '#artplayer-container',
        url: url,
        theme: '#FF9F76',
        volume: 0.5,
        isLive: false,
        muted: false,
        autoplay: true,
        pip: true,
        autoSize: true,
        autoMini: true,
        screenshot: true,
        setting: true,
        loop: true,
        flip: true,
        playbackRate: true,
        aspectRatio: true,
        fullscreen: true,
        fullscreenWeb: true,
        subtitleOffset: true,
        miniProgressBar: false,
        mutex: true,
        backdrop: true,
        playsInline: true,
        autoPlayback: true,
        airplay: true,
        // HLS 支持
        customType: {
          m3u8: function (video: HTMLVideoElement, m3u8Url: string) {
            if (Hls.isSupported()) {
              const hls = new Hls();
              hls.loadSource(m3u8Url);
              hls.attachMedia(video);
            } else if (video.canPlayType('application/vnd.apple.mpegurl')) {
              video.src = m3u8Url;
            } else {
              if (art && art.notice) {
                art.notice.show = '不支持播放此视频格式';
              }
            }
          },
        },
        icons: {
          state: '<svg width="48" height="48" viewBox="0 0 48 48" fill="none"><path d="M16 12v24l20-12-20-12z" fill="#ffffff"/></svg>',
        },
      });
      
      art.on('play', () => {
        console.log('开始播放');
      });
      
      art.on('pause', () => {
        console.log('暂停播放');
      });
    }
  };
  
  loadScripts();
}

// 切换控制栏显示
function toggleControls() {
  controlsVisible.value = !controlsVisible.value;
}

// 退出房间
function exitRoom() {
  router.push('/');
}

// 自动播放点击处理
function handleAutoplayClick() {
  showAutoplayGuide.value = false;
  if (art) {
    art.play();
  }
}

// 模拟同步状态变化
const SYNC_LEVELS: Record<string, SyncLevel> = {
  perfect: { text: '同步精度：完美', color: '#A3C9A8' },
  gentle: { text: '同步精度：微调', color: '#FFD93D' },
  ghost: { text: '同步精度：追赶', color: '#FF9F76' },
  forced: { text: '同步精度：对齐', color: '#FF6B6B' }
};

function updateSyncStatus() {
  const levels = Object.keys(SYNC_LEVELS);
  const randomLevel = levels[Math.floor(Math.random() * levels.length)];
  const syncData = SYNC_LEVELS[randomLevel];
  syncStatusText.value = syncData.text;
}

// 组件挂载时初始化
onMounted(() => {
  // 检查是否首次进房
  const hasVisited = localStorage.getItem('xiaowo_visited')
  if (!hasVisited) {
    // 显示引导遮罩
    setTimeout(() => {
      showGuide.value = true
    }, 1000)
    
    // 标记为已访问
    localStorage.setItem('xiaowo_visited', 'true')
  } else {
    // 非首次访问，不显示引导
    showGuide.value = false
  }
  
  // 模拟定期更新成员列表
  setInterval(() => {
    members.value = generateRandomMembers()
  }, 10000)
  
  // 每5秒模拟一次同步状态变化
  setInterval(updateSyncStatus, 5000)
  
  // 显示自动播放引导
  showAutoplayGuide.value = true
});

// 组件卸载时清理
onBeforeUnmount(() => {
  if (art) {
    art.destroy()
  }
})
</script>

<style scoped>
/* Artplayer 自定义样式 - 温暖居家风 */
.art-video-player .art-bottom {
  /* 强制定位 */
  position: absolute !important;
  bottom: 24px !important;
  top: auto !important; /* 必须重置 top，防止拉伸 */
  height: auto !important;
  max-height: 120px !important; /* 进一步收窄 */
  left: 50% !important;
  transform: translateX(-50%) !important;
  width: 90% !important;
  max-width: 640px !important; /* 进一步收窄 */
  
  /* 视觉风格 */
  background: rgba(247, 238, 221, 0.85) !important;
  backdrop-filter: blur(12px) !important;
  -webkit-backdrop-filter: blur(12px) !important;
  border: 1px solid rgba(255, 255, 255, 0.3) !important;
  border-radius: 20px !important;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.2) !important;
  
  /* 内部布局 */
  padding: 10px 20px !important;
  display: flex !important;
  flex-direction: column !important;
  justify-content: center !important;
  align-items: center !important;
  gap: 6px !important;
  
  pointer-events: auto !important;
  transition: all 0.3s ease !important;
  z-index: 50 !important;
}

/* 悬浮时稍微上浮 */
.art-video-player .art-bottom:hover {
  transform: translateX(-50%) translateY(-2px) !important;
  box-shadow: 0 15px 50px rgba(0, 0, 0, 0.3) !important;
}

/* 进度条容器 - 整合在胶囊内部 */
.art-video-player .art-progress {
  position: relative !important;
  width: 100% !important;
  height: 6px !important;
  bottom: auto !important;
  left: auto !important;
  border-radius: 3px !important;
  background: rgba(74, 59, 50, 0.1) !important; /* 深色轨道，对比度更高 */
  overflow: visible !important;
  margin: 0 !important;
  padding: 0 !important;
  order: 1 !important; /* 确保在上方 */
}

/* 进度条播放部分 */
.art-video-player .art-control-progress-played {
  background: #FF9F76 !important; /* primary */
  border-radius: 3px !important;
}

/* 进度条缓冲部分 */
.art-video-player .art-control-progress-loaded {
  background: rgba(255, 159, 118, 0.2) !important;
  border-radius: 3px !important;
}

/* 鼠标悬停/拖拽时的幻影进度条 (修复黑色横线问题) */
.art-video-player .art-control-progress-hover {
  background: rgba(74, 59, 50, 0.08) !important; /* 极淡的深色 */
  border-radius: 3px !important;
  height: 100% !important;
  border: none !important; /* 确保无边框 */
  transform: scaleY(1) !important; /* 防止被默认样式压缩 */
}

/* 彻底隐藏缩略图容器 (防止出现黑线) */
.art-video-player .art-thumb {
  display: none !important;
  width: 0 !important;
  height: 0 !important;
  opacity: 0 !important;
  visibility: hidden !important;
  background: none !important;
  border: none !important;
  box-shadow: none !important;
  pointer-events: none !important;
}

/* 隐藏高亮标记 */
.art-video-player .art-highlight {
  display: none !important;
}

/* 拖拽球 */
.art-video-player .art-control-progress-indicator {
  width: 14px !important;
  height: 14px !important;
  background: #fff !important;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.2) !important;
  border: 2px solid #FF9F76 !important; /* 增加边框，更清晰 */
  border-radius: 50% !important;
  top: 50% !important;
  transform: translateY(-50%) !important;
  display: block !important;
}

/* 控制栏 (按钮区域) */
.art-video-player .art-controls {
  background: transparent !important;
  border: none !important;
  box-shadow: none !important;
  backdrop-filter: none !important;
  width: 100% !important;
  height: 36px !important;
  padding: 0 !important;
  margin: 0 !important;
  display: flex !important;
  align-items: center !important;
  order: 2 !important; /* 确保在下方 */
}

/* 图标颜色 */
.art-video-player .art-icon path {
  fill: #4A3B32 !important; /* text-main */
  transition: fill 0.2s !important;
}

.art-video-player .art-icon:hover path {
  fill: #FF9F76 !important; /* primary-hover */
}

.art-video-player .art-control-time {
  color: #4A3B32 !important;
  font-family: 'Inter', monospace !important;
  font-weight: 600 !important;
  font-size: 13px !important;
  margin: 0 10px !important;
}

/* 隐藏默认的背景渐变 */
.art-video-player .art-mask {
  background: none !important;
}

/* 悬浮提示框样式 (时间气泡) */
.art-video-player .art-tip {
  background: rgba(62, 50, 40, 0.9) !important;
  border-radius: 8px !important;
  padding: 6px 10px !important;
  font-size: 12px !important;
  color: #fff !important;
  box-shadow: 0 4px 12px rgba(0,0,0,0.2) !important;
  border: 1px solid rgba(255,255,255,0.1) !important;
  bottom: 20px !important; /* 稍微抬高，避免遮挡进度条 */
}

/* 隐藏提示框的小箭头 */
.art-video-player .art-tip::after {
  display: none !important;
}

/* Loading 状态优化 - 拒绝死黑，使用磨砂暖咖色 */
.art-video-player .art-layer-loading {
  background: rgba(44, 36, 31, 0.4) !important; /* dark.bg with opacity */
  backdrop-filter: blur(20px) !important;
  -webkit-backdrop-filter: blur(20px) !important;
}

/* Loading 图标颜色 */
.art-video-player .art-loading-indicator {
  color: #FF9F76 !important; /* primary */
}

/* 动画效果定义 */
@keyframes pulse-animation {
  0%, 100% {
    opacity: 1;
    transform: scale(1);
  }
  50% {
    opacity: 0.7;
    transform: scale(1.1);
  }
}

@keyframes fade-in {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 应用动画类 */
.pulse-animation {
  animation: pulse-animation 2s infinite;
}

.fade-in {
  animation: fade-in 0.5s ease-out;
}

/* 确保样式正确应用 */
:deep(.pulse-animation) {
  animation: pulse-animation 2s infinite;
}

:deep(.fade-in) {
  animation: fade-in 0.5s ease-out;
}
</style>