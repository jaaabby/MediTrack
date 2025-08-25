<template>
  <div class="space-y-6">
    <!-- Toast Notification -->
    <transition name="fade">
      <div v-if="toast.show" class="fixed top-6 right-6 z-50 bg-blue-600 text-white px-4 py-2 rounded shadow-lg flex items-center">
        <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01"/></svg>
        <span>{{ toast.message }}</span>
      </div>
    </transition>
    <!-- Breadcrumb y acciones -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between">
      <nav class="flex mb-4 sm:mb-0" aria-label="Breadcrumb">
        <ol class="inline-flex items-center space-x-1 md:space-x-3">
          <li class="inline-flex items-center">
            <router-link to="/qr" class="text-gray-700 hover:text-blue-600 transition-colors">
              <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h2M4 4h5l2 3h3l2-3h5v5M4 4v5m16-5v5" />
              </svg>
              Escáner QR
            </router-link>
          </li>
          <li>
            <div class="flex items-center">
              <svg class="w-6 h-6 text-gray-400" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clip-rule="evenodd"></path>
              </svg>
              <span class="ml-1 text-gray-500 md:ml-2">Detalles</span>
            </div>
          </li>
        </ol>
      </nav>
      
      <div class="flex flex-wrap gap-2">
        <button @click="refreshData" class="btn-secondary text-sm" :disabled="loading" title="Actualizar información">
          <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
          Actualizar
        </button>
  <button @click="exportData" class="btn-primary text-sm" title="Exportar datos en JSON">
          <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
          </svg>
          Exportar
        </button>
  <button @click="printQRCode" class="btn-secondary text-sm" title="Imprimir etiqueta QR">
          <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 17h2a2 2 0 002-2v-4a2 2 0 00-2-2H5a2 2 0 00-2 2v4a2 2 0 002 2h2m2 4h6a2 2 0 002-2v-4a2 2 0 00-2-2H9a2 2 0 00-2 2v4a2 2 0 002 2zm8-12V5a2 2 0 00-2-2H9a2 2 0 00-2 2v4h10z" />
          </svg>
          Imprimir
        </button>
      </div>
    </div>

    <!-- Loading state -->
    <div v-if="loading" class="flex justify-center items-center py-12">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      <span class="ml-3 text-gray-600">Cargando información...</span>
    </div>

    <!-- Error state -->
    <div v-else-if="error" class="bg-red-50 border border-red-200 rounded-lg p-6">
      <div class="flex items-start">
        <div class="flex-shrink-0">
          <svg class="h-6 w-6 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </div>
        <div class="ml-3 flex-1">
          <h3 class="text-lg font-medium text-red-800">Error al cargar información</h3>
          <div class="mt-2 text-sm text-red-700">
            <p>{{ error }}</p>
          </div>
          <div class="mt-4 flex space-x-3">
            <button @click="refreshData" class="btn-secondary text-sm">
              Intentar de Nuevo
            </button>
            <router-link to="/qr" class="btn-secondary text-sm">
              Volver al Escáner
            </router-link>
          </div>
        </div>
      </div>
    </div>

    <!-- Content -->
    <div v-else-if="qrInfo" class="space-y-6">
      
      <!-- QR Information Display -->
      <QRInfoDisplay 
        :qr-info="qrInfo" 
        @view-details="handleViewDetails" 
        @consume-supply="consumeSupply"
        @view-batch="viewBatch"
      />

      <!-- Additional Actions Section -->
      <div class="bg-white rounded-lg shadow-sm border p-6">
        <h3 class="text-lg font-medium text-gray-900 mb-4 flex items-center">
          <svg class="h-5 w-5 mr-2 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
          </svg>
          Acciones Disponibles
        </h3>

        <div class="grid sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
          <!-- Ver en Inventario -->
          <button @click="viewInInventory" class="action-card" title="Ver en Inventario">
            <div class="action-icon bg-blue-100 text-blue-600">
              <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
              </svg>
            </div>
            <div class="action-content">
              <div class="font-medium">Ver en Inventario</div>
              <div class="text-sm text-gray-500">Buscar en lista completa</div>
            </div>
          </button>

          <!-- Descargar QR -->
          <button @click="downloadQR" class="action-card" title="Descargar QR en alta calidad">
            <div class="action-icon bg-green-100 text-green-600">
              <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
            </div>
            <div class="action-content">
              <div class="font-medium">Descargar QR</div>
              <div class="text-sm text-gray-500">Imagen en alta calidad</div>
            </div>
          </button>

          <!-- Imprimir QR -->
          <button @click="printQRCode" class="action-card" title="Imprimir QR físico">
            <div class="action-icon bg-purple-100 text-purple-600">
              <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 17h2a2 2 0 002-2v-4a2 2 0 00-2-2H5a2 2 0 00-2 2v4a2 2 0 002 2h2m2 4h6a2 2 0 002-2v-4a2 2 0 00-2-2H9a2 2 0 00-2 2v4a2 2 0 002 2zm8-12V5a2 2 0 00-2-2H9a2 2 0 00-2 2v4h10z" />
              </svg>
            </div>
            <div class="action-content">
              <div class="font-medium">Imprimir QR</div>
              <div class="text-sm text-gray-500">Etiqueta física</div>
            </div>
          </button>

          <!-- Compartir -->
          <button @click="shareQR" class="action-card" title="Compartir por email/chat">
            <div class="action-icon bg-yellow-100 text-yellow-600">
              <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.684 13.342C8.886 12.938 9 12.482 9 12c0-.482-.114-.938-.316-1.342m0 2.684a3 3 0 110-2.684m0 2.684l6.632 3.316m-6.632-6l6.632-3.316m0 0a3 3 0 105.367-2.684 3 3 0 00-5.367 2.684zm0 9.316a3 3 0 105.367 2.684 3 3 0 00-5.367-2.684z" />
              </svg>
            </div>
            <div class="action-content">
              <div class="font-medium">Compartir</div>
              <div class="text-sm text-gray-500">Enviar por email/chat</div>
            </div>
          </button>

          <!-- Ver Historial -->
          <button @click="openHistoryModal" class="action-card" title="Ver historial de movimientos">
            <div class="action-icon bg-indigo-100 text-indigo-600">
              <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <div class="action-content">
              <div class="font-medium">Ver Historial</div>
              <div class="text-sm text-gray-500">Movimientos completos</div>
            </div>
          </button>

          <!-- Registrar Movimiento -->
          <button 
            v-if="qrInfo.type === 'medical_supply' && !qrInfo.is_consumed" 
            @click="addMovement" 
            class="action-card"
            title="Registrar consumo de producto"
          >
            <div class="action-icon bg-red-100 text-red-600">
              <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
              </svg>
            </div>
            <div class="action-content">
              <div class="font-medium">Registrar Consumo</div>
              <div class="text-sm text-gray-500">Marcar como usado</div>
            </div>
          </button>

          <!-- Sincronizar -->
          <button v-if="qrInfo.type === 'batch'" @click="syncBatch" class="action-card" title="Sincronizar cantidades">
            <div class="action-icon bg-orange-100 text-orange-600">
              <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
              </svg>
            </div>
            <div class="action-content">
              <div class="font-medium">Sincronizar</div>
              <div class="text-sm text-gray-500">Actualizar cantidades</div>
            </div>
          </button>

          <!-- Generar Reporte -->
          <button @click="generateReport" class="action-card" title="Generar reporte detallado">
            <div class="action-icon bg-teal-100 text-teal-600">
              <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 17v-2m3 2v-4m3 4v-6m2 10H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
            </div>
            <div class="action-content">
              <div class="font-medium">Generar Reporte</div>
              <div class="text-sm text-gray-500">Análisis detallado</div>
            </div>
          </button>
        </div>
      </div>

      <!-- Modal Historial -->
      <transition name="fade">
        <div v-if="showHistoryModal" class="fixed inset-0 z-40 flex items-center justify-center bg-black bg-opacity-40">
          <div class="bg-white rounded-lg shadow-lg p-6 w-full max-w-lg relative">
            <button @click="closeHistoryModal" class="absolute top-2 right-2 text-gray-500 hover:text-gray-700" title="Cerrar">
              <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/></svg>
            </button>
            <h2 class="text-xl font-bold mb-4">Historial de Movimientos</h2>
            <div v-if="historyLoading" class="text-center py-4">Cargando...</div>
            <div v-else-if="historyError" class="text-red-600">{{ historyError }}</div>
            <div v-else>
              <div v-if="historyData && historyData.length">
                <ul class="space-y-2">
                  <li v-for="(item, idx) in historyData" :key="idx" class="border-b pb-2">
                    <div><strong>Fecha:</strong> {{ formatDate(item.date) }}</div>
                    <div><strong>Acción:</strong> {{ item.action }}</div>
                    <div><strong>Usuario:</strong> {{ item.user }}</div>
                    <div v-if="item.details"><strong>Detalles:</strong> {{ item.details }}</div>
                  </li>
                </ul>
              </div>
              <div v-else class="text-gray-500">No hay movimientos registrados.</div>
            </div>
          </div>
        </div>
      </transition>
      <!-- QR Code Print Preview (hidden) -->
      <div ref="printArea" class="print-only">
        <div class="print-qr-card">
          <div class="print-header">
            <h1>MediTrack - {{ getTypeLabel(qrInfo.type) }}</h1>
            <p>{{ formatDate(new Date()) }}</p>
          </div>
          
          <div class="print-content">
            <div class="print-qr-section">
              <img 
                v-if="qrInfo.qr_code" 
                :src="getQRImageUrl(qrInfo.qr_code)" 
                alt="Código QR"
                class="print-qr-image"
              />
              <p class="print-qr-text">{{ qrInfo.qr_code }}</p>
            </div>
            
            <div class="print-info-section">
              <div v-if="qrInfo.type === 'batch' && qrInfo.batch_info" class="print-info-group">
                <h3>Información del Lote</h3>
                <p><strong>ID:</strong> {{ qrInfo.batch_info.id }}</p>
                <p><strong>Proveedor:</strong> {{ qrInfo.batch_info.supplier }}</p>
                <p><strong>Vencimiento:</strong> {{ formatDate(qrInfo.batch_info.expiration_date) }}</p>
                <p><strong>Cantidad:</strong> {{ qrInfo.batch_info.amount }} unidades</p>
              </div>
              
              <div v-if="qrInfo.type === 'medical_supply' && qrInfo.supply_info" class="print-info-group">
                <h3>Información del Producto</h3>
                <p><strong>Nombre:</strong> {{ qrInfo.supply_info.supply_code_name }}</p>
                <p><strong>Código:</strong> {{ qrInfo.supply_info.code }}</p>
                <p><strong>Proveedor:</strong> {{ qrInfo.supply_info.supplier }}</p>
                <p><strong>Estado:</strong> {{ qrInfo.supply_info.is_consumed ? 'Consumido' : 'Disponible' }}</p>
              </div>
            </div>
          </div>
          
          <div class="print-footer">
            <p>Generado automáticamente por MediTrack</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { format } from 'date-fns'
