<template>
  <div class="space-y-8">
    <!-- Header -->
    <div class="bg-white rounded-lg shadow-sm border p-6">
      <h2 class="text-2xl font-semibold text-gray-900">Dashboard de Inventario</h2>
      <p class="text-gray-600 mt-2">Inventario organizado por tipo de cirugía</p>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex justify-center items-center py-12">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-brand-blue-dark"></div>
      <span class="ml-3 text-gray-600">Cargando inventario...</span>
    </div>

    <!-- Error -->
    <div v-else-if="error" class="card">
      <div class="bg-red-50 border border-red-200 rounded-md p-4">
        <div class="flex">
          <svg class="h-5 w-5 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <div class="ml-3">
            <h3 class="text-sm font-medium text-red-800">Error al cargar inventario</h3>
            <div class="mt-2 text-sm text-red-700">{{ error }}</div>
            <button @click="loadDashboard" class="btn-secondary mt-4 text-sm">Reintentar</button>
          </div>
        </div>
      </div>
    </div>

    <!-- Dashboard Content -->
    <div v-else class="space-y-8">
      <!-- Inventario por Tipo de Cirugía -->
      <div class="card">
        <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 mb-6">
          <div>
            <h3 class="text-xl font-semibold text-gray-900">Inventario por Tipo de Cirugía</h3>
            <p class="text-gray-600 text-sm mt-2">Stock disponible organizado por procedimiento</p>
          </div>
          <button @click="reloadSurgeryInventory" class="btn-secondary">
            <svg class="h-5 w-5 mr-2 inline" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
            Actualizar
          </button>
        </div>

        <div v-if="surgeryInventoryLoading" class="flex justify-center items-center py-8">
          <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-brand-blue-dark"></div>
          <span class="ml-3 text-gray-600">Cargando...</span>
        </div>

        <!-- Error al actualizar (solo aparece al usar el botón Actualizar, no en carga inicial) -->
        <div v-else-if="surgeryInventoryError" class="bg-red-50 border border-red-200 rounded-md p-4">
          <div class="flex">
            <svg class="h-5 w-5 text-red-400 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <div class="ml-3">
              <h3 class="text-sm font-medium text-red-800">Error al actualizar el inventario</h3>
              <div class="mt-1 text-sm text-red-700">{{ surgeryInventoryError }}</div>
              <button @click="reloadSurgeryInventory" class="btn-secondary mt-3 text-sm">Reintentar</button>
            </div>
          </div>
        </div>

        <div v-else-if="surgeryInventory.length === 0" class="text-center py-8">
          <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4" />
          </svg>
          <p class="mt-2 text-sm text-gray-500">No hay inventario organizado por tipo de cirugía</p>
        </div>

        <div v-else class="overflow-x-auto">
          <div class="max-h-[600px] overflow-y-auto">
          <table class="min-w-full divide-y divide-gray-200">
            <thead class="bg-gray-50 sticky top-0 z-10">
              <tr>
                <th @click="sortSurgeryBy('surgery_name')"
                  class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100 transition-colors select-none">
                  <div class="flex items-center space-x-1">
                    <span>Tipo de Cirugía</span>
                    <span class="flex flex-col -space-y-1">
                      <svg class="h-3 w-3 transition-colors" 
                        :class="surgerySortKey === 'surgery_name' && surgerySortOrder === 'asc' ? 'text-brand-blue-dark' : 'text-gray-300'" 
                        fill="currentColor" viewBox="0 0 20 20">
                        <path d="M5.293 9.707a1 1 0 010-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 01-1.414 1.414L11 7.414V15a1 1 0 11-2 0V7.414L6.707 9.707a1 1 0 01-1.414 0z"/>
                      </svg>
                      <svg class="h-3 w-3 transition-colors" 
                        :class="surgerySortKey === 'surgery_name' && surgerySortOrder === 'desc' ? 'text-blue-600' : 'text-gray-300'" 
                        fill="currentColor" viewBox="0 0 20 20">
                        <path d="M14.707 10.293a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 111.414-1.414L9 12.586V5a1 1 0 012 0v7.586l2.293-2.293a1 1 0 011.414 0z"/>
                      </svg>
                    </span>
                  </div>
                </th>
                <th @click="sortSurgeryBy('total_in_store')"
                  class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100 transition-colors select-none">
                  <div class="flex items-center space-x-1">
                    <span>Stock en Bodega</span>
                    <span class="flex flex-col -space-y-1">
                      <svg class="h-3 w-3 transition-colors" 
                        :class="surgerySortKey === 'total_in_store' && surgerySortOrder === 'asc' ? 'text-blue-600' : 'text-gray-300'" 
                        fill="currentColor" viewBox="0 0 20 20">
                        <path d="M5.293 9.707a1 1 0 010-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 01-1.414 1.414L11 7.414V15a1 1 0 11-2 0V7.414L6.707 9.707a1 1 0 01-1.414 0z"/>
                      </svg>
                      <svg class="h-3 w-3 transition-colors" 
                        :class="surgerySortKey === 'total_in_store' && surgerySortOrder === 'desc' ? 'text-blue-600' : 'text-gray-300'" 
                        fill="currentColor" viewBox="0 0 20 20">
                        <path d="M14.707 10.293a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 111.414-1.414L9 12.586V5a1 1 0 012 0v7.586l2.293-2.293a1 1 0 011.414 0z"/>
                      </svg>
                    </span>
                  </div>
                </th>
                <th @click="sortSurgeryBy('total_transferred')"
                  class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100 transition-colors select-none">
                  <div class="flex items-center space-x-1">
                    <span>Total Transferido</span>
                    <span class="flex flex-col -space-y-1">
                      <svg class="h-3 w-3 transition-colors" 
                        :class="surgerySortKey === 'total_transferred' && surgerySortOrder === 'asc' ? 'text-blue-600' : 'text-gray-300'" 
                        fill="currentColor" viewBox="0 0 20 20">
                        <path d="M5.293 9.707a1 1 0 010-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 01-1.414 1.414L11 7.414V15a1 1 0 11-2 0V7.414L6.707 9.707a1 1 0 01-1.414 0z"/>
                      </svg>
                      <svg class="h-3 w-3 transition-colors" 
                        :class="surgerySortKey === 'total_transferred' && surgerySortOrder === 'desc' ? 'text-blue-600' : 'text-gray-300'" 
                        fill="currentColor" viewBox="0 0 20 20">
                        <path d="M14.707 10.293a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 111.414-1.414L9 12.586V5a1 1 0 012 0v7.586l2.293-2.293a1 1 0 011.414 0z"/>
                      </svg>
                    </span>
                  </div>
                </th>
                <th @click="sortSurgeryBy('batch_count')"
                  class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100 transition-colors select-none">
                  <div class="flex items-center space-x-1">
                    <span>N° de Lotes</span>
                    <span class="flex flex-col -space-y-1">
                      <svg class="h-3 w-3 transition-colors" 
                        :class="surgerySortKey === 'batch_count' && surgerySortOrder === 'asc' ? 'text-blue-600' : 'text-gray-300'" 
                        fill="currentColor" viewBox="0 0 20 20">
                        <path d="M5.293 9.707a1 1 0 010-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 01-1.414 1.414L11 7.414V15a1 1 0 11-2 0V7.414L6.707 9.707a1 1 0 01-1.414 0z"/>
                      </svg>
                      <svg class="h-3 w-3 transition-colors" 
                        :class="surgerySortKey === 'batch_count' && surgerySortOrder === 'desc' ? 'text-blue-600' : 'text-gray-300'" 
                        fill="currentColor" viewBox="0 0 20 20">
                        <path d="M14.707 10.293a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 111.414-1.414L9 12.586V5a1 1 0 012 0v7.586l2.293-2.293a1 1 0 011.414 0z"/>
                      </svg>
                    </span>
                  </div>
                </th>
                <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Acciones</th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
              <tr v-for="item in sortedSurgeryInventory" :key="item.surgery_id" class="hover:bg-gray-50">
                <td class="px-6 py-4 whitespace-nowrap">
                  <div class="flex items-center">
                    <div class="flex-shrink-0 h-10 w-10 rounded-full bg-gradient-to-br from-brand-blue-light to-brand-blue-dark flex items-center justify-center">
                      <svg class="h-6 w-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
                      </svg>
                    </div>
                    <div class="ml-4">
                      <div class="text-sm font-medium text-gray-900">{{ item.surgery_name || 'Sin tipo de cirugía' }}</div>
                      <div class="text-sm text-gray-500">ID: {{ item.surgery_id || 'N/A' }}</div>
                    </div>
                  </div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <span class="text-sm font-semibold text-gray-900">{{ item.total_in_store || 0 }}</span>
                  <span class="text-xs text-gray-500 ml-1">unidades</span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <span class="text-sm text-gray-900">{{ item.total_transferred || 0 }}</span>
                  <span class="text-xs text-gray-500 ml-1">unidades</span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <span class="px-2 py-1 text-xs font-semibold rounded-full bg-brand-blue-light bg-opacity-30 text-brand-blue-dark">
                    {{ item.batch_count || 0 }} lotes
                  </span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                  <router-link :to="`/inventory/store?surgery_id=${item.surgery_id}`" class="text-brand-blue-dark hover:text-brand-blue-medium">
                    Ver detalles →
                  </router-link>
                </td>
              </tr>
            </tbody>
          </table>
          </div>
        </div>
      </div>

      <!-- Accesos Rápidos -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4 md:gap-6">
        <!-- Ver Inventario de Bodegas -->
        <router-link to="/inventory/store" class="card hover:shadow-lg transition-shadow cursor-pointer group">
          <div class="flex items-center">
            <div class="flex-shrink-0 h-12 w-12 rounded-lg bg-brand-blue-light bg-opacity-30 flex items-center justify-center group-hover:bg-brand-blue-light group-hover:bg-opacity-50 transition-colors">
              <svg class="h-6 w-6 text-brand-blue-dark" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
              </svg>
            </div>
            <div class="ml-4 flex-1">
              <h3 class="text-lg font-semibold text-gray-900 group-hover:text-brand-blue-dark transition-colors">Inventario de Bodegas</h3>
              <p class="text-sm text-gray-600">Ver stock detallado por bodega</p>
            </div>
            <svg class="h-5 w-5 text-gray-400 group-hover:text-brand-blue-dark transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
          </div>
        </router-link>

        <!-- Ver Inventario de Pabellones -->
        <router-link to="/inventory/pavilion" class="card hover:shadow-lg transition-shadow cursor-pointer group">
          <div class="flex items-center">
            <div class="flex-shrink-0 h-12 w-12 rounded-lg bg-green-100 flex items-center justify-center group-hover:bg-green-200 transition-colors">
              <svg class="h-6 w-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
              </svg>
            </div>
            <div class="ml-4 flex-1">
              <h3 class="text-lg font-semibold text-gray-900 group-hover:text-green-600 transition-colors">Inventario de Pabellones</h3>
              <p class="text-sm text-gray-600">Ver stock en cada pabellón</p>
            </div>
            <svg class="h-5 w-5 text-gray-400 group-hover:text-green-600 transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
          </div>
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import inventoryService from '@/services/inventory/inventoryService'

