<template>
  <div class="max-w-4xl mx-auto p-3 sm:p-6 bg-white rounded-lg shadow-lg">
    <!-- Título principal -->
    <div class="mb-6">
      <h2 class="text-2xl font-bold text-gray-900 mb-2">
        {{ props.editMode ? 'Editar Solicitud de Insumo' : 'Nueva Solicitud de Insumo' }}
      </h2>
      <p class="text-gray-600">
        {{ props.editMode ? 'Modificar y reenviar solicitud devuelta' : 'Crear solicitud con trazabilidad completa' }}
      </p>
    </div>

    <!-- Mostrar errores generales -->
    <div v-if="errors.length > 0" class="mb-4 sm:mb-6 p-3 sm:p-4 bg-red-50 border border-red-200 rounded-md">
      <div class="flex">
        <svg class="h-5 w-5 text-red-400 mr-2 mt-0.5" fill="currentColor" viewBox="0 0 20 20">
          <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
        </svg>
        <div>
          <h3 class="text-sm font-medium text-red-800">Se encontraron los siguientes errores:</h3>
          <ul class="mt-2 text-sm text-red-700 list-disc list-inside">
            <li v-for="error in errors" :key="error">{{ error }}</li>
          </ul>
        </div>
      </div>
    </div>

    <form @submit.prevent="submitRequest" class="space-y-4 sm:space-y-6">
      <!-- Información básica -->
      <div class="bg-gray-50 p-3 sm:p-4 rounded-lg">
        <h3 class="text-base sm:text-lg font-semibold text-gray-900 mb-3 sm:mb-4">Información Básica</h3>
        
        <!-- Información del solicitante (automática) -->
        <div class="mb-4">
          <label class="block text-sm font-medium text-gray-700 mb-1">
            Solicitante
          </label>
          <div class="w-full px-3 py-2 bg-gray-50 border border-gray-300 rounded-md text-gray-700">
            <div class="font-medium">{{ authStore.getUserName || 'Usuario' }}</div>
            <div class="text-sm text-gray-500">{{ authStore.getUserRut || 'RUT no disponible' }}</div>
            <div v-if="authStore.getUserSpecialty" class="font-medium">
              <span class="inline-flex items-center">
                {{ authStore.getUserSpecialty }}
              </span>
            </div>
          </div>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <!-- Fecha y hora de cirugía -->
          <div>
            <label for="surgery_datetime" class="block text-sm font-medium text-gray-700 mb-1">
              Fecha y Hora de Cirugía <span class="text-red-500">*</span>
            </label>
            <!-- Modo edición: mostrar como solo lectura -->
            <div v-if="props.editMode" class="w-full px-3 py-2 bg-gray-200 border border-gray-400 rounded-md text-gray-600 font-medium cursor-not-allowed">
              {{ originalRequestData.surgery_datetime_display }}
            </div>
            <!-- Modo creación: campo editable -->
            <input
              v-else
              type="datetime-local"
              id="surgery_datetime"
              v-model="requestForm.surgery_datetime"
              required
              :min="minDateTime"
              class="form-input"
              :class="{ 'border-orange-500': isUrgentRequest, 'border-red-500': isEmergencyRequest }"
              @change="checkAdvanceNotice"
            />
            <p class="text-xs text-gray-500 mt-1">
              {{ props.editMode ? 'Fecha y hora programada para la cirugía' : 'Selecciona la fecha y hora programada para la cirugía' }}
            </p>
            <!-- Advertencia de anticipación mínima -->
            <div v-if="!props.editMode && advanceNoticeWarning" class="mt-2 p-2 rounded-md" 
                 :class="advanceNoticeWarning.type === 'emergency' ? 'bg-red-50 border border-red-200' : 
                         advanceNoticeWarning.type === 'urgent' ? 'bg-orange-50 border border-orange-200' : 
                         'bg-yellow-50 border border-yellow-200'">
              <div class="flex items-start gap-2">
                <svg class="h-5 w-5 mt-0.5 flex-shrink-0" 
                     :class="advanceNoticeWarning.type === 'emergency' ? 'text-red-600' : 
                             advanceNoticeWarning.type === 'urgent' ? 'text-orange-600' : 
                             'text-yellow-600'"
                     fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                </svg>
                <div class="flex-1">
                  <p class="text-sm font-medium" 
                     :class="advanceNoticeWarning.type === 'emergency' ? 'text-red-800' : 
                             advanceNoticeWarning.type === 'urgent' ? 'text-orange-800' : 
                             'text-yellow-800'">
                    {{ advanceNoticeWarning.title }}
                  </p>
                  <p class="text-xs mt-1" 
                     :class="advanceNoticeWarning.type === 'emergency' ? 'text-red-700' : 
                             advanceNoticeWarning.type === 'urgent' ? 'text-orange-700' : 
                             'text-yellow-700'">
                    {{ advanceNoticeWarning.message }}
                  </p>
                </div>
              </div>
            </div>
          </div>

          <!-- Selección de Pabellón -->
          <div>
            <label for="pavilion" class="block text-xs sm:text-sm font-medium text-gray-700 mb-1">
              Pabellón <span class="text-red-500">*</span>
            </label>
            <!-- Modo edición: mostrar como solo lectura -->
            <div v-if="props.editMode" class="w-full px-3 py-2 bg-gray-200 border border-gray-400 rounded-md text-gray-600 font-medium cursor-not-allowed">
              {{ originalRequestData.pavilion_name }}
            </div>
            <!-- Modo creación: campo editable -->
            <select
              v-else
              id="pavilion"
              v-model="requestForm.pavilion_id"
              required
              class="form-select text-sm"
              :disabled="loadingPavilions"
            >
              <option value="">Seleccionar pabellón</option>
              <option
                v-for="pavilion in pavilions"
                :key="pavilion.id"
                :value="pavilion.id"
              >
                {{ pavilion.name }}
              </option>
            </select>
            <p v-if="loadingPavilions && !props.editMode" class="text-xs text-gray-500 mt-1">Cargando pabellones...</p>
            <p v-else class="text-xs text-gray-500 mt-1">
              {{ props.editMode ? 'Pabellón donde se realizará la cirugía' : 'Selecciona el pabellón donde se realizará la cirugía' }}
            </p>
          </div>

          <!-- Tipo de Cirugía -->
          <div>
            <label for="surgery_id" class="block text-sm font-medium text-gray-700 mb-1">
              Tipo de Cirugía <span class="text-red-500">*</span>
            </label>
            <select
              id="surgery_id"
              v-model="requestForm.surgery_id"
              required
              class="form-select text-sm"
              :disabled="loadingSurgeries"
            >
              <option :value="null">Seleccionar cirugía</option>
              <option v-for="surgery in filteredSurgeries" :key="surgery.id" :value="surgery.id">
                {{ surgery.name }}
              </option>
            </select>
            <p v-if="loadingSurgeries" class="text-xs text-gray-500 mt-1">Cargando cirugías...</p>
            <p v-else-if="filteredSurgeries.length === 0" class="text-xs text-yellow-600 mt-1">
              No hay cirugías disponibles para tu especialidad
            </p>
            <p v-else class="text-xs text-gray-500 mt-1">
              Cirugías de {{ authStore.getUserSpecialty || 'tu especialidad' }}
            </p>
          </div>
        </div>

        <!-- Observaciones generales -->
        <div class="mt-3 sm:mt-4">
          <label for="notes" class="block text-xs sm:text-sm font-medium text-gray-700 mb-1">
            Observaciones Generales
          </label>
          <textarea
            id="notes"
            v-model="requestForm.notes"
            rows="3"
            :placeholder="props.editMode ? 'Observaciones sobre los cambios realizados en esta solicitud...' : 'Observaciones adicionales sobre la solicitud...'"
                  class="form-input text-sm"
          ></textarea>
        </div>
      </div>

      <!-- Insumos Típicos de la Cirugía -->
      <div v-if="!props.editMode && showTypicalSupplies" class="bg-green-50 p-3 sm:p-4 rounded-lg border border-green-200">
        <div class="flex flex-col sm:flex-row sm:justify-between sm:items-center gap-3 mb-3">
          <div class="flex items-center gap-2 flex-1">
            <button
              type="button"
              @click="typicalSuppliesExpanded = !typicalSuppliesExpanded"
              class="text-gray-600 hover:text-gray-900 transition-colors"
              title="Minimizar/Maximizar"
            >
              <svg 
                v-if="typicalSuppliesExpanded" 
                class="h-5 w-5" 
                fill="none" 
                stroke="currentColor" 
                viewBox="0 0 24 24"
              >
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
              </svg>
              <svg 
                v-else 
                class="h-5 w-5" 
                fill="none" 
                stroke="currentColor" 
                viewBox="0 0 24 24"
              >
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
              </svg>
            </button>
            <div>
              <h3 class="text-base sm:text-lg font-semibold text-gray-900">Insumos Típicos de esta Cirugía</h3>
              <p v-if="typicalSuppliesExpanded" class="text-xs text-gray-600 mt-1">
                Estos son los insumos típicamente utilizados para esta cirugía. Puedes agregarlos todos o seleccionar los que necesites.
              </p>
            </div>
          </div>
          <button
            type="button"
            @click="addAllTypicalSupplies"
            class="btn-primary text-sm whitespace-nowrap"
            :disabled="loadingTypicalSupplies || typicalSupplies.length === 0"
          >
            <svg class="h-4 w-4 mr-1 inline" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
            </svg>
            Agregar Todos
          </button>
        </div>

        <div v-if="typicalSuppliesExpanded">
          <div v-if="loadingTypicalSupplies" class="text-center py-4">
            <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-green-600 mx-auto"></div>
            <p class="text-sm text-gray-600 mt-2">Cargando insumos típicos...</p>
          </div>

          <div v-else-if="typicalSupplies.length === 0" class="text-center py-4 text-gray-500">
            <p class="text-sm">No hay insumos típicos configurados para esta cirugía</p>
          </div>

          <div v-else class="grid grid-cols-1 md:grid-cols-2 gap-3">
            <div
              v-for="typicalSupply in typicalSupplies"
              :key="typicalSupply.id"
              class="bg-white p-3 rounded-lg border border-green-300 hover:border-green-400 transition-colors"
            >
              <div class="flex justify-between items-start gap-2">
                <div class="flex-1">
                  <div class="flex items-center gap-2 mb-1">
                    <span class="font-medium text-gray-900">{{ getSupplyNameForTypical(typicalSupply.supply_code) }}</span>
                    <span v-if="typicalSupply.is_required" class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-red-100 text-red-800">
                      Requerido
                    </span>
                  </div>
                  <div class="text-sm text-gray-600 mb-2">
                    <span class="font-medium">Código:</span> {{ typicalSupply.supply_code }}
                  </div>
                  <div class="text-sm text-gray-700 mb-2">
                    <span class="font-medium">Cantidad típica:</span> {{ typicalSupply.typical_quantity || 1 }}
                  </div>
                  <div v-if="typicalSupply.notes" class="text-xs text-gray-500 italic">
                    {{ typicalSupply.notes }}
                  </div>
                </div>
                <button
                  type="button"
                  @click="addTypicalSupply(typicalSupply)"
                  class="flex-shrink-0 px-3 py-1.5 bg-green-600 text-white text-xs font-medium rounded-md hover:bg-green-700 transition-colors"
                  title="Agregar este insumo"
                >
                  <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
                  </svg>
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Insumos solicitados -->
      <div class="bg-blue-50 p-3 sm:p-4 rounded-lg">
        <div class="flex flex-col sm:flex-row sm:justify-between sm:items-center gap-3 sm:gap-0 mb-3 sm:mb-4">
          <h3 class="text-base sm:text-lg font-semibold text-gray-900">Insumos Necesarios</h3>
          <button
            type="button"
            @click="addSupplyItem"
            class="btn-primary w-full sm:w-auto"
          >
            <svg class="h-4 w-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
            </svg>
            Agregar Insumo
          </button>
        </div>

        <!-- Lista de insumos -->
        <div v-if="requestForm.items.length === 0" class="text-center py-6 sm:py-8 text-gray-500">
          <svg class="h-12 w-12 mx-auto mb-4 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
          </svg>
          <p>No hay insumos agregados</p>
          <p class="text-sm">Haz clic en "Agregar Insumo" para comenzar</p>
        </div>

        <div class="space-y-3 sm:space-y-4">
          <div
            v-for="(item, index) in requestForm.items"
            :key="index"
            class="bg-white p-3 sm:p-4 rounded-lg border relative"
            :class="{
              'border-orange-400 bg-orange-50': isDuplicateItem(index),
              'border-green-300 bg-green-50': props.editMode && item.item_status === 'aceptado',
              'border-gray-200': !isDuplicateItem(index) && !(props.editMode && item.item_status === 'aceptado')
            }"
          >
            <!-- Advertencia de duplicado -->
            <div v-if="isDuplicateItem(index)" class="mb-3 p-2 bg-orange-100 border border-orange-300 rounded-md">
              <div class="flex items-center gap-2 text-orange-800 text-sm">
                <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                </svg>
                <span class="font-medium">Este insumo está duplicado. Las cantidades se consolidarán automáticamente al enviar.</span>
              </div>
            </div>
            <!-- Botón eliminar (solo para items no aceptados) -->
            <button
              v-if="!props.editMode || item.item_status !== 'aceptado'"
              type="button"
              @click="removeSupplyItem(index)"
              class="absolute top-2 right-2 p-1.5 text-red-500 hover:text-red-700 hover:bg-red-50 rounded-full z-10"
              title="Eliminar insumo"
            >
              <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
            
            <!-- Indicador de item aceptado (solo lectura) -->
            <div v-if="props.editMode && item.item_status === 'aceptado'" class="absolute top-2 right-2 px-2 py-1 bg-green-100 text-green-800 text-xs font-semibold rounded z-10 flex items-center gap-1">
              <svg class="h-3 w-3" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
              </svg>
              Aceptado (solo lectura)
            </div>
            
            <!-- Mensaje informativo para items aceptados -->
            <div v-if="props.editMode && item.item_status === 'aceptado'" class="mb-3 p-2 bg-green-50 border border-green-200 rounded-md">
              <div class="flex items-center gap-2 text-green-800 text-xs">
                <svg class="h-4 w-4 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd" />
                </svg>
                <span>Este insumo ya fue aceptado por bodega y no puede ser modificado ni eliminado.</span>
              </div>
            </div>

            <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3 sm:gap-4 pr-8">
              <!-- Código del insumo -->
              <div class="relative">
                <label class="block text-xs sm:text-sm font-medium text-gray-700 mb-1">
                  Código Insumo <span class="text-red-500">*</span>
                </label>
                <input
                  type="text"
                  v-model="item.supply_code"
                  @input="onSupplyCodeChange(index, $event.target.value)"
                  @focus="onSupplyCodeFocus(index)"
                  @blur="onSupplyCodeBlur(index)"
                  placeholder="Escribir código..."
                  class="form-select text-sm"
                  :disabled="props.editMode && item.item_status === 'aceptado'"
                  :readonly="props.editMode && item.item_status === 'aceptado'"
                />
                <!-- Dropdown de sugerencias por código -->
                <div 
                  v-if="showCodeDropdowns[index] && medicalSupplies.length > 0"
                  class="absolute z-10 w-full mt-1 bg-white border border-gray-300 rounded-md shadow-lg max-h-60 overflow-y-auto"
                >
                  <div 
                    v-for="supply in getFilteredSuppliesByCode(index)" 
                    :key="supply.id"
                    @click="selectSupply(index, supply)"
                    class="px-3 py-2 hover:bg-blue-50 cursor-pointer border-b border-gray-100 last:border-b-0"
                  >
                    <div class="font-medium text-gray-900">Código: {{ supply.code }}</div>
                    <div class="text-sm text-gray-500">{{ supply.name }}</div>
                  </div>
                  <div v-if="getFilteredSuppliesByCode(index).length === 0" class="px-3 py-2 text-gray-500 text-center">
                    No se encontraron insumos
                  </div>
                </div>
              </div>

              <!-- Nombre del insumo con autocompletado -->
              <div class="relative">
                <label class="block text-sm font-medium text-gray-700 mb-1">
                  Nombre Insumo <span class="text-red-500">*</span>
                </label>
                <input
                  type="text"
                  :value="supplySearchTerms[index] || item.supply_name"
                  @input="onSupplyInputChange(index, $event.target.value)"
                  @focus="onSupplyInputFocus(index)"
                  @blur="onSupplyInputBlur(index)"
                  placeholder="Escribir nombre..."
                  class="form-select text-sm"
                  :class="{ 'bg-gray-100 cursor-not-allowed': props.editMode && item.item_status === 'aceptado' }"
                  autocomplete="off"
                  :disabled="props.editMode && item.item_status === 'aceptado'"
                  :readonly="props.editMode && item.item_status === 'aceptado'"
                />
                
                <!-- Dropdown de sugerencias por nombre -->
                <div 
                  v-if="showSupplyDropdowns[index] && medicalSupplies.length > 0"
                  class="absolute z-10 w-full mt-1 bg-white border border-gray-300 rounded-md shadow-lg max-h-60 overflow-y-auto"
                >
                  <div 
                    v-for="supply in getFilteredSupplies(index)" 
                    :key="supply.id"
                    @click="selectSupply(index, supply)"
                    class="px-3 py-2 hover:bg-blue-50 cursor-pointer border-b border-gray-100 last:border-b-0"
                  >
                    <div class="font-medium text-gray-900">{{ supply.name }}</div>
                    <div class="text-sm text-gray-500">Código: {{ supply.code }}</div>
                  </div>
                  
                  <div v-if="getFilteredSupplies(index).length === 0" class="px-3 py-2 text-gray-500 text-center">
                    No se encontraron insumos
                  </div>
                </div>
                
                <!-- Indicador de carga -->
                <div v-if="loadingSupplies" class="absolute right-3 top-9 text-gray-400">
                  <svg class="animate-spin h-4 w-4" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                </div>
              </div>

              <!-- Cantidad solicitada -->
              <div>
                <label class="block text-xs sm:text-sm font-medium text-gray-700 mb-1">
                  Cantidad <span class="text-red-500">*</span>
                </label>
                <input
                  type="number"
                  v-model.number="item.quantity_requested"
                  required
                  min="1"
                  placeholder="1"
                  class="form-input text-sm"
                  :disabled="props.editMode && item.item_status === 'aceptado'"
                  :readonly="props.editMode && item.item_status === 'aceptado'"
                />
              </div>
            </div>

            <!-- Observaciones del item devuelto (solo en modo edición cuando item fue devuelto) -->
            <div v-if="props.editMode && item.item_status === 'devuelto'" class="mt-4 space-y-3">
              <!-- Mostrar observaciones anteriores del encargado -->
              <div v-if="item.item_notes" class="bg-orange-50 border border-orange-200 rounded-md p-3">
                <label class="block text-xs font-semibold text-orange-800 mb-1">
                  Motivo de devolución del encargado:
                </label>
                <p class="text-xs text-orange-900 whitespace-pre-wrap">{{ item.item_notes }}</p>
              </div>
              
              <!-- Campo para nuevas observaciones del doctor -->
              <div>
                <label class="block text-xs sm:text-sm font-medium text-gray-700 mb-1">
                  Observaciones Insumo
                </label>
                <textarea
                  v-model="item.resubmit_notes"
                  rows="2"
                  placeholder="Agregue sus observaciones..."
                  class="w-full px-3 py-2 text-sm border border-orange-300 bg-white rounded-md focus:outline-none focus:ring-2 focus:ring-orange-500 focus:border-orange-500"
                ></textarea>
                <p class="mt-1 text-xs text-orange-600">
                  Opcional: Agregue observaciones si modificó la cantidad o desea aclarar algo sobre este insumo.
                </p>
              </div>
            </div>

            <!-- Especificaciones técnicas -->
            <div class="mt-4 grid grid-cols-1 md:grid-cols-2 gap-4">
              <!-- Medidas/Tamaño 
              <div>
                <label class="block text-xs sm:text-sm font-medium text-gray-700 mb-1">
                  Medidas/Tamaño
                </label>
                <input
                  type="text"
                  v-model="item.size"
                  placeholder="Ej: Grande, Mediano, 20cm, etc."
                  class="form-select text-sm"
                />
              </div>-->

              <!-- Marca 
              <div>
                <label class="block text-xs sm:text-sm font-medium text-gray-700 mb-1">
                  Marca
                </label>
                <input
                  type="text"
                  v-model="item.brand"
                  placeholder="Marca preferida"
                  class="form-select text-sm"
                />
              </div>
            </div>-->

            <!-- Características especiales
            <div class="mt-4 grid grid-cols-1 md:grid-cols-2 gap-4">-->
              <!-- Es pediátrico -->
              <div class="flex items-center">
                <input
                  type="checkbox"
                  :id="`pediatric-${index}`"
                  v-model="item.is_pediatric"
                  class="h-4 w-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
                  :disabled="props.editMode && item.item_status === 'aceptado'"
                />
                <label :for="`pediatric-${index}`" class="ml-2 text-xs sm:text-sm font-medium text-gray-700"
                  :class="{ 'text-gray-400': props.editMode && item.item_status === 'aceptado' }">
                  Es insumo pediátrico
                </label>
              </div>
            </div>

            <!-- Especificaciones y observaciones del insumo 
            <div class="mt-4">
              <label class="block text-sm font-medium text-gray-700 mb-1">
                Especificaciones Técnicas
              </label>
              <textarea
                v-model="item.specifications"
                rows="2"
                placeholder="Especificaciones técnicas del insumo (material, dimensiones exactas, características especiales, etc.)"
                class="form-select text-sm"
              ></textarea>
            </div>-->

            <!-- Solicitudes especiales
            <div class="mt-4">
              <label class="block text-sm font-medium text-gray-700 mb-1">
                Solicitudes Especiales
              </label>
              <textarea
                v-model="item.special_requests"
                rows="2"
                placeholder="Solicitudes especiales para este insumo (entrega urgente, manipulación especial, etc.)"
                class="form-select text-sm"
              ></textarea>
            </div>-->
          </div>
        </div>
      </div>

      <!-- Modal de lista de insumos (fuera del loop para evitar múltiples modales) -->
      <Teleport to="body">
        <div 
          v-if="currentSupplyListIndex !== null"
          class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
          @click.self="closeSupplyList(currentSupplyListIndex)"
        >
          <div class="bg-white rounded-lg shadow-xl p-6 max-w-3xl w-full mx-4 max-h-[90vh] overflow-hidden flex flex-col">
            <div class="flex justify-between items-center mb-4">
              <h3 class="text-xl font-bold text-gray-800">Seleccionar Insumo</h3>
              <button
                @click="closeSupplyList(currentSupplyListIndex)"
                class="text-gray-400 hover:text-gray-600"
              >
                <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
            
            <!-- Buscador en modal -->
            <div class="mb-4">
              <input
                type="text"
                v-model="supplyListSearchTerms[currentSupplyListIndex]"
                placeholder="Buscar por código o nombre..."
                class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              />
            </div>
            
            <!-- Lista de insumos -->
            <div class="flex-1 overflow-y-auto border border-gray-200 rounded-lg">
              <div 
                v-for="supply in getFilteredSuppliesForList(currentSupplyListIndex)" 
                :key="supply.id"
                @click="selectSupply(currentSupplyListIndex, supply)"
                class="px-4 py-3 hover:bg-blue-50 cursor-pointer border-b border-gray-100 last:border-b-0 transition-colors"
              >
                <div class="flex justify-between items-start">
                  <div class="flex-1">
                    <div class="font-medium text-gray-900">{{ supply.name }}</div>
                    <div class="text-sm text-gray-500 mt-1">Código: {{ supply.code }}</div>
                  </div>
                  <svg class="h-5 w-5 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
                  </svg>
                </div>
              </div>
              <div v-if="getFilteredSuppliesForList(currentSupplyListIndex).length === 0" class="px-4 py-8 text-center text-gray-500">
                No se encontraron insumos
              </div>
            </div>
            
            <div class="mt-4 flex justify-end">
              <button
                @click="closeSupplyList(currentSupplyListIndex)"
                class="px-4 py-2 bg-gray-300 text-gray-700 rounded-lg hover:bg-gray-400 transition-colors"
              >
                Cerrar
              </button>
            </div>
          </div>
        </div>
      </Teleport>

      <!-- Botones de acción -->
      <div class="flex flex-col sm:flex-row sm:justify-between gap-3 sm:gap-0 pt-3 sm:pt-4 border-t border-gray-200">
        <button
          type="button"
          @click="resetForm"
          class="inline-flex items-center justify-center px-4 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-blue-500 w-full sm:w-auto order-3 sm:order-1"
          :disabled="submitting"
        >
          <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
          Limpiar Formulario
        </button>
        
        <div class="flex flex-col sm:flex-row gap-2 sm:gap-3 sm:space-x-0 order-1 sm:order-2">
          <button
            type="button"
            @click="cancelForm"
            class="inline-flex items-center justify-center px-4 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-gray-500 w-full sm:w-auto order-2"
            :disabled="submitting"
          >
            Cancelar
          </button>
          
          <button
            type="submit"
            class="btn-primary disabled:opacity-50 disabled:cursor-not-allowed w-full sm:w-auto order-1"
            :disabled="submitting"
          >
            <svg v-if="submitting" class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            <svg v-if="!submitting" class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
            </svg>
            <span v-if="props.editMode">
              {{ submitting ? 'Reenviando Solicitud...' : 'Reenviar Solicitud' }}
            </span>
            <span v-else>
              {{ submitting ? 'Creando Solicitud...' : 'Crear Solicitud' }}
            </span>
          </button>
        </div>
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, nextTick, watch } from 'vue'
import { useAuthStore } from '@/stores/auth'
import supplyRequestService from '@/services/requests/supplyRequestService'
import pavilionService from '@/services/config/pavilionService'
import inventoryService from '@/services/inventory/inventoryService'
import surgeryService from '@/services/management/surgeryService'
import surgeryTypicalSupplyService from '@/services/management/surgeryTypicalSupplyService'
import { useNotification } from '@/composables/useNotification'
import Swal from 'sweetalert2'