import { es } from 'date-fns/locale'
import qrService from '@/services/qrService'
import QRInfoDisplay from '@/components/QRInfoDisplay.vue'

const route = useRoute()
const router = useRouter()

// Referencias
const printArea = ref(null)

// Estado reactivo
const qrInfo = ref(null)
const loading = ref(false)
const error = ref(null)

// Toast notification
const toast = ref({ show: false, message: '' })
function showToast(msg, duration = 2500) {
  toast.value = { show: true, message: msg }
  setTimeout(() => { toast.value.show = false }, duration)
}

// Modal historial
const showHistoryModal = ref(false)
const historyLoading = ref(false)
const historyError = ref(null)
const historyData = ref([])

// Computed
const qrCode = computed(() => route.params.qrcode)

// Métodos principales
const loadQRInfo = async () => {
  loading.value = true
  error.value = null
  try {
    const result = await qrService.scanQRCode(qrCode.value)
    if (result) {
      qrInfo.value = result
    } else {
      error.value = 'No se pudo cargar la información del QR'
    }
  } catch (err) {
    console.error('Error loading QR info:', err)
    error.value = err.message || err.response?.data?.error || 'Error de conexión'
  } finally {
    loading.value = false
  }
}

const refreshData = () => {
  loadQRInfo()
  showToast('Datos actualizados')
}

