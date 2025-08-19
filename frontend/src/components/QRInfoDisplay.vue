<template>
  <div class="space-y-4">
    <!-- Información básica del QR -->
    <div class="bg-gray-50 rounded-lg p-4">
      <div class="flex items-center justify-between mb-3">
        <div class="flex items-center">
          <div :class="getTypeIconClass(qrInfo.type)" class="p-2 rounded-lg mr-3">
            <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path v-if="qrInfo.type === 'batch'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
              <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2zM9 9h6v6H9V9z" />
            </svg>
          </div>
          <div>
            <h3 class="font-semibold text-gray-900">{{ getTypeLabel(qrInfo.type) }}</h3>
            <p class="text-sm text-gray-600">ID: {{ qrInfo.id }}</p>
          </div>
        </div>
        <span :class="getTypeBadgeClass(qrInfo.type)" class="badge">
          {{ qrInfo.type.toUpperCase() }}
        </span>
      </div>
      
      <div class="border-t pt-3">
        <p class="text-sm text-gray-700 mb-2">
          <span class="font-medium">Código QR:</span>
        </p>
        <code class="block bg-white px-3 py-2 rounded border text-sm text-gray-800 font-mono">
          {{ qrInfo.qr_code }}
        </code>
      </div>
    </div>

    <!-- Información específica según el tipo -->
    <div v-if="qrInfo.type === 'batch' && qrInfo.batch_info" class="space-y-4">
      <h4 class="font-medium text-gray-900 flex items-center">
        <svg class="h-5 w-5 mr-2 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
        </svg>
        Información del Lote
      </h4>
      
      <div class="grid md:grid-cols-2 gap-4">
        <div class="bg-white border rounded-lg p-4">
          <dl class="space-y-2">
            <div>
              <dt class="text-sm font-medium text-gray-500">Fecha de Vencimiento</dt>
              <dd :class="getExpirationClass(qrInfo.batch_info.expiration_date)" class="text-sm font-semibold">
                {{ formatDate(qrInfo.batch_info.expiration_date) }}
              </dd>
            </div>
            <div>
              <dt class="text-sm font-medium text-gray-500">Cantidad</dt>
              <dd class="text-sm text-gray-900">{{ qrInfo.batch_info.amount }} unidades</dd>
            </div>
          </dl>
        </div>
        
        <div class="bg-white border rounded-lg p-4">
          <dl class="space-y-2">
            <div>
              <dt class="text-sm font-medium text-gray-500">Proveedor</dt>
              <dd class="text-sm text-gray-900">{{ qrInfo.batch_info.supplier }}</dd>
            </div>
            <div>
              <dt class="text-sm font-medium text-gray-500">ID de Bodega</dt>
              <dd class="text-sm text-gray-900">{{ qrInfo.batch_info.store_id }}</dd>
            </div>
          </dl>
        </div>
      </div>
    </div>

    <div v-else-if="qrInfo.type === 'medical_supply' && qrInfo.supply_info" class="space-y-4">
      <h4 class="font-medium text-gray-900 flex items-center">
        <svg class="h-5 w-5 mr-2 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2zM9 9h6v6H9V9z" />
        </svg>
        Información del Insumo Médico
      </h4>
      
      <div class="grid md:grid-cols-2 gap-4">
        <div class="bg-white border rounded-lg p-4">
          <dl class="space-y-2">
            <div>
              <dt class="text-sm font-medium text-gray-500">Nombre del Insumo</dt>
              <dd class="text-sm font-semibold text-gray-900">{{ qrInfo.supply_info.supply_code_name }}</dd>
            </div>
            <div>
              <dt class="text-sm font-medium text-gray-500">Código Interno</dt>
              <dd class="text-sm text-gray-900">{{ qrInfo.supply_info.code }}</dd>
            </div>
            <div>
              <dt class="text-sm font-medium text-gray-500">ID del Lote</dt>
              <dd class="text-sm text-gray-900">{{ qrInfo.supply_info.batch_id }}</dd>
            </div>
          </dl>
        </div>
        
        <div class="bg-white border rounded-lg p-4">
          <dl class="space-y-2">
            <div>
              <dt class="text-sm font-medium text-gray-500">Proveedor</dt>
              <dd class="text-sm text-gray-900">{{ qrInfo.supply_info.supplier }}</dd>
            </div>
            <div>
              <dt class="text-sm font-medium text-gray-500">Fecha de Vencimiento</dt>
              <dd :class="getExpirationClass(qrInfo.supply_info.expiration_date)" class="text-sm font-semibold">
                {{ formatDate(qrInfo.supply_info.expiration_date) }}
              </dd>
            </div>
            <div>
              <dt class="text-sm font-medium text-gray-500">Bodega</dt>
              <dd class="text-sm text-gray-900">{{ qrInfo.supply_info.store_name }}</dd>
            </div>
          </dl>
        </div>
      </div>
    </div>

    <!-- Historial (si está disponible) -->
    <div v-if="qrInfo.history && qrInfo.history.length > 0" class="space-y-2">
      <h4 class="font-medium text-gray-900 flex items-center">
        <svg class="h-5 w-5 mr-2 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        Últimos Movimientos ({{ qrInfo.history.length }})
      </h4>
      
      <div class="bg-white border rounded-lg overflow-hidden">
        <div class="max-h-40 overflow-y-auto">
          <table class="min-w-full">
            <thead class="bg-gray-50">
              <tr>
                <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Fecha</th>
                <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Estado</th>
                <th class="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Destino</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-200">
              <tr v-for="movement in qrInfo.history.slice(0, 3)" :key="movement.id" class="hover:bg-gray-50">
                <td class="px-4 py-2 text-sm text-gray-900">{{ formatDateTime(movement.date_time) }}</td>
                <td class="px-4 py-2">
                  <span :class="getStatusBadgeClass(movement.status)" class="badge text-xs">
                    {{ movement.status }}
                  </span>
                </td>
                <td class="px-4 py-2 text-sm text-gray-700">
                  {{ movement.destination_type }}: {{ movement.destination_id }}
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- Botones de acción -->
    <div class="flex space-x-3 pt-4 border-t">
      <button 
        @click="$emit('view-details', qrInfo)"
        class="btn-primary flex-1"
      >
        <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
        </svg>
        Ver Detalles Completos
      </button>
      
      <button 
        @click="copyQRCode"
        class="btn-secondary px-4"
        :title="copied ? 'Copiado!' : 'Copiar código QR'"
      >
        <svg v-if="!copied" class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
        </svg>
        <svg v-else class="h-4 w-4 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
        </svg>
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { format } from 'date-fns'
import { es } from 'date-fns/locale'

