<template>
  <div class="space-y-6">
    <!-- Header del inventario -->
    <div>
      <h1 class="text-xl sm:text-2xl font-bold text-gray-900">Inventario de Insumos Médicos</h1>
      <p class="text-sm sm:text-base text-gray-600 mt-1">Gestión y control de stock médico</p>
    </div>

    <!-- Filtros y búsqueda -->
    <div class="card">
      <div class="flex flex-col sm:flex-row sm:items-end gap-4">
        <!-- Buscador único -->
        <div class="flex-1">
          <label class="block text-sm font-medium text-gray-700 mb-2">Buscar insumo</label>
          <div class="relative">
            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <svg class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
            </div>
            <input type="text" placeholder="Buscar por número de lote, nombre, código o proveedor..."
              class="form-input pl-10 w-full" v-model="searchTerm" />
          </div>
        </div>

        <!-- Botón de limpiar búsqueda -->
        <div class="w-full sm:w-auto">
          <button class="btn-secondary px-4 py-2 h-10 w-full sm:w-auto" @click="clearSearch" :disabled="!searchTerm">
            Limpiar
          </button>
        </div>
      </div>
    </div>

    <!-- Tabla de inventario -->
    <div class="card">
      <div class="card-header">
        <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
          <div>
            <h2 class="card-title">Inventario de Insumos Médicos</h2>
            <p class="text-sm text-gray-600">Total: {{ filteredSupplies.length }} lotes</p>
          </div>
          <div class="flex flex-col sm:flex-row gap-2 w-full sm:w-auto">
            <button 
              @click="exportToExcel" 
              :disabled="loading || filteredSupplies.length === 0"
              class="btn-secondary flex items-center justify-center px-3 py-2 text-sm"
              :class="{ 'opacity-50 cursor-not-allowed': loading || filteredSupplies.length === 0 }"
            >
              <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
              <span class="hidden sm:inline">Exportar a Excel</span>
              <span class="sm:hidden">Excel</span>
            </button>
            <button v-if="authStore.canViewAllRequests"
              class="btn-primary flex items-center justify-center px-3 py-2 text-sm w-full sm:w-auto"
              @click="openGlobalHistoryModal">
              <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <span class="hidden sm:inline">Historial de Movimientos</span>
              <span class="sm:hidden">Historial</span>
            </button>
          </div>
        </div>
      </div>

      <!-- Indicador de carga -->
      <div v-if="loading" class="flex justify-center items-center py-8">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
        <span class="ml-2 text-gray-600">Cargando inventario...</span>
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
            <h3 class="text-sm font-medium text-red-800">Error al cargar inventario</h3>
            <div class="mt-2 text-sm text-red-700">{{ error }}</div>
            <div class="mt-4">
              <button @click="loadInventory" class="btn-secondary text-sm">
                Reintentar
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Tabla de datos -->
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
              <th
                class="px-3 sm:px-6 py-2 sm:py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider whitespace-nowrap min-w-[100px] sm:min-w-[120px]">
                <div class="flex items-center justify-between">
                  <span>N° de lote</span>
                  <div class="flex flex-col ml-2">
                    <button @click="sortBy('batch_id', 'asc')" class="text-gray-400 hover:text-gray-600 p-1"
                      :class="{ 'text-blue-600': sortField === 'batch_id' && sortDirection === 'asc' }">
                      <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7" />
                      </svg>
                    </button>
                    <button @click="sortBy('batch_id', 'desc')" class="text-gray-400 hover:text-gray-600 p-1"
                      :class="{ 'text-blue-600': sortField === 'batch_id' && sortDirection === 'desc' }">
                      <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                      </svg>
                    </button>
                  </div>
                </div>
              </th>
              <th class="table-header-cell">
                <div class="flex items-center justify-between">
                  <span>Nombre del Insumo</span>
                  <div class="flex flex-col ml-2">
                    <button @click="sortBy('name', 'asc')" class="text-gray-400 hover:text-gray-600 p-1"
                      :class="{ 'text-blue-600': sortField === 'name' && sortDirection === 'asc' }">
                      <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7" />
                      </svg>
                    </button>
                    <button @click="sortBy('name', 'desc')" class="text-gray-400 hover:text-gray-600 p-1"
                      :class="{ 'text-blue-600': sortField === 'name' && sortDirection === 'desc' }">
                      <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                      </svg>
                    </button>
                  </div>
                </div>
              </th>
              <th class="table-header-cell">
                <div class="flex items-center justify-between">
                  <span>Código Interno</span>
                  <div class="flex flex-col ml-2">
                    <button @click="sortBy('code', 'asc')" class="text-gray-400 hover:text-gray-600 p-1"
                      :class="{ 'text-blue-600': sortField === 'code' && sortDirection === 'asc' }">
                      <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7" />
                      </svg>
                    </button>
                    <button @click="sortBy('code', 'desc')" class="text-gray-400 hover:text-gray-600 p-1"
                      :class="{ 'text-blue-600': sortField === 'code' && sortDirection === 'desc' }">
                      <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                      </svg>
                    </button>
                  </div>
                </div>
              </th>
              <th class="table-header-cell">
                <div class="flex items-center justify-between">
                  <span>F. Vencimiento</span>
                  <div class="flex flex-col ml-2">
                    <button @click="sortBy('expiration_date', 'asc')" class="text-gray-400 hover:text-gray-600 p-1"
                      :class="{ 'text-blue-600': sortField === 'expiration_date' && sortDirection === 'asc' }">
                      <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7" />
                      </svg>
                    </button>
                    <button @click="sortBy('expiration_date', 'desc')" class="text-gray-400 hover:text-gray-600 p-1"
                      :class="{ 'text-blue-600': sortField === 'expiration_date' && sortDirection === 'desc' }">
                      <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                      </svg>
                    </button>
                  </div>
                </div>
              </th>
              <th class="table-header-cell">
                <div class="flex items-center justify-between">
                  <span>Cantidad</span>
                  <div class="flex flex-col ml-2">
                    <button @click="sortBy('amount', 'asc')" class="text-gray-400 hover:text-gray-600 p-1"
                      :class="{ 'text-blue-600': sortField === 'amount' && sortDirection === 'asc' }">
                      <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7" />
                      </svg>
                    </button>
                    <button @click="sortBy('amount', 'desc')" class="text-gray-400 hover:text-gray-600 p-1"
                      :class="{ 'text-blue-600': sortField === 'amount' && sortDirection === 'desc' }">
                      <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                      </svg>
                    </button>
                  </div>
                </div>
              </th>
              <th class="table-header-cell">
                <div class="flex items-center justify-between">
                  <span>Proveedor</span>
                  <div class="flex flex-col ml-2">
                    <button @click="sortBy('supplier', 'asc')" class="text-gray-400 hover:text-gray-600 p-1"
                      :class="{ 'text-blue-600': sortField === 'supplier' && sortDirection === 'asc' }">
                      <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7" />
                      </svg>
                    </button>
                    <button @click="sortBy('supplier', 'desc')" class="text-gray-400 hover:text-gray-600 p-1"
                      :class="{ 'text-blue-600': sortField === 'supplier' && sortDirection === 'desc' }">
                      <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                      </svg>
                    </button>
                  </div>
                </div>
              </th>
              <th class="table-header-cell">Acciones</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-for="supply in paginatedSupplies" :key="supply.batch_id" class="hover:bg-gray-50 cursor-pointer"
              @click="openBatchDetailsModal(supply)">
              <td class="px-3 sm:px-6 py-3 sm:py-4 whitespace-nowrap text-xs sm:text-sm font-mono text-gray-900">{{
                supply.batch_id }}</td>
              <td class="px-3 sm:px-6 py-3 sm:py-4 whitespace-nowrap">
                <div>
                  <div class="font-medium text-gray-900 text-xs sm:text-sm">{{ supply.name }}</div>
                </div>
              </td>
              <td class="px-3 sm:px-6 py-3 sm:py-4 whitespace-nowrap">
                <span class="text-gray-700 text-xs sm:text-sm">{{ supply.code }}</span>
              </td>
              <td class="px-3 sm:px-6 py-3 sm:py-4 whitespace-nowrap">
                <span :class="getExpirationClass(supply.expiration_date)" class="text-xs sm:text-sm">
                  {{ formatDate(supply.expiration_date) }}
                </span>
              </td>
              <td class="px-3 sm:px-6 py-3 sm:py-4 whitespace-nowrap">
                <span :class="getAmountClass(supply.amount)" class="font-medium text-xs sm:text-sm">{{ supply.amount
                  }}</span>
                <span :class="getAmountClass(supply.amount)" class="text-xs ml-1">unidades</span>
              </td>
              <td class="px-3 sm:px-6 py-3 sm:py-4 whitespace-nowrap">
                <span class="text-gray-700 text-xs sm:text-sm">{{ supply.supplier }}</span>
              </td>
              <td class="px-3 sm:px-6 py-3 sm:py-4 whitespace-nowrap" @click.stop>
                <div class="flex space-x-1.5 sm:space-x-2">
                  <button class="text-warning-600 hover:text-warning-800" @click.stop="editSupply(supply)">
                    <svg class="h-4 w-4 sm:h-5 sm:w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                    </svg>
                  </button>
                  <button class="text-danger-600 hover:text-danger-800" @click.stop="deleteSupply(supply)">
                    <svg class="h-4 w-4 sm:h-5 sm:w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                    </svg>
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Modal de edición -->
      <!-- Modal de detalles de lote con QR -->
      <div v-if="showBatchDetailsModal"
        class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50">
        <div
          class="relative top-4 sm:top-20 mx-2 sm:mx-auto p-4 sm:p-8 border w-auto sm:w-full max-w-2xl shadow-lg rounded-xl bg-white">
          <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 sm:gap-0 mb-4 sm:mb-6">
            <h3 class="text-lg sm:text-2xl font-bold text-gray-900">Detalles del lote de {{ supplyName || '...' }}</h3>
            <button @click="closeBatchDetailsModal"
              class="btn-secondary px-4 py-2 rounded-lg text-sm sm:text-base w-full sm:w-auto">Cerrar</button>
          </div>
          <div v-if="batchDetailsLoading" class="flex justify-center items-center py-8">
            <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
            <span class="ml-2 text-gray-600">Cargando insumos del lote...</span>
          </div>
          <div v-else-if="batchDetailsError" class="bg-red-50 border border-red-200 rounded-md p-4 mb-4">
            <div class="text-red-700">{{ batchDetailsError }}</div>
          </div>
          <div v-else-if="batchDetails.length === 0" class="text-center py-8">
            No se encontraron insumos para este lote.
          </div>
          <div v-else>
            <!-- Información de paginación -->
            <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-2 sm:gap-3 mb-4">
              <div class="text-xs sm:text-sm text-gray-700 text-center sm:text-left">
                Mostrando {{ batchDetailsStartIndex + 1 }} a {{ batchDetailsEndIndex }} de {{ batchDetails.length }}
                insumos
              </div>
              <div class="text-xs sm:text-sm text-gray-600 text-center sm:text-left">
                Página {{ batchDetailsCurrentPage }} de {{ batchDetailsTotalPages }}
              </div>
            </div>

            <div class="overflow-x-auto -mx-4 sm:mx-0">
              <table class="min-w-full divide-y divide-gray-200 mb-6 rounded-lg overflow-hidden shadow">
                <thead class="bg-gray-100">
                  <tr>
                    <th class="px-3 sm:px-6 py-2 sm:py-3 text-left text-xs sm:text-sm font-semibold text-gray-700">ID
                    </th>
                    <th class="px-3 sm:px-6 py-2 sm:py-3 text-left text-xs sm:text-sm font-semibold text-gray-700">
                      Código</th>
                    <th class="px-3 sm:px-6 py-2 sm:py-3 text-left text-xs sm:text-sm font-semibold text-gray-700">
                      Estado</th>
                    <th class="px-3 sm:px-6 py-2 sm:py-3 text-center text-xs sm:text-sm font-semibold text-gray-700">QR
                    </th>
                  </tr>
                </thead>
                <tbody>
                  <tr v-for="item in paginatedBatchDetails" :key="item.id" class="hover:bg-gray-50">
                    <td class="px-3 sm:px-6 py-3 sm:py-4 text-gray-900 font-mono text-xs sm:text-sm">{{ item.id }}</td>
                    <td class="px-3 sm:px-6 py-3 sm:py-4 text-gray-900 text-xs sm:text-sm">{{ item.code }}</td>
                    <td class="px-3 sm:px-6 py-3 sm:py-4">
                      <span
                        class="inline-flex items-center px-2 sm:px-2.5 py-0.5 rounded-full text-xs font-medium whitespace-nowrap"
                        :class="getStatusClass(item.status)">
                        {{ getStatusDescription(item.status) }}
                      </span>
                    </td>
                    <td class="px-3 sm:px-6 py-3 sm:py-4">
                      <div class="flex flex-col items-center justify-center text-center space-y-2">
                        <qrcode-vue :value="item.qr_code" :size="40" class="mx-auto" />
                        <button @click="handleDownloadQR(item.qr_code)"
                          class="btn-secondary text-xs px-2 sm:px-3 py-1 rounded whitespace-nowrap">Descargar
                          QR</button>
                      </div>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>

            <!-- Paginación -->
            <div v-if="batchDetails.length > batchDetailsItemsPerPage"
              class="flex flex-col sm:flex-row items-center justify-between mt-4 gap-3">
              <div class="text-xs sm:text-sm text-gray-700 text-center sm:text-left">
                Mostrando {{ batchDetailsStartIndex + 1 }} a {{ batchDetailsEndIndex }} de {{ batchDetails.length }}
                insumos
              </div>
              <div class="flex items-center gap-2">
                <button class="btn-secondary px-3 py-2 text-sm min-w-[80px]" :disabled="batchDetailsCurrentPage === 1"
                  @click="batchDetailsCurrentPage--">
                  <span class="hidden sm:inline">Anterior</span>
                  <span class="sm:hidden">Ant.</span>
                </button>
                <span class="px-3 py-2 text-sm text-gray-700 bg-gray-100 rounded-md min-w-[100px] text-center">
                  Página {{ batchDetailsCurrentPage }} de {{ batchDetailsTotalPages }}
                </span>
                <button class="btn-secondary px-3 py-2 text-sm min-w-[80px]"
                  :disabled="batchDetailsCurrentPage === batchDetailsTotalPages" @click="batchDetailsCurrentPage++">
                  <span class="hidden sm:inline">Siguiente</span>
                  <span class="sm:hidden">Sig.</span>
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
      <div v-if="showEditModal" class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50 px-4">
        <div
          class="relative top-4 sm:top-20 mx-auto p-4 sm:p-5 border w-full max-w-md shadow-lg rounded-md bg-white my-8">
          <div class="mt-2 sm:mt-3">
            <div class="flex items-center justify-between mb-4">
              <h3 class="text-base sm:text-lg font-medium text-gray-900">Editar Lote de Insumo</h3>
              <button @click="closeEditModal" class="text-gray-400 hover:text-gray-600">
                <svg class="h-5 w-5 sm:h-6 sm:w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>

            <form @submit.prevent="saveEdit" class="space-y-3 sm:space-y-4">
              <div>
                <label class="block text-xs sm:text-sm font-medium text-gray-700 mb-1">N° de lote</label>
                <input type="text" v-model="editingSupply.batch_id"
                  class="form-input w-full text-sm bg-gray-100 text-gray-600 cursor-not-allowed" readonly disabled />
                <p class="text-xs text-gray-500 mt-1">El número de lote no se puede editar</p>
              </div>

              <div>
                <label class="block text-xs sm:text-sm font-medium text-gray-700 mb-1">Nombre del Insumo</label>
                <input type="text" v-model="editingSupply.name"
                  class="form-input w-full text-sm bg-gray-100 text-gray-600 cursor-not-allowed" readonly disabled />
                <p class="text-xs text-gray-500 mt-1">El nombre no se puede editar</p>
              </div>

              <div>
                <label class="block text-xs sm:text-sm font-medium text-gray-700 mb-1">Código Interno</label>
                <input type="text" v-model="editingSupply.code"
                  class="form-input w-full text-sm bg-gray-100 text-gray-600 cursor-not-allowed" readonly disabled />
                <p class="text-xs text-gray-500 mt-1">El código interno no se puede editar</p>
              </div>

              <div>
                <label class="block text-xs sm:text-sm font-medium text-gray-700 mb-1">Fecha de Vencimiento</label>
                <input type="date" v-model="editingSupply.expiration_date" class="form-input w-full text-sm" required />
              </div>

              <div>
                <label class="block text-xs sm:text-sm font-medium text-gray-700 mb-1">Cantidad</label>
                <input type="number" v-model="editingSupply.amount" class="form-input w-full text-sm" min="0"
                  required />
                <span class="text-xs text-gray-500">unidades</span>
              </div>

              <div>
                <label class="block text-xs sm:text-sm font-medium text-gray-700 mb-1">Proveedor</label>
                <input type="text" v-model="editingSupply.supplier" class="form-input w-full text-sm" required />
              </div>

              <div class="flex flex-col sm:flex-row sm:space-x-3 space-y-2 sm:space-y-0 pt-3 sm:pt-4">
                <button type="button" @click="closeEditModal"
                  class="btn-secondary flex-1 text-sm sm:text-base py-2.5 sm:py-2">
                  Cancelar
                </button>
                <button type="submit" class="btn-primary flex-1 text-sm sm:text-base py-2.5 sm:py-2" :disabled="saving">
                  <span v-if="saving" class="flex items-center justify-center">
                    <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
                      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                      <path class="opacity-75" fill="currentColor"
                        d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z">
                      </path>
                    </svg>
                    Guardando...
                  </span>
                  <span v-else>Guardar Cambios</span>
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>

      <!-- Modal de historial -->
      <div v-if="showHistoryModal"
        class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50 px-2 sm:px-0">
        <div
          class="relative top-2 sm:top-4 md:top-20 mx-auto p-3 sm:p-5 border w-full max-w-4xl shadow-lg rounded-md bg-white mb-4">
          <div class="mt-2 sm:mt-3">
            <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-2 sm:gap-4 mb-3 sm:mb-4">
              <div class="flex-1 min-w-0">
                <h3 class="text-base sm:text-lg md:text-xl font-medium text-gray-900 truncate">Historial de Movimientos
                </h3>
                <p class="text-xs sm:text-sm text-gray-600 truncate mt-1">{{ selectedSupplyForHistory?.name }}</p>
              </div>
              <button @click="closeHistoryModal"
                class="text-gray-400 hover:text-gray-600 self-end sm:self-auto p-1 flex-shrink-0">
                <svg class="h-5 w-5 sm:h-6 sm:w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>

            <!-- Filtros de historial -->
            <div class="flex flex-col sm:flex-row sm:items-end gap-2 sm:gap-4 mb-3 sm:mb-4">
              <div class="flex-1">
                <label class="block text-xs sm:text-sm font-medium text-gray-700 mb-1 sm:mb-2">Buscar movimiento</label>
                <div class="relative">
                  <div class="absolute inset-y-0 left-0 pl-2 sm:pl-3 flex items-center pointer-events-none">
                    <svg class="h-3.5 w-3.5 sm:h-4 sm:w-4 text-gray-400" fill="none" stroke="currentColor"
                      viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                    </svg>
                  </div>
                  <input type="text" placeholder="Buscar..."
                    class="form-input pl-8 sm:pl-9 w-full text-xs sm:text-sm min-h-[40px]"
                    v-model="historySearchTerm" />
                </div>
              </div>
              <div class="w-full sm:w-auto">
                <button class="btn-secondary px-3 sm:px-4 py-2 h-10 text-xs sm:text-sm w-full sm:w-auto"
                  @click="clearHistorySearch" :disabled="!historySearchTerm">
                  Limpiar
                </button>
              </div>
            </div>

            <!-- Indicador de carga -->
            <div v-if="historyLoading" class="flex justify-center items-center py-6 sm:py-8">
              <div class="animate-spin rounded-full h-6 w-6 sm:h-8 sm:w-8 border-b-2 border-blue-600"></div>
              <span class="ml-2 text-gray-600 text-xs sm:text-sm">Cargando...</span>
            </div>

            <!-- Mensaje de error -->
            <div v-else-if="historyError" class="bg-red-50 border border-red-200 rounded-md p-3 sm:p-4 mb-3 sm:mb-4">
              <div class="flex">
                <div class="flex-shrink-0">
                  <svg class="h-4 w-4 sm:h-5 sm:w-5 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                </div>
                <div class="ml-2 sm:ml-3 flex-1">
                  <h3 class="text-xs sm:text-sm font-medium text-red-800">Error al cargar historial</h3>
                  <div class="mt-1 sm:mt-2 text-xs sm:text-sm text-red-700">{{ historyError }}</div>
                  <div class="mt-3 sm:mt-4">
                    <button @click="loadHistory" class="btn-secondary text-xs sm:text-sm px-3 py-1.5">
                      Reintentar
                    </button>
                  </div>
                </div>
              </div>
            </div>

            <!-- Tabla de datos -->
            <div v-else-if="!historyData.length" class="text-center py-6 sm:py-8">
              <svg class="mx-auto h-10 w-10 sm:h-12 sm:w-12 text-gray-400" fill="none" stroke="currentColor"
                viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <h3 class="mt-2 text-xs sm:text-sm font-medium text-gray-900">No hay historial disponible</h3>
              <p class="mt-1 text-xs sm:text-sm text-gray-500">Este lote no tiene movimientos registrados.</p>
            </div>

            <!-- Vista de tarjetas para móviles y tabla para desktop -->
            <div v-else>
              <!-- Vista de tarjetas para móviles -->
              <div class="md:hidden">
                <div class="space-y-3">
                  <div v-for="movement in paginatedHistory" :key="movement.id || movement.date_time"
                    class="bg-white border border-gray-200 rounded-lg p-3 shadow-sm">
                    <!-- Fecha y tipo -->
                    <div class="flex items-start justify-between mb-2 pb-2 border-b border-gray-100">
                      <div class="flex-1 min-w-0">
                        <div class="flex items-center gap-2 mb-1">
                          <svg class="h-3.5 w-3.5 text-gray-400 flex-shrink-0" fill="none" stroke="currentColor"
                            viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                              d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                          </svg>
                          <span class="text-xs font-medium text-gray-900">{{ formatDate(movement.date_time ||
                            movement.date) }}</span>
                        </div>
                        <span :class="getStatusBadgeClass(movement.type || movement.action)"
                          class="inline-block text-xs px-2 py-0.5 rounded-full">
                          {{ movement.action || movement.type || 'N/A' }}
                        </span>
                      </div>
                    </div>

                    <!-- Detalles -->
                    <div class="space-y-1.5">
                      <!-- Cantidad -->
                      <div class="flex items-center justify-between text-xs">
                        <span class="text-gray-500 font-medium">Cantidad:</span>
                        <span :class="getAmountClass(movement.details?.amount || movement.amount || 0)"
                          class="font-semibold">
                          {{ movement.details?.amount || movement.amount || 'N/A' }} unidades
                        </span>
                      </div>

                      <!-- Usuario -->
                      <div class="flex items-center justify-between text-xs">
                        <span class="text-gray-500 font-medium">Usuario:</span>
                        <span class="text-gray-900">{{ movement.user_rut || movement.user || 'N/A' }}</span>
                      </div>
                    </div>

                    <!-- Botón de ver detalles -->
                    <div class="mt-3 pt-2 border-t border-gray-100">
                      <button @click="viewMovementDetails(movement)"
                        class="w-full text-center text-xs font-medium text-blue-600 hover:text-blue-800 py-1.5">
                        Ver detalles completos
                      </button>
                    </div>
                  </div>
                </div>

              </div>

              <!-- Vista de tabla para desktop -->
              <div class="hidden md:block table-container">
                <table class="min-w-full divide-y divide-gray-200" style="min-width: 700px;">
                  <thead class="table-header">
                    <tr>
                      <th
                        class="px-3 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider whitespace-nowrap min-w-[100px] sm:min-w-[120px]">
                        <div class="flex items-center justify-between">
                          <span class="text-xs sm:text-sm">Fecha</span>
                          <div class="flex flex-col ml-1 sm:ml-2">
                            <button @click="sortHistoryBy('date', 'asc')" class="text-gray-400 hover:text-gray-600 p-1"
                              :class="{ 'text-blue-600': sortHistoryField === 'date' && sortHistoryDirection === 'asc' }">
                              <svg class="w-2 h-2 sm:w-3 sm:h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                  d="M5 15l7-7 7 7" />
                              </svg>
                            </button>
                            <button @click="sortHistoryBy('date', 'desc')" class="text-gray-400 hover:text-gray-600 p-1"
                              :class="{ 'text-blue-600': sortHistoryField === 'date' && sortHistoryDirection === 'desc' }">
                              <svg class="w-2 h-2 sm:w-3 sm:h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                  d="M19 9l-7 7-7-7" />
                              </svg>
                            </button>
                          </div>
                        </div>
                      </th>
                      <th
                        class="px-3 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider whitespace-nowrap min-w-[120px] sm:min-w-[150px]">
                        <div class="flex items-center justify-between">
                          <span class="text-xs sm:text-sm">Tipo</span>
                          <div class="flex flex-col ml-1 sm:ml-2">
                            <button @click="sortHistoryBy('type', 'asc')" class="text-gray-400 hover:text-gray-600 p-1"
                              :class="{ 'text-blue-600': sortHistoryField === 'type' && sortHistoryDirection === 'asc' }">
                              <svg class="w-2 h-2 sm:w-3 sm:h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                  d="M5 15l7-7 7 7" />
                              </svg>
                            </button>
                            <button @click="sortHistoryBy('type', 'desc')" class="text-gray-400 hover:text-gray-600 p-1"
                              :class="{ 'text-blue-600': sortHistoryField === 'type' && sortHistoryDirection === 'desc' }">
                              <svg class="w-2 h-2 sm:w-3 sm:h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                  d="M19 9l-7 7-7-7" />
                              </svg>
                            </button>
                          </div>
                        </div>
                      </th>
                      <th
                        class="px-3 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider whitespace-nowrap min-w-[80px] sm:min-w-[100px]">
                        <div class="flex items-center justify-between">
                          <span class="text-xs sm:text-sm">Cantidad</span>
                          <div class="flex flex-col ml-1 sm:ml-2">
                            <button @click="sortHistoryBy('amount', 'asc')"
                              class="text-gray-400 hover:text-gray-600 p-1"
                              :class="{ 'text-blue-600': sortHistoryField === 'amount' && sortHistoryDirection === 'asc' }">
                              <svg class="w-2 h-2 sm:w-3 sm:h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                  d="M5 15l7-7 7 7" />
                              </svg>
                            </button>
                            <button @click="sortHistoryBy('amount', 'desc')"
                              class="text-gray-400 hover:text-gray-600 p-1"
                              :class="{ 'text-blue-600': sortHistoryField === 'amount' && sortHistoryDirection === 'desc' }">
                              <svg class="w-2 h-2 sm:w-3 sm:h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                  d="M19 9l-7 7-7-7" />
                              </svg>
                            </button>
                          </div>
                        </div>
                      </th>
                      <th
                        class="px-3 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider whitespace-nowrap min-w-[100px] sm:min-w-[120px]">
                        <div class="flex items-center justify-between">
                          <span class="text-xs sm:text-sm">Usuario</span>
                          <div class="flex flex-col ml-1 sm:ml-2">
                            <button @click="sortHistoryBy('user', 'asc')" class="text-gray-400 hover:text-gray-600 p-1"
                              :class="{ 'text-blue-600': sortHistoryField === 'user' && sortHistoryDirection === 'asc' }">
                              <svg class="w-2 h-2 sm:w-3 sm:h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                  d="M5 15l7-7 7 7" />
                              </svg>
                            </button>
                            <button @click="sortHistoryBy('user', 'desc')" class="text-gray-400 hover:text-gray-600 p-1"
                              :class="{ 'text-blue-600': sortHistoryField === 'user' && sortHistoryDirection === 'desc' }">
                              <svg class="w-2 h-2 sm:w-3 sm:h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                  d="M19 9l-7 7-7-7" />
                              </svg>
                            </button>
                          </div>
                        </div>
                      </th>
                      <th
                        class="px-3 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider whitespace-nowrap min-w-[80px] sm:min-w-[100px]">
                        Acciones</th>
                    </tr>
                  </thead>
                  <tbody class="bg-white divide-y divide-gray-200">
                    <tr v-for="movement in paginatedHistory" :key="movement.id || movement.date_time"
                      class="hover:bg-gray-50">
                      <td
                        class="px-3 sm:px-6 py-3 sm:py-4 whitespace-nowrap text-xs sm:text-sm font-mono text-gray-900">
                        {{ formatDate(movement.date_time || movement.date) }}
                      </td>
                      <td class="px-3 sm:px-6 py-3 sm:py-4 whitespace-nowrap">
                        <span :class="getStatusBadgeClass(movement.type || movement.action)" class="text-xs sm:text-sm">
                          {{ movement.action || movement.type || 'N/A' }}
                        </span>
                      </td>
                      <td class="px-3 sm:px-6 py-3 sm:py-4 whitespace-nowrap">
                        <span :class="getAmountClass(movement.details?.amount || movement.amount || 0)"
                          class="font-medium text-xs sm:text-sm">{{ movement.details?.amount || movement.amount || 'N/A'
                          }}</span>
                        <span v-if="movement.details?.amount || movement.amount"
                          :class="getAmountClass(movement.details?.amount || movement.amount)"
                          class="text-xs ml-1">unidades</span>
                      </td>
                      <td class="px-3 sm:px-6 py-3 sm:py-4 whitespace-nowrap">
                        <span class="text-gray-700 text-xs sm:text-sm">{{ movement.user_rut || movement.user || 'N/A'
                          }}</span>
                      </td>
                      <td class="px-3 sm:px-6 py-3 sm:py-4 whitespace-nowrap">
                        <div class="flex space-x-2">
                          <button class="text-blue-600 hover:text-blue-800"
                            @click="viewMovementDetails(movement)">
                            <svg class="h-4 w-4 sm:h-5 sm:w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                            </svg>
                          </button>
                        </div>
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>

            <!-- Paginación -->
            <div class="flex flex-col sm:flex-row items-center justify-between mt-4 sm:mt-6 gap-3 sm:gap-4">
              <div class="text-xs sm:text-sm text-gray-700 text-center sm:text-left">
                Mostrando {{ filteredHistory.length > 0 ? historyStartIndex + 1 : 0 }} a {{ historyEndIndex }} de {{
                filteredHistory.length }} resultados
              </div>
              <div class="flex items-center gap-2 w-full sm:w-auto">
                <button class="btn-secondary px-3 sm:px-4 py-2 text-xs sm:text-sm flex-1 sm:flex-none sm:min-w-[80px]"
                  :disabled="historyCurrentPage === 1" @click="historyCurrentPage--">
                  <span class="hidden sm:inline">Anterior</span>
                  <span class="sm:hidden">Ant.</span>
                </button>
                <span
                  class="px-3 py-2 text-xs sm:text-sm text-gray-700 bg-gray-100 rounded-md text-center whitespace-nowrap">
                  <span class="hidden sm:inline">Página {{ historyCurrentPage }} de {{ historyTotalPages }}</span>
                  <span class="sm:hidden">{{ historyCurrentPage }} / {{ historyTotalPages }}</span>
                </span>
                <button class="btn-secondary px-3 sm:px-4 py-2 text-xs sm:text-sm flex-1 sm:flex-none sm:min-w-[80px]"
                  :disabled="historyCurrentPage === historyTotalPages" @click="historyCurrentPage++">
                  <span class="hidden sm:inline">Siguiente</span>
                  <span class="sm:hidden">Sig.</span>
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Modal de historial global -->
      <div v-if="showGlobalHistoryModal"
        class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50 px-2 sm:px-0">
        <div
          class="relative top-2 sm:top-4 md:top-20 mx-auto p-3 sm:p-5 border w-full max-w-4xl shadow-lg rounded-md bg-white mb-4">
          <div class="mt-2 sm:mt-3">
            <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-2 sm:gap-4 mb-3 sm:mb-4">
              <div class="flex-1 min-w-0">
                <h3 class="text-base sm:text-lg md:text-xl font-medium text-gray-900">Historial de Lotes</h3>
                <p v-if="hasActiveFilters" class="text-xs sm:text-sm text-gray-600 mt-1">
                  {{ filteredGlobalHistory.length }} resultado{{ filteredGlobalHistory.length !== 1 ? 's' : '' }}
                </p>
              </div>
              <button @click="closeGlobalHistoryModal"
                class="text-gray-400 hover:text-gray-600 self-end sm:self-auto p-1 flex-shrink-0">
                <svg class="h-5 w-5 sm:h-6 sm:w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>

            <!-- Filtros de historial global -->
            <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-2 sm:gap-3 lg:gap-4 mb-3 sm:mb-4 items-end">
              <!-- Filtro por tipo de cambio -->
              <div>
                <label class="block text-xs sm:text-sm font-medium text-gray-700 mb-1 sm:mb-2">Tipo de Cambio</label>
                <select v-model="globalHistoryChangeTypeFilter"
                  class="form-select w-full text-xs sm:text-sm min-h-[40px]" @change="applyFilters">
                  <option value="">Todos los tipos</option>
                  <option value="cantidad">Cantidad</option>
                  <option value="proveedor">Proveedor</option>
                  <option value="fecha de expiración">Fecha de Expiración</option>
                  <option value="creado">Creado</option>
                  <option value="eliminado">Eliminado</option>
                </select>
              </div>

              <!-- Buscador por número de lote -->
              <div>
                <label class="block text-xs sm:text-sm font-medium text-gray-700 mb-1 sm:mb-2">N° de Lote</label>
                <div class="relative">
                  <div class="absolute inset-y-0 left-0 pl-2 sm:pl-3 flex items-center pointer-events-none">
                    <svg class="h-3.5 w-3.5 sm:h-4 sm:w-4 text-gray-400" fill="none" stroke="currentColor"
                      viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                    </svg>
                  </div>
                  <input type="text" placeholder="N° de lote..."
                    class="form-input pl-8 sm:pl-9 w-full text-xs sm:text-sm min-h-[40px]"
                    v-model="globalHistoryBatchFilter" @input="applyFilters" />
                </div>
              </div>

              <!-- Buscador por usuario -->
              <div>
                <label class="block text-xs sm:text-sm font-medium text-gray-700 mb-1 sm:mb-2">Usuario</label>
                <div class="relative">
                  <div class="absolute inset-y-0 left-0 pl-2 sm:pl-3 flex items-center pointer-events-none">
                    <svg class="h-3.5 w-3.5 sm:h-4 sm:w-4 text-gray-400" fill="none" stroke="currentColor"
                      viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                        d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                    </svg>
                  </div>
                  <input type="text" placeholder="RUT o nombre..."
                    class="form-input pl-8 sm:pl-9 w-full text-xs sm:text-sm min-h-[40px]"
                    v-model="globalHistoryUserFilter" @input="applyFilters" />
                </div>
              </div>

              <!-- Botón de limpiar filtros -->
              <div class="flex justify-end items-end h-full">
                <button
                  class="btn-secondary flex items-center justify-center px-3 sm:px-4 py-2 h-10 w-full sm:w-auto text-xs sm:text-sm"
                  @click="clearGlobalHistoryFilters" :disabled="!hasActiveFilters">
                  <svg v-if="hasActiveFilters" class="h-3.5 w-3.5 sm:h-4 sm:w-4 mr-1.5 sm:mr-2 flex-shrink-0"
                    fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                  </svg>
                  <span>Limpiar</span>
                </button>
              </div>
            </div>

            <!-- Indicador de carga -->
            <div v-if="globalHistoryLoading" class="flex justify-center items-center py-6 sm:py-8">
              <div class="animate-spin rounded-full h-6 w-6 sm:h-8 sm:w-8 border-b-2 border-blue-600"></div>
              <span class="ml-2 text-gray-600 text-xs sm:text-sm">Cargando...</span>
            </div>

            <!-- Mensaje de error -->
            <div v-else-if="globalHistoryError"
              class="bg-red-50 border border-red-200 rounded-md p-3 sm:p-4 mb-3 sm:mb-4">
              <div class="flex">
                <div class="flex-shrink-0">
                  <svg class="h-4 w-4 sm:h-5 sm:w-5 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                      d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                </div>
                <div class="ml-2 sm:ml-3 flex-1">
                  <h3 class="text-xs sm:text-sm font-medium text-red-800">Error al cargar historial</h3>
                  <div class="mt-1 sm:mt-2 text-xs sm:text-sm text-red-700">{{ globalHistoryError }}</div>
                  <div class="mt-3 sm:mt-4">
                    <button @click="loadGlobalHistory" class="btn-secondary text-xs sm:text-sm px-3 py-1.5">
                      Reintentar
                    </button>
                  </div>
                </div>
              </div>
            </div>

            <!-- Tabla de datos -->
            <div v-else-if="!globalHistoryData.length" class="text-center py-6 sm:py-8">
              <svg class="mx-auto h-10 w-10 sm:h-12 sm:w-12 text-gray-400" fill="none" stroke="currentColor"
                viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <h3 class="mt-2 text-xs sm:text-sm font-medium text-gray-900">No hay historial disponible</h3>
              <p class="mt-1 text-xs sm:text-sm text-gray-500">No hay movimientos registrados.</p>
            </div>

            <!-- Mensaje de búsqueda sin resultados -->
            <div v-else-if="hasActiveFilters && !filteredGlobalHistory.length" class="text-center py-6 sm:py-8">
              <svg class="mx-auto h-10 w-10 sm:h-12 sm:w-12 text-gray-400" fill="none" stroke="currentColor"
                viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                  d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
              <h3 class="mt-2 text-xs sm:text-sm font-medium text-gray-900">No se encontraron resultados</h3>
              <p class="mt-1 text-xs sm:text-sm text-gray-500 px-4">
                No hay movimientos que coincidan con los filtros aplicados.
              </p>
            </div>

            <!-- Vista de tarjetas para móviles y tabla para desktop (Historial Global) -->
            <div v-else-if="filteredGlobalHistory.length > 0">
              <!-- Vista de tarjetas para móviles -->
              <div class="md:hidden">
                <div class="space-y-3">
                  <div v-for="movement in paginatedGlobalHistory" :key="movement.id"
                    class="bg-white border border-gray-200 rounded-lg p-3 shadow-sm">
                    <!-- Header: Fecha y N° Lote -->
                    <div class="flex items-start justify-between mb-2 pb-2 border-b border-gray-100">
                      <div class="flex-1 min-w-0">
                        <div class="flex items-center gap-2 mb-1">
                          <svg class="h-3.5 w-3.5 text-gray-400 flex-shrink-0" fill="none" stroke="currentColor"
                            viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                              d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                          </svg>
                          <span class="text-xs font-medium text-gray-900">{{ formatDate(movement.date_time) }}</span>
                        </div>
                        <div class="flex items-center gap-2">
                          <span class="text-[10px] text-gray-500 font-medium">Lote:</span>
                          <span class="text-xs font-mono text-gray-900">{{ movement.batch_id }}</span>
                        </div>
                      </div>
                    </div>

                    <!-- Tipo de cambio -->
                    <div class="mb-2">
                      <div class="text-[10px] text-gray-500 font-medium mb-1">Tipo de cambio:</div>
                      <div class="text-xs font-medium text-gray-900 bg-gray-50 rounded px-2 py-1.5">
                        {{ movement.change_details }}
                      </div>
                    </div>

                    <!-- Valores Anteriores y Nuevos -->
                    <div class="grid grid-cols-2 gap-2 mb-2">
                      <!-- Valores Anteriores -->
                      <div class="bg-red-50 rounded p-2">
                        <div class="text-[10px] text-red-700 font-semibold mb-1 uppercase">Anterior</div>
                        <div v-if="movement.previous_values && Object.keys(movement.previous_values).length > 0"
                          class="text-[10px] text-gray-700 space-y-0.5">
                          <div v-if="movement.previous_values.amount" class="truncate">
                            <span class="font-medium">Cant:</span> {{ movement.previous_values.amount }}
                          </div>
                          <div v-if="movement.previous_values.supplier" class="truncate">
                            <span class="font-medium">Prov:</span> {{ movement.previous_values.supplier }}
                          </div>
                          <div v-if="movement.previous_values.expiration_date" class="truncate">
                            <span class="font-medium">Venc:</span> {{
                              formatDate(movement.previous_values.expiration_date) }}
                          </div>
                        </div>
                        <span v-else class="text-[10px] text-gray-400">N/A</span>
                      </div>

                      <!-- Valores Nuevos -->
                      <div class="bg-green-50 rounded p-2">
                        <div class="text-[10px] text-green-700 font-semibold mb-1 uppercase">Nuevo</div>
                        <div v-if="movement.new_values && Object.keys(movement.new_values).length > 0"
                          class="text-[10px] text-gray-700 space-y-0.5">
                          <div v-if="movement.new_values.amount" class="truncate">
                            <span class="font-medium">Cant:</span> {{ movement.new_values.amount }}
                          </div>
                          <div v-if="movement.new_values.supplier" class="truncate">
                            <span class="font-medium">Prov:</span> {{ movement.new_values.supplier }}
                          </div>
                          <div v-if="movement.new_values.expiration_date" class="truncate">
                            <span class="font-medium">Venc:</span> {{ formatDate(movement.new_values.expiration_date) }}
                          </div>
                        </div>
                        <span v-else class="text-[10px] text-gray-400">N/A</span>
                      </div>
                    </div>

                    <!-- Usuario -->
                    <div class="flex items-center justify-between text-xs pt-2 border-t border-gray-100">
                      <div class="flex items-center gap-1.5">
                        <svg class="h-3.5 w-3.5 text-gray-400 flex-shrink-0" fill="none" stroke="currentColor"
                          viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                            d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                        </svg>
                        <span class="text-gray-500 text-[10px]">{{ movement.user_rut }}</span>
                      </div>
                      <span class="text-gray-900 font-medium truncate ml-2">{{ movement.user_name || 'N/A' }}</span>
                    </div>
                  </div>
                </div>

              </div>

              <!-- Vista de tabla para desktop -->
              <div class="hidden md:block table-container">
                <table class="min-w-full divide-y divide-gray-200" style="min-width: 900px;">
                  <thead class="table-header">
                    <tr>
                      <th
                        class="px-2 sm:px-3 md:px-6 py-2 sm:py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider whitespace-nowrap min-w-[80px] sm:min-w-[100px] md:min-w-[120px]">
                        <div class="flex items-center justify-between">
                          <span class="text-xs sm:text-sm">Fecha</span>
                          <div class="flex flex-col ml-1 sm:ml-2">
                            <button @click="sortGlobalHistoryBy('date', 'asc')"
                              class="text-gray-400 hover:text-gray-600 p-1"
                              :class="{ 'text-blue-600': globalHistorySortField === 'date' && globalHistorySortDirection === 'asc' }">
                              <svg class="w-2 h-2 sm:w-3 sm:h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                  d="M5 15l7-7 7 7" />
                              </svg>
                            </button>
                            <button @click="sortGlobalHistoryBy('date', 'desc')"
                              class="text-gray-400 hover:text-gray-600 p-1"
                              :class="{ 'text-blue-600': globalHistorySortField === 'date' && globalHistorySortDirection === 'desc' }">
                              <svg class="w-2 h-2 sm:w-3 sm:h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                  d="M19 9l-7 7-7-7" />
                              </svg>
                            </button>
                          </div>
                        </div>
                      </th>
                      <th
                        class="px-3 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider whitespace-nowrap min-w-[80px] sm:min-w-[100px]">
                        <div class="flex items-center justify-between">
                          <span class="text-xs sm:text-sm">N° Lote</span>
                          <div class="flex flex-col ml-1 sm:ml-2">
                            <button @click="sortGlobalHistoryBy('batch_id', 'asc')"
                              class="text-gray-400 hover:text-gray-600 p-1"
                              :class="{ 'text-blue-600': globalHistorySortField === 'batch_id' && globalHistorySortDirection === 'asc' }">
                              <svg class="w-2 h-2 sm:w-3 sm:h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                  d="M5 15l7-7 7 7" />
                              </svg>
                            </button>
                            <button @click="sortGlobalHistoryBy('batch_id', 'desc')"
                              class="text-gray-400 hover:text-gray-600 p-1"
                              :class="{ 'text-blue-600': globalHistorySortField === 'batch_id' && globalHistorySortDirection === 'desc' }">
                              <svg class="w-2 h-2 sm:w-3 sm:h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                  d="M19 9l-7 7-7-7" />
                              </svg>
                            </button>
                          </div>
                        </div>
                      </th>
                      <th
                        class="px-3 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider whitespace-nowrap min-w-[120px] sm:min-w-[150px]">
                        <div class="flex items-center justify-between">
                          <span class="text-xs sm:text-sm">Detalles</span>
                          <div class="flex flex-col ml-1 sm:ml-2">
                            <button @click="sortGlobalHistoryBy('change_details', 'asc')"
                              class="text-gray-400 hover:text-gray-600 p-1"
                              :class="{ 'text-blue-600': globalHistorySortField === 'change_details' && globalHistorySortDirection === 'asc' }">
                              <svg class="w-2 h-2 sm:w-3 sm:h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                  d="M5 15l7-7 7 7" />
                              </svg>
                            </button>
                            <button @click="sortGlobalHistoryBy('change_details', 'desc')"
                              class="text-gray-400 hover:text-gray-600 p-1"
                              :class="{ 'text-blue-600': globalHistorySortField === 'change_details' && globalHistorySortDirection === 'desc' }">
                              <svg class="w-2 h-2 sm:w-3 sm:h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                  d="M19 9l-7 7-7-7" />
                              </svg>
                            </button>
                          </div>
                        </div>
                      </th>
                      <th
                        class="px-3 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider whitespace-nowrap min-w-[100px] sm:min-w-[120px]">
                        <div class="flex items-center justify-between">
                          <span class="text-xs sm:text-sm">Valores Ant.</span>
                          <div class="flex flex-col ml-1 sm:ml-2">
                            <button @click="sortGlobalHistoryBy('previous_values', 'asc')"
                              class="text-gray-400 hover:text-gray-600 p-1"
                              :class="{ 'text-blue-600': globalHistorySortField === 'previous_values' && globalHistorySortDirection === 'asc' }">
                              <svg class="w-2 h-2 sm:w-3 sm:h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                  d="M5 15l7-7 7 7" />
                              </svg>
                            </button>
                            <button @click="sortGlobalHistoryBy('previous_values', 'desc')"
                              class="text-gray-400 hover:text-gray-600 p-1"
                              :class="{ 'text-blue-600': globalHistorySortField === 'previous_values' && globalHistorySortDirection === 'desc' }">
                              <svg class="w-2 h-2 sm:w-3 sm:h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                  d="M19 9l-7 7-7-7" />
                              </svg>
                            </button>
                          </div>
                        </div>
                      </th>
                      <th
                        class="px-3 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider whitespace-nowrap min-w-[100px] sm:min-w-[120px]">
                        <div class="flex items-center justify-between">
                          <span class="text-xs sm:text-sm">Valores Nuevos</span>
                          <div class="flex flex-col ml-1 sm:ml-2">
                            <button @click="sortGlobalHistoryBy('new_values', 'asc')"
                              class="text-gray-400 hover:text-gray-600 p-1"
                              :class="{ 'text-blue-600': globalHistorySortField === 'new_values' && globalHistorySortDirection === 'asc' }">
                              <svg class="w-2 h-2 sm:w-3 sm:h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                  d="M5 15l7-7 7 7" />
                              </svg>
                            </button>
                            <button @click="sortGlobalHistoryBy('new_values', 'desc')"
                              class="text-gray-400 hover:text-gray-600 p-1"
                              :class="{ 'text-blue-600': globalHistorySortField === 'new_values' && globalHistorySortDirection === 'desc' }">
                              <svg class="w-2 h-2 sm:w-3 sm:h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                  d="M19 9l-7 7-7-7" />
                              </svg>
                            </button>
                          </div>
                        </div>
                      </th>
                      <th
                        class="px-3 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider whitespace-nowrap min-w-[80px] sm:min-w-[100px]">
                        <div class="flex items-center justify-between">
                          <span class="text-xs sm:text-sm">RUT</span>
                          <div class="flex flex-col ml-1 sm:ml-2">
                            <button @click="sortGlobalHistoryBy('user_rut', 'asc')"
                              class="text-gray-400 hover:text-gray-600 p-1"
                              :class="{ 'text-blue-600': globalHistorySortField === 'user_rut' && globalHistorySortDirection === 'asc' }">
                              <svg class="w-2 h-2 sm:w-3 sm:h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                  d="M5 15l7-7 7 7" />
                              </svg>
                            </button>
                            <button @click="sortGlobalHistoryBy('user_rut', 'desc')"
                              class="text-gray-400 hover:text-gray-600 p-1"
                              :class="{ 'text-blue-600': globalHistorySortField === 'user_rut' && globalHistorySortDirection === 'desc' }">
                              <svg class="w-2 h-2 sm:w-3 sm:h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                  d="M19 9l-7 7-7-7" />
                              </svg>
                            </button>
                          </div>
                        </div>
                      </th>
                      <th
                        class="px-3 sm:px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider whitespace-nowrap min-w-[100px] sm:min-w-[120px]">
                        <div class="flex items-center justify-between">
                          <span class="text-xs sm:text-sm">Usuario</span>
                          <div class="flex flex-col ml-1 sm:ml-2">
                            <button @click="sortGlobalHistoryBy('user_name', 'asc')"
                              class="text-gray-400 hover:text-gray-600 p-1"
                              :class="{ 'text-blue-600': globalHistorySortField === 'user_name' && globalHistorySortDirection === 'asc' }">
                              <svg class="w-2 h-2 sm:w-3 sm:h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                  d="M5 15l7-7 7 7" />
                              </svg>
                            </button>
                            <button @click="sortGlobalHistoryBy('user_name', 'desc')"
                              class="text-gray-400 hover:text-gray-600 p-1"
                              :class="{ 'text-blue-600': globalHistorySortField === 'user_name' && globalHistorySortDirection === 'desc' }">
                              <svg class="w-2 h-2 sm:w-3 sm:h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                                  d="M19 9l-7 7-7-7" />
                              </svg>
                            </button>
                          </div>
                        </div>
                      </th>
                    </tr>
                  </thead>
                  <tbody class="bg-white divide-y divide-gray-200">
                    <tr v-for="movement in paginatedGlobalHistory" :key="movement.id" class="hover:bg-gray-50">
                      <td
                        class="px-2 sm:px-3 md:px-6 py-2 sm:py-3 md:py-4 whitespace-nowrap text-[10px] sm:text-xs md:text-sm font-mono text-gray-900">
                        {{ formatDate(movement.date_time) }}
                      </td>
                      <td
                        class="px-2 sm:px-3 md:px-6 py-2 sm:py-3 md:py-4 whitespace-nowrap text-[10px] sm:text-xs md:text-sm font-mono text-gray-900">
                        {{ movement.batch_id }}
                      </td>
                      <td class="px-2 sm:px-3 md:px-6 py-2 sm:py-3 md:py-4 whitespace-nowrap">
                        <div class="font-medium text-gray-900 text-[10px] sm:text-xs md:text-sm">{{
                          movement.change_details }}</div>
                      </td>
                      <td class="px-2 sm:px-3 md:px-6 py-2 sm:py-3 md:py-4 whitespace-nowrap text-gray-700">
                        <div v-if="movement.previous_values && Object.keys(movement.previous_values).length > 0"
                          class="text-[9px] sm:text-[10px] md:text-xs space-y-0.5">
                          <div v-if="movement.previous_values.amount" class="truncate">Cant: {{
                            movement.previous_values.amount }}</div>
                          <div v-if="movement.previous_values.supplier" class="truncate">Prov: {{
                            movement.previous_values.supplier }}</div>
                          <div v-if="movement.previous_values.expiration_date" class="truncate">Venc: {{
                            formatDate(movement.previous_values.expiration_date) }}</div>
                        </div>
                        <span v-else class="text-gray-400 text-[10px] sm:text-xs">N/A</span>
                      </td>
                      <td class="px-2 sm:px-3 md:px-6 py-2 sm:py-3 md:py-4 whitespace-nowrap text-gray-700">
                        <div v-if="movement.new_values && Object.keys(movement.new_values).length > 0"
                          class="text-[9px] sm:text-[10px] md:text-xs space-y-0.5">
                          <div v-if="movement.new_values.amount" class="truncate">Cant: {{ movement.new_values.amount }}
                          </div>
                          <div v-if="movement.new_values.supplier" class="truncate">Prov: {{
                            movement.new_values.supplier }}</div>
                          <div v-if="movement.new_values.expiration_date" class="truncate">Venc: {{
                            formatDate(movement.new_values.expiration_date) }}</div>
                        </div>
                        <span v-else class="text-gray-400 text-[10px] sm:text-xs">N/A</span>
                      </td>
                      <td
                        class="px-2 sm:px-3 md:px-6 py-2 sm:py-3 md:py-4 whitespace-nowrap text-[10px] sm:text-xs md:text-sm text-gray-700">
                        {{ movement.user_rut }}
                      </td>
                      <td
                        class="px-2 sm:px-3 md:px-6 py-2 sm:py-3 md:py-4 whitespace-nowrap text-[10px] sm:text-xs md:text-sm text-gray-700">
                        {{ movement.user_name || 'N/A' }}
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>
            </div>

            <!-- Paginación -->
            <div class="flex flex-col sm:flex-row items-center justify-between mt-4 sm:mt-6 gap-3 sm:gap-4">
              <div class="text-xs sm:text-sm text-gray-700 text-center sm:text-left">
                Mostrando {{ filteredGlobalHistory.length > 0 ? globalHistoryStartIndex + 1 : 0 }} a {{
                globalHistoryEndIndex }} de {{
                  filteredGlobalHistory.length }} resultados
              </div>
              <div class="flex items-center gap-2 w-full sm:w-auto">
                <button class="btn-secondary px-3 sm:px-4 py-2 text-xs sm:text-sm flex-1 sm:flex-none sm:min-w-[80px]"
                  :disabled="globalHistoryCurrentPage === 1" @click="globalHistoryCurrentPage--">
                  <span class="hidden sm:inline">Anterior</span>
                  <span class="sm:hidden">Ant.</span>
                </button>
                <span
                  class="px-3 py-2 text-xs sm:text-sm text-gray-700 bg-gray-100 rounded-md text-center whitespace-nowrap">
                  <span class="hidden sm:inline">Página {{ globalHistoryCurrentPage }} de {{ globalHistoryTotalPages
                    }}</span>
                  <span class="sm:hidden">{{ globalHistoryCurrentPage }} / {{ globalHistoryTotalPages }}</span>
                </span>
                <button class="btn-secondary px-3 sm:px-4 py-2 text-xs sm:text-sm flex-1 sm:flex-none sm:min-w-[80px]"
                  :disabled="globalHistoryCurrentPage === globalHistoryTotalPages" @click="globalHistoryCurrentPage++">
                  <span class="hidden sm:inline">Siguiente</span>
                  <span class="sm:hidden">Sig.</span>
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>


      <!-- Paginación -->
      <div class="flex flex-col sm:flex-row items-center justify-between mt-6 gap-4">
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
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router';
import { useAuthStore } from '@/stores/auth'
import { format } from 'date-fns'
import { es } from 'date-fns/locale'
import inventoryService from '@/services/inventory/inventoryService'
import qrService from '@/services/qr/qrService'
import QrcodeVue from 'qrcode.vue'
import { exportToExcel as exportExcel, formatDateForExcel } from '@/utils/excelExport'
import { useNotification } from '@/composables/useNotification'

