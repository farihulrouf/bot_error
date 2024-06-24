<template>
    <div class="relative w-[400px] overflow-hidden bg-gray-50 py-4 ">
      <div
        class="group relative cursor-pointer overflow-hidden bg-white px-6 pt-10 pb-8 shadow-xl ring-1 ring-gray-900/5 transition-all duration-300 hover:-translate-y-1 hover:shadow-2xl sm:mx-auto sm:max-w-sm sm:rounded-lg sm:px-10"
      >
        <span
          class="absolute top-10 z-0 h-20 w-20 rounded-full bg-sky-500 transition-all duration-300 group-hover:scale-[10]"
        ></span>
        <div class="relative z-10 mx-auto max-w-md">
          <span
            class="grid h-20 w-20 place-items-center rounded-full bg-sky-500 transition-all duration-300 group-hover:bg-sky-400"
          >
            
          </span>
          <div
            class="space-y-6 pt-5 text-base leading-7 text-gray-600 transition-all duration-300 group-hover:text-white/90"
          >
            <p>
                device-1
            </p>
            <p>
                Phone: 6289534449417
            </p>
            <p>
                Name: Device1
            </p>
            <p>
                Status: ready
            </p>
            <p>
                Process: getMessages
            </p>
          </div>
          <div class="pt-5 text-base font-semibold leading-7">
           
          </div>
        </div>
      </div>
    </div>
  </template>
  
  <script>

  export default {
    name: 'QrcodeCard'
  };
  </script>
  
  <style scoped>
  /* CSS menggunakan Tailwind CSS atau CSS kustom */
  </style>
  