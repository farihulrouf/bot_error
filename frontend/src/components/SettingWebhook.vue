<template>
    <div class="flex items-center justify-between min-h-screen bg-gray-50 p-6">
      
      <div class="w-1/2 pr-6">
        <h2 class="text-xl font-semibold mb-4">Use webhook field</h2>
        <button class="bg-blue-500 hover:bg-blue-600 text-white py-2 px-4 rounded-lg shadow-md mr-4">Save</button>
        <button class="bg-gray-300 hover:bg-gray-400 text-gray-700 py-2 px-4 rounded-lg shadow-md">Close</button>
      </div>
      
      <div class="w-1/2 pl-6">
        <h2 class="text-xl font-semibold mb-4">API Token</h2>
        <div class="bg-white shadow-md rounded-lg p-4">
          <p class="font-medium">Token:</p>
          <div class="flex items-center mt-2">
            <input type="text" class="w-full bg-gray-100 border-gray-300 rounded-md py-2 px-4 focus:outline-none focus:ring-2 focus:ring-blue-500" value="NtENWuyQ7nZMDFVP" readonly>
            <button class="ml-2 bg-blue-500 hover:bg-blue-600 text-white py-2 px-4 rounded-lg shadow-md">Copy</button>
            <button class="ml-2 bg-blue-500 hover:bg-blue-600 text-white py-2 px-4 rounded-lg shadow-md">New</button>
          </div>
        </div>
      </div>
    </div>
  </template>
  
  <script>
  export default {
    name: 'SettingWebhook',
  };
  </script>
  
  <style scoped>
 
  </style>
  