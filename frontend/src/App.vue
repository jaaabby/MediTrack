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
          
          <!-- Navegación principal (desktop) -->
          <nav class="hidden md:flex space-x-8">
            <router-link
              to="/"
              class="text-white hover:text-blue-200 px-3 py-2 rounded-md text-sm font-medium transition-colors"
              :class="{ 'bg-blue-700': $route.path === '/' }"
            >
              Inicio
            </router-link>
            
            <router-link
              to="/inventory"
              class="text-white hover:text-blue-200 px-3 py-2 rounded-md text-sm font-medium transition-colors"
              :class="{ 'bg-blue-700': $route.path.startsWith('/inventory') }"
            >
              Inventario
            </router-link>
            
            <router-link
              to="/supply-requests"
              class="text-white hover:text-blue-200 px-3 py-2 rounded-md text-sm font-medium transition-colors"
              :class="{ 'bg-blue-700': $route.path.startsWith('/supply-requests') }"
            >
              Solicitudes
            </router-link>
            
            <router-link
              to="/qr"
              class="text-white hover:text-blue-200 px-3 py-2 rounded-md text-sm font-medium transition-colors"
              :class="{ 'bg-blue-700': $route.path.startsWith('/qr') }"
            >
              Escáner QR
            </router-link>
            
            <router-link
              to="/reports"
              class="text-white hover:text-blue-200 px-3 py-2 rounded-md text-sm font-medium transition-colors"
              :class="{ 'bg-blue-700': $route.path.startsWith('/reports') }"
            >
              Reportes
            </router-link>
          </nav>
          
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
            
            <!-- Notificaciones (futuro) -->
            <button class="p-2 rounded-lg hover:bg-blue-700 transition-colors relative">
              <svg class="h-6 w-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-5-5 5-5H15m-5-3v5l-3 3 3 3v-5h5" />
              </svg>
              <!-- Badge de notificaciones -->
              <!-- <span class="absolute top-0 right-0 block h-2 w-2 rounded-full bg-red-400"></span> -->
            </button>
            
            <!-- Menú hamburguesa (móvil) -->
            <button 
              @click="mobileMenuOpen = !mobileMenuOpen"
              class="md:hidden p-2 rounded-lg hover:bg-blue-700 transition-colors"
            >
              <svg class="h-6 w-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path v-if="!mobileMenuOpen" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
                <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
        </div>
      </div>

      <!-- Menú móvil -->
      <div v-if="mobileMenuOpen" class="md:hidden">
        <div class="px-2 pt-2 pb-3 space-y-1 sm:px-3 bg-blue-700">
          <router-link
            to="/"
            @click="mobileMenuOpen = false"
            class="text-white hover:text-blue-200 block px-3 py-2 rounded-md text-base font-medium"
            :class="{ 'bg-blue-800': $route.path === '/' }"
          >
            Inicio
          </router-link>
          
          <router-link
            to="/inventory"
            @click="mobileMenuOpen = false"
            class="text-white hover:text-blue-200 block px-3 py-2 rounded-md text-base font-medium"
            :class="{ 'bg-blue-800': $route.path.startsWith('/inventory') }"
          >
            Inventario
          </router-link>
          
          <router-link
            to="/supply-requests"
            @click="mobileMenuOpen = false"
            class="text-white hover:text-blue-200 block px-3 py-2 rounded-md text-base font-medium"
            :class="{ 'bg-blue-800': $route.path.startsWith('/supply-requests') }"
          >
            Solicitudes de Insumo
          </router-link>
          
          <router-link
            to="/qr"
            @click="mobileMenuOpen = false"
            class="text-white hover:text-blue-200 block px-3 py-2 rounded-md text-base font-medium"
            :class="{ 'bg-blue-800': $route.path.startsWith('/qr') }"
          >
            Escáner QR
          </router-link>
          
          <router-link
            to="/reports"
            @click="mobileMenuOpen = false"
            class="text-white hover:text-blue-200 block px-3 py-2 rounded-md text-base font-medium"
            :class="{ 'bg-blue-800': $route.path.startsWith('/reports') }"
          >
            Reportes
          </router-link>
        </div>
      </div>
    </header>

    <!-- Contenido principal -->
    <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
      <router-view />
    </main>

    <!-- Navegación inferior (inspirada en la app móvil) -->
    <nav class="fixed bottom-0 left-0 right-0 bg-white border-t border-gray-200 md:hidden z-40">
      <div class="flex justify-around">
        <router-link
          to="/"
          class="flex flex-col items-center py-3 px-4 text-sm font-medium transition-colors"
          :class="$route.path === '/' ? 'text-blue-600' : 'text-gray-500 hover:text-blue-600'"
        >
          <svg class="h-6 w-6 mb-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
          </svg>
          <span class="text-xs">Inicio</span>
        </router-link>
        
        <router-link
          to="/inventory"
          class="flex flex-col items-center py-3 px-4 text-sm font-medium transition-colors"
          :class="$route.path.startsWith('/inventory') ? 'text-blue-600' : 'text-gray-500 hover:text-blue-600'"
        >
          <svg class="h-6 w-6 mb-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
          </svg>
          <span class="text-xs">Inventario</span>
        </router-link>
        
        <router-link
          to="/supply-requests"
          class="flex flex-col items-center py-3 px-4 text-sm font-medium transition-colors relative"
          :class="$route.path.startsWith('/supply-requests') ? 'text-blue-600' : 'text-gray-500 hover:text-blue-600'"
        >
          <svg class="h-6 w-6 mb-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01" />
          </svg>
          <span class="text-xs">Solicitudes</span>
          <!-- Badge para solicitudes pendientes (futuro) -->
          <!-- <span v-if="pendingRequests > 0" class="absolute top-1 right-2 block h-4 w-4 rounded-full bg-red-500 text-white text-xs flex items-center justify-center">{{ pendingRequests }}</span> -->
        </router-link>
        
        <router-link
          to="/qr"
          class="flex flex-col items-center py-3 px-4 text-sm font-medium transition-colors"
          :class="$route.path.startsWith('/qr') ? 'text-blue-600' : 'text-gray-500 hover:text-blue-600'"
        >
          <svg class="h-6 w-6 mb-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h.01M5 8h2a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v2a2 2 0 002 2zm0 0h2a2 2 0 012 2v2a2 2 0 01-2 2H5a2 2 0 01-2-2v-2a2 2 0 012-2z" />
          </svg>
          <span class="text-xs">Escáner</span>
        </router-link>
        
        <router-link
          to="/reports"
          class="flex flex-col items-center py-3 px-4 text-sm font-medium transition-colors"
          :class="$route.path.startsWith('/reports') ? 'text-blue-600' : 'text-gray-500 hover:text-blue-600'"
        >
          <svg class="h-6 w-6 mb-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
          </svg>
          <span class="text-xs">Reportes</span>
        </router-link>
      </div>
    </nav>

    <!-- Padding para la navegación inferior en móviles -->
    <div class="h-20 md:hidden"></div>

    <!-- Toast/Notificación global (futuro) -->
    <!-- <div v-if="globalNotification" class="fixed bottom-4 right-4 bg-green-500 text-white p-4 rounded-lg shadow-lg z-50">
      {{ globalNotification }}
    </div> -->
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()
const searchQuery = ref('')
const mobileMenuOpen = ref(false)

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
// Estado para notificaciones (futuro)
// const globalNotification = ref(null)
// const pendingRequests = ref(0)

