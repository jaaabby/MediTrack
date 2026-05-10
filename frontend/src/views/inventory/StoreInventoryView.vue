<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="bg-white rounded-lg shadow-sm border p-6">
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h2 class="text-2xl font-semibold text-gray-900">Inventario de Bodegas</h2>
          <p class="text-gray-600 mt-1">Stock detallado en cada bodega del sistema</p>
        </div>
      </div>
    </div>

    <!-- Filtros -->
    <FilterPanel
      :key="filterPanelKey"
      :filters="filterConfig"
      :result-count="filteredInventory.length"
      :show-clear="false"
      @filter-change="onFilterChange"
    >
      <template #filter-surgery_search="{ setValue }">
        <div class="relative">
          <label class="block text-sm font-medium text-gray-700 mb-2">Tipo de Cirugía</label>
          <input
            type="text"
            :value="surgerySearch"
            placeholder="Buscar tipo de cirugía..."
            class="form-input"
            @input="setValue($event.target.value)"
            @focus="showSurgeryOptions = true"
            @blur="hideSurgeryOptions"
            autocomplete="off"
          />

          <div
            v-if="showSurgeryOptions && filteredSurgeries.length > 0"
            class="absolute z-50 w-full mt-1 bg-white border border-gray-300 rounded-md shadow-lg max-h-60 overflow-auto"
          >
            <button
              v-for="surgery in filteredSurgeries"
              :key="surgery.id"
              @mousedown.prevent="selectSurgery(surgery)"
              class="w-full text-left px-4 py-2 hover:bg-blue-50 cursor-pointer text-sm border-b last:border-b-0"
            >
              {{ surgery.name }}
            </button>
          </div>
        </div>
      </template>

      <template #filter-supplier_search="{ setValue }">
        <div class="relative">
          <label class="block text-sm font-medium text-gray-700 mb-2">Proveedor</label>
          <input
            type="text"
            :value="supplierSearch"
            placeholder="Buscar proveedor..."
            class="form-input"
            @input="setValue($event.target.value)"
            @focus="showSupplierOptions = true"
            @blur="hideSupplierOptions"
            autocomplete="off"
          />

          <div
            v-if="showSupplierOptions && filteredSuppliers.length > 0"
            class="absolute z-50 w-full mt-1 bg-white border border-gray-300 rounded-md shadow-lg max-h-60 overflow-auto"
          >
            <button
              v-for="supplier in filteredSuppliers"
              :key="supplier"
              @mousedown.prevent="selectSupplier(supplier)"
              class="w-full text-left px-4 py-2 hover:bg-blue-50 cursor-pointer text-sm border-b last:border-b-0"
            >
              {{ supplier }}
            </button>
          </div>
        </div>
      </template>

      <template #actions>
        <button @click="clearFilters" :disabled="!hasActiveFilters" class="btn-secondary disabled:opacity-40 disabled:cursor-not-allowed">Limpiar filtros</button>
      </template>
    </FilterPanel>

    <!-- Loading -->
    <div v-if="loading" class="card">
      <div class="flex justify-center items-center py-12">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
        <span class="ml-3 text-gray-600">Cargando inventario...</span>
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

    <!-- Tabla -->
    <div v-else-if="inventory.length === 0" class="card text-center py-12">
      <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4" />
      </svg>
      <h3 class="mt-2 text-sm font-medium text-gray-900">No se encontró inventario</h3>
      <p class="mt-1 text-sm text-gray-500">
        No hay lotes que coincidan con los filtros aplicados
      </p>
    </div>

    <DataTable
      v-else
      :columns="tableColumns"
      :rows="filteredInventory"
      default-sort-key="created_at"
      default-sort-order="desc"
      :items-per-page="10"
    >









      <template #cell-store_name="{ row }">
        {{ row.store_name }}
      </template>

      <template #cell-created_at="{ row }">
        {{ formatDate(row.created_at) }}
      </template>

      <template #cell-supply_name="{ row }">
        <div>{{ row.supply_name }}</div>
        <div class="text-xs text-gray-500">Código: {{ row.supply_code }}</div>
        <div class="text-xs text-gray-500">Lote: {{ row.batch_id }}</div>
      </template>

      <template #cell-surgery_name="{ row }">
        {{ row.surgery_name || 'Sin asignar' }}
      </template>

      <template #cell-current_in_store="{ row }">
        <span :class="getStockClass(row.current_in_store, row.critical_stock)" class="font-medium">{{ row.current_in_store }}</span>
        <span :class="getStockClass(row.current_in_store, row.critical_stock)" class="text-xs ml-1">unidades</span>
      </template>

      <template #cell-original_amount="{ row }">
        {{ row.original_amount }} unidades
      </template>

      <template #cell-total_transferred_out="{ row }">
        <div>Transferidos: {{ row.total_transferred_out || 0 }}</div>
        <div class="text-xs text-gray-500">Consumidos: {{ (row.total_consumed_in_store || 0) + (row.total_consumed_from_pavilions || 0) }}</div>
        <div class="text-xs text-gray-500">Devueltos: {{ row.total_returned_in || 0 }}</div>
      </template>

      <template #cell-batch_supplier="{ row }">
        {{ row.batch_supplier }}
      </template>

      <template #cell-expiration_date="{ row }">
        <span :class="getExpirationClass(row.expiration_date, row.expiration_alert_days)">
          {{ formatDate(row.expiration_date) }}
        </span>
      </template>
    </DataTable>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import FilterPanel from '@/components/common/FilterPanel.vue'
