import axios from 'axios'
import inventoryService from '@/services/inventory/inventoryService.js'
import { useAuthStore } from '@/stores/auth.js'
import { getApiBaseUrl } from '@/config/api.js'

const API_BASE_URL = getApiBaseUrl()

class QRService {
  constructor() {
    this.api = axios.create({
      baseURL: API_BASE_URL,
      headers: {
        'Content-Type': 'application/json'
      }
    })

    // Información de sesión y contexto para trazabilidad
    this.sessionId = this.generateSessionId()
    this.deviceInfo = this.detectDeviceInfo()
    this.browserInfo = this.detectBrowserInfo()
    this.setupRequestInterceptors()
  }

  // ========================
  // CONFIGURACIÓN DE INTERCEPTORS PARA TRAZABILIDAD
  // ========================

  setupRequestInterceptors() {
    // Interceptor para agregar headers de trazabilidad y autenticación
    this.api.interceptors.request.use((config) => {
      // Agregar headers de trazabilidad
      config.headers['X-Session-ID'] = this.sessionId
      config.headers['X-Device-Info'] = JSON.stringify(this.deviceInfo)
      config.headers['X-Browser-Info'] = JSON.stringify(this.browserInfo)

      // Agregar token de autenticación
      const token = localStorage.getItem('auth_token')
      if (token) {
        config.headers.Authorization = `Bearer ${token}`
      }

      return config
    })
  }

  // ========================
  // ESCANEO CON TRAZABILIDAD AUTOMÁTICA
  // ========================

  // Escanear un código QR con contexto completo para trazabilidad
  async scanQRCode(qrCode, scanContext = {}) {
    try {

      // Construir contexto completo de escaneo
      const fullContext = {
        scan_purpose: 'lookup',
        scan_source: 'web',
        session_id: this.sessionId,
        user_rut: this.getCurrentUserRUT(),
        user_name: this.getCurrentUserName(),
        pavilion_id: this.getCurrentPavilionId(),
        medical_center_id: this.getCurrentMedicalCenterId(),
        notes: '',
        ...scanContext // Override con contexto específico
      }

      // Construir query params para contexto
      const queryParams = new URLSearchParams()
      Object.keys(fullContext).forEach(key => {
        if (fullContext[key] !== null && fullContext[key] !== undefined && fullContext[key] !== '') {
          queryParams.append(key, fullContext[key])
        }
      })

      const response = await this.api.get(`/qr/scan/${encodeURIComponent(qrCode)}?${queryParams.toString()}`)

      // Extraer datos de la respuesta del backend
      let result = response.data

      // Si la respuesta viene en formato {success: true, data: {...}}, preservar la estructura
      if (result && result.success && result.data) {
        const originalSuccess = result.success
        const originalMessage = result.message
        const dataContent = result.data

        // Crear nuevo objeto que mantenga la estructura esperada
        result = {
          ...dataContent,
          success: originalSuccess,
          ...(originalMessage && { backend_message: originalMessage })
        }
      }

      // Mapear tipos del backend al frontend si es necesario
      if (result && result.qr_type && !result.type) {
        if (result.qr_type === 'SUPPLY') {
          result.type = 'medical_supply'
        } else if (result.qr_type === 'BATCH') {
          result.type = 'batch'
        }
      }

      // Añadir información adicional para mejor UX
      if (result) {
        result.scan_timestamp = new Date()
        result.is_individual_supply = result.type === 'medical_supply'
        result.is_batch = result.type === 'batch'
        result.scan_context = fullContext

        // Para insumos individuales, verificar estado de consumo
        if (result.type === 'medical_supply') {
          result.can_be_consumed = !result.is_consumed && result.available_for_use !== false

          // Normalizar el estado del insumo desde la respuesta del backend
          if (result.supply_info && result.supply_info.Status) {
            result.status = result.supply_info.Status
            result.current_status = result.supply_info.Status
          } else if (result.status) {
            result.current_status = result.status
          } else if (result.current_status) {
            result.status = result.current_status
          }

          // Añadir información del lote si está disponible
          if (result.batch_status) {
            result.batch_context = {
              has_stock: result.batch_status.current_amount > 0,
              low_stock: result.batch_status.current_amount < 10,
              batch_id: result.batch_status.batch_id
            }
          }
        }

        // Para lotes, añadir información sobre insumos individuales
        if (result.type === 'batch') {
          if (result.batch_status) {
            result.individual_supplies_info = {
              total: result.batch_status.total_individual_supplies,
              available: result.batch_status.available_supplies,
              consumed: result.batch_status.consumed_supplies,
              consumption_rate: result.batch_status.total_individual_supplies > 0
                ? Math.round((result.batch_status.consumed_supplies / result.batch_status.total_individual_supplies) * 100)
                : 0
            }
          }
        }

        // Agregar información de trazabilidad si está disponible
        result.traceability_available = !!(result.traceability || result.scan_events || result.scan_statistics)
      }

      // Asegurar que siempre retornemos la estructura esperada por el frontend
      if (!result.hasOwnProperty('success')) {
        return {
          success: true,
          data: result,
          message: result.backend_message || 'QR escaneado correctamente'
        }
      }

      return result
    } catch (error) {
      console.error('Error al escanear código QR:', error)

      // Mejorar mensajes de error según el tipo
      if (error.response?.status === 404) {
        const qrType = this.getQRType(qrCode)
        if (qrType === 'SUPPLY') {
          throw new Error('Insumo individual no encontrado. Verifique que el código QR sea correcto.')
        } else if (qrType === 'BATCH') {
          throw new Error('Lote no encontrado. Verifique que el código QR sea correcto.')
        } else {
          throw new Error('Código QR no encontrado o inválido.')
        }
      }

      throw error
    }
  }