// Props
const props = defineProps({
  id: {
    type: Number,
    default: null
  },
  editMode: {
    type: Boolean,
    default: false
  }
})

// Emits
const emit = defineEmits(['success', 'cancel', 'error'])

// Stores
const authStore = useAuthStore()
const { success: showSuccess, error: showError, info: showInfo, warning: showWarning } = useNotification()

// Estado reactivo
const submitting = ref(false)
const errors = ref([])
const pavilions = ref([])
const loadingPavilions = ref(false)
const medicalSupplies = ref([])
const loadingSupplies = ref(false)
const supplySearchTerms = ref([])
const showSupplyDropdowns = ref([])
const showCodeDropdowns = ref([])
const showSupplyListModals = ref([])
const supplyListSearchTerms = ref([])
const currentSupplyListIndex = ref(null)
const surgeries = ref([])
const loadingSurgeries = ref(false)
const typicalSupplies = ref([])
const loadingTypicalSupplies = ref(false)
const showTypicalSupplies = ref(false)
const typicalSuppliesExpanded = ref(true)

// Fecha mínima para la cirugía (hoy)
const minDateTime = computed(() => {
  const now = new Date()
  const year = now.getFullYear()
  const month = String(now.getMonth() + 1).padStart(2, '0')
  const day = String(now.getDate()).padStart(2, '0')
  const hours = String(now.getHours()).padStart(2, '0')
  const minutes = String(now.getMinutes()).padStart(2, '0')
  return `${year}-${month}-${day}T${hours}:${minutes}`
})

