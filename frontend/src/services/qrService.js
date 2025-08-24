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
  // ESCANEO Y VALIDACIÓN
  // ========================

  // Escanear un código QR y obtener toda su información
  async scanQRCode(qrCode) {
    try {
      const response = await this.api.get(`/qr/scan/${encodeURIComponent(qrCode)}`)
      return response.data
    } catch (error) {
      console.error('Error al escanear código QR:', error)
      throw error
    }
  }

  // Validar si un código QR es válido
  async validateQRCode(qrCode) {
    try {
      const response = await this.api.get(`/qr/validate/${encodeURIComponent(qrCode)}`)
      return response.data
    } catch (error) {
      console.error('Error al validar código QR:', error)
      throw error
    }
  }

  // Obtener historial de un insumo por código QR
  async getSupplyHistory(qrCode) {
    try {
      const response = await this.api.get(`/qr/history/${encodeURIComponent(qrCode)}`)
      return response.data
    } catch (error) {
      console.error('Error al obtener historial:', error)
      throw error
    }
  }

  // Obtener información detallada de un insumo con datos del lote
  async getSupplyDetails(qrCode) {
    try {
      const response = await this.api.get(`/qr/details/${encodeURIComponent(qrCode)}`)
      return response.data
    } catch (error) {
      console.error('Error al obtener detalles del insumo:', error)
      throw error
    }
  }

  // ========================
  // GENERACIÓN DE QR CODES
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

  // Generar código QR para insumo médico con imagen
  async generateSupplyQR() {
    try {
      const response = await this.api.post('/qr/generate/supply')
      return response.data
    } catch (error) {
      console.error('Error al generar código QR de insumo:', error)
      throw error
    }
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

  // Descargar imagen QR
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
      
      const filename = resolution === 'high' ? `${qrCode}_qr_hd.png` : `${qrCode}_qr.png`
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
  // CONSUMO DE PRODUCTOS
  // ========================

  // Consumir un insumo por su código QR
  async consumeSupply(consumptionData) {
    try {
      const response = await this.api.post('/qr/consume', consumptionData)
      return response.data
    } catch (error) {
      console.error('Error al consumir insumo:', error)
      throw error
    }
  }

  // Consumir múltiples insumos en lote
  async bulkConsumeSupplies(bulkConsumptionData) {
    try {
      const response = await this.api.post('/qr/consume/bulk', bulkConsumptionData)
      return response.data
    } catch (error) {
      console.error('Error al consumir múltiples insumos:', error)
      throw error
    }
  }

  // Verificar disponibilidad de un insumo para consumo
  async verifySupplyAvailability(qrCode) {
    try {
      const response = await this.api.get(`/qr/verify/${encodeURIComponent(qrCode)}`)
      return response.data
    } catch (error) {
      console.error('Error al verificar disponibilidad:', error)
      throw error
    }
  }

  // ========================
  // ADMINISTRACIÓN
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

  // Obtener estadísticas generales de uso de QR codes
  async getQRStats() {
    try {
      const response = await this.api.get('/qr/stats')
      return response.data
    } catch (error) {
      console.error('Error al obtener estadísticas:', error)
      throw error
    }
  }

  // ========================
  // UTILIDADES
  // ========================

  // Validar formato de código QR
  isValidQRFormat(qrCode) {
    if (!qrCode || typeof qrCode !== 'string') return false
    
    // Verificar que tenga el formato esperado: PREFIX_TIMESTAMP_RANDOM
    const qrPattern = /^(BATCH|SUPPLY)_\d+_[a-f0-9]+$/i
    return qrPattern.test(qrCode)
  }

  // Extraer información básica del código QR
  parseQRCode(qrCode) {
    if (!this.isValidQRFormat(qrCode)) return null
    
    const parts = qrCode.split('_')
    if (parts.length !== 3) return null
    
    return {
      type: parts[0].toLowerCase(),
      timestamp: parseInt(parts[1]),
      randomId: parts[2],
      date: new Date(parseInt(parts[1]) * 1000)
    }
  }

  // Formatear información de QR para mostrar
  formatQRInfo(qrInfo) {
    if (!qrInfo) return null
    
    const formatted = {
      ...qrInfo,
      typeLabel: this.getTypeLabel(qrInfo.type),
      statusLabel: this.getStatusLabel(qrInfo),
      statusColor: this.getStatusColor(qrInfo)
    }
    
    return formatted
  }

  // Obtener etiqueta del tipo de QR
  getTypeLabel(type) {
    const labels = {
      'batch': 'Lote',
      'medical_supply': 'Insumo Individual'
    }
    return labels[type] || type
  }

  // Obtener etiqueta de estado
  getStatusLabel(qrInfo) {
    if (!qrInfo) return 'Desconocido'
    
    if (qrInfo.type === 'batch') {
      const batchStatus = qrInfo.batch_status
      if (batchStatus && batchStatus.available_supplies === 0) {
        return 'Agotado'
      } else if (batchStatus && batchStatus.available_supplies < 5) {
        return 'Stock Bajo'
      } else {
        return 'Disponible'
      }
    } else if (qrInfo.type === 'medical_supply') {
      return qrInfo.is_consumed ? 'Consumido' : 'Disponible'
    }
    
    return 'Desconocido'
  }

  // Obtener color de estado
  getStatusColor(qrInfo) {
    if (!qrInfo) return 'gray'
    
    if (qrInfo.type === 'batch') {
      const batchStatus = qrInfo.batch_status
      if (batchStatus && batchStatus.available_supplies === 0) {
        return 'red'
      } else if (batchStatus && batchStatus.available_supplies < 5) {
        return 'yellow'
      } else {
        return 'green'
      }
    } else if (qrInfo.type === 'medical_supply') {
      return qrInfo.is_consumed ? 'red' : 'green'
    }
    
    return 'gray'
  }
}

export default new QRService()