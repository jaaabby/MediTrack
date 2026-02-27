<template>
  <div v-if="show" class="fixed inset-0 z-50 flex items-center justify-center p-4" aria-labelledby="modal-title" role="dialog" aria-modal="true">
    <!-- Overlay -->
    <div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" aria-hidden="true" @click="closeModal"></div>

      <!-- Contenido del modal -->
      <div class="relative bg-white rounded-lg text-left shadow-xl transform transition-all w-full max-w-4xl flex flex-col max-h-[90vh]">
        <!-- Sticky header -->
        <div class="bg-white px-4 pt-5 pb-4 sm:px-6 border-b border-gray-200 flex-shrink-0">
          <!-- Fila principal: icono + título + botones masivos -->
          <div class="flex items-center justify-between gap-4">
            <div class="flex items-center gap-3 min-w-0">
              <div class="flex-shrink-0 flex items-center justify-center h-10 w-10 rounded-full bg-blue-100">
                <svg class="h-6 w-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4" />
                </svg>
              </div>
              <h3 class="text-lg font-medium text-gray-900 truncate" id="modal-title">
                Revisar Insumos de la Solicitud
              </h3>
            </div>
            <!-- Botones masivos -->
            <div v-if="hasPendingItems" class="flex flex-shrink-0 gap-2">
              <button
                @click="openBulkActionModal('aceptado')"
                class="inline-flex items-center px-3 py-1.5 text-xs font-medium rounded-md text-white bg-green-600 hover:bg-green-700 transition-colors"
              >
                Aprobar todo
              </button>
              <button
                @click="openBulkActionModal('rechazado')"
                class="inline-flex items-center px-3 py-1.5 text-xs font-medium rounded-md text-white bg-red-600 hover:bg-red-700 transition-colors"
              >
                Rechazar todo
              </button>
              <button
                @click="openBulkActionModal('devuelto')"
                class="inline-flex items-center px-3 py-1.5 text-xs font-medium rounded-md text-white bg-yellow-600 hover:bg-yellow-700 transition-colors"
              >
                Devolver todo
              </button>
            </div>
          </div>
          <!-- Meta info -->
          <div class="mt-2 ml-13 pl-1 flex flex-wrap gap-x-6 gap-y-1 text-sm text-gray-500">
            <span>Solicitud: <span class="font-semibold text-gray-900">{{ request?.request_number }}</span></span>
            <span>Solicitante: <span class="font-semibold text-gray-900">{{ request?.requested_by_name }}</span></span>
          </div>
        </div>

        <!-- Scrollable body -->
        <div class="overflow-y-auto flex-1 px-4 sm:px-6 py-4">
          <!-- Tabla de items -->
          <div class="overflow-x-auto">
            <table class="min-w-full divide-y divide-gray-200">
              <thead class="bg-gray-50">
                <tr>
                  <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Insumo</th>
                  <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Cantidad</th>
                  <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Estado</th>
                  <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Acciones</th>
                </tr>
              </thead>
              <tbody class="bg-white divide-y divide-gray-200">
                <tr v-for="item in items" :key="item.id" :class="getRowClass(item.item_status)">
                  <td class="px-4 py-4 whitespace-nowrap">
                    <div class="text-sm font-medium text-gray-900">{{ item.supply_name }}</div>
                    <div class="text-sm text-gray-500">Código: {{ item.supply_code }}</div>
                  </td>
                  <td class="px-4 py-4 whitespace-nowrap text-sm text-gray-900">
                    {{ item.quantity_requested }}
                  </td>
                  <td class="px-4 py-4 whitespace-nowrap">
                    <span :class="getStatusBadgeClass(item.item_status)" class="inline-flex px-2 py-1 text-xs font-semibold rounded-full">
                      {{ getStatusLabel(item.item_status) }}
                    </span>
                  </td>
                  <td class="px-4 py-4 whitespace-nowrap text-sm font-medium">
                    <div v-if="item.item_status === 'pendiente'" class="flex space-x-2">
                      <button
                        @click="openActionModal(item, 'aceptado')"
                        class="text-green-600 hover:text-green-900"
                        title="Aceptar"
                      >
                        <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                        </svg>
                      </button>
                      <button
                        @click="openActionModal(item, 'rechazado')"
                        class="text-red-600 hover:text-red-900"
                        title="Rechazar"
                      >
                        <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                        </svg>
                      </button>
                      <button
                        @click="openActionModal(item, 'devuelto')"
                        class="text-yellow-600 hover:text-yellow-900"
                        title="Devolver"
                      >
                        <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 10h10a8 8 0 018 8v2M3 10l6 6m-6-6l6-6" />
                        </svg>
                      </button>
                    </div>
                    <div v-else class="text-gray-500 text-xs">
                      <div v-if="item.reviewed_by_name">Por: {{ item.reviewed_by_name }}</div>
                      <div v-if="item.item_notes" :class="getItemNotesClass(item.item_status)">{{ item.item_notes }}</div>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>

          <!-- Mensaje de error -->
          <div v-if="errorMessage" class="mt-4 bg-red-50 border border-red-200 rounded-md p-3">
            <div class="flex">
              <div class="flex-shrink-0">
                <svg class="h-5 w-5 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </div>
              <div class="ml-3">
                <p class="text-sm text-red-800">{{ errorMessage }}</p>
              </div>
            </div>
          </div>
        </div>

        <!-- Sticky footer -->
        <div class="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse border-t border-gray-200 flex-shrink-0">
          <button
            type="button"
            @click="closeModal"
            class="w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 sm:ml-3 sm:w-auto sm:text-sm"
          >
            Cerrar
          </button>
        </div>
      </div>

    <!-- Modal de confirmación de acción -->
    <div v-if="showActionModal" class="fixed inset-0 z-[60] overflow-y-auto">
      <div class="flex items-center justify-center min-h-screen px-4">
        <div class="fixed inset-0 bg-gray-500 bg-opacity-75" @click="closeActionModal"></div>
        
        <div class="relative bg-white rounded-lg shadow-xl max-w-md w-full p-6">
          <h4 class="text-lg font-medium text-gray-900 mb-4">
            {{ getActionTitle() }}
          </h4>
          
          <div class="mb-4">
            <template v-if="isBulkAction">
              <p class="text-sm text-gray-600 mb-2">Se aplicará a los siguientes insumos:</p>
              <ul class="divide-y divide-gray-100 border border-gray-200 rounded-md overflow-hidden text-sm">
                <li
                  v-for="item in items.filter(i => i.item_status === 'pendiente')"
                  :key="item.id"
                  class="flex items-center justify-between px-3 py-2 bg-white"
                >
                  <span class="text-gray-800 font-medium">{{ item.supply_name }}</span>
                  <span class="text-gray-500 ml-4 whitespace-nowrap">Cant: {{ item.quantity_requested }}</span>
                </li>
              </ul>
            </template>
            <template v-else>
              <p class="text-sm text-gray-600 mb-2">
                <span class="font-semibold">Insumo:</span> {{ selectedItem?.supply_name }}
              </p>
              <p class="text-sm text-gray-600">
                <span class="font-semibold">Cantidad:</span> {{ selectedItem?.quantity_requested }}
              </p>
            </template>
          </div>

          <!-- Campo de comentario (obligatorio para rechazar/devolver) -->
          <div v-if="selectedAction !== 'aceptado'" class="mb-4">
            <label class="block text-sm font-medium text-gray-700 mb-2">
              {{ selectedAction === 'rechazado' ? 'Motivo del rechazo' : 'Comentario para el solicitante' }}
              <span class="text-red-500">*</span>
            </label>
            <textarea
              v-model="actionComment"
              rows="3"
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              :placeholder="selectedAction === 'rechazado' ? 'Explique por qué se rechaza este insumo...' : 'Explique por qué debe revisar/cambiar esta solicitud...'"
            ></textarea>
            <p v-if="commentError" class="mt-1 text-sm text-red-600">{{ commentError }}</p>
          </div>

          <div class="flex justify-end space-x-3">
            <button
              @click="closeActionModal"
              class="px-4 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-500 transition-colors"
            >
              Cancelar
            </button>
            <button
              @click="confirmAction"
              :disabled="loading"
              :class="getActionButtonClass()"
              class="px-6 py-2 rounded-md text-sm font-medium text-white focus:outline-none focus:ring-2 focus:ring-offset-2 transition-all disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {{ loading ? 'Procesando...' : 'Confirmar' }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useNotification } from '@/composables/useNotification'
