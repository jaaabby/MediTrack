<!-- QRDetails.vue CORREGIDO - CÓDIGO COMPLETO -->
<template>
  <div class="container mx-auto px-4 py-6">
    <!-- Header con breadcrumb -->
    <div class="mb-6">
      <nav class="flex mb-4" aria-label="Breadcrumb">
        <ol class="flex items-center space-x-2">
          <li><router-link to="/qr" class="text-blue-600 hover:text-blue-800">Escáner QR</router-link></li>
          <li><span class="text-gray-500">></span></li>
          <li><span class="text-gray-700 font-medium">{{ qrCode }}</span></li>
        </ol>
      </nav>
      <h1 class="text-2xl font-bold text-gray-900">
        Detalles Completos - {{ qrInfo ? getTypeLabel(qrInfo.type) : 'Cargando...' }}
      </h1>
    </div>

    <!-- Loading state -->
    <div v-if="loading" class="flex justify-center items-center h-64">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
    </div>

    <!-- Error state -->
    <div v-else-if="error" class="bg-red-50 border border-red-200 rounded-lg p-6 text-center">
      <svg class="h-12 w-12 text-red-400 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>
      <h3 class="text-lg font-medium text-red-800 mb-2">Error al cargar información</h3>
      <p class="text-red-600 mb-4">{{ error }}</p>
      <div class="space-x-3">
        <button @click="refreshData" class="btn-secondary text-sm">Intentar de Nuevo</button>
        <router-link to="/qr" class="btn-secondary text-sm">Volver al Escáner</router-link>
      </div>
    </div>

    <!-- Main content - VISTA EXPANDIDA -->
    <div v-else-if="qrInfo" class="space-y-6">
      
      <!-- Panel principal con información expandida -->
      <div class="grid lg:grid-cols-2 gap-6">
        
        <!-- Información básica del QR (lado izquierdo) -->
        <div class="bg-white rounded-lg shadow-sm border p-6">
          <h2 class="text-xl font-semibold text-gray-900 mb-4 flex items-center">
            <svg class="h-6 w-6 mr-2 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
            Información Detallada
          </h2>
          
          <!-- QR Image más grande -->
          <div class="text-center mb-6">
            <img 
              v-if="qrInfo.qr_code" 
              :src="getQRImageUrl(qrInfo.qr_code)" 
              alt="Código QR"
              class="mx-auto w-32 h-32 border rounded-lg shadow-sm"
            />
            <p class="mt-2 text-sm font-mono text-gray-600 bg-gray-50 px-3 py-1 rounded">
              {{ qrInfo.qr_code }}
            </p>
          </div>

          <!-- Información según tipo -->
          <div v-if="qrInfo.type === 'batch' && qrInfo.batch_info" class="space-y-4">
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="text-sm font-medium text-gray-600">ID:</label>
                <p class="text-gray-900 font-semibold">{{ qrInfo.batch_info.id || qrInfo.id || 'N/A' }}</p>
              </div>
              <div>
                <label class="text-sm font-medium text-gray-600">Cantidad:</label>
                <p class="text-gray-900">{{ qrInfo.batch_info.amount || 'N/A' }}</p>
              </div>
              <div>
                <label class="text-sm font-medium text-gray-600">Proveedor:</label>
                <p class="text-gray-900">{{ qrInfo.batch_info.supplier || 'N/A' }}</p>
              </div>
              <div>
                <label class="text-sm font-medium text-gray-600">Almacén:</label>
                <p class="text-gray-900">{{ qrInfo.batch_info.store_id || 'N/A' }}</p>
              </div>
              <div>
                <label class="text-sm font-medium text-gray-600">Vencimiento:</label>
                <p class="text-gray-900">{{ formatDate(qrInfo.batch_info.expiration_date) }}</p>
              </div>
              <div>
                <label class="text-sm font-medium text-gray-600">Creado:</label>
                <p class="text-gray-900">{{ formatDate(qrInfo.batch_info.created_at) }}</p>
              </div>
            </div>
          </div>

          <div v-if="qrInfo.type === 'medical_supply' && qrInfo.supply_info" class="space-y-4">
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="text-sm font-medium text-gray-600">Código:</label>
                <p class="text-gray-900 font-mono">{{ qrInfo.supply_info.SupplyCode?.code || qrInfo.supply_code || qrInfo.supply_info.supply_code || 'N/A' }}</p>
              </div>
              <div>
                <label class="text-sm font-medium text-gray-600">Estado:</label>
                <span :class="getStatusClass(qrInfo.current_status || qrInfo.status)" 
                      class="px-2 py-1 rounded-full text-xs font-medium">
                  {{ getStatusText(qrInfo.current_status || qrInfo.status) }}
                </span>
              </div>
              <div>
                <label class="text-sm font-medium text-gray-600">Nombre:</label>
                <p class="text-gray-900">{{ qrInfo.supply_info.name || qrInfo.supply_info.SupplyCode?.name || qrInfo.supply_info.supply_code_name || 'N/A' }}</p>
              </div>
              <div>
                <label class="text-sm font-medium text-gray-600">Proveedor:</label>
                <p class="text-gray-900">{{ qrInfo.supply_info.batch?.supplier || qrInfo.supply_info.supplier || 'N/A' }}</p>
              </div>
              <div>
                <label class="text-sm font-medium text-gray-600">Almacén:</label>
                <p class="text-gray-900">{{ qrInfo.supply_info.store_name || 'N/A' }}</p>
              </div>
              <div>
                <label class="text-sm font-medium text-gray-600">Vencimiento:</label>
                <p class="text-gray-900">{{ formatDate(qrInfo.supply_info.batch?.expiration_date || qrInfo.supply_info.expiration_date) }}</p>
              </div>
            </div>
          </div>
        </div>

        <!-- Panel de estadísticas y estado (lado derecho) -->
        <div class="space-y-4">
          
          <!-- Estadísticas del lote -->
          <div v-if="qrInfo.type === 'batch' && qrInfo.batch_status" class="bg-gradient-to-r from-blue-50 to-blue-100 rounded-lg p-6">
            <h3 class="text-lg font-semibold text-gray-900 mb-4">Estadísticas del Lote</h3>
            <div class="grid grid-cols-2 gap-4">
              <div class="text-center">
                <div class="text-3xl font-bold text-blue-600">{{ qrInfo.batch_status.total_individual_supplies || 0 }}</div>
                <div class="text-sm text-gray-600">Total Insumos</div>
              </div>
              <div class="text-center">
                <div class="text-3xl font-bold text-green-600">{{ qrInfo.batch_status.available_supplies || 0 }}</div>
                <div class="text-sm text-gray-600">Disponibles</div>
              </div>
              <div class="text-center">
                <div class="text-3xl font-bold text-red-600">{{ qrInfo.batch_status.consumed_supplies || 0 }}</div>
                <div class="text-sm text-gray-600">Consumidos</div>
              </div>
              <div class="text-center">
                <div class="text-3xl font-bold text-purple-600">{{ getUsagePercentage() }}%</div>
                <div class="text-sm text-gray-600">Utilización</div>
              </div>
            </div>
            <!-- Barra de progreso -->
            <div class="mt-4">
              <div class="bg-gray-200 rounded-full h-3">
                <div 
                  class="bg-gradient-to-r from-green-500 to-blue-500 h-3 rounded-full transition-all duration-300"
                  :style="`width: ${getUsagePercentage()}%`"
                ></div>
              </div>
            </div>
          </div>

          <!-- Panel de estado del insumo individual -->
          <div v-if="qrInfo.type === 'medical_supply'" class="bg-gradient-to-r from-green-50 to-green-100 rounded-lg p-6">
            <h3 class="text-lg font-semibold text-gray-900 mb-4">Estado del Insumo</h3>
            <div class="text-center">
              <div :class="getStatusIconClass(qrInfo.current_status || qrInfo.status)" 
                   class="text-4xl font-bold mb-2">
                {{ getStatusIcon(qrInfo.current_status || qrInfo.status) }}
              </div>
              <div class="text-lg font-medium text-gray-900">
                {{ getStatusText(qrInfo.current_status || qrInfo.status) }}
              </div>
              <div v-if="(qrInfo.current_status || qrInfo.status) === 'consumido'" class="text-sm text-gray-600 mt-2">
                Este insumo ya no está disponible para uso
              </div>
              <div v-else-if="qrInfo.can_consume || (qrInfo.current_status || qrInfo.status) === 'recepcionado'" class="text-sm text-gray-600 mt-2">
                Listo para ser consumido
              </div>
            </div>
          </div>

        </div>
      </div>

      <!-- Información del lote relacionado (para insumos individuales) -->
      <div v-if="qrInfo.type === 'medical_supply' && (qrInfo.supply_info?.batch || qrInfo.batch_status)" class="bg-blue-50 rounded-lg border border-blue-200 p-6">
        <h2 class="text-xl font-semibold text-gray-900 mb-4 flex items-center">
          <svg class="h-6 w-6 mr-2 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
          </svg>
          Información del Lote Relacionado
        </h2>
        <div class="grid grid-cols-2 gap-4 mb-4">
          <div>
            <label class="text-sm font-medium text-gray-600">ID del Lote:</label>
            <p class="text-gray-900 font-mono">{{ qrInfo.supply_info?.batch?.id || qrInfo.supply_info?.BatchID || qrInfo.batch_status?.batch_id || 'N/A' }}</p>
          </div>
          <div>
            <label class="text-sm font-medium text-gray-600">Cantidad Disponible:</label>
            <p class="text-gray-900 font-bold" :class="(qrInfo.supply_info?.batch?.amount || qrInfo.batch_status?.current_amount || 0) > 0 ? 'text-green-600' : 'text-red-600'">
              {{ qrInfo.supply_info?.batch?.amount || qrInfo.batch_status?.current_amount || 0 }} unidades
            </p>
          </div>
          <div>
            <label class="text-sm font-medium text-gray-600">Proveedor:</label>
            <p class="text-gray-900">{{ qrInfo.supply_info?.batch?.supplier || qrInfo.batch_status?.supplier || 'N/A' }}</p>
          </div>
          <div>
            <label class="text-sm font-medium text-gray-600">Vencimiento:</label>
            <p class="text-gray-900">{{ formatDate(qrInfo.supply_info?.batch?.expiration_date || qrInfo.batch_status?.expiration_date) }}</p>
          </div>
          <div>
            <label class="text-sm font-medium text-gray-600">Stock:</label>
            <span :class="(qrInfo.supply_info?.batch?.amount || qrInfo.batch_status?.current_amount || 0) > 0 ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'"
                  class="px-2 py-1 rounded-full text-xs font-medium">
              {{ (qrInfo.supply_info?.batch?.amount || qrInfo.batch_status?.current_amount || 0) > 0 ? 'Disponible' : 'Agotado' }}
            </span>
          </div>
          <div v-if="qrInfo.batch_status?.batch_qr_code || qrInfo.supply_info?.batch?.qr_code">
            <label class="text-sm font-medium text-gray-600">QR del Lote:</label>
            <p class="text-gray-900 font-mono text-xs">{{ qrInfo.batch_status?.batch_qr_code || qrInfo.supply_info?.batch?.qr_code }}</p>
          </div>
        </div>
        <button 
          v-if="qrInfo.supply_info?.batch?.id || qrInfo.supply_info?.BatchID || qrInfo.batch_status?.batch_id"
          @click="viewBatch(qrInfo.supply_info?.batch?.id || qrInfo.supply_info?.BatchID || qrInfo.batch_status?.batch_id)" 
          class="btn-primary"
        >
          Ver Lote Completo en Inventario
        </button>
      </div>

      <!-- Historial SIEMPRE VISIBLE (no en modal) -->
      <div class="bg-white rounded-lg shadow-sm border p-6">
        <div class="flex items-center justify-between mb-4">
          <h2 class="text-xl font-semibold text-gray-900 flex items-center">
            <svg class="h-6 w-6 mr-2 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            Historial de Movimientos
          </h2>
          <button @click="loadHistory" class="flex items-center px-3 py-2 text-sm bg-gray-100 hover:bg-gray-200 text-gray-700 rounded-lg transition-colors">
            <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
            Actualizar
          </button>
        </div>
        
        <div v-if="historyLoading" class="text-center py-8">
          <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>
          <p class="mt-2 text-gray-600">Cargando historial...</p>
        </div>
        
        <div v-else-if="historyError" class="text-center py-8">
          <svg class="h-12 w-12 text-red-400 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <p class="text-red-600">{{ historyError }}</p>
          <button @click="loadHistory" class="mt-2 btn-secondary text-sm">Reintentar</button>
        </div>
        
        <div v-else-if="historyData && historyData.length > 0" class="space-y-3">
          <div 
            v-for="(item, index) in historyData" 
            :key="index"
            class="flex items-center justify-between p-4 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors"
          >
            <div class="flex items-center">
              <div :class="getHistoryIconClass(item)" class="p-2 rounded-full mr-4">
                <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path v-if="getHistoryStatus(item) === 'consumido'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                  <path v-else-if="getHistoryStatus(item) === 'creado'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
                  <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </div>
              <div>
                <p class="font-semibold text-gray-900">{{ getHistoryStatusFormatted(item) }}</p>
                <p class="text-sm text-gray-600">{{ formatDate(item.date_time) }}</p>
                <p v-if="item.notes" class="text-xs text-gray-500 mt-1">{{ item.notes }}</p>
              </div>
            </div>
            <div class="text-right">
              <p class="text-sm font-medium text-gray-700">{{ item.user_rut || 'N/A' }}</p>
              <p class="text-xs text-gray-500">
                {{ getDestinationLabel(item.destination_type) }}: {{ item.destination_id || 'N/A' }}
              </p>
            </div>
          </div>
        </div>
        
        <div v-else class="text-center py-8 text-gray-500">
          <svg class="h-16 w-16 mx-auto mb-4 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
          </svg>
          <p class="mb-4">No hay movimientos registrados para este QR</p>
          <button @click="addMovement" class="btn-primary">
            Agregar Primer Movimiento
          </button>
        </div>
      </div>

      <!-- Acciones administrativas expandidas -->
      <div class="bg-white rounded-lg shadow-sm border p-6">
        <h2 class="text-xl font-semibold text-gray-900 mb-4">Acciones Administrativas</h2>
        <div class="grid sm:grid-cols-2 lg:grid-cols-4 gap-4">
          <button @click="viewInInventory" class="action-card-detailed">
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

          <button v-if="!qrInfo.is_consumed" @click="addMovement" class="action-card-detailed">
            <div class="action-icon bg-orange-100 text-orange-600">
              <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
              </svg>
            </div>
            <div class="action-content">
              <div class="font-medium">Consumir Insumo</div>
              <div class="text-sm text-gray-500">Registrar consumo</div>
            </div>
          </button>

          <button v-if="qrInfo.type === 'batch'" @click="syncBatch" class="action-card-detailed">
            <div class="action-icon bg-purple-100 text-purple-600">
              <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
              </svg>
            </div>
            <div class="action-content">
              <div class="font-medium">Sincronizar</div>
              <div class="text-sm text-gray-500">Actualizar cantidades</div>
            </div>
          </button>

          <button @click="refreshData" class="action-card-detailed">
            <div class="action-icon bg-gray-100 text-gray-600">
              <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
              </svg>
            </div>
            <div class="action-content">
              <div class="font-medium">Actualizar</div>
              <div class="text-sm text-gray-500">Refrescar datos</div>
            </div>
          </button>
        </div>
      </div>
    </div>

    <!-- Toast Notification -->
    <transition name="fade">
      <div v-if="toast.show" class="fixed bottom-4 right-4 bg-gray-800 text-white px-6 py-3 rounded-lg shadow-lg z-50">
        {{ toast.message }}
      </div>
    </transition>

    <!-- Área de impresión (oculta) -->
    <div ref="printArea" class="print-only">
      <div class="print-qr-card">
        <div class="print-header">
          <h1>MediTrack - {{ qrInfo ? getTypeLabel(qrInfo.type) : '' }}</h1>
          <p>{{ formatDate(new Date()) }}</p>
        </div>
        
        <div class="print-content">
          <div class="print-qr-section">
            <img 
              v-if="qrInfo?.qr_code" 
              :src="getQRImageUrl(qrInfo.qr_code)" 
              alt="Código QR"
              class="print-qr-image"
            />
            <p class="print-qr-text">{{ qrInfo?.qr_code }}</p>
          </div>
          
          <div class="print-info-section">
            <div v-if="qrInfo?.type === 'batch' && qrInfo.batch_info" class="print-info-group">
              <h3>Información del Lote</h3>
              <p><strong>ID:</strong> {{ qrInfo.batch_info.id }}</p>
              <p><strong>Proveedor:</strong> {{ qrInfo.batch_info.supplier }}</p>
              <p><strong>Vencimiento:</strong> {{ formatDate(qrInfo.batch_info.expiration_date) }}</p>
              <p><strong>Cantidad:</strong> {{ qrInfo.batch_info.amount }} unidades</p>
            </div>
            
            <div v-if="qrInfo?.type === 'medical_supply' && qrInfo.supply_info" class="print-info-group">
              <h3>Información del Producto</h3>
              <p><strong>Nombre:</strong> {{ qrInfo.supply_info.supply_code_name }}</p>
              <p><strong>Código:</strong> {{ qrInfo.supply_info.supply_code }}</p>
              <p><strong>Proveedor:</strong> {{ qrInfo.supply_info.supplier }}</p>
              <p><strong>Estado:</strong> {{ qrInfo.is_consumed ? 'Consumido' : 'Disponible' }}</p>
            </div>
          </div>
        </div>
        
        <div class="print-footer">
          <p>Generado automáticamente por MediTrack</p>
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

