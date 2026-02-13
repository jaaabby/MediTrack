<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="bg-white rounded-lg shadow-sm border p-6">
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h2 class="text-2xl font-semibold text-gray-900">Historial de Insumos</h2>
          <p class="text-gray-600 mt-1">Rastrea todos los movimientos y cambios de estado de los insumos</p>
          <p v-if="!loading" class="text-sm text-gray-500 mt-1">Total: {{ filteredHistory.length }} registros</p>
        </div>
        <button 
          @click="exportToExcel" 
          :disabled="loading || sortedHistory.length === 0"
          class="btn-secondary flex items-center justify-center"
          :class="{ 'opacity-50 cursor-not-allowed': loading || sortedHistory.length === 0 }"
        >
          <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
          </svg>
          Exportar a Excel
        </button>
      </div>
    </div>

    <!-- Filtros -->
    <div class="card">
      <div class="flex flex-col md:flex-row gap-4 items-end">
        <div class="flex-1">
          <label class="block text-sm font-medium text-gray-700 mb-2">Buscar</label>
          <input type="text" v-model="searchTerm" placeholder="Nombre, QR, estado, usuario..." class="form-input" />
        </div>
        <div class="flex-1">
          <label class="block text-sm font-medium text-gray-700 mb-2">Desde</label>
          <input type="date" v-model="filters.from_date" class="form-input" />
        </div>
        <div class="flex-1">
          <label class="block text-sm font-medium text-gray-700 mb-2">Hasta</label>
          <input type="date" v-model="filters.to_date" class="form-input" />
        </div>
        <div>
          <button 
            @click="clearFilters" 
            :disabled="!hasActiveFilters"
            class="btn-secondary flex items-center justify-center whitespace-nowrap"
            :class="{ 'opacity-50 cursor-not-allowed': !hasActiveFilters }"
          >
            <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
            Limpiar Filtros
          </button>
        </div>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="card">
      <div class="flex justify-center items-center py-12">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
        <span class="ml-3">Cargando historial...</span>
      </div>
    </div>

    <!-- Tabla -->
    <div v-else class="card overflow-hidden">
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th scope="col" @click="sortBy('id')"
                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100 transition-colors select-none">
                <div class="flex items-center space-x-1">
                  <span>ID</span>
                  <span class="flex flex-col -space-y-1">
                    <svg class="h-3 w-3 transition-colors" 
                      :class="sortKey === 'id' && sortOrder === 'asc' ? 'text-blue-600' : 'text-gray-300'" 
                      fill="currentColor" viewBox="0 0 20 20">
                      <path d="M5.293 9.707a1 1 0 010-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 01-1.414 1.414L11 7.414V15a1 1 0 11-2 0V7.414L6.707 9.707a1 1 0 01-1.414 0z"/>
                    </svg>
                    <svg class="h-3 w-3 transition-colors" 
                      :class="sortKey === 'id' && sortOrder === 'desc' ? 'text-blue-600' : 'text-gray-300'" 
                      fill="currentColor" viewBox="0 0 20 20">
                      <path d="M14.707 10.293a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 111.414-1.414L9 12.586V5a1 1 0 012 0v7.586l2.293-2.293a1 1 0 011.414 0z"/>
                    </svg>
                  </span>
                </div>
              </th>
              <th scope="col" @click="sortBy('supply_name')"
                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100 transition-colors select-none">
                <div class="flex items-center space-x-1">
                  <span>Nombre Insumo</span>
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
              <th scope="col" @click="sortBy('qr_code')"
                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100 transition-colors select-none">
                <div class="flex items-center space-x-1">
                  <span>QR Code</span>
                  <span class="flex flex-col -space-y-1">
                    <svg class="h-3 w-3 transition-colors" 
                      :class="sortKey === 'qr_code' && sortOrder === 'asc' ? 'text-blue-600' : 'text-gray-300'" 
                      fill="currentColor" viewBox="0 0 20 20">
                      <path d="M5.293 9.707a1 1 0 010-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 01-1.414 1.414L11 7.414V15a1 1 0 11-2 0V7.414L6.707 9.707a1 1 0 01-1.414 0z"/>
                    </svg>
                    <svg class="h-3 w-3 transition-colors" 
                      :class="sortKey === 'qr_code' && sortOrder === 'desc' ? 'text-blue-600' : 'text-gray-300'" 
                      fill="currentColor" viewBox="0 0 20 20">
                      <path d="M14.707 10.293a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 111.414-1.414L9 12.586V5a1 1 0 012 0v7.586l2.293-2.293a1 1 0 011.414 0z"/>
                    </svg>
                  </span>
                </div>
              </th>
              <th scope="col" @click="sortBy('status')"
                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100 transition-colors select-none">
                <div class="flex items-center space-x-1">
                  <span>Estado</span>
                  <span class="flex flex-col -space-y-1">
                    <svg class="h-3 w-3 transition-colors" 
                      :class="sortKey === 'status' && sortOrder === 'asc' ? 'text-blue-600' : 'text-gray-300'" 
                      fill="currentColor" viewBox="0 0 20 20">
                      <path d="M5.293 9.707a1 1 0 010-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 01-1.414 1.414L11 7.414V15a1 1 0 11-2 0V7.414L6.707 9.707a1 1 0 01-1.414 0z"/>
                    </svg>
                    <svg class="h-3 w-3 transition-colors" 
                      :class="sortKey === 'status' && sortOrder === 'desc' ? 'text-blue-600' : 'text-gray-300'" 
                      fill="currentColor" viewBox="0 0 20 20">
                      <path d="M14.707 10.293a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 111.414-1.414L9 12.586V5a1 1 0 012 0v7.586l2.293-2.293a1 1 0 011.414 0z"/>
                    </svg>
                  </span>
                </div>
              </th>
              <th scope="col" @click="sortBy('destination_type')"
                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100 transition-colors select-none">
                <div class="flex items-center space-x-1">
                  <span>Destino</span>
                  <span class="flex flex-col -space-y-1">
                    <svg class="h-3 w-3 transition-colors" 
                      :class="sortKey === 'destination_type' && sortOrder === 'asc' ? 'text-blue-600' : 'text-gray-300'" 
                      fill="currentColor" viewBox="0 0 20 20">
                      <path d="M5.293 9.707a1 1 0 010-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 01-1.414 1.414L11 7.414V15a1 1 0 11-2 0V7.414L6.707 9.707a1 1 0 01-1.414 0z"/>
                    </svg>
                    <svg class="h-3 w-3 transition-colors" 
                      :class="sortKey === 'destination_type' && sortOrder === 'desc' ? 'text-blue-600' : 'text-gray-300'" 
                      fill="currentColor" viewBox="0 0 20 20">
                      <path d="M14.707 10.293a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 111.414-1.414L9 12.586V5a1 1 0 012 0v7.586l2.293-2.293a1 1 0 011.414 0z"/>
                    </svg>
                  </span>
                </div>
              </th>
              <th scope="col" @click="sortBy('user_rut')"
                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100 transition-colors select-none">
                <div class="flex items-center space-x-1">
                  <span>Usuario</span>
                  <span class="flex flex-col -space-y-1">
                    <svg class="h-3 w-3 transition-colors" 
                      :class="sortKey === 'user_rut' && sortOrder === 'asc' ? 'text-blue-600' : 'text-gray-300'" 
                      fill="currentColor" viewBox="0 0 20 20">
                      <path d="M5.293 9.707a1 1 0 010-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 01-1.414 1.414L11 7.414V15a1 1 0 11-2 0V7.414L6.707 9.707a1 1 0 01-1.414 0z"/>
                    </svg>
                    <svg class="h-3 w-3 transition-colors" 
                      :class="sortKey === 'user_rut' && sortOrder === 'desc' ? 'text-blue-600' : 'text-gray-300'" 
                      fill="currentColor" viewBox="0 0 20 20">
                      <path d="M14.707 10.293a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 111.414-1.414L9 12.586V5a1 1 0 012 0v7.586l2.293-2.293a1 1 0 011.414 0z"/>
                    </svg>
                  </span>
                </div>
              </th>
              <th scope="col" @click="sortBy('date_time')"
                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100 transition-colors select-none">
                <div class="flex items-center space-x-1">
                  <span>Fecha</span>
                  <span class="flex flex-col -space-y-1">
                    <svg class="h-3 w-3 transition-colors" 
                      :class="sortKey === 'date_time' && sortOrder === 'asc' ? 'text-blue-600' : 'text-gray-300'" 
                      fill="currentColor" viewBox="0 0 20 20">
                      <path d="M5.293 9.707a1 1 0 010-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 01-1.414 1.414L11 7.414V15a1 1 0 11-2 0V7.414L6.707 9.707a1 1 0 01-1.414 0z"/>
                    </svg>
                    <svg class="h-3 w-3 transition-colors" 
                      :class="sortKey === 'date_time' && sortOrder === 'desc' ? 'text-blue-600' : 'text-gray-300'" 
                      fill="currentColor" viewBox="0 0 20 20">
                      <path d="M14.707 10.293a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 111.414-1.414L9 12.586V5a1 1 0 012 0v7.586l2.293-2.293a1 1 0 011.414 0z"/>
                    </svg>
                  </span>
                </div>
              </th>
              <th scope="col" class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider sticky right-0 bg-gray-50 z-10 shadow-[-2px_0_4px_rgba(0,0,0,0.05)]">
                Detalles
              </th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-for="(item, index) in paginatedHistory" :key="item.id" class="hover:bg-gray-50 transition-colors duration-150">
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">#{{ startIndex + index + 1 }}</td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="text-sm">
                  <div class="font-medium text-gray-900">{{ item.supply_name || 'Sin nombre' }}</div>
                  <div class="text-gray-500 text-xs">ID: {{ item.medical_supply_id }}</div>
                </div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm font-mono text-gray-600">{{ item.qr_code || 'N/A' }}</td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span class="px-2 py-1 text-xs font-semibold rounded-full" :class="getStatusClass(item.status)">
                  {{ formatStatus(item.status) }}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="text-sm">
                  <div class="font-medium text-gray-900">{{ formatDestinationType(item.destination_type) }}</div>
                  <div class="text-gray-500">ID: {{ item.destination_id }}</div>
                </div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ item.user_rut || 'SYSTEM' }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ formatDateTime(item.date_time) }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium sticky right-0 bg-white z-10 shadow-[-2px_0_4px_rgba(0,0,0,0.05)]">
                <button @click="viewDetails(item)" 
                  class="text-blue-600 hover:text-blue-800 hover:bg-blue-50 p-1.5 rounded transition-colors"
                  title="Ver detalles">
                  <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Paginación -->
    <div v-if="!loading && sortedHistory.length > 0" class="card">
      <div class="flex flex-col sm:flex-row items-center justify-between gap-4">
        <div class="text-sm text-gray-700 text-center sm:text-left">
          Mostrando {{ startIndex + 1 }} a {{ endIndex }} de {{ sortedHistory.length }} registros
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

    <!-- Modal de detalles -->
    <Teleport to="body">
      <div v-if="showDetailsModal" class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50" @click.self="closeDetailsModal">
        <div class="relative top-20 mx-auto p-5 border w-full max-w-2xl shadow-lg rounded-md bg-white">
        <div class="space-y-4">
          <div class="flex justify-between items-center border-b pb-3">
            <h3 class="text-xl font-semibold text-gray-900">
              Detalles del Historial #{{ selectedItem?.id }}
            </h3>
            <button @click="closeDetailsModal" class="text-gray-400 hover:text-gray-600">
              <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
          
          <div v-if="selectedItem" class="space-y-4 max-h-96 overflow-y-auto">
            <!-- Información básica -->
            <div class="bg-gray-50 rounded-lg p-4">
              <h4 class="text-sm font-semibold text-gray-700 mb-3 flex items-center">
                <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                Información Básica
              </h4>
              <div class="grid grid-cols-2 gap-3">
                <div>
                  <p class="text-xs text-gray-500">ID Historial</p>
                  <p class="text-sm font-medium text-gray-900">{{ selectedItem.id }}</p>
                </div>
                <div>
                  <p class="text-xs text-gray-500">Estado</p>
                  <p class="text-sm">
                    <span class="px-2 py-1 text-xs font-semibold rounded-full" :class="getStatusClass(selectedItem.status)">
                      {{ formatStatus(selectedItem.status) }}
                    </span>
                  </p>
                </div>
                <div class="col-span-2">
                  <p class="text-xs text-gray-500">Nombre del Insumo</p>
                  <p class="text-sm font-medium text-gray-900">{{ selectedItem.supply_name || 'Sin nombre' }}</p>
                </div>
                <div>
                  <p class="text-xs text-gray-500">Código QR</p>
                  <p class="text-sm font-medium text-gray-900 font-mono">{{ selectedItem.qr_code || 'N/A' }}</p>
                </div>
                <div>
                  <p class="text-xs text-gray-500">ID Insumo</p>
                  <p class="text-sm font-medium text-gray-900">{{ selectedItem.medical_supply_id }}</p>
                </div>
                <div>
                  <p class="text-xs text-gray-500">Usuario</p>
                  <p class="text-sm font-medium text-gray-900">{{ selectedItem.user_rut || 'SYSTEM' }}</p>
                </div>
                <div>
                  <p class="text-xs text-gray-500">Fecha y Hora</p>
                  <p class="text-sm font-medium text-gray-900">{{ formatDateTime(selectedItem.date_time) }}</p>
                </div>
              </div>
            </div>

            <!-- Destino -->
            <div class="bg-blue-50 rounded-lg p-4">
              <h4 class="text-sm font-semibold text-gray-700 mb-3 flex items-center">
                <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
                </svg>
                Destino
              </h4>
              <div class="grid grid-cols-2 gap-3">
                <div>
                  <p class="text-xs text-gray-500">Tipo</p>
                  <p class="text-sm font-medium text-gray-900">{{ formatDestinationType(selectedItem.destination_type) }}</p>
                </div>
                <div>
                  <p class="text-xs text-gray-500">ID</p>
                  <p class="text-sm font-medium text-gray-900">{{ selectedItem.destination_id }}</p>
                </div>
              </div>
            </div>

            <!-- Origen (si existe) -->
            <div v-if="selectedItem.origin_type || selectedItem.origin_id" class="bg-purple-50 rounded-lg p-4">
              <h4 class="text-sm font-semibold text-gray-700 mb-3 flex items-center">
                <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4" />
                </svg>
                Origen
              </h4>
              <div class="grid grid-cols-2 gap-3">
                <div>
                  <p class="text-xs text-gray-500">Tipo</p>
                  <p class="text-sm font-medium text-gray-900">{{ selectedItem.origin_type ? formatDestinationType(selectedItem.origin_type) : 'N/A' }}</p>
                </div>
                <div>
                  <p class="text-xs text-gray-500">ID</p>
                  <p class="text-sm font-medium text-gray-900">{{ selectedItem.origin_id || 'N/A' }}</p>
                </div>
              </div>
            </div>

            <!-- Confirmación (si existe) -->
            <div v-if="selectedItem.confirmed_by || selectedItem.confirmation_date" class="bg-green-50 rounded-lg p-4">
              <h4 class="text-sm font-semibold text-gray-700 mb-3 flex items-center">
                <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                Confirmación
              </h4>
              <div class="grid grid-cols-2 gap-3">
                <div>
                  <p class="text-xs text-gray-500">Confirmado por</p>
                  <p class="text-sm font-medium text-gray-900">{{ selectedItem.confirmed_by || 'N/A' }}</p>
                </div>
                <div>
                  <p class="text-xs text-gray-500">Fecha de confirmación</p>
                  <p class="text-sm font-medium text-gray-900">{{ selectedItem.confirmation_date ? formatDateTime(selectedItem.confirmation_date) : 'N/A' }}</p>
                </div>
              </div>
            </div>

            <!-- Notas -->
            <div v-if="selectedItem.notes" class="bg-yellow-50 rounded-lg p-4">
              <h4 class="text-sm font-semibold text-gray-700 mb-2 flex items-center">
                <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                </svg>
                Notas
              </h4>
              <p class="text-sm text-gray-700 whitespace-pre-wrap">{{ selectedItem.notes }}</p>
            </div>

            <!-- Notas de transferencia -->
            <div v-if="selectedItem.transfer_notes" class="bg-orange-50 rounded-lg p-4">
              <h4 class="text-sm font-semibold text-gray-700 mb-2 flex items-center">
                <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                Notas de Transferencia
              </h4>
              <p class="text-sm text-gray-700 whitespace-pre-wrap">{{ selectedItem.transfer_notes }}</p>
            </div>
          </div>

          <div class="flex justify-end pt-4 border-t">
            <button @click="closeDetailsModal" class="btn-secondary">Cerrar</button>
          </div>
        </div>
      </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import supplyHistoryService from '@/services/inventory/supplyHistoryService'