  // Escanear para consumo específico con registro automático
  async scanForConsumption(qrCode, consumptionContext = {}) {
    return this.scanQRCode(qrCode, {
      scan_purpose: 'consume',
      ...consumptionContext
    })
  }

  // Escanear para verificación con registro automático
  async scanForVerification(qrCode, verificationContext = {}) {
    return this.scanQRCode(qrCode, {
      scan_purpose: 'verify',
      ...verificationContext
    })
  }

  // Escanear para inventario con registro automático
  async scanForInventory(qrCode, inventoryContext = {}) {
    return this.scanQRCode(qrCode, {
      scan_purpose: 'inventory_check',
      ...inventoryContext
    })
  }

  // ========================
  // NUEVAS FUNCIONALIDADES DE TRAZABILIDAD COMPLETA
  // ========================

  // Obtener trazabilidad completa de un QR
  async getCompleteTraceability(qrCode) {
    try {
      const response = await this.api.get(`/qr/traceability/${encodeURIComponent(qrCode)}`)

      if (response.data && response.data.success) {
        return response.data.data
      }

      return response.data
    } catch (error) {
      console.error('Error obteniendo trazabilidad completa:', error)
      throw error
    }
  }

  // Obtener historial de escaneos específicamente
  async getScanHistory(qrCode, limit = 50) {
    try {
      const response = await this.api.get(`/qr/scan-history/${encodeURIComponent(qrCode)}?limit=${limit}`)

      if (response.data && response.data.success) {
        return response.data.data
      }

      return response.data
    } catch (error) {
      console.error('Error obteniendo historial de escaneos:', error)

      // Si es un 404, devolver array vacío en lugar de lanzar el error
      if (error.response?.status === 404) {
        return []
      }

      throw error
    }
  }

  // Obtener estadísticas de escaneo
  async getScanStatistics(qrCode) {
    try {
      const response = await this.api.get(`/qr/scan-stats/${encodeURIComponent(qrCode)}`)

      if (response.data && response.data.success) {
        return response.data.data
      }

      return response.data
    } catch (error) {
      console.error('Error obteniendo estadísticas de escaneo:', error)

      // Si es un 404, devolver datos por defecto en lugar de lanzar el error
      if (error.response?.status === 404) {
        return {
          total_scans: 0,
          last_scan: null,
          scan_sources: [],
          scan_purposes: []
        }
      }

      throw error
    }
  }

  // Registrar evento de escaneo manualmente
  async registerScanEvent(qrCode, eventContext = {}) {
    try {
      const payload = {
        qr_code: qrCode,
        user_rut: this.getCurrentUserRUT(),
        user_name: this.getCurrentUserName(),
        pavilion_id: this.getCurrentPavilionId(),
        medical_center_id: this.getCurrentMedicalCenterId(),
        scan_source: 'web',
        ...eventContext
      }

      const response = await this.api.post('/qr/register-scan', payload)

      if (response.data && response.data.success) {
        return response.data.data
      }

      return response.data
    } catch (error) {
      console.error('Error registrando evento de escaneo:', error)
      throw error
    }
  }

  // ========================
  // VALIDACIÓN Y DETECCIÓN DE CONTEXTO
  // ========================

  // Validar si un código QR es válido con información detallada
  async validateQRCode(qrCode) {
    try {
      const response = await this.api.get(`/qr/validate/${encodeURIComponent(qrCode)}`)
      return {
        ...response.data,
        qr_type: this.getQRType(qrCode),
        is_individual_supply: this.getQRType(qrCode) === 'SUPPLY',
        is_batch: this.getQRType(qrCode) === 'BATCH'
      }
    } catch (error) {
      console.error('Error al validar código QR:', error)
      throw error
    }
  }

  // Obtener información del usuario actual
  getCurrentUserInfo() {
    return {
      rut: this.getCurrentUserRUT(),
      name: this.getCurrentUserName(),
      pavilionId: this.getCurrentPavilionId(),
      medicalCenterId: this.getCurrentMedicalCenterId()
    }
  }

  // Obtener RUT del usuario actual
  getCurrentUserRUT() {
    try {
      const authStore = useAuthStore()
      return authStore.user?.rut ||
        localStorage.getItem('user_rut') ||
        null
    } catch (error) {
      return localStorage.getItem('user_rut') || null
    }
  }

  // Obtener nombre del usuario actual
  getCurrentUserName() {
    try {
      const authStore = useAuthStore()
      return authStore.user?.name ||
        localStorage.getItem('user_name') ||
        null
    } catch (error) {
      return localStorage.getItem('user_name') || null
    }
  }

