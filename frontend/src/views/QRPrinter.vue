<template>
  <!-- Modal de configuración de impresión -->
  <div v-if="showModal" class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50">
    <div class="relative top-20 mx-auto p-5 border w-11/12 max-w-md shadow-lg rounded-md bg-white">
      <div class="mt-3">
        <!-- Header del modal -->
        <div class="flex items-center justify-between pb-3 border-b">
          <h3 class="text-lg font-medium text-gray-900">
            Configurar Impresión
          </h3>
          <button 
            @click="closeModal"
            class="text-gray-400 hover:text-gray-600 focus:outline-none"
          >
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <!-- Configuraciones -->
        <div class="mt-4 space-y-4">
          <!-- Tamaño de etiqueta -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Tamaño de Etiqueta:
            </label>
            <select v-model="printConfig.labelSize" class="form-select w-full">
              <option value="small">Pequeña (2x2 cm)</option>
              <option value="medium">Mediana (4x4 cm)</option>
              <option value="large">Grande (6x6 cm)</option>
              <option value="custom">Personalizado</option>
            </select>
          </div>

          <!-- Información adicional -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Información a incluir:
            </label>
            <div class="space-y-2">
              <label class="flex items-center">
                <input 
                  type="checkbox" 
                  v-model="printConfig.includeCode"
                  class="form-checkbox"
                />
                <span class="ml-2 text-sm">Código QR (texto)</span>
              </label>
              <label class="flex items-center">
                <input 
                  type="checkbox" 
                  v-model="printConfig.includeName"
                  class="form-checkbox"
                />
                <span class="ml-2 text-sm">Nombre del producto</span>
              </label>
              <label class="flex items-center">
                <input 
                  type="checkbox" 
                  v-model="printConfig.includeDate"
                  class="form-checkbox"
                />
                <span class="ml-2 text-sm">Fecha de impresión</span>
              </label>
              <label class="flex items-center">
                <input 
                  type="checkbox" 
                  v-model="printConfig.includeExpiry"
                  class="form-checkbox"
                />
                <span class="ml-2 text-sm">Fecha de vencimiento</span>
              </label>
            </div>
          </div>

          <!-- Cantidad -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Cantidad de etiquetas:
            </label>
            <input 
              v-model="printConfig.quantity"
              type="number"
              min="1"
              max="100"
              class="form-input w-full"
            />
          </div>

          <!-- Formato -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">
              Orientación:
            </label>
            <div class="flex space-x-4">
              <label class="flex items-center">
                <input 
                  type="radio" 
                  value="portrait" 
                  v-model="printConfig.orientation"
                  class="form-radio"
                />
                <span class="ml-2 text-sm">Vertical</span>
              </label>
              <label class="flex items-center">
                <input 
                  type="radio" 
                  value="landscape" 
                  v-model="printConfig.orientation"
                  class="form-radio"
                />
                <span class="ml-2 text-sm">Horizontal</span>
              </label>
            </div>
          </div>
        </div>

        <!-- Vista previa -->
        <div class="mt-6 p-4 border rounded-lg bg-gray-50">
          <h4 class="text-sm font-medium text-gray-700 mb-2">Vista Previa:</h4>
          <div class="preview-container" :class="getPreviewClass()">
            <div class="preview-label">
              <div class="preview-qr">
                <div class="qr-placeholder"></div>
              </div>
              <div v-if="printConfig.includeCode" class="preview-text code">
                {{ qrData.qr_code }}
              </div>
              <div v-if="printConfig.includeName" class="preview-text name">
                {{ qrData.name || 'Nombre del Producto' }}
              </div>
              <div v-if="printConfig.includeDate" class="preview-text date">
                {{ formatDate(new Date()) }}
              </div>
              <div v-if="printConfig.includeExpiry && qrData.expiry_date" class="preview-text expiry">
                Vence: {{ formatDate(qrData.expiry_date) }}
              </div>
            </div>
          </div>
        </div>

        <!-- Botones de acción -->
        <div class="flex justify-end space-x-3 mt-6 pt-4 border-t">
          <button @click="closeModal" class="btn-secondary">
            Cancelar
          </button>
          <button @click="printQR" class="btn-primary">
            <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 17h2a2 2 0 002-2v-4a2 2 0 00-2-2H5a2 2 0 00-2 2v4a2 2 0 002 2h2m2 4h6a2 2 0 002-2v-4a2 2 0 00-2-2H9a2 2 0 00-2 2v4a2 2 0 002 2zm8-12V5a2 2 0 00-2-2H9a2 2 0 00-2 2v4h10z" />
            </svg>
            Imprimir
          </button>
        </div>
      </div>
    </div>
  </div>

  <!-- Área de impresión (oculta) -->
  <div ref="printArea" class="print-only">
    <div 
      v-for="n in printConfig.quantity" 
      :key="n" 
      class="print-label"
      :class="getPrintLabelClass()"
    >
      <div class="print-content">
        <div class="print-qr-section">
          <img 
            :src="qrImageUrl" 
            :alt="`QR Code ${qrData.qr_code}`"
            class="print-qr-img"
          />
        </div>
        
        <div v-if="hasTextInfo()" class="print-info-section">
          <div v-if="printConfig.includeCode" class="print-code">
            {{ qrData.qr_code }}
          </div>
          <div v-if="printConfig.includeName" class="print-name">
            {{ qrData.name }}
          </div>
          <div v-if="printConfig.includeDate" class="print-date">
            {{ formatDate(new Date()) }}
          </div>
          <div v-if="printConfig.includeExpiry && qrData.expiry_date" class="print-expiry">
            Vence: {{ formatDate(qrData.expiry_date) }}
          </div>
        </div>
      </div>
      
      <!-- Separador entre etiquetas (excepto la última) -->
      <div v-if="n < printConfig.quantity" class="print-separator"></div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { format } from 'date-fns'
