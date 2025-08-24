<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="bg-white rounded-lg shadow-sm border p-6">
      <h2 class="text-2xl font-semibold text-gray-900">Escáner de Insumos Médicos</h2>
      <p class="text-gray-600 mt-1">Escanea códigos QR de insumos individuales para ver su información y estado</p>
      
      <!-- Info Panel -->
      <div class="mt-4 p-4 bg-blue-50 rounded-lg border border-blue-200">
        <div class="flex items-start space-x-3">
          <svg class="h-5 w-5 text-blue-600 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <div>
            <h4 class="text-sm font-medium text-blue-800">Enfoque en Insumos Individuales</h4>
            <p class="text-sm text-blue-700 mt-1">
              Este escáner está optimizado para insumos individuales. Cada producto de un lote tiene su propio código QR único que permite trazabilidad completa.
            </p>
          </div>
        </div>
      </div>
    </div>

    <!-- Scanner Input -->
    <div class="bg-white rounded-lg shadow-sm border p-6">
      <h3 class="text-lg font-medium text-gray-900 mb-4 flex items-center">
        <svg class="h-5 w-5 mr-2 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h2M4 4h5l2 3h3l2-3h5v5M4 4v5m16-5v5" />
        </svg>
        Escanear Insumo Individual
      </h3>

      <div class="grid lg:grid-cols-2 gap-6">
        <!-- Manual Input -->
        <div>
          <label for="qr-input" class="block text-sm font-medium text-gray-700 mb-2">
            Ingrese el Código QR del Insumo:
          </label>
          <div class="flex space-x-3">
            <input
              id="qr-input"
              v-model="qrInput"
              type="text"
              placeholder="SUPPLY_1755580808_def456"
              class="form-input flex-1"
              @keyup.enter="scanQRCode"
              @paste="handlePaste"
            />
            <button
              @click="scanQRCode"
              :disabled="!qrInput.trim() || loading"
              class="btn-primary"
            >
              <svg v-if="loading" class="animate-spin h-4 w-4 mr-2" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              {{ loading ? 'Escaneando...' : 'Escanear' }}
            </button>
          </div>
          
          <!-- Format Helper - Focused on Individual Supplies -->
          <div class="mt-3 p-3 bg-green-50 rounded-lg">
            <p class="text-sm text-green-800">
              <svg class="h-4 w-4 inline mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <strong>Prioridad: Insumos Individuales</strong>
            </p>
            <ul class="text-sm text-green-700 mt-1 space-y-1">
              <li>• <code class="bg-green-100 px-1 rounded">SUPPLY_...</code> - Insumos individuales (recomendado)</li>
              <li>• <code class="bg-blue-100 px-1 rounded">BATCH_...</code> - Información del lote (solo consulta)</li>
            </ul>
          </div>
        </div>

        <!-- Camera Scanner con Efectos Mejorados -->
        <div class="space-y-4">
          <label class="block text-sm font-medium text-gray-700">
            Escanear con Cámara:
          </label>
          
          <div class="relative">
            <!-- Camera View con Múltiples Efectos -->
            <div 
              v-if="cameraActive" 
              class="bg-gray-900 rounded-lg overflow-hidden aspect-video flex items-center justify-center relative scanner-container"
            >
              <!-- Indicador de estado mejorado -->
              <div class="absolute top-4 left-1/2 transform -translate-x-1/2 z-20">
                <div class="bg-green-500 text-white px-4 py-2 rounded-full text-sm font-medium glow-effect">
                  <div class="flex items-center space-x-2">
                    <div class="w-3 h-3 bg-white rounded-full animate-ping"></div>
                    <span>{{ detecting ? 'Detectando QR de Insumo...' : 'Cámara Activa' }}</span>
                  </div>
                </div>
              </div>
              
              <!-- Video element -->
              <video 
                ref="videoElement" 
                autoplay 
                muted 
                playsinline
                class="w-full h-full object-cover transform scale-x-[-1]"
              ></video>
              
              <!-- Grid overlay para sensación tech -->
              <div class="absolute inset-0 detection-grid opacity-30 pointer-events-none"></div>
              
              <!-- Marco de detección QR con efectos avanzados -->
              <div class="absolute inset-0 flex items-center justify-center pointer-events-none">
                <div class="relative w-56 h-56">
                  <!-- Marco principal con glow -->
                  <div class="w-full h-full border-2 border-green-400 rounded-lg bg-transparent glow-effect"></div>
                  
                  <!-- Esquinas animadas -->
                  <div class="absolute -top-1 -left-1 w-8 h-8 border-t-4 border-l-4 border-green-300 corner-animation"></div>
                  <div class="absolute -top-1 -right-1 w-8 h-8 border-t-4 border-r-4 border-green-300 corner-animation corner-delay-1"></div>
                  <div class="absolute -bottom-1 -left-1 w-8 h-8 border-b-4 border-l-4 border-green-300 corner-animation corner-delay-2"></div>
                  <div class="absolute -bottom-1 -right-1 w-8 h-8 border-b-4 border-r-4 border-green-300 corner-animation corner-delay-3"></div>
                  
                  <!-- Línea de escaneo -->
                  <div class="absolute inset-0 overflow-hidden rounded-lg">
                    <div class="w-full h-1 scanning-line"></div>
                  </div>
                  
                  <!-- Efecto radar -->
                  <div class="absolute inset-0 overflow-hidden rounded-lg">
                    <div class="absolute top-1/2 left-1/2 w-full h-0.5 radar-line transform -translate-y-1/2 origin-left"></div>
                  </div>
                  
                  <!-- Partículas flotantes -->
                  <div class="absolute inset-0">
                    <div class="absolute top-4 left-4 w-2 h-2 bg-green-400 rounded-full floating-particle opacity-60"></div>
                    <div class="absolute top-8 right-6 w-1.5 h-1.5 bg-blue-400 rounded-full floating-particle particle-delay-1 opacity-50"></div>
                    <div class="absolute bottom-6 left-8 w-1 h-1 bg-yellow-400 rounded-full floating-particle particle-delay-2 opacity-40"></div>
                    <div class="absolute bottom-4 right-4 w-2 h-2 bg-purple-400 rounded-full floating-particle particle-delay-3 opacity-70"></div>
                    <div class="absolute top-1/2 left-2 w-1.5 h-1.5 bg-pink-400 rounded-full floating-particle particle-delay-4 opacity-60"></div>
                  </div>
                  
                  <!-- Crosshair central -->
                  <div class="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2">
                    <div class="w-6 h-6 border-2 border-green-300 rounded-full flex items-center justify-center animate-pulse">
                      <div class="w-2 h-2 bg-green-300 rounded-full animate-ping"></div>
                    </div>
                  </div>
                </div>
              </div>
              
              <!-- Indicadores de estado adicionales -->
              <div class="absolute top-4 right-4 pointer-events-none">
                <div class="bg-black bg-opacity-50 text-green-400 px-3 py-1 rounded text-xs font-mono">
                  SCANNER ACTIVE
                </div>
              </div>
              
              <div class="absolute bottom-4 left-4 pointer-events-none">
                <div class="bg-black bg-opacity-50 text-blue-400 px-3 py-1 rounded text-xs font-mono">
                  AI: {{ detecting ? 'DETECTANDO...' : 'ESPERA' }}
                </div>
              </div>
              
              <div class="absolute bottom-4 right-4 pointer-events-none">
                <div class="bg-black bg-opacity-50 text-yellow-400 px-3 py-1 rounded text-xs font-mono">
                  MODO QR
                </div>
              </div>
              
              <!-- Controles -->
              <div class="absolute bottom-4 left-1/2 transform -translate-x-1/2 z-20">
                <button @click="stopCamera" class="bg-red-600 hover:bg-red-700 text-white px-6 py-2 rounded-lg text-sm font-medium shadow-lg transition-all duration-200 transform hover:scale-105">
                  Detener Cámara
                </button>
              </div>
            </div>
            
            <!-- Botón de activación de cámara -->
            <button
              v-if="!cameraActive"
              @click="startCameraScanner"
              :disabled="cameraStarting"
              class="w-full h-32 border-2 border-dashed border-gray-300 hover:border-gray-400 rounded-lg flex items-center justify-center space-y-2 flex-col bg-gray-50 hover:bg-gray-100 transition-colors"
            >
              <svg class="h-10 w-10 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" />
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z" />
              </svg>
              <span class="text-sm font-medium text-gray-700">
                {{ cameraStarting ? 'Iniciando Cámara...' : 'Activar Cámara' }}
              </span>
            </button>
          </div>
          
          <!-- Error de cámara -->
          <div v-if="cameraError" class="text-sm text-red-600 bg-red-50 p-3 rounded border border-red-200">
            <div class="flex items-start space-x-2">
              <svg class="h-4 w-4 text-red-500 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <div>
                <p class="font-medium">Error de cámara:</p>
                <p>{{ cameraError }}</p>
                <button @click="startCameraScanner" class="mt-2 text-blue-600 hover:text-blue-800 underline text-sm">
                  Intentar de nuevo
                </button>
              </div>
            </div>
          </div>
          
          <!-- Instrucciones mejoradas -->
          <div class="text-sm text-gray-600 bg-gradient-to-r from-blue-50 to-green-50 p-4 rounded-lg border border-blue-200">
            <p class="font-medium mb-2 text-blue-800">Instrucciones para escanear insumos:</p>
            <ul class="space-y-1 text-xs">
              <li>• Coloque el QR frente a la cámara para escanear su código QR.</li>
              <li>• Si lo prefiere, puede ingresar el código QR manualmente en el campo de texto.</li>
              <li>• El sistema detecta automáticamente si el QR corresponde a un insumo individual o a un lote.</li>
              <li>• Para trazabilidad completa, escanee los códigos QR individuales de cada insumo.</li>
              <li>• Si tiene problemas con la cámara, revise los permisos o intente ingresar el código manualmente.</li>
            </ul>
          </div>
        </div>
      </div>
    </div>

    <!-- Error Display -->
    <div v-if="error" class="bg-red-50 border border-red-200 rounded-lg p-4">
      <div class="flex items-start space-x-3">
        <svg class="h-5 w-5 text-red-600 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <div class="flex-1">
          <h4 class="text-sm font-medium text-red-800">Error al escanear</h4>
          <p class="text-sm text-red-700 mt-1">{{ error }}</p>
          <button @click="clearError" class="text-sm text-red-600 hover:text-red-800 mt-2 underline">
            Limpiar error
          </button>
        </div>
      </div>
    </div>

    <!-- Resto del componente igual que el original... -->
    <!-- [El resto del template permanece igual] -->

    <!-- Scanned Supply Info Display -->
    <div v-if="scannedInfo && !error" class="bg-white rounded-lg shadow-sm border overflow-hidden">
      <!-- [Contenido igual al original] -->
      <!-- Individual Supply Display -->
  <div v-if="scannedInfo.type === 'medical_supply' || scannedInfo.type === 'supply'" class="divide-y divide-gray-200">
        <!-- Header -->
        <div class="p-6 bg-gradient-to-r from-green-50 to-blue-50">
          <div class="flex items-center justify-between">
            <div class="flex items-center space-x-3">
              <div class="p-2 bg-green-100 rounded-full">
                <svg class="h-6 w-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </div>
              <div>
                <h3 class="text-lg font-medium text-gray-900">Insumo Individual Encontrado</h3>
                <p class="text-sm text-gray-600">{{ scannedInfo.supply_info?.supply_code_name || 'Información del producto' }}</p>
              </div>
            </div>
            <div class="flex items-center space-x-2">
              <span :class="[
                'px-2 py-1 text-xs font-medium rounded-full',
                scannedInfo.is_consumed 
                  ? 'bg-red-100 text-red-800' 
                  : 'bg-green-100 text-green-800'
              ]">
                {{ scannedInfo.is_consumed ? 'Consumido' : 'Disponible' }}
              </span>
            </div>
          </div>
        </div>

        <!-- Supply Details -->
        <div class="p-6 grid md:grid-cols-2 gap-6">
          <!-- Basic Info -->
          <div class="space-y-4">
            <h4 class="font-medium text-gray-900 border-b pb-2">Información del Insumo</h4>
            
            <div class="grid grid-cols-1 gap-3">
              <div>
                <label class="text-xs font-medium text-gray-500 uppercase tracking-wide">Código QR</label>
                <p class="text-sm font-mono text-gray-900 bg-gray-50 px-2 py-1 rounded border">{{ scannedInfo.qr_code }}</p>
              </div>
              
              <div v-if="scannedInfo.supply_info?.supply_code_name">
                <label class="text-xs font-medium text-gray-500 uppercase tracking-wide">Nombre del Producto</label>
                <p class="text-sm font-medium text-gray-900">{{ scannedInfo.supply_info.supply_code_name }}</p>
              </div>
              
              <div v-if="scannedInfo.supply_info?.supplier">
                <label class="text-xs font-medium text-gray-500 uppercase tracking-wide">Proveedor</label>
                <p class="text-sm text-gray-900">{{ scannedInfo.supply_info.supplier }}</p>
              </div>
              
              <div v-if="scannedInfo.supply_info?.expiration_date">
                <label class="text-xs font-medium text-gray-500 uppercase tracking-wide">Fecha de Vencimiento</label>
                <p class="text-sm text-gray-900">{{ formatDate(scannedInfo.supply_info.expiration_date) }}</p>
              </div>
            </div>
          </div>

          <!-- Batch Info -->
          <div class="space-y-4">
            <h4 class="font-medium text-gray-900 border-b pb-2">Información del Lote</h4>
            
            <div class="grid grid-cols-1 gap-3">
              <div v-if="scannedInfo.batch_status?.batch_id">
                <label class="text-xs font-medium text-gray-500 uppercase tracking-wide">ID del Lote</label>
                <p class="text-sm text-gray-900">#{{ scannedInfo.batch_status.batch_id }}</p>
              </div>
              
              <div v-if="scannedInfo.batch_status?.current_amount !== undefined">
                <label class="text-xs font-medium text-gray-500 uppercase tracking-wide">Stock Actual del Lote</label>
                <p class="text-sm text-gray-900">{{ scannedInfo.batch_status.current_amount }} unidades</p>
              </div>
              
              <div v-if="scannedInfo.supply_info?.store_name">
                <label class="text-xs font-medium text-gray-500 uppercase tracking-wide">Ubicación</label>
                <p class="text-sm text-gray-900">{{ scannedInfo.supply_info.store_name }} ({{ scannedInfo.supply_info.store_type }})</p>
              </div>
            </div>
          </div>
        </div>

        <!-- Action Buttons -->
        <div class="px-6 py-4 bg-gray-50 flex flex-wrap gap-3">
          <button
            v-if="!scannedInfo.is_consumed"
            @click="consumeSupply(scannedInfo)"
            class="btn-danger flex items-center space-x-2"
          >
            <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
            </svg>
            <span>Consumir Insumo</span>
          </button>
          
          <button @click="viewDetails(scannedInfo)" class="btn-primary flex items-center space-x-2">
            <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <span>Ver Historial Completo</span>
          </button>
          
          <button v-if="scannedInfo.batch_status?.batch_id" @click="viewBatch(scannedInfo.batch_status.batch_id)" class="btn-secondary flex items-center space-x-2">
            <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
            </svg>
            <span>Ver Lote Completo</span>
          </button>
        </div>
      </div>

      <!-- Batch Display (When batch QR is scanned) -->
      <div v-else-if="scannedInfo.type === 'batch'" class="divide-y divide-gray-200">
        <div class="p-6 bg-gradient-to-r from-blue-50 to-purple-50">
          <div class="flex items-center justify-between">
            <div class="flex items-center space-x-3">
              <div class="p-2 bg-blue-100 rounded-full">
                <svg class="h-6 w-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
                </svg>
              </div>
              <div>
                <h3 class="text-lg font-medium text-gray-900">Información del Lote</h3>
                <p class="text-sm text-gray-600">QR de lote detectado - Solo consulta</p>
              </div>
            </div>
          </div>
        </div>

        <div class="p-6">
          <div class="bg-yellow-50 border border-yellow-200 rounded-lg p-4 mb-4">
            <div class="flex items-start space-x-3">
              <svg class="h-5 w-5 text-yellow-600 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.882 16.5c-.77.833.192 2.5 1.732 2.5z" />
              </svg>
              <div>
                <h4 class="text-sm font-medium text-yellow-800">Código QR de Lote</h4>
                <p class="text-sm text-yellow-700 mt-1">
                  Has escaneado un código QR de lote completo. Para trazabilidad individual, escanea el código QR específico de cada insumo del lote.
                </p>
              </div>
            </div>
          </div>

          <div class="grid md:grid-cols-2 gap-6">
            <div>
              <label class="text-xs font-medium text-gray-500 uppercase tracking-wide">ID del Lote</label>
              <p class="text-sm font-medium text-gray-900">#{{ scannedInfo.id }}</p>
            </div>
            
            <div v-if="scannedInfo.batch_status?.available_supplies !== undefined">
              <label class="text-xs font-medium text-gray-500 uppercase tracking-wide">Insumos Disponibles</label>
              <p class="text-sm text-gray-900">{{ scannedInfo.batch_status.available_supplies }} unidades</p>
            </div>
          </div>

          <div class="mt-6 flex flex-wrap gap-3">
            <button @click="viewBatch(scannedInfo.id)" class="btn-primary flex items-center space-x-2">
              <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <span>Ver Detalles del Lote</span>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Scan History -->
    <div v-if="scanHistory.length > 0" class="bg-white rounded-lg shadow-sm border p-6">
      <div class="flex items-center justify-between mb-4">
        <h3 class="text-lg font-medium text-gray-900 flex items-center">
          <svg class="h-5 w-5 mr-2 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          Historial de Escaneos
        </h3>
        <button @click="clearHistory" class="text-sm text-gray-500 hover:text-gray-700">
          Limpiar Historial
        </button>
      </div>
      
      <div class="space-y-2 max-h-60 overflow-y-auto">
        <div 
          v-for="(item, index) in scanHistory" 
          :key="index"
          class="flex items-center justify-between p-3 bg-gray-50 rounded-lg hover:bg-gray-100 cursor-pointer transition-colors"
          @click="quickRescan(item.qr_code)"
        >
          <div class="flex items-center space-x-3">
            <div :class="[
              'p-1 rounded',
              item.type === 'medical_supply' ? 'bg-green-100 text-green-600' : 'bg-blue-100 text-blue-600'
            ]">
              <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path v-if="item.type === 'medical_supply'" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
              </svg>
            </div>
            <div>
              <p class="text-sm font-medium text-gray-900">{{ item.qr_code }}</p>
              <p class="text-xs text-gray-500">{{ formatDate(item.scanned_at) }}</p>
            </div>
          </div>
          <span :class="[
            'px-2 py-1 text-xs font-medium rounded',
            item.type === 'medical_supply' ? 'bg-green-100 text-green-700' : 'bg-blue-100 text-blue-700'
          ]">
            {{ item.type === 'medical_supply' ? 'Insumo' : 'Lote' }}
          </span>
        </div>
      </div>
    </div>

    <!-- Quick Actions -->
    <div class="bg-white rounded-lg shadow-sm border p-6">
      <h3 class="text-lg font-medium text-gray-900 mb-4">Acciones Rápidas</h3>
      <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
        <router-link to="/consume" class="p-4 text-center block text-white bg-red-600 hover:bg-red-700 rounded-md shadow font-medium flex flex-col items-center justify-center transition-colors">
          <svg class="h-6 w-6 mx-auto mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
          </svg>
          <div class="text-sm font-medium">Consumir</div>
        </router-link>

        <router-link to="/inventory" class="p-4 text-center block text-gray-700 bg-white border border-gray-300 hover:bg-gray-100 rounded-md shadow font-medium flex flex-col items-center justify-center transition-colors">
          <svg class="h-6 w-6 mx-auto mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
          </svg>
          <div class="text-sm font-medium">Inventario</div>
        </router-link>

        <button @click="clearAll" class="p-4 text-center block text-gray-700 bg-white border border-gray-300 hover:bg-gray-100 rounded-md shadow font-medium flex flex-col items-center justify-center transition-colors">
          <svg class="h-6 w-6 mx-auto mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
          <div class="text-sm font-medium">Limpiar</div>
        </button>

        <router-link to="/" class="p-4 text-center block text-gray-700 bg-white border border-gray-300 hover:bg-gray-100 rounded-md shadow font-medium flex flex-col items-center justify-center transition-colors">
          <svg class="h-6 w-6 mx-auto mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
          </svg>
          <div class="text-sm font-medium">Inicio</div>
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { format } from 'date-fns'
import { es } from 'date-fns/locale'
import qrService from '@/services/qrService'
import jsQR from 'jsqr'

