import { useNotification } from '@kyvg/vue3-notification';

// Fungsi untuk menampilkan notifikasi
export function showNotification(title, text, type = 'success', duration = 3000) {
  const { notify } = useNotification();

  notify({
    title: title,
    text: text,
    type: type,
    duration: duration,
  });
}