const route = useRoute()
const authStore = useAuthStore()
const { success: showSuccess, error: showError } = useNotification()

// Estado reactivo
const supplies = ref([])
const loading = ref(false)
const error = ref(null)
const searchTerm = ref('')
const sortField = ref('none')
const sortDirection = ref('asc')
const currentPage = ref(1)
const itemsPerPage = 10

// Estado del modal de edición
const showEditModal = ref(false)
const editingSupply = ref({})
const saving = ref(false)

// Estado del modal de historial global
const showGlobalHistoryModal = ref(false)
const globalHistoryChangeTypeFilter = ref('')
const globalHistoryBatchFilter = ref('')
const globalHistoryUserFilter = ref('')
const globalHistorySortField = ref('date')
const globalHistorySortDirection = ref('desc')
const globalHistoryCurrentPage = ref(1)
const globalHistoryItemsPerPage = 5
const globalHistoryLoading = ref(false)
const globalHistoryError = ref(null)
const globalHistoryData = ref([])

// Estado del modal de historial individual
const showHistoryModal = ref(false)
const historySearchTerm = ref('')
const historySortField = ref('date')
const historySortDirection = ref('desc')
const historyCurrentPage = ref(1)
const historyItemsPerPage = 10
const historyLoading = ref(false)
const historyError = ref(null)
const historyData = ref([])
const selectedSupplyForHistory = ref(null)

