import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import './style.css'

// Crear la aplicación Vue
const app = createApp(App)

// Configurar Pinia para manejo de estado
const pinia = createPinia()
app.use(pinia)

// Restaurar sesión antes de montar la app
import { useAuthStore } from '@/stores/auth'
const authStore = useAuthStore()
authStore.initializeAuth()

// Configurar Vue Router
app.use(router)

// Montar la aplicación
app.mount('#app')