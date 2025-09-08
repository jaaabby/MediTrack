<template>
  <div class="space-y-6">
    <!-- Header con estadísticas mejoradas -->
    <div class="bg-white rounded-lg shadow-sm border p-6">
      <div class="flex items-center justify-between">
        <div>
          <h2 class="text-xl font-semibold text-gray-900">Consumir Insumos Médicos</h2>
          <p class="text-gray-600 mt-1">Escanea o ingresa códigos QR para registrar el consumo de productos con trazabilidad completa</p>
        </div>
        <div class="grid grid-cols-2 gap-4 text-right">
          <div>
            <p class="text-sm text-gray-500">Productos consumidos hoy</p>
            <p class="text-2xl font-bold text-blue-600">{{ consumptionStats.today }}</p>
          </div>
          <div>
            <p class="text-sm text-gray-500">Total esta semana</p>
            <p class="text-2xl font-bold text-green-600">{{ consumptionStats.week }}</p>
          </div>
        </div>
      </div>
      
      <!-- Información de usuario y ubicación actual -->
      <div class="mt-4 flex items-center justify-between bg-gray-50 rounded-lg p-3">
        <div class="flex items-center space-x-4">
          <div>
            <span class="text-sm font-medium text-gray-700">Usuario:</span>
            <span class="text-sm text-gray-900 ml-1">{{ currentUser?.name || 'No identificado' }}</span>
          </div>
          <div>
            <span class="text-sm font-medium text-gray-700">Ubicación:</span>
            <span class="text-sm text-gray-900 ml-1">{{ currentLocation?.name || 'No seleccionada' }}</span>
          </div>
        </div>
        <div class="text-xs text-gray-500">
          Sesión iniciada: {{ formatDate(sessionStart) }}
        </div>
      </div>
    </div>

    <!-- Selector de propósito de consumo -->
    <div v-if="!scannedProduct || isQuickMode" class="bg-white rounded-lg shadow-sm border p-6">
      <h3 class="text-lg font-medium text-gray-900 mb-4">Tipo de Consumo</h3>
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-3">
        <button
          v-for="purpose in consumptionPurposes"
          :key="purpose.value"
          @click="selectedConsumptionPurpose = purpose.value"
          :class="getPurposeButtonClass(purpose.value)"
        >
          <div class="flex flex-col items-center">
            <svg class="h-6 w-6 mb-2" :class="purpose.iconClass" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" :d="purpose.icon" />
            </svg>
            <span class="font-medium">{{ purpose.label }}</span>
            <span class="text-xs text-gray-500 mt-1">{{ purpose.description }}</span>
          </div>
        </button>
      </div>
    </div>

    <!-- Scanner/Input Section -->
    <div v-if="!scannedProduct" class="bg-white rounded-lg shadow-sm border p-6">
      <h3 class="text-lg font-medium text-gray-900 mb-4">
        <svg class="h-5 w-5 inline mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h2M4 4h5l2 3h3l2-3h5v5M4 4v5m16-5v5" />
        </svg>
        Escanear Código QR
      </h3>

      <div class="grid md:grid-cols-2 gap-6">
        <!-- QR Input -->
        <div>
          <label for="qrInput" class="block text-sm font-medium text-gray-700 mb-2">
            Código QR del Producto:
          </label>
          <div class="flex space-x-2">
            <input
              id="qrInput"
              v-model="qrInput"
              type="text"
              placeholder="Ej: SUPPLY_1755580808_abc123def"
              class="form-input flex-1"
              @keyup.enter="scanQRCode"
              :disabled="scanning"
            />
            <button
              @click="scanQRCode"
              :disabled="!qrInput.trim() || scanning"
              class="btn-primary"
            >
              <LoadingIcon v-if="scanning" />
              {{ scanning ? 'Escaneando...' : 'Escanear' }}
            </button>
          </div>
        </div>

        <!-- Camera Scanner -->
        <div class="text-center">
          <button
            @click="startCameraScanner"
            class="btn-secondary w-full h-20 border-2 border-dashed border-gray-300 hover:border-gray-400"
            :disabled="cameraActive"
          >
            <svg class="h-8 w-8 mx-auto mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
            <span class="text-sm">{{ cameraActive ? 'Cámara Activa' : 'Usar Cámara' }}</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Consumption Form -->
    <div v-if="scannedProduct && !error" class="bg-white rounded-lg shadow-sm border p-6">
      <!-- Indicador de progreso para modo rápido -->
      <QuickModeProgress 
        v-if="isQuickMode" 
        :current-location="currentLocation"
      />

      <h3 class="text-lg font-medium text-gray-900 mb-4">
        <svg class="h-5 w-5 inline mr-2 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        Producto Encontrado - Registrar Consumo
      </h3>

      <!-- Product Info Display -->
      <ProductInfoDisplay 
        :product="scannedProduct" 
        :scan-context="lastScanContext"
        :current-user="currentUser"
        :current-location="currentLocation"
      />

      <!-- Consumption Form -->
      <form v-if="canConsumeProduct" @submit.prevent="consumeProduct">
        <div class="grid md:grid-cols-2 gap-6">
          <div>
            <label for="userRUT" class="block text-sm font-medium text-gray-700 mb-2">
              RUT del Usuario <span class="text-red-500">*</span>
            </label>
            <input
              id="userRUT"
              v-model="consumptionForm.userRUT"
              type="text"
              placeholder="12.345.678-9"
              class="form-input"
              required
            />
          </div>

          <div>
            <label for="destinationType" class="block text-sm font-medium text-gray-700 mb-2">
              Tipo de Destino <span class="text-red-500">*</span>
            </label>
            <select
              id="destinationType"
              v-model="consumptionForm.destinationType"
              @change="onDestinationTypeChange"
              class="form-select"
              required
            >
              <option value="">Seleccionar tipo</option>
              <option value="pavilion">Pabellón</option>
              <option value="store">Almacén</option>
            </select>
          </div>

          <!-- Selector de Centro Médico (solo si destination es pavilion) -->
          <div v-if="consumptionForm.destinationType === 'pavilion'">
            <label for="medicalCenter" class="block text-sm font-medium text-gray-700 mb-2">
              Centro Médico <span class="text-red-500">*</span>
            </label>
            <select
              id="medicalCenter"
              v-model="consumptionForm.medicalCenterId"
              @change="onMedicalCenterChange"
              class="form-select"
              :disabled="loadingMedicalCenters"
              required
            >
              <option value="">Seleccionar centro médico</option>
              <option
                v-for="center in medicalCenters"
                :key="center.id"
                :value="center.id"
              >
                {{ center.name }}
              </option>
            </select>
            <p v-if="loadingMedicalCenters" class="text-xs text-gray-500 mt-1">Cargando centros médicos...</p>
          </div>

          <!-- Selector de Pabellón (solo si destination es pavilion y se ha seleccionado centro médico) -->
          <div v-if="consumptionForm.destinationType === 'pavilion' && consumptionForm.medicalCenterId">
            <label for="pavilion" class="block text-sm font-medium text-gray-700 mb-2">
              Pabellón <span class="text-red-500">*</span>
            </label>
            <select
              id="pavilion"
              v-model="consumptionForm.destinationID"
              class="form-select"
              :disabled="loadingPavilions"
              required
            >
              <option value="">Seleccionar pabellón</option>
              <option
                v-for="pavilion in availablePavilions"
                :key="pavilion.id"
                :value="pavilion.id"
              >
                {{ pavilion.name }}
              </option>
            </select>
            <p v-if="loadingPavilions" class="text-xs text-gray-500 mt-1">Cargando pabellones...</p>
          </div>

          <!-- Selector de Almacén (solo si destination es store) -->
          <div v-if="consumptionForm.destinationType === 'store'">
            <label for="store" class="block text-sm font-medium text-gray-700 mb-2">
              Almacén <span class="text-red-500">*</span>
            </label>
            <select
              id="store"
              v-model="consumptionForm.destinationID"
              class="form-select"
              :disabled="loadingStores"
              required
            >
              <option value="">Seleccionar almacén</option>
              <option
                v-for="store in stores"
                :key="store.id"
                :value="store.id"
              >
                {{ store.name }} - {{ store.location }}
              </option>
            </select>
            <p v-if="loadingStores" class="text-xs text-gray-500 mt-1">Cargando almacenes...</p>
          </div>

          <div>
            <label for="notes" class="block text-sm font-medium text-gray-700 mb-2">
              Notas (Opcional)
            </label>
            <input
              id="notes"
              v-model="consumptionForm.notes"
              type="text"
              placeholder="Observaciones adicionales"
              class="form-input"
            />
          </div>
        </div>

        <!-- Submit Button -->
        <div class="mt-6 flex justify-end space-x-3">
          <button type="button" @click="clearScannedProduct" class="btn-secondary">
            Cancelar
          </button>
          <button type="submit" :disabled="consuming" class="btn-danger">
            <LoadingIcon v-if="consuming" />
            <svg v-else class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
            </svg>
            {{ consuming ? 'Consumiendo...' : 'Confirmar Consumo' }}
          </button>
        </div>
      </form>
    </div>

    <!-- Success Message -->
    <SuccessMessage 
      v-if="consumptionSuccess" 
      :success-data="consumptionSuccess"
      @clear-all="clearAll"
      @view-traceability="viewTraceability"
      @view-batch-history="viewBatchHistory"
    />

    <!-- Error Message -->
    <ErrorMessage 
      v-if="error" 
      :error-text="error"
      :is-quick-mode="isQuickMode"
      :qr-input="qrInput"
      @clear-error="clearError"
      @retry-scan="retryQRScan"
    />

    <!-- Recent Consumptions -->
    <RecentConsumptions 
      v-if="recentConsumptions.length > 0"
      :consumptions="recentConsumptions"
      @view-traceability="viewTraceability"
      @view-batch-history="viewBatchHistory"
    />

    <!-- Location Selector Component -->
    <LocationSelector @location-changed="onLocationChanged" />
  </div>
