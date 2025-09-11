<template>
  <div class="space-y-6">
    <!-- Header con título -->
    <div class="bg-gradient-to-r from-purple-600 to-purple-700 rounded-xl p-6 text-white">
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-3xl font-bold">Estadísticas</h1>
          <p class="text-purple-100 mt-1">Panel de análisis y métricas del sistema</p>
        </div>
        <div class="text-right">
          <div class="bg-purple-500 bg-opacity-30 rounded-lg p-3">
            <svg class="h-8 w-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
            </svg>
          </div>
        </div>
      </div>
    </div>

    <!-- Métricas principales -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
      <!-- Total de Insumos -->
      <div class="bg-white rounded-lg shadow-sm border p-6 hover:shadow-md transition-shadow">
        <div class="flex items-center">
          <div class="p-3 bg-blue-100 rounded-lg">
            <svg class="h-8 w-8 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
            </svg>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600">Total Insumos</p>
            <p class="text-2xl font-bold text-gray-900">1,247</p>
            <p class="text-xs text-green-600">+12% vs mes anterior</p>
          </div>
        </div>
      </div>

      <!-- Stock Bajo -->
      <div class="bg-white rounded-lg shadow-sm border p-6 hover:shadow-md transition-shadow">
        <div class="flex items-center">
          <div class="p-3 bg-red-100 rounded-lg">
            <svg class="h-8 w-8 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.664-.833-2.464 0L4.35 16.5c-.77.833.192 2.5 1.732 2.5z" />
            </svg>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600">Stock Bajo</p>
            <p class="text-2xl font-bold text-gray-900">23</p>
            <p class="text-xs text-red-600">Requiere atención</p>
          </div>
        </div>
      </div>

      <!-- Solicitudes Pendientes -->
      <div class="bg-white rounded-lg shadow-sm border p-6 hover:shadow-md transition-shadow">
        <div class="flex items-center">
          <div class="p-3 bg-yellow-100 rounded-lg">
            <svg class="h-8 w-8 text-yellow-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
            </svg>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600">Solicitudes Pendientes</p>
            <p class="text-2xl font-bold text-gray-900">8</p>
            <p class="text-xs text-yellow-600">En proceso</p>
          </div>
        </div>
      </div>

      <!-- Eficiencia QR -->
      <div class="bg-white rounded-lg shadow-sm border p-6 hover:shadow-md transition-shadow">
        <div class="flex items-center">
          <div class="p-3 bg-green-100 rounded-lg">
            <svg class="h-8 w-8 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600">Eficiencia QR</p>
            <p class="text-2xl font-bold text-gray-900">94.2%</p>
            <p class="text-xs text-green-600">Óptimo rendimiento</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Gráficos principales -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Gráfico de inventario por mes -->
      <div class="bg-white rounded-lg shadow-sm border p-6">
        <div class="flex items-center justify-between mb-6">
          <h3 class="text-lg font-semibold text-gray-900">Movimientos de Inventario</h3>
          <div class="flex space-x-2">
            <button class="px-3 py-1 text-sm bg-purple-100 text-purple-700 rounded-md">7 días</button>
            <button class="px-3 py-1 text-sm bg-gray-100 text-gray-700 rounded-md">30 días</button>
            <button class="px-3 py-1 text-sm bg-gray-100 text-gray-700 rounded-md">90 días</button>
          </div>
        </div>
        
        <!-- Simulación de gráfico de barras -->
        <div class="space-y-3">
          <div class="flex items-center justify-between">
            <span class="text-sm text-gray-600">Entradas</span>
            <div class="flex-1 mx-4 bg-gray-200 rounded-full h-2">
              <div class="bg-green-500 h-2 rounded-full" style="width: 75%;"></div>
            </div>
            <span class="text-sm font-medium text-gray-900">156</span>
          </div>
          
          <div class="flex items-center justify-between">
            <span class="text-sm text-gray-600">Salidas</span>
            <div class="flex-1 mx-4 bg-gray-200 rounded-full h-2">
              <div class="bg-blue-500 h-2 rounded-full" style="width: 60%;"></div>
            </div>
            <span class="text-sm font-medium text-gray-900">124</span>
          </div>
          
          <div class="flex items-center justify-between">
            <span class="text-sm text-gray-600">Transferencias</span>
            <div class="flex-1 mx-4 bg-gray-200 rounded-full h-2">
              <div class="bg-purple-500 h-2 rounded-full" style="width: 45%;"></div>
            </div>
            <span class="text-sm font-medium text-gray-900">93</span>
          </div>
          
          <div class="flex items-center justify-between">
            <span class="text-sm text-gray-600">Consumos</span>
            <div class="flex-1 mx-4 bg-gray-200 rounded-full h-2">
              <div class="bg-orange-500 h-2 rounded-full" style="width: 80%;"></div>
            </div>
            <span class="text-sm font-medium text-gray-900">167</span>
          </div>
        </div>
      </div>

      <!-- Top 5 Insumos más utilizados -->
      <div class="bg-white rounded-lg shadow-sm border p-6">
        <h3 class="text-lg font-semibold text-gray-900 mb-6">Insumos Más Utilizados</h3>
        
        <div class="space-y-4">
          <div class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
            <div class="flex items-center space-x-3">
              <div class="w-8 h-8 bg-blue-100 rounded-full flex items-center justify-center">
                <span class="text-sm font-bold text-blue-600">1</span>
              </div>
              <div>
                <p class="font-medium text-gray-900">Guantes de Nitrilo</p>
                <p class="text-sm text-gray-500">Protección médica</p>
              </div>
            </div>
            <div class="text-right">
              <p class="font-bold text-gray-900">234</p>
              <p class="text-xs text-gray-500">usos</p>
            </div>
          </div>
          
          <div class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
            <div class="flex items-center space-x-3">
              <div class="w-8 h-8 bg-green-100 rounded-full flex items-center justify-center">
                <span class="text-sm font-bold text-green-600">2</span>
              </div>
              <div>
                <p class="font-medium text-gray-900">Mascarillas Quirúrgicas</p>
                <p class="text-sm text-gray-500">Protección respiratoria</p>
              </div>
            </div>
            <div class="text-right">
              <p class="font-bold text-gray-900">189</p>
              <p class="text-xs text-gray-500">usos</p>
            </div>
          </div>
          
          <div class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
            <div class="flex items-center space-x-3">
              <div class="w-8 h-8 bg-yellow-100 rounded-full flex items-center justify-center">
                <span class="text-sm font-bold text-yellow-600">3</span>
              </div>
              <div>
                <p class="font-medium text-gray-900">Jeringas Desechables</p>
                <p class="text-sm text-gray-500">Administración de medicamentos</p>
              </div>
            </div>
            <div class="text-right">
              <p class="font-bold text-gray-900">156</p>
              <p class="text-xs text-gray-500">usos</p>
            </div>
          </div>
          
          <div class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
            <div class="flex items-center space-x-3">
              <div class="w-8 h-8 bg-purple-100 rounded-full flex items-center justify-center">
                <span class="text-sm font-bold text-purple-600">4</span>
              </div>
              <div>
                <p class="font-medium text-gray-900">Gasas Estériles</p>
                <p class="text-sm text-gray-500">Cuidado de heridas</p>
              </div>
            </div>
            <div class="text-right">
              <p class="font-bold text-gray-900">142</p>
              <p class="text-xs text-gray-500">usos</p>
            </div>
          </div>
          
          <div class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
            <div class="flex items-center space-x-3">
              <div class="w-8 h-8 bg-red-100 rounded-full flex items-center justify-center">
                <span class="text-sm font-bold text-red-600">5</span>
              </div>
              <div>
                <p class="font-medium text-gray-900">Alcohol en Gel</p>
                <p class="text-sm text-gray-500">Desinfección</p>
              </div>
            </div>
            <div class="text-right">
              <p class="font-bold text-gray-900">128</p>
              <p class="text-xs text-gray-500">usos</p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Estadísticas por ubicación -->
    <div class="bg-white rounded-lg shadow-sm border p-6">
      <h3 class="text-lg font-semibold text-gray-900 mb-6">Distribución por Ubicación</h3>
      
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div class="text-center">
          <div class="relative w-32 h-32 mx-auto mb-4">
            <!-- Simulación de gráfico circular -->
            <div class="w-32 h-32 rounded-full bg-gradient-to-r from-blue-500 to-purple-500 flex items-center justify-center">
              <div class="w-20 h-20 bg-white rounded-full flex items-center justify-center">
                <span class="text-lg font-bold text-gray-900">67%</span>
              </div>
            </div>
          </div>
          <h4 class="font-semibold text-gray-900">Centro Principal</h4>
          <p class="text-sm text-gray-500">834 insumos</p>
        </div>
        
        <div class="text-center">
          <div class="relative w-32 h-32 mx-auto mb-4">
            <div class="w-32 h-32 rounded-full bg-gradient-to-r from-green-500 to-blue-500 flex items-center justify-center">
              <div class="w-20 h-20 bg-white rounded-full flex items-center justify-center">
                <span class="text-lg font-bold text-gray-900">23%</span>
              </div>
            </div>
          </div>
          <h4 class="font-semibold text-gray-900">Pabellón Norte</h4>
          <p class="text-sm text-gray-500">287 insumos</p>
        </div>
        
        <div class="text-center">
          <div class="relative w-32 h-32 mx-auto mb-4">
            <div class="w-32 h-32 rounded-full bg-gradient-to-r from-yellow-500 to-orange-500 flex items-center justify-center">
              <div class="w-20 h-20 bg-white rounded-full flex items-center justify-center">
                <span class="text-lg font-bold text-gray-900">10%</span>
              </div>
            </div>
          </div>
          <h4 class="font-semibold text-gray-900">Pabellón Sur</h4>
          <p class="text-sm text-gray-500">126 insumos</p>
        </div>
      </div>
    </div>

    <!-- Próximas funcionalidades -->
    <div class="bg-gradient-to-r from-gray-50 to-gray-100 rounded-lg p-6 border border-dashed border-gray-300">
      <div class="text-center">
        <div class="mx-auto w-12 h-12 bg-gray-200 rounded-full flex items-center justify-center mb-4">
          <svg class="h-6 w-6 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
              d="M13 10V3L4 14h7v7l9-11h-7z" />
          </svg>
        </div>
        <h3 class="text-lg font-medium text-gray-900 mb-2">Más funcionalidades próximamente</h3>
        <p class="text-gray-600 mb-4">
          Estamos trabajando en reportes avanzados, alertas inteligentes, y análisis predictivo para optimizar tu gestión de inventario.
        </p>
        <div class="flex flex-wrap justify-center gap-2">
          <span class="px-3 py-1 bg-blue-100 text-blue-700 rounded-full text-sm">Reportes Automáticos</span>
          <span class="px-3 py-1 bg-green-100 text-green-700 rounded-full text-sm">Alertas Inteligentes</span>
          <span class="px-3 py-1 bg-purple-100 text-purple-700 rounded-full text-sm">Análisis Predictivo</span>
          <span class="px-3 py-1 bg-orange-100 text-orange-700 rounded-full text-sm">Dashboard Personalizable</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

// Datos simulados para la demostración
const stats = ref({
  totalSupplies: 1247,
  lowStock: 23,
  pendingRequests: 8,
  qrEfficiency: 94.2
})

// Simular carga de datos
onMounted(() => {
  // Aquí se cargarían los datos reales desde la API
  console.log('Cargando estadísticas...')
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