// Estado reactivo para advertencias de anticipación
const advanceNoticeWarning = ref(null)

// Constantes de anticipación
const MINIMUM_ADVANCE_DAYS = 3
const URGENT_HOURS = 48
const EMERGENCY_HOURS = 12

// Calcular días hasta la cirugía
const daysUntilSurgery = computed(() => {
  if (!requestForm.surgery_datetime) return null
  const surgeryDate = new Date(requestForm.surgery_datetime)
  const now = new Date()
  const diffTime = surgeryDate - now
  const diffDays = diffTime / (1000 * 60 * 60 * 24)
  return diffDays
})

// Calcular horas hasta la cirugía
const hoursUntilSurgery = computed(() => {
  if (!requestForm.surgery_datetime) return null
  const surgeryDate = new Date(requestForm.surgery_datetime)
  const now = new Date()
  const diffTime = surgeryDate - now
  const diffHours = diffTime / (1000 * 60 * 60)
  return diffHours
})

// Verificar si es urgente (menos de 48 horas)
const isUrgentRequest = computed(() => {
  const hours = hoursUntilSurgery.value
  return hours !== null && hours > 0 && hours <= URGENT_HOURS
})

// Verificar si es emergencia (menos de 12 horas)
const isEmergencyRequest = computed(() => {
  const hours = hoursUntilSurgery.value
  return hours !== null && hours > 0 && hours <= EMERGENCY_HOURS
})