import { exportToExcel as exportExcel, formatDateForExcel, formatStatusForExcel } from '@/utils/excelExport'
import { format } from 'date-fns'
import { es } from 'date-fns/locale'
import { useNotification } from '@/composables/useNotification'

const history = ref([])
const loading = ref(false)
const searchTerm = ref('')
const showDetailsModal = ref(false)
const selectedItem = ref(null)
const { success: showSuccess, error: showError } = useNotification()

const filters = ref({
  from_date: '',
  to_date: ''
})

// Estado de ordenamiento (por defecto ordenado por fecha más reciente primero)
const sortKey = ref('date_time')
const sortOrder = ref('desc')

// Estado de paginación
const currentPage = ref(1)
const itemsPerPage = 10

// Computed para verificar si hay filtros activos
// Helper para parsear fechas consistentemente
const parseDateTimeToLocal = (dateStr) => {
  if (!dateStr) return null
  if (typeof dateStr === 'string') {
    // Si no tiene información de zona horaria, asumir que es hora local
    if (!dateStr.includes('Z') && !dateStr.includes('+') && !dateStr.includes('-', 10)) {
      // Reemplazar espacio por 'T' para formato ISO compatible con hora local
      return new Date(dateStr.replace(' ', 'T'))
    }
  }
  return new Date(dateStr)
}