import { es } from 'date-fns/locale'
import qrService from '@/services/qrService'

// Props
const props = defineProps({
  qrData: {
    type: Object,
    required: true,
    default: () => ({})
  },
  visible: {
    type: Boolean,
    default: false
  }
})

// Emits
const emit = defineEmits(['close', 'printed'])

// Referencias
const printArea = ref(null)

// Estado reactivo
const showModal = ref(false)
const printConfig = ref({
  labelSize: 'medium',
  includeCode: true,
  includeName: true,
  includeDate: false,
  includeExpiry: false,
  quantity: 1,
  orientation: 'portrait'
})

// Computed
const qrImageUrl = computed(() => {
  if (!props.qrData?.qr_code) return ''
  return qrService.getQRImageUrl(props.qrData.qr_code)
})

// Watchers
import { watch } from 'vue'

watch(() => props.visible, (newValue) => {
  showModal.value = newValue
})

// Métodos principales
const openModal = () => {
  showModal.value = true
}

const closeModal = () => {
  showModal.value = false
  emit('close')
}

const printQR = () => {
  if (!printArea.value) return
  
  // Crear una nueva ventana para imprimir
  const printWindow = window.open('', '_blank')
  const printContent = generatePrintHTML()
  
  printWindow.document.write(printContent)
  printWindow.document.close()
  
  // Esperar a que se carguen las imágenes y luego imprimir
  printWindow.onload = () => {
    setTimeout(() => {
      printWindow.print()
      printWindow.close()
      
      emit('printed')
      closeModal()
    }, 500)
  }
}

