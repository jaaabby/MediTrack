<template>
  <div class="max-w-4xl mx-auto p-3 sm:p-6">
    <!-- Header -->
    <div class="mb-4 sm:mb-6">
      <h1 class="text-xl sm:text-2xl font-bold text-gray-900">Escáner QR</h1>
      <p class="text-sm sm:text-base text-gray-600 mt-2">Escanea códigos QR para gestionar insumos médicos</p>
    </div>

    <!-- Sección de Escaneo -->
    <div class="bg-white rounded-lg shadow-sm border p-3 sm:p-6 mb-4 sm:mb-6">
      <h3 class="text-base sm:text-lg font-medium text-gray-900 mb-3 sm:mb-4 flex items-center">
        <svg class="h-5 w-5 mr-2 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h2M4 4h5l2 3h3l2-3h5v5M4 4v5m16-5v5" />
        </svg>
        Escanear Código QR
      </h3>

      <div class="grid grid-cols-1 lg:grid-cols-2 gap-4 sm:gap-6">
        <!-- Entrada Manual -->
        <div class="space-y-3 sm:space-y-4">
          <label for="qr-input" class="block text-xs sm:text-sm font-medium text-gray-700">
            Código QR del Insumo:
          </label>
          <div class="flex flex-col sm:flex-row gap-2 sm:gap-3">
            <input
              id="qr-input"
              v-model="qrInput"
              type="text"
              placeholder="SUPPLY_1755580808_def456"
              class="form-input flex-1 text-sm"
              @keyup.enter="scanQRCode"
              @paste="handlePaste"
              :disabled="loading"
            />
            <button
              @click="scanQRCode"
              :disabled="!qrInput.trim() || loading"
              class="btn-primary w-full sm:w-auto justify-center"
            >
              <svg v-if="loading" class="animate-spin h-4 w-4 mr-2" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              {{ loading ? 'Escaneando...' : 'Escanear' }}
            </button>
          </div>
          
          <!-- Ayuda de formato -->
          <div class="p-2 sm:p-3 bg-green-50 rounded-lg">
            <p class="text-xs sm:text-sm text-green-800">
              <svg class="h-4 w-4 inline mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <strong>Formatos válidos:</strong>
            </p>
            <ul class="text-xs sm:text-sm text-green-700 mt-1 space-y-1">
              <li>• <code class="bg-green-100 px-1 rounded text-xs">SUPPLY_...</code> - Insumos individuales</li>
              <li>• <code class="bg-blue-100 px-1 rounded text-xs">BATCH_...</code> - Información del lote</li>
            </ul>
          </div>
        </div>

        <!-- Escáner con Cámara -->
        <div class="space-y-3 sm:space-y-4">
          <label class="block text-xs sm:text-sm font-medium text-gray-700">
            Usar Cámara:
          </label>
          
          <div class="relative">
            <!-- Vista de la cámara -->
            <div v-if="cameraActive" class="bg-gray-900 rounded-lg overflow-hidden aspect-video relative">
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
                class="w-full h-full object-cover"
              ></video>
              
              <div class="absolute inset-0 flex items-center justify-center">
                <div class="w-48 h-48 border-2 border-white opacity-50 rounded"></div>
              </div>
              
              <div class="absolute bottom-2 sm:bottom-4 left-1/2 transform -translate-x-1/2">
                <button @click="stopCamera" class="btn-danger text-xs sm:text-sm px-3 py-1.5 sm:py-1">
                  Detener Cámara
                </button>
              </div>
            </div>
            
            <!-- Botón para activar cámara -->
            <button
              v-else
              @click="startCameraScanner"
              class="w-full h-28 sm:h-32 border-2 border-dashed border-gray-300 hover:border-blue-400 rounded-lg flex flex-col items-center justify-center transition-colors"
              :disabled="cameraStarting"
            >
              <svg class="h-7 w-7 sm:h-8 sm:w-8 text-gray-400 mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" />
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z" />
              </svg>
              <span class="text-xs sm:text-sm text-gray-600">
                {{ cameraStarting ? 'Iniciando Cámara...' : 'Activar Cámara' }}
              </span>
            </button>
          </div>
          
          <!-- Error de cámara -->
          <div v-if="cameraError" class="text-xs sm:text-sm text-red-600 bg-red-50 p-2 rounded">
            {{ cameraError }}
          </div>
          
          <!-- Instrucciones -->
          <div class="text-xs sm:text-sm text-gray-600 bg-gray-50 p-2 sm:p-3 rounded-lg">
            <p class="font-medium mb-1">Instrucciones:</p>
            <ul class="space-y-1 text-xs">
              <li>• Posicione el código QR dentro del marco</li>
              <li>• Mantenga el código estable y bien iluminado</li>
              <li>• La detección es automática</li>
            </ul>
          </div>
        </div>
      </div>
    </div>

    <!-- Mensaje de Error -->
    <div v-if="error" class="bg-red-50 border border-red-200 rounded-lg p-3 sm:p-4 mb-4 sm:mb-6">
      <div class="flex items-start space-x-2 sm:space-x-3">
        <svg class="h-5 w-5 text-red-600 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <div class="flex-1 min-w-0">
          <h4 class="text-xs sm:text-sm font-medium text-red-800">Error al escanear</h4>
          <p class="text-xs sm:text-sm text-red-700 mt-1 break-words">{{ error }}</p>
          <button @click="clearError" class="text-xs sm:text-sm text-red-600 hover:text-red-800 mt-2 underline">
            Limpiar error
          </button>
        </div>
      </div>
    </div>

    <!-- Información del Insumo Escaneado -->
    <div v-if="scannedInfo && !error" class="bg-white rounded-lg shadow-sm border overflow-hidden mb-4 sm:mb-6">
      <QRInfoDisplay 
        :qr-info="scannedInfo"
        :show-traceability="true"
        :scan-context="lastScanContext"
        @view-details="viewDetails"
        @view-batch="viewBatch"
        @consume-supply="consumeSupply"
      />
      
      <!-- Gestión del Carrito -->
      <div v-if="canAddToCart(scannedInfo) && availableCarts.length > 0" class="p-3 sm:p-4 border-t border-gray-200 bg-blue-50">
        <div class="flex items-start gap-3">
          <div class="flex-shrink-0">
            <svg class="h-6 w-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z" />
            </svg>
          </div>
          <div class="flex-1">
            <h3 class="text-sm font-medium text-blue-900">Agregar a Carrito</h3>
            <p class="text-sm text-blue-700 mt-1">
              Este insumo está disponible. Puedes agregarlo a un carrito activo.
            </p>
            <div class="mt-3">
              <label class="block text-sm font-medium text-blue-900 mb-2">Seleccionar Carrito:</label>
              <select 
                v-model="selectedCartForAdd" 
                class="form-select w-full mb-3"
              >
                <option value="">-- Seleccione un carrito --</option>
                <option 
                  v-for="cart in availableCarts" 
                  :key="cart.id" 
                  :value="cart.id"
                >
                  {{ cart.cart_number }} - {{ cart.supply_request?.request_number }}
                </option>
              </select>
              <button
                @click="addScannedSupplyToCart"
                :disabled="!selectedCartForAdd || addingToCart"
                class="btn-primary w-full"
              >
                <svg v-if="addingToCart" class="animate-spin h-4 w-4 mr-2" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                {{ addingToCart ? 'Agregando...' : 'Agregar al Carrito' }}
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- NUEVA LÓGICA: State-specific recommendations -->
      <div v-if="getStateRecommendation(scannedInfo)" class="p-3 sm:p-4 border-t border-gray-200">
        <div :class="getRecommendationClass(scannedInfo)" class="rounded-md p-3 sm:p-4">
          <div class="flex">
            <div class="flex-shrink-0">
              <svg :class="getRecommendationIconClass(scannedInfo)" class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" :d="getRecommendationIcon(scannedInfo)" />
              </svg>
            </div>
            <div class="ml-2 sm:ml-3">
              <h3 :class="getRecommendationTitleClass(scannedInfo)" class="text-xs sm:text-sm font-medium">
                {{ getStateRecommendation(scannedInfo).title }}
              </h3>
              <div :class="getRecommendationTextClass(scannedInfo)" class="mt-1 sm:mt-2 text-xs sm:text-sm">
                {{ getStateRecommendation(scannedInfo).message }}
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <!-- NUEVA LÓGICA: Acciones basadas en estado -->
      <div v-if="scannedInfo.supply_info && !scannedInfo.is_consumed" class="p-3 sm:p-4 border-t border-gray-200 bg-gray-50">
        <div class="flex flex-col sm:flex-row sm:flex-wrap gap-2 sm:gap-3">
          <!-- NUEVA LÓGICA: Solo mostrar consumir si estado es "recepcionado" -->
          <router-link 
            v-if="canBeConsumed(scannedInfo)"
            :to="{ name: 'QRConsumer', query: { qr: scannedInfo.qr_code } }" 
            class="btn-primary text-sm flex items-center justify-center"
          >
            <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
            Consumir
          </router-link>
          
          <!-- NUEVA LÓGICA: Solo mostrar transferir si estado es "disponible" -->
          <router-link 
            v-if="canBeTransferred(scannedInfo)"
            :to="{ name: 'QRTransfer', query: { qr: scannedInfo.qr_code } }" 
            class="btn-primary text-sm flex items-center justify-center"
          >
            <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4" />
            </svg>
            Transferir
          </router-link>
          
          <!-- NUEVA LÓGICA: Solo mostrar recepcionar si estado es "en_camino_a_pabellon" -->
          <router-link 
            v-if="canBeReceived(scannedInfo)"
            :to="{ name: 'QRReception', query: { qr: scannedInfo.qr_code } }" 
            class="btn-primary text-sm flex items-center justify-center"
          >
            <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            Recepcionar
          </router-link>
          
          <!-- NUEVA LÓGICA: Botón rojo para regresar a bodega -->
          <button 
            v-if="canBeReturnedToStore(scannedInfo)"
            @click="returnToStore(scannedInfo.qr_code)"
            :disabled="returningToStore"
            class="btn-danger text-sm flex items-center justify-center"
          >
            <svg v-if="returningToStore" class="animate-spin h-4 w-4 mr-2" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            <svg v-else class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 10h10a8 8 0 018 8v2M3 10l6 6m-6-6l6-6" />
            </svg>
            {{ returningToStore ? 'Regresando...' : 'Regresar a Bodega' }}
          </button>

          <!-- NUEVO: Botón para confirmar llegada a bodega -->
          <button 
            v-if="isOnRouteToStore(scannedInfo)"
            @click="confirmArrivalToStore(scannedInfo.qr_code)"
            :disabled="confirmingArrival"
            class="btn-success text-sm flex items-center justify-center"
          >
            <svg v-if="confirmingArrival" class="animate-spin h-4 w-4 mr-2" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            <svg v-else class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
            {{ confirmingArrival ? 'Confirmando...' : 'Confirmar Llegada a Bodega' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Historial de Escaneos -->
    <div v-if="scanHistory.length > 0" class="bg-white rounded-lg shadow-sm border p-3 sm:p-6">
      <div class="flex items-center justify-between mb-3 sm:mb-4">
        <h3 class="text-base sm:text-lg font-medium text-gray-900 flex items-center">
          <svg class="h-4 w-4 sm:h-5 sm:w-5 mr-2 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          Historial Reciente
        </h3>
        <button @click="clearHistory" class="text-xs sm:text-sm text-gray-500 hover:text-gray-700">
          Limpiar
        </button>
      </div>
      
      <div class="space-y-2 sm:space-y-3 max-h-80 overflow-y-auto">
        <button 
          v-for="(item, index) in scanHistory" 
          :key="index"
          @click="quickRescan(item.qr_code)"
          class="w-full text-left p-2 sm:p-3 bg-gray-50 rounded-lg hover:bg-blue-50 transition-colors border border-gray-200 hover:border-blue-300"
        >
          <div class="flex items-center space-x-2 sm:space-x-3">
            <!-- Icono según tipo de código -->
            <div class="flex-shrink-0">
              <div v-if="item.type === 'medical_supply'" class="w-8 h-8 sm:w-10 sm:h-10 bg-green-100 rounded-full flex items-center justify-center">
                <svg class="h-4 w-4 sm:h-6 sm:w-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                </svg>
              </div>
              <div v-else class="w-8 h-8 sm:w-10 sm:h-10 bg-blue-100 rounded-full flex items-center justify-center">
                <svg class="h-4 w-4 sm:h-6 sm:w-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
                </svg>
              </div>
            </div>
            
            <!-- Contenido principal -->
            <div class="flex-1 min-w-0">
              <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-1 sm:gap-0">
                <div class="flex items-center space-x-1 sm:space-x-2 flex-wrap">
                  <code class="text-xs sm:text-sm font-mono text-gray-900 truncate">{{ item.qr_code }}</code>
                  <span :class="[
                    'px-1.5 sm:px-2 py-0.5 sm:py-1 text-xs font-medium rounded-full whitespace-nowrap',
                    item.type === 'medical_supply' ? 'bg-green-100 text-green-700' : 'bg-blue-100 text-blue-700'
                  ]">
                    {{ item.type === 'medical_supply' ? 'Insumo' : 'Lote' }}
                  </span>
                </div>
                <span class="text-xs text-gray-500">{{ formatDate(item.scanned_at) }}</span>
              </div>
              <div v-if="item.supply_name" class="text-xs sm:text-sm text-gray-700 mt-1 truncate font-medium">
                {{ item.supply_name }}
              </div>
            </div>
            
            <!-- Indicador de acción -->
            <div class="flex-shrink-0">
              <svg class="h-4 w-4 sm:h-5 sm:w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
              </svg>
            </div>
          </div>
        </button>
      </div>
    </div>
  </div>

  <!-- Sistema de Notificaciones -->
  <div v-if="notification" class="fixed top-4 right-4 left-4 sm:left-auto z-50 max-w-sm sm:w-full">
    <div :class="[
      'rounded-lg p-3 sm:p-4 shadow-lg border transition-all duration-300',
      notification.type === 'success' ? 'bg-green-50 border-green-200' : 'bg-red-50 border-red-200'
    ]">
      <div class="flex items-start">
        <div class="flex-shrink-0">
          <svg v-if="notification.type === 'success'" class="h-5 w-5 text-green-500" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
          </svg>
          <svg v-else class="h-5 w-5 text-red-500" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
          </svg>
        </div>
        <div class="ml-2 sm:ml-3 flex-1 min-w-0">
          <p :class="[
            'text-xs sm:text-sm font-medium break-words',
            notification.type === 'success' ? 'text-green-800' : 'text-red-800'
          ]">
            {{ notification.message }}
          </p>
        </div>
        <div class="ml-2 sm:ml-4 flex-shrink-0">
          <button @click="closeNotification" :class="[
            'rounded-md inline-flex text-sm focus:outline-none focus:ring-2 focus:ring-offset-2',
            notification.type === 'success' ? 'text-green-500 hover:text-green-600 focus:ring-green-500' : 'text-red-500 hover:text-red-600 focus:ring-red-500'
          ]">
            <svg class="h-4 w-4" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
            </svg>
          </button>
        </div>
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
import qrService from '@/services/qr/qrService'
import returnToBodegaService from '@/services/management/returnToBodegaService'
import jsQR from 'jsqr'
import QRInfoDisplay from '@/components/qr/QRInfoDisplay.vue'
import Swal from 'sweetalert2'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