// Verificar si no tiene anticipación mínima (menos de 3 días)
const isNotProgrammed = computed(() => {
  const days = daysUntilSurgery.value
  return days !== null && days > 0 && days < MINIMUM_ADVANCE_DAYS
})

// Función para verificar anticipación y mostrar advertencia
const checkAdvanceNotice = () => {
  if (!requestForm.surgery_datetime) {
    advanceNoticeWarning.value = null
    return
  }

  const hours = hoursUntilSurgery.value
  const days = daysUntilSurgery.value

  if (hours === null || hours < 0) {
    advanceNoticeWarning.value = {
      type: 'error',
      title: 'Fecha inválida',
      message: 'La fecha de cirugía no puede ser en el pasado.'
    }
    return
  }

  if (hours <= EMERGENCY_HOURS) {
    advanceNoticeWarning.value = {
      type: 'emergency',
      title: 'SOLICITUD DE EMERGENCIA',
      message: `La cirugía está programada en menos de ${EMERGENCY_HOURS} horas. Esta solicitud será procesada con máxima prioridad.`
    }
  } else if (hours <= URGENT_HOURS) {
    advanceNoticeWarning.value = {
      type: 'urgent',
      title: 'SOLICITUD URGENTE',
      message: `La cirugía está programada en menos de ${URGENT_HOURS} horas. Esta solicitud será procesada con alta prioridad.`
    }
  } else if (days < MINIMUM_ADVANCE_DAYS) {
    advanceNoticeWarning.value = {
      type: 'warning',
      title: 'Anticipación recomendada',
      message: `Se recomienda programar con al menos ${MINIMUM_ADVANCE_DAYS} días de anticipación. La solicitud será marcada como urgente.`
    }
  } else {
    advanceNoticeWarning.value = null
  }
}