</template>

<script setup>
import { ref, onMounted, computed, watchEffect, reactive, toRefs } from 'vue'
import { useRouter } from 'vue-router'
import { format } from 'date-fns'
import { es } from 'date-fns/locale'
import { useAuthStore } from '@/stores/auth'
import qrService from '@/services/qrService'
import medicalCenterService from '@/services/medicalCenterService'
import pavilionService from '@/services/pavilionService'
import storeService from '@/services/storeService'
import LocationSelector from '@/components/LocationSelector.vue'

// Componentes optimizados para reutilización
const LoadingIcon = {
  template: `
    <svg class="animate-spin h-4 w-4 mr-2" fill="none" viewBox="0 0 24 24">
      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
      <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
    </svg>
  `
}

const QuickModeProgress = {
  props: ['currentLocation'],
  template: `
    <div class="mb-6 bg-green-50 border border-green-200 rounded-lg p-4">
      <h4 class="text-green-800 font-medium mb-3 flex items-center">
        <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        Modo Consumo Rápido - Pasos Restantes
      </h4>
      <div class="space-y-2 text-sm">
        <ProgressStep :completed="true" text="Producto escaneado automáticamente" />
        <ProgressStep :completed="!!currentLocation" :text="currentLocation ? 'Ubicación seleccionada' : 'Selecciona tu ubicación actual abajo'" />
        <ProgressStep :completed="false" text="Completa los datos del consumo y confirma" />
      </div>
    </div>
  `,
  components: {
    ProgressStep: {
      props: ['completed', 'text'],
      template: `
        <div class="flex items-center">
          <svg v-if="completed" class="h-4 w-4 text-green-600 mr-2" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
          </svg>
          <svg v-else class="h-4 w-4 text-yellow-500 mr-2" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd" />
          </svg>
          <span :class="completed ? 'text-green-700' : 'text-yellow-700'">
            {{ completed ? '✓' : '⚠' }} {{ text }}
          </span>
        </div>
      `
    }
  }
}

