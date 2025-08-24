<template>
  <div class="space-y-6">
    <!-- Header con saludo -->
    <div class="bg-gradient-to-r from-blue-600 to-blue-700 rounded-xl p-6 text-white">
      <div class="flex items-center justify-between">
        <div>
          <p class="text-blue-100 mt-1">Sistema de gestión de inventario médico</p>
        </div>
        <div class="text-right">
          <p class="text-sm text-blue-100">Fecha actual</p>
          <p class="text-xl font-semibold">{{ currentDate }}</p>
        </div>
      </div>
    </div>

    <!-- Barra de búsqueda principal -->
    <div class="card">
      <div class="flex items-center gap-4">
        <div class="flex-1">
          <label for="search" class="block text-sm font-medium text-gray-700 mb-2">
            Buscar insumo médico
          </label>
          <div class="relative">
            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <svg class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
            </div>
            <input
              id="search"
              type="text"
              placeholder="Ingrese nombre, código o lote del insumo..."
              class="form-input pl-10 w-full"
              v-model="searchQuery"
              @keyup.enter="handleSearch"
            />
          </div>
        </div>
        <div class="flex gap-3 pt-6">
          <button @click="handleSearch" class="btn-primary">
            <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
            Buscar
          </button>
          <button class="btn-secondary">
            <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 4a1 1 0 011-1h16a1 1 0 011 1v2.586a1 1 0 01-.293.707l-6.414 6.414a1 1 0 00-.293.707V17l-4 4v-6.586a1 1 0 00-.293-.707L3.293 7.207A1 1 0 013 6.5V4z" />
            </svg>
            Filtros
          </button>
        </div>
      </div>
    </div>

    <!-- Funcionalidades principales -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <!-- Gestión de inventario -->
      <div class="card hover:shadow-lg transition-all duration-200 cursor-pointer" @click="navigateTo('/inventory')">
        <div class="flex items-center gap-4">
          <div class="bg-blue-100 p-3 rounded-lg">
            <svg class="h-8 w-8 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
            </svg>
          </div>
          <div>
            <h3 class="text-lg font-semibold text-gray-900">Gestión de Inventario</h3>
            <p class="text-gray-600 text-sm">Administrar insumos médicos</p>
          </div>
        </div>
      </div>

      <!-- Agregar nuevo insumo -->
      <div class="card hover:shadow-lg transition-all duration-200 cursor-pointer" @click="navigateTo('/add-supply')">
        <div class="flex items-center gap-4">
          <div class="bg-green-100 p-3 rounded-lg">
            <svg class="h-8 w-8 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
            </svg>
          </div>
          <div>
            <h3 class="text-lg font-semibold text-gray-900">Agregar Insumo</h3>
            <p class="text-gray-600 text-sm">Registrar nuevo producto médico</p>
          </div>
        </div>
      </div>

      <!-- Escanear QR -->
      <div class="card hover:shadow-lg transition-all duration-200 cursor-pointer" @click="navigateTo('/qr')">
        <div class="flex items-center gap-4">
          <div class="bg-yellow-100 p-3 rounded-lg">
            <svg class="h-8 w-8 text-yellow-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h.01M5 8h2a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v2a2 2 0 002 2zm0 0h2a2 2 0 012 2v2a2 2 0 01-2 2H5a2 2 0 01-2-2v-2a2 2 0 012-2z" />
            </svg>
          </div>
          <div>
            <h3 class="text-lg font-semibold text-gray-900">Escáner QR</h3>
            <p class="text-gray-600 text-sm">Trazabilidad y verificación</p>
          </div>
        </div>
      </div>

      <!-- Reportes de inventario -->
      <div class="card hover:shadow-lg transition-all duration-200 cursor-pointer opacity-75" @click="navigateTo('/inventory/reports')">
        <div class="flex items-center gap-4">
          <div class="bg-purple-100 p-3 rounded-lg">
            <svg class="h-8 w-8 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
            </svg>
          </div>
          <div>
            <h3 class="text-lg font-semibold text-gray-900">Reportes</h3>
            <p class="text-gray-600 text-sm">Análisis y estadísticas <span class="text-yellow-600">(Próximamente)</span></p>
          </div>
        </div>
      </div>

      <!-- Movimientos de stock -->
      <div class="card hover:shadow-lg transition-all duration-200 cursor-pointer opacity-75" @click="navigateTo('/inventory/movements')">
        <div class="flex items-center gap-4">
          <div class="bg-indigo-100 p-3 rounded-lg">
            <svg class="h-8 w-8 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
          </div>
          <div>
            <h3 class="text-lg font-semibold text-gray-900">Movimientos</h3>
            <p class="text-gray-600 text-sm">Entradas y salidas de stock <span class="text-yellow-600">(Próximamente)</span></p>
          </div>
        </div>
      </div>

      <!-- Configuración -->
      <div class="card hover:shadow-lg transition-all duration-200 cursor-pointer opacity-75" @click="navigateTo('/inventory/settings')">
        <div class="flex items-center gap-4">
          <div class="bg-gray-100 p-3 rounded-lg">
            <svg class="h-8 w-8 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
          </div>
          <div>
            <h3 class="text-lg font-semibold text-gray-900">Configuración</h3>
            <p class="text-gray-600 text-sm">Ajustes del sistema <span class="text-yellow-600">(Próximamente)</span></p>
          </div>
        </div>
      </div>
    </div>

    
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const searchQuery = ref('')

const currentDate = computed(() => {
  return new Date().toLocaleDateString('es-ES', {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
})

const handleSearch = () => {
  if (searchQuery.value.trim()) {
    router.push({
      name: 'Inventory',
      query: { search: searchQuery.value }
    })
  } else {
    router.push('/inventory')
  }
}

const navigateTo = (path) => {
  router.push(path)
}
</script>

<style scoped>
/* Efectos hover para las tarjetas */
.card:hover {
  transform: translateY(-2px);
  transition: all 0.2s ease-in-out;
}

/* Gradientes personalizados */
.bg-gradient-to-r {
  background-image: linear-gradient(to right, var(--tw-gradient-stops));
}

/* Transiciones suaves */
.transition-all {
  transition: all 0.2s ease-in-out;
}
</style>