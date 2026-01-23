<template>
  <div class="max-w-7xl mx-auto p-3 sm:p-6">
    <!-- Encabezado -->
    <div class="mb-4 sm:mb-6">
      <div class="flex flex-col sm:flex-row sm:justify-between sm:items-center gap-3">
        <div>
          <h1 class="text-xl sm:text-2xl font-bold text-gray-900">Solicitudes de Insumo</h1>
          <p class="mt-1 text-sm sm:text-base text-gray-600">Gestión de solicitudes con trazabilidad QR</p>
        </div>
        <div class="flex flex-col sm:flex-row gap-2 w-full sm:w-auto">
          <button 
            @click="exportToExcel" 
            :disabled="loading || filteredRequests.length === 0"
            class="btn-secondary flex items-center justify-center"
            :class="{ 'opacity-50 cursor-not-allowed': loading || filteredRequests.length === 0 }"
          >
            <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
            Exportar a Excel
          </button>
          <router-link
            v-if="authStore.canCreateRequests"
            to="/supply-requests/new"
            class="btn-primary w-full sm:w-auto"
          >
            <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
            </svg>
            Nueva Solicitud
          </router-link>
        </div>
      </div>
    </div>

    <!-- Filtros y búsqueda -->
    <div class="mb-4 sm:mb-6 card p-3 sm:p-4">
      <!-- Badge de filtro activo -->
      <div v-if="filters.statusCategory" class="mb-3 flex items-center gap-2">
        <span class="text-xs font-medium text-gray-600">Filtro activo:</span>
        <span :class="getFilterBadgeClass(filters.statusCategory)" class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium">
          {{ getFilterLabel(filters.statusCategory) }}
          <button @click="filterByStatCategory('')" class="ml-1.5 hover:bg-white hover:bg-opacity-20 rounded-full">
            <svg class="h-3 w-3" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
            </svg>
          </button>
        </span>
      </div>
      
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-3 sm:gap-4">
        <!-- Filtro por estado -->
        <div>
          <label for="statusFilter" class="block text-sm font-medium text-gray-700 mb-1">
            Estado
          </label>
          <select
            id="statusFilter"
            v-model="filters.status"
            @change="loadSupplyRequests"
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
          >
            <option value="">Todos los estados</option>
            <option value="pendiente_pavedad">Pendiente Pavedad</option>
            <option value="asignado_bodega">Asignado a Bodega</option>
            <option value="en_proceso">En Proceso</option>
            <option value="pendiente_revision">Pendiente Revisión</option>
            <option value="devuelto">Devuelto al Solicitante</option>
            <option value="devuelto_al_encargado">Devuelto al Encargado</option>
            <option value="aprobado">Aprobado</option>
            <option value="parcialmente_aprobado">Parcialmente Aprobado</option>
            <option value="rechazado">Rechazado</option>
            <option value="completado">Completado</option>
            <option value="cancelado">Cancelado</option>
          </select>
        </div>

        <!-- Filtro por fecha de cirugía -->
        <div>
          <label for="surgeryDateFilter" class="block text-sm font-medium text-gray-700 mb-1">
            Fecha de Cirugía
          </label>
          <input
            type="date"
            id="surgeryDateFilter"
            v-model="filters.surgeryDate"
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
          />
        </div>

        <!-- Búsqueda por número de solicitud -->
        <div>
          <label for="searchNumber" class="block text-sm font-medium text-gray-700 mb-1">
            Número de Solicitud
          </label>
          <input
            type="text"
            id="searchNumber"
            v-model="filters.search"
            @input="filterRequests"
            placeholder="SOL-XXXXXX"
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
          />
        </div>

        <!-- Botón refrescar -->
        <div class="flex items-end">
          <button
            @click="refreshRequests"
            :disabled="loading"
            class="w-full px-3 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50"
          >
            <svg class="h-4 w-4 inline mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
            Refrescar
          </button>
        </div>
      </div>
    </div>

    <!-- Estadísticas rápidas -->
    <div v-if="stats" class="mb-4 sm:mb-6 grid grid-cols-2 lg:grid-cols-4 gap-3 sm:gap-4">
      <div @click="filterByStatCategory('pending')" class="bg-white p-3 sm:p-4 rounded-lg shadow border cursor-pointer hover:shadow-md transition-shadow">
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <div class="w-7 h-7 sm:w-8 sm:h-8 bg-yellow-100 rounded-full flex items-center justify-center">
              <svg class="w-3.5 h-3.5 sm:w-4 sm:h-4 text-yellow-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
          </div>
          <div class="ml-2 sm:ml-3 min-w-0">
            <p class="text-xs sm:text-sm font-medium text-gray-500 truncate">Pendientes</p>
            <p class="text-base sm:text-lg font-semibold text-gray-900">{{ stats.pending || 0 }}</p>
          </div>
        </div>
      </div>

      <div @click="filterByStatCategory('approved')" class="bg-white p-3 sm:p-4 rounded-lg shadow border cursor-pointer hover:shadow-md transition-shadow">
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <div class="w-7 h-7 sm:w-8 sm:h-8 bg-green-100 rounded-full flex items-center justify-center">
              <svg class="w-3.5 h-3.5 sm:w-4 sm:h-4 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
            </div>
          </div>
          <div class="ml-2 sm:ml-3 min-w-0">
            <p class="text-xs sm:text-sm font-medium text-gray-500 truncate">Aprobadas</p>
            <p class="text-base sm:text-lg font-semibold text-gray-900">{{ stats.approved || 0 }}</p>
          </div>
        </div>
      </div>

      <div @click="filterByStatCategory('rejected')" class="bg-white p-3 sm:p-4 rounded-lg shadow border cursor-pointer hover:shadow-md transition-shadow">
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <div class="w-7 h-7 sm:w-8 sm:h-8 bg-red-100 rounded-full flex items-center justify-center">
              <svg class="w-3.5 h-3.5 sm:w-4 sm:h-4 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </div>
          </div>
          <div class="ml-2 sm:ml-3 min-w-0">
            <p class="text-xs sm:text-sm font-medium text-gray-500 truncate">Rechazadas</p>
            <p class="text-base sm:text-lg font-semibold text-gray-900">{{ stats.rejected || 0 }}</p>
          </div>
        </div>
      </div>

      <div @click="filterByStatCategory('all')" class="bg-white p-3 sm:p-4 rounded-lg shadow border cursor-pointer hover:shadow-md transition-shadow">
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <div class="w-7 h-7 sm:w-8 sm:h-8 bg-gray-100 rounded-full flex items-center justify-center">
              <svg class="w-3.5 h-3.5 sm:w-4 sm:h-4 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
              </svg>
            </div>
          </div>
          <div class="ml-2 sm:ml-3 min-w-0">
            <p class="text-xs sm:text-sm font-medium text-gray-500 truncate">Total</p>
            <p class="text-base sm:text-lg font-semibold text-gray-900">{{ totalRequests }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Lista de solicitudes -->
    <div class="bg-white shadow-lg rounded-lg border">
      <!-- Loading state -->
      <div v-if="loading" class="p-8 text-center">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
        <p class="text-gray-600">Cargando solicitudes...</p>
      </div>

      <!-- Error state -->
      <div v-else-if="error" class="p-8 text-center">
        <svg class="h-12 w-12 text-red-400 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <h3 class="text-lg font-medium text-gray-900 mb-2">Error al cargar solicitudes</h3>
        <p class="text-gray-600 mb-4">{{ error }}</p>
        <button
          @click="loadSupplyRequests"
          class="btn-primary"
        >
          Reintentar
        </button>
      </div>

      <!-- Empty state -->
      <div v-else-if="filteredRequests.length === 0 && !loading" class="p-8 text-center">
        <svg class="h-12 w-12 text-gray-400 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
        </svg>
        <h3 class="text-lg font-medium text-gray-900 mb-2">No hay solicitudes</h3>
        <p class="text-gray-600 mb-4">
          {{ (requests.length === 0 && !filters.status && !filters.statusCategory && !filters.urgency && !filters.search && !filters.surgeryDate) ? 'No se han creado solicitudes aún' : 'No se encontraron solicitudes con los filtros aplicados' }}
        </p>
        <router-link
          v-if="requests.length === 0 && authStore.canCreateRequests && !filters.status && !filters.statusCategory"
          to="/supply-requests/new"
          class="btn-primary"
        >
          <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          Crear Primera Solicitud
        </router-link>
      </div>

      <!-- Contenido con datos -->
      <div v-else>
        <!-- Vista de tarjetas para móviles -->
        <div class="md:hidden">
          <div class="space-y-3 p-3">
            <div 
              v-for="request in paginatedRequests" 
              :key="request.id"
              @click="viewRequest(request.id)"
              class="bg-white border rounded-lg p-4 shadow-sm hover:shadow-md transition-shadow cursor-pointer"
              :class="isEmergency(request) ? 'border-red-400 border-2' : isUrgent(request) ? 'border-orange-400 border-2' : 'border-gray-200'"
            >
              <!-- Header con número y fecha -->
              <div class="flex items-start justify-between mb-3 pb-3 border-b border-gray-100">
                <div class="flex-1 min-w-0">
                  <div class="flex items-center gap-2 mb-1">
                    <svg class="h-4 w-4 text-gray-400 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                    </svg>
                    <span class="text-sm font-semibold text-gray-900 truncate">{{ request.request_number }}</span>
                    <!-- Badge de urgencia -->
                    <span v-if="isEmergency(request) || isUrgent(request)" 
                          :class="getUrgencyBadgeClass(request)" 
                          class="inline-flex items-center px-2 py-0.5 text-xs font-bold rounded-full border">
                      {{ getUrgencyLabel(request) }}
                    </span>
                  </div>
                  <div class="flex items-center gap-2 text-xs text-gray-500">
                    <svg class="h-3.5 w-3.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                    </svg>
                    <span>{{ formatDate(request.request_date) }}</span>
                  </div>
                </div>
              </div>

              <!-- Información principal -->
              <div class="space-y-2.5">
                <!-- Pabellón -->
                <div class="flex items-start">
                  <div class="w-20 text-xs font-medium text-gray-500 flex-shrink-0">Pabellón:</div>
                  <div class="text-sm text-gray-900 font-medium">{{ getPavilionName(request.pavilion_id) }}</div>
                </div>

                <!-- Solicitante -->
                <div class="flex items-start">
                  <div class="w-20 text-xs font-medium text-gray-500 flex-shrink-0">Solicitante:</div>
                  <div class="flex-1 min-w-0">
                    <div class="text-sm font-medium text-gray-900 truncate">{{ request.requested_by_name }}</div>
                    <div class="text-xs text-gray-500 truncate">{{ request.requested_by }}</div>
                  </div>
                </div>

                <!-- Items -->
                <div class="flex items-center">
                  <div class="w-20 text-xs font-medium text-gray-500 flex-shrink-0">Items:</div>
                  <div class="text-sm text-gray-900">{{ request.total_items || 'N/A' }} items</div>
                </div>

                <!-- Fecha de cirugía con indicador de urgencia -->
                <div v-if="request.surgery_datetime" class="flex items-start">
                  <div class="w-20 text-xs font-medium text-gray-500 flex-shrink-0">Cirugía:</div>
                  <div class="flex-1 min-w-0">
                    <div class="text-sm text-gray-900 font-medium">{{ formatDate(request.surgery_datetime) }}</div>
                    <div v-if="isEmergency(request) || isUrgent(request)" class="text-xs mt-1" 
                         :class="isEmergency(request) ? 'text-red-600 font-bold' : 'text-orange-600 font-semibold'">
                      ⚠️ {{ getUrgencyLabel(request) }}
                    </div>
                  </div>
                </div>

                <!-- Estado y Prioridad -->
                <div class="flex items-center gap-2 pt-2">
                  <span :class="getStatusBadgeClass(request.status)" class="inline-flex px-2.5 py-1 text-xs font-semibold rounded-full">
                    {{ getStatusLabel(request.status) }}
                  </span>
                </div>
              </div>

              <!-- Acciones -->
              <div class="flex gap-2 pt-3 mt-3 border-t border-gray-100" @click.stop>
                <button
                  @click.stop="viewRequest(request.id)"
                  class="flex-1 inline-flex items-center justify-center px-3 py-2 border border-gray-300 rounded-md text-xs font-medium text-gray-700 bg-white hover:bg-gray-50"
                >
                  <svg class="h-4 w-4 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                  </svg>
                  Ver
                </button>
                
                <!-- Botón de Asignar para Pavedad -->
                <button
                  v-if="authStore.isPavedad"
                  @click.stop="request.status === 'pendiente_pavedad' ? openAssignModal(request) : null"
                  :disabled="request.status !== 'pendiente_pavedad'"
                  :class="[
                    'inline-flex items-center justify-center px-3 py-2 border rounded-md text-xs font-medium',
                    request.status === 'pendiente_pavedad' 
                      ? 'border-purple-300 text-purple-700 bg-purple-50 hover:bg-purple-100 cursor-pointer'
                      : 'border-gray-300 text-gray-400 bg-gray-50 cursor-not-allowed opacity-60'
                  ]"
                  :title="request.status === 'pendiente_pavedad' ? 'Asignar' : 'Ya asignada'"
                >
                  <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 5v2m0 4v2m0 4v2M5 5a2 2 0 00-2 2v3a2 2 0 110 4v3a2 2 0 002 2h14a2 2 0 002-2v-3a2 2 0 110-4V7a2 2 0 00-2-2H5z" />
                  </svg>
                  <span v-if="request.status !== 'pendiente_pavedad'" class="ml-1">✓</span>
                </button>
                
                <!-- Botón Revisar Items para Encargado de Bodega -->
                <button
                  v-if="authStore.isWarehouseManager && request.assigned_to === authStore.getUserRut"
                  @click.stop="(request.status === 'asignado_bodega' || request.status === 'en_proceso' || request.status === 'devuelto_al_encargado') ? openReviewItemsModal(request) : null"
                  :disabled="request.status !== 'asignado_bodega' && request.status !== 'en_proceso' && request.status !== 'devuelto_al_encargado'"
                  :class="[
                    'flex-1 inline-flex items-center justify-center px-3 py-2 border rounded-md text-xs font-medium',
                    (request.status === 'asignado_bodega' || request.status === 'en_proceso' || request.status === 'devuelto_al_encargado')
                      ? 'border-blue-300 text-blue-700 bg-blue-50 hover:bg-blue-100 cursor-pointer'
                      : 'border-gray-300 text-gray-400 bg-gray-50 cursor-not-allowed opacity-60'
                  ]"
                  :title="(request.status === 'asignado_bodega' || request.status === 'en_proceso' || request.status === 'devuelto_al_encargado') ? 'Revisar Insumos' : 'Ya revisada'"
                >
                  <svg class="h-4 w-4 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4" />
                  </svg>
                  {{ (request.status === 'asignado_bodega' || request.status === 'en_proceso' || request.status === 'devuelto_al_encargado') ? 'Revisar Items' : 'Revisado ✓' }}
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- Vista de tabla para desktop -->
        <div class="hidden md:block overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th
                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer select-none"
                @click="toggleSort('request_number')"
              >
                <div class="flex items-center gap-1">
                  <span>Solicitud</span>
                  <span v-if="sortBy === 'request_number'" class="text-gray-400 text-[10px]">
                    {{ sortDirection === 'asc' ? '▲' : '▼' }}
                  </span>
                </div>
              </th>
              <th
                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer select-none"
                @click="toggleSort('pavilion')"
              >
                <div class="flex items-center gap-1">
                  <span>Pabellón</span>
                  <span v-if="sortBy === 'pavilion'" class="text-gray-400 text-[10px]">
                    {{ sortDirection === 'asc' ? '▲' : '▼' }}
                  </span>
                </div>
              </th>
              <th
                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer select-none"
                @click="toggleSort('requester')"
              >
                <div class="flex items-center gap-1">
                  <span>Solicitante</span>
                  <span v-if="sortBy === 'requester'" class="text-gray-400 text-[10px]">
                    {{ sortDirection === 'asc' ? '▲' : '▼' }}
                  </span>
                </div>
              </th>
              <th
                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer select-none"
                @click="toggleSort('status')"
              >
                <div class="flex items-center gap-1">
                  <span>Estado</span>
                  <span v-if="sortBy === 'status'" class="text-gray-400 text-[10px]">
                    {{ sortDirection === 'asc' ? '▲' : '▼' }}
                  </span>
                </div>
              </th>
              <th
                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer select-none"
                @click="toggleSort('surgery_datetime')"
              >
                <div class="flex items-center gap-1">
                  <span>Fecha de Cirugía</span>
                  <span v-if="sortBy === 'surgery_datetime'" class="text-gray-400 text-[10px]">
                    {{ sortDirection === 'asc' ? '▲' : '▼' }}
                  </span>
                </div>
              </th>
              <th
                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer select-none"
                @click="toggleSort('items')"
              >
                <div class="flex items-center gap-1">
                  <span>Items</span>
                  <span v-if="sortBy === 'items'" class="text-gray-400 text-[10px]">
                    {{ sortDirection === 'asc' ? '▲' : '▼' }}
                  </span>
                </div>
              </th>
              <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                Acciones
              </th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr
              v-for="request in paginatedRequests"
              :key="request.id"
              class="hover:bg-gray-50 cursor-pointer"
              @click="viewRequest(request.id)"
            >
              <!-- Número de solicitud -->
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="flex items-center gap-2">
                <div class="text-sm font-medium text-gray-900">
                  {{ request.request_number }}
                  </div>
                  <!-- Badge de urgencia -->
                  <span v-if="isEmergency(request) || isUrgent(request)" 
                        :class="getUrgencyBadgeClass(request)" 
                        class="inline-flex items-center px-2 py-0.5 text-xs font-bold rounded-full border">
                    {{ getUrgencyLabel(request) }}
                  </span>
                </div>
              </td>

              <!-- Pabellón -->
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="text-sm text-gray-900">
                  {{ getPavilionName(request.pavilion_id) }}
                </div>
              </td>

              <!-- Solicitante -->
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="text-sm font-medium text-gray-900">
                  {{ request.requested_by_name }}
                </div>
                <div class="text-sm text-gray-500">
                  {{ request.requested_by }}
                </div>
              </td>

              <!-- Estado -->
              <td class="px-6 py-4 whitespace-nowrap">
                <span :class="getStatusBadgeClass(request.status)" class="inline-flex px-2 py-1 text-xs font-semibold rounded-full">
                  {{ getStatusLabel(request.status) }}
                </span>
              </td>

              <!-- Fecha de Cirugía -->
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="text-sm text-gray-900">
                  {{ formatDate(request.surgery_datetime) }}
                </div>
                <span v-if="isEmergency(request) || isUrgent(request)" 
                      :class="getUrgencyBadgeClass(request)" 
                      class="inline-flex px-2 py-1 text-xs font-semibold rounded-full mt-1 border">
                  {{ getUrgencyLabel(request) }}
                </span>
              </td>

              <!-- Número de items -->
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                {{ request.total_items || 'N/A' }} items
              </td>

              <!-- Acciones -->
              <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                <div class="flex justify-end space-x-2">
                  <button
                    @click.stop="viewRequest(request.id)"
                    class="text-blue-600 hover:text-blue-900 p-1"
                    title="Ver detalles"
                  >
                    <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                    </svg>
                  </button>

                  <!-- Botón de Asignar para Pavedad -->
                  <button
                    v-if="authStore.isPavedad"
                    @click.stop="request.status === 'pendiente_pavedad' ? openAssignModal(request) : null"
                    :disabled="request.status !== 'pendiente_pavedad'"
                    :class="[
                      'p-1',
                      request.status === 'pendiente_pavedad'
                        ? 'text-purple-600 hover:text-purple-900 cursor-pointer'
                        : 'text-gray-400 cursor-not-allowed opacity-60'
                    ]"
                    :title="request.status === 'pendiente_pavedad' ? 'Asignar' : 'Ya asignada'"
                  >
                    <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 5v2m0 4v2m0 4v2M5 5a2 2 0 00-2 2v3a2 2 0 110 4v3a2 2 0 002 2h14a2 2 0 002-2v-3a2 2 0 110-4V7a2 2 0 00-2-2H5z" />
                    </svg>
                    <span v-if="request.status !== 'pendiente_pavedad'" class="text-xs">✓</span>
                  </button>

                  <!-- Botón Revisar Items para Encargado de Bodega -->
                  <button
                    v-if="authStore.isWarehouseManager && request.assigned_to === authStore.getUserRut"
                    @click.stop="(request.status === 'asignado_bodega' || request.status === 'en_proceso' || request.status === 'devuelto_al_encargado') ? openReviewItemsModal(request) : null"
                    :disabled="request.status !== 'asignado_bodega' && request.status !== 'en_proceso' && request.status !== 'devuelto_al_encargado'"
                    :class="[
                      'p-1',
                      (request.status === 'asignado_bodega' || request.status === 'en_proceso' || request.status === 'devuelto_al_encargado')
                        ? 'text-blue-600 hover:text-blue-900 cursor-pointer'
                        : 'text-gray-400 cursor-not-allowed opacity-60'
                    ]"
                    :title="(request.status === 'asignado_bodega' || request.status === 'en_proceso' || request.status === 'devuelto_al_encargado') ? 'Revisar insumos' : 'Ya revisada'"
                  >
                    <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4" />
                    </svg>
                    <span v-if="request.status !== 'asignado_bodega' && request.status !== 'en_proceso' && request.status !== 'devuelto_al_encargado'" class="text-xs">✓</span>
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
        </div>
      </div>

      <!-- Paginación -->
      <div class="px-3 sm:px-6 py-3 sm:py-4 border-t border-gray-200 bg-gray-50">
        <div class="flex flex-col sm:flex-row items-center justify-between gap-3 sm:gap-4">
          <div class="text-xs sm:text-sm text-gray-700 text-center sm:text-left">
            Mostrando {{ filteredRequests.length > 0 ? startIndex + 1 : 0 }} a {{ endIndex }} de {{ totalRequests }} solicitudes
          </div>
          <div class="flex items-center gap-2 w-full sm:w-auto">
            <button
              @click="previousPage"
              :disabled="currentPage === 1"
              class="btn-secondary px-3 py-2 text-xs sm:text-sm min-w-[60px] sm:min-w-[70px] disabled:opacity-50 disabled:cursor-not-allowed flex-1 sm:flex-initial"
            >
              <span class="hidden sm:inline">Anterior</span>
              <span class="sm:hidden">Ant.</span>
            </button>
            <span class="px-2 sm:px-3 py-2 text-xs sm:text-sm text-gray-700 bg-gray-100 rounded-md min-w-[80px] sm:min-w-[90px] text-center font-medium whitespace-nowrap">
              Pág {{ currentPage }} / {{ totalPages }}
            </span>
            <button
              @click="nextPage"
              :disabled="currentPage === totalPages"
              class="btn-secondary px-3 py-2 text-xs sm:text-sm min-w-[60px] sm:min-w-[70px] disabled:opacity-50 disabled:cursor-not-allowed flex-1 sm:flex-initial"
            >
              <span class="hidden sm:inline">Siguiente</span>
              <span class="sm:hidden">Sig.</span>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Modal de Asignación -->
    <AssignRequestModal
      :show="showAssignModal"
      :request="selectedRequestForAssignment"
      @close="showAssignModal = false"
      @assigned="handleRequestAssigned"
    />

    <!-- Modal de Revisión de Items -->
    <ReviewItemsModal
      :show="showReviewItemsModal"
      :request="selectedRequestForReview"
      @close="showReviewItemsModal = false"
      @itemsReviewed="handleItemsReviewed"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import supplyRequestService from '@/services/requests/supplyRequestService'
