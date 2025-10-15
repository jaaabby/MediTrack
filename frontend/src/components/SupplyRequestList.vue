<template>
  <div class="max-w-7xl mx-auto p-6">
    <!-- Encabezado -->
    <div class="mb-6">
      <div class="flex justify-between items-center">
        <div>
          <h1 class="text-2xl font-bold text-gray-900">Solicitudes de Insumo</h1>
          <p class="mt-1 text-gray-600">Gestión de solicitudes con trazabilidad QR</p>
        </div>
        <router-link
          v-if="authStore.canCreateRequests"
          to="/supply-requests/new"
          class="inline-flex items-center px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500"
        >
          <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          Nueva Solicitud
        </router-link>
      </div>
    </div>

    <!-- Filtros y búsqueda -->
    <div class="mb-6 bg-white p-4 rounded-lg shadow border">
      <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
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
            <option value="pending">Pendiente</option>
            <option value="approved">Aprobada</option>
            <option value="rejected">Rechazada</option>
            <option value="in_process">En Proceso</option>
            <option value="completed">Completada</option>
            <option value="cancelled">Cancelada</option>
          </select>
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
            @click="loadSupplyRequests"
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
    <div v-if="stats" class="mb-6 grid grid-cols-1 md:grid-cols-4 gap-4">
      <div class="bg-white p-4 rounded-lg shadow border">
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <div class="w-8 h-8 bg-yellow-100 rounded-full flex items-center justify-center">
              <svg class="w-4 h-4 text-yellow-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
          </div>
          <div class="ml-3">
            <p class="text-sm font-medium text-gray-500">Pendientes</p>
            <p class="text-lg font-semibold text-gray-900">{{ stats.pending || 0 }}</p>
          </div>
        </div>
      </div>

      <div class="bg-white p-4 rounded-lg shadow border">
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <div class="w-8 h-8 bg-green-100 rounded-full flex items-center justify-center">
              <svg class="w-4 h-4 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
              </svg>
            </div>
          </div>
          <div class="ml-3">
            <p class="text-sm font-medium text-gray-500">Aprobadas</p>
            <p class="text-lg font-semibold text-gray-900">{{ stats.approved || 0 }}</p>
          </div>
        </div>
      </div>

      <div class="bg-white p-4 rounded-lg shadow border">
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <div class="w-8 h-8 bg-red-100 rounded-full flex items-center justify-center">
              <svg class="w-4 h-4 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </div>
          </div>
          <div class="ml-3">
            <p class="text-sm font-medium text-gray-500">Rechazadas</p>
            <p class="text-lg font-semibold text-gray-900">{{ stats.rejected || 0 }}</p>
          </div>
        </div>
      </div>

      <div class="bg-white p-4 rounded-lg shadow border">
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <div class="w-8 h-8 bg-gray-100 rounded-full flex items-center justify-center">
              <svg class="w-4 h-4 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
              </svg>
            </div>
          </div>
          <div class="ml-3">
            <p class="text-sm font-medium text-gray-500">Total</p>
            <p class="text-lg font-semibold text-gray-900">{{ totalRequests }}</p>
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
          class="inline-flex items-center px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700"
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
          {{ requests.length === 0 ? 'No se han creado solicitudes aún' : 'No se encontraron solicitudes con los filtros aplicados' }}
        </p>
        <router-link
          v-if="requests.length === 0 && authStore.canCreateRequests"
          to="/supply-requests/new"
          class="inline-flex items-center px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700"
        >
          <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          Crear Primera Solicitud
        </router-link>
      </div>

      <!-- Table -->
      <div v-else class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Solicitud
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Pabellón
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Solicitante
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Estado
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Fecha de Cirugía
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Fecha
              </th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Items
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
                <div class="text-sm font-medium text-gray-900">
                  {{ request.request_number }}
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
                  {{ formatSurgeryDateTime(request.surgery_datetime) }}
                </div>
                <span :class="getUrgencyBadgeClass(request.surgery_datetime)" class="inline-flex px-2 py-1 text-xs font-semibold rounded-full mt-1">
                  {{ getUrgencyLabel(request.surgery_datetime) }}
                </span>
              </td>

              <!-- Fecha -->
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                {{ formatDate(request.request_date) }}
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

                  <button
                    v-if="request.status === 'pending' && authStore.canApproveRequests"
                    @click.stop="approveRequest(request.id)"
                    class="text-green-600 hover:text-green-900 p-1"
                    title="Aprobar solicitud"
                  >
                    <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                    </svg>
                  </button>

                  <button
                    v-if="request.status === 'pending' && authStore.canApproveRequests"
                    @click.stop="rejectRequest(request.id)"
                    class="text-red-600 hover:text-red-900 p-1"
                    title="Rechazar solicitud"
                  >
                    <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                    </svg>
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Paginación -->
      <div v-if="totalPages > 1" class="px-6 py-4 border-t border-gray-200 bg-gray-50">
        <div class="flex items-center justify-between">
          <div class="text-sm text-gray-500">
            Mostrando {{ (currentPage - 1) * pageSize + 1 }} a {{ Math.min(currentPage * pageSize, totalRequests) }} de {{ totalRequests }} solicitudes
          </div>
          <div class="flex space-x-2">
            <button
              @click="previousPage"
              :disabled="currentPage === 1"
              class="px-3 py-1 border border-gray-300 rounded text-sm text-gray-700 hover:bg-gray-100 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              Anterior
            </button>
            <button
              @click="nextPage"
              :disabled="currentPage === totalPages"
              class="px-3 py-1 border border-gray-300 rounded text-sm text-gray-700 hover:bg-gray-100 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              Siguiente
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import supplyRequestService from '../services/supplyRequestService'
import pavilionService from '../services/pavilionService'
import { format } from 'date-fns'
import { es } from 'date-fns/locale'
import Swal from 'sweetalert2'

