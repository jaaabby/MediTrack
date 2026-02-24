<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="bg-white rounded-lg shadow-sm border p-6">
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h2 class="text-2xl font-semibold text-gray-900">Inventario de Bodegas</h2>
          <p class="text-gray-600 mt-1">Stock detallado en cada bodega del sistema</p>
          <p v-if="!loading && inventory.length > 0" class="text-sm text-gray-500 mt-1">
            Total: {{ inventory.length }} lotes encontrados
          </p>
        </div>
        <router-link to="/inventory/dashboard" class="btn-secondary flex items-center justify-center">
          <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
          </svg>
          Volver al Dashboard
        </router-link>
      </div>
    </div>

    <!-- Filtros -->
    <div class="card">
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Bodega</label>
          <select v-model="filters.store_id" class="form-input" @change="applyFilters">
            <option value="">Todas las bodegas</option>
            <option v-for="store in stores" :key="store.id" :value="store.id">
              {{ store.name }}
            </option>
          </select>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Tipo de Cirugía</label>
          <select v-model="filters.surgery_id" class="form-input" @change="applyFilters">
            <option value="">Todos los tipos</option>
            <option v-for="surgery in surgeries" :key="surgery.id" :value="surgery.id">
              {{ surgery.name }}
            </option>
          </select>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Proveedor</label>
          <input
            type="text"
            v-model="filters.supplier"
            placeholder="Buscar proveedor..."
            class="form-input"
            @input="debouncedApplyFilters"
          />
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Alertas</label>
          <div class="flex flex-col items-start gap-5">
            <!-- ID 7: inline-flex restringe el área clickeable a la casilla y su etiqueta -->
            <label class="inline-flex items-center gap-2 cursor-pointer">
              <input
                type="checkbox"
                v-model="filters.low_stock"
                @change="applyFilters"
                class="h-4 w-4 rounded border-gray-300 text-blue-600 shadow-sm focus:border-blue-300 focus:ring focus:ring-blue-200 focus:ring-opacity-50"
              />
              <span class="text-sm text-gray-700">Stock bajo</span>
            </label>
            <label class="inline-flex items-center gap-2 cursor-pointer">
              <input
                type="checkbox"
                v-model="filters.near_expiration"
                @change="applyFilters"
                class="h-4 w-4 rounded border-gray-300 text-blue-600 shadow-sm focus:border-blue-300 focus:ring focus:ring-blue-200 focus:ring-opacity-50"
              />
              <span class="text-sm text-gray-700">Próximo a vencer</span>
            </label>
          </div>
        </div>
      </div>
      <div class="mt-4 flex justify-end space-x-2">
        <button @click="clearFilters" class="btn-secondary">Limpiar Filtros</button>
        <button @click="applyFilters" class="btn-primary">Aplicar Filtros</button>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="card">
      <div class="flex justify-center items-center py-12">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
        <span class="ml-3 text-gray-600">Cargando inventario...</span>
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

    <!-- Tabla -->
    <div v-else-if="inventory.length === 0" class="card text-center py-12">
      <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4" />
      </svg>
      <h3 class="mt-2 text-sm font-medium text-gray-900">No se encontró inventario</h3>
      <p class="mt-1 text-sm text-gray-500">
        No hay lotes que coincidan con los filtros aplicados
      </p>
    </div>

    <div v-else class="card overflow-hidden">
      <div class="overflow-x-auto">
        <div class="max-h-[600px] overflow-y-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50 sticky top-0 z-10">
            <tr>
              <th @click="sortBy('store_name')"
                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100 transition-colors select-none">
                <div class="flex items-center space-x-1">
                  <span>Bodega</span>
                  <span class="flex flex-col -space-y-1">
                    <svg class="h-3 w-3 transition-colors" 
                      :class="sortKey === 'store_name' && sortOrder === 'asc' ? 'text-blue-600' : 'text-gray-300'" 
                      fill="currentColor" viewBox="0 0 20 20">
                      <path d="M5.293 9.707a1 1 0 010-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 01-1.414 1.414L11 7.414V15a1 1 0 11-2 0V7.414L6.707 9.707a1 1 0 01-1.414 0z"/>
                    </svg>
                    <svg class="h-3 w-3 transition-colors" 
                      :class="sortKey === 'store_name' && sortOrder === 'desc' ? 'text-blue-600' : 'text-gray-300'" 
                      fill="currentColor" viewBox="0 0 20 20">
                      <path d="M14.707 10.293a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 111.414-1.414L9 12.586V5a1 1 0 012 0v7.586l2.293-2.293a1 1 0 011.414 0z"/>
                    </svg>
                  </span>
                </div>
              </th>
              <th @click="sortBy('created_at')"
                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100 transition-colors select-none">
                <div class="flex items-center space-x-1">
                  <span>F. Ingreso</span>
                  <span class="flex flex-col -space-y-1">
                    <svg class="h-3 w-3 transition-colors"
                      :class="sortKey === 'created_at' && sortOrder === 'asc' ? 'text-blue-600' : 'text-gray-300'"
                      fill="currentColor" viewBox="0 0 20 20">
                      <path d="M5.293 9.707a1 1 0 010-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 01-1.414 1.414L11 7.414V15a1 1 0 11-2 0V7.414L6.707 9.707a1 1 0 01-1.414 0z"/>
                    </svg>
                    <svg class="h-3 w-3 transition-colors"
                      :class="sortKey === 'created_at' && sortOrder === 'desc' ? 'text-blue-600' : 'text-gray-300'"
                      fill="currentColor" viewBox="0 0 20 20">
                      <path d="M14.707 10.293a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 111.414-1.414L9 12.586V5a1 1 0 012 0v7.586l2.293-2.293a1 1 0 011.414 0z"/>
                    </svg>
                  </span>
                </div>
              </th>
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
              <th @click="sortBy('surgery_name')"
                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100 transition-colors select-none">
                <div class="flex items-center space-x-1">
                  <span>Tipo Cirugía</span>
                  <span class="flex flex-col -space-y-1">
                    <svg class="h-3 w-3 transition-colors" 
                      :class="sortKey === 'surgery_name' && sortOrder === 'asc' ? 'text-blue-600' : 'text-gray-300'" 
                      fill="currentColor" viewBox="0 0 20 20">
                      <path d="M5.293 9.707a1 1 0 010-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 01-1.414 1.414L11 7.414V15a1 1 0 11-2 0V7.414L6.707 9.707a1 1 0 01-1.414 0z"/>
                    </svg>
                    <svg class="h-3 w-3 transition-colors" 
                      :class="sortKey === 'surgery_name' && sortOrder === 'desc' ? 'text-blue-600' : 'text-gray-300'" 
                      fill="currentColor" viewBox="0 0 20 20">
                      <path d="M14.707 10.293a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 111.414-1.414L9 12.586V5a1 1 0 012 0v7.586l2.293-2.293a1 1 0 011.414 0z"/>
                    </svg>
                  </span>
                </div>
              </th>
              <th @click="sortBy('current_in_store')"
                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100 transition-colors select-none">
                <div class="flex items-center space-x-1">
                  <span>Stock Actual</span>
                  <span class="flex flex-col -space-y-1">
                    <svg class="h-3 w-3 transition-colors" 
                      :class="sortKey === 'current_in_store' && sortOrder === 'asc' ? 'text-blue-600' : 'text-gray-300'" 
                      fill="currentColor" viewBox="0 0 20 20">
                      <path d="M5.293 9.707a1 1 0 010-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 01-1.414 1.414L11 7.414V15a1 1 0 11-2 0V7.414L6.707 9.707a1 1 0 01-1.414 0z"/>
                    </svg>
                    <svg class="h-3 w-3 transition-colors" 
                      :class="sortKey === 'current_in_store' && sortOrder === 'desc' ? 'text-blue-600' : 'text-gray-300'" 
                      fill="currentColor" viewBox="0 0 20 20">
                      <path d="M14.707 10.293a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 111.414-1.414L9 12.586V5a1 1 0 012 0v7.586l2.293-2.293a1 1 0 011.414 0z"/>
                    </svg>
                  </span>
                </div>
              </th>
              <th @click="sortBy('original_amount')"
                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100 transition-colors select-none">
                <div class="flex items-center space-x-1">
                  <span>Stock Original</span>
                  <span class="flex flex-col -space-y-1">
                    <svg class="h-3 w-3 transition-colors" 
                      :class="sortKey === 'original_amount' && sortOrder === 'asc' ? 'text-blue-600' : 'text-gray-300'" 
                      fill="currentColor" viewBox="0 0 20 20">
                      <path d="M5.293 9.707a1 1 0 010-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 01-1.414 1.414L11 7.414V15a1 1 0 11-2 0V7.414L6.707 9.707a1 1 0 01-1.414 0z"/>
                    </svg>
                    <svg class="h-3 w-3 transition-colors" 
                      :class="sortKey === 'original_amount' && sortOrder === 'desc' ? 'text-blue-600' : 'text-gray-300'" 
                      fill="currentColor" viewBox="0 0 20 20">
                      <path d="M14.707 10.293a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 111.414-1.414L9 12.586V5a1 1 0 012 0v7.586l2.293-2.293a1 1 0 011.414 0z"/>
                    </svg>
                  </span>
                </div>
              </th>
              <th @click="sortBy('total_transferred_out')"
                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100 transition-colors select-none">
                <div class="flex items-center space-x-1">
                  <span>Movimientos</span>
                  <span class="flex flex-col -space-y-1">
                    <svg class="h-3 w-3 transition-colors" 
                      :class="sortKey === 'total_transferred_out' && sortOrder === 'asc' ? 'text-blue-600' : 'text-gray-300'" 
                      fill="currentColor" viewBox="0 0 20 20">
                      <path d="M5.293 9.707a1 1 0 010-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 01-1.414 1.414L11 7.414V15a1 1 0 11-2 0V7.414L6.707 9.707a1 1 0 01-1.414 0z"/>
                    </svg>
                    <svg class="h-3 w-3 transition-colors" 
                      :class="sortKey === 'total_transferred_out' && sortOrder === 'desc' ? 'text-blue-600' : 'text-gray-300'" 
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
            <!-- ID 22/23: fondo de fila según nivel de stock -->
            <tr v-for="item in paginatedInventory" :key="item.id" :class="getRowClass(item.current_in_store, item.original_amount)">
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="flex items-center">
                  <div class="flex-shrink-0 h-10 w-10 rounded-lg bg-blue-100 flex items-center justify-center">
                    <svg class="h-5 w-5 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
                    </svg>
                  </div>
                  <div class="ml-3">
                    <div class="text-sm font-medium text-gray-900">{{ item.store_name }}</div>
                    <div class="text-sm text-gray-500">ID: {{ item.store_id }}</div>
                  </div>
                </div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                {{ formatDate(item.created_at) }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="text-sm font-medium text-gray-900">{{ item.supply_name }}</div>
                <div class="text-sm text-gray-500">Código: {{ item.supply_code }}</div>
                <div class="text-sm text-gray-500">Lote: {{ item.batch_id }}</div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span v-if="item.surgery_name" class="px-2 py-1 text-xs font-semibold rounded-full bg-purple-100 text-purple-800">
                  {{ item.surgery_name }}
                </span>
                <span v-else class="text-sm text-gray-400">Sin asignar</span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="flex items-center">
                  <span :class="getStockClass(item.current_in_store, item.original_amount)" class="text-lg font-bold">
                    {{ item.current_in_store }}
                  </span>
                  <span class="text-xs text-gray-500 ml-1">unidades</span>
                </div>
                <!-- ID 22/23: badges diferenciados por nivel de stock -->
                <div v-if="isLowStock(item.current_in_store, item.original_amount)" class="text-xs text-red-600 font-medium">
                  ⚠️ Stock bajo
                </div>
                <div v-else-if="isMediumStock(item.current_in_store, item.original_amount)" class="text-xs text-orange-600 font-medium">
                  ⚠️ Stock medio
                </div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                {{ item.original_amount }} unidades
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="text-sm text-gray-900">Transferidos: {{ item.total_transferred_out || 0 }}</div>
                <div class="text-xs text-gray-500">
                  Consumidos: {{ (item.total_consumed_in_store || 0) + (item.total_consumed_from_pavilions || 0) }}
                </div>
                <div class="text-xs text-gray-500">Devueltos: {{ item.total_returned_in || 0 }}</div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                {{ item.batch_supplier }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span :class="getExpirationClass(item.expiration_date)" class="text-sm font-medium">
                  {{ formatDate(item.expiration_date) }}
                </span>
                <!-- ID 26: badges distintos para vencido vs próximo a vencer -->
                <div v-if="isExpired(item.expiration_date)" class="text-xs text-red-700 font-semibold">
                  🔴 Vencido
                </div>
                <div v-else-if="isNearExpiration(item.expiration_date)" class="text-xs text-orange-600 font-medium">
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
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import inventoryService from '@/services/inventory/inventoryService'
import surgeryService from '@/services/management/surgeryService'

const route = useRoute()

const loading = ref(false)
const error = ref(null)
const inventory = ref([])
const stores = ref([])
const surgeries = ref([])

const filters = ref({
  store_id: '',
  surgery_id: '',
  supplier: '',
  low_stock: false,
  near_expiration: false
})

// Paginación
const currentPage = ref(1)
const itemsPerPage = 10

// Estado de ordenamiento – ID 13: orden inicial por fecha de ingreso, más reciente primero
const sortKey = ref('created_at')
const sortOrder = ref('desc')

// ID 6: normaliza texto eliminando tildes y convirtiendo a minúsculas
const normalizeText = (text) => {
  if (!text) return ''
  return text
    .toLowerCase()
    .normalize('NFD')
    .replace(/[\u0300-\u036f]/g, '')
}

// Computed para ordenar inventario
const sortedInventory = computed(() => {
  if (!inventory.value || inventory.value.length === 0) return []

  // ID 6: filtro client-side de proveedor, insensible a mayúsculas y tildes
  const supplierQuery = normalizeText(filters.value.supplier)
  const base = supplierQuery
    ? inventory.value.filter(item =>
        normalizeText(item.batch_supplier).includes(supplierQuery)
      )
    : inventory.value

  const sorted = [...base].sort((a, b) => {
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

let debounceTimeout = null
const debouncedApplyFilters = () => {
  if (debounceTimeout) clearTimeout(debounceTimeout)
  debounceTimeout = setTimeout(() => {
    applyFilters()
  }, 500)
}

const loadInventory = async () => {
  loading.value = true
  error.value = null
  
  try {
    // El filtro de proveedor se aplica client-side con normalización de tildes y mayúsculas (ID 6)
    // No se envía al backend para evitar el LIKE case-sensitive de PostgreSQL
    const backendFilters = { ...filters.value, supplier: '' }
    const data = await inventoryService.getStoreInventory(backendFilters)
    // Asegurarse de que inventory.value siempre sea un array
    inventory.value = Array.isArray(data) ? data : []
  } catch (err) {
    error.value = err.message || 'Error al cargar inventario de bodegas'
    console.error('Error loading store inventory:', err)
    // En caso de error, asegurar que sea un array vacío
    inventory.value = []
  } finally {
    loading.value = false
  }
}

const loadStores = async () => {
  try {
    const data = await inventoryService.getAllStores()
    stores.value = Array.isArray(data) ? data : []
  } catch (err) {
    console.error('Error loading stores:', err)
    stores.value = []
  }
}

const loadSurgeries = async () => {
  try {
    const data = await surgeryService.getAllSurgeries()
    surgeries.value = Array.isArray(data) ? data : []
  } catch (err) {
    console.error('Error loading surgeries:', err)
    surgeries.value = []
  }
}

const applyFilters = () => {
  currentPage.value = 1
  loadInventory()
}

const clearFilters = () => {
  filters.value = {
    store_id: '',
    surgery_id: '',
    supplier: '',
    low_stock: false,
    near_expiration: false
  }
  applyFilters()
}

const formatDate = (dateString) => {
  if (!dateString) return 'N/A'
  try {
    return new Date(dateString).toLocaleDateString('es-CL')
  } catch {
    return dateString
  }
}

// ID 22: stock bajo = menos del 20% del original
const isLowStock = (current, original) => {
  if (!original) return false
  return current < original * 0.2
}

// ID 23: stock medio = entre 20% y 50% del original
const isMediumStock = (current, original) => {
  if (!original) return false
  return current >= original * 0.2 && current < original * 0.5
}

// ID 22/23: color del número de stock (prioridad: rojo > naranja > verde)
const getStockClass = (current, original) => {
  if (isLowStock(current, original)) return 'text-red-600'
  if (isMediumStock(current, original)) return 'text-orange-600'
  return 'text-green-600'
}

// ID 22/23: color de fondo de la fila según nivel de stock
const getRowClass = (current, original) => {
  if (isLowStock(current, original)) return 'bg-red-50 hover:bg-red-100'
  if (isMediumStock(current, original)) return 'bg-orange-50 hover:bg-orange-100'
  return 'hover:bg-gray-50'
}

// ID 26: lote con fecha ya vencida
const isExpired = (expirationDate) => {
  if (!expirationDate) return false
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  return new Date(expirationDate) < today
}

// ID 26: próximo a vencer = dentro de los próximos 30 días (sin estar vencido)
const isNearExpiration = (expirationDate) => {
  if (!expirationDate) return false
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  const expDate = new Date(expirationDate)
  const days = Math.ceil((expDate - today) / (1000 * 60 * 60 * 24))
  return days >= 0 && days <= 30
}

// ID 26: colores diferenciados: rojo oscuro = vencido, naranja = próximo a vencer
const getExpirationClass = (expirationDate) => {
  if (!expirationDate) return 'text-gray-900'
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  const expDate = new Date(expirationDate)
  const days = Math.ceil((expDate - today) / (1000 * 60 * 60 * 24))

  if (days < 0) return 'text-red-700 font-semibold'   // Vencido
  if (days <= 30) return 'text-orange-600 font-semibold' // Próximo a vencer
  if (days <= 90) return 'text-yellow-600'               // Vence en 1-3 meses
  return 'text-gray-900'
}

onMounted(async () => {
  // Cargar datos auxiliares
  await Promise.all([loadStores(), loadSurgeries()])
  
  // Aplicar filtros de la URL si existen
  if (route.query.surgery_id) {
    filters.value.surgery_id = route.query.surgery_id
  }
  if (route.query.store_id) {
    filters.value.store_id = route.query.store_id
  }
  
  // Cargar inventario
  await loadInventory()
})
</script>