import pavilionService from '@/services/config/pavilionService'
import AssignRequestModal from '@/components/requests/AssignRequestModal.vue'
import ReviewItemsModal from '@/components/requests/ReviewItemsModal.vue'
import { format } from 'date-fns'
import { es } from 'date-fns/locale'
import { useNotification } from '@/composables/useNotification'
import { useAlert } from '@/composables/useAlert'
import { exportToExcel as exportExcel, formatDateForExcel, formatStatusForExcel } from '@/utils/excelExport'

const router = useRouter()
const authStore = useAuthStore()
const { success: showSuccess, error: showError } = useNotification()
const { prompt } = useAlert()

// Estado del modal de asignación
const showAssignModal = ref(false)
const selectedRequestForAssignment = ref(null)

// Estado del modal de revisión de items
const showReviewItemsModal = ref(false)
const selectedRequestForReview = ref(null)

// Estado reactivo
const loading = ref(false)
const requests = ref([])
const pavilions = ref([])
const currentPage = ref(1)
const pageSize = ref(10)

// Ordenamiento de tabla
const sortBy = ref('') // 'request_number', 'pavilion', 'requester', 'status', 'surgery_datetime', 'items'
const sortDirection = ref('asc') // 'asc' | 'desc'

// Filtros
const filters = ref({
  status: '',
  urgency: '',
  search: '',
  statusCategory: '', // pending, approved, rejected, o '' para todos
  surgeryDate: '' // fecha de cirugía en formato YYYY-MM-DD
})

