<template>
  <div class="bg-white rounded-lg shadow-md p-6">
    <!-- Header del carrito -->
    <div class="flex justify-between items-center mb-6">
      <div class="flex items-center gap-3 flex-1">
        <button
          @click="collapsed = !collapsed"
          class="text-gray-500 hover:text-gray-700 transition-colors"
          :title="collapsed ? 'Expandir' : 'Colapsar'"
        >
          <svg
            class="h-5 w-5 transition-transform duration-200"
            :class="{ 'transform rotate-180': !collapsed }"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
          </svg>
        </button>
      <div>
        <h2 class="text-2xl font-bold text-gray-800">
          🛒 Carrito de Insumos
        </h2>
        <p class="text-sm text-gray-500 mt-1">
          {{ cart?.cart_number || 'Cargando...' }}
        </p>
        </div>
      </div>
      <div class="flex items-center gap-3">
        <span 
          class="px-4 py-2 rounded-full text-sm font-semibold"
          :class="getStatusClass(cart?.status)"
        >
          {{ getStatusLabel(cart?.status) }}
        </span>
        <button
          v-if="cart?.status === 'active' && canTransferToPavilion"
          @click="handleTransferToPavilion"
          class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors flex items-center gap-2"
          :disabled="loading"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4" />
          </svg>
          Transferir al Pabellón
        </button>
        <button
          v-if="cart?.status === 'active' && canClose"
          @click="handleCloseCart"
          class="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors"
          :disabled="loading"
        >
          Cerrar Carrito
        </button>
      </div>
    </div>

    <!-- Contenido colapsable -->
    <div v-show="!collapsed">
    <!-- Información del carrito -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6 p-4 bg-gray-50 rounded-lg">
      <div>
        <p class="text-sm text-gray-600">Solicitud</p>
        <p class="font-semibold text-gray-800">
          {{ cart?.supply_request?.request_number || '-' }}
        </p>
      </div>
      <div>
        <p class="text-sm text-gray-600">Creado por</p>
        <p class="font-semibold text-gray-800">
          {{ cart?.created_by_name || '-' }}
        </p>
      </div>
      <div>
        <p class="text-sm text-gray-600">Fecha de creación</p>
        <p class="font-semibold text-gray-800">
          {{ formatDate(cart?.created_at) }}
        </p>
      </div>
    </div>

    <!-- Loading state -->
    <div v-if="loading" class="flex justify-center items-center py-12">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
    </div>

    <!-- Items del carrito -->
    <div v-else-if="cart && activeItems.length > 0" class="space-y-4">
      <div class="flex justify-between items-center mb-4">
        <h3 class="text-lg font-semibold text-gray-800">
          Items en el Carrito ({{ activeItems.length }})
        </h3>
        <div class="flex gap-2">
          <!-- Botones de operación múltiple -->
          <div v-if="cart?.status === 'active' && canManageItems && selectedItems.length > 0" class="flex gap-2">
            <button
              @click="showBatchOperationModal = true"
              class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors flex items-center gap-2"
              :disabled="loading"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01" />
              </svg>
              Operar Seleccionados ({{ selectedItems.length }})
            </button>
            <button
              @click="clearSelection"
              class="px-3 py-2 bg-gray-400 text-white rounded-lg hover:bg-gray-500 transition-colors"
              title="Limpiar selección"
            >
              ✕
            </button>
          </div>
        <button
          v-if="showAddButton"
          @click="$emit('add-item')"
          class="px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors"
        >
          + Agregar Item
        </button>
        </div>
      </div>

      <!-- Seleccionar todos -->
      <div v-if="cart?.status === 'active' && canManageItems && selectableItems.length > 0" class="mb-3">
        <label class="flex items-center gap-2 text-sm text-gray-700 cursor-pointer">
          <input
            type="checkbox"
            :checked="allSelected"
            @change="toggleSelectAll"
            class="w-4 h-4 text-blue-600 rounded focus:ring-blue-500"
          />
          <span class="font-medium">Seleccionar todos los items disponibles ({{ selectableItems.length }})</span>
        </label>
      </div>

      <div class="grid grid-cols-1 gap-4">
        <div
          v-for="item in activeItems"
          :key="item.id"
          class="border rounded-lg p-4 hover:shadow-md transition-shadow"
          :class="isItemSelected(item.id) ? 'border-blue-500 bg-blue-50' : 'border-gray-200'"
        >
          <div class="flex justify-between items-start">
            <div class="flex-1">
              <!-- Checkbox para selección múltiple -->
              <div v-if="cart?.status === 'active' && canManageItems && !isItemProcessed(item)" class="mb-3">
                <label class="flex items-center gap-2 cursor-pointer" :class="!isItemReceivedInPavilion(item) ? 'opacity-50 cursor-not-allowed' : ''">
                  <input
                    type="checkbox"
                    :checked="isItemSelected(item.id)"
                    @change="toggleItemSelection(item.id)"
                    :disabled="!isItemReceivedInPavilion(item)"
                    class="w-4 h-4 text-blue-600 rounded focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
                  />
                  <span class="text-sm font-medium text-gray-700">
                    Seleccionar para operación múltiple
                    <span v-if="!isItemReceivedInPavilion(item)" class="text-red-600 text-xs ml-1">(No recibido)</span>
                  </span>
                </label>
              </div>
              
              <!-- Información del insumo -->
              <div class="flex items-center gap-3 mb-3">
                <div class="bg-blue-100 p-2 rounded-lg">
                  <svg class="w-6 h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
                  </svg>
                </div>
                <div>
                  <h4 class="font-semibold text-gray-800">
                    {{ item.supply_request_qr_assignment?.supply_request_item?.supply_name || 'Insumo' }}
                  </h4>
                  <p class="text-sm text-gray-600">
                    Código: {{ item.supply_request_qr_assignment?.supply_request_item?.supply_code || '-' }}
                  </p>
                </div>
              </div>

              <!-- Código QR -->
              <div class="bg-gray-50 p-3 rounded-lg mb-3">
                <div class="flex items-center gap-2">
                  <svg class="w-5 h-5 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h.01M5 8h2a1 1 0 001-1V5a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1zm12 0h2a1 1 0 001-1V5a1 1 0 00-1-1h-2a1 1 0 00-1 1v2a1 1 0 001 1zM5 20h2a1 1 0 001-1v-2a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1z" />
                  </svg>
                  <span class="font-mono text-sm text-gray-800">
                    {{ item.supply_request_qr_assignment?.qr_code || '-' }}
                  </span>
                  <button
                    v-if="item.supply_request_qr_assignment?.qr_code"
                    @click="copyToClipboard(item.supply_request_qr_assignment.qr_code)"
                    class="ml-auto p-1 text-gray-500 hover:text-blue-600 transition-colors"
                    title="Copiar QR"
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
                    </svg>
                  </button>
                </div>
              </div>

              <!-- Estado y fechas -->
              <div class="grid grid-cols-2 gap-3 text-sm">
                <div>
                  <span class="text-gray-600">Agregado:</span>
                  <span class="text-gray-800 ml-1">{{ formatDate(item.added_at) }}</span>
                </div>
                <div>
                  <span class="text-gray-600">Por:</span>
                  <span class="text-gray-800 ml-1">{{ item.added_by_name }}</span>
                </div>
              </div>

              <!-- Estado del insumo -->
              <div v-if="getItemStatus(item) || getItemLocationStatus(item)" class="mt-3 flex gap-2 flex-wrap">
                <span 
                  v-if="getItemStatus(item)"
                  class="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium"
                  :class="getItemStatusClass(item)"
                >
                  {{ getItemStatus(item) }}
                </span>
                <span 
                  v-if="getItemLocationStatus(item)"
                  class="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium"
                  :class="getItemLocationStatusClass(item)"
                >
                  {{ getItemLocationStatus(item) }}
                </span>
              </div>

              <!-- Notas si existen -->
              <div v-if="item.notes" class="mt-3 p-2 bg-yellow-50 border-l-4 border-yellow-400 text-sm text-gray-700">
                <strong>Nota:</strong> {{ item.notes }}
              </div>

              <!-- Mensaje si el insumo no está recibido en el pabellón -->
              <div v-if="cart?.status === 'active' && canManageItems && !isItemProcessed(item) && !isItemReceivedInPavilion(item)" class="mt-4 p-3 bg-yellow-50 border border-yellow-200 rounded-lg">
                <p class="text-sm text-yellow-800">
                  <svg class="w-4 h-4 inline mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                  </svg>
                  Este insumo debe ser transferido al pabellón y recibido antes de poder ser utilizado o devuelto.
                </p>
              </div>

              <!-- Botones de acción para gestión del item -->
              <div v-if="cart?.status === 'active' && canManageItems && !isItemProcessed(item) && isItemReceivedInPavilion(item)" class="mt-4 flex gap-2">
                <button
                  @click="handleMarkAsUsed(item.id)"
                  class="flex-1 px-3 py-2 bg-green-600 text-white text-sm font-medium rounded-lg hover:bg-green-700 transition-colors flex items-center justify-center gap-2"
                  :disabled="loading"
                >
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                  </svg>
                  Utilizado
                </button>
                <button
                  @click="handleMarkForReturn(item.id)"
                  class="flex-1 px-3 py-2 bg-orange-600 text-white text-sm font-medium rounded-lg hover:bg-orange-700 transition-colors flex items-center justify-center gap-2"
                  :disabled="loading"
                >
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 10h10a8 8 0 018 8v2M3 10l6 6m-6-6l6-6" />
                  </svg>
                  Devolver
                </button>
              </div>
            </div>

            <!-- Acciones de admin -->
            <div v-if="cart?.status === 'active' && canRemoveItems" class="ml-4">
              <button
                @click="handleRemoveItem(item.id)"
                class="p-2 text-red-600 hover:bg-red-50 rounded-lg transition-colors"
                title="Remover del carrito"
                :disabled="loading"
              >
                <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                </svg>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Carrito no existe todavía -->
    <div v-else-if="!cart && !error" class="text-center py-12 bg-blue-50 rounded-lg border-2 border-dashed border-blue-200">
      <svg class="w-16 h-16 mx-auto text-blue-400 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z" />
      </svg>
      <p class="text-gray-600 text-lg font-semibold">Carrito no disponible</p>
      <p class="text-gray-500 text-sm mt-2">El carrito se creará automáticamente cuando se asigne el primer código QR a esta solicitud</p>
    </div>

    <!-- Estado vacío -->
    <div v-else-if="cart && activeItems.length === 0" class="text-center py-12">
      <svg class="w-16 h-16 mx-auto text-gray-400 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z" />
      </svg>
      <p class="text-gray-600 text-lg">El carrito está vacío</p>
      <p class="text-gray-500 text-sm mt-2">No hay items agregados a este carrito</p>
    </div>

    <!-- Error state -->
    <div v-else-if="error" class="text-center py-12">
      <svg class="w-16 h-16 mx-auto text-red-400 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>
      <p class="text-gray-600 text-lg">{{ error }}</p>
      </div>
    </div>

    <!-- Modales personalizados -->
    <!-- Modal de confirmación -->
    <div v-if="showConfirmModal" class="fixed inset-0 z-[60] overflow-y-auto" aria-labelledby="modal-title" role="dialog" aria-modal="true">
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
    <div v-if="showMessageModal" class="fixed inset-0 z-[60] overflow-y-auto" aria-labelledby="modal-title" role="dialog" aria-modal="true">
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

    <!-- Modal de input (para motivo de devolución) -->
    <div v-if="showInputModal" class="fixed inset-0 z-[60] overflow-y-auto" aria-labelledby="modal-title" role="dialog" aria-modal="true">
      <div class="flex items-center justify-center min-h-screen px-4 pt-4 pb-20 text-center sm:block sm:p-0">
        <div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" aria-hidden="true" @click="closeInputModal"></div>
        <span class="hidden sm:inline-block sm:align-middle sm:h-screen" aria-hidden="true">&#8203;</span>
        <div class="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full">
          <div class="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
            <div class="sm:flex sm:items-start">
              <div class="mx-auto flex-shrink-0 flex items-center justify-center h-12 w-12 rounded-full bg-orange-100 sm:mx-0 sm:h-10 sm:w-10">
                <svg class="h-6 w-6 text-orange-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 10h10a8 8 0 018 8v2M3 10l6 6m-6-6l6-6" />
                </svg>
              </div>
              <div class="mt-3 text-center sm:mt-0 sm:ml-4 sm:text-left flex-1 w-full">
                <h3 class="text-lg leading-6 font-medium text-gray-900" id="modal-title">
                  {{ inputModalTitle }}
                </h3>
                <div class="mt-4">
                  <label class="block text-sm font-medium text-gray-700 mb-2">{{ inputModalLabel }}</label>
                  <textarea
                    v-model="inputModalValue"
                    rows="3"
                    class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
                    :placeholder="inputModalPlaceholder"
                  ></textarea>
                </div>
              </div>
            </div>
          </div>
          <div class="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse gap-2">
            <button
              type="button"
              @click="confirmInputModal"
              class="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-blue-600 text-base font-medium text-white hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 sm:ml-3 sm:w-auto sm:text-sm"
            >
              Confirmar
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

    <!-- Estado vacío -->
    <div v-else-if="cart && activeItems.length === 0" class="text-center py-12">
      <svg class="w-16 h-16 mx-auto text-gray-400 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z" />
      </svg>
      <p class="text-gray-600 text-lg">El carrito está vacío</p>
      <p class="text-gray-500 text-sm mt-2">No hay items agregados a este carrito</p>
    </div>

    <!-- Error state -->
    <div v-else-if="error" class="text-center py-12">
      <svg class="w-16 h-16 mx-auto text-red-400 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>
      <p class="text-gray-600 text-lg">{{ error }}</p>
    </div>

    <!-- Modal de operación múltiple -->
    <div v-if="showBatchOperationModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg shadow-xl p-6 max-w-2xl w-full mx-4 max-h-[90vh] overflow-y-auto">
        <h3 class="text-xl font-bold text-gray-800 mb-4">
          Operación Múltiple - {{ selectedItems.length }} items seleccionados
        </h3>
        
        <p class="text-sm text-gray-600 mb-4">
          Marque cada item como <strong>usado</strong> o <strong>devolución</strong>. Puede procesar todos en un solo paso.
        </p>

        <!-- Advertencia si hay items no recibidos en el pabellón -->
        <div v-if="hasNonReceivedItems" class="mb-4 p-3 bg-yellow-50 border border-yellow-200 rounded-lg">
          <p class="text-sm text-yellow-800">
            <svg class="w-4 h-4 inline mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
            </svg>
            Algunos items seleccionados no están recibidos en el pabellón. Solo los items recibidos pueden ser utilizados o devueltos.
          </p>
        </div>

        <div class="space-y-3 mb-6">
          <div
            v-for="itemId in selectedItems"
            :key="itemId"
            class="border rounded-lg p-4 transition-all duration-200"
            :class="getItemOperation(itemId) === 'use' 
              ? 'border-green-400 bg-green-50' 
              : getItemOperation(itemId) === 'return' 
              ? 'border-orange-400 bg-orange-50' 
              : 'border-gray-200 bg-white'"
          >
            <div class="flex items-center justify-between">
              <div class="flex-1">
                <div class="flex items-center gap-2 mb-1">
                  <p class="font-semibold text-gray-800">
                    {{ getItemName(itemId) }}
                  </p>
                  <span 
                    v-if="getItemOperation(itemId)"
                    class="inline-flex items-center px-2 py-1 text-xs font-semibold rounded-full"
                    :class="getItemOperation(itemId) === 'use' 
                      ? 'bg-green-100 text-green-800' 
                      : 'bg-orange-100 text-orange-800'"
                  >
                    {{ getItemOperation(itemId) === 'use' ? '✓ Usado' : '↩ Devolver' }}
                  </span>
                  <span 
                    v-if="!isItemReceivedInPavilionById(itemId)"
                    class="inline-flex items-center px-2 py-1 text-xs font-semibold rounded-full bg-red-100 text-red-800"
                  >
                    ⚠ No recibido
                  </span>
                </div>
                <p class="text-sm text-gray-600">
                  QR: {{ getItemQR(itemId) }}
                </p>
                <p v-if="!isItemReceivedInPavilionById(itemId)" class="text-xs text-red-600 mt-1">
                  Este insumo debe ser transferido y recibido en el pabellón antes de poder ser operado.
                </p>
              </div>
              <div class="flex gap-2">
                <button
                  @click="setItemOperation(itemId, 'use')"
                  type="button"
                  :disabled="!isItemReceivedInPavilionById(itemId)"
                  class="px-4 py-2 rounded-lg text-sm font-medium transition-all duration-200 transform active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed"
                  :class="getItemOperation(itemId) === 'use' 
                    ? 'bg-green-600 text-white shadow-md ring-2 ring-green-300 ring-offset-2' 
                    : isItemReceivedInPavilionById(itemId)
                    ? 'bg-gray-200 text-gray-700 hover:bg-green-100 hover:text-green-700 hover:border-green-300 border-2 border-transparent'
                    : 'bg-gray-200 text-gray-400 border-2 border-transparent'"
                >
                  <span class="flex items-center gap-1">
                    <svg v-if="getItemOperation(itemId) === 'use'" class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                      <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                    </svg>
                    ✓ Usado
                  </span>
                </button>
                <button
                  @click="setItemOperation(itemId, 'return')"
                  type="button"
                  :disabled="!isItemReceivedInPavilionById(itemId)"
                  class="px-4 py-2 rounded-lg text-sm font-medium transition-all duration-200 transform active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed"
                  :class="getItemOperation(itemId) === 'return' 
                    ? 'bg-orange-600 text-white shadow-md ring-2 ring-orange-300 ring-offset-2' 
                    : isItemReceivedInPavilionById(itemId)
                    ? 'bg-gray-200 text-gray-700 hover:bg-orange-100 hover:text-orange-700 hover:border-orange-300 border-2 border-transparent'
                    : 'bg-gray-200 text-gray-400 border-2 border-transparent'"
                >
                  <span class="flex items-center gap-1">
                    <svg v-if="getItemOperation(itemId) === 'return'" class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                      <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                    </svg>
                    ↩ Devolver
                  </span>
                </button>
              </div>
            </div>
            <!-- Campo de motivo para devolución -->
            <div v-if="getItemOperation(itemId) === 'return'" class="mt-3 transition-all duration-300 ease-in-out">
              <label class="block text-sm font-medium text-gray-700 mb-1">Motivo de devolución (opcional):</label>
              <input
                v-model="itemOperations[itemId].reason"
                type="text"
                placeholder="Ingrese el motivo de la devolución..."
                class="w-full px-3 py-2 border-2 border-orange-300 rounded-lg text-sm focus:ring-2 focus:ring-orange-500 focus:border-orange-500 transition-all"
                @focus="$event.target.select()"
              />
            </div>
          </div>
        </div>

        <div class="flex justify-end gap-3">
          <button
            @click="showBatchOperationModal = false"
            class="px-4 py-2 bg-gray-300 text-gray-700 rounded-lg hover:bg-gray-400 transition-colors"
          >
            Cancelar
          </button>
          <button
            @click="handleBatchOperation"
            class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors flex items-center gap-2"
            :disabled="loading || !hasValidOperations"
          >
            <svg v-if="loading" class="animate-spin h-4 w-4" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            {{ loading ? 'Procesando...' : 'Procesar Operación Múltiple' }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import cartService from '@/services/cartService'

const props = defineProps({
  requestId: {
    type: Number,
    default: null
  },
  cartId: {
    type: Number,
    default: null
  },
  qrCode: {
    type: String,
    default: null
  },
  canClose: {
    type: Boolean,
    default: false
  },
  canRemoveItems: {
    type: Boolean,
    default: false
  },
  canManageItems: {
    type: Boolean,
    default: true  // Por defecto permite marcar como usado/devolver
  },
  showAddButton: {
    type: Boolean,
    default: false
  },
  defaultCollapsed: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['cart-loaded', 'cart-closed', 'item-removed', 'item-used', 'item-returned', 'add-item', 'error'])

const cart = ref(null)
const loading = ref(false)
const error = ref(null)
const selectedItems = ref([])
const showBatchOperationModal = ref(false)
const itemOperations = ref({}) // { itemId: { action: 'use'|'return', reason: '' } }
const collapsed = ref(props.defaultCollapsed)

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
const inputModalCallback = ref(null)

const activeItems = computed(() => {
  if (!cart.value?.items) return []
  return cart.value.items.filter(item => item.is_active)
})

const selectableItems = computed(() => {
  // Solo items que no estén procesados Y que estén recibidos en el pabellón
  return activeItems.value.filter(item => 
    !isItemProcessed(item) && isItemReceivedInPavilion(item)
  )
})

const allSelected = computed(() => {
  return selectableItems.value.length > 0 && 
         selectableItems.value.every(item => selectedItems.value.includes(item.id))
})

const hasValidOperations = computed(() => {
  return selectedItems.value.length > 0 && 
         selectedItems.value.every(itemId => {
           // El item debe estar recibido en el pabellón
           if (!isItemReceivedInPavilionById(itemId)) {
             return false
           }
           const op = itemOperations.value[itemId]
           return op && (op.action === 'use' || op.action === 'return')
         })
})

// Verificar si hay items no recibidos en el pabellón
const hasNonReceivedItems = computed(() => {
  return selectedItems.value.some(itemId => !isItemReceivedInPavilionById(itemId))
})

// Verificar si un item está recibido en el pabellón por ID
const isItemReceivedInPavilionById = (itemId) => {
  const item = activeItems.value.find(i => i.id === itemId)
  if (!item) return false
  return isItemReceivedInPavilion(item)
}

const loadCart = async () => {
  loading.value = true
  error.value = null
  
  try {
    let response
    
    if (props.requestId) {
      response = await cartService.getCartByRequestId(props.requestId)
    } else if (props.cartId) {
      response = await cartService.getCartById(props.cartId)
    } else if (props.qrCode) {
      response = await cartService.getCartByQRCode(props.qrCode)
    } else {
      throw new Error('Se requiere requestId, cartId o qrCode')
    }
    
    if (response.success) {
      cart.value = response.data
      emit('cart-loaded', cart.value)
    } else {
      throw new Error(response.message || 'Error al cargar carrito')
    }
  } catch (err) {
    console.error('Error al cargar carrito:', err)
    
    // Si el error es 404 (carrito no encontrado), no es un error real
    // El carrito simplemente no existe todavía
    if (err.response?.status === 404) {
      cart.value = null
      error.value = null // No mostrar error, es un estado válido
    } else {
      error.value = err.response?.data?.message || err.message || 'Error al cargar el carrito'
      emit('error', error.value)
    }
  } finally {
    loading.value = false
  }
}

const handleCloseCart = async () => {
  openConfirmModal(
    'Cerrar Carrito',
    '¿Está seguro de cerrar este carrito? Esta acción no se puede deshacer.',
    async () => {
  loading.value = true
  try {
    const response = await cartService.closeCart(cart.value.id)
    if (response.success) {
      await loadCart() // Recargar para ver el estado actualizado
      emit('cart-closed', cart.value)
          showMessage('success', 'Carrito Cerrado', 'El carrito ha sido cerrado exitosamente')
    }
  } catch (err) {
    console.error('Error al cerrar carrito:', err)
        showMessage('error', 'Error', err.response?.data?.message || 'Error al cerrar el carrito')
  } finally {
    loading.value = false
  }
}
  )
  }
  
const handleRemoveItem = async (itemId) => {
  openConfirmModal(
    'Remover Item',
    '¿Está seguro de remover este item del carrito?',
    async () => {
  loading.value = true
  try {
    const response = await cartService.removeItemFromCart(cart.value.id, itemId)
    if (response.success) {
      await loadCart() // Recargar para ver los cambios
      emit('item-removed', itemId)
          showMessage('success', 'Item Removido', 'El item ha sido removido del carrito exitosamente')
    }
  } catch (err) {
    console.error('Error al remover item:', err)
        showMessage('error', 'Error', err.response?.data?.message || 'Error al remover el item')
  } finally {
    loading.value = false
  }
}
  )
  }
  
const handleMarkAsUsed = async (itemId) => {
  openConfirmModal(
    'Marcar como Utilizado',
    '¿Confirma que este insumo fue utilizado?',
    async () => {
  loading.value = true
  try {
    const response = await cartService.markItemAsUsed(cart.value.id, itemId)
    if (response.success) {
          const previousStatus = cart.value?.status
      await loadCart() // Recargar para ver los cambios
          
          // Verificar si el carrito se cerró automáticamente
          if (previousStatus === 'active' && cart.value?.status === 'closed') {
            emit('cart-closed', cart.value)
            showMessage('info', 'Carrito Cerrado Automáticamente', 'Todos los items han sido procesados. El carrito se cerró automáticamente.')
          } else {
      emit('item-used', itemId)
            showMessage('success', 'Éxito', 'Insumo marcado como utilizado exitosamente')
          }
    }
  } catch (err) {
    console.error('Error al marcar item como utilizado:', err)
        showMessage('error', 'Error', err.response?.data?.message || 'Error al marcar el item como utilizado')
  } finally {
    loading.value = false
  }
    }
  )
}

const handleMarkForReturn = async (itemId) => {
  openInputModal(
    'Motivo de Devolución',
    'Motivo de la devolución (opcional):',
    'Ingrese el motivo de la devolución...',
    async (reason) => {
  loading.value = true
  try {
        const response = await cartService.markItemForReturn(cart.value.id, itemId, reason || '')
    if (response.success) {
          const previousStatus = cart.value?.status
      await loadCart() // Recargar para ver los cambios
          
          // Verificar si el carrito se cerró automáticamente
          if (previousStatus === 'active' && cart.value?.status === 'closed') {
            emit('cart-closed', cart.value)
            showMessage('info', 'Carrito Cerrado Automáticamente', 'Todos los items han sido procesados. El carrito se cerró automáticamente.')
          } else {
      emit('item-returned', itemId)
            showMessage('success', 'Éxito', 'Insumo marcado para devolución exitosamente')
          }
    }
  } catch (err) {
    console.error('Error al marcar item para devolución:', err)
        showMessage('error', 'Error', err.response?.data?.message || 'Error al marcar el item para devolución')
  } finally {
    loading.value = false
  }
    }
  )
}

const getItemStatus = (item) => {
  const status = item.supply_request_qr_assignment?.status
  if (status === 'consumed') return '✓ Utilizado'
  if (status === 'returned') return '↩ Para Devolver'
  return null
}

const getItemStatusClass = (item) => {
  const status = item.supply_request_qr_assignment?.status
  if (status === 'consumed') return 'bg-green-100 text-green-800'
  if (status === 'returned') return 'bg-orange-100 text-orange-800'
  return ''
}

const isItemProcessed = (item) => {
  const status = item.supply_request_qr_assignment?.status
  return status === 'consumed' || status === 'returned'
}

// Verificar si un item está recibido en el pabellón
const isItemReceivedInPavilion = (item) => {
  const medicalSupply = item.supply_request_qr_assignment?.medical_supply
  if (!medicalSupply) return false
  
  // El insumo debe estar en el pabellón y con estado "recepcionado"
  return medicalSupply.location_type === 'pavilion' && 
         medicalSupply.status === 'recepcionado'
}

// Verificar si todos los items están en bodega (para mostrar botón de transferir)
const canTransferToPavilion = computed(() => {
  if (!cart.value || cart.value.status !== 'active') return false
  if (activeItems.value.length === 0) return false
  
  // Verificar que todos los items activos estén en bodega
  return activeItems.value.every(item => {
    const medicalSupply = item.supply_request_qr_assignment?.medical_supply
    if (!medicalSupply) return false
    
    // Debe estar en bodega y no consumido
    return medicalSupply.location_type === 'store' && 
           medicalSupply.status !== 'consumido' &&
           !medicalSupply.in_transit
  })
})

// Obtener estado de ubicación del item
const getItemLocationStatus = (item) => {
  const medicalSupply = item.supply_request_qr_assignment?.medical_supply
  if (!medicalSupply) return null
  
  if (medicalSupply.location_type === 'store') {
    if (medicalSupply.in_transit) {
      return 'En tránsito al pabellón'
    }
    return 'En bodega'
  } else if (medicalSupply.location_type === 'pavilion') {
    if (medicalSupply.status === 'recepcionado') {
      return 'Recibido en pabellón'
    } else if (medicalSupply.status === 'en_camino_a_pabellon') {
      return 'En tránsito al pabellón'
    } else if (medicalSupply.status === 'en_camino_a_bodega') {
      return 'En tránsito a bodega'
    }
    return 'En pabellón'
  }
  
  return null
}

// Obtener clase CSS para el estado de ubicación
const getItemLocationStatusClass = (item) => {
  const medicalSupply = item.supply_request_qr_assignment?.medical_supply
  if (!medicalSupply) return ''
  
  if (medicalSupply.location_type === 'store') {
    return 'bg-blue-100 text-blue-800'
  } else if (medicalSupply.location_type === 'pavilion') {
    if (medicalSupply.status === 'recepcionado') {
      return 'bg-green-100 text-green-800'
    } else if (medicalSupply.in_transit) {
      return 'bg-yellow-100 text-yellow-800'
    }
    return 'bg-purple-100 text-purple-800'
  }
  
  return ''
}

// Manejar transferencia al pabellón
const handleTransferToPavilion = async () => {
  openConfirmModal(
    'Transferir al Pabellón',
    `¿Está seguro de transferir todos los items del carrito al pabellón? Esta acción marcará los insumos como en tránsito. El pabellón deberá confirmar la recepción.`,
    async () => {
      loading.value = true
      try {
        const response = await cartService.transferCartToPavilion(cart.value.id)
        if (response.success) {
          await loadCart() // Recargar para ver el estado actualizado
          showMessage('success', 'Transferencia Iniciada', response.message || 'Los insumos han sido transferidos al pabellón. El pabellón debe confirmar la recepción.')
        }
      } catch (err) {
        console.error('Error al transferir carrito:', err)
        showMessage('error', 'Error', err.response?.data?.message || 'Error al transferir el carrito al pabellón')
      } finally {
        loading.value = false
      }
    }
  )
}

// Funciones de selección múltiple
const isItemSelected = (itemId) => {
  return selectedItems.value.includes(itemId)
}

const toggleItemSelection = (itemId) => {
  // Verificar que el item esté recibido en el pabellón antes de permitir selección
  if (!isItemReceivedInPavilionById(itemId)) {
    showMessage('warning', 'Item no recibido', 'Este insumo debe ser transferido y recibido en el pabellón antes de poder ser seleccionado para operación.')
    return
  }
  
  const index = selectedItems.value.indexOf(itemId)
  if (index > -1) {
    selectedItems.value.splice(index, 1)
    delete itemOperations.value[itemId]
  } else {
    selectedItems.value.push(itemId)
    // NO inicializar con acción por defecto - el usuario debe elegir explícitamente
    if (!itemOperations.value[itemId]) {
      itemOperations.value[itemId] = { action: null, reason: '' }
    }
  }
}

const toggleSelectAll = () => {
  if (allSelected.value) {
    selectedItems.value = []
    itemOperations.value = {}
  } else {
    selectedItems.value = selectableItems.value.map(item => item.id)
    selectableItems.value.forEach(item => {
      // NO inicializar con acción por defecto - el usuario debe elegir explícitamente
      if (!itemOperations.value[item.id]) {
        itemOperations.value[item.id] = { action: null, reason: '' }
      }
    })
  }
}

const clearSelection = () => {
  selectedItems.value = []
  itemOperations.value = {}
}

const getItemName = (itemId) => {
  const item = activeItems.value.find(i => i.id === itemId)
  return item?.supply_request_qr_assignment?.supply_request_item?.supply_name || 'Insumo'
}

const getItemQR = (itemId) => {
  const item = activeItems.value.find(i => i.id === itemId)
  return item?.supply_request_qr_assignment?.qr_code || '-'
}

const getItemOperation = (itemId) => {
  const action = itemOperations.value[itemId]?.action
  // Retornar null si action es null explícitamente o si no existe
  return action || null
}

const setItemOperation = (itemId, action) => {
  // Asegurar que el item tiene una entrada en itemOperations
  if (!itemOperations.value[itemId]) {
    itemOperations.value[itemId] = { action: null, reason: '' }
  }
  
  // Si se hace clic en el mismo botón, deseleccionar (toggle)
  if (itemOperations.value[itemId].action === action) {
    itemOperations.value[itemId].action = null
    itemOperations.value[itemId].reason = ''
  } else {
    itemOperations.value[itemId].action = action
    if (action !== 'return') {
      itemOperations.value[itemId].reason = ''
    }
  }
  
  // Forzar actualización reactiva
  itemOperations.value = { ...itemOperations.value }
}

const handleBatchOperation = async () => {
  if (!hasValidOperations.value) {
    showMessage('error', 'Error', 'Por favor, seleccione una acción (usado o devolver) para cada item')
    return
  }

  openConfirmModal(
    'Procesar Operación Múltiple',
    `¿Confirma procesar ${selectedItems.value.length} items en una sola operación?`,
    async () => {
      loading.value = true
      try {
        // Preparar items para la operación múltiple
        const itemsToProcess = selectedItems.value.map(itemId => ({
          itemId: itemId,
          action: itemOperations.value[itemId].action,
          reason: itemOperations.value[itemId].reason || ''
        }))

        const response = await cartService.batchOperationItems(cart.value.id, itemsToProcess)
        
        if (response.success || response.data) {
          // Mostrar resultado
          const result = response.data || response
          const successCount = result.success_count || result.SuccessCount || 0
          const errorCount = result.error_count || result.ErrorCount || 0
          
          let message = `Operación completada: ${successCount} items procesados exitosamente`
          const messageType = errorCount > 0 ? 'info' : 'success'
          
          if (errorCount > 0) {
            message += `, ${errorCount} fallaron`
            if (result.errors && result.errors.length > 0) {
              console.error('Errores en operación múltiple:', result.errors)
            }
          }
          
          showMessage(messageType, 'Operación Múltiple', message)
          
          // Guardar información antes de limpiar
          const processedItems = [...selectedItems.value]
          const processedOperations = { ...itemOperations.value }
          
          // Limpiar selección
          clearSelection()
          showBatchOperationModal.value = false
          
          // Recargar carrito
          const previousStatus = cart.value?.status
          await loadCart()
          
          // Verificar si el carrito se cerró automáticamente
          if (previousStatus === 'active' && cart.value?.status === 'closed') {
            emit('cart-closed', cart.value)
            showMessage('info', 'Carrito Cerrado Automáticamente', 'Todos los items han sido procesados. El carrito se cerró automáticamente.')
          }
          
          // Emitir eventos
          processedItems.forEach(itemId => {
            const op = processedOperations[itemId]
            if (op?.action === 'use') {
              emit('item-used', itemId)
            } else if (op?.action === 'return') {
              emit('item-returned', itemId)
            }
          })
        } else {
          throw new Error(response.message || 'Error al procesar operación múltiple')
        }
      } catch (err) {
        console.error('Error en operación múltiple:', err)
        showMessage('error', 'Error', err.response?.data?.message || err.message || 'Error al procesar la operación múltiple')
      } finally {
        loading.value = false
      }
    }
  )
}

const copyToClipboard = async (text) => {
  try {
    await navigator.clipboard.writeText(text)
    showMessage('success', 'Copiado', 'Código QR copiado al portapapeles')
  } catch (err) {
    console.error('Error al copiar:', err)
    showMessage('error', 'Error', 'Error al copiar al portapapeles')
  }
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
  // Convertir 'warning' a 'info' si no está soportado
  const modalType = type === 'warning' ? 'info' : type
  openMessageModal(modalType, title, message)
}

const openInputModal = (title, label, placeholder, callback) => {
  inputModalTitle.value = title
  inputModalLabel.value = label
  inputModalPlaceholder.value = placeholder
  inputModalValue.value = ''
  inputModalCallback.value = callback
  showInputModal.value = true
}

const closeInputModal = () => {
  showInputModal.value = false
  inputModalTitle.value = ''
  inputModalLabel.value = ''
  inputModalPlaceholder.value = ''
  inputModalValue.value = ''
  inputModalCallback.value = null
}

const confirmInputModal = () => {
  if (inputModalCallback.value) {
    inputModalCallback.value(inputModalValue.value)
  }
  closeInputModal()
}

const formatDate = (dateString) => {
  return cartService.formatDate(dateString)
}

const getStatusLabel = (status) => {
  return cartService.getStatusLabel(status)
}

const getStatusClass = (status) => {
  return cartService.getStatusClass(status)
}

// Método público para recargar el carrito
const refresh = () => {
  loadCart()
}

// Exponer método para uso externo
defineExpose({
  refresh,
  loadCart
})

onMounted(() => {
  loadCart()
})
</script>

<style scoped>
/* Animaciones adicionales si se necesitan */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
