<template>
  <div class="max-w-4xl mx-auto p-6">
    <!-- Header con estadísticas 
    <div class="mb-6">
      <h1 class="text-2xl font-bold text-gray-900 mb-2">Escáner QR con Trazabilidad</h1>
      <div class="grid grid-cols-1 md:grid-cols-4 gap-4 mb-4">
        <div class="bg-blue-50 border border-blue-200 rounded-lg p-3">
          <div class="text-xs font-medium text-blue-600 uppercase">Escaneos Hoy</div>
          <div class="text-xl font-bold text-blue-900">{{ todayScans }}</div>
        </div>
        <div class="bg-green-50 border border-green-200 rounded-lg p-3">
          <div class="text-xs font-medium text-green-600 uppercase">Exitosos</div>
          <div class="text-xl font-bold text-green-900">{{ successfulScans }}</div>
        </div>
        <div class="bg-purple-50 border border-purple-200 rounded-lg p-3">
          <div class="text-xs font-medium text-purple-600 uppercase">Ubicación</div>
          <div class="text-sm font-medium text-purple-900">{{ currentLocation?.name || 'No seleccionada' }}</div>
        </div>
        <div class="bg-gray-50 border border-gray-200 rounded-lg p-3">
          <div class="text-xs font-medium text-gray-600 uppercase">Usuario</div>
          <div class="text-sm font-medium text-gray-900">{{ currentUser?.name || 'No logueado' }}</div>
        </div>
      </div>
    </div>-->

    <!-- Selector de propósito de escaneo 
    <div class="bg-white rounded-lg shadow-sm border p-6 mb-6">
      <h3 class="text-lg font-medium text-gray-900 mb-4">Propósito del Escaneo</h3>
      <div class="grid grid-cols-2 md:grid-cols-5 gap-3">
        <button
          v-for="purpose in scanPurposes"
          :key="purpose.value"
          @click="selectedPurpose = purpose.value"
          :class="[
            'p-3 rounded-lg border-2 text-sm font-medium transition-all',
            selectedPurpose === purpose.value
              ? 'border-blue-500 bg-blue-50 text-blue-700'
              : 'border-gray-200 bg-white text-gray-700 hover:border-gray-300 hover:bg-gray-50'
          ]"
        >
          <div class="flex flex-col items-center">
            <svg class="h-6 w-6 mb-1" :class="purpose.iconClass" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" :d="purpose.icon" />
            </svg>
            <span>{{ purpose.label }}</span>
          </div>
        </button>
      </div>
    </div>-->

    <!-- Scanner Input -->
    <div class="bg-white rounded-lg shadow-sm border p-6 mb-6">
      <h3 class="text-lg font-medium text-gray-900 mb-4 flex items-center">
        <svg class="h-5 w-5 mr-2 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h2M4 4h5l2 3h3l2-3h5v5M4 4v5m16-5v5" />
        </svg>
        Escanear Insumo Individual
      </h3>

      <div class="grid lg:grid-cols-2 gap-6">
        <!-- Manual Input -->
        <div>
          <label for="qr-input" class="block text-sm font-medium text-gray-700 mb-2">
            Ingrese el Código QR del Insumo:
          </label>
          <div class="flex space-x-3">
            <input
              id="qr-input"
              v-model="qrInput"
              type="text"
              placeholder="SUPPLY_1755580808_def456"
              class="form-input flex-1"
              @keyup.enter="scanQRCode"
              @paste="handlePaste"
              :disabled="loading"
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
          <div class="mt-3 p-3 bg-green-50 rounded-lg">
            <p class="text-sm text-green-800">
              <svg class="h-4 w-4 inline mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <strong>Prioridad: Insumos Individuales</strong>
            </p>
            <ul class="text-sm text-green-700 mt-1 space-y-1">
              <li>• <code class="bg-green-100 px-1 rounded">SUPPLY_...</code> - Insumos individuales (trazabilidad)</li>
              <li>• <code class="bg-blue-100 px-1 rounded">BATCH_...</code> - Información del lote (solo consulta)</li>
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
                <div class="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 w-48 h-48 border-2 border-white opacity-50 rounded"></div>
                
                <!-- Indicador de escaneo -->
                <div class="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 w-48 h-48">
                  <div class="absolute inset-0 border-2 border-blue-400 rounded animate-pulse"></div>
                  <div class="absolute inset-0 border-2 border-blue-300 rounded animate-ping"></div>
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
              @click="startCameraScanner"
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

    <!-- Error Display -->
    <div v-if="error" class="bg-red-50 border border-red-200 rounded-lg p-4 mb-6">
      <div class="flex items-start space-x-3">
        <svg class="h-5 w-5 text-red-600 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <div class="flex-1">
          <h4 class="text-sm font-medium text-red-800">Error al escanear</h4>
          <p class="text-sm text-red-700 mt-1">{{ error }}</p>
          <button @click="clearError" class="text-sm text-red-600 hover:text-red-800 mt-2 underline">
            Limpiar error
          </button>
        </div>
      </div>
    </div>

    <!-- Scanned Supply Info Display -->
    <div v-if="scannedInfo && !error" class="bg-white rounded-lg shadow-sm border overflow-hidden mb-6">
      <QRInfoDisplay 
        :qr-info="scannedInfo"
        :show-traceability="true"
        :scan-context="lastScanContext"
        @view-details="viewDetails"
        @view-batch="viewBatch"
        @consume-supply="consumeSupply"
      />
      
      <!-- Botones de acción adicionales con trazabilidad -->
      <div class="p-4 border-t border-gray-200 bg-gray-50">
        <div class="flex flex-wrap gap-3">
          
          <button
            v-if="scannedInfo.supply_info && !scannedInfo.is_consumed"
            @click="quickConsume"
            class="btn-primary text-sm"
          >
            <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
            Consumir Rápido
          </button>
          
          <button
            v-if="scannedInfo.supply_info && !scannedInfo.is_consumed"
            @click="quickTransfer"
            class="btn-secondary text-sm"
          >
            <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4" />
            </svg>
            Transferir Rápido
          </button>
        </div>
      </div>
    </div>

    <!-- Scan History con información de trazabilidad -->
    <div v-if="scanHistory.length > 0" class="bg-white rounded-lg shadow-sm border p-6 mb-6">
      <div class="flex items-center justify-between mb-4">
        <h3 class="text-lg font-medium text-gray-900 flex items-center">
          <svg class="h-5 w-5 mr-2 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          Historial de Escaneos Recientes
        </h3>
        <button @click="clearHistory" class="text-sm text-gray-500 hover:text-gray-700">
          Limpiar Historial
        </button>
      </div>
      
      <div class="space-y-3 max-h-80 overflow-y-auto">
        <button 
          v-for="(item, index) in scanHistory" 
          :key="index"
          @click="quickRescan(item.qr_code)"
          class="w-full text-left p-3 bg-gray-50 rounded-lg hover:bg-blue-50 transition-colors flex flex-col"
        >
          <div class="flex items-center space-x-2">
            <span class="p-1.5 rounded-full bg-gray-100 text-gray-600">
              <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </span>
            <code class="text-sm font-mono text-gray-900 truncate">{{ item.qr_code }}</code>
            <span 
              :class="[
                'px-2 py-0.5 rounded-full text-xs font-medium',
                item.scan_purpose === 'consume' ? 'bg-red-100 text-red-800' :
                item.scan_purpose === 'verify' ? 'bg-blue-100 text-blue-800' :
                item.scan_purpose === 'lookup' ? 'bg-green-100 text-green-800' :
                item.scan_purpose === 'transfer' ? 'bg-orange-100 text-orange-800' :
                'bg-gray-100 text-gray-800'
              ]"
            >
              {{ item.scan_purpose || 'lookup' }}
            </span>
            <span :class="[
              'px-2 py-1 text-xs font-medium rounded',
              item.type === 'medical_supply' ? 'bg-green-100 text-green-700' : 'bg-blue-100 text-blue-700'
            ]">
              {{ item.type === 'medical_supply' ? 'Insumo' : 'Lote' }}
            </span>
          </div>
          <div class="flex items-center space-x-4 text-xs text-gray-500 mt-1">
            <span>{{ formatDate(item.scanned_at) }}</span>
            <span v-if="item.location">📍 {{ item.location }}</span>
            <span v-if="item.user_name">👤 {{ item.user_name }}</span>
          </div>
          <div v-if="item.supply_name" class="text-sm text-gray-700 mt-1 truncate">
            {{ item.supply_name }}
          </div>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { format } from 'date-fns'
