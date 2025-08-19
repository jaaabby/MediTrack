<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="text-center">
      <h1 class="text-3xl font-bold text-gray-900 mb-2">Escáner de Códigos QR</h1>
      <p class="text-gray-600">Escanea o ingresa el código QR para obtener información del insumo médico</p>
    </div>

    <!-- Métodos de escaneo -->
    <div class="grid md:grid-cols-2 gap-6">
      <!-- Escáner manual -->
      <div class="card">
        <div class="card-header">
          <h2 class="card-title flex items-center">
            <svg class="h-6 w-6 mr-2 text-primary-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
            </svg>
            Ingreso Manual
          </h2>
        </div>
        <div class="space-y-4">
          <div>
            <label for="manual-qr" class="block text-sm font-medium text-gray-700 mb-2">
              Código QR
            </label>
            <input
              id="manual-qr"
              type="text"
              placeholder="Ej: BATCH_1723456789_abc123"
              class="form-input w-full"
              v-model="manualQRCode"
              @keyup.enter="scanManualQR"
              :disabled="loading"
            />
          </div>
          <button 
            @click="scanManualQR"
            :disabled="!manualQRCode.trim() || loading"
            class="btn-primary w-full"
          >
            <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
            {{ loading ? 'Escaneando...' : 'Escanear Código' }}
          </button>
        </div>
      </div>

      <!-- Escáner con cámara (simulado) -->
      <div class="card">
        <div class="card-header">
          <h2 class="card-title flex items-center">
            <svg class="h-6 w-6 mr-2 text-primary-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
            Escáner con Cámara
          </h2>
        </div>
        <div class="text-center space-y-4">
          <div class="bg-gray-100 rounded-lg p-8 border-2 border-dashed border-gray-300">
            <svg class="h-16 w-16 mx-auto text-gray-400 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h.01M5 8h2a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v2a2 2 0 002 2zm0 0h2a2 2 0 012 2v2a2 2 0 01-2 2H5a2 2 0 01-2-2v-2a2 2 0 012-2z" />
            </svg>
            <p class="text-gray-500 mb-4">Función de cámara próximamente</p>
          </div>
          <button class="btn-secondary w-full" disabled>
            <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" />
            </svg>
            Abrir Cámara
          </button>
        </div>
      </div>
    </div>

    <!-- Códigos QR de ejemplo -->
    <div class="card">
      <div class="card-header">
        <h2 class="card-title">Códigos QR de Ejemplo</h2>
        <p class="text-sm text-gray-600">Usa estos códigos para probar el sistema</p>
      </div>
      <div class="grid md:grid-cols-2 gap-4">
        <div class="space-y-2">
          <h3 class="font-medium text-gray-900">Códigos de Lotes (BATCH)</h3>
          <div class="space-y-2">
            <button 
              v-for="batch in exampleBatchCodes" 
              :key="batch"
              @click="scanExampleCode(batch)"
              class="w-full text-left p-3 bg-blue-50 border border-blue-200 rounded-lg hover:bg-blue-100 transition-colors"
            >
              <code class="text-sm text-blue-800">{{ batch }}</code>
            </button>
          </div>
        </div>
        <div class="space-y-2">
          <h3 class="font-medium text-gray-900">Códigos de Insumos (SUPPLY)</h3>
          <div class="space-y-2">
            <button 
              v-for="supply in exampleSupplyCodes" 
              :key="supply"
              @click="scanExampleCode(supply)"
              class="w-full text-left p-3 bg-green-50 border border-green-200 rounded-lg hover:bg-green-100 transition-colors"
            >
              <code class="text-sm text-green-800">{{ supply }}</code>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Resultados del escaneo -->
    <div v-if="scanResult" class="card">
      <div class="card-header">
        <h2 class="card-title flex items-center">
          <svg class="h-6 w-6 mr-2 text-success-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          Resultado del Escaneo
        </h2>
      </div>
      
      <QRInfoDisplay :qr-info="scanResult" @view-details="viewDetails" />
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
              Reintentar
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Acciones rápidas -->
    <div class="grid md:grid-cols-3 gap-4">
      <button @click="generateBatchQR" class="btn-secondary p-4 h-auto">
        <svg class="h-6 w-6 mx-auto mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
        </svg>
        <div class="text-sm font-medium">Generar QR de Lote</div>
      </button>
      
      <button @click="generateSupplyQR" class="btn-secondary p-4 h-auto">
        <svg class="h-6 w-6 mx-auto mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2zM9 9h6v6H9V9z" />
        </svg>
        <div class="text-sm font-medium">Generar QR de Insumo</div>
      </button>
      
      <router-link to="/inventory" class="btn-secondary p-4 h-auto text-center">
        <svg class="h-6 w-6 mx-auto mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
        </svg>
        <div class="text-sm font-medium">Ver Inventario</div>
      </router-link>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import qrService from '@/services/qrService'
import QRInfoDisplay from '@/components/QRInfoDisplay.vue'

const router = useRouter()

// Estado reactivo
const manualQRCode = ref('')
const loading = ref(false)
const error = ref(null)
const scanResult = ref(null)

// Códigos de ejemplo para testing
const exampleBatchCodes = [
  'BATCH_1723456789_abc123',
  'BATCH_1723456790_def456',
  'BATCH_1723456791_ghi789'
]

const exampleSupplyCodes = [
  'SUPPLY_1723456789_xyz123',
  'SUPPLY_1723456790_uvw456',
  'SUPPLY_1723456791_rst789'
]

// Métodos
const scanManualQR = async () => {
  if (!manualQRCode.value.trim()) return
  
  await scanQRCode(manualQRCode.value.trim())
}

const scanExampleCode = async (code) => {
  manualQRCode.value = code
  await scanQRCode(code)
}

const scanQRCode = async (qrCode) => {
  loading.value = true
  error.value = null
  scanResult.value = null
  
  try {
    const result = await qrService.scanQRCode(qrCode)
    if (result.success) {
      scanResult.value = result.data
    } else {
      error.value = result.error || 'Error desconocido al escanear código QR'
    }
  } catch (err) {
    error.value = err.response?.data?.error || 'Error de conexión al escanear código QR'
  } finally {
    loading.value = false
  }
}

const viewDetails = (qrInfo) => {
  // Navegar a la vista de detalles con la información del QR
  router.push({
    name: 'QRDetails',
    params: { qrcode: qrInfo.qr_code },
    state: { qrInfo }
  })
}

const generateBatchQR = async () => {
  try {
    const result = await qrService.generateBatchQR()
    if (result.success) {
      manualQRCode.value = result.data.qr_code
      await scanQRCode(result.data.qr_code)
    }
  } catch (err) {
    error.value = 'Error al generar código QR de lote'
  }
}

const generateSupplyQR = async () => {
  try {
    const result = await qrService.generateSupplyQR()
    if (result.success) {
      manualQRCode.value = result.data.qr_code
      await scanQRCode(result.data.qr_code)
    }
  } catch (err) {
    error.value = 'Error al generar código QR de insumo'
  }
}

const clearError = () => {
  error.value = null
}
</script>