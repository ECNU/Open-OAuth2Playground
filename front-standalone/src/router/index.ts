import { createRouter, createWebHistory } from 'vue-router';
import Layout from '/@/views/Layout.vue';

const routes = [
    {
        path: '/',
        name: 'Layout',
        component: Layout,
        meta: {title: 'OO2P'}
    }
];

const router = createRouter({
    history: createWebHistory("/playground/"),
    routes: routes
});

router.beforeEach((to, from, next) => {
    /* 路由发生变化修改页面title */
    if(to.meta?.title){
        document.title = to.meta?.title as string;
    }
    next();
});

export default router;
