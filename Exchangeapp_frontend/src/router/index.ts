import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router';
import HomeView from '../views/HomeView.vue';
import NewsView from '../views/NewsView.vue';
import NewsDetailView from '../views/NewsDetailView.vue';
import Login from '../components/Login.vue';
import Register from '../components/Register.vue';
import VideoUpload from '../components/VideoUpload.vue';
import Test from '../components/Test.vue';
import History from '../components/History.vue'; // 导入新组件

const routes: RouteRecordRaw[] = [
  { path: '/', name: 'Home', component: HomeView },
  { path: '/news', name: 'News', component: NewsView },
  { path: '/news/:id', name: 'NewsDetail', component: NewsDetailView },
  { path: '/login', name: 'Login', component: Login },
  { path: '/register', name: 'Register', component: Register },
  { path: '/video-upload', name: 'VideoUpload', component: VideoUpload , meta: { keepAlive: true }},
  { path: '/videoget', name: 'Test', component: Test },
  { path: '/videoHistory', name: 'History', component: History,meta: { keepAlive: true } }, // 添加新路由
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;