const router = useRouter()
const authStore = useAuthStore()

// Estado reactivo
const loading = ref(false)
const error = ref(null)
const requests = ref([])
const pavilions = ref([])
const currentPage = ref(1)
const pageSize = ref(10)

// Filtros
const filters = ref({
  status: '',
  urgency: '',
  search: ''
})

// Computed properties
const filteredRequests = computed(() => {
  let filtered = requests.value

  if (filters.value.urgency) {
    filtered = filtered.filter(request => 
      supplyRequestService.calculateUrgencyFromSurgeryDate(request.surgery_datetime) === filters.value.urgency)
  }

  if (filters.value.search) {
    const search = filters.value.search.toLowerCase()
    filtered = filtered.filter(request =>
      request.request_number.toLowerCase().includes(search)
    )
  }

  return filtered
})

const paginatedRequests = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredRequests.value.slice(start, end)
})

const totalRequests = computed(() => filteredRequests.value.length)
const totalPages = computed(() => Math.ceil(totalRequests.value / pageSize.value))

// Calcular estadísticas localmente basándose en las solicitudes filtradas
const stats = computed(() => {
  const requests = filteredRequests.value
  
  return {
    pending: requests.filter(r => r.status === 'pending').length,
    approved: requests.filter(r => r.status === 'approved').length,
    rejected: requests.filter(r => r.status === 'rejected').length,
    total: requests.length
  }
})

// Métodos
const loadSupplyRequests = async () => {
  loading.value = true
  error.value = null

  try {
    const result = await supplyRequestService.getAllSupplyRequests(100, 0, filters.value.status)
    
    if (result.success && result.data) {
      let allRequests = result.data.requests || []
      
      // Filtrar solicitudes según el rol del usuario
      if (authStore.canViewAllRequests) {
        // Admin y encargado de bodega ven todas las solicitudes
        requests.value = allRequests
      } else if (authStore.canCreateRequests) {
        // Enfermera y doctor solo ven sus propias solicitudes
        const userRut = authStore.getUserRut
        requests.value = allRequests.filter(request => request.requested_by === userRut)
      } else {
        // Otros roles no deberían ver solicitudes (esto no debería pasar por la navegación)
        requests.value = []
      }
      
      console.log('Solicitudes cargadas:', requests.value)
    } else {
      requests.value = []
      error.value = result.error || 'Error al cargar solicitudes'
    }
  } catch (err) {
    console.error('Error cargando solicitudes:', err)
    error.value = 'Error al conectar con el servidor'
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

const filterRequests = () => {
  currentPage.value = 1
}

const viewRequest = (id) => {
  router.push(`/supply-requests/${id}`)
}

const approveRequest = async (id) => {
  try {
    const approvalData = {
      approved_by: 'ADMIN',
      approved_by_name: 'Sistema Admin',
      approval_notes: 'Aprobada desde interfaz web'
    }
    
    await supplyRequestService.approveSupplyRequest(id, approvalData)
    await loadSupplyRequests()
  } catch (err) {
    console.error('Error aprobando solicitud:', err)
    Swal.fire({
      icon: 'error',
      title: 'Error al aprobar la solicitud',
      text: err.response?.data?.error || err.message
    })
  }
}

const rejectRequest = async (id) => {
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

  try {
    const rejectionData = {
      rejected_by: 'ADMIN',
      rejected_by_name: 'Sistema Admin',
      rejection_reason: reason
    }
    
    await supplyRequestService.rejectSupplyRequest(id, rejectionData)
    await loadSupplyRequests()
  } catch (err) {
    console.error('Error rechazando solicitud:', err)
    Swal.fire({
      icon: 'error',
      title: 'Error al rechazar la solicitud',
      text: err.response?.data?.error || err.message
    })
  }
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

const formatSurgeryDateTime = (surgeryDateTime) => {
  return supplyRequestService.formatSurgeryDateTime(surgeryDateTime)
}

const getUrgencyLabel = (surgeryDateTime) => {
  const urgency = supplyRequestService.calculateUrgencyFromSurgeryDate(surgeryDateTime)
  const labels = {
    'critical': 'Crítica',
    'high': 'Alta',
    'normal': 'Normal',
    'low': 'Baja'
  }
  return labels[urgency] || 'Normal'
}

const getUrgencyBadgeClass = (surgeryDateTime) => {
  const color = supplyRequestService.getUrgencyColor(surgeryDateTime)
  const classes = {
    'red': 'bg-red-100 text-red-800',
    'orange': 'bg-orange-100 text-orange-800',
    'blue': 'bg-blue-100 text-blue-800',
    'gray': 'bg-gray-100 text-gray-800'
  }
  return classes[color] || classes.blue
}

// Lifecycle
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
</style>