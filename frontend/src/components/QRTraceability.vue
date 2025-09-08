<template>
  <div class="max-w-7xl mx-auto p-6">
    <!-- Header -->
    <div class="mb-6">
      <div class="flex items-center mb-4">
        <button
          @click="$router.go(-1)"
          class="mr-4 p-2 text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded-full"
        >
          <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
          </svg>
        </button>
        <div class="flex-1">
          <h1 class="text-2xl font-bold text-gray-900">Trazabilidad Completa</h1>
          <p class="text-gray-600 mt-1">Seguimiento completo del código QR {{ qrCode }}</p>
        </div>
        <div class="flex space-x-2">
          <button
            @click="loadTraceability"
            :disabled="loading"
            class="btn-secondary"
          >
            <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
            Actualizar
          </button>
          <button
            @click="exportData"
            :disabled="!traceabilityData"
            class="btn-primary"
          >
            <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
            Exportar
          </button>
        </div>
      </div>
      
      <!-- QR Code Info Card con estadísticas -->
      <div class="bg-gradient-to-r from-blue-50 to-indigo-50 border border-blue-200 rounded-lg p-6">
        <div class="flex items-center justify-between">
          <div class="flex items-center">
            <div class="bg-white p-3 rounded-lg border border-blue-300 shadow-sm">
              <svg class="h-8 w-8 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h.01M5 8h2a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v2a2 2 0 002 2zm0 0h2a2 2 0 012 2v2a2 2 0 01-2 2H5a2 2 0 01-2-2v-2a2 2 0 012-2z" />
              </svg>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-700">Código QR</p>
              <code class="text-lg font-mono text-blue-900 bg-white px-3 py-1 rounded border">{{ qrCode }}</code>
            </div>
          </div>
          
          <!-- Estadísticas Rápidas 
          <div v-if="scanStatistics" class="grid grid-cols-1 md:grid-cols-3 gap-4">
            <div class="text-center">
              <div class="text-2xl font-bold text-blue-900">{{ scanStatistics.total_scans }}</div>
              <div class="text-xs text-blue-700">Total Escaneos</div>
            </div>
            <div class="text-center">
              <div class="text-2xl font-bold text-green-900">{{ scanStatistics.unique_scanners }}</div>
              <div class="text-xs text-green-700">Usuarios Únicos</div>
            </div>
            <div class="text-center">
              <div class="text-2xl font-bold text-purple-900">{{ scanStatistics.locations_visited }}</div>
              <div class="text-xs text-purple-700">Ubicaciones</div>
            </div>
          </div> -->
        </div>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="flex items-center justify-center py-12">
      <div class="text-center">
        <svg class="animate-spin -ml-1 mr-3 h-8 w-8 text-blue-600 mx-auto" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
        <p class="text-gray-600 mt-2">Cargando trazabilidad completa...</p>
      </div>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="bg-red-50 border border-red-200 rounded-lg p-6">
      <div class="flex items-center">
        <svg class="h-6 w-6 text-red-400 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <div>
          <h3 class="text-red-800 font-medium">Error al cargar trazabilidad</h3>
          <p class="text-red-700 text-sm mt-1">{{ error }}</p>
        </div>
      </div>
    </div>

    <!-- Content -->
    <div v-else-if="traceabilityData">
      <!-- Información general del insumo -->
      <div v-if="qrInfo" class="bg-white rounded-lg shadow border p-6 mb-6">
        <h2 class="text-lg font-semibold text-gray-900 mb-4">Información del Insumo</h2>
        
        <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
          <!-- Información básica -->
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700">Tipo</label>
              <span :class="getTypeIconClass(qrInfo.type)" class="inline-flex items-center px-2 py-1 text-sm font-medium rounded-full mt-1">
                <svg class="h-4 w-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path v-if="qrInfo.type === 'medical_supply'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                  <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
                </svg>
                {{ getTypeLabel(qrInfo.type) }}
              </span>
            </div>
            
            <div v-if="qrInfo.supply_info">
              <label class="block text-sm font-medium text-gray-700">Estado Actual</label>
              <span :class="getStatusBadgeClass()" class="inline-flex px-2 py-1 text-sm font-semibold rounded-full mt-1">
                {{ getStatusLabel() }}
              </span>
            </div>
          </div>

          <!-- Información del código de insumo -->
          <div v-if="qrInfo.supply_code" class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700">Código Insumo</label>
              <p class="text-sm text-gray-900 mt-1">{{ qrInfo.supply_code.code }}</p>
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700">Nombre</label>
              <p class="text-sm text-gray-900 mt-1">{{ qrInfo.supply_code.name }}</p>
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700">Código Proveedor</label>
              <p class="text-sm text-gray-900 mt-1">{{ qrInfo.supply_code.code_supplier }}</p>
            </div>
          </div>

          <!-- Información del lote -->
          <div v-if="qrInfo.supply_info?.batch" class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700">Lote</label>
              <p class="text-sm text-gray-900 mt-1">ID: {{ qrInfo.supply_info.batch.id }}</p>
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700">Proveedor</label>
              <p class="text-sm text-gray-900 mt-1">{{ qrInfo.supply_info.batch.supplier }}</p>
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700">Fecha de Vencimiento</label>
              <p class="text-sm text-gray-900 mt-1">{{ formatDate(qrInfo.supply_info.batch.expiration_date) }}</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Filtros -->
      <div class="bg-white rounded-lg shadow-sm border p-6 mb-6">
        <h3 class="text-lg font-medium text-gray-900 mb-4">Filtros de Eventos</h3>
        <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Tipo de Evento</label>
            <select v-model="selectedEventType" class="form-select w-full">
              <option value="">Todos los eventos</option>
              <option value="scan">Escaneos</option>
              <option value="movement">Movimientos</option>
              <option value="status_change">Cambios de Estado</option>
              <option value="request">Solicitudes</option>
            </select>
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Usuario</label>
            <select v-model="selectedUser" class="form-select w-full">
              <option value="">Todos los usuarios</option>
              <option
                v-for="user in uniqueUsers"
                :key="user"
                :value="user"
              >
                {{ user }}
              </option>
            </select>
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Ubicación</label>
            <select v-model="selectedLocation" class="form-select w-full">
              <option value="">Todas las ubicaciones</option>
              <option
                v-for="location in uniqueLocations"
                :key="location"
                :value="location"
              >
                {{ location }}
              </option>
            </select>
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Período</label>
            <select v-model="selectedPeriod" class="form-select w-full">
              <option value="">Todo el tiempo</option>
              <option value="today">Hoy</option>
              <option value="week">Última semana</option>
              <option value="month">Último mes</option>
            </select>
          </div>
        </div>
      </div>

      <!-- Estadísticas Detalladas -->
      <div v-if="scanStatistics" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-6">
        <div class="bg-white rounded-lg shadow-sm border p-6">
          <div class="flex items-center">
            <div class="bg-blue-100 rounded-lg p-3">
              <svg class="h-6 w-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
              </svg>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600">Total Escaneos</p>
              <p class="text-2xl font-bold text-gray-900">{{ scanStatistics.total_scans }}</p>
            </div>
          </div>
          <div class="mt-4">
            <div class="flex items-center text-sm">
              <span class="text-green-600 font-medium">{{ scanStatistics.successful_scans }}</span>
              <span class="text-gray-500 ml-1">exitosos</span>
              <span class="text-red-600 font-medium ml-4">{{ scanStatistics.error_scans }}</span>
              <span class="text-gray-500 ml-1">errores</span>
            </div>
          </div>
        </div>

        <div class="bg-white rounded-lg shadow-sm border p-6">
          <div class="flex items-center">
            <div class="bg-green-100 rounded-lg p-3">
              <svg class="h-6 w-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.5 2.5 0 11-5 0 2.5 2.5 0 015 0z" />
              </svg>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600">Usuarios Únicos</p>
              <p class="text-2xl font-bold text-gray-900">{{ scanStatistics.unique_scanners }}</p>
            </div>
          </div>
          <div class="mt-4">
            <div class="text-sm text-gray-500">
              {{ Math.round(scanStatistics.total_scans / scanStatistics.unique_scanners * 10) / 10 }} escaneos promedio por usuario
            </div>
          </div>
        </div>

        <div class="bg-white rounded-lg shadow-sm border p-6">
          <div class="flex items-center">
            <div class="bg-purple-100 rounded-lg p-3">
              <svg class="h-6 w-6 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
              </svg>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600">Ubicaciones</p>
              <p class="text-2xl font-bold text-gray-900">{{ scanStatistics.locations_visited }}</p>
            </div>
          </div>
          <div class="mt-4">
            <div class="text-sm text-gray-500">
              Visitadas desde el primer escaneo
            </div>
          </div>
        </div>

        <div class="bg-white rounded-lg shadow-sm border p-6">
          <div class="flex items-center">
            <div class="bg-yellow-100 rounded-lg p-3">
              <svg class="h-6 w-6 text-yellow-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600">Tiempo en Sistema</p>
              <p class="text-2xl font-bold text-gray-900">{{ Math.round(scanStatistics.hours_in_system) }}h</p>
            </div>
          </div>
          <div class="mt-4">
            <div class="text-sm text-gray-500">
              Desde: {{ formatDate(scanStatistics.first_scan) }}
            </div>
          </div>
        </div>
      </div>

      <!-- Estado actual y resumen 
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
        <!-- Estado actual 
        <div class="bg-white rounded-lg shadow border p-6">
          <h3 class="text-lg font-semibold text-gray-900 mb-4">Estado Actual</h3>
          
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700">Estado</label>
              <span :class="getCurrentStatusBadgeClass()" class="inline-flex px-3 py-1 text-sm font-semibold rounded-full mt-1">
                {{ traceabilityData.current_status || 'Desconocido' }}
              </span>
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700">Asignado a Solicitud</label>
              <p class="text-sm mt-1">
                {{ traceabilityData.is_assigned_to_request ? 'Sí' : 'No' }}
              </p>
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700">Última Actualización</label>
              <p class="text-sm text-gray-900 mt-1">{{ formatDate(traceabilityData.last_updated) }}</p>
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700">Fecha de Creación</label>
              <p class="text-sm text-gray-900 mt-1">{{ formatDate(traceabilityData.created_date) }}</p>
            </div>
          </div>
        </div>

        <!-- Estadísticas
        <div class="bg-white rounded-lg shadow border p-6">
          <h3 class="text-lg font-semibold text-gray-900 mb-4">Estadísticas de Historial</h3>
          
          <div class="space-y-4">
            <div class="flex justify-between items-center">
              <span class="text-sm text-gray-600">Total de Asignaciones:</span>
              <span class="text-sm font-medium text-gray-900">{{ traceabilityData.request_history?.length || 0 }}</span>
            </div>
            
            <div class="flex justify-between items-center">
              <span class="text-sm text-gray-600">Movimientos de Supply:</span>
              <span class="text-sm font-medium text-gray-900">{{ traceabilityData.supply_history?.length || 0 }}</span>
            </div>
            
            <div class="flex justify-between items-center">
              <span class="text-sm text-gray-600">Eventos de Escaneo:</span>
              <span class="text-sm font-medium text-gray-900">{{ traceabilityData.scan_history?.length || 0 }}</span>
            </div>
            
            <div class="flex justify-between items-center">
              <span class="text-sm text-gray-600">Tiempo en el Sistema:</span>
              <span class="text-sm font-medium text-gray-900">{{ calculateTimeInSystem() }}</span>
            </div>
          </div>
        </div>
      </div> -->

      <!-- Timeline de Eventos -->
      <div class="bg-white rounded-lg shadow-sm border">
        <div class="p-6 border-b border-gray-200">
          <h3 class="text-lg font-medium text-gray-900">Timeline de Eventos</h3>
          <p class="text-sm text-gray-600 mt-1">
            {{ filteredEvents.length }} eventos
            <span v-if="selectedEventType || selectedUser || selectedLocation || selectedPeriod">
              (filtrados de {{ allEvents.length }} total)
            </span>
          </p>
        </div>
        
        <div class="p-6">
          <!-- Historial de Escaneos -->
          <ScanEventHistory 
            :events="filteredScanEvents"
            :show-filters="false"
            class="mb-8"
            @view-details="handleScanEventDetails"
            @scan-again="handleScanAgain"
          />
          
          <!-- Timeline Completo -->
          <div class="flow-root">
            <ul class="-mb-8">
              <li
                v-for="(event, index) in filteredEvents"
                :key="index"
                class="relative pb-8"
              >
                <span
                  v-if="index !== filteredEvents.length - 1"
                  class="absolute top-4 left-4 -ml-px h-full w-0.5 bg-gray-200"
                  aria-hidden="true"
                ></span>
                
                <div class="relative flex space-x-3">
                  <!-- Icon -->
                  <div>
                    <span
                      :class="[
                        'h-8 w-8 rounded-full flex items-center justify-center ring-8 ring-white',
                        getEventIconClass(event)
                      ]"
                    >
                      <svg class="h-4 w-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" :d="getEventIcon(event)" />
                      </svg>
                    </span>
                  </div>
                  
                  <!-- Content -->
                  <div class="flex-1 min-w-0">
                    <div class="flex items-center justify-between">
                      <div>
                        <p class="text-sm font-medium text-gray-900">
                          {{ getEventTitle(event) }}
                        </p>
                        <p class="text-sm text-gray-500">
                          {{ getEventDescription(event) }}
                        </p>
                      </div>
                      <div class="text-right">
                        <time class="text-sm text-gray-500">
                          {{ formatDate(event.date_time || event.scanned_at || event.assigned_date) }}
                        </time>
                        <div v-if="event.user_name || event.scanned_by_name" class="text-xs text-gray-400">
                          {{ event.user_name || event.scanned_by_name }}
                        </div>
                      </div>
                    </div>
                    
                    <!-- Detalles adicionales -->
                    <div v-if="event.current_location || event.pavilion_name" class="mt-2">
                      <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-800">
                        <svg class="h-3 w-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                        </svg>
                        {{ event.current_location || event.pavilion_name }}
                      </span>
                    </div>
                    
                    <!-- Contexto del escaneo -->
                    <div v-if="event.scan_purpose" class="mt-2">
                      <span 
                        :class="[
                          'inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium',
                          getScanPurposeBadgeClass(event.scan_purpose)
                        ]"
                      >
                        {{ getScanPurposeLabel(event.scan_purpose) }}
                      </span>
                    </div>
                    
                    <!-- Error message para escaneos fallidos -->
                    <div v-if="event.error_message" class="mt-2 p-2 bg-red-50 border border-red-200 rounded text-sm">
                      <div class="text-red-800 font-medium">Error:</div>
                      <div class="text-red-700">{{ event.error_message }}</div>
                    </div>
                  </div>
                </div>
              </li>
            </ul>
          </div>
        </div>
      </div>

      <!-- Información adicional si no hay historial -->
      <div v-if="allEvents.length === 0" class="bg-yellow-50 border border-yellow-200 rounded-lg p-6 mt-6">
        <div class="flex">
          <svg class="h-5 w-5 text-yellow-400 mr-3 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <div>
            <h3 class="text-sm font-medium text-yellow-800">Sin historial disponible</h3>
            <p class="text-sm text-yellow-700 mt-2">
              Este código QR no tiene eventos registrados aún. Esto puede significar que es un insumo recién creado que no ha sido procesado en el sistema.
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { format, isToday, isThisWeek, isThisMonth, differenceInDays, differenceInHours } from 'date-fns'
import { es } from 'date-fns/locale'
import qrService from '@/services/qrService'
import ScanEventHistory from '@/components/ScanEventHistory.vue'