// Computed properties
const filteredRequests = computed(() => {
  let filtered = [...requests.value]

  // Filtrar por estado específico (del select)
  if (filters.value.status) {
    filtered = filtered.filter(r => r.status === filters.value.status)
  }

  // Filtrar por categoría de estadística (si está activo)
  if (filters.value.statusCategory) {
    if (filters.value.statusCategory === 'pending') {
      // Pendientes: todas las que NO están aprobadas, parcialmente aprobadas ni rechazadas
      filtered = filtered.filter(r => 
        !['aprobado', 'parcialmente_aprobado', 'rechazado'].includes(r.status)
      )
    } else if (filters.value.statusCategory === 'approved') {
      // Aprobadas: aprobadas + parcialmente aprobadas
      filtered = filtered.filter(r => 
        r.status === 'aprobado' || r.status === 'parcialmente_aprobado'
      )
    } else if (filters.value.statusCategory === 'rejected') {
      // Rechazadas: solo rechazadas
      filtered = filtered.filter(r => 
        r.status === 'rechazado'
      )
    }
  }

  if (filters.value.urgency) {
    filtered = filtered.filter(request => 
      getUrgencyLevel(request) === filters.value.urgency)
  }

  if (filters.value.surgeryDate) {
    filtered = filtered.filter(request => {
      if (!request.surgery_datetime) return false
      // Extraer solo la fecha sin considerar hora ni zona horaria
      const requestDate = request.surgery_datetime.split('T')[0]
      return requestDate === filters.value.surgeryDate
    })
  }

  if (filters.value.search) {
    const search = filters.value.search.toLowerCase()
    filtered = filtered.filter(request =>
      request.request_number.toLowerCase().includes(search)
    )
  }

  // Ordenar: si hay sortBy activo, usarlo; si no, usar orden por urgencia parecido al backend
  filtered.sort((a, b) => {
    const direction = sortDirection.value === 'desc' ? -1 : 1

    if (sortBy.value === 'request_number') {
      return direction * a.request_number.localeCompare(b.request_number)
    }

    if (sortBy.value === 'pavilion') {
      return direction * getPavilionName(a.pavilion_id).localeCompare(getPavilionName(b.pavilion_id))
    }

    if (sortBy.value === 'requester') {
      return direction * a.requested_by_name.localeCompare(b.requested_by_name)
    }

    if (sortBy.value === 'status') {
      return direction * getStatusLabel(a.status).localeCompare(getStatusLabel(b.status))
    }

    if (sortBy.value === 'surgery_datetime') {
      const aDate = a.surgery_datetime ? new Date(a.surgery_datetime).getTime() : Infinity
      const bDate = b.surgery_datetime ? new Date(b.surgery_datetime).getTime() : Infinity
      return direction * (aDate - bDate)
    }

    if (sortBy.value === 'items') {
      const aItems = a.total_items || 0
      const bItems = b.total_items || 0
      return direction * (aItems - bItems)
    }

    // Orden por defecto (sin sortBy): mismas reglas que el backend
    // 1) Completadas siempre al final
    const aCompleted = a.status === 'completado'
    const bCompleted = b.status === 'completado'
    if (aCompleted !== bCompleted) {
      return aCompleted ? 1 : -1
    }

    // 2) Emergencias (<12h) primero, luego urgentes (<48h), luego resto
    const aHours = getHoursUntilSurgery(a.surgery_datetime)
    const bHours = getHoursUntilSurgery(b.surgery_datetime)

    const getLevel = (hours) => {
      if (hours === null || hours < 0) return 2
      if (hours <= 12) return 0
      if (hours <= 48) return 1
      return 2
    }

    const aLevel = getLevel(aHours)
    const bLevel = getLevel(bHours)

    if (aLevel !== bLevel) {
      return aLevel - bLevel
    }

    // 3) Como desempate, ordenar por fecha de cirugía ascendente
    const aDate = a.surgery_datetime ? new Date(a.surgery_datetime).getTime() : Infinity
    const bDate = b.surgery_datetime ? new Date(b.surgery_datetime).getTime() : Infinity
    return aDate - bDate
  })

  return filtered
})

