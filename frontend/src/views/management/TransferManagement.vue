<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="bg-white rounded-lg shadow-sm border p-6">
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h2 class="text-2xl font-semibold text-gray-900">Historial de Transferencias</h2>
          <p class="text-gray-600 mt-1">Consulta y gestiona las transferencias realizadas desde el escáner QR</p>
        </div>
        <div class="flex flex-col sm:flex-row gap-2">
          <button 
            @click="exportToExcel" 
            :disabled="loading || filteredTransfers.length === 0"
            class="btn-secondary flex items-center justify-center"
            :class="{ 'opacity-50 cursor-not-allowed': loading || filteredTransfers.length === 0 }"
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
    <FilterPanel :filters="filterConfig" :result-count="filteredTransfers.length" @filter-change="onFilterChange" />

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
    <DataTable
      v-else
      :columns="tableColumns"
      :rows="filteredTransfers"
      default-sort-key="created_at"
      default-sort-order="desc"
      empty-message="No hay transferencias registradas"
      :table-actions="[{ type: 'view', onClick: (row) => viewDetails(row) }]"
    >
      <template #cell-transfer_code="{ row }">
        <span>{{ row.transfer_code || row.code || 'N/A' }}</span>
      </template>
      <template #cell-origin_name="{ row }">{{ getOriginName(row) }}</template>
      <template #cell-destination_name="{ row }">{{ getDestinationName(row) }}</template>
      <template #cell-cantidad="{ row }">1</template>
      <template #cell-status="{ row }">
        <span class="px-2 py-1 text-xs font-semibold rounded-full" :class="getStatusClass(row.status)">
          {{ getStatusLabel(row.status) }}
        </span>
      </template>
      <template #cell-created_at="{ row }">{{ formatDate(row.created_at) }}</template>

    </DataTable>

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
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useNotification } from '@/composables/useNotification'
import supplyTransferService from '@/services/management/supplyTransferService'
import FilterPanel from '@/components/common/FilterPanel.vue'
import DataTable from '@/components/common/DataTable.vue'
import { exportToExcel as exportExcel, formatDateForExcel, formatStatusForExcel } from '@/utils/excelExport'
import { TRANSFER_STATUS_OPTIONS } from '@/config/statuses'
import { format } from 'date-fns'
import { es } from 'date-fns/locale'

const { success: showSuccess, error: showError, warning: showWarning, info: showInfo } = useNotification()
const route = useRoute()

const transfers = ref([])
const loading = ref(false)
const error = ref(null)
const showDetailsModal = ref(false)
const selectedTransfer = ref(null)
const activeTab = ref('general')

const normalizeTransferStatus = (status) => {
  const statusMap = {
    pending: 'pendiente',
    in_transit: 'en_transito',
    completed: 'recibido',
    cancelled: 'cancelado'
  }
  return statusMap[status] || status || ''
}

const filterState = reactive({
  code: route.query.status ? '' : '',
  status: normalizeTransferStatus(route.query.status),
  from_date: '',
  to_date: '',
  pavilion_id: '',
  store_id: ''
})

const filterConfig = [
  { type: 'text', key: 'code', label: 'Buscar por código', placeholder: 'Código de transferencia' },
  {
    type: 'select', key: 'status', label: 'Estado', default: filterState.status,
    options: [
      { value: '', label: 'Todos' },
      ...TRANSFER_STATUS_OPTIONS
    ]
  },
  { type: 'date', key: 'from_date', label: 'Desde' },
  { type: 'date', key: 'to_date', label: 'Hasta' }
]

const onFilterChange = (key, value) => {
  filterState[key] = value
  if (key === 'status' || key === 'from_date' || key === 'to_date') {
    loadTransfers()
  }
}

const loadTransfers = async () => {
  loading.value = true
  error.value = null
  try {
    const result = await supplyTransferService.getTransfers({
      status: filterState.status,
      from_date: filterState.from_date,
      to_date: filterState.to_date,
      pavilion_id: filterState.pavilion_id,
      store_id: filterState.store_id
    })
    transfers.value = Array.isArray(result) ? result : (result?.transfers || result?.results || result?.data || [])
  } catch (err) {
    error.value = err?.response?.data?.message || err.message || 'Error al cargar transferencias'
  } finally {
    loading.value = false
  }
}

const tableColumns = [
  { key: 'transfer_code', label: 'Código' },
  { key: 'origin_name', label: 'Origen', sortable: false },
  { key: 'destination_name', label: 'Destino', sortable: false },
  { key: 'cantidad', label: 'Cantidad', sortable: false },
  { key: 'status', label: 'Estado' },
  { key: 'created_at', label: 'Fecha' }
]

const filteredTransfers = computed(() => {
  if (!transfers.value || transfers.value.length === 0) return []
  const codeFilter = filterState.code?.trim().toLowerCase()
  if (!codeFilter) return transfers.value
  return transfers.value.filter(t => {
    const code = (t.transfer_code || t.code || '').toLowerCase()
    return code.includes(codeFilter)
  })
})

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
    
    await exportExcel(filteredTransfers.value, columns, 'transferencias')
    showNotification('Exportación a Excel completada exitosamente', 'success')
  } catch (error) {
    console.error('Error al exportar:', error)
    showNotification('Error al exportar a Excel: ' + error.message, 'error')
  }
}

const viewDetails = (transfer) => {
  selectedTransfer.value = transfer
  showDetailsModal.value = true
  activeTab.value = 'general'
}

const closeDetailsModal = () => {
  showDetailsModal.value = false
  selectedTransfer.value = null
}

const formatDate = (dateString) => {
  if (!dateString) return 'No disponible'
  try {
    return format(new Date(dateString), 'dd/MM/yyyy HH:mm', { locale: es })
  } catch {
    return dateString
  }
}

onMounted(() => {
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
