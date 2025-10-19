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
        class="inline-flex items-center px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700"
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
            <div class="flex flex-col sm:flex-row sm:items-center gap-2 sm:gap-4">
              <div class="flex items-center gap-2 flex-wrap">
                <span :class="getStatusBadgeClass(request.status)" class="inline-flex px-2 sm:px-3 py-1 text-xs sm:text-sm font-semibold rounded-full whitespace-nowrap">
                  {{ getStatusLabel(request.status) }}
                </span>
                <span :class="getPriorityBadgeClass(request.priority)" class="inline-flex px-2 sm:px-3 py-1 text-xs sm:text-sm font-semibold rounded-full whitespace-nowrap">
                  {{ getPriorityLabel(request.priority) }}
                </span>
              </div>
              <span class="text-xs sm:text-sm text-gray-600">
                Creada: {{ formatDate(request.request_date) }}
              </span>
            </div>
          </div>

          <!-- Acciones -->
          <div class="flex flex-col sm:flex-row gap-2 sm:flex-shrink-0">
            <button
              v-if="request.status === 'pending' && authStore.canViewAllRequests"
              @click="approveRequest"
              :disabled="processing"
              class="inline-flex items-center justify-center px-3 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-green-600 hover:bg-green-700 disabled:opacity-50 w-full sm:w-auto"
            >
              <svg class="h-4 w-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
              Aprobar
            </button>

            <button
              v-if="request.status === 'pending' && authStore.canViewAllRequests"
              @click="rejectRequest"
              :disabled="processing"
              class="inline-flex items-center justify-center px-3 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-red-600 hover:bg-red-700 disabled:opacity-50 w-full sm:w-auto"
            >
              <svg class="h-4 w-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
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
                <label class="block text-xs sm:text-sm font-medium text-gray-700">Pabellón</label>
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
              <div v-if="request.approval_date">
                <label class="block text-xs sm:text-sm font-medium text-gray-700">Fecha de Aprobación</label>
                <p class="text-sm text-gray-900 mt-1">{{ formatDate(request.approval_date) }}</p>
                <p class="text-xs text-gray-500">{{ request.approved_by_name }}</p>
              </div>
            </div>
            
            <!-- Observaciones -->
            <div v-if="request.notes" class="mt-3 sm:mt-4">
              <label class="block text-xs sm:text-sm font-medium text-gray-700">Observaciones</label>
              <p class="text-sm text-gray-900 mt-1 bg-gray-50 p-2 sm:p-3 rounded">{{ request.notes }}</p>
            </div>
          </div>

          <!-- Items solicitados -->
          <div class="bg-white rounded-lg shadow border">
            <div class="px-4 sm:px-6 py-3 sm:py-4 border-b border-gray-200">
              <h3 class="text-base sm:text-lg font-semibold text-gray-900">Insumos Solicitados</h3>
            </div>
            <div class="divide-y divide-gray-200">
              <div v-for="(item, index) in items" :key="item.id" class="p-4 sm:p-6">
                <div class="flex flex-col sm:flex-row sm:justify-between sm:items-start gap-2 sm:gap-4 mb-3 sm:mb-4">
                  <div class="flex-1 min-w-0">
                    <h4 class="text-sm font-medium text-gray-900 truncate">{{ item.supply_name }}</h4>
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
                  <div class="sm:col-span-2">
                    <span :class="getUrgencyBadgeClass(item.urgency_level)" class="inline-flex px-2 py-1 text-xs font-semibold rounded-full whitespace-nowrap">
                      Urgencia: {{ getPriorityLabel(item.urgency_level) }}
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

                <!-- QRs asignados a este item -->
                <div v-if="getItemAssignments(item.id).length > 0" class="mt-3 sm:mt-4">
                  <label class="block text-xs sm:text-sm font-medium text-gray-700 mb-2">QRs Asignados</label>
                  <div class="grid grid-cols-1 gap-2">
                    <div
                      v-for="assignment in getItemAssignments(item.id)"
                      :key="assignment.id"
                      class="bg-blue-50 p-2 sm:p-3 rounded border"
                    >
                      <div class="flex flex-col sm:flex-row sm:justify-between sm:items-center gap-2">
                        <code class="text-xs font-mono break-all">{{ assignment.qr_code }}</code>
                        <span :class="getAssignmentStatusBadgeClass(assignment.status)" class="inline-flex px-2 py-1 text-xs font-semibold rounded-full whitespace-nowrap self-start sm:self-auto">
                          {{ getAssignmentStatusLabel(assignment.status) }}
                        </span>
                      </div>
                      <div v-if="assignment.assigned_date" class="text-xs text-gray-500 mt-1">
                        Asignado: {{ formatDate(assignment.assigned_date) }}
                      </div>
                      <button
                        @click="viewQRTraceability(assignment.qr_code)"
                        class="text-xs text-blue-600 hover:text-blue-800 mt-2 block"
                      >
                        Ver trazabilidad →
                      </button>
                    </div>
                  </div>
                </div>

                <!-- Botón para asignar QR -->
                <div v-if="request.status === 'approved' && getItemAssignments(item.id).length < item.quantity_approved && authStore.canViewAllRequests" class="mt-3 sm:mt-4">
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
        </div>

        <!-- Panel lateral -->
        <div class="space-y-4 sm:space-y-6">
          <!-- Resumen -->
          <div class="bg-white rounded-lg shadow border p-4 sm:p-6">
            <h3 class="text-base sm:text-lg font-semibold text-gray-900 mb-3 sm:mb-4">Resumen</h3>
            <div class="space-y-2 sm:space-y-3">
              <div class="flex justify-between">
                <span class="text-xs sm:text-sm text-gray-600">Total Items:</span>
                <span class="text-xs sm:text-sm font-medium text-gray-900">{{ items.length }}</span>
              </div>
              <!--<div class="flex justify-between">
                <span class="text-xs sm:text-sm text-gray-600">QRs Asignados:</span>
                <span class="text-xs sm:text-sm font-medium text-gray-900">{{ assignments.length }}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-xs sm:text-sm text-gray-600">QRs Entregados:</span>
                <span class="text-xs sm:text-sm font-medium text-green-600">{{ getDeliveredAssignments().length }}</span>
              </div>-->
              <div class="flex justify-between">
                <span class="text-xs sm:text-sm text-gray-600">Progreso:</span>
                <span class="text-xs sm:text-sm font-medium text-blue-600">{{ getProgressPercentage() }}%</span>
              </div>
            </div>
            
            <!-- Barra de progreso -->
            <div class="mt-3 sm:mt-4">
              <div class="bg-gray-200 rounded-full h-2">
                <div
                  class="bg-blue-600 h-2 rounded-full transition-all duration-300"
                  :style="`width: ${getProgressPercentage()}%`"
                ></div>
              </div>
            </div>
          </div>

          <!-- Trazabilidad QR -->
          <div v-if="selectedQRTraceability" class="bg-white rounded-lg shadow border p-4 sm:p-6">
            <div class="flex justify-between items-center mb-3 sm:mb-4">
              <h3 class="text-base sm:text-lg font-semibold text-gray-900">Trazabilidad QR</h3>
              <button
                @click="selectedQRTraceability = null"
                class="text-gray-400 hover:text-gray-600"
              >
                <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
            
            <div class="space-y-2 sm:space-y-3">
              <div>
                <label class="block text-xs sm:text-sm font-medium text-gray-700">Código QR</label>
                <code class="text-xs sm:text-sm bg-gray-100 px-2 py-1 rounded break-all">{{ selectedQRTraceability.qr_code }}</code>
              </div>
              <div>
                <label class="block text-xs sm:text-sm font-medium text-gray-700">Estado Actual</label>
                <span :class="getAssignmentStatusBadgeClass(selectedQRTraceability.current_status)" class="inline-flex px-2 py-1 text-xs font-semibold rounded-full mt-1">
                  {{ getAssignmentStatusLabel(selectedQRTraceability.current_status) }}
                </span>
              </div>
              
              <!-- Historial de movimientos -->
              <div v-if="selectedQRTraceability.assignments && selectedQRTraceability.assignments.length > 0">
                <label class="block text-xs sm:text-sm font-medium text-gray-700 mb-2">Historial</label>
                <div class="space-y-2">
                  <div
                    v-for="assignment in selectedQRTraceability.assignments"
                    :key="assignment.id"
                    class="border-l-2 border-blue-200 pl-2 sm:pl-3 pb-2"
                  >
                    <div class="text-xs sm:text-sm font-medium text-gray-900">{{ assignment.status }}</div>
                    <div class="text-xs text-gray-500">{{ formatDate(assignment.assigned_date) }}</div>
                    <div v-if="assignment.notes" class="text-xs text-gray-600 mt-1">{{ assignment.notes }}</div>
                  </div>
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
              class="w-full sm:w-auto px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed order-1 sm:order-2"
            >
              {{ assigningQR ? 'Asignando...' : 'Asignar' }}
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
import supplyRequestService from '../services/supplyRequestService'
import pavilionService from '../services/pavilionService'
import { format } from 'date-fns'
import { es } from 'date-fns/locale'
import Swal from 'sweetalert2'

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
const showAssignQRModal = ref(false)
const selectedItem = ref(null)