const paginatedRequests = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredRequests.value.slice(start, end)
})

const totalRequests = computed(() => filteredRequests.value.length)
const totalPages = computed(() => Math.ceil(totalRequests.value / pageSize.value))

const startIndex = computed(() => (currentPage.value - 1) * pageSize.value)
const endIndex = computed(() => Math.min(startIndex.value + pageSize.value, totalRequests.value))

// Calcular estadísticas localmente basándose en las solicitudes filtradas
const stats = computed(() => {
  const allRequests = requests.value
  
  // Pendientes: todas las que NO están aprobadas, parcialmente aprobadas ni rechazadas
  const pending = allRequests.filter(r => 
    !['aprobado', 'parcialmente_aprobado', 'rechazado'].includes(r.status)
  ).length
  
  // Aprobadas: aprobadas + parcialmente aprobadas
  const approved = allRequests.filter(r => 
    r.status === 'aprobado' || r.status === 'parcialmente_aprobado'
  ).length
  
  // Rechazadas: solo rechazadas
  const rejected = allRequests.filter(r => 
    r.status === 'rechazado'
  ).length
  
  return {
    pending,
    approved,
    rejected,
    total: allRequests.length
  }
})

// Métodos
const loadSupplyRequests = async () => {
  loading.value = true

  try {
    let result
    
    // Si es encargado de bodega, obtener solo sus solicitudes asignadas
    if (authStore.isWarehouseManager) {
      const userRut = authStore.getUserRut
      result = await supplyRequestService.getAssignedRequestsForWarehouseManager(userRut)
    } 
    // Si es Pavedad, obtener solicitudes pendientes para asignar
    else if (authStore.isPavedad) {
      result = await supplyRequestService.getPendingRequestsForPavedad()
    }
    // Admin ve todas las solicitudes
    else if (authStore.isAdmin) {
      result = await supplyRequestService.getAllSupplyRequests(100, 0, filters.value.status)
    }
    // Consignación ve todas las solicitudes (mismo comportamiento que encargado de bodega)
    else if (authStore.isConsignation) {
      result = await supplyRequestService.getAllSupplyRequests(100, 0, filters.value.status)
    }
    // Doctor ve solo sus solicitudes
    else if (authStore.canCreateRequests) {
      result = await supplyRequestService.getAllSupplyRequests(100, 0, filters.value.status)
    }
    else {
      result = { success: true, data: { requests: [] } }
    }
    
    if (result.success) {
      // Si result.success es true, tratar como éxito incluso si no hay data
      let allRequests = []
      
      if (result.data) {
        allRequests = result.data.requests || result.data || []
      }
      
      // Asegurar que allRequests sea un array
      if (!Array.isArray(allRequests)) {
        allRequests = []
      }
      
      // Filtrar por doctor si corresponde
      if (authStore.canCreateRequests && !authStore.isAdmin && !authStore.isPavedad && !authStore.isWarehouseManager) {
        const userRut = authStore.getUserRut
        requests.value = allRequests.filter(request => request.requested_by === userRut)
      } else {
        requests.value = allRequests
      }
      
      console.log('Solicitudes cargadas:', requests.value)
    } else {
      // Solo establecer error si realmente hay un error (no solo falta de datos)
      requests.value = []
      // Solo mostrar error si hay un mensaje de error específico
      if (result.error && result.error !== 'No hay solicitudes' && result.error !== 'No se encontraron solicitudes') {
        showError(result.error)
      }
    }
  } catch (err) {
    console.error('Error cargando solicitudes:', err)
    showError('Error al conectar con el servidor')
    requests.value = []
  } finally {
    loading.value = false
  }
}

