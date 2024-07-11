<template>
  <div class="p-4 bg-[#dcf8c6] flex justify-center relative">
    <div class="absolute z-20 right-8" v-if="!device.qr">
      <!-- Meneruskan properti bgColor ke DropDown -->
      <DropDown
        :menuItems="menuItems"
        :profile="profile"
        @menu-click="handleMenuClick"
        :bgColor="bgColor"
      />
    </div>
    <div
      class="block w-[300px] rounded-lg bg-warning text-black shadow-secondary-1 p-4"
>
      <div class="w-full flex-grow space-y-8 text-[#075e54]" v-if="!device.qr">
        <p v-if="device.id">{{ device.id }}</p>
        <p v-if="device.number">Phone: {{ device.number.split(":")[0] }}</p>
        <p v-if="device.name">Name: {{ device.name }}</p>
        <p>Status: ready</p>
        <p>Process: getMessages</p>
      </div>
      <div v-if="device.qr">
      <QRCodeVue3
        :width="500"
        :height="500"
        :value="device.qr"
        :qrOptions="{ errorCorrectionLevel: 'H' }"
        :dotsOptions="{ type: 'dots', color: '#34B7F1', gradient: { type: 'linear', rotation: 0, colorStops: [{ offset: 0, color: '#4F46E5' }, { offset: 1, color: '#075e54' }] } }"
        :imageOptions="{ hideBackgroundDots: true, imageSize: 0.4, margin: 10 }"
        :cornersSquareOptions="{ type: 'dot', color: '#25D366' }"
        :cornersDotOptions="{ type: undefined, color: '#41B883' }"
        :backgroundOptions="{ color: '#dcf8c6' }"
      />
    </div>
    </div>
    <Modal v-if="showModal" @close="showModal = false">
      <template #header>
        <h3 class="text-lg font-bold">Messages</h3>
      </template>
      <template #body>
        <p>This is the messages modal content.</p>
      </template>
      <template #footer>
        <button class="bg-red-500 text-white px-4 py-2 rounded" @click="showModal = false">Close</button>
      </template>
    </Modal>
  </div>
</template>
<script>
import QRCodeVue3 from "qrcode-vue3";
import DropDown from "./DropDown.vue";
import api from "../api/api.js";
import { mdiDotsVertical } from "@mdi/js";
import Modal from "./ModalMessages.vue";

export default {
  name: "QrcodeCard",
  components: {
    QRCodeVue3,
    DropDown,
    Modal,
  },
  props: {
    device: {
      type: Object,
      required: true,
    },
    bgColor: {
      type: String,
      default:
        "bg-[#25d366] rounded-full w-6 h-6 flex justify-center absolute left-4 text-white",
      // Warna latar belakang default (jika tidak disediakan oleh parent)
    },
  },
  data() {
    return {
      isLoading: false, 
      isHovering: false,
      showModal: false,
      menuItems: [
        { name: "Settings", action: "settings" },
        { name: "Messages", action: "Messages" },
        { name: "Logout", action: "logout" },
      ],
      profile: {
        name: "", // Ganti dengan data profil yang sesuai
        // Tambahkan properti lain sesuai kebutuhan, seperti email, avatar, dll.
        comp: mdiDotsVertical, // Contoh nilai untuk properti comp
      },
    };
  },
  methods: {
    handleMenuClick(action) {
      if (action === "settings") {
        this.navigateToSettings();
      }
      else if (action === "Messages") {
        this.showModal = true;
      } 
      else if (action === "logout") {
        this.logout();
      }
    },
    navigateToSettings() {
      // Logika navigasi ke pengaturan di sini
      console.log("Navigating to settings");
    },
    async logout() {
      try {
        console.log("Nomor perangkat sebelum logout:", this.device.id);
        const response = await api.delete("/system/logout/" + this.device.id);

        console.log("Respon logout:", response);

        if (response.status === "success") {
          console.log("Logout berhasil");
          this.$emit('device-logged-out'); // Emit event to parent component
        } else {
          console.log("Logout gagal: ", response.data.message);
        }
      } catch (error) {
        console.error("Error saat logout:", error);
      }
    },
    getMessage() {
      // Implementasikan logika untuk mendapatkan pesan status berdasarkan status perangkat
      return "Aktif";
    },
  },
};
</script>

<style scoped>
/* CSS menggunakan Tailwind CSS atau CSS kustom */
</style>
