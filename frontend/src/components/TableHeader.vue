<template>
  <div class="table-header mt-4 mb-4">
    <table>
      <thead class="bg-white">
        <tr>
          <th>ACCOUNT</th>
          <th>DEVICES</th>
          <th>EXPIRED</th>
          <th>BALANCE</th>
          <th>SETTING</th>
          <th>API</th>
        </tr>
      </thead>
      <tbody>
        <tr class="table-row">
          <td @click="handleAccountClick" class="table-cell">{{ account }}</td>
          <td class="table-cell">{{ devicesCount }}</td>
          <td class="table-cell">{{ expiredAt }}</td>
          <td class="table-cell">{{ balance }}</td>
          <td @click="handleSettingClick" class="table-cell">{{ setting }}</td>
          <td class="table-cell">{{ apiVersion }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script>
export default {
  name: 'TableHeader',
  props: {
    account: {
      type: String,
      required: true
    },
    devicesCount: {
      type: Number,
      required: true
    },
    expiredAt: {
      type: String,
      required: true
    },
    balance: {
      type: String,
      required: true
    },
    setting: {
      type: String,
      required: true
    },
    apiVersion: {
      type: String,
      required: true
    }
  },
  methods: {
    handleAccountClick() {
      this.$emit('show-update-profile');
    },
    handleSettingClick() {
      this.$emit('toggle-token-webhook');
    }
  }
};
</script>

<style scoped>
.table-header {
  width: 100%;
  overflow-x: auto; /* Membuat tabel horizontal scroll jika terlalu lebar */
}

table {
  width: 100%;
  border-collapse: collapse; /* Gabungkan garis-garis tabel */
}

th, td {
  padding: 10px;
  text-align: center;
  font-weight: bold;
  text-transform: uppercase;
  color: #6b6868;
  cursor: pointer; /* Add cursor pointer for clickable cells */
  border-bottom: 1px solid #ccc; /* Garis bawah untuk setiap baris */
}

th {
  background-color: #f3f4f6; /* Warna latar belakang header */
}

.table-cell {
  min-width: 100px; /* Lebar minimum untuk setiap sel */
}
</style>