const loadPavilions = async () => {
  try {
    pavilions.value = await pavilionService.getAllPavilions()
  } catch (err) {
    console.error('Error cargando pabellones:', err)
  }
}

const loadStats = async () => {
  // Las estadísticas se calculan localmente basándose en las solicitudes filtradas
  // No necesitamos cargar estadísticas del backend
}

const filterByStatCategory = (category) => {
  // Limpiar búsqueda
  filters.value.search = ''
  
  // Aplicar filtro según la categoría
  if (category === 'pending') {
    // Mostrar todas las que NO están aprobadas, parcialmente aprobadas ni rechazadas
    filters.value.statusCategory = 'pending'
  } else if (category === 'approved') {
    // Mostrar aprobadas + parcialmente aprobadas
    filters.value.statusCategory = 'approved'
  } else if (category === 'rejected') {
    // Mostrar solo rechazadas
    filters.value.statusCategory = 'rejected'
  } else {
    // Mostrar todas
    filters.value.statusCategory = ''
  }
  
  currentPage.value = 1
}

const filterRequests = () => {
  currentPage.value = 1
}

const refreshRequests = () => {
  // Limpiar todos los filtros
  filters.value = {
    status: '',
    urgency: '',
    search: '',
    statusCategory: '',
    surgeryDate: ''
  }
  // Recargar solicitudes
  loadSupplyRequests()
}

