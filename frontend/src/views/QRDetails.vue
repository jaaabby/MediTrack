<template>
  <div class="space-y-6">
    <!-- Breadcrumb y acciones -->
    <div class="flex items-center justify-between">
      <nav class="flex" aria-label="Breadcrumb">
        <ol class="inline-flex items-center space-x-1 md:space-x-3">
          <li class="inline-flex items-center">
            <router-link to="/qr" class="text-gray-700 hover:text-primary-600">
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
      
      <div class="flex space-x-2">
        <button @click="refreshData" class="btn-secondary" :disabled="loading">
          <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
          Actualizar
        </button>
        <button @click="exportData" class="btn-primary">
          <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
          </svg>
          Exportar
        </button>
      </div>
    </div>

    <!-- Loading state -->
    <div v-if="loading" class="flex justify-center items-center py-12">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-600"></div>
      <span class="ml-2 text-gray-600">Cargando información...</span>
    </div>

    <!-- Error state -->
    <div v-else-if="error" class="bg-red-50 border border-red-200 rounded-md p-4">
      <div class="flex">
        <div class="flex-shrink-0">
          <svg class="h-5 w-5 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </div>
        <div class="ml-3">
          <h3 class="text-sm font-medium text-red-800">Error al cargar información</h3>
          <div class="mt-2 text-sm text-red-700">{{ error }}</div>
          <div class="mt-4">
            <button @click="refreshData" class="btn-secondary text-sm">
              Reintentar
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Contenido principal -->
    <div v-else-if="qrInfo" class="space-y-6">
      <!-- Header con información básica -->
      <div class="bg-gradient-to-r from-primary-600 to-primary-700 rounded-xl p-6 text-white">
        <div class="flex items-center justify-between">
          <div class="flex items-center">
            <div class="bg-white bg-opacity-20 p-3 rounded-lg mr-4">
              <svg class="h-8 w-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path v-if="qrInfo.type === 'batch'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
                <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2zM9 9h6v6H9V9z" />
              </svg>
            </div>
            <div>
              <h1 class="text-2xl font-bold">{{ getTypeLabel(qrInfo.type) }}</h1>
              <p class="text-primary-100">ID: {{ qrInfo.id }} | Tipo: {{ qrInfo.type.toUpperCase() }}</p>
            </div>
          </div>
          <div class="text-right">
            <p class="text-primary-100 text-sm">Código QR</p>
            <code class="bg-white bg-opacity-20 px-3 py-1 rounded font-mono text-sm">
              {{ qrInfo.qr_code }}
            </code>
          </div>
        </div>
      </div>

      <!-- Información detallada -->
      <div class="grid lg:grid-cols-2 gap-6">
        <!-- Información principal -->
        <div class="card">
          <div class="card-header">
            <h2 class="card-title">Información Principal</h2>
          </div>
          
          <div v-if="qrInfo.type === 'batch' && qrInfo.batch_info" class="space-y-4">
            <dl class="grid grid-cols-1 gap-4">
              <div>
                <dt class="text-sm font-medium text-gray-500">Fecha de Vencimiento</dt>
                <dd :class="getExpirationClass(qrInfo.batch_info.expiration_date)" class="mt-1 text-lg font-semibold">
                  {{ formatDate(qrInfo.batch_info.expiration_date) }}
                </dd>
                <dd class="text-sm text-gray-500">
                  {{ getExpirationWarning(qrInfo.batch_info.expiration_date) }}
                </dd>
              </div>
              
              <div>
                <dt class="text-sm font-medium text-gray-500">Cantidad Total</dt>
                <dd class="mt-1 text-lg font-semibold text-gray-900">
                  {{ qrInfo.batch_info.amount }} unidades
                </dd>
              </div>
              
              <div>
                <dt class="text-sm font-medium text-gray-500">Proveedor</dt>
                <dd class="mt-1 text-lg font-semibold text-gray-900">
                  {{ qrInfo.batch_info.supplier }}
                </dd>
              </div>
              
              <div>
                <dt class="text-sm font-medium text-gray-500">ID de Bodega</dt>
                <dd class="mt-1 text-lg font-semibold text-gray-900">
                  {{ qrInfo.batch_info.store_id }}
                </dd>
              </div>
            </dl>
          </div>

          <div v-else-if="qrInfo.type === 'medical_supply' && qrInfo.supply_info" class="space-y-4">
            <dl class="grid grid-cols-1 gap-4">
              <div>
                <dt class="text-sm font-medium text-gray-500">Nombre del Insumo</dt>
                <dd class="mt-1 text-lg font-semibold text-gray-900">
                  {{ qrInfo.supply_info.supply_code_name }}
                </dd>
              </div>
              
              <div>
                <dt class="text-sm font-medium text-gray-500">Código Interno</dt>
                <dd class="mt-1 text-lg font-semibold text-gray-900">
                  {{ qrInfo.supply_info.code }}
                </dd>
              </div>
              
              <div>
                <dt class="text-sm font-medium text-gray-500">Fecha de Vencimiento</dt>
                <dd :class="getExpirationClass(qrInfo.supply_info.expiration_date)" class="mt-1 text-lg font-semibold">
                  {{ formatDate(qrInfo.supply_info.expiration_date) }}
                </dd>
                <dd class="text-sm text-gray-500">
                  {{ getExpirationWarning(qrInfo.supply_info.expiration_date) }}
                </dd>
              </div>
              
              <div>
                <dt class="text-sm font-medium text-gray-500">Proveedor</dt>
                <dd class="mt-1 text-lg font-semibold text-gray-900">
                  {{ qrInfo.supply_info.supplier }}
                </dd>
              </div>
              
              <div>
                <dt class="text-sm font-medium text-gray-500">Bodega</dt>
                <dd class="mt-1 text-lg font-semibold text-gray-900">
                  {{ qrInfo.supply_info.store_name }}
                </dd>
              </div>
              
              <div>
                <dt class="text-sm font-medium text-gray-500">ID del Lote</dt>
                <dd class="mt-1 text-lg font-semibold text-gray-900">
                  {{ qrInfo.supply_info.batch_id }}
                </dd>
              </div>
            </dl>
          </div>
        </div>

        <!-- Código de insumo (si está disponible) -->
        <div v-if="qrInfo.supply_code" class="card">
          <div class="card-header">
            <h2 class="card-title">Información del Código de Insumo</h2>
          </div>
          
          <dl class="space-y-4">
            <div>
              <dt class="text-sm font-medium text-gray-500">Código</dt>
              <dd class="mt-1 text-lg font-semibold text-gray-900">
                {{ qrInfo.supply_code.code }}
              </dd>
            </div>
            
            <div>
              <dt class="text-sm font-medium text-gray-500">Nombre</dt>
              <dd class="mt-1 text-lg font-semibold text-gray-900">
                {{ qrInfo.supply_code.name }}
              </dd>
            </div>
            
            <div>
              <dt class="text-sm font-medium text-gray-500">Código del Proveedor</dt>
              <dd class="mt-1 text-lg font-semibold text-gray-900">
                {{ qrInfo.supply_code.code_supplier }}
              </dd>
            </div>
            
            <div>
              <dt class="text-sm font-medium text-gray-500">ID del Lote</dt>
              <dd class="mt-1 text-lg font-semibold text-gray-900">
                {{ qrInfo.supply_code.batch_id }}
              </dd>
            </div>
          </dl>
        </div>
      </div>

      <!-- Historial completo -->
      <div v-if="qrInfo.history && qrInfo.history.length > 0" class="card">
        <div class="card-header">
          <h2 class="card-title flex items-center">
            <svg class="h-6 w-6 mr-2 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            Historial Completo de Movimientos
          </h2>
          <p class="text-sm text-gray-600">{{ qrInfo.history.length }} movimientos registrados</p>
        </div>
        
        <div class="overflow-x-auto">
          <table class="table">
            <thead class="table-header">
              <tr>
                <th class="table-header-cell">Fecha y Hora</th>
                <th class="table-header-cell">Estado</th>
                <th class="table-header-cell">Tipo de Destino</th>
                <th class="table-header-cell">ID Destino</th>
                <th class="table-header-cell">Usuario</th>
              </tr>
            </thead>
            <tbody class="table-body">
              <tr v-for="movement in qrInfo.history" :key="movement.id" class="table-row">
                <td class="table-cell">
                  <div>
                    <div class="font-medium">{{ formatDate(movement.date_time) }}</div>
                    <div class="text-sm text-gray-500">{{ formatTime(movement.date_time) }}</div>
                  </div>
                </td>
                <td class="table-cell">
                  <span :class="getStatusBadgeClass(movement.status)" class="badge">
                    {{ movement.status }}
                  </span>
                </td>
                <td class="table-cell">
                  <span class="font-medium">{{ movement.destination_type }}</span>
                </td>
                <td class="table-cell">{{ movement.destination_id }}</td>
                <td class="table-cell">
                  <code class="text-sm">{{ movement.user_rut }}</code>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- Acciones disponibles -->
      <div class="card">
        <div class="card-header">
          <h2 class="card-title">Acciones Disponibles</h2>
        </div>
        
        <div class="grid md:grid-cols-3 gap-4">
          <button @click="viewInInventory" class="btn-secondary p-4 h-auto">
            <svg class="h-6 w-6 mx-auto mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
            </svg>
            <div class="text-sm font-medium">Ver en Inventario</div>
          </button>
          
          <button @click="printQRCode" class="btn-secondary p-4 h-auto">
            <svg class="h-6 w-6 mx-auto mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 17h2a2 2 0 002-2v-4a2 2 0 00-2-2H5a2 2 0 00-2 2v4a2 2 0 002 2h2m2 4h6a2 2 0 002-2v-4a2 2 0 00-2-2H9a2 2 0 00-2 2v4a2 2 0 002 2zm8-12V5a2 2 0 00-2-2H9a2 2 0 00-2 2v4h10z" />
            </svg>
            <div class="text-sm font-medium">Imprimir QR</div>
          </button>
          
          <button @click="addMovement" class="btn-primary p-4 h-auto">
            <svg class="h-6 w-6 mx-auto mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
            </svg>
            <div class="text-sm font-medium">Registrar Movimiento</div>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { format } from 'date-fns'
