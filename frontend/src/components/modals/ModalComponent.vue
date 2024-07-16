<template>
  <div class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50">
    <div class="bg-white rounded-lg shadow-lg w-lg p-2 px-6">
      <div class="mb-2">
        <p>{{ id }} : {{ phone }}</p>
      </div>
      <div class="mb-2">
        <label class="block mb-2 font-bold" for="to">To</label>
        <input type="text" id="to" v-model="to" class="w-full p-2 border border-gray-300 rounded" />
      </div>
      <div class="mb-2">
        <label class="block mb-2 font-bold" for="type">Type</label>
        <select id="type" v-model="type" class="w-full p-2 border border-gray-300 rounded">
          <option value="text">Text</option>
          <option value="media">Media</option>
          <option value="audio">Audio</option>
          <option value="doc">Doc</option>
        </select>
      </div>
      <div class="mb-2">
        <label class="block mb-2 font-bold" for="message">Message</label>
        <textarea id="message" v-model="text" class="w-full p-2 border border-gray-300 rounded" rows="4"></textarea>
      </div>
      <div class="mb-2">
        <label class="block mb-2 font-bold" for="caption">Caption</label>
        <input type="text" id="caption" v-model="caption" class="w-full p-2 border border-gray-300 rounded" />
      </div>
      <div class="mb-2 flex items-center">
        <div class="w-1/2">
          <label class="block mb-2 font-bold" for="file">Upload Image</label>
          <input type="file" id="file" @change="handleFileUpload" class="w-full p-2 border border-gray-300 rounded" />
          <p v-if="selectedFile" class="mt-2 text-sm text-gray-500">Selected File: {{ selectedFile.name }}</p>
        </div>
        <div class="w-1/2 pl-2">
          <label class="block mb-2 font-bold" for="url">File URL</label>
          <input type="text" id="url" v-model="fileUrl" class="w-full p-2 border border-gray-300 rounded" />
        </div>
      </div>
      <div class="flex justify-end items-center">
        <button class="bg-whatsapp_teal text-white px-4 py-2 rounded flex items-center" @click="send">
          <span>Send</span>
          <svg-icon type="mdi" :path="sendPath" class="ml-2"></svg-icon>
        </button>
        <button class="ml-2 bg-gray-500 text-white px-4 py-2 rounded flex items-center" @click="$emit('close')">
          <span>Close</span>
          <svg-icon type="mdi" :path="closePath" class="ml-2"></svg-icon>
        </button>
      </div>
    </div>
  </div>
</template>

<script>
import api from '@/api/api.js'; // Sesuaikan path sesuai struktur aplikasi Anda
import { showNotification } from '@/utils/notification.js'; // Import fungsi showNotification
import SvgIcon from '@jamescoyle/vue-icon';
import { mdiSendCheck, mdiClose } from '@mdi/js';

export default {
  name: "ModalComponent",
  components: {
    SvgIcon
  },
  props: {
    id: {
      type: String,
      required: true,
    },
    phone: {
      type: String,
      required: true,
    },
    name: {
      type: String,
      required: true,
    },
  },
  data() {
    return {
      to: "",
      type: "text",
      text: "",
      caption: "",
      selectedFile: null,
      fileUrl: "",
      sendPath: mdiSendCheck,
      closePath: mdiClose
    };
  },
  methods: {
    send() {
      // Validasi bidang yang diperlukan
      if (!this.to || !this.type || !this.text || !this.phone) {
        alert("Missing required fields: 'to', 'type', 'text', or 'from'");
        return;
      }

      const token = localStorage.getItem('token');
      const payload = {
        to: this.to,
        type: this.type.toLowerCase(),
        text: this.text,
        caption: this.caption,
        url: this.fileUrl,
        from: this.phone
      };

      // Mengirim permintaan POST ke API
      api.post('/api/messages', payload, {
        headers: {
          'Authorization': `Bearer ${token}`
        }
      })
      .then(response => {
        // Handle respon dari server
        console.log('Response:', response);
        // Menampilkan notifikasi sukses
        showNotification('Success', 'Message sent successfully', 'success');
        // Bersihkan form setelah berhasil dikirim
        this.clearForm();
        // Emit event untuk menutup modal
        this.$emit('close');
      })
      .catch(error => {
        // Tangani error jika permintaan gagal
        console.error('Error:', error);
        // Anda dapat menambahkan logika penanganan error sesuai kebutuhan
      });
    },
    handleFileUpload(event) {
      const file = event.target.files[0];
      if (file) {
        this.selectedFile = file;

        // Simulasi pengunggahan file ke server
        setTimeout(() => {
          this.fileUrl = file.name; // Simulasi URL file yang diunggah
        }, 1000);
      }
    },
    clearForm() {
      this.to = "";
      this.type = "text";
      this.text = "";
      this.caption = "";
      this.selectedFile = null;
      this.fileUrl = "";
    }
  }
};
</script>

<style scoped>
/* Optional: Tambahkan styling khusus di sini jika diperlukan */
</style>