const toggleSort = (column) => {
  if (sortBy.value === column) {
    // Alternar asc/desc
    sortDirection.value = sortDirection.value === 'asc' ? 'desc' : 'asc'
  } else {
    sortBy.value = column
    sortDirection.value = 'asc'
  }
}

const viewRequest = (id) => {
  router.push(`/supply-requests/${id}`)
}

const approveRequest = async (id) => {
  try {
    const approvalData = {
      approved_by: authStore.getUserRut || 'ADMIN',
      approved_by_name: authStore.getUserName || 'Sistema Admin',
      approval_notes: 'Aprobada desde interfaz web'
    }
    
    await supplyRequestService.approveSupplyRequest(id, approvalData)
    
    showSuccess('La solicitud ha sido aprobada exitosamente')
    
    await loadSupplyRequests()
  } catch (err) {
    console.error('Error aprobando solicitud:', err)
    showError(err.response?.data?.error || err.message || 'Error al aprobar la solicitud')
  }
}

const rejectRequest = async (id) => {
  const reason = await prompt(
    'Ingrese el motivo del rechazo:',
    'Rechazar Solicitud',
    {
      placeholder: 'Motivo del rechazo...',
      inputType: 'textarea',
      inputValidator: (value) => {
        if (!value || !value.trim()) {
          return 'Debe ingresar un motivo para rechazar la solicitud'
        }
        return null
      }
    }
  )
  if (!reason) {
    return
  }

  try {
    const rejectionData = {
      rejected_by: authStore.getUserRut || 'ADMIN',
      rejected_by_name: authStore.getUserName || 'Sistema Admin',
      notes: reason
    }
    
    await supplyRequestService.rejectSupplyRequest(id, rejectionData)
    
    showSuccess('La solicitud ha sido rechazada')
    
    await loadSupplyRequests()
  } catch (err) {
    console.error('Error rechazando solicitud:', err)
    showError(err.response?.data?.error || err.message || 'Error al rechazar la solicitud')
  }
}

