<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="bg-white rounded-lg shadow-sm border p-6">
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <div class="flex items-center gap-3">
            <h2 class="text-2xl font-semibold text-gray-900">Inventario de Pabellones</h2>
          </div>
          <p class="text-gray-600 mt-1">Stock disponible en cada pabellón del hospital</p>
        </div>
      </div>
    </div>

    <!-- Resumen del Pabellón (solo cuando hay datos) -->
    <div v-if="filterState.pavilion && !loading && !error && inventory.length > 0" class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <div class="card">
        <div class="flex items-center">
          <div class="flex-shrink-0 h-12 w-12 rounded-lg bg-green-100 flex items-center justify-center">
            <svg class="h-6 w-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
          </div>
          <div class="ml-4">
            <p class="text-sm text-gray-600">Stock Disponible</p>
            <p class="text-2xl font-bold text-gray-900">{{ getTotalAvailable() }}</p>
            <p class="text-xs text-gray-500">unidades</p>
          </div>
        </div>
      </div>
      <div class="card">
        <div class="flex items-center">
          <div class="flex-shrink-0 h-12 w-12 rounded-lg bg-blue-100 flex items-center justify-center">
            <svg class="h-6 w-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4" />
            </svg>
          </div>
          <div class="ml-4">
            <p class="text-sm text-gray-600">Total Recibido</p>
            <p class="text-2xl font-bold text-gray-900">{{ getTotalReceived() }}</p>
            <p class="text-xs text-gray-500">unidades</p>
          </div>
        </div>
      </div>
      <div class="card">
        <div class="flex items-center">
          <div class="flex-shrink-0 h-12 w-12 rounded-lg bg-purple-100 flex items-center justify-center">
            <svg class="h-6 w-6 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
            </svg>
          </div>
          <div class="ml-4">
            <p class="text-sm text-gray-600">Total Consumido</p>
            <p class="text-2xl font-bold text-gray-900">{{ getTotalConsumed() }}</p>
            <p class="text-xs text-gray-500">unidades</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Selector de Pabellón y Filtros -->
    <FilterPanel :filters="filterConfig" :result-count="filterState.pavilion ? filteredInventory.length : null" @filter-change="onFilterChange" />

    <!-- Sin Selección -->
    <div v-if="!filterState.pavilion" class="card text-center py-12">
      <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
      </svg>
      <h3 class="mt-2 text-sm font-medium text-gray-900">Seleccione un pabellón</h3>
      <p class="mt-1 text-sm text-gray-500">
        Seleccione un pabellón del menú desplegable para ver su inventario
      </p>
    </div>

    <!-- Loading -->
    <div v-else-if="loading" class="card">
      <div class="flex justify-center items-center py-12">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
        <span class="ml-3 text-gray-600">Cargando inventario del pabellón...</span>
      </div>
    </div>

    <!-- Error -->
    <div v-else-if="error" class="card">
      <div class="bg-red-50 border border-red-200 rounded-md p-4">
        <div class="flex">
          <svg class="h-5 w-5 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <div class="ml-3">
            <h3 class="text-sm font-medium text-red-800">Error al cargar inventario</h3>
            <div class="mt-2 text-sm text-red-700">{{ error }}</div>
            <button @click="loadInventory" class="btn-secondary mt-4 text-sm">Reintentar</button>
          </div>
        </div>
      </div>
    </div>

    <!-- Sin Inventario -->
    <div v-else-if="inventory.length === 0" class="card text-center py-12">
      <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4" />
      </svg>
      <h3 class="mt-2 text-sm font-medium text-gray-900">No hay inventario en este pabellón</h3>
      <p class="mt-1 text-sm text-gray-500">
        Este pabellón no tiene insumos registrados actualmente
      </p>
    </div>

    <!-- Tabla de Inventario -->
    <div v-else>
      <!-- Tabla -->
      <DataTable
        :columns="tableColumns"
        :rows="filteredInventory"
        default-sort-key="supply_name"
        max-height="600px"
        empty-message="Este pabellón no tiene insumos registrados actualmente"
        :row-class="(row) => row.in_transit ? 'bg-blue-50 hover:bg-blue-100' : 'hover:bg-gray-50'"
        :table-actions="[{ type: 'view', label: 'Ver detalles', onClick: (row) => openDetailModal(row) }]"
      >
        <!-- Columna Insumo -->
        <template #cell-supply_name="{ row }">
          <div>{{ row.supply_name }}</div>
          <div class="text-xs text-gray-500">{{ row.qr_code || '—' }}</div>
        </template>

        <!-- Columna Disponible -->
        <template #cell-current_available="{ row }">
          <div>{{ row.current_available }} unidades</div>
        </template>

        <!-- Columna Total Recibido -->
        <template #cell-total_received="{ row }">
          {{ row.total_received }} unidades
        </template>

        <!-- Columna Consumido -->
        <template #cell-total_consumed="{ row }">
          {{ row.total_consumed || 0 }} unidades
        </template>

        <!-- Columna Devuelto -->
        <template #cell-total_returned="{ row }">
          {{ row.total_returned || 0 }} unidades
        </template>

        <!-- Columna F. Vencimiento -->
        <template #cell-expiration_date="{ row }">
          <div :class="getExpirationClass(row.expiration_date)">{{ formatDate(row.expiration_date) }}</div>
        </template>
      </DataTable>
    </div>
  </div>

  <!-- Modal de detalles -->
  <Teleport to="body">
    <div v-if="detailModal.show" class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50 p-4" @click.self="closeDetailModal">
      <div class="bg-white rounded-xl shadow-2xl w-full max-w-lg">
        <!-- Header -->
        <div class="flex items-center justify-between px-6 py-4 border-b">
          <div>
            <h2 class="text-lg font-semibold text-gray-900">{{ detailModal.row?.supply_name }}</h2>
            <p class="text-sm text-gray-500 mt-0.5">Código: {{ detailModal.row?.supply_code }}</p>
          </div>
          <button @click="closeDetailModal" class="text-gray-400 hover:text-gray-600 p-1 rounded-full hover:bg-gray-100">
            <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <!-- Body -->
        <div class="p-6 space-y-6">
          <!-- QR Image -->
          <div class="flex flex-col items-center">
            <div v-if="detailModal.qrLoading" class="h-40 flex items-center justify-center">
              <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-brand-blue-dark"></div>
            </div>
            <img v-else-if="detailModal.qrImage" :src="detailModal.qrImage" alt="Código QR" class="h-40 w-40" />
            <div v-else class="h-40 flex items-center justify-center text-gray-400 text-sm">Sin código QR</div>
            <p class="mt-2 text-sm text-gray-600 font-mono">{{ detailModal.row?.qr_code || '—' }}</p>
          </div>

          <!-- Detalles -->
          <div class="grid grid-cols-2 gap-4">
            <div class="bg-gray-50 rounded-lg p-3">
              <p class="text-xs text-gray-500">Disponible</p>
              <p class="text-lg font-semibold text-gray-900">{{ detailModal.row?.current_available }} unidades</p>
            </div>
            <div class="bg-gray-50 rounded-lg p-3">
              <p class="text-xs text-gray-500">Total Recibido</p>
              <p class="text-lg font-semibold text-gray-900">{{ detailModal.row?.total_received }} unidades</p>
            </div>
            <div class="bg-gray-50 rounded-lg p-3">
              <p class="text-xs text-gray-500">Consumido</p>
              <p class="text-lg font-semibold text-gray-900">{{ detailModal.row?.total_consumed || 0 }} unidades</p>
            </div>
            <div class="bg-gray-50 rounded-lg p-3">
              <p class="text-xs text-gray-500">Devuelto</p>
              <p class="text-lg font-semibold text-gray-900">{{ detailModal.row?.total_returned || 0 }} unidades</p>
            </div>
            <div class="bg-gray-50 rounded-lg p-3">
              <p class="text-xs text-gray-500">Proveedor</p>
              <p class="text-sm font-medium text-gray-900">{{ detailModal.row?.batch_supplier || '—' }}</p>
            </div>
            <div class="bg-gray-50 rounded-lg p-3">
              <p class="text-xs text-gray-500">Vencimiento</p>
              <p class="text-sm font-medium" :class="getExpirationClass(detailModal.row?.expiration_date)">
                {{ formatDate(detailModal.row?.expiration_date) }}
              </p>
            </div>
          </div>
        </div>

        <!-- Footer -->
        <div class="px-6 py-3 border-t flex justify-end">
          <button @click="closeDetailModal" class="btn-secondary">Cerrar</button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import QRCode from 'qrcode'