// Métodos auxiliares
const generatePrintHTML = () => {
  const labels = []
  
  for (let i = 0; i < printConfig.value.quantity; i++) {
    labels.push(`
      <div class="print-label ${getPrintLabelClass()}">
        <div class="print-content">
          <div class="print-qr-section">
            <img src="${qrImageUrl.value}" alt="QR Code" class="print-qr-img" />
          </div>
          ${hasTextInfo() ? `
          <div class="print-info-section">
            ${printConfig.value.includeCode ? `<div class="print-code">${props.qrData.qr_code}</div>` : ''}
            ${printConfig.value.includeName ? `<div class="print-name">${props.qrData.name || 'N/A'}</div>` : ''}
            ${printConfig.value.includeDate ? `<div class="print-date">${formatDate(new Date())}</div>` : ''}
            ${printConfig.value.includeExpiry && props.qrData.expiry_date ? `<div class="print-expiry">Vence: ${formatDate(props.qrData.expiry_date)}</div>` : ''}
          </div>
          ` : ''}
        </div>
      </div>
    `)
  }
  
  return `
    <!DOCTYPE html>
    <html>
    <head>
      <title>Impresión QR - ${props.qrData.qr_code}</title>
      <style>
        ${getPrintStyles()}
      </style>
    </head>
    <body>
      <div class="print-container">
        ${labels.join('')}
      </div>
    </body>
    </html>
  `
}

const getPrintStyles = () => {
  const sizes = {
    small: { width: '2cm', height: '2cm', fontSize: '6px' },
    medium: { width: '4cm', height: '4cm', fontSize: '8px' },
    large: { width: '6cm', height: '6cm', fontSize: '10px' }
  }
  
  const size = sizes[printConfig.value.labelSize] || sizes.medium
  
  return `
    * { margin: 0; padding: 0; box-sizing: border-box; }
    
    body {
      font-family: Arial, sans-serif;
      background: white;
      color: black;
    }
    
    .print-container {
      display: flex;
      flex-wrap: wrap;
      gap: 2mm;
      padding: 5mm;
    }
    
    .print-label {
      width: ${size.width};
      height: ${size.height};
      border: 1px solid #333;
      page-break-inside: avoid;
      display: flex;
      flex-direction: column;
      background: white;
    }
    
    .print-label.portrait {
      flex-direction: column;
    }
    
    .print-label.landscape {
      flex-direction: row;
      width: ${size.height};
      height: ${size.width};
    }
    
    .print-content {
      display: flex;
      flex-direction: inherit;
      height: 100%;
      padding: 1mm;
    }
    
    .print-qr-section {
      flex: ${hasTextInfo() ? '0 0 70%' : '1'};
      display: flex;
      align-items: center;
      justify-content: center;
    }
    
    .print-qr-img {
      max-width: 100%;
      max-height: 100%;
      object-fit: contain;
    }
    
    .print-info-section {
      flex: 1;
      display: flex;
      flex-direction: column;
      justify-content: center;
      padding: 0.5mm;
      font-size: ${size.fontSize};
      line-height: 1.2;
      text-align: center;
    }
    
    .print-code {
      font-family: monospace;
      font-weight: bold;
      word-break: break-all;
      margin-bottom: 0.5mm;
    }
    
    .print-name {
      font-weight: bold;
      margin-bottom: 0.5mm;
    }
    
    .print-date, .print-expiry {
      font-size: calc(${size.fontSize} - 1px);
      color: #666;
    }
    
    @media print {
      body { margin: 0; }
      .print-container { padding: 0; gap: 1mm; }
    }
  `
}

const getPreviewClass = () => {
  return [
    'preview-orientation',
    printConfig.value.orientation,
    printConfig.value.labelSize
  ]
}

const getPrintLabelClass = () => {
  return [
    printConfig.value.orientation,
    printConfig.value.labelSize
  ].join(' ')
}

const hasTextInfo = () => {
  return printConfig.value.includeCode || 
         printConfig.value.includeName || 
         printConfig.value.includeDate || 
         printConfig.value.includeExpiry
}

const formatDate = (date) => {
  try {
    return format(new Date(date), 'dd/MM/yyyy', { locale: es })
  } catch (error) {
    return 'N/A'
  }
}

// Exponer métodos públicos
defineExpose({
  openModal,
  closeModal,
  printQR
})
</script>