const handleViewDetails = (qrInfo) => {
  // Ya estamos en la vista de detalles, mostrar mensaje informativo
  showToast('Ya estás viendo los detalles completos de este QR')
}

// Métodos de acción
const viewInInventory = () => {
  if (qrInfo.value?.type === 'batch') {
    router.push({ name: 'Inventory', query: { batch: qrInfo.value.id } })
  } else if (qrInfo.value?.type === 'medical_supply') {
    router.push({ name: 'Inventory', query: { supply: qrInfo.value.id } })
  }
}

const downloadQR = async () => {
  if (!qrInfo.value?.qr_code) return
  try {
    await qrService.downloadQRImage(qrInfo.value.qr_code, 'high')
    showToast('QR descargado correctamente')
  } catch (error) {
    console.error('Error downloading QR:', error)
    showToast('Error al descargar el código QR')
  }
}

const printQRCode = () => {
  if (!printArea.value) return
  const printContent = printArea.value.innerHTML
  const originalContent = document.body.innerHTML
  document.body.innerHTML = printContent
  window.print()
  document.body.innerHTML = originalContent
  window.location.reload()
  showToast('Etiqueta enviada a impresión')
}

const shareQR = async () => {
  const shareData = {
    title: `MediTrack - ${getTypeLabel(qrInfo.value?.type)}`,
    text: `Código QR: ${qrInfo.value?.qr_code}`,
    url: window.location.href
  }
  if (navigator.share) {
    try {
      await navigator.share(shareData)
      showToast('Enlace compartido correctamente')
    } catch (error) {
      console.log('Error sharing:', error)
      await copyToClipboard(window.location.href)
    }
  } else {
    await copyToClipboard(window.location.href)
  }
}