import supplyRequestService from '@/services/requests/supplyRequestService'

const props = defineProps({
  show: {
    type: Boolean,
    required: true
  },
  request: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['close', 'itemsReviewed'])

const authStore = useAuthStore()
const { success: showSuccess, error: showError } = useNotification()

// Estado
const loading = ref(false)
const items = ref([])
const errorMessage = ref('')
const showActionModal = ref(false)
const selectedItem = ref(null)
const selectedAction = ref('')
const actionComment = ref('')
const commentError = ref('')
const isBulkAction = ref(false)

// Computed
const hasPendingItems = computed(() => items.value.some(i => i.item_status === 'pendiente'))
const pendingItemsCount = computed(() => items.value.filter(i => i.item_status === 'pendiente').length)

// Métodos
const loadItems = async () => {
  if (!props.request) return
  
  loading.value = true
  errorMessage.value = ''
  
  try {
    const result = await supplyRequestService.getSupplyRequestItems(props.request.id)
    if (result.success && result.data) {
      items.value = result.data
    } else {
      errorMessage.value = result.error || 'No se pudieron cargar los items'
    }
  } catch (error) {
    console.error('Error cargando items:', error)
    errorMessage.value = 'Error al cargar los items de la solicitud'
  } finally {
    loading.value = false
  }
}

const openActionModal = (item, action) => {
  isBulkAction.value = false
  selectedItem.value = item
  selectedAction.value = action
  actionComment.value = ''
  commentError.value = ''
  showActionModal.value = true
}

const openBulkActionModal = (action) => {
  isBulkAction.value = true
  selectedItem.value = null
  selectedAction.value = action
  actionComment.value = ''
  commentError.value = ''
  showActionModal.value = true
}

const closeActionModal = () => {
  showActionModal.value = false
  selectedItem.value = null
  selectedAction.value = ''
  actionComment.value = ''
  commentError.value = ''
  isBulkAction.value = false
}

const getActionTitle = () => {
  if (isBulkAction.value) {
    const actions = {
      'aceptado': 'Aprobar todos los insumos pendientes',
      'rechazado': 'Rechazar todos los insumos pendientes',
      'devuelto': 'Devolver todos los insumos pendientes'
    }
    return actions[selectedAction.value] || 'Acción masiva'
  }
  const actions = {
    'aceptado': 'Aceptar Insumo',
    'rechazado': 'Rechazar Insumo',
    'devuelto': 'Devolver para Revisión'
  }
  return actions[selectedAction.value] || 'Acción'
}

const getActionButtonClass = () => {
  const classes = {
    'aceptado': 'bg-green-600 hover:bg-green-700 focus:ring-green-500',
    'rechazado': 'bg-red-600 hover:bg-red-700 focus:ring-red-500',
    'devuelto': 'bg-yellow-600 hover:bg-yellow-700 focus:ring-yellow-500'
  }
  return classes[selectedAction.value] || 'bg-blue-600 hover:bg-blue-700 focus:ring-blue-500'
}

const getItemNotesClass = (status) => {
  const classes = {
    'rechazado': 'text-red-600',
    'devuelto': 'text-yellow-600',
    'aceptado': 'text-green-600'
  }
  return classes[status] || 'text-gray-600'
}

const confirmAction = async () => {
  // Validar comentario obligatorio para rechazar/devolver
  if (selectedAction.value !== 'aceptado' && !actionComment.value.trim()) {
    commentError.value = 'El comentario es obligatorio'
    return
  }

  loading.value = true
  commentError.value = ''

  try {
    const reviewData = {
      item_status: selectedAction.value,
      reviewed_by: authStore.getUserRut,
      reviewed_by_name: authStore.getUserName,
      comment: actionComment.value.trim() || null
    }

    if (isBulkAction.value) {
      const pendingItems = items.value.filter(i => i.item_status === 'pendiente')
      // Procesar secuencialmente para evitar race conditions en el backend
      // (transacciones concurrentes sobre la misma solicitud causan conflictos en
      //  la creación del carrito y en la actualización del estado de la solicitud)
      for (const item of pendingItems) {
        await supplyRequestService.reviewSupplyRequestItem(item.id, reviewData)
      }
      showSuccess(`${pendingItems.length} insumo${pendingItems.length !== 1 ? 's' : ''} ${getStatusLabel(selectedAction.value).toLowerCase()}${pendingItems.length !== 1 ? 's' : ''}`)
    } else {
      await supplyRequestService.reviewSupplyRequestItem(selectedItem.value.id, reviewData)
      showSuccess(`El insumo ha sido ${getStatusLabel(selectedAction.value).toLowerCase()}`)
    }

    closeActionModal()
    await loadItems() // Recargar items
    
    // Verificar si todos los items han sido resueltos
    const allItemsResolved = items.value.every(item => 
      item.item_status !== 'pendiente'
    )
    
    // Siempre emitir el evento para que el componente padre pueda refrescar
    emit('itemsReviewed')
    
    // Si todos los items están resueltos, cerrar el modal y recargar la página después de un breve delay
    // para asegurar que el estado se actualice en el backend
    if (allItemsResolved) {
      setTimeout(() => {
        emit('close') // Cerrar el modal principal
        window.location.reload()
      }, 1000)
    }
    // Si no todos están resueltos, el modal permanece abierto para seguir revisando items
  } catch (error) {
    console.error('Error revisando item:', error)
    showError(error.response?.data?.error || error.message || 'Error al revisar el item')
  } finally {
    loading.value = false
  }
}

const closeModal = () => {
  if (!loading.value) {
    emit('close')
  }
}

const getStatusLabel = (status) => {
  const labels = {
    'pendiente': 'Pendiente',
    'aceptado': 'Aceptado',
    'rechazado': 'Rechazado',
    'devuelto': 'Devuelto'
  }
  return labels[status] || status
}

const getStatusBadgeClass = (status) => {
  const classes = {
    'pendiente': 'bg-yellow-100 text-yellow-800',
    'aceptado': 'bg-green-100 text-green-800',
    'rechazado': 'bg-red-100 text-red-800',
    'devuelto': 'bg-orange-100 text-orange-800'
  }
  return classes[status] || 'bg-gray-100 text-gray-800'
}

const getRowClass = (status) => {
  const classes = {
    'aceptado': 'bg-green-50',
    'rechazado': 'bg-red-50',
    'devuelto': 'bg-orange-50'
  }
  return classes[status] || ''
}

// Watchers
watch(() => props.show, (newValue) => {
  if (newValue) {
    loadItems()
  }
})
</script>

<style scoped>
/* Animación de entrada del modal */
@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.inline-block {
  animation: fadeIn 0.3s ease-out;
}
</style>
