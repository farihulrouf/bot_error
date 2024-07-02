<template>
  <div class="relative flex justify-between items-center">
    <!-- Bagian untuk menampilkan label dan input field Webhook -->
    <div class="w-1/2 pr-4">
      <h3 class="font-bold">Webhook</h3>
      <input type="text" v-model="webhook" class="border-gray-300 border rounded px-3 py-2 w-full mt-2">
    </div>

    <!-- Bagian untuk menampilkan label dan nilai token -->
    <div class="w-1/2 pl-4">
      <div class="flex justify-between items-center">
        <h3 class="font-bold">API Token</h3>
        <p>Token: {{ displayedToken }}</p>
      </div>

      <!-- Tombol Copy dan New -->
      <div class="flex mt-2 items-center">
        <button @click="copyToken" class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded mr-2 focus:outline-none focus:shadow-outline">
          Copy
        </button>
        <button @click="generateNewToken" class="bg-green-500 hover:bg-green-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline">
          New
        </button>
      </div>
    </div>

    <!-- LoadingSpin ditampilkan di tengah -->
    <div v-if="isLoading" class="absolute inset-0 flex items-center justify-center bg-white bg-opacity-75">
      <LoadingSpin />
    </div>
  </div>
</template>

<script>
import api from '../api/api.js';
import LoadingSpin from './LoadingSpin.vue';

export default {
  name: 'WebhookAndToken',
  components: {
    LoadingSpin
  },
  data() {
    return {
      token: '', // Variabel untuk menyimpan token
      webhook: 'https://example.com/webhook', // Ganti dengan nilai webhook yang sesuai
      displayedToken: '', // Variabel untuk menampilkan token yang dipotong
      isLoading: false // Variabel untuk menunjukkan status loading
    };
  },
  created() {
    // Ambil token dari localStorage saat komponen dibuat
    this.token = localStorage.getItem('token');
    // Potong token menjadi 10 karakter
    this.displayedToken = this.token ? this.token.slice(0, 10) : '';
  },
  methods: {
    copyToken() {
      // Copy token lengkap ke clipboard
      const textarea = document.createElement('textarea');
      textarea.value = this.token;
      document.body.appendChild(textarea);
      textarea.select();
      document.execCommand('copy');
      document.body.removeChild(textarea);
      alert('Token copied to clipboard');
    },
    async generateNewToken() {
      try {
        this.isLoading = true; // Set loading menjadi true saat memulai proses
        const response = await api.post('/token', {});
        const newToken = response.token;

        // Perbarui token di data komponen
        this.token = newToken;
        this.displayedToken = newToken.slice(0, 10);

        // Menunda perubahan status isLoading menjadi false
        setTimeout(() => {
          this.isLoading = false; // Set loading menjadi false setelah beberapa detik
        }, 1000); // Delay 2 detik (2000 milidetik)
      } catch (error) {
        console.error('Error generating new token:', error);
        alert('Failed to generate new token');
        this.isLoading = false; // Set loading menjadi false jika terjadi error
      }
    }
  }
};
</script>

<style scoped>
/* Tambahkan styling khusus jika diperlukan */
</style>
