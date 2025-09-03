<template>
  <div id="app" class="min-h-screen bg-gray-50">
    <!-- Header principal -->
    <header class="bg-blue-600 shadow-sm border-b border-blue-700">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center h-16">
          <!-- Logo y título -->
          <div class="flex items-center">
            <div class="flex-shrink-0">
              <router-link to="/" class="cursor-pointer hover:opacity-80 transition-opacity">
                <h1 class="text-2xl font-bold text-white">MediTrack</h1>
              </router-link>
            </div>
          </div>
          
          <!-- Barra de búsqueda -->
          <div class="flex-1 max-w-lg mx-8">
            <div class="relative">
              <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              </div>
            </div>
          </div>
          
          <!-- Menú de usuario -->
          <div class="flex items-center space-x-4">
            <!-- Información del usuario autenticado -->
            <div v-if="authStore.isAuthenticated" class="flex items-center space-x-3">
              <div class="text-right text-white">
                <p class="text-sm font-medium">{{ authStore.getUserName }}</p>
                <p class="text-xs text-blue-200">{{ authStore.getUserRole }}</p>
              </div>
              <button 
                @click="handleLogout"
                class="p-2 rounded-lg hover:bg-blue-700 transition-colors"
                title="Cerrar sesión"
              >
                <svg class="h-6 w-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
                </svg>
              </button>
            </div>
            
          </div>
        </div>
      </div>
    </header>

    <!-- Contenido principal -->
    <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
      <router-view />
    </main>

    <!-- Navegación inferior (inspirada en la app móvil) -->
    <nav class="fixed bottom-0 left-0 right-0 bg-white border-t border-gray-200 md:hidden">
      <div class="flex justify-around">
        <router-link
          to="/"
          class="flex flex-col items-center py-3 px-4 text-sm font-medium transition-colors"
          :class="$route.path === '/' ? 'text-primary-600' : 'text-gray-500 hover:text-primary-600'"
        >
          <svg class="h-6 w-6 mb-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
          </svg>
          Inicio
        </router-link>
        
        <router-link
          to="/inventory"
          class="flex flex-col items-center py-3 px-4 text-sm font-medium transition-colors"
          :class="$route.path === '/inventory' ? 'text-primary-600' : 'text-gray-500 hover:text-primary-600'"
        >
          <svg class="h-6 w-6 mb-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
          </svg>
          Inventario
        </router-link>
        
        <router-link
          to="/qr"
          class="flex flex-col items-center py-3 px-4 text-sm font-medium transition-colors"
          :class="$route.path.startsWith('/qr') ? 'text-primary-600' : 'text-gray-500 hover:text-primary-600'"
        >
          <svg class="h-6 w-6 mb-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h.01M5 8h2a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v2a2 2 0 002 2zm0 0h2a2 2 0 012 2v2a2 2 0 01-2 2H5a2 2 0 01-2-2v-2a2 2 0 012-2z" />
          </svg>
          Escáner QR
        </router-link>
        
        <!--<router-link
          to="/statistics"
          class="flex flex-col items-center py-3 px-4 text-sm font-medium transition-colors"
          :class="$route.path === '/statistics' ? 'text-primary-600' : 'text-gray-500 hover:text-primary-600'"
        >
          <svg class="h-6 w-6 mb-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
          </svg>
          Estadísticas
        </router-link> -->
        
        <!--<router-link
          to="/movements"
          class="flex flex-col items-center py-3 px-4 text-sm font-medium transition-colors"
          :class="$route.path === '/movements' ? 'text-primary-600' : 'text-gray-500 hover:text-primary-600'"
        >
          <svg class="h-6 w-6 mb-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
          Movimientos
        </router-link> -->
        
        <!--<router-link
          to="/profile"
          class="flex flex-col items-center py-3 px-4 text-sm font-medium transition-colors"
          :class="$route.path === '/profile' ? 'text-primary-600' : 'text-gray-500 hover:text-primary-600'"
        >
          <svg class="h-6 w-6 mb-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
          </svg>
          Perfil
        </router-link> -->
      </div>
    </nav>

    <!-- Padding para la navegación inferior en móviles -->
    <div class="h-20 md:hidden"></div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()
const searchQuery = ref('')

// Inicializar autenticación al montar el componente
onMounted(() => {
  authStore.initializeAuth()
})

const handleSearch = () => {
  // TODO: Implementar lógica de búsqueda
  console.log('Buscando:', searchQuery.value)
}

const handleLogout = () => {
  authStore.logout()
  router.push('/login')
}
</script>

<style scoped>
/* Estilos específicos del componente */
.router-link-active {
  color: #2563eb; /* Use Tailwind's primary-600 color or your project's equivalent */
}

/* Transiciones suaves para la navegación */
.router-link-active svg {
  transform: scale(1.10);
  transition: transform 0.2s;
}
</style>