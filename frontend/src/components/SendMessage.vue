<template>
    <div class="max-w-md mx-auto p-6 border border-gray-300 rounded-lg bg-white shadow-md">
      <div class="mb-4">
        <label for="to" class="block text-sm font-medium text-gray-700">To:</label>
        <input type="text" id="to" v-model="to" class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:border-indigo-500 focus:ring-indigo-500" />
      </div>
      <div class="mb-4">
        <label for="type" class="block text-sm font-medium text-gray-700">Type:</label>
        <select id="type" v-model="type" class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:border-indigo-500 focus:ring-indigo-500">
          <option value="text">Text</option>
          <option value="video">Video</option>
          <option value="document">Document</option>
          <option value="image">Image</option>
        </select>
      </div>
      <div class="mb-4">
        <label for="content" class="block text-sm font-medium text-gray-700">Content Message:</label>
        <textarea id="content" v-model="content" class="mt-1 block w-full border-gray-300 rounded-md shadow-sm focus:border-indigo-500 focus:ring-indigo-500"></textarea>
      </div>
      <div class="flex justify-end space-x-2">
        <button @click="close" class="py-2 px-4 bg-gray-500 text-white rounded-md shadow-sm hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-500">Close</button>
        <button @click="send" class="py-2 px-4 bg-blue-500 text-white rounded-md shadow-sm hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500">Send</button>
      </div>
    </div>
  </template>
  
  <script>
  export default {
    data() {
      return {
        to: '',
        type: 'text',
        content: ''
      };
    },
    methods: {
      close() {
        // Implementasi logika tombol close
        this.to = '';
        this.type = 'text';
        this.content = '';
        console.log('Form closed');
      },
      send() {
        // Implementasi logika tombol send
        const message = {
          to: this.to,
          type: this.type,
          content: this.content
        };
        console.log('Message sent:', message);
        // Tambahkan logika untuk mengirim pesan, misalnya panggilan API
      }
    }
  };
  </script>
  
  <style scoped>
  /* Tidak ada style khusus, menggunakan Tailwind CSS */
  </style>
  