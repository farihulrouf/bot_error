export default {
  beforeMount(el, binding) {
    el.__clickOutsideHandler__ = (event) => {
      if (!(el === event.target || el.contains(event.target))) {
        if (typeof binding.value === 'function') {
          binding.value(event);
        }
      }
    };
    document.body.addEventListener('click', el.__clickOutsideHandler__);
  },
  unmounted(el) {
    document.body.removeEventListener('click', el.__clickOutsideHandler__);
    delete el.__clickOutsideHandler__;
  }
};