// Formatear fecha para mostrar (en modo edición)
const formatDateTimeForDisplay = (dateString) => {
  if (!dateString) return 'No especificada'
  
  try {
    const date = new Date(dateString)
    if (isNaN(date.getTime())) return 'Fecha inválida'
    
    const options = {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
      hour12: false
    }
    return date.toLocaleDateString('es-CL', options)
  } catch (error) {
    console.error('Error formateando fecha:', error)
    return dateString
  }
}

// Formulario de solicitud
const requestForm = reactive({
  pavilion_id: '',
  surgery_datetime: '',
  notes: '',
  items: [],
  surgery_id: null
})

// Datos originales para mostrar en modo edición (solo lectura)
const originalRequestData = ref({
  pavilion_name: '',
  surgery_datetime_display: '',
  requester_name: '',
  requester_rut: ''
})

// Computed: Filtrar cirugías según la especialidad del doctor logueado
const filteredSurgeries = computed(() => {
  const userSpecialtyId = authStore.getUserSpecialtyId
  
  // Si el usuario no tiene especialidad, mostrar todas las cirugías
  if (!userSpecialtyId) {
    return surgeries.value
  }
  
  // Filtrar cirugías que coincidan con la especialidad del usuario
  return surgeries.value.filter(surgery => surgery.specialty_id === userSpecialtyId)
})

// Template para nuevo item
const newSupplyItem = () => ({
  supply_code: '',
  supply_name: '',
  quantity_requested: 1,
  specifications: '',
  is_pediatric: false,
  special_requests: '',
  urgency_level: 'normal',
  size: '',
  brand: ''
})

// Métodos
const loadPavilions = async () => {
  loadingPavilions.value = true
  try {
    const result = await pavilionService.getAllPavilions()
    pavilions.value = result
  } catch (error) {
    console.error('Error cargando pabellones:', error)
    errors.value.push('Error al cargar la lista de pabellones')
    emit('error', error)
  } finally {
    loadingPavilions.value = false
  }
}

const loadSurgeries = async () => {
  loadingSurgeries.value = true
  try {
    const result = await surgeryService.getAllSurgeries()
    surgeries.value = result
  } catch (error) {
    console.error('Error cargando cirugías:', error)
  } finally {
    loadingSurgeries.value = false
  }
}

const loadTypicalSupplies = async (surgeryId) => {
  if (!surgeryId) {
    typicalSupplies.value = []
    showTypicalSupplies.value = false
    return
  }

  loadingTypicalSupplies.value = true
  try {
    const result = await surgeryTypicalSupplyService.getTypicalSuppliesBySurgeryId(surgeryId)
    typicalSupplies.value = result
    showTypicalSupplies.value = result.length > 0
  } catch (error) {
    console.error('Error cargando insumos típicos:', error)
    typicalSupplies.value = []
    showTypicalSupplies.value = false
  } finally {
    loadingTypicalSupplies.value = false
  }
}

// Watcher para cargar insumos típicos cuando se selecciona una cirugía
watch(() => requestForm.surgery_id, (newSurgeryId) => {
  if (newSurgeryId && !props.editMode) {
    loadTypicalSupplies(newSurgeryId)
  } else {
    typicalSupplies.value = []
    showTypicalSupplies.value = false
  }
})

const loadMedicalSupplies = async () => {
  loadingSupplies.value = true
  try {
    const result = await inventoryService.getAllSupplyCodes()
    medicalSupplies.value = result
  } catch (error) {
    console.error('Error cargando códigos de insumo:', error)
    errors.value.push('Error al cargar la lista de códigos de insumo')
    emit('error', error)
  } finally {
    loadingSupplies.value = false
  }
}

const addSupplyItem = () => {
  requestForm.items.unshift(newSupplyItem())
  // Inicializar estados de búsqueda para el nuevo item al principio
  supplySearchTerms.value.unshift('')
  showSupplyDropdowns.value.unshift(false)
  showCodeDropdowns.value.unshift(false)
  showSupplyListModals.value.unshift(false)
  supplyListSearchTerms.value.unshift('')
}

const addTypicalSupply = (typicalSupply) => {
  // Verificar si el insumo ya está en la lista
  const existingIndex = requestForm.items.findIndex(
    item => item.supply_code && item.supply_code.toString() === typicalSupply.supply_code.toString()
  )

  if (existingIndex !== -1) {
    // Si ya existe, actualizar la cantidad sumando la cantidad típica
    requestForm.items[existingIndex].quantity_requested += typicalSupply.typical_quantity || 1
    showInfo(`Se sumó la cantidad típica (${typicalSupply.typical_quantity || 1}) al insumo existente`)
  } else {
    // Buscar el nombre del insumo en la lista de medicalSupplies
    const supplyInfo = medicalSupplies.value.find(s => s.code === typicalSupply.supply_code)
    const supplyName = supplyInfo ? supplyInfo.name : `Insumo ${typicalSupply.supply_code}`

    // Buscar si hay un item vacío (sin código y sin nombre)
    const emptyItemIndex = requestForm.items.findIndex(
      item => (!item.supply_code || item.supply_code === '') && 
              (!item.supply_name || item.supply_name === '')
    )

    // Crear el nuevo item con la cantidad típica
    const newItem = {
      supply_code: typicalSupply.supply_code,
      supply_name: supplyName,
      quantity_requested: typicalSupply.typical_quantity || 1,
      specifications: '',
      is_pediatric: false,
      special_requests: typicalSupply.notes || '',
      urgency_level: 'normal',
      size: '',
      brand: ''
    }

    if (emptyItemIndex !== -1) {
      // Reemplazar el item vacío
      requestForm.items[emptyItemIndex] = newItem
      supplySearchTerms.value[emptyItemIndex] = supplyName
      showSupplyDropdowns.value[emptyItemIndex] = false
      showCodeDropdowns.value[emptyItemIndex] = false
      showSupplyListModals.value[emptyItemIndex] = false
      supplyListSearchTerms.value[emptyItemIndex] = ''
    } else {
      // No hay items vacíos, agregar uno nuevo
      requestForm.items.push(newItem)
      
      // Inicializar estados de búsqueda
      supplySearchTerms.value.push(supplyName)
      showSupplyDropdowns.value.push(false)
      showCodeDropdowns.value.push(false)
      showSupplyListModals.value.push(false)
      supplyListSearchTerms.value.push('')
    }
  }
}