const route = useRoute()
const router = useRouter()

// Props
const props = defineProps({
  qrCode: {
    type: String,
    required: true
  }
})

// Estado reactivo
const loading = ref(false)
const error = ref(null)
const traceabilityData = ref(null)
const scanStatistics = ref(null)
const qrInfo = ref(null)

// Filtros
const selectedEventType = ref('')
const selectedUser = ref('')
const selectedLocation = ref('')
const selectedPeriod = ref('')

// Computed properties
const qrCode = computed(() => props.qrCode || route.params.qrcode)

const allEvents = computed(() => {
  if (!traceabilityData.value) return []
  
  const events = []
  
  // Agregar escaneos
  if (traceabilityData.value.scan_history) {
    events.push(...traceabilityData.value.scan_history.map(scan => ({
      ...scan,
      event_type: 'scan',
      date_time: scan.scanned_at
    })))
  }
  
  // Agregar historial de movimientos
  if (traceabilityData.value.supply_history) {
    events.push(...traceabilityData.value.supply_history.map(movement => ({
      ...movement,
      event_type: 'movement'
    })))
  }
  
  // Agregar historial de solicitudes
  if (traceabilityData.value.request_history) {
    events.push(...traceabilityData.value.request_history.map(request => ({
      ...request,
      event_type: 'request',
      date_time: request.assigned_date
    })))
  }
  
  // Ordenar por fecha (más reciente primero)
  return events.sort((a, b) => {
    const dateA = new Date(a.date_time || a.scanned_at || a.assigned_date)
    const dateB = new Date(b.date_time || b.scanned_at || b.assigned_date)
    return dateB - dateA
  })
})

