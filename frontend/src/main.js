import { createApp } from 'vue';
import App from './App.vue';
import router from './router'; // Import Vue Router instance
import './assets/css/index.css'; // Pastikan pathnya sesuai dengan letak file Anda

createApp(App).use(router).mount('#app');
