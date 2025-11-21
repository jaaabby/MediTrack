<template>
  <div id="app" class="min-h-screen bg-gray-50">
    <!-- Header principal - Oculto en login y registro -->
    <header v-if="$route.name !== 'Login' && $route.name !== 'Register'" class="bg-brand-blue-dark shadow-sm border-b border-brand-blue-medium">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center h-16">
          <!-- Logo y título -->
          <div class="flex items-center flex-shrink-0">
            <router-link to="/" class="cursor-pointer hover:opacity-80 transition-opacity flex items-center">
              <img 
                :src="logoImage" 
                alt="MediTrack" 
                class="h-12 sm:h-14 md:h-16 w-auto brightness-0 invert"
              />
            </router-link>
          </div>
          
          <!-- Navegación principal (desktop y tablet) - Oculta en login y registro -->
          <nav v-if="$route.name !== 'Login' && $route.name !== 'Register'" class="hidden lg:flex items-center space-x-4 xl:space-x-8 flex-1 justify-center max-w-4xl">
            <router-link v-if="authStore.isAuthenticated && authStore.canViewHome"
              to="/"
              class="text-white hover:text-brand-blue-light px-2 xl:px-3 py-2 rounded-md text-sm font-medium transition-colors whitespace-nowrap"
              :class="{ 'bg-brand-blue-medium': $route.path === '/' }"
              @click.stop
            >
              Inicio
            </router-link>
            
            <!-- Menú desplegable Inventario (solo admin y encargado de bodega) -->
            <div v-if="authStore.isAuthenticated && authStore.canViewInventoryMenu" class="relative">
              <button
                @click="inventoryMenuOpen = !inventoryMenuOpen"
                class="text-white hover:text-brand-blue-light px-2 xl:px-3 py-2 rounded-md text-sm font-medium transition-colors whitespace-nowrap flex items-center"
                :class="{ 'bg-brand-blue-medium': $route.path.startsWith('/inventory') }"
              >
                Inventario
                <svg class="ml-1 h-4 w-4" :class="{ 'rotate-180': inventoryMenuOpen }" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                </svg>
              </button>
              
              <!-- Dropdown menu -->
              <div v-if="inventoryMenuOpen" class="absolute left-0 mt-2 w-56 bg-white rounded-md shadow-lg z-50 border border-gray-200">
                <router-link
                  to="/inventory"
                  class="block px-4 py-2 text-sm text-gray-700 hover:bg-brand-blue-light hover:bg-opacity-20 hover:text-brand-blue-dark transition-colors"
                  :class="{ 'bg-brand-blue-light bg-opacity-20 text-brand-blue-dark': $route.path === '/inventory' }"
                  @click="inventoryMenuOpen = false"
                >
                  <div class="flex items-center">
                    <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 10h16M4 14h16M4 18h16" />
                    </svg>
                    Ver Lotes
                  </div>
                </router-link>
                <router-link
                  to="/inventory/dashboard"
                  class="block px-4 py-2 text-sm text-gray-700 hover:bg-brand-blue-light hover:bg-opacity-20 hover:text-brand-blue-dark transition-colors"
                  :class="{ 'bg-brand-blue-light bg-opacity-20 text-brand-blue-dark': $route.path === '/inventory/dashboard' }"
                  @click="inventoryMenuOpen = false"
                >
                  <div class="flex items-center">
                    <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
                    </svg>
                    Dashboard
                  </div>
                </router-link>
                <router-link
                  to="/inventory/store"
                  class="block px-4 py-2 text-sm text-gray-700 hover:bg-brand-blue-light hover:bg-opacity-20 hover:text-brand-blue-dark transition-colors"
                  :class="{ 'bg-brand-blue-light bg-opacity-20 text-brand-blue-dark': $route.path === '/inventory/store' }"
                  @click="inventoryMenuOpen = false"
                >
                  <div class="flex items-center">
                    <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
                    </svg>
                    Bodegas
                  </div>
                </router-link>
                <router-link
                  to="/inventory/pavilion"
                  class="block px-4 py-2 text-sm text-gray-700 hover:bg-brand-blue-light hover:bg-opacity-20 hover:text-brand-blue-dark transition-colors"
                  :class="{ 'bg-brand-blue-light bg-opacity-20 text-brand-blue-dark': $route.path === '/inventory/pavilion' }"
                  @click="inventoryMenuOpen = false"
                >
                  <div class="flex items-center">
                    <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
                    </svg>
                    Pabellones
                  </div>
                </router-link>
                <router-link
                  v-if="authStore.canAddSupplies"
                  to="/inventory/add"
                  class="block px-4 py-2 text-sm text-gray-700 hover:bg-brand-blue-light hover:bg-opacity-20 hover:text-brand-blue-dark transition-colors border-t border-gray-200"
                  :class="{ 'bg-brand-blue-light bg-opacity-20 text-brand-blue-dark': $route.path === '/inventory/add' }"
                  @click="inventoryMenuOpen = false"
                >
                  <div class="flex items-center">
                    <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
                    </svg>
                    Agregar Insumo
                  </div>
                </router-link>
              </div>
            </div>
            
            <!-- Enlace directo a Pabellones para enfermeras -->
            <router-link v-if="authStore.isAuthenticated && authStore.isNurse"
              to="/inventory/pavilion"
              class="text-white hover:text-brand-blue-light px-2 xl:px-3 py-2 rounded-md text-sm font-medium transition-colors whitespace-nowrap"
              :class="{ 'bg-brand-blue-medium': $route.path === '/inventory/pavilion' }"
              @click.stop
            >
              Pabellones
            </router-link>
            
            <router-link v-if="authStore.isAuthenticated && authStore.canViewRequests"
              to="/supply-requests"
              class="text-white hover:text-brand-blue-light px-2 xl:px-3 py-2 rounded-md text-sm font-medium transition-colors whitespace-nowrap"
              :class="{ 'bg-brand-blue-medium': $route.path.startsWith('/supply-requests') }"
              @click.stop
            >
              Solicitudes
            </router-link>
            
            <router-link v-if="authStore.isAuthenticated && authStore.canViewQR"
              to="/qr"
              class="text-white hover:text-brand-blue-light px-2 xl:px-3 py-2 rounded-md text-sm font-medium transition-colors whitespace-nowrap"
              :class="{ 'bg-brand-blue-medium': $route.path.startsWith('/qr') }"
              @click.stop
            >
              Escaner QR
            </router-link>

            <router-link v-if="authStore.isAuthenticated && authStore.canViewStatistics"
              to="/statistics"
              class="text-white hover:text-brand-blue-light px-2 xl:px-3 py-2 rounded-md text-sm font-medium transition-colors whitespace-nowrap"
              :class="{ 'bg-brand-blue-medium': $route.path === '/statistics' }"
              @click.stop
            >
              Estadisticas
            </router-link>

            <!-- Menú desplegable Gestión -->
            <div v-if="authStore.isAuthenticated && !authStore.isDoctor && !authStore.isPavedad && !authStore.isNurse" class="relative">
              <button
                @click="managementMenuOpen = !managementMenuOpen"
                class="text-white hover:text-brand-blue-light px-2 xl:px-3 py-2 rounded-md text-sm font-medium transition-colors whitespace-nowrap flex items-center"
                :class="{ 'bg-brand-blue-medium': ['/transfers', '/surgeries', '/supply-history', '/return-management', '/medical-specialties', '/surgery-typical-supplies', '/doctor-info', '/supplier-configs', '/supply-codes'].some(path => $route.path.startsWith(path)) }"
              >
                Gestión
                <svg class="ml-1 h-4 w-4" :class="{ 'rotate-180': managementMenuOpen }" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                </svg>
              </button>
              
              <!-- Dropdown menu -->
              <div v-if="managementMenuOpen" class="absolute left-0 mt-2 w-56 bg-white rounded-md shadow-lg z-50 border border-gray-200 max-h-96 overflow-y-auto">
                <router-link
                  v-if="authStore.canManageTransfers"
                  to="/transfers"
                  class="block px-4 py-2 text-sm text-gray-700 hover:bg-brand-blue-light hover:bg-opacity-20 hover:text-brand-blue-dark transition-colors"
                  :class="{ 'bg-brand-blue-light bg-opacity-20 text-brand-blue-dark': $route.path === '/transfers' }"
                  @click="managementMenuOpen = false"
                >
                  <div class="flex items-center">
                    <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7h12m0 0l-4-4m4 4l-4 4m0 6H4m0 0l4 4m-4-4l4-4" />
                    </svg>
                    Transferencias
                  </div>
                </router-link>
                <router-link
                  v-if="authStore.canManageSurgeries"
                  to="/surgeries"
                  class="block px-4 py-2 text-sm text-gray-700 hover:bg-brand-blue-light hover:bg-opacity-20 hover:text-brand-blue-dark transition-colors"
                  :class="{ 'bg-brand-blue-light bg-opacity-20 text-brand-blue-dark': $route.path === '/surgeries' }"
                  @click="managementMenuOpen = false"
                >
                  <div class="flex items-center">
                    <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                    </svg>
                    Cirugías
                  </div>
                </router-link>
                <router-link
                  v-if="authStore.canViewSupplyHistory"
                  to="/supply-history"
                  class="block px-4 py-2 text-sm text-gray-700 hover:bg-brand-blue-light hover:bg-opacity-20 hover:text-brand-blue-dark transition-colors"
                  :class="{ 'bg-brand-blue-light bg-opacity-20 text-brand-blue-dark': $route.path === '/supply-history' }"
                  @click="managementMenuOpen = false"
                >
                  <div class="flex items-center">
                    <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                    Historial de Insumos
                  </div>
                </router-link>
                <router-link
                  v-if="authStore.canManageReturns"
                  to="/return-management"
                  class="block px-4 py-2 text-sm text-gray-700 hover:bg-brand-blue-light hover:bg-opacity-20 hover:text-brand-blue-dark transition-colors"
                  :class="{ 'bg-brand-blue-light bg-opacity-20 text-brand-blue-dark': $route.path === '/return-management' }"
                  @click="managementMenuOpen = false"
                >
                  <div class="flex items-center">
                    <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 10h10a8 8 0 018 8v2M3 10l6 6m-6-6l6-6" />
                    </svg>
                    Retornos a Bodega
                  </div>
                </router-link>
                
                <!-- Separador -->
                <div v-if="authStore.canManageMedicalConfig" class="border-t border-gray-200 my-1"></div>
                
                <div v-if="authStore.canManageMedicalConfig" class="px-4 py-2 text-xs font-semibold text-gray-500 uppercase tracking-wider">
                  Configuración Médica
                </div>
                
                <router-link
                  v-if="authStore.canManageMedicalConfig"
                  to="/medical-specialties"
                  class="block px-4 py-2 text-sm text-gray-700 hover:bg-brand-blue-light hover:bg-opacity-20 hover:text-brand-blue-dark transition-colors"
                  :class="{ 'bg-brand-blue-light bg-opacity-20 text-brand-blue-dark': $route.path === '/medical-specialties' }"
                  @click="managementMenuOpen = false"
                >
                  <div class="flex items-center">
                    <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                    </svg>
                    Especialidades Médicas
                  </div>
                </router-link>
                <router-link
                  v-if="authStore.canManageMedicalConfig"
                  to="/surgery-typical-supplies"
                  class="block px-4 py-2 text-sm text-gray-700 hover:bg-brand-blue-light hover:bg-opacity-20 hover:text-brand-blue-dark transition-colors"
                  :class="{ 'bg-brand-blue-light bg-opacity-20 text-brand-blue-dark': $route.path === '/surgery-typical-supplies' }"
                  @click="managementMenuOpen = false"
                >
                  <div class="flex items-center">
                    <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
                    </svg>
                    Insumos Típicos
                  </div>
                </router-link>
                <router-link
                  v-if="authStore.canManageMedicalConfig"
                  to="/doctor-info"
                  class="block px-4 py-2 text-sm text-gray-700 hover:bg-brand-blue-light hover:bg-opacity-20 hover:text-brand-blue-dark transition-colors"
                  :class="{ 'bg-brand-blue-light bg-opacity-20 text-brand-blue-dark': $route.path === '/doctor-info' }"
                  @click="managementMenuOpen = false"
                >
                  <div class="flex items-center">
                    <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                    </svg>
                    Información de Doctores
                  </div>
                </router-link>
                
                <!-- Separador -->
                <div v-if="authStore.canManageSystemConfig" class="border-t border-gray-200 my-1"></div>
                
                <div v-if="authStore.canManageSystemConfig" class="px-4 py-2 text-xs font-semibold text-gray-500 uppercase tracking-wider">
                  Configuración del Sistema
                </div>
                
                <router-link
                  v-if="authStore.canManageSystemConfig"
                  to="/supplier-configs"
                  class="block px-4 py-2 text-sm text-gray-700 hover:bg-brand-blue-light hover:bg-opacity-20 hover:text-brand-blue-dark transition-colors"
                  :class="{ 'bg-brand-blue-light bg-opacity-20 text-brand-blue-dark': $route.path === '/supplier-configs' }"
                  @click="managementMenuOpen = false"
                >
                  <div class="flex items-center">
                    <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                    </svg>
                    Configuración de Proveedores
                  </div>
                </router-link>
                <router-link
                  v-if="authStore.canManageSystemConfig"
                  to="/supply-codes"
                  class="block px-4 py-2 text-sm text-gray-700 hover:bg-brand-blue-light hover:bg-opacity-20 hover:text-brand-blue-dark transition-colors"
                  :class="{ 'bg-brand-blue-light bg-opacity-20 text-brand-blue-dark': $route.path === '/supply-codes' }"
                  @click="managementMenuOpen = false"
                >
                  <div class="flex items-center">
                    <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z" />
                    </svg>
                    Códigos de Insumos
                  </div>
                </router-link>
              </div>
            </div>
          </nav>

          <!-- Navegación tablet (oculta en desktop y mobile) -->
          <nav class="hidden md:flex lg:hidden items-center space-x-2 flex-1 justify-center">
            <router-link v-if="authStore.isAuthenticated && authStore.canViewHome"
              to="/"
              class="text-white hover:text-brand-blue-light p-2 rounded-md transition-colors"
              :class="{ 'bg-brand-blue-medium': $route.path === '/' }"
              @click.stop
              title="Inicio"
            >
              <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
              </svg>
            </router-link>
            
            <router-link v-if="authStore.isAuthenticated && authStore.canViewInventory"
              to="/inventory"
              class="text-white hover:text-brand-blue-light p-2 rounded-md transition-colors"
              :class="{ 'bg-brand-blue-medium': $route.path === '/inventory' }"
              @click.stop
              title="Inventario"
            >
              <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
              </svg>
            </router-link>
            
            <router-link v-if="authStore.isAuthenticated && authStore.canViewRequests"
              to="/supply-requests"
              class="text-white hover:text-brand-blue-light p-2 rounded-md transition-colors"
              :class="{ 'bg-brand-blue-medium': $route.path.startsWith('/supply-requests') }"
              @click.stop
              title="Solicitudes"
            >
              <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01" />
              </svg>
            </router-link>
            
            <router-link v-if="authStore.isAuthenticated && authStore.canViewQR"
              to="/qr"
              class="text-white hover:text-brand-blue-light p-2 rounded-md transition-colors"
              :class="{ 'bg-brand-blue-medium': $route.path.startsWith('/qr') }"
              @click.stop
              title="Escaner QR"
            >
              <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h.01M5 8h2a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v2a2 2 0 002 2zm0 0h2a2 2 0 012 2v2a2 2 0 01-2 2H5a2 2 0 01-2-2v-2a2 2 0 012-2z" />
              </svg>
            </router-link>

            <router-link v-if="authStore.isAuthenticated && authStore.canViewStatistics"
              to="/statistics"
              class="text-white hover:text-brand-blue-light p-2 rounded-md transition-colors"
              :class="{ 'bg-brand-blue-medium': $route.path === '/statistics' }"
              @click.stop
              title="Estadisticas"
            >
              <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
              </svg>
            </router-link>
          </nav>
          
          <!-- Menu de usuario - Oculto en login y registro -->
          <div v-if="$route.name !== 'Login' && $route.name !== 'Register'" class="flex items-center space-x-2 sm:space-x-4 flex-shrink-0">
            <!-- Informacion del usuario autenticado -->
            <div v-if="authStore.isAuthenticated" class="flex items-center space-x-2 sm:space-x-3">
              <div class="text-right text-white hidden sm:block">
                <router-link 
                  to="/profile"
                  class="text-sm font-medium hover:text-brand-blue-light transition-colors cursor-pointer block"
                  title="Ver mi perfil"
                >
                  {{ authStore.getUserName }}
                </router-link>
                <p class="text-xs text-brand-blue-light">{{ authStore.getUserRole }}</p>
              </div>
              
              <!-- Icono de usuario para pantallas pequeñas -->
              <router-link 
                to="/profile"
                class="sm:hidden p-2 rounded-lg hover:bg-brand-blue-medium transition-colors"
                title="Ver mi perfil"
              >
                <svg class="h-6 w-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                </svg>
              </router-link>
              
              <button 
                @click="handleLogout"
                class="p-2 rounded-lg hover:bg-brand-blue-medium transition-colors"
                title="Cerrar sesion"
              >
                <svg class="h-5 w-5 sm:h-6 sm:w-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
                </svg>
              </button>
            </div>
            
            <!-- Menu hamburguesa (movil) -->
            <button 
              v-if="authStore.isAuthenticated"
              @click="mobileMenuOpen = !mobileMenuOpen"
              class="md:hidden p-2 rounded-lg hover:bg-brand-blue-medium transition-colors"
              aria-label="Abrir menu"
            >
              <svg class="h-6 w-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path v-if="!mobileMenuOpen" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
                <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
        </div>
      </div>

      <!-- Menu movil -->
      <div v-if="mobileMenuOpen" class="md:hidden">
        <div class="px-2 pt-2 pb-3 space-y-1 sm:px-3 bg-brand-blue-medium">
          <router-link
            v-if="authStore.isAuthenticated && authStore.canViewHome"
            to="/"
            @click.stop="mobileMenuOpen = false"
            class="text-white hover:text-brand-blue-light block px-3 py-2 rounded-md text-base font-medium transition-colors"
            :class="{ 'bg-brand-blue-dark': $route.path === '/' }"
          >
            Inicio
          </router-link>
          
          <!-- Sección de Inventario (solo admin y encargado de bodega) -->
          <div v-if="authStore.isAuthenticated && authStore.canViewInventoryMenu" class="space-y-1">
            <div class="text-brand-blue-light px-3 py-2 text-sm font-semibold uppercase tracking-wide">
              Inventario
            </div>
            <router-link
              to="/inventory"
              @click.stop="mobileMenuOpen = false"
              class="text-white hover:text-brand-blue-light block px-3 py-2 rounded-md text-base font-medium transition-colors pl-6"
              :class="{ 'bg-brand-blue-dark': $route.path === '/inventory' }"
            >
              Ver Lotes
            </router-link>
            <router-link
              to="/inventory/dashboard"
              @click.stop="mobileMenuOpen = false"
              class="text-white hover:text-brand-blue-light block px-3 py-2 rounded-md text-base font-medium transition-colors pl-6"
              :class="{ 'bg-brand-blue-dark': $route.path === '/inventory/dashboard' }"
            >
              Dashboard
            </router-link>
            <router-link
              to="/inventory/store"
              @click.stop="mobileMenuOpen = false"
              class="text-white hover:text-brand-blue-light block px-3 py-2 rounded-md text-base font-medium transition-colors pl-6"
              :class="{ 'bg-brand-blue-dark': $route.path === '/inventory/store' }"
            >
              Bodegas
            </router-link>
            <router-link
              to="/inventory/pavilion"
              @click.stop="mobileMenuOpen = false"
              class="text-white hover:text-brand-blue-light block px-3 py-2 rounded-md text-base font-medium transition-colors pl-6"
              :class="{ 'bg-brand-blue-dark': $route.path === '/inventory/pavilion' }"
            >
              Pabellones
            </router-link>
            <router-link
              v-if="authStore.canAddSupplies"
              to="/inventory/add"
              @click.stop="mobileMenuOpen = false"
              class="text-white hover:text-brand-blue-light block px-3 py-2 rounded-md text-base font-medium transition-colors pl-6"
              :class="{ 'bg-brand-blue-dark': $route.path === '/inventory/add' }"
            >
              Agregar Insumo
            </router-link>
          </div>
          
          <!-- Separador -->
          <div class="border-t border-brand-blue-medium my-2"></div>
          
          <router-link
            v-if="authStore.isAuthenticated && authStore.canViewRequests"
            to="/supply-requests"
            @click.stop="mobileMenuOpen = false"
            class="text-white hover:text-brand-blue-light block px-3 py-2 rounded-md text-base font-medium transition-colors"
            :class="{ 'bg-brand-blue-dark': $route.path.startsWith('/supply-requests') }"
          >
            Solicitudes de Insumo
          </router-link>
          
          <router-link
            v-if="authStore.isAuthenticated && authStore.canViewQR"
            to="/qr"
            @click.stop="mobileMenuOpen = false"
            class="text-white hover:text-brand-blue-light block px-3 py-2 rounded-md text-base font-medium transition-colors"
            :class="{ 'bg-brand-blue-dark': $route.path.startsWith('/qr') }"
          >
            Escaner QR
          </router-link>

          <router-link
            v-if="authStore.isAuthenticated && authStore.canViewStatistics"
            to="/statistics"
            @click.stop="mobileMenuOpen = false"
            class="text-white hover:text-brand-blue-light block px-3 py-2 rounded-md text-base font-medium transition-colors"
            :class="{ 'bg-brand-blue-dark': $route.path === '/statistics' }"
          >
            Estadisticas
          </router-link>

          <!-- Separador -->
          <div class="border-t border-brand-blue-medium my-2"></div>
          
          <!-- Enlace directo a Pabellones para enfermeras en menú móvil -->
          <router-link v-if="authStore.isAuthenticated && authStore.isNurse"
            to="/inventory/pavilion"
            @click.stop="mobileMenuOpen = false"
            class="text-white hover:text-brand-blue-light block px-3 py-2 rounded-md text-base font-medium transition-colors"
            :class="{ 'bg-brand-blue-dark': $route.path === '/inventory/pavilion' }"
          >
            Pabellones
          </router-link>
          
          <!-- Sección de Gestión -->
          <div v-if="authStore.isAuthenticated && !authStore.isDoctor && !authStore.isPavedad && !authStore.isNurse" class="space-y-1">
            <div class="text-brand-blue-light px-3 py-2 text-sm font-semibold uppercase tracking-wide">
              Gestión
            </div>
            <router-link
              v-if="authStore.canManageTransfers"
              to="/transfers"
              @click.stop="mobileMenuOpen = false"
              class="text-white hover:text-brand-blue-light block px-3 py-2 rounded-md text-base font-medium transition-colors pl-6"
              :class="{ 'bg-brand-blue-dark': $route.path === '/transfers' }"
            >
              Transferencias
            </router-link>

            <router-link
              v-if="authStore.canManageSurgeries"
              to="/surgeries"
              @click.stop="mobileMenuOpen = false"
              class="text-white hover:text-brand-blue-light block px-3 py-2 rounded-md text-base font-medium transition-colors pl-6"
              :class="{ 'bg-brand-blue-dark': $route.path === '/surgeries' }"
            >
              Cirugías
            </router-link>

            <router-link
              v-if="authStore.canViewSupplyHistory"
              to="/supply-history"
              @click.stop="mobileMenuOpen = false"
              class="text-white hover:text-brand-blue-light block px-3 py-2 rounded-md text-base font-medium transition-colors pl-6"
              :class="{ 'bg-brand-blue-dark': $route.path === '/supply-history' }"
            >
              Historial de Insumos
            </router-link>
            
            <router-link
              v-if="authStore.canManageReturns"
              to="/return-management"
              @click.stop="mobileMenuOpen = false"
              class="text-white hover:text-brand-blue-light block px-3 py-2 rounded-md text-base font-medium transition-colors pl-6"
              :class="{ 'bg-brand-blue-dark': $route.path === '/return-management' }"
            >
              Retornos a Bodega
            </router-link>
            
            <!-- Separador -->
            <div v-if="authStore.canManageMedicalConfig" class="border-t border-brand-blue-medium my-2"></div>
            
            <div v-if="authStore.canManageMedicalConfig" class="text-brand-blue-light px-3 py-2 text-xs font-semibold uppercase tracking-wide">
              Configuración Médica
            </div>
            
            <router-link
              v-if="authStore.canManageMedicalConfig"
              to="/medical-specialties"
              @click.stop="mobileMenuOpen = false"
              class="text-white hover:text-brand-blue-light block px-3 py-2 rounded-md text-base font-medium transition-colors pl-6"
              :class="{ 'bg-brand-blue-dark': $route.path === '/medical-specialties' }"
            >
              Especialidades Médicas
            </router-link>
            
            <router-link
              v-if="authStore.canManageMedicalConfig"
              to="/surgery-typical-supplies"
              @click.stop="mobileMenuOpen = false"
              class="text-white hover:text-brand-blue-light block px-3 py-2 rounded-md text-base font-medium transition-colors pl-6"
              :class="{ 'bg-brand-blue-dark': $route.path === '/surgery-typical-supplies' }"
            >
              Insumos Típicos
            </router-link>
            
            <router-link
              v-if="authStore.canManageMedicalConfig"
              to="/doctor-info"
              @click.stop="mobileMenuOpen = false"
              class="text-white hover:text-brand-blue-light block px-3 py-2 rounded-md text-base font-medium transition-colors pl-6"
              :class="{ 'bg-brand-blue-dark': $route.path === '/doctor-info' }"
            >
              Información de Doctores
            </router-link>
            
            <!-- Separador -->
            <div v-if="authStore.canManageSystemConfig" class="border-t border-brand-blue-medium my-2"></div>
            
            <div v-if="authStore.canManageSystemConfig" class="text-brand-blue-light px-3 py-2 text-xs font-semibold uppercase tracking-wide">
              Configuración del Sistema
            </div>
            
            <router-link
              v-if="authStore.canManageSystemConfig"
              to="/supplier-configs"
              @click.stop="mobileMenuOpen = false"
              class="text-white hover:text-brand-blue-light block px-3 py-2 rounded-md text-base font-medium transition-colors pl-6"
              :class="{ 'bg-brand-blue-dark': $route.path === '/supplier-configs' }"
            >
              Configuración de Proveedores
            </router-link>
            <router-link
              v-if="authStore.canManageSystemConfig"
              to="/supply-codes"
              @click.stop="mobileMenuOpen = false"
              class="text-white hover:text-brand-blue-light block px-3 py-2 rounded-md text-base font-medium transition-colors pl-6"
              :class="{ 'bg-brand-blue-dark': $route.path === '/supply-codes' }"
            >
              Códigos de Insumos
            </router-link>
          </div>

        </div>
      </div>
    </header>

    <!-- Contenido principal -->
    <main :class="$route.name === 'Login' || $route.name === 'Register' ? '' : 'max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6'">
      <router-view />
    </main>

    <!-- Navegacion inferior (inspirada en la app movil) -->
    <nav class="fixed bottom-0 left-0 right-0 bg-white border-t border-gray-200 md:hidden z-40" v-if="authStore.isAuthenticated && $route.name !== 'Login' && $route.name !== 'Register'">
      <!-- Para doctores y pavedad: solo mostrar inicio y solicitudes -->
      <div v-if="authStore.isDoctor || authStore.isPavedad" class="grid grid-cols-2 gap-1">
        <router-link
          to="/"
          class="flex flex-col items-center py-3 px-6 text-sm font-medium transition-colors"
          :class="$route.path === '/' ? 'text-brand-blue-dark' : 'text-gray-500 hover:text-brand-blue-dark'"
        >
          <svg class="h-8 w-8 mb-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
          </svg>
          <span class="text-sm leading-tight">Inicio</span>
        </router-link>
        <router-link
          to="/supply-requests"
          class="flex flex-col items-center py-3 px-6 text-sm font-medium transition-colors"
          :class="$route.path.startsWith('/supply-requests') ? 'text-brand-blue-dark' : 'text-gray-500 hover:text-brand-blue-dark'"
        >
          <svg class="h-8 w-8 mb-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01" />
          </svg>
          <span class="text-sm leading-tight">Solicitudes</span>
        </router-link>
      </div>
      
      <!-- Para otros roles: navegación completa -->
      <div v-else class="grid grid-cols-5 gap-1">
        <router-link
          v-if="authStore.canViewHome"
          to="/"
          class="flex flex-col items-center py-2 px-1 text-xs font-medium transition-colors"
          :class="$route.path === '/' ? 'text-brand-blue-dark' : 'text-gray-500 hover:text-brand-blue-dark'"
        >
          <svg class="h-6 w-6 mb-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
          </svg>
          <span class="text-xs leading-tight">Inicio</span>
        </router-link>
        
        <router-link
          v-if="authStore.canViewInventory"
          to="/inventory"
          class="flex flex-col items-center py-2 px-1 text-xs font-medium transition-colors"
          :class="$route.path.startsWith('/inventory') ? 'text-brand-blue-dark' : 'text-gray-500 hover:text-brand-blue-dark'"
        >
          <svg class="h-6 w-6 mb-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
          </svg>
          <span class="text-xs leading-tight">Inventario</span>
        </router-link>
        
        <router-link
          v-if="authStore.canViewRequests"
          to="/supply-requests"
          class="flex flex-col items-center py-2 px-1 text-xs font-medium transition-colors relative"
          :class="$route.path.startsWith('/supply-requests') ? 'text-brand-blue-dark' : 'text-gray-500 hover:text-brand-blue-dark'"
        >
          <svg class="h-6 w-6 mb-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01" />
          </svg>
          <span class="text-xs leading-tight">Solicitudes</span>
        </router-link>
        
        <router-link
          v-if="authStore.canViewQR"
          to="/qr"
          class="flex flex-col items-center py-2 px-1 text-xs font-medium transition-colors"
          :class="$route.path.startsWith('/qr') ? 'text-brand-blue-dark' : 'text-gray-500 hover:text-brand-blue-dark'"
        >
          <svg class="h-6 w-6 mb-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h.01M5 8h2a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v2a2 2 0 002 2zm0 0h2a2 2 0 012 2v2a2 2 0 01-2 2H5a2 2 0 01-2-2v-2a2 2 0 012-2z" />
          </svg>
          <span class="text-xs leading-tight">Escaner</span>
        </router-link>

        <router-link
          v-if="authStore.canViewStatistics"
          to="/statistics"
          class="flex flex-col items-center py-2 px-1 text-xs font-medium transition-colors"
          :class="mobileMenuOpen ? 'text-brand-blue-dark' : 'text-gray-500 hover:text-brand-blue-dark'"
        >
          <svg class="h-6 w-6 mb-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
          </svg>
          <span class="text-xs leading-tight">Más</span>
        </router-link>
      </div>
    </nav>

    <!-- Padding para la navegacion inferior en moviles -->
    <div v-if="authStore.isAuthenticated && $route.name !== 'Login' && $route.name !== 'Register'" class="h-20 md:hidden"></div>

    <!-- Toast/Notificacion global (futuro) -->
    <!-- <div v-if="globalNotification" class="fixed bottom-4 right-4 bg-green-500 text-white p-4 rounded-lg shadow-lg z-50">
      {{ globalNotification }}
    </div> -->
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import logoImage from '@/assets/images/MEDITRACK_LOGO.svg'