const ProductInfoDisplay = {
  props: ['product', 'scanContext', 'currentUser', 'currentLocation'],
  methods: {
    formatDate(dateString) {
      if (!dateString) return 'No disponible'
      try {
        return format(new Date(dateString), 'dd/MM/yyyy HH:mm', { locale: es })
      } catch (error) {
        return dateString
      }
    }
  },
  template: `
    <div class="bg-gray-50 rounded-lg p-4 mb-6">
      <div class="grid md:grid-cols-2 gap-4">
        <div>
          <label class="text-sm font-medium text-gray-600">Nombre del Producto:</label>
          <p class="text-gray-900 font-medium">{{ product.supply_code?.name || 'N/A' }}</p>
        </div>
        <div>
          <label class="text-sm font-medium text-gray-600">Código QR:</label>
          <p class="text-sm font-mono text-gray-800 bg-white px-2 py-1 rounded border">{{ product.qr_code }}</p>
        </div>
        <div>
          <label class="text-sm font-medium text-gray-600">Código del Producto:</label>
          <p class="text-gray-900">{{ product.supply_code?.code || product.supply_info?.code || 'N/A' }}</p>
        </div>
        <div>
          <label class="text-sm font-medium text-gray-600">ID Proveedor:</label>
          <p class="text-gray-900">{{ product.supply_code?.code_supplier || 'N/A' }}</p>
        </div>
      </div>

      <!-- Información de lote -->
      <div v-if="batchInfo" class="mt-4 pt-4 border-t border-gray-200">
        <h5 class="text-sm font-semibold text-gray-700 mb-2">Información del Lote:</h5>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4 text-sm">
          <div v-if="batchInfo.batch_number">
            <label class="font-medium text-gray-600">Número de Lote:</label>
            <p class="text-gray-900">{{ batchInfo.batch_number }}</p>
          </div>
          <div v-if="batchInfo.supplier">
            <label class="font-medium text-gray-600">Proveedor del Lote:</label>
            <p class="text-gray-900">{{ batchInfo.supplier }}</p>
          </div>
          <div v-if="batchInfo.expiration_date">
            <label class="font-medium text-gray-600">Fecha de Vencimiento:</label>
            <p class="text-gray-900">{{ formatDate(batchInfo.expiration_date) }}</p>
          </div>
        </div>
      </div>

      <!-- Información de trazabilidad -->
      <div v-if="scanContext" class="mt-4 pt-4 border-t border-gray-200">
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4 text-sm">
          <div>
            <label class="font-medium text-gray-600">Escaneado por:</label>
            <p class="text-gray-900">{{ scanContext.user_name || currentUser?.name || 'N/A' }}</p>
          </div>
          <div>
            <label class="font-medium text-gray-600">Ubicación de escaneo:</label>
            <p class="text-gray-900">{{ scanContext.location || currentLocation?.name || 'N/A' }}</p>
          </div>
          <div>
            <label class="font-medium text-gray-600">Hora de escaneo:</label>
            <p class="text-gray-900">{{ formatDate(scanContext.scanned_at) }}</p>
          </div>
        </div>
      </div>

      <!-- Status Alert -->
      <StatusAlert :product="product" />
    </div>
  `,
  computed: {
    batchInfo() {
      return this.product.supply_info?.batch_info || this.product.batch_info
    }
  },
  components: {
    StatusAlert: {
      props: ['product'],
      template: `
        <div v-if="product.is_consumed" class="mt-4 bg-red-50 border border-red-200 rounded p-3">
          <div class="flex items-center">
            <svg class="h-5 w-5 text-red-400 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <span class="text-red-800 font-medium">Este producto ya ha sido consumido anteriormente</span>
          </div>
        </div>
        <div v-else-if="!product.can_consume" class="mt-4 bg-yellow-50 border border-yellow-200 rounded p-3">
          <div class="flex items-center">
            <svg class="h-5 w-5 text-yellow-400 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <span class="text-yellow-800 font-medium">Este producto no está disponible para consumo</span>
          </div>
        </div>
      `
    }
  }
}