const router = useRouter()
const route = useRoute()

// Referencias
const videoElement = ref(null)

// Estado reactivo
const loading = ref(false)
const error = ref(null)
const qrInput = ref('')
const scannedInfo = ref(null)
const cameraActive = ref(false)
const cameraStarting = ref(false)
const cameraError = ref(null)
const scanHistory = ref([])
const detecting = ref(false)

let mediaStream = null
let animationFrameId = null

// Camera Functions - Manteniendo la lógica original pero agregando efectos
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
    
    // Mostrar toast de activación con efectos
  showDetectionToast('Cámara activada - Efectos visuales cargados', 'success')
    
    // Esperar a que el elemento video esté disponible
    await new Promise(resolve => setTimeout(resolve, 100))
    
    if (videoElement.value) {
      videoElement.value.srcObject = stream
      
      // Iniciar detección QR
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
  if (!cameraActive.value || !videoElement.value) {
    console.log('Detección cancelada: cámara no activa o elemento no disponible')
    return
  }
  
  console.log('🔍 Iniciando detección QR con efectos visuales...')
  detecting.value = true
  
  const detectQR = () => {
    if (!cameraActive.value || !videoElement.value) {
      detecting.value = false
      console.log('Detección detenida: cámara no activa')
      return
    }
    
    try {
      // Verificar que el video tiene dimensiones válidas antes de procesar
      if (videoElement.value.videoWidth === 0 || videoElement.value.videoHeight === 0) {
        // El video aún no está listo, intentar de nuevo en el siguiente frame
        animationFrameId = requestAnimationFrame(detectQR)
        return
      }
      
      // Crear canvas localmente (no como variable global)
      const canvas = document.createElement('canvas')
      const context = canvas.getContext('2d')
      const video = videoElement.value
      
      canvas.width = video.videoWidth
      canvas.height = video.videoHeight
      
      // Dibujar el frame actual del video en el canvas
      context.drawImage(video, 0, 0, canvas.width, canvas.height)
      
      // Obtener los datos de la imagen
      const imageData = context.getImageData(0, 0, canvas.width, canvas.height)
      
      // Intentar detectar QR
      const qrCode = jsQR(imageData.data, imageData.width, imageData.height, {
        inversionAttempts: "dontInvert"
      })
      
      if (qrCode && qrCode.data) {
        const qrData = qrCode.data.trim()
        console.log('QR detectado:', qrData)
        
        // Verificar que sea un código QR válido (SUPPLY_ o BATCH_)
        if (qrData.startsWith('SUPPLY_') || qrData.startsWith('BATCH_')) {
          console.log('✅ QR válido detectado:', qrData)
          
          // Detener la detección
          detecting.value = false
          
          // Establecer el valor en el input
          qrInput.value = qrData
          
          // Detener la cámara
          stopCamera()
          
          // Mostrar notificación de éxito con efectos mejorados
          showDetectionSuccess(qrData)
          
          // Auto-escanear después de un breve delay
          setTimeout(() => {
            scanQRCode()
          }, 300)
          
          return
        } else {
          console.log('QR detectado pero formato inválido:', qrData)
        }
      }
      
    } catch (canvasError) {
      console.warn('Error en procesamiento de canvas:', canvasError)
      // No detener la detección por errores de canvas
    }
    
    // Continuar detectando si llegamos hasta aquí
    if (cameraActive.value) {
      animationFrameId = requestAnimationFrame(detectQR)
    } else {
      detecting.value = false
    }
  }
  
  // Iniciar el bucle de detección
  detectQR()
}