// Props
const props = defineProps({
  qrInfo: {
    type: Object,
    required: true
  }
})

// Emits
const emit = defineEmits(['view-details'])

// Estado reactivo
const copied = ref(false)

// Métodos
const getTypeLabel = (type) => {
  return type === 'batch' ? 'Lote de Insumos' : 'Insumo Médico Individual'
}

const getTypeIconClass = (type) => {
  return type === 'batch' 
    ? 'bg-blue-100 text-blue-600' 
    : 'bg-green-100 text-green-600'
}

const getTypeBadgeClass = (type) => {
  return type === 'batch' 
    ? 'badge-info' 
    : 'badge-success'
}

const getExpirationClass = (expirationDate) => {
  const today = new Date()
  const expDate = new Date(expirationDate)
  const daysUntilExpiration = Math.ceil((expDate - today) / (1000 * 60 * 60 * 24))
  
  if (daysUntilExpiration < 0) return 'text-red-600'
  if (daysUntilExpiration <= 30) return 'text-orange-600'
  return 'text-green-600'
}

const getStatusBadgeClass = (status) => {
  const statusClasses = {
    'activo': 'badge-success',
    'en_transito': 'badge-warning',
    'utilizado': 'badge-info',
    'vencido': 'badge-danger'
  }
  return statusClasses[status] || 'badge-info'
}

const formatDate = (dateString) => {
  try {
    return format(new Date(dateString), 'dd/MM/yyyy', { locale: es })
  } catch {
    return dateString
  }
}

const formatDateTime = (dateString) => {
  try {
    return format(new Date(dateString), 'dd/MM/yyyy HH:mm', { locale: es })
  } catch {
    return dateString
  }
}

const copyQRCode = async () => {
  try {
    await navigator.clipboard.writeText(props.qrInfo.qr_code)
    copied.value = true
    setTimeout(() => {
      copied.value = false
    }, 2000)
  } catch (err) {
    console.error('Error al copiar:', err)
  }
}
</script>