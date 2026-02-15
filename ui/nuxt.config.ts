import Aura from '@primeuix/themes/aura';
import tailwindcss from "@tailwindcss/vite";


export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },
  runtimeConfig: {
    public: {
      appUrl: '', 
      apiUrl: '',
      appName: '',
      appVersion:''
    }
  },
  modules: [
    '@primevue/nuxt-module',
    '@pinia/nuxt',
    '@nuxtjs/color-mode',
    '@nuxtjs/i18n',
  ],
  primevue: {
    options: {
      theme: {
        preset: Aura,
        options: {
          prefix: 'p',
          darkModeSelector: '.dark',
        },
        components: {
          exclude: []
        }
      },
      
    }
  },
  css: [
    "./app/assets/css/main.css",
    'primeicons/primeicons.css'
  ],
  vite: {
    plugins: [
      tailwindcss(),
    ],
  },
  nitro: {
    routeRules: {
      '/_nuxt/**': { cors: true }
    }
  },
  pinia: {
    storesDirs: ['./stores/**'],
  },
  colorMode: {
    preference: 'system',
    fallback: 'light', 
  },
  app: {
    baseURL: '/',
  },
  i18n: {
    baseUrl: "http://localhost:3000",
    strategy: "prefix_except_default",
    defaultLocale: "en",
    locales: [
      { code: "en", iso: "en-US", file: "en.json", name: "English", language: "en-US" },
      //{ code: "fr", iso: "fr-FR", file: "fr.json", name: "Français", language: "fr-FR" },
    ],
    experimental: {
      strictSeo: true,
    },
    detectBrowserLanguage: {
      useCookie: true,
      cookieKey: "i18n_redirected",
      redirectOn: "root",
    },
  },
})