const loading = ref(false)
const error = ref(null)
const surgeryInventory = ref([])
const surgeryInventoryLoading = ref(false)
const surgeryInventoryError = ref(null)

// Estado de ordenamiento para tabla de cirugías
const surgerySortKey = ref('surgery_name')
const surgerySortOrder = ref('asc')

const loadDashboard = async () => {
  loading.value = true
  error.value = null
  
  try {
    // Cargar inventario por cirugía
    await loadSurgeryInventory()
  } catch (err) {
    error.value = err.message || 'Error al cargar el dashboard'
    console.error('Error loading dashboard:', err)
  } finally {
    loading.value = false
  }
}

const loadSurgeryInventory = async () => {
  surgeryInventoryLoading.value = true
  try {
    const data = await inventoryService.getInventoryBySurgeryType()
    // Asegurarse de que siempre sea un array
    surgeryInventory.value = Array.isArray(data) ? data : []
  } catch (err) {
    console.error('Error loading surgery inventory:', err)
    surgeryInventory.value = []
    throw err // Propagar para que el llamador decida cómo manejarlo
  } finally {
    surgeryInventoryLoading.value = false
  }
}

// Wrapper para el botón "Actualizar": muestra error inline sin ocultar el dashboard
const reloadSurgeryInventory = async () => {
  surgeryInventoryError.value = null
  try {
    await loadSurgeryInventory()
  } catch (err) {
    surgeryInventoryError.value = err.message || 'Error al actualizar el inventario'
  }
}