const addAllTypicalSupplies = () => {
  if (typicalSupplies.value.length === 0) return

  // Primero, eliminar items vacíos (sin código ni nombre)
  const itemsToRemove = []
  requestForm.items.forEach((item, index) => {
    if (!item.supply_code || !item.supply_name || 
        item.supply_code === '' || item.supply_name === '') {
      itemsToRemove.push(index)
    }
  })
  
  // Eliminar items vacíos en orden inverso para no afectar los índices
  itemsToRemove.reverse().forEach(index => {
    requestForm.items.splice(index, 1)
    supplySearchTerms.value.splice(index, 1)
    showSupplyDropdowns.value.splice(index, 1)
    showCodeDropdowns.value.splice(index, 1)
    showSupplyListModals.value.splice(index, 1)
    supplyListSearchTerms.value.splice(index, 1)
  })

  let addedCount = 0
  let updatedCount = 0
  const newItems = []
  const newSearchTerms = []
  const newShowDropdowns = []
  const newCodeDropdowns = []
  const newShowListModals = []
  const newListSearchTerms = []

  typicalSupplies.value.forEach(typicalSupply => {
    const existingIndex = requestForm.items.findIndex(
      item => item.supply_code && item.supply_code.toString() === typicalSupply.supply_code.toString()
    )

    if (existingIndex !== -1) {
      requestForm.items[existingIndex].quantity_requested += typicalSupply.typical_quantity || 1
      updatedCount++
    } else {
      const supplyInfo = medicalSupplies.value.find(s => s.code === typicalSupply.supply_code)
      const supplyName = supplyInfo ? supplyInfo.name : `Insumo ${typicalSupply.supply_code}`

      const newItem = {
        supply_code: typicalSupply.supply_code,
        supply_name: supplyName,
        quantity_requested: typicalSupply.typical_quantity || 1,
        specifications: '',
        is_pediatric: false,
        special_requests: typicalSupply.notes || '',
        urgency_level: 'normal',
        size: '',
        brand: ''
      }
      
      newItems.push(newItem)
      newSearchTerms.push(supplyName)
      newShowDropdowns.push(false)
      newCodeDropdowns.push(false)
      newShowListModals.push(false)
      newListSearchTerms.push('')
      addedCount++
    }
  })

  // Agregar todos los nuevos items al principio
  if (newItems.length > 0) {
    requestForm.items.unshift(...newItems)
    supplySearchTerms.value.unshift(...newSearchTerms)
    showSupplyDropdowns.value.unshift(...newShowDropdowns)
    showCodeDropdowns.value.unshift(...newCodeDropdowns)
    showSupplyListModals.value.unshift(...newShowListModals)
    supplyListSearchTerms.value.unshift(...newListSearchTerms)
  }

  showSuccess(`${addedCount} insumo(s) agregado(s)${updatedCount > 0 ? `. ${updatedCount} insumo(s) actualizado(s)` : ''}`)
  
  // Minimizar la sección después de agregar todos
  typicalSuppliesExpanded.value = false
}

const removeSupplyItem = (index) => {
  // No permitir eliminar items aceptados en modo edición
  if (props.editMode && requestForm.items[index]?.item_status === 'aceptado') {
    showWarning('Este insumo ya fue aceptado por bodega y no puede ser eliminado.')
    return
  }
  
  requestForm.items.splice(index, 1)
  // También remover los estados de búsqueda correspondientes
  supplySearchTerms.value.splice(index, 1)
  showSupplyDropdowns.value.splice(index, 1)
  showCodeDropdowns.value.splice(index, 1)
  showSupplyListModals.value.splice(index, 1)
  supplyListSearchTerms.value.splice(index, 1)
}

// Funciones para autocompletado de insumos
const getFilteredSupplies = (index) => {
  const searchTerm = supplySearchTerms.value[index] || ''
  if (!searchTerm) return medicalSupplies.value.slice(0, 10) // Mostrar primeros 10 si no hay búsqueda
  
  return medicalSupplies.value.filter(supply => 
    supply.name?.toLowerCase().includes(searchTerm.toLowerCase()) ||
    supply.code?.toString().includes(searchTerm)
  ).slice(0, 10) // Limitar a 10 resultados
}

const getFilteredSuppliesByCode = (index) => {
  const codeTerm = requestForm.items[index]?.supply_code || ''
  if (!codeTerm) return medicalSupplies.value.slice(0, 10)
  
  const codeStr = codeTerm.toString()
  return medicalSupplies.value.filter(supply => 
    supply.code?.toString().includes(codeStr)
  ).slice(0, 10)
}

const getFilteredSuppliesForList = (index) => {
  const searchTerm = supplyListSearchTerms.value[index] || ''
  if (!searchTerm) return medicalSupplies.value
  
  const term = searchTerm.toLowerCase()
  return medicalSupplies.value.filter(supply => 
    supply.name?.toLowerCase().includes(term) ||
    supply.code?.toString().includes(searchTerm)
  )
}

const getSupplyNameForTypical = (supplyCode) => {
  const supply = medicalSupplies.value.find(s => s.code === supplyCode)
  return supply ? supply.name : `Insumo ${supplyCode}`
}

const selectSupply = (index, supply) => {
  requestForm.items[index].supply_name = supply.name
  requestForm.items[index].supply_code = supply.code
  supplySearchTerms.value[index] = supply.name
  showSupplyDropdowns.value[index] = false
  showCodeDropdowns.value[index] = false
  if (currentSupplyListIndex.value === index) {
    currentSupplyListIndex.value = null
  }
}

const onSupplyInputFocus = (index) => {
  showSupplyDropdowns.value[index] = true
  // Inicializar el término de búsqueda si no existe
  if (!supplySearchTerms.value[index]) {
    supplySearchTerms.value[index] = requestForm.items[index].supply_name || ''
  }
}

const onSupplyInputBlur = (index) => {
  // Delay para permitir clicks en el dropdown
  setTimeout(() => {
    showSupplyDropdowns.value[index] = false
  }, 200)
}

const onSupplyInputChange = (index, value) => {
  supplySearchTerms.value[index] = value
  requestForm.items[index].supply_name = value
  showSupplyDropdowns.value[index] = true
  // Si se escribe el nombre manualmente, intentar buscar el código
  if (value && medicalSupplies.value.length > 0) {
    const foundSupply = medicalSupplies.value.find(s => 
      s.name?.toLowerCase() === value.toLowerCase()
    )
    if (foundSupply) {
      requestForm.items[index].supply_code = foundSupply.code
    }
  }
}

const onSupplyCodeChange = (index, value) => {
  requestForm.items[index].supply_code = value
  showCodeDropdowns.value[index] = true
  // Si se escribe el código, intentar buscar el nombre
  if (value && medicalSupplies.value.length > 0) {
    const foundSupply = medicalSupplies.value.find(s => 
      s.code?.toString() === value.toString()
    )
    if (foundSupply) {
      requestForm.items[index].supply_name = foundSupply.name
      supplySearchTerms.value[index] = foundSupply.name
    }
  }
}

const onSupplyCodeFocus = (index) => {
  showCodeDropdowns.value[index] = true
  if (!supplySearchTerms.value[index]) {
    supplySearchTerms.value[index] = requestForm.items[index].supply_name || ''
  }
}

const onSupplyCodeBlur = (index) => {
  // Delay para permitir clicks en el dropdown
  setTimeout(() => {
    showCodeDropdowns.value[index] = false
  }, 200)
}

const toggleSupplyList = (index) => {
  if (currentSupplyListIndex.value === index) {
    currentSupplyListIndex.value = null
  } else {
    currentSupplyListIndex.value = index
    if (!supplyListSearchTerms.value[index]) {
      supplyListSearchTerms.value[index] = ''
    }
  }
}