// Estado del modal de detalles de lote con paginación
const batchDetailsCurrentPage = ref(1)
const batchDetailsItemsPerPage = 5 // Constante, no ref

// Sistema de notificaciones unificado - usa useNotification()

// Computed properties
const filteredSupplies = computed(() => {
  let filtered = [...supplies.value]

  if (searchTerm.value) {
    filtered = filtered.filter(supply =>
      supply.batch_id.toString().includes(searchTerm.value) ||
      supply.code.toString().includes(searchTerm.value) ||
      supply.name.toLowerCase().includes(searchTerm.value.toLowerCase()) ||
      supply.supplier.toLowerCase().includes(searchTerm.value.toLowerCase())
    )
  }

  // Solo aplicar ordenamiento si se ha seleccionado explícitamente un campo
  if (sortField.value && sortField.value !== 'none') {
    filtered.sort((a, b) => {
      let result = 0

      switch (sortField.value) {
        case 'batch_id':
          result = a.batch_id - b.batch_id
          break
        case 'code':
          result = a.code - b.code
          break
        case 'expiration_date':
          result = new Date(a.expiration_date) - new Date(b.expiration_date)
          break
        case 'amount':
          result = a.amount - b.amount
          break
        case 'supplier':
          result = (a.supplier || '').localeCompare(b.supplier || '')
          break
        default:
          result = a.name.localeCompare(b.name)
          break
      }

      return sortDirection.value === 'asc' ? result : -result
    })
  }

  return filtered
})