// Referencias DOM
const videoElement = ref(null)

// Estado del componente
const loading = ref(false)
const error = ref(null)
const qrInput = ref('')
const scannedInfo = ref(null)
const scanHistory = ref([])
const lastScanContext = ref(null)

// Estado de la cámara
const cameraActive = ref(false)
const cameraStarting = ref(false)
const cameraError = ref(null)
const detecting = ref(false)

// Estado del retorno a bodega
const returningToStore = ref(false)
const confirmingArrival = ref(false)

// Estado del carrito
const availableCarts = ref([])
const selectedCartForAdd = ref('')
const addingToCart = ref(false)

// Variables para manejo de cámara
let mediaStream = null
let animationFrameId = null

// Computed properties
const currentUser = computed(() => authStore.user)
const isAuthenticated = computed(() => authStore.isAuthenticated)

// ===== NUEVA LÓGICA: Funciones de estado y recomendaciones =====
const canBeConsumed = (info) => {
  if (!info || info.type !== 'medical_supply') return false
  if (info.is_consumed) return false
  // Si está asociado a un carrito, no mostrar botón de consumir
  if (isInCart(info)) return false
  
  const status = info.supply_info?.Status || info.supply_info?.status || info.status || info.current_status
  return status === 'recepcionado'
}

const canBeTransferred = (info) => {
  if (!info || info.type !== 'medical_supply') return false
  if (info.is_consumed) return false
  
  const status = info.supply_info?.Status || info.supply_info?.status || info.status || info.current_status
  return status === 'disponible'
}