const router = useRouter()
const authStore = useAuthStore()
const searchQuery = ref('')
const mobileMenuOpen = ref(false)
const managementMenuOpen = ref(false)
const inventoryMenuOpen = ref(false)

// Inicializar autenticacion al montar el componente
onMounted(() => {
  authStore.initializeAuth()
})

const handleSearch = () => {
  // TODO: Implementar logica de busqueda
  console.log('Buscando:', searchQuery.value)
}

const handleLogout = () => {
  authStore.logout()
  router.push('/login')
}
// Estado para notificaciones (futuro)
// const globalNotification = ref(null)
// const pendingRequests = ref(0)

// Cerrar menu movil cuando se cambia de ruta
router.afterEach(() => {
  mobileMenuOpen.value = false
  managementMenuOpen.value = false
  inventoryMenuOpen.value = false
})

// Metodos para futuras funcionalidades
// const showNotification = (message) => {
//   globalNotification.value = message
//   setTimeout(() => {
//     globalNotification.value = null
//   }, 5000)
// }

// const loadPendingRequestsCount = async () => {
//   // Cargar numero de solicitudes pendientes
// }

// Lifecycle hooks para futuras funcionalidades
// onMounted(() => {
//   loadPendingRequestsCount()
// })
</script>

<style scoped>
/* Estilos especificos del componente */
.router-link-active {
  color: #0C70CD; /* text-brand-blue-dark */
}

