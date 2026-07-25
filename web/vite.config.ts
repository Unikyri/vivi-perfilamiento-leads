import { defineConfig } from 'vite';

export default defineConfig({
  build: {
    outDir: '../internal/adapters/http/estaticos',
    emptyOutDir: true,
  },
  server: {
    port: 5173,
    // En desarrollo, /api va al backend Go local.
    proxy: {
      '/api': { target: 'http://localhost:8080', changeOrigin: true },
      '/salud': { target: 'http://localhost:8080', changeOrigin: true },
    },
  },
});