import DataTable from '@/components/common/DataTable.vue'
import inventoryService from '@/services/inventory/inventoryService'
import surgeryService from '@/services/management/surgeryService'

const route = useRoute()
const initialStoreId = Array.isArray(route.query.store_id) ? route.query.store_id[0] : route.query.store_id
const initialSurgeryId = Array.isArray(route.query.surgery_id) ? route.query.surgery_id[0] : route.query.surgery_id

const loading = ref(false)
const error = ref(null)
const inventory = ref([])
const stores = ref([])
const surgeries = ref([])

const filters = ref({
  store_id: initialStoreId ? String(initialStoreId) : '',
  surgery_id: initialSurgeryId ? String(initialSurgeryId) : '',
  supplier: '',
  low_stock: '',
  near_expiration: ''
})
const filterPanelKey = ref(0)

// Estados para autocompletado
const surgerySearch = ref('')
const showSurgeryOptions = ref(false)
const supplierSearch = ref('')
const showSupplierOptions = ref(false)

// Paginación
// (manejada por DataTable)

const tableColumns = [
  { key: 'store_name', label: 'Bodega' },
  { key: 'created_at', label: 'F. Ingreso' },
  { key: 'supply_name', label: 'Insumo' },
  { key: 'surgery_name', label: 'Tipo Cirugía' },
  { key: 'current_in_store', label: 'Stock Actual' },
  { key: 'original_amount', label: 'Stock Original' },
  { key: 'total_transferred_out', label: 'Movimientos' },
  { key: 'batch_supplier', label: 'Proveedor' },
  { key: 'expiration_date', label: 'F. Vencimiento' },
]

// ID 6: normaliza texto eliminando tildes y convirtiendo a minúsculas
const normalizeText = (text) => {
  if (!text) return ''
  return text
    .toLowerCase()
    .normalize('NFD')
    .replace(/[\u0300-\u036f]/g, '')
}

const filterConfig = computed(() => [
  {
    type: 'select',
    key: 'store_id',
    label: 'Bodega',
    default: '',
    initialValue: filters.value.store_id,
    options: [
      { value: '', label: 'Todas las bodegas' },
      ...stores.value.map(store => ({ value: String(store.id), label: store.name }))
    ]
  },
  {
    type: 'custom',
    key: 'surgery_search',
    label: 'Tipo de Cirugía',
    default: ''
  },
  {
    type: 'custom',
    key: 'supplier_search',
    label: 'Proveedor',
    default: ''
  },
  {
    type: 'toggle',
    key: 'low_stock',
    label: 'Stock crítico',
    default: '',
    options: [
      { value: '', label: 'Todos' },
      { value: 'true', label: 'Stock Crítico', activeClass: 'bg-red-600 text-white' }
    ]
  },
  {
    type: 'toggle',
    key: 'near_expiration',
    label: 'Por vencer',
    default: '',
    options: [
      { value: '', label: 'Todos' },
      { value: 'true', label: 'Por Vencer', activeClass: 'bg-orange-500 text-white' }
    ]
  }
])

