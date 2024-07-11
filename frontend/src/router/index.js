import { createRouter, createWebHistory } from 'vue-router';
import Login from '../views/LoginBot.vue'; // Import your login component
import Dashboard from '../views/DashboardBot.vue'; // Import your dashboard component

const routes = [
  { path: '/', redirect: '/dashboard' }, // Redirect to dashboard if no path matches
  { path: '/login', component: Login, meta: { requiresGuest: true } },
  { path: '/dashboard', component: Dashboard, meta: { requiresAuth: true } },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

router.beforeEach((to, from, next) => {
  const isLoggedIn = !!localStorage.getItem('token'); // Memeriksa token di localStorage

  if (to.meta.requiresAuth && !isLoggedIn) {
    next('/login'); // Redirect ke login jika mencoba mengakses halaman yang memerlukan login
  } else if (to.meta.requiresGuest && isLoggedIn) {
    next('/dashboard'); // Redirect ke dashboard jika sudah login dan mencoba mengakses halaman login
  } else {
    next(); // Lanjutkan navigasi
  }
});

export default router;