/* Transiciones suaves para la navegacion */
.router-link-active svg {
  transform: scale(1.05);
  transition: transform 0.2s;
}

/* Hover effects mejorados */
.transition-colors {
  transition: color 0.2s ease-in-out, background-color 0.2s ease-in-out;
}

/* Rotación del icono del dropdown */
.rotate-180 {
  transform: rotate(180deg);
  transition: transform 0.3s ease-in-out;
}

/* Efecto de resaltado para navegacion activa */
.bg-brand-blue-medium {
  background-color: #2B88D0;
}

.bg-brand-blue-dark {
  background-color: #0C70CD;
}

/* Sombras para el header */
.shadow-sm {
  box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.05);
}

/* Z-index para elementos superpuestos */
.z-40 {
  z-index: 40;
}

.z-50 {
  z-index: 50;
}

/* Estilos responsivos mejorados */
@media (max-width: 640px) {
  .text-xl {
    font-size: 1.125rem;
  }
}

@media (min-width: 640px) {
  .sm\:text-2xl {
    font-size: 1.5rem;
  }
}

@media (min-width: 768px) {
  .md\:hidden {
    display: none;
  }
  
  .md\:flex {
    display: flex;
  }
}

@media (min-width: 1024px) {
  .lg\:hidden {
    display: none;
  }
  
  .lg\:flex {
    display: flex;
  }
}

