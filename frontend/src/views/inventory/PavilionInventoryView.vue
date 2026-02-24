<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="bg-white rounded-lg shadow-sm border p-6">
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <div class="flex items-center gap-3">
            <h2 class="text-2xl font-semibold text-gray-900">Inventario de Pabellones</h2>
            <span v-if="!loading && inventory.length > 0"
              class="inline-flex items-center px-2.5 py-1 rounded-full text-sm font-semibold bg-blue-100 text-blue-800">
              {{ sortedInventory.length }} lotes
            </span>
          </div>
          <p class="text-gray-600 mt-1">Stock disponible en cada pabellón del hospital</p>
        </div>
        <router-link to="/inventory/dashboard" class="btn-secondary flex items-center justify-center">
          <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
          </svg>
          Volver al Dashboard
        </router-link>
      </div>
    </div>

    <!-- Selector de Pabellón y Filtros -->
    <div class="card">
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mb-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Seleccionar Pabellón</label>
          <select v-model="selectedPavilionId" class="form-input" @change="loadInventory">
            <option value="">Seleccione un pabellón...</option>
            <option v-for="pavilion in pavilions" :key="pavilion.id" :value="pavilion.id">
              {{ pavilion.name }}
            </option>
          </select>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Proveedor</label>
          <input
            type="text"
            v-model="filterSupplier"
            placeholder="Buscar proveedor..."
            class="form-input"
            @input="debouncedApplyFilters"
            :disabled="!selectedPavilionId"
          />
        </div>
        <div class="flex items-end">
          <label class="flex items-center space-x-2">
            <input
              type="checkbox"
              v-model="includeInTransit"
              @change="loadInventory"
              :disabled="!selectedPavilionId"
              class="rounded border-gray-300 text-blue-600 shadow-sm focus:border-blue-300 focus:ring focus:ring-blue-200 focus:ring-opacity-50"
            />
            <span class="text-sm text-gray-700">Incluir en tránsito</span>
          </label>
        </div>
      </div>
      <div class="flex justify-end space-x-2">
        <button @click="clearFilters" class="btn-secondary" :disabled="!selectedPavilionId">Limpiar Filtros</button>
        <button @click="applyFilters" class="btn-primary" :disabled="!selectedPavilionId">Aplicar Filtros</button>
      </div>
    </div>

    <!-- Sin Selección -->
    <div v-if="!selectedPavilionId" class="card text-center py-12">
      <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
      </svg>
      <h3 class="mt-2 text-sm font-medium text-gray-900">Seleccione un pabellón</h3>
      <p class="mt-1 text-sm text-gray-500">
        Seleccione un pabellón del menú desplegable para ver su inventario
      </p>
    </div>

    <!-- Loading -->
    <div v-else-if="loading" class="card">
      <div class="flex justify-center items-center py-12">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
        <span class="ml-3 text-gray-600">Cargando inventario del pabellón...</span>
      </div>
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
            <button @click="loadInventory" class="btn-secondary mt-4 text-sm">Reintentar</button>
          </div>
        </div>
      </div>
    </div>

    <!-- Sin Inventario -->
    <div v-else-if="inventory.length === 0" class="card text-center py-12">
      <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4" />
      </svg>
      <h3 class="mt-2 text-sm font-medium text-gray-900">No hay inventario en este pabellón</h3>
      <p class="mt-1 text-sm text-gray-500">
        Este pabellón no tiene insumos registrados actualmente
      </p>
    </div>

    <!-- Tabla de Inventario -->
    <div v-else>
      <!-- Resumen del Pabellón -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-6">
        <div class="card">
          <div class="flex items-center">
            <div class="flex-shrink-0 h-12 w-12 rounded-lg bg-green-100 flex items-center justify-center">
              <svg class="h-6 w-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
            </div>
            <div class="ml-4">
              <p class="text-sm text-gray-600">Stock Disponible</p>
              <p class="text-2xl font-bold text-gray-900">{{ getTotalAvailable() }}</p>
              <p class="text-xs text-gray-500">unidades</p>
            </div>
          </div>
        </div>

        <div class="card">
          <div class="flex items-center">
            <div class="flex-shrink-0 h-12 w-12 rounded-lg bg-blue-100 flex items-center justify-center">
              <svg class="h-6 w-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4" />
              </svg>
            </div>
            <div class="ml-4">
              <p class="text-sm text-gray-600">Total Recibido</p>
              <p class="text-2xl font-bold text-gray-900">{{ getTotalReceived() }}</p>
              <p class="text-xs text-gray-500">unidades</p>
            </div>
          </div>
        </div>

        <div class="card">
          <div class="flex items-center">
            <div class="flex-shrink-0 h-12 w-12 rounded-lg bg-purple-100 flex items-center justify-center">
              <svg class="h-6 w-6 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
              </svg>
            </div>
            <div class="ml-4">
              <p class="text-sm text-gray-600">Total Consumido</p>
              <p class="text-2xl font-bold text-gray-900">{{ getTotalConsumed() }}</p>
              <p class="text-xs text-gray-500">unidades</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Tabla -->
      <div class="card overflow-hidden">
        <div class="overflow-x-auto">
          <div class="max-h-[600px] overflow-y-auto">
          <table class="min-w-full divide-y divide-gray-200">
            <thead class="bg-gray-50 sticky top-0 z-10">
              <tr>
                <th @click="sortBy('supply_name')"
                  class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100 transition-colors select-none">
                  <div class="flex items-center space-x-1">
                    <span>Insumo</span>
                    <span class="flex flex-col -space-y-1">
                      <svg class="h-3 w-3 transition-colors" 
                        :class="sortKey === 'supply_name' && sortOrder === 'asc' ? 'text-blue-600' : 'text-gray-300'" 
                        fill="currentColor" viewBox="0 0 20 20">
                        <path d="M5.293 9.707a1 1 0 010-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 01-1.414 1.414L11 7.414V15a1 1 0 11-2 0V7.414L6.707 9.707a1 1 0 01-1.414 0z"/>
                      </svg>
                      <svg class="h-3 w-3 transition-colors" 
                        :class="sortKey === 'supply_name' && sortOrder === 'desc' ? 'text-blue-600' : 'text-gray-300'" 
                        fill="currentColor" viewBox="0 0 20 20">
                        <path d="M14.707 10.293a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 111.414-1.414L9 12.586V5a1 1 0 012 0v7.586l2.293-2.293a1 1 0 011.414 0z"/>
                      </svg>
                    </span>
                  </div>
                </th>
                <th @click="sortBy('current_available')"
                  class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100 transition-colors select-none">
                  <div class="flex items-center space-x-1">
                    <span>Disponible</span>
                    <span class="flex flex-col -space-y-1">
                      <svg class="h-3 w-3 transition-colors" 
                        :class="sortKey === 'current_available' && sortOrder === 'asc' ? 'text-blue-600' : 'text-gray-300'" 
                        fill="currentColor" viewBox="0 0 20 20">
                        <path d="M5.293 9.707a1 1 0 010-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 01-1.414 1.414L11 7.414V15a1 1 0 11-2 0V7.414L6.707 9.707a1 1 0 01-1.414 0z"/>
                      </svg>
                      <svg class="h-3 w-3 transition-colors" 
                        :class="sortKey === 'current_available' && sortOrder === 'desc' ? 'text-blue-600' : 'text-gray-300'" 
                        fill="currentColor" viewBox="0 0 20 20">
                        <path d="M14.707 10.293a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 111.414-1.414L9 12.586V5a1 1 0 012 0v7.586l2.293-2.293a1 1 0 011.414 0z"/>
                      </svg>
                    </span>
                  </div>
                </th>
                <th @click="sortBy('total_received')"
                  class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100 transition-colors select-none">
                  <div class="flex items-center space-x-1">
                    <span>Total Recibido</span>
                    <span class="flex flex-col -space-y-1">
                      <svg class="h-3 w-3 transition-colors" 
                        :class="sortKey === 'total_received' && sortOrder === 'asc' ? 'text-blue-600' : 'text-gray-300'" 
                        fill="currentColor" viewBox="0 0 20 20">
                        <path d="M5.293 9.707a1 1 0 010-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 01-1.414 1.414L11 7.414V15a1 1 0 11-2 0V7.414L6.707 9.707a1 1 0 01-1.414 0z"/>
                      </svg>
                      <svg class="h-3 w-3 transition-colors" 
                        :class="sortKey === 'total_received' && sortOrder === 'desc' ? 'text-blue-600' : 'text-gray-300'" 
                        fill="currentColor" viewBox="0 0 20 20">
                        <path d="M14.707 10.293a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 111.414-1.414L9 12.586V5a1 1 0 012 0v7.586l2.293-2.293a1 1 0 011.414 0z"/>
                      </svg>
                    </span>
                  </div>
                </th>
                <th @click="sortBy('total_consumed')"
                  class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100 transition-colors select-none">
                  <div class="flex items-center space-x-1">
                    <span>Consumido</span>
                    <span class="flex flex-col -space-y-1">
                      <svg class="h-3 w-3 transition-colors" 
                        :class="sortKey === 'total_consumed' && sortOrder === 'asc' ? 'text-blue-600' : 'text-gray-300'" 
                        fill="currentColor" viewBox="0 0 20 20">
                        <path d="M5.293 9.707a1 1 0 010-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 01-1.414 1.414L11 7.414V15a1 1 0 11-2 0V7.414L6.707 9.707a1 1 0 01-1.414 0z"/>
                      </svg>
                      <svg class="h-3 w-3 transition-colors" 
                        :class="sortKey === 'total_consumed' && sortOrder === 'desc' ? 'text-blue-600' : 'text-gray-300'" 
                        fill="currentColor" viewBox="0 0 20 20">
                        <path d="M14.707 10.293a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 111.414-1.414L9 12.586V5a1 1 0 012 0v7.586l2.293-2.293a1 1 0 011.414 0z"/>
                      </svg>
                    </span>
                  </div>
                </th>
                <th @click="sortBy('total_returned')"
                  class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100 transition-colors select-none">
                  <div class="flex items-center space-x-1">
                    <span>Devuelto</span>
                    <span class="flex flex-col -space-y-1">
                      <svg class="h-3 w-3 transition-colors" 
                        :class="sortKey === 'total_returned' && sortOrder === 'asc' ? 'text-blue-600' : 'text-gray-300'" 
                        fill="currentColor" viewBox="0 0 20 20">
                        <path d="M5.293 9.707a1 1 0 010-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 01-1.414 1.414L11 7.414V15a1 1 0 11-2 0V7.414L6.707 9.707a1 1 0 01-1.414 0z"/>
                      </svg>
                      <svg class="h-3 w-3 transition-colors" 
                        :class="sortKey === 'total_returned' && sortOrder === 'desc' ? 'text-blue-600' : 'text-gray-300'" 
                        fill="currentColor" viewBox="0 0 20 20">
                        <path d="M14.707 10.293a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 111.414-1.414L9 12.586V5a1 1 0 012 0v7.586l2.293-2.293a1 1 0 011.414 0z"/>
                      </svg>
                    </span>
                  </div>
                </th>
                <th @click="sortBy('batch_supplier')"
                  class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100 transition-colors select-none">
                  <div class="flex items-center space-x-1">
                    <span>Proveedor</span>
                    <span class="flex flex-col -space-y-1">
                      <svg class="h-3 w-3 transition-colors" 
                        :class="sortKey === 'batch_supplier' && sortOrder === 'asc' ? 'text-blue-600' : 'text-gray-300'" 
                        fill="currentColor" viewBox="0 0 20 20">
                        <path d="M5.293 9.707a1 1 0 010-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 01-1.414 1.414L11 7.414V15a1 1 0 11-2 0V7.414L6.707 9.707a1 1 0 01-1.414 0z"/>
                      </svg>
                      <svg class="h-3 w-3 transition-colors" 
                        :class="sortKey === 'batch_supplier' && sortOrder === 'desc' ? 'text-blue-600' : 'text-gray-300'" 
                        fill="currentColor" viewBox="0 0 20 20">
                        <path d="M14.707 10.293a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 111.414-1.414L9 12.586V5a1 1 0 012 0v7.586l2.293-2.293a1 1 0 011.414 0z"/>
                      </svg>
                    </span>
                  </div>
                </th>
                <th @click="sortBy('expiration_date')"
                  class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100 transition-colors select-none">
                  <div class="flex items-center space-x-1">
                    <span>F. Vencimiento</span>
                    <span class="flex flex-col -space-y-1">
                      <svg class="h-3 w-3 transition-colors" 
                        :class="sortKey === 'expiration_date' && sortOrder === 'asc' ? 'text-blue-600' : 'text-gray-300'" 
                        fill="currentColor" viewBox="0 0 20 20">
                        <path d="M5.293 9.707a1 1 0 010-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 01-1.414 1.414L11 7.414V15a1 1 0 11-2 0V7.414L6.707 9.707a1 1 0 01-1.414 0z"/>
                      </svg>
                      <svg class="h-3 w-3 transition-colors" 
                        :class="sortKey === 'expiration_date' && sortOrder === 'desc' ? 'text-blue-600' : 'text-gray-300'" 
                        fill="currentColor" viewBox="0 0 20 20">
                        <path d="M14.707 10.293a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 111.414-1.414L9 12.586V5a1 1 0 012 0v7.586l2.293-2.293a1 1 0 011.414 0z"/>
                      </svg>
                    </span>
                  </div>
                </th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
              <tr v-for="item in paginatedInventory" :key="item.id"
                :class="item.in_transit ? 'bg-blue-50 hover:bg-blue-100' : 'hover:bg-gray-50'">
                <td class="px-6 py-4 whitespace-nowrap">
                  <div class="flex items-center">
                    <div class="flex-shrink-0 h-10 w-10 rounded-lg bg-gradient-to-br from-green-400 to-green-600 flex items-center justify-center">
                      <svg class="h-5 w-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
                      </svg>
                    </div>
                    <div class="ml-3">
                      <div class="flex items-center gap-2">
                        <span class="text-sm font-medium text-gray-900">{{ item.supply_name }}</span>
                        <span v-if="item.in_transit"
                          class="inline-flex items-center px-1.5 py-0.5 rounded text-xs font-medium bg-blue-100 text-blue-700">
                          🚚 En tránsito
                        </span>
                      </div>
                      <div class="text-sm text-gray-500">Código: {{ item.supply_code }}</div>
                      <div class="text-sm text-gray-500">Lote: {{ item.batch_id }}</div>
                    </div>
                  </div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <div class="flex items-center">
                    <span :class="getAvailabilityClass(item.current_available)" class="text-lg font-bold">
                      {{ item.current_available }}
                    </span>
                    <span class="text-xs text-gray-500 ml-1">unidades</span>
                  </div>
                  <div v-if="item.current_available < 5" class="text-xs text-red-600 font-medium">
                    ⚠️ Stock bajo
                  </div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {{ item.total_received }} unidades
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {{ item.total_consumed || 0 }} unidades
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {{ item.total_returned || 0 }} unidades
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {{ item.batch_supplier }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <span :class="getExpirationClass(item.expiration_date)" class="text-sm font-medium">
                    {{ formatDate(item.expiration_date) }}
                  </span>
                  <div v-if="isNearExpiration(item.expiration_date)" class="text-xs text-orange-600 font-medium">
                    ⚠️ Vence pronto
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
          </div>
        </div>
      </div>

      <!-- Paginación -->
      <div v-if="!loading && sortedInventory.length > 0" class="card">
        <div class="flex flex-col sm:flex-row items-center justify-between gap-4">
          <div class="text-sm text-gray-700 text-center sm:text-left">
            Mostrando {{ startIndex + 1 }} a {{ endIndex }} de {{ sortedInventory.length }} lotes
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
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import inventoryService from '@/services/inventory/inventoryService'

const loading = ref(false)
const error = ref(null)
const inventory = ref([])
const pavilions = ref([])
const selectedPavilionId = ref('')
const includeInTransit = ref(false)
const filterSupplier = ref('')

// Paginación
const currentPage = ref(1)
const itemsPerPage = 10

// Estado de ordenamiento
const sortKey = ref('supply_name')
const sortOrder = ref('asc')

// Normaliza texto para comparación case-insensitive y sin tildes (ID 6)
const normalizeText = (text) => {
  if (!text) return ''
  return text.toString().toLowerCase().normalize('NFD').replace(/[\u0300-\u036f]/g, '')
}

// Computed para ordenar y filtrar inventario
const sortedInventory = computed(() => {
  if (!inventory.value || inventory.value.length === 0) return []

  // Filtro cliente por proveedor (ID 6)
  let result = [...inventory.value]
  if (filterSupplier.value.trim()) {
    const needle = normalizeText(filterSupplier.value)
    result = result.filter(item => normalizeText(item.batch_supplier || '').includes(needle))
  }

  const sorted = result.sort((a, b) => {
    let aVal = a[sortKey.value]
    let bVal = b[sortKey.value]
    
    // Manejo de valores nulos/undefined
    if (aVal === null || aVal === undefined) aVal = ''
    if (bVal === null || bVal === undefined) bVal = ''
    
    // Manejo de strings (comparación case-insensitive)
    if (typeof aVal === 'string') {
      aVal = aVal.toLowerCase()
      bVal = (bVal || '').toString().toLowerCase()
    }
    
    // Comparación
    if (aVal < bVal) return sortOrder.value === 'asc' ? -1 : 1
    if (aVal > bVal) return sortOrder.value === 'asc' ? 1 : -1
    return 0
  })

  return sorted
})

const totalPages = computed(() => Math.ceil(sortedInventory.value.length / itemsPerPage))
const startIndex = computed(() => (currentPage.value - 1) * itemsPerPage)
const endIndex = computed(() => Math.min(startIndex.value + itemsPerPage, sortedInventory.value.length))

const paginatedInventory = computed(() => {
  return sortedInventory.value.slice(startIndex.value, endIndex.value)
})

// Función para ordenar por columna
const sortBy = (key) => {
  if (sortKey.value === key) {
    sortOrder.value = sortOrder.value === 'asc' ? 'desc' : 'asc'
  } else {
    sortKey.value = key
    sortOrder.value = 'asc'
  }
  currentPage.value = 1
}

const loadPavilions = async () => {
  try {
    const data = await inventoryService.getAllPavilions()
    pavilions.value = Array.isArray(data) ? data : []
  } catch (err) {
    console.error('Error loading pavilions:', err)
    pavilions.value = []
  }
}

const loadInventory = async () => {
  if (!selectedPavilionId.value) {
    inventory.value = []
    return
  }

  loading.value = true
  error.value = null
  currentPage.value = 1
  
  try {
    // El proveedor se filtra en cliente para evitar problemas de mayúsculas/tildes (ID 6)
    const data = await inventoryService.getPavilionInventory(
      selectedPavilionId.value,
      includeInTransit.value,
      null
    )
    // Asegurarse de que inventory.value siempre sea un array
    inventory.value = Array.isArray(data) ? data : []
  } catch (err) {
    error.value = err.message || 'Error al cargar inventario del pabellón'
    console.error('Error loading pavilion inventory:', err)
    // En caso de error, asegurar que sea un array vacío
    inventory.value = []
  } finally {
    loading.value = false
  }
}

// Función para aplicar filtros
const applyFilters = () => {
  loadInventory()
}

// Función para limpiar filtros — resetea TODOS los filtros incluyendo pabellón (ID 10)
const clearFilters = () => {
  filterSupplier.value = ''
  selectedPavilionId.value = ''
  includeInTransit.value = false
  inventory.value = []
  error.value = null
  currentPage.value = 1
}

// Debounce para el filtro de proveedor.
// El filtrado es reactivo (computed), solo resetea la paginación (ID 6).
let debounceTimer = null
const debouncedApplyFilters = () => {
  clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => {
    currentPage.value = 1
  }, 300)
}