  // Obtener ID del pabellón actual
  getCurrentPavilionId() {
    const pavilionId = localStorage.getItem('current_pavilion_id') ||
      this.$store?.state?.location?.pavilionId ||
      null
    return pavilionId ? parseInt(pavilionId) : null
  }

  // Obtener ID del centro médico actual
  getCurrentMedicalCenterId() {
    const medicalCenterId = localStorage.getItem('current_medical_center_id') ||
      this.$store?.state?.location?.medicalCenterId ||
      null
    return medicalCenterId ? parseInt(medicalCenterId) : null
  }

  // Establecer ubicación actual para todos los escaneos
  setCurrentLocation(pavilionId, medicalCenterId) {
    if (pavilionId) {
      localStorage.setItem('current_pavilion_id', pavilionId.toString())
    }
    if (medicalCenterId) {
      localStorage.setItem('current_medical_center_id', medicalCenterId.toString())
    }
  }

  // Limpiar ubicación actual
  clearCurrentLocation() {
    localStorage.removeItem('current_pavilion_id')
    localStorage.removeItem('current_medical_center_id')
  }

  // ========================
  // DETECCIÓN DE DISPOSITIVO Y NAVEGADOR
  // ========================

  // Detectar información del dispositivo
  detectDeviceInfo() {
    const ua = navigator.userAgent

    return {
      platform: this.detectPlatform(ua),
      device_type: this.detectDeviceType(ua),
      screen_size: `${screen.width}x${screen.height}`,
      touch_enabled: 'ontouchstart' in window,
      language: navigator.language || navigator.userLanguage
    }
  }

  // Detectar información del navegador
  detectBrowserInfo() {
    const ua = navigator.userAgent

    return {
      name: this.detectBrowserName(ua),
      version: this.detectBrowserVersion(ua),
      engine: this.detectBrowserEngine(ua),
      cookies_enabled: navigator.cookieEnabled,
      javascript_enabled: true
    }
  }

  // Detectar plataforma
  detectPlatform(ua) {
    if (ua.includes('Windows')) return 'Windows'
    if (ua.includes('Mac')) return 'Mac'
    if (ua.includes('Linux')) return 'Linux'
    if (ua.includes('Android')) return 'Android'
    if (ua.includes('iOS') || ua.includes('iPhone') || ua.includes('iPad')) return 'iOS'
    return 'Unknown'
  }

  // Detectar tipo de dispositivo
  detectDeviceType(ua) {
    if (ua.includes('Mobile') || ua.includes('Android')) return 'mobile'
    if (ua.includes('Tablet') || ua.includes('iPad')) return 'tablet'
    return 'desktop'
  }

  // Detectar nombre del navegador
  detectBrowserName(ua) {
    if (ua.includes('Chrome')) return 'Chrome'
    if (ua.includes('Firefox')) return 'Firefox'
    if (ua.includes('Safari') && !ua.includes('Chrome')) return 'Safari'
    if (ua.includes('Edge')) return 'Edge'
    if (ua.includes('Opera')) return 'Opera'
    return 'Unknown'
  }

  // Detectar versión del navegador
  detectBrowserVersion(ua) {
    const match = ua.match(/(Chrome|Firefox|Safari|Edge|Opera)\/([0-9.]+)/)
    return match ? match[2] : 'Unknown'
  }

  // Detectar motor del navegador
  detectBrowserEngine(ua) {
    if (ua.includes('Blink') || ua.includes('Chrome')) return 'Blink'
    if (ua.includes('Gecko') && ua.includes('Firefox')) return 'Gecko'
    if (ua.includes('WebKit') && ua.includes('Safari')) return 'WebKit'
    return 'Unknown'
  }

  // Generar ID de sesión único
  generateSessionId() {
    return 'qr_session_' + Date.now() + '_' + Math.random().toString(36).substr(2, 9)
  }

  // ========================
  // HISTORIAL Y DETALLE DE INSUMOS (FUNCIONES EXISTENTES MEJORADAS)
  // ========================

  // Obtener historial de un insumo individual por código QR
  async getSupplyHistory(qrCode) {
    try {
      if (!this.isIndividualSupply(qrCode)) {
        throw new Error('Este endpoint está diseñado para insumos individuales. Use códigos QR que comiencen con SUPPLY_')
      }

      const response = await this.api.get(`/qr/history/${encodeURIComponent(qrCode)}`)

      console.log('Raw history response:', response.data)

      // Manejar la estructura de respuesta del backend
      let historyData = null

      if (response.data) {
        if (response.data.success && response.data.data) {
          historyData = response.data.data
        } else if (Array.isArray(response.data)) {
          historyData = response.data
        } else {
          historyData = response.data
        }
      }

      // Ordenar historial por fecha descendente si es array
      if (Array.isArray(historyData)) {
        historyData.sort((a, b) => {
          const dateA = new Date(a.date_time || a.DateTime)
          const dateB = new Date(b.date_time || b.DateTime)
          return dateB - dateA
        })
      }

      console.log('Processed history data:', historyData)

      return historyData
    } catch (error) {
      console.error('Error al obtener historial:', error)
      throw error
    }
  }

