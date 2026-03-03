<template>
  <div class="max-w-5xl mx-auto px-4 py-6 space-y-5">

    <!-- Breadcrumb -->
    <nav class="flex items-center space-x-2 text-sm" aria-label="Breadcrumb">
      <router-link to="/qr" class="text-blue-600 hover:text-blue-800 font-medium">Escáner QR</router-link>
      <svg class="h-4 w-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
      </svg>
      <span class="text-gray-500 font-mono">{{ qrCode }}</span>
    </nav>

    <!-- Loading -->
    <div v-if="loading" class="flex justify-center items-center h-64">
      <div class="animate-spin rounded-full h-10 w-10 border-b-2 border-blue-600"></div>
    </div>

    <!-- Error -->
    <div v-else-if="error" class="bg-red-50 border border-red-200 rounded-xl p-8 text-center">
      <svg class="h-12 w-12 text-red-400 mx-auto mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>
      <h3 class="text-base font-semibold text-red-800 mb-1">Error al cargar información</h3>
      <p class="text-sm text-red-600 mb-4">{{ error }}</p>
      <div class="flex justify-center gap-3">
        <button @click="refreshData" class="btn-secondary text-sm">Intentar de Nuevo</button>
        <router-link to="/qr" class="btn-secondary text-sm">Volver al Escáner</router-link>
      </div>
    </div>

    <!-- Contenido principal -->
    <template v-else-if="qrInfo">

      <!-- ── HERO CARD ─────────────────────────────────────────── -->
      <div class="bg-white rounded-xl shadow-sm border overflow-hidden">
        <!-- Franja de color según tipo -->
        <div :class="qrInfo.type === 'batch' ? 'bg-indigo-600' : 'bg-blue-600'" class="h-1.5 w-full"></div>

        <div class="p-6 flex flex-col sm:flex-row gap-6">
          <!-- Imagen QR -->
          <div class="flex-shrink-0 flex flex-col items-center gap-2">
            <div class="border-2 border-gray-200 rounded-xl p-2 bg-gray-50">
              <img
                v-if="qrInfo.qr_code"
                :src="getQRImageUrl(qrInfo.qr_code)"
                alt="Código QR"
                class="w-28 h-28 object-contain"
              />
            </div>
            <span class="text-xs font-mono text-gray-500 bg-gray-100 px-2 py-0.5 rounded select-all">
              {{ qrInfo.qr_code }}
            </span>
          </div>

          <!-- Info principal -->
          <div class="flex-1 min-w-0">
            <div class="flex flex-wrap items-start justify-between gap-3 mb-3">
              <div>
                <p class="text-xs font-medium text-gray-400 uppercase tracking-wide mb-0.5">
                  {{ getTypeLabel(qrInfo.type) }}
                </p>
                <h1 class="text-xl sm:text-2xl font-bold text-gray-900 leading-tight">
                  {{ qrInfo.supply_info?.name || qrInfo.supply_info?.SupplyCode?.name || qrInfo.batch_info?.supplier || 'Sin nombre' }}
                </h1>
              </div>
              <!-- Badge estado -->
              <span :class="getStatusClass(qrInfo.current_status || qrInfo.status)"
                    class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-full text-sm font-semibold flex-shrink-0">
                <span class="text-base">{{ getStatusIcon(qrInfo.current_status || qrInfo.status) }}</span>
                {{ getStatusText(qrInfo.current_status || qrInfo.status) }}
              </span>
            </div>

            <!-- Datos clave en fila de chips -->
            <div class="flex flex-wrap gap-2 mb-4">
              <!-- Para insumos individuales -->
              <template v-if="qrInfo.type === 'medical_supply' && qrInfo.supply_info">
                <span class="detail-chip">
                  <span class="detail-chip-label">Código</span>
                  <span class="font-mono">{{ qrInfo.supply_info.SupplyCode?.code || qrInfo.supply_info.Code || '—' }}</span>
                </span>
                <span class="detail-chip">
                  <span class="detail-chip-label">Proveedor</span>
                  {{ qrInfo.supply_info.batch?.supplier || qrInfo.supply_info.supplier || '—' }}
                </span>
                <span class="detail-chip">
                  <span class="detail-chip-label">Almacén</span>
                  {{ qrInfo.supply_info.store_name || qrInfo.supply_info.batch?.store_name || '—' }}
                </span>
                <span class="detail-chip">
                  <span class="detail-chip-label">Vence</span>
                  {{ formatDate(qrInfo.supply_info.batch?.expiration_date || qrInfo.supply_info.expiration_date) }}
                </span>
              </template>
              <!-- Para lotes -->
              <template v-else-if="qrInfo.type === 'batch' && qrInfo.batch_info">
                <span class="detail-chip">
                  <span class="detail-chip-label">ID Lote</span>
                  {{ qrInfo.batch_info.id || '—' }}
                </span>
                <span class="detail-chip">
                  <span class="detail-chip-label">Cantidad</span>
                  {{ qrInfo.batch_info.amount }} u.
                </span>
                <span class="detail-chip">
                  <span class="detail-chip-label">Proveedor</span>
                  {{ qrInfo.batch_info.supplier || '—' }}
                </span>
                <span class="detail-chip">
                  <span class="detail-chip-label">Almacén</span>
                  {{ qrInfo.batch_info.store_name || qrInfo.batch_info.store_id || '—' }}
                </span>
                <span class="detail-chip">
                  <span class="detail-chip-label">Vence</span>
                  {{ formatDate(qrInfo.batch_info.expiration_date) }}
                </span>
              </template>
            </div>

            <!-- Acciones -->
            <div class="flex flex-wrap gap-2">
              <button
                @click="downloadQRAsPDF"
                :disabled="isGeneratingPDF"
                class="inline-flex items-center gap-1.5 px-4 py-2 rounded-lg border-2 border-green-600 text-green-700 text-sm font-semibold bg-white hover:bg-green-50 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
              >
                <svg v-if="isGeneratingPDF" class="animate-spin h-4 w-4" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                <svg v-else class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                </svg>
                {{ isGeneratingPDF ? 'Generando...' : 'Descargar PDF' }}
              </button>
              <button
                @click="refreshData"
                :disabled="loading"
                class="inline-flex items-center gap-1.5 px-4 py-2 rounded-lg border border-gray-300 text-gray-600 text-sm font-medium bg-white hover:bg-gray-50 transition-colors"
              >
                <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                </svg>
                Actualizar
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- ── GRID PRINCIPAL ────────────────────────────────────── -->
      <div class="grid lg:grid-cols-3 gap-5">

        <!-- Col izquierda (2/3): Lote relacionado + Historial -->
        <div class="lg:col-span-2 space-y-5">

          <!-- Lote relacionado (solo para insumos individuales) -->
          <div
            v-if="qrInfo.type === 'medical_supply' && (qrInfo.supply_info?.batch || qrInfo.supply_info?.BatchID)"
            class="bg-white rounded-xl shadow-sm border overflow-hidden"
          >
            <div class="flex items-center justify-between px-5 py-4 border-b bg-gray-50">
              <h2 class="text-sm font-semibold text-gray-700 flex items-center gap-2">
                <svg class="h-4 w-4 text-indigo-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
                </svg>
                Lote Relacionado
              </h2>
              <button
                v-if="qrInfo.supply_info?.batch?.id || qrInfo.supply_info?.BatchID"
                @click="viewBatch(qrInfo.supply_info?.batch?.id || qrInfo.supply_info?.BatchID)"
                class="text-xs text-blue-600 hover:text-blue-800 font-medium"
              >
                Ver en Inventario →
              </button>
            </div>
            <div class="px-5 py-4 grid grid-cols-2 sm:grid-cols-3 gap-x-6 gap-y-4">
              <div>
                <p class="text-xs text-gray-400 mb-0.5">ID del Lote</p>
                <p class="text-sm font-semibold font-mono text-gray-800">{{ qrInfo.supply_info?.batch?.id || qrInfo.supply_info?.BatchID || '—' }}</p>
              </div>
              <div>
                <p class="text-xs text-gray-400 mb-0.5">Proveedor</p>
                <p class="text-sm font-medium text-gray-800">{{ qrInfo.supply_info?.batch?.supplier || '—' }}</p>
              </div>
              <div>
                <p class="text-xs text-gray-400 mb-0.5">Vencimiento</p>
                <p class="text-sm font-medium text-gray-800">{{ formatDate(qrInfo.supply_info?.batch?.expiration_date) }}</p>
              </div>
              <div>
                <p class="text-xs text-gray-400 mb-0.5">Cantidad en lote</p>
                <p class="text-sm font-bold" :class="(qrInfo.supply_info?.batch?.amount || 0) > 0 ? 'text-green-600' : 'text-red-600'">
                  {{ qrInfo.supply_info?.batch?.amount || 0 }} unidades
                </p>
              </div>
              <div>
                <p class="text-xs text-gray-400 mb-0.5">Stock</p>
                <span :class="(qrInfo.supply_info?.batch?.amount || 0) > 0 ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'"
                      class="inline-block px-2 py-0.5 rounded-full text-xs font-semibold">
                  {{ (qrInfo.supply_info?.batch?.amount || 0) > 0 ? 'Disponible' : 'Agotado' }}
                </span>
              </div>
              <div v-if="qrInfo.supply_info?.batch?.qr_code">
                <p class="text-xs text-gray-400 mb-0.5">QR del Lote</p>
                <p class="text-xs font-mono text-gray-600 break-all">{{ qrInfo.supply_info?.batch?.qr_code }}</p>
              </div>
            </div>
          </div>

          <!-- Estadísticas del lote (solo para tipo batch) -->
          <div v-if="qrInfo.type === 'batch' && qrInfo.batch_status" class="bg-white rounded-xl shadow-sm border overflow-hidden">
            <div class="px-5 py-4 border-b bg-gray-50">
              <h2 class="text-sm font-semibold text-gray-700">Estadísticas del Lote</h2>
            </div>
            <div class="px-5 py-4">
              <div class="grid grid-cols-4 gap-4 mb-4">
                <div class="text-center">
                  <p class="text-2xl font-bold text-blue-600">{{ qrInfo.batch_status.total_individual_supplies || 0 }}</p>
                  <p class="text-xs text-gray-500 mt-0.5">Total</p>
                </div>
                <div class="text-center">
                  <p class="text-2xl font-bold text-green-600">{{ qrInfo.batch_status.available_supplies || 0 }}</p>
                  <p class="text-xs text-gray-500 mt-0.5">Disponibles</p>
                </div>
                <div class="text-center">
                  <p class="text-2xl font-bold text-red-500">{{ qrInfo.batch_status.consumed_supplies || 0 }}</p>
                  <p class="text-xs text-gray-500 mt-0.5">Consumidos</p>
                </div>
                <div class="text-center">
                  <p class="text-2xl font-bold text-purple-600">{{ getUsagePercentage() }}%</p>
                  <p class="text-xs text-gray-500 mt-0.5">Utilización</p>
                </div>
              </div>
              <div class="bg-gray-100 rounded-full h-2">
                <div
                  class="bg-blue-500 h-2 rounded-full transition-all duration-500"
                  :style="`width: ${getUsagePercentage()}%`"
                ></div>
              </div>
            </div>
          </div>

          <!-- Historial de movimientos -->
          <div class="bg-white rounded-xl shadow-sm border overflow-hidden">
            <div class="flex items-center justify-between px-5 py-4 border-b bg-gray-50">
              <h2 class="text-sm font-semibold text-gray-700 flex items-center gap-2">
                <svg class="h-4 w-4 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                Historial de Movimientos
              </h2>
              <button @click="loadHistory" class="text-xs text-gray-500 hover:text-gray-700 font-medium flex items-center gap-1">
                <svg class="h-3.5 w-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                </svg>
                Actualizar
              </button>
            </div>

            <div v-if="historyLoading" class="flex justify-center items-center py-10">
              <div class="animate-spin rounded-full h-7 w-7 border-b-2 border-blue-500"></div>
            </div>

            <div v-else-if="historyError" class="text-center py-8 px-5">
              <p class="text-sm text-red-500 mb-2">{{ historyError }}</p>
              <button @click="loadHistory" class="btn-secondary text-xs">Reintentar</button>
            </div>

            <!-- Timeline -->
            <div v-else-if="historyData && historyData.length > 0" class="px-5 py-4">
              <div class="relative">
                <!-- Línea vertical -->
                <div class="absolute left-4 top-2 bottom-2 w-0.5 bg-gray-200"></div>

                <div class="space-y-5">
                  <div v-for="(item, index) in historyData" :key="index" class="relative flex gap-4">
                    <!-- Punto del timeline -->
                    <div :class="getHistoryIconClass(item)"
                         class="relative z-10 flex-shrink-0 w-8 h-8 rounded-full flex items-center justify-center shadow-sm border-2 border-white">
                      <svg class="h-3.5 w-3.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path v-if="getHistoryStatus(item) === 'consumido'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                        <path v-else-if="getHistoryStatus(item) === 'recepcionado'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                        <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 9l3 3-3 3M6 12h10" />
                      </svg>
                    </div>

                    <!-- Contenido -->
                    <div class="flex-1 bg-gray-50 rounded-lg p-3 min-w-0">
                      <div class="flex items-start justify-between gap-2 flex-wrap">
                        <div>
                          <p class="text-sm font-semibold text-gray-800 capitalize">{{ item.status || 'Movimiento' }}</p>
                          <p class="text-xs text-gray-400 mt-0.5">{{ formatDate(item.date_time) }}</p>
                        </div>
                        <div class="text-right flex-shrink-0">
                          <p class="text-xs font-medium text-gray-600 font-mono">{{ item.user_rut || '—' }}</p>
                          <p class="text-xs text-gray-400">
                            {{ getDestinationLabel(item.destination_type) }}
                            <span v-if="item.destination_id"> #{{ item.destination_id }}</span>
                          </p>
                        </div>
                      </div>
                      <p v-if="item.notes" class="text-xs text-gray-500 mt-2 border-t border-gray-200 pt-2 leading-relaxed">
                        {{ item.notes }}
                      </p>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <div v-else class="text-center py-10 px-5">
              <svg class="h-12 w-12 mx-auto mb-3 text-gray-200" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
              </svg>
              <p class="text-sm text-gray-400">Sin movimientos registrados</p>
            </div>
          </div>
        </div>

        <!-- Col derecha (1/3): Info adicional -->
        <div class="space-y-5">

          <!-- Días hasta vencimiento -->
          <div v-if="qrInfo.supply_info?.DaysToExpire != null" class="bg-white rounded-xl shadow-sm border p-5">
            <p class="text-xs font-semibold text-gray-400 uppercase tracking-wide mb-3">Vencimiento</p>
            <div class="text-center">
              <p class="text-4xl font-bold" :class="qrInfo.supply_info.DaysToExpire <= 30 ? 'text-red-500' : 'text-green-600'">
                {{ qrInfo.supply_info.DaysToExpire }}
              </p>
              <p class="text-sm text-gray-500 mt-1">días restantes</p>
              <div class="mt-3 h-1.5 bg-gray-100 rounded-full overflow-hidden">
                <div
                  class="h-1.5 rounded-full transition-all"
                  :class="qrInfo.supply_info.DaysToExpire <= 30 ? 'bg-red-400' : 'bg-green-400'"
                  :style="`width: ${Math.min(100, Math.round(qrInfo.supply_info.DaysToExpire / 3.65))}%`"
                ></div>
              </div>
            </div>
          </div>

          <!-- Metadata del QR -->
          <div class="bg-white rounded-xl shadow-sm border p-5">
            <p class="text-xs font-semibold text-gray-400 uppercase tracking-wide mb-3">Información del QR</p>
            <dl class="space-y-2.5 text-sm">
              <div class="flex justify-between gap-2">
                <dt class="text-gray-500">Tipo</dt>
                <dd class="font-medium text-gray-800 text-right">{{ getTypeLabel(qrInfo.type) }}</dd>
              </div>
              <div class="flex justify-between gap-2">
                <dt class="text-gray-500">Escaneos totales</dt>
                <dd class="font-medium text-gray-800">{{ qrInfo.scan_statistics?.total_scans || 0 }}</dd>
              </div>
              <div v-if="qrInfo.scan_statistics?.first_scan" class="flex justify-between gap-2">
                <dt class="text-gray-500">Primer escaneo</dt>
                <dd class="font-medium text-gray-800 text-right">{{ formatDate(qrInfo.scan_statistics.first_scan) }}</dd>
              </div>
              <div v-if="qrInfo.scan_statistics?.last_scan" class="flex justify-between gap-2">
                <dt class="text-gray-500">Último escaneo</dt>
                <dd class="font-medium text-gray-800 text-right">{{ formatDate(qrInfo.scan_statistics.last_scan) }}</dd>
              </div>
              <div v-if="qrInfo.traceability?.current_location" class="flex justify-between gap-2">
                <dt class="text-gray-500">Ubicación actual</dt>
                <dd class="font-medium text-gray-800 text-right">{{ qrInfo.traceability.current_location.name }}</dd>
              </div>
            </dl>
          </div>

        </div>
      </div>

    </template>

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
import qrService from '@/services/qr/qrService'
import { useNotification } from '@/composables/useNotification'
import { useQRPdfDownload } from '@/composables/useQRPdfDownload'