const closeSupplyList = (index) => {
  if (currentSupplyListIndex.value === index) {
    currentSupplyListIndex.value = null
  }
}

// Verificar si un item es duplicado (mismo código que otro item)
const isDuplicateItem = (index) => {
  const currentItem = requestForm.items[index]
  if (!currentItem?.supply_code) return false
  
  // Contar cuántas veces aparece este código
  const count = requestForm.items.filter(item => 
    item.supply_code && item.supply_code.toString() === currentItem.supply_code.toString()
  ).length
  
  return count > 1
}

// Detectar y consolidar items duplicados
const consolidateDuplicateItems = () => {
  const consolidated = []
  const seen = new Map() // Map<supply_code, index_en_consolidated>
  
  requestForm.items.forEach((item, originalIndex) => {
    if (!item.supply_code) {
      // Si no tiene código, agregarlo sin consolidar
      consolidated.push({
        ...item,
        _originalIndex: originalIndex
      })
      return
    }
    
    const codeKey = item.supply_code.toString()
    
    if (seen.has(codeKey)) {
      // Ya existe este código, consolidar
      const existingIndex = seen.get(codeKey)
      const existingItem = consolidated[existingIndex]
      
      // Sumar cantidades
      existingItem.quantity_requested += item.quantity_requested || 0
      
      // Si uno es pediátrico, ambos deben ser pediátricos
      if (item.is_pediatric) {
        existingItem.is_pediatric = true
      }
      
      // Mostrar advertencia si hay diferencias en el nombre
      if (existingItem.supply_name !== item.supply_name) {
        console.warn(`Advertencia: El insumo ${codeKey} tiene nombres diferentes: "${existingItem.supply_name}" vs "${item.supply_name}"`)
      }
    } else {
      // Primera vez que vemos este código
      seen.set(codeKey, consolidated.length)
      consolidated.push({
        ...item,
        _originalIndex: originalIndex
      })
    }
  })
  
  return consolidated
}

const validateForm = () => {
  // Primero consolidar items duplicados
  const consolidated = consolidateDuplicateItems()
  
  // Verificar si hay duplicados
  const hasDuplicates = consolidated.length < requestForm.items.length
  
  if (hasDuplicates) {
    // Consolidar automáticamente sin preguntar
    const duplicateCount = requestForm.items.length - consolidated.length
    
    // Consolidar items
    const originalSearchTerms = [...supplySearchTerms.value]
    
    // Reconstruir arrays basados en items consolidados
    requestForm.items = consolidated.map(item => {
      const { _originalIndex, ...cleanItem } = item
      return cleanItem
    })
    
    // Reconstruir arrays de estados
    supplySearchTerms.value = requestForm.items.map((item, index) => {
      const originalIndex = consolidated[index]._originalIndex
      return originalSearchTerms[originalIndex] || item.supply_name || ''
    })
    
    showCodeDropdowns.value = requestForm.items.map(() => false)
    showSupplyDropdowns.value = requestForm.items.map(() => false)
    showSupplyListModals.value = requestForm.items.map(() => false)
    supplyListSearchTerms.value = requestForm.items.map(() => '')
    
    // Mostrar mensaje informativo
    showInfo(`Se detectaron ${duplicateCount} insumo(s) duplicado(s). Las cantidades fueron sumadas automáticamente.`)
  }
  
  // Validar el formulario (con items consolidados si se aplicó)
  const validation = supplyRequestService.validateSupplyRequest(requestForm)
  errors.value = validation.errors
  return validation.isValid
}

const resetForm = () => {
  Object.assign(requestForm, {
    pavilion_id: '',
    surgery_datetime: '',
    notes: '',
    items: [],
    surgery_id: null
  })
  errors.value = []
}

const cancelForm = () => {
  emit('cancel')
}

const submitRequest = async () => {
  if (!validateForm()) {
    return
  }

  // Verificar anticipación mínima antes de enviar
  if (!props.editMode && isNotProgrammed.value) {
    const days = Math.ceil(daysUntilSurgery.value)
    const result = await Swal.fire({
      icon: 'warning',
      title: 'Anticipación recomendada',
      html: `
        <p>La cirugía está programada en <strong>${days} día(s)</strong>, menos de los ${MINIMUM_ADVANCE_DAYS} días recomendados.</p>
        <p class="mt-2">¿Deseas continuar con la solicitud?</p>
        <p class="text-sm text-gray-600 mt-2">Esta solicitud será marcada como urgente y procesada con prioridad.</p>
      `,
      showCancelButton: true,
      confirmButtonText: 'Sí, continuar',
      cancelButtonText: 'Cancelar',
      confirmButtonColor: '#3085d6',
      cancelButtonColor: '#d33'
    })

    if (!result.isConfirmed) {
      return
    }
  }

  submitting.value = true
  errors.value = []

  try {
    // Si está en modo edición, reenviar la solicitud devuelta
    if (props.editMode && props.id) {
      await resubmitRequest()
    } else {
      // Crear nueva solicitud
      await createNewRequest()
    }
  } catch (error) {
    console.error('Error al enviar solicitud:', error)
    
    let errorMessage = 'Error desconocido al procesar la solicitud'
    if (error.response?.data?.error) {
      errorMessage = 'Error del servidor: ' + error.response.data.error
    } else if (error.message) {
      errorMessage = 'Error: ' + error.message
    }
    
    errors.value.push(errorMessage)
    emit('error', error)
  } finally {
    submitting.value = false
  }
}

// Crear nueva solicitud
const createNewRequest = async () => {
  const formattedData = supplyRequestService.formatSupplyRequestForAPI(requestForm)
  console.log('Enviando solicitud:', formattedData)
  
  const result = await supplyRequestService.createSupplyRequest(formattedData)
  
  if (result.success) {
    console.log('Solicitud creada exitosamente:', result)
    emit('success', result.data?.request || result.data)
  } else {
    const errorMessage = 'Error al crear la solicitud: ' + (result.error || 'Error desconocido')
    errors.value.push(errorMessage)
    emit('error', new Error(errorMessage))
  }
}

// Reenviar solicitud devuelta
const resubmitRequest = async () => {
  // Preparar items: los devueltos (con ID) y los nuevos (sin ID)
  const updatedItems = requestForm.items
    .filter(item => {
      // Incluir items devueltos (con ID y status devuelto) o nuevos items (sin ID)
      return (item.item_status === 'devuelto' && item.id) || (!item.id && item.supply_code && item.supply_name)
    })
    .map(item => {
      const itemData = {
        quantity: item.quantity_requested,
        notes: item.resubmit_notes || '' // Nuevas observaciones del doctor sobre este item
      }
      
      // Si tiene ID, es un item existente devuelto
      if (item.id) {
        itemData.item_id = item.id
      }
      // Si no tiene ID, es un nuevo item (se creará)
      
      // Incluir supply_code y supply_name si tienen valores válidos
      if (item.supply_code && item.supply_code.toString().trim() !== '') {
        itemData.supply_code = parseInt(item.supply_code)
      }
      if (item.supply_name && item.supply_name.trim() !== '') {
        itemData.supply_name = item.supply_name.trim()
      }
      
      // Incluir is_pediatric
      if (item.is_pediatric !== undefined && item.is_pediatric !== null) {
        itemData.is_pediatric = Boolean(item.is_pediatric)
      }
      
      return itemData
    })
  
  console.log('Reenviando solicitud con items:', updatedItems)
  console.log('Notas del solicitante:', requestForm.notes)
  
  const result = await supplyRequestService.resubmitReturnedRequest(props.id, updatedItems, requestForm.notes)
  
  if (result.success) {
    console.log('Solicitud reenviada exitosamente')
    
    showSuccess('La solicitud ha sido reenviada al encargado de bodega')
    
    emit('success', { id: props.id })
  } else {
    const errorMessage = 'Error al reenviar la solicitud: ' + (result.error || 'Error desconocido')
    errors.value.push(errorMessage)
    emit('error', new Error(errorMessage))
  }
}

