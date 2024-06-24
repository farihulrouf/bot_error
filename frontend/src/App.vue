<template>
  <div id="app">
    <Navbar />
    <div class="container mx-auto my-4">
      <ListAccount :headers="tableHeaders" :rows="tableData" />
    </div>
    <div class="container mx-auto my-4 w-full flex relative">
      <QrcodeCard />
      <QrcodeCard />
    </div>
  </div>
</template>

<script>
import Navbar from './components/Navbar.vue';
import ListAccount from './components/ListAccount.vue';
import './assets/css/index.css'; // Mengimpor file tailwind.css di sini
import QrcodeCard from './components/QrcodeCard.vue';

export default {
  name: 'App',
  components: {
    Navbar,
    ListAccount,  // Menambahkan ListTable ke dalam objek components
    QrcodeCard
  },
  data() {
    return {
      tableHeaders: ['Account', 'Devices', 'Expired', 'Balance', 'Settings', 'API'],
      tableData: [
        { account: 'Account 1', devices: 3, expired: '2024-12-31', balance: '$100.00', settings: 'Edit', api: 'Access' },
       
        // Data lainnya
      ]
    };
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
