// vite.config.ts
import { defineConfig } from 'vitest/config'
import react from '@vitejs/plugin-react'
import { TanStackRouterVite } from '@tanstack/router-plugin/vite'
// https://vitejs.dev/config/
export default defineConfig({
        test: {
                environment: "jsdom",
                globals: true,
                setupFiles: "./src/testSetup.ts",
        },
        plugins: [
                TanStackRouterVite({ target: 'react', autoCodeSplitting: true }),
                react(),
        ],
        server: {
                host: '127.0.0.1',
                port: 5173,
                proxy: {
                        '/api':{
                                target: 'http://127.0.0.1:8081',
                                changeOrigin: true,
                                rewrite: path => path.replace(/^\/api/, '')
                                

                        }
                }
        },
        
});