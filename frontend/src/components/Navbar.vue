<template>
  <nav class="bg-teal-500 p-4">
    <div class="max-w-7xl mx-auto flex justify-between items-center">
      <div class="flex items-center">
        <router-link to="/" class="text-white text-2xl font-semibold">Bot Optimasi</router-link>
      </div>
      <div class="hidden md:flex items-center space-x-4">
        <span class="text-white hover:text-gray-300">{{ firstName }} {{ lastName }}</span>
        <router-link to="/settings" class="text-white hover:text-gray-300">Settings</router-link>
        <button @click="handleLogout" class="text-white hover:text-gray-300">Logout</button>
      </div>
      <div class="md:hidden flex items-center">
        <button @click="toggleMobileMenu" class="text-white focus:outline-none">
          <svg class="h-6 w-6" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
            <path d="M4 6H20M4 12H20M4 18H20" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
          </svg>
        </button>
      </div>
    </div>
    <div v-if="isMobileMenuOpen" class="md:hidden absolute top-0 left-0 w-full bg-gray-800">
      <div class="flex flex-col items-center py-4 space-y-4">
        <router-link to="/" class="text-white hover:text-gray-300">Home</router-link>
        <router-link to="/settings" class="text-white hover:text-gray-300">Settings</router-link>
        <button @click="handleLogout" class="text-white hover:text-gray-300">Logout</button>
      </div>
    </div>
  </nav>
</template>

<script>
export default {
  name: 'NavbarBot',
  props: ['firstName', 'lastName'],
  data() {
    return {
      isMobileMenuOpen: false
    };
  },
  methods: {
    toggleMobileMenu() {
      this.isMobileMenuOpen = !this.isMobileMenuOpen;
    },
    handleLogout() {
      localStorage.removeItem('token'); // Hapus token dari local storage saat logout
      this.$router.push('/login'); // Redirect ke halaman login setelah logout
    }
  }
};
</script>

<style scoped>
/* Tambahkan styling khusus jika diperlukan */
</style>