const SuccessMessage = {
  props: ['successData'],
  emits: ['clear-all', 'view-traceability', 'view-batch-history'],
  methods: {
    formatDate(dateString) {
      if (!dateString) return 'No disponible'
      try {
        return format(new Date(dateString), 'dd/MM/yyyy HH:mm', { locale: es })
      } catch (error) {
        return dateString
      }
    }
  },
  template: `
    <div class="bg-green-50 border border-green-200 rounded-lg p-6">
      <div class="flex items-start">
        <div class="flex-shrink-0">
          <svg class="h-6 w-6 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </div>
        <div class="ml-3 flex-1">
          <h3 class="text-lg font-medium text-green-800">Producto Consumido Exitosamente</h3>
          <div class="mt-2 text-sm text-green-700">
            <p>{{ successData.message }}</p>
            <p v-if="successData.remaining_amount !== undefined" class="mt-1">
              <strong>Cantidad restante en lote:</strong> {{ successData.remaining_amount }} unidades
            </p>
            <p v-if="successData.batch_history_updated" class="mt-1">
              <strong>Historial del lote actualizado automáticamente</strong>
            </p>
          </div>
          
          <!-- Información de trazabilidad -->
          <div v-if="successData.traceability_info" class="mt-4 bg-green-100 rounded p-3">
            <h4 class="text-sm font-medium text-green-800 mb-2">Información de Trazabilidad</h4>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-2 text-xs text-green-700">
              <div>
                <span class="font-medium">Evento registrado:</span> {{ successData.traceability_info.event_id }}
              </div>
              <div>
                <span class="font-medium">Timestamp:</span> {{ formatDate(successData.traceability_info.timestamp) }}
              </div>
              <div>
                <span class="font-medium">Usuario:</span> {{ successData.traceability_info.user || 'N/A' }}
              </div>
              <div>
                <span class="font-medium">Ubicación:</span> {{ successData.traceability_info.location || 'N/A' }}
              </div>
            </div>
          </div>
          
          <div class="mt-4 flex space-x-3">
            <button @click="$emit('clear-all')" class="btn-secondary text-sm">
              Consumir Otro Producto
            </button>
            <button 
              v-if="successData.qr_code"
              @click="$emit('view-traceability', successData.qr_code)" 
              class="btn-primary text-sm"
            >
              Ver Trazabilidad Completa
            </button>
            <button 
              v-if="successData.batch_id"
              @click="$emit('view-batch-history', successData.batch_id)" 
              class="btn-secondary text-sm"
            >
              Ver Historial del Lote
            </button>
          </div>
        </div>
      </div>
    </div>
  `
}

const ErrorMessage = {
  props: ['errorText', 'isQuickMode', 'qrInput'],
  emits: ['clear-error', 'retry-scan'],
  template: `
    <div class="bg-red-50 border border-red-200 rounded-lg p-6">
      <div class="flex items-start">
        <div class="flex-shrink-0">
          <svg class="h-6 w-6 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </div>
        <div class="ml-3 flex-1">
          <h3 class="text-lg font-medium text-red-800">Error</h3>
          <div class="mt-2 text-sm text-red-700">
            <p>{{ errorText }}</p>
            <p v-if="isQuickMode" class="mt-2 text-sm text-red-600">
              <strong>Modo rápido detectado:</strong> El código QR "{{ qrInput }}" no se pudo escanear automáticamente. 
              Puedes intentar escanearlo manualmente o verificar que el producto existe en el sistema.
            </p>
          </div>
          <div class="mt-4 flex space-x-2">
            <button @click="$emit('clear-error')" class="btn-secondary text-sm">
              Intentar de Nuevo
            </button>
            <button v-if="isQuickMode && qrInput" @click="$emit('retry-scan')" class="btn-primary text-sm">
              Reintentar Escaneo Automático
            </button>
          </div>
        </div>
      </div>
    </div>
  `
}

