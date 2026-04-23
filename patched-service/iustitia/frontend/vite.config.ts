// Patch F5 (Secrets leak в JS-бандл).
// Raw-версия содержала:
//   define: {
//     __SERVICE_TOKEN__: JSON.stringify(process.env.VITE_SERVICE_TOKEN ?? ''),
//     __INTERNAL_HMAC_KEY__: JSON.stringify(process.env.VITE_INTERNAL_HMAC_KEY ?? ''),
//   }
// Vite заменяет эти идентификаторы на JSON-литералы на этапе сборки,
// поэтому SECRET_MARKER_B_* и SECRET_MARKER_K_* оказывались в конечном
// /assets/index-*.js и становились публично-читаемыми через nginx.
// Фикс - define-блок удалён; соответствующие TS-модули (shared/config,
// shared/crypto/hmac) заменены на экспорты пустых строк.
import react from '@vitejs/plugin-react';
import path from 'node:path';
import { defineConfig } from 'vite';

export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      '@': path.resolve(__dirname, 'src'),
      '@app': path.resolve(__dirname, 'src/app'),
      '@pages': path.resolve(__dirname, 'src/pages'),
      '@widgets': path.resolve(__dirname, 'src/widgets'),
      '@features': path.resolve(__dirname, 'src/features'),
      '@entities': path.resolve(__dirname, 'src/entities'),
      '@shared': path.resolve(__dirname, 'src/shared'),
    },
  },
  server: {
    host: '0.0.0.0',
    port: 3000,
  },
  build: {
    outDir: 'dist',
    sourcemap: false,
    target: 'es2022',
  },
});
