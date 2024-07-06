<template>
  <div class="relative w-[350px] h-[370px] overflow-hidden py-4 p-4 rounded-xl">
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
          <div v-if="device.qr">
            <qrcode-vue
              :value="device.qr"
              :size="240"
              level="H"
              render-as="canvas"
            ></qrcode-vue>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import QrcodeVue from "qrcode.vue";
import DropDown from "./DropDown.vue";
import api from "../api/api.js";
import { mdiDotsVertical } from "@mdi/js";
export default {
  name: "QrcodeCard",
  components: {
    QrcodeVue,
    DropDown,
  },
  props: {
    device: {
      type: Object,
      required: true,
    },
    bgColor: {
      type: String,
      default: "bg-indigo-600 rounded-full w-6 h-6 flex justify-center absolute left-4 text-white", // Default background color (if not provided by parent)
    },
  },
  data() {
    return {
      isHovering: false,
      menuItems: [
        { name: "Settings", action: "settings" },
        { name: "Messages", action: "Messages" },
        { name: "Logout", action: "logout" },
      ],
      profile: {
        name: "Profile", // Ganti dengan data profil yang sesuai
        // Tambahkan properti lain sesuai kebutuhan, seperti email, avatar, dll.
        comp: mdiDotsVertical, // Contoh nilai untuk properti comp
      },
    };
  },
  methods: {
    handleMenuClick(action) {
      if (action === "settings") {
        this.navigateToSettings();
      } else if (action === "logout") {
        this.logout();
      }
    },
    navigateToSettings() {
      // Logika navigasi ke settings di sini
      console.log("Navigating to settings");
    },
    async logout() {
      try {
        console.log("Device number before logout:", this.device.number);
        const response = api.delete(
          "/system/logout/" + this.device.number
        );

        console.log("Logout response:", response);

        if (response.data.status === "success") {
          console.log("Logout berhasil");
        } else {
          console.log("Logout gagal: ", response.data.message);
        }
      } catch (error) {
        console.error("Error saat logout:", error);
      }
    },
    getMessage() {
      // Implement logic to get status message based on device status
      return "Active";
    },
  },
};
</script>

<style scoped>
/* CSS menggunakan Tailwind CSS atau CSS kustom */
</style>