// Formulario de asignación QR
const qrAssignmentForm = reactive({
  qrCode: '',
  assignedBy: 'ADMIN',
  assignedByName: 'Sistema Admin',
  notes: ''
})

// Computed
const requestId = computed(() => parseInt(route.params.id))

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
      console.log('Solicitud cargada:', result.data)
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
  const result = await Swal.fire({
    title: '¿Está seguro de aprobar esta solicitud?',
    icon: 'question',
    showCancelButton: true,
    confirmButtonText: 'Sí, aprobar',
    cancelButtonText: 'Cancelar',
  })
  if (!result.isConfirmed) return

  processing.value = true
  try {
    const approvalData = {
      approved_by: 'ADMIN',
      approved_by_name: 'Sistema Admin',
      approval_notes: 'Aprobada desde interfaz web'
    }
    
    await supplyRequestService.approveSupplyRequest(requestId.value, approvalData)
    await loadSupplyRequest()
  } catch (err) {
    console.error('Error aprobando solicitud:', err)
    Swal.fire({
      icon: 'error',
      title: 'Error al aprobar la solicitud',
      text: err.response?.data?.error || err.message
    })
  } finally {
    processing.value = false
  }
}

const rejectRequest = async () => {
  const { value: reason } = await Swal.fire({
    title: 'Motivo del rechazo',
    input: 'text',
    inputLabel: 'Ingrese el motivo del rechazo:',
    inputPlaceholder: 'Motivo...',
    showCancelButton: true,
    confirmButtonText: 'Rechazar',
    cancelButtonText: 'Cancelar',
    inputValidator: (value) => {
      if (!value) return 'Debe ingresar un motivo';
    }
  })
  if (!reason) return

  processing.value = true
  try {
    const rejectionData = {
      rejected_by: 'ADMIN',
      rejected_by_name: 'Sistema Admin',
      rejection_reason: reason
    }
    
    await supplyRequestService.rejectSupplyRequest(requestId.value, rejectionData)
    await loadSupplyRequest()
  } catch (err) {
    console.error('Error rechazando solicitud:', err)
    Swal.fire({
      icon: 'error',
      title: 'Error al rechazar la solicitud',
      text: err.response?.data?.error || err.message
    })
  } finally {
    processing.value = false
  }
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
  } catch (err) {
    console.error('Error asignando QR:', err)
    Swal.fire({
      icon: 'error',
      title: 'Error al asignar QR',
      text: err.response?.data?.error || err.message
    })
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
  const totalItems = items.value.reduce((sum, item) => sum + (item.quantity_approved || item.quantity_requested), 0)
  const deliveredItems = getDeliveredAssignments().length
  return totalItems > 0 ? Math.round((deliveredItems / totalItems) * 100) : 0
}

const getPavilionName = (pavilionId) => {
  const pavilion = pavilions.value.find(p => p.id === pavilionId)
  return pavilion ? pavilion.name : `Pabellón ${pavilionId}`
}

const formatDate = (dateString) => {
  if (!dateString) return 'N/A'
  try {
    return format(new Date(dateString), 'dd/MM/yyyy HH:mm', { locale: es })
  } catch {
    return dateString
  }
}

// Métodos de estilo
const getStatusLabel = (status) => supplyRequestService.getStatusLabel(status)
const getPriorityLabel = (priority) => supplyRequestService.getPriorityLabel(priority)

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

const getPriorityBadgeClass = (priority) => {
  const color = supplyRequestService.getPriorityColor(priority)
  const classes = {
    'gray': 'bg-gray-100 text-gray-800',
    'blue': 'bg-blue-100 text-blue-800',
    'orange': 'bg-orange-100 text-orange-800',
    'red': 'bg-red-100 text-red-800'
  }
  return classes[color] || classes.blue
}

const getUrgencyBadgeClass = (urgency) => {
  return getPriorityBadgeClass(urgency)
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