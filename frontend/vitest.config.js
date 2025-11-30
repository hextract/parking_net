import { defineConfig } from 'vitest/config'
import react from '@vitejs/plugin-react'

export default defineConfig({
  plugins: [react()],
  test: {
    globals: true,
    environment: 'jsdom',
    setupFiles: './src/tests/setup.js',
    css: true,
    exclude: [
      '**/node_modules/**',
      '**/dist/**',
      '**/e2e/**', // Exclude Playwright e2e tests - run them separately with npx playwright test
      '**/.{idea,git,cache,output,temp}/**',
    ],
  },
  resolve: {
    extensions: ['.js', '.jsx', '.json'],
  },
})
