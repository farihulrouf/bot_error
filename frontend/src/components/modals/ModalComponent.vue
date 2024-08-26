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
          <option value="image">Image</option>
          <option value="audio">Audio</option>
          <option value="video">Video</option>
          <option value="doc">Doc</option>
        </select>
      </div>
      <div class="mb-2">
        <label class="block mb-2 font-bold" for="message">Message</label>
        <textarea id="message" v-model="text" class="w-full p-2 border border-gray-300 rounded" rows="4"></textarea>
      </div>
      <div class="mb-2 flex">
        <div class="w-1/2 pr-2">
          <label class="block mb-2 font-bold" for="caption">Caption</label>
          <input type="text" id="caption" v-model="caption" class="w-full p-2 border border-gray-300 rounded" />
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
import api from '@/api/api.js'; // Adjust path as per your application's structure
import { showNotification } from '@/utils/notification.js'; // Import showNotification function
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
      fileUrl: "",
      sendPath: mdiSendCheck,
      closePath: mdiClose
    };
  },
  methods: {
    send() {
      // Validate required fields
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

      // Send POST request to API
      api.post('/api/messages', payload, {
        headers: {
          'Authorization': `Bearer ${token}`
        }
      })
      .then(response => {
        // Handle response from server
        console.log('Response:', response);
        // Display success notification
        showNotification('Success', 'Message sent successfully', 'success');
        // Clear form after successful submission
        this.clearForm();
        // Emit event to close modal
        this.$emit('close');
      })
      .catch(error => {
        // Handle error if request fails
        console.error('Error:', error);
        // Add error handling logic as needed
      });
    },
    clearForm() {
      this.to = "";
      this.type = "text";
      this.text = "";
      this.caption = "";
      this.fileUrl = "";
    }
  }
};
</script>

<style scoped>
/* Optional: Add specific styling here if needed */
</style>
