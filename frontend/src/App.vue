<template>
  <div id="app" class="min-h-screen bg-gray-50">
    <!-- Header principal -->
    <header class="bg-blue-600 shadow-sm border-b border-blue-700">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center h-16">
          <!-- Logo y título -->
          <div class="flex items-center flex-shrink-0">
            <router-link to="/" class="cursor-pointer hover:opacity-80 transition-opacity">
              <h1 class="text-xl sm:text-2xl font-bold text-white">MediTrack</h1>
            </router-link>
          </div>
          
          <!-- Navegación principal (desktop y tablet) -->
          <nav class="hidden lg:flex items-center space-x-4 xl:space-x-8 flex-1 justify-center max-w-4xl">
            <router-link v-if="authStore.isAuthenticated"
              to="/"
              class="text-white hover:text-blue-200 px-2 xl:px-3 py-2 rounded-md text-sm font-medium transition-colors whitespace-nowrap"
              :class="{ 'bg-blue-700': $route.path === '/' }"
              @click.stop
            >
              Inicio
            </router-link>
            
            <router-link v-if="authStore.isAuthenticated && authStore.canViewInventory"
              to="/inventory"
              class="text-white hover:text-blue-200 px-2 xl:px-3 py-2 rounded-md text-sm font-medium transition-colors whitespace-nowrap"
              :class="{ 'bg-blue-700': $route.path === '/inventory' }"
              @click.stop
            >
              Inventario
            </router-link>
            
            <router-link v-if="authStore.isAuthenticated && authStore.canViewRequests"
              to="/supply-requests"
              class="text-white hover:text-blue-200 px-2 xl:px-3 py-2 rounded-md text-sm font-medium transition-colors whitespace-nowrap"
              :class="{ 'bg-blue-700': $route.path.startsWith('/supply-requests') }"
              @click.stop
            >
              Solicitudes
            </router-link>
            
            <router-link v-if="authStore.isAuthenticated"
              to="/qr"
              class="text-white hover:text-blue-200 px-2 xl:px-3 py-2 rounded-md text-sm font-medium transition-colors whitespace-nowrap"
              :class="{ 'bg-blue-700': $route.path.startsWith('/qr') }"
              @click.stop
            >
              Escaner QR
            </router-link>

            <router-link v-if="authStore.isAuthenticated"
              to="/statistics"
              class="text-white hover:text-blue-200 px-2 xl:px-3 py-2 rounded-md text-sm font-medium transition-colors whitespace-nowrap"
              :class="{ 'bg-blue-700': $route.path === '/statistics' }"
              @click.stop
            >
              Estadisticas
            </router-link>
            
            <router-link v-if="authStore.isAuthenticated && authStore.canViewInventory"
              to="/inventory/add"
              class="text-white hover:text-blue-200 px-2 xl:px-3 py-2 rounded-md text-sm font-medium transition-colors whitespace-nowrap"
              :class="{ 'bg-blue-700': $route.path === '/inventory/add' }"
              @click.stop
            >
              Agregar Insumo
            </router-link>
          </nav>

          <!-- Navegación tablet (oculta en desktop y mobile) -->
          <nav class="hidden md:flex lg:hidden items-center space-x-2 flex-1 justify-center">
            <router-link v-if="authStore.isAuthenticated"
              to="/"
              class="text-white hover:text-blue-200 p-2 rounded-md transition-colors"
              :class="{ 'bg-blue-700': $route.path === '/' }"
              @click.stop
              title="Inicio"
            >
              <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
              </svg>
            </router-link>
            
            <router-link v-if="authStore.isAuthenticated && authStore.canViewInventory"
              to="/inventory"
              class="text-white hover:text-blue-200 p-2 rounded-md transition-colors"
              :class="{ 'bg-blue-700': $route.path === '/inventory' }"
              @click.stop
              title="Inventario"
            >
              <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
              </svg>
            </router-link>
            
            <router-link v-if="authStore.isAuthenticated && authStore.canViewRequests"
              to="/supply-requests"
              class="text-white hover:text-blue-200 p-2 rounded-md transition-colors"
              :class="{ 'bg-blue-700': $route.path.startsWith('/supply-requests') }"
              @click.stop
              title="Solicitudes"
            >
              <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01" />
              </svg>
            </router-link>
            
            <router-link v-if="authStore.isAuthenticated"
              to="/qr"
              class="text-white hover:text-blue-200 p-2 rounded-md transition-colors"
              :class="{ 'bg-blue-700': $route.path.startsWith('/qr') }"
              @click.stop
              title="Escaner QR"
            >
              <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h.01M5 8h2a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v2a2 2 0 002 2zm0 0h2a2 2 0 012 2v2a2 2 0 01-2 2H5a2 2 0 01-2-2v-2a2 2 0 012-2z" />
              </svg>
            </router-link>

            <router-link v-if="authStore.isAuthenticated"
              to="/statistics"
              class="text-white hover:text-blue-200 p-2 rounded-md transition-colors"
              :class="{ 'bg-blue-700': $route.path === '/statistics' }"
              @click.stop
              title="Estadisticas"
            >
              <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
              </svg>
            </router-link>
            
            <router-link v-if="authStore.isAuthenticated && authStore.canViewInventory"
              to="/inventory/add"
              class="text-white hover:text-blue-200 p-2 rounded-md transition-colors"
              :class="{ 'bg-blue-700': $route.path === '/inventory/add' }"
              @click.stop
              title="Agregar Insumo"
            >
              <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
              </svg>
            </router-link>
          </nav>
          
          <!-- Menu de usuario -->
          <div class="flex items-center space-x-2 sm:space-x-4 flex-shrink-0">
            <!-- Informacion del usuario autenticado -->
            <div v-if="authStore.isAuthenticated" class="flex items-center space-x-2 sm:space-x-3">
              <div class="text-right text-white hidden sm:block">
                <router-link 
                  to="/profile"
                  class="text-sm font-medium hover:text-blue-200 transition-colors cursor-pointer block"
                  title="Ver mi perfil"
                >
                  {{ authStore.getUserName }}
                </router-link>
                <p class="text-xs text-blue-200">{{ authStore.getUserRole }}</p>
              </div>
              
              <!-- Icono de usuario para pantallas pequeñas -->
              <router-link 
                to="/profile"
                class="sm:hidden p-2 rounded-lg hover:bg-blue-700 transition-colors"
                title="Ver mi perfil"
              >
                <svg class="h-6 w-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                </svg>
              </router-link>
              
              <button 
                @click="handleLogout"
                class="p-2 rounded-lg hover:bg-blue-700 transition-colors"
                title="Cerrar sesion"
              >
                <svg class="h-5 w-5 sm:h-6 sm:w-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
                </svg>
              </button>
            </div>
            
            <!-- Menu hamburguesa (movil) -->
            <button 
              @click="mobileMenuOpen = !mobileMenuOpen"
              class="md:hidden p-2 rounded-lg hover:bg-blue-700 transition-colors"
              aria-label="Abrir menu"
            >
              <svg class="h-6 w-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path v-if="!mobileMenuOpen" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
                <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
        </div>
      </div>

      <!-- Menu movil -->
      <div v-if="mobileMenuOpen" class="md:hidden">
        <div class="px-2 pt-2 pb-3 space-y-1 sm:px-3 bg-blue-700">
          <router-link
            v-if="authStore.isAuthenticated"
            to="/"
            @click.stop="mobileMenuOpen = false"
            class="text-white hover:text-blue-200 block px-3 py-2 rounded-md text-base font-medium transition-colors"
            :class="{ 'bg-blue-800': $route.path === '/' }"
          >
            Inicio
          </router-link>
          
          <router-link
            v-if="authStore.isAuthenticated && authStore.canViewInventory"
            to="/inventory"
            @click.stop="mobileMenuOpen = false"
            class="text-white hover:text-blue-200 block px-3 py-2 rounded-md text-base font-medium transition-colors"
            :class="{ 'bg-blue-800': $route.path === '/inventory' }"
          >
            Inventario
          </router-link>
          
          <router-link
            v-if="authStore.isAuthenticated && authStore.canViewRequests"
            to="/supply-requests"
            @click.stop="mobileMenuOpen = false"
            class="text-white hover:text-blue-200 block px-3 py-2 rounded-md text-base font-medium transition-colors"
            :class="{ 'bg-blue-800': $route.path.startsWith('/supply-requests') }"
          >
            Solicitudes de Insumo
          </router-link>
          
          <router-link
            v-if="authStore.isAuthenticated"
            to="/qr"
            @click.stop="mobileMenuOpen = false"
            class="text-white hover:text-blue-200 block px-3 py-2 rounded-md text-base font-medium transition-colors"
            :class="{ 'bg-blue-800': $route.path.startsWith('/qr') }"
          >
            Escaner QR
          </router-link>

          <router-link
            v-if="authStore.isAuthenticated"
            to="/statistics"
            @click.stop="mobileMenuOpen = false"
            class="text-white hover:text-blue-200 block px-3 py-2 rounded-md text-base font-medium transition-colors"
            :class="{ 'bg-blue-800': $route.path === '/statistics' }"
          >
            Estadisticas
          </router-link>
          
          <router-link
            v-if="authStore.isAuthenticated"
            to="/reports"
            @click.stop="mobileMenuOpen = false"
            class="text-white hover:text-blue-200 block px-3 py-2 rounded-md text-base font-medium transition-colors"
            :class="{ 'bg-blue-800': $route.path.startsWith('/reports') }"
          >
            Reportes
          </router-link>
          
          <router-link
            v-if="authStore.isAuthenticated && authStore.canViewInventory"
            to="/inventory/add"
            @click.stop="mobileMenuOpen = false"
            class="text-white hover:text-blue-200 block px-3 py-2 rounded-md text-base font-medium transition-colors"
            :class="{ 'bg-blue-800': $route.path === '/inventory/add' }"
          >
            Agregar Insumo
          </router-link>
        </div>
      </div>
    </header>

    <!-- Contenido principal -->
    <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
      <router-view />
    </main>

    <!-- Navegacion inferior (inspirada en la app movil) -->
    <nav class="fixed bottom-0 left-0 right-0 bg-white border-t border-gray-200 md:hidden z-40">
      <div class="grid grid-cols-5 gap-1">
        <router-link
          v-if="authStore.isAuthenticated"
          to="/"
          class="flex flex-col items-center py-2 px-1 text-xs font-medium transition-colors"
          :class="$route.path === '/' ? 'text-blue-600' : 'text-gray-500 hover:text-blue-600'"
        >
          <svg class="h-6 w-6 mb-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
          </svg>
          <span class="text-xs leading-tight">Inicio</span>
        </router-link>
        
        <router-link
          v-if="authStore.isAuthenticated && authStore.canViewInventory"
          to="/inventory"
          class="flex flex-col items-center py-2 px-1 text-xs font-medium transition-colors"
          :class="$route.path === '/inventory' ? 'text-blue-600' : 'text-gray-500 hover:text-blue-600'"
        >
          <svg class="h-6 w-6 mb-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
          </svg>
          <span class="text-xs leading-tight">Inventario</span>
        </router-link>
        
        <router-link
          v-if="authStore.isAuthenticated && authStore.canViewRequests"
          to="/supply-requests"
          class="flex flex-col items-center py-2 px-1 text-xs font-medium transition-colors relative"
          :class="$route.path.startsWith('/supply-requests') ? 'text-blue-600' : 'text-gray-500 hover:text-blue-600'"
        >
          <svg class="h-6 w-6 mb-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01" />
          </svg>
          <span class="text-xs leading-tight">Solicitudes</span>
        </router-link>
        
        <router-link
          v-if="authStore.isAuthenticated"
          to="/qr"
          class="flex flex-col items-center py-2 px-1 text-xs font-medium transition-colors"
          :class="$route.path.startsWith('/qr') ? 'text-blue-600' : 'text-gray-500 hover:text-blue-600'"
        >
          <svg class="h-6 w-6 mb-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h.01M5 8h2a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v2a2 2 0 002 2zm0 0h2a2 2 0 012 2v2a2 2 0 01-2 2H5a2 2 0 01-2-2v-2a2 2 0 012-2z" />
          </svg>
          <span class="text-xs leading-tight">Escaner</span>
        </router-link>

        <router-link
          v-if="authStore.isAuthenticated"
          to="/statistics"
          class="flex flex-col items-center py-2 px-1 text-xs font-medium transition-colors"
          :class="$route.path === '/statistics' ? 'text-blue-600' : 'text-gray-500 hover:text-blue-600'"
        >
          <svg class="h-6 w-6 mb-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
          </svg>
          <span class="text-xs leading-tight">Estadisticas</span>
        </router-link>
      </div>
    </nav>

    <!-- Padding para la navegacion inferior en moviles -->
    <div class="h-20 md:hidden"></div>

    <!-- Toast/Notificacion global (futuro) -->
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