const RecentConsumptions = {
  props: ['consumptions'],
  emits: ['view-traceability', 'view-batch-history'],
  methods: {
    formatDate(dateString) {
      if (!dateString) return 'No disponible'
      try {
        return format(new Date(dateString), 'dd/MM/yyyy HH:mm', { locale: es })
      } catch (error) {
        return dateString
      }
    },
    getConsumptionPurposeLabel(purpose) {
      const labels = {
        'routine': 'Rutina',
        'emergency': 'Emergencia',
        'maintenance': 'Mantenimiento',
        'transfer': 'Transferir'
      }
      return labels[purpose] || purpose
    },
    getPurposeBadgeClass(purpose) {
      return purpose === 'emergency' ? 'bg-red-100 text-red-800' :
             purpose === 'routine' ? 'bg-blue-100 text-blue-800' :
             purpose === 'transfer' ? 'bg-purple-100 text-purple-800' :
             'bg-green-100 text-green-800'
    }
  },
  template: `
    <div class="bg-white rounded-lg shadow-sm border p-6">
      <h3 class="text-lg font-medium text-gray-900 mb-4">
        <svg class="h-5 w-5 inline mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        Consumos Recientes
      </h3>

      <div class="space-y-3">
        <div 
          v-for="(consumption, index) in consumptions" 
          :key="index"
          class="flex items-center justify-between p-3 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors"
        >
          <div class="flex-1">
            <div class="flex items-center justify-between">
              <p class="font-medium text-gray-900">{{ consumption.product_name }}</p>
              <span 
                :class="['px-2 py-1 rounded-full text-xs font-medium', getPurposeBadgeClass(consumption.consumption_purpose)]"
              >
                {{ getConsumptionPurposeLabel(consumption.consumption_purpose) }}
              </span>
            </div>
            <p class="text-sm text-gray-600">QR: {{ consumption.qr_code }}</p>
            <div class="text-xs text-gray-500 mt-1">
              <span>{{ consumption.user_rut }}</span>
              <span class="mx-2">•</span>
              <span>{{ consumption.location || 'Sin ubicación' }}</span>
              <span class="mx-2">•</span>
              <span>{{ formatDate(consumption.consumed_at) }}</span>
            </div>
          </div>
          <div class="ml-4 flex space-x-2">
            <button 
              @click="$emit('view-traceability', consumption.qr_code)"
              class="text-blue-600 hover:text-blue-800 p-1"
              title="Ver trazabilidad"
            >
              <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
              </svg>
            </button>
            <button 
              v-if="consumption.batch_id"
              @click="$emit('view-batch-history', consumption.batch_id)"
              class="text-green-600 hover:text-green-800 p-1"
              title="Ver historial del lote"
            >
              <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
            </button>
          </div>
        </div>
      </div>
    </div>
  `
}

const router = useRouter()
const authStore = useAuthStore()

// Estado reactivo consolidado
const state = reactive({
  qrInput: '',
  scanning: false,
  consuming: false,
  cameraActive: false,
  scannedProduct: null,
  consumptionSuccess: null,
  error: null,
  recentConsumptions: [],
  selectedConsumptionPurpose: 'routine',
  currentLocation: null,
  lastScanContext: null,
  sessionStart: new Date(),
  isQuickMode: false,
  // Nuevos estados para las listas de destinos
  medicalCenters: [],
  pavilions: [],
  stores: [],
  loadingMedicalCenters: false,
  loadingPavilions: false,
  loadingStores: false
})

// Usar destructuring para mantener compatibilidad
const {
  qrInput, scanning, consuming, cameraActive, scannedProduct,
  consumptionSuccess, error, recentConsumptions, selectedConsumptionPurpose,
  currentLocation, lastScanContext, sessionStart, isQuickMode,
  medicalCenters, pavilions, stores,
  loadingMedicalCenters, loadingPavilions, loadingStores
} = toRefs(state)

// Estadísticas consolidadas
const consumptionStats = reactive({
  today: 0,
  week: 0
})

// Computed properties
const currentUser = computed(() => authStore.user)
const isAuthenticated = computed(() => authStore.isAuthenticated)
const canConsumeProduct = computed(() => 
  scannedProduct.value && !scannedProduct.value.is_consumed && scannedProduct.value.can_consume
)

// Pabellones disponibles según el centro médico seleccionado
const availablePavilions = computed(() => {
  if (!consumptionForm.value.medicalCenterId) return []
  return pavilions.value.filter(p => p.medical_center_id === parseInt(consumptionForm.value.medicalCenterId))
})