const openAssignModal = (request) => {
  selectedRequestForAssignment.value = request
  showAssignModal.value = true
}

const handleRequestAssigned = () => {
  showAssignModal.value = false
  selectedRequestForAssignment.value = null
  loadSupplyRequests()
}

const openReviewItemsModal = (request) => {
  selectedRequestForReview.value = request
  showReviewItemsModal.value = true
}

const handleItemsReviewed = () => {
  loadSupplyRequests()
}


const previousPage = () => {
  if (currentPage.value > 1) {
    currentPage.value--
  }
}

const nextPage = () => {
  if (currentPage.value < totalPages.value) {
    currentPage.value++
  }
}

// Utilidades
const formatDate = (dateString) => {
  if (!dateString) return 'N/A'
  try {
    return format(new Date(dateString), 'dd/MM/yyyy HH:mm', { locale: es })
  } catch {
    return dateString
  }
}

const getPavilionName = (pavilionId) => {
  const pavilion = pavilions.value.find(p => p.id === pavilionId)
  return pavilion ? pavilion.name : `Pabellón ${pavilionId}`
}

const getStatusLabel = (status) => {
  return supplyRequestService.getStatusLabel(status)
}

const getFilterBadgeClass = (category) => {
  const classes = {
    'pending': 'bg-yellow-100 text-yellow-800',
    'approved': 'bg-green-100 text-green-800',
    'rejected': 'bg-red-100 text-red-800'
  }
  return classes[category] || 'bg-gray-100 text-gray-800'
}

const getFilterLabel = (category) => {
  const labels = {
    'pending': 'Pendientes',
    'approved': 'Aprobadas',
    'rejected': 'Rechazadas'
  }
  return labels[category] || category
}

const getStatusBadgeClass = (status) => {
  const color = supplyRequestService.getStatusColor(status)
  const classes = {
    'yellow': 'bg-yellow-100 text-yellow-800',
    'green': 'bg-green-100 text-green-800',
    'red': 'bg-red-100 text-red-800',
    'blue': 'bg-blue-100 text-blue-800',
    'purple': 'bg-purple-100 text-purple-800',
    'orange': 'bg-orange-100 text-orange-800',
    'gray': 'bg-gray-100 text-gray-800'
  }
  return classes[color] || classes.gray
}

const formatSurgeryDateTime = (surgeryDateTime) => {
  if (!surgeryDateTime) return 'N/A'
  try {
    return format(new Date(surgeryDateTime), 'dd/MM/yyyy HH:mm', { locale: es })
  } catch {
    return surgeryDateTime
  }
}

// Funciones de urgencia
const getHoursUntilSurgery = (surgeryDatetime) => {
  if (!surgeryDatetime) return null
  const surgeryDate = new Date(surgeryDatetime)
  const now = new Date()
  const diffTime = surgeryDate - now
  const diffHours = diffTime / (1000 * 60 * 60)
  return diffHours
}

const getDaysUntilSurgery = (surgeryDatetime) => {
  if (!surgeryDatetime) return null
  const surgeryDate = new Date(surgeryDatetime)
  const now = new Date()
  const diffTime = surgeryDate - now
  const diffDays = diffTime / (1000 * 60 * 60 * 24)
  return diffDays
}

const isUrgent = (request) => {
  // No marcar urgencia visual para solicitudes completadas
  if (request.status === 'completado') return false
  const hours = getHoursUntilSurgery(request.surgery_datetime)
  return hours !== null && hours > 0 && hours <= 48
}

const isEmergency = (request) => {
  // No marcar emergencia visual para solicitudes completadas
  if (request.status === 'completado') return false
  const hours = getHoursUntilSurgery(request.surgery_datetime)
  return hours !== null && hours > 0 && hours <= 12
}