@media (min-width: 1280px) {
  .xl\:space-x-8 > :not([hidden]) ~ :not([hidden]) {
    margin-left: 2rem;
  }
  
  .xl\:px-3 {
    padding-left: 0.75rem;
    padding-right: 0.75rem;
  }
}

/* Grid para navegacion inferior */
.grid-cols-5 {
  grid-template-columns: repeat(5, minmax(0, 1fr));
}

.gap-1 {
  gap: 0.25rem;
}

/* Espaciado consistente */
.space-x-2 > :not([hidden]) ~ :not([hidden]) {
  margin-left: 0.5rem;
}

.space-x-4 > :not([hidden]) ~ :not([hidden]) {
  margin-left: 1rem;
}

.space-y-1 > :not([hidden]) ~ :not([hidden]) {
  margin-top: 0.25rem;
}

/* Texto que no se corta */
.whitespace-nowrap {
  white-space: nowrap;
}

/* Flex shrink */
.flex-shrink-0 {
  flex-shrink: 0;
}

/* Leading tight para texto compacto */
.leading-tight {
  line-height: 1.25;
}

/* Animacion suave para el menu hamburguesa */
.md\:hidden svg {
  transition: transform 0.3s ease-in-out;
}

/* Estilos para elementos flexibles */
.flex-1 {
  flex: 1 1 0%;
}

.justify-center {
  justify-content: center;
}

.items-center {
  align-items: center;
}

.max-w-4xl {
  max-width: 56rem;
}

/* Responsive utilities adicionales */
@media (max-width: 767px) {
  .mobile-only {
    display: block;
  }
}

@media (min-width: 768px) and (max-width: 1023px) {
  .tablet-only {
    display: block;
  }
}

@media (min-width: 1024px) {
  .desktop-only {
    display: block;
  }
}

/* Ajustes para pantallas muy pequeñas */
@media (max-width: 375px) {
  .text-xs {
    font-size: 0.625rem;
  }
  
  .h-6 {
    height: 1.25rem;
  }
  
  .w-6 {
    width: 1.25rem;
  }
}
</style>