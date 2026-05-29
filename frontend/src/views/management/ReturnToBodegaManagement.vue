<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="bg-white rounded-lg shadow-sm border p-6">
      <div class="flex items-center justify-between gap-4">
        <div class="min-w-0">
          <h1 class="text-2xl font-semibold text-gray-900">Gestión de Retornos a Bodega</h1>
          <p class="text-gray-600 mt-1">Monitoreo y gestión de insumos que deben regresar a bodega después de 8 horas laborales sin consumir</p>
        </div>
        <button @click="notifyAllPavilions" :disabled="processingReturns || criticalSupplies.length === 0" class="inline-flex items-center whitespace-nowrap flex-shrink-0 px-3 py-2 text-sm font-medium rounded-lg bg-amber-500 hover:bg-amber-600 text-white transition-colors disabled:opacity-50 disabled:cursor-not-allowed">
          <svg v-if="processingReturns" class="animate-spin h-4 w-4 mr-1.5" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          <svg v-else class="h-4 w-4 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
          </svg>
          {{ processingReturns ? 'Notificando...' : `Notificar Retornos Automáticos (${criticalSupplies.length})` }}
        </button>
      </div>
    </div>

    <!-- Stats Cards -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <div class="bg-white rounded-lg shadow-sm border p-6">
        <div class="flex items-center">
          <div class="p-2 bg-yellow-100 rounded-lg">
            <svg class="h-8 w-8 text-yellow-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600">Pendientes de Retorno</p>
            <p class="text-2xl font-bold text-gray-900">{{ suppliesForReturn.length }}</p>
          </div>
        </div>
      </div>

      <div class="bg-white rounded-lg shadow-sm border p-6">
        <div class="flex items-center">
          <div class="p-2 bg-red-100 rounded-lg">
            <svg class="h-8 w-8 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v3.75m9-.75a9 9 0 11-18 0 9 9 0 0118 0zm-9 3.75h.008v.008H12v-.008z" />
            </svg>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600">Críticos (>8 horas laborales)</p>
            <p class="text-2xl font-bold text-gray-900">{{ criticalSupplies.length }}</p>
          </div>
        </div>
      </div>

      <div class="bg-white rounded-lg shadow-sm border p-6">
        <div class="flex items-center">
          <div class="p-2 bg-green-100 rounded-lg">
            <svg class="h-8 w-8 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600">Último Proceso</p>
            <p class="text-sm font-bold text-gray-900">{{ lastProcessDate || 'Sin procesar' }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Filtros -->
    <FilterPanel
      :filters="filterConfig"
      :result-count="filteredSupplies.length"
      @filter-change="onFilterChange"
    />

    <!-- Indicador de carga -->
      <div v-if="loading" class="flex justify-center items-center py-8">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
        <span class="ml-2 text-gray-600">Cargando datos...</span>
      </div>

      <!-- Mensaje de error -->
      <div v-else-if="error" class="bg-red-50 border border-red-200 rounded-md p-4 mx-4 mb-4">
        <div class="flex">
          <div class="flex-shrink-0">
            <svg class="h-5 w-5 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
                d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <div class="ml-3">
            <h3 class="text-sm font-medium text-red-800">Error al cargar datos</h3>
            <div class="mt-2 text-sm text-red-700">{{ error }}</div>
            <div class="mt-4">
              <button @click="refreshData" class="btn-secondary text-sm">
                Reintentar
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Tabla de datos (siempre visible; el estado vacío va dentro del tbody vía slot #empty) -->
      <DataTable
        v-else
        :columns="tableColumns"
        :rows="filteredSupplies"
        default-sort-key="businessHoursElapsed"
        default-sort-order="desc"
        :table-actions="[
          {
            type: 'notify',
            label: 'Notificar',
            disabled: (row) => !!notifyingSupplies[row.qrCode],
            loading: (row) => !!notifyingSupplies[row.qrCode],
            title: (row) => notifyingSupplies[row.qrCode] ? 'Enviando...' : 'Notificar Pabellón por correo',
            onClick: (row) => notifyPavilion(row)
          },
          { type: 'view', label: 'Ver detalles', onClick: (row) => openDetailsModal(row) }
        ]"
      >
        <template #empty>
          <svg class="h-12 w-12 text-gray-400 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <p class="text-gray-600">No hay insumos pendientes de retorno</p>
          <p class="text-sm text-gray-500 mt-1">Todos los insumos están siendo utilizados correctamente</p>
        </template>

        <template #cell-name="{ row }">
          <div class="text-xs sm:text-sm font-medium text-gray-900">{{ row.name }}</div>
        </template>
        <template #cell-status="{ row }">
          <span :class="getStatusClass(row.status)" class="inline-flex px-2 py-1 text-xs font-semibold rounded-full">
            {{ getStatusText(row.status) }}
          </span>
        </template>
        <template #cell-businessHoursElapsed="{ row }">
          <span class="text-xs sm:text-sm text-gray-600">
            {{ formatBusinessHours(row.businessHoursElapsed) }}
          </span>
        </template>
        <template #cell-receivedAt="{ row }">{{ formatDate(row.receivedAt) }}</template>
        <template #cell-storeName="{ row }">{{ row.storeName }}</template>
        <template #cell-pavilionName="{ row }">
          <span v-if="row.pavilionName || row.pavilionId" class="text-sm text-gray-900">
            {{ row.pavilionName || ('Pabellón ' + row.pavilionId) }}
          </span>
          <span v-else class="text-gray-400">—</span>
        </template>
      </DataTable>
  </div>

  <!-- Modal de detalles del insumo -->
  <Teleport to="body">
    <div v-if="showDetailsModal" class="fixed inset-0 z-50 flex items-center justify-center p-4" role="dialog" aria-modal="true">
      <div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" @click="closeDetailsModal"></div>
      <div class="relative bg-white rounded-lg shadow-xl w-full max-w-2xl flex flex-col max-h-[90vh]">
        <!-- Header -->
        <div class="px-6 pt-5 pb-4 border-b border-gray-200 flex-shrink-0">
          <div class="flex items-center justify-between">
            <div class="flex items-center gap-3">
              <div class="flex-shrink-0 flex items-center justify-center h-10 w-10 rounded-full bg-blue-100">
                <svg class="h-6 w-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </div>
              <h3 class="text-lg font-medium text-gray-900">Detalle del Insumo</h3>
            </div>
            <button @click="closeDetailsModal" class="text-gray-400 hover:text-gray-600 p-1 rounded">
              <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
        </div>

        <!-- Body -->
        <div v-if="selectedSupplyDetail" class="overflow-y-auto flex-1 px-6 py-5">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <!-- Código QR -->
            <div class="flex flex-col items-center justify-center bg-gray-50 rounded-lg p-4 border">
              <qrcode-vue :value="selectedSupplyDetail.qrCode" :size="160" level="M" class="mb-3" />
              <p class="text-xs text-gray-500 font-mono break-all text-center">{{ selectedSupplyDetail.qrCode }}</p>
              <span :class="getStatusClass(selectedSupplyDetail.status)" class="inline-flex px-2 py-1 text-xs font-semibold rounded-full mt-2">
                {{ getStatusText(selectedSupplyDetail.status) }}
              </span>
            </div>

            <!-- Información -->
            <div class="space-y-3">
              <div>
                <p class="text-xs font-semibold text-gray-500 uppercase tracking-wider">Insumo</p>
                <p class="text-sm font-medium text-gray-900 mt-0.5">{{ selectedSupplyDetail.name }}</p>
              </div>
              <div>
                <p class="text-xs font-semibold text-gray-500 uppercase tracking-wider">Tiempo sin consumir</p>
                <p class="text-sm mt-0.5" :class="isSupplyCritical(selectedSupplyDetail) ? 'text-red-600 font-bold' : 'text-gray-900'">
                  {{ formatBusinessHours(selectedSupplyDetail.businessHoursElapsed) }}
                  <span v-if="isSupplyCritical(selectedSupplyDetail)" class="ml-1 text-xs bg-red-100 text-red-700 px-1.5 py-0.5 rounded-full">Crítico</span>
                </p>
              </div>
              <div>
                <p class="text-xs font-semibold text-gray-500 uppercase tracking-wider">Recepcionado</p>
                <p class="text-sm text-gray-900 mt-0.5">{{ formatDate(selectedSupplyDetail.receivedAt) }}</p>
              </div>
              <div>
                <p class="text-xs font-semibold text-gray-500 uppercase tracking-wider">Bodega destino</p>
                <p class="text-sm text-gray-900 mt-0.5">{{ selectedSupplyDetail.storeName || '—' }}</p>
              </div>
              <div>
                <p class="text-xs font-semibold text-gray-500 uppercase tracking-wider">Pabellón actual</p>
                <p class="text-sm text-gray-900 mt-0.5">{{ selectedSupplyDetail.pavilionName || (selectedSupplyDetail.pavilionId ? 'Pabellón ' + selectedSupplyDetail.pavilionId : '—') }}</p>
              </div>
              <div v-if="selectedSupplyDetail.supplier">
                <p class="text-xs font-semibold text-gray-500 uppercase tracking-wider">Proveedor</p>
                <p class="text-sm text-gray-900 mt-0.5">{{ selectedSupplyDetail.supplier }}</p>
              </div>
              <div v-if="selectedSupplyDetail.expirationDate">
                <p class="text-xs font-semibold text-gray-500 uppercase tracking-wider">Vencimiento</p>
                <p class="text-sm text-gray-900 mt-0.5">{{ formatDate(selectedSupplyDetail.expirationDate) }}</p>
              </div>
              <div v-if="selectedSupplyDetail.batchId">
                <p class="text-xs font-semibold text-gray-500 uppercase tracking-wider">Lote ID</p>
                <p class="text-sm text-gray-900 mt-0.5">#{{ selectedSupplyDetail.batchId }}</p>
              </div>
            </div>
          </div>
        </div>

        <!-- Footer -->
        <div class="bg-gray-50 px-6 py-3 border-t border-gray-200 flex justify-end flex-shrink-0">
          <button @click="closeDetailsModal" class="px-4 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500">
            Cerrar
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed, watch, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { format } from 'date-fns'
import { es } from 'date-fns/locale'
import QrcodeVue from 'qrcode.vue'
import FilterPanel from '@/components/common/FilterPanel.vue'
import DataTable from '@/components/common/DataTable.vue'
import returnToBodegaService from '@/services/management/returnToBodegaService'
import { useNotification } from '@/composables/useNotification'
import { useAlert } from '@/composables/useAlert'
import { useQRPdfDownload } from '@/composables/useQRPdfDownload'
import { normalize } from '@/utils/normalize'
import { parseDbDate } from '@/utils/dateUtils'

const { success: showSuccess, error: showError, warning: showWarning, info: showInfo } = useNotification()
const { confirm } = useAlert()
const { generateQRPdfAsBase64 } = useQRPdfDownload()

const router = useRouter()

const tableColumns = [
  { key: 'name', label: 'Insumo' },
  { key: 'qrCode', label: 'Código QR' },
  { key: 'businessHoursElapsed', label: 'Tiempo Sin Consumir' },
  { key: 'receivedAt', label: 'Recepcionado' },
  { key: 'storeName', label: 'Bodega Destino' },
  { key: 'pavilionName', label: 'Pabellón Actual', sortable: false },
]

// Estado del componente
const loading = ref(false)
const error = ref(null)
const suppliesForReturn = ref([])
const processingReturns = ref(false)
const returningSupplies = ref({})
const notifyingSupplies = ref({})
const lastProcessDate = ref(null)
const autoRefreshInterval = ref(null)
const autoRefreshEnabled = ref(true)
const refreshIntervalSeconds = 30 // Actualizar cada 30 segundos

// Estado de filtros centralizado
const filterState = reactive({ search: '', bodega: '', pavilion: '', critical: '' })

const onFilterChange = (key, value) => {
  filterState[key] = value
}

// Opciones dinámicas para el FilterPanel
const uniqueBodegas = computed(() =>
  [...new Set(suppliesForReturn.value.map(s => s.storeName).filter(Boolean))].sort()
)
const uniquePavilions = computed(() =>
  [...new Set(
    suppliesForReturn.value
      .map(s => s.pavilionName || (s.pavilionId ? 'Pabellón ' + s.pavilionId : null))
      .filter(Boolean)
  )].sort()
)

const filterConfig = computed(() => [
  {
    type: 'text',
    key: 'search',
    label: 'Buscar',
    placeholder: 'Nombre, código QR o proveedor...',
    colSpan: 2
  },
  {
    type: 'select',
    key: 'bodega',
    label: 'Bodega Destino',
    options: [
      { value: '', label: 'Todas' },
      ...uniqueBodegas.value.map(b => ({ value: b, label: b }))
    ]
  },
  {
    type: 'select',
    key: 'pavilion',
    label: 'Pabellón Actual',
    options: [
      { value: '', label: 'Todos' },
      ...uniquePavilions.value.map(p => ({ value: p, label: p }))
    ]
  },
  {
    type: 'toggle',
    key: 'critical',
    label: 'Mostrar:',
    options: [
      { value: '', label: 'Todos', activeClass: 'bg-blue-600 text-white' },
      { value: 'critical', label: 'Críticos', activeClass: 'bg-red-600 text-white' }
    ]
  }
])

// Computed properties
const criticalSupplies = computed(() => {
  return suppliesForReturn.value.filter(supply => supply.shouldReturn || (supply.businessHoursElapsed || 0) >= 8)
})

const filteredSupplies = computed(() => {
  let filtered = [...suppliesForReturn.value]

  // Filtro de texto
  if (filterState.search.trim()) {
    const q = normalize(filterState.search)
    filtered = filtered.filter(s =>
      normalize(s.name).includes(q) ||
      normalize(s.qrCode).includes(q) ||
      normalize(s.supplier).includes(q) ||
      String(s.batchId || '').includes(q)
    )
  }

  // Filtro bodega
  if (filterState.bodega) {
    filtered = filtered.filter(s => s.storeName === filterState.bodega)
  }

  // Filtro pabellón
  if (filterState.pavilion) {
    filtered = filtered.filter(s => {
      const pName = s.pavilionName || (s.pavilionId ? 'Pabellón ' + s.pavilionId : '')
      return pName === filterState.pavilion
    })
  }

  // Filtro criticidad
  if (filterState.critical === 'critical') {
    filtered = filtered.filter(s => s.shouldReturn || (s.businessHoursElapsed || 0) >= 8)
  }

  return filtered
})


// Funciones principales
const refreshData = async () => {
  loading.value = true
  error.value = null
  
  try {
    const supplies = await returnToBodegaService.getSuppliesForReturn()
    
    // Verificar si supplies es un array, si no convertirlo
    let suppliesArray = []
    if (Array.isArray(supplies)) {
      suppliesArray = supplies
    } else if (supplies && typeof supplies === 'object') {
      // Si supplies es un objeto, intentar extraer el array
      if (supplies.data && Array.isArray(supplies.data)) {
        suppliesArray = supplies.data
      } else if (supplies.supplies && Array.isArray(supplies.supplies)) {
        suppliesArray = supplies.supplies
      } else {
        suppliesArray = []
      }
    } else {
      suppliesArray = []
    }
    
    suppliesForReturn.value = suppliesArray.map(supply => returnToBodegaService.formatSupplyForUI(supply))
    lastProcessDate.value = new Date().toLocaleString()
    
  } catch (err) {
    console.error('Error refreshing data:', err)
    error.value = err.message || 'Error al cargar los datos'
    // Asegurar que suppliesForReturn sea un array vacío en caso de error
    suppliesForReturn.value = []
  } finally {
    loading.value = false
  }
}

const notifyAllPavilions = async () => {
  if (criticalSupplies.value.length === 0) return

  const confirmed = await confirm(
    `¿Enviar notificaciones de devolución a todos los pabellones?\n\nSe notificará a ${criticalSupplies.value.length} insumo(s) crítico(s) para que cada pabellón realice la devolución a bodega.`,
    'Notificar Retornos Automáticos'
  )
  if (!confirmed) return

  processingReturns.value = true
  error.value = null

  let notified = 0
  let failed = 0
  let lastError = ''

  for (const supply of criticalSupplies.value) {
    try {
      let pdfBase64 = ''
      try {
        const qrInfoForPdf = {
          qr_code: supply.qrCode,
          type: 'medical_supply',
          supply_info: {
            name: supply.name,
            store_name: supply.storeName,
            pavilion_name: supply.pavilionName || (supply.pavilionId ? ('Pabellón ' + supply.pavilionId) : null),
            SupplyCode: {
              code_supplier: supply.supplyCodeSupplier,
              name: supply.name
            },
            batch: {
              id: supply.batchId,
              supplier: supply.supplier,
              expiration_date: supply.expirationDate,
              amount: supply.batchAmount
            }
          }
        }
        pdfBase64 = await generateQRPdfAsBase64(qrInfoForPdf)
      } catch (pdfErr) {
        console.warn(`PDF no generado para ${supply.qrCode}:`, pdfErr)
      }
      await returnToBodegaService.notifyPavilionForReturn(supply.qrCode, pdfBase64 || '')
      notified++
    } catch (err) {
      console.error(`Error notificando ${supply.qrCode}:`, err)
      failed++
      lastError = err.message || lastError
    }
  }

  processingReturns.value = false

  if (failed === 0) {
    showSuccess(`Notificaciones enviadas a ${notified} pabellón(es) exitosamente`)
  } else if (notified === 0) {
    // Todas fallaron: mostrar el motivo amigable
    showError(lastError || 'No se pudieron enviar las notificaciones por correo.')
  } else {
    showError(`Se notificaron ${notified} pabellón(es), pero ${failed} no se pudieron enviar. ${lastError}`)
  }

  await refreshData()
}

const notifyPavilion = async (supply) => {
  if (notifyingSupplies.value[supply.qrCode]) return

  const pavilionLabel = supply.pavilionName || `Pabellón ID ${supply.pavilionId}` || 'el pabellón'
  const confirmed = await confirm(
    `¿Enviar notificación de devolución al ${pavilionLabel}?\n\nInsumo: ${supply.name}\nQR: ${supply.qrCode}\n\nSe enviará un correo con el PDF del insumo a los usuarios del pabellón para que realicen la devolución.`,
    'Notificar Pabellón'
  )
  if (!confirmed) return

  notifyingSupplies.value[supply.qrCode] = true

  try {
    // Generar el PDF del insumo para adjuntarlo al correo
    let pdfBase64 = ''
    try {
      // Construir el objeto qrInfo mínimo para la generación del PDF
      const qrInfoForPdf = {
        qr_code: supply.qrCode,
        type: 'medical_supply',
        supply_info: {
          name: supply.name,
          store_name: supply.storeName,
          pavilion_name: supply.pavilionName || (supply.pavilionId ? ('Pabellón ' + supply.pavilionId) : null),
          SupplyCode: {
            code_supplier: supply.supplyCodeSupplier,
            name: supply.name
          },
          batch: {
            id: supply.batchId,
            supplier: supply.supplier,
            expiration_date: supply.expirationDate,
            amount: supply.batchAmount
          }
        }
      }
      pdfBase64 = await generateQRPdfAsBase64(qrInfoForPdf)
    } catch (pdfErr) {
      console.warn('No se pudo generar el PDF, se enviará el correo sin adjunto:', pdfErr)
    }

    const result = await returnToBodegaService.notifyPavilionForReturn(supply.qrCode, pdfBase64 || '')
    const emailCount = result.data?.email_count || 0
    showSuccess(
      result.message || `Notificación enviada a ${emailCount} usuario(s) del pabellón`
    )
  } catch (err) {
    console.error('Error notificando al pabellón:', err)
    showError(err.message || 'Error al enviar la notificación al pabellón')
  } finally {
    delete notifyingSupplies.value[supply.qrCode]
  }
}

const returnIndividualSupply = async (supply) => {
  if (returningSupplies.value[supply.qrCode]) return
  
  const confirmed = await confirm(
    `¿Regresar ${supply.name} a bodega?\n\nQR: ${supply.qrCode}\nTiempo sin moverse: ${formatBusinessHours(supply.businessHoursElapsed || 0)}`,
    'Confirmar retorno'
  )
  if (!confirmed) {
    return
  }
  
  returningSupplies.value[supply.qrCode] = true
  
  try {
    await returnToBodegaService.returnSupplyToStore(
      supply.qrCode, 
      `Retorno manual - ${formatBusinessHours(supply.businessHoursElapsed || 0)} sin moverse`
    )
    
    showSuccess('Insumo regresado a bodega exitosamente')
    
    // Recargar datos
    await refreshData()
    
  } catch (err) {
    console.error('Error returning supply:', err)
    const errorMessage = err.message || 'Error al regresar insumo a bodega'
    error.value = errorMessage
    // ID 20-22: Notificar el error al usuario
    showError(errorMessage)
  } finally {
    delete returningSupplies.value[supply.qrCode]
  }
}

const showDetailsModal = ref(false)
const selectedSupplyDetail = ref(null)

const openDetailsModal = (supply) => {
  selectedSupplyDetail.value = supply
  showDetailsModal.value = true
}

const closeDetailsModal = () => {
  showDetailsModal.value = false
  selectedSupplyDetail.value = null
}

// ID 17/18: Helper centralizado para determinar si un insumo es crítico
// Un insumo es crítico si el backend lo marca como tal O si lleva >= 8 horas laborales
const isSupplyCritical = (supply) => {
  return supply.shouldReturn || (supply.businessHoursElapsed || 0) >= 8
}

// Funciones auxiliares
const getStatusClass = (status) => {
  switch (status) {
    case 'recepcionado':
      return 'bg-blue-100 text-blue-800'
    case 'disponible':
      return 'bg-green-100 text-green-800'
    case 'en_camino_a_pabellon':
      return 'bg-yellow-100 text-yellow-800'
    case 'consumido':
      return 'bg-gray-100 text-gray-800'
    default:
      return 'bg-gray-100 text-gray-800'
  }
}

const getStatusText = (status) => {
  switch (status) {
    case 'recepcionado':
      return 'Recepcionado'
    case 'disponible':
      return 'Disponible'
    case 'en_camino_a_pabellon':
      return 'En Camino'
    case 'consumido':
      return 'Consumido'
    default:
      return status || 'Desconocido'
  }
}

const formatDate = (dateString) => {
  if (!dateString) return 'No disponible'
  try {
    return format(parseDbDate(dateString), 'dd/MM/yyyy HH:mm', { locale: es })
  } catch (error) {
    return dateString
  }
}

const formatBusinessHours = (hours) => {
  if (!hours && hours !== 0) return '0 horas laborales'
  
  // Si es menos de 1 hora, mostrar en minutos
  if (hours < 1) {
    const minutes = Math.round(hours * 60)
    return `${minutes} min${minutes !== 1 ? 's' : ''} laborales`
  }
  
  // Si es menos de 8 horas, mostrar en horas
  if (hours < 8) {
    const roundedHours = Math.round(hours * 10) / 10
    return `${roundedHours} hora${roundedHours !== 1 ? 's' : ''} laborales`
  }
  
  // Más de 8 horas, mostrar en días aproximados
  const days = Math.floor(hours / 8)
  const remainingHours = hours % 8
  if (remainingHours === 0) {
    return `${days} día${days !== 1 ? 's' : ''} laborales`
  }
  return `${days} día${days !== 1 ? 's' : ''} y ${Math.round(remainingHours)} hora${Math.round(remainingHours) !== 1 ? 's' : ''} laborales`
}

const clearError = () => {
  error.value = null
}

// Funciones para actualización automática
const startAutoRefresh = () => {
  if (autoRefreshInterval.value) {
    clearInterval(autoRefreshInterval.value)
  }
  
  autoRefreshInterval.value = setInterval(() => {
    if (autoRefreshEnabled.value && !loading.value && !processingReturns.value) {
      refreshData()
    }
  }, refreshIntervalSeconds * 1000)
}

const stopAutoRefresh = () => {
  if (autoRefreshInterval.value) {
    clearInterval(autoRefreshInterval.value)
    autoRefreshInterval.value = null
  }
}

const toggleAutoRefresh = () => {
  autoRefreshEnabled.value = !autoRefreshEnabled.value
  if (autoRefreshEnabled.value) {
    startAutoRefresh()
  } else {
    stopAutoRefresh()
  }
}

// Lifecycle
onMounted(() => {
  refreshData()
  startAutoRefresh()
})

onUnmounted(() => {
  stopAutoRefresh()
})
</script>

<style scoped>
/* Estilos similares a Inventory.vue */

.card {
  @apply bg-white rounded-lg shadow-sm border p-4 sm:p-6;
}

.card-header {
  @apply pb-4 mb-4 border-b border-gray-200;
}

.card-title {
  @apply text-lg sm:text-xl font-semibold text-gray-900;
}

.table-container {
  overflow-x: auto;
  overflow-y: auto;
  max-height: 62vh;
  border: 1px solid #e5e7eb;
  border-radius: 0.5rem;
  box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1);
  -webkit-overflow-scrolling: touch;
  position: relative;
  max-width: 100%;
  width: 100%;
}