const hasActiveFilters = computed(() =>
  filters.value.store_id !== '' ||
  filters.value.surgery_id !== '' ||
  filters.value.supplier.trim() !== '' ||
  filters.value.low_stock !== '' ||
  filters.value.near_expiration !== '' ||
  surgerySearch.value.trim() !== '' ||
  supplierSearch.value.trim() !== ''
)

// Computed para proveedores únicos del inventario
const uniqueSuppliers = computed(() => {
  if (!inventory.value || inventory.value.length === 0) return []
  const suppliers = inventory.value
    .map(item => item.batch_supplier)
    .filter(supplier => supplier && supplier.trim())
  return [...new Set(suppliers)].sort()
})

// Computed para cirugías filtradas
const filteredSurgeries = computed(() => {
  if (!surgerySearch.value.trim()) {
    return surgeries.value.slice(0, 10)
  }
  const search = surgerySearch.value.toLowerCase().trim()
  return surgeries.value.filter(surgery => 
    surgery.name.toLowerCase().includes(search)
  ).slice(0, 10)
})

// Computed para proveedores filtrados
const filteredSuppliers = computed(() => {
  if (!supplierSearch.value.trim()) {
    return uniqueSuppliers.value.slice(0, 10)
  }
  const search = normalizeText(supplierSearch.value)
  return uniqueSuppliers.value.filter(supplier => 
    normalizeText(supplier).includes(search)
  ).slice(0, 10)
})

// Computed para filtrar inventario (filtros client-side de proveedor y stock bajo)
const filteredInventory = computed(() => {
  if (!inventory.value || inventory.value.length === 0) return []

  let result = inventory.value

  // Filtro client-side de stock bajo usando critical_stock (regla consolidada del frontend)
  if (filters.value.low_stock === 'true') {
    result = result.filter(item => isLowStock(item.current_in_store, item.critical_stock))
  }

  // ID 6: filtro client-side de proveedor, insensible a mayúsculas y tildes
  const supplierQuery = normalizeText(filters.value.supplier)
  if (!supplierQuery) return result
  return result.filter(item =>
    normalizeText(item.batch_supplier).includes(supplierQuery)
  )
})

let debounceTimeout = null
const debouncedApplyFilters = () => {
  if (debounceTimeout) clearTimeout(debounceTimeout)
  debounceTimeout = setTimeout(() => {
    applyFilters()
  }, 500)
}

const loadInventory = async () => {
  loading.value = true
  error.value = null
  
  try {
    // El filtro de proveedor se aplica client-side con normalización de tildes y mayúsculas (ID 6)
    // No se envía al backend para evitar el LIKE case-sensitive de PostgreSQL
    const backendFilters = { ...filters.value, supplier: '' }
    const data = await inventoryService.getStoreInventory(backendFilters)
    // Asegurarse de que inventory.value siempre sea un array
    inventory.value = Array.isArray(data) ? data : []
  } catch (err) {
    error.value = err.message || 'Error al cargar inventario de bodegas'
    console.error('Error loading store inventory:', err)
    // En caso de error, asegurar que sea un array vacío
    inventory.value = []
  } finally {
    loading.value = false
  }
}

const loadStores = async () => {
  try {
    const data = await inventoryService.getAllStores()
    stores.value = Array.isArray(data) ? data : []
  } catch (err) {
    console.error('Error loading stores:', err)
    stores.value = []
  }
}

const loadSurgeries = async () => {
  try {
    const data = await surgeryService.getAllSurgeries()
    surgeries.value = Array.isArray(data) ? data : []
  } catch (err) {
    console.error('Error loading surgeries:', err)
    surgeries.value = []
  }
}

const applyFilters = () => {
  loadInventory()
}

const onFilterChange = (key, value) => {
  switch (key) {
    case 'store_id':
      filters.value.store_id = value
      applyFilters()
      break
    case 'surgery_search':
      surgerySearch.value = value
      onSurgerySearch()
      break
    case 'supplier_search':
      supplierSearch.value = value
      onSupplierSearch()
      break
    case 'low_stock':
    case 'near_expiration':
      filters.value[key] = value
      applyFilters()
      break
    default:
      break
  }
}