// Datos de configuración
const consumptionPurposes = [
  {
    value: 'routine',
    label: 'Rutina',
    description: 'Consumo programado normal',
    icon: 'M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2',
    iconClass: 'text-blue-600'
  },
  {
    value: 'emergency',
    label: 'Emergencia',
    description: 'Consumo de urgencia',
    icon: 'M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z',
    iconClass: 'text-red-600'
  },
  {
    value: 'maintenance',
    label: 'Mantenimiento',
    description: 'Uso en mantenimiento',
    icon: 'M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z M15 12a3 3 0 11-6 0 3 3 0 016 0z',
    iconClass: 'text-green-600'
  },
  {
    value: 'transfer',
    label: 'Transferir',
    description: 'Envío a otro pabellón o centro médico',
    icon: 'M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4',
    iconClass: 'text-purple-600'
  }
]

// Formulario de consumo
const consumptionForm = ref({
  userRUT: '',
  destinationType: '',
  destinationID: '',
  medicalCenterId: '',
  notes: ''
})

// Auto-completar RUT del usuario actual
watchEffect(() => {
  if (currentUser.value?.rut && !consumptionForm.value.userRUT) {
    consumptionForm.value.userRUT = currentUser.value.rut
  }
})

// Métodos optimizados
const handleError = (err, defaultMessage = 'Error de conexión') => {
  console.error(defaultMessage + ':', err)
  error.value = err.response?.data?.error || err.message || defaultMessage
}

const buildScanContext = () => ({
  scan_purpose: 'consume',
  consumption_purpose: selectedConsumptionPurpose.value,
  pavilion_id: currentLocation.value?.pavilion_id,
  medical_center_id: currentLocation.value?.medical_center_id,
  scan_source: 'web',
  user_agent: navigator.userAgent,
  device_info: {
    platform: navigator.platform,
    language: navigator.language,
    screen_resolution: `${screen.width}x${screen.height}`
  }
})

const validateConsumptionForm = () => {
  const errors = []
  
  if (!scannedProduct.value?.qr_code) errors.push('Código QR requerido')
  if (!consumptionForm.value.userRUT) errors.push('RUT del usuario requerido')
  if (!consumptionForm.value.destinationType) errors.push('Tipo de destino requerido')
  
  // Validaciones específicas según el tipo de destino
  if (consumptionForm.value.destinationType === 'pavilion') {
    if (!consumptionForm.value.medicalCenterId) errors.push('Centro médico requerido')
    if (!consumptionForm.value.destinationID) errors.push('Pabellón requerido')
  } else if (consumptionForm.value.destinationType === 'store') {
    if (!consumptionForm.value.destinationID) errors.push('Almacén requerido')
  }
  
  const destinationIdNum = parseInt(consumptionForm.value.destinationID)
  if (isNaN(destinationIdNum)) errors.push('Destino debe ser un número válido')
  
  return { isValid: errors.length === 0, errors, destinationIdNum }
}

const updateConsumptionStats = (consumptionData) => {
  // Determinar el nombre del destino
  let destinationName = `${consumptionForm.value.destinationType} ${consumptionForm.value.destinationID}`
  
  if (consumptionForm.value.destinationType === 'pavilion') {
    const pavilion = availablePavilions.value.find(p => p.id === parseInt(consumptionForm.value.destinationID))
    const medicalCenter = medicalCenters.value.find(c => c.id === parseInt(consumptionForm.value.medicalCenterId))
    destinationName = pavilion ? `${medicalCenter?.name || 'Centro'} - ${pavilion.name}` : destinationName
  } else if (consumptionForm.value.destinationType === 'store') {
    const store = stores.value.find(s => s.id === parseInt(consumptionForm.value.destinationID))
    destinationName = store ? `Almacén ${store.name}` : destinationName
  }

  recentConsumptions.value.unshift({
    qr_code: scannedProduct.value.qr_code,
    product_name: scannedProduct.value.supply_info?.supply_code_name || 'N/A',
    user_rut: consumptionForm.value.userRUT,
    user_name: consumptionData.user_name,
    consumed_at: new Date().toISOString(),
    consumption_purpose: selectedConsumptionPurpose.value,
    location: currentLocation.value?.name,
    destination: destinationName,
    batch_id: scannedProduct.value.supply_info?.batch?.id
  })
  
  if (recentConsumptions.value.length > 15) {
    recentConsumptions.value = recentConsumptions.value.slice(0, 15)
  }
  
  consumptionStats.today += 1
  consumptionStats.week += 1
  saveConsumptionStats()
}

// Métodos para cargar datos de destinos
const loadMedicalCenters = async () => {
  loadingMedicalCenters.value = true
  try {
    const response = await medicalCenterService.getAll()
    medicalCenters.value = response.data || []
  } catch (err) {
    console.error('Error cargando centros médicos:', err)
    handleError(err, 'Error al cargar centros médicos')
  } finally {
    loadingMedicalCenters.value = false
  }
}

const loadPavilions = async () => {
  loadingPavilions.value = true
  try {
    pavilions.value = await pavilionService.getAllPavilions()
  } catch (err) {
    console.error('Error cargando pabellones:', err)
    handleError(err, 'Error al cargar pabellones')
  } finally {
    loadingPavilions.value = false
  }
}

