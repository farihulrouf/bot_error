<template>
  <div class="relative flex justify-between items-center">
    <!-- Bagian untuk menampilkan label dan input field Webhook -->
    <div class="w-1/2 pr-4">
      <h3 class="font-bold">Webhook</h3>
      <input type="text" v-model="webhook" class="border-gray-300 border rounded px-3 py-2 w-full mt-2">
      <div class="flex justify-end mt-2">
        <button @click="saveWebhook" :disabled="isLoading" class="bg-green-500 hover:bg-green-700 text-white font-bold py-2 px-4 rounded mr-2 focus:outline-none focus:shadow-outline">
          Save
        </button>
        <button @click="closeWebhook" class="bg-gray-500 hover:bg-gray-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline">
          Close
        </button>
      </div>
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

    <!-- LoadingSpin ditampilkan di tengah saat isLoading true -->
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
      webhook: '', // Variabel untuk menyimpan URL webhook
      displayedToken: '', // Variabel untuk menampilkan 15 karakter terakhir token
      isLoading: false // Variabel untuk menunjukkan status loading
    };
  },
  created() {
    // Ambil token dari localStorage saat komponen dibuat
    this.token = localStorage.getItem('token');
    // Potong token menjadi 15 karakter terakhir
    this.displayedToken = this.token ? this.token.slice(-15) : '';
    this.fetchUserData();
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
    async fetchUserData() {
      try {
        const token = localStorage.getItem('token');
        const config = {
          headers: { Authorization: `Bearer ${token}` }
        };
        const response = await api.get('/api/user/detail', config);
        
        // Assign nilai URL webhook dari response data ke variabel webhook
        this.webhook = response.url; // Pastikan 'url' sesuai dengan field yang mengandung URL

        console.log('Fetched User Data:', response);
      } catch (error) {
        console.error('Error fetching user data:', error);
      }
    },
    async generateNewToken() {
      try {
        this.isLoading = true; // Set loading menjadi true saat memulai proses
        const response = await api.post('/token', {});
        const newToken = response.token;

        // Perbarui token di data komponen
        this.token = newToken;
        this.displayedToken = newToken.slice(-15);

        // Menunda perubahan status isLoading menjadi false
        setTimeout(() => {
          this.isLoading = false; // Set loading menjadi false setelah beberapa detik
        }, 1000); // Delay 2 detik (2000 milidetik)
      } catch (error) {
        console.error('Error generating new token:', error);
        alert('Failed to generate new token');
        this.isLoading = false; // Set loading menjadi false jika terjadi error
      }
    },
    async saveWebhook() {
      try {
        this.isLoading = true; // Set loading menjadi true saat memulai proses
        const token = localStorage.getItem('token');
        const config = {
          headers: { Authorization: `Bearer ${token}` }
        };
        const data = { url: this.webhook }; // Siapkan data untuk dikirim ke API
        const response = await api.put('webhook/update', data, config); // Menggunakan method PUT untuk update

        console.log('Webhook updated successfully:', response);
        alert('Webhook updated successfully'); // Tampilkan alert atau feedback sukses

        // Menunda perubahan status isLoading menjadi false
        setTimeout(() => {
          this.isLoading = false; // Set loading menjadi false setelah beberapa detik
        }, 1000); // Delay 2 detik (2000 milidetik)
      } catch (error) {
        console.error('Error updating webhook:', error);
        alert('Failed to update webhook'); // Tampilkan alert atau feedback error
        this.isLoading = false; // Set loading menjadi false jika terjadi error
      }
    },
    closeWebhook() {
      // Tambahkan logika jika diperlukan saat tombol Close ditekan
      console.log('Close button clicked');
    }
  }
};
</script>

<style scoped>
/* Tambahkan styling khusus jika diperlukan */
</style>
