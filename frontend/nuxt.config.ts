export default defineNuxtConfig({
  compatibilityDate: '2026-07-08',
  
  modules: [
    '@nuxt/ui',
    '@pinia/nuxt'
  ],

  devtools: {
    enabled: true
  },

  css: [
    '~/assets/css/main.css'
  ]
})