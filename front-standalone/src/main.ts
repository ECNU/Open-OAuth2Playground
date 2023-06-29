import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import 'highlight.js/styles/atom-one-light.css';
// import 'highlight.js/styles/mono-blue.css';
import 'highlight.js/lib/common';
import hljsVuePlugin from '@highlightjs/vue-plugin';

const app = createApp(App);
app.use(ElementPlus);
app.use(router);
app.use(store);
app.use(hljsVuePlugin);
app.mount('#app');