<template>
  <div class="space-y-4 px-3 sm:px-0 sm:space-y-6">
    <!-- Header con título -->
    <div class="bg-gradient-to-r from-brand-pink to-brand-pink rounded-lg sm:rounded-xl p-3 sm:p-6 text-gray-900 shadow" style="background: linear-gradient(135deg, #FA92B9 0%, #f57ba8 100%);">
      <div class="flex items-center justify-between gap-3">
        <div class="min-w-0 flex-1">
          <h1 class="text-xl sm:text-2xl md:text-3xl font-bold truncate text-gray-900">Estadísticas</h1>
          <p class="text-gray-900 opacity-80 mt-1 text-xs sm:text-sm md:text-base truncate">Panel de análisis y métricas del sistema</p>
        </div>
        <div class="flex-shrink-0">
          <div class="bg-white bg-opacity-50 rounded-lg p-2 sm:p-3 shadow-inner">
            <svg class="h-5 w-5 sm:h-6 sm:w-6 md:h-8 md:w-8 text-gray-900" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
            </svg>
          </div>
        </div>
      </div>
    </div>

    <!-- Loading / Error -->
    <div v-if="loading" class="bg-white rounded-lg shadow border p-4 text-sm sm:text-base">Cargando estadísticas...</div>
    <div v-else-if="error" class="bg-red-50 border border-red-200 text-red-800 rounded-lg p-3 sm:p-4 text-sm sm:text-base">{{ error }}</div>

    <template v-else>
      <!-- Métricas principales -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-3 sm:gap-4 md:gap-6">
        <!-- Total de Insumos -->
        <button @click="goToInventory" class="group bg-white rounded-lg sm:rounded-xl shadow-sm border-2 border-transparent hover:border-brand-blue-medium hover:shadow-lg transition-all duration-200 cursor-pointer text-left w-full p-3 sm:p-4 md:p-6 relative overflow-hidden">
          <div class="absolute top-2 right-2 opacity-0 group-hover:opacity-100 transition-opacity">
            <svg class="h-5 w-5 text-brand-blue-dark" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
          </div>
          <div class="flex items-start gap-3">
            <div class="p-2 sm:p-3 bg-gradient-to-br from-brand-blue-light from-opacity-30 to-brand-blue-light to-opacity-50 rounded-lg sm:rounded-xl shadow-sm flex-shrink-0 group-hover:from-brand-blue-light group-hover:to-brand-blue-medium group-hover:to-opacity-30 transition-colors">
              <svg class="h-6 w-6 sm:h-7 sm:w-7 md:h-8 md:w-8 text-brand-blue-dark" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                  d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
              </svg>
            </div>
            <div class="min-w-0 flex-1">
              <p class="text-xs sm:text-sm font-medium text-gray-600 truncate flex items-center gap-1">
                Total Insumos
                <span class="text-brand-blue-dark opacity-0 group-hover:opacity-100 transition-opacity text-xs">Ver →</span>
              </p>
              <p class="text-lg sm:text-xl md:text-2xl font-bold text-gray-900 group-hover:text-brand-blue-dark transition-colors">{{ mainMetrics.totalSupplies.toLocaleString('es-CL') }}</p>
              <p class="text-xs text-brand-green truncate">Inventario total</p>
            </div>
          </div>
        </button>

        <!-- Stock Bajo -->
        <button @click="goToInventory('lowStock')" class="group bg-white rounded-lg sm:rounded-xl shadow-sm border-2 border-transparent hover:border-brand-pink hover:shadow-lg transition-all duration-200 cursor-pointer text-left w-full p-3 sm:p-4 md:p-6 relative overflow-hidden">
          <div class="absolute top-2 right-2 opacity-0 group-hover:opacity-100 transition-opacity">
            <svg class="h-5 w-5 text-brand-pink" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
          </div>
          <div class="flex items-start gap-3">
            <div class="p-2 sm:p-3 bg-gradient-to-br from-brand-pink from-opacity-20 to-brand-pink from-opacity-40 rounded-lg sm:rounded-xl shadow-sm flex-shrink-0 group-hover:from-brand-pink group-hover:to-brand-pink group-hover:to-opacity-50 transition-colors">
              <svg class="h-6 w-6 sm:h-7 sm:w-7 md:h-8 md:w-8 text-gray-900" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                  d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.664-.833-2.464 0L4.35 16.5c-.77.833.192 2.5 1.732 2.5z" />
              </svg>
            </div>
            <div class="min-w-0 flex-1">
              <p class="text-xs sm:text-sm font-medium text-gray-600 truncate flex items-center gap-1">
                Stock Bajo
                <span class="text-brand-pink opacity-0 group-hover:opacity-100 transition-opacity text-xs">Ver →</span>
              </p>
              <p class="text-lg sm:text-xl md:text-2xl font-bold text-gray-900 group-hover:text-brand-pink transition-colors">{{ mainMetrics.lowStock.toLocaleString('es-CL') }}</p>
              <p class="text-xs text-brand-pink truncate">Bajo 20% del stock</p>
            </div>
          </div>
        </button>

        <!-- Transferencias Pendientes -->
        <button @click="goToTransfers" class="group bg-white rounded-lg sm:rounded-xl shadow-sm border-2 border-transparent hover:border-brand-blue-light hover:shadow-lg transition-all duration-200 cursor-pointer text-left w-full p-3 sm:p-4 md:p-6 relative overflow-hidden">
          <div class="absolute top-2 right-2 opacity-0 group-hover:opacity-100 transition-opacity">
            <svg class="h-5 w-5 text-brand-blue-dark" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
          </div>
          <div class="flex items-start gap-3">
            <div class="p-2 sm:p-3 bg-gradient-to-br from-brand-blue-light from-opacity-20 to-brand-blue-light from-opacity-40 rounded-lg sm:rounded-xl shadow-sm flex-shrink-0 group-hover:from-brand-blue-light group-hover:to-brand-blue-light group-hover:to-opacity-50 transition-colors">
              <svg class="h-6 w-6 sm:h-7 sm:w-7 md:h-8 md:w-8 text-brand-blue-dark" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                  d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
              </svg>
            </div>
            <div class="min-w-0 flex-1">
              <p class="text-xs sm:text-sm font-medium text-gray-600 truncate flex items-center gap-1">
                Transferencias
                <span class="text-brand-blue-light opacity-0 group-hover:opacity-100 transition-opacity text-xs">Ver →</span>
              </p>
              <p class="text-lg sm:text-xl md:text-2xl font-bold text-gray-900 group-hover:text-brand-blue-light transition-colors">{{ mainMetrics.pendingTransfers.toLocaleString('es-CL') }}</p>
              <p class="text-xs text-brand-blue-light truncate">En tránsito</p>
            </div>
          </div>
        </button>

        <!-- Próximos a vencer -->
        <button @click="goToInventory('expiring')" class="group bg-white rounded-lg sm:rounded-xl shadow-sm border-2 border-transparent hover:border-brand-green hover:shadow-lg transition-all duration-200 cursor-pointer text-left w-full p-3 sm:p-4 md:p-6 relative overflow-hidden">
          <div class="absolute top-2 right-2 opacity-0 group-hover:opacity-100 transition-opacity">
            <svg class="h-5 w-5 text-brand-green" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
          </div>
          <div class="flex items-start gap-3">
            <div class="p-2 sm:p-3 bg-gradient-to-br from-brand-green from-opacity-30 to-brand-green from-opacity-50 rounded-lg sm:rounded-xl shadow-sm flex-shrink-0 group-hover:from-brand-green group-hover:to-brand-green group-hover:to-opacity-60 transition-colors">
              <svg class="h-6 w-6 sm:h-7 sm:w-7 md:h-8 md:w-8 text-gray-900" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                  d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <div class="min-w-0 flex-1">
              <p class="text-xs sm:text-sm font-medium text-gray-600 truncate flex items-center gap-1">
                Por Vencer
                <span class="text-brand-green opacity-0 group-hover:opacity-100 transition-opacity text-xs">Ver →</span>
              </p>
              <p class="text-lg sm:text-xl md:text-2xl font-bold text-gray-900 group-hover:text-brand-green transition-colors">{{ Number(summary.near_expiration || 0).toLocaleString('es-CL') }}</p>
              <p class="text-xs text-brand-green truncate">Dentro 90 días</p>
            </div>
          </div>
        </button>
      </div>

      <!-- Métricas de Configuración Médica -->
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-3 sm:gap-4 md:gap-6">
        <!-- Total de Especialidades -->
        <button @click="goToMedicalSpecialties" class="group bg-white rounded-lg sm:rounded-xl shadow-sm border-2 border-transparent hover:border-brand-blue-medium hover:shadow-lg transition-all duration-200 cursor-pointer text-left w-full p-3 sm:p-4 md:p-6 relative overflow-hidden">
          <div class="absolute top-2 right-2 opacity-0 group-hover:opacity-100 transition-opacity">
            <svg class="h-5 w-5 text-brand-blue-medium" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
          </div>
          <div class="flex items-start gap-3">
            <div class="p-2 sm:p-3 bg-gradient-to-br from-brand-blue-medium from-opacity-20 to-brand-blue-medium from-opacity-40 rounded-lg sm:rounded-xl shadow-sm flex-shrink-0 group-hover:from-brand-blue-medium group-hover:to-brand-blue-medium group-hover:to-opacity-50 transition-colors">
              <svg class="h-6 w-6 sm:h-7 sm:w-7 md:h-8 md:w-8 text-gray-900" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                  d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
              </svg>
            </div>
            <div class="min-w-0 flex-1">
              <p class="text-xs sm:text-sm font-medium text-gray-600 truncate flex items-center gap-1">
                Especialidades
                <span class="text-brand-blue-medium opacity-0 group-hover:opacity-100 transition-opacity text-xs">Ver →</span>
              </p>
              <p class="text-lg sm:text-xl md:text-2xl font-bold text-gray-900 group-hover:text-brand-blue-medium transition-colors">{{ medicalMetrics.totalSpecialties.toLocaleString('es-CL') }}</p>
              <p class="text-xs text-brand-blue-medium truncate">Médicas registradas</p>
            </div>
          </div>
        </button>

        <!-- Total de Cirugías -->
        <button @click="goToSurgeries" class="group bg-white rounded-lg sm:rounded-xl shadow-sm border-2 border-transparent hover:border-brand-blue-dark hover:shadow-lg transition-all duration-200 cursor-pointer text-left w-full p-3 sm:p-4 md:p-6 relative overflow-hidden">
          <div class="absolute top-2 right-2 opacity-0 group-hover:opacity-100 transition-opacity">
            <svg class="h-5 w-5 text-brand-blue-dark" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
          </div>
          <div class="flex items-start gap-3">
            <div class="p-2 sm:p-3 bg-gradient-to-br from-brand-blue-dark from-opacity-20 to-brand-blue-dark from-opacity-40 rounded-lg sm:rounded-xl shadow-sm flex-shrink-0 group-hover:from-brand-blue-dark group-hover:to-brand-blue-dark group-hover:to-opacity-50 transition-colors">
              <svg class="h-6 w-6 sm:h-7 sm:w-7 md:h-8 md:w-8 text-gray-900" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                  d="M12 2v7M5 10h14" />
              </svg>
            </div>
            <div class="min-w-0 flex-1">
              <p class="text-xs sm:text-sm font-medium text-gray-600 truncate flex items-center gap-1">
                Tipos de Cirugía
                <span class="text-brand-blue-dark opacity-0 group-hover:opacity-100 transition-opacity text-xs">Ver →</span>
              </p>
              <p class="text-lg sm:text-xl md:text-2xl font-bold text-gray-900 group-hover:text-brand-blue-dark transition-colors">{{ medicalMetrics.totalSurgeries.toLocaleString('es-CL') }}</p>
              <p class="text-xs text-brand-blue-dark truncate">Procedimientos</p>
            </div>
          </div>
        </button>

        <!-- Total de Doctores -->
        <button @click="goToDoctors" class="group bg-white rounded-lg sm:rounded-xl shadow-sm border-2 border-transparent hover:border-brand-blue-light hover:shadow-lg transition-all duration-200 cursor-pointer text-left w-full p-3 sm:p-4 md:p-6 relative overflow-hidden">
          <div class="absolute top-2 right-2 opacity-0 group-hover:opacity-100 transition-opacity">
            <svg class="h-5 w-5 text-brand-blue-dark" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
          </div>
          <div class="flex items-start gap-3">
            <div class="p-2 sm:p-3 bg-gradient-to-br from-brand-blue-light from-opacity-20 to-brand-blue-light from-opacity-40 rounded-lg sm:rounded-xl shadow-sm flex-shrink-0 group-hover:from-brand-blue-light group-hover:to-brand-blue-light group-hover:to-opacity-50 transition-colors">
              <svg class="h-6 w-6 sm:h-7 sm:w-7 md:h-8 md:w-8 text-brand-blue-dark" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                  d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
              </svg>
            </div>
            <div class="min-w-0 flex-1">
              <p class="text-xs sm:text-sm font-medium text-gray-600 truncate flex items-center gap-1">
                Doctores
                <span class="text-brand-blue-light opacity-0 group-hover:opacity-100 transition-opacity text-xs">Ver →</span>
              </p>
              <p class="text-lg sm:text-xl md:text-2xl font-bold text-gray-900 group-hover:text-brand-blue-light transition-colors">{{ medicalMetrics.totalDoctors.toLocaleString('es-CL') }}</p>
              <p class="text-xs text-brand-blue-light truncate">Registrados</p>
            </div>
          </div>
        </button>

        <!-- Total de Insumos Típicos -->
        <button @click="goToTypicalSupplies" class="group bg-white rounded-lg sm:rounded-xl shadow-sm border-2 border-transparent hover:border-brand-pink hover:shadow-lg transition-all duration-200 cursor-pointer text-left w-full p-3 sm:p-4 md:p-6 relative overflow-hidden">
          <div class="absolute top-2 right-2 opacity-0 group-hover:opacity-100 transition-opacity">
            <svg class="h-5 w-5 text-brand-pink" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
          </div>
          <div class="flex items-start gap-3">
            <div class="p-2 sm:p-3 bg-gradient-to-br from-brand-pink from-opacity-20 to-brand-pink from-opacity-40 rounded-lg sm:rounded-xl shadow-sm flex-shrink-0 group-hover:from-brand-pink group-hover:to-brand-pink group-hover:to-opacity-50 transition-colors">
              <svg class="h-6 w-6 sm:h-7 sm:w-7 md:h-8 md:w-8 text-gray-900" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                  d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
              </svg>
            </div>
            <div class="min-w-0 flex-1">
              <p class="text-xs sm:text-sm font-medium text-gray-600 truncate flex items-center gap-1">
                Insumos Típicos
                <span class="text-brand-pink opacity-0 group-hover:opacity-100 transition-opacity text-xs">Ver →</span>
              </p>
              <p class="text-lg sm:text-xl md:text-2xl font-bold text-gray-900 group-hover:text-brand-pink transition-colors">{{ medicalMetrics.totalTypicalSupplies.toLocaleString('es-CL') }}</p>
              <p class="text-xs text-brand-pink truncate">Asociaciones</p>
            </div>
          </div>
        </button>
      </div>

      <!-- Reportes de stock y tendencia -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-3 sm:gap-4 md:gap-6">
        <!-- Tendencia de transferencias -->
        <div class="bg-white rounded-lg sm:rounded-xl shadow-sm border p-3 sm:p-4 md:p-6">
          <div class="flex flex-col gap-2 mb-3 sm:mb-4">
            <div class="flex items-center gap-2">
              <div class="w-7 h-7 sm:w-8 sm:h-8 rounded-full bg-brand-pink bg-opacity-20 flex items-center justify-center flex-shrink-0">
                <svg class="w-3 h-3 sm:w-4 sm:h-4 text-gray-900" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 3v18h18"/></svg>
              </div>
              <h3 class="text-sm sm:text-base md:text-lg font-semibold text-gray-900 truncate">Tendencia de transferencias</h3>
            </div>
            <div class="flex gap-1.5 flex-wrap">
              <button class="px-2 py-1 text-xs rounded border flex-1 sm:flex-initial whitespace-nowrap" :class="transferRange==='7d' ? 'bg-brand-pink bg-opacity-20 border-brand-pink text-brand-pink' : 'bg-white text-gray-700'" @click="transferRange='7d'">7d</button>
              <button class="px-2 py-1 text-xs rounded border flex-1 sm:flex-initial whitespace-nowrap" :class="transferRange==='30d' ? 'bg-brand-pink bg-opacity-20 border-brand-pink text-brand-pink' : 'bg-white text-gray-700'" @click="transferRange='30d'">30d</button>
              <button class="px-2 py-1 text-xs rounded border flex-1 sm:flex-initial whitespace-nowrap" :class="transferRange==='6m' ? 'bg-brand-pink bg-opacity-20 border-brand-pink text-brand-pink' : 'bg-white text-gray-700'" @click="transferRange='6m'">6m</button>
              <button class="px-2 py-1 text-xs rounded border flex-1 sm:flex-initial whitespace-nowrap" :class="transferRange==='1y' ? 'bg-brand-pink bg-opacity-20 border-brand-pink text-brand-pink' : 'bg-white text-gray-700'" @click="transferRange='1y'">1a</button>
              <button class="px-2 py-1 text-xs rounded border flex-1 sm:flex-initial whitespace-nowrap" :class="transferRange==='all' ? 'bg-brand-pink bg-opacity-20 border-brand-pink text-brand-pink' : 'bg-white text-gray-700'" @click="transferRange='all'">Todo</button>
            </div>
          </div>
          <div class="h-20 sm:h-24 flex items-center justify-center" v-if="!hasTrendData">
            <div class="text-center">
              <div class="text-xs sm:text-sm text-gray-500">Sin transferencias</div>
              <div class="text-xs text-gray-400">en el periodo seleccionado</div>
            </div>
          </div>
          <div class="h-20 sm:h-24" v-else>
            <svg width="100%" height="100%" viewBox="0 0 220 48" preserveAspectRatio="none" style="overflow: visible;">
              <defs>
                <linearGradient id="spark" x1="0" x2="0" y1="0" y2="1">
                  <stop offset="0%" stop-color="#FA92B9" stop-opacity="0.4"/>
                  <stop offset="100%" stop-color="#FA92B9" stop-opacity="0"/>
                </linearGradient>
              </defs>
              <path 
                v-if="sparklinePath && sparklinePath.length > 0"
                :d="sparklinePath" 
                stroke="#FA92B9" 
                fill="none" 
                stroke-width="2" 
                stroke-linecap="round"
                stroke-linejoin="round"
              />
            </svg>
            <div class="text-xs text-gray-500 mt-2">Total: {{ totalTransfers }}</div>
          </div>
        </div>

        <!-- Stock y movimientos -->
        <div class="bg-white rounded-lg sm:rounded-xl shadow-sm border p-3 sm:p-4 md:p-6">
          <div class="flex items-center gap-2 mb-3 sm:mb-4">
            <div class="w-7 h-7 sm:w-8 sm:h-8 rounded-full bg-brand-blue-dark bg-opacity-20 flex items-center justify-center flex-shrink-0">
              <svg class="w-3 h-3 sm:w-4 sm:h-4 text-brand-blue-dark" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 12h18"/><path d="M3 6h18"/><path d="M3 18h18"/></svg>
            </div>
            <h3 class="text-sm sm:text-base md:text-lg font-semibold text-gray-900 truncate">Stock y movimientos</h3>
          </div>
          <div class="space-y-2.5 sm:space-y-3">
            <div class="flex items-center gap-2 sm:gap-3">
              <span class="text-xs sm:text-sm text-gray-600 w-16 sm:w-20 flex-shrink-0">Stock</span>
              <div class="flex-1 bg-gray-200 rounded-full h-1.5 sm:h-2">
                <div class="bg-brand-blue-dark h-1.5 sm:h-2 rounded-full" :style="{ width: Math.min(100, mainMetrics.totalSupplies ? 90 : 0) + '%' }"></div>
              </div>
              <span class="text-xs sm:text-sm font-medium text-gray-900 w-10 sm:w-12 text-right">{{ mainMetrics.totalSupplies }}</span>
            </div>
            <div class="flex items-center gap-2 sm:gap-3">
              <span class="text-xs sm:text-sm text-gray-600 w-16 sm:w-20 flex-shrink-0">Entradas</span>
              <div class="flex-1 bg-gray-200 rounded-full h-1.5 sm:h-2">
                <div class="bg-brand-green h-1.5 sm:h-2 rounded-full" :style="{ width: Math.min(100, movementBars.entradas ? 70 : 0) + '%' }"></div>
              </div>
              <span class="text-xs sm:text-sm font-medium text-gray-900 w-10 sm:w-12 text-right">{{ movementBars.entradas }}</span>
            </div>
            <div class="flex items-center gap-2 sm:gap-3">
              <span class="text-xs sm:text-sm text-gray-600 w-16 sm:w-20 flex-shrink-0">Salidas</span>
              <div class="flex-1 bg-gray-200 rounded-full h-1.5 sm:h-2">
                <div class="bg-brand-blue-medium h-1.5 sm:h-2 rounded-full" :style="{ width: Math.min(100, movementBars.salidas ? 60 : 0) + '%' }"></div>
              </div>
              <span class="text-xs sm:text-sm font-medium text-gray-900 w-10 sm:w-12 text-right">{{ movementBars.salidas }}</span>
            </div>
            <div class="flex items-center gap-2 sm:gap-3">
              <span class="text-xs sm:text-sm text-gray-600 w-16 sm:w-20 flex-shrink-0">Consumos</span>
              <div class="flex-1 bg-gray-200 rounded-full h-1.5 sm:h-2">
                <div class="bg-brand-blue-light h-1.5 sm:h-2 rounded-full" :style="{ width: Math.min(100, movementBars.consumos ? 80 : 0) + '%' }"></div>
              </div>
              <span class="text-xs sm:text-sm font-medium text-gray-900 w-10 sm:w-12 text-right">{{ movementBars.consumos }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Cirugías -->
      <div class="bg-white rounded-lg sm:rounded-xl shadow-sm border p-3 sm:p-4 md:p-6">
        <div class="flex flex-col gap-2 sm:gap-3 mb-3 sm:mb-4">
          <div class="flex items-center justify-between gap-2">
            <div class="flex items-center gap-2 min-w-0 flex-1">
              <div class="w-7 h-7 sm:w-8 sm:h-8 rounded-full bg-brand-pink bg-opacity-20 flex items-center justify-center flex-shrink-0">
                <svg class="w-3 h-3 sm:w-4 sm:h-4 text-gray-900" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M12 2v7"/><path d="M5 10h14"/></svg>
              </div>
              <h3 class="text-sm sm:text-base md:text-lg font-semibold text-gray-900 truncate">Cirugías</h3>
              <span class="text-xs bg-brand-pink bg-opacity-20 text-brand-pink px-1.5 sm:px-2 py-0.5 rounded-full flex-shrink-0">{{ surgeriesWithTotals.length }}</span>
            </div>
            <div class="text-xs sm:text-sm text-gray-500 flex-shrink-0">Prom: {{ avgSurgeryDuration }}h</div>
          </div>
          <input v-model="surgerySearch" type="text" placeholder="Buscar..." class="text-xs sm:text-sm border rounded px-2 sm:px-3 py-1.5 sm:py-2 focus:outline-none focus:ring-2 focus:ring-brand-pink focus:ring-opacity-30 w-full" />
        </div>
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-3 sm:gap-4 md:gap-6">
          <!-- Lista cirugías -->
          <div class="space-y-2 max-h-52 sm:max-h-60 md:max-h-80 overflow-y-auto pr-1 sm:pr-2">
            <div
              v-for="item in surgeriesWithTotals"
              :key="item.surgery_id"
              class="flex items-center justify-between py-2 border-b border-gray-100 last:border-0"
              @mouseenter="(e) => showTooltip(e, item.surgery_name)"
              @mousemove="moveTooltip"
              @mouseleave="hideTooltip"
            >
              <div class="flex items-center min-w-0 mr-2 sm:mr-3 gap-2 sm:gap-3 flex-1">
                <div class="w-6 h-6 sm:w-8 sm:h-8 rounded-full bg-brand-pink bg-opacity-20 flex items-center justify-center flex-shrink-0 shadow-inner">
                  <svg class="w-3 h-3 sm:w-4 sm:h-4 text-gray-900" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M6 3v6a4 4 0 0 0 8 0V3"/>
                    <path d="M4 3h4"/>
                    <path d="M12 3h4"/>
                    <circle cx="20" cy="10" r="2"/>
                    <path d="M20 12v2a6 6 0 0 1-6 6h-1a3 3 0 0 1-3-3v-2"/>
                  </svg>
                </div>
                <div class="min-w-0 flex-1">
                  <div class="truncate font-medium text-gray-900 text-xs sm:text-sm md:text-base">{{ item.surgery_name }}</div>
                  <div class="text-xs text-gray-500 truncate">{{ item.duration }}h · {{ item.total_transferred }} transf.</div>
                </div>
              </div>
              <div class="w-16 sm:w-20 md:w-28 bg-gray-200 rounded-full h-1.5 sm:h-2 flex-shrink-0">
                <div class="bg-brand-pink h-1.5 sm:h-2 rounded-full" :style="{ width: Math.min(100, Number(item.total_transferred || 0)) + '%' }"></div>
              </div>
            </div>
            <div v-if="surgeriesWithTotals.length === 0" class="text-xs sm:text-sm text-gray-500 text-center py-4">Sin datos de cirugías</div>
          </div>
          
          <!-- Top Insumos -->
          <div class="space-y-2 max-h-52 sm:max-h-60 md:max-h-80 overflow-y-auto pr-1 sm:pr-2">
            <div class="text-xs sm:text-sm font-semibold text-gray-900 mb-2 sticky top-0 bg-white py-1">Top Insumos más utilizados</div>
            <div v-for="s in topSupplies" :key="s.code" class="p-2 sm:p-3 bg-gray-50 rounded-lg shadow-sm">
              <div class="flex items-center justify-between gap-2 mb-1.5 sm:mb-2">
                <div class="flex items-center gap-2 min-w-0 flex-1">
                  <div class="w-6 h-6 sm:w-7 sm:h-7 rounded bg-brand-blue-dark bg-opacity-20 flex items-center justify-center flex-shrink-0 shadow-inner">
                    <svg class="w-3 h-3 sm:w-4 sm:h-4 text-brand-blue-dark" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                      <rect x="6" y="3" width="12" height="6" rx="1"/>
                      <path d="M8 9v9a3 3 0 0 0 3 3h2a3 3 0 0 0 3-3V9"/>
                      <path d="M9 6h6"/>
                    </svg>
                  </div>
                  <div class="min-w-0 flex-1">
                    <div class="font-medium text-gray-900 truncate text-xs sm:text-sm md:text-base">{{ s.name }}</div>
                    <div class="text-xs text-gray-500 truncate">{{ s.code }}</div>
                  </div>
                </div>
                <div class="flex items-center gap-1 flex-shrink-0">
                  <span class="text-xs text-gray-500 hidden sm:inline">movidos</span>
                  <span class="text-xs sm:text-sm font-semibold text-gray-900">{{ s.total }}</span>
                </div>
              </div>
              <div class="h-1.5 sm:h-2 w-full bg-gray-200 rounded-full">
                <div class="h-1.5 sm:h-2 bg-brand-blue-dark rounded-full" :style="{ width: Math.min(100, Math.round((Number(s.total||0) / maxTopSupply) * 100)) + '%' }"></div>
              </div>
            </div>
            <div v-if="!topSupplies.length" class="text-xs sm:text-sm text-gray-500 text-center py-4">Sin datos disponibles</div>
          </div>
        </div>
      </div>

      <!-- Bodegas -->
      <div class="bg-white rounded-lg sm:rounded-xl shadow-sm border p-3 sm:p-4 md:p-6">
        <div class="flex items-center gap-2 mb-3 sm:mb-4 md:mb-6">
          <div class="w-7 h-7 sm:w-8 sm:h-8 rounded-full bg-brand-blue-dark bg-opacity-20 flex items-center justify-center flex-shrink-0">
            <svg class="w-3 h-3 sm:w-4 sm:h-4 text-brand-blue-dark" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 6l9-4 9 4v12l-9 4-9-4z"/></svg>
          </div>
          <h3 class="text-sm sm:text-base md:text-lg font-semibold text-gray-900">Bodegas</h3>
          <span class="text-xs bg-brand-blue-dark bg-opacity-20 text-brand-blue-dark px-1.5 sm:px-2 py-0.5 rounded-full">{{ storeList.length }}</span>
        </div>
        <div class="max-h-52 sm:max-h-60 md:max-h-80 overflow-y-auto pr-1 sm:pr-2">
          <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-3 sm:gap-4 md:gap-6">
            <template v-if="storeList.length">
              <div class="text-center" v-for="st in storeList" :key="st.id">
                <div class="relative w-20 h-20 sm:w-24 sm:h-24 md:w-32 md:h-32 mx-auto mb-2 sm:mb-3 md:mb-4">
                  <div class="w-full h-full rounded-full bg-gradient-to-r from-brand-blue-dark to-brand-blue-medium flex items-center justify-center shadow">
                    <div class="w-14 h-14 sm:w-16 sm:h-16 md:w-20 md:h-20 bg-white rounded-full flex items-center justify-center">
                      <span class="text-sm sm:text-base md:text-lg font-bold text-gray-900">{{ st.current_in_store }}</span>
                    </div>
                  </div>
                </div>
                <h4 class="font-semibold text-gray-900 text-xs sm:text-sm md:text-base truncate px-1">{{ st.name }}</h4>
                <p class="text-xs text-gray-500 truncate">Stock actual</p>
              </div>
            </template>
            <template v-else>
              <div class="col-span-2 sm:col-span-3 md:col-span-4 text-center text-xs sm:text-sm text-gray-500 py-4">No hay datos de bodegas</div>
            </template>
          </div>
        </div>
      </div>

      <!-- Pabellones -->
      <div class="bg-white rounded-lg sm:rounded-xl shadow-sm border p-3 sm:p-4 md:p-6">
        <div class="flex items-center gap-2 mb-3 sm:mb-4 md:mb-6">
          <div class="w-7 h-7 sm:w-8 sm:h-8 rounded-full bg-brand-green bg-opacity-30 flex items-center justify-center flex-shrink-0">
            <svg class="w-3 h-3 sm:w-4 sm:h-4 text-gray-900" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M21 10H7"/><path d="M21 6H3"/><path d="M21 14H3"/><path d="M21 18H7"/></svg>
          </div>
          <h3 class="text-sm sm:text-base md:text-lg font-semibold text-gray-900">Pabellones</h3>
          <span class="text-xs bg-brand-green bg-opacity-30 text-gray-900 px-1.5 sm:px-2 py-0.5 rounded-full">{{ pavilionList.length }}</span>
        </div>
        <div class="max-h-52 sm:max-h-60 md:max-h-80 overflow-y-auto pr-1 sm:pr-2">
          <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 gap-3 sm:gap-4 md:gap-6">
            <template v-if="pavilionList.length">
              <div class="text-center" v-for="pv in pavilionList" :key="pv.id">
                <div class="relative w-20 h-20 sm:w-24 sm:h-24 md:w-32 md:h-32 mx-auto mb-2 sm:mb-3 md:mb-4">
                  <div class="w-full h-full rounded-full bg-gradient-to-r from-brand-green to-brand-blue-light flex items-center justify-center shadow">
                    <div class="w-14 h-14 sm:w-16 sm:h-16 md:w-20 md:h-20 bg-white rounded-full flex items-center justify-center">
                      <span class="text-sm sm:text-base md:text-lg font-bold text-gray-900">{{ pv.current_available }}</span>
                    </div>
                  </div>
                </div>
                <h4 class="font-semibold text-gray-900 text-xs sm:text-sm md:text-base truncate px-1">{{ pv.name }}</h4>
                <p class="text-xs text-gray-500 truncate">Disponibles</p>
              </div>
            </template>
            <template v-else>
              <div class="col-span-2 sm:col-span-3 md:col-span-4 text-center text-xs sm:text-sm text-gray-500 py-4">No hay datos de pabellones</div>
            </template>
          </div>
        </div>
      </div>

      <!-- Alertas críticas -->
      <div class="bg-gradient-to-r from-gray-50 to-gray-100 rounded-lg sm:rounded-xl p-3 sm:p-4 md:p-6 border border-dashed border-gray-300 shadow-sm">
        <div class="text-center">
          <div class="mx-auto w-9 h-9 sm:w-10 sm:h-10 md:w-12 md:h-12 bg-gray-200 rounded-full flex items-center justify-center mb-2 sm:mb-3 md:mb-4 shadow-inner">
            <svg class="h-4 w-4 sm:h-5 sm:w-5 md:h-6 md:w-6 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                d="M13 10V3L4 14h7v7l9-11h-7z" />
            </svg>
          </div>
          <h3 class="text-sm sm:text-base md:text-lg font-medium text-gray-900 mb-1 sm:mb-2">Alertas críticas</h3>
          <p class="text-gray-600 mb-2 sm:mb-3 md:mb-4 text-xs sm:text-sm md:text-base">
            {{ lowStockList.length ? 'Insumos con stock crítico' : 'Sin alertas críticas' }}
          </p>
          <div class="flex flex-wrap justify-center gap-1.5 sm:gap-2" v-if="lowStockList.length">
            <span
              v-for="item in lowStockList"
              :key="item.batch_id"
              class="px-2 py-1 bg-brand-pink bg-opacity-20 text-brand-pink rounded-full text-xs shadow-sm"
            >{{ item.supply_name }} (#{{ item.supply_code }})</span>
          </div>
        </div>
      </div>
    </template>
    
    <!-- Tooltip -->
    <teleport to="body">
      <div
        v-if="tooltipVisible"
        :style="{ position: 'fixed', top: tooltipY + 'px', left: tooltipX + 'px', zIndex: 9999 }"
        class="px-2 py-1 text-xs bg-gray-900 text-white rounded shadow-lg max-w-xs pointer-events-none"
      >
        {{ tooltipText }}
      </div>
    </teleport>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import inventoryService from '@/services/inventory/inventoryService'
import supplyRequestService from '@/services/requests/supplyRequestService'
import surgeryService from '@/services/management/surgeryService'
import medicalSpecialtyService from '@/services/config/medicalSpecialtyService'
import doctorInfoService from '@/services/config/doctorInfoService'
import surgeryTypicalSupplyService from '@/services/management/surgeryTypicalSupplyService'

const router = useRouter()

const loading = ref(false)
const error = ref('')
const autoRefreshInterval = ref(null)
const autoRefreshEnabled = ref(true)
const refreshIntervalSeconds = 60 // Actualizar cada 60 segundos

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

const bySurgery = ref([])
const surgeries = ref([])
const avgSurgeryDuration = computed(() => {
  if (!surgeries.value.length) return 0
  const sum = surgeries.value.reduce((a, s) => a + Number(s.duration || 0), 0)
  return Math.round((sum / surgeries.value.length) * 10) / 10
})

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

const lowStockList = ref([])

const movementBars = ref({
  entradas: 0,
  salidas: 0,
  transferencias: 0,
  consumos: 0,
})
const transferTrend = ref([])
const totalTransfers = computed(() => transferTrend.value.reduce((a, x) => a + (Number(x.count)||0), 0))
const hasTrendData = computed(() => {
  // Tener datos si hay puntos Y si el total es mayor a 0
  return transferTrend.value.length > 0 && totalTransfers.value > 0
})
const transferRange = ref('7d')

const pavilionList = ref([])
const storeList = ref([])

const topSupplies = ref([])
const maxTopSupply = computed(() => {
  if (!topSupplies.value.length) return 1
  return Math.max(...topSupplies.value.map(s => Number(s.total || 0)), 1)
})

// Métricas de configuración médica
const medicalMetrics = ref({
  totalSpecialties: 0,
  totalSurgeries: 0,
  totalDoctors: 0,
  totalTypicalSupplies: 0,
  requestsBySpecialty: [],
  mostRequestedSurgeries: []
})

// Funciones de navegación
const goToInventory = (filter = null) => {
  if (filter === 'lowStock') {
    router.push({ path: '/inventory', query: { lowStock: 'true' } })
  } else if (filter === 'expiring') {
    router.push({ path: '/inventory', query: { expiring: 'true' } })
  } else {
    router.push('/inventory')
  }
}

const goToSupplyRequests = () => {
  router.push('/supply-requests')
}

const goToTransfers = () => {
  router.push('/transfers')
}

const goToMedicalSpecialties = () => {
  router.push('/medical-specialties')
}

const goToSurgeries = () => {
  router.push('/surgeries')
}

const goToDoctors = () => {
  router.push('/doctor-info')
}

const goToTypicalSupplies = () => {
  router.push('/surgery-typical-supplies')
}

function formatISODate(d) {
  return new Date(d).toISOString().slice(0, 10)
}

function buildSparklinePath(points, width = 220, height = 48) {
  if (!points.length) return ''
  
  // Obtener todos los counts y encontrar el máximo
  const counts = points.map(p => Number(p.count) || 0)
  const maxY = Math.max(...counts, 1) // Al menos 1 para evitar división por cero
  
  // Si solo hay un punto, duplicarlo al inicio y al final para que se vea la línea
  let processedPoints = points
  if (points.length === 1) {
    processedPoints = [
      { count: 0, date: '' },
      { ...points[0] },
      { count: 0, date: '' }
    ]
  }
  
  const pointCount = processedPoints.length
  const stepX = pointCount > 1 ? width / (pointCount - 1) : width
  const padding = 4 // Padding vertical para que no toque los bordes
  const scaleY = (val) => {
    const scaled = (val / maxY) * (height - padding * 2)
    return height - scaled - padding
  }
  
  // Construir el path
  let path = ''
  processedPoints.forEach((p, i) => {
    const x = i * stepX
    const y = scaleY(Number(p.count) || 0)
    if (i === 0) {
      path = `M ${x} ${y}`
    } else {
      path += ` L ${x} ${y}`
    }
  })
  
  return path
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
    // Sumar transferencias, consumos y devoluciones para el total de movimientos
    const transferred = Number(it.total_transferred_out || 0)
    const consumed = Number(it.total_consumed_in_store || 0)
    const returned = Number(it.total_returned_in || 0)
    const moved = transferred + consumed + returned
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
    // El backend devuelve {report: [...], ...}, necesitamos extraer el array de report
    const reportData = transferReport?.report || transferReport || []
    transferTrend.value = (Array.isArray(reportData) ? reportData : []).map(r => ({
      date: r.transfer_date || r.date || '',
      count: Number(r.transfer_count || 0),
    }))
    movementBars.value.transferencias = totalTransfers.value
  } catch (err) {
    console.error('Error cargando tendencia de transferencias:', err)
    transferTrend.value = []
    movementBars.value.transferencias = 0
  }
}

async function loadMedicalMetrics() {
  try {
    const [specialties, surgeries, doctors, typicalSuppliesCount] = await Promise.all([
      medicalSpecialtyService.getAllSpecialties().catch(() => []),
      surgeryService.getAllSurgeries().catch(() => []),
      doctorInfoService.getAllDoctors().catch(() => []),
      surgeryTypicalSupplyService.getTypicalSuppliesCount().catch(() => 0)
    ])

    medicalMetrics.value.totalSpecialties = Array.isArray(specialties) ? specialties.filter(s => s.is_active).length : 0
    medicalMetrics.value.totalSurgeries = Array.isArray(surgeries) ? surgeries.length : 0
    medicalMetrics.value.totalDoctors = Array.isArray(doctors) ? doctors.filter(d => d.is_available).length : 0
    medicalMetrics.value.totalTypicalSupplies = Number(typicalSuppliesCount) || 0

    // Cargar solicitudes por especialidad (si hay datos disponibles)
    try {
      const requests = await supplyRequestService.getAllSupplyRequests(1000, 0, '')
      if (requests.data && Array.isArray(requests.data.requests)) {
        const bySpecialty = new Map()
        requests.data.requests.forEach(req => {
          if (req.specialty_id) {
            const count = bySpecialty.get(req.specialty_id) || 0
            bySpecialty.set(req.specialty_id, count + 1)
          }
        })
        medicalMetrics.value.requestsBySpecialty = Array.from(bySpecialty.entries()).map(([specialtyId, count]) => {
          const specialty = specialties.find(s => s.id === specialtyId)
          return {
            specialty_id: specialtyId,
            specialty_name: specialty ? specialty.name : 'Sin especialidad',
            count: count
          }
        }).sort((a, b) => b.count - a.count)
      }
    } catch (err) {
      console.error('Error cargando solicitudes por especialidad:', err)
    }
  } catch (err) {
    console.error('Error cargando métricas médicas:', err)
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

    await Promise.all([
      loadTrend(), 
      loadPavilionDistribution(), 
      loadStoreDistribution(), 
      loadTopSupplies(),
      loadMedicalMetrics()
    ])

    movementBars.value.consumos = Number(summary.value.total_consumed || 0)
    movementBars.value.entradas = Number(summary.value.total_in_stores || 0)
    movementBars.value.salidas = Number(summary.value.total_transferred || 0)
  } catch (e) {
    error.value = e?.message || 'Error cargando estadísticas'
  } finally {
    loading.value = false
  }
}

// Funciones para actualización automática
const startAutoRefresh = () => {
  if (autoRefreshInterval.value) {
    clearInterval(autoRefreshInterval.value)
  }
  
  autoRefreshInterval.value = setInterval(() => {
    if (autoRefreshEnabled.value && !loading.value) {
      loadData()
    }
  }, refreshIntervalSeconds * 1000)
}

const stopAutoRefresh = () => {
  if (autoRefreshInterval.value) {
    clearInterval(autoRefreshInterval.value)
    autoRefreshInterval.value = null
  }
}

watch(transferRange, () => {
  loadTrend()
})

onMounted(() => {
  loadData()
  startAutoRefresh()
})

onUnmounted(() => {
  stopAutoRefresh()
})
</script>

<style scoped>
.hover\:shadow-md:hover {
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
}

.transition-shadow {
  transition: box-shadow 0.2s ease-in-out;
}

.bg-gradient-to-r {
  background-image: linear-gradient(to right, var(--tw-gradient-stops));
}

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

.space-y-4 > *, .space-y-6 > * {
  animation: fadeInUp 0.6s ease-out;
}
</style>