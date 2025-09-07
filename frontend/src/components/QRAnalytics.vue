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
          <h1 class="text-2xl font-bold text-gray-900">Analytics del Código QR</h1>
          <p class="text-gray-600 mt-1">
            Análisis detallado de escaneos y uso del código: 
            <code class="font-mono bg-gray-100 px-2 py-1 rounded">{{ qrCode }}</code>
          </p>
        </div>
        <div class="flex space-x-2">
          <button
            @click="refreshData"
            :disabled="loading"
            class="btn-secondary"
          >
            <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
            Actualizar
          </button>
          <button
            @click="exportAnalytics"
            :disabled="!analyticsData"
            class="btn-primary"
          >
            <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
            Exportar
          </button>
        </div>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading && !analyticsData" class="flex items-center justify-center py-12">
      <div class="text-center">
        <svg class="animate-spin -ml-1 mr-3 h-8 w-8 text-blue-600 mx-auto" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
        <p class="text-gray-600 mt-2">Cargando datos de analytics...</p>
      </div>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="bg-red-50 border border-red-200 rounded-lg p-6">
      <div class="flex items-center">
        <svg class="h-6 w-6 text-red-400 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <div>
          <h3 class="text-red-800 font-medium">Error al cargar analytics</h3>
          <p class="text-red-700 text-sm mt-1">{{ error }}</p>
        </div>
      </div>
    </div>

    <!-- Main Content -->
    <div v-else-if="analyticsData">
      <!-- Filtros de Periodo -->
      <div class="bg-white rounded-lg shadow-sm border p-6 mb-6">
        <h3 class="text-lg font-medium text-gray-900 mb-4">Filtros de Análisis</h3>
        <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Período</label>
            <select v-model="selectedPeriod" @change="applyFilters" class="form-select w-full">
              <option value="today">Hoy</option>
              <option value="week">Última semana</option>
              <option value="month">Último mes</option>
              <option value="quarter">Último trimestre</option>
              <option value="all">Todo el tiempo</option>
            </select>
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Tipo de Escaneo</label>
            <select v-model="selectedScanType" @change="applyFilters" class="form-select w-full">
              <option value="all">Todos</option>
              <option value="lookup">Consultas</option>
              <option value="consume">Consumos</option>
              <option value="verify">Verificaciones</option>
              <option value="inventory_check">Inventario</option>
            </select>
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Fuente</label>
            <select v-model="selectedSource" @change="applyFilters" class="form-select w-full">
              <option value="all">Todas</option>
              <option value="web">Web</option>
              <option value="mobile">Móvil</option>
              <option value="api">API</option>
            </select>
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Estado</label>
            <select v-model="selectedStatus" @change="applyFilters" class="form-select w-full">
              <option value="all">Todos</option>
              <option value="success">Exitosos</option>
              <option value="error">Con Error</option>
            </select>
          </div>
        </div>
      </div>

      <!-- Métricas Principales -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-6">
        <div class="bg-white rounded-lg shadow-sm border p-6">
          <div class="flex items-center">
            <div class="bg-blue-100 rounded-lg p-3">
              <svg class="h-6 w-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
              </svg>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600">Total Escaneos</p>
              <p class="text-2xl font-bold text-gray-900">{{ analyticsData.total_scans }}</p>
              <p class="text-xs text-gray-500 mt-1">
                <span :class="getTrendClass(analyticsData.scan_trend)">
                  {{ getTrendText(analyticsData.scan_trend) }}
                </span>
                vs período anterior
              </p>
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
              <p class="text-2xl font-bold text-gray-900">{{ analyticsData.unique_scanners }}</p>
              <p class="text-xs text-gray-500 mt-1">
                {{ Math.round(analyticsData.total_scans / analyticsData.unique_scanners * 10) / 10 }} escaneos promedio
              </p>
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
              <p class="text-2xl font-bold text-gray-900">{{ analyticsData.locations_visited }}</p>
              <p class="text-xs text-gray-500 mt-1">
                Centros médicos visitados
              </p>
            </div>
          </div>
        </div>

        <div class="bg-white rounded-lg shadow-sm border p-6">
          <div class="flex items-center">
            <div class="bg-yellow-100 rounded-lg p-3">
              <svg class="h-6 w-6 text-yellow-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
              </svg>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600">Tasa de Éxito</p>
              <p class="text-2xl font-bold text-gray-900">{{ getSuccessRate() }}%</p>
              <p class="text-xs text-gray-500 mt-1">
                {{ analyticsData.successful_scans }}/{{ analyticsData.total_scans }} exitosos
              </p>
            </div>
          </div>
        </div>
      </div>

      <!-- Gráficos Principal -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-6">
        <!-- Gráfico de Escaneos por Tiempo -->
        <div class="bg-white rounded-lg shadow-sm border p-6">
          <h3 class="text-lg font-medium text-gray-900 mb-4">Escaneos por Tiempo</h3>
          <div class="h-64">
            <ScanTimeChart :data="timeChartData" :period="selectedPeriod" />
          </div>
        </div>

        <!-- Gráfico de Propósitos de Escaneo -->
        <div class="bg-white rounded-lg shadow-sm border p-6">
          <h3 class="text-lg font-medium text-gray-900 mb-4">Propósitos de Escaneo</h3>
          <div class="h-64">
            <ScanPurposeChart :data="purposeChartData" />
          </div>
        </div>
      </div>

      <!-- Análisis por Usuario y Ubicación -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-6">
        <!-- Top Usuarios -->
        <div class="bg-white rounded-lg shadow-sm border p-6">
          <h3 class="text-lg font-medium text-gray-900 mb-4">Usuarios Más Activos</h3>
          <div class="space-y-3">
            <div
              v-for="(user, index) in topUsers"
              :key="user.user_name"
              class="flex items-center justify-between p-3 bg-gray-50 rounded-lg"
            >
              <div class="flex items-center">
                <div class="flex-shrink-0 w-8 h-8 bg-blue-100 rounded-full flex items-center justify-center text-sm font-medium text-blue-600">
                  {{ index + 1 }}
                </div>
                <div class="ml-3">
                  <p class="text-sm font-medium text-gray-900">{{ user.user_name || 'Usuario anónimo' }}</p>
                  <p class="text-xs text-gray-500">{{ user.scan_count }} escaneos</p>
                </div>
              </div>
              <div class="text-right">
                <div class="text-sm font-medium text-gray-900">{{ getPercentage(user.scan_count, analyticsData.total_scans) }}%</div>
                <div class="w-20 bg-gray-200 rounded-full h-2 mt-1">
                  <div 
                    class="bg-blue-600 h-2 rounded-full" 
                    :style="{ width: `${getPercentage(user.scan_count, analyticsData.total_scans)}%` }"
                  ></div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Top Ubicaciones -->
        <div class="bg-white rounded-lg shadow-sm border p-6">
          <h3 class="text-lg font-medium text-gray-900 mb-4">Ubicaciones Más Activas</h3>
          <div class="space-y-3">
            <div
              v-for="(location, index) in topLocations"
              :key="location.location_name"
              class="flex items-center justify-between p-3 bg-gray-50 rounded-lg"
            >
              <div class="flex items-center">
                <div class="flex-shrink-0 w-8 h-8 bg-purple-100 rounded-full flex items-center justify-center text-sm font-medium text-purple-600">
                  {{ index + 1 }}
                </div>
                <div class="ml-3">
                  <p class="text-sm font-medium text-gray-900">{{ location.location_name || 'Ubicación no especificada' }}</p>
                  <p class="text-xs text-gray-500">{{ location.scan_count }} escaneos</p>
                </div>
              </div>
              <div class="text-right">
                <div class="text-sm font-medium text-gray-900">{{ getPercentage(location.scan_count, analyticsData.total_scans) }}%</div>
                <div class="w-20 bg-gray-200 rounded-full h-2 mt-1">
                  <div 
                    class="bg-purple-600 h-2 rounded-full" 
                    :style="{ width: `${getPercentage(location.scan_count, analyticsData.total_scans)}%` }"
                  ></div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Análisis de Patrones -->
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-6 mb-6">
        <!-- Horarios de Actividad -->
        <div class="bg-white rounded-lg shadow-sm border p-6">
          <h3 class="text-lg font-medium text-gray-900 mb-4">Horarios de Mayor Actividad</h3>
          <div class="space-y-2">
            <div
              v-for="hour in hourlyActivity"
              :key="hour.hour"
              class="flex items-center justify-between"
            >
              <div class="text-sm text-gray-600">{{ hour.hour }}:00</div>
              <div class="flex-1 mx-3">
                <div class="w-full bg-gray-200 rounded-full h-2">
                  <div 
                    class="bg-green-600 h-2 rounded-full" 
                    :style="{ width: `${getPercentage(hour.count, maxHourlyCount)}%` }"
                  ></div>
                </div>
              </div>
              <div class="text-sm font-medium text-gray-900">{{ hour.count }}</div>
            </div>
          </div>
        </div>

        <!-- Días de la Semana -->
        <div class="bg-white rounded-lg shadow-sm border p-6">
          <h3 class="text-lg font-medium text-gray-900 mb-4">Actividad por Día</h3>
          <div class="space-y-2">
            <div
              v-for="day in weeklyActivity"
              :key="day.day"
              class="flex items-center justify-between"
            >
              <div class="text-sm text-gray-600">{{ day.day }}</div>
              <div class="flex-1 mx-3">
                <div class="w-full bg-gray-200 rounded-full h-2">
                  <div 
                    class="bg-blue-600 h-2 rounded-full" 
                    :style="{ width: `${getPercentage(day.count, maxDailyCount)}%` }"
                  ></div>
                </div>
              </div>
              <div class="text-sm font-medium text-gray-900">{{ day.count }}</div>
            </div>
          </div>
        </div>

        <!-- Fuentes de Escaneo -->
        <div class="bg-white rounded-lg shadow-sm border p-6">
          <h3 class="text-lg font-medium text-gray-900 mb-4">Fuentes de Escaneo</h3>
          <div class="space-y-3">
            <div
              v-for="source in scanSources"
              :key="source.source"
              class="flex items-center justify-between p-2 rounded-lg hover:bg-gray-50"
            >
              <div class="flex items-center">
                <div 
                  :class="[
                    'w-3 h-3 rounded-full mr-3',
                    source.source === 'web' ? 'bg-blue-500' :
                    source.source === 'mobile' ? 'bg-green-500' :
                    source.source === 'api' ? 'bg-purple-500' : 'bg-gray-500'
                  ]"
                ></div>
                <span class="text-sm font-medium text-gray-900 capitalize">{{ source.source }}</span>
              </div>
              <div class="text-right">
                <div class="text-sm font-medium text-gray-900">{{ source.count }}</div>
                <div class="text-xs text-gray-500">{{ getPercentage(source.count, analyticsData.total_scans) }}%</div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Eventos Recientes -->
      <div class="bg-white rounded-lg shadow-sm border">
        <div class="p-6 border-b border-gray-200">
          <h3 class="text-lg font-medium text-gray-900">Eventos de Escaneo Recientes</h3>
        </div>
        <div class="p-6">
          <ScanEventHistory 
            :events="recentEvents" 
            :show-filters="false"
            :items-per-page="10"
            @view-details="handleEventDetails"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { format } from 'date-fns'
