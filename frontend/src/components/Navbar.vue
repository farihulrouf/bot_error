<template>
  <nav class="bg-[#128c7e] p-4">
    <div class="max-w-7xl mx-auto flex justify-between items-center">
      <div class="flex items-center">
        <router-link to="/" class="text-white text-md flex items-center">
          <svg-icon type="mdi" :path="path"></svg-icon>
          <span class="text-xl">ptimasi bot</span>
        </router-link>
      </div>
      <div class="hidden md:flex items-center space-x-4 mr-10">
        <DropDown
          :menuItems="menuItems"
          :profile="profile"
          :additionalText="additionalText"
          @menu-click="handleMenuClick"
           :bgColor="bgColor"
        />
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
      </div>
    </div>
  </nav>
</template>

<script>
import { mdiMenuDown } from '@mdi/js';
import DropDown from './DropDown.vue'; // Pastikan jalur ini sesuai dengan struktur proyek Anda

import SvgIcon from '@jamescoyle/vue-icon';
import { mdiWhatsapp } from '@mdi/js';

export default {
  name: 'NavbarBot',
  components: {
    DropDown,
    SvgIcon
  },
  props: {
    firstName: {
      type: String,
      required: true,
    },
    bgColor: {
      type: String,
      default: 'flex justify-center absolute left-4 text-white',
      // Warna latar belakang default (jika tidak disediakan oleh parent)
    },
    lastName: {
      type: String,
      required: true,
    },
    additionalText: {
      type: String,
      default: '', // Default jika tidak ada teks tambahan yang dilempar dari parent
    },
    menuItems: {
      type: Array,
      default: () => [
       { name: 'profile', action: 'profile' },
        { name: 'Settings', action: 'settings' },
        { name: 'Logout', action: 'logout' },
      ],
    },
  },
  data() {
    return {
      profile: {
        name: 'Profile', // Ganti dengan data profil yang sesuai
        // Tambahkan properti lain sesuai kebutuhan, seperti email, avatar, dll.
        comp: mdiMenuDown, // Contoh nilai untuk properti comp
      },
      path: mdiWhatsapp,
      isMobileMenuOpen: false,
    };
  },
  methods: {
    toggleMobileMenu() {
      this.isMobileMenuOpen = !this.isMobileMenuOpen;
    },
    handleLogout() {
      localStorage.removeItem('token'); // Hapus token dari local storage saat logout
      this.$router.push('/login'); // Redirect ke halaman login setelah logout
    },
    handleMenuClick(action) {
      if (action === 'settings') {
        this.$router.push('/settings');
      } else if (action === 'logout') {
        this.handleLogout();
      }
      // Menangani aksi lainnya jika diperlukan
    },
  },
};
</script>

<style scoped>
/* Tambahkan styling khusus jika diperlukan */
</style>