  async getBatchById(batchId) {
    try {
      const response = await this.api.get(`/batches/${batchId}`)

      // Extraer datos de la respuesta del backend
      let result = response.data

      // Si la respuesta viene en formato {success: true, data: {...}}
      if (result && result.success && result.data) {
        result = result.data
      }

      return result
    } catch (error) {
      console.error('Error al obtener información del lote:', error)
      throw error
    }
  }

  // Obtener información detallada de un insumo individual
  async getSupplyDetails(qrCode) {
    try {
      const response = await this.api.get(`/qr/details/${encodeURIComponent(qrCode)}`)

      if (response.data) {
        // Añadir información calculada útil para el frontend
        response.data.ui_info = {
          type_label: this.getTypeLabel(response.data.type),
          status_label: this.getStatusLabel(response.data),
          status_color: this.getStatusColor(response.data),
          can_consume: !response.data.is_consumed && response.data.available_for_use,
          is_priority_scan: response.data.type === 'medical_supply'
        }

        // Para insumos individuales, calcular días hasta vencimiento
        if (response.data.type === 'medical_supply' && response.data.supply_info?.expiration_date) {
          const expDate = new Date(response.data.supply_info.expiration_date)
          const today = new Date()
          const diffTime = expDate - today
          const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))

          response.data.expiration_info = {
            days_until_expiration: diffDays,
            is_expired: diffDays < 0,
            is_expiring_soon: diffDays <= 30 && diffDays >= 0,
            expiration_status: diffDays < 0 ? 'expired' :
              diffDays <= 7 ? 'critical' :
                diffDays <= 30 ? 'warning' : 'good'
          }
        }
      }

