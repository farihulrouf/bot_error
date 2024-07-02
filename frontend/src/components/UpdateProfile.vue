<template>
    <div class="flex justify-between items-start">
      <!-- Bagian untuk memperbarui profil -->
      <div class="w-1/2 pr-4">
        <h2 class="text-2xl font-bold mb-4">Update Profile</h2>
        <form @submit.prevent="updateProfile">
          <div class="mb-4">
            <label for="firstName" class="block text-gray-700 text-sm font-bold mb-2">First Name</label>
            <input type="text" id="firstName" v-model="firstName" class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline">
          </div>
          <div class="mb-4">
            <label for="lastName" class="block text-gray-700 text-sm font-bold mb-2">Last Name</label>
            <input type="text" id="lastName" v-model="lastName" class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline">
          </div>
          <div class="mb-4">
            <label for="email" class="block text-gray-700 text-sm font-bold mb-2">Email</label>
            <input type="email" id="email" v-model="email" class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline">
          </div>
          <div class="flex items-center justify-between">
            <button type="submit" class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline">
              Save
            </button>
          </div>
        </form>
      </div>
  
      <!-- Bagian untuk mengubah kata sandi -->
      <div class="w-1/2 pl-4">
        <h2 class="text-2xl font-bold mb-4">Change Password</h2>
        <form @submit.prevent="changePassword">
          <div class="mb-4">
            <label for="currentPassword" class="block text-gray-700 text-sm font-bold mb-2">Current Password</label>
            <input type="password" id="currentPassword" v-model="currentPassword" class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline">
          </div>
          <div class="mb-4">
            <label for="newPassword" class="block text-gray-700 text-sm font-bold mb-2">New Password</label>
            <input type="password" id="newPassword" v-model="newPassword" class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline">
          </div>
          <div class="mb-4">
            <label for="retypePassword" class="block text-gray-700 text-sm font-bold mb-2">Retype New Password</label>
            <input type="password" id="retypePassword" v-model="retypePassword" class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline">
          </div>
          <div class="flex items-center justify-between">
            <button type="submit" class="bg-green-500 hover:bg-green-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline">
              Save
            </button>
          </div>
        </form>
      </div>
    </div>
  </template>
  
  <script>
  export default {
    name: 'UpdateProfile',
    data() {
      return {
        firstName: '',
        lastName: '',
        email: '',
        currentPassword: '',
        newPassword: '',
        retypePassword: ''
      };
    },
    methods: {
      updateProfile() {
        // Tambahkan logika untuk memperbarui profil pengguna
        console.log('Profile updated:', {
          firstName: this.firstName,
          lastName: this.lastName,
          email: this.email
        });
        alert('Profile updated successfully!');
      },
      changePassword() {
        // Tambahkan logika untuk mengubah kata sandi
        if (this.newPassword !== this.retypePassword) {
          alert('New password and retype password do not match!');
          return;
        }
  
        console.log('Password changed:', {
          currentPassword: this.currentPassword,
          newPassword: this.newPassword
        });
        alert('Password changed successfully!');
      }
    }
  };
  </script>
  
  <style scoped>
  /* Tambahkan styling khusus jika diperlukan */
  </style>
  