const webpack = require('webpack');
const dotenv = require('dotenv');

// Memuat variabel lingkungan dari file .env
dotenv.config();

module.exports = {
  configureWebpack: {
    plugins: [
      new webpack.DefinePlugin({
        'process.env': Object.keys(process.env).reduce((acc, key) => {
          acc[key] = JSON.stringify(process.env[key]);
          return acc;
        }, {})
      })
    ]
  },
  devServer: {
    port: process.env.PORT || 8080,
    proxy: {
      '/api': {
        target: process.env.VUE_APP_API_URL,
        changeOrigin: true
      }
    }
  }
};
