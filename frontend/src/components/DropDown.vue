<template>
  <div :class="dropdownClasses">
    <div class="dropdown-button" @click="toggleDropdown">
      <svg-icon v-if="iconType === 'mdi'" :type="iconType" :path="profile.comp" size="24" :fill="iconFill" />
    </div>
    <div class="dropdown-content" v-if="isOpen">
      <a v-for="item in menuItems" :key="item.name" href="#" @click.prevent="handleClick(item.action)">
        {{ item.name }}
      </a>
    </div>
  </div>
</template>

<script>
import SvgIcon from '@jamescoyle/vue-icon';

export default {
  name: 'DropDown',
  components: {
    SvgIcon
  },
  props: {
    menuItems: {
      type: Array,
      required: true
    },
    profile: {
      type: Object,
      required: true
    },
    bgColor: {
      type: String,
      default: 'bg-gray-500' // Default background color (if not provided by parent)
    }
  },
  data() {
    return {
      isOpen: false,
      iconType: 'mdi', // Properti untuk menentukan jenis ikon (misalnya 'mdi', 'material', dll.)
      iconFill: '#000' // Warna latar belakang ikon (hitam)
    };
  },
  computed: {
    dropdownClasses() {
      return [
        'dropdown',
        this.bgColor // Menambahkan kelas warna latar belakang dari properti bgColor
      ];
    }
  },
  methods: {
    toggleDropdown() {
      this.isOpen = !this.isOpen;
    },
    handleClick(action) {
      this.$emit('menu-click', action);
      this.isOpen = false; // Menutup dropdown setelah klik
    }
  }
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
  position: absolute;
  background-color: #f9f9f9;
  min-width: 160px;
  box-shadow: 0px 8px 16px 0px rgba(0,0,0,0.2);
  z-index: 1;
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
