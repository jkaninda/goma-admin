import Aura from '@primeuix/themes/aura';

export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },
  runtimeConfig: {
    public: {
      appUrl: '', 
      apiUrl: ''
    }
  },
  modules: ['@primevue/nuxt-module', '@pinia/nuxt'],
  primevue: {
    options: {
      theme: {
        preset: Aura
      }
    }
  },
  pinia: {
    storesDirs: ['./stores/**'],
  },
})