const totalPages = computed(() => Math.ceil(filteredSupplies.value.length / itemsPerPage))

const startIndex = computed(() => (currentPage.value - 1) * itemsPerPage)
const endIndex = computed(() => Math.min(startIndex.value + itemsPerPage, filteredSupplies.value.length))

const paginatedSupplies = computed(() => {
  return filteredSupplies.value.slice(startIndex.value, endIndex.value)
})

// Computed properties para historial
const filteredHistory = computed(() => {
  let filtered = [...historyData.value]

  if (historySearchTerm.value) {
    filtered = filtered.filter(movement =>
      formatDate(movement.date_time || movement.date).includes(historySearchTerm.value) ||
      (movement.action || movement.type || '').toLowerCase().includes(historySearchTerm.value.toLowerCase()) ||
      (movement.details?.amount?.toString() || movement.amount?.toString() || '').includes(historySearchTerm.value) ||
      (movement.user_rut || movement.user || '').toLowerCase().includes(historySearchTerm.value.toLowerCase())
    )
  }

  // Ordenamiento
  filtered.sort((a, b) => {
    let result = 0

    switch (historySortField.value) {
      case 'date':
        result = new Date(a.date_time || a.date) - new Date(b.date_time || b.date)
        break
      case 'type':
        result = (a.action || a.type || '').localeCompare(b.action || b.type || '')
        break
      case 'amount':
        const amountA = a.details?.amount || a.amount || 0
        const amountB = b.details?.amount || b.amount || 0
        result = amountA - amountB
        break
      case 'user':
        result = (a.user_rut || a.user || '').localeCompare(b.user_rut || b.user || '')
        break
      default:
        result = new Date(a.date_time || a.date) - new Date(b.date_time || b.date) // Fallback to date sort
        break
    }

    return historySortDirection.value === 'asc' ? result : -result
  })

  return filtered
})

