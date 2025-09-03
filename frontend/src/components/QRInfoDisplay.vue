<template>
  <div class="space-y-6">
    <!-- Información principal del QR -->
    <div class="bg-white rounded-lg shadow border overflow-hidden">
      <!-- Header del QR -->
      <div class="px-6 py-4 bg-gray-50 border-b border-gray-200">
        <div class="flex items-center justify-between">
          <div class="flex items-center">
            <div :class="getTypeIconClass(qrInfo.type)" class="w-10 h-10 rounded-full flex items-center justify-center mr-3">
              <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path v-if="qrInfo.type === 'medical_supply'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
              </svg>
            </div>
            <div>
              <h3 class="text-lg font-semibold text-gray-900">
                {{ getTypeLabel(qrInfo.type) }}
              </h3>
              <p class="text-sm text-gray-600">Código: {{ qrInfo.qr_code }}</p>
            </div>
          </div>
          
          <div class="flex items-center space-x-2">
            <span :class="getStatusBadgeClass()" class="inline-flex px-3 py-1 text-sm font-semibold rounded-full">
              {{ getStatusLabel() }}
            </span>
            
            <!-- Nuevo botón de trazabilidad avanzada -->
            <router-link
              :to="`/qr/${qrInfo.qr_code}/traceability`"
              class="inline-flex items-center px-3 py-2 border border-blue-300 text-sm font-medium rounded-md text-blue-700 bg-blue-50 hover:bg-blue-100 transition-colors"
            >
              <svg class="h-4 w-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
              </svg>
              Ver Trazabilidad
            </router-link>
          </div>
        </div>
      </div>

      <!-- Contenido del QR -->
      <div class="p-6">
        <!-- Información de asignación a solicitud -->
        <div v-if="qrInfo.request_assignment" class="mb-6 bg-blue-50 border border-blue-200 rounded-lg p-4">
          <div class="flex items-center mb-3">
            <svg class="h-5 w-5 text-blue-600 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
            </svg>
            <h4 class="font-semibold text-blue-900">Asignado a Solicitud</h4>
          </div>
          
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <label class="text-sm font-medium text-blue-700">Número de Solicitud:</label>
              <router-link
                v-if="qrInfo.supply_request"
                :to="`/supply-requests/${qrInfo.supply_request.id}`"
                class="block text-sm text-blue-900 hover:text-blue-700 underline"
              >
                {{ qrInfo.supply_request.request_number }}
              </router-link>
              <p v-else class="text-sm text-blue-900">N/A</p>
            </div>
            
            <div>
              <label class="text-sm font-medium text-blue-700">Estado de Asignación:</label>
              <span :class="getAssignmentStatusBadgeClass(qrInfo.request_assignment.status)" class="inline-flex px-2 py-1 text-xs font-semibold rounded-full ml-1">
                {{ getAssignmentStatusLabel(qrInfo.request_assignment.status) }}
              </span>
            </div>
            
            <div>
              <label class="text-sm font-medium text-blue-700">Fecha de Asignación:</label>
              <p class="text-sm text-blue-900">{{ formatDate(qrInfo.request_assignment.assigned_date) }}</p>
            </div>
            
            <div v-if="qrInfo.request_assignment.delivered_date">
              <label class="text-sm font-medium text-blue-700">Fecha de Entrega:</label>
              <p class="text-sm text-blue-900">{{ formatDate(qrInfo.request_assignment.delivered_date) }}</p>
            </div>
          </div>
          
          <div v-if="qrInfo.request_assignment.notes" class="mt-3">
            <label class="text-sm font-medium text-blue-700">Notas:</label>
            <p class="text-sm text-blue-900 bg-white p-2 rounded border mt-1">{{ qrInfo.request_assignment.notes }}</p>
          </div>
        </div>

        <!-- Resto del contenido existente -->
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <!-- Información del insumo individual -->
          <div v-if="qrInfo.supply_info" class="space-y-4">
            <h4 class="font-semibold text-gray-900 flex items-center">
              <svg class="h-5 w-5 mr-2 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
              Información del Insumo
            </h4>

            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 text-sm">
              <div>
                <label class="font-medium text-gray-600">ID Insumo:</label>
                <p class="text-gray-900">{{ qrInfo.supply_info.id }}</p>
              </div>
              <div>
                <label class="font-medium text-gray-600">Código:</label>
                <p class="text-gray-900">{{ qrInfo.supply_info.code }}</p>
              </div>
              <div>
                <label class="font-medium text-gray-600">Estado:</label>
                <span :class="qrInfo.supply_info.is_consumed ? 'text-red-600' : 'text-green-600'" class="font-medium">
                  {{ qrInfo.supply_info.is_consumed ? 'Consumido' : 'Disponible' }}
                </span>
              </div>
              <div v-if="qrInfo.supply_info.batch">
                <label class="font-medium text-gray-600">ID Lote:</label>
                <p class="text-gray-900">{{ qrInfo.supply_info.batch.id }}</p>
              </div>
            </div>
          </div>

          <!-- Información del lote -->
          <div v-if="qrInfo.batch_info" class="space-y-4">
            <h4 class="font-semibold text-gray-900 flex items-center">
              <svg class="h-5 w-5 mr-2 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
              </svg>
              Información del Lote
            </h4>

            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 text-sm">
              <div>
                <label class="font-medium text-gray-600">ID Lote:</label>
                <p class="text-gray-900">{{ qrInfo.batch_info.id }}</p>
              </div>
              <div>
                <label class="font-medium text-gray-600">Proveedor:</label>
                <p class="text-gray-900">{{ qrInfo.batch_info.supplier }}</p>
              </div>
              <div>
                <label class="font-medium text-gray-600">Cantidad Total:</label>
                <p class="text-gray-900">{{ qrInfo.batch_info.amount }}</p>
              </div>
              <div>
                <label class="font-medium text-gray-600">Vencimiento:</label>
                <p class="text-gray-900">{{ formatDate(qrInfo.batch_info.expiration_date) }}</p>
              </div>
            </div>
          </div>
        </div>

        <!-- Información del código de insumo -->
        <div v-if="qrInfo.supply_code" class="mt-6 bg-gray-50 rounded-lg p-4">
          <h4 class="font-semibold text-gray-900 mb-3 flex items-center">
            <svg class="h-5 w-5 mr-2 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 20l4-16m2 16l4-16M6 9h14M4 15h14" />
            </svg>
            Código de Insumo
          </h4>

          <div class="grid grid-cols-1 sm:grid-cols-2 gap-4 text-sm">
            <div>
              <label class="font-medium text-gray-600">Código:</label>
              <p class="text-gray-900">{{ qrInfo.supply_code.code }}</p>
            </div>
            <div class="sm:col-span-2">
              <label class="font-medium text-gray-600">Nombre:</label>
              <p class="text-gray-900">{{ qrInfo.supply_code.name }}</p>
            </div>
            <div class="sm:col-span-2">
              <label class="font-medium text-gray-600">Código de Proveedor:</label>
              <p class="text-gray-900">{{ qrInfo.supply_code.code_supplier }}</p>
            </div>
          </div>
        </div>

        <!-- Acciones rápidas -->
        <div v-if="qrInfo.type === 'medical_supply' && !qrInfo.supply_info?.is_consumed" class="mt-6">
          <h4 class="font-semibold text-gray-900 mb-3">Acciones Disponibles</h4>
          <div class="flex flex-wrap gap-3">
            <!-- Crear solicitud para este insumo -->
            <router-link
              to="/supply-requests/new"
              :state="{ preSelectedSupplyCode: qrInfo.supply_info?.code, preSelectedSupplyName: qrInfo.supply_code?.name }"
              class="inline-flex items-center px-4 py-2 border border-green-300 text-sm font-medium rounded-md text-green-700 bg-green-50 hover:bg-green-100 transition-colors"
            >
              <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
              </svg>
              Crear Solicitud
            </router-link>

            <!-- Consumir insumo (acción existente) -->
            <button
              v-if="!qrInfo.request_assignment"
              @click="$emit('consume-supply', qrInfo.qr_code)"
              class="inline-flex items-center px-4 py-2 border border-red-300 text-sm font-medium rounded-md text-red-700 bg-red-50 hover:bg-red-100 transition-colors"
            >
              <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
              </svg>
              Consumir Insumo
            </button>

            <!-- Marcar como entregado (si está asignado pero no entregado) -->
            <button
              v-if="qrInfo.request_assignment && qrInfo.request_assignment.status === 'assigned'"
              @click="$emit('mark-as-delivered', qrInfo.qr_code)"
              class="inline-flex items-center px-4 py-2 border border-purple-300 text-sm font-medium rounded-md text-purple-700 bg-purple-50 hover:bg-purple-100 transition-colors"
            >
              <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
              Marcar como Entregado
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Historial de movimientos (versión resumida) -->
    <div v-if="qrInfo.history && qrInfo.history.length > 0" class="bg-white rounded-lg shadow border">
      <div class="px-6 py-4 border-b border-gray-200 flex justify-between items-center">
        <h4 class="font-semibold text-gray-900">Historial Reciente</h4>
        <router-link
          :to="`/qr/${qrInfo.qr_code}/traceability`"
          class="text-sm text-blue-600 hover:text-blue-800"
        >
          Ver historial completo →
        </router-link>
      </div>
      
      <div class="p-6">
        <div class="flow-root">
          <ul class="-mb-8">
            <li
              v-for="(movement, index) in qrInfo.history.slice(0, 3)"
              :key="movement.id || index"
              class="relative pb-8"
            >
              <span
                v-if="index !== Math.min(qrInfo.history.length - 1, 2)"
                class="absolute top-4 left-4 -ml-px h-full w-0.5 bg-gray-200"
                aria-hidden="true"
              ></span>
              
              <div class="relative flex space-x-3">
                <div>
                  <span :class="getHistoryIconClass(movement.status || movement.Status)" class="h-8 w-8 rounded-full flex items-center justify-center ring-8 ring-white">
                    <svg class="h-4 w-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                  </span>
                </div>
                <div class="min-w-0 flex-1 pt-1.5">
                  <div>
                    <p class="text-sm text-gray-900">
                      <span class="font-medium">{{ movement.status || movement.Status }}</span>
                    </p>
                    <p class="text-sm text-gray-500">{{ formatDate(movement.date_time || movement.DateTime) }}</p>
                  </div>
                </div>
              </div>
            </li>
          </ul>
        </div>
        
        <div v-if="qrInfo.history.length > 3" class="mt-4 text-center">
          <p class="text-sm text-gray-500">{{ qrInfo.history.length - 3 }} movimientos adicionales</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { format } from 'date-fns'
