<template>
  <div class="space-y-8">
    <!-- Header -->
    <div class="bg-white rounded-lg shadow-sm border p-6">
      <h2 class="text-2xl font-semibold text-gray-900">Dashboard de Inventario</h2>
      <p class="text-gray-600 mt-2">Resumen general del inventario por ubicaciones</p>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex justify-center items-center py-12">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600"></div>
      <span class="ml-3 text-gray-600">Cargando estadísticas...</span>
    </div>

    <!-- Error -->
    <div v-else-if="error" class="card">
      <div class="bg-red-50 border border-red-200 rounded-md p-4">
        <div class="flex">
          <svg class="h-5 w-5 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <div class="ml-3">
            <h3 class="text-sm font-medium text-red-800">Error al cargar estadísticas</h3>
            <div class="mt-2 text-sm text-red-700">{{ error }}</div>
            <button @click="loadDashboard" class="btn-secondary mt-4 text-sm">Reintentar</button>
          </div>
        </div>
      </div>
    </div>

    <!-- Dashboard Content -->
    <div v-else class="space-y-8">
      <!-- Estadísticas Generales -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 md:gap-6">
        <!-- Total en Bodegas -->
        <div class="bg-gradient-to-br from-blue-500 to-blue-600 rounded-lg shadow-lg p-6 text-white">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-blue-100 text-sm font-medium">Total en Bodegas</p>
              <p class="text-3xl font-bold mt-2">{{ summary.total_in_stores || 0 }}</p>
              <p class="text-blue-100 text-xs mt-1">unidades</p>
            </div>
            <div class="bg-blue-400 bg-opacity-30 rounded-full p-3">
              <svg class="h-8 w-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
              </svg>
            </div>
          </div>
        </div>

        <!-- Total en Pabellones -->
        <div class="bg-gradient-to-br from-green-500 to-green-600 rounded-lg shadow-lg p-6 text-white">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-green-100 text-sm font-medium">Total en Pabellones</p>
              <p class="text-3xl font-bold mt-2">{{ summary.total_in_pavilions || 0 }}</p>
              <p class="text-green-100 text-xs mt-1">unidades</p>
            </div>
            <div class="bg-green-400 bg-opacity-30 rounded-full p-3">
              <svg class="h-8 w-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
              </svg>
            </div>
          </div>
        </div>

        <!-- Total Consumido -->
        <div class="bg-gradient-to-br from-purple-500 to-purple-600 rounded-lg shadow-lg p-6 text-white">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-purple-100 text-sm font-medium">Total Consumido</p>
              <p class="text-3xl font-bold mt-2">{{ summary.total_consumed || 0 }}</p>
              <p class="text-purple-100 text-xs mt-1">unidades</p>
            </div>
            <div class="bg-purple-400 bg-opacity-30 rounded-full p-3">
              <svg class="h-8 w-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
          </div>
        </div>

        <!-- Total Transferido -->
        <div class="bg-gradient-to-br from-orange-500 to-orange-600 rounded-lg shadow-lg p-6 text-white">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-orange-100 text-sm font-medium">Total Transferido</p>
              <p class="text-3xl font-bold mt-2">{{ summary.total_transferred || 0 }}</p>
              <p class="text-orange-100 text-xs mt-1">unidades</p>
            </div>
            <div class="bg-orange-400 bg-opacity-30 rounded-full p-3">
              <svg class="h-8 w-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4" />
              </svg>
            </div>
          </div>
        </div>
      </div>

      <!-- Alertas -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4 md:gap-6">
        <!-- Stock Bajo -->
        <div class="card">
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-lg font-semibold text-gray-900">Stock Bajo</h3>
            <span class="flex items-center justify-center h-10 w-10 rounded-full bg-red-100">
              <svg class="h-6 w-6 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
              </svg>
            </span>
          </div>
          <p class="text-3xl font-bold text-gray-900">{{ summary.low_stock_stores || 0 }}</p>
          <p class="text-sm text-gray-600 mt-1">lotes con stock bajo</p>
        </div>

        <!-- Próximos a Vencer -->
        <div class="card">
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-lg font-semibold text-gray-900">Próximos a Vencer</h3>
            <span class="flex items-center justify-center h-10 w-10 rounded-full bg-yellow-100">
              <svg class="h-6 w-6 text-yellow-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </span>
          </div>
          <p class="text-3xl font-bold text-gray-900">{{ summary.near_expiration || 0 }}</p>
          <p class="text-sm text-gray-600 mt-1">lotes vencen en 90 días</p>
        </div>

        <!-- Transferencias Pendientes -->
        <div class="card">
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-lg font-semibold text-gray-900">Transferencias Pendientes</h3>
            <span class="flex items-center justify-center h-10 w-10 rounded-full bg-blue-100">
              <svg class="h-6 w-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6" />
              </svg>
            </span>
          </div>
          <p class="text-3xl font-bold text-gray-900">{{ summary.pending_transfers || 0 }}</p>
          <p class="text-sm text-gray-600 mt-1">transferencias en proceso</p>
        </div>
      </div>

      <!-- Inventario por Tipo de Cirugía -->
      <div class="card">
        <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 mb-6">
          <div>
            <h3 class="text-xl font-semibold text-gray-900">Inventario por Tipo de Cirugía</h3>
            <p class="text-gray-600 text-sm mt-2">Stock disponible organizado por procedimiento</p>
          </div>
          <button @click="loadSurgeryInventory" class="btn-secondary">
            <svg class="h-5 w-5 mr-2 inline" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
            Actualizar
          </button>
        </div>

        <div v-if="surgeryInventoryLoading" class="flex justify-center items-center py-8">
          <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
          <span class="ml-3 text-gray-600">Cargando...</span>
        </div>

        <div v-else-if="surgeryInventory.length === 0" class="text-center py-8">
          <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4" />
          </svg>
          <p class="mt-2 text-sm text-gray-500">No hay inventario organizado por tipo de cirugía</p>
        </div>

        <div v-else class="overflow-x-auto">
          <table class="min-w-full divide-y divide-gray-200">
            <thead class="bg-gray-50">
              <tr>
                <th @click="sortSurgeryBy('surgery_name')"
                  class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100 transition-colors select-none">
                  <div class="flex items-center space-x-1">
                    <span>Tipo de Cirugía</span>
                    <span class="flex flex-col -space-y-1">
                      <svg class="h-3 w-3 transition-colors" 
                        :class="surgerySortKey === 'surgery_name' && surgerySortOrder === 'asc' ? 'text-blue-600' : 'text-gray-300'" 
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
              <tr v-for="item in paginatedSurgeryInventory" :key="item.surgery_id" class="hover:bg-gray-50">
                <td class="px-6 py-4 whitespace-nowrap">
                  <div class="flex items-center">
                    <div class="flex-shrink-0 h-10 w-10 rounded-full bg-gradient-to-br from-blue-400 to-blue-600 flex items-center justify-center">
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
                  <span class="px-2 py-1 text-xs font-semibold rounded-full bg-blue-100 text-blue-800">
                    {{ item.batch_count || 0 }} lotes
                  </span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                  <router-link :to="`/inventory/store?surgery_id=${item.surgery_id}`" class="text-blue-600 hover:text-blue-900">
                    Ver detalles →
                  </router-link>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Paginación para tabla de cirugías -->
      <div v-if="!surgeryInventoryLoading && sortedSurgeryInventory.length > 0" class="card">
        <div class="flex flex-col sm:flex-row items-center justify-between gap-4">
          <div class="text-sm text-gray-700 text-center sm:text-left">
            Mostrando {{ surgeryStartIndex + 1 }} a {{ surgeryEndIndex }} de {{ sortedSurgeryInventory.length }} tipos de cirugía
          </div>
          <div class="flex items-center gap-2">
            <button class="btn-secondary px-3 py-2 text-sm min-w-[70px]" :disabled="surgeryCurrentPage === 1"
              @click="surgeryCurrentPage--">
              <span class="hidden sm:inline">Anterior</span>
              <span class="sm:hidden">Ant.</span>
            </button>
            <span class="px-3 py-2 text-sm text-gray-700 bg-gray-100 rounded-md min-w-[90px] text-center">
              Página {{ surgeryCurrentPage }} de {{ surgeryTotalPages }}
            </span>
            <button class="btn-secondary px-3 py-2 text-sm min-w-[70px]" :disabled="surgeryCurrentPage === surgeryTotalPages"
              @click="surgeryCurrentPage++">
              <span class="hidden sm:inline">Siguiente</span>
              <span class="sm:hidden">Sig.</span>
            </button>
          </div>
        </div>
      </div>

      <!-- Accesos Rápidos -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4 md:gap-6">
        <!-- Ver Inventario de Bodegas -->
        <router-link to="/inventory/store" class="card hover:shadow-lg transition-shadow cursor-pointer group">
          <div class="flex items-center">
            <div class="flex-shrink-0 h-12 w-12 rounded-lg bg-blue-100 flex items-center justify-center group-hover:bg-blue-200 transition-colors">
              <svg class="h-6 w-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
              </svg>
            </div>
            <div class="ml-4 flex-1">
              <h3 class="text-lg font-semibold text-gray-900 group-hover:text-blue-600 transition-colors">Inventario de Bodegas</h3>
              <p class="text-sm text-gray-600">Ver stock detallado por bodega</p>
            </div>
            <svg class="h-5 w-5 text-gray-400 group-hover:text-blue-600 transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24">
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
import inventoryService from '@/services/inventoryService'

