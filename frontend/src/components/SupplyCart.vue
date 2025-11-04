<template>
  <div class="bg-white rounded-lg shadow-md p-6">
    <!-- Header del carrito -->
    <div class="flex justify-between items-center mb-6">
      <div>
        <h2 class="text-2xl font-bold text-gray-800">
          🛒 Carrito de Insumos
        </h2>
        <p class="text-sm text-gray-500 mt-1">
          {{ cart?.cart_number || 'Cargando...' }}
        </p>
      </div>
      <div class="flex items-center gap-3">
        <span 
          class="px-4 py-2 rounded-full text-sm font-semibold"
          :class="getStatusClass(cart?.status)"
        >
          {{ getStatusLabel(cart?.status) }}
        </span>
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
        <button
          v-if="showAddButton"
          @click="$emit('add-item')"
          class="px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors"
        >
          + Agregar Item
        </button>
      </div>

      <div class="grid grid-cols-1 gap-4">
        <div
          v-for="item in activeItems"
          :key="item.id"
          class="border border-gray-200 rounded-lg p-4 hover:shadow-md transition-shadow"
        >
          <div class="flex justify-between items-start">
            <div class="flex-1">
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
              <div v-if="getItemStatus(item)" class="mt-3">
                <span 
                  class="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium"
                  :class="getItemStatusClass(item)"
                >
                  {{ getItemStatus(item) }}
                </span>
              </div>

              <!-- Notas si existen -->
              <div v-if="item.notes" class="mt-3 p-2 bg-yellow-50 border-l-4 border-yellow-400 text-sm text-gray-700">
                <strong>Nota:</strong> {{ item.notes }}
              </div>

              <!-- Botones de acción para gestión del item -->
              <div v-if="cart?.status === 'active' && canManageItems && !isItemProcessed(item)" class="mt-4 flex gap-2">
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
  }
})

const emit = defineEmits(['cart-loaded', 'cart-closed', 'item-removed', 'item-used', 'item-returned', 'add-item', 'error'])

const cart = ref(null)
const loading = ref(false)
const error = ref(null)

const activeItems = computed(() => {
  if (!cart.value?.items) return []
  return cart.value.items.filter(item => item.is_active)
})

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
  if (!confirm('¿Está seguro de cerrar este carrito? Esta acción no se puede deshacer.')) {
    return
  }
  
  loading.value = true
  try {
    const response = await cartService.closeCart(cart.value.id)
    if (response.success) {
      await loadCart() // Recargar para ver el estado actualizado
      emit('cart-closed', cart.value)
    }
  } catch (err) {
    console.error('Error al cerrar carrito:', err)
    alert(err.response?.data?.message || 'Error al cerrar el carrito')
  } finally {
    loading.value = false
  }
}

const handleRemoveItem = async (itemId) => {
  if (!confirm('¿Está seguro de remover este item del carrito?')) {
    return
  }
  
  loading.value = true
  try {
    const response = await cartService.removeItemFromCart(cart.value.id, itemId)
    if (response.success) {
      await loadCart() // Recargar para ver los cambios
      emit('item-removed', itemId)
    }
  } catch (err) {
    console.error('Error al remover item:', err)
    alert(err.response?.data?.message || 'Error al remover el item')
  } finally {
    loading.value = false
  }
}

const handleMarkAsUsed = async (itemId) => {
  if (!confirm('¿Confirma que este insumo fue utilizado?')) {
    return
  }
  
  loading.value = true
  try {
    const response = await cartService.markItemAsUsed(cart.value.id, itemId)
    if (response.success) {
      await loadCart() // Recargar para ver los cambios
      emit('item-used', itemId)
      alert('Insumo marcado como utilizado exitosamente')
    }
  } catch (err) {
    console.error('Error al marcar item como utilizado:', err)
    alert(err.response?.data?.message || 'Error al marcar el item como utilizado')
  } finally {
    loading.value = false
  }
}

const handleMarkForReturn = async (itemId) => {
  const reason = prompt('Motivo de la devolución (opcional):')
  if (reason === null) return // Usuario canceló
  
  loading.value = true
  try {
    const response = await cartService.markItemForReturn(cart.value.id, itemId, reason)
    if (response.success) {
      await loadCart() // Recargar para ver los cambios
      emit('item-returned', itemId)
      alert('Insumo marcado para devolución exitosamente')
    }
  } catch (err) {
    console.error('Error al marcar item para devolución:', err)
    alert(err.response?.data?.message || 'Error al marcar el item para devolución')
  } finally {
    loading.value = false
  }
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

const copyToClipboard = async (text) => {
  try {
    await navigator.clipboard.writeText(text)
    alert('Código QR copiado al portapapeles')
  } catch (err) {
    console.error('Error al copiar:', err)
  }
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
