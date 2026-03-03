<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="bg-white rounded-lg shadow-sm border p-6">
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h2 class="text-2xl font-semibold text-gray-900">Historial de Transferencias</h2>
          <p class="text-gray-600 mt-1">Consulta y gestiona las transferencias realizadas desde el escáner QR</p>
          <p v-if="!loading" class="text-sm text-gray-500 mt-1">Total: {{ sortedTransfers.length }} transferencias</p>
        </div>
        <div class="flex flex-col sm:flex-row gap-2">
          <button 
            @click="exportToExcel" 
            :disabled="loading || sortedTransfers.length === 0"
            class="btn-secondary flex items-center justify-center"
            :class="{ 'opacity-50 cursor-not-allowed': loading || sortedTransfers.length === 0 }"
          >
            <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
            Exportar a Excel
          </button>
          <router-link to="/qr" class="btn-primary flex items-center justify-center">
            <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h.01M5 8h2a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v2a2 2 0 002 2zm0 0h2a2 2 0 012 2v2a2 2 0 01-2 2H5a2 2 0 01-2-2v-2a2 2 0 012-2z" />
            </svg>
            Ir a Escáner QR
          </router-link>
        </div>
      </div>
    </div>

    <!-- Filtros -->
    <div class="card space-y-4">
      <!-- Buscador -->
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-2">Buscar por código</label>
        <div class="relative">
          <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
            <svg class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
          </div>
          <input type="text" v-model="filters.code" placeholder="Código de transferencia" class="form-input pl-10 w-full" />
        </div>
      </div>
      <!-- Filtros adicionales + Limpiar -->
      <div class="flex flex-wrap gap-4 items-end">
        <div class="flex-1 min-w-[160px]">
          <label class="block text-sm font-medium text-gray-700 mb-2">Estado</label>
          <select v-model="filters.status" class="form-input">
            <option value="">Todos</option>
            <option value="pending">Pendiente</option>
            <option value="en_transito">En Tránsito</option>
            <option value="completed">Completado</option>
            <option value="cancelled">Cancelado</option>
          </select>
        </div>
        <div class="flex-1 min-w-[140px]">
          <label class="block text-sm font-medium text-gray-700 mb-2">Desde</label>
          <input type="date" v-model="filters.from_date" class="form-input" />
        </div>
        <div class="flex-1 min-w-[140px]">
          <label class="block text-sm font-medium text-gray-700 mb-2">Hasta</label>
          <input type="date" v-model="filters.to_date" class="form-input" />
        </div>
        <div class="flex-shrink-0">
          <label class="block text-sm font-medium text-gray-700 mb-2 sm:invisible">Acción</label>
          <button class="btn-secondary px-4 py-2 h-10 w-full sm:w-auto" @click="clearFilters"
            :disabled="!filters.code && !filters.status && !filters.from_date && !filters.to_date"
          >
            <svg class="h-4 w-4 mr-2 inline" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
            Limpiar
          </button>
        </div>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="card">
      <div class="flex justify-center items-center py-12">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-brand-blue-dark"></div>
        <span class="ml-3 text-gray-600">Cargando transferencias...</span>
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
            <h3 class="text-sm font-medium text-red-800">Error al cargar transferencias</h3>
            <div class="mt-2 text-sm text-red-700">{{ error }}</div>
            <button @click="loadTransfers" class="btn-secondary mt-4 text-sm">Reintentar</button>
          </div>
        </div>
      </div>
    </div>

    <!-- Tabla de transferencias -->
    <div v-else class="card overflow-hidden">
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th scope="col" @click="sortBy('transfer_code')"
                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100 transition-colors select-none">
                <div class="flex items-center space-x-1">
                  <span>Código</span>
                  <span class="flex flex-col -space-y-1">
                    <svg class="h-3 w-3 transition-colors" 
                      :class="sortKey === 'transfer_code' && sortOrder === 'asc' ? 'text-brand-blue-dark' : 'text-gray-300'" 
                      fill="currentColor" viewBox="0 0 20 20">
                      <path d="M5.293 9.707a1 1 0 010-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 01-1.414 1.414L11 7.414V15a1 1 0 11-2 0V7.414L6.707 9.707a1 1 0 01-1.414 0z"/>
                    </svg>
                    <svg class="h-3 w-3 transition-colors" 
                      :class="sortKey === 'transfer_code' && sortOrder === 'desc' ? 'text-blue-600' : 'text-gray-300'" 
                      fill="currentColor" viewBox="0 0 20 20">
                      <path d="M14.707 10.293a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 111.414-1.414L9 12.586V5a1 1 0 012 0v7.586l2.293-2.293a1 1 0 011.414 0z"/>
                    </svg>
                  </span>
                </div>
              </th>
              <th scope="col" @click="sortBy('origin_name')"
                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100 transition-colors select-none">
                <div class="flex items-center space-x-1">
                  <span>Origen</span>
                  <span class="flex flex-col -space-y-1">
                    <svg class="h-3 w-3 transition-colors" 
                      :class="sortKey === 'origin_name' && sortOrder === 'asc' ? 'text-blue-600' : 'text-gray-300'" 
                      fill="currentColor" viewBox="0 0 20 20">
                      <path d="M5.293 9.707a1 1 0 010-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 01-1.414 1.414L11 7.414V15a1 1 0 11-2 0V7.414L6.707 9.707a1 1 0 01-1.414 0z"/>
                    </svg>
                    <svg class="h-3 w-3 transition-colors" 
                      :class="sortKey === 'origin_name' && sortOrder === 'desc' ? 'text-blue-600' : 'text-gray-300'" 
                      fill="currentColor" viewBox="0 0 20 20">
                      <path d="M14.707 10.293a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 111.414-1.414L9 12.586V5a1 1 0 012 0v7.586l2.293-2.293a1 1 0 011.414 0z"/>
                    </svg>
                  </span>
                </div>
              </th>
              <th scope="col" @click="sortBy('destination_name')"
                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100 transition-colors select-none">
                <div class="flex items-center space-x-1">
                  <span>Destino</span>
                  <span class="flex flex-col -space-y-1">
                    <svg class="h-3 w-3 transition-colors" 
                      :class="sortKey === 'destination_name' && sortOrder === 'asc' ? 'text-blue-600' : 'text-gray-300'" 
                      fill="currentColor" viewBox="0 0 20 20">
                      <path d="M5.293 9.707a1 1 0 010-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 01-1.414 1.414L11 7.414V15a1 1 0 11-2 0V7.414L6.707 9.707a1 1 0 01-1.414 0z"/>
                    </svg>
                    <svg class="h-3 w-3 transition-colors" 
                      :class="sortKey === 'destination_name' && sortOrder === 'desc' ? 'text-blue-600' : 'text-gray-300'" 
                      fill="currentColor" viewBox="0 0 20 20">
                      <path d="M14.707 10.293a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 111.414-1.414L9 12.586V5a1 1 0 012 0v7.586l2.293-2.293a1 1 0 011.414 0z"/>
                    </svg>
                  </span>
                </div>
              </th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Cantidad
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
              <th scope="col" @click="sortBy('created_at')"
                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100 transition-colors select-none">
                <div class="flex items-center space-x-1">
                  <span>Fecha</span>
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
              <th scope="col" class="px-4 py-3 text-center text-xs font-medium text-gray-500 uppercase tracking-wider w-16 sticky right-0 bg-gray-50 z-10 shadow-[-2px_0_4px_rgba(0,0,0,0.05)]">
                Acciones
              </th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-for="transfer in paginatedTransfers" :key="transfer.transfer_code || transfer.id" class="hover:bg-gray-50 transition-colors duration-150">
              <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900 font-mono">
                {{ transfer.transfer_code || transfer.code || 'N/A' }}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ getOriginName(transfer) }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ getDestinationName(transfer) }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">1</td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span class="px-2 py-1 text-xs font-semibold rounded-full" :class="getStatusClass(transfer.status)">
                  {{ getStatusLabel(transfer.status) }}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ formatDate(transfer.created_at) }}</td>
              <td class="px-4 py-4 text-center text-sm font-medium w-16 sticky right-0 bg-white z-10 shadow-[-2px_0_4px_rgba(0,0,0,0.05)]">
                <button @click="viewDetails(transfer)" 
                  class="text-blue-600 hover:text-blue-800 hover:bg-blue-50 p-1.5 rounded transition-colors"
                  title="Ver detalles">
                  <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                  </svg>
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Paginación -->
    <div v-if="!loading && sortedTransfers.length > 0" class="card">
      <div class="flex flex-col sm:flex-row items-center justify-between gap-4">
        <div class="text-sm text-gray-700 text-center sm:text-left">
          Mostrando {{ startIndex + 1 }} a {{ endIndex }} de {{ sortedTransfers.length }} transferencias
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

    <!-- Estado vacío -->
    <div v-if="!loading && !error && transfers.length === 0" class="card">
      <div class="text-center py-12">
        <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4" />
        </svg>
        <h3 class="mt-2 text-sm font-medium text-gray-900">No hay transferencias</h3>
        <p class="mt-1 text-sm text-gray-500">Las transferencias se crean desde el Escáner QR al transferir insumos entre ubicaciones.</p>
        <div class="mt-6">
          <router-link to="/qr" class="btn-primary inline-flex items-center">
            <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h.01M5 8h2a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v2a2 2 0 002 2zm0 0h2a2 2 0 012 2v2a2 2 0 01-2 2H5a2 2 0 01-2-2v-2a2 2 0 012-2z" />
            </svg>
            Ir a Escáner QR
          </router-link>
        </div>
      </div>
    </div>

    <!-- Modal de detalles de transferencia -->
    <Teleport to="body">
      <div v-if="showDetailsModal && selectedTransfer" class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50 flex items-start justify-center pt-20 pb-8" @click.self="closeDetailsModal">
        <div class="w-full max-w-4xl p-6 border shadow-lg rounded-lg bg-white my-auto">
        <!-- Header del modal -->
        <div class="flex justify-between items-start mb-6">
          <div class="flex items-center space-x-4">
            <div class="h-12 w-12 rounded-full bg-gradient-to-br from-brand-blue-dark to-brand-blue-medium flex items-center justify-center">
              <svg class="h-7 w-7 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4" />
              </svg>
            </div>
            <div>
              <h3 class="text-2xl font-bold text-gray-900">Transferencia</h3>
              <p class="text-sm text-gray-500 mt-1">{{ selectedTransfer.transfer_code }}</p>
            </div>
          </div>
          <button @click="closeDetailsModal" class="text-gray-400 hover:text-gray-600 transition-colors">
            <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <!-- Estado destacado -->
        <div class="mb-6 p-4 rounded-lg" :class="getStatusBgClass(selectedTransfer.status)">
          <div class="flex items-center justify-between">
            <div class="flex items-center space-x-3">
              <div class="h-10 w-10 rounded-full bg-white/80 flex items-center justify-center">
                <svg v-if="selectedTransfer.status === 'recibido' || selectedTransfer.status === 'completed'" class="h-6 w-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                </svg>
                <svg v-else-if="selectedTransfer.status === 'cancelado' || selectedTransfer.status === 'rechazado'" class="h-6 w-6 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
                <svg v-else class="h-6 w-6 text-yellow-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </div>
              <div>
                <p class="text-sm font-medium opacity-80">Estado actual</p>
                <p class="text-lg font-bold">{{ getStatusLabel(selectedTransfer.status) }}</p>
              </div>
            </div>
          </div>
        </div>

        <!-- Tabs Navigation -->
        <div class="border-b border-gray-200 mb-6">
          <nav class="flex space-x-4" aria-label="Tabs">
            <button
              @click="activeTab = 'general'"
              :class="[
                'pb-3 px-1 border-b-2 font-medium text-sm transition-colors',
                activeTab === 'general'
                  ? 'border-blue-500 text-blue-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
              ]"
            >
              <span class="flex items-center space-x-2">
                <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                <span>General</span>
              </span>
            </button>
            <button
              @click="activeTab = 'ubicacion'"
              :class="[
                'pb-3 px-1 border-b-2 font-medium text-sm transition-colors',
                activeTab === 'ubicacion'
                  ? 'border-blue-500 text-blue-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
              ]"
            >
              <span class="flex items-center space-x-2">
                <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
                </svg>
                <span>Ubicación</span>
              </span>
            </button>
            <button
              @click="activeTab = 'personal'"
              :class="[
                'pb-3 px-1 border-b-2 font-medium text-sm transition-colors',
                activeTab === 'personal'
                  ? 'border-blue-500 text-blue-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
              ]"
            >
              <span class="flex items-center space-x-2">
                <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
                </svg>
                <span>Personal</span>
              </span>
            </button>
            <button
              @click="activeTab = 'adicional'"
              :class="[
                'pb-3 px-1 border-b-2 font-medium text-sm transition-colors',
                activeTab === 'adicional'
                  ? 'border-blue-500 text-blue-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
              ]"
            >
              <span class="flex items-center space-x-2">
                <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                </svg>
                <span>Adicional</span>
              </span>
            </button>
          </nav>
        </div>

        <!-- Tab Content -->
        <div class="min-h-[300px]">
          <!-- Tab: General -->
          <div v-show="activeTab === 'general'" class="space-y-6 animate-fade-in">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div class="bg-gradient-to-br from-brand-blue-light from-opacity-20 to-brand-blue-light from-opacity-40 p-5 rounded-xl border border-brand-blue-medium border-opacity-30">
                <div class="flex items-center space-x-3 mb-2">
                  <svg class="h-5 w-5 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h.01M5 8h2a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v2a2 2 0 002 2zm0 0h2a2 2 0 012 2v2a2 2 0 01-2 2H5a2 2 0 01-2-2v-2a2 2 0 012-2z" />
                  </svg>
                  <p class="text-sm font-semibold text-blue-900">Código QR</p>
                </div>
                <p class="text-xl font-bold text-blue-700">{{ selectedTransfer.qr_code }}</p>
              </div>

              <div class="bg-gradient-to-br from-purple-50 to-purple-100 p-5 rounded-xl border border-purple-200">
                <div class="flex items-center space-x-3 mb-2">
                  <svg class="h-5 w-5 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
                  </svg>
                  <p class="text-sm font-semibold text-purple-900">ID Insumo</p>
                </div>
                <p class="text-xl font-bold text-purple-700">#{{ selectedTransfer.medical_supply_id }}</p>
              </div>
            </div>

            <div class="bg-gray-50 p-5 rounded-xl">
              <h4 class="text-sm font-semibold text-gray-700 mb-4 flex items-center">
                <svg class="h-5 w-5 mr-2 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                </svg>
                Fechas importantes
              </h4>
              <div class="grid grid-cols-2 gap-4">
                <div class="bg-white p-4 rounded-lg border border-gray-200">
                  <p class="text-xs text-gray-500 mb-1">Fecha de envío</p>
                  <p class="text-sm font-semibold text-gray-900">{{ formatDate(selectedTransfer.send_date) }}</p>
                </div>
                <div class="bg-white p-4 rounded-lg border" :class="selectedTransfer.receive_date ? 'border-green-200' : 'border-gray-200'">
                  <p class="text-xs mb-1" :class="selectedTransfer.receive_date ? 'text-green-600' : 'text-gray-500'">Fecha de recepción</p>
                  <p class="text-sm font-semibold" :class="selectedTransfer.receive_date ? 'text-green-900' : 'text-gray-400'">
                    {{ selectedTransfer.receive_date ? formatDate(selectedTransfer.receive_date) : 'Pendiente' }}
                  </p>
                </div>
              </div>
            </div>
          </div>

          <!-- Tab: Ubicación -->
          <div v-show="activeTab === 'ubicacion'" class="space-y-6 animate-fade-in">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <!-- Origen -->
              <div class="border-2 border-gray-300 p-6 rounded-xl bg-white hover:shadow-md transition-shadow">
                <div class="flex items-center mb-4">
                  <div class="h-12 w-12 rounded-full bg-gray-100 flex items-center justify-center mr-3">
                    <svg class="h-6 w-6 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
                    </svg>
                  </div>
                  <div>
                    <p class="text-xs text-gray-500 uppercase tracking-wide">Origen</p>
                    <p class="text-sm text-gray-600">{{ selectedTransfer.origin_type === 'store' ? 'Bodega' : 'Pabellón' }}</p>
                  </div>
                </div>
                <p class="text-2xl font-bold text-gray-900">{{ getOriginName(selectedTransfer) }}</p>
              </div>

              <!-- Destino -->
              <div class="border-2 border-green-300 bg-gradient-to-br from-green-50 to-green-100 p-6 rounded-xl hover:shadow-md transition-shadow">
                <div class="flex items-center mb-4">
                  <div class="h-12 w-12 rounded-full bg-green-200 flex items-center justify-center mr-3">
                    <svg class="h-6 w-6 text-green-700" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16l-4-4m0 0l4-4m-4 4h18" />
                    </svg>
                  </div>
                  <div>
                    <p class="text-xs text-green-700 uppercase tracking-wide">Destino</p>
                    <p class="text-sm text-green-600">{{ selectedTransfer.destination_type === 'store' ? 'Bodega' : 'Pabellón' }}</p>
                  </div>
                </div>
                <p class="text-2xl font-bold text-green-900">{{ getDestinationName(selectedTransfer) }}</p>
              </div>
            </div>

            <!-- Flecha visual de transferencia -->
            <div class="flex justify-center">
              <div class="flex items-center space-x-4 bg-blue-50 px-6 py-3 rounded-full border border-blue-200">
                <span class="text-sm font-medium text-blue-700">{{ getOriginName(selectedTransfer) }}</span>
                <svg class="h-6 w-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7l5 5m0 0l-5 5m5-5H6" />
                </svg>
                <span class="text-sm font-medium text-blue-700">{{ getDestinationName(selectedTransfer) }}</span>
              </div>
            </div>
          </div>

          <!-- Tab: Personal -->
          <div v-show="activeTab === 'personal'" class="space-y-6 animate-fade-in">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <!-- Enviado por -->
              <div class="bg-blue-50 p-6 rounded-xl border border-blue-200">
                <div class="flex items-start space-x-4">
                  <div class="h-14 w-14 rounded-full bg-blue-100 flex items-center justify-center flex-shrink-0">
                    <svg class="h-8 w-8 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                    </svg>
                  </div>
                  <div class="flex-1">
                    <p class="text-xs text-blue-600 uppercase tracking-wide mb-1">Enviado por</p>
                    <p class="text-lg font-bold text-gray-900">{{ selectedTransfer.sent_by_name || 'No registrado' }}</p>
                    <p class="text-sm text-gray-600 mt-1">RUT: {{ selectedTransfer.sent_by || 'N/A' }}</p>
                  </div>
                </div>
              </div>

              <!-- Recibido por -->
              <div v-if="selectedTransfer.received_by" class="bg-green-50 p-6 rounded-xl border border-green-200">
                <div class="flex items-start space-x-4">
                  <div class="h-14 w-14 rounded-full bg-green-100 flex items-center justify-center flex-shrink-0">
                    <svg class="h-8 w-8 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                    </svg>
                  </div>
                  <div class="flex-1">
                    <p class="text-xs text-green-600 uppercase tracking-wide mb-1">Recibido por</p>
                    <p class="text-lg font-bold text-gray-900">{{ selectedTransfer.received_by_name || 'No registrado' }}</p>
                    <p class="text-sm text-gray-600 mt-1">RUT: {{ selectedTransfer.received_by || 'N/A' }}</p>
                  </div>
                </div>
              </div>
              <div v-else class="bg-gray-50 p-6 rounded-xl border border-gray-200 flex items-center justify-center">
                <div class="text-center">
                  <svg class="h-12 w-12 text-gray-300 mx-auto mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  <p class="text-sm text-gray-500">Pendiente de recepción</p>
                  <p class="text-xs text-gray-400 mt-1">Aún no ha sido recibido</p>
                </div>
              </div>
            </div>
          </div>

          <!-- Tab: Adicional -->
          <div v-show="activeTab === 'adicional'" class="space-y-6 animate-fade-in">
            <div v-if="selectedTransfer.transfer_reason || selectedTransfer.notes || selectedTransfer.rejection_reason" class="space-y-4">
              <div v-if="selectedTransfer.transfer_reason" class="bg-blue-50 p-5 rounded-xl border border-blue-200">
                <div class="flex items-center mb-3">
                  <svg class="h-5 w-5 text-blue-600 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
                  </svg>
                  <p class="text-sm font-semibold text-blue-900">Motivo de transferencia</p>
                </div>
                <p class="text-sm text-gray-700 leading-relaxed">{{ selectedTransfer.transfer_reason }}</p>
              </div>

              <div v-if="selectedTransfer.notes" class="bg-gray-50 p-5 rounded-xl border border-gray-200">
                <div class="flex items-center mb-3">
                  <svg class="h-5 w-5 text-gray-600 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 8h10M7 12h4m1 8l-4-4H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-3l-4 4z" />
                  </svg>
                  <p class="text-sm font-semibold text-gray-900">Notas</p>
                </div>
                <p class="text-sm text-gray-700 leading-relaxed">{{ selectedTransfer.notes }}</p>
              </div>

              <div v-if="selectedTransfer.rejection_reason" class="bg-red-50 p-5 rounded-xl border-2 border-red-300">
                <div class="flex items-center mb-3">
                  <svg class="h-5 w-5 text-red-600 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                  </svg>
                  <p class="text-sm font-semibold text-red-900">Motivo de rechazo</p>
                </div>
                <p class="text-sm text-red-700 leading-relaxed">{{ selectedTransfer.rejection_reason }}</p>
              </div>
            </div>

            <div v-else class="text-center py-12">
              <svg class="h-16 w-16 text-gray-300 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
              <p class="text-gray-500">No hay información adicional</p>
            </div>

            <!-- Timestamps del sistema -->
            <div class="border-t pt-4 text-xs text-gray-500 space-y-1">
              <p class="flex items-center">
                <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
                </svg>
                Creado: {{ formatDate(selectedTransfer.created_at) }}
              </p>
              <p class="flex items-center">
                <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                </svg>
                Última actualización: {{ formatDate(selectedTransfer.updated_at) }}
              </p>
            </div>
          </div>
        </div>

        <!-- Footer del modal -->
        <div class="flex justify-end pt-6 border-t mt-6">
          <button @click="closeDetailsModal" class="btn-secondary">Cerrar</button>
        </div>
      </div>
      </div>
    </Teleport>



  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useNotification } from '@/composables/useNotification'