const canBeReceived = (info) => {
  if (!info || info.type !== 'medical_supply') return false
  if (info.is_consumed) return false
  
  const status = info.supply_info?.Status || info.supply_info?.status || info.status || info.current_status
  return status === 'en_camino_a_pabellon'
}

const canBeReturnedToStore = (info) => {
  if (!info || info.type !== 'medical_supply') return false
  if (info.is_consumed) return false
  // Si está asociado a un carrito, no mostrar botón de devolver a bodega
  if (isInCart(info)) return false
  
  const status = info.supply_info?.Status || info.supply_info?.status || info.status || info.current_status
  // Puede ser regresado si está recepcionado (hace tiempo sin consumir)
  return status === 'recepcionado'
}

const isOnRouteToStore = (info) => {
  if (!info || info.type !== 'medical_supply') return false
  
  const status = info.supply_info?.Status || info.supply_info?.status || info.status || info.current_status
  // Puede confirmar llegada si está en camino a bodega
  return status === 'en_camino_a_bodega'
}

const getStateRecommendation = (info) => {
  if (!info || info.type !== 'medical_supply') {
    if (info?.type === 'batch') {
      return {
        title: 'Código de Lote',
        message: 'Este código corresponde a un lote completo. Use los códigos QR individuales de los insumos para realizar acciones.',
        type: 'info'
      }
    }
    return null
  }

  if (info.is_consumed) {
    return {
      title: 'Insumo Consumido',
      message: 'Este insumo ya ha sido consumido. Solo puede consultar su información y trazabilidad.',
      type: 'info'
    }
  }

  const status = info.supply_info?.Status || info.supply_info?.status || info.status || info.current_status

  switch (status) {
    case 'disponible':
      return {
        title: 'Listo para Transferir',
        message: 'Este insumo está disponible para ser transferido.',
        type: 'success'
      }
    case 'recepcionado':
      return {
        title: 'Listo para Consumir',
        message: 'Este insumo tiene estado "recepcionado" y puede ser consumido en un procedimiento médico.',
        type: 'success'
      }
    case 'en_camino_a_pabellon':
      return {
        title: 'Listo para Recepcionar',
        message: 'Este insumo está en camino al pabellón. Puede ser recepcionado para cambiar su estado a "recepcionado".',
        type: 'info'
      }
    case 'transferido':
      return {
        title: 'Insumo Transferido',
        message: 'Este insumo ha sido transferido. Consulte su trazabilidad para ver el destino.',
        type: 'info'
      }
    default:
      return {
        title: 'De regreso a Bodega',
        message: `El insumo tiene estado "${status}". Consulte con el administrador del sistema.`,
        type: 'warning'
      }
  }
}

