<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="bg-white rounded-lg shadow-sm border p-6">
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h2 class="text-2xl font-semibold text-gray-900">Historial de Insumos</h2>
          <p class="text-gray-600 mt-1">Rastrea todos los movimientos y cambios de estado de los insumos</p>
        </div>
        <button 
          @click="exportToExcel" 
          :disabled="loading || filteredHistory.length === 0"
          class="btn-secondary flex items-center justify-center"
          :class="{ 'opacity-50 cursor-not-allowed': loading || filteredHistory.length === 0 }"
        >
          <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
          </svg>
          Exportar a Excel
        </button>
      </div>
    </div>

    <!-- Filtros -->
    <FilterPanel :filters="filterConfig" :result-count="filteredHistory.length" @filter-change="onFilterChange" />

    <!-- Loading -->
    <div v-if="loading" class="card">
      <div class="flex justify-center items-center py-12">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
        <span class="ml-3">Cargando historial...</span>
      </div>
    </div>

    <!-- Tabla -->
    <DataTable
      v-else
      :columns="tableColumns"
      :rows="filteredHistory"
      default-sort-key="date_time"
      default-sort-order="desc"
      empty-message="No hay registros en el historial"
      :table-actions="[{ type: 'view', label: 'Ver detalles', onClick: (row) => viewDetails(row) }]">
      <template #cell-supply_name="{ row }">
        <div class="font-medium text-gray-900">{{ row.supply_name || 'Sin nombre' }}</div>
      </template>
      <template #cell-qr_code="{ row }">
        <span class="text-gray-600">{{ row.qr_code || 'N/A' }}</span>
      </template>
      <template #cell-status="{ row }">
        <div class="flex flex-wrap items-center gap-1">
          <span class="px-2 py-1 text-xs font-semibold rounded-full" :class="getStatusClass(row.status)">
            {{ formatStatus(row.status) }}
          </span>
          <!-- Etiqueta visual complementaria; el estado real del insumo es el de arriba -->
          <span v-if="isInRequest(row)"
            class="px-2 py-1 text-xs font-medium rounded-full bg-purple-100 text-purple-800"
            title="El insumo está dentro de una solicitud">
            En solicitud
          </span>
        </div>
      </template>
      <template #cell-destination_type="{ row }">
        <div class="text-sm font-medium text-gray-900">{{ row.destination_name || formatDestinationType(row.destination_type) }}</div>
      </template>
      <template #cell-date_time="{ row }">
        {{ formatDateTime(row.date_time) }}
      </template>

    </DataTable>



    <!-- Modal de detalles -->
    <Teleport to="body">
      <div v-if="showDetailsModal" class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50" @click.self="closeDetailsModal">
        <div class="relative top-20 mx-auto p-5 border w-full max-w-2xl shadow-lg rounded-md bg-white">
        <div class="space-y-4">
          <div class="flex justify-between items-center border-b pb-3">
            <h3 class="text-xl font-semibold text-gray-900">
              Detalles del Historial #{{ selectedItem?.id }}
            </h3>
            <button @click="closeDetailsModal" class="text-gray-400 hover:text-gray-600">
              <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
          
          <div v-if="selectedItem" class="space-y-4 max-h-96 overflow-y-auto">
            <!-- Información básica -->
            <div class="bg-gray-50 rounded-lg p-4">
              <h4 class="text-sm font-semibold text-gray-700 mb-3 flex items-center">
                <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                Información Básica
              </h4>
              <div class="grid grid-cols-2 gap-3">
                <div>
                  <p class="text-xs text-gray-500">ID Historial</p>
                  <p class="text-sm font-medium text-gray-900">{{ selectedItem.id }}</p>
                </div>
                <div>
                  <p class="text-xs text-gray-500">Estado</p>
                  <p class="text-sm">
                    <span class="px-2 py-1 text-xs font-semibold rounded-full" :class="getStatusClass(selectedItem.status)">
                      {{ formatStatus(selectedItem.status) }}
                    </span>
                  </p>
                </div>
                <div class="col-span-2">
                  <p class="text-xs text-gray-500">Nombre del Insumo</p>
                  <p class="text-sm font-medium text-gray-900">{{ selectedItem.supply_name || 'Sin nombre' }}</p>
                </div>
                <div>
                  <p class="text-xs text-gray-500">Código QR</p>
                  <p class="text-sm font-medium text-gray-900 font-mono">{{ selectedItem.qr_code || 'N/A' }}</p>
                </div>
                <div>
                  <p class="text-xs text-gray-500">ID Insumo</p>
                  <p class="text-sm font-medium text-gray-900">{{ selectedItem.medical_supply_id }}</p>
                </div>
                <div>
                  <p class="text-xs text-gray-500">Usuario</p>
                  <p class="text-sm font-medium text-gray-900">{{ selectedItem.user_rut || 'SYSTEM' }}</p>
                </div>
                <div>
                  <p class="text-xs text-gray-500">Fecha y Hora</p>
                  <p class="text-sm font-medium text-gray-900">{{ formatDateTime(selectedItem.date_time) }}</p>
                </div>
              </div>
            </div>

            <!-- Destino -->
            <div class="bg-blue-50 rounded-lg p-4">
              <h4 class="text-sm font-semibold text-gray-700 mb-3 flex items-center">
                <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
                </svg>
                Destino
              </h4>
              <div class="grid grid-cols-2 gap-3">
                <div>
                  <p class="text-xs text-gray-500">Tipo</p>
                  <p class="text-sm font-medium text-gray-900">{{ formatDestinationType(selectedItem.destination_type) }}</p>
                </div>
                <div>
                  <p class="text-xs text-gray-500">ID</p>
                  <p class="text-sm font-medium text-gray-900">{{ selectedItem.destination_id }}</p>
                </div>
              </div>
            </div>

            <!-- Origen (si existe) -->
            <div v-if="selectedItem.origin_type || selectedItem.origin_id" class="bg-purple-50 rounded-lg p-4">
              <h4 class="text-sm font-semibold text-gray-700 mb-3 flex items-center">
                <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4" />
                </svg>
                Origen
              </h4>
              <div class="grid grid-cols-2 gap-3">
                <div>
                  <p class="text-xs text-gray-500">Tipo</p>
                  <p class="text-sm font-medium text-gray-900">{{ selectedItem.origin_type ? formatDestinationType(selectedItem.origin_type) : 'N/A' }}</p>
                </div>
                <div>
                  <p class="text-xs text-gray-500">ID</p>
                  <p class="text-sm font-medium text-gray-900">{{ selectedItem.origin_id || 'N/A' }}</p>
                </div>
              </div>
            </div>

            <!-- Confirmación (si existe) -->
            <div v-if="selectedItem.confirmed_by || selectedItem.confirmation_date" class="bg-green-50 rounded-lg p-4">
              <h4 class="text-sm font-semibold text-gray-700 mb-3 flex items-center">
                <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                Confirmación
              </h4>
              <div class="grid grid-cols-2 gap-3">
                <div>
                  <p class="text-xs text-gray-500">Confirmado por</p>
                  <p class="text-sm font-medium text-gray-900">{{ selectedItem.confirmed_by || 'N/A' }}</p>
                </div>
                <div>
                  <p class="text-xs text-gray-500">Fecha de confirmación</p>
                  <p class="text-sm font-medium text-gray-900">{{ selectedItem.confirmation_date ? formatDateTime(selectedItem.confirmation_date) : 'N/A' }}</p>
                </div>
              </div>
            </div>

            <!-- Notas -->
            <div v-if="selectedItem.notes" class="bg-yellow-50 rounded-lg p-4">
              <h4 class="text-sm font-semibold text-gray-700 mb-2 flex items-center">
                <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                </svg>
                Notas
              </h4>
              <p class="text-sm text-gray-700 whitespace-pre-wrap">{{ selectedItem.notes }}</p>
            </div>

            <!-- Notas de transferencia -->
            <div v-if="selectedItem.transfer_notes" class="bg-orange-50 rounded-lg p-4">
              <h4 class="text-sm font-semibold text-gray-700 mb-2 flex items-center">
                <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                Notas de Transferencia
              </h4>
              <p class="text-sm text-gray-700 whitespace-pre-wrap">{{ selectedItem.transfer_notes }}</p>
            </div>
          </div>

          <div class="flex justify-end pt-4 border-t">
            <button @click="closeDetailsModal" class="btn-secondary">Cerrar</button>
          </div>
        </div>
      </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import supplyHistoryService from '@/services/inventory/supplyHistoryService'
