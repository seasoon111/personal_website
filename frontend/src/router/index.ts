import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(),
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) return savedPosition
    if (to.hash) return { el: to.hash, behavior: 'smooth' }
    return { top: 0 }
  },
  routes: [
    {
      path: '/',
      component: () => import('@/layouts/DefaultLayout.vue'),
      children: [
        { path: '', name: 'home', component: () => import('@/pages/Home.vue') },
        { path: 'articles', name: 'articles', component: () => import('@/pages/Articles.vue') },
        { path: 'articles/:slug', name: 'article-detail', component: () => import('@/pages/ArticleDetail.vue') },
        { path: 'archives', name: 'archives', component: () => import('@/pages/Archives.vue') },
        { path: 'projects', name: 'projects', component: () => import('@/pages/Projects.vue') },
        { path: 'projects/:slug', name: 'project-detail', component: () => import('@/pages/ProjectDetail.vue') },
        { path: 'moments', name: 'moments', component: () => import('@/pages/Moments.vue') },
        { path: 'about', name: 'about', component: () => import('@/pages/About.vue') },
      ],
    },
    { path: '/admin/login', name: 'admin-login', component: () => import('@/pages/admin/Login.vue') },
    {
      path: '/admin',
      component: () => import('@/layouts/AdminLayout.vue'),
      meta: { requiresAuth: true },
      children: [
        { path: '', name: 'admin-dashboard', component: () => import('@/pages/admin/Dashboard.vue') },
        { path: 'articles', name: 'admin-articles', component: () => import('@/pages/admin/Articles.vue') },
        { path: 'articles/new', name: 'admin-article-new', component: () => import('@/pages/admin/ArticleEdit.vue') },
        { path: 'articles/:id/edit', name: 'admin-article-edit', component: () => import('@/pages/admin/ArticleEdit.vue') },
        { path: 'categories', name: 'admin-categories', component: () => import('@/pages/admin/Categories.vue') },
        { path: 'tags', name: 'admin-tags', component: () => import('@/pages/admin/Tags.vue') },
        { path: 'projects', name: 'admin-projects', component: () => import('@/pages/admin/Projects.vue') },
        { path: 'projects/new', name: 'admin-project-new', component: () => import('@/pages/admin/ProjectEdit.vue') },
        { path: 'projects/:id/edit', name: 'admin-project-edit', component: () => import('@/pages/admin/ProjectEdit.vue') },
        { path: 'moments', name: 'admin-moments', component: () => import('@/pages/admin/Moments.vue') },
        { path: 'about', name: 'admin-about', component: () => import('@/pages/admin/About.vue') },
        { path: 'images', name: 'admin-images', component: () => import('@/pages/admin/Images.vue') },
        { path: 'messages', name: 'admin-messages', component: () => import('@/pages/admin/Messages.vue') },
      ],
    },
    { path: '/:pathMatch(.*)*', name: 'not-found', component: () => import('@/pages/NotFound.vue') },
  ],
})

router.beforeEach((to) => {
  if (to.meta.requiresAuth && !localStorage.getItem('pw_token')) {
    return { name: 'admin-login', query: { redirect: to.fullPath } }
  }
  return true
})

export default router
