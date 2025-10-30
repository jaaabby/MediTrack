<template>
  <div class="space-y-6">
    <!-- Header con título -->
    <div class="bg-gradient-to-r from-purple-600 to-purple-700 rounded-xl p-6 text-white shadow">
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-3xl font-bold">Estadísticas</h1>
          <p class="text-purple-100 mt-1">Panel de análisis y métricas del sistema</p>
        </div>
        <div class="text-right">
          <div class="bg-purple-500 bg-opacity-30 rounded-lg p-3 shadow-inner">
            <svg class="h-8 w-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
            </svg>
          </div>
        </div>
      </div>
    </div>

    <!-- Loading / Error -->
    <div v-if="loading" class="bg-white rounded-lg shadow border p-6">Cargando estadísticas...</div>
    <div v-else-if="error" class="bg-red-50 border border-red-200 text-red-800 rounded-lg p-4">{{ error }}</div>

    <template v-else>
      <!-- Métricas principales -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <!-- Total de Insumos -->
        <div class="bg-white rounded-xl shadow-sm border p-6 hover:shadow-md transition-shadow">
          <div class="flex items-center">
            <div class="p-3 bg-gradient-to-br from-blue-100 to-blue-200 rounded-xl shadow-sm">
              <svg class="h-8 w-8 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                  d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
              </svg>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600">Total Insumos</p>
              <p class="text-2xl font-bold text-gray-900">{{ mainMetrics.totalSupplies.toLocaleString('es-CL') }}</p>
              <p class="text-xs text-green-600">Inventario total (bodega + pabellones)</p>
            </div>
          </div>
        </div>

        <!-- Stock Bajo -->
        <div class="bg-white rounded-xl shadow-sm border p-6 hover:shadow-md transition-shadow">
          <div class="flex items-center">
            <div class="p-3 bg-gradient-to-br from-red-100 to-red-200 rounded-xl shadow-sm">
              <svg class="h-8 w-8 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                  d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.664-.833-2.464 0L4.35 16.5c-.77.833.192 2.5 1.732 2.5z" />
              </svg>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600">Stock Bajo</p>
              <p class="text-2xl font-bold text-gray-900">{{ mainMetrics.lowStock.toLocaleString('es-CL') }}</p>
              <p class="text-xs text-red-600">Batches bajo 20% del stock</p>
            </div>
          </div>
        </div>

        <!-- Transferencias Pendientes -->
        <div class="bg-white rounded-xl shadow-sm border p-6 hover:shadow-md transition-shadow">
          <div class="flex items-center">
            <div class="p-3 bg-gradient-to-br from-yellow-100 to-yellow-200 rounded-xl shadow-sm">
              <svg class="h-8 w-8 text-yellow-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                  d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
              </svg>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600">Transferencias Pendientes</p>
              <p class="text-2xl font-bold text-gray-900">{{ mainMetrics.pendingTransfers.toLocaleString('es-CL') }}</p>
              <p class="text-xs text-yellow-600">En tránsito o por confirmar</p>
            </div>
          </div>
        </div>

        <!-- Próximos a vencer -->
        <div class="bg-white rounded-xl shadow-sm border p-6 hover:shadow-md transition-shadow">
          <div class="flex items-center">
            <div class="p-3 bg-gradient-to-br from-green-100 to-green-200 rounded-xl shadow-sm">
              <svg class="h-8 w-8 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                  d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600">Próximos a Vencer</p>
              <p class="text-2xl font-bold text-gray-900">{{ Number(summary.near_expiration || 0).toLocaleString('es-CL') }}</p>
              <p class="text-xs text-green-600">Dentro de 90 días</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Reportes de stock y tendencia -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div class="bg-white rounded-xl shadow-sm border p-6">
          <div class="flex items-center justify-between mb-4">
            <div class="flex items-center gap-2">
              <div class="w-8 h-8 rounded-full bg-purple-100 flex items-center justify-center">
                <svg class="w-4 h-4 text-purple-600" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 3v18h18"/></svg>
              </div>
              <h3 class="text-lg font-semibold text-gray-900">Tendencia de transferencias</h3>
            </div>
            <div class="flex gap-2">
              <button class="px-2 py-1 text-xs rounded border" :class="transferRange==='7d' ? 'bg-purple-100 border-purple-300 text-purple-700' : 'bg-white text-gray-700'" @click="transferRange='7d'">7 días</button>
              <button class="px-2 py-1 text-xs rounded border" :class="transferRange==='30d' ? 'bg-purple-100 border-purple-300 text-purple-700' : 'bg-white text-gray-700'" @click="transferRange='30d'">30 días</button>
              <button class="px-2 py-1 text-xs rounded border" :class="transferRange==='6m' ? 'bg-purple-100 border-purple-300 text-purple-700' : 'bg-white text-gray-700'" @click="transferRange='6m'">6 meses</button>
              <button class="px-2 py-1 text-xs rounded border" :class="transferRange==='1y' ? 'bg-purple-100 border-purple-300 text-purple-700' : 'bg-white text-gray-700'" @click="transferRange='1y'">1 año</button>
              <button class="px-2 py-1 text-xs rounded border" :class="transferRange==='all' ? 'bg-purple-100 border-purple-300 text-purple-700' : 'bg-white text-gray-700'" @click="transferRange='all'">Todo</button>
            </div>
          </div>
          <div class="h-24 flex items-center justify-center" v-if="!hasTrendData">
            <div class="text-center">
              <div class="text-sm text-gray-500">Sin transferencias en el periodo seleccionado</div>
              <div class="text-xs text-gray-400">Realiza una transferencia para ver la tendencia aquí</div>
            </div>
          </div>
          <div class="h-24" v-else>
            <svg :width="220" :height="48">
              <defs>
                <linearGradient id="spark" x1="0" x2="0" y1="0" y2="1">
                  <stop offset="0%" stop-color="#7C3AED" stop-opacity="0.4"/>
                  <stop offset="100%" stop-color="#7C3AED" stop-opacity="0"/>
                </linearGradient>
              </defs>
              <path :d="sparklinePath" stroke="#7C3AED" fill="none" stroke-width="2" />
            </svg>
            <div class="text-xs text-gray-500 mt-2">Total: {{ totalTransfers }}</div>
          </div>
        </div>
        <div class="bg-white rounded-xl shadow-sm border p-6">
          <div class="flex items-center justify-between mb-4">
            <div class="flex items-center gap-2">
              <div class="w-8 h-8 rounded-full bg-indigo-100 flex items-center justify-center">
                <svg class="w-4 h-4 text-indigo-600" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 12h18"/><path d="M3 6h18"/><path d="M3 18h18"/></svg>
              </div>
              <h3 class="text-lg font-semibold text-gray-900">Stock y movimientos</h3>
            </div>
          </div>
          <div class="space-y-3">
            <div class="flex items-center justify-between">
              <span class="text-sm text-gray-600">Stock actual</span>
              <div class="flex-1 mx-4 bg-gray-200 rounded-full h-2">
                <div class="bg-indigo-600 h-2 rounded-full" :style="{ width: Math.min(100, mainMetrics.totalSupplies ? 90 : 0) + '%' }"></div>
              </div>
              <span class="text-sm font-medium text-gray-900">{{ mainMetrics.totalSupplies }}</span>
            </div>
            <div class="flex items-center justify-between">
              <span class="text-sm text-gray-600">Entradas (aprox)</span>
              <div class="flex-1 mx-4 bg-gray-200 rounded-full h-2">
                <div class="bg-green-500 h-2 rounded-full" :style="{ width: Math.min(100, movementBars.entradas ? 70 : 0) + '%' }"></div>
              </div>
              <span class="text-sm font-medium text-gray-900">{{ movementBars.entradas }}</span>
            </div>
            <div class="flex items-center justify-between">
              <span class="text-sm text-gray-600">Salidas (transfer)</span>
              <div class="flex-1 mx-4 bg-gray-200 rounded-full h-2">
                <div class="bg-blue-500 h-2 rounded-full" :style="{ width: Math.min(100, movementBars.salidas ? 60 : 0) + '%' }"></div>
              </div>
              <span class="text-sm font-medium text-gray-900">{{ movementBars.salidas }}</span>
            </div>
            <div class="flex items-center justify-between">
              <span class="text-sm text-gray-600">Consumos</span>
              <div class="flex-1 mx-4 bg-gray-200 rounded-full h-2">
                <div class="bg-orange-500 h-2 rounded-full" :style="{ width: Math.min(100, movementBars.consumos ? 80 : 0) + '%' }"></div>
              </div>
              <span class="text-sm font-medium text-gray-900">{{ movementBars.consumos }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Cirugías: duración promedio y actividad -->
      <div class="bg-white rounded-xl shadow-sm border p-6">
        <div class="flex items-center justify-between mb-4">
          <div class="flex items-center gap-2">
            <div class="w-8 h-8 rounded-full bg-purple-100 flex items-center justify-center">
              <svg class="w-4 h-4 text-purple-600" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 2v7"/><path d="M5 10h14"/></svg>
            </div>
            <h3 class="text-lg font-semibold text-gray-900">Estadísticas por tipo de cirugía</h3>
            <span class="ml-2 text-xs bg-purple-100 text-purple-700 px-2 py-0.5 rounded-full">{{ surgeriesWithTotals.length }} tipos</span>
          </div>
          <div class="flex items-center gap-4">
            <span class="text-sm text-gray-500 hidden sm:inline">Duración promedio: {{ avgSurgeryDuration }} h</span>
            <input v-model="surgerySearch" type="text" placeholder="Buscar cirugía..." class="text-sm border rounded px-2 py-1 focus:outline-none focus:ring-2 focus:ring-purple-200" />
          </div>
        </div>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <!-- Lista con scroll -->
          <div class="space-y-3 max-h-80 overflow-y-auto pr-2 divide-y divide-gray-100">
            <div
              v-for="item in surgeriesWithTotals"
              :key="item.surgery_id"
              class="flex items-center justify-between py-2"
              @mouseenter="(e) => showTooltip(e, item.surgery_name)"
              @mousemove="moveTooltip"
              @mouseleave="hideTooltip"
            >
              <div class="flex items-center min-w-0 mr-3 gap-3">
                <!-- Icono médico: estetoscopio -->
                <div class="w-8 h-8 rounded-full bg-purple-100 flex items-center justify-center flex-shrink-0 shadow-inner">
                  <svg class="w-4 h-4 text-purple-600" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M6 3v6a4 4 0 0 0 8 0V3"/>
                    <path d="M4 3h4"/>
                    <path d="M12 3h4"/>
                    <circle cx="20" cy="10" r="2"/>
                    <path d="M20 12v2a6 6 0 0 1-6 6h-1a3 3 0 0 1-3-3v-2"/>
                  </svg>
                </div>
                <div class="min-w-0">
                  <div class="truncate font-medium text-gray-900">{{ item.surgery_name }}</div>
                  <div class="text-xs text-gray-500">Duración: {{ item.duration }} h · Transferidos: {{ item.total_transferred }}</div>
                </div>
              </div>
              <div class="w-40 bg-gray-200 rounded-full h-2 flex-shrink-0">
                <div class="bg-purple-500 h-2 rounded-full" :style="{ width: Math.min(100, Number(item.total_transferred || 0)) + '%' }"></div>
              </div>
            </div>
            <div v-if="surgeriesWithTotals.length === 0" class="text-sm text-gray-500">Sin datos de cirugías para mostrar</div>
          </div>
          <!-- Top Insumos con scroll -->
          <div class="space-y-3 max-h-80 overflow-y-auto pr-2">
            <div class="text-sm font-semibold text-gray-900 mb-2 sticky top-0 bg-white">Top Insumos más utilizados</div>
            <div v-for="s in topSupplies" :key="s.code" class="p-3 bg-gray-50 rounded-lg shadow-sm">
              <div class="flex items-center justify-between">
                <div class="truncate mr-3 flex items-center gap-2">
                  <!-- Icono médico: frasco/píldora -->
                  <div class="w-7 h-7 rounded bg-blue-100 flex items-center justify-center flex-shrink-0 shadow-inner">
                    <svg class="w-4 h-4 text-blue-600" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                      <rect x="6" y="3" width="12" height="6" rx="1"/>
                      <path d="M8 9v9a3 3 0 0 0 3 3h2a3 3 0 0 0 3-3V9"/>
                      <path d="M9 6h6"/>
                    </svg>
                  </div>
                  <div>
                    <div class="font-medium text-gray-900 truncate">{{ s.name }}</div>
                    <div class="text-xs text-gray-500">Código: {{ s.code }}</div>
                  </div>
                </div>
                <div class="flex items-center gap-2">
                  <span class="text-xs text-gray-500">movidos</span>
                  <span class="text-sm font-semibold text-gray-900">{{ s.total }}</span>
                </div>
              </div>
              <div class="mt-2 h-2 w-full bg-gray-200 rounded-full">
                <div class="h-2 bg-blue-500 rounded-full" :style="{ width: Math.min(100, Math.round((Number(s.total||0) / maxTopSupply) * 100)) + '%' }"></div>
              </div>
            </div>
            <div v-if="!topSupplies.length" class="text-sm text-gray-500">Sin datos disponibles</div>
          </div>
        </div>
      </div>

      <!-- Distribución por Ubicación: Bodegas -->
      <div class="bg-white rounded-xl shadow-sm border p-6">
        <div class="flex items-center justify-between mb-6">
          <div class="flex items-center gap-2">
            <div class="w-8 h-8 rounded-full bg-blue-100 flex items-center justify-center">
              <svg class="w-4 h-4 text-blue-600" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 6l9-4 9 4v12l-9 4-9-4z"/></svg>
            </div>
            <h3 class="text-lg font-semibold text-gray-900">Bodegas</h3>
            <span class="ml-2 text-xs bg-blue-100 text-blue-700 px-2 py-0.5 rounded-full">{{ storeList.length }}</span>
          </div>
        </div>
        <div class="max-h-80 overflow-y-auto pr-2">
          <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
            <template v-if="storeList.length">
              <div class="text-center" v-for="st in storeList" :key="st.id">
                <div class="relative w-32 h-32 mx-auto mb-4">
                  <div class="w-32 h-32 rounded-full bg-gradient-to-r from-blue-500 to-purple-500 flex items-center justify-center shadow">
                    <div class="w-20 h-20 bg-white rounded-full flex items-center justify-center">
                      <span class="text-lg font-bold text-gray-900">{{ st.current_in_store }}</span>
                    </div>
                  </div>
                </div>
                <h4 class="font-semibold text-gray-900">{{ st.name }}</h4>
                <p class="text-sm text-gray-500">Stock actual</p>
              </div>
            </template>
            <template v-else>
              <div class="col-span-3 text-center text-sm text-gray-500">No hay datos de bodegas para mostrar</div>
            </template>
          </div>
        </div>
      </div>

      <!-- Distribución por Ubicación: Pabellones -->
      <div class="bg-white rounded-xl shadow-sm border p-6">
        <div class="flex items-center justify-between mb-6">
          <div class="flex items-center gap-2">
            <div class="w-8 h-8 rounded-full bg-green-100 flex items-center justify-center">
              <svg class="w-4 h-4 text-green-600" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 10H7"/><path d="M21 6H3"/><path d="M21 14H3"/><path d="M21 18H7"/></svg>
            </div>
            <h3 class="text-lg font-semibold text-gray-900">Pabellones</h3>
            <span class="ml-2 text-xs bg-green-100 text-green-700 px-2 py-0.5 rounded-full">{{ pavilionList.length }}</span>
          </div>
        </div>
        <div class="max-h-80 overflow-y-auto pr-2">
          <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
            <template v-if="pavilionList.length">
              <div class="text-center" v-for="pv in pavilionList" :key="pv.id">
                <div class="relative w-32 h-32 mx-auto mb-4">
                  <div class="w-32 h-32 rounded-full bg-gradient-to-r from-green-500 to-blue-500 flex items-center justify-center shadow">
                    <div class="w-20 h-20 bg-white rounded-full flex items-center justify-center">
                      <span class="text-lg font-bold text-gray-900">{{ pv.current_available }}</span>
                    </div>
                  </div>
                </div>
                <h4 class="font-semibold text-gray-900">{{ pv.name }}</h4>
                <p class="text-sm text-gray-500">Disponibles</p>
              </div>
            </template>
            <template v-else>
              <div class="col-span-3 text-center text-sm text-gray-500">No hay datos de pabellones para mostrar</div>
            </template>
          </div>
        </div>
      </div>

      <!-- Alertas críticas -->
      <div class="bg-gradient-to-r from-gray-50 to-gray-100 rounded-xl p-6 border border-dashed border-gray-300 shadow-sm">
        <div class="text-center">
          <div class="mx-auto w-12 h-12 bg-gray-200 rounded-full flex items-center justify-center mb-4 shadow-inner">
            <svg class="h-6 w-6 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                d="M13 10V3L4 14h7v7l9-11h-7z" />
            </svg>
          </div>
          <h3 class="text-lg font-medium text-gray-900 mb-2">Alertas críticas</h3>
          <p class="text-gray-600 mb-4">
            {{ lowStockList.length ? 'Insumos con stock crítico en bodega' : 'Sin alertas críticas por ahora' }}
          </p>
          <div class="flex flex-wrap justify-center gap-2" v-if="lowStockList.length">
            <span
              v-for="item in lowStockList"
              :key="item.batch_id"
              class="px-3 py-1 bg-red-100 text-red-700 rounded-full text-sm shadow-sm"
            >{{ item.supply_name }} (#{{ item.supply_code }})</span>
          </div>
        </div>
      </div>
    </template>
    <!-- Tooltip global -->
    <teleport to="body">
      <div
        v-if="tooltipVisible"
        :style="{ position: 'fixed', top: tooltipY + 'px', left: tooltipX + 'px', zIndex: 9999 }"
        class="px-2 py-1 text-xs bg-gray-900 text-white rounded shadow-lg max-w-xs"
      >
        {{ tooltipText }}
      </div>
    </teleport>
  </div>
</template>

<script setup>
import { ref, onMounted, computed, watch } from 'vue'
import inventoryService from '@/services/inventoryService'
import supplyRequestService from '@/services/supplyRequestService'
import surgeryService from '@/services/surgeryService'

const loading = ref(false)
const error = ref('')

// Tooltip global para nombres completos
const tooltipVisible = ref(false)
const tooltipText = ref('')
const tooltipX = ref(0)
const tooltipY = ref(0)

function showTooltip(e, text) {
  tooltipText.value = text
  tooltipX.value = e.clientX + 12
  tooltipY.value = e.clientY + 12
  tooltipVisible.value = true
}
function moveTooltip(e) {
  if (!tooltipVisible.value) return
  tooltipX.value = e.clientX + 12
  tooltipY.value = e.clientY + 12
}
function hideTooltip() {
  tooltipVisible.value = false
  tooltipText.value = ''
}

// Resumen principal
const summary = ref({
  total_in_stores: 0,
  total_in_pavilions: 0,
  total_transferred: 0,
  total_consumed: 0,
  low_stock_stores: 0,
  near_expiration: 0,
  pending_transfers: 0,
})

const mainMetrics = computed(() => ({
  totalSupplies: Number(summary.value.total_in_stores || 0) + Number(summary.value.total_in_pavilions || 0),
  lowStock: Number(summary.value.low_stock_stores || 0),
  pendingTransfers: Number(summary.value.pending_transfers || 0),
}))

// Top por tipo de cirugía
const bySurgery = ref([])
const surgeries = ref([])
const avgSurgeryDuration = computed(() => {
  if (!surgeries.value.length) return 0
  const sum = surgeries.value.reduce((a, s) => a + Number(s.duration || 0), 0)
  return Math.round((sum / surgeries.value.length) * 10) / 10
})

// Mapear todas las cirugías con totales (0 por defecto)
const surgerySearch = ref('')
const surgeriesWithTotals = computed(() => {
  const totalsMap = new Map((bySurgery.value || []).map(it => [it.surgery_id, Number(it.total_transferred || 0)]))
  const all = (surgeries.value || []).map(s => ({
    surgery_id: s.id,
    surgery_name: s.name,
    duration: Number(s.duration || 0),
    total_transferred: totalsMap.get(s.id) || 0,
  }))
  if (!surgerySearch.value) return all
  const term = surgerySearch.value.toLowerCase()
  return all.filter(x => (x.surgery_name || '').toLowerCase().includes(term))
})

// Alertas críticas
const lowStockList = ref([])

// Movimientos
const movementBars = ref({
  entradas: 0,
  salidas: 0,
  transferencias: 0,
  consumos: 0,
})
const transferTrend = ref([]) // [{date, count}]
const totalTransfers = computed(() => transferTrend.value.reduce((a, x) => a + (Number(x.count)||0), 0))
const hasTrendData = computed(() => transferTrend.value.length > 0 && totalTransfers.value > 0)
const transferRange = ref('7d') // '7d' | '30d' | '6m' | '1y'

// Distribución por ubicación: pabellones y bodegas
const pavilionList = ref([])
const storeList = ref([])

// Top insumos más utilizados
const topSupplies = ref([])
const maxTopSupply = computed(() => {
  if (!topSupplies.value.length) return 1
  return Math.max(...topSupplies.value.map(s => Number(s.total || 0)), 1)
})

function formatISODate(d) {
  return new Date(d).toISOString().slice(0, 10)
}

function buildSparklinePath(points, width = 220, height = 48) {
  if (!points.length) return ''
  const maxY = Math.max(...points.map(p => p.count), 1)
  const stepX = width / Math.max(points.length - 1, 1)
  const scaleY = (val) => height - (val / maxY) * (height - 4)
  return points.map((p, i) => `${i === 0 ? 'M' : 'L'} ${i * stepX} ${scaleY(p.count)}`).join(' ')
}

const sparklinePath = computed(() => buildSparklinePath(transferTrend.value))

async function loadPavilionDistribution() {
  const pavilions = await inventoryService.getAllPavilions()
  const enriched = await Promise.all(
    (Array.isArray(pavilions) ? pavilions : []).map(async (p) => {
      try {
        const items = await inventoryService.getPavilionInventory(p.id, false)
        const current = Array.isArray(items) ? items.reduce((acc, it) => acc + Number(it.current_available || 0), 0) : 0
        const consumed = Array.isArray(items) ? items.reduce((acc, it) => acc + Number(it.total_consumed || 0), 0) : 0
        return { ...p, current_available: current, total_consumed: consumed }
      } catch {
        return { ...p, current_available: 0, total_consumed: 0 }
      }
    })
  )
  pavilionList.value = enriched
}

async function loadStoreDistribution() {
  const stores = await inventoryService.getAllStores()
  const enriched = await Promise.all(
    (Array.isArray(stores) ? stores : []).map(async (s) => {
      try {
        const items = await inventoryService.getStoreInventory({ store_id: s.id, page: 1, page_size: 1000 })
        const current = Array.isArray(items) ? items.reduce((acc, it) => acc + Number(it.current_in_store || 0), 0) : 0
        return { ...s, current_in_store: current }
      } catch {
        return { ...s, current_in_store: 0 }
      }
    })
  )
  storeList.value = enriched
}

async function loadTopSupplies() {
  const items = await inventoryService.getStoreInventory({ page: 1, page_size: 1000 })
  const map = new Map()
  ;(Array.isArray(items) ? items : []).forEach(it => {
    const key = it.supply_code || it.code
    const name = it.supply_name || it.name || `Código ${key}`
    const moved = Number(it.total_transferred_out || 0)
    if (!map.has(key)) map.set(key, { code: key, name, total: 0 })
    map.get(key).total += moved
  })
  const arr = Array.from(map.values()).sort((a, b) => b.total - a.total)
  topSupplies.value = arr
}

function getStartDateByRange(range) {
  const now = new Date()
  const start = new Date(now)
  if (range === '7d') start.setDate(now.getDate() - 7)
  else if (range === '30d') start.setDate(now.getDate() - 30)
  else if (range === '6m') start.setMonth(now.getMonth() - 6)
  else if (range === '1y') start.setFullYear(now.getFullYear() - 1)
  else if (range === 'all') start.setFullYear(1970, 0, 1)
  return { start, now }
}

async function loadTrend() {
  try {
    const { start, now } = getStartDateByRange(transferRange.value)
    const transferReport = await inventoryService.getTransferReport(formatISODate(start), formatISODate(now), 'date')
    transferTrend.value = (Array.isArray(transferReport) ? transferReport : []).map(r => ({
      date: r.transfer_date || r.date || '',
      count: Number(r.transfer_count || 0),
    }))
    movementBars.value.transferencias = totalTransfers.value
  } catch {
    transferTrend.value = []
    movementBars.value.transferencias = 0
  }
}

async function loadData() {
  loading.value = true
  error.value = ''
  try {
    const [summaryResp, bySurgResp, surgResp, lowStockResp] = await Promise.all([
      inventoryService.getInventorySummary(),
      inventoryService.getInventoryBySurgeryType(),
      surgeryService.getAllSurgeries(),
      inventoryService.getStoreInventory({ low_stock: true, page: 1, page_size: 8 }),
    ])

    summary.value = summaryResp || summary.value
    bySurgery.value = Array.isArray(bySurgResp) ? bySurgResp : []
    surgeries.value = Array.isArray(surgResp) ? surgResp : []
    lowStockList.value = Array.isArray(lowStockResp) ? lowStockResp : []

    await Promise.all([loadTrend(), loadPavilionDistribution(), loadStoreDistribution(), loadTopSupplies()])

    movementBars.value.consumos = Number(summary.value.total_consumed || 0)
    movementBars.value.entradas = Number(summary.value.total_in_stores || 0)
    movementBars.value.salidas = Number(summary.value.total_transferred || 0)
  } catch (e) {
    error.value = e?.message || 'Error cargando estadísticas'
  } finally {
    loading.value = false
  }
}

watch(transferRange, () => {
  loadTrend()
})

onMounted(() => {
  loadData()
})
</script>

<style scoped>
/* Animaciones para las métricas */
.hover\:shadow-md:hover {
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
}

/* Transiciones suaves */
.transition-shadow {
  transition: box-shadow 0.2s ease-in-out;
}

/* Efectos de gradiente para los gráficos circulares */
.bg-gradient-to-r {
  background-image: linear-gradient(to right, var(--tw-gradient-stops));
}

/* Animación sutil para los elementos */
@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.space-y-6 > * {
  animation: fadeInUp 0.6s ease-out;
}

.space-y-6 > *:nth-child(2) {
  animation-delay: 0.1s;
}

.space-y-6 > *:nth-child(3) {
  animation-delay: 0.2s;
}

.space-y-6 > *:nth-child(4) {
  animation-delay: 0.3s;
}
</style>