// Computed para ordenar inventario de cirugías
const sortedSurgeryInventory = computed(() => {
  if (!surgeryInventory.value || surgeryInventory.value.length === 0) return []
  
  const sorted = [...surgeryInventory.value].sort((a, b) => {
    let aVal = a[surgerySortKey.value]
    let bVal = b[surgerySortKey.value]
    
    // Manejo de valores nulos/undefined
    if (aVal === null || aVal === undefined) aVal = ''
    if (bVal === null || bVal === undefined) bVal = ''
    
    // Manejo de strings (comparación case-insensitive)
    if (typeof aVal === 'string') {
      aVal = aVal.toLowerCase()
      bVal = (bVal || '').toString().toLowerCase()
    }
    
    // Comparación
    if (aVal < bVal) return surgerySortOrder.value === 'asc' ? -1 : 1
    if (aVal > bVal) return surgerySortOrder.value === 'asc' ? 1 : -1
    return 0
  })
  
  return sorted
})

// Función para ordenar por columna
const sortSurgeryBy = (key) => {
  if (surgerySortKey.value === key) {
    surgerySortOrder.value = surgerySortOrder.value === 'asc' ? 'desc' : 'asc'
  } else {
    surgerySortKey.value = key
    surgerySortOrder.value = 'asc'
  }
}

onMounted(async () => {
  await loadDashboard()
  // Prefetch del componente StoreInventoryView para acelerar la navegación
  // al hacer clic en "Ver detalles" de cualquier fila (ID 11)
  import('@/views/inventory/StoreInventoryView.vue')
})
</script>

