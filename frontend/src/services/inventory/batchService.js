import axios from 'axios'
import { getApiBaseUrl } from '@/config/api.js'

const API_BASE_URL = getApiBaseUrl()

class BatchService {
  constructor() {
    this.api = axios.create({
      baseURL: API_BASE_URL,
      headers: {
        'Content-Type': 'application/json'
      }
    })

    // Interceptor para agregar token de autenticación
    this.api.interceptors.request.use(
      (config) => {
        const token = localStorage.getItem('auth_token')
        if (token) {
          config.headers.Authorization = `Bearer ${token}`
        }
        return config
      },
      (error) => {
        return Promise.reject(error)
      }
    )
  }

  // ========================
  // CRUD BÁSICO
  // ========================

  async createBatch(batchData) {
    try {
      const response = await this.api.post('/batches', batchData)
      return response.data
    } catch (error) {
      console.error('Error al crear lote:', error)
      throw error
    }
  }

  async getAllBatches() {
    try {
      const response = await this.api.get('/batches')
      return response.data.data || response.data || []
    } catch (error) {
      console.error('Error al obtener lotes:', error)
      throw error
    }
  }

  async getBatchByID(id) {
    try {
      const response = await this.api.get(`/batches/${id}`)
      return response.data.data || response.data
    } catch (error) {
      console.error('Error al obtener lote:', error)
      throw error
    }
  }

  async updateBatch(id, batchData) {
    try {
      const response = await this.api.put(`/batches/${id}`, batchData)
      return response.data
    } catch (error) {
      console.error('Error al actualizar lote:', error)
      throw error
    }
  }

  async deleteBatch(id) {
    try {
      const response = await this.api.delete(`/batches/${id}`)
      return response.data
    } catch (error) {
      console.error('Error al eliminar lote:', error)
      throw error
    }
  }

  // ========================
  // RUTAS MEJORADAS
  // ========================

  async createBatchWithIndividualSupplies(batchData) {
    try {
      const response = await this.api.post('/batches/create-with-supplies', batchData)
      return response.data
    } catch (error) {
      console.error('Error al crear lote con insumos:', error)
      throw error
    }
  }

  async getBatchWithSupplyInfo(id) {
    try {
      const response = await this.api.get(`/batches/${id}/with-supplies`)
      return response.data.data || response.data
    } catch (error) {
      console.error('Error al obtener lote con información de insumos:', error)
      throw error
    }
  }

  async getBatchByQR(qrcode) {
    try {
      const response = await this.api.get(`/batches/qr/${qrcode}`)
      return response.data.data || response.data
    } catch (error) {
      console.error('Error al obtener lote por QR:', error)
      throw error
    }
  }

  async getBatchesNeedingSync() {
    try {
      const response = await this.api.get('/batches/sync/needed')
      return response.data.data || response.data || []
    } catch (error) {
      console.error('Error al obtener lotes que necesitan sincronización:', error)
      throw error
    }
  }

  // ========================
  // ACTUALIZACIÓN Y MANTENIMIENTO
  // ========================

  async updateBatchAmount(id, amount) {
    try {
      const response = await this.api.patch(`/batches/${id}/amount`, { amount })
      return response.data
    } catch (error) {
      console.error('Error al actualizar cantidad del lote:', error)
      throw error
    }
  }

  async syncAllBatchAmounts() {
    try {
      const response = await this.api.post('/batches/sync/all')
      return response.data
    } catch (error) {
      console.error('Error al sincronizar cantidades de lotes:', error)
      throw error
    }
  }

  // ========================
  // ALERTAS
  // ========================

  async checkLowStockAlert(id) {
    try {
      const response = await this.api.post(`/batches/${id}/check-low-stock`)
      return response.data
    } catch (error) {
      console.error('Error al verificar alerta de bajo stock:', error)
      throw error
    }
  }

  async checkExpirationAlert(id, days = null) {
    try {
      const params = days ? { days } : {}
      const response = await this.api.post(`/batches/${id}/check-expiration`, null, { params })
      return response.data
    } catch (error) {
      console.error('Error al verificar alerta de expiración:', error)
      throw error
    }
  }

  async checkAllBatchesExpiration() {
    try {
      const response = await this.api.post('/batches/check-all-expiration')
      return response.data
    } catch (error) {
      console.error('Error al verificar alertas de vencimiento para todos los lotes:', error)
      throw error
    }
  }

  async checkAllBatchesLowStock() {
    try {
      const response = await this.api.post('/batches/check-all-low-stock')
      return response.data
    } catch (error) {
      console.error('Error al verificar alertas de stock bajo para todos los lotes:', error)
      throw error
    }
  }

  async getLowStockSummary() {
    try {
      const response = await this.api.get('/batches/low-stock-summary')
      return response.data.data || response.data || []
    } catch (error) {
      console.error('Error al obtener resumen de stock bajo:', error)
      throw error
    }
  }
}

export default new BatchService()
