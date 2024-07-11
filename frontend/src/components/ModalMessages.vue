<template>
  <div v-if="show" class="fixed inset-0 flex items-center justify-center bg-gray-900 bg-opacity-50" @click="closeModal">
    <div class="bg-white rounded-lg shadow-lg max-w-lg w-full" @click.stop>
      <!-- Header modal -->
      <div class="flex justify-between items-center p-4 border-b">
        <slot name="header">Default Header</slot>
        <button class="text-gray-500 hover:text-gray-700" @click="closeModal">&times;</button>
      </div>
      <!-- Body modal -->
      <div class="p-4">
        <slot name="body">Default Body</slot>
      </div>
      <!-- Footer modal -->
      <div class="flex justify-end p-4 border-t">
        <slot name="footer">
          <button class="bg-blue-500 text-white px-4 py-2 rounded" @click="closeModal">Close</button>
        </slot>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    show: {
      type: Boolean,
      required: true
    }
  },
  methods: {
    closeModal() {
      this.$emit('close');
    }
  }
}
</script>