const historyTotalPages = computed(() => Math.max(1, Math.ceil(filteredHistory.value.length / historyItemsPerPage)))

const historyStartIndex = computed(() => {
  if (filteredHistory.value.length === 0) return 0
  return (historyCurrentPage.value - 1) * historyItemsPerPage
})
const historyEndIndex = computed(() => {
  if (filteredHistory.value.length === 0) return 0
  return Math.min(historyStartIndex.value + historyItemsPerPage, filteredHistory.value.length)
})

const paginatedHistory = computed(() => {
  return filteredHistory.value.slice(historyStartIndex.value, historyEndIndex.value)
})

// Computed properties para historial global
const filteredGlobalHistory = computed(() => {
  let filtered = [...globalHistoryData.value]

  // Filtro por tipo de cambio
  if (globalHistoryChangeTypeFilter.value) {
    filtered = filtered.filter(movement => {
      const changeDetails = (movement.change_details || '').toLowerCase()
      const filterValue = globalHistoryChangeTypeFilter.value.toLowerCase()

      // Buscar en el formato "Lote actualizado: [campo]"
      if (filterValue === 'cantidad') {
        return changeDetails.includes('cantidad')
      } else if (filterValue === 'proveedor') {
        return changeDetails.includes('proveedor')
      } else if (filterValue === 'fecha de expiración') {
        return changeDetails.includes('fecha de expiración')
      } else if (filterValue === 'creado') {
        return changeDetails.includes('creado') && !changeDetails.includes('actualizado')
      } else if (filterValue === 'eliminado') {
        return changeDetails.includes('eliminado')
      }

      // Si no hay filtro específico, mostrar todos
      return true
    })
  }

  // Filtro por número de lote
  if (globalHistoryBatchFilter.value) {
    filtered = filtered.filter(movement =>
      (movement.batch_id?.toString() || '').includes(globalHistoryBatchFilter.value)
    )
  }

  // Filtro por usuario (RUT o nombre)
  if (globalHistoryUserFilter.value) {
    filtered = filtered.filter(movement =>
      (movement.user_rut || '').toLowerCase().includes(globalHistoryUserFilter.value.toLowerCase()) ||
      (movement.user_name || '').toLowerCase().includes(globalHistoryUserFilter.value.toLowerCase())
    )
  }

  // Ordenamiento
  filtered.sort((a, b) => {
    let result = 0

    switch (globalHistorySortField.value) {
      case 'date':
        result = new Date(a.date_time) - new Date(b.date_time)
        break
      case 'batch_id':
        result = (a.batch_id || 0) - (b.batch_id || 0)
        break
      case 'change_details':
        result = (a.change_details || '').localeCompare(b.change_details || '')
        break
      case 'previous_values':
        result = JSON.stringify(a.previous_values || {}).localeCompare(JSON.stringify(b.previous_values || {}))
        break
      case 'new_values':
        result = JSON.stringify(a.new_values || {}).localeCompare(JSON.stringify(b.new_values || {}))
        break
      case 'user_rut':
        result = (a.user_rut || '').localeCompare(b.user_rut || '')
        break
      case 'user_name':
        result = (a.user_name || '').localeCompare(b.user_name || '')
        break
      default:
        result = new Date(a.date_time) - new Date(b.date_time) // Fallback to date sort
        break
    }

    return globalHistorySortDirection.value === 'asc' ? result : -result
  })

  return filtered
})

