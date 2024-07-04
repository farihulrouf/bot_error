<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-100">
    <div class="bg-white p-8 rounded-lg shadow-md w-full max-w-md">
      <h2 class="text-2xl font-bold mb-6 text-center">Login</h2>
      <form @submit.prevent="handleLogin">
        <div class="mb-4">
          <label class="block text-gray-700 text-sm font-bold mb-2" for="username">Username</label>
          <input
            type="text"
            id="username"
            v-model="username"
            class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
            required
          />
        </div>
        <div class="mb-6">
          <label class="block text-gray-700 text-sm font-bold mb-2" for="password">Password</label>
          <input
            type="password"
            id="password"
            v-model="password"
            class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 mb-3 leading-tight focus:outline-none focus:shadow-outline"
            required
          />
        </div>
        <div class="flex items-center justify-between">
          <button
            type="submit"
            class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline"
          >
            Sign In
          </button>
        </div>
        <p v-if="errorMessage" class="text-red-500 text-xs italic mt-4">{{ errorMessage }}</p>
      </form>
    </div>
  </div>
</template>

<script>
import api from '../api/api.js';

export default {
  data() {
    return {
      username: '',
      password: '',
      errorMessage: ''
    };
  },
  methods: {
    handleLogin() {
      // Kirim permintaan login ke server
      api.post('/login', {
        username: this.username,
        password: this.password
      })
      .then(response => {
        // Tanggapan sukses dari server, simpan token ke localStorage
        const token = response.data.token;
        localStorage.setItem('token', token);
       
        this.$router.push('/dashboard');  // Redirect ke dashboard setelah login berhasil
      })
      .catch(error => {
        // Tangani error dari server
        if (error.response && error.response.status === 401) {
          this.errorMessage = 'Invalid username or password';
        } else {
          this.errorMessage = 'Failed to login';
        }
      });
    }
  }
};
</script>

<style scoped>
/* Tambahkan CSS khusus jika diperlukan */
</style>
