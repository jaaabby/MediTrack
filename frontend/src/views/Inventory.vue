<template>
  <div class="space-y-6">
    <!-- Header del inventario -->
    <div>
      <h1 class="text-2xl font-bold text-gray-900">Inventario de Insumos Médicos</h1>
      <p class="text-gray-600 mt-1">Gestión y control de stock médico</p>
    </div>

    <!-- Filtros y búsqueda -->
    <div class="card">
      <div class="flex items-end space-x-4">
        <!-- Buscador único -->
        <div class="flex-1">
          <label class="block text-sm font-medium text-gray-700 mb-2">Buscar insumo</label>
          <div class="relative">
            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <svg class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
            </div>
            <input
              type="text"
              placeholder="Buscar por número de lote, nombre, código o proveedor..."
              class="form-input pl-10 w-full"
              v-model="searchTerm"
            />
          </div>
        </div>

        <!-- Botón de limpiar búsqueda -->
        <div>
          <button 
            class="btn-secondary px-4 py-2 h-10" 
            @click="clearSearch"
            :disabled="!searchTerm"
          >
            Limpiar
          </button>
        </div>
      </div>
    </div>

    <!-- Tabla de inventario -->
    <div class="card">
      <div class="card-header">
        <div class="flex items-center justify-between">
          <div>
            <h2 class="card-title">Inventario de Insumos Médicos</h2>
            <p class="text-sm text-gray-600">Total: {{ filteredSupplies.length }} lotes</p>
          </div>
          <button class="btn-primary flex items-center" @click="openGlobalHistoryModal">
            <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            Historial de Movimientos
          </button>
        </div>
      </div>

      <!-- Indicador de carga -->
      <div v-if="loading" class="flex justify-center items-center py-8">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
        <span class="ml-2 text-gray-600">Cargando inventario...</span>
      </div>

      <!-- Mensaje de error -->
      <div v-else-if="error" class="bg-red-50 border border-red-200 rounded-md p-4 mx-4 mb-4">
        <div class="flex">
          <div class="flex-shrink-0">
            <svg class="h-5 w-5 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
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
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="table-header">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider whitespace-nowrap min-w-[120px]">
                <div class="flex items-center justify-between">
                  <span>N° de lote</span>
                  <div class="flex flex-col ml-2">
                    <button 
                      @click="sortBy('batch_id', 'asc')" 
                      class="text-gray-400 hover:text-gray-600 p-1"
                      :class="{ 'text-primary-600': sortField === 'batch_id' && sortDirection === 'asc' }"
                    >
                      <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7" />
                      </svg>
                    </button>
                    <button 
                      @click="sortBy('batch_id', 'desc')" 
                      class="text-gray-400 hover:text-gray-600 p-1"
                      :class="{ 'text-primary-600': sortField === 'batch_id' && sortDirection === 'desc' }"
                    >
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
                    <button 
                      @click="sortBy('name', 'asc')" 
                      class="text-gray-400 hover:text-gray-600 p-1"
                      :class="{ 'text-primary-600': sortField === 'name' && sortDirection === 'asc' }"
                    >
                      <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7" />
                      </svg>
                    </button>
                    <button 
                      @click="sortBy('name', 'desc')" 
                      class="text-gray-400 hover:text-gray-600 p-1"
                      :class="{ 'text-primary-600': sortField === 'name' && sortDirection === 'desc' }"
                    >
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
                    <button 
                      @click="sortBy('code', 'asc')" 
                      class="text-gray-400 hover:text-gray-600 p-1"
                      :class="{ 'text-primary-600': sortField === 'code' && sortDirection === 'asc' }"
                    >
                      <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7" />
                      </svg>
                    </button>
                    <button 
                      @click="sortBy('code', 'desc')" 
                      class="text-gray-400 hover:text-gray-600 p-1"
                      :class="{ 'text-primary-600': sortField === 'code' && sortDirection === 'desc' }"
                    >
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
                    <button 
                      @click="sortBy('expiration_date', 'asc')" 
                      class="text-gray-400 hover:text-gray-600 p-1"
                      :class="{ 'text-primary-600': sortField === 'expiration_date' && sortDirection === 'asc' }"
                    >
                      <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7" />
                      </svg>
                    </button>
                    <button 
                      @click="sortBy('expiration_date', 'desc')" 
                      class="text-gray-400 hover:text-gray-600 p-1"
                      :class="{ 'text-primary-600': sortField === 'expiration_date' && sortDirection === 'desc' }"
                    >
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
                    <button 
                      @click="sortBy('amount', 'asc')" 
                      class="text-gray-400 hover:text-gray-600 p-1"
                      :class="{ 'text-primary-600': sortField === 'amount' && sortDirection === 'asc' }"
                    >
                      <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7" />
                      </svg>
                    </button>
                    <button 
                      @click="sortBy('amount', 'desc')" 
                      class="text-gray-400 hover:text-gray-600 p-1"
                      :class="{ 'text-primary-600': sortField === 'amount' && sortDirection === 'desc' }"
                    >
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
                    <button 
                      @click="sortBy('supplier', 'asc')" 
                      class="text-gray-400 hover:text-gray-600 p-1"
                      :class="{ 'text-primary-600': sortField === 'supplier' && sortDirection === 'asc' }"
                    >
                      <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7" />
                      </svg>
                    </button>
                    <button 
                      @click="sortBy('supplier', 'desc')" 
                      class="text-gray-400 hover:text-gray-600 p-1"
                      :class="{ 'text-primary-600': sortField === 'supplier' && sortDirection === 'desc' }"
                    >
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
            <tr v-for="supply in paginatedSupplies" :key="supply.batch_id" class="hover:bg-gray-50">
              <td class="px-6 py-4 whitespace-nowrap text-sm font-mono text-gray-900">{{ supply.batch_id }}</td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div>
                  <div class="font-medium text-gray-900">{{ supply.name }}</div>
                </div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span class="text-gray-700">{{ supply.code }}</span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span :class="getExpirationClass(supply.expiration_date)">
                  {{ formatDate(supply.expiration_date) }}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span :class="getAmountClass(supply.amount)" class="font-medium">{{ supply.amount }}</span>
                <span :class="getAmountClass(supply.amount)" class="text-sm ml-1">unidades</span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span class="text-gray-700">{{ supply.supplier }}</span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="flex space-x-2">
                  <button class="text-primary-600 hover:text-primary-800" @click="viewSupply(supply)">
                    <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                    </svg>
                  </button>
                  <button class="text-warning-600 hover:text-warning-800" @click="editSupply(supply)">
                    <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2.5 2.5 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                    </svg>
                  </button>
                  <button class="text-danger-600 hover:text-danger-800" @click="deleteSupply(supply)">
                    <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                    </svg>
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Modal de edición -->
      <div v-if="showEditModal" class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50">
        <div class="relative top-20 mx-auto p-5 border w-96 shadow-lg rounded-md bg-white">
          <div class="mt-3">
            <div class="flex items-center justify-between mb-4">
              <h3 class="text-lg font-medium text-gray-900">Editar Lote de Insumo</h3>
              <button 
                @click="closeEditModal" 
                class="text-gray-400 hover:text-gray-600"
              >
                <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
            
            <form @submit.prevent="saveEdit" class="space-y-4">
                             <div>
                 <label class="block text-sm font-medium text-gray-700 mb-1">N° de lote</label>
                 <input
                   type="text"
                   v-model="editingSupply.batch_id"
                   class="form-input w-full bg-gray-100 text-gray-600 cursor-not-allowed"
                   readonly
                   disabled
                 />
                 <p class="text-xs text-gray-500 mt-1">El número de lote no se puede editar</p>
               </div>
              
                             <div>
                 <label class="block text-sm font-medium text-gray-700 mb-1">Nombre del Insumo</label>
                 <input
                   type="text"
                   v-model="editingSupply.name"
                   class="form-input w-full bg-gray-100 text-gray-600 cursor-not-allowed"
                   readonly
                   disabled
                 />
                 <p class="text-xs text-gray-500 mt-1">El nombre no se puede editar</p>
               </div>
              
                             <div>
                 <label class="block text-sm font-medium text-gray-700 mb-1">Código Interno</label>
                 <input
                   type="text"
                   v-model="editingSupply.code"
                   class="form-input w-full bg-gray-100 text-gray-600 cursor-not-allowed"
                   readonly
                   disabled
                 />
                 <p class="text-xs text-gray-500 mt-1">El código interno no se puede editar</p>
               </div>
              
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Fecha de Vencimiento</label>
                <input
                  type="date"
                  v-model="editingSupply.expiration_date"
                  class="form-input w-full"
                  required
                />
              </div>
              
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Cantidad</label>
                <input
                  type="number"
                  v-model="editingSupply.amount"
                  class="form-input w-full"
                  min="0"
                  required
                />
                <span class="text-xs text-gray-500">unidades</span>
              </div>
              
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Proveedor</label>
                <input
                  type="text"
                  v-model="editingSupply.supplier"
                  class="form-input w-full"
                  required
                />
              </div>
              
              <div class="flex space-x-3 pt-4">
                <button
                  type="button"
                  @click="closeEditModal"
                  class="btn-secondary flex-1"
                >
                  Cancelar
                </button>
                <button
                  type="submit"
                  class="btn-primary flex-1"
                  :disabled="saving"
                >
                  <span v-if="saving" class="flex items-center justify-center">
                    <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
                      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                      <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
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
      <div v-if="showHistoryModal" class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50">
        <div class="relative top-20 mx-auto p-5 border w-full max-w-4xl shadow-lg rounded-md bg-white">
          <div class="mt-3">
            <div class="flex items-center justify-between mb-4">
              <h3 class="text-lg font-medium text-gray-900">Historial de Movimientos para {{ selectedSupplyForHistory?.name }}</h3>
              <button 
                @click="closeHistoryModal" 
                class="text-gray-400 hover:text-gray-600"
              >
                <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
            
            <!-- Filtros de historial -->
            <div class="flex items-end space-x-4 mb-4">
              <div class="flex-1">
                <label class="block text-sm font-medium text-gray-700 mb-2">Buscar movimiento</label>
                <div class="relative">
                  <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                    <svg class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                    </svg>
                  </div>
                  <input
                    type="text"
                    placeholder="Buscar por fecha, tipo, cantidad o usuario..."
                    class="form-input pl-10 w-full"
                    v-model="historySearchTerm"
                  />
                </div>
              </div>
              <div>
                <button 
                  class="btn-secondary px-4 py-2 h-10" 
                  @click="clearHistorySearch"
                  :disabled="!historySearchTerm"
                >
                  Limpiar
                </button>
              </div>
            </div>

            <!-- Indicador de carga -->
            <div v-if="historyLoading" class="flex justify-center items-center py-8">
              <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
              <span class="ml-2 text-gray-600">Cargando historial...</span>
            </div>

            <!-- Mensaje de error -->
            <div v-else-if="historyError" class="bg-red-50 border border-red-200 rounded-md p-4 mx-4 mb-4">
              <div class="flex">
                <div class="flex-shrink-0">
                  <svg class="h-5 w-5 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                </div>
                <div class="ml-3">
                  <h3 class="text-sm font-medium text-red-800">Error al cargar historial</h3>
                  <div class="mt-2 text-sm text-red-700">{{ historyError }}</div>
                  <div class="mt-4">
                    <button @click="loadHistory" class="btn-secondary text-sm">
                      Reintentar
                    </button>
                  </div>
                </div>
              </div>
            </div>

            <!-- Tabla de datos -->
            <div v-else-if="!historyData.length" class="text-center py-8">
              <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <h3 class="mt-2 text-sm font-medium text-gray-900">No hay historial disponible</h3>
              <p class="mt-1 text-sm text-gray-500">Este lote no tiene movimientos registrados en el historial.</p>
            </div>
            
            <div v-else class="table-container">
              <table class="min-w-full divide-y divide-gray-200">
                <thead class="table-header">
                  <tr>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider whitespace-nowrap min-w-[120px]">
                      <div class="flex items-center justify-between">
                        <span>Fecha</span>
                        <div class="flex flex-col ml-2">
                          <button 
                            @click="sortHistoryBy('date', 'asc')" 
                            class="text-gray-400 hover:text-gray-600 p-1"
                            :class="{ 'text-primary-600': sortHistoryField === 'date' && sortHistoryDirection === 'asc' }"
                          >
                            <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7" />
                            </svg>
                          </button>
                          <button 
                            @click="sortHistoryBy('date', 'desc')" 
                            class="text-gray-400 hover:text-gray-600 p-1"
                            :class="{ 'text-primary-600': sortHistoryField === 'date' && sortHistoryDirection === 'desc' }"
                          >
                            <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                            </svg>
                          </button>
                        </div>
                      </div>
                    </th>
                    <th class="table-header-cell">
                      <div class="flex items-center justify-between">
                        <span>Tipo de Movimiento</span>
                        <div class="flex flex-col ml-2">
                          <button 
                            @click="sortHistoryBy('type', 'asc')" 
                            class="text-gray-400 hover:text-gray-600 p-1"
                            :class="{ 'text-primary-600': sortHistoryField === 'type' && sortHistoryDirection === 'asc' }"
                          >
                            <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7" />
                            </svg>
                          </button>
                          <button 
                            @click="sortHistoryBy('type', 'desc')" 
                            class="text-gray-400 hover:text-gray-600 p-1"
                            :class="{ 'text-primary-600': sortHistoryField === 'type' && sortHistoryDirection === 'desc' }"
                          >
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
                          <button 
                            @click="sortHistoryBy('amount', 'asc')" 
                            class="text-gray-400 hover:text-gray-600 p-1"
                            :class="{ 'text-primary-600': sortHistoryField === 'amount' && sortHistoryDirection === 'asc' }"
                          >
                            <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7" />
                            </svg>
                          </button>
                          <button 
                            @click="sortHistoryBy('amount', 'desc')" 
                            class="text-gray-400 hover:text-gray-600 p-1"
                            :class="{ 'text-primary-600': sortHistoryField === 'amount' && sortHistoryDirection === 'desc' }"
                          >
                            <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                            </svg>
                          </button>
                        </div>
                      </div>
                    </th>
                    <th class="table-header-cell">
                      <div class="flex items-center justify-between">
                        <span>Usuario</span>
                        <div class="flex flex-col ml-2">
                          <button 
                            @click="sortHistoryBy('user', 'asc')" 
                            class="text-gray-400 hover:text-gray-600 p-1"
                            :class="{ 'text-primary-600': sortHistoryField === 'user' && sortHistoryDirection === 'asc' }"
                          >
                            <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7" />
                            </svg>
                          </button>
                          <button 
                            @click="sortHistoryBy('user', 'desc')" 
                            class="text-gray-400 hover:text-gray-600 p-1"
                            :class="{ 'text-primary-600': sortHistoryField === 'user' && sortHistoryDirection === 'desc' }"
                          >
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
                  <tr v-for="movement in paginatedHistory" :key="movement.id || movement.date_time" class="hover:bg-gray-50">
                    <td class="px-6 py-4 whitespace-nowrap text-sm font-mono text-gray-900">{{ formatDate(movement.date_time || movement.date) }}</td>
                    <td class="px-6 py-4 whitespace-nowrap">
                      <span v-if="movement.change_details">{{ movement.change_details }}</span>
                      <span v-else class="badge badge-info">N/A</span>
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap">
                      <span v-if="movement.new_values?.amount !== undefined">{{ movement.new_values.amount }}</span>
                      <span v-else-if="movement.previous_values?.amount !== undefined">{{ movement.previous_values.amount }}</span>
                      <span v-else class="text-red-600 font-semibold">N/A</span>
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap">{{ movement.user_rut || movement.user || movement.user_name || 'N/A' }}</td>
                    <td class="px-6 py-4 whitespace-nowrap">
                      <span class="text-blue-600">●</span>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>

            <!-- Paginación -->
            <div v-if="historyData.length > historyItemsPerPage" class="flex flex-col sm:flex-row items-center justify-between mt-6 space-y-4 sm:space-y-0">
              <div class="text-sm text-gray-700 text-center sm:text-left">
                Mostrando {{ historyStartIndex + 1 }} a {{ historyEndIndex }} de {{ filteredHistory.length }} resultados
              </div>
              <div class="flex items-center space-x-2">
                <button
                  class="btn-secondary px-4 py-2 text-sm min-w-[80px]"
                  :disabled="historyCurrentPage === 1"
                  @click="historyCurrentPage--"
                >
                  Anterior
                </button>
                <span class="px-3 py-2 text-sm text-gray-700 bg-gray-100 rounded-md min-w-[100px] text-center">
                  Página {{ historyCurrentPage }} de {{ historyTotalPages }}
                </span>
                <button
                  class="btn-secondary px-4 py-2 text-sm min-w-[80px]"
                  :disabled="historyCurrentPage === historyTotalPages"
                  @click="historyCurrentPage++"
                >
                  Siguiente
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Modal de historial global -->
      <div v-if="showGlobalHistoryModal" class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50">
        <div class="relative top-20 mx-auto p-5 border w-full max-w-4xl shadow-lg rounded-md bg-white">
          <div class="mt-3">
            <div class="flex items-center justify-between mb-4">
              <div>
                <h3 class="text-lg font-medium text-gray-900">Historial de Movimientos de Lotes</h3>
                <p v-if="hasActiveFilters" class="text-sm text-gray-600 mt-1">
                  {{ filteredGlobalHistory.length }} resultado{{ filteredGlobalHistory.length !== 1 ? 's' : '' }} encontrado{{ filteredGlobalHistory.length !== 1 ? 's' : '' }}
                </p>
              </div>
              <button 
                @click="closeGlobalHistoryModal" 
                class="text-gray-400 hover:text-gray-600"
              >
                <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
            
                                      <!-- Filtros de historial global -->
              <div class="grid grid-cols-1 md:grid-cols-4 gap-4 mb-4 items-end">
                              <!-- Filtro por tipo de cambio -->
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2">Tipo de Cambio</label>
                  <select 
                    v-model="globalHistoryChangeTypeFilter" 
                    class="form-select w-full"
                    @change="applyFilters"
                  >
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
                 <label class="block text-sm font-medium text-gray-700 mb-2">N° de Lote</label>
                 <div class="relative">
                   <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                     <svg class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                       <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                     </svg>
                   </div>
                   <input
                     type="text"
                     placeholder="N° de lote..."
                     class="form-input pl-10 w-full"
                     v-model="globalHistoryBatchFilter"
                     @input="applyFilters"
                   />
                 </div>
               </div>

               <!-- Buscador por usuario -->
               <div>
                 <label class="block text-sm font-medium text-gray-700 mb-2">Usuario</label>
                 <div class="relative">
                   <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                     <svg class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                       <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                     </svg>
                   </div>
                   <input
                     type="text"
                     placeholder="RUT o nombre..."
                     class="form-input pl-10 w-full"
                     v-model="globalHistoryUserFilter"
                     @input="applyFilters"
                   />
                 </div>
               </div>

                               <!-- Botón de limpiar filtros -->
                <div class="flex items-end h-full">
                  <button 
                    class="btn-secondary px-4 py-2 h-10 w-full" 
                    @click="clearGlobalHistoryFilters"
                    :disabled="!hasActiveFilters"
                  >
                    <svg v-if="hasActiveFilters" class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                    </svg>
                    Limpiar Filtros
                  </button>
                </div>
             </div>

            <!-- Indicador de carga -->
            <div v-if="globalHistoryLoading" class="flex justify-center items-center py-8">
              <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
              <span class="ml-2 text-gray-600">Cargando historial...</span>
            </div>

            <!-- Mensaje de error -->
            <div v-else-if="globalHistoryError" class="bg-red-50 border border-red-200 rounded-md p-4 mx-4 mb-4">
              <div class="flex">
                <div class="flex-shrink-0">
                  <svg class="h-5 w-5 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                </div>
                <div class="ml-3">
                  <h3 class="text-sm font-medium text-red-800">Error al cargar historial</h3>
                  <div class="mt-2 text-sm text-red-700">{{ globalHistoryError }}</div>
                  <div class="mt-4">
                    <button @click="loadGlobalHistory" class="btn-secondary text-sm">
                      Reintentar
                    </button>
                  </div>
                </div>
              </div>
            </div>

            <!-- Tabla de datos -->
            <div v-else-if="!globalHistoryData.length" class="text-center py-8">
              <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <h3 class="mt-2 text-sm font-medium text-gray-900">No hay historial disponible</h3>
              <p class="mt-1 text-sm text-gray-500">No hay movimientos registrados en el historial.</p>
            </div>
            
            <!-- Mensaje de búsqueda sin resultados -->
            <div v-else-if="hasActiveFilters && !filteredGlobalHistory.length" class="text-center py-8">
              <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
              <h3 class="mt-2 text-sm font-medium text-gray-900">No se encontraron resultados</h3>
              <p class="mt-1 text-sm text-gray-500">
                No hay movimientos que coincidan con los filtros aplicados. 
                Intenta con otros criterios de búsqueda.
              </p>
            </div>
            
            <div v-else-if="filteredGlobalHistory.length > 0" class="table-container">
              <table class="min-w-full divide-y divide-gray-200">
                <thead class="table-header">
                  <tr>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider whitespace-nowrap min-w-[120px]">
                      <div class="flex items-center justify-between">
                        <span>Fecha</span>
                        <div class="flex flex-col ml-2">
                          <button 
                            @click="sortGlobalHistoryBy('date', 'asc')" 
                            class="text-gray-400 hover:text-gray-600 p-1"
                            :class="{ 'text-primary-600': globalHistorySortField === 'date' && globalHistorySortDirection === 'asc' }"
                          >
                            <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7" />
                            </svg>
                          </button>
                          <button 
                            @click="sortGlobalHistoryBy('date', 'desc')" 
                            class="text-gray-400 hover:text-gray-600 p-1"
                            :class="{ 'text-primary-600': globalHistorySortField === 'date' && globalHistorySortDirection === 'desc' }"
                          >
                            <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                            </svg>
                          </button>
                        </div>
                      </div>
                    </th>
                    <th class="table-header-cell">
                      <div class="flex items-center justify-between">
                        <span>N° de Lote</span>
                        <div class="flex flex-col ml-2">
                          <button 
                            @click="sortGlobalHistoryBy('batch_id', 'asc')" 
                            class="text-gray-400 hover:text-gray-600 p-1"
                            :class="{ 'text-primary-600': globalHistorySortField === 'batch_id' && globalHistorySortDirection === 'asc' }"
                          >
                            <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7" />
                            </svg>
                          </button>
                          <button 
                            @click="sortGlobalHistoryBy('batch_id', 'desc')" 
                            class="text-gray-400 hover:text-gray-600 p-1"
                            :class="{ 'text-primary-600': globalHistorySortField === 'batch_id' && globalHistorySortDirection === 'desc' }"
                          >
                            <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                            </svg>
                          </button>
                        </div>
                      </div>
                    </th>
                    <th class="table-header-cell">
                      <div class="flex items-center justify-between">
                        <span>Detalles del Cambio</span>
                        <div class="flex flex-col ml-2">
                          <button 
                            @click="sortGlobalHistoryBy('change_details', 'asc')" 
                            class="text-gray-400 hover:text-gray-600 p-1"
                            :class="{ 'text-primary-600': globalHistorySortField === 'change_details' && globalHistorySortDirection === 'asc' }"
                          >
                            <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7" />
                            </svg>
                          </button>
                          <button 
                            @click="sortGlobalHistoryBy('change_details', 'desc')" 
                            class="text-gray-400 hover:text-gray-600 p-1"
                            :class="{ 'text-primary-600': globalHistorySortField === 'change_details' && globalHistorySortDirection === 'desc' }"
                          >
                            <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                            </svg>
                          </button>
                        </div>
                      </div>
                    </th>
                    <th class="table-header-cell">
                      <div class="flex items-center justify-between">
                        <span>Valores Anteriores</span>
                        <div class="flex flex-col ml-2">
                          <button 
                            @click="sortGlobalHistoryBy('previous_values', 'asc')" 
                            class="text-gray-400 hover:text-gray-600 p-1"
                            :class="{ 'text-primary-600': globalHistorySortField === 'previous_values' && globalHistorySortDirection === 'asc' }"
                          >
                            <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7" />
                            </svg>
                          </button>
                          <button 
                            @click="sortGlobalHistoryBy('previous_values', 'desc')" 
                            class="text-gray-400 hover:text-gray-600 p-1"
                            :class="{ 'text-primary-600': globalHistorySortField === 'previous_values' && globalHistorySortDirection === 'desc' }"
                          >
                            <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                            </svg>
                          </button>
                        </div>
                      </div>
                    </th>
                    <th class="table-header-cell">
                      <div class="flex items-center justify-between">
                        <span>Valores Nuevos</span>
                        <div class="flex flex-col ml-2">
                          <button 
                            @click="sortGlobalHistoryBy('new_values', 'asc')" 
                            class="text-gray-400 hover:text-gray-600 p-1"
                            :class="{ 'text-primary-600': globalHistorySortField === 'new_values' && globalHistorySortDirection === 'asc' }"
                          >
                            <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7" />
                            </svg>
                          </button>
                          <button 
                            @click="sortGlobalHistoryBy('new_values', 'desc')" 
                            class="text-gray-400 hover:text-gray-600 p-1"
                            :class="{ 'text-primary-600': globalHistorySortField === 'new_values' && globalHistorySortDirection === 'desc' }"
                          >
                            <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                            </svg>
                          </button>
                        </div>
                      </div>
                    </th>

                    <th class="table-header-cell">
                      <div class="flex items-center justify-between">
                        <span>RUT Usuario</span>
                        <div class="flex flex-col ml-2">
                          <button 
                            @click="sortGlobalHistoryBy('user_rut', 'asc')" 
                            class="text-gray-400 hover:text-gray-600 p-1"
                            :class="{ 'text-primary-600': globalHistorySortField === 'user_rut' && globalHistorySortDirection === 'asc' }"
                          >
                            <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7" />
                            </svg>
                          </button>
                          <button 
                            @click="sortGlobalHistoryBy('user_rut', 'desc')" 
                            class="text-gray-400 hover:text-gray-600 p-1"
                            :class="{ 'text-primary-600': globalHistorySortField === 'user_rut' && globalHistorySortDirection === 'desc' }"
                          >
                            <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                            </svg>
                          </button>
                        </div>
                      </div>
                    </th>
                    <th class="table-header-cell">
                      <div class="flex items-center justify-between">
                        <span>Nombre Usuario</span>
                        <div class="flex flex-col ml-2">
                          <button 
                            @click="sortGlobalHistoryBy('user_name', 'asc')" 
                            class="text-gray-400 hover:text-gray-600 p-1"
                            :class="{ 'text-primary-600': globalHistorySortField === 'user_name' && globalHistorySortDirection === 'asc' }"
                          >
                            <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 15l7-7 7 7" />
                            </svg>
                          </button>
                          <button 
                            @click="sortGlobalHistoryBy('user_name', 'desc')" 
                            class="text-gray-400 hover:text-gray-600 p-1"
                            :class="{ 'text-primary-600': globalHistorySortField === 'user_name' && globalHistorySortDirection === 'desc' }"
                          >
                            <svg class="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                            </svg>
                          </button>
                        </div>
                      </div>
                    </th>
                  </tr>
                </thead>
                <tbody class="bg-white divide-y divide-gray-200">
                  <tr v-for="movement in paginatedGlobalHistory" :key="movement.id" class="hover:bg-gray-50">
                    <td class="px-6 py-4 whitespace-nowrap text-sm font-mono text-gray-900">{{ formatDate(movement.date_time) }}</td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm font-mono text-gray-900">{{ movement.batch_id }}</td>
                    <td class="px-6 py-4 whitespace-nowrap">
                      <div class="font-medium text-gray-900">{{ movement.change_details }}</div>
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-700">
                      <div v-if="movement.previous_values" class="text-xs">
                        <div v-if="movement.previous_values.amount">Cantidad: {{ movement.previous_values.amount }}</div>
                        <div v-if="movement.previous_values.supplier">Proveedor: {{ movement.previous_values.supplier }}</div>
                        <div v-if="movement.previous_values.expiration_date">Vencimiento: {{ formatDate(movement.previous_values.expiration_date) }}</div>
                      </div>
                      <span v-else class="text-gray-400">N/A</span>
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-700">
                      <div v-if="movement.new_values" class="text-xs">
                        <div v-if="movement.new_values.amount">Cantidad: {{ movement.new_values.amount }}</div>
                        <div v-if="movement.new_values.supplier">Proveedor: {{ movement.new_values.supplier }}</div>
                        <div v-if="movement.new_values.expiration_date">Vencimiento: {{ formatDate(movement.new_values.expiration_date) }}</div>
                      </div>
                      <span v-else class="text-gray-400">N/A</span>
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-700">{{ movement.user_rut }}</td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-700">{{ movement.user_name || 'N/A' }}</td>
                  </tr>
                </tbody>
              </table>
            </div>

            <!-- Paginación -->
            <div v-if="globalHistoryData.length > globalHistoryItemsPerPage" class="flex flex-col sm:flex-row items-center justify-between mt-6 space-y-4 sm:space-y-0">
              <div class="text-sm text-gray-700 text-center sm:text-left">
                Mostrando {{ globalHistoryStartIndex + 1 }} a {{ globalHistoryEndIndex }} de {{ filteredGlobalHistory.length }} resultados
              </div>
              <div class="flex items-center space-x-2">
                <button
                  class="btn-secondary px-4 py-2 text-sm min-w-[80px]"
                  :disabled="globalHistoryCurrentPage === 1"
                  @click="globalHistoryCurrentPage--"
                >
                  Anterior
                </button>
                <span class="px-3 py-2 text-sm text-gray-700 bg-gray-100 rounded-md min-w-[100px] text-center">
                  Página {{ globalHistoryCurrentPage }} de {{ globalHistoryTotalPages }}
                </span>
                <button
                  class="btn-secondary px-4 py-2 text-sm min-w-[80px]"
                  :disabled="globalHistoryCurrentPage === globalHistoryTotalPages"
                  @click="globalHistoryCurrentPage++"
                >
                  Siguiente
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Sistema de notificaciones -->
      <div v-if="notification.show" class="fixed top-4 right-4 z-50">
        <div 
          class="px-6 py-4 rounded-lg shadow-lg text-white max-w-sm"
          :class="notification.type === 'success' ? 'bg-green-500' : 'bg-red-500'"
        >
          <div class="flex items-center justify-between">
            <div class="flex items-center">
              <svg v-if="notification.type === 'success'" class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
              <svg v-else class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
              <span class="font-medium">{{ notification.message }}</span>
            </div>
            <button @click="hideNotification" class="ml-4 text-white hover:text-gray-200">
              <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
        </div>
      </div>

      <!-- Paginación -->
      <div class="flex flex-col sm:flex-row items-center justify-between mt-6 space-y-4 sm:space-y-0">
        <div class="text-sm text-gray-700 text-center sm:text-left">
          Mostrando {{ startIndex + 1 }} a {{ endIndex }} de {{ filteredSupplies.length }} resultados
        </div>
        <div class="flex items-center space-x-2">
          <button
            class="btn-secondary px-4 py-2 text-sm min-w-[80px]"
            :disabled="currentPage === 1"
            @click="currentPage--"
          >
            Anterior
          </button>
          <span class="px-3 py-2 text-sm text-gray-700 bg-gray-100 rounded-md min-w-[100px] text-center">
            Página {{ currentPage }} de {{ totalPages }}
          </span>
          <button
            class="btn-secondary px-4 py-2 text-sm min-w-[80px]"
            :disabled="currentPage === totalPages"
            @click="currentPage++"
          >
            Siguiente
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router';
import { format } from 'date-fns'
import { es } from 'date-fns/locale'
import inventoryService from '@/services/inventoryService'