const showDetectionSuccess = (qrData) => {
  // Crear notificación mejorada con más efectos
  const notification = document.createElement('div')
  notification.className = 'fixed top-4 right-4 z-50 transform transition-all duration-500'
  notification.innerHTML = `
    <div class="bg-gradient-to-r from-green-500 to-emerald-500 text-white px-6 py-4 rounded-lg shadow-2xl">
      <div class="flex items-center space-x-3">
        <div class="flex-shrink-0">
          <div class="w-8 h-8 bg-white bg-opacity-20 rounded-full flex items-center justify-center animate-bounce">
            <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
          </div>
        </div>
        <div class="flex-1">
          <h4 class="text-sm font-bold">QR Detectado con Éxito</h4>
          <p class="text-xs opacity-90 font-mono">${qrData}</p>
        </div>
        <div class="flex-shrink-0">
          <div class="w-2 h-2 bg-white rounded-full animate-ping"></div>
        </div>
      </div>
    </div>
  `
  
  // Animar entrada
  notification.style.transform = 'translateX(100%) scale(0.8)'
  notification.style.opacity = '0'
  
  document.body.appendChild(notification)
  
  setTimeout(() => {
    notification.style.transform = 'translateX(0) scale(1)'
    notification.style.opacity = '1'
  }, 100)
  
  // Remover notificación después de 3 segundos con animación
  setTimeout(() => {
    notification.style.transform = 'translateX(100%) scale(0.8)'
    notification.style.opacity = '0'
    setTimeout(() => {
      if (notification.parentNode) {
        notification.remove()
      }
    }, 500)
  }, 3000)
}