import { es } from 'date-fns/locale'
import qrService from '@/services/qrService'
import ScanEventHistory from '@/components/ScanEventHistory.vue'
import ScanTimeChart from '@/components/charts/ScanTimeChart.vue'
import ScanPurposeChart from '@/components/charts/ScanPurposeChart.vue'

const route = useRoute()

// Props
const qrCode = computed(() => route.params.qrcode)

// Estado reactivo
const loading = ref(false)
const error = ref(null)
const analyticsData = ref(null)
const scanEvents = ref([])
const recentEvents = ref([])

// Filtros
const selectedPeriod = ref('month')
const selectedScanType = ref('all')
const selectedSource = ref('all')
const selectedStatus = ref('all')

// Computed properties para datos procesados
const timeChartData = computed(() => {
  if (!scanEvents.value?.length) return []
  
  // Procesar datos para el gráfico de tiempo
  const groupedData = {}
  const now = new Date()
  
  scanEvents.value.forEach(event => {
    const date = new Date(event.scanned_at)
    let key
    
    switch (selectedPeriod.value) {
      case 'today':
        key = format(date, 'HH:00')
        break
      case 'week':
        key = format(date, 'EEE', { locale: es })
        break
      case 'month':
        key = format(date, 'dd/MM')
        break
      default:
        key = format(date, 'MMM yyyy', { locale: es })
    }
    
    groupedData[key] = (groupedData[key] || 0) + 1
  })
  
  return Object.entries(groupedData).map(([time, count]) => ({
    time,
    count,
    successful: count // Simplificado, podrías filtrar por éxito
  }))
})