// Cerrar menú móvil cuando se cambia de ruta
router.afterEach(() => {
  mobileMenuOpen.value = false
})

// Métodos para futuras funcionalidades
// const showNotification = (message) => {
//   globalNotification.value = message
//   setTimeout(() => {
//     globalNotification.value = null
//   }, 5000)
// }

// const loadPendingRequestsCount = async () => {
//   // Cargar número de solicitudes pendientes
// }

// Lifecycle hooks para futuras funcionalidades
// onMounted(() => {
//   loadPendingRequestsCount()
// })
</script>

<style scoped>
/* Estilos específicos del componente */
.router-link-active {
  color: #2563eb; /* text-blue-600 */
}

/* Transiciones suaves para la navegación */
.router-link-active svg {
  transform: scale(1.05);
  transition: transform 0.2s;
}

/* Hover effects mejorados */
.transition-colors {
  transition: color 0.2s ease-in-out, background-color 0.2s ease-in-out;
}

/* Efecto de resaltado para navegación activa */
.bg-blue-700 {
  background-color: #1d4ed8;
}

.bg-blue-800 {
  background-color: #1e40af;
}

/* Sombras para el header */
.shadow-sm {
  box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
}

/* Z-index para elementos superpuestos */
.z-40 {
  z-index: 40;
}

.z-50 {
  z-index: 50;
}

/* Estilos para el menú móvil */
@media (max-width: 768px) {
  .md\:hidden {
    display: block;
  }
}

/* Animación suave para el menú hamburguesa */
.md\:hidden svg {
  transition: transform 0.3s ease-in-out;
}

/* Badge styles para notificaciones futuras */
.absolute {
  position: absolute;
}

.top-1 {
  top: 0.25rem;
}

.right-2 {
  right: 0.5rem;
}

.h-4 {
  height: 1rem;
}

.w-4 {
  width: 1rem;
}

.rounded-full {
  border-radius: 9999px;
}

.bg-red-500 {
  background-color: #ef4444;
}

.text-white {
  color: #ffffff;
}

.text-xs {
  font-size: 0.75rem;
}

.flex {
  display: flex;
}

.items-center {
  align-items: center;
}

.justify-center {
  justify-content: center;
}
</style>