const filteredEvents = computed(() => {
  let events = allEvents.value
  
  // Filtrar por tipo de evento
  if (selectedEventType.value) {
    events = events.filter(event => event.event_type === selectedEventType.value)
  }
  
  // Filtrar por usuario
  if (selectedUser.value) {
    events = events.filter(event => 
      event.scanned_by_name === selectedUser.value || 
      event.user_name === selectedUser.value ||
      event.assigned_by_name === selectedUser.value
    )
  }
  
  // Filtrar por ubicación
  if (selectedLocation.value) {
    events = events.filter(event => 
      event.current_location === selectedLocation.value ||
      event.pavilion_name === selectedLocation.value
    )
  }
  
  // Filtrar por período
  if (selectedPeriod.value) {
    events = events.filter(event => {
      const eventDate = new Date(event.date_time || event.scanned_at || event.assigned_date)
      switch (selectedPeriod.value) {
        case 'today':
          return isToday(eventDate)
        case 'week':
          return isThisWeek(eventDate)
        case 'month':
          return isThisMonth(eventDate)
        default:
          return true
      }
    })
  }
  
  return events
})

const filteredScanEvents = computed(() => {
  return filteredEvents.value.filter(event => event.event_type === 'scan')
})

const uniqueUsers = computed(() => {
  const users = new Set()
  allEvents.value.forEach(event => {
    if (event.scanned_by_name) users.add(event.scanned_by_name)
    if (event.user_name) users.add(event.user_name)
    if (event.assigned_by_name) users.add(event.assigned_by_name)
  })
  return Array.from(users).sort()
})