import { es } from 'date-fns/locale'
import qrService from '@/services/qrService'

const route = useRoute()
const router = useRouter()

// Estado reactivo
const qrInfo = ref(null)
const loading = ref(false)
const error = ref(null)

// Métodos
const loadQRInfo = async () => {
  loading.value = true
  error.value = null
  
  try {
    const qrCode = route.params.qrcode
    const result = await qrService.scanQRCode(qrCode)
    
    if (result.success) {
      qrInfo.value = result.data
    } else {
      error.value = result.error || 'Error al cargar información del QR'
    }
  } catch (err) {
    error.value = err.response?.data?.error || 'Error de conexión'
  } finally {
    loading.value = false
  }
}

const refreshData = () => {
  loadQRInfo()
}

const getTypeLabel = (type) => {
  return type === 'batch' ? 'Lote de Insumos Médicos' : 'Insumo Médico Individual'
}

const getExpirationClass = (expirationDate) => {
  const today = new Date()
  const expDate = new Date(expirationDate)
  const daysUntilExpiration = Math.ceil((expDate - today) / (1000 * 60 * 60 * 24))
  
  if (daysUntilExpiration < 0) return 'text-red-600'
  if (daysUntilExpiration <= 30) return 'text-orange-600'
  return 'text-green-600'
}