import { exportToExcel as exportExcel, formatDateForExcel, formatStatusForExcel } from '@/utils/excelExport'
import { format } from 'date-fns'
import { es } from 'date-fns/locale'
import { useNotification } from '@/composables/useNotification'
import FilterPanel from '@/components/common/FilterPanel.vue'
import DataTable from '@/components/common/DataTable.vue'
import { normalize } from '@/utils/normalize'
import { SUPPLY_HISTORY_STATUS_OPTIONS } from '@/config/statuses'

const history = ref([])
const loading = ref(false)
const showDetailsModal = ref(false)
const selectedItem = ref(null)
const route = useRoute()
const { success: showSuccess, error: showError } = useNotification()

const filterState = reactive({
  search: route.query.search || '',
  status: route.query.status || '',
  from_date: route.query.from_date || '',
  to_date: route.query.to_date || ''
})

const filterConfig = [
  { type: 'text', key: 'search', label: 'Buscar', placeholder: 'Nombre, QR, estado, usuario...', default: filterState.search },
  {
    type: 'select',
    key: 'status',
    label: 'Estado',
    default: filterState.status,
    options: [
      { value: '', label: 'Todos' },
      ...SUPPLY_HISTORY_STATUS_OPTIONS
    ]
  },
  { type: 'date', key: 'from_date', label: 'Desde', default: filterState.from_date },
  { type: 'date', key: 'to_date', label: 'Hasta', default: filterState.to_date }
]