const globalHistoryTotalPages = computed(() => Math.max(1, Math.ceil(filteredGlobalHistory.value.length / globalHistoryItemsPerPage)))

const globalHistoryStartIndex = computed(() => {
  if (filteredGlobalHistory.value.length === 0) return 0
  return (globalHistoryCurrentPage.value - 1) * globalHistoryItemsPerPage
})
const globalHistoryEndIndex = computed(() => {
  if (filteredGlobalHistory.value.length === 0) return 0
  return Math.min(globalHistoryStartIndex.value + globalHistoryItemsPerPage, filteredGlobalHistory.value.length)
})

const paginatedGlobalHistory = computed(() => {
  return filteredGlobalHistory.value.slice(globalHistoryStartIndex.value, globalHistoryEndIndex.value)
})

// Verificar si hay filtros activos
const hasActiveFilters = computed(() => {
  return globalHistoryChangeTypeFilter.value ||
    globalHistoryBatchFilter.value ||
    globalHistoryUserFilter.value
})

// Computed properties para paginación del modal de detalles de lote
const batchDetailsTotalPages = computed(() => Math.ceil(batchDetails.value.length / batchDetailsItemsPerPage))

const batchDetailsStartIndex = computed(() => (batchDetailsCurrentPage.value - 1) * batchDetailsItemsPerPage)
const batchDetailsEndIndex = computed(() => Math.min(batchDetailsStartIndex.value + batchDetailsItemsPerPage, batchDetails.value.length))

const paginatedBatchDetails = computed(() => {
  return batchDetails.value.slice(batchDetailsStartIndex.value, batchDetailsEndIndex.value)
})

const clearSearch = () => {
  searchTerm.value = ''
  sortField.value = 'none'
  sortDirection.value = 'asc'
  currentPage.value = 1
}