const getRecommendationClass = (info) => {
  const rec = getStateRecommendation(info)
  if (!rec) return ''
  
  switch (rec.type) {
    case 'success': return 'bg-green-50 border border-green-200'
    case 'warning': return 'bg-yellow-50 border border-yellow-200'
    case 'info': return 'bg-blue-50 border border-blue-200'
    default: return 'bg-gray-50 border border-gray-200'
  }
}

const getRecommendationIconClass = (info) => {
  const rec = getStateRecommendation(info)
  if (!rec) return 'text-gray-400'
  
  switch (rec.type) {
    case 'success': return 'text-green-400'
    case 'warning': return 'text-yellow-400'
    case 'info': return 'text-blue-400'
    default: return 'text-gray-400'
  }
}

const getRecommendationTitleClass = (info) => {
  const rec = getStateRecommendation(info)
  if (!rec) return 'text-gray-800'
  
  switch (rec.type) {
    case 'success': return 'text-green-800'
    case 'warning': return 'text-yellow-800'
    case 'info': return 'text-blue-800'
    default: return 'text-gray-800'
  }
}

const getRecommendationTextClass = (info) => {
  const rec = getStateRecommendation(info)
  if (!rec) return 'text-gray-700'
  
  switch (rec.type) {
    case 'success': return 'text-green-700'
    case 'warning': return 'text-yellow-700'
    case 'info': return 'text-blue-700'
    default: return 'text-gray-700'
  }
}