const hasActiveFilters = computed(() => {
  return searchTerm.value !== '' || filters.value.from_date !== '' || filters.value.to_date !== ''
})

const filteredHistory = computed(() => {
  let filtered = [...history.value]
  
  // Filtrar por búsqueda de texto
  if (searchTerm.value) {
    const term = searchTerm.value.toLowerCase()
    filtered = filtered.filter(item => 
      item.supply_name?.toLowerCase().includes(term) ||
      item.qr_code?.toLowerCase().includes(term) ||
      item.medical_supply_id?.toString().includes(term) ||
      item.status?.toLowerCase().includes(term) ||
      item.user_rut?.toLowerCase().includes(term) ||
      item.destination_type?.toLowerCase().includes(term)
    )
  }
  
  // Filtrar por rango de fechas
  if (filters.value.from_date) {
    // Crear fecha local a las 00:00:00
    const fromDate = new Date(filters.value.from_date + 'T00:00:00')
    filtered = filtered.filter(item => {
      if (!item.date_time) return false
      const itemDate = parseDateTimeToLocal(item.date_time)
      return itemDate >= fromDate
    })
  }
  
  if (filters.value.to_date) {
    // Crear fecha local a las 23:59:59.999
    const toDate = new Date(filters.value.to_date + 'T23:59:59.999')
    filtered = filtered.filter(item => {
      if (!item.date_time) return false
      const itemDate = parseDateTimeToLocal(item.date_time)
      return itemDate <= toDate
    })
  }
  
  return filtered
})