const getTotalAvailable = () => {
  return inventory.value.reduce((sum, item) => sum + (item.current_available || 0), 0)
}

const getTotalReceived = () => {
  return inventory.value.reduce((sum, item) => sum + (item.total_received || 0), 0)
}

const getTotalConsumed = () => {
  return inventory.value.reduce((sum, item) => sum + (item.total_consumed || 0), 0)
}

const getAvailabilityClass = (available) => {
  if (available === 0) return 'text-red-600'
  if (available < 5) return 'text-orange-600'
  if (available < 10) return 'text-yellow-600'
  return 'text-green-600'
}

const formatDate = (dateString) => {
  if (!dateString) return 'N/A'
  try {
    return new Date(dateString).toLocaleDateString('es-CL')
  } catch {
    return dateString
  }
}

const isNearExpiration = (expirationDate) => {
  if (!expirationDate) return false
  const today = new Date()
  const expDate = new Date(expirationDate)
  const daysUntilExpiration = Math.ceil((expDate - today) / (1000 * 60 * 60 * 24))
  return daysUntilExpiration >= 0 && daysUntilExpiration <= 90
}

const getExpirationClass = (expirationDate) => {
  if (!expirationDate) return 'text-gray-900'
  const today = new Date()
  const expDate = new Date(expirationDate)
  const daysUntilExpiration = Math.ceil((expDate - today) / (1000 * 60 * 60 * 24))

  if (daysUntilExpiration < 0) return 'text-red-600'
  if (daysUntilExpiration <= 30) return 'text-red-600'
  if (daysUntilExpiration <= 90) return 'text-orange-600'
  return 'text-gray-900'
}

onMounted(async () => {
  await loadPavilions()
})
</script>