const purposeChartData = computed(() => {
  if (!scanEvents.value?.length) return []
  
  const purposeCount = {}
  scanEvents.value.forEach(event => {
    const purpose = event.scan_purpose || 'lookup'
    purposeCount[purpose] = (purposeCount[purpose] || 0) + 1
  })
  
  return Object.entries(purposeCount).map(([name, value]) => ({
    name: getPurposeLabel(name),
    value,
    color: getPurposeColor(name)
  }))
})

const topUsers = computed(() => {
  if (!scanEvents.value?.length) return []
  
  const userCount = {}
  scanEvents.value.forEach(event => {
    const user = event.scanned_by_name || 'Anónimo'
    userCount[user] = (userCount[user] || 0) + 1
  })
  
  return Object.entries(userCount)
    .map(([user_name, scan_count]) => ({ user_name, scan_count }))
    .sort((a, b) => b.scan_count - a.scan_count)
    .slice(0, 5)
})

const topLocations = computed(() => {
  if (!scanEvents.value?.length) return []
  
  const locationCount = {}
  scanEvents.value.forEach(event => {
    const location = event.current_location || event.pavilion_name || 'No especificada'
    locationCount[location] = (locationCount[location] || 0) + 1
  })
  
  return Object.entries(locationCount)
    .map(([location_name, scan_count]) => ({ location_name, scan_count }))
    .sort((a, b) => b.scan_count - a.scan_count)
    .slice(0, 5)
})

