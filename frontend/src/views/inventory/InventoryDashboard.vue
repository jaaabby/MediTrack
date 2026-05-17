<template>
  <div class="space-y-8">
    <!-- Header -->
    <div class="bg-white rounded-lg shadow-sm border p-6">
      <h2 class="text-2xl font-semibold text-gray-900">Dashboard de Inventario</h2>
      <p class="text-gray-600 mt-2">Inventario organizado por tipo de cirugía</p>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="flex justify-center items-center py-12">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-brand-blue-dark"></div>
      <span class="ml-3 text-gray-600">Cargando inventario...</span>
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
            <button @click="loadDashboard" class="btn-secondary mt-4 text-sm">Reintentar</button>
          </div>
        </div>
      </div>
    </div>

    <!-- Dashboard Content -->
    <div v-else class="space-y-4">
      <!-- Inventario por Tipo de Cirugía -->
      <div v-if="surgeryInventoryLoading" class="card flex justify-center items-center py-8">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-brand-blue-dark"></div>
        <span class="ml-3 text-gray-600">Cargando...</span>
      </div>

      <div v-else-if="surgeryInventoryError" class="card">
        <div class="bg-red-50 border border-red-200 rounded-md p-4">
          <div class="flex">
            <svg class="h-5 w-5 text-red-400 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <div class="ml-3">
              <h3 class="text-sm font-medium text-red-800">Error al actualizar el inventario</h3>
              <div class="mt-1 text-sm text-red-700">{{ surgeryInventoryError }}</div>
              <button @click="reloadSurgeryInventory" class="btn-secondary mt-3 text-sm">Reintentar</button>
            </div>
          </div>
        </div>
      </div>

      <template v-else>
        <FilterPanel
          :key="mainFilterKey"
          :filters="mainFilterConfig"
          :result-count="filteredSurgeryInventory.length"
          :show-clear="false"
          @filter-change="onMainFilterChange"
        >
          <template #actions>
            <button
              class="btn-secondary text-sm disabled:opacity-40 disabled:cursor-not-allowed"
              :disabled="!mainSearch"
              @click="clearMainFilter"
            >
              Limpiar filtros
            </button>
          </template>
        </FilterPanel>

        <DataTable
          :columns="tableColumns"
          :rows="filteredSurgeryInventory"
          default-sort-key="surgery_name"
          default-sort-order="asc"
          max-height="600px"
          empty-message="No hay cirugías que coincidan con el filtro"
          :table-actions="[{ type: 'view', onClick: (row) => openSurgeryModal(row) }]"
        >
          <template #cell-surgery_name="{ row }">
            {{ row.surgery_name || 'Sin tipo de cirugía' }}
          </template>
          <template #cell-total_in_store="{ row }">
            <span :class="getSurgeryStockClass(row)" class="font-semibold">{{ row.total_in_store || 0 }}</span>
            <span :class="getSurgeryStockClass(row)" class="text-xs ml-1">unidades</span>
            <span v-if="row.critical_batch_count > 0" class="block text-xs text-red-500 mt-0.5">{{ row.critical_batch_count }} lote{{ row.critical_batch_count > 1 ? 's' : '' }} crítico{{ row.critical_batch_count > 1 ? 's' : '' }}</span>
            <span v-else-if="row.warning_batch_count > 0" class="block text-xs text-orange-500 mt-0.5">{{ row.warning_batch_count }} lote{{ row.warning_batch_count > 1 ? 's' : '' }} en alerta</span>
          </template>
          <template #cell-total_transferred="{ row }">
            {{ row.total_transferred || 0 }} unidades
          </template>
          <template #cell-batch_count="{ row }">
            {{ row.batch_count || 0 }} lotes
          </template>
        </DataTable>
      </template>
    </div>
  </div>

  <!-- Modal de detalle por tipo de cirugía -->
  <div v-if="surgeryModal.show" class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50 p-4" @click.self="closeSurgeryModal">
    <div class="bg-white rounded-xl shadow-2xl w-full max-w-4xl max-h-[90vh] flex flex-col">
      <!-- Header -->
      <div class="flex items-center justify-between px-6 py-4 border-b">
        <div>
          <h2 class="text-lg font-semibold text-gray-900">{{ surgeryModal.surgeryName }}</h2>
          <p class="text-sm text-gray-500 mt-0.5">Insumos disponibles por bodega</p>
        </div>
        <button @click="closeSurgeryModal" class="text-gray-400 hover:text-gray-600 p-1 rounded-full hover:bg-gray-100">
          <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
          </svg>
        </button>
      </div>

      <!-- Body -->
      <div class="overflow-y-auto flex-1 p-6 space-y-4">
        <!-- Loading -->
        <div v-if="surgeryModal.loading" class="flex justify-center items-center py-12">
          <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-brand-blue-dark"></div>
          <span class="ml-3 text-gray-600">Cargando inventario...</span>
        </div>

        <!-- Error -->
        <div v-else-if="surgeryModal.error" class="bg-red-50 border border-red-200 rounded-md p-4 text-sm text-red-700">
          {{ surgeryModal.error }}
        </div>

        <!-- Sin datos -->
        <div v-else-if="surgeryModal.items.length === 0" class="text-center py-12 text-gray-500">
          <svg class="mx-auto h-10 w-10 text-gray-300 mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4"/>
          </svg>
          No hay insumos disponibles para este tipo de cirugía.
        </div>

        <!-- Filtro + Tabla -->
        <template v-else>
          <FilterPanel
            :key="modalFilterKey"
            :filters="modalFilterConfig"
            :result-count="filteredModalItems.length"
            :show-clear="false"
            @filter-change="onModalFilterChange"
          >
            <template #actions>
              <button
                class="btn-secondary text-sm disabled:opacity-40 disabled:cursor-not-allowed"
                :disabled="!modalSearch && !modalStoreFilter"
                @click="clearModalFilter"
              >
                Limpiar filtros
              </button>
            </template>
          </FilterPanel>

          <div v-if="filteredModalItems.length === 0" class="text-center py-8 text-gray-500 text-sm">
            No hay insumos que coincidan con el filtro.
          </div>

          <div v-else class="overflow-x-auto">
            <table class="min-w-full divide-y divide-gray-200 text-sm">
              <thead class="bg-gray-50">
                <tr>
                  <th class="px-6 py-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Insumo</th>
                  <th class="px-6 py-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Bodega</th>
                  <th class="px-6 py-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Stock</th>
                  <th class="px-6 py-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Vencimiento</th>
                </tr>
              </thead>
              <tbody class="bg-white divide-y divide-gray-100">
                <tr v-for="row in filteredModalItems" :key="row.id || (row.supply_name + row.store_name)" class="hover:bg-gray-50">
                  <td class="px-6 py-4">
                    <div class="font-medium text-gray-900">{{ row.supply_name || row.name || '—' }}</div>
                    <div v-if="row.supply_code" class="text-xs text-gray-400">Cód. {{ row.supply_code }}</div>
                  </td>
                  <td class="px-6 py-4 text-gray-700">{{ row.store_name || row.store || '—' }}</td>
                  <td class="px-6 py-4">
                    <span class="font-semibold text-sm" :class="stockColorClass(row.current_in_store ?? 0, row.critical_stock ?? 1)">{{ row.current_in_store ?? 0 }}</span>
                    <span class="text-xs ml-1" :class="stockColorClass(row.current_in_store ?? 0, row.critical_stock ?? 1)">unidades</span>
                  </td>
                  <td class="px-6 py-4">
                    <span v-if="row.expiration_date" :class="expiryColorClass(row.expiration_date)">
                      {{ formatModalDate(row.expiration_date) }}
                    </span>
                    <span v-else class="text-gray-400">—</span>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </template>
      </div>

      <!-- Footer -->
      <div class="px-6 py-3 border-t flex justify-end">
        <button @click="closeSurgeryModal" class="btn-secondary">Cerrar</button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import inventoryService from '@/services/inventory/inventoryService'
