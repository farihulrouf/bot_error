<template>
  <div class="container mx-auto my-4">
    <Navbar />
    <div class="w-full mx-auto mt-4 mb-4">
      <TableHeader 
        :account="account"
        :devicesCount="devicesCount"
        :expiredAt="expiredAt"
        :balance="balance"
        :setting="setting"
        :apiVersion="apiVersion"
        @toggle-token-webhook="toggleTokenWebhook"
        @show-update-profile="showUpdateProfile"
      />
    </div>
    <div v-if="showTokenWebhook" class="max-w-screen-md mx-auto mt-4 mb-4">
      <GenerateTokenVue />
    </div>
    <div v-else-if="showProfile" class="max-w-screen-md mx-auto mt-4 mb-4">
      <UpdateProfile />
    </div>
    <div v-else class="p-4 max-w-screen-md mt-4 mb-4 mx-auto flex justify-center flex-wrap gap-4 relative">
      <LoadingSpin v-if="isLoading" />
      <QrcodeCard v-else v-for="device in deviceData" :key="device.id || device.qr" :device="device" />
    </div>
  </div>
</template>

<script>
import jwtDecode from 'jwt-decode';  // Import jwt-decode library
import Navbar from '@/components/Navbar.vue';
import QrcodeCard from '@/components/QrcodeCard.vue';
import LoadingSpin from '@/components/LoadingSpin.vue';
import TableHeader from '@/components/TableHeader.vue';
import GenerateTokenVue from '../components/GenerateToken.vue';
import UpdateProfile from '../components/UpdateProfile.vue'; // Import UpdateProfile component
import api from "../api/api.js"; // Import Api.js untuk melakukan request HTTP

export default {
  name: 'DashboardBot',
  components: {
    Navbar,
    QrcodeCard,
    LoadingSpin,
    TableHeader,
    GenerateTokenVue,
    UpdateProfile
  },
  data() {
    return {
      deviceData: [],
      isLoading: true,
      account: 'farihul', // Updated to be set dynamically
      devicesCount: 1,
      expiredAt: '2030-12-12',
      balance: '30000',
      setting: 'token, webhook',
      apiVersion: 'v1.14.3',
      showTokenWebhook: false, // State to control visibility of GenerateTokenVue
      showProfile: false // State to control visibility 
    };
  },
  created() {
    this.decodeToken();
    this.fetchDeviceData();
  },
  methods: {
    async decodeToken() {
      try {
        const token = localStorage.getItem('token');
        if (token) {
          const decoded = jwtDecode(token);
          this.account = decoded.username; // Set account from token's username
        }
      } catch (error) {
        console.error('Error decoding token:', error);
      }
    },
    async fetchDeviceData() {
      try {
        const token = localStorage.getItem('token');
        const config = {
          headers: { Authorization: `Bearer ${token}` }
        };
        const response = await api.get('/system/devices', config);
        this.deviceData = response;
        this.isLoading = false;
        console.log('Fetched Device Data:', this.deviceData);
      } catch (error) {
        console.error('Error fetching device data:', error);
        this.isLoading = false;
      }
    },
    toggleTokenWebhook() {
      this.showTokenWebhook = !this.showTokenWebhook;
      this.showProfile = false; // Hide profile if token webhook is shown
    },
    showUpdateProfile() {
      this.showProfile = true;
      this.showTokenWebhook = false; // Hide token webhook if profile is shown
    }
  }
};
</script>

<style>
/* Tambahkan styling khusus jika diperlukan */
</style>
