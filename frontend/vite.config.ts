import react from '@vitejs/plugin-react'
import path from 'path'
import { defineConfig } from 'vite'
import check from 'vite-plugin-checker'

// https://vitejs.dev/config/
export default defineConfig({
  server: {
    proxy: {
      '/api': {
        target: 'https://localhost',
        changeOrigin: true,
        secure: false,
      },
    },
  },
  plugins: [
    react(),
    check({
      overlay: {
        initialIsOpen: false,
      },
      typescript: true,
      eslint: {
        lintCommand: 'eslint "./src/**/*.{ts,tsx}"',
      },
    }),
  ],
  resolve: {
    alias: {
      '#': path.resolve(__dirname, './src'),
    },
  },
})