import inventoryService from '@/services/inventory/inventoryService'
import FilterPanel from '@/components/common/FilterPanel.vue'
import DataTable from '@/components/common/DataTable.vue'
import { useInventoryAlerts } from '@/composables/useInventoryAlerts'
import { parseDbDate } from '@/utils/dateUtils'

const { isExpired, isNearExpiration, getExpirationClass } = useInventoryAlerts()

const loading = ref(false)
const error = ref(null)
const inventory = ref([])
const pavilions = ref([])

const filterState = reactive({ pavilion: '', supplier: '', transit: 'false' })

const tableColumns = [
  { key: 'supply_name', label: 'Insumo', wrap: true },
  { key: 'current_available', label: 'Disponible' },
  { key: 'total_received', label: 'Total Recibido' },
  { key: 'total_consumed', label: 'Consumido' },
  { key: 'total_returned', label: 'Devuelto' },
  { key: 'batch_supplier', label: 'Proveedor' },
  { key: 'expiration_date', label: 'F. Vencimiento' }
]

const filterConfig = computed(() => [
  {
    type: 'select', key: 'pavilion', label: 'Seleccionar Pabellón',
    options: [
      { value: '', label: 'Seleccione un pabellón...' },
      ...pavilions.value.map(p => ({ value: String(p.id), label: p.name }))
    ]
  },
  { type: 'text', key: 'supplier', label: 'Proveedor', placeholder: 'Buscar proveedor...' },
  {
    type: 'toggle', key: 'transit', label: 'En tránsito', default: 'false',
    options: [
      { value: 'false', label: 'Todos' },
      { value: 'true', label: 'En tránsito', activeClass: 'bg-blue-600 text-white' }
    ]
  }
])