// Computed para obtener la lista ordenada
const sortedHistory = computed(() => {
  if (!filteredHistory.value || filteredHistory.value.length === 0) return []
  
  const sorted = [...filteredHistory.value].sort((a, b) => {
    let aVal = a[sortKey.value]
    let bVal = b[sortKey.value]
    
    // Manejo de valores null/undefined
    if (aVal == null) aVal = ''
    if (bVal == null) bVal = ''
    
    // Manejo de strings (comparación case-insensitive)
    if (typeof aVal === 'string' && sortKey.value !== 'date_time') {
      aVal = aVal.toLowerCase()
      bVal = typeof bVal === 'string' ? bVal.toLowerCase() : ''
    }
    
    // Manejo de fechas
    if (sortKey.value === 'date_time') {
      aVal = aVal ? new Date(aVal).getTime() : 0
      bVal = bVal ? new Date(bVal).getTime() : 0
    }
    
    // Comparación
    if (aVal < bVal) return sortOrder.value === 'asc' ? -1 : 1
    if (aVal > bVal) return sortOrder.value === 'asc' ? 1 : -1
    return 0
  })
  
  return sorted
})

// Computed properties para paginación
const totalPages = computed(() => Math.ceil(sortedHistory.value.length / itemsPerPage))

const startIndex = computed(() => (currentPage.value - 1) * itemsPerPage)
const endIndex = computed(() => Math.min(startIndex.value + itemsPerPage, sortedHistory.value.length))

