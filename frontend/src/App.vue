<template>
  <div id="app">
    <Navbar />
    <div class="container mx-auto my-4">
    </div>
    <div class="container mx-auto my-4 w-full flex flex-wrap gap-4">
      <!-- Show LoadingSpinner component if deviceData is empty -->
      <LoadingSpin class="absolute" v-if="deviceData.length === 0" />
      <!-- Show QrcodeCard components once deviceData is fetched -->
      <QrcodeCard v-else v-for="device in deviceData" :key="device.id || device.qr" :device="device" />
    </div>
  </div>
</template>

<script>
import Navbar from './components/Navbar.vue';
import QrcodeCard from './components/QrcodeCard.vue';
import LoadingSpin from './components/LoadingSpin.vue'; // Updated import for renamed component
import './assets/css/index.css'; // Import Tailwind CSS here
import axios from 'axios';

export default {
  name: 'App',
  components: {
    Navbar,
    QrcodeCard,
    LoadingSpin // Updated registration for renamed component
  },
  data() {
    return {
      deviceData: []  // State to store fetched data from API
    };
  },
  created() {
    this.fetchDeviceData();
  },
  methods: {
    async fetchDeviceData() {
      try {
        const response = await axios.get('http://localhost:8080/api/system/devices');
        this.deviceData = response.data;
        console.log('Fetched Device Data:', this.deviceData);  // Log data to console
      } catch (error) {
        console.error('Error fetching device data:', error);
      }
    }
  }
};
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  text-align: center;
  color: #2c3e50;
}
</style>
