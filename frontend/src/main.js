import { createApp } from 'vue';
import App from './App.vue';
import router from './router'; // Import Vue Router instance
import Notifications from '@kyvg/vue3-notification'; // Import Notifications plugin
import './assets/css/index.css'; 
const app = createApp(App);

// Gunakan plugin Notifications
app.use(Notifications);

app.use(router).mount('#app');
