import axios from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1'

class QRService {
  constructor() {
    this.api = axios.create({
      baseURL: API_BASE_URL,
      headers: {
        'Content-Type': 'application/json'
      }
    })
  }

  // ========================
  // ESCANEO Y VALIDACIÓN - ENFOCADO EN INSUMOS INDIVIDUALES
  // ========================

  // Escanear un código QR con prioridad para insumos individuales
  async scanQRCode(qrCode) {
    try {
      const response = await this.api.get(`/qr/scan/${encodeURIComponent(qrCode)}`)
      
      // Añadir información adicional para mejor UX
      if (response.data) {
        response.data.scan_timestamp = new Date()
        response.data.is_individual_supply = response.data.type === 'medical_supply'
        response.data.is_batch = response.data.type === 'batch'
        
        // Para insumos individuales, verificar estado de consumo
        if (response.data.type === 'medical_supply') {
          response.data.can_be_consumed = !response.data.is_consumed && 
                                         response.data.available_for_use !== false
          
          // Añadir información del lote si está disponible
          if (response.data.batch_status) {
            response.data.batch_context = {
              has_stock: response.data.batch_status.current_amount > 0,
              low_stock: response.data.batch_status.current_amount < 10,
              batch_id: response.data.batch_status.batch_id
            }
          }
        }
        
        // Para lotes, añadir información sobre insumos individuales
        if (response.data.type === 'batch') {
          if (response.data.batch_status) {
            response.data.individual_supplies_info = {
              total: response.data.batch_status.total_individual_supplies,
              available: response.data.batch_status.available_supplies,
              consumed: response.data.batch_status.consumed_supplies,
              consumption_rate: response.data.batch_status.total_individual_supplies > 0 
                ? Math.round((response.data.batch_status.consumed_supplies / response.data.batch_status.total_individual_supplies) * 100)
                : 0
            }
          }
        }
      }
      
      return response.data
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

  // Obtener historial de un insumo individual por código QR
  async getSupplyHistory(qrCode) {
    try {
      if (!this.isIndividualSupply(qrCode)) {
        throw new Error('Este endpoint está diseñado para insumos individuales. Use códigos QR que comiencen con SUPPLY_')
      }
      
      const response = await this.api.get(`/qr/history/${encodeURIComponent(qrCode)}`)
      
      // Ordenar historial por fecha descendente
      if (response.data && Array.isArray(response.data)) {
        response.data.sort((a, b) => new Date(b.date_time) - new Date(a.date_time))
      }
      
      return response.data
    } catch (error) {
      console.error('Error al obtener historial:', error)
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
  // CONSUMO DE INSUMOS INDIVIDUALES
  // ========================

  // Consumir un insumo individual por su código QR
  async consumeSupply(consumptionData) {
    try {
      // Validar que sea un insumo individual
      if (!this.isIndividualSupply(consumptionData.qr_code)) {
        throw new Error('Solo se pueden consumir insumos individuales. Use códigos QR que comiencen con SUPPLY_')
      }
      
      const response = await this.api.post('/qr/consume', {
        ...consumptionData,
        consumption_timestamp: new Date().toISOString()
      })
      
      return {
        ...response.data,
        ui_message: 'Insumo individual consumido exitosamente',
        next_actions: [
          { label: 'Ver historial', action: 'view_history' },
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

  // Consumir múltiples insumos individuales en lote
  async bulkConsumeSupplies(bulkConsumptionData) {
    try {
      // Validar que todos sean insumos individuales
      const invalidQRs = bulkConsumptionData.qr_codes?.filter(qr => !this.isIndividualSupply(qr))
      if (invalidQRs && invalidQRs.length > 0) {
        throw new Error(`Los siguientes códigos QR no son insumos individuales: ${invalidQRs.join(', ')}`)
      }
      
      const response = await this.api.post('/qr/consume/bulk', {
        ...bulkConsumptionData,
        consumption_timestamp: new Date().toISOString()
      })
      
      return {
        ...response.data,
        ui_message: `${response.data.successful_consumptions || 0} insumos individuales consumidos exitosamente`
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
      priority_level: type === 'supply' ? 'high' : 'medium' // Prioridad para insumos individuales
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
      is_recommended_scan: qrInfo.type === 'medical_supply', // Prioridad para insumos individuales
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
        return 'Insumo Listo para Usar'
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
      return 'green' // Listo para usar
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

  // Generar sugerencias de QR para testing (solo desarrollo)
  generateTestQRCodes() {
    const timestamp = Math.floor(Date.now() / 1000)
    const randomId = Math.random().toString(16).substr(2, 8)
    
    return {
      supply: `SUPPLY_${timestamp}_${randomId}`,
      batch: `BATCH_${timestamp}_${randomId}`
    }
  }

  // Validar entrada de usuario con sugerencias
  validateUserInput(input) {
    if (!input || typeof input !== 'string') {
      return {
        valid: false,
        error: 'Por favor ingrese un código QR válido',
        suggestions: ['Escanee con la cámara o pegue el código QR']
      }
    }
    
    const trimmed = input.trim().toUpperCase()
    
    if (!this.isValidQRFormat(trimmed)) {
      const suggestions = ['Formato válido: SUPPLY_timestamp_random', 'Formato válido: BATCH_timestamp_random']
      
      if (trimmed.includes('SUPPLY') || trimmed.includes('BATCH')) {
        suggestions.unshift('Verifique que el código esté completo')
      } else {
        suggestions.unshift('Los códigos QR deben comenzar con SUPPLY_ o BATCH_')
      }
      
      return {
        valid: false,
        error: 'Formato de código QR inválido',
        suggestions
      }
    }
    
    return {
      valid: true,
      formatted: trimmed,
      type: this.getQRType(trimmed),
      is_individual_supply: this.isIndividualSupply(trimmed)
    }
  }
}

// Crear y exportar instancia singleton
const qrService = new QRService()
export default qrService