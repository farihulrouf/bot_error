<template>
  <div :class="dropdownClasses" @click="toggleDropdown" v-click-outside="closeDropdown">
    <div :class="['flex', { 'justify-center items-center': !profile.name }]">
      <span v-if="profile.name" class="ml-1">{{ profile.name }}</span>
      <svg-icon
        v-if="iconType === 'mdi'"
        :type="iconType"
        :path="profile.comp"
        size="24"
        :fill="iconFill"
      />
    </div>
  </div>
  <div class="dropdown-content absolute w-44 bg-white z-40" v-if="isOpen" style="right: 0;">
    <a
      v-for="item in menuItems"
      :key="item.name"
      href="#"
      @click.prevent="handleClick(item.action)"
    >
      {{ item.name }}
    </a>
  </div>
</template>

<script>
import SvgIcon from "@jamescoyle/vue-icon";
import clickOutside from "../directives/v-click-outside";

export default {
  name: "DropDown",
  components: {
    SvgIcon,
  },
  directives: {
    clickOutside,
  },
  props: {
    menuItems: {
      type: Array,
      required: true,
    },
    profile: {
      type: Object,
      required: true,
    },
    bgColor: {
      type: String,
      default: "bg-gray-500", // Warna latar belakang default (jika tidak disediakan oleh parent)
    },
  },
  data() {
    return {
      isOpen: false,
      iconType: "mdi", // Properti untuk menentukan jenis ikon (misalnya 'mdi', 'material', dll.)
      iconFill: "#000", // Warna latar belakang ikon (hitam)
    };
  },
  computed: {
    dropdownClasses() {
      return [
        "dropdown",
        this.bgColor, // Menambahkan kelas warna latar belakang dari properti bgColor
      ];
    },
  },
  emits: ['menuClick'],  // Deklarasikan event menuClick
  methods: {
    toggleDropdown() {
      this.isOpen = !this.isOpen;
    },
    closeDropdown() { // Metode untuk menutup dropdown
      this.isOpen = false;
    },
    handleClick(action) {
      this.$emit("menu-click", action);
      this.isOpen = false; // Menutup dropdown setelah klik
    },
  },
};
</script>

<style scoped>
.dropdown {
  position: relative;
  display: inline-block;
}
.dropdown-button {
  color: white;
  padding: 10px 20px;
  font-size: 16px;
  border: none;
  cursor: pointer;
}
.dropdown-content {
  display: block;
  box-shadow: 0px 8px 16px 0px rgba(0, 0, 0, 0.2);
  right: 0; /* Menyesuaikan posisi dropdown ke kanan */
}
.dropdown-content a {
  color: black;
  padding: 12px 16px;
  text-decoration: none;
  display: block;
}
.dropdown-content a:hover {
  background-color: #f1f1f1;
}
</style>