// Lifecycle
onMounted(async () => {
  // Cargar listas primero
  await Promise.all([
    loadPavilions(),
    loadMedicalSupplies(),
    loadSurgeries()
  ])
  
  // Si está en modo edición, cargar la solicitud DESPUÉS de que las listas estén listas
  if (props.editMode && props.id) {
    await loadRequestForEdit()
  } else {
    // Agregar un insumo por defecto solo si no está en modo edición
    addSupplyItem()
  }
})

// Función auxiliar para extraer solo las notas originales del solicitante
const extractOriginalNotes = (fullNotes) => {
  if (!fullNotes) return ''
  
  // Buscar el marcador de devolución
  const devolucionIndex = fullNotes.indexOf('[Devolución por')
  
  // Si no hay marcador, devolver todas las notas
  if (devolucionIndex === -1) {
    return fullNotes.trim()
  }
  
  // Extraer solo la parte antes del marcador
  return fullNotes.substring(0, devolucionIndex).trim()
}

// Cargar solicitud para editar
const loadRequestForEdit = async () => {
  try {
    console.log('🔄 Cargando solicitud para editar, ID:', props.id)
    const response = await supplyRequestService.getSupplyRequestById(props.id)
    console.log('📦 Respuesta de la solicitud:', response)
    
    if (response.success && response.data) {
      // El backend devuelve data.request, data.items, data.assignments
      const request = response.data.request || response.data
      console.log('✅ Datos de la solicitud completa:', JSON.stringify(request, null, 2))
      console.log('🏥 Pabellones disponibles:', JSON.stringify(pavilions.value, null, 2))
      console.log('🔍 Datos directos - pavilion_id:', request.pavilion_id, 'surgery_datetime:', request.surgery_datetime)
      
      const pavilionId = request.pavilion_id
      const surgeryDatetime = request.surgery_datetime
      
      console.log('📊 Valores finales - pavilionId:', pavilionId, 'surgeryDatetime:', surgeryDatetime)
      
      // Cargar datos básicos
      requestForm.pavilion_id = pavilionId ? pavilionId.toString() : ''
      requestForm.surgery_datetime = surgeryDatetime ? formatDateTimeForInput(surgeryDatetime) : ''
      // En modo edición, dejar las observaciones vacías para que el usuario ingrese nuevas
      requestForm.notes = ''
      
      // Cargar tipo de cirugía - convertir a número si existe
      if (request.surgery_id !== null && request.surgery_id !== undefined) {
        requestForm.surgery_id = parseInt(request.surgery_id)
      } else {
        requestForm.surgery_id = null
      }
      console.log('🔍 Surgery ID cargado:', requestForm.surgery_id, 'desde:', request.surgery_id)
      
      // Guardar datos originales para mostrar en modo solo lectura
      // Buscar el pabellón por ID - probar con conversión de tipos
      const pavilionIdToFind = parseInt(pavilionId || 0)
      let pavilionName = 'No especificado'
      
      console.log('🔍 Buscando pabellón con ID:', pavilionIdToFind)
      console.log('📋 Pabellones disponibles:', pavilions.value)
      
      if (pavilionIdToFind > 0) {
        const foundPavilion = pavilions.value.find(p => {
          const pId = parseInt(p.id)
          console.log(`  Comparando: ${pId} === ${pavilionIdToFind}?`, pId === pavilionIdToFind)
          return pId === pavilionIdToFind
        })
        
        if (foundPavilion) {
          pavilionName = foundPavilion.name
          console.log('✅ Pabellón encontrado:', pavilionName)
        } else {
          console.warn('⚠️ Pabellón no encontrado para ID:', pavilionIdToFind)
          pavilionName = `Pabellón ID: ${pavilionIdToFind}`
        }
      }
      
      originalRequestData.value = {
        pavilion_name: pavilionName,
        surgery_datetime_display: surgeryDatetime ? formatDateTimeForDisplay(surgeryDatetime) : 'No especificada',
        requester_name: request.requested_by_name || 'No disponible',
        requester_rut: request.requested_by || 'No disponible'
      }
      
      console.log('📝 Datos originales asignados:', originalRequestData.value)
      
      console.log('📝 Datos originales guardados:', originalRequestData.value)
      
      // Forzar actualización del DOM
      await nextTick()
      
      console.log('📝 Formulario actualizado:', {
        pavilion_id: requestForm.pavilion_id,
        surgery_datetime: requestForm.surgery_datetime,
        notes: requestForm.notes
      })
      
      // Cargar items
      const itemsResponse = await supplyRequestService.getSupplyRequestItems(props.id)
      if (itemsResponse.success && itemsResponse.data) {
        requestForm.items = itemsResponse.data.map(item => ({
          id: item.id, // Incluir ID para tracking
          supply_code: item.supply_code,
          supply_name: item.supply_name,
          quantity_requested: item.quantity_requested,
          is_pediatric: item.is_pediatric,
          item_status: item.item_status, // Para saber cuáles son editables
          item_notes: item.item_notes, // Observaciones anteriores del encargado (solo lectura)
          resubmit_notes: '', // Nuevas observaciones del doctor (editable)
          specifications: '',
          special_requests: '',
          urgency_level: 'normal',
          size: '',
          brand: ''
        }))
        
        // Inicializar arrays de búsqueda
        supplySearchTerms.value = requestForm.items.map(() => '')
        showSupplyDropdowns.value = requestForm.items.map(() => false)
        showCodeDropdowns.value = requestForm.items.map(() => false)
        showSupplyListModals.value = requestForm.items.map(() => false)
        supplyListSearchTerms.value = requestForm.items.map(() => '')
      }
    }
  } catch (error) {
    console.error('Error cargando solicitud:', error)
    showError('No se pudo cargar la solicitud para editar')
  }
}

// Función auxiliar para formatear fecha para input datetime-local
const formatDateTimeForInput = (dateString) => {
  const date = new Date(dateString)
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  return `${year}-${month}-${day}T${hours}:${minutes}`
}

</script>

<style scoped>
/* Estilos adicionales si son necesarios */
.bg-red-50 {
  background-color: #fef2f2;
}

.border-red-200 {
  border-color: #fecaca;
}

.text-red-400 {
  color: #f87171;
}

.text-red-800 {
  color: #991b1b;
}

.text-red-700 {
  color: #b91c1c;
}

.text-red-500 {
  color: #ef4444;
}

.hover\:bg-red-50:hover {
  background-color: #fef2f2;
}

.hover\:text-red-700:hover {
  color: #b91c1c;
}

/* Mejorar experiencia táctil en móviles */
@media (max-width: 640px) {
  /* Aumentar área táctil de botones */
  button {
    min-height: 44px;
  }
  
  /* Mejorar áreas de entrada */
  input,
  select,
  textarea {
    min-height: 44px;
    font-size: 16px; /* Prevenir zoom automático en iOS */
  }
  
  /* Mejorar área de checkboxes */
  input[type="checkbox"] {
    min-width: 20px;
    min-height: 20px;
  }

  /* Suavizar transiciones */
  * {
    -webkit-tap-highlight-color: transparent;
  }
}

/* Transiciones suaves */
button,
input,
select,
textarea {
  transition: all 0.2s ease-in-out;
}
</style>