import { es } from 'date-fns/locale'
import { useAuthStore } from '@/stores/auth'
import qrService from '@/services/qrService'
import jsQR from 'jsqr'
import QRInfoDisplay from '@/components/QRInfoDisplay.vue'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

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
const selectedPurpose = ref('lookup')
const currentLocation = ref(null)
const lastScanContext = ref(null)

// Computed properties para datos del usuario desde el store
const currentUser = computed(() => authStore.user)
const isAuthenticated = computed(() => authStore.isAuthenticated)

// Estadísticas
const todayScans = ref(0)
const successfulScans = ref(0)

let mediaStream = null
let animationFrameId = null

// Propósitos de escaneo
const scanPurposes = [
  {
    value: 'lookup',
    label: 'Consultar',
    icon: 'M15 12a3 3 0 11-6 0 3 3 0 016 0z M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z',
    iconClass: 'text-green-600'
  },
  {
    value: 'consume',
    label: 'Consumir',
    icon: 'M5 13l4 4L19 7',
    iconClass: 'text-red-600'
  },
  {
    value: 'verify',
    label: 'Verificar',
    icon: 'M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z',
    iconClass: 'text-blue-600'
  },
  {
    value: 'transfer',
    label: 'Transferir',
    icon: 'M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4',
    iconClass: 'text-orange-600'
  },
  {
    value: 'inventory_check',
    label: 'Inventario',
    icon: 'M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2',
    iconClass: 'text-purple-600'
  }
]