import DataTable from '@/components/common/DataTable.vue'
import FilterPanel from '@/components/common/FilterPanel.vue'
import { normalize } from '@/utils/normalize'
import { useInventoryAlerts } from '@/composables/useInventoryAlerts'

const { getExpirationClass, getStockClass } = useInventoryAlerts()

const getSurgeryStockClass = (row) => {
  if ((row.critical_batch_count ?? 0) > 0) return 'text-red-600'
  if ((row.warning_batch_count ?? 0) > 0) return 'text-orange-600'
  return 'text-gray-900'
}

const tableColumns = [
  { key: 'surgery_name', label: 'Tipo de Cirugía' },
  { key: 'total_in_store', label: 'Stock en Bodega' },
  { key: 'total_transferred', label: 'Total Transferido' },
  { key: 'batch_count', label: 'N° de Lotes', align: 'center' },
]

const loading = ref(false)
const error = ref(null)
const surgeryInventory = ref([])
const surgeryInventoryLoading = ref(false)
const surgeryInventoryError = ref(null)

// --- Filtro vista principal ---
const mainSearch = ref('')
const mainFilterKey = ref(0)

const mainFilterConfig = [
  {
    type: 'text',
    key: 'surgery_search',
    label: 'Buscar tipo de cirugía',
    placeholder: 'Filtrar por nombre...',
    default: ''
  }
]