// Inicializar autenticacion al montar el componente
onMounted(() => {
  authStore.initializeAuth()
})

const handleSearch = () => {
  // TODO: Implementar logica de busqueda
  console.log('Buscando:', searchQuery.value)
}

const handleLogout = () => {
  authStore.logout()
  router.push('/login')
}
// Estado para notificaciones (futuro)
// const globalNotification = ref(null)
// const pendingRequests = ref(0)

// Cerrar menu movil cuando se cambia de ruta
router.afterEach(() => {
  mobileMenuOpen.value = false
})

// Metodos para futuras funcionalidades
// const showNotification = (message) => {
//   globalNotification.value = message
//   setTimeout(() => {
//     globalNotification.value = null
//   }, 5000)
// }

// const loadPendingRequestsCount = async () => {
//   // Cargar numero de solicitudes pendientes
// }

// Lifecycle hooks para futuras funcionalidades
// onMounted(() => {
//   loadPendingRequestsCount()
// })
</script>

<style scoped>
/* Estilos especificos del componente */
.router-link-active {
  color: #2563eb; /* text-blue-600 */
}

/* Transiciones suaves para la navegacion */
.router-link-active svg {
  transform: scale(1.05);
  transition: transform 0.2s;
}

/* Hover effects mejorados */
.transition-colors {
  transition: color 0.2s ease-in-out, background-color 0.2s ease-in-out;
}