const paginatedHistory = computed(() => {
  return sortedHistory.value.slice(startIndex.value, endIndex.value)
})

// Función para ordenar por columna
const sortBy = (key) => {
  if (sortKey.value === key) {
    // Si ya estamos ordenando por esta columna, cambiar dirección
    sortOrder.value = sortOrder.value === 'asc' ? 'desc' : 'asc'
  } else {
    // Nueva columna, ordenar ascendente por defecto
    sortKey.value = key
    sortOrder.value = 'asc'
  }
  currentPage.value = 1 // Resetear a la primera página al ordenar
}

const loadHistory = async () => {
  loading.value = true
  try {
    const data = await supplyHistoryService.getAllSupplyHistoryWithDetails()
    history.value = data
  } catch (err) {
    console.error('Error al cargar historial:', err)
    showError(err.message || 'Ocurrió un error al cargar el historial')
  } finally {
    loading.value = false
  }
}

const clearFilters = () => {
  searchTerm.value = ''
  filters.value.from_date = ''
  filters.value.to_date = ''
  currentPage.value = 1
}

const viewDetails = (item) => {
  selectedItem.value = item
  showDetailsModal.value = true
}

const closeDetailsModal = () => {
  showDetailsModal.value = false
  selectedItem.value = null
}

const formatDateTime = (date) => {
  if (!date) return 'N/A'
  try {
    // Crear objeto Date desde el string
    // Si la fecha viene sin zona horaria, asumir que es UTC y convertir a hora local
    let dateObj
    if (typeof date === 'string') {
      // Si no tiene información de zona horaria, asumir que es UTC
      if (!date.includes('Z') && !date.includes('+') && !date.includes('-', 10)) {
        // Formato: "2024-01-15 10:30:00" - asumir UTC
        dateObj = new Date(date + 'Z')
      } else {
        dateObj = new Date(date)
      }
    } else {
      dateObj = new Date(date)
    }
    
    // Verificar que la fecha sea válida
    if (isNaN(dateObj.getTime())) {
      console.warn('Fecha inválida:', date)
      return 'Fecha inválida'
    }
    
    // Usar date-fns para formatear consistentemente con el resto de la aplicación
    return format(dateObj, 'dd/MM/yyyy HH:mm:ss', { locale: es })
  } catch (error) {
    console.error('Error formateando fecha:', error, date)
    return 'Error en fecha'
  }
}

