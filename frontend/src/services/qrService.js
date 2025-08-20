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

  // Obtener URL de imagen QR
  getQRImageUrl(qrCode) {
    return `${API_BASE_URL}/qr/image/${encodeURIComponent(qrCode)}`
  }

  // Obtener URL de descarga de imagen QR
  getQRDownloadUrl(qrCode) {
    return `${API_BASE_URL}/qr/download/${encodeURIComponent(qrCode)}`
  }

  // Descargar imagen QR directamente
  async downloadQRImage(qrCode, filename = null) {
    try {
      const response = await this.api.get(`/qr/download/${encodeURIComponent(qrCode)}`, {
        responseType: 'blob'
      })
      
      // Crear enlace de descarga
      const blob = new Blob([response.data], { type: 'image/png' })
      const url = window.URL.createObjectURL(blob)
      
      const link = document.createElement('a')
      link.href = url
      link.download = filename || `qr_${qrCode}.png`
      document.body.appendChild(link)
      link.click()
      
      // Limpiar
      document.body.removeChild(link)
      window.URL.revokeObjectURL(url)
      
      return true
    } catch (error) {
      console.error('Error al descargar imagen QR:', error)
      throw error
    }
  }
}

export default new QRService()