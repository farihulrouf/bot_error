<template>
  <div class="relative w-[300px] h-[300px] overflow-hidden py-2 p-2 rounded-xl">
    <div class="absolute z-20 left-0" v-if="!device.qr">
      <!-- Meneruskan properti bgColor ke DropDown -->
      <DropDown
        :menuItems="menuItems"
        :profile="profile"
        @menu-click="handleMenuClick"
        :bgColor="bgColor"
      />
    </div>
    <div
      class="group relative cursor-pointer overflow-hidden bg-white px-6 pt-2 pb-4 shadow-xl ring-1 ring-gray-900/5 transition-all duration-300 hover:-translate-y-1 hover:shadow-2xl sm:mx-auto sm:max-w-sm sm:rounded-lg sm:px-10"
      @mouseenter="isHovering = true"
      @mouseleave="isHovering = false"
      :class="{ 'bg-blue-100': isHovering }"
    >
      <div class="relative z-10 mx-auto max-w-md">
        <div
          class="space-y-6 pt-2 text-base leading-7 text-gray-600 transition-all duration-300"
        >
          <template v-if="!device.qr">
            <p v-if="device.id">{{ device.id }}</p>
            <p v-if="device.number">Phone: {{ device.number.split(':')[0] }}</p>
            <p v-if="device.name">Name: {{ device.name }}</p>
            <p>Status: ready</p>
            <p>Process: getMessages</p>
          </template>
          <div class="flex justify-center w-full h-[240px]" v-if="device.qr">
            <QRCodeVue3
              :width="240"
              :height="300"
              :value="device.qr"
              :qrOptions="{ errorCorrectionLevel: 'H' }"
              :dotsOptions="{ type: 'dots', color: '#4F46E5', gradient: { type: 'linear', rotation: 0, colorStops: [{ offset: 0, color: '#4F46E5' }, { offset: 1, color: '#34495E' }] } }"
              :imageOptions="{ hideBackgroundDots: true, imageSize: 0.4, margin: 10 }"
              :cornersSquareOptions="{ type: 'dot', color: '#34495E' }"
              :cornersDotOptions="{ type: undefined, color: '#41B883' }"
              :backgroundOptions="{ color: '#FFFFFF' }"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import QRCodeVue3 from 'qrcode-vue3';
import DropDown from './DropDown.vue';
import api from '../api/api.js';
import { mdiDotsVertical } from '@mdi/js';

export default {
  name: 'QrcodeCard',
  components: {
    QRCodeVue3,
    DropDown,
  },
  props: {
    device: {
      type: Object,
      required: true,
    },
    bgColor: {
      type: String,
      default: 'bg-indigo-600 rounded-full w-6 h-6 flex justify-center absolute left-4 text-white',
      // Warna latar belakang default (jika tidak disediakan oleh parent)
    },
  },
  data() {
    return {
      isHovering: false,
      menuItems: [
        { name: 'Settings', action: 'settings' },
        { name: 'Messages', action: 'Messages' },
        { name: 'Logout', action: 'logout' },
      ],
      profile: {
        name: '', // Ganti dengan data profil yang sesuai
        // Tambahkan properti lain sesuai kebutuhan, seperti email, avatar, dll.
        comp: mdiDotsVertical, // Contoh nilai untuk properti comp
      },
    };
  },
  methods: {
    handleMenuClick(action) {
      if (action === 'settings') {
        this.navigateToSettings();
      } else if (action === 'logout') {
        this.logout();
      }
    },
    navigateToSettings() {
      // Logika navigasi ke pengaturan di sini
      console.log('Navigating to settings');
    },
    async logout() {
      try {
        console.log('Nomor perangkat sebelum logout:', this.device.number);
        const response = await api.delete('/system/logout/' + this.device.number);

        console.log('Respon logout:', response);

        if (response.data.status === 'success') {
          console.log('Logout berhasil');
        } else {
          console.log('Logout gagal: ', response.data.message);
        }
      } catch (error) {
        console.error('Error saat logout:', error);
      }
    },
    getMessage() {
      // Implementasikan logika untuk mendapatkan pesan status berdasarkan status perangkat
      return 'Aktif';
    },
  },
};
</script>

<style scoped>
/* CSS menggunakan Tailwind CSS atau CSS kustom */
</style>