const getExpirationWarning = (expirationDate) => {
  const today = new Date()
  const expDate = new Date(expirationDate)
  const daysUntilExpiration = Math.ceil((expDate - today) / (1000 * 60 * 60 * 24))
  
  if (daysUntilExpiration < 0) return `Vencido hace ${Math.abs(daysUntilExpiration)} días`
  if (daysUntilExpiration === 0) return 'Vence hoy'
  if (daysUntilExpiration <= 30) return `Vence en ${daysUntilExpiration} días`
  return `${daysUntilExpiration} días hasta el vencimiento`
}

const getStatusBadgeClass = (status) => {
  const statusClasses = {
    'activo': 'badge-success',
    'en_transito': 'badge-warning',
    'utilizado': 'badge-info',
    'vencido': 'badge-danger'
  }
  return statusClasses[status] || 'badge-info'
}

const formatDate = (dateString) => {
  try {
    return format(new Date(dateString), 'dd/MM/yyyy', { locale: es })
  } catch {
    return dateString
  }
}

const formatTime = (dateString) => {
  try {
    return format(new Date(dateString), 'HH:mm:ss', { locale: es })
  } catch {
    return ''
  }
}

const exportData = () => {
  const dataStr = JSON.stringify(qrInfo.value, null, 2)
  const blob = new Blob([dataStr], { type: 'application/json' })
  const url = URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = `qr-${qrInfo.value.qr_code}-${new Date().toISOString().split('T')[0]}.json`
  link.click()
  URL.revokeObjectURL(url)
}

const viewInInventory = () => {
  router.push('/inventory')
}

const printQRCode = () => {
  // Implementar impresión del código QR
  console.log('Imprimir QR:', qrInfo.value.qr_code)
}

const addMovement = () => {
  // Implementar formulario para agregar movimiento
  console.log('Agregar movimiento para:', qrInfo.value.qr_code)
}

// Lifecycle
onMounted(() => {
  // Verificar si hay información en el state de la navegación
  if (history.state?.qrInfo) {
    qrInfo.value = history.state.qrInfo
  } else {
    loadQRInfo()
  }
})
</script>