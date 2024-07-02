import axios from 'axios';

const baseURL = 'http://localhost:8080/api';  

const api = axios.create({
  baseURL: baseURL,
  timeout: 10000,  // Timeout 10 detik
  headers: {
    'Content-Type': 'application/json',
  },
});


// Interceptor untuk mengatur Authorization header dengan token JWT jika tersedia
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => {
    return Promise.reject(error);
  }
);
////https://chatgpt.com/c/46121827-36e9-4c0f-a9cf-7f2db1ea9b5e

// Handler untuk mengambil error response dari server
const errorHandler = (error) => {
  if (error.response) {
    // Tangani error dari response server (response tidak dalam range 2xx)
    return Promise.reject(error.response.data);
  } else if (error.request) {
    // Tangani error dari request tanpa response dari server (misalnya timeout atau jaringan down)
    return Promise.reject({ message: 'Network Error', error });
  } else {
    // Tangani error lainnya
    return Promise.reject(error);
  }
};

export default {
  // Method untuk mengirimkan request GET
  get(url) {
    return api.get(url).then(response => response.data).catch(errorHandler);
  },

  // Method untuk mengirimkan request POST
  post(url, data) {
    return api.post(url, data).then(response => response.data).catch(errorHandler);
  },

  // Method untuk mengirimkan request PUT
  put(url, data) {
    return api.put(url, data).then(response => response.data).catch(errorHandler);
  },

  // Method untuk mengirimkan request DELETE
  delete(url) {
    return api.delete(url).then(response => response.data).catch(errorHandler);
  },

  // Method untuk mengambil baseURL
  getBaseURL() {
    return baseURL;
  },
};