/* Efecto de resaltado para navegacion activa */
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

/* Estilos responsivos mejorados */
@media (max-width: 640px) {
  .text-xl {
    font-size: 1.125rem;
  }
}

@media (min-width: 640px) {
  .sm\:text-2xl {
    font-size: 1.5rem;
  }
}

@media (min-width: 768px) {
  .md\:hidden {
    display: none;
  }
  
  .md\:flex {
    display: flex;
  }
}

@media (min-width: 1024px) {
  .lg\:hidden {
    display: none;
  }
  
  .lg\:flex {
    display: flex;
  }
}

@media (min-width: 1280px) {
  .xl\:space-x-8 > :not([hidden]) ~ :not([hidden]) {
    margin-left: 2rem;
  }
  
  .xl\:px-3 {
    padding-left: 0.75rem;
    padding-right: 0.75rem;
  }
}

/* Grid para navegacion inferior */
.grid-cols-5 {
  grid-template-columns: repeat(5, minmax(0, 1fr));
}

.gap-1 {
  gap: 0.25rem;
}

/* Espaciado consistente */
.space-x-2 > :not([hidden]) ~ :not([hidden]) {
  margin-left: 0.5rem;
}

.space-x-4 > :not([hidden]) ~ :not([hidden]) {
  margin-left: 1rem;
}

