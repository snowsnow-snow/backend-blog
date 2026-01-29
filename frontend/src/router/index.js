import { createRouter, createWebHashHistory } from 'vue-router'

const router = createRouter({
    history: createWebHashHistory(import.meta.env.BASE_URL),
    scrollBehavior(to, from, savedPosition) {
        if (savedPosition) {
            return savedPosition
        }
        // 如果去往首页或详情页，都不自动滚动，交由组件手动处理或保持现状
        if (to.name === 'home' || to.name === 'post-detail') {
            return false
        }
        return { top: 0 }
    },
    routes: [
        {
            path: '/',
            component: () => import('../layouts/PublicLayout.vue'),
            children: [
                {
                    path: '',
                    name: 'home',
                    component: () => import('../views/HomeView.vue'),
                    meta: { keepAlive: true }
                },
                {
                    path: 'post/:id',
                    name: 'post-detail',
                    component: () => import('../views/PostDetailView.vue'),
                    props: true
                }
            ]
        },
        {
            path: '/login',
            name: 'login',
            component: () => import('../views/LoginView.vue')
        },
        {
            path: '/admin',
            component: () => import('../layouts/AdminLayout.vue'),
            meta: { requiresAuth: true },
            children: [
                {
                    path: '',
                    name: 'admin-dashboard',
                    component: () => import('../views/admin/DashboardView.vue')
                },
                {
                    path: 'posts/new',
                    name: 'post-create',
                    component: () => import('../views/admin/EditorView.vue')
                },
                {
                    path: 'posts/:id/edit',
                    name: 'post-edit',
                    component: () => import('../views/admin/EditorView.vue'),
                    props: true
                },
                {
                    path: 'categories',
                    name: 'admin-categories',
                    component: () => import('../views/admin/CategoryView.vue')
                }
            ]
        }
    ]
})

// Navigation guard for auth
router.beforeEach((to, from, next) => {
    const token = localStorage.getItem('token')
    if (to.meta.requiresAuth && !token) {
        next({ name: 'login' })
    } else {
        next()
    }
})

export default router
