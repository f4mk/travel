import react from '@vitejs/plugin-react'
import path from 'path'
import { defineConfig } from 'vite'
import check from 'vite-plugin-checker'
import ssr from 'vite-plugin-ssr/plugin'

export default defineConfig({
  plugins: [
    react(),
    check({
      overlay: {
        initialIsOpen: false
      },
      typescript: true,
      eslint: {
        lintCommand: 'eslint "./src/**/*.{ts,tsx}"'
      }
    }),
    ssr()
  ],
  resolve: {
    alias: {
      '#': path.resolve(__dirname, './src')
    }
  }
})