const uniqueLocations = computed(() => {
  const locations = new Set()
  allEvents.value.forEach(event => {
    if (event.current_location) locations.add(event.current_location)
    if (event.pavilion_name) locations.add(event.pavilion_name)
  })
  return Array.from(locations).sort()
})

// Funciones principales
const loadTraceability = async () => {
  if (!qrCode.value) {
    error.value = 'Código QR requerido'
    return
  }
  
  loading.value = true
  error.value = null
  
  try {
    // Cargar información básica del QR
    const qrResult = await qrService.scanQRCode(qrCode.value)
    qrInfo.value = qrResult

    // Cargar trazabilidad completa
    const traceabilityResponse = await qrService.getCompleteTraceability(qrCode.value)
    traceabilityData.value = traceabilityResponse
    
    // Cargar estadísticas de escaneo
    const statsResponse = await qrService.getScanStatistics(qrCode.value)
    scanStatistics.value = statsResponse
    
  } catch (err) {
    console.error('Error loading traceability:', err)
    error.value = err.message || 'Error al cargar la trazabilidad'
  } finally {
    loading.value = false
  }
}

const exportData = () => {
  if (!traceabilityData.value) return
  
  const data = {
    qr_code: qrCode.value,
    export_date: new Date().toISOString(),
    scan_statistics: scanStatistics.value,
    traceability: traceabilityData.value,
    qr_info: qrInfo.value,
    filtered_events: filteredEvents.value,
    filters_applied: {
      event_type: selectedEventType.value,
      user: selectedUser.value,
      location: selectedLocation.value,
      period: selectedPeriod.value
    }
  }
  
  const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `traceability-${qrCode.value}-${format(new Date(), 'yyyy-MM-dd')}.json`
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  URL.revokeObjectURL(url)
}