const showDetectionToast = (message, type = 'info') => {
  const toast = document.createElement('div')
  const colors = {
    success: 'from-green-500 to-emerald-500',
    error: 'from-red-500 to-rose-500', 
    warning: 'from-yellow-500 to-amber-500',
    info: 'from-blue-500 to-indigo-500'
  }
  
  toast.className = `fixed top-4 left-1/2 transform -translate-x-1/2 z-50 transition-all duration-300`
  toast.innerHTML = `
    <div class="bg-gradient-to-r ${colors[type]} text-white px-6 py-3 rounded-full shadow-lg">
      <div class="flex items-center space-x-2">
        <div class="w-2 h-2 bg-white rounded-full animate-pulse"></div>
        <span class="text-sm font-medium">${message}</span>
      </div>
    </div>
  `
  
  toast.style.transform = 'translateX(-50%) translateY(-100%)'
  toast.style.opacity = '0'
  
  document.body.appendChild(toast)
  
  setTimeout(() => {
    toast.style.transform = 'translateX(-50%) translateY(0)'
    toast.style.opacity = '1'
  }, 100)
  
  setTimeout(() => {
    toast.style.transform = 'translateX(-50%) translateY(-100%)'
    toast.style.opacity = '0'
    setTimeout(() => {
      if (toast.parentNode) {
        toast.remove()
      }
    }, 300)
  }, 2500)
}