const route = useRoute()

// Estado reactivo
const supplies = ref([])
const loading = ref(false)
const error = ref(null)
const searchTerm = ref('')
const sortField = ref('name')
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
const globalHistoryItemsPerPage = 10
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

// Estado de notificación
const notification = ref({
  show: false,
  message: '',
  type: 'info' // 'success', 'error', 'info'
})

// Computed properties
const filteredSupplies = computed(() => {
  let filtered = [...supplies.value]
  
  if (searchTerm.value) {
    filtered = filtered.filter(supply => 
      supply.code.toString().includes(searchTerm.value) ||
      supply.name.toLowerCase().includes(searchTerm.value.toLowerCase()) ||
      supply.supplier.toLowerCase().includes(searchTerm.value.toLowerCase())
    )
  }
  
  // Ordenamiento
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

const historyTotalPages = computed(() => Math.ceil(filteredHistory.value.length / historyItemsPerPage))

const historyStartIndex = computed(() => (historyCurrentPage.value - 1) * historyItemsPerPage)
const historyEndIndex = computed(() => Math.min(historyStartIndex.value + historyItemsPerPage, filteredHistory.value.length))

const paginatedHistory = computed(() => {
  return filteredHistory.value.slice(historyStartIndex.value, historyEndIndex.value)
})

// Computed properties para historial global
const filteredGlobalHistory = computed(() => {
  let filtered = [...globalHistoryData.value]

  // Filtro por tipo de cambio
  if (globalHistoryChangeTypeFilter.value) {
         console.log('Filtrando por tipo de cambio:', globalHistoryChangeTypeFilter.value)
    
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
    
    console.log('Resultados filtrados:', filtered.length)
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

const globalHistoryTotalPages = computed(() => Math.ceil(filteredGlobalHistory.value.length / globalHistoryItemsPerPage))

const globalHistoryStartIndex = computed(() => (globalHistoryCurrentPage.value - 1) * globalHistoryItemsPerPage)
const globalHistoryEndIndex = computed(() => Math.min(globalHistoryStartIndex.value + globalHistoryItemsPerPage, filteredGlobalHistory.value.length))

const paginatedGlobalHistory = computed(() => {
  return filteredGlobalHistory.value.slice(globalHistoryStartIndex.value, globalHistoryEndIndex.value)
})

// Verificar si hay filtros activos
const hasActiveFilters = computed(() => {
  return globalHistoryChangeTypeFilter.value || 
         globalHistoryBatchFilter.value || 
         globalHistoryUserFilter.value
})


const clearSearch = () => {
  searchTerm.value = ''
  sortField.value = 'name'
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
  switch (status) {
    case 'active':
      return 'badge badge-success'
    case 'expired':
      return 'badge badge-danger'
    case 'low_stock':
      return 'badge badge-warning'
    default:
      return 'badge badge-info'
  }
}

const editSupply = (supply) => {
  // Preparar los datos para edición
  editingSupply.value = {
    ...supply,
    // Convertir la fecha al formato requerido por el input date (YYYY-MM-DD)
    expiration_date: supply.expiration_date ? new Date(supply.expiration_date).toISOString().split('T')[0] : ''
  }
  showEditModal.value = true
}

const closeEditModal = () => {
  showEditModal.value = false
  editingSupply.value = {}
  saving.value = false
}

const saveEdit = async () => {
  if (!editingSupply.value.batch_id) {
    alert('Error: ID de lote no encontrado')
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

    const batchData = {
      expiration_date: formattedDate,
      amount: parseInt(editingSupply.value.amount),
      supplier: editingSupply.value.supplier,
      store_id: 1 // Por defecto, se puede hacer configurable después
    }

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
    showNotification('Lote actualizado exitosamente', 'success')
    
  } catch (error) {
    console.error('Error al actualizar lote:', error)
    showNotification('Error al actualizar el lote: ' + (error.response?.data?.error || error.message), 'error')
  } finally {
    saving.value = false
  }
}

const deleteSupply = async (supply) => {
  if (confirm(`¿Está seguro de que desea eliminar el lote de ${supply.name}?`)) {
    try {
      await inventoryService.deleteBatch(supply.batch_id)
      showNotification('Lote eliminado exitosamente', 'success')
      await loadInventory()
    } catch (error) {
      console.error('Error al eliminar lote:', error)
      showNotification('Error al eliminar el lote: ' + (error.response?.data?.error || error.message), 'error')
    }
  }
}

const scanQR = () => {
  console.log('Escanear código QR')
  // TODO: Implementar escáner QR
}

const exportInventory = () => {
  console.log('Exportar inventario')
  // TODO: Implementar exportación
}

const showNotification = (message, type = 'info') => {
  notification.value = {
    show: true,
    message: message,
    type: type
  }
  setTimeout(() => {
    notification.value.show = false
  }, 3000) // 3 segundos de duración
}

const hideNotification = () => {
  notification.value.show = false
}



// Métodos para historial global
const openGlobalHistoryModal = () => {
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

const debugGlobalHistoryData = () => {
  console.log('=== DEBUG: Datos del historial global ===')
  console.log('Total de registros:', globalHistoryData.value.length)
  console.log('Filtros activos:', {
    changeType: globalHistoryChangeTypeFilter.value,
    batch: globalHistoryBatchFilter.value,
    user: globalHistoryUserFilter.value
  })
  
  // Mostrar los primeros 5 registros con sus detalles
  globalHistoryData.value.slice(0, 5).forEach((movement, index) => {
    console.log(`Registro ${index + 1}:`, {
      batch_id: movement.batch_id,
      change_details: movement.change_details,
      user_rut: movement.user_rut,
      user_name: movement.user_name,
      date_time: movement.date_time
    })
  })
  
  // Mostrar todos los tipos de change_details únicos
  const uniqueChangeDetails = [...new Set(globalHistoryData.value.map(m => m.change_details))]
  console.log('Tipos de cambio únicos disponibles:', uniqueChangeDetails)
  
  // Mostrar estadísticas de filtrado
  console.log('Registros filtrados:', filteredGlobalHistory.value.length)
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

const viewMovementDetails = (movement) => {
  // TODO: Implementar vista de detalles del movimiento
  console.log('Ver detalles del movimiento:', movement)
}

// Métodos
const loadInventory = async () => {
  loading.value = true
  error.value = null
  
  try {
  const data = await inventoryService.getInventory()
  supplies.value = Array.isArray(data) ? data : (data.inventory_items || [])
  } catch (err) {
    error.value = 'Error al cargar el inventario: ' + err.message
    console.error('Error al cargar inventario:', err)
  } finally {
    loading.value = false
  }
}

// Lifecycle
onMounted(() => {
  // Si viene con un término de búsqueda desde Home, aplicarlo
  if (route.query.search) {
    searchTerm.value = route.query.search
  }
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
    font-size: 0.875rem;
  }
  
  .table-header th:first-child,
  .table-header th:nth-child(2) {
    min-width: 100px;
  }
  
  .table-header th:nth-child(3),
  .table-header th:nth-child(4),
  .table-header th:nth-child(5),
  .table-header th:nth-child(6) {
    min-width: 80px;
  }
  
  .table-header th:last-child {
    min-width: 80px;
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
  @apply block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-primary-500 focus:border-primary-500 sm:text-sm;
}
</style> 
