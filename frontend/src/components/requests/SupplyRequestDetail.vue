<template>
  <div class="max-w-7xl mx-auto p-3 sm:p-6">
    <!-- Loading State -->
    <div v-if="loading" class="text-center py-8 sm:py-12">
      <div class="animate-spin rounded-full h-16 w-16 border-b-2 border-blue-600 mx-auto mb-4"></div>
      <p class="text-gray-600">Cargando detalles de la solicitud...</p>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="text-center py-8 sm:py-12">
      <svg class="h-12 w-12 sm:h-16 sm:w-16 text-red-400 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>
      <h3 class="text-base sm:text-lg font-medium text-gray-900 mb-2">Error al cargar solicitud</h3>
      <p class="text-sm sm:text-base text-gray-600 mb-4 px-4">{{ error }}</p>
      <button
        @click="loadSupplyRequest"
        class="btn-primary"
      >
        Reintentar
      </button>
    </div>

    <!-- Contenido principal -->
    <div v-else-if="request">
      <!-- Encabezado -->
      <div class="mb-4 sm:mb-6">
        <div class="flex flex-col sm:flex-row sm:justify-between sm:items-start gap-4">
          <div class="flex-1">
            <div class="flex items-center mb-2 sm:mb-3">
              <button
                @click="$router.go(-1)"
                class="mr-2 sm:mr-4 p-2 text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded-full"
              >
                <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
                </svg>
              </button>
              <h1 class="text-lg sm:text-2xl font-bold text-gray-900 truncate">Solicitud {{ request.request_number }}</h1>
            </div>
            <div class="flex items-center space-x-4">
              <span :class="getStatusBadgeClass(request.status)" class="inline-flex px-3 py-1 text-sm font-semibold rounded-full">
                {{ getStatusLabel(request.status) }}
              </span>
              <span class="text-sm text-gray-600">
                Creada: {{ formatDate(request.request_date) }}
              </span>
            </div>
          </div>

          <!-- Acciones -->
          <div class="flex flex-col sm:flex-row gap-2 sm:flex-shrink-0">
            <button
              v-if="request.status === 'pending' && authStore.canApproveRequests"
              @click="approveRequest"
              :disabled="processing"
              class="btn-success disabled:opacity-50 w-full sm:w-auto"
            >
              <svg class="h-4 w-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
              Aprobar
            </button>

            <button
              v-if="request.status === 'pending' && authStore.canApproveRequests"
              @click="rejectRequest"
              :disabled="processing"
              class="btn-danger disabled:opacity-50 w-full sm:w-auto"
            >
              <svg class="h-4 w-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-weight="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
              Rechazar
            </button>

            <button
              @click="loadSupplyRequest"
              :disabled="loading"
              class="inline-flex items-center justify-center px-3 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 w-full sm:w-auto"
            >
              <svg class="h-4 w-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
              </svg>
              Refrescar
            </button>
          </div>
        </div>
      </div>

      <!-- Grid principal -->
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-4 sm:gap-6">
        <!-- Información de la solicitud -->
        <div class="lg:col-span-2 space-y-4 sm:space-y-6">
          <!-- Información básica -->
          <div class="bg-white rounded-lg shadow border p-4 sm:p-6">
            <h3 class="text-base sm:text-lg font-semibold text-gray-900 mb-3 sm:mb-4">Información Básica</h3>
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-3 sm:gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700">Centro Médico</label>
                <p class="text-sm text-gray-900 mt-1">{{ getMedicalCenterName() }}</p>
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700">Pabellón</label>
                <p class="text-sm text-gray-900 mt-1">{{ getPavilionName(request.pavilion_id) }}</p>
              </div>
              <div>
                <label class="block text-xs sm:text-sm font-medium text-gray-700">Solicitante</label>
                <p class="text-sm text-gray-900 mt-1">{{ request.requested_by_name }}</p>
                <p class="text-xs text-gray-500">{{ request.requested_by }}</p>
              </div>
              <div>
                <label class="block text-xs sm:text-sm font-medium text-gray-700">Fecha de Solicitud</label>
                <p class="text-sm text-gray-900 mt-1">{{ formatDate(request.request_date) }}</p>
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700">Fecha de Cirugía</label>
                <p class="text-sm text-gray-900 mt-1">{{ formatDate(request.surgery_datetime) }}</p>
                <p v-if="isSurgeryUrgent(request.surgery_datetime)" class="text-xs text-red-600 font-medium">
                  <svg class="inline h-3 w-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z" />
                  </svg>
                  Cirugía próxima
                </p>
              </div>
              <div v-if="request.surgery">
                <label class="block text-sm font-medium text-gray-700">Tipo de Cirugía</label>
                <p class="text-sm text-gray-900 mt-1">{{ request.surgery.name }}</p>
                <p v-if="request.surgery.specialty" class="text-xs text-gray-500">
                  {{ request.surgery.specialty.name }}
                </p>
              </div>
              <div v-if="request.approval_date">
                <label class="block text-xs sm:text-sm font-medium text-gray-700">Fecha de Aprobación</label>
                <p class="text-sm text-gray-900 mt-1">{{ formatDate(request.approval_date) }}</p>
                <p class="text-xs text-gray-500">{{ request.approved_by_name }}</p>
              </div>
              <div v-if="request.completed_date">
                <label class="block text-sm font-medium text-gray-700">Fecha de Completado</label>
                <p class="text-sm text-gray-900 mt-1">{{ formatDate(request.completed_date) }}</p>
              </div>
            </div>
          </div>

          <!-- Configuración de Retiro (solo para encargado de bodega asignado) -->
          <div v-if="(request.status === 'asignado_bodega' || request.status === 'en_proceso' || request.status === 'aprobado' || request.status === 'devuelto_al_encargado') && authStore.isWarehouseManager && request.assigned_to === authStore.getUserRut" 
               class="bg-white rounded-lg shadow border p-4 sm:p-6">
            <h3 class="text-base sm:text-lg font-semibold text-gray-900 mb-3 sm:mb-4">Configuración de Retiro</h3>
            <p class="text-sm text-gray-600 mb-4">Configure quién puede retirar los insumos de esta solicitud desde bodega.</p>
            
            <div class="space-y-4">
              <!-- Checkbox para permitir a cualquiera -->
              <div class="flex items-start">
                <input
                  type="checkbox"
                  id="allowAnyonePickup"
                  v-model="pickupConfig.allow_anyone_to_pickup"
                  @change="onCheckboxChange"
                  class="mt-1 h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
                />
                <label for="allowAnyonePickup" class="ml-3 text-sm text-gray-700">
                  <span class="font-medium">Permitir que cualquier persona retire los insumos</span>
                  <p class="text-xs text-gray-500 mt-1">Si está marcado, cualquier persona puede venir a buscar los insumos escaneando el QR.</p>
                </label>
              </div>

              <!-- Campos para persona específica (solo si no está marcado "cualquiera") -->
              <div v-if="!pickupConfig.allow_anyone_to_pickup" class="border-t pt-4 space-y-4">
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-2">
                    Buscar persona autorizada <span class="text-red-500">*</span>
                  </label>
                  <!-- Campo de búsqueda con autocompletado -->
                  <div class="relative">
                    <input
                      type="text"
                      v-model="userSearchQuery"
                      @input="onUserSearch"
                      @focus="showUserSuggestions = true"
                      @blur="handleUserSearchBlur"
                      placeholder="Escriba el nombre o RUT de la persona..."
                      class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                    />
                    <!-- Lista de sugerencias -->
                    <div v-if="showUserSuggestions && filteredUsers.length > 0" 
                         class="absolute z-50 w-full mt-1 bg-white border border-gray-300 rounded-md shadow-lg max-h-60 overflow-auto">
                      <div
                        v-for="user in filteredUsers"
                        :key="user.rut"
                        @mousedown="selectUser(user)"
                        class="px-4 py-2 hover:bg-blue-50 cursor-pointer border-b border-gray-100 last:border-b-0"
                      >
                        <div class="font-medium text-gray-900">{{ user.name }}</div>
                        <div class="text-sm text-gray-500">{{ user.rut }}</div>
                        <div v-if="user.email" class="text-xs text-gray-400">{{ user.email }}</div>
                      </div>
                    </div>
                    <div v-if="showUserSuggestions && userSearchQuery && filteredUsers.length === 0" 
                         class="absolute z-50 w-full mt-1 bg-white border border-gray-300 rounded-md shadow-lg p-4">
                      <p class="text-sm text-gray-500">No se encontraron usuarios</p>
                    </div>
                  </div>
                  <!-- Usuario seleccionado -->
                  <div v-if="pickupConfig.authorized_pickup_rut" class="mt-2 p-2 bg-green-50 border border-green-200 rounded-md">
                    <div class="flex items-center justify-between">
                      <div>
                        <p class="text-sm font-medium text-green-900">{{ pickupConfig.authorized_pickup_name || 'Usuario seleccionado' }}</p>
                        <p class="text-xs text-green-700">{{ pickupConfig.authorized_pickup_rut }}</p>
                      </div>
                      <button
                        @click="clearSelectedUser"
                        class="text-red-600 hover:text-red-800"
                        title="Quitar selección"
                      >
                        <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                        </svg>
                      </button>
                    </div>
                  </div>
                  <p class="text-xs text-gray-500 mt-1">Solo esta persona podrá retirar los insumos de esta solicitud.</p>
                </div>
              </div>

              <!-- Estado actual -->
              <div v-if="request.allow_anyone_to_pickup !== undefined" class="mt-4 p-3 bg-blue-50 rounded-lg">
                <p class="text-sm text-gray-700">
                  <span class="font-medium">Estado actual:</span>
                  <span v-if="request.allow_anyone_to_pickup" class="text-green-700 ml-2">
                    ✓ Cualquier persona puede retirar
                  </span>
                  <span v-else-if="request.authorized_pickup_rut" class="text-orange-700 ml-2">
                    ✓ Solo {{ request.authorized_pickup_name || request.authorized_pickup_rut }} puede retirar
                  </span>
                </p>
              </div>
            </div>
          </div>

          <!-- Items solicitados -->
          <div class="bg-white rounded-lg shadow border">
            <div class="px-4 sm:px-6 py-3 sm:py-4 border-b border-gray-200 flex items-center justify-between">
              <h3 class="text-base sm:text-lg font-semibold text-gray-900">Insumos Solicitados</h3>
              <button
                v-if="hasActiveCart"
                @click="toggleItemsCollapsed"
                class="text-gray-500 hover:text-gray-700"
                :title="itemsCollapsed ? 'Expandir' : 'Colapsar'"
              >
                <svg 
                  class="h-5 w-5 transition-transform duration-200"
                  :class="{ 'transform rotate-180': !itemsCollapsed }"
                  fill="none" 
                  stroke="currentColor" 
                  viewBox="0 0 24 24"
                >
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                </svg>
              </button>
            </div>
            <div v-show="!itemsCollapsed || !hasActiveCart" class="divide-y divide-gray-200">
              <div v-for="(item, index) in items" :key="item.id" class="p-4 sm:p-6">
                <div class="flex flex-col sm:flex-row sm:justify-between sm:items-start gap-2 sm:gap-4 mb-3 sm:mb-4">
                  <div class="flex-1 min-w-0">
                    <div class="flex items-center gap-2 mb-1">
                      <h4 class="text-sm font-medium text-gray-900 truncate">{{ item.supply_name }}</h4>
                      <!-- Badge de estado del item -->
                      <span v-if="item.item_status && item.item_status !== 'pendiente'" :class="getItemStatusBadgeClass(item.item_status)" class="inline-flex px-2 py-0.5 text-xs font-semibold rounded-full whitespace-nowrap">
                        {{ getItemStatusLabel(item.item_status) }}
                      </span>
                    </div>
                    <p class="text-xs sm:text-sm text-gray-500">Código: {{ item.supply_code }}</p>
                  </div>
                  <div class="text-left sm:text-right flex-shrink-0">
                    <p class="text-sm font-medium text-gray-900">
                      Cantidad: {{ item.quantity_requested }}
                    </p>
                    <p v-if="item.quantity_approved" class="text-sm text-green-600">
                      Aprobado: {{ item.quantity_approved }}
                    </p>
                  </div>
                </div>

                <!-- Especificaciones -->
                <div class="grid grid-cols-1 sm:grid-cols-2 gap-3 sm:gap-4 text-xs sm:text-sm">
                  <div v-if="item.size">
                    <span class="font-medium text-gray-700">Tamaño:</span>
                    <span class="ml-1 text-gray-900">{{ item.size }}</span>
                  </div>
                  <div v-if="item.brand">
                    <span class="font-medium text-gray-700">Marca:</span>
                    <span class="ml-1 text-gray-900">{{ item.brand }}</span>
                  </div>
                  <div v-if="item.is_pediatric" class="sm:col-span-2">
                    <span class="inline-flex items-center px-2 py-1 text-xs font-medium bg-purple-100 text-purple-800 rounded-full">
                      <svg class="h-3 w-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" />
                      </svg>
                      Pediátrico
                    </span>
                  </div>
                </div>

                <!-- Especificaciones técnicas -->
                <div v-if="item.specifications" class="mt-3 sm:mt-4">
                  <label class="block text-xs sm:text-sm font-medium text-gray-700">Especificaciones Técnicas</label>
                  <p class="text-xs sm:text-sm text-gray-900 mt-1 bg-gray-50 p-2 rounded">{{ item.specifications }}</p>
                </div>

                <!-- Solicitudes especiales -->
                <div v-if="item.special_requests" class="mt-3 sm:mt-4">
                  <label class="block text-xs sm:text-sm font-medium text-gray-700">Solicitudes Especiales</label>
                  <p class="text-xs sm:text-sm text-gray-900 mt-1 bg-yellow-50 p-2 rounded">{{ item.special_requests }}</p>
                </div>

                <!-- Estado y comentarios del encargado de bodega -->
                <!-- Item Rechazado -->
                <div v-if="item.item_status === 'rechazado' && item.item_notes" class="mt-3 sm:mt-4">
                  <div class="bg-red-50 border border-red-200 rounded-lg p-3">
                    <div class="flex items-start">
                      <svg class="h-5 w-5 text-red-400 mr-2 flex-shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
                      </svg>
                      <div class="flex-1">
                        <label class="block text-xs sm:text-sm font-medium text-red-800 mb-1">
                          Motivo del rechazo
                          <span v-if="item.reviewed_by_name" class="text-xs font-normal text-red-600">
                            (por {{ item.reviewed_by_name }})
                          </span>
                        </label>
                        <p class="text-xs sm:text-sm text-red-900">{{ item.item_notes }}</p>
                        <p v-if="item.reviewed_at" class="text-xs text-red-600 mt-1">
                          {{ formatDate(item.reviewed_at) }}
                        </p>
                      </div>
                    </div>
                  </div>
                </div>

                <!-- Item Devuelto -->
                <div v-if="item.item_status === 'devuelto' && item.item_notes" class="mt-3 sm:mt-4">
                  <div class="bg-orange-50 border border-orange-200 rounded-lg p-3">
                    <div class="flex items-start">
                      <svg class="h-5 w-5 text-orange-400 mr-2 flex-shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                      </svg>
                      <div class="flex-1">
                        <label class="block text-xs sm:text-sm font-medium text-orange-800 mb-1">
                          Motivo de devolución
                          <span v-if="item.reviewed_by_name" class="text-xs font-normal text-orange-600">
                            (por {{ item.reviewed_by_name }})
                          </span>
                        </label>
                        <p class="text-xs sm:text-sm text-orange-900">{{ item.item_notes }}</p>
                        <p v-if="item.reviewed_at" class="text-xs text-orange-600 mt-1">
                          {{ formatDate(item.reviewed_at) }}
                        </p>
                      </div>
                    </div>
                  </div>
                </div>

                <!-- QRs asignados a este item -->
                <div v-if="getItemAssignments(item.id).length > 0" class="mt-3 sm:mt-4">
                  <label class="block text-xs sm:text-sm font-medium text-gray-700 mb-2">QRs Asignados</label>
                  <div class="grid grid-cols-1 gap-2">
                    <div
                      v-for="assignment in getItemAssignments(item.id)"
                      :key="assignment.id"
                      class="bg-blue-50 p-2 sm:p-3 rounded border"
                    >
                      <div class="flex flex-col sm:flex-row sm:justify-between sm:items-start gap-3">
                        <!-- Imagen del QR Code -->
                        <div class="flex items-center gap-3 flex-1">
                          <div class="flex-shrink-0 bg-white p-2 rounded border border-gray-200 shadow-sm">
                            <QrcodeVue 
                              :value="assignment.qr_code" 
                              :size="80" 
                              level="M"
                              render-as="canvas"
                            />
                          </div>
                          <div class="flex-1 min-w-0">
                            <code class="text-xs font-mono break-all block mb-1 font-semibold text-gray-700">{{ assignment.qr_code }}</code>
                            <div v-if="assignment.assigned_date" class="text-xs text-gray-500 mb-2">
                              Asignado: {{ formatDate(assignment.assigned_date) }}
                            </div>
                            <!-- Botones de acción -->
                            <div class="flex flex-wrap gap-2 mt-2">
                              <button
                                v-if="authStore.user?.role === 'encargado de bodega'"
                                @click="goToScannerWithQR(assignment.qr_code)"
                                class="text-xs text-blue-600 hover:text-blue-800 font-medium hover:underline"
                              >
                                Ver trazabilidad →
                              </button>
                            </div>
                          </div>
                        </div>
                        <!-- Badge de estado -->
                        <span :class="getAssignmentStatusBadgeClass(assignment.status)" class="inline-flex px-2 py-1 text-xs font-semibold rounded-full whitespace-nowrap self-start">
                          {{ getAssignmentStatusLabel(assignment.status) }}
                        </span>
                      </div>
                    </div>
                  </div>
                </div>

                <!-- Botón para asignar QR -->
                <div v-if="request.status === 'approved' && getItemAssignments(item.id).length < item.quantity_approved && authStore.canApproveRequests" class="mt-4">
                  <button
                    @click="openAssignQRModal(item)"
                    class="inline-flex items-center justify-center px-3 py-2 border border-blue-300 text-xs sm:text-sm font-medium rounded-md text-blue-700 bg-blue-50 hover:bg-blue-100 w-full sm:w-auto"
                  >
                    <svg class="h-4 w-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
                    </svg>
                    Asignar QR
                  </button>
                </div>
              </div>
            </div>
          </div>

          <!-- Carrito de Insumos (mostrar si hay carrito o si la solicitud está aprobada) -->
          <div v-if="['aprobado', 'en_proceso', 'completado', 'asignado_bodega', 'parcialmente_aprobado', 'devuelto_al_encargado'].includes(request.status)">
            <SupplyCart 
              :request-id="request.id"
              :can-close="authStore.canApproveRequests"
              :can-remove-items="authStore.canApproveRequests"
              @cart-loaded="onCartLoaded"
              @cart-closed="onCartClosed"
              @item-removed="onItemRemoved"
              @transfer-to-pavilion="handleTransferCartToPavilion"
            />
          </div>

          <!-- Cronología de la solicitud -->
          <div class="bg-white rounded-lg shadow border p-6">
            <h3 class="text-lg font-semibold text-gray-900 mb-4">Cronología del Proceso</h3>
            <div class="flow-root">
              <ul class="-mb-8">
                <!-- Solicitud creada -->
                <li>
                  <div class="relative pb-8">
                    <span 
                      v-if="request.assigned_date || request.approval_date || request.status === 'rechazado' || request.status === 'devuelto' || request.status === 'devuelto_al_encargado' || request.completed_date || parseAllComments().some(c => c.type === 'devolución' || c.type === 'reenvío')" 
                      class="absolute top-4 left-4 -ml-px h-full w-0.5 bg-gray-200" 
                      aria-hidden="true"
                    ></span>
                    <div class="relative flex space-x-3">
                      <div>
                        <span class="h-8 w-8 rounded-full bg-blue-500 flex items-center justify-center ring-8 ring-white">
                          <svg class="h-4 w-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
                          </svg>
                        </span>
                      </div>
                      <div class="min-w-0 flex-1 pt-1.5">
                        <p class="text-sm font-medium text-gray-900">Solicitud creada</p>
                        <p class="text-sm text-gray-500">{{ formatDate(request.request_date) }}</p>
                        <p class="text-xs text-gray-400">por {{ request.requested_by_name }}</p>
                        <p class="text-xs mt-1 italic" :class="getOriginalNotes() && getOriginalNotes().trim() === '' ? 'text-gray-400' : 'text-gray-600'">
                          "{{ getOriginalNotes() && getOriginalNotes().trim() === '' ? 'No hay observaciones' : getOriginalNotes() }}"
                        </p>
                      </div>
                    </div>
                  </div>
                </li>

                <!-- Asignación por Pavedad (si existe) -->
                <li v-if="request.assigned_date">
                  <div class="relative pb-8">
                    <span 
                      v-if="request.approval_date || request.status === 'rechazado' || request.status === 'devuelto' || request.status === 'devuelto_al_encargado' || request.status === 'cancelado' || request.completed_date || parseAllComments().some(c => c.type === 'devolución' || c.type === 'reenvío')" 
                      class="absolute top-4 left-4 -ml-px h-full w-0.5 bg-gray-200" 
                      aria-hidden="true"
                    ></span>
                    <div class="relative flex space-x-3">
                      <div>
                        <span class="h-8 w-8 rounded-full bg-purple-500 flex items-center justify-center ring-8 ring-white">
                          <svg class="h-4 w-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 5v2m0 4v2m0 4v2M5 5a2 2 0 00-2 2v3a2 2 0 110 4v3a2 2 0 002 2h14a2 2 0 002-2v-3a2 2 0 110-4V7a2 2 0 00-2-2H5z" />
                          </svg>
                        </span>
                      </div>
                      <div class="min-w-0 flex-1 pt-1.5">
                        <p class="text-sm font-medium text-gray-900">Asignada</p>
                        <p class="text-sm text-gray-500">{{ formatDate(request.assigned_date) }}</p>
                        <p class="text-xs text-gray-400">por {{ request.assigned_by_pavedad_name }}</p>
                        <p class="text-xs text-gray-600 mt-1">Asignada a: {{ request.assigned_to_name }}</p>
                        <template v-for="comment in parseAllComments().filter(c => c.type === 'asignación')" :key="comment.date">
                          <p class="text-xs mt-1 italic" :class="comment.content.trim() === '' ? 'text-gray-400' : 'text-gray-600'">
                            "{{ comment.content.trim() === '' ? 'No hay observaciones' : comment.content }}"
                          </p>
                        </template>
                      </div>
                    </div>
                  </div>
                </li>

                <!-- Rechazada (si existe) -->
                <li v-if="request.status === 'rechazado'">
                  <div class="relative pb-8">
                    <span v-if="request.completed_date" class="absolute top-4 left-4 -ml-px h-full w-0.5 bg-gray-200" aria-hidden="true"></span>
                    <div class="relative flex space-x-3">
                      <div>
                        <span class="h-8 w-8 rounded-full bg-red-500 flex items-center justify-center ring-8 ring-white">
                          <svg class="h-4 w-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                          </svg>
                        </span>
                      </div>
                      <div class="min-w-0 flex-1 pt-1.5">
                        <p class="text-sm font-medium text-gray-900">Solicitud rechazada</p>
                      </div>
                    </div>
                  </div>
                </li>

                <!-- Items devueltos (si existe) -->
                <li v-if="request.status === 'devuelto' || request.status === 'devuelto_al_encargado' || parseAllComments().some(c => c.type === 'devolución')">
                  <div class="relative pb-8">
                    <span 
                      v-if="request.status === 'devuelto_al_encargado' || request.status === 'cancelado' || request.completed_date || parseAllComments().some(c => c.type === 'reenvío') || (request.approval_date && !['rechazado', 'devuelto', 'devuelto_al_encargado'].includes(request.status))" 
                      class="absolute top-4 left-4 -ml-px h-full w-0.5 bg-gray-200" 
                      aria-hidden="true"
                    ></span>
                    <div class="relative flex space-x-3">
                      <div>
                        <span class="h-8 w-8 rounded-full bg-orange-500 flex items-center justify-center ring-8 ring-white">
                          <svg class="h-4 w-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 10h10a8 8 0 018 8v2M3 10l6 6m-6-6l6-6" />
                          </svg>
                        </span>
                      </div>
                      <div class="min-w-0 flex-1 pt-1.5">
                        <p class="text-sm font-medium text-gray-900">Items devueltos</p>
                        <template v-for="comment in parseAllComments().filter(c => c.type === 'devolución')" :key="comment.date">
                          <div class="mt-2">
                            <p class="text-sm text-gray-500">{{ comment.date }}</p>
                            <p class="text-xs text-gray-400">por {{ comment.author }}</p>
                            <p class="text-xs mt-1 italic text-gray-600 whitespace-pre-wrap">"{{ comment.content }}"</p>
                          </div>
                        </template>
                      </div>
                    </div>
                  </div>
                </li>

                <!-- Reenviado al encargado (si existe) -->
                <li v-if="request.status === 'devuelto_al_encargado' || parseAllComments().some(c => c.type === 'reenvío')">
                  <div class="relative pb-8">
                    <span 
                      v-if="request.status === 'cancelado' || request.completed_date || (request.approval_date && !['rechazado', 'devuelto', 'devuelto_al_encargado'].includes(request.status))" 
                      class="absolute top-4 left-4 -ml-px h-full w-0.5 bg-gray-200" 
                      aria-hidden="true"
                    ></span>
                    <div class="relative flex space-x-3">
                      <div>
                        <span class="h-8 w-8 rounded-full bg-teal-500 flex items-center justify-center ring-8 ring-white">
                          <svg class="h-4 w-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                          </svg>
                        </span>
                      </div>
                      <div class="min-w-0 flex-1 pt-1.5">
                        <p class="text-sm font-medium text-gray-900">Reenviado al encargado</p>
                        <template v-for="comment in parseAllComments().filter(c => c.type === 'reenvío')" :key="comment.date">
                          <div class="mt-2">
                            <p class="text-sm text-gray-500">{{ comment.date }}</p>
                            <p class="text-xs text-gray-400">por {{ comment.author }}</p>
                            <p class="text-xs mt-1 italic text-gray-600 whitespace-pre-wrap">"{{ comment.content }}"</p>
                          </div>
                        </template>
                      </div>
                    </div>
                  </div>
                </li>

                <!-- Aprobada (si existe) - AHORA AL FINAL -->
                <li v-if="request.approval_date && !['rechazado', 'devuelto', 'devuelto_al_encargado'].includes(request.status)">
                  <div class="relative pb-8">
                    <span 
                      v-if="request.completed_date" 
                      class="absolute top-4 left-4 -ml-px h-full w-0.5 bg-gray-200" 
                      aria-hidden="true"
                    ></span>
                    <div class="relative flex space-x-3">
                      <div>
                        <span class="h-8 w-8 rounded-full bg-green-500 flex items-center justify-center ring-8 ring-white">
                          <svg class="h-4 w-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                          </svg>
                        </span>
                      </div>
                      <div class="min-w-0 flex-1 pt-1.5">
                        <p class="text-sm font-medium text-gray-900">
                          {{ request.status === 'parcialmente_aprobado' ? 'Solicitud parcialmente aprobada' : 'Solicitud aprobada' }}
                        </p>
                        <p class="text-sm text-gray-500">{{ formatDate(request.approval_date) }}</p>
                        <p class="text-xs text-gray-400">por {{ request.approved_by_name }}</p>
                      </div>
                    </div>
                  </div>
                </li>
              </ul>
            </div>
          </div>
        </div>

        <!-- Panel lateral -->
        <div class="space-y-4 sm:space-y-6">
          <!-- Resumen -->
          <div class="bg-white rounded-lg shadow border p-4 sm:p-6">
            <h3 class="text-base sm:text-lg font-semibold text-gray-900 mb-3 sm:mb-4">Resumen</h3>
            <div class="space-y-2 sm:space-y-3">
              <div class="flex justify-between">
                <span class="text-sm text-gray-600">Items Solicitados:</span>
                <span class="text-sm font-medium text-gray-900">{{ items.length }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-sm text-gray-600">Items Aprobados:</span>
                <span class="text-sm font-medium text-green-600">{{ getApprovedItemsCount() }}</span>
              </div>
              <!--<div class="flex justify-between">
                <span class="text-sm text-gray-600">QRs Asignados:</span>
                <span class="text-sm font-medium text-blue-600">{{ assignments.length }}</span>
              </div>-->
              <!--<div class="flex justify-between">
                <span class="text-sm text-gray-600">Progreso:</span>
                <span class="text-sm font-medium text-blue-600">{{ getProgressPercentage() }}%</span>
              </div>-->
            </div>
            
            <!-- Barra de progreso 
            <div class="mt-4">
              <div class="bg-gray-200 rounded-full h-2">
                <div
                  class="bg-blue-600 h-2 rounded-full transition-all duration-300"
                  :style="`width: ${getProgressPercentage()}%`"
                ></div>
              </div>
            </div>-->

            <!-- Estado y Acciones -->
            <div class="mt-4 pt-4 border-t border-gray-200">
              <h4 class="text-sm font-semibold text-gray-900 mb-3">Estado y Acciones</h4>
              
              <!-- Estado actual -->
              <div class="mb-3">
                <label class="block text-xs font-medium text-gray-700 mb-1">Estado Actual</label>
                <span :class="getStatusBadgeClass(request.status)" class="inline-flex px-2.5 py-1 text-xs font-semibold rounded-full">
                  {{ getStatusLabel(request.status) }}
                </span>
              </div>

              <!-- Acciones disponibles según rol y estado -->
              <div class="space-y-2">
                <label class="block text-xs font-medium text-gray-700">Acciones Disponibles</label>
                
                <!-- Pavedad puede asignar -->
                <button
                  v-if="request.status === 'pendiente_pavedad' && authStore.isPavedad"
                  @click="openAssignModal"
                  class="w-full inline-flex items-center justify-center px-3 py-2 border border-purple-300 rounded-md text-xs font-medium text-purple-700 bg-purple-50 hover:bg-purple-100"
                >
                  <svg class="h-4 w-4 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 5v2m0 4v2m0 4v2M5 5a2 2 0 00-2 2v3a2 2 0 110 4v3a2 2 0 002 2h14a2 2 0 002-2v-3a2 2 0 110-4V7a2 2 0 00-2-2H5z" />
                  </svg>
                  Asignar
                </button>

                <!-- Encargado de bodega puede revisar items -->
                <button
                  v-if="(request.status === 'asignado_bodega' || request.status === 'en_proceso' || request.status === 'devuelto_al_encargado') && authStore.isWarehouseManager && request.assigned_to === authStore.getUserRut"
                  @click="openReviewItemsModal"
                  class="w-full inline-flex items-center justify-center px-3 py-2 border border-blue-300 rounded-md text-xs font-medium text-blue-700 bg-blue-50 hover:bg-blue-100"
                >
                  <svg class="h-4 w-4 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4" />
                  </svg>
                  Revisar Insumos
                </button>

                <!-- Botón Listo para retiro (solo para encargado de bodega cuando el carrito está activo) -->
                <button
                  v-if="currentCart && currentCart.status === 'active' && canTransferCartToPavilion"
                  @click="handleTransferCartToPavilion"
                  class="w-full inline-flex items-center justify-center px-3 py-2 border border-green-300 rounded-md text-xs font-medium text-green-700 bg-green-50 hover:bg-green-100"
                >
                  <svg class="h-4 w-4 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  Listo para retiro
                </button>

                <!-- Doctor puede editar y reenviar solicitud devuelta -->
                <button
                  v-if="hasReturnedItems && (request.status === 'devuelto' || request.status === 'parcialmente_aprobado') && (authStore.isDoctor || authStore.isNurse) && request.requested_by === authStore.getUserRut"
                  @click="goToEditRequest"
                  class="w-full inline-flex items-center justify-center px-3 py-2 border border-orange-300 rounded-md text-xs font-medium text-orange-700 bg-orange-50 hover:bg-orange-100"
                >
                  <svg class="h-4 w-4 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                  </svg>
                  Editar y Reenviar
                </button>

                <!-- Si ya está procesada -->
                <div v-if="['aprobado', 'rechazado', 'parcialmente_aprobado', 'completado'].includes(request.status)" class="p-2 bg-gray-50 rounded text-center">
                  <p class="text-xs text-gray-600">
                    Solicitud procesada
                  </p>
                </div>
              </div>
            </div>

            <!-- Comentarios y Observaciones -->
            <div class="mt-4 pt-4 border-t border-gray-200">
              <h4 class="text-sm font-semibold text-gray-900 mb-3">Comentarios y Observaciones</h4>
              
              <div class="space-y-3">
                <!-- Renderizar todos los comentarios parseados -->
                <div 
                  v-for="(comment, index) in parseAllComments()" 
                  :key="index"
                  :class="[
                    'p-2 rounded-lg',
                    comment.type === 'solicitante' ? 'bg-blue-50' : 
                    comment.type === 'asignación' ? 'bg-purple-50' : 
                    comment.type === 'reenvío' ? 'bg-green-50' : 'bg-orange-50'
                  ]"
                >
                  <div class="flex items-start">
                    <!-- Icono según tipo -->
                    <svg 
                      :class="[
                        'h-4 w-4 mt-0.5 mr-1.5 flex-shrink-0',
                        comment.type === 'solicitante' ? 'text-blue-600' :
                        comment.type === 'asignación' ? 'text-purple-600' :
                        comment.type === 'reenvío' ? 'text-green-600' : 'text-orange-600'
                      ]"
                      fill="none" 
                      stroke="currentColor" 
                      viewBox="0 0 24 24"
                    >
                      <path 
                        v-if="comment.type === 'solicitante'"
                        stroke-linecap="round" 
                        stroke-linejoin="round" 
                        stroke-width="2" 
                        d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" 
                      />
                      <path 
                        v-else-if="comment.type === 'asignación'"
                        stroke-linecap="round" 
                        stroke-linejoin="round" 
                        stroke-width="2" 
                        d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" 
                      />
                      <path 
                        v-else-if="comment.type === 'reenvío'"
                        stroke-linecap="round" 
                        stroke-linejoin="round" 
                        stroke-width="2" 
                        d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" 
                      />
                      <path 
                        v-else
                        stroke-linecap="round" 
                        stroke-linejoin="round" 
                        stroke-width="2" 
                        d="M3 10h10a8 8 0 018 8v2M3 10l6 6m-6-6l6-6" 
                      />
                    </svg>
                    <div class="flex-1">
                      <div class="flex items-baseline justify-between">
                        <p class="text-xs font-medium text-gray-700">{{ comment.author }}</p>
                        <p v-if="comment.date" class="text-xs text-gray-500 ml-2">{{ comment.date }}</p>
                      </div>
                      <p class="text-xs whitespace-pre-wrap mt-0.5" :class="comment.content.trim() === '' ? 'text-gray-400 italic' : 'text-gray-900'">
                        {{ comment.content.trim() === '' ? 'No hay observaciones' : comment.content }}
                      </p>
                    </div>
                  </div>
                </div>

                <!-- Si no hay comentarios -->
                <div v-if="parseAllComments().length === 0" class="p-3 text-center text-gray-400 text-xs">
                  <svg class="h-6 w-6 mx-auto mb-1 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
                  </svg>
                  <p>No hay comentarios u observaciones para esta solicitud</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Modal para asignar QR -->
    <div v-if="showAssignQRModal" class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50 px-4">
      <div class="relative top-10 sm:top-20 mx-auto p-4 sm:p-5 border w-full max-w-md shadow-lg rounded-md bg-white my-8">
        <div class="mt-2 sm:mt-3">
          <h3 class="text-base sm:text-lg font-medium text-gray-900 mb-3 sm:mb-4">Asignar QR a Insumo</h3>
          
          <div class="mb-3 sm:mb-4">
            <label class="block text-xs sm:text-sm font-medium text-gray-700">Insumo:</label>
            <p class="text-sm text-gray-900 truncate">{{ selectedItem?.supply_name }}</p>
          </div>
          
          <div class="mb-3 sm:mb-4">
            <label for="qrCode" class="block text-xs sm:text-sm font-medium text-gray-700 mb-1">
              Código QR <span class="text-red-500">*</span>
            </label>
            <input
              id="qrCode"
              type="text"
              v-model="qrAssignmentForm.qrCode"
              placeholder="Escanear o escribir código QR"
              class="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            />
          </div>
          
          <div class="mb-3 sm:mb-4">
            <label for="assignedBy" class="block text-xs sm:text-sm font-medium text-gray-700 mb-1">
              Asignado por
            </label>
            <input
              id="assignedBy"
              type="text"
              v-model="qrAssignmentForm.assignedByName"
              placeholder="Nombre completo"
              class="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            />
          </div>
          
          <div class="mb-4 sm:mb-6">
            <label for="notes" class="block text-xs sm:text-sm font-medium text-gray-700 mb-1">
              Notas
            </label>
            <textarea
              id="notes"
              v-model="qrAssignmentForm.notes"
              rows="3"
              placeholder="Notas adicionales (opcional)"
              class="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            ></textarea>
          </div>
          
          <div class="flex flex-col sm:flex-row sm:justify-end gap-2 sm:gap-3">
            <button
              @click="closeAssignQRModal"
              type="button"
              class="w-full sm:w-auto px-4 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 hover:bg-gray-50 order-2 sm:order-1"
            >
              Cancelar
            </button>
            <button
              @click="assignQR"
              :disabled="!qrAssignmentForm.qrCode || assigningQR"
              type="button"
              class="btn-primary disabled:opacity-50 disabled:cursor-not-allowed order-1 sm:order-2"
            >
              {{ assigningQR ? 'Asignando...' : 'Asignar' }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Modal de Asignación -->
    <AssignRequestModal
      :show="showAssignModal"
      :request="request"
      @close="showAssignModal = false"
    />

    <!-- Modal de Revisión de Items -->
    <ReviewItemsModal
      :show="showReviewItemsModal"
      :request="request"
      @close="showReviewItemsModal = false"
      @itemsReviewed="handleItemsReviewed"
    />

    <!-- Modales personalizados -->
    <!-- Modal de confirmación -->
    <div v-if="showConfirmModal" class="fixed inset-0 z-50 overflow-y-auto" aria-labelledby="modal-title" role="dialog" aria-modal="true">
      <div class="flex items-center justify-center min-h-screen px-4 pt-4 pb-20 text-center sm:block sm:p-0">
        <div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" aria-hidden="true" @click="closeConfirmModal"></div>
        <span class="hidden sm:inline-block sm:align-middle sm:h-screen" aria-hidden="true">&#8203;</span>
        <div class="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full">
          <div class="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
            <div class="sm:flex sm:items-start">
              <div class="mx-auto flex-shrink-0 flex items-center justify-center h-12 w-12 rounded-full bg-blue-100 sm:mx-0 sm:h-10 sm:w-10">
                <svg class="h-6 w-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </div>
              <div class="mt-3 text-center sm:mt-0 sm:ml-4 sm:text-left flex-1">
                <h3 class="text-lg leading-6 font-medium text-gray-900" id="modal-title">
                  {{ confirmModalTitle }}
                </h3>
                <div class="mt-2">
                  <p class="text-sm text-gray-500">{{ confirmModalMessage }}</p>
                </div>
              </div>
            </div>
          </div>
          <div class="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse gap-2">
            <button
              type="button"
              @click="executeConfirmCallback"
              class="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-blue-600 text-base font-medium text-white hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 sm:ml-3 sm:w-auto sm:text-sm"
            >
              Confirmar
            </button>
            <button
              type="button"
              @click="closeConfirmModal"
              class="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm"
            >
              Cancelar
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Modal de mensaje -->
    <div v-if="showMessageModal" class="fixed inset-0 z-50 overflow-y-auto" aria-labelledby="modal-title" role="dialog" aria-modal="true">
      <div class="flex items-center justify-center min-h-screen px-4 pt-4 pb-20 text-center sm:block sm:p-0">
        <div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" aria-hidden="true" @click="closeMessageModal"></div>
        <span class="hidden sm:inline-block sm:align-middle sm:h-screen" aria-hidden="true">&#8203;</span>
        <div class="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full">
          <div class="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
            <div class="sm:flex sm:items-start">
              <div class="mx-auto flex-shrink-0 flex items-center justify-center h-12 w-12 rounded-full sm:mx-0 sm:h-10 sm:w-10"
                :class="messageModalType === 'success' ? 'bg-green-100' : messageModalType === 'error' ? 'bg-red-100' : 'bg-blue-100'">
                <svg v-if="messageModalType === 'success'" class="h-6 w-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                </svg>
                <svg v-else-if="messageModalType === 'error'" class="h-6 w-6 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
                <svg v-else class="h-6 w-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </div>
              <div class="mt-3 text-center sm:mt-0 sm:ml-4 sm:text-left flex-1">
                <h3 class="text-lg leading-6 font-medium text-gray-900" id="modal-title">
                  {{ messageModalTitle }}
                </h3>
                <div class="mt-2">
                  <p class="text-sm text-gray-500">{{ messageModalMessage }}</p>
                </div>
              </div>
            </div>
          </div>
          <div class="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
            <button
              type="button"
              @click="closeMessageModal"
              class="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-blue-600 text-base font-medium text-white hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 sm:ml-3 sm:w-auto sm:text-sm"
            >
              Aceptar
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Modal de input (para motivo de rechazo) -->
    <div v-if="showInputModal" class="fixed inset-0 z-50 overflow-y-auto" aria-labelledby="modal-title" role="dialog" aria-modal="true">
      <div class="flex items-center justify-center min-h-screen px-4 pt-4 pb-20 text-center sm:block sm:p-0">
        <div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" aria-hidden="true" @click="closeInputModal"></div>
        <span class="hidden sm:inline-block sm:align-middle sm:h-screen" aria-hidden="true">&#8203;</span>
        <div class="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full">
          <div class="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
            <div class="sm:flex sm:items-start">
              <div class="mx-auto flex-shrink-0 flex items-center justify-center h-12 w-12 rounded-full bg-red-100 sm:mx-0 sm:h-10 sm:w-10">
                <svg class="h-6 w-6 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </div>
              <div class="mt-3 text-center sm:mt-0 sm:ml-4 sm:text-left flex-1 w-full">
                <h3 class="text-lg leading-6 font-medium text-gray-900" id="modal-title">
                  {{ inputModalTitle }}
                </h3>
                <div class="mt-4">
                  <label class="block text-sm font-medium text-gray-700 mb-2">{{ inputModalLabel }}</label>
                  <input
                    v-model="inputModalValue"
                    type="text"
                    class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                    :placeholder="inputModalPlaceholder"
                    @keyup.enter="confirmInputModal"
                  />
                  <p v-if="inputModalError" class="mt-1 text-sm text-red-600">{{ inputModalError }}</p>
                </div>
              </div>
            </div>
          </div>
          <div class="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse gap-2">
            <button
              type="button"
              @click="confirmInputModal"
              :disabled="!inputModalValue.trim()"
              class="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-red-600 text-base font-medium text-white hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 disabled:opacity-50 disabled:cursor-not-allowed sm:ml-3 sm:w-auto sm:text-sm"
            >
              {{ inputModalConfirmText }}
            </button>
            <button
              type="button"
              @click="closeInputModal"
              class="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm"
            >
              Cancelar
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import supplyRequestService from '@/services/requests/supplyRequestService'
import pavilionService from '@/services/config/pavilionService'
import { userService } from '@/services/common/userService'
import cartService from '@/services/requests/cartService'
import { format } from 'date-fns'
import { es } from 'date-fns/locale'
import AssignRequestModal from '@/components/requests/AssignRequestModal.vue'
import ReviewItemsModal from '@/components/requests/ReviewItemsModal.vue'
import QrcodeVue from 'qrcode.vue'
import SupplyCart from '@/components/requests/SupplyCart.vue'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

// Estado reactivo
const loading = ref(false)
const processing = ref(false)
const assigningQR = ref(false)
const error = ref(null)
const request = ref(null)
const items = ref([])
const assignments = ref([])
const pavilions = ref([])
const selectedQRTraceability = ref(null)
const pickupConfig = ref({
  allow_anyone_to_pickup: true,
  authorized_pickup_rut: null,
  authorized_pickup_name: null
})
const savingPickupConfig = ref(false)
const userSearchQuery = ref('')
const allUsers = ref([])
const filteredUsers = ref([])
const showUserSuggestions = ref(false)
const showAssignQRModal = ref(false)
const selectedItem = ref(null)
const showAssignModal = ref(false)
const showReviewItemsModal = ref(false)
const hasActiveCart = ref(false)
const itemsCollapsed = ref(false)
const currentCart = ref(null)
const canTransferCartToPavilion = ref(false)

// Estados para modales personalizados
const showConfirmModal = ref(false)
const confirmModalTitle = ref('')
const confirmModalMessage = ref('')
const confirmModalCallback = ref(null)

const showMessageModal = ref(false)
const messageModalType = ref('info') // 'success', 'error', 'info'
const messageModalTitle = ref('')
const messageModalMessage = ref('')

const showInputModal = ref(false)
const inputModalTitle = ref('')
const inputModalLabel = ref('')
const inputModalPlaceholder = ref('')
const inputModalValue = ref('')
const inputModalError = ref('')
const inputModalConfirmText = ref('Confirmar')
const inputModalCallback = ref(null)

// Formulario de asignación QR
const qrAssignmentForm = reactive({
  qrCode: '',
  assignedBy: 'ADMIN',
  assignedByName: 'Sistema Admin',
  notes: ''
})

// Computed
const requestId = computed(() => parseInt(route.params.id))

// Verificar si hay items devueltos en la solicitud
const hasReturnedItems = computed(() => {
  return items.value.some(item => item.item_status === 'devuelto')
})

// Métodos principales
const loadSupplyRequest = async () => {
  loading.value = true
  error.value = null

  try {
    const result = await supplyRequestService.getSupplyRequestById(requestId.value)
    
    if (result.success && result.data) {
      request.value = result.data.request
      items.value = result.data.items || []
      assignments.value = result.data.assignments || []
      
      // Inicializar configuración de retiro
      if (request.value) {
        pickupConfig.value = {
          allow_anyone_to_pickup: request.value.allow_anyone_to_pickup !== undefined ? request.value.allow_anyone_to_pickup : true,
          authorized_pickup_rut: request.value.authorized_pickup_rut || null,
          authorized_pickup_name: request.value.authorized_pickup_name || null
        }
        // Si hay un usuario autorizado, mostrar su nombre en el campo de búsqueda
        if (request.value.authorized_pickup_rut && request.value.authorized_pickup_name) {
          userSearchQuery.value = request.value.authorized_pickup_name
        }
      }
      
      console.log('Solicitud cargada:', result.data)
      console.log('Estado:', request.value.status)
      console.log('Fecha de aprobación:', request.value.approval_date)
      console.log('Aprobado por:', request.value.approved_by_name)
    } else {
      error.value = result.error || 'Solicitud no encontrada'
    }
  } catch (err) {
    console.error('Error cargando solicitud:', err)
    error.value = 'Error al conectar con el servidor'
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

const approveRequest = async () => {
  openConfirmModal(
    'Aprobar Solicitud',
    '¿Está seguro de aprobar esta solicitud?',
    async () => {
      processing.value = true
      try {
        const approvalData = {
          approved_by: 'ADMIN',
          approved_by_name: 'Sistema Admin',
          approval_notes: 'Aprobada desde interfaz web'
        }
        
        await supplyRequestService.approveSupplyRequest(requestId.value, approvalData)
        await loadSupplyRequest()
        showMessage('success', 'Solicitud Aprobada', 'La solicitud ha sido aprobada exitosamente')
      } catch (err) {
        console.error('Error aprobando solicitud:', err)
        showMessage('error', 'Error al aprobar la solicitud', err.response?.data?.error || err.message)
      } finally {
        processing.value = false
      }
    }
  )
}

const rejectRequest = async () => {
  openInputModal(
    'Motivo del Rechazo',
    'Ingrese el motivo del rechazo:',
    'Motivo...',
    async (reason) => {
      if (!reason || !reason.trim()) {
        inputModalError.value = 'Debe ingresar un motivo'
        return
      }
      
      processing.value = true
      try {
        const rejectionData = {
          rejected_by: 'ADMIN',
          rejected_by_name: 'Sistema Admin',
          notes: reason
        }
        
        await supplyRequestService.rejectSupplyRequest(requestId.value, rejectionData)
        await loadSupplyRequest()
        showMessage('success', 'Solicitud Rechazada', 'La solicitud ha sido rechazada exitosamente')
      } catch (err) {
        console.error('Error rechazando solicitud:', err)
        showMessage('error', 'Error al rechazar la solicitud', err.response?.data?.error || err.message)
      } finally {
        processing.value = false
      }
    },
    'Rechazar'
  )
}

const openAssignModal = () => {
  showAssignModal.value = true
}

const openReviewItemsModal = () => {
  showReviewItemsModal.value = true
}

const handleItemsReviewed = async () => {
  // NO cerrar el modal aquí - el modal se cerrará automáticamente cuando todos los items estén resueltos
  // Solo recargar la solicitud para obtener el estado actualizado y mostrar el carrito si se aprobó
  await new Promise(resolve => setTimeout(resolve, 500))
  await loadSupplyRequest()
  console.log('Solicitud recargada después de revisar items. Estado:', request.value?.status)
}

const openEditReturnedModal = () => {
  showEditReturnedModal.value = true
}

const handleRequestResubmitted = async () => {
  showEditReturnedModal.value = false
  await loadSupplyRequest()
}

const goToEditRequest = async () => {
  // Obtener el ID de la solicitud actual
  const id = request.value?.id || requestId.value
  
  if (!id) {
    console.error('No se pudo obtener el ID de la solicitud')
    showMessage('error', 'Error', 'No se pudo determinar el ID de la solicitud')
    return
  }
  
  console.log('Navegando a editar solicitud con ID:', id)
  console.log('Nombre de ruta:', 'EditSupplyRequest')
  console.log('Params:', { id: id.toString() })
  
  try {
    await router.push({
      name: 'EditSupplyRequest',
      params: { id: id.toString() }
    })
    console.log('Navegación exitosa')
  } catch (error) {
    console.error('Error en navegación:', error)
    showMessage('error', 'Error de Navegación', 'No se pudo navegar a la página de edición: ' + error.message)
  }
}

const goToScannerWithQR = (qrCode) => {
  // Navegar al escáner QR con el código pre-escaneado
  router.push({
    name: 'QRScanner',
    query: { 
      qr: qrCode // El componente QRScanner espera el parámetro 'qr'
    }
  })
}

const viewQRTraceability = async (qrCode) => {
  try {
    const result = await supplyRequestService.getQRTraceability(qrCode)
    if (result.success) {
      selectedQRTraceability.value = result.data
    }
  } catch (err) {
    console.error('Error obteniendo trazabilidad:', err)
  }
}

const openAssignQRModal = (item) => {
  selectedItem.value = item
  showAssignQRModal.value = true
  // Reset form
  Object.assign(qrAssignmentForm, {
    qrCode: '',
    assignedBy: 'ADMIN',
    assignedByName: 'Sistema Admin',
    notes: ''
  })
}

const closeAssignQRModal = () => {
  showAssignQRModal.value = false
  selectedItem.value = null
}

const assignQR = async () => {
  if (!qrAssignmentForm.qrCode || !selectedItem.value) return

  assigningQR.value = true
  try {
    const assignmentData = {
      supply_request_id: requestId.value,
      supply_request_item_id: selectedItem.value.id,
      qr_code: qrAssignmentForm.qrCode,
      assigned_by: qrAssignmentForm.assignedBy,
      assigned_by_name: qrAssignmentForm.assignedByName,
      notes: qrAssignmentForm.notes
    }
    
    await supplyRequestService.assignQRToRequest(assignmentData)
    await loadSupplyRequest()
    closeAssignQRModal()
    showMessage('success', 'QR Asignado', 'El código QR ha sido asignado exitosamente')
  } catch (err) {
    console.error('Error asignando QR:', err)
    showMessage('error', 'Error al asignar QR', err.response?.data?.error || err.message)
  } finally {
    assigningQR.value = false
  }
}

// Métodos auxiliares
const getItemAssignments = (itemId) => {
  return assignments.value.filter(assignment => assignment.supply_request_item_id === itemId)
}

const getDeliveredAssignments = () => {
  return assignments.value.filter(assignment => assignment.status === 'delivered')
}

const getProgressPercentage = () => {
  const totalApproved = getTotalApprovedQuantity()
  const assignedCount = assignments.value.length
  return totalApproved > 0 ? Math.round((assignedCount / totalApproved) * 100) : 0
}

const getTotalRequestedQuantity = () => {
  return items.value.reduce((sum, item) => sum + item.quantity_requested, 0)
}

const getTotalApprovedQuantity = () => {
  return items.value.reduce((sum, item) => {
    // Si el item está aceptado y tiene quantity_approved, usar ese valor
    if (item.quantity_approved && item.quantity_approved > 0) {
      return sum + item.quantity_approved
    }
    // Si el item está aceptado pero no tiene quantity_approved, usar quantity_requested
    if (item.item_status === 'aceptado') {
      return sum + item.quantity_requested
    }
    return sum
  }, 0)
}

const getApprovedItemsCount = () => {
  // Contar cuántos items (tipos diferentes) están aprobados/aceptados
  return items.value.filter(item => 
    item.item_status === 'aceptado' || 
    (item.quantity_approved && item.quantity_approved > 0)
  ).length
}

const getSurgeryUrgencyClass = (surgeryDatetime) => {
  if (!surgeryDatetime) return 'text-gray-500'
  
  const surgeryDate = new Date(surgeryDatetime)
  const now = new Date()
  const hoursUntilSurgery = (surgeryDate.getTime() - now.getTime()) / (1000 * 60 * 60)
  
  if (hoursUntilSurgery < 0) return 'text-gray-500' // Ya pasó
  if (hoursUntilSurgery <= 4) return 'text-red-600' // Muy urgente
  if (hoursUntilSurgery <= 12) return 'text-orange-600' // Urgente
  if (hoursUntilSurgery <= 24) return 'text-yellow-600' // Próximo
  return 'text-green-600' // Programado
}

const getSurgeryUrgencyText = (surgeryDatetime) => {
  if (!surgeryDatetime) return 'No programada'
  
  const surgeryDate = new Date(surgeryDatetime)
  const now = new Date()
  const hoursUntilSurgery = (surgeryDate.getTime() - now.getTime()) / (1000 * 60 * 60)
  
  if (hoursUntilSurgery < 0) return 'Completada'
  if (hoursUntilSurgery <= 4) return `En ${Math.round(hoursUntilSurgery)}h - MUY URGENTE`
  if (hoursUntilSurgery <= 12) return `En ${Math.round(hoursUntilSurgery)}h - Urgente`
  if (hoursUntilSurgery <= 24) return `En ${Math.round(hoursUntilSurgery)}h - Próxima`
  
  const daysUntilSurgery = Math.round(hoursUntilSurgery / 24)
  return `En ${daysUntilSurgery} día${daysUntilSurgery > 1 ? 's' : ''}`
}

const getPavilionName = (pavilionId) => {
  const pavilion = pavilions.value.find(p => p.id === pavilionId)
  return pavilion ? pavilion.name : `Pabellón ${pavilionId}`
}

const getMedicalCenterName = () => {
  const pavilion = pavilions.value.find(p => p.id === request.value?.pavilion_id)
  return pavilion?.medical_center?.name || 'Centro Médico'
}

const formatDate = (dateString) => {
  if (!dateString) return 'N/A'
  try {
    return format(new Date(dateString), 'dd/MM/yyyy HH:mm', { locale: es })
  } catch {
    return dateString
  }
}

const isSurgeryUrgent = (surgeryDatetime) => {
  if (!surgeryDatetime) return false
  const surgeryDate = new Date(surgeryDatetime)
  const now = new Date()
  const hoursUntilSurgery = (surgeryDate.getTime() - now.getTime()) / (1000 * 60 * 60)
  return hoursUntilSurgery <= 24 && hoursUntilSurgery > 0 // Próxima en las siguientes 24 horas
}

// Métodos de estilo
const getStatusLabel = (status) => supplyRequestService.getStatusLabel(status)
const getPriorityLabel = (priority) => supplyRequestService.getPriorityLabel(priority)

// Parsear comentarios del campo notes y separarlos por tipo usando los campos de la BD
const parseAllComments = () => {
  const comments = []
  
  if (!request.value) return comments
  
  // Dividir el campo notes por doble salto de línea - NO FILTRAR espacios para mantener posiciones
  const notesArray = request.value.notes ? request.value.notes.split('\n\n') : []
  
  if (notesArray.length === 0) return comments
  
  // Construir una línea de tiempo con todos los eventos que pueden tener comentarios
  const timeline = []
  
  // 1. Comentario original del solicitante (request_date) - SIEMPRE existe
  timeline.push({
    type: 'solicitante',
    author: request.value.requested_by_name || 'Solicitante',
    date: request.value.request_date ? new Date(request.value.request_date) : null,
    hasComment: true
  })
  
  // 2. Comentario de asignación por Pavedad (assigned_date)
  if (request.value.assigned_date && request.value.assigned_by_pavedad_name) {
    timeline.push({
      type: 'asignación',
      author: request.value.assigned_by_pavedad_name,
      date: new Date(request.value.assigned_date),
      hasComment: true
    })
  }
  
  // 3. Comentarios de devolución (cuando hay items que fueron devueltos)
  // Verificar si hay algún item que tenga notas de revisión
  const itemsWithNotes = items.value?.filter(item => 
    item.item_notes && item.item_notes.trim() !== ''
  ) || []
  
  // Si hay items con notas de revisión, significa que hubo una devolución
  if (itemsWithNotes.length > 0) {
    // Buscar el item con la fecha de revisión más antigua para la devolución
    const firstReviewedItem = itemsWithNotes.reduce((earliest, item) => {
      if (!earliest || !earliest.reviewed_at) return item
      if (!item.reviewed_at) return earliest
      return new Date(item.reviewed_at) < new Date(earliest.reviewed_at) ? item : earliest
    }, null)
    
    if (firstReviewedItem && firstReviewedItem.reviewed_at) {
      timeline.push({
        type: 'devolución',
        author: firstReviewedItem.reviewed_by_name || 'Encargado de Bodega',
        date: new Date(firstReviewedItem.reviewed_at),
        hasComment: true
      })
    }
  }
  
  // 4. Comentario de reenvío (cuando hubo items devueltos y luego reenviados)
  // Si hay items con notas Y el número de comentarios en notes sugiere que hay un reenvío
  // O si el status actual es devuelto_al_encargado o aprobado (después de devolución)
  const hasDevolucion = itemsWithNotes.length > 0
  const probablyHasReenvio = notesArray.length > timeline.length // Hay más comentarios de los que hemos procesado
  
  if (hasDevolucion && probablyHasReenvio) {
    // Calcular la fecha del reenvío basándose en el contexto
    let reenvioDate = new Date()
    let reenvioAuthor = request.value.requested_by_name || 'Solicitante'
    
    // Si el status es devuelto_al_encargado, usar updated_at
    if (request.value.status === 'devuelto_al_encargado' && request.value.updated_at) {
      reenvioDate = new Date(request.value.updated_at)
    } 
    // Si hay items que fueron devueltos pero ahora están en otro estado, usar la fecha más reciente
    else if (itemsWithNotes.length > 0) {
      const lastReviewedItem = itemsWithNotes.reduce((latest, item) => {
        if (!latest || !latest.reviewed_at) return item
        if (!item.reviewed_at) return latest
        return new Date(item.reviewed_at) > new Date(latest.reviewed_at) ? item : latest
      }, null)
      if (lastReviewedItem && lastReviewedItem.reviewed_at) {
        reenvioDate = new Date(lastReviewedItem.reviewed_at)
      }
    }
    
    timeline.push({
      type: 'reenvío',
      author: reenvioAuthor,
      date: reenvioDate,
      hasComment: true
    })
  }
  
  // Ordenar la línea de tiempo por fecha
  timeline.sort((a, b) => {
    if (!a.date) return -1
    if (!b.date) return 1
    return a.date - b.date
  })
  
  // Asignar los comentarios del array notesArray a la línea de tiempo
  // IMPORTANTE: No filtrar, respetar el orden y las posiciones
  timeline.forEach((event, index) => {
    if (index < notesArray.length) {
      const content = notesArray[index]
      comments.push({
        type: event.type,
        author: event.author,
        content: content, // Mantener el contenido tal cual (puede ser espacio)
        date: formatDate(event.date)
      })
    }
  })
  
  return comments
}

// Funciones de compatibilidad (mantener para no romper el código existente)
const getOriginalNotes = () => {
  const comments = parseAllComments()
  const solicitanteComment = comments.find(c => c.type === 'solicitante')
  return solicitanteComment ? solicitanteComment.content : ''
}

const getReturnComments = () => {
  const comments = parseAllComments()
  const devolucionComments = comments.filter(c => c.type === 'devolución')
  return devolucionComments.length > 0 ? devolucionComments[devolucionComments.length - 1].content : ''
}

const getReturnCommentsAuthor = () => {
  const comments = parseAllComments()
  const devolucionComments = comments.filter(c => c.type === 'devolución')
  return devolucionComments.length > 0 ? devolucionComments[devolucionComments.length - 1].author : 'Encargado de Bodega'
}

const getStatusBadgeClass = (status) => {
  const color = supplyRequestService.getStatusColor(status)
  const classes = {
    'yellow': 'bg-yellow-100 text-yellow-800',
    'green': 'bg-green-100 text-green-800',
    'red': 'bg-red-100 text-red-800',
    'blue': 'bg-blue-100 text-blue-800',
    'gray': 'bg-gray-100 text-gray-800'
  }
  return classes[color] || classes.gray
}

const getItemStatusBadgeClass = (status) => {
  const classes = {
    'pendiente': 'bg-yellow-100 text-yellow-800',
    'aceptado': 'bg-green-100 text-green-800',
    'rechazado': 'bg-red-100 text-red-800',
    'devuelto': 'bg-orange-100 text-orange-800'
  }
  return classes[status] || 'bg-gray-100 text-gray-800'
}

const getItemStatusLabel = (status) => {
  const labels = {
    'pendiente': 'Pendiente',
    'aceptado': 'Aceptado',
    'rechazado': 'Rechazado',
    'devuelto': 'Devuelto'
  }
  return labels[status] || status
}

const getUrgencyBadgeClass = (urgency) => {
  // Implementación directa para urgencia de items
  const urgencyColors = {
    'low': 'bg-gray-100 text-gray-800',
    'normal': 'bg-blue-100 text-blue-800',
    'high': 'bg-orange-100 text-orange-800',
    'critical': 'bg-red-100 text-red-800'
  }
  return urgencyColors[urgency] || urgencyColors.normal
}

const getAssignmentStatusBadgeClass = (status) => {
  const classes = {
    'assigned': 'bg-blue-100 text-blue-800',
    'delivered': 'bg-green-100 text-green-800',
    'consumed': 'bg-gray-100 text-gray-800',
    'returned': 'bg-yellow-100 text-yellow-800'
  }
  return classes[status] || 'bg-gray-100 text-gray-800'
}

const getAssignmentStatusLabel = (status) => {
  const labels = {
    'assigned': 'Asignado',
    'delivered': 'Entregado',
    'consumed': 'Consumido',
    'returned': 'Devuelto'
  }
  return labels[status] || status
}

// Métodos para el carrito
const onCartLoaded = (cart) => {
  console.log('Carrito cargado:', cart)
  currentCart.value = cart
  canTransferCartToPavilion.value = cart?.canTransferToPavilion || false
  
  // Si el carrito está activo, colapsar automáticamente la sección de insumos
  if (cart && cart.status === 'active') {
    hasActiveCart.value = true
    itemsCollapsed.value = true // Colapsar automáticamente cuando hay carrito activo
  } else {
    hasActiveCart.value = false
  }
}

const toggleItemsCollapsed = () => {
  itemsCollapsed.value = !itemsCollapsed.value
}

const onCartClosed = (cart) => {
  hasActiveCart.value = false
  itemsCollapsed.value = false // Expandir la sección cuando se cierra el carrito
  showMessage('success', 'Carrito cerrado', 'El carrito ha sido cerrado exitosamente')
  // Recargar la solicitud para reflejar el nuevo estado (por ejemplo, "Completado")
  loadSupplyRequest()
}

const onItemRemoved = (itemId) => {
  console.log('Item removido del carrito:', itemId)
  showMessage('info', 'Item removido', 'El item ha sido removido del carrito')
}

const handleTransferCartToPavilion = async () => {
  if (!currentCart.value) return
  
  openConfirmModal(
    'Listo para retiro',
    `¿Está seguro de marcar este carrito como "Listo para retiro"? Esto indicará al pabellón que puede proceder a retirar los insumos.`,
    async () => {
      processing.value = true
      try {
        const response = await cartService.transferCartToPavilion(currentCart.value.id)
        if (response.success) {
          showMessage('success', 'Listo para retiro', response.message || 'Los insumos han sido transferidos al pabellón. El pabellón debe confirmar la recepción.')
          // Recargar la solicitud para actualizar el estado del carrito
          await loadSupplyRequest()
        }
      } catch (err) {
        console.error('Error al transferir carrito:', err)
        showMessage('error', 'Error', err.response?.data?.message || 'Error al transferir el carrito al pabellón')
      } finally {
        processing.value = false
      }
    }
  )
}

// Funciones para manejar modales
const openConfirmModal = (title, message, callback) => {
  confirmModalTitle.value = title
  confirmModalMessage.value = message
  confirmModalCallback.value = callback
  showConfirmModal.value = true
}

const closeConfirmModal = () => {
  showConfirmModal.value = false
  confirmModalTitle.value = ''
  confirmModalMessage.value = ''
  confirmModalCallback.value = null
}

const executeConfirmCallback = () => {
  if (confirmModalCallback.value) {
    confirmModalCallback.value()
  }
  closeConfirmModal()
}

const openMessageModal = (type, title, message) => {
  messageModalType.value = type
  messageModalTitle.value = title
  messageModalMessage.value = message
  showMessageModal.value = true
}

const closeMessageModal = () => {
  showMessageModal.value = false
  messageModalType.value = 'info'
  messageModalTitle.value = ''
  messageModalMessage.value = ''
}

const showMessage = (type, title, message) => {
  openMessageModal(type, title, message)
}

const openInputModal = (title, label, placeholder, callback, confirmText = 'Confirmar') => {
  inputModalTitle.value = title
  inputModalLabel.value = label
  inputModalPlaceholder.value = placeholder
  inputModalValue.value = ''
  inputModalError.value = ''
  inputModalConfirmText.value = confirmText
  inputModalCallback.value = callback
  showInputModal.value = true
}

const closeInputModal = () => {
  showInputModal.value = false
  inputModalTitle.value = ''
  inputModalLabel.value = ''
  inputModalPlaceholder.value = ''
  inputModalValue.value = ''
  inputModalError.value = ''
  inputModalConfirmText.value = 'Confirmar'
  inputModalCallback.value = null
}

const confirmInputModal = () => {
  if (!inputModalValue.value.trim() && inputModalConfirmText.value === 'Rechazar') {
    inputModalError.value = 'Debe ingresar un motivo'
    return
  }
  
  if (inputModalCallback.value) {
    inputModalError.value = '' // Limpiar error antes de ejecutar
    inputModalCallback.value(inputModalValue.value)
    closeInputModal() // Cerrar después de ejecutar el callback
  }
}

// Función para normalizar texto (quitar tildes y convertir a minúsculas)
const normalizeText = (text) => {
  if (!text) return ''
  return text
    .toLowerCase()
    .normalize('NFD')
    .replace(/[\u0300-\u036f]/g, '') // Quitar diacríticos
    .replace(/ñ/g, 'n')
    .replace(/ü/g, 'u')
}

// Buscar usuarios según búsqueda (con debounce)
let searchTimeout = null
const onUserSearch = async () => {
  const query = userSearchQuery.value.trim()
  
  // Limpiar timeout anterior
  if (searchTimeout) {
    clearTimeout(searchTimeout)
  }
  
  if (!query || query.length < 2) {
    filteredUsers.value = []
    showUserSuggestions.value = false
    return
  }
  
  // Debounce: esperar 300ms antes de buscar
  searchTimeout = setTimeout(async () => {
    try {
      const result = await userService.searchUsers(query)
      if (result.success && result.data) {
        filteredUsers.value = result.data.slice(0, 10) // Limitar a 10 resultados
        showUserSuggestions.value = true
      } else {
        filteredUsers.value = []
        showUserSuggestions.value = false
      }
    } catch (err) {
      console.error('Error buscando usuarios:', err)
      filteredUsers.value = []
      showUserSuggestions.value = false
    }
  }, 300)
}

// Seleccionar usuario de la lista
const selectUser = (user) => {
  // Asegurar que allow_anyone_to_pickup esté en false cuando se selecciona un usuario
  pickupConfig.value.allow_anyone_to_pickup = false
  pickupConfig.value.authorized_pickup_rut = user.rut
  pickupConfig.value.authorized_pickup_name = user.name
  userSearchQuery.value = user.name
  showUserSuggestions.value = false
  // Guardar automáticamente al seleccionar
  savePickupConfig()
}

// Limpiar usuario seleccionado
const clearSelectedUser = () => {
  pickupConfig.value.authorized_pickup_rut = null
  pickupConfig.value.authorized_pickup_name = null
  userSearchQuery.value = ''
  showUserSuggestions.value = false
  // Guardar automáticamente al limpiar
  savePickupConfig()
}

// Manejar blur del campo de búsqueda
const handleUserSearchBlur = () => {
  // Esperar un poco antes de ocultar para permitir clicks en las sugerencias
  setTimeout(() => {
    showUserSuggestions.value = false
  }, 200)
}

// Manejar cambio del checkbox (solo mostrar/ocultar campos, no guardar)
const onCheckboxChange = () => {
  // Si se marca "cualquiera", limpiar la selección y guardar
  if (pickupConfig.value.allow_anyone_to_pickup) {
    pickupConfig.value.authorized_pickup_rut = null
    pickupConfig.value.authorized_pickup_name = null
    userSearchQuery.value = ''
    savePickupConfig()
  }
  // Si se desmarca, solo mostrar los campos (no guardar hasta que se seleccione un usuario)
}

// Guardar configuración de retiro
const savePickupConfig = async () => {
  if (savingPickupConfig.value) return
  
  // Asegurar que allow_anyone_to_pickup siempre sea un booleano explícito
  // Si es undefined o null, usar true por defecto
  const allowAnyone = pickupConfig.value.allow_anyone_to_pickup === false ? false : true
  
  // Validar que si no permite a cualquiera, debe tener un RUT
  if (!allowAnyone && !pickupConfig.value.authorized_pickup_rut) {
    // No mostrar error si el usuario está escribiendo
    return
  }

  savingPickupConfig.value = true
  try {
    // Construir objeto de configuración con valores explícitos
    // IMPORTANTE: Siempre incluir allow_anyone_to_pickup como booleano explícito
    const configData = {
      allow_anyone_to_pickup: allowAnyone, // Siempre un booleano explícito (true o false)
      authorized_pickup_rut: allowAnyone ? null : (pickupConfig.value.authorized_pickup_rut || null),
      authorized_pickup_name: allowAnyone ? null : (pickupConfig.value.authorized_pickup_name || null)
    }
    

    const result = await supplyRequestService.configurePickupAuthorization(requestId.value, configData)
    
    if (result.success) {
      // Actualizar la solicitud local
      if (request.value) {
        request.value.allow_anyone_to_pickup = configData.allow_anyone_to_pickup
        request.value.authorized_pickup_rut = configData.authorized_pickup_rut
        request.value.authorized_pickup_name = configData.authorized_pickup_name
      }
      showMessage('success', 'Configuración guardada', 'La configuración de retiro se ha actualizado correctamente')
    } else {
      showMessage('error', 'Error', result.error || 'No se pudo guardar la configuración')
      // Revertir cambios
      if (request.value) {
        pickupConfig.value = {
          allow_anyone_to_pickup: request.value.allow_anyone_to_pickup !== undefined ? request.value.allow_anyone_to_pickup : true,
          authorized_pickup_rut: request.value.authorized_pickup_rut || null,
          authorized_pickup_name: request.value.authorized_pickup_name || null
        }
      }
    }
  } catch (err) {
    console.error('Error guardando configuración de retiro:', err)
    showMessage('error', 'Error', err.response?.data?.error || 'Error al guardar la configuración')
    // Revertir cambios
    if (request.value) {
      pickupConfig.value = {
        allow_anyone_to_pickup: request.value.allow_anyone_to_pickup !== undefined ? request.value.allow_anyone_to_pickup : true,
        authorized_pickup_rut: request.value.authorized_pickup_rut || null,
        authorized_pickup_name: request.value.authorized_pickup_name || null
      }
    }
  } finally {
    savingPickupConfig.value = false
  }
}

// Lifecycle
onMounted(async () => {
  await Promise.all([
    loadSupplyRequest(),
    loadPavilions()
  ])
})
</script>

<style scoped>
/* Mejoras para dispositivos móviles */
@media (max-width: 640px) {
  /* Inputs más grandes para fácil toque en móviles */
  input[type="text"],
  textarea {
    min-height: 42px;
    font-size: 16px; /* Previene zoom en iOS */
  }

  /* Mejorar la altura mínima de los botones para toque fácil */
  button {
    min-height: 44px; /* Tamaño recomendado por Apple para touch targets */
  }

  /* Asegurar que todos los elementos interactivos sean táctiles */
  button,
  a {
    touch-action: manipulation;
    -webkit-tap-highlight-color: transparent;
  }

  /* Mejorar la legibilidad de los badges */
  .inline-flex.rounded-full {
    white-space: nowrap;
  }

  /* Ajustar modales para móviles */
  .fixed.inset-0 > div {
    max-height: 95vh;
    overflow-y: auto;
  }
}

/* Animaciones suaves */
@media (prefers-reduced-motion: no-preference) {
  button {
    transition: all 0.15s ease-in-out;
  }

  .hover\:bg-gray-50:hover {
    transition: background-color 0.15s ease-in-out;
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

/* Asegurar que las tarjetas sean scrolleables en móviles */
.overflow-y-auto {
  -webkit-overflow-scrolling: touch;
}

/* Mejoras para códigos QR largos */
code {
  word-break: break-all;
  overflow-wrap: break-word;
}
</style>