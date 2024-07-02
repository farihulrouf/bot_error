import { createRouter, createWebHistory } from 'vue-router';
import Login from '../views/LoginBot.vue'; // Import your login component
import Dashboard from '../views/DashboardBot.vue'; // Import your dashboard component

const routes = [
  { path: '/', redirect: '/dashboard' }, // Redirect to login if no path matches
  { path: '/login', component: Login },
  { path: '/dashboard', component: Dashboard, meta: { requiresAuth: true } },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router; // Ensure you export the router instance
