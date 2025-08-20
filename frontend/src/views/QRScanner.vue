<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="text-center">
      <h1 class="text-3xl font-bold text-gray-900 mb-2">Escáner de Códigos QR</h1>
      <p class="text-gray-600">Escanea códigos QR para ver información de insumos médicos</p>
    </div>

    <!-- Escáner principal -->
    <div class="card max-w-2xl mx-auto">
      <div class="card-header text-center">
        <h3 class="card-title text-blue-700">
          <svg class="h-6 w-6 inline mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" />
          </svg>
          Escanear Código QR
        </h3>
      </div>
      
      <div class="space-y-6">
        <!-- Área del escáner -->
        <div class="bg-gray-100 rounded-lg p-8 border-2 border-dashed border-gray-300 text-center">
          <svg class="h-20 w-20 mx-auto text-gray-400 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h.01M5 8h2a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v2a2 2 0 002 2zm0 0h2a2 2 0 012 2v2a2 2 0 01-2 2H5a2 2 0 01-2-2v-2a2 2 0 012-2z" />
          </svg>
          <p class="text-gray-500 mb-6">Enfoca la cámara hacia el código QR del insumo médico</p>
          <button class="btn-primary mb-4" disabled>
            <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" />
            </svg>
            Activar Cámara (Próximamente)
          </button>
        </div>
        
        <!-- Entrada manual para desarrollo -->
        <div class="border-t pt-6">
          <label class="block text-sm font-medium text-gray-700 mb-3">
            <svg class="h-4 w-4 inline mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
            </svg>
            Ingresa el código QR manualmente
          </label>
          <div class="flex space-x-3">
            <input
              type="text"
              placeholder="Ejemplo: SUPPLY_1755580808_1"
              class="form-input flex-1"
              v-model="qrInput"
              @keyup.enter="scanQRCode"
            />
            <button 
              @click="scanQRCode" 
              :disabled="!qrInput.trim() || loading" 
              class="btn-primary px-6"
            >
              <svg v-if="!loading" class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
              <div v-else class="h-4 w-4 mr-2 spinner"></div>
              {{ loading ? 'Escaneando...' : 'Escanear' }}
            </button>
          </div>
          <p class="text-xs text-gray-500 mt-2">
            💡 Tip: Puedes usar códigos como SUPPLY_1755580808_1, SUPPLY_1755580808_2, etc.
          </p>
        </div>
      </div>
    </div>

    <!-- Información del insumo escaneado -->
    <div v-if="scannedInfo && !error" class="card">
      <div class="card-header">
        <h3 class="card-title text-green-700">
          <svg class="h-6 w-6 inline mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          Información del Insumo Médico
        </h3>
      </div>
      <QRInfoDisplay :qr-info="scannedInfo" @view-details="viewDetails" />
    </div>

    <!-- Mensaje de error -->
    <div v-if="error" class="bg-red-50 border border-red-200 rounded-md p-4">
      <div class="flex">
        <div class="flex-shrink-0">
          <svg class="h-5 w-5 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </div>
        <div class="ml-3">
          <h3 class="text-sm font-medium text-red-800">Error al escanear código QR</h3>
          <div class="mt-2 text-sm text-red-700">{{ error }}</div>
          <div class="mt-4">
            <button @click="clearError" class="btn-secondary text-sm">
              Intentar de nuevo
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Acciones rápidas -->
    <div class="bg-gray-50 rounded-lg p-6">
      <h3 class="font-medium text-gray-900 mb-4 text-center">Acciones Rápidas</h3>
      <div class="grid md:grid-cols-4 gap-4">
        <router-link to="/inventory" class="btn-secondary p-4 text-center block hover:bg-gray-200 transition-colors">
          <svg class="h-6 w-6 mx-auto mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
          </svg>
          <div class="text-sm font-medium">Ver Inventario</div>
        </router-link>
        
        <router-link to="/add-supply" class="btn-success p-4 text-center block hover:bg-green-700 transition-colors">
          <svg class="h-6 w-6 mx-auto mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
          </svg>
          <div class="text-sm font-medium">Agregar Insumo</div>
        </router-link>
        
        <router-link to="/" class="btn-secondary p-4 text-center block hover:bg-gray-200 transition-colors">
          <svg class="h-6 w-6 mx-auto mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
          </svg>
          <div class="text-sm font-medium">Ir al Inicio</div>
        </router-link>
        
        <button @click="clearAll" class="btn-secondary p-4 hover:bg-gray-200 transition-colors">
          <svg class="h-6 w-6 mx-auto mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
          <div class="text-sm font-medium">Limpiar</div>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import qrService from '@/services/qrService'
import QRInfoDisplay from '@/components/QRInfoDisplay.vue'

const router = useRouter()
const route = useRoute()

// Estado reactivo
const loading = ref(false)
const error = ref(null)
const qrInput = ref('')
const scannedInfo = ref(null)

// Método de escaneo
const scanQRCode = async () => {
  if (!qrInput.value.trim()) return
  
  loading.value = true
  error.value = null
  scannedInfo.value = null
  
  try {
    const result = await qrService.scanQRCode(qrInput.value.trim())
    if (result.success) {
      scannedInfo.value = result.data
    } else {
      error.value = result.error || 'Error desconocido al escanear código QR'
    }
  } catch (err) {
    error.value = err.response?.data?.error || 'Error de conexión al escanear código QR'
  } finally {
    loading.value = false
  }
}

// Métodos de utilidad
const viewDetails = (qrInfo) => {
  router.push({
    name: 'QRDetails',
    params: { qrcode: qrInfo.qr_code },
    state: { qrInfo }
  })
}

const clearError = () => {
  error.value = null
}

const clearAll = () => {
  qrInput.value = ''
  scannedInfo.value = null
  error.value = null
}

// Inicializar con QR de prueba si viene en query params
onMounted(() => {
  if (route.query.test) {
    qrInput.value = route.query.test
    scanQRCode()
  }
})
</script>

<style scoped>
/* Efectos hover mejorados */
.card:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
  transition: all 0.2s ease-in-out;
}

/* Colores específicos */
.text-blue-700 {
  color: #1d4ed8;
}

.text-green-700 {
  color: #15803d;
}
</style>