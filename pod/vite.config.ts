import { defineConfig } from 'vite'
import react from '@vitejs/plugin-react'
import { resolve } from 'path'

export default defineConfig({
  plugins: [react()],
  build: {
    outDir: '../frontend/dist/app',
    lib: {
      entry: resolve(__dirname, 'src/pod-wc.tsx'),
      name: 'PodWebComponent',
      fileName: 'pod',
      formats: ['es']
    }
  }
})