const onFilterChange = (key, value) => { filterState[key] = value }

const tableColumns = [
  { key: 'supply_name', label: 'Nombre Insumo' },
  { key: 'qr_code', label: 'QR Code' },
  { key: 'status', label: 'Estado' },
  { key: 'destination_type', label: 'Destino' },
  { key: 'user_rut', label: 'Usuario' },
  { key: 'date_time', label: 'Fecha' }
]

const parseDateTimeToLocal = (dateStr) => {
  if (!dateStr) return null
  if (typeof dateStr === 'string') {
    if (!dateStr.includes('Z') && !dateStr.includes('+') && !dateStr.includes('-', 10)) {
      return new Date(dateStr.replace(' ', 'T'))
    }
  }
  return new Date(dateStr)
}

const hasActiveFilters = computed(() =>
  filterState.search !== '' || filterState.status !== '' || filterState.from_date !== '' || filterState.to_date !== ''
)

const filteredHistory = computed(() => {
  let filtered = [...history.value]

  if (filterState.status) {
    filtered = filtered.filter(item => item.status === filterState.status)
  }

  if (filterState.search) {
    const term = normalize(filterState.search)
    filtered = filtered.filter(item =>
      normalize(item.supply_name).includes(term) ||
      item.qr_code?.toLowerCase().includes(filterState.search.toLowerCase()) ||
      item.medical_supply_id?.toString().includes(filterState.search) ||
      normalize(item.status).includes(term) ||
      item.user_rut?.toLowerCase().includes(filterState.search.toLowerCase()) ||
      normalize(item.destination_type).includes(term)
    )
  }

  if (filterState.from_date) {
    const fromDate = new Date(filterState.from_date + 'T00:00:00')
    filtered = filtered.filter(item => {
      if (!item.date_time) return false
      return parseDateTimeToLocal(item.date_time) >= fromDate
    })
  }

  if (filterState.to_date) {
    const toDate = new Date(filterState.to_date + 'T23:59:59.999')
    filtered = filtered.filter(item => {
      if (!item.date_time) return false
      return parseDateTimeToLocal(item.date_time) <= toDate
    })
  }

  return filtered
})

const loadHistory = async () => {
  loading.value = true
  try {
    // Get-all una sola vez; los filtros (estado, búsqueda, fechas) son client-side.
    const data = await supplyHistoryService.getAllSupplyHistoryWithDetails()
    history.value = data
  } catch (err) {
    console.error('Error al cargar historial:', err)
    showError(err.message || 'Ocurrió un error al cargar el historial')
  } finally {
    loading.value = false
  }
}