import { es } from 'date-fns/locale'
import qrService from '../services/qrService'

// Props
const props = defineProps({
  qrInfo: {
    type: Object,
    required: true
  }
})

// Emits
const emit = defineEmits(['consume-supply', 'mark-as-delivered'])

// Métodos de formato
const formatDate = (dateString) => {
  if (!dateString) return 'No disponible'
  try {
    return format(new Date(dateString), 'dd/MM/yyyy HH:mm', { locale: es })
  } catch (error) {
    return dateString
  }
}

const getTypeLabel = (type) => {
  const labels = {
    'batch': 'Lote de Productos',
    'medical_supply': 'Insumo Individual'
  }
  return labels[type] || type
}

const getStatusLabel = () => {
  return qrService.getStatusLabel(props.qrInfo)
}

const getStatusBadgeClass = () => {
  const color = qrService.getStatusColor(props.qrInfo)
  const classes = {
    'green': 'bg-green-100 text-green-800',
    'red': 'bg-red-100 text-red-800',
    'yellow': 'bg-yellow-100 text-yellow-800',
    'gray': 'bg-gray-100 text-gray-800'
  }
  return classes[color] || classes.gray
}

const getTypeIconClass = (type) => {
  const classes = {
    'batch': 'bg-blue-100 text-blue-600',
    'medical_supply': 'bg-green-100 text-green-600'
  }
  return classes[type] || 'bg-gray-100 text-gray-600'
}

const getHistoryIconClass = (status) => {
  const classes = {
    'consumido': 'bg-red-500',
    'creado': 'bg-green-500',
    'movido': 'bg-blue-500'
  }
  return classes[status] || 'bg-gray-500'
}

// Nuevos métodos para solicitudes
const getAssignmentStatusLabel = (status) => {
  const labels = {
    'assigned': 'Asignado',
    'delivered': 'Entregado',
    'consumed': 'Consumido',
    'returned': 'Devuelto'
  }
  return labels[status] || status
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
</script>

<style scoped>
/* Estilos específicos */
.flow-root {
  overflow: hidden;
}

.transition-colors {
  transition: color 0.2s ease-in-out, background-color 0.2s ease-in-out;
}

/* Ring styles for timeline */
.ring-8 {
  box-shadow: 0 0 0 8px #fff;
}
</style>