const stopCamera = () => {
  console.log('🛑 Deteniendo cámara...')
  
  // Detener todos los tracks del stream
  if (mediaStream) {
    mediaStream.getTracks().forEach(track => {
      track.stop()
      console.log('Track detenido:', track.kind)
    })
    mediaStream = null
  }
  
  // Cancelar la animación
  if (animationFrameId) {
    cancelAnimationFrame(animationFrameId)
    animationFrameId = null
  }
  
  // Limpiar el elemento video
  if (videoElement.value) {
    videoElement.value.srcObject = null
  }
  
  // Actualizar estado
  cameraActive.value = false
  detecting.value = false
  cameraError.value = null
  
  //showDetectionToast('Cámara detenida', 'info')
  
  console.log('Cámara detenida correctamente')
}

// Scanner Functions - mantener igual
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
    
    console.log('Escaneando QR:', qrInput.value.trim())
    const result = await qrService.scanQRCode(qrInput.value.trim())
    scannedInfo.value = result
    
    // Añadir al historial con prioridad para insumos individuales
    addToHistory(qrInput.value.trim(), result.type)
    
    // Mostrar mensaje diferente según el tipo
    if (result.type === 'medical_supply') {
      console.log('✅ Insumo individual escaneado correctamente')
  showDetectionToast('Insumo individual procesado', 'success')
    } else if (result.type === 'batch') {
      console.log('ℹ️ Lote escaneado - Para trazabilidad individual, escanee códigos QR específicos de cada insumo')
  showDetectionToast('Lote detectado - Use QR individuales para trazabilidad', 'info')
    }
    
  } catch (err) {
    console.error('Error scanning QR:', err)
    error.value = err.response?.data?.error || err.message || 'Error al escanear el código QR'
  showDetectionToast('Error al escanear QR', 'error')
  } finally {
    loading.value = false
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

// History Management - mantener igual
const addToHistory = (qrCode, type) => {
  const existing = scanHistory.value.findIndex(item => item.qr_code === qrCode)
  
  const historyItem = {
    qr_code: qrCode,
    type: type,
    scanned_at: new Date()
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

// Navigation - mantener igual
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

// Utilities - mantener igual
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

// Lifecycle - mantener igual
onMounted(() => {
  loadHistory()
  
  // Auto-escanear si viene QR en query params
  if (route.query.qr || route.query.test) {
    const testQR = route.query.qr || route.query.test
    qrInput.value = testQR
    scanQRCode()
  }
})

onUnmounted(() => {
  console.log('Componente desmontándose, limpiando cámara...')
  stopCamera()
})
</script>

<style scoped>
.form-input {
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

/* ===== EFECTOS VISUALES MEJORADOS ===== */

/* Animaciones personalizadas */
@keyframes scanning {
  0% { 
    transform: translateY(-100%); 
    opacity: 0; 
  }
  10% { 
    opacity: 1; 
  }
  90% { 
    opacity: 1; 
  }
  100% { 
    transform: translateY(200px); 
    opacity: 0; 
  }
}

@keyframes corner-pulse {
  0%, 100% { 
    transform: scale(1);
    opacity: 0.8;
  }
  50% { 
    transform: scale(1.1);
    opacity: 1;
  }
}

@keyframes radar-sweep {
  0% { 
    transform: rotate(0deg) translateX(0); 
  }
  100% { 
    transform: rotate(360deg) translateX(0); 
  }
}

@keyframes float-particle {
  0%, 100% { 
    transform: translateY(0px) translateX(0px);
    opacity: 0.3;
  }
  50% { 
    transform: translateY(-20px) translateX(10px);
    opacity: 0.7;
  }
}

@keyframes glow {
  0%, 100% { 
    box-shadow: 0 0 10px rgba(34, 197, 94, 0.5);
  }
  50% { 
    box-shadow: 0 0 25px rgba(34, 197, 94, 0.8), 0 0 40px rgba(34, 197, 94, 0.3);
  }
}

/* Clases de efectos */
.scanning-line {
  animation: scanning 2s infinite;
  background: linear-gradient(to bottom, 
    transparent 0%, 
    rgba(34, 197, 94, 0.8) 45%, 
    rgba(34, 197, 94, 1) 50%, 
    rgba(34, 197, 94, 0.8) 55%, 
    transparent 100%);
}

.corner-animation {
  animation: corner-pulse 1.5s infinite;
}

.corner-delay-1 {
  animation-delay: 0.2s;
}

.corner-delay-2 {
  animation-delay: 0.4s;
}

.corner-delay-3 {
  animation-delay: 0.6s;
}

.radar-line {
  animation: radar-sweep 3s linear infinite;
  background: linear-gradient(90deg, 
    transparent 0%, 
    rgba(34, 197, 94, 0.1) 30%, 
    rgba(34, 197, 94, 0.6) 50%, 
    rgba(34, 197, 94, 0.1) 70%, 
    transparent 100%);
}

.floating-particle {
  animation: float-particle 3s infinite ease-in-out;
}

.particle-delay-1 {
  animation-delay: 1s;
}

.particle-delay-2 {
  animation-delay: 2s;
}

.particle-delay-3 {
  animation-delay: 1.5s;
}

.particle-delay-4 {
  animation-delay: 0.5s;
}

.glow-effect {
  animation: glow 2s infinite alternate;
}

.detection-grid {
  background-image: 
    linear-gradient(rgba(34, 197, 94, 0.1) 1px, transparent 1px),
    linear-gradient(90deg, rgba(34, 197, 94, 0.1) 1px, transparent 1px);
  background-size: 20px 20px;
  animation: float-particle 4s infinite ease-in-out;
}

.scanner-container {
  position: relative;
  overflow: hidden;
}

.scanner-container::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(
    90deg,
    transparent,
    rgba(34, 197, 94, 0.1),
    transparent
  );
  animation: sweep 2s infinite;
  z-index: 1;
  pointer-events: none;
}

@keyframes sweep {
  0% {
    left: -100%;
  }
  100% {
    left: 100%;
  }
}
</style>