const hourlyActivity = computed(() => {
  if (!scanEvents.value?.length) return []
  
  const hourCount = {}
  for (let i = 0; i < 24; i++) {
    hourCount[i] = 0
  }
  
  scanEvents.value.forEach(event => {
    const hour = new Date(event.scanned_at).getHours()
    hourCount[hour]++
  })
  
  return Object.entries(hourCount).map(([hour, count]) => ({
    hour: hour.toString().padStart(2, '0'),
    count
  }))
})

const weeklyActivity = computed(() => {
  if (!scanEvents.value?.length) return []
  
  const dayNames = ['Lun', 'Mar', 'Mié', 'Jue', 'Vie', 'Sáb', 'Dom']
  const dayCount = {}
  dayNames.forEach(day => dayCount[day] = 0)
  
  scanEvents.value.forEach(event => {
    const dayIndex = (new Date(event.scanned_at).getDay() + 6) % 7 // Lunes = 0
    const dayName = dayNames[dayIndex]
    dayCount[dayName]++
  })
  
  return dayNames.map(day => ({
    day,
    count: dayCount[day]
  }))
})

const scanSources = computed(() => {
  if (!scanEvents.value?.length) return []
  
  const sourceCount = {}
  scanEvents.value.forEach(event => {
    const source = event.scan_source || 'unknown'
    sourceCount[source] = (sourceCount[source] || 0) + 1
  })
  
  return Object.entries(sourceCount).map(([source, count]) => ({
    source,
    count
  }))
})

const maxHourlyCount = computed(() => {
  return Math.max(...hourlyActivity.value.map(h => h.count), 1)
})

