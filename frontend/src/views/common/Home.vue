<template>
  <div class="space-y-6">
    <!-- Header con saludo -->
    <div class="bg-gradient-to-r from-brand-blue-dark to-brand-blue-medium rounded-xl p-4 sm:p-6 text-white">
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 class="text-2xl font-bold text-white mb-1">
            Bienvenido{{ authStore.getUserName ? ', ' + authStore.getUserName : '' }}
          </h1>
          <p class="text-blue-100 mt-1">
            {{ authStore.isDoctor ? 'Panel de Solicitudes Médicas' : authStore.isPavedad ? 'Panel de Visualización de Solicitudes' : authStore.isPavilionUser ? 'Panel de Gestión de Pabellón' : authStore.isNurse ? 'Panel de Gestión de Inventario' : authStore.isConsignation ? 'Panel de Gestión y Configuración' : 'Sistema de gestión de inventario médico' }}
          </p>
        </div>
        <div class="sm:text-right">
          <p class="text-xs sm:text-sm text-blue-100">Fecha actual</p>
          <p class="text-base sm:text-xl font-semibold">{{ currentDate }}</p>
        </div>
      </div>
    </div>

    <!-- Barra de búsqueda principal - Solo para roles que pueden ver inventario (excepto pabellón y consignación) -->
    <div v-if="!authStore.isDoctor && !authStore.isPavedad && !authStore.isNurse && !authStore.isPavilionUser && !authStore.isConsignation" class="card">
      <label for="search" class="block text-sm font-medium text-gray-700 mb-2">
        Buscar insumo médico
      </label>
      <div class="flex items-end gap-4">
        <div class="flex-1">
          <div class="relative">
            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <svg class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
            </div>
            <input
              id="search"
              type="text"
              placeholder="Ingrese nombre, código o lote del insumo..."
              class="form-input pl-10 w-full"
              v-model="searchQuery"
              @keyup.enter="handleSearch"
            />
          </div>
        </div>
        <div class="flex flex-col sm:flex-row gap-3 md:flex-shrink-0">
          <button @click="handleSearch" class="btn-primary w-full sm:w-auto flex items-center justify-center whitespace-nowrap">
            <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
            Buscar
          </button>
          <button class="btn-secondary w-full sm:w-auto flex items-center justify-center whitespace-nowrap">
            <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 4a1 1 0 011-1h16a1 1 0 011 1v2.586a1 1 0 01-.293.707l-6.414 6.414a1 1 0 00-.293.707V17l-4 4v-6.586a1 1 0 00-.293-.707L3.293 7.207A1 1 0 013 6.5V4z" />
            </svg>
            Filtros
          </button>
        </div>
      </div>
    </div>

    <!-- Funcionalidades principales -->
    <!-- Vista específica para doctores -->
    <div v-if="authStore.isDoctor" class="max-w-2xl mx-auto">
      <div class="text-center mb-8">
        <p class="text-gray-600">Gestiona tus solicitudes de insumos médicos de manera rápida y eficiente.</p>
      </div>
      
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <!-- Nueva Solicitud de Insumo -->
        <div class="card hover:shadow-lg transition-all duration-200 cursor-pointer transform hover:scale-105" @click="navigateTo('/supply-requests/new')">
          <div class="text-center p-6">
            <div class="bg-teal-100 p-4 rounded-full mx-auto w-20 h-20 flex items-center justify-center mb-4">
              <svg class="h-10 w-10 text-teal-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
            </div>
            <h3 class="text-xl font-semibold text-gray-900 mb-2">Nueva Solicitud</h3>
            <p class="text-gray-600">Crear una nueva solicitud de insumos médicos para tu área de trabajo</p>
          </div>
        </div>

        <!-- Gestión de Solicitudes -->
        <div class="card hover:shadow-lg transition-all duration-200 cursor-pointer transform hover:scale-105" @click="navigateTo('/supply-requests')">
          <div class="text-center p-6">
            <div class="bg-indigo-100 p-4 rounded-full mx-auto w-20 h-20 flex items-center justify-center mb-4">
              <svg class="h-10 w-10 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01" />
              </svg>
            </div>
            <h3 class="text-xl font-semibold text-gray-900 mb-2">Mis Solicitudes</h3>
            <p class="text-gray-600">Ver el estado y gestionar todas tus solicitudes de insumos</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Vista específica para enfermeras -->
    <div v-else-if="authStore.isNurse" class="max-w-2xl mx-auto">
      <div class="text-center mb-8">
        <p class="text-gray-600">Gestiona el inventario de pabellones y escanea códigos QR de manera rápida y eficiente.</p>
      </div>
      
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <!-- Inventario de Pabellones -->
        <div class="card hover:shadow-lg transition-all duration-200 cursor-pointer transform hover:scale-105" @click="navigateTo('/inventory/pavilion')">
          <div class="text-center p-6">
            <div class="bg-purple-100 p-4 rounded-full mx-auto w-20 h-20 flex items-center justify-center mb-4">
              <svg class="h-10 w-10 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
              </svg>
            </div>
            <h3 class="text-xl font-semibold text-gray-900 mb-2">Inventario de Pabellones</h3>
            <p class="text-gray-600">Ver el stock disponible en cada pabellón del hospital</p>
          </div>
        </div>

        <!-- Escaner QR -->
        <div class="card hover:shadow-lg transition-all duration-200 cursor-pointer transform hover:scale-105" @click="navigateTo('/qr')">
          <div class="text-center p-6">
            <div class="bg-yellow-100 p-4 rounded-full mx-auto w-20 h-20 flex items-center justify-center mb-4">
              <svg class="h-10 w-10 text-yellow-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h.01M5 8h2a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v2a2 2 0 002 2zm0 0h2a2 2 0 012 2v2a2 2 0 01-2 2H5a2 2 0 01-2-2v-2a2 2 0 012-2z" />
              </svg>
            </div>
            <h3 class="text-xl font-semibold text-gray-900 mb-2">Escáner QR</h3>
            <p class="text-gray-600">Escanear códigos QR para trazabilidad y verificación de insumos</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Vista específica para perfil de pabellón -->
    <div v-else-if="authStore.isPavilionUser" class="max-w-2xl mx-auto">
      <div class="text-center mb-8">
        <p class="text-gray-600">Gestiona el inventario de tu pabellón y escanea códigos QR de manera rápida y eficiente.</p>
      </div>
      
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <!-- Inventario de Pabellones -->
        <div class="card hover:shadow-lg transition-all duration-200 cursor-pointer transform hover:scale-105" @click="navigateTo('/inventory/pavilion')">
          <div class="text-center p-6">
            <div class="bg-purple-100 p-4 rounded-full mx-auto w-20 h-20 flex items-center justify-center mb-4">
              <svg class="h-10 w-10 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
              </svg>
            </div>
            <h3 class="text-xl font-semibold text-gray-900 mb-2">Inventario del Pabellón</h3>
            <p class="text-gray-600">Ver el stock disponible en tu pabellón</p>
          </div>
        </div>

        <!-- Escaner QR -->
        <div class="card hover:shadow-lg transition-all duration-200 cursor-pointer transform hover:scale-105" @click="navigateTo('/qr')">
          <div class="text-center p-6">
            <div class="bg-yellow-100 p-4 rounded-full mx-auto w-20 h-20 flex items-center justify-center mb-4">
              <svg class="h-10 w-10 text-yellow-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h.01M5 8h2a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v2a2 2 0 002 2zm0 0h2a2 2 0 012 2v2a2 2 0 01-2 2H5a2 2 0 01-2-2v-2a2 2 0 012-2z" />
              </svg>
            </div>
            <h3 class="text-xl font-semibold text-gray-900 mb-2">Escáner QR</h3>
            <p class="text-gray-600">Escanear códigos QR para trazabilidad y verificación de insumos</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Vista específica para consignación -->
    <div v-else-if="authStore.isConsignation" class="max-w-4xl mx-auto">
      <div class="text-center mb-8">
        <p class="text-gray-600">Gestiona solicitudes, configuración y escanea códigos QR de manera rápida y eficiente.</p>
      </div>
      
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <!-- Ver Todas las Solicitudes -->
        <div class="card hover:shadow-lg transition-all duration-200 cursor-pointer transform hover:scale-105" @click="navigateTo('/supply-requests')">
          <div class="text-center p-6">
            <div class="bg-indigo-100 p-4 rounded-full mx-auto w-20 h-20 flex items-center justify-center mb-4">
              <svg class="h-10 w-10 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01" />
              </svg>
            </div>
            <h3 class="text-xl font-semibold text-gray-900 mb-2">Ver Solicitudes</h3>
            <p class="text-gray-600">Ver todas las solicitudes de insumos médicos del sistema</p>
          </div>
        </div>

        <!-- Escaner QR -->
        <div class="card hover:shadow-lg transition-all duration-200 cursor-pointer transform hover:scale-105" @click="navigateTo('/qr')">
          <div class="text-center p-6">
            <div class="bg-yellow-100 p-4 rounded-full mx-auto w-20 h-20 flex items-center justify-center mb-4">
              <svg class="h-10 w-10 text-yellow-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h.01M5 8h2a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v2a2 2 0 002 2zm0 0h2a2 2 0 012 2v2a2 2 0 01-2 2H5a2 2 0 01-2-2v-2a2 2 0 012-2z" />
              </svg>
            </div>
            <h3 class="text-xl font-semibold text-gray-900 mb-2">Escáner QR</h3>
            <p class="text-gray-600">Escanear códigos QR para trazabilidad y verificación de insumos</p>
          </div>
        </div>

        <!-- Configuración de Proveedores -->
        <div class="card hover:shadow-lg transition-all duration-200 cursor-pointer transform hover:scale-105" @click="navigateTo('/supplier-configs')">
          <div class="text-center p-6">
            <div class="bg-blue-100 p-4 rounded-full mx-auto w-20 h-20 flex items-center justify-center mb-4">
              <svg class="h-10 w-10 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
              </svg>
            </div>
            <h3 class="text-xl font-semibold text-gray-900 mb-2">Configuración de Proveedores</h3>
            <p class="text-gray-600">Gestionar alertas de vencimiento por proveedor</p>
          </div>
        </div>

        <!-- Códigos de Insumos -->
        <div class="card hover:shadow-lg transition-all duration-200 cursor-pointer transform hover:scale-105" @click="navigateTo('/supply-codes')">
          <div class="text-center p-6">
            <div class="bg-green-100 p-4 rounded-full mx-auto w-20 h-20 flex items-center justify-center mb-4">
              <svg class="h-10 w-10 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z" />
              </svg>
            </div>
            <h3 class="text-xl font-semibold text-gray-900 mb-2">Códigos de Insumos</h3>
            <p class="text-gray-600">Gestionar códigos de insumos y niveles críticos de stock</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Vista específica para pavedad -->
    <div v-else-if="authStore.isPavedad" class="max-w-2xl mx-auto">
      <div class="text-center mb-8">
        <p class="text-gray-600">Visualiza y monitorea todas las solicitudes de insumos médicos del sistema.</p>
      </div>
      
      <div class="grid grid-cols-1 gap-6">
        <!-- Ver Todas las Solicitudes -->
        <div class="card hover:shadow-lg transition-all duration-200 cursor-pointer transform hover:scale-105" @click="navigateTo('/supply-requests')">
          <div class="text-center p-8">
            <div class="bg-indigo-100 p-4 rounded-full mx-auto w-24 h-24 flex items-center justify-center mb-4">
              <svg class="h-12 w-12 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01" />
              </svg>
            </div>
            <h3 class="text-2xl font-semibold text-gray-900 mb-3">Ver Todas las Solicitudes</h3>
            <p class="text-gray-600 text-lg">Monitorea el estado de todas las solicitudes de insumos médicos del sistema</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Vista completa para otros roles -->
    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <!-- Gestión de inventario (solo admin y encargado de bodega) -->
      <div v-if="authStore.canViewInventoryMenu" class="card hover:shadow-lg transition-all duration-200 cursor-pointer" @click="navigateTo('/inventory')">
        <div class="flex items-center gap-4">
          <div class="bg-blue-100 p-3 rounded-lg">
            <svg class="h-8 w-8 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
            </svg>
          </div>
          <div class="min-w-0 flex-1">
            <h3 class="text-base sm:text-lg font-semibold text-gray-900 truncate">Gestión de Inventario</h3>
            <p class="text-gray-600 text-xs sm:text-sm truncate">Administrar insumos médicos</p>
          </div>
        </div>
      </div>

      <!-- Agregar nuevo insumo -->
      <div v-if="authStore.canAddSupplies" class="card hover:shadow-lg transition-all duration-200 cursor-pointer" @click="navigateTo('/inventory/add')">
        <div class="flex items-center gap-4">
          <div class="bg-green-100 p-3 rounded-lg">
            <svg class="h-8 w-8 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
            </svg>
          </div>
          <div class="min-w-0 flex-1">
            <h3 class="text-base sm:text-lg font-semibold text-gray-900 truncate">Agregar Insumo</h3>
            <p class="text-gray-600 text-xs sm:text-sm truncate">Registrar nuevo producto médico</p>
          </div>
        </div>
      </div>

      <!-- Escanear QR -->
      <div v-if="authStore.canViewQR" class="card hover:shadow-lg transition-all duration-200 cursor-pointer" @click="navigateTo('/qr')">
        <div class="flex items-center gap-4">
          <div class="bg-yellow-100 p-3 rounded-lg">
            <svg class="h-8 w-8 text-yellow-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h.01M5 8h2a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v2a2 2 0 002 2zm0 0h2a2 2 0 012 2v2a2 2 0 01-2 2H5a2 2 0 01-2-2v-2a2 2 0 012-2z" />
            </svg>
          </div>
          <div class="min-w-0 flex-1">
            <h3 class="text-base sm:text-lg font-semibold text-gray-900 truncate">Escáner QR</h3>
            <p class="text-gray-600 text-xs sm:text-sm truncate">Trazabilidad y verificación</p>
          </div>
        </div>
      </div>

      <!-- Estadísticas -->
      <div v-if="authStore.canViewStatistics" class="card hover:shadow-lg transition-all duration-200 cursor-pointer" @click="navigateTo('/statistics')">
        <div class="flex items-center gap-4">
          <div class="bg-purple-100 p-3 rounded-lg">
            <svg class="h-8 w-8 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
            </svg>
          </div>
          <div class="min-w-0 flex-1">
            <h3 class="text-base sm:text-lg font-semibold text-gray-900 truncate">Estadísticas</h3>
            <p class="text-gray-600 text-xs sm:text-sm truncate">Análisis y métricas del sistema</p>
          </div>
        </div>
      </div>

      <!-- Nueva Solicitud de Insumo -->
      <div v-if="authStore.canCreateRequests" class="card hover:shadow-lg transition-all duration-200 cursor-pointer" @click="navigateTo('/supply-requests/new')">
        <div class="flex items-center gap-4">
          <div class="bg-teal-100 p-3 rounded-lg">
            <svg class="h-8 w-8 text-teal-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
          </div>
          <div class="min-w-0 flex-1">
            <h3 class="text-base sm:text-lg font-semibold text-gray-900 truncate">Nueva Solicitud</h3>
            <p class="text-gray-600 text-xs sm:text-sm truncate">Crear solicitud de insumos</p>
          </div>
        </div>
      </div>


      <!-- Gestión de Solicitudes -->
      <div v-if="authStore.canViewRequests" class="card hover:shadow-lg transition-all duration-200 cursor-pointer" @click="navigateTo('/supply-requests')">
        <div class="flex items-center gap-4">
          <div class="bg-indigo-100 p-3 rounded-lg">
            <svg class="h-8 w-8 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01" />
            </svg>
          </div>
          <div class="min-w-0 flex-1">
            <h3 class="text-base sm:text-lg font-semibold text-gray-900 truncate">Gestión de Solicitudes</h3>
            <p class="text-gray-600 text-xs sm:text-sm truncate">Ver y gestionar solicitudes</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()
const searchQuery = ref('')

const userName = computed(() => {
  return authStore.getUserName || 'Usuario'
})

const currentDate = computed(() => {
  return new Date().toLocaleDateString('es-ES', {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
})

const handleSearch = () => {
  // Los doctores y pavedad no pueden acceder al inventario, redireccionar a solicitudes
  if (authStore.isDoctor || authStore.isPavedad) {
    router.push('/supply-requests')
    return
  }
  
  if (searchQuery.value.trim()) {
    router.push({
      name: 'Inventory',
      query: { search: searchQuery.value }
    })
  } else {
    router.push('/inventory')
  }
}

const navigateTo = (path) => {
  router.push(path)
}
</script>

<style scoped>
/* Efectos hover para las tarjetas */
.card:hover {
  transform: translateY(-2px);
  transition: all 0.2s ease-in-out;
}

/* Gradientes personalizados */
.bg-gradient-to-r {
  background-image: linear-gradient(to right, var(--tw-gradient-stops));
}

/* Transiciones suaves */
.transition-all {
  transition: all 0.2s ease-in-out;
}
</style>