const getRecommendationIcon = (info) => {
  const rec = getStateRecommendation(info)
  if (!rec) return 'M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z'
  
  switch (rec.type) {
    case 'success': return 'M5 13l4 4L19 7'
    case 'warning': return 'M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z'
    case 'info': return 'M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z'
    default: return 'M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z'
  }
}

// ===== FUNCIONES DE CÁMARA =====
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
    
    // Esperar a que el elemento video esté disponible
    await new Promise(resolve => setTimeout(resolve, 100))
    
    if (videoElement.value) {
      videoElement.value.srcObject = stream
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
  if (!cameraActive.value || !videoElement.value) return
  
  detecting.value = true
  
  const detectQR = () => {
    if (!cameraActive.value || !videoElement.value) {
      detecting.value = false
      return
    }
    
    try {
      // Verificar que el video tiene dimensiones válidas
      if (videoElement.value.videoWidth === 0 || videoElement.value.videoHeight === 0) {
        animationFrameId = requestAnimationFrame(detectQR)
        return
      }
      
      // Crear canvas para procesar la imagen
      const canvas = document.createElement('canvas')
      const context = canvas.getContext('2d')
      const video = videoElement.value
      
      canvas.width = video.videoWidth
      canvas.height = video.videoHeight
      
      // Aplicar transformación espejo al canvas para que coincida con la vista del video
      context.save()
      context.scale(-1, 1)
      context.drawImage(video, -canvas.width, 0, canvas.width, canvas.height)
      context.restore()
      
      const imageData = context.getImageData(0, 0, canvas.width, canvas.height)
      const qrCode = jsQR(imageData.data, imageData.width, imageData.height, {
        inversionAttempts: "dontInvert"
      })
      
      if (qrCode && qrCode.data) {
        const qrData = qrCode.data.trim()
        
        if (qrData.startsWith('SUPPLY_') || qrData.startsWith('BATCH_')) {
          detecting.value = false
          qrInput.value = qrData
          stopCamera()
          
          // Auto-escanear después de un breve delay
          setTimeout(() => {
            scanQRCode()
          }, 300)
          
          return
        }
      }
      
    } catch (canvasError) {
      console.warn('Error en procesamiento de canvas:', canvasError)
    }
    
    if (cameraActive.value) {
      animationFrameId = requestAnimationFrame(detectQR)
    } else {
      detecting.value = false
    }
  }
  
  detectQR()
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
  
  if (videoElement.value) {
    videoElement.value.srcObject = null
  }
  
  cameraActive.value = false
  detecting.value = false
  cameraError.value = null
}

