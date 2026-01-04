import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import MainLayout from '@/layouts/MainLayout.vue'
import RoomLayout from '@/layouts/RoomLayout.vue'

const routes: RouteRecordRaw[] = [
  {
    path: '/',
    name: 'Home',
    component: () => import('@/views/HomeView.vue'),
    meta: {
      title: '首页 - 小窝同步观影',
      requiresAuth: false,
    },
  },
  {
    path: '/app',
    component: MainLayout,
    children: [
      // 这里可以放置其他需要MainLayout的页面，比如个人中心等
    ],
  },
  {
    path: '/room/:roomId?',
    name: 'Room',
    component: () => import('@/views/RoomPage.vue'),
    meta: {
      title: '房间 - 小窝同步观影',
      requiresAuth: true,
    },
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/NotFound.vue'),
    meta: {
      title: '页面未找到 - 小窝同步观影',
    },
  },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    } else {
      return { top: 0 }
    }
  },
})

import { useUserStore } from '@/store/modules/user'

// 路由守卫
router.beforeEach((to, from, next) => {
  // 设置页面标题
  if (to.meta.title) {
    document.title = to.meta.title as string
  }
  
  // 权限检查
  if (to.meta.requiresAuth) {
    // 暂时允许所有用户访问，后续添加完整的登录逻辑
    // const userStore = useUserStore()
    // if (!userStore.isLoggedIn) {
    //   next('/')
    //   return
    // }
  }
  
  next()
})

export default router