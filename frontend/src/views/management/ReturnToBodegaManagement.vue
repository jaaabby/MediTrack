<template>
  <div class="space-y-6">
    <!-- Header -->
    <div>
      <h1 class="text-xl sm:text-2xl font-bold text-gray-900">Gestión de Retornos a Bodega</h1>
      <p class="text-sm sm:text-base text-gray-600 mt-1">
        Monitoreo y gestión de insumos que deben regresar a bodega después de 8 horas laborales sin consumir
      </p>
    </div>

    <!-- Stats Cards -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <div class="bg-white rounded-lg shadow-sm border p-6">
        <div class="flex items-center">
          <div class="p-2 bg-yellow-100 rounded-lg">
            <svg class="h-8 w-8 text-yellow-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600">Pendientes de Retorno</p>
            <p class="text-2xl font-bold text-gray-900">{{ suppliesForReturn.length }}</p>
          </div>
        </div>
      </div>

      <div class="bg-white rounded-lg shadow-sm border p-6">
        <div class="flex items-center">
          <div class="p-2 bg-red-100 rounded-lg">
            <svg class="h-8 w-8 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v3.75m9-.75a9 9 0 11-18 0 9 9 0 0118 0zm-9 3.75h.008v.008H12v-.008z" />
            </svg>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600">Críticos (>8 horas laborales)</p>
            <p class="text-2xl font-bold text-gray-900">{{ criticalSupplies.length }}</p>
          </div>
        </div>
      </div>

      <div class="bg-white rounded-lg shadow-sm border p-6">
        <div class="flex items-center">
          <div class="p-2 bg-green-100 rounded-lg">
            <svg class="h-8 w-8 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600">Último Proceso</p>
            <p class="text-sm font-bold text-gray-900">{{ lastProcessDate || 'Sin procesar' }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Action Buttons -->
    <div class="card">
      <div class="flex flex-col sm:flex-row sm:items-center gap-4">
        <button 
          @click="refreshData"
          :disabled="loading"
          class="btn-primary flex items-center justify-center px-3 py-2 text-sm w-full sm:w-auto"
        >
          <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
          {{ loading ? 'Actualizando...' : 'Actualizar Datos' }}
        </button>

        <button 
          @click="processAutomaticReturns"
          :disabled="processingReturns || criticalSupplies.length === 0"
          class="btn-danger flex items-center justify-center px-3 py-2 text-sm w-full sm:w-auto"
        >
          <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 10h10a8 8 0 018 8v2M3 10l6 6m-6-6l6-6" />
          </svg>
          {{ processingReturns ? 'Procesando...' : `Procesar Retornos Automáticos (${criticalSupplies.length})` }}
        </button>

        <button 
          @click="toggleAutoRefresh"
          :class="autoRefreshEnabled ? 'btn-secondary' : 'btn-secondary opacity-50'"
          class="flex items-center justify-center px-3 py-2 text-sm w-full sm:w-auto"
          :title="autoRefreshEnabled ? 'Actualización automática activa (cada 30s)' : 'Actualización automática desactivada'"
        >
          <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path v-if="autoRefreshEnabled" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          {{ autoRefreshEnabled ? 'Actualización: ON' : 'Actualización: OFF' }}
        </button>
      </div>
    </div>

    <!-- Supplies Table -->
    <div class="card">
      <div class="card-header">
        <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
          <div>
            <h2 class="card-title">Insumos Pendientes de Retorno</h2>
            <p class="text-sm text-gray-600">Total: {{ filteredSupplies.length }} insumos</p>
          </div>
        </div>
      </div>

      <!-- Indicador de carga -->
      <div v-if="loading" class="flex justify-center items-center py-8">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
        <span class="ml-2 text-gray-600">Cargando datos...</span>
      </div>

      <!-- Mensaje de error -->
      <div v-else-if="error" class="bg-red-50 border border-red-200 rounded-md p-4 mx-4 mb-4">
        <div class="flex">
          <div class="flex-shrink-0">
            <svg class="h-5 w-5 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <div class="ml-3">
            <h3 class="text-sm font-medium text-red-800">Error al cargar datos</h3>
            <div class="mt-2 text-sm text-red-700">{{ error }}</div>
            <div class="mt-4">
              <button @click="refreshData" class="btn-secondary text-sm">
                Reintentar
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Tabla de datos -->
      <div v-else-if="suppliesForReturn.length === 0" class="p-8 text-center">
        <svg class="h-12 w-12 text-gray-400 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <p class="text-gray-600">No hay insumos pendientes de retorno</p>
        <p class="text-sm text-gray-500 mt-1">Todos los insumos están siendo utilizados correctamente</p>
      </div>

      <div v-else class="table-container">
        <!-- Indicador de scroll horizontal para móviles -->
        <div class="md:hidden bg-blue-50 border-b border-blue-200 px-3 py-2 text-center sticky left-0 z-10">
          <div class="flex items-center justify-center text-blue-700 text-xs">
            <svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 9l4-4 4 4m0 6l-4 4-4-4" />
            </svg>
            Desliza horizontalmente para ver todas las columnas
          </div>
        </div>

        <table class="min-w-full divide-y divide-gray-200" style="min-width: 900px;">
          <thead class="table-header">
            <tr>
              <th class="table-header-cell">
                Insumo
              </th>
              <th class="table-header-cell">
                Estado
              </th>
              <th class="table-header-cell">
                Tiempo Sin Consumir
              </th>
              <th class="table-header-cell">
                Recepcionado
              </th>
              <th class="table-header-cell">
                Bodega Destino
              </th>
              <th class="table-header-cell">Acciones</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-for="supply in paginatedSupplies" :key="supply.id" 
                class="hover:bg-gray-50 cursor-pointer"
                :class="{ 'bg-red-50': supply.shouldReturn }">
              <td class="px-3 sm:px-6 py-3 sm:py-4 whitespace-nowrap">
                <div class="flex items-center">
                  <div class="flex-shrink-0 h-8 w-8 sm:h-10 sm:w-10">
                    <div class="h-8 w-8 sm:h-10 sm:w-10 bg-blue-100 rounded-full flex items-center justify-center">
                      <svg class="h-4 w-4 sm:h-6 sm:w-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19.428 15.428a2 2 0 00-1.022-.547l-2.387-.477a6 6 0 00-3.86.517l-.318.158a6 6 0 01-3.86.517L6.05 15.21a2 2 0 00-1.806.547M8 4h8l-1 1v5.172a2 2 0 00.586 1.414l5 5c1.26 1.26.367 3.414-1.415 3.414H4.828c-1.782 0-2.674-2.154-1.414-3.414l5-5A2 2 0 009 9.172V5L8 4z" />
                      </svg>
                    </div>
                  </div>
                  <div class="ml-3 sm:ml-4">
                    <div class="text-xs sm:text-sm font-medium text-gray-900">{{ supply.name }}</div>
                    <div class="text-xs sm:text-sm text-gray-500">{{ supply.qrCode }}</div>
                  </div>
                </div>
              </td>
              <td class="px-3 sm:px-6 py-3 sm:py-4 whitespace-nowrap">
                <span :class="getStatusClass(supply.status)" class="inline-flex px-2 py-1 text-xs font-semibold rounded-full">
                  {{ getStatusText(supply.status) }}
                </span>
              </td>
              <td class="px-3 sm:px-6 py-3 sm:py-4 whitespace-nowrap">
                <div class="flex items-center">
                  <span :class="{ 'text-red-600 font-bold': supply.shouldReturn, 'text-gray-900': !supply.shouldReturn }" class="text-xs sm:text-sm">
                    {{ formatBusinessHours(supply.businessHoursElapsed) }}
                  </span>
                  <svg v-if="supply.shouldReturn" class="h-3 w-3 sm:h-4 sm:w-4 text-red-600 ml-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-2.694-.833-3.464 0L3.35 16.5c-.77.833.192 2.5 1.732 2.5z" />
                  </svg>
                </div>
              </td>
              <td class="px-3 sm:px-6 py-3 sm:py-4 whitespace-nowrap text-xs sm:text-sm text-gray-900">
                {{ formatDate(supply.receivedAt) }}
              </td>
              <td class="px-3 sm:px-6 py-3 sm:py-4 whitespace-nowrap text-xs sm:text-sm text-gray-900">
                {{ supply.storeName }}
              </td>
              <td class="px-3 sm:px-6 py-3 sm:py-4 whitespace-nowrap" @click.stop>
                <div class="flex flex-col sm:flex-row gap-1 sm:gap-2">
                <button 
                    @click.stop="returnIndividualSupply(supply)"
                  :disabled="returningSupplies[supply.qrCode]"
                    class="btn-danger-sm text-xs w-full sm:w-auto min-w-[100px] sm:min-w-[90px]"
                >
                  {{ returningSupplies[supply.qrCode] ? 'Regresando...' : 'Regresar' }}
                </button>
                <button 
                    @click.stop="viewSupplyDetails(supply.qrCode)"
                    class="btn-secondary-sm text-xs w-full sm:w-auto min-w-[100px] sm:min-w-[90px]"
                >
                  Ver Detalles
                </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Paginación -->
      <div v-if="!loading && suppliesForReturn.length > 0" class="flex flex-col sm:flex-row items-center justify-between mt-6 gap-4 px-4 pb-4">
        <div class="text-sm text-gray-700 text-center sm:text-left">
          Mostrando {{ startIndex + 1 }} a {{ endIndex }} de {{ filteredSupplies.length }} resultados
    </div>
        <div class="flex items-center gap-2">
          <button class="btn-secondary px-3 py-2 text-sm min-w-[70px]" :disabled="currentPage === 1"
            @click="currentPage--">
            <span class="hidden sm:inline">Anterior</span>
            <span class="sm:hidden">Ant.</span>
          </button>
          <span class="px-3 py-2 text-sm text-gray-700 bg-gray-100 rounded-md min-w-[90px] text-center">
            Página {{ currentPage }} de {{ totalPages }}
          </span>
          <button class="btn-secondary px-3 py-2 text-sm min-w-[70px]" :disabled="currentPage === totalPages"
            @click="currentPage++">
            <span class="hidden sm:inline">Siguiente</span>
            <span class="sm:hidden">Sig.</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { format } from 'date-fns'
import { es } from 'date-fns/locale'
import returnToBodegaService from '@/services/management/returnToBodegaService'
import Swal from 'sweetalert2'

const router = useRouter()

// Estado del componente
const loading = ref(false)
const error = ref(null)
const suppliesForReturn = ref([])
const processingReturns = ref(false)
const returningSupplies = ref({})
const lastProcessDate = ref(null)
const currentPage = ref(1)
const itemsPerPage = 10
const autoRefreshInterval = ref(null)
const autoRefreshEnabled = ref(true)
const refreshIntervalSeconds = 30 // Actualizar cada 30 segundos

// Computed properties
const criticalSupplies = computed(() => {
  return suppliesForReturn.value.filter(supply => supply.shouldReturn || (supply.businessHoursElapsed || 0) >= 8)
})

const filteredSupplies = computed(() => {
  return [...suppliesForReturn.value]
})

const totalPages = computed(() => Math.ceil(filteredSupplies.value.length / itemsPerPage))

const startIndex = computed(() => (currentPage.value - 1) * itemsPerPage)
const endIndex = computed(() => Math.min(startIndex.value + itemsPerPage, filteredSupplies.value.length))

const paginatedSupplies = computed(() => {
  return filteredSupplies.value.slice(startIndex.value, endIndex.value)
})

// Funciones principales
const refreshData = async () => {
  loading.value = true
  error.value = null
  
  try {
    console.log('🔍 Iniciando carga de datos...')
    const supplies = await returnToBodegaService.getSuppliesForReturn()
    console.log('🔍 Datos recibidos del servicio:', supplies)
    console.log('🔍 Tipo de datos:', typeof supplies)
    console.log('🔍 Es array:', Array.isArray(supplies))
    
    // Verificar si supplies es un array, si no convertirlo
    let suppliesArray = []
    if (Array.isArray(supplies)) {
      suppliesArray = supplies
    } else if (supplies && typeof supplies === 'object') {
      // Si supplies es un objeto, intentar extraer el array
      if (supplies.data && Array.isArray(supplies.data)) {
        suppliesArray = supplies.data
      } else if (supplies.supplies && Array.isArray(supplies.supplies)) {
        suppliesArray = supplies.supplies
      } else {
        console.warn('⚠️ Estructura de respuesta no esperada:', supplies)
        suppliesArray = []
      }
    } else {
      console.warn('⚠️ Respuesta no válida del servidor:', supplies)
      suppliesArray = []
    }
    
    console.log('🔍 Array final a procesar:', suppliesArray)
    suppliesForReturn.value = suppliesArray.map(supply => returnToBodegaService.formatSupplyForUI(supply))
    lastProcessDate.value = new Date().toLocaleString()
    
    console.log('✅ Datos cargados exitosamente:', suppliesForReturn.value.length, 'insumos')
    currentPage.value = 1 // Resetear a la primera página al cargar nuevos datos
  } catch (err) {
    console.error('❌ Error refreshing data:', err)
    error.value = err.message || 'Error al cargar los datos'
    // Asegurar que suppliesForReturn sea un array vacío en caso de error
    suppliesForReturn.value = []
  } finally {
    loading.value = false
  }
}

const processAutomaticReturns = async () => {
  if (criticalSupplies.value.length === 0) return
  
  const result = await Swal.fire({
    title: `¿Está seguro de que desea procesar automáticamente ${criticalSupplies.value.length} retornos?`,
    text: 'Esta acción regresará todos los insumos críticos (15+ días) a bodega.',
    icon: 'warning',
    showCancelButton: true,
    confirmButtonText: 'Sí, procesar',
    cancelButtonText: 'Cancelar',
  })
  if (!result.isConfirmed) return
  
  processingReturns.value = true
  error.value = null
  
  try {
    const result = await returnToBodegaService.processAutomaticReturns()
    
    Swal.fire({
      icon: 'success',
      title: 'Proceso completado',
      text: result.message
    })
    
    // Recargar datos
    await refreshData()
    
  } catch (err) {
    console.error('Error processing automatic returns:', err)
    const errorMessage = err.response?.data?.error || err.message || 'Error al procesar retornos automáticos'
    error.value = errorMessage
    
    Swal.fire({
      icon: 'error',
      title: 'Error al procesar retornos',
      text: errorMessage,
      confirmButtonText: 'Aceptar'
    })
  } finally {
    processingReturns.value = false
  }
}

const returnIndividualSupply = async (supply) => {
  if (returningSupplies.value[supply.qrCode]) return
  
  const result = await Swal.fire({
    title: `¿Regresar ${supply.name} a bodega?`,
    html: `QR: ${supply.qrCode}<br>Días sin consumir: ${supply.daysElapsed}`,
    icon: 'question',
    showCancelButton: true,
    confirmButtonText: 'Sí, regresar',
    cancelButtonText: 'Cancelar',
  })
  if (!result.isConfirmed) return
  
  returningSupplies.value[supply.qrCode] = true
  
  try {
    await returnToBodegaService.returnSupplyToStore(
      supply.qrCode, 
      `Retorno manual - ${supply.daysElapsed} días sin consumir`
    )
    
    Swal.fire({
      icon: 'success',
      title: 'Éxito',
      text: 'Insumo regresado a bodega exitosamente'
    })
    
    // Recargar datos
    await refreshData()
    
  } catch (err) {
    console.error('Error returning supply:', err)
    error.value = err.message || 'Error al regresar insumo a bodega'
  } finally {
    delete returningSupplies.value[supply.qrCode]
  }
}

const viewSupplyDetails = (qrCode) => {
  router.push({
    name: 'QRDetails',
    params: { qrCode: qrCode }
  })
}

// Funciones auxiliares
const getStatusClass = (status) => {
  switch (status) {
    case 'recepcionado':
      return 'bg-blue-100 text-blue-800'
    case 'disponible':
      return 'bg-green-100 text-green-800'
    case 'en_camino_a_pabellon':
      return 'bg-yellow-100 text-yellow-800'
    case 'consumido':
      return 'bg-gray-100 text-gray-800'
    default:
      return 'bg-gray-100 text-gray-800'
  }
}

const getStatusText = (status) => {
  switch (status) {
    case 'recepcionado':
      return 'Recepcionado'
    case 'disponible':
      return 'Disponible'
    case 'en_camino_a_pabellon':
      return 'En Camino'
    case 'consumido':
      return 'Consumido'
    default:
      return status || 'Desconocido'
  }
}

const formatDate = (dateString) => {
  if (!dateString) return 'No disponible'
  try {
    return format(new Date(dateString), 'dd/MM/yyyy HH:mm', { locale: es })
  } catch (error) {
    return dateString
  }
}

const formatBusinessHours = (hours) => {
  if (!hours && hours !== 0) return '0 horas laborales'
  
  // Si es menos de 1 hora, mostrar en minutos
  if (hours < 1) {
    const minutes = Math.round(hours * 60)
    return `${minutes} min${minutes !== 1 ? 's' : ''} laborales`
  }
  
  // Si es menos de 8 horas, mostrar en horas
  if (hours < 8) {
    const roundedHours = Math.round(hours * 10) / 10
    return `${roundedHours} hora${roundedHours !== 1 ? 's' : ''} laborales`
  }
  
  // Más de 8 horas, mostrar en días aproximados
  const days = Math.floor(hours / 8)
  const remainingHours = hours % 8
  if (remainingHours === 0) {
    return `${days} día${days !== 1 ? 's' : ''} laborales`
  }
  return `${days} día${days !== 1 ? 's' : ''} y ${Math.round(remainingHours)} hora${Math.round(remainingHours) !== 1 ? 's' : ''} laborales`
}

const clearError = () => {
  error.value = null
}

// Funciones para actualización automática
const startAutoRefresh = () => {
  if (autoRefreshInterval.value) {
    clearInterval(autoRefreshInterval.value)
  }
  
  autoRefreshInterval.value = setInterval(() => {
    if (autoRefreshEnabled.value && !loading.value && !processingReturns.value) {
      refreshData()
    }
  }, refreshIntervalSeconds * 1000)
}

const stopAutoRefresh = () => {
  if (autoRefreshInterval.value) {
    clearInterval(autoRefreshInterval.value)
    autoRefreshInterval.value = null
  }
}

const toggleAutoRefresh = () => {
  autoRefreshEnabled.value = !autoRefreshEnabled.value
  if (autoRefreshEnabled.value) {
    startAutoRefresh()
  } else {
    stopAutoRefresh()
  }
}

// Lifecycle
onMounted(() => {
  refreshData()
  startAutoRefresh()
})

onUnmounted(() => {
  stopAutoRefresh()
})
</script>

<style scoped>
/* Estilos similares a Inventory.vue */

.card {
  @apply bg-white rounded-lg shadow-sm border p-4 sm:p-6;
}

.card-header {
  @apply pb-4 mb-4 border-b border-gray-200;
}

.card-title {
  @apply text-lg sm:text-xl font-semibold text-gray-900;
}

.table-container {
  overflow-x: auto;
  border: 1px solid #e5e7eb;
  border-radius: 0.5rem;
  box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1);
  -webkit-overflow-scrolling: touch;
  position: relative;
}

