import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import { readFileSync } from 'fs'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [vue()],
  server: {
    https: {
      key: readFileSync('../Golang/ssl/localhost/localhost.decrypted.key'),
      cert: readFileSync('../Golang/ssl/localhost/localhost.crt'),
    },
    port: 10001
  }
})
