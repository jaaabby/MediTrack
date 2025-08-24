<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="bg-white rounded-lg shadow-sm border p-6">
      <h2 class="text-2xl font-semibold text-gray-900">Escáner de Códigos QR</h2>
      <p class="text-gray-600 mt-1">Escanea códigos QR de lotes o productos individuales para ver su información completa</p>
    </div>

    <!-- Scanner Input -->
    <div class="bg-white rounded-lg shadow-sm border p-6">
      <h3 class="text-lg font-medium text-gray-900 mb-4 flex items-center">
        <svg class="h-5 w-5 mr-2 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h2M4 4h5l2 3h3l2-3h5v5M4 4v5m16-5v5" />
        </svg>
        Escanear Código QR
      </h3>

      <div class="grid lg:grid-cols-2 gap-6">
        <!-- Manual Input -->
        <div>
          <label for="qr-input" class="block text-sm font-medium text-gray-700 mb-2">
            Ingrese el Código QR:
          </label>
          <div class="flex space-x-3">
            <input
              id="qr-input"
              v-model="qrInput"
              type="text"
              placeholder="BATCH_1755580808_abc123 o SUPPLY_1755580808_def456"
              class="form-input flex-1"
              @keyup.enter="scanQRCode"
              @paste="handlePaste"
            />
            <button
              @click="scanQRCode"
              :disabled="!qrInput.trim() || loading"
              class="btn-primary"
            >
              <svg v-if="loading" class="animate-spin h-4 w-4 mr-2" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              {{ loading ? 'Escaneando...' : 'Escanear' }}
            </button>
          </div>
          
          <!-- Format Helper -->
          <div class="mt-3 p-3 bg-blue-50 rounded-lg">
            <p class="text-sm text-blue-800">
              <svg class="h-4 w-4 inline mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <strong>Formatos válidos:</strong>
            </p>
            <ul class="text-sm text-blue-700 mt-1 space-y-1">
              <li>• <code class="bg-blue-100 px-1 rounded">BATCH_...</code> - Para lotes completos</li>
              <li>• <code class="bg-blue-100 px-1 rounded">SUPPLY_...</code> - Para productos individuales</li>
            </ul>
          </div>
        </div>

        <!-- Camera Scanner -->
        <div class="space-y-4">
          <label class="block text-sm font-medium text-gray-700">
            Escanear con Cámara:
          </label>
          
          <div class="relative">
            <!-- Camera View -->
            <div 
              v-if="cameraActive" 
              class="bg-gray-900 rounded-lg overflow-hidden aspect-video flex items-center justify-center relative"
            >
              <!-- Indicador de estado superior -->
              <div v-if="detecting" class="absolute top-2 left-1/2 transform -translate-x-1/2 bg-blue-500 text-white px-3 py-1 rounded-full text-xs font-medium z-10">
                <div class="flex items-center space-x-2">
                  <div class="w-2 h-2 bg-white rounded-full animate-pulse"></div>
                  <span>Detectando QR...</span>
                </div>
              </div>
              
              <video 
                ref="videoElement" 
                autoplay 
                muted 
                playsinline
                class="w-full h-full object-cover transform scale-x-[-1]"
              ></video>
              
              <!-- Overlay con indicador de detección -->
              <div class="absolute inset-0 border-2 border-blue-500 rounded-lg">
                <!-- Marco de detección -->
                <div class="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 w-32 h-32 border-2 border-white opacity-50 rounded"></div>
                
                <!-- Indicador de escaneo -->
                <div class="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 w-32 h-32">
                  <div class="absolute inset-0 border-2 border-blue-400 rounded animate-pulse"></div>
                  <div class="absolute inset-0 border-2 border-blue-300 rounded animate-ping"></div>
                </div>
                
                <!-- Texto de estado -->
                <div class="absolute bottom-16 left-1/2 transform -translate-x-1/2 bg-black bg-opacity-50 text-white px-3 py-1 rounded text-sm">
                  {{ detecting ? 'Detectando...' : 'Posicione el código QR en el marco' }}
                </div>
              </div>
              
              <!-- Controls -->
              <div class="absolute bottom-4 left-1/2 transform -translate-x-1/2 flex space-x-2">
                <button @click="stopCamera" class="bg-red-600 hover:bg-red-700 text-white px-3 py-1 rounded text-sm">
                  Detener
                </button>
              </div>
            </div>
            
            <!-- Camera Button -->
            <button
              v-else
              @click="startCamera"
              class="w-full h-32 border-2 border-dashed border-gray-300 hover:border-blue-400 rounded-lg flex flex-col items-center justify-center transition-colors"
              :disabled="cameraStarting"
            >
              <svg class="h-8 w-8 text-gray-400 mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" />
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z" />
              </svg>
              <span class="text-sm text-gray-600">
                {{ cameraStarting ? 'Iniciando Cámara...' : 'Activar Cámara para Escanear' }}
              </span>
            </button>
          </div>
          
          <!-- Camera Error -->
          <div v-if="cameraError" class="text-sm text-red-600 bg-red-50 p-2 rounded">
            {{ cameraError }}
          </div>
          
          <!-- Camera Instructions -->
          <div class="text-sm text-gray-600 bg-gray-50 p-3 rounded-lg">
            <p class="font-medium mb-1">Instrucciones:</p>
            <ul class="space-y-1 text-xs">
              <li>• Posicione el código QR dentro del marco azul</li>
              <li>• Mantenga el código estable y bien iluminado</li>
              <li>• La detección es automática</li>
              <li>• La cámara se detendrá automáticamente al detectar</li>
            </ul>
          </div>
        </div>
      </div>
    </div>

    <!-- Scan History -->
    <div v-if="scanHistory.length > 0" class="bg-white rounded-lg shadow-sm border p-6">
      <div class="flex items-center justify-between mb-4">
        <h3 class="text-lg font-medium text-gray-900 flex items-center">
          <svg class="h-5 w-5 mr-2 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          Historial de Escaneos
        </h3>
        <button @click="clearHistory" class="text-sm text-gray-500 hover:text-gray-700">
          Limpiar Historial
        </button>
      </div>
      
      <div class="space-y-2">
        <div 
          v-for="(item, index) in scanHistory" 
          :key="index"
          class="flex items-center justify-between p-3 bg-gray-50 rounded-lg hover:bg-gray-100 cursor-pointer transition-colors"
          @click="quickScan(item.qr_code)"
        >
          <div class="flex items-center space-x-3">
            <div :class="getTypeIconClass(item.type)" class="p-2 rounded-full">
              <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path v-if="item.type === 'batch'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
                <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2zM9 9h6v6H9V9z" />
              </svg>
            </div>
            <div>
              <p class="font-medium text-gray-900">{{ getTypeLabel(item.type) }}</p>
              <p class="text-sm text-gray-600">{{ item.qr_code }}</p>
            </div>
          </div>
          <div class="text-right">
            <p class="text-sm text-gray-500">{{ formatDate(item.scanned_at) }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- QR Information Display -->
    <div v-if="scannedInfo && !error">
      <QRInfoDisplay 
        :qr-info="scannedInfo" 
        @view-details="viewDetails" 
        @consume-supply="consumeSupply"
        @view-batch="viewBatch"
      />
    </div>

    <!-- Error Message -->
    <div v-if="error" class="bg-red-50 border border-red-200 rounded-lg p-6">
      <div class="flex items-start">
        <div class="flex-shrink-0">
          <svg class="h-6 w-6 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </div>
        <div class="ml-3 flex-1">
          <h3 class="text-lg font-medium text-red-800">Error al escanear código QR</h3>
          <div class="mt-2 text-sm text-red-700">
            <p>{{ error }}</p>
          </div>
          <div class="mt-4 space-x-3">
            <button @click="clearError" class="btn-secondary text-sm">
              Intentar de Nuevo
            </button>
            <button @click="validateQRFormat" class="btn-secondary text-sm">
              Validar Formato
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Quick Actions -->
    <div class="bg-gray-50 rounded-lg p-6">
      <h3 class="font-medium text-gray-900 mb-4 text-center">Acciones Rápidas</h3>
      <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
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
        
        <router-link to="/consume" class="btn-danger p-4 text-center block hover:bg-red-700 transition-colors">
          <svg class="h-6 w-6 mx-auto mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
          </svg>
          <div class="text-sm font-medium">Consumir</div>
        </router-link>
        
        <router-link to="/" class="btn-secondary p-4 text-center block hover:bg-gray-200 transition-colors">
          <svg class="h-6 w-6 mx-auto mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
          </svg>
          <div class="text-sm font-medium">Ir al Inicio</div>
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { format } from 'date-fns'
import { es } from 'date-fns/locale'
import qrService from '@/services/qrService'
import QRInfoDisplay from '@/components/QRInfoDisplay.vue'
import jsQR from 'jsqr'

const router = useRouter()
const route = useRoute()

// Referencias
const videoElement = ref(null)

// Estado reactivo
const loading = ref(false)
const error = ref(null)
const qrInput = ref('')
const scannedInfo = ref(null)
const cameraActive = ref(false)
const cameraStarting = ref(false)
const cameraError = ref(null)
const scanHistory = ref([])
const detecting = ref(false)

let mediaStream = null
let animationFrameId = null
let canvas = null
let canvasContext = null

// Métodos principales
const scanQRCode = async () => {
  if (!qrInput.value.trim()) return
  
  loading.value = true
  error.value = null
  scannedInfo.value = null
  
  try {
    const result = await qrService.scanQRCode(qrInput.value.trim())
    
    if (result.success) {
      scannedInfo.value = result.data
      
      // Agregar al historial
      addToHistory(result.data)
      
    } else {
      error.value = result.error || 'Error desconocido al escanear código QR'
    }
  } catch (err) {
    console.error('Error scanning QR:', err)
    error.value = err.response?.data?.error || 'Error de conexión al escanear código QR'
  } finally {
    loading.value = false
  }
}

const quickScan = async (qrCode) => {
  qrInput.value = qrCode
  await scanQRCode()
}

// Manejo de cámara
const startCamera = async () => {
  cameraStarting.value = true
  cameraError.value = null
  
  try {
    const stream = await navigator.mediaDevices.getUserMedia({ 
      video: { 
        facingMode: 'environment', // Usar cámara trasera si está disponible
        width: { ideal: 1280 },
        height: { ideal: 720 }
      } 
    })
    
    mediaStream = stream
    cameraActive.value = true
    
    // Esperar a que el elemento video esté disponible
    await new Promise(resolve => setTimeout(resolve, 100))
    
    if (videoElement.value) {
      videoElement.value.srcObject = stream
      
      // Crear canvas para procesar frames
      canvas = document.createElement('canvas')
      canvasContext = canvas.getContext('2d')
      
      // Iniciar detección QR
      startQRDetection()
    }
    
  } catch (err) {
    console.error('Error accessing camera:', err)
    cameraError.value = 'No se puede acceder a la cámara. Verifique los permisos.'
  } finally {
    cameraStarting.value = false
  }
}

const startQRDetection = () => {
  detecting.value = true
  
  const detectQR = () => {
    if (!videoElement.value || !canvas || !canvasContext || !cameraActive.value) {
      detecting.value = false
      return
    }
    
    try {
      // Configurar canvas con las dimensiones del video
      canvas.width = videoElement.value.videoWidth
      canvas.height = videoElement.value.videoHeight
      
      // Dibujar frame del video en el canvas
      canvasContext.drawImage(videoElement.value, 0, 0, canvas.width, canvas.height)
      
      // Obtener datos de imagen para procesar
      const imageData = canvasContext.getImageData(0, 0, canvas.width, canvas.height)
      
      // Detectar código QR
      const code = jsQR(imageData.data, imageData.width, imageData.height, {
        inversionAttempts: "dontInvert",
      })
      
      if (code) {
        console.log('Código QR detectado:', code.data)
        detecting.value = false
        
        // Procesar el código detectado
        handleDetectedQR(code.data)
        
        // Detener la cámara después de detectar
        stopCamera()
        return
      }
      
      // Continuar detectando si no se encontró código
      if (cameraActive.value) {
        animationFrameId = requestAnimationFrame(detectQR)
      } else {
        detecting.value = false
      }
      
    } catch (error) {
      console.error('Error en detección QR:', error)
      detecting.value = false
      // Continuar detectando incluso si hay error
      if (cameraActive.value) {
        animationFrameId = requestAnimationFrame(detectQR)
      }
    }
  }
  
  // Iniciar detección
  detectQR()
}

const handleDetectedQR = async (qrData) => {
  // Limpiar input y establecer el código detectado
  qrInput.value = qrData
  
  // Mostrar notificación de éxito
  showDetectionSuccess(qrData)
  
  // Auto-escanear después de un breve delay
  setTimeout(() => {
    scanQRCode()
  }, 500)
}

const showDetectionSuccess = (qrData) => {
  // Crear notificación temporal
  const notification = document.createElement('div')
  notification.className = 'fixed top-4 right-4 bg-green-500 text-white px-6 py-3 rounded-lg shadow-lg z-50 transform transition-all duration-300'
  notification.innerHTML = `
    <div class="flex items-center space-x-2">
      <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
      </svg>
      <span>Código QR detectado: ${qrData.substring(0, 20)}...</span>
    </div>
  `
  
  document.body.appendChild(notification)
  
  // Remover notificación después de 3 segundos
  setTimeout(() => {
    notification.remove()
  }, 3000)
}

const stopCamera = () => {
  if (mediaStream) {
    mediaStream.getTracks().forEach(track => track.stop())
    mediaStream = null
  }
  
  if (animationFrameId) {
    cancelAnimationFrame(animationFrameId)
    animationFrameId = null
  }
  
  cameraActive.value = false
  cameraError.value = null
  detecting.value = false
  
  if (videoElement.value) {
    videoElement.value.srcObject = null
  }
  
  // Limpiar canvas
  if (canvas) {
    canvas = null
    canvasContext = null
  }
}

// Utilidades
const handlePaste = (event) => {
  // Auto-escanear después de pegar
  setTimeout(() => {
    if (qrInput.value.trim()) {
      scanQRCode()
    }
  }, 100)
}

const addToHistory = (qrInfo) => {
  const historyItem = {
    qr_code: qrInfo.qr_code,
    type: qrInfo.type,
    scanned_at: new Date().toISOString()
  }
  
  // Evitar duplicados
  const exists = scanHistory.value.some(item => item.qr_code === qrInfo.qr_code)
  if (!exists) {
    scanHistory.value.unshift(historyItem)
    
    // Mantener solo los últimos 20
    if (scanHistory.value.length > 20) {
      scanHistory.value = scanHistory.value.slice(0, 20)
    }
    
    // Guardar en localStorage
    localStorage.setItem('qr_scan_history', JSON.stringify(scanHistory.value))
  }
}

const loadHistory = () => {
  try {
    const saved = localStorage.getItem('qr_scan_history')
    if (saved) {
      scanHistory.value = JSON.parse(saved)
    }
  } catch (error) {
    console.error('Error loading scan history:', error)
  }
}

const clearHistory = () => {
  scanHistory.value = []
  localStorage.removeItem('qr_scan_history')
}

const validateQRFormat = () => {
  if (!qrInput.value.trim()) return
  
  const isValid = qrService.isValidQRFormat(qrInput.value.trim())
  
  if (isValid) {
    error.value = null
    alert('El formato del código QR es válido.')
  } else {
    error.value = 'El formato del código QR no es válido. Debe tener el formato: BATCH_timestamp_random o SUPPLY_timestamp_random'
  }
}

// Navegación
const viewDetails = (qrInfo) => {
  router.push({
    name: 'QRDetails',
    params: { qrcode: qrInfo.qr_code }
  })
}

const consumeSupply = (qrInfo) => {
  router.push({
    name: 'QRConsumer',
    query: { qr: qrInfo.qr_code }
  })
}

const viewBatch = (batchId) => {
  router.push({
    name: 'Inventory',
    query: { batch: batchId }
  })
}

// Utilidades de UI
const getTypeLabel = (type) => {
  return qrService.getTypeLabel(type)
}

const getTypeIconClass = (type) => {
  const classes = {
    'batch': 'bg-blue-100 text-blue-600',
    'medical_supply': 'bg-green-100 text-green-600'
  }
  return classes[type] || 'bg-gray-100 text-gray-600'
}

const formatDate = (dateString) => {
  if (!dateString) return 'No disponible'
  try {
    return format(new Date(dateString), 'dd/MM HH:mm', { locale: es })
  } catch (error) {
    return dateString
  }
}

const clearError = () => {
  error.value = null
}

const clearAll = () => {
  qrInput.value = ''
  scannedInfo.value = null
  error.value = null
  stopCamera()
}

// Lifecycle
onMounted(() => {
  loadHistory()
  
  // Auto-escanear si viene QR en query params
  if (route.query.qr || route.query.test) {
    const testQR = route.query.qr || route.query.test
    qrInput.value = testQR
    scanQRCode()
  }
})

onUnmounted(() => {
  stopCamera()
})
</script>