const loadStores = async () => {
  loadingStores.value = true
  try {
    stores.value = await storeService.getAllStores()
  } catch (err) {
    console.error('Error cargando almacenes:', err)
    handleError(err, 'Error al cargar almacenes')
  } finally {
    loadingStores.value = false
  }
}

// Métodos de manejo de cambios en el formulario
const onDestinationTypeChange = () => {
  // Limpiar campos relacionados cuando cambia el tipo de destino
  consumptionForm.value.destinationID = ''
  consumptionForm.value.medicalCenterId = ''
  
  // Cargar datos según el tipo de destino
  if (consumptionForm.value.destinationType === 'pavilion') {
    if (medicalCenters.value.length === 0) {
      loadMedicalCenters()
    }
    if (pavilions.value.length === 0) {
      loadPavilions()
    }
  } else if (consumptionForm.value.destinationType === 'store') {
    if (stores.value.length === 0) {
      loadStores()
    }
  }
}

const onMedicalCenterChange = () => {
  // Limpiar pabellón seleccionado cuando cambia el centro médico
  consumptionForm.value.destinationID = ''
}

// Métodos principales
const scanQRCode = async () => {
  if (!qrInput.value.trim()) return
  
  scanning.value = true
  error.value = null
  scannedProduct.value = null
  
  try {
    const scanContext = buildScanContext()
    const result = await qrService.scanQRCode(qrInput.value.trim(), scanContext)
    
    if (!result) {
      error.value = 'No se recibió respuesta del servicio de QR'
      return
    }
    
    if (result.success && (result.data || result.type)) {
      const productData = result.data || result
      
      if (productData.type && productData.type !== 'medical_supply') {
        error.value = 'Solo se pueden consumir productos individuales (no lotes)'
        return
      }
      
      scannedProduct.value = productData
      lastScanContext.value = {
        ...scanContext,
        scanned_at: new Date(),
        user_name: currentUser.value?.name,
        location: currentLocation.value?.name
      }
    } else {
      error.value = result.error || result.message || 'Error desconocido al escanear código QR'
    }
  } catch (err) {
    handleError(err, 'Error al escanear código QR')
  } finally {
    scanning.value = false
  }
}

const consumeProduct = async () => {
  if (!scannedProduct.value) return
  
  consuming.value = true
  error.value = null
  
  try {
    const validation = validateConsumptionForm()
    
    if (!validation.isValid) {
      error.value = 'Faltan campos obligatorios: ' + validation.errors.join(', ')
      return
    }
    
    const consumptionData = {
      qr_code: scannedProduct.value.qr_code,
      user_rut: consumptionForm.value.userRUT,
      user_name: currentUser.value?.name || 'Encargado Bodega',
      destination_type: consumptionForm.value.destinationType,
      destination_id: validation.destinationIdNum,
      notes: consumptionForm.value.notes,
      consumption_purpose: selectedConsumptionPurpose.value,
      consumption_context: {
        pavilion_id: currentLocation.value?.pavilion_id,
        medical_center_id: currentLocation.value?.medical_center_id,
        user_agent: navigator.userAgent,
        scan_source: 'web'
      }
    }
    
    const result = await qrService.consumeIndividualSupply(consumptionData)
    
    if (result.success) {
      consumptionSuccess.value = {
        ...result.data,
        qr_code: scannedProduct.value.qr_code,
        batch_id: scannedProduct.value.supply_info?.batch?.id,
        traceability_info: result.traceability_info,
        batch_history_updated: true
      }
      
      updateConsumptionStats(consumptionData)
      clearScannedProduct()
      qrInput.value = ''
    } else {
      error.value = result.error || 'Error al consumir el producto'
    }
  } catch (err) {
    handleError(err, 'Error al consumir producto')
  } finally {
    consuming.value = false
  }
}

const startCameraScanner = () => {
  cameraActive.value = true
  alert('Funcionalidad de cámara en desarrollo. Por favor usa el input manual.')
  cameraActive.value = false
}

const viewTraceability = async (qrCode) => {
  try {
    await router.push({
      name: 'QRTraceability',
      params: { qrCode: qrCode }
    })
  } catch (err) {
    handleError(err, 'Error al navegar a la trazabilidad')
  }
}

const viewBatchHistory = async (batchId) => {
  if (!batchId) {
    console.warn('No se proporcionó ID de lote')
    return
  }

  try {
    await router.push({
      name: 'BatchHistory',
      params: { batchId: batchId }
    })
  } catch (err) {
    try {
      const history = await qrService.getBatchHistoryFormatted(batchId)
      const formattedHistory = history.map(entry => {
        const display = entry.display_format
        return `${display.date}\n${display.action}\nCant: ${display.previous_amount} → ${display.new_amount}\n${display.user_rut} ${display.user_name}`
      }).join('\n\n')
      
      alert(`Historial del Lote:\n\n${formattedHistory}`)
    } catch (historyError) {
      handleError(historyError, 'Error al cargar historial del lote')
    }
  }
}