const sortBy = (field, direction) => {
  sortField.value = field
  sortDirection.value = direction
  currentPage.value = 1
}

const formatDate = (dateString) => {
  try {
    return format(new Date(dateString), 'dd/MM/yyyy', { locale: es })
  } catch {
    return dateString
  }
}

const getExpirationClass = (expirationDate) => {
  const today = new Date()
  const expDate = new Date(expirationDate)
  const daysUntilExpiration = Math.ceil((expDate - today) / (1000 * 60 * 60 * 24))

  if (daysUntilExpiration < 0) return 'text-red-600 font-semibold'
  if (daysUntilExpiration <= 15) return 'text-red-600 font-semibold'
  if (daysUntilExpiration <= 30) return 'text-orange-600 font-semibold'
  return 'text-gray-900'
}

const getAmountClass = (amount) => {

  if (amount < 5) return 'text-red-600 font-semibold'
  if (amount < 10) return 'text-orange-600 font-semibold'
  return 'text-gray-900'
}

const getStatusBadgeClass = (status) => {
  const statusLower = (status || '').toLowerCase()

  if (statusLower.includes('creado') || statusLower.includes('crear')) {
    return 'bg-green-100 text-green-800'
  } else if (statusLower.includes('eliminado') || statusLower.includes('eliminar')) {
    return 'bg-red-100 text-red-800'
  } else if (statusLower.includes('actualizado') || statusLower.includes('actualizar')) {
    return 'bg-blue-100 text-blue-800'
  } else if (statusLower.includes('cantidad')) {
    return 'bg-purple-100 text-purple-800'
  } else if (statusLower.includes('proveedor')) {
    return 'bg-yellow-100 text-yellow-800'
  } else if (statusLower.includes('vencimiento') || statusLower.includes('expiración')) {
    return 'bg-orange-100 text-orange-800'
  }

  return 'bg-gray-100 text-gray-800'
}

const viewMovementDetails = (movement) => {
  // Mostrar detalles completos del movimiento
  let detailsText = `Fecha: ${formatDate(movement.date_time || movement.date)}\n`
  detailsText += `Tipo: ${movement.action || movement.type || 'N/A'}\n`
  detailsText += `Cantidad: ${movement.details?.amount || movement.amount || 'N/A'} unidades\n`
  detailsText += `Usuario: ${movement.user_rut || movement.user || 'N/A'}\n`

  if (movement.details) {
    detailsText += `\nDetalles adicionales:\n${JSON.stringify(movement.details, null, 2)}`
  }

  alert(`Detalles del Movimiento\n\n${detailsText}`)
}

// Función para obtener la descripción del estado en español
const getStatusDescription = (status) => {
  const statusMap = {
    'disponible': 'Disponible',
    'en_camino_a_pabellon': 'En camino a pabellón',
    'recepcionado': 'Recepcionado',
    'consumido': 'Consumido',
    'en_camino_a_bodega': 'En camino a bodega'
  }
  return statusMap[status] || 'Estado desconocido'
}

// Función para obtener la clase CSS del estado
const getStatusClass = (status) => {
  switch (status) {
    case 'disponible':
      return 'text-green-700 bg-green-100'
    case 'en_camino_a_pabellon':
      return 'text-blue-700 bg-blue-100'
    case 'recepcionado':
      return 'text-purple-700 bg-purple-100'
    case 'consumido':
      return 'text-gray-700 bg-gray-100'
    case 'en_camino_a_bodega':
      return 'text-orange-700 bg-orange-100'
    default:
      return 'text-gray-700 bg-gray-100'
  }
}

const editSupply = (supply) => {
  // Preparar los datos para edición
  editingSupply.value = {
    ...supply,
    // Convertir la fecha al formato requerido por el input date (YYYY-MM-DD)
    expiration_date: supply.expiration_date ? new Date(supply.expiration_date).toISOString().split('T')[0] : ''
  }

  // Log para debugging: verificar que el store_id se esté recibiendo correctamente
  console.log('Editando lote:', supply)
  console.log('Store ID del lote:', supply.store_id)

  showEditModal.value = true
}

const closeEditModal = () => {
  showEditModal.value = false
  editingSupply.value = {}
  saving.value = false
}

const saveEdit = async () => {
  if (!editingSupply.value.batch_id) {
    showError('ID de lote no encontrado')
    return
  }

  saving.value = true

  try {
    // Preparar los datos para enviar al backend
    // Convertir la fecha al formato ISO que espera Go
    let formattedDate = editingSupply.value.expiration_date
    if (formattedDate) {
      // Crear un objeto Date y convertirlo al formato ISO completo
      const date = new Date(formattedDate)
      formattedDate = date.toISOString()
    }

    // IMPORTANTE: Mantener el store_id original del lote para evitar cambiar su ubicación
    const batchData = {
      expiration_date: formattedDate,
      amount: parseInt(editingSupply.value.amount),
      supplier: editingSupply.value.supplier,
      store_id: editingSupply.value.store_id || 1 // Mantener el store_id original del lote
    }

    // Log para debugging: verificar que se esté enviando el store_id correcto
    console.log('Enviando datos del lote:', batchData)
    console.log('Store ID que se mantiene:', batchData.store_id)

    // Actualizar el batch usando el servicio
    await inventoryService.updateBatch(editingSupply.value.batch_id, batchData)

    // Actualizar la lista local
    const index = supplies.value.findIndex(s => s.batch_id === editingSupply.value.batch_id)
    if (index !== -1) {
      supplies.value[index] = {
        ...supplies.value[index],
        expiration_date: editingSupply.value.expiration_date,
        amount: editingSupply.value.amount,
        supplier: editingSupply.value.supplier
      }
    }

    // Cerrar el modal y mostrar mensaje de éxito
    closeEditModal()
    showSuccess('Lote actualizado exitosamente')

  } catch (error) {
    console.error('Error al actualizar lote:', error)
    showError('Error al actualizar el lote: ' + (error.response?.data?.error || error.message))
  } finally {
    saving.value = false
  }
}

const deleteSupply = async (supply) => {
  if (!confirm(`¿Está seguro de que desea eliminar el lote de ${supply.name}?`)) {
    return
  }
  
  try {
    await inventoryService.deleteBatch(supply.batch_id)
    showSuccess('Lote eliminado exitosamente')
    await loadInventory()
  } catch (error) {
    console.error('Error al eliminar lote:', error)
    showError('Error al eliminar el lote: ' + (error.response?.data?.error || error.message))
  }
}

// Funciones de notificación - usa el sistema unificado
const { info: showInfo } = useNotification()
const showNotification = (message, type = 'info') => {
  if (type === 'success') {
    showSuccess(message)
  } else if (type === 'error') {
    showError(message)
  } else {
    showInfo(message)
  }
}

// Modal de detalles de lote
const showBatchDetailsModal = ref(false)
const batchDetails = ref([])
const batchDetailsLoading = ref(false)
const batchDetailsError = ref(null)
const selectedBatch = ref(null)
const supplyName = ref('')
const openBatchDetailsModal = async (batch) => {
  showBatchDetailsModal.value = true
  batchDetailsLoading.value = true
  batchDetailsError.value = null
  batchDetailsCurrentPage.value = 1 // Resetear paginación
  selectedBatch.value = batch
  supplyName.value = batch.name || 'Insumo' // Usar el nombre del batch directamente
  try {
    // Obtener insumos individuales por lote
    const supplies = await inventoryService.getAvailableSuppliesByBatch(batch.batch_id)
    batchDetails.value = supplies

    // Solo intentar obtener el nombre si no lo tenemos y hay insumos
    if (!supplyName.value || supplyName.value === 'Insumo') {
      if (supplies.length > 0) {
        try {
          const code = supplies[0].code
          const supplyCodes = await inventoryService.getAllSupplyCodes()
          const found = supplyCodes.find(sc => sc.code === code)
          if (found) {
            supplyName.value = found.name
          }
        } catch (codeError) {
          console.warn('No se pudo obtener el nombre del insumo:', codeError)
          // Mantener el nombre del batch si falla
        }
      }
    }
  } catch (err) {
    batchDetailsError.value = 'Error al cargar insumos del lote: ' + err.message
  } finally {
    batchDetailsLoading.value = false
  }
}

const closeBatchDetailsModal = () => {
  showBatchDetailsModal.value = false
  batchDetails.value = []
  selectedBatch.value = null
  batchDetailsError.value = null
  batchDetailsCurrentPage.value = 1 // Resetear paginación
}

const getQrDataUrl = (code) => {
  return `https://api.qrserver.com/v1/create-qr-code/?data=${encodeURIComponent(code)}&size=200x200`
}

const handleDownloadQR = async (qrCode) => {
  try {
    await qrService.downloadQRImage(qrCode, 'normal')
    showSuccess('QR descargado correctamente')
  } catch (error) {
    showError('Error al descargar el QR')
  }
}

// Métodos para historial global
const openGlobalHistoryModal = () => {
  // Verificar permisos - solo admin y encargado de bodega pueden ver historial
  if (!authStore.canViewAllRequests) {
    showError('No tienes permisos para ver el historial de movimientos del inventario')
    return
  }

  globalHistoryChangeTypeFilter.value = ''
  globalHistoryBatchFilter.value = ''
  globalHistoryUserFilter.value = ''
  globalHistoryCurrentPage.value = 1
  globalHistorySortField.value = 'date'
  globalHistorySortDirection.value = 'desc'
  showGlobalHistoryModal.value = true

  // Cargar historial global
  loadGlobalHistory()
}

const closeGlobalHistoryModal = () => {
  showGlobalHistoryModal.value = false
  globalHistoryData.value = []
  globalHistoryError.value = null
}

const clearGlobalHistoryFilters = () => {
  globalHistoryChangeTypeFilter.value = ''
  globalHistoryBatchFilter.value = ''
  globalHistoryUserFilter.value = ''
  globalHistoryCurrentPage.value = 1
}

const applyFilters = () => {
  globalHistoryCurrentPage.value = 1
}

const sortGlobalHistoryBy = (field, direction) => {
  globalHistorySortField.value = field
  globalHistorySortDirection.value = direction
  globalHistoryCurrentPage.value = 1
}

const loadGlobalHistory = async () => {
  globalHistoryLoading.value = true
  globalHistoryError.value = null

  try {
    // Usar getBatchHistoryWithDetails para obtener el historial completo
    const data = await inventoryService.getBatchHistoryWithDetails()
    globalHistoryData.value = data
  } catch (err) {
    globalHistoryError.value = 'Error al cargar el historial global: ' + err.message
    console.error('Error al cargar historial global:', err)
  } finally {
    globalHistoryLoading.value = false
  }
}

// Métodos para historial individual
const viewSupply = (supply) => {
  // Verificar permisos - solo admin y encargado de bodega pueden ver historial
  if (!authStore.canViewAllRequests) {
    showError('No tienes permisos para ver el historial de movimientos del inventario')
    return
  }

  selectedSupplyForHistory.value = supply
  historySearchTerm.value = ''
  historyCurrentPage.value = 1
  historySortField.value = 'date'
  historySortDirection.value = 'desc'
  showHistoryModal.value = true

  // Cargar historial del insumo
  loadHistory()
}