const clearFilters = () => {
  filters.value = {
    store_id: '',
    surgery_id: '',
    supplier: '',
    low_stock: '',
    near_expiration: ''
  }
  surgerySearch.value = ''
  supplierSearch.value = ''
  showSurgeryOptions.value = false
  showSupplierOptions.value = false
  filterPanelKey.value += 1
  applyFilters()
}

// Funciones para manejar autocompletado de cirugía
const onSurgerySearch = () => {
  showSurgeryOptions.value = true
  // Si el texto coincide exactamente con una cirugía, seleccionarla automáticamente
  const exactMatch = surgeries.value.find(surgery => 
    surgery.name.toLowerCase() === surgerySearch.value.toLowerCase()
  )
  if (exactMatch) {
    filters.value.surgery_id = exactMatch.id.toString()
    debouncedApplyFilters()
  } else if (surgerySearch.value.trim() === '') {
    filters.value.surgery_id = ''
    debouncedApplyFilters()
  }
}

const selectSurgery = (surgery) => {
  filters.value.surgery_id = surgery.id.toString()
  surgerySearch.value = surgery.name
  showSurgeryOptions.value = false
  applyFilters()
}

const hideSurgeryOptions = () => {
  setTimeout(() => {
    showSurgeryOptions.value = false
    // Si hay una cirugía seleccionada, mostrar su nombre completo
    if (filters.value.surgery_id) {
      const selectedSurgery = surgeries.value.find(s => s.id === parseInt(filters.value.surgery_id))
      if (selectedSurgery) {
        surgerySearch.value = selectedSurgery.name
      }
    } else if (!surgerySearch.value.trim()) {
      surgerySearch.value = ''
    }
  }, 200)
}

// Funciones para manejar autocompletado de proveedor
const onSupplierSearch = () => {
  showSupplierOptions.value = true
  filters.value.supplier = supplierSearch.value
  debouncedApplyFilters()
}

const selectSupplier = (supplier) => {
  supplierSearch.value = supplier
  filters.value.supplier = supplier
  showSupplierOptions.value = false
  applyFilters()
}

const hideSupplierOptions = () => {
  setTimeout(() => {
    showSupplierOptions.value = false
  }, 200)
}

const formatDate = (dateString) => {
  if (!dateString) return 'N/A'
  try {
    return new Date(dateString).toLocaleDateString('es-CL')
  } catch {
    return dateString
  }
}

// Stock bajo = en o por debajo del umbral crítico del insumo (misma regla que Inventory.vue)
const isLowStock = (current, critical) => {
  if (critical == null) return false
  const threshold = critical || 1
  return current <= threshold
}

// Stock medio = entre crítico y 2x crítico (misma regla que Inventory.vue)
const isMediumStock = (current, critical) => {
  if (critical == null) return false
  const threshold = critical || 1
  return current > threshold && current <= threshold * 2
}

// Color del número de stock (prioridad: rojo > naranja > gris)
const getStockClass = (current, critical) => {
  if (isLowStock(current, critical)) return 'text-red-600 font-semibold'
  if (isMediumStock(current, critical)) return 'text-orange-600 font-semibold'
  return 'text-gray-900'
}

// Color de fondo de la fila según nivel de stock
const getRowClass = (current, critical) => {
  if (isLowStock(current, critical)) return 'bg-red-50 hover:bg-red-100'
  if (isMediumStock(current, critical)) return 'bg-orange-50 hover:bg-orange-100'
  return 'hover:bg-gray-50'
}


const getExpirationClass = (expirationDate, alertDays) => {
  if (!expirationDate) return 'text-gray-900'
  const today = new Date()
  const expDate = new Date(expirationDate)
  const days = Math.ceil((expDate - today) / (1000 * 60 * 60 * 24))
  const threshold = alertDays && alertDays > 0 ? alertDays : 90

  if (days < 0) return 'text-red-600 font-semibold'
  if (days <= threshold) return 'text-red-600 font-semibold'
  return 'text-gray-900'
}

onMounted(async () => {
  // Cargar datos auxiliares
  await Promise.all([loadStores(), loadSurgeries()])
  
  // Inicializar nombre de cirugía si vino desde query params
  if (filters.value.surgery_id) {
    const selectedSurgery = surgeries.value.find(s => s.id === parseInt(filters.value.surgery_id))
    if (selectedSurgery) {
      surgerySearch.value = selectedSurgery.name
    }
  }
  
  // Cargar inventario
  await loadInventory()
})
</script>