const copyToClipboard = async (text) => {
  try {
    await navigator.clipboard.writeText(text)
    showToast('Enlace copiado al portapapeles')
  } catch (error) {
    console.error('Error copying to clipboard:', error)
    showToast('No se pudo copiar el enlace')
  }
}

const openHistoryModal = async () => {
  showHistoryModal.value = true
  historyLoading.value = true
  historyError.value = null
  historyData.value = []
  try {
    const result = await qrService.getSupplyHistory(qrInfo.value.qr_code)
    if (result) {
      historyData.value = Array.isArray(result) ? result : [result]
    } else {
      historyError.value = 'No se pudo obtener el historial'
    }
  } catch (error) {
    console.error('Error getting history:', error)
    historyError.value = 'Error al obtener el historial'
  } finally {
    historyLoading.value = false
  }
}

const closeHistoryModal = () => {
  showHistoryModal.value = false
}

const addMovement = () => {
  router.push({ name: 'QRConsumer', query: { qr: qrInfo.value?.qr_code } })
}

const consumeSupply = (qrInfo) => {
  router.push({ name: 'QRConsumer', query: { qr: qrInfo.qr_code } })
}

const viewBatch = (batchId) => {
  router.push({ name: 'Inventory', query: { batch: batchId } })
}

const syncBatch = async () => {
  try {
    await qrService.syncBatchAmounts()
    showToast('Lote sincronizado correctamente')
    await refreshData()
  } catch (error) {
    console.error('Error syncing batch:', error)
    showToast('Error al sincronizar el lote')
  }
}

const generateReport = () => {
  const reportData = {
    qr_code: qrInfo.value?.qr_code,
    type: qrInfo.value?.type,
    generated_at: new Date().toISOString(),
    data: qrInfo.value
  }
  const blob = new Blob([JSON.stringify(reportData, null, 2)], { type: 'application/json' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = `reporte_${qrInfo.value?.qr_code}_${format(new Date(), 'yyyy-MM-dd')}.json`
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  URL.revokeObjectURL(url)
  showToast('Reporte generado y descargado')
}

const exportData = () => {
  generateReport()
}

// Utilidades
const getTypeLabel = (type) => {
  return qrService.getTypeLabel(type)
}

const getQRImageUrl = (qrCode) => {
  return qrService.getQRImageUrl(qrCode)
}

const formatDate = (dateString) => {
  if (!dateString) return 'No disponible'
  try {
    return format(new Date(dateString), 'dd/MM/yyyy HH:mm', { locale: es })
  } catch (error) {
    return dateString
  }
}

// Lifecycle
onMounted(() => {
  if (qrCode.value) {
    loadQRInfo()
  } else {
    error.value = 'Código QR no válido'
  }
})
</script>