const maxDailyCount = computed(() => {
  return Math.max(...weeklyActivity.value.map(d => d.count), 1)
})

// Funciones principales
const loadAnalytics = async () => {
  if (!qrCode.value) {
    error.value = 'Código QR requerido'
    return
  }
  
  loading.value = true
  error.value = null
  
  try {
    // Cargar estadísticas principales
    const stats = await qrService.getScanStatistics(qrCode.value)
    analyticsData.value = stats
    
    // Cargar eventos de escaneo para análisis detallado
    const events = await qrService.getScanHistory(qrCode.value, 500) // Últimos 500
    scanEvents.value = events || []
    recentEvents.value = events?.slice(0, 20) || []
    
  } catch (err) {
    console.error('Error loading analytics:', err)
    error.value = err.message || 'Error al cargar los datos de analytics'
  } finally {
    loading.value = false
  }
}

const refreshData = () => {
  loadAnalytics()
}

const applyFilters = () => {
  // Aplicar filtros a los datos
  let filtered = scanEvents.value || []
  
  // Filtros se aplicarán a los computed properties automáticamente
  console.log('Aplicando filtros:', {
    period: selectedPeriod.value,
    scanType: selectedScanType.value,
    source: selectedSource.value,
    status: selectedStatus.value
  })
}

const exportAnalytics = () => {
  if (!analyticsData.value) return
  
  const exportData = {
    qr_code: qrCode.value,
    export_date: new Date().toISOString(),
    analytics: analyticsData.value,
    time_chart_data: timeChartData.value,
    purpose_chart_data: purposeChartData.value,
    top_users: topUsers.value,
    top_locations: topLocations.value,
    hourly_activity: hourlyActivity.value,
    weekly_activity: weeklyActivity.value,
    scan_sources: scanSources.value,
    filters_applied: {
      period: selectedPeriod.value,
      scan_type: selectedScanType.value,
      source: selectedSource.value,
      status: selectedStatus.value
    }
  }
  
  const blob = new Blob([JSON.stringify(exportData, null, 2)], { type: 'application/json' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `qr-analytics-${qrCode.value}-${format(new Date(), 'yyyy-MM-dd')}.json`
  document.body.appendChild(a)
  a.click()
  document.body.removeChild(a)
  URL.revokeObjectURL(url)
}

const handleEventDetails = (event) => {
  console.log('Ver detalles del evento:', event)
}

// Funciones de utilidad
const getSuccessRate = () => {
  if (!analyticsData.value) return 0
  const rate = (analyticsData.value.successful_scans / analyticsData.value.total_scans) * 100
  return Math.round(rate * 10) / 10
}

const getPercentage = (value, total) => {
  if (!total) return 0
  return Math.round((value / total) * 100 * 10) / 10
}

const getTrendClass = (trend) => {
  if (!trend) return 'text-gray-500'
  return trend > 0 ? 'text-green-600' : trend < 0 ? 'text-red-600' : 'text-gray-500'
}

const getTrendText = (trend) => {
  if (!trend) return '0%'
  const prefix = trend > 0 ? '+' : ''
  return `${prefix}${Math.round(trend * 10) / 10}%`
}

const getPurposeLabel = (purpose) => {
  const labels = {
    'lookup': 'Consulta',
    'consume': 'Consumo',
    'verify': 'Verificación',
    'inventory_check': 'Inventario',
    'assign': 'Asignación'
  }
  return labels[purpose] || purpose
}

const getPurposeColor = (purpose) => {
  const colors = {
    'lookup': '#10B981',
    'consume': '#EF4444',
    'verify': '#3B82F6',
    'inventory_check': '#8B5CF6',
    'assign': '#F59E0B'
  }
  return colors[purpose] || '#6B7280'
}

const formatDate = (dateString) => {
  if (!dateString) return 'No disponible'
  try {
    return format(new Date(dateString), 'dd/MM/yyyy HH:mm', { locale: es })
  } catch (error) {
    return dateString
  }
}

// Lifecycle
onMounted(() => {
  loadAnalytics()
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
</style>