const handleScanEventDetails = (event) => {
  // Mostrar detalles del evento de escaneo en un modal o navegación
  console.log('Ver detalles del evento:', event)
}

const handleScanAgain = (qrCode) => {
  // Redirigir al escáner con el QR code
  router.push({
    name: 'QRScanner',
    query: { qr: qrCode }
  })
}

// Funciones de utilidad
const calculateTimeInSystem = () => {
  if (!traceabilityData.value?.created_date) return 'N/A'
  
  try {
    const created = new Date(traceabilityData.value.created_date)
    const now = new Date()
    const days = differenceInDays(now, created)
    
    if (days === 0) {
      const hours = differenceInHours(now, created)
      return hours === 0 ? 'Menos de 1 hora' : `${hours} hora${hours > 1 ? 's' : ''}`
    } else {
      return `${days} día${days > 1 ? 's' : ''}`
    }
  } catch {
    return 'N/A'
  }
}

const formatDate = (dateString) => {
  if (!dateString) return 'N/A'
  try {
    return format(new Date(dateString), 'dd/MM/yyyy HH:mm', { locale: es })
  } catch {
    return dateString
  }
}

// Funciones de estilo y etiquetas (manteniendo las originales)
const getTypeLabel = (type) => {
  const labels = {
    'batch': 'Lote de Productos',
    'medical_supply': 'Insumo Individual'
  }
  return labels[type] || type
}