const onLocationChanged = (location) => {
  currentLocation.value = location
  
  if (scannedProduct.value && location && !consumptionForm.value.destinationType) {
    // Auto-asignar tipo de destino, centro médico y pabellón si hay ubicación seleccionada
    if (location.pavilion_id) {
      consumptionForm.value.destinationType = 'pavilion'
      consumptionForm.value.medicalCenterId = location.medical_center_id?.toString() || ''
      consumptionForm.value.destinationID = location.pavilion_id.toString()
      
      // Asegurar que los datos estén cargados
      if (medicalCenters.value.length === 0) {
        loadMedicalCenters()
      }
      if (pavilions.value.length === 0) {
        loadPavilions()
      }
    }
  }
}

// Funciones de gestión de datos
const saveConsumptionStats = () => {
  const stats = {
    today_consumptions: consumptionStats.today,
    week_consumptions: consumptionStats.week,
    recent_consumptions: recentConsumptions.value,
    last_updated: new Date().toISOString()
  }
  localStorage.setItem('consumption-stats', JSON.stringify(stats))
}

const loadConsumptionStats = () => {
  try {
    const saved = localStorage.getItem('consumption-stats')
    if (saved) {
      const stats = JSON.parse(saved)
      consumptionStats.today = stats.today_consumptions || 0
      consumptionStats.week = stats.week_consumptions || 0
      recentConsumptions.value = stats.recent_consumptions || []
    }
  } catch (err) {
    console.error('Error loading consumption stats:', err)
  }
}

// Funciones de limpieza
const clearScannedProduct = () => {
  scannedProduct.value = null
  lastScanContext.value = null
  consumptionForm.value = {
    userRUT: currentUser.value?.rut || '',
    destinationType: '',
    destinationID: '',
    medicalCenterId: '',
    notes: ''
  }
}

const clearError = () => {
  error.value = null
}

const retryQRScan = async () => {
  if (qrInput.value) {
    clearError()
    await scanQRCode()
  }
}

const clearAll = () => {
  Object.assign(state, {
    qrInput: '',
    scannedProduct: null,
    consumptionSuccess: null,
    error: null,
    lastScanContext: null,
    isQuickMode: false
  })
  clearScannedProduct()
}

// Funciones de utilidad
const getPurposeButtonClass = (value) => {
  const baseClass = 'p-4 rounded-lg border-2 text-sm font-medium transition-all'
  return selectedConsumptionPurpose.value === value
    ? `${baseClass} border-red-500 bg-red-50 text-red-700`
    : `${baseClass} border-gray-200 bg-white text-gray-700 hover:border-gray-300 hover:bg-gray-50`
}

const formatDate = (dateString) => {
  if (!dateString) return 'No disponible'
  try {
    return format(new Date(dateString), 'dd/MM/yyyy HH:mm', { locale: es })
  } catch (error) {
    return dateString
  }
}

// Procesamiento de parámetros de URL
const processUrlParams = async () => {
  const urlParams = new URLSearchParams(window.location.search)
  const qrCode = urlParams.get('qr')
  const isQuick = urlParams.get('quick') === 'true'
  const purpose = urlParams.get('purpose')
  
  if (qrCode && isQuick && purpose === 'consume') {
    isQuickMode.value = true
    
    const purposeParam = urlParams.get('consumption_purpose')
    if (purposeParam && ['routine', 'emergency', 'maintenance', 'transfer'].includes(purposeParam)) {
      selectedConsumptionPurpose.value = purposeParam
    }
    
    qrInput.value = qrCode
    
    await new Promise(resolve => setTimeout(resolve, 500))
    
    try {
      await scanQRCode()
    } catch (err) {
      console.error('Error en escaneo automático:', err)
    }
    
    // Limpiar parámetros de URL
    const url = new URL(window.location)
    ;['qr', 'quick', 'purpose', 'consumption_purpose'].forEach(param => 
      url.searchParams.delete(param)
    )
    window.history.replaceState({}, document.title, url.pathname + 
      (url.searchParams.toString() ? '?' + url.searchParams.toString() : ''))
  }
}

// Lifecycle
onMounted(async () => {
  loadConsumptionStats()
  
  if (!authStore.isAuthenticated) {
    authStore.initializeAuth()
  }
  
  if (authStore.isAuthenticated && authStore.token && (!authStore.user || !authStore.user.name)) {
    try {
      await authStore.fetchProfile()
    } catch (err) {
      console.error('Error obteniendo perfil:', err)
    }
  }

  // Cargar datos de destinos iniciales
  await Promise.all([
    loadMedicalCenters(),
    loadPavilions(),
    loadStores()
  ])
  
  await processUrlParams()
})
</script>

<style scoped>
.form-input {
  @apply block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm;
}

.form-select {
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
</style>