.space-y-1 > :not([hidden]) ~ :not([hidden]) {
  margin-top: 0.25rem;
}

/* Texto que no se corta */
.whitespace-nowrap {
  white-space: nowrap;
}

/* Flex shrink */
.flex-shrink-0 {
  flex-shrink: 0;
}

/* Leading tight para texto compacto */
.leading-tight {
  line-height: 1.25;
}

/* Animacion suave para el menu hamburguesa */
.md\:hidden svg {
  transition: transform 0.3s ease-in-out;
}

/* Estilos para elementos flexibles */
.flex-1 {
  flex: 1 1 0%;
}

.justify-center {
  justify-content: center;
}

.items-center {
  align-items: center;
}

.max-w-4xl {
  max-width: 56rem;
}

/* Responsive utilities adicionales */
@media (max-width: 767px) {
  .mobile-only {
    display: block;
  }
}

@media (min-width: 768px) and (max-width: 1023px) {
  .tablet-only {
    display: block;
  }
}

@media (min-width: 1024px) {
  .desktop-only {
    display: block;
  }
}

/* Ajustes para pantallas muy pequeñas */
@media (max-width: 375px) {
  .text-xs {
    font-size: 0.625rem;
  }
  
  .h-6 {
    height: 1.25rem;
  }
  
  .w-6 {
    width: 1.25rem;
  }
}
</style>