const getTypeIconClass = (type) => {
  const classes = {
    'batch': 'bg-blue-100 text-blue-600',
    'medical_supply': 'bg-green-100 text-green-600'
  }
  return classes[type] || 'bg-gray-100 text-gray-600'
}

const getStatusLabel = () => {
  if (!qrInfo.value) return 'Desconocido'
  return qrInfo.value.is_consumed ? 'Consumido' : 'Disponible'
}

const getStatusBadgeClass = () => {
  if (!qrInfo.value) return 'bg-gray-100 text-gray-800'
  return qrInfo.value.is_consumed ? 'bg-red-100 text-red-800' : 'bg-green-100 text-green-800'
}

const getCurrentStatusBadgeClass = () => {
  const status = traceabilityData.value?.current_status || 'unknown'
  const classes = {
    'available': 'bg-green-100 text-green-800',
    'assigned': 'bg-blue-100 text-blue-800',
    'delivered': 'bg-purple-100 text-purple-800',
    'consumed': 'bg-gray-100 text-gray-800',
    'expired': 'bg-red-100 text-red-800',
    'unknown': 'bg-gray-100 text-gray-800'
  }
  return classes[status] || classes.unknown
}

// Funciones para eventos del timeline
const getEventTitle = (event) => {
  switch (event.event_type) {
    case 'scan':
      return 'Código QR Escaneado'
    case 'movement':
      if (event.status === 'consumido' && event.destination_name) {
        return 'Producto Consumido'
      }
      return `Movimiento: ${event.movement_type || event.status || 'Cambio'}`
    case 'request':
      return 'Asignado a Solicitud'
    default:
      return 'Evento'
  }
}

const getEventDescription = (event) => {
  switch (event.event_type) {
    case 'scan':
      return `Escaneado por ${event.scanned_by_name || 'Usuario desconocido'} desde ${event.scan_source || 'web'}`
    case 'movement':
      // Si hay información de destino, mostrarla
      if (event.destination_name) {
        const destinationType = event.destination_type === 'pavilion' ? 'Pabellón' : 'Almacén'
        let description = `${event.status || 'Procesado'} - Enviado a ${destinationType}: ${event.destination_name}`
        if (event.medical_center_name) {
          description += ` (${event.medical_center_name})`
        }
        return description
      }
      return event.observations || `Estado: ${event.status}`
    case 'request':
      return `Solicitud #${event.supply_request?.request_number || event.id || 'N/A'}`
    default:
      return 'Información no disponible'
  }
}

const getEventIcon = (event) => {
  switch (event.event_type) {
    case 'scan':
      return 'M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h.01M5 8h2a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v2a2 2 0 002 2zm0 0h2a2 2 0 012 2v2a2 2 0 01-2 2H5a2 2 0 01-2-2v-2a2 2 0 012-2z'
    case 'movement':
      return 'M13 10V3L4 14h7v7l9-11h-7z'
    case 'request':
      return 'M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2'
    default:
      return 'M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z'
  }
}

const getEventIconClass = (event) => {
  switch (event.event_type) {
    case 'scan':
      return 'bg-blue-500'
    case 'movement':
      return 'bg-green-500'
    case 'request':
      return 'bg-purple-500'
    default:
      return 'bg-gray-500'
  }
}

const getScanPurposeBadgeClass = (purpose) => {
  const classes = {
    'lookup': 'bg-green-100 text-green-800',
    'consume': 'bg-red-100 text-red-800',
    'verify': 'bg-blue-100 text-blue-800',
    'inventory_check': 'bg-purple-100 text-purple-800',
    'assign': 'bg-yellow-100 text-yellow-800'
  }
  return classes[purpose] || 'bg-gray-100 text-gray-800'
}

const getScanPurposeLabel = (purpose) => {
  const labels = {
    'lookup': 'Consulta',
    'consume': 'Consumo',
    'verify': 'Verificación',
    'inventory_check': 'Inventario',
    'assign': 'Asignación'
  }
  return labels[purpose] || purpose
}

// Lifecycle
onMounted(() => {
  if (qrCode.value) {
    loadTraceability()
  } else {
    error.value = 'No se proporcionó código QR'
  }
})
</script>

<style scoped>
.form-select {
  @apply block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm;
}

.btn-primary {
  @apply inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed;
}

.btn-secondary {
  @apply inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md shadow-sm text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed;
}

.flow-root {
  overflow: hidden;
}
</style>