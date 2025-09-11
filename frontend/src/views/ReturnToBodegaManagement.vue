<template>
  <div class="max-w-6xl mx-auto p-6">
    <!-- Header -->
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-gray-900">Gestión de Retornos a Bodega</h1>
      <p class="text-gray-600 mt-2">
        Monitoreo y gestión de insumos que deben regresar a bodega después de 15 días sin consumir
      </p>
    </div>

    <!-- Stats Cards -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-6">
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
            <p class="text-sm font-medium text-gray-600">Críticos (>15 días)</p>
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
    <div class="bg-white rounded-lg shadow-sm border p-6 mb-6">
      <div class="flex flex-wrap gap-4">
        <button 
          @click="refreshData"
          :disabled="loading"
          class="btn-primary"
        >
          <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
          {{ loading ? 'Actualizando...' : 'Actualizar Datos' }}
        </button>

        <button 
          @click="processAutomaticReturns"
          :disabled="processingReturns || criticalSupplies.length === 0"
          class="btn-danger"
        >
          <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 10h10a8 8 0 018 8v2M3 10l6 6m-6-6l6-6" />
          </svg>
          {{ processingReturns ? 'Procesando...' : `Procesar Retornos Automáticos (${criticalSupplies.length})` }}
        </button>
      </div>
    </div>

    <!-- Supplies Table -->
    <div class="bg-white rounded-lg shadow-sm border overflow-hidden">
      <div class="px-6 py-4 border-b border-gray-200">
        <h3 class="text-lg font-medium text-gray-900">Insumos Pendientes de Retorno</h3>
        <p class="text-sm text-gray-600 mt-1">
          Insumos recepcionados que llevan tiempo sin consumir
        </p>
      </div>

      <div v-if="loading" class="p-8 text-center">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>
        <p class="text-gray-600 mt-2">Cargando datos...</p>
      </div>

      <div v-else-if="suppliesForReturn.length === 0" class="p-8 text-center">
        <svg class="h-12 w-12 text-gray-400 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <p class="text-gray-600">No hay insumos pendientes de retorno</p>
        <p class="text-sm text-gray-500 mt-1">Todos los insumos están siendo utilizados correctamente</p>
      </div>

      <div v-else class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Insumo
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Estado
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Días Sin Consumir
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Recepcionado
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Bodega Destino
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Acciones
              </th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-for="supply in suppliesForReturn" :key="supply.supply_id" 
                :class="{ 'bg-red-50': supply.days_elapsed >= 15 }">
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="flex items-center">
                  <div class="flex-shrink-0 h-10 w-10">
                    <div class="h-10 w-10 bg-blue-100 rounded-full flex items-center justify-center">
                      <svg class="h-6 w-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19.428 15.428a2 2 0 00-1.022-.547l-2.387-.477a6 6 0 00-3.86.517l-.318.158a6 6 0 01-3.86.517L6.05 15.21a2 2 0 00-1.806.547M8 4h8l-1 1v5.172a2 2 0 00.586 1.414l5 5c1.26 1.26.367 3.414-1.415 3.414H4.828c-1.782 0-2.674-2.154-1.414-3.414l5-5A2 2 0 009 9.172V5L8 4z" />
                      </svg>
                    </div>
                  </div>
                  <div class="ml-4">
                    <div class="text-sm font-medium text-gray-900">{{ supply.supply_name }}</div>
                    <div class="text-sm text-gray-500">{{ supply.qr_code }}</div>
                  </div>
                </div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span :class="getStatusClass(supply.status)" class="inline-flex px-2 py-1 text-xs font-semibold rounded-full">
                  {{ getStatusText(supply.status) }}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="flex items-center">
                  <span :class="{ 'text-red-600 font-bold': supply.days_elapsed >= 15, 'text-gray-900': supply.days_elapsed < 15 }">
                    {{ supply.days_elapsed }} días
                  </span>
                  <svg v-if="supply.days_elapsed >= 15" class="h-4 w-4 text-red-600 ml-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-2.694-.833-3.464 0L3.35 16.5c-.77.833.192 2.5 1.732 2.5z" />
                  </svg>
                </div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                {{ formatDate(supply.received_at) }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                {{ supply.store_name }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
                <button 
                  @click="returnIndividualSupply(supply)"
                  :disabled="returningSupplies[supply.qr_code]"
                  class="btn-danger-sm mr-2"
                >
                  {{ returningSupplies[supply.qr_code] ? 'Regresando...' : 'Regresar' }}
                </button>
                <button 
                  @click="viewSupplyDetails(supply.qr_code)"
                  class="btn-secondary-sm"
                >
                  Ver Detalles
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Error Message -->
    <div v-if="error" class="mt-6 bg-red-50 border border-red-200 rounded-lg p-4">
      <div class="flex items-start space-x-3">
        <svg class="h-5 w-5 text-red-600 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <div class="flex-1">
          <h4 class="text-sm font-medium text-red-800">Error</h4>
          <p class="text-sm text-red-700 mt-1">{{ error }}</p>
          <button @click="clearError" class="text-sm text-red-600 hover:text-red-800 mt-2 underline">
            Limpiar error
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { format } from 'date-fns'
import { es } from 'date-fns/locale'
import returnToBodegaService from '@/services/returnToBodegaService'

const router = useRouter()

// Estado del componente
const loading = ref(false)
const error = ref(null)
const suppliesForReturn = ref([])
const processingReturns = ref(false)
const returningSupplies = ref({})
const lastProcessDate = ref(null)

// Computed properties
const criticalSupplies = computed(() => {
  return suppliesForReturn.value.filter(supply => supply.days_elapsed >= 15)
})

// Funciones principales
const refreshData = async () => {
  loading.value = true
  error.value = null
  
  try {
    const supplies = await returnToBodegaService.getSuppliesForReturn()
    suppliesForReturn.value = supplies.map(supply => returnToBodegaService.formatSupplyForUI(supply))
    lastProcessDate.value = new Date().toLocaleString()
  } catch (err) {
    console.error('Error refreshing data:', err)
    error.value = err.message || 'Error al cargar los datos'
  } finally {
    loading.value = false
  }
}

const processAutomaticReturns = async () => {
  if (criticalSupplies.value.length === 0) return
  
  const confirmed = confirm(
    `¿Está seguro de que desea procesar automáticamente ${criticalSupplies.value.length} retornos?\n\n` +
    'Esta acción regresará todos los insumos críticos (15+ días) a bodega.'
  )
  
  if (!confirmed) return
  
  processingReturns.value = true
  error.value = null
  
  try {
    const result = await returnToBodegaService.processAutomaticReturns()
    
    alert(`✅ Proceso completado: ${result.message}`)
    
    // Recargar datos
    await refreshData()
    
  } catch (err) {
    console.error('Error processing automatic returns:', err)
    error.value = err.message || 'Error al procesar retornos automáticos'
  } finally {
    processingReturns.value = false
  }
}

const returnIndividualSupply = async (supply) => {
  if (returningSupplies.value[supply.qrCode]) return
  
  const confirmed = confirm(
    `¿Regresar ${supply.name} a bodega?\n\n` +
    `QR: ${supply.qrCode}\n` +
    `Días sin consumir: ${supply.daysElapsed}`
  )
  
  if (!confirmed) return
  
  returningSupplies.value[supply.qrCode] = true
  
  try {
    await returnToBodegaService.returnSupplyToStore(
      supply.qrCode, 
      `Retorno manual - ${supply.daysElapsed} días sin consumir`
    )
    
    alert('✅ Insumo regresado a bodega exitosamente')
    
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
    params: { qrcode: qrCode }
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

const clearError = () => {
  error.value = null
}

// Lifecycle
onMounted(() => {
  refreshData()
})
</script>

<style scoped>
.btn-primary {
  @apply inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed;
}

.btn-danger {
  @apply inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-red-600 hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 disabled:opacity-50 disabled:cursor-not-allowed;
}

.btn-danger-sm {
  @apply inline-flex items-center px-3 py-1 border border-transparent text-xs font-medium rounded text-white bg-red-600 hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 disabled:opacity-50 disabled:cursor-not-allowed;
}

.btn-secondary-sm {
  @apply inline-flex items-center px-3 py-1 border border-gray-300 text-xs font-medium rounded text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500;
}
</style>