      return response.data
    } catch (error) {
      console.error('Error al obtener detalles del insumo:', error)
      throw error
    }
  }

  // ========================
  // CONSUMO DE INSUMOS CON TRAZABILIDAD Y HISTORIAL DE LOTES
  // ========================

  // Consumir un insumo individual por su código QR con trazabilidad y actualización de lote
  async consumeSupply(consumptionData) {
    try {
      // Validar que sea un insumo individual
      if (!this.isIndividualSupply(consumptionData.qr_code)) {
        throw new Error('Solo se pueden consumir insumos individuales. Use códigos QR que comiencen con SUPPLY_')
      }

      // Obtener información del insumo antes del consumo para saber su lote
      const supplyInfo = await this.getQRInfo(consumptionData.qr_code)

      if (!supplyInfo.supply_info || !supplyInfo.supply_info.batch) {
        throw new Error('No se pudo obtener información del lote del insumo')
      }

      const batchId = supplyInfo.supply_info.batch.id
      const currentBatchAmount = supplyInfo.supply_info.batch.amount

      // Consumir el insumo individual
      const response = await this.api.post('/qr/consume', {
        ...consumptionData,
        consumption_timestamp: new Date().toISOString()
      })

      // Si el consumo fue exitoso, actualizar el historial del lote
      if (response.data && response.data.success) {
        try {
          await this.updateBatchHistoryAfterConsumption({
            batchId: batchId,
            previousAmount: currentBatchAmount,
            newAmount: currentBatchAmount - 1,
            userRUT: consumptionData.user_rut,
            userName: consumptionData.user_name || 'Usuario',
            consumedQRCode: consumptionData.qr_code,
            destination: `${consumptionData.destination_type} ${consumptionData.destination_id}`
          })
        } catch (historyError) {
          console.warn('Error al actualizar historial del lote:', historyError)
        }
      }

      return {
        ...response.data,
        ui_message: 'Insumo individual consumido exitosamente y historial actualizado',
        next_actions: [
          { label: 'Ver historial del lote', action: 'view_batch_history' },
          { label: 'Ver historial del insumo', action: 'view_supply_history' },
          { label: 'Escanear otro', action: 'scan_next' },
          { label: 'Ver inventario', action: 'view_inventory' }
        ]
      }
    } catch (error) {
      console.error('Error al consumir insumo:', error)

      // Mejorar mensajes de error para insumos individuales
      if (error.response?.data?.message?.includes('ya ha sido consumido')) {
        throw new Error('Este insumo individual ya ha sido consumido anteriormente.')
      } else if (error.response?.data?.message?.includes('no encontrado')) {
        throw new Error('Insumo individual no encontrado. Verifique el código QR.')
      }

      throw error
    }
  }

  // Método para actualizar el historial del lote después del consumo
  async updateBatchHistoryAfterConsumption(historyData) {
    try {
      const {
        batchId,
        previousAmount,
        newAmount,
        userRUT,
        userName,
        consumedQRCode,
        destination
      } = historyData

      // Formatear la fecha en el formato específico
      const now = new Date()
      const formattedDate = now.toLocaleDateString('es-CL', {
        day: '2-digit',
        month: '2-digit',
        year: 'numeric'
      }).replace(/\//g, '/')

      // Crear el registro de historial con el formato específico
      const historyEntry = {
        date_time: now.toISOString(),
        change_details: 'Cantidad actualizada',
        previous_values: JSON.stringify({
          amount: previousAmount,
          observaciones: `Antes del consumo del insumo ${consumedQRCode}`
        }),
        new_values: JSON.stringify({
          amount: newAmount,
          observaciones: `Después del consumo del insumo ${consumedQRCode}`,
          destino: destination
        }),
        user_name: userName,
        user_rut: userRUT,
        batch_id: batchId,
        batch_number: batchId
      }

      // Crear el registro usando el servicio de inventario existente
      const response = await this.api.post('/batch-history/', historyEntry)

      console.log('Historial del lote actualizado exitosamente:', response.data)
      return response.data

    } catch (error) {
      console.error('Error al crear historial del lote:', error)
      throw error
    }
  }

  // Método para obtener el historial formateado de un lote con los consumos
  async getBatchHistoryFormatted(batchId) {
    try {
      const history = await inventoryService.getBatchHistory(batchId)

      // Formatear el historial para mostrar en el formato específico
      return history.map(entry => {
        const date = new Date(entry.date_time)
        const formattedDate = date.toLocaleDateString('es-CL', {
          day: '2-digit',
          month: '2-digit',
          year: 'numeric'
        })

        let previousValues = {}
        let newValues = {}

        try {
          previousValues = JSON.parse(entry.previous_values || '{}')
          newValues = JSON.parse(entry.new_values || '{}')
        } catch (parseError) {
          console.warn('Error parsing values:', parseError)
        }

        return {
          ...entry,
          formatted_date: formattedDate,
          display_format: {
            date: formattedDate,
            action: entry.change_details,
            previous_amount: previousValues.cantidad || '',
            new_amount: newValues.cantidad || '',
            user_rut: entry.user_rut,
            user_name: entry.user_name,
            destination: newValues.destino || '',
            observations: newValues.observaciones || ''
          }
        }
      })
    } catch (error) {
      console.error('Error al obtener historial formateado:', error)
      throw error
    }
  }

  // Obtener información QR como alias para scanQRCode
  async getQRInfo(qrCode) {
    return this.scanQRCode(qrCode)
  }

  // Obtener ID y detalles del insumo desde su código QR
  async getSupplyIdFromQR(qrCode) {
    try {
      // Primero intentamos con la información QR completa
      const qrInfo = await this.scanQRCode(qrCode)
      
      if (qrInfo && qrInfo.supply_info && qrInfo.supply_info.id) {
        return {
          id: qrInfo.supply_info.id,
          status: qrInfo.supply_info.status || qrInfo.status,
          batch_id: qrInfo.supply_info.batch?.id,
          supply_code_id: qrInfo.supply_info.supply_code_id || qrInfo.supply_info.SupplyCode?.id
        }
      }

      // Si no funciona, intentamos con el endpoint directo de medical supplies
      const response = await this.api.get(`/medical-supplies/qr/${qrCode}`)
      if (response.data && response.data.data) {
        const supply = response.data.data
        return {
          id: supply.id,
          status: supply.status,
          batch_id: supply.batch_id,
          supply_code_id: supply.supply_code_id
        }
      }

      throw new Error('No se pudo obtener información del insumo')
    } catch (error) {
      console.error('Error al obtener ID del insumo:', error)
      throw new Error('No se pudo obtener la información del insumo para la transferencia')
    }
  }

  // Consumir insumo individual con endpoint específico y trazabilidad
  async consumeIndividualSupply(consumeData) {
    try {
      const payload = {
        qr_code: consumeData.qr_code,
        user_rut: consumeData.user_rut || this.getCurrentUserRUT(),
        user_name: consumeData.user_name || 'Encargado Bodega',
        destination_type: consumeData.destination_type,
        destination_id: consumeData.destination_id,
        notes: consumeData.notes || '',
        consumed_at: consumeData.consumed_at || new Date().toISOString()
      }

      // Usar el método principal que incluye actualización de historial
      return this.consumeSupply(payload)
    } catch (error) {
      console.error('Error consuming individual supply:', error)
      throw error
    }
  }

  // Consumir múltiples insumos individuales en lote con actualización de historial
  async bulkConsumeSupplies(bulkConsumptionData) {
    try {
      // Validar que todos sean insumos individuales
      const invalidQRs = bulkConsumptionData.qr_codes?.filter(qr => !this.isIndividualSupply(qr))
      if (invalidQRs && invalidQRs.length > 0) {
        throw new Error(`Los siguientes códigos QR no son insumos individuales: ${invalidQRs.join(', ')}`)
      }

      const results = []
      const batchUpdates = new Map()

      // Procesar cada QR individualmente
      for (const qrCode of bulkConsumptionData.qr_codes) {
        try {
          // Obtener info del insumo para conocer su lote
          const supplyInfo = await this.getQRInfo(qrCode)
          const batchId = supplyInfo.supply_info?.batch?.id

          if (batchId) {
            // Acumular cambios por lote
            if (!batchUpdates.has(batchId)) {
              batchUpdates.set(batchId, {
                batchId: batchId,
                currentAmount: supplyInfo.supply_info.batch.amount,
                consumedCount: 0,
                consumedQRCodes: []
              })
            }
            batchUpdates.get(batchId).consumedCount++
            batchUpdates.get(batchId).consumedQRCodes.push(qrCode)
          }

          // Consumir el insumo individual
          const consumptionResult = await this.api.post('/qr/consume', {
            qr_code: qrCode,
            user_rut: bulkConsumptionData.user_rut,
            user_name: bulkConsumptionData.user_name || 'Usuario',
            destination_type: bulkConsumptionData.destination_type,
            destination_id: bulkConsumptionData.destination_id,
            notes: bulkConsumptionData.notes || '',
            consumed_at: new Date().toISOString()
          })

          results.push({
            qr_code: qrCode,
            success: true,
            data: consumptionResult.data
          })

        } catch (error) {
          results.push({
            qr_code: qrCode,
            success: false,
            error: error.message
          })
        }
      }

      // Actualizar historial de lotes agrupado
      for (const [batchId, updateInfo] of batchUpdates) {
        try {
          await this.updateBatchHistoryAfterConsumption({
            batchId: updateInfo.batchId,
            previousAmount: updateInfo.currentAmount,
            newAmount: updateInfo.currentAmount - updateInfo.consumedCount,
            userRUT: bulkConsumptionData.user_rut,
            userName: bulkConsumptionData.user_name || 'Usuario',
            consumedQRCode: `Múltiples: ${updateInfo.consumedQRCodes.join(', ')}`,
            destination: `${bulkConsumptionData.destination_type} ${bulkConsumptionData.destination_id}`
          })
        } catch (historyError) {
          console.warn(`Error al actualizar historial del lote ${batchId}:`, historyError)
        }
      }

      const successfulConsumptions = results.filter(r => r.success).length

      return {
        success: true,
        successful_consumptions: successfulConsumptions,
        failed_consumptions: results.length - successfulConsumptions,
        results: results,
        batch_history_updated: batchUpdates.size,
        ui_message: `${successfulConsumptions} insumos consumidos exitosamente. Historial de ${batchUpdates.size} lote(s) actualizado.`
      }

    } catch (error) {
      console.error('Error al consumir múltiples insumos:', error)
      throw error
    }
  }

  // Verificar disponibilidad de un insumo individual para consumo
  async verifySupplyAvailability(qrCode) {
    try {
      if (!this.isIndividualSupply(qrCode)) {
        return {
          available: false,
          reason: 'not_individual_supply',
          message: 'Solo se pueden verificar insumos individuales'
        }
      }

      const response = await this.api.get(`/qr/verify/${encodeURIComponent(qrCode)}`)

      return {
        ...response.data,
        is_individual_supply: true,
        ui_recommendations: this.getAvailabilityRecommendations(response.data)
      }
    } catch (error) {
      console.error('Error al verificar disponibilidad:', error)
      throw error
    }
  }

  // ========================
  // GENERACIÓN DE QR CODES PARA INSUMOS INDIVIDUALES
  // ========================

  // Generar código QR para lote con imagen
  async generateBatchQR() {
    try {
      const response = await this.api.post('/qr/generate/batch')
      return response.data
    } catch (error) {
      console.error('Error al generar código QR de lote:', error)
      throw error
    }
  }

  // Generar código QR para insumo médico individual con imagen
  async generateSupplyQR() {
    try {
      const response = await this.api.post('/qr/generate/supply')
      return response.data
    } catch (error) {
      console.error('Error al generar código QR de insumo:', error)
      throw error
    }
  }

  // Generar múltiples códigos QR para insumos individuales
  async generateMultipleSupplyQRs(quantity = 1) {
    try {
      const qrCodes = []

      for (let i = 0; i < quantity; i++) {
        const response = await this.generateSupplyQR()
        qrCodes.push(response)
      }

      return {
        success: true,
        quantity: quantity,
        qr_codes: qrCodes,
        message: `${quantity} códigos QR de insumos generados exitosamente`
      }
    } catch (error) {
      console.error('Error al generar múltiples códigos QR:', error)
      throw error
    }
  }

  // ========================
  // UTILIDADES ESPECÍFICAS PARA INSUMOS INDIVIDUALES
  // ========================

  // Verificar si un código QR corresponde a un insumo individual
  isIndividualSupply(qrCode) {
    return this.getQRType(qrCode) === 'SUPPLY'
  }

  // Verificar si un código QR corresponde a un lote
  isBatchQR(qrCode) {
    return this.getQRType(qrCode) === 'BATCH'
  }

  // Obtener tipo de QR (SUPPLY o BATCH)
  getQRType(qrCode) {
    if (!qrCode || typeof qrCode !== 'string') return null

    if (qrCode.startsWith('SUPPLY_')) return 'SUPPLY'
    if (qrCode.startsWith('BATCH_')) return 'BATCH'

    return null
  }

  // Validar formato de código QR con mensajes específicos
  isValidQRFormat(qrCode) {
    if (!qrCode || typeof qrCode !== 'string') return false

    // Verificar que tenga el formato esperado: PREFIX_TIMESTAMP_RANDOM
    const qrPattern = /^(BATCH|SUPPLY)_\d+_[a-f0-9]+$/i
    return qrPattern.test(qrCode)
  }

  // Extraer información básica del código QR con detalles del tipo
  parseQRCode(qrCode) {
    if (!this.isValidQRFormat(qrCode)) return null

    const parts = qrCode.split('_')
    if (parts.length !== 3) return null

    const type = parts[0].toLowerCase()
    const timestamp = parseInt(parts[1])

    return {
      type: type,
      type_label: type === 'supply' ? 'Insumo Individual' : 'Lote',
      timestamp: timestamp,
      randomId: parts[2],
      date: new Date(timestamp * 1000),
      is_individual_supply: type === 'supply',
      is_batch: type === 'batch',
      priority_level: type === 'supply' ? 'high' : 'medium'
    }
  }

  // Formatear información de QR para mostrar con enfoque en insumos individuales
  formatQRInfo(qrInfo) {
    if (!qrInfo) return null

    const formatted = {
      ...qrInfo,
      typeLabel: this.getTypeLabel(qrInfo.type),
      statusLabel: this.getStatusLabel(qrInfo),
      statusColor: this.getStatusColor(qrInfo),
      is_recommended_scan: qrInfo.type === 'medical_supply',
      ui_priority: qrInfo.type === 'medical_supply' ? 'high' : 'medium'
    }

    return formatted
  }

  // Obtener etiqueta del tipo de QR con énfasis en insumos individuales
  getTypeLabel(type) {
    const labels = {
      'batch': 'Lote (Solo Consulta)',
      'medical_supply': 'Insumo Individual',
      'supply': 'Insumo Individual'
    }
    return labels[type] || type
  }

  // Obtener etiqueta de estado con información detallada
  getStatusLabel(qrInfo) {
    if (!qrInfo) return 'Desconocido'

    if (qrInfo.type === 'batch') {
      const batchStatus = qrInfo.batch_status
      if (batchStatus && batchStatus.available_supplies === 0) {
        return 'Lote Agotado'
      } else if (batchStatus && batchStatus.available_supplies < 5) {
        return 'Lote con Stock Bajo'
      } else {
        return 'Lote Disponible'
      }
    } else if (qrInfo.type === 'medical_supply') {
      if (qrInfo.is_consumed) {
        return 'Insumo Consumido'
      } else if (qrInfo.available_for_use === false) {
        return 'Insumo No Disponible'
      } else {
        return 'Insumo Disponible'
      }
    }

    return 'Estado Desconocido'
  }

  // Obtener color de estado con códigos específicos
  getStatusColor(qrInfo) {
    if (!qrInfo) return 'gray'

    if (qrInfo.type === 'medical_supply') {
      if (qrInfo.is_consumed) return 'red'
      if (qrInfo.available_for_use === false) return 'orange'
      return 'green'
    } else if (qrInfo.type === 'batch') {
      const batchStatus = qrInfo.batch_status
      if (batchStatus && batchStatus.available_supplies === 0) return 'red'
      if (batchStatus && batchStatus.available_supplies < 5) return 'yellow'
      return 'blue'
    }

    return 'gray'
  }

  // Obtener recomendaciones de disponibilidad para insumos individuales
  getAvailabilityRecommendations(availabilityData) {
    const recommendations = []

    if (availabilityData.available && availabilityData.is_individual_supply) {
      recommendations.push({
        type: 'success',
        message: 'Insumo listo para consumir',
        action: 'consume'
      })
    }

    if (!availabilityData.available && availabilityData.reason === 'already_consumed') {
      recommendations.push({
        type: 'info',
        message: 'Este insumo ya fue consumido. Escanee otro del mismo lote.',
        action: 'scan_another'
      })
    }

    if (availabilityData.batch_info && availabilityData.batch_info.low_stock) {
      recommendations.push({
        type: 'warning',
        message: 'Stock bajo en este lote. Considere solicitar reposición.',
        action: 'alert_low_stock'
      })
    }

    return recommendations
  }

  // ========================
  // MANEJO DE IMÁGENES QR
  // ========================

  // Obtener URL de imagen QR
  getQRImageUrl(qrCode) {
    return `${API_BASE_URL}/qr/image/${encodeURIComponent(qrCode)}`
  }

  // Obtener URL de descarga de imagen QR
  getQRDownloadUrl(qrCode, resolution = 'normal') {
    return `${API_BASE_URL}/qr/download/${encodeURIComponent(qrCode)}?resolution=${resolution}`
  }

  // Descargar imagen QR con nombre descriptivo según el tipo
  async downloadQRImage(qrCode, resolution = 'normal') {
    try {
      const response = await this.api.get(`/qr/download/${encodeURIComponent(qrCode)}`, {
        params: { resolution },
        responseType: 'blob'
      })

      // Crear URL para descarga
      const url = window.URL.createObjectURL(new Blob([response.data]))
      const link = document.createElement('a')
      link.href = url

      // Generar nombre de archivo descriptivo
      const qrType = this.getQRType(qrCode)
      const prefix = qrType === 'SUPPLY' ? 'insumo' : qrType === 'BATCH' ? 'lote' : 'qr'
      const suffix = resolution === 'high' ? '_hd' : ''
      const filename = `${prefix}_${qrCode}${suffix}.png`

      link.setAttribute('download', filename)
      document.body.appendChild(link)
      link.click()
      link.remove()
      window.URL.revokeObjectURL(url)

      return true
    } catch (error) {
      console.error('Error al descargar imagen QR:', error)
      throw error
    }
  }

  // Verificar si una imagen QR existe
  async checkQRImageExists(qrCode) {
    try {
      const response = await fetch(this.getQRImageUrl(qrCode), { method: 'HEAD' })
      return response.ok
    } catch (error) {
      return false
    }
  }

  // ========================
  // ADMINISTRACIÓN Y ESTADÍSTICAS
  // ========================

  // Sincronizar cantidades de lotes con productos individuales
  async syncBatchAmounts() {
    try {
      const response = await this.api.post('/qr/sync/batch-amounts')
      return response.data
    } catch (error) {
      console.error('Error al sincronizar cantidades:', error)
      throw error
    }
  }

  // Obtener estadísticas generales con enfoque en insumos individuales
  async getQRStats() {
    try {
      const response = await this.api.get('/qr/stats')

      if (response.data) {
        // Añadir cálculos adicionales útiles para el frontend
        const data = response.data

        response.data.individual_supplies_stats = {
          total_supplies: data.total_supplies,
          consumed_supplies: data.consumed_supplies,
          available_supplies: data.available_supplies,
          consumption_rate: data.consumption_rate,
          average_supplies_per_batch: data.total_batches > 0 ?
            Math.round(data.total_supplies / data.total_batches) : 0
        }

        response.data.ui_insights = {
          primary_metric: `${data.available_supplies} insumos disponibles`,
          consumption_trend: data.consumption_rate > 50 ? 'high' :
            data.consumption_rate > 25 ? 'medium' : 'low',
          needs_attention: data.available_supplies < 100
        }
      }

      return response.data
    } catch (error) {
      console.error('Error al obtener estadísticas:', error)
      throw error
    }
  }

  // ========================
  // UTILIDADES ADICIONALES PARA UX
  // ========================


  // Transferir un insumo individual por su código QR con trazabilidad
  async transferSupply(transferData) {
    try {
      // Validar que sea un insumo individual
      if (!this.isIndividualSupply(transferData.qr_code)) {
        throw new Error('Solo se pueden transferir insumos individuales. Use códigos QR que comiencen con SUPPLY_')
      }

      // Usar el nuevo endpoint específico para transferencias
      const response = await this.api.post('/qr/transfer', {
        qr_code: transferData.qr_code,
        user_rut: transferData.user_rut,
        receiver_rut: transferData.receiver_rut,
        destination_type: transferData.destination_type,
        destination_id: transferData.destination_id,
        notes: transferData.notes
      })

      // Si la respuesta viene en formato {success: true, data: {...}}
      let result = response.data
      if (result && result.success && result.data) {
        result = result.data
      }

      return {
        ...result,
        ui_message: `Insumo transferido exitosamente. Estado cambiado a "${result.new_status}". La trazabilidad ha sido actualizada.`,
        next_actions: [
          { label: 'Ver trazabilidad', action: 'view_traceability' },
          { label: 'Transferir otro', action: 'transfer_next' },
          { label: 'Ver inventario', action: 'view_inventory' }
        ],
        transfer_mode: true,
        status_change: result.status_change || {
          from: result.old_status,
          to: result.new_status
        }
      }
    } catch (error) {
      console.error('Error al transferir insumo:', error)

      // Mejorar mensajes de error para transferencias
      if (error.response?.data?.message?.includes('ya ha sido consumido')) {
        throw new Error('Este insumo ya ha sido consumido y no puede ser transferido.')
      } else if (error.response?.data?.message?.includes('no encontrado')) {
        throw new Error('Insumo no encontrado. Verifique que el código QR sea correcto.')
      } else if (error.response?.data?.message?.includes('no está disponible para transferencia')) {
        throw new Error('El insumo no está disponible para transferencia.')
      } else if (error.response?.data?.message?.includes('recepcionado')) {
        throw new Error('Este insumo tiene estado "recepcionado" y solo puede ser consumido, no transferido.')
      }

      throw error
    }
  }




  // Recepcionar un insumo que está en estado "en_camino_a_pabellon"
  async receiveSupply(qrCode, userRUT, destinationType, destinationID, notes = '') {
    try {
      // Validar que el QR code tenga el formato correcto
      if (!qrCode || !qrCode.startsWith('SUPPLY_')) {
        throw new Error('El código QR debe ser de un insumo individual (SUPPLY_)')
      }

      const response = await this.api.post('/qr/receive', {
        qr_code: qrCode,
        user_rut: userRUT,
        destination_type: destinationType,
        destination_id: destinationID,
        notes: notes
      })

      if (response.data.success) {
        return {
          success: true,
          data: response.data.data,
          message: response.data.message
        }
      } else {
        throw new Error(response.data.error || 'Error al recepcionar el insumo')
      }
    } catch (error) {
      console.error('Error receiving supply:', error)
      throw error
    }
  }
}

// Crear y exportar instancia singleton
const qrService = new QRService()
export default qrService