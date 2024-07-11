<template>
  <div class="flex justify-between items-start mt-12 relative">
    <!-- Bagian untuk memperbarui profil -->
    <div class="w-1/2 pr-4">
      <h2 class="text-2xl font-bold mb-8">Update Profile</h2>
      <form @submit.prevent="updateProfile">
        <div class="mb-4">
          <div class="relative">
            <input autocomplete="off" id="username" v-model="username" type="text" class="peer placeholder-transparent h-10 w-full border-b-2 border-gray-300 text-gray-900 focus:outline-none focus:border-rose-600" required>
            <label for="username" class="absolute left-0 -top-3.5 text-gray-600 text-sm peer-placeholder-shown:text-base peer-placeholder-shown:text-gray-440 peer-placeholder-shown:top-2 transition-all peer-focus:-top-3.5 peer-focus:text-gray-600 peer-focus:text-sm">Username</label>
          </div>
        </div>
        <div class="mb-4">
          <div class="relative">
            <input autocomplete="off" id="firstName" v-model="firstName" type="text" class="peer placeholder-transparent h-10 w-full border-b-2 border-gray-300 text-gray-900 focus:outline-none focus:border-rose-600" required>
            <label for="firstName" class="absolute left-0 -top-3.5 text-gray-600 text-sm peer-placeholder-shown:text-base peer-placeholder-shown:text-gray-440 peer-placeholder-shown:top-2 transition-all peer-focus:-top-3.5 peer-focus:text-gray-600 peer-focus:text-sm">First Name</label>
          </div>
        </div>
        <div class="mb-4">
          <div class="relative">
            <input autocomplete="off" id="lastName" v-model="lastName" type="text" class="peer placeholder-transparent h-10 w-full border-b-2 border-gray-300 text-gray-900 focus:outline-none focus:border-rose-600" required>
            <label for="lastName" class="absolute left-0 -top-3.5 text-gray-600 text-sm peer-placeholder-shown:text-base peer-placeholder-shown:text-gray-440 peer-placeholder-shown:top-2 transition-all peer-focus:-top-3.5 peer-focus:text-gray-600 peer-focus:text-sm">Last Name</label>
          </div>
        </div>
        <div class="mb-4">
          <div class="relative">
            <input autocomplete="off" id="email" v-model="email" type="email" class="peer placeholder-transparent h-10 w-full border-b-2 border-gray-300 text-gray-900 focus:outline-none focus:border-rose-600" required>
            <label for="email" class="absolute left-0 -top-3.5 text-gray-600 text-sm peer-placeholder-shown:text-base peer-placeholder-shown:text-gray-440 peer-placeholder-shown:top-2 transition-all peer-focus:-top-3.5 peer-focus:text-gray-600 peer-focus:text-sm">Email</label>
          </div>
        </div>
        <div class="flex items-center justify-between">
          <button type="submit" :disabled="isLoadingProfile" class="bg-whatsapp_teal hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline">
            <span v-if="!isLoadingProfile">Save Profile</span>
            <loading-spin v-if="isLoadingProfile" class="ml-2"></loading-spin>
          </button>
        </div>
      </form>
    </div>

    <!-- Bagian untuk mengubah kata sandi -->
    <div class="w-1/2 pl-4">
      <h2 class="text-2xl font-bold mb-8">Change Password</h2>
      <form @submit.prevent="changePassword">
        <div class="mb-4">
          <div class="relative">
            <input autocomplete="off" id="currentPassword" v-model="currentPassword" type="password" class="peer placeholder-transparent h-10 w-full border-b-2 border-gray-300 text-gray-900 focus:outline-none focus:border-rose-600" required>
            <label for="currentPassword" class="absolute left-0 -top-3.5 text-gray-600 text-sm peer-placeholder-shown:text-base peer-placeholder-shown:text-gray-440 peer-placeholder-shown:top-2 transition-all peer-focus:-top-3.5 peer-focus:text-gray-600 peer-focus:text-sm">Current Password</label>
          </div>
        </div>
        <div class="mb-4">
          <div class="relative">
            <input autocomplete="off" id="newPassword" v-model="newPassword" type="password" class="peer placeholder-transparent h-10 w-full border-b-2 border-gray-300 text-gray-900 focus:outline-none focus:border-rose-600" required>
            <label for="newPassword" class="absolute left-0 -top-3.5 text-gray-600 text-sm peer-placeholder-shown:text-base peer-placeholder-shown:text-gray-440 peer-placeholder-shown:top-2 transition-all peer-focus:-top-3.5 peer-focus:text-gray-600 peer-focus:text-sm">New Password</label>
          </div>
        </div>
        <div class="mb-4">
          <div class="relative">
            <input autocomplete="off" id="retypePassword" v-model="retypePassword" type="password" class="peer placeholder-transparent h-10 w-full border-b-2 border-gray-300 text-gray-900 focus:outline-none focus:border-rose-600" required>
            <label for="retypePassword" class="absolute left-0 -top-3.5 text-gray-600 text-sm peer-placeholder-shown:text-base peer-placeholder-shown:text-gray-440 peer-placeholder-shown:top-2 transition-all peer-focus:-top-3.5 peer-focus:text-gray-600 peer-focus:text-sm">Retype New Password</label>
          </div>
        </div>
        <div class="flex items-center justify-between">
          <button type="submit" :disabled="isLoadingPassword" class="bg-whatsapp_teal hover:bg-green-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline">
            <span v-if="!isLoadingPassword">Save Password</span>
            <loading-spin v-if="isLoadingPassword" class="ml-2"></loading-spin>
          </button>
        </div>
      </form>
    </div>

    <!-- Loading Spinners -->
    <div class="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2" v-if="isLoadingProfile || isLoadingPassword">
      <loading-spin></loading-spin>
    </div>
  </div>
