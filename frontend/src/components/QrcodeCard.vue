<template>
  <div class="relative w-[350px] h-[450px] overflow-hidden bg-gray-50 py-4">
    <div
      class="group relative cursor-pointer overflow-hidden bg-white px-6 pt-10 pb-8 shadow-xl ring-1 ring-gray-900/5 transition-all duration-300 hover:-translate-y-1 hover:shadow-2xl sm:mx-auto sm:max-w-sm sm:rounded-lg sm:px-10"
      @mouseenter="isHovering = true"
      @mouseleave="isHovering = false"
      :class="{ 'bg-blue-100': isHovering }"
    >
      <div class="relative z-10 mx-auto max-w-md">
        <div
          class="space-y-6 pt-5 text-base leading-7 text-gray-600 transition-all duration-300"
        >
          <template v-if="!device.qr">
            <p v-if="device.id">ID: {{ device.id }}</p>
            <p v-if="device.number">Phone: {{ device.number }}</p>
            <p v-if="device.name">Name: {{ device.name }}</p>
            <p>Status: getMessage</p>
          </template>
          <div v-if="device.qr">
            <img :src="generateQRCodeURL(device.qr)" alt="QR Code">
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'QrcodeCard',
  props: {
    device: {
      type: Object,
      required: true
    }
  },
  data() {
    return {
      isHovering: false
    };
  },
  methods: {
    generateQRCodeURL(qrCodeData) {
      // Example: Replace with actual method to generate QR code URL or image
      return `https://api.qrserver.com/v1/create-qr-code/?data=${encodeURIComponent(qrCodeData)}`;
    }
  }
};
</script>

<style scoped>
/* CSS using Tailwind CSS or custom CSS */
</style>