// ===== FUNCIONES DE ESCANEO =====
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
    
    // Crear contexto de escaneo
    const scanContext = buildScanContext()
    
    // Escanear código QR
    const result = await qrService.scanQRCode(qrInput.value.trim(), scanContext)
    
    scannedInfo.value = result
    lastScanContext.value = scanContext
    
    // Añadir al historial
    addToHistory(qrInput.value.trim(), result.type, result, scanContext)
    
  } catch (err) {
    console.error('Error scanning QR:', err)
    error.value = err.response?.data?.error || err.message || 'Error al escanear el código QR'
  } finally {
    loading.value = false
  }
}

const buildScanContext = () => {
  return {
    scan_purpose: 'lookup',
    scan_source: 'web',
    user_agent: navigator.userAgent,
    device_info: {
      platform: navigator.platform,
      language: navigator.language,
      screen_resolution: `${screen.width}x${screen.height}`
    }
  }
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

// ===== FUNCIONES DE HISTORIAL =====
const addToHistory = (qrCode, type, qrInfo = null, scanContext = null) => {
  const existing = scanHistory.value.findIndex(item => item.qr_code === qrCode)
  
  const historyItem = {
    qr_code: qrCode,
    type: type,
    scanned_at: new Date(),
    supply_name: qrInfo?.supply_info?.supply_code_name,
    success: true
  }
  
  if (existing >= 0) {
    scanHistory.value.splice(existing, 1)
  }
  
  scanHistory.value.unshift(historyItem)
  
  // Mantener solo los últimos 10 escaneos
  if (scanHistory.value.length > 10) {
    scanHistory.value = scanHistory.value.slice(0, 10)
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
  saveHistory()
}

// ===== FUNCIONES DE NAVEGACIÓN =====
const quickConsume = () => {
  if (scannedInfo.value) {
    router.push({
      name: 'QRConsumer',
      query: { 
        qr: scannedInfo.value.qr_code,
        quick: 'true'
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
        purpose: 'transfer'
      }
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

// ===== NUEVA FUNCIÓN: REGRESAR A BODEGA =====
const returnToStore = async (qrCode) => {
  if (!qrCode || returningToStore.value) return
  
  // Confirmar la acción
  const result = await Swal.fire({
    title: '¿Está seguro de que desea regresar este insumo a bodega?',
    html: 'Esta acción cambiará el estado del insumo a <b>en_camino_a_bodega</b>. Deberá confirmar la llegada cuando el insumo llegue físicamente a bodega.',
    icon: 'question',
    showCancelButton: true,
    confirmButtonText: 'Sí, marcar como en camino',
    cancelButtonText: 'Cancelar',
  })
  if (!result.isConfirmed) return
  
  returningToStore.value = true
  error.value = null
  
  try {
    const result = await returnToBodegaService.returnSupplyToStore(
      qrCode, 
      'Retorno manual desde escáner QR'
    )
    
    // Mostrar notificación de éxito
    showSuccessNotification('Insumo marcado como en camino a bodega. Confirme la llegada cuando el insumo llegue físicamente.')
    
    // Volver a escanear para actualizar la información
    await scanQRCode()
    
  } catch (err) {
    console.error('Error regresando a bodega:', err)
    error.value = err.message || 'Error al regresar el insumo a bodega'
  } finally {
    returningToStore.value = false
  }
}

// ===== NUEVA FUNCIÓN: CONFIRMAR LLEGADA A BODEGA =====
const confirmArrivalToStore = async (qrCode) => {
  if (!qrCode || confirmingArrival.value) return
  
  // Confirmar la acción
  const result = await Swal.fire({
    title: '¿Confirma que este insumo ha llegado a bodega?',
    html: 'Esta acción cambiará el estado del insumo a <b>disponible</b> y será registrada en el historial de trazabilidad.',
    icon: 'question',
    showCancelButton: true,
    confirmButtonText: 'Sí, confirmar',
    cancelButtonText: 'Cancelar',
  })
  if (!result.isConfirmed) return
  
  confirmingArrival.value = true
  error.value = null
  
  try {
    // Usar el servicio real para confirmar llegada a bodega
    const result = await returnToBodegaService.confirmArrivalToStore(
      qrCode, 
      'Llegada confirmada desde escáner QR'
    )
    
    // Mostrar notificación de éxito
    showSuccessNotification('Llegada a bodega confirmada exitosamente')
    
    // Volver a escanear para actualizar la información
    await scanQRCode()
    
  } catch (err) {
    console.error('Error confirmando llegada a bodega:', err)
    error.value = err.message || 'Error al confirmar llegada a bodega'
  } finally {
    confirmingArrival.value = false
  }
}

// ===== GESTIÓN DE CARRITOS =====
const canAddToCart = (info) => {
  if (!info || info.type !== 'medical_supply') return false
  const status = info.supply_info?.Status || info.supply_info?.status || info.status || info.current_status
  // Solo puede agregar si está disponible y no está ya asignado a un carrito
  return status === 'disponible' && !info.request_assignment
}

const isInCart = (info) => {
  if (!info || info.type !== 'medical_supply') return false
  // Está en un carrito si tiene una asignación de request Y un carrito asociado
  if (!info.request_assignment) return false
  
  // Verificar si hay un carrito activo asociado
  const cart = info.request_assignment?.cart || info.request_assignment?.Cart
  if (!cart) return false
  
  // Solo considerar que está en un carrito si el carrito está activo
  return cart.status === 'active' || cart.Status === 'active'
}

const loadAvailableCarts = async () => {
  try {
    const cartService = (await import('@/services/requests/cartService')).default
    const response = await cartService.getAllCarts(1, 100, 'active')
    if (response.success) {
      availableCarts.value = response.data || []
    }
  } catch (err) {
    console.error('Error cargando carritos:', err)
  }
}

const addScannedSupplyToCart = async () => {
  if (!selectedCartForAdd.value || !scannedInfo.value || addingToCart.value) return
  
  addingToCart.value = true
  error.value = null
  
  try {
    const cartService = (await import('@/services/requests/cartService')).default
    const supplyRequestService = (await import('@/services/requests/supplyRequestService')).default
    
    // Primero, obtener el carrito para saber su supply_request_id
    const cart = availableCarts.value.find(c => c.id === parseInt(selectedCartForAdd.value))
    if (!cart) {
      throw new Error('Carrito no encontrado')
    }
    
    // Necesitamos crear una asignación QR primero
    // Para esto, necesitamos el supply_request_id y un supply_request_item_id
    // Vamos a obtener los items de la solicitud
    const itemsResponse = await supplyRequestService.getSupplyRequestItems(cart.supply_request_id)
    if (!itemsResponse.success || !itemsResponse.data || itemsResponse.data.length === 0) {
      throw new Error('No se encontraron items en la solicitud')
    }
    
    // Buscar un item que coincida con el supply_code del insumo escaneado
    const supplyCode = scannedInfo.value.supply_info?.supply_code || scannedInfo.value.supply_info?.code
    let targetItem = itemsResponse.data.find(item => item.supply_code === supplyCode)
    
    // Si no hay un item específico para este código, usar el primero disponible
    if (!targetItem) {
      targetItem = itemsResponse.data[0]
    }
    
    // Asignar el QR a la solicitud
    const assignmentData = {
      supply_request_id: cart.supply_request_id,
      supply_request_item_id: targetItem.id,
      qr_code: scannedInfo.value.qr_code,
      notes: 'Agregado mediante escáner QR'
    }
    
    await supplyRequestService.assignQRToRequest(assignmentData)
    
    // Mostrar notificación de éxito
    showSuccessNotification(`Insumo agregado al carrito ${cart.cart_number} exitosamente`)
    
    // Limpiar selección y reescanear
    selectedCartForAdd.value = ''
    await scanQRCode()
    
  } catch (err) {
    console.error('Error agregando al carrito:', err)
    const errorMsg = err.response?.data?.error || err.message || 'Error al agregar el insumo al carrito'
    Swal.fire({
      icon: 'error',
      title: 'Error',
      text: errorMsg
    })
  } finally {
    addingToCart.value = false
  }
}

// ===== SISTEMA DE NOTIFICACIONES =====
const notification = ref(null)

const showSuccessNotification = (message) => {
  notification.value = {
    type: 'success',
    message: message,
    visible: true
  }
  setTimeout(() => {
    notification.value = null
  }, 4000)
}

const showErrorNotification = (message) => {
  notification.value = {
    type: 'error',
    message: message,
    visible: true
  }
  setTimeout(() => {
    notification.value = null
  }, 5000)
}

const closeNotification = () => {
  notification.value = null
}

// ===== FUNCIONES AUXILIARES =====
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

// ===== LIFECYCLE =====
onMounted(() => {
  loadHistory()
  loadAvailableCarts()
  
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
  stopCamera()
})
</script>

<style scoped>
.form-input {
  @apply block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm;
}

/* Usar .btn-primary de style.css global */

/* Usar clases de botones de style.css global */

/* Mejorar experiencia táctil en móviles */
@media (max-width: 640px) {
  /* Aumentar área táctil de botones */
  button,
  .btn-primary,
  .btn-secondary,
  .btn-danger,
  .btn-success {
    min-height: 44px;
  }
  
  /* Mejorar áreas de entrada */
  input,
  select,
  textarea,
  .form-input {
    min-height: 44px;
    font-size: 16px; /* Prevenir zoom automático en iOS */
  }

  /* Mejorar scrolling en historial */
  .overflow-y-auto {
    -webkit-overflow-scrolling: touch;
  }

  /* Suavizar transiciones */
  * {
    -webkit-tap-highlight-color: transparent;
  }

  /* Mejorar códigos QR en móvil */
  code {
    word-break: break-all;
  }
}

/* Transiciones suaves */
button,
input,
select,
textarea {
  transition: all 0.2s ease-in-out;
}

/* Mejorar contraste en modo oscuro de cámara */
video {
  max-width: 100%;
}
</style>