</template>

<script>
import api from '../api/api.js';
import LoadingSpin from './LoadingSpin.vue';
import { showNotification } from '../utils/notification'; // Import utilitas notifikasi

export default {
  name: 'UpdateProfile',
  components: {
    LoadingSpin,
  },
  data() {
    return {
      username: '',
      firstName: '',
      lastName: '',
      email: '',
      currentPassword: '',
      newPassword: '',
      retypePassword: '',
      isLoadingProfile: false, // Status loading untuk update profil
      isLoadingPassword: false, // Status loading untuk ubah kata sandi
    };
  },
  created() {
    this.fetchUserData();
  },
  methods: {
    async fetchUserData() {
      try {
        // Mengambil data pengguna dari API
        const token = localStorage.getItem('token');
        const config = {
          headers: { Authorization: `Bearer ${token}` }
        };
        const response = await api.get('/user/detail', config);

        // Mengisi data dari respons ke properti data
        this.username = response.username;
        this.firstName = response.first_name;
        this.lastName = response.last_name;
        this.email = response.email;
        console.log('Fetched User Data:', response);
      } catch (error) {
        console.error('Error fetching user data:', error);
      }
    },
    async updateProfile() {
      this.isLoadingProfile = true;
      try {
        // Panggil API untuk update profil
        const payload = {
          first_name: this.firstName,
          last_name: this.lastName,
          email: this.email,
        };
        const response = await api.put('/api/user/update', payload);
        console.log('Profile updated:', response);
      } catch (error) {
        console.error('Failed to update profile:', error);
      } finally {
        this.isLoadingProfile = false;
        showNotification('Success', 'Profile updated successfully!', 'success');
      }
    },
    async changePassword() {
      this.isLoadingPassword = true;
      try {
        // Panggil API untuk mengubah password
        const payload = {
          first_name: this.firstName,
          last_name: this.lastName,
          email: this.email,
          current_password: this.currentPassword,
          new_password: this.newPassword
        };
        const response = await api.put('/api/user/update', payload);
        console.log('Password changed:', response);
      } catch (error) {
        console.error('Failed to change password:', error);
      } finally {
        this.isLoadingPassword = false;
        showNotification('Success', 'Password changed successfully!', 'success');
      }
    }
  }
};
</script>

<style scoped>
/* Tambahkan styling khusus jika diperlukan */
</style>