const closeHistoryModal = () => {
  showHistoryModal.value = false
  selectedSupplyForHistory.value = null
  historyData.value = []
  historyError.value = null
}

const clearHistorySearch = () => {
  historySearchTerm.value = ''
  historyCurrentPage.value = 1
}

const sortHistoryBy = (field, direction) => {
  historySortField.value = field
  historySortDirection.value = direction
  historyCurrentPage.value = 1
}

const loadHistory = async () => {
  if (!selectedSupplyForHistory.value) return

  historyLoading.value = true
  historyError.value = null

  try {
    // Cargar historial del lote específico
    const data = await inventoryService.getBatchHistory(selectedSupplyForHistory.value.batch_id)
    historyData.value = data
  } catch (err) {
    historyError.value = 'Error al cargar el historial: ' + err.message
    console.error('Error al cargar historial:', err)
  } finally {
    historyLoading.value = false
  }
}

// Métodos
const loadInventory = async () => {
  loading.value = true
  error.value = null

  try {
    const data = await inventoryService.getInventory()
    // Si la respuesta tiene inventory_items, usar ese array. Si no, usar data si es array.
    if (Array.isArray(data)) {
      supplies.value = data
    } else if (Array.isArray(data.inventory_items)) {
      supplies.value = data.inventory_items
    } else {
      supplies.value = []
    }
  } catch (err) {
    error.value = 'Error al cargar el inventario: ' + err.message
    console.error('Error al cargar inventario:', err)
  } finally {
    loading.value = false
  }
}

// Lifecycle
const exportToExcel = () => {
  try {
    const columns = [
      { key: 'batch_id', label: 'N° de Lote' },
      { key: 'name', label: 'Nombre del Insumo' },
      { key: 'code', label: 'Código Interno' },
      { key: 'expiration_date', label: 'Fecha de Vencimiento', formatter: formatDateForExcel },
      { key: 'amount', label: 'Cantidad' },
      { key: 'supplier', label: 'Proveedor' },
      { key: 'location_type', label: 'Tipo de Ubicación' },
      { key: 'location_name', label: 'Ubicación' },
      { key: 'created_at', label: 'Fecha de Creación', formatter: formatDateForExcel },
      { key: 'updated_at', label: 'Última Actualización', formatter: formatDateForExcel }
    ]
    
    exportExcel(filteredSupplies.value, columns, 'inventario_insumos')
    showSuccess('Exportación a Excel completada exitosamente')
  } catch (error) {
    console.error('Error al exportar:', error)
    showError('Error al exportar a Excel: ' + error.message)
  }
}

onMounted(() => {
  // Si viene con un término de búsqueda desde Home, aplicarlo
  if (route.query.search) {
    searchTerm.value = route.query.search
  }
  // No aplicar ordenamiento por defecto, mantener el orden de la base de datos
  loadInventory()
})

// Watchers para los filtros del historial global
watch([globalHistoryChangeTypeFilter, globalHistoryBatchFilter, globalHistoryUserFilter], () => {
  if (showGlobalHistoryModal.value) {
    // Resetear la página cuando se cambian los filtros
    globalHistoryCurrentPage.value = 1
  }
})
</script>

<style scoped>
/* Estilos responsivos para la tabla */
.table-container {
  overflow-x: auto;
  border: 1px solid #e5e7eb;
  border-radius: 0.5rem;
  box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1);
  -webkit-overflow-scrolling: touch;
  /* Mejor scroll en iOS */
  position: relative;
}

/* Scrollbar visible pero elegante en móviles */
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

.table-header th:first-child {
  min-width: 120px;
}

.table-header th:nth-child(2) {
  min-width: 200px;
}

.table-header th:nth-child(3) {
  min-width: 140px;
}

.table-header th:nth-child(4) {
  min-width: 150px;
}

.table-header th:nth-child(5) {
  min-width: 120px;
}

.table-header th:nth-child(6) {
  min-width: 150px;
}

.table-header th:last-child {
  min-width: 120px;
}

.table-body tr {
  background-color: white;
  border-bottom: 1px solid #e5e7eb;
}

.table-body tr:hover {
  background-color: #f9fafb;
}

.table-body td {
  padding: 1rem 1.5rem;
  white-space: nowrap;
  vertical-align: middle;
}

/* Responsive para móviles */
@media (max-width: 640px) {

  .table-header th,
  .table-body td {
    padding: 0.5rem 0.75rem;
    font-size: 0.75rem;
  }

  .table-header th:first-child {
    min-width: 90px;
  }

  .table-header th:nth-child(2) {
    min-width: 140px;
  }

  .table-header th:nth-child(3) {
    min-width: 100px;
  }

  .table-header th:nth-child(4) {
    min-width: 110px;
  }

  .table-header th:nth-child(5) {
    min-width: 90px;
  }

  .table-header th:nth-child(6) {
    min-width: 110px;
  }

  .table-header th:last-child {
    min-width: 90px;
  }

  /* Ajustes para modales en móviles */
  .fixed.inset-0>div {
    max-height: 95vh;
    overflow-y: auto;
  }

  /* Mejor visualización de botones en móviles */
  button {
    touch-action: manipulation;
    -webkit-tap-highlight-color: transparent;
  }
}

/* Responsive para tablets */
@media (max-width: 768px) {
  .table-header th:nth-child(3) {
    min-width: 120px;
  }

  .table-header th:nth-child(4) {
    min-width: 120px;
  }

  .table-header th:nth-child(6) {
    min-width: 120px;
  }
}

@media (max-width: 1024px) {

  .table-header th:nth-child(4),
  .table-header th:nth-child(6) {
    min-width: 120px;
  }
}

/* Botones de paginación responsivos */
.pagination-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: space-between;
  margin-top: 1.5rem;
  gap: 1rem;
}

@media (min-width: 640px) {
  .pagination-container {
    flex-direction: row;
    gap: 0;
  }
}

.pagination-button {
  padding: 0.5rem 1rem;
  font-size: 0.875rem;
  min-width: 80px;
  border-radius: 0.375rem;
  font-weight: 500;
  transition: all 0.2s;
}

.pagination-info {
  padding: 0.5rem 0.75rem;
  font-size: 0.875rem;
  background-color: #f3f4f6;
  border-radius: 0.375rem;
  min-width: 100px;
  text-align: center;
}

/* Estilos para el select de filtros */
.form-select {
  @apply block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm;
}

/* Mejoras para dispositivos móviles */
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

  .btn-primary,
  .btn-secondary {
    padding: 0.625rem 1rem;
    font-size: 0.875rem;
    min-height: 40px;
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

  /* Inputs más grandes para fácil toque en móviles */
  input[type="text"],
  input[type="number"],
  input[type="date"],
  select {
    min-height: 42px;
    font-size: 16px;
    /* Previene zoom en iOS */
  }

  /* Modal ocupando casi toda la pantalla en móviles */
  .modal-container {
    max-height: 90vh;
  }
}

/* Mejoras para tablets */
@media (max-width: 768px) {
  .card {
    padding: 1.5rem;
  }

  .btn-primary,
  .btn-secondary {
    padding: 0.625rem 1.25rem;
  }
}

/* Estilos específicos para el historial responsivo */
.history-modal {
  max-height: 90vh;
  overflow-y: auto;
}

.history-filters {
  display: grid;
  gap: 1rem;
}

@media (max-width: 640px) {
  .history-filters {
    grid-template-columns: 1fr;
  }

  .history-table {
    font-size: 0.75rem;
  }

  .history-table th,
  .history-table td {
    padding: 0.5rem 0.75rem;
  }

  .history-pagination {
    flex-direction: column;
    gap: 1rem;
    align-items: center;
  }

  .history-pagination .btn-secondary {
    min-width: 60px;
    padding: 0.5rem 0.75rem;
    font-size: 0.75rem;
  }

  /* Mejoras adicionales para textos en móviles */
  h1 {
    font-size: 1.25rem !important;
  }

  h2 {
    font-size: 1.125rem !important;
  }

  h3 {
    font-size: 1rem !important;
  }

  /* Espaciado optimizado */
  .space-y-6>*+* {
    margin-top: 1rem;
  }
}

@media (min-width: 641px) and (max-width: 1024px) {
  .history-filters {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (min-width: 1025px) {
  .history-filters {
    grid-template-columns: repeat(4, 1fr);
  }
}

/* Mejoras para el scroll horizontal en móviles */
.history-table-container {
  -webkit-overflow-scrolling: touch;
  scrollbar-width: thin;
}

.history-table-container::-webkit-scrollbar {
  height: 6px;
}

.history-table-container::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 3px;
}

.history-table-container::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 3px;
}

.history-table-container::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}

/* Optimización adicional para modales en móviles */
@media (max-width: 640px) {

  /* Asegurar que los modales sean scrolleables */
  .fixed.inset-0 {
    padding: 0;
  }

  /* Mejorar la visualización de badges y etiquetas */
  .inline-flex.items-center {
    font-size: 0.7rem;
    padding: 0.25rem 0.5rem;
  }

  /* Optimizar el espaciado de los formularios */
  form .space-y-4>*+* {
    margin-top: 0.75rem;
  }

  /* Mejorar la altura mínima de los botones para toque fácil */
  button,
  .btn-primary,
  .btn-secondary {
    min-height: 44px;
    /* Tamaño recomendado por Apple para touch targets */
  }

  /* Ajustar el tamaño de los iconos SVG para mejor visibilidad */
  svg {
    flex-shrink: 0;
  }
}

/* Mejoras para tablets y pantallas medianas */
@media (min-width: 641px) and (max-width: 1024px) {
  .card {
    padding: 1.25rem;
  }

  .table-header th {
    font-size: 0.8rem;
  }

  .table-body td {
    font-size: 0.875rem;
  }
}

/* Animaciones suaves para mejorar la experiencia */
@media (prefers-reduced-motion: no-preference) {

  button,
  .btn-primary,
  .btn-secondary {
    transition: all 0.15s ease-in-out;
  }

  .hover\:bg-gray-50:hover {
    transition: background-color 0.15s ease-in-out;
  }
}

/* Estilos para vista de tarjetas en móviles */
@media (max-width: 768px) {

  /* Animación suave al aparecer las tarjetas */
  .space-y-3>div {
    animation: fadeInUp 0.3s ease-out;
  }

  @keyframes fadeInUp {
    from {
      opacity: 0;
      transform: translateY(10px);
    }

    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  /* Efecto hover/active en las tarjetas */
  .space-y-3>div:active {
    transform: scale(0.98);
    transition: transform 0.1s ease-in-out;
  }

  /* Mejorar la legibilidad de los badges */
  .inline-block.rounded-full {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    max-width: 100%;
  }

  /* Estilos mejorados para los botones de paginación móvil */
  .btn-secondary:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }

  /* Mejorar el contraste del indicador de página actual */
  .bg-gray-100 {
    font-weight: 600;
  }
}
</style>