const loading = ref(false)
const error = ref(null)
const summary = ref({})
const surgeryInventory = ref([])
const surgeryInventoryLoading = ref(false)

// Estado de ordenamiento para tabla de cirugías
const surgerySortKey = ref('surgery_name')
const surgerySortOrder = ref('asc')

// Estado de paginación para tabla de cirugías
const surgeryCurrentPage = ref(1)
const surgeryItemsPerPage = 10

const loadDashboard = async () => {
  loading.value = true
  error.value = null
  
  try {
    // Cargar resumen general
    const summaryData = await inventoryService.getInventorySummary()
    summary.value = summaryData
    
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
  } finally {
    surgeryInventoryLoading.value = false
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

// Computed properties para paginación de cirugías
const surgeryTotalPages = computed(() => Math.ceil(sortedSurgeryInventory.value.length / surgeryItemsPerPage))
const surgeryStartIndex = computed(() => (surgeryCurrentPage.value - 1) * surgeryItemsPerPage)
const surgeryEndIndex = computed(() => Math.min(surgeryStartIndex.value + surgeryItemsPerPage, sortedSurgeryInventory.value.length))

const paginatedSurgeryInventory = computed(() => {
  return sortedSurgeryInventory.value.slice(surgeryStartIndex.value, surgeryEndIndex.value)
})

// Función para ordenar por columna
const sortSurgeryBy = (key) => {
  if (surgerySortKey.value === key) {
    surgerySortOrder.value = surgerySortOrder.value === 'asc' ? 'desc' : 'asc'
  } else {
    surgerySortKey.value = key
    surgerySortOrder.value = 'asc'
  }
  surgeryCurrentPage.value = 1
}

onMounted(() => {
  loadDashboard()
})
</script>