const onFilterChange = (key, value) => {
  filterState[key] = value
  if (key === 'pavilion') {
    loadInventory()
  }
}

// Normaliza texto para comparación case-insensitive y sin tildes
const normalizeText = (text) => {
  if (!text) return ''
  return text.toString().toLowerCase().normalize('NFD').replace(/[\u0300-\u036f]/g, '')
}

// Solo filtra (transit + proveedor); DataTable se encarga del ordenamiento y paginación
const filteredInventory = computed(() => {
  if (!inventory.value || inventory.value.length === 0) return []
  let result = [...inventory.value]
  if (filterState.transit === 'true') {
    result = result.filter(item => item.in_transit)
  }
  if (filterState.supplier.trim()) {
    const needle = normalizeText(filterState.supplier)
    result = result.filter(item => normalizeText(item.batch_supplier || '').includes(needle))
  }
  return result
})

const loadPavilions = async () => {
  try {
    const data = await inventoryService.getAllPavilions()
    pavilions.value = Array.isArray(data) ? data : []
  } catch (err) {
    console.error('Error loading pavilions:', err)
    pavilions.value = []
  }
}

const loadInventory = async () => {
  if (!filterState.pavilion) {
    inventory.value = []
    return
  }

  loading.value = true
  error.value = null

  try {
    const data = await inventoryService.getPavilionInventory(
      filterState.pavilion,
      true,
      null
    )
    inventory.value = Array.isArray(data) ? data : []
  } catch (err) {
    error.value = err.message || 'Error al cargar inventario del pabellón'
    console.error('Error loading pavilion inventory:', err)
    inventory.value = []
  } finally {
    loading.value = false
  }
}

const getTotalAvailable = () => {
  return inventory.value.reduce((sum, item) => sum + (item.current_available || 0), 0)
}

const getTotalReceived = () => {
  return inventory.value.reduce((sum, item) => sum + (item.total_received || 0), 0)
}

const getTotalConsumed = () => {
  return inventory.value.reduce((sum, item) => sum + (item.total_consumed || 0), 0)
}

const formatDate = (dateString) => {
  if (!dateString) return 'N/A'
  try {
    return parseDbDate(dateString).toLocaleDateString('es-CL')
  } catch {
    return dateString
  }
}


const detailModal = reactive({ show: false, row: null, qrImage: null, qrLoading: false })

const openDetailModal = async (row) => {
  detailModal.row = row
  detailModal.qrImage = null
  detailModal.show = true
  if (row.qr_code) {
    detailModal.qrLoading = true
    try {
      detailModal.qrImage = await QRCode.toDataURL(row.qr_code, { width: 200, margin: 2 })
    } catch {
      detailModal.qrImage = null
    } finally {
      detailModal.qrLoading = false
    }
  }
}

const closeDetailModal = () => {
  detailModal.show = false
  detailModal.row = null
  detailModal.qrImage = null
}

onMounted(async () => {
  await loadPavilions()
})
</script>