@media (min-width: 1024px) {
  .table-container {
    overflow-x: visible;
  }
  
  .table-container table {
    width: 100%;
    min-width: 100%;
  }
}

.table-container::-webkit-scrollbar {
  height: 6px;
}

@media (min-width: 768px) {
  .table-container::-webkit-scrollbar {
    height: 8px;
  }
}

.table-container::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 4px;
}

.table-container::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 4px;
}

.table-container::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}

.table-header {
  background-color: #f9fafb;
}

.table-header th {
  padding: 0.75rem 1.5rem;
  text-align: left;
  font-size: 0.75rem;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: #6b7280;
  white-space: nowrap;
  min-width: 120px;
}

.table-header-cell {
  @apply px-3 sm:px-6 py-2 sm:py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider whitespace-nowrap min-w-[100px] sm:min-w-[120px];
}

.btn-danger {
  @apply inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-red-600 hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 disabled:opacity-50 disabled:cursor-not-allowed;
}

.btn-danger-sm {
  @apply inline-flex items-center justify-center px-2 sm:px-3 py-1.5 sm:py-2 border border-transparent text-xs font-medium rounded text-white bg-red-600 hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 disabled:opacity-50 disabled:cursor-not-allowed;
}

.btn-secondary-sm {
  @apply inline-flex items-center justify-center px-2 sm:px-3 py-1.5 sm:py-2 border border-gray-300 text-xs font-medium rounded text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500;
}

/* Responsive para móviles */
@media (max-width: 640px) {
  .card {
    padding: 0.875rem;
    margin-left: -0.5rem;
    margin-right: -0.5rem;
    border-radius: 0.5rem;
  }

  .card-header {
    padding: 0.875rem;
  }

  .table-container {
    margin: 0 -0.875rem;
    border-radius: 0;
    border-left: none;
    border-right: none;
  }

  .table-header th,
  .table-body td {
    padding: 0.5rem 0.75rem;
  }

  button,
  .btn-primary,
  .btn-secondary {
    min-height: 44px;
  }
}

/* Responsive para tablets */
@media (max-width: 768px) {
  .card {
    padding: 1.5rem;
  }
}
</style>