// Funciones de cámara con efectos visuales originales
const startCameraScanner = async () => {
  if (cameraActive.value) return
  
  cameraStarting.value = true
  cameraError.value = null
  
  try {
    const stream = await navigator.mediaDevices.getUserMedia({ 
      video: { 
        facingMode: 'environment',
        width: { ideal: 1280 },
        height: { ideal: 720 }
      } 
    })
    
    mediaStream = stream
    cameraActive.value = true
    
    // Mostrar toast de activación con efectos
    showDetectionToast('Cámara activada', 'success')
    
    // Esperar a que el elemento video esté disponible
    await new Promise(resolve => setTimeout(resolve, 100))
    
    if (videoElement.value) {
      videoElement.value.srcObject = stream
      
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
  if (!cameraActive.value || !videoElement.value) {
    console.log('Detección cancelada: cámara no activa o elemento no disponible')
    return
  }
  
  console.log('🔍 Iniciando detección QR con efectos visuales...')
  detecting.value = true
  
  const detectQR = () => {
    if (!cameraActive.value || !videoElement.value) {
      detecting.value = false
      console.log('Detección detenida: cámara no activa')
      return
    }
    
    try {
      // Verificar que el video tiene dimensiones válidas antes de procesar
      if (videoElement.value.videoWidth === 0 || videoElement.value.videoHeight === 0) {
        // El video aún no está listo, intentar de nuevo en el siguiente frame
        animationFrameId = requestAnimationFrame(detectQR)
        return
      }
      
      // Crear canvas localmente (no como variable global)
      const canvas = document.createElement('canvas')
      const context = canvas.getContext('2d')
      const video = videoElement.value
      
      canvas.width = video.videoWidth
      canvas.height = video.videoHeight
      
      // Dibujar el frame actual del video en el canvas
      context.drawImage(video, 0, 0, canvas.width, canvas.height)
      
      // Obtener los datos de la imagen
      const imageData = context.getImageData(0, 0, canvas.width, canvas.height)
      
      // Intentar detectar QR
      const qrCode = jsQR(imageData.data, imageData.width, imageData.height, {
        inversionAttempts: "dontInvert"
      })
      
      if (qrCode && qrCode.data) {
        const qrData = qrCode.data.trim()
        console.log('QR detectado:', qrData)
        
        // Verificar que sea un código QR válido (SUPPLY_ o BATCH_)
        if (qrData.startsWith('SUPPLY_') || qrData.startsWith('BATCH_')) {
          console.log('✅ QR válido detectado:', qrData)
          
          // Detener la detección
          detecting.value = false
          
          // Establecer el valor en el input
          qrInput.value = qrData
          
          // Detener la cámara
          stopCamera()
          
          // Mostrar notificación de éxito con efectos mejorados
          showDetectionSuccess(qrData)
          
          // Auto-escanear después de un breve delay
          setTimeout(() => {
            scanQRCode()
          }, 300)
          
          return
        } else {
          console.log('QR detectado pero formato inválido:', qrData)
        }
      }
      
    } catch (canvasError) {
      console.warn('Error en procesamiento de canvas:', canvasError)
      // No detener la detección por errores de canvas
    }
    
    // Continuar detectando si llegamos hasta aquí
    if (cameraActive.value) {
      animationFrameId = requestAnimationFrame(detectQR)
    } else {
      detecting.value = false
    }
  }
  
  // Iniciar el bucle de detección
  detectQR()
}

const showDetectionSuccess = (qrData) => {
  // Crear notificación mejorada con más efectos
  const notification = document.createElement('div')
  notification.className = 'fixed top-4 right-4 z-50 transform transition-all duration-500'
  notification.innerHTML = `
    <div class="bg-gradient-to-r from-green-500 to-emerald-500 text-white px-6 py-4 rounded-lg shadow-2xl">
      <div class="flex items-center space-x-3">
        <div class="flex-shrink-0">
          <div class="w-8 h-8 bg-white bg-opacity-20 rounded-full flex items-center justify-center animate-bounce">
            <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
          </div>
        </div>
        <div class="flex-1">
          <h4 class="text-sm font-bold">QR Detectado con Éxito</h4>
          <p class="text-xs opacity-90 font-mono">${qrData}</p>
        </div>
        <div class="flex-shrink-0">
          <div class="w-2 h-2 bg-white rounded-full animate-ping"></div>
        </div>
      </div>
    </div>
  `
  
  // Animar entrada
  notification.style.transform = 'translateX(100%) scale(0.8)'
  notification.style.opacity = '0'
  
  document.body.appendChild(notification)
  
  setTimeout(() => {
    notification.style.transform = 'translateX(0) scale(1)'
    notification.style.opacity = '1'
  }, 100)
  
  // Remover notificación después de 3 segundos con animación
  setTimeout(() => {
    notification.style.transform = 'translateX(100%) scale(0.8)'
    notification.style.opacity = '0'
    setTimeout(() => {
      if (notification.parentNode) {
        notification.remove()
      }
    }, 500)
  }, 3000)
}

const showDetectionToast = (message, type = 'info') => {
  const toast = document.createElement('div')
  const colors = {
    success: 'from-green-500 to-emerald-500',
    error: 'from-red-500 to-rose-500', 
    warning: 'from-yellow-500 to-amber-500',
    info: 'from-blue-500 to-indigo-500'
  }
  
  toast.className = `fixed top-4 left-1/2 transform -translate-x-1/2 z-50 transition-all duration-300`
  toast.innerHTML = `
    <div class="bg-gradient-to-r ${colors[type]} text-white px-6 py-3 rounded-full shadow-lg">
      <div class="flex items-center space-x-2">
        <div class="w-2 h-2 bg-white rounded-full animate-pulse"></div>
        <span class="text-sm font-medium">${message}</span>
      </div>
    </div>
  `
  
  toast.style.transform = 'translateX(-50%) translateY(-100%)'
  toast.style.opacity = '0'
  
  document.body.appendChild(toast)
  
  setTimeout(() => {
    toast.style.transform = 'translateX(-50%) translateY(0)'
    toast.style.opacity = '1'
  }, 100)
  
  setTimeout(() => {
    toast.style.transform = 'translateX(-50%) translateY(-100%)'
    toast.style.opacity = '0'
    setTimeout(() => {
      if (toast.parentNode) {
        toast.remove()
      }
    }, 300)
  }, 2500)
}

const stopCamera = () => {
  console.log('🛑 Deteniendo cámara...')
  
  // Detener todos los tracks del stream
  if (mediaStream) {
    mediaStream.getTracks().forEach(track => {
      track.stop()
      console.log('Track detenido:', track.kind)
    })
    mediaStream = null
  }
  
  // Cancelar la animación
  if (animationFrameId) {
    cancelAnimationFrame(animationFrameId)
    animationFrameId = null
  }
  
  // Limpiar el elemento video
  if (videoElement.value) {
    videoElement.value.srcObject = null
  }
  
  // Actualizar estado
  cameraActive.value = false
  detecting.value = false
  cameraError.value = null
  
  console.log('Cámara detenida correctamente')
}

// Funciones principales de escaneo
const scanQRCode = async () => {
  if (!qrInput.value.trim() || loading.value) return
  
  loading.value = true
  error.value = null
  scannedInfo.value = null
  
  try {
    // Validar formato
    if (!isValidQRFormat(qrInput.value.trim())) {
      throw new Error('Formato de código QR inválido. Use SUPPLY_... para insumos individuales o BATCH_... para lotes.')
    }
    
    // Crear contexto de escaneo con trazabilidad
    const scanContext = buildScanContext()
    
    // Usar el método de escaneo con trazabilidad automática
    const result = await qrService.scanQRCode(qrInput.value.trim(), scanContext)
    
    scannedInfo.value = result
    lastScanContext.value = scanContext
    
    // Actualizar estadísticas
    updateScanStatistics(result, true)
    
    // Añadir al historial con información de trazabilidad
    addToHistory(qrInput.value.trim(), result.type, result, scanContext)
    
    // Mostrar mensaje diferente según el tipo
    if (result.type === 'medical_supply') {
      showDetectionToast('Insumo individual procesado', 'success')
    } else if (result.type === 'batch') {
      showDetectionToast('Lote detectado - Use QR individuales para trazabilidad', 'info')
    }
    
  } catch (err) {
    console.error('Error scanning QR:', err)
    error.value = err.response?.data?.error || err.message || 'Error al escanear el código QR'
    showDetectionToast('Error al escanear QR', 'error')
    updateScanStatistics(null, false)
  } finally {
    loading.value = false
  }
}

const buildScanContext = () => {
  return {
    scan_purpose: selectedPurpose.value,
    pavilion_id: currentLocation.value?.pavilion_id,
    medical_center_id: currentLocation.value?.medical_center_id,
    scan_source: 'web',
    user_agent: navigator.userAgent,
    device_info: {
      platform: navigator.platform,
      language: navigator.language,
      screen_resolution: `${screen.width}x${screen.height}`
    }
  }
}

const updateScanStatistics = (result, success) => {
  if (success) {
    successfulScans.value++
  }
  todayScans.value++
  
  // Guardar estadísticas locales
  const stats = {
    today_scans: todayScans.value,
    successful_scans: successfulScans.value,
    last_updated: new Date().toISOString()
  }
  localStorage.setItem('qr-scan-stats', JSON.stringify(stats))
}

const handlePaste = (event) => {
  setTimeout(() => {
    if (qrInput.value.trim()) {
      scanQRCode()
    }
  }, 100)
}

const isValidQRFormat = (qrCode) => {
  if (!qrCode || typeof qrCode !== 'string') return false
  const qrPattern = /^(BATCH|SUPPLY)_\d+_[a-f0-9]+$/i
  return qrPattern.test(qrCode)
}

const quickRescan = (qrCode) => {
  qrInput.value = qrCode
  scanQRCode()
}

// Funciones de historial mejoradas
const addToHistory = (qrCode, type, qrInfo = null, scanContext = null) => {
  const existing = scanHistory.value.findIndex(item => item.qr_code === qrCode)
  
  const historyItem = {
    qr_code: qrCode,
    type: type,
    scanned_at: new Date(),
    scan_purpose: scanContext?.scan_purpose || 'lookup',
    location: currentLocation.value?.name,
    user_name: currentUser.value?.name,
    supply_name: qrInfo?.supply_info?.supply_code_name,
    success: true,
    scan_context: scanContext
  }
  
  if (existing >= 0) {
    scanHistory.value.splice(existing, 1)
  }
  
  scanHistory.value.unshift(historyItem)
  
  // Mantener solo los últimos 20 escaneos
  if (scanHistory.value.length > 20) {
    scanHistory.value = scanHistory.value.slice(0, 20)
  }
  
  saveHistory()
}

const loadHistory = () => {
  try {
    const saved = localStorage.getItem('qr-scan-history')
    if (saved) {
      const parsed = JSON.parse(saved)
      scanHistory.value = parsed.map(item => ({
        ...item,
        scanned_at: new Date(item.scanned_at)
      }))
    }
    
    // Cargar estadísticas
    const stats = localStorage.getItem('qr-scan-stats')
    if (stats) {
      const parsed = JSON.parse(stats)
      todayScans.value = parsed.today_scans || 0
      successfulScans.value = parsed.successful_scans || 0
    }
  } catch (error) {
    console.error('Error loading scan history:', error)
  }
}

const saveHistory = () => {
  try {
    localStorage.setItem('qr-scan-history', JSON.stringify(scanHistory.value))
  } catch (error) {
    console.error('Error saving scan history:', error)
  }
}

const clearHistory = () => {
  scanHistory.value = []
  todayScans.value = 0
  successfulScans.value = 0
  saveHistory()
  localStorage.removeItem('qr-scan-stats')
}

// Funciones de navegación y acciones con trazabilidad
const viewCompleteTraceability = () => {
  if (scannedInfo.value) {
    router.push({
      name: 'QRTraceability',
      params: { qrcode: scannedInfo.value.qr_code }
    })
  }
}

const quickConsume = () => {
  if (scannedInfo.value) {
    router.push({
      name: 'QRConsumer',
      query: { 
        qr: scannedInfo.value.qr_code,
        quick: 'true',
        purpose: 'consume'
      }
    })
  }
}

const quickTransfer = () => {
  if (scannedInfo.value) {
    router.push({
      name: 'QRConsumer',
      query: { 
        qr: scannedInfo.value.qr_code,
        quick: 'true',
        purpose: 'consume',
        consumption_purpose: 'transfer'
      }
    })
  }
}

const viewScanStatistics = () => {
  if (scannedInfo.value) {
    router.push({
      name: 'QRAnalytics',
      params: { qrcode: scannedInfo.value.qr_code }
    })
  }
}

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

const onLocationChanged = (location) => {
  currentLocation.value = location
}

// Utilidades
const formatDate = (dateString) => {
  if (!dateString) return 'No disponible'
  try {
    return format(new Date(dateString), 'dd/MM/yyyy HH:mm', { locale: es })
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
  
  // Inicializar autenticación si no está inicializada
  if (!authStore.isAuthenticated) {
    authStore.initializeAuth()
  }
})

onUnmounted(() => {
  console.log('Componente desmontándose, limpiando cámara...')
  stopCamera()
})
</script>

<style scoped>
.form-input {
  @apply block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm;
}

.btn-primary {
  @apply inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed;
}

.btn-secondary {
  @apply inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md shadow-sm text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed;
}

.btn-danger {
  @apply inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-red-600 hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 disabled:opacity-50 disabled:cursor-not-allowed;
}

/* ===== ESTILOS SIMPLES PARA LA CÁMARA ===== */

/* Efecto espejo para la cámara */
.mirror-effect {
  transform: scaleX(-1);
}

/* Contenedor del scanner */
.scanner-container {
  position: relative;
  overflow: hidden;
}
</style>