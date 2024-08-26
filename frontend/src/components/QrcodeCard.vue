<template>
  <div class="p-4 bg-[#dcf8c6] flex justify-center relative">
    <div class="absolute z-20 right-8" v-if="!device.qr">
      <DropDown
        :menuItems="menuItems"
        :profile="profile"
        @menu-click="handleMenuClick"
        :bgColor="bgColor"
      />
    </div>
    <div class="block w-[300px] rounded-lg bg-warning text-black shadow-secondary-1 p-4">
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
      <!-- Tampilkan FileInput.vue jika showFileInput bernilai true -->
      <ModalComponent
        v-if="showModal"
        @close="closeModalComponent"
        :id="device.id"
        :phone="device.number ? device.number.split(':')[0] : ''"
        :name="device.name"
      />
    </div>
  </div>
</template>

<script>
import QRCodeVue3 from "qrcode-vue3";
import DropDown from "./DropDown.vue";
import api from "../api/api.js";
import { mdiDotsVertical } from "@mdi/js";
import { showNotification } from '../utils/notification';
import ModalComponent from "./modals/ModalComponent.vue";

export default {
  name: "QrcodeCard",
  components: {
    QRCodeVue3,
    DropDown, 
    ModalComponent,
  },
  props: {
    device: {
      type: Object,
      required: true,
    },
    bgColor: {
      type: String,
      default:
        "bg-whatsapp_teal rounded-full w-6 h-6 flex justify-center absolute left-4 text-white",
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
        name: "",
        comp: mdiDotsVertical,
      },
    };
  },
  methods: {
    handleMenuClick(action) {
      if (action === "settings") {
        this.navigateToSettings();
      } else if (action === "Messages") {
        this.showModalComponent();
      } else if (action === "logout") {
        this.logout();
      }
    },
    showModalComponent() {
    this.showModal = true;
  },
  closeModalComponent() {
    this.showModal = false;
  },
    navigateToSettings() {
      console.log("Navigating to settings");
    },
    async logout() {
      try {
        console.log("Nomor perangkat sebelum logout:", this.device.id);
        const response = await api.delete("/system/logout/" + this.device.id);

        console.log("Respon logout:", response);

        if (response.status === "success") {
          console.log("Logout berhasil");
          showNotification('Success', 'Logout successfully!', 'success');
          this.$emit('device-logged-out');
        } else {
          console.log("Logout gagal: ", response.data.message);
        }
      } catch (error) {
        showNotification('Error', 'Failed to Logout', 'error');
        console.error("Error saat logout:", error);
      }
    },
    getMessage() {
      return "Aktif";
    },
    showMessageComponent() {
      this.showMessageBot = true;
      this.$emit('show-hello-world'); // Emit event untuk menunjukkan HelloWorld.vue di komponen induk
    },
  },
};
</script>

<style scoped>
/* Styling menggunakan Tailwind CSS atau CSS kustom */
</style>
