<template>
  <div class="items-center justify-center">
    <div class="overflow-x-auto">
      <table class="min-w-full bg-white shadow-md rounded-xl">
        <thead>
          <tr class="bg-blue-gray-100 text-gray-700">
            <th class="py-3 px-4 text-left">Account</th>
            <th class="py-3 px-4 text-left">Devices</th>
            <th class="py-3 px-4 text-left">Expired</th>
            <th class="py-3 px-4 text-left">Balance</th>
            <th class="py-3 px-4 text-left">Settings</th>
            <th class="py-3 px-4 text-left">API</th>
            <th class="py-3 px-4 text-left">Action</th>
          </tr>
        </thead>
        <tbody class="text-blue-gray-900">
          <!-- Single Row with Data -->
          <tr class="border-b border-blue-gray-200">
            <td class="py-3 px-4">Multi</td>
            <td class="py-3 px-4">2</td>
            <td class="py-3 px-4">2030-12-12</td>
            <td class="py-3 px-4">$300,000.00</td>
            <td class="py-3 px-4">Token, Webhook</td>
            <td class="py-3 px-4">API v1.4.3</td>
            <td class="py-3 px-4">
              <a href="#" class="font-medium text-blue-600 hover:text-blue-800">Edit</a>
            </td>
          </tr>
        </tbody>
      </table>
     
    </div>
  </div>
</template>

<script>
export default {
  name: 'ListAccount',
  props: {
    headers: {
      type: Array,
      required: true
    }
  }
};
</script>

<style scoped>
/* Styling menggunakan Tailwind CSS atau CSS kustom */
</style>