import supplyTransferService from '@/services/management/supplyTransferService'
import { exportToExcel as exportExcel, formatDateForExcel, formatStatusForExcel } from '@/utils/excelExport'

const { success: showSuccess, error: showError, warning: showWarning, info: showInfo } = useNotification()
const route = useRoute()

const transfers = ref([])
const loading = ref(false)
const error = ref(null)
const showDetailsModal = ref(false)
const selectedTransfer = ref(null)
const activeTab = ref('general')

const filters = ref({
  code: '',
  status: '',
  from_date: '',
  to_date: '',
  pavilion_id: '',
  store_id: ''
})

// Estado de ordenamiento
const sortKey = ref('created_at')
const sortOrder = ref('desc')

// Estado de paginación
const currentPage = ref(1)
const itemsPerPage = 10

// Computed para obtener la lista ordenada
const sortedTransfers = computed(() => {
  if (!transfers.value || transfers.value.length === 0) return []

  // Filtro client-side por código
  const codeFilter = filters.value.code?.trim().toLowerCase()
  const filtered = codeFilter
    ? transfers.value.filter(t => {
        const code = (t.transfer_code || t.code || '').toLowerCase()
        return code.includes(codeFilter)
      })
    : transfers.value
  
  const sorted = [...filtered].sort((a, b) => {
    let aVal = a[sortKey.value]
    let bVal = b[sortKey.value]
    
    // Manejo de valores null/undefined
    if (aVal == null) aVal = ''
    if (bVal == null) bVal = ''
    
    // Manejo de strings (comparación case-insensitive)
    if (typeof aVal === 'string' && sortKey.value !== 'created_at' && sortKey.value !== 'send_date' && sortKey.value !== 'receive_date') {
      aVal = aVal.toLowerCase()
      bVal = typeof bVal === 'string' ? bVal.toLowerCase() : ''
    }
    
    // Manejo de fechas
    if (sortKey.value === 'created_at' || sortKey.value === 'send_date' || sortKey.value === 'receive_date') {
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
const totalPages = computed(() => Math.ceil(sortedTransfers.value.length / itemsPerPage))

const startIndex = computed(() => (currentPage.value - 1) * itemsPerPage)
const endIndex = computed(() => Math.min(startIndex.value + itemsPerPage, sortedTransfers.value.length))

const paginatedTransfers = computed(() => {
  return sortedTransfers.value.slice(startIndex.value, endIndex.value)
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

// Resetear página al cambiar el filtro de código (búsqueda en tiempo real, client-side)
watch(() => filters.value.code, () => { currentPage.value = 1 })

// Aplicar filtros automáticamente al cambiar estado o fechas (requieren llamada al backend)
watch(() => filters.value.status, () => loadTransfers())
watch(() => filters.value.from_date, () => loadTransfers())
watch(() => filters.value.to_date, () => loadTransfers())

const loadTransfers = async () => {
  loading.value = true
  error.value = null
  currentPage.value = 1 // Resetear a la primera página al cargar
  try {
    const data = await supplyTransferService.getTransfers(filters.value)
    // El backend devuelve { transfers: [], total: n, page: n, page_size: n }
    transfers.value = data.transfers || data.data?.transfers || data || []
    console.log('Transferencias cargadas:', transfers.value)
  } catch (err) {
    error.value = err.message || 'Error al cargar transferencias'
    console.error('Error loading transfers:', err)
  } finally {
    loading.value = false
  }
}

const viewDetails = (transfer) => {
  selectedTransfer.value = transfer
  activeTab.value = 'general' // Reset a la primera pestaña
  showDetailsModal.value = true
}

const closeDetailsModal = () => {
  showDetailsModal.value = false
  selectedTransfer.value = null
  activeTab.value = 'general'
}

const clearFilters = () => {
  filters.value = {
    code: '',
    status: '',
    from_date: '',
    to_date: '',
    pavilion_id: '',
    store_id: ''
  }
  currentPage.value = 1
  loadTransfers()
}

const formatDate = (date) => {
  if (!date) return 'N/A'
  return new Date(date).toLocaleDateString('es-CL', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const getOriginName = (transfer) => {
  if (transfer.origin_name) return transfer.origin_name
  return `${transfer.origin_type === 'store' ? 'Bodega' : 'Pabellón'} #${transfer.origin_id}`
}

const getDestinationName = (transfer) => {
  if (transfer.destination_name) return transfer.destination_name
  return `${transfer.destination_type === 'store' ? 'Bodega' : 'Pabellón'} #${transfer.destination_id}`
}

const getStatusClass = (status) => {
  const classes = {
    pendiente: 'bg-yellow-100 text-yellow-800',
    en_transito: 'bg-blue-100 text-blue-800',
    recibido: 'bg-green-100 text-green-800',
    cancelado: 'bg-red-100 text-red-800',
    rechazado: 'bg-red-100 text-red-800',
    // Inglés para compatibilidad
    pending: 'bg-yellow-100 text-yellow-800',
    in_transit: 'bg-blue-100 text-blue-800',
    completed: 'bg-green-100 text-green-800',
    cancelled: 'bg-red-100 text-red-800'
  }
  return classes[status] || 'bg-gray-100 text-gray-800'
}

const getStatusBgClass = (status) => {
  const classes = {
    pendiente: 'bg-gradient-to-r from-yellow-100 to-yellow-200 text-yellow-900',
    en_transito: 'bg-gradient-to-r from-blue-100 to-blue-200 text-blue-900',
    recibido: 'bg-gradient-to-r from-green-100 to-green-200 text-green-900',
    cancelado: 'bg-gradient-to-r from-red-100 to-red-200 text-red-900',
    rechazado: 'bg-gradient-to-r from-red-100 to-red-200 text-red-900',
    // Inglés para compatibilidad
    pending: 'bg-gradient-to-r from-yellow-100 to-yellow-200 text-yellow-900',
    in_transit: 'bg-gradient-to-r from-brand-blue-light from-opacity-30 to-brand-blue-light from-opacity-50 text-brand-blue-dark',
    completed: 'bg-gradient-to-r from-green-100 to-green-200 text-green-900',
    cancelled: 'bg-gradient-to-r from-red-100 to-red-200 text-red-900'
  }
  return classes[status] || 'bg-gradient-to-r from-gray-100 to-gray-200 text-gray-900'
}

const getStatusLabel = (status) => {
  const labels = {
    pendiente: 'Pendiente',
    en_transito: 'En Tránsito',
    recibido: 'Recibido',
    cancelado: 'Cancelado',
    rechazado: 'Rechazado',
    // Inglés para compatibilidad
    pending: 'Pendiente',
    in_transit: 'En Tránsito',
    completed: 'Completado',
    cancelled: 'Cancelado'
  }
  return labels[status] || status
}

// Función helper para mantener compatibilidad con código existente
const showNotification = (message, type = 'info') => {
  if (type === 'success') {
    showSuccess(message)
  } else if (type === 'error') {
    showError(message)
  } else if (type === 'warning') {
    showWarning(message)
  } else {
    showInfo(message)
  }
}

const exportToExcel = async () => {
  try {
    const columns = [
      { key: 'transfer_code', label: 'Código de Transferencia' },
      { key: 'qr_code', label: 'Código QR' },
      { key: 'origin_name', label: 'Origen', formatter: (val, item) => getOriginName(item) },
      { key: 'destination_name', label: 'Destino', formatter: (val, item) => getDestinationName(item) },
      { key: 'status', label: 'Estado', formatter: (val) => formatStatusForExcel(val) },
      { key: 'created_at', label: 'Fecha de Creación', formatter: formatDateForExcel },
      { key: 'send_date', label: 'Fecha de Envío', formatter: formatDateForExcel },
      { key: 'receive_date', label: 'Fecha de Recepción', formatter: formatDateForExcel },
      { key: 'sent_by_name', label: 'Enviado por' },
      { key: 'picked_up_by_name', label: 'Retirado por' },
      { key: 'picked_up_date', label: 'Fecha de Retiro', formatter: formatDateForExcel },
      { key: 'received_by_name', label: 'Recibido por' },
      { key: 'transfer_reason', label: 'Motivo' },
      { key: 'notes', label: 'Notas' },
      { key: 'rejection_reason', label: 'Motivo de Rechazo' }
    ]
    
    await exportExcel(sortedTransfers.value, columns, 'transferencias')
    showNotification('Exportación a Excel completada exitosamente', 'success')
  } catch (error) {
    console.error('Error al exportar:', error)
    showNotification('Error al exportar a Excel: ' + error.message, 'error')
  }
}

onMounted(() => {
  if (route.query.status) {
    filters.value.status = route.query.status
  }
  loadTransfers()
})
</script>

<style scoped>
/* Animaciones para las notificaciones */
.notification-enter-active,
.notification-leave-active {
  transition: all 0.3s ease;
}

.notification-enter-from {
  transform: translateX(100%);
  opacity: 0;
}

.notification-leave-to {
  transform: translateX(100%);
  opacity: 0;
}

.notification-toast {
  transition: all 0.3s ease;
  animation: slideIn 0.3s ease-out;
}

@keyframes slideIn {
  from {
    transform: translateX(100%);
    opacity: 0;
  }
  to {
    transform: translateX(0);
    opacity: 1;
  }
}

/* Animación de fade-in para las pestañas */
.animate-fade-in {
  animation: fadeIn 0.3s ease-in;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