const { success: showSuccess, error: showError, info: showInfo } = useNotification()

// PDF download
const { downloadQRAsPDF: downloadPDF, isGenerating: isGeneratingPDF } = useQRPdfDownload()
const downloadQRAsPDF = async () => {
  if (!qrInfo.value) return
  try {
    const ok = await downloadPDF(qrInfo.value, {
      filename: `QR_${qrInfo.value.qr_code}_${Date.now()}.pdf`,
      includeInfo: true
    })
    if (ok) showSuccess('PDF descargado exitosamente')
    else showError('Error al generar el PDF')
  } catch (err) {
    showError('Error al descargar el PDF: ' + (err.message || 'Error desconocido'))
  }
}

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
  showInfo('Datos actualizados')
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
    showSuccess('QR descargado correctamente')
  } catch (error) {
    console.error('Error downloading QR:', error)
    showError('Error al descargar el código QR')
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
  showInfo('Etiqueta enviada a impresión')
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
      showInfo('Enlace compartido correctamente')
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
    showSuccess('Enlace copiado al portapapeles')
  } catch (error) {
    showError('No se pudo copiar el enlace')
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
    showInfo('Lote sincronizado correctamente')
    await refreshData()
  } catch (error) {
    console.error('Error syncing batch:', error)
    showError('Error al sincronizar el lote')
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
  showSuccess('Reporte detallado generado y descargado')
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
  if (!status) return 'Sin estado'
  
  const statusLower = status.toLowerCase()
  
  switch (statusLower) {
    case 'disponible':
    case 'available':
      return 'Insumo Disponible'
    case 'pendiente_retiro':
      return 'Pendiente de Retiro'
    case 'en_camino_a_pabellon':
      return 'En Camino a Pabellón'
    case 'recepcionado':
    case 'received':
      return 'Recepcionado'
    case 'en_camino_a_bodega':
      return 'En Camino a Bodega'
    case 'consumido':
    case 'consumed':
      return 'Consumido'
    case 'reservado':
      return 'Reservado'
    case 'transferido':
      return 'Transferido'
    default:
      return status || 'Sin estado'
  }
}

const getStatusClass = (status) => {
  if (!status) return 'bg-gray-100 text-gray-800'
  const statusLower = status.toLowerCase()
  const statusColors = {
    'disponible': 'bg-green-100 text-green-800',
    'available': 'bg-green-100 text-green-800',
    'pendiente_retiro': 'bg-yellow-100 text-yellow-800',
    'recepcionado': 'bg-blue-100 text-blue-800',
    'received': 'bg-blue-100 text-blue-800',
    'en_camino_a_pabellon': 'bg-orange-100 text-orange-800',
    'en_camino_a_bodega': 'bg-orange-100 text-orange-800',
    'consumido': 'bg-red-100 text-red-800',
    'consumed': 'bg-red-100 text-red-800',
    'reservado': 'bg-purple-100 text-purple-800',
    'transferido': 'bg-purple-100 text-purple-800'
  }
  return statusColors[statusLower] || 'bg-gray-100 text-gray-800'
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

/* Chip de dato clave en el hero */
.detail-chip {
  @apply inline-flex items-center gap-1 px-2.5 py-1 bg-gray-100 rounded-full text-sm text-gray-700;
}
.detail-chip-label {
  @apply text-xs text-gray-400 mr-0.5;
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