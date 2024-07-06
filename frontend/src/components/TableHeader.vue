<template>
  <div class="table-header">
    <div class="table-row">
      <div class="table-cell" @click="handleAccountClick">{{ account }}</div>
      <div class="table-cell">{{ devicesCount }}</div>
      <div class="table-cell">{{ expiredAt }}</div>
      <div class="table-cell">{{ balance }}</div>
      <div class="table-cell" @click="handleSettingClick">{{ setting }}</div>
      <div class="table-cell">{{ apiVersion }}</div>
    </div>
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
  display: block;
  width: 100%;
  padding: 10px;
  border-radius: 4px;
}

.table-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid #ccc;
  padding: 8px 0;
}

.table-cell {
  flex: 1;
  text-align: center;
  font-weight: bold;
  text-transform: uppercase;
  color: #333;
  cursor: pointer; /* Add cursor pointer for clickable cell */
}
</style>