const viewDetails = (item) => {
  selectedItem.value = item
  showDetailsModal.value = true
}

const closeDetailsModal = () => {
  showDetailsModal.value = false
  selectedItem.value = null
}

const formatDateTime = (date) => {
  if (!date) return 'N/A'
  try {
    // date_time llega como wall-clock local naive (sin zona). Se interpreta como local,
    // igual que el filtro de fechas (parseDateTimeToLocal). NO se convierte desde UTC.
    const dateObj = parseDateTimeToLocal(date)
    if (!dateObj || isNaN(dateObj.getTime())) {
      console.warn('Fecha inválida:', date)
      return 'Fecha inválida'
    }
    return format(dateObj, 'dd/MM/yyyy HH:mm:ss', { locale: es })
  } catch (error) {
    console.error('Error formateando fecha:', error, date)
    return 'Error en fecha'
  }
}

const formatStatus = (status) => {
  const statusMap = {
    'disponible': 'Disponible',
    'pendiente_retiro': 'Pendiente de retiro',
    'en_camino_a_pabellon': 'En Camino a Pabellón',
    'recepcionado': 'Recepcionado',
    'consumido': 'Consumido',
    'en_camino_a_bodega': 'En Camino a Bodega',
    'devuelto': 'Devuelto',
    'vencido': 'Vencido'
  }
  return statusMap[status] || status
}

const formatDestinationType = (type) => {
  const typeMap = {
    'store': 'Bodega',
    'pavilion': 'Pabellón'
  }
  return typeMap[type] || type
}

// Etiqueta visual complementaria: el movimiento corresponde a una reserva por solicitud.
// Se deriva de la nota (no es un estado del insumo).
const isInRequest = (row) => (row?.notes || '').toLowerCase().includes('solicitud')

const getStatusClass = (status) => {
  // Colores consistentes con el resto del front (QRDetails.vue: badge de estados de insumo
  // + mapa de historial para creado/devuelto). Creado (verde) ya no comparte color con Consumido (rojo).
  const classes = {
    'disponible': 'bg-green-100 text-green-800',
    'recepcionado': 'bg-blue-100 text-blue-800',
    'en_camino_a_pabellon': 'bg-orange-100 text-orange-800',
    'en_camino_a_bodega': 'bg-orange-100 text-orange-800',
    'pendiente_retiro': 'bg-yellow-100 text-yellow-800',
    'devuelto': 'bg-yellow-100 text-yellow-800',
    'transferido': 'bg-purple-100 text-purple-800',
    'consumido': 'bg-red-100 text-red-800',
    'vencido': 'bg-red-100 text-red-800'
  }
  return classes[status?.toLowerCase()] || 'bg-gray-100 text-gray-800'
}

const exportToExcel = async () => {
  try {
    const columns = [
      { key: 'id', label: 'ID' },
      { key: 'supply_name', label: 'Nombre del Insumo' },
      { key: 'medical_supply_id', label: 'ID Insumo' },
      { key: 'qr_code', label: 'Código QR' },
      { key: 'status', label: 'Estado', formatter: (val) => formatStatusForExcel(val) },
      { key: 'destination_type', label: 'Tipo de Destino', formatter: (val) => formatDestinationType(val) },
      { key: 'destination_id', label: 'ID Destino' },
      { key: 'destination_name', label: 'Nombre Destino' },
      { key: 'origin_type', label: 'Tipo de Origen', formatter: (val) => val ? formatDestinationType(val) : '' },
      { key: 'origin_id', label: 'ID Origen' },
      { key: 'origin_name', label: 'Nombre Origen' },
      { key: 'user_rut', label: 'Usuario RUT' },
      { key: 'date_time', label: 'Fecha y Hora', formatter: formatDateForExcel },
      { key: 'notes', label: 'Notas' }
    ]
    
    await exportExcel(filteredHistory.value, columns, 'historial_insumos')
    showSuccess('El archivo Excel se ha descargado exitosamente')
  } catch (error) {
    console.error('Error al exportar:', error)
    showError('Ocurrió un error al exportar a Excel: ' + error.message)
  }
}

onMounted(() => {
  loadHistory()
})
</script>