.table-container table {
  width: 100%;
  min-width: 900px;
  table-layout: auto;
}

.table-container::-webkit-scrollbar {
  height: 6px;
}

@media (min-width: 768px) {
  .table-container::-webkit-scrollbar {
    height: 8px;
  }
}

.table-container::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 4px;
}

.table-container::-webkit-scrollbar-thumb {
  background: #c1c1c1;
  border-radius: 4px;
}

.table-container::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8;
}

.table-header {
  background-color: #f9fafb;
}

.table-header th {
  padding: 0.75rem 1.5rem;
  text-align: left;
  font-size: 0.75rem;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: #6b7280;
  white-space: nowrap;
}

.table-header-cell {
  @apply px-3 sm:px-6 py-2 sm:py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider;
}

.table-header-cell:first-child {
  min-width: 180px;
}

.table-header-cell:nth-child(2) {
  min-width: 120px;
}

.table-header-cell:nth-child(3) {
  min-width: 150px;
}

.table-header-cell:nth-child(4) {
  min-width: 150px;
}

.table-header-cell:nth-child(5) {
  min-width: 150px;
}

.table-header-cell:last-child {
  min-width: 200px;
}

.btn-danger {
  @apply inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-red-600 hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 disabled:opacity-50 disabled:cursor-not-allowed;
}

.btn-danger-sm {
  @apply inline-flex items-center justify-center px-2 sm:px-3 py-1.5 sm:py-2 border border-transparent text-xs font-medium rounded text-white bg-red-600 hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 disabled:opacity-50 disabled:cursor-not-allowed;
}

.btn-secondary-sm {
  @apply inline-flex items-center justify-center px-2 sm:px-3 py-1.5 sm:py-2 border border-gray-300 text-xs font-medium rounded text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500;
}

/* Responsive para móviles */
@media (max-width: 640px) {
  .card {
    padding: 0.875rem;
    margin-left: -0.5rem;
    margin-right: -0.5rem;
    border-radius: 0.5rem;
  }

  .card-header {
    padding: 0.875rem;
  }

  .table-container {
    margin: 0 -0.875rem;
    border-radius: 0;
    border-left: none;
    border-right: none;
    overflow-x: auto;
  }
  
  .table-container table {
    min-width: 800px;
  }

  .table-header th,
  .table-body td {
    padding: 0.5rem 0.75rem;
  }

  button,
  .btn-primary,
  .btn-secondary {
    min-height: 44px;
  }
}

/* Responsive para tablets */
@media (max-width: 768px) {
  .card {
    padding: 1.5rem;
  }
}
</style>