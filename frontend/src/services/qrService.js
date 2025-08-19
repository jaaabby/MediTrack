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

  // Generar código QR para lote
  async generateBatchQR() {
    try {
      const response = await this.api.post('/qr/generate/batch')
      return response.data
    } catch (error) {
      console.error('Error al generar código QR de lote:', error)
      throw error
    }
  }

  // Generar código QR para insumo médico
  async generateSupplyQR() {
    try {
      const response = await this.api.post('/qr/generate/supply')
      return response.data
    } catch (error) {
      console.error('Error al generar código QR de insumo:', error)
      throw error
    }
  }
}

export default new QRService()