const route = useRoute()
const router = useRouter()

// Referencias
const printArea = ref(null)

// Estado reactivo
const qrInfo = ref(null)
const loading = ref(false)
const error = ref(null)
const historyLoading = ref(false)
const historyError = ref(null)
const historyData = ref([])

// Toast notification
const toast = ref({ show: false, message: '' })
function showToast(msg, duration = 2500) {
  toast.value = { show: true, message: msg }
  setTimeout(() => { toast.value.show = false }, duration)
}

// Computed
const qrCode = computed(() => route.params.qrCode || route.params.qrcode)

// Métodos principales
const loadQRInfo = async () => {
  loading.value = true
  error.value = null
  try {
    const result = await qrService.scanQRCode(qrCode.value)
    console.log('QR Info loaded:', result) // Debug log
    if (result) {
      qrInfo.value = result
      
      // Si es un insumo individual, intentar obtener el QR del lote
      if (result.type === 'medical_supply' && result.batch_status?.batch_id && !result.batch_status.batch_qr_code) {
        try {
          // Buscar el QR del lote por ID
          const batchInfo = await qrService.getBatchById(result.batch_status.batch_id)
          if (batchInfo?.qr_code) {
            qrInfo.value.batch_status.batch_qr_code = batchInfo.qr_code
          }
        } catch (err) {
          console.log('No se pudo obtener QR del lote:', err)
        }
      }
      
      // Cargar historial automáticamente
      await loadHistory()
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

const loadHistory = async () => {
  if (!qrInfo.value?.qr_code) return
  
  historyLoading.value = true
  historyError.value = null
  try {
    const result = await qrService.getSupplyHistory(qrInfo.value.qr_code)
    console.log('History loaded:', result) // Debug log
    
    let historyArray = []
    
    if (result) {
      // Manejar diferentes estructuras de respuesta
      if (Array.isArray(result)) {
        historyArray = result
      } else if (result.data && Array.isArray(result.data)) {
        historyArray = result.data
      } else if (result.success && Array.isArray(result.data)) {
        historyArray = result.data
      } else if (typeof result === 'object') {
        historyArray = [result]
      }
      
      console.log('Processed history array:', historyArray) // Debug log
      
      // Debug cada item del historial
      historyArray.forEach((item, index) => {
        console.log(`History item ${index}:`, {
          status: item.status,
          date_time: item.date_time,
          user_rut: item.user_rut,
          destination_type: item.destination_type,
          destination_id: item.destination_id,
          medical_supply_id: item.medical_supply_id,
          id: item.id
        })
      })
      
      historyData.value = historyArray
    } else {
      historyData.value = []
    }
  } catch (error) {
    console.error('Error getting history:', error)
    historyError.value = 'Error al obtener el historial'
    historyData.value = []
  } finally {
    historyLoading.value = false
  }
}

const refreshData = async () => {
  await loadQRInfo()
  showToast('Datos actualizados')
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
    title: `MediTrack - ${qrInfo.value ? getTypeLabel(qrInfo.value.type) : 'QR'}`,
    text: `Código QR: ${qrInfo.value?.qr_code}`,
    url: window.location.href
  }
  if (navigator.share) {
    try {
      await navigator.share(shareData)
      showToast('Enlace compartido correctamente')
    } catch (error) {
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
    showToast('No se pudo copiar el enlace')
  }
}

const addMovement = () => {
  router.push({ name: 'QRConsumer', query: { qr: qrInfo.value?.qr_code } })
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
    data: qrInfo.value,
    history: historyData.value
  }
  const blob = new Blob([JSON.stringify(reportData, null, 2)], { type: 'application/json' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = `reporte_detallado_${qrInfo.value?.qr_code}_${format(new Date(), 'yyyy-MM-dd')}.json`
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  URL.revokeObjectURL(url)
  showToast('Reporte detallado generado y descargado')
}

// Utilidades
const getTypeLabel = (type) => {
  if (!type) return 'Desconocido'
  const labels = {
    'medical_supply': 'Insumo Individual',
    'batch': 'Lote',
    'supply': 'Insumo Individual'
  }
  return labels[type] || 'Desconocido'
}

const getQRImageUrl = (qrCode) => {
  return qrService.getQRImageUrl(qrCode)
}

const getHistoryStatus = (item) => {
  return item?.status || item?.action || 'desconocido'
}

const getHistoryStatusFormatted = (item) => {
  const status = getHistoryStatus(item)
  if (!status || status === 'desconocido') return 'Acción no especificada'
  
  const statusLabels = {
    'consumido': 'Consumido',
    'creado': 'Creado',
    'recibido': 'Recibido',
    'entregado': 'Entregado',
    'devuelto': 'Devuelto',
    'perdido': 'Perdido',
    'dañado': 'Dañado'
  }
  
  return statusLabels[status.toLowerCase()] || (status.charAt(0).toUpperCase() + status.slice(1))
}

const getHistoryIconClass = (item) => {
  const status = getHistoryStatus(item)
  
  const statusColors = {
    'consumido': 'bg-red-100 text-red-600',
    'creado': 'bg-green-100 text-green-600',
    'recibido': 'bg-blue-100 text-blue-600',
    'entregado': 'bg-orange-100 text-orange-600',
    'devuelto': 'bg-yellow-100 text-yellow-600',
    'perdido': 'bg-gray-100 text-gray-600',
    'dañado': 'bg-red-100 text-red-600'
  }
  
  return statusColors[status?.toLowerCase()] || 'bg-gray-100 text-gray-600'
}

const getDestinationLabel = (destinationType) => {
  if (!destinationType) return 'Destino'
  
  const labels = {
    'pavilion': 'Pabellón',
    'store': 'Almacén',
    'almacen': 'Almacén'
  }
  
  return labels[destinationType.toLowerCase()] || destinationType
}

const getUsagePercentage = () => {
  if (!qrInfo.value?.batch_status) return 0
  const total = qrInfo.value.batch_status.total_individual_supplies || 0
  const consumed = qrInfo.value.batch_status.consumed_supplies || 0
  if (total === 0) return 0
  return Math.round((consumed / total) * 100)
}

const formatDate = (dateString) => {
  if (!dateString) return 'No disponible'
  try {
    return format(new Date(dateString), 'dd/MM/yyyy HH:mm', { locale: es })
  } catch (error) {
    return dateString
  }
}

// Funciones helper para estado
const getStatusText = (status) => {
  if (!status) return 'Desconocido'
  const statusLabels = {
    'disponible': 'Disponible',
    'recepcionado': 'Recepcionado',
    'en_camino_a_pabellon': 'En Camino a Pabellón',
    'en_camino_a_bodega': 'En Camino a Bodega',
    'consumido': 'Consumido'
  }
  return statusLabels[status.toLowerCase()] || status.charAt(0).toUpperCase() + status.slice(1)
}

const getStatusClass = (status) => {
  if (!status) return 'bg-gray-100 text-gray-800'
  const statusColors = {
    'disponible': 'bg-green-100 text-green-800',
    'recepcionado': 'bg-blue-100 text-blue-800',
    'en_camino_a_pabellon': 'bg-yellow-100 text-yellow-800',
    'en_camino_a_bodega': 'bg-orange-100 text-orange-800',
    'consumido': 'bg-red-100 text-red-800'
  }
  return statusColors[status.toLowerCase()] || 'bg-gray-100 text-gray-800'
}

const getStatusIcon = (status) => {
  if (!status) return '❓'
  const statusIcons = {
    'disponible': '✅',
    'recepcionado': '📥',
    'en_camino_a_pabellon': '🚚',
    'en_camino_a_bodega': '📤',
    'consumido': '❌'
  }
  return statusIcons[status.toLowerCase()] || '❓'
}

const getStatusIconClass = (status) => {
  if (!status) return 'text-gray-600'
  const statusColors = {
    'disponible': 'text-green-600',
    'recepcionado': 'text-blue-600',
    'en_camino_a_pabellon': 'text-yellow-600',
    'en_camino_a_bodega': 'text-orange-600',
    'consumido': 'text-red-600'
  }
  return statusColors[status.toLowerCase()] || 'text-gray-600'
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

<style scoped>
.action-btn {
  @apply flex items-center justify-center space-x-2 px-4 py-2 text-white font-medium rounded-lg transition-colors;
}

.action-card-detailed {
  @apply flex items-center p-4 border rounded-lg hover:bg-gray-50 transition-colors cursor-pointer;
}

.action-icon {
  @apply flex items-center justify-center w-12 h-12 rounded-lg mr-4;
}

.action-content {
  @apply flex-1;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* Usar clases de botones de style.css global */

/* Estilos de impresión */
@media print {
  .print-only {
    display: block !important;
  }
  
  .print-qr-card {
    width: 100%;
    max-width: 600px;
    margin: 0 auto;
    padding: 20px;
    border: 2px solid #000;
    background: white;
  }
  
  .print-header {
    text-align: center;
    margin-bottom: 20px;
    border-bottom: 1px solid #ccc;
    padding-bottom: 10px;
  }
  
  .print-header h1 {
    font-size: 24px;
    font-weight: bold;
    margin: 0;
  }
  
  .print-content {
    display: flex;
    align-items: center;
    gap: 20px;
  }
  
  .print-qr-section {
    flex-shrink: 0;
    text-align: center;
  }
  
  .print-qr-image {
    width: 150px;
    height: 150px;
    border: 1px solid #ddd;
  }
  
  .print-qr-text {
    font-family: monospace;
    font-size: 12px;
    margin-top: 10px;
    word-break: break-all;
  }
  
  .print-info-section {
    flex: 1;
  }
  
  .print-info-group {
    margin-bottom: 15px;
  }
  
  .print-info-group h3 {
    font-size: 16px;
    font-weight: bold;
    margin-bottom: 8px;
    color: #333;
  }
  
  .print-info-group p {
    font-size: 12px;
    margin: 2px 0;
    color: #666;
  }
  
  .print-footer {
    text-align: center;
    margin-top: 20px;
    padding-top: 10px;
    border-top: 1px solid #ccc;
    font-size: 10px;
    color: #999;
  }
}

.print-only {
  display: none;
}

/* Container responsivo */
.container {
  max-width: 1200px;
}

/* Asegurar que las tarjetas se vean bien en móvil */
@media (max-width: 768px) {
  .grid {
    grid-template-columns: 1fr;
  }
  
  .action-card-detailed {
    flex-direction: column;
    text-align: center;
  }
  
  .action-icon {
    margin: 0 0 8px 0;
  }
}
</style>