const getUrgencyLevel = (request) => {
  const hours = getHoursUntilSurgery(request.surgery_datetime)
  if (hours === null || hours < 0) return 'completed'
  if (hours <= 12) return 'emergency'
  if (hours <= 48) return 'urgent'
  if (hours <= 72) return 'normal'
  return 'low'
}

const getUrgencyBadgeClass = (request) => {
  const level = getUrgencyLevel(request)
  const classes = {
    emergency: 'bg-red-100 text-red-800 border-red-300',
    urgent: 'bg-orange-100 text-orange-800 border-orange-300',
    normal: 'bg-yellow-100 text-yellow-800 border-yellow-300',
    low: 'bg-gray-100 text-gray-800 border-gray-300',
    completed: 'bg-gray-100 text-gray-500 border-gray-300'
  }
  return classes[level] || classes.low
}

const getUrgencyLabel = (request) => {
  const level = getUrgencyLevel(request)
  const hours = getHoursUntilSurgery(request.surgery_datetime)
  const days = getDaysUntilSurgery(request.surgery_datetime)
  
  if (level === 'emergency') {
    return `EMERGENCIA (${Math.ceil(hours)}h)`
  }
  if (level === 'urgent') {
    return `URGENTE (${Math.ceil(hours)}h)`
  }
  if (level === 'normal') {
    return `Normal (${Math.ceil(days)}d)`
  }
  if (level === 'low') {
    return `Baja (${Math.ceil(days)}d)`
  }
  return 'Completada'
}

// Extraer motivo de rechazo desde las notas
const getRejectionReason = (request) => {
  if (!request.notes) return ''
  const marker = 'MOTIVO DEL RECHAZO:'
  const idx = request.notes.indexOf(marker)
  if (idx === -1) return ''
  return request.notes.substring(idx + marker.length).trim()
}

// Lifecycle
const exportToExcel = async () => {
  try {
    const columns = [
      { key: 'request_number', label: 'Número de Solicitud' },
      { 
        key: 'pavilion_id', 
        label: 'Pabellón',
        formatter: (_, row) => getPavilionName(row.pavilion_id)
      },
      { key: 'surgery_datetime', label: 'Fecha de Cirugía', formatter: formatDateForExcel },
      { key: 'status', label: 'Estado', formatter: (val) => formatStatusForExcel(val) },
      { 
        key: 'urgency', 
        label: 'Urgencia',
        formatter: (_, row) => getUrgencyLabel(row)
      },
      { key: 'total_items', label: 'Total de Items' },
      { key: 'created_at', label: 'Fecha de Creación', formatter: formatDateForExcel },
      { key: 'requested_by_name', label: 'Creado por' },
      { key: 'approved_by_name', label: 'Aprobado por' },
      { key: 'approval_date', label: 'Fecha de Aprobación', formatter: formatDateForExcel },
      { key: 'assigned_to_name', label: 'Asignado a' },
      { key: 'assigned_date', label: 'Fecha de Asignación', formatter: formatDateForExcel },
      { key: 'notes', label: 'Observaciones' },
      { 
        key: 'rejection_reason', 
        label: 'Motivo de Rechazo',
        formatter: (_, row) => getRejectionReason(row)
      }
    ]
    
    await exportExcel(filteredRequests.value, columns, 'solicitudes_insumos')
    showSuccess('El archivo Excel se ha descargado exitosamente')
  } catch (error) {
    console.error('Error al exportar:', error)
    showError('Ocurrió un error al exportar a Excel: ' + error.message)
  }
}

onMounted(async () => {
  await Promise.all([
    loadSupplyRequests(),
    loadPavilions()
  ])
})
</script>

<style scoped>
/* Estilos adicionales */
.hover\:bg-gray-50:hover {
  background-color: #f9fafb;
}

/* Estilos para botones secundarios */
/* Usar .btn-secondary de style.css global */

/* Mejoras para dispositivos móviles */
@media (max-width: 640px) {
  /* Inputs más grandes para fácil toque en móviles */
  input[type="text"],
  select {
    min-height: 42px;
    font-size: 16px; /* Previene zoom en iOS */
  }

  /* Mejorar la altura mínima de los botones para toque fácil */
  button,
  .btn-secondary {
    min-height: 44px; /* Tamaño recomendado por Apple para touch targets */
  }

  /* Animación suave al aparecer las tarjetas */
  .space-y-3 > div {
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
  .space-y-3 > div:active {
    transform: scale(0.98);
    transition: transform 0.1s ease-in-out;
  }

  /* Mejorar la legibilidad de los badges */
  .inline-flex.rounded-full {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    max-width: 100%;
  }
}

/* Paginación responsiva */
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

/* Animaciones suaves para mejorar la experiencia */
@media (prefers-reduced-motion: no-preference) {
  button,
  .btn-secondary {
    transition: all 0.15s ease-in-out;
  }
}

/* Estilos para vista de tarjetas */
.space-y-3 > * + * {
  margin-top: 0.75rem;
}

/* Mejoras de accesibilidad táctil */
@media (max-width: 768px) {
  /* Asegurar que todos los elementos interactivos sean táctiles */
  button,
  a {
    touch-action: manipulation;
    -webkit-tap-highlight-color: transparent;
  }

  /* Mejorar el espaciado de las tarjetas */
  .space-y-3 {
    padding: 0.75rem;
  }

  /* Optimizar badges para móviles */
  .inline-flex.px-2 {
    padding-left: 0.5rem;
    padding-right: 0.5rem;
  }
}

/* Ajustes para tablets */
@media (min-width: 641px) and (max-width: 1024px) {
  /* Mantener un buen equilibrio entre compacto y legible */
  .text-sm {
    font-size: 0.875rem;
  }

  .text-xs {
    font-size: 0.75rem;
  }
}
</style>