const formatStatus = (status) => {
  const statusMap = {
    'disponible': 'Disponible',
    'en_camino_a_pabellon': 'En Camino a Pabellón',
    'en_camino_a_bodega': 'En Camino a Bodega',
    'recepcionado': 'Recepcionado',
    'consumido': 'Consumido',
    'vencido': 'Vencido',
    'reservado': 'Reservado'
  }
  return statusMap[status] || status
}

const formatDestinationType = (type) => {
  const typeMap = {
    'store': 'Bodega',
    'pavilion': 'Pabellón'
  }
  return typeMap[type] || type
}

const getStatusClass = (status) => {
  const classes = {
    'disponible': 'bg-green-100 text-green-800',
    'en_camino_a_pabellon': 'bg-blue-100 text-blue-800',
    'en_camino_a_bodega': 'bg-yellow-100 text-yellow-800',
    'recepcionado': 'bg-purple-100 text-purple-800',
    'consumido': 'bg-gray-100 text-gray-800',
    'vencido': 'bg-red-100 text-red-800',
    'reservado': 'bg-orange-100 text-orange-800'
  }
  return classes[status?.toLowerCase()] || 'bg-gray-100 text-gray-800'
}

const exportToExcel = async () => {
  try {
    const columns = [
      { key: 'id', label: 'ID' },
      { key: 'supply_name', label: 'Nombre del Insumo' },
      { key: 'medical_supply_id', label: 'ID Insumo' },
      { key: 'qr_code', label: 'Código QR' },
      { key: 'status', label: 'Estado', formatter: (val) => formatStatusForExcel(val) },
      { key: 'destination_type', label: 'Tipo de Destino', formatter: (val) => formatDestinationType(val) },
      { key: 'destination_id', label: 'ID Destino' },
      { key: 'destination_name', label: 'Nombre Destino' },
      { key: 'origin_type', label: 'Tipo de Origen', formatter: (val) => val ? formatDestinationType(val) : '' },
      { key: 'origin_id', label: 'ID Origen' },
      { key: 'origin_name', label: 'Nombre Origen' },
      { key: 'user_rut', label: 'Usuario RUT' },
      { key: 'date_time', label: 'Fecha y Hora', formatter: formatDateForExcel },
      { key: 'notes', label: 'Notas' }
    ]
    
    await exportExcel(sortedHistory.value, columns, 'historial_insumos')
    showSuccess('El archivo Excel se ha descargado exitosamente')
  } catch (error) {
    console.error('Error al exportar:', error)
    showError('Ocurrió un error al exportar a Excel: ' + error.message)
  }
}

onMounted(() => {
  loadHistory()
})
</script>
