/** @type {import('tailwindcss').Config} */
module.exports = {
  purge: ['./index.html', './src/**/*.{vue,js,ts,jsx,tsx}'], // Menggunakan purge untuk menghapus kode yang tidak digunakan di produksi
  darkMode: false, 
  theme: {
    extend: {
      colors: {
        whatsapp_teal: '#128c7e', // Menambahkan warna khusus untuk WhatsApp
        whatsapp_primary: '#25D366',
        //whatsapp_third: '#34B7F1',
        whatsapp_thidr: '#34B7F1'
      },
    },
  },
  variants: {
    extend: {}, 
  },
  plugins: [], 
};
