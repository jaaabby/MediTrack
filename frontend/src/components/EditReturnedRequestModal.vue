<template>
  <div v-if="show" class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50 flex items-center justify-center p-4">
    <div class="relative bg-white rounded-lg shadow-xl max-w-4xl w-full max-h-[90vh] overflow-hidden">
      <!-- Header -->
      <div class="bg-gradient-to-r from-blue-600 to-blue-700 px-6 py-4">
        <div class="flex justify-between items-center">
          <div>
            <h3 class="text-xl font-semibold text-white">Editar Solicitud Devuelta</h3>
            <p class="text-blue-100 text-sm mt-1">Revisa y ajusta los items devueltos</p>
          </div>
          <button
            @click="closeModal"
            class="text-white hover:text-gray-200 transition-colors"
          >
            <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
      </div>

      <!-- Body -->
      <div class="p-6 overflow-y-auto max-h-[calc(90vh-200px)]">
        <!-- Mensaje informativo -->
        <div class="mb-4 p-4 bg-yellow-50 border-l-4 border-yellow-400 rounded">
          <div class="flex">
            <svg class="h-5 w-5 text-yellow-400 mr-3 flex-shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-1.964-1.333-2.732 0L3.732 16c-.77 1.333.192 3 1.732 3z" />
            </svg>
            <div class="flex-1">
              <p class="text-sm text-yellow-700">
                <strong>Items aceptados se mantendrán aprobados.</strong> Solo necesitas ajustar los items devueltos.
              </p>
            </div>
          </div>
        </div>

        <!-- Lista de items -->
        <div class="space-y-3">
          <div
            v-for="item in items"
            :key="item.id"
            :class="getItemCardClass(item)"
            class="p-4 rounded-lg border-2"
          >
            <div class="flex items-start justify-between">
              <div class="flex-1">
                <div class="flex items-center space-x-2 mb-2">
                  <span :class="getItemStatusBadgeClass(item.item_status)" class="px-2 py-1 text-xs font-semibold rounded">
                    {{ getItemStatusLabel(item.item_status) }}
                  </span>
                  <span v-if="item.is_pediatric" class="px-2 py-1 text-xs font-semibold rounded bg-pink-100 text-pink-800">
                    Pediátrico
                  </span>
                </div>
                
                <h4 class="font-medium text-gray-900">{{ item.supply_name }}</h4>
                <p class="text-sm text-gray-600">Código: {{ item.supply_code }}</p>
                
                <!-- Mostrar comentario de devolución -->
                <div v-if="item.item_status === 'devuelto' && item.item_notes" class="mt-2 p-2 bg-yellow-50 rounded text-sm">
                  <p class="text-xs font-medium text-gray-700 mb-1">Motivo de devolución:</p>
                  <p class="text-gray-900">{{ item.item_notes }}</p>
                </div>
              </div>

              <!-- Input de cantidad solo para items devueltos -->
              <div v-if="item.item_status === 'devuelto'" class="ml-4">
                <label class="block text-xs font-medium text-gray-700 mb-1">Cantidad</label>
                <input
                  v-model.number="item.new_quantity"
                  type="number"
                  min="1"
                  class="w-24 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>

              <!-- Mostrar cantidad para items aceptados -->
              <div v-else-if="item.item_status === 'aceptado'" class="ml-4 text-right">
                <p class="text-xs text-gray-600">Cantidad aprobada</p>
                <p class="text-lg font-semibold text-green-600">{{ item.quantity_approved || item.quantity_requested }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Footer -->
      <div class="bg-gray-50 px-6 py-4 flex justify-end space-x-3 border-t">
        <button
          @click="closeModal"
          :disabled="loading"
          class="px-4 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 hover:bg-gray-50 disabled:opacity-50"
        >
          Cancelar
        </button>
        <button
          @click="resubmitRequest"
          :disabled="loading || !hasReturnedItems"
          class="btn-primary disabled:opacity-50"
        >
          {{ loading ? 'Reenviando...' : 'Reenviar Solicitud' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import supplyRequestService from '@/services/supplyRequestService'
import Swal from 'sweetalert2'

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

const emit = defineEmits(['close', 'resubmitted'])

const loading = ref(false)
const items = ref([])

const hasReturnedItems = computed(() => {
  return items.value.some(item => item.item_status === 'devuelto')
})

watch(() => props.show, async (newVal) => {
  if (newVal && props.request) {
    await loadItems()
  }
})

const loadItems = async () => {
  try {
    const response = await supplyRequestService.getSupplyRequestItems(props.request.id)
    items.value = response.data.map(item => ({
      ...item,
      new_quantity: item.quantity_requested
    }))
  } catch (error) {
    console.error('Error cargando items:', error)
    Swal.fire({
      icon: 'error',
      title: 'Error',
      text: 'No se pudieron cargar los items de la solicitud'
    })
  }
}

const getItemCardClass = (item) => {
  if (item.item_status === 'aceptado') {
    return 'bg-green-50 border-green-200'
  } else if (item.item_status === 'devuelto') {
    return 'bg-yellow-50 border-yellow-300'
  }
  return 'bg-gray-50 border-gray-200'
}

const getItemStatusBadgeClass = (status) => {
  const classes = {
    'aceptado': 'bg-green-100 text-green-800',
    'devuelto': 'bg-yellow-100 text-yellow-800',
    'rechazado': 'bg-red-100 text-red-800',
    'pendiente': 'bg-gray-100 text-gray-800'
  }
  return classes[status] || 'bg-gray-100 text-gray-800'
}

const getItemStatusLabel = (status) => {
  const labels = {
    'aceptado': 'Aceptado ✓',
    'devuelto': 'Devuelto',
    'rechazado': 'Rechazado',
    'pendiente': 'Pendiente'
  }
  return labels[status] || status
}

const resubmitRequest = async () => {
  // Validar cantidades
  const returnedItems = items.value.filter(item => item.item_status === 'devuelto')
  const invalidItems = returnedItems.filter(item => !item.new_quantity || item.new_quantity < 1)
  
  if (invalidItems.length > 0) {
    Swal.fire({
      icon: 'warning',
      title: 'Cantidades inválidas',
      text: 'Todos los items devueltos deben tener una cantidad válida (mayor a 0)'
    })
    return
  }

  loading.value = true

  try {
    const updatedItems = returnedItems.map(item => ({
      item_id: item.id,
      quantity: item.new_quantity
    }))

    await supplyRequestService.resubmitReturnedRequest(props.request.id, updatedItems)

    await Swal.fire({
      icon: 'success',
      title: 'Solicitud Reenviada',
      text: 'La solicitud ha sido reenviada al encargado de bodega',
      timer: 2000,
      showConfirmButton: false
    })

    emit('resubmitted')
    closeModal()
  } catch (error) {
    console.error('Error reenviando solicitud:', error)
    Swal.fire({
      icon: 'error',
      title: 'Error',
      text: error.response?.data?.error || 'No se pudo reenviar la solicitud'
    })
  } finally {
    loading.value = false
  }
}

const closeModal = () => {
  if (!loading.value) {
    emit('close')
  }
}
</script>