const filteredSurgeryInventory = computed(() => {
  if (!mainSearch.value.trim()) return surgeryInventory.value
  const q = normalize(mainSearch.value)
  return surgeryInventory.value.filter(item =>
    normalize(item.surgery_name).includes(q)
  )
})

const onMainFilterChange = (_key, value) => {
  mainSearch.value = value
}

const clearMainFilter = () => {
  mainSearch.value = ''
  mainFilterKey.value += 1
}

// --- Modal de detalle por cirugía ---
const surgeryModal = ref({
  show: false,
  loading: false,
  error: null,
  surgeryName: '',
  items: []
})

// --- Filtro del modal ---
const modalSearch = ref('')
const modalStoreFilter = ref('')
const modalFilterKey = ref(0)

const modalStoreOptions = computed(() => {
  const names = surgeryModal.value.items
    .map(row => row.store_name || row.store || '')
    .filter(Boolean)
  const unique = [...new Set(names)].sort()
  return [
    { value: '', label: 'Todas las bodegas' },
    ...unique.map(name => ({ value: name, label: name }))
  ]
})

const modalFilterConfig = computed(() => [
  {
    type: 'text',
    key: 'modal_search',
    label: 'Buscar insumo',
    placeholder: 'Filtrar por nombre...',
    default: ''
  },
  {
    type: 'select',
    key: 'modal_store',
    label: 'Bodega',
    default: '',
    options: modalStoreOptions.value
  }
])

const filteredModalItems = computed(() => {
  return surgeryModal.value.items.filter(row => {
    const matchesSearch = !modalSearch.value.trim() ||
      normalize(row.supply_name || row.name || '').includes(normalize(modalSearch.value))
    const matchesStore = !modalStoreFilter.value ||
      (row.store_name || row.store || '') === modalStoreFilter.value
    return matchesSearch && matchesStore
  })
})

const onModalFilterChange = (key, value) => {
  if (key === 'modal_search') modalSearch.value = value
  if (key === 'modal_store') modalStoreFilter.value = value
}

const clearModalFilter = () => {
  modalSearch.value = ''
  modalStoreFilter.value = ''
  modalFilterKey.value += 1
}

const openSurgeryModal = async (item) => {
  modalSearch.value = ''
  modalStoreFilter.value = ''
  modalFilterKey.value += 1
  surgeryModal.value = { show: true, loading: true, error: null, surgeryName: item.surgery_name || 'Sin tipo de cirugía', items: [] }
  try {
    const data = await inventoryService.getStoreInventory({ surgery_id: item.surgery_id })
    surgeryModal.value.items = data
  } catch (err) {
    surgeryModal.value.error = err.message || 'Error al cargar el inventario'
  } finally {
    surgeryModal.value.loading = false
  }
}

const closeSurgeryModal = () => {
  surgeryModal.value.show = false
}

const formatModalDate = (dateStr) => {
  if (!dateStr) return '—'
  const d = new Date(dateStr)
  return `${String(d.getDate()).padStart(2,'0')}/${String(d.getMonth()+1).padStart(2,'0')}/${d.getFullYear()}`
}

const expiryColorClass = getExpirationClass
const stockColorClass = getStockClass

const loadDashboard = async () => {
  loading.value = true
  error.value = null

  try {
    await loadSurgeryInventory()
  } catch (err) {
    error.value = err.message || 'Error al cargar el dashboard'
    console.error('Error loading dashboard:', err)
  } finally {
    loading.value = false
  }
}

const loadSurgeryInventory = async () => {
  surgeryInventoryLoading.value = true
  try {
    const data = await inventoryService.getInventoryBySurgeryType()
    surgeryInventory.value = Array.isArray(data) ? data : []
  } catch (err) {
    console.error('Error loading surgery inventory:', err)
    surgeryInventory.value = []
    throw err
  } finally {
    surgeryInventoryLoading.value = false
  }
}

const reloadSurgeryInventory = async () => {
  surgeryInventoryError.value = null
  try {
    await loadSurgeryInventory()
  } catch (err) {
    surgeryInventoryError.value = err.message || 'Error al actualizar el inventario'
  }
}

onMounted(async () => {
  await loadDashboard()
  import('@/views/inventory/StoreInventoryView.vue')
})
</script>
