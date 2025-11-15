import axios from 'axios'
import { getApiBaseUrl } from '@/config/api.js'

const API_BASE_URL = getApiBaseUrl()

class BatchHistoryService {
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

  async createBatchHistory(batchHistoryData) {
    try {
      const response = await this.api.post('/batch-history', batchHistoryData)
      return response.data
    } catch (error) {
      console.error('Error al crear historial de lote:', error)
      throw error
    }
  }

  async getAllBatchHistories() {
    try {
      const response = await this.api.get('/batch-history')
      return response.data.data || response.data || []
    } catch (error) {
      console.error('Error al obtener historiales de lotes:', error)
      throw error
    }
  }

  async getBatchHistoryByID(id) {
    try {
      const response = await this.api.get(`/batch-history/${id}`)
      return response.data.data || response.data
    } catch (error) {
      console.error('Error al obtener historial de lote:', error)
      throw error
    }
  }

  async updateBatchHistory(id, batchHistoryData) {
    try {
      const response = await this.api.put(`/batch-history/${id}`, batchHistoryData)
      return response.data
    } catch (error) {
      console.error('Error al actualizar historial de lote:', error)
      throw error
    }
  }

  async deleteBatchHistory(id) {
    try {
      const response = await this.api.delete(`/batch-history/${id}`)
      return response.data
    } catch (error) {
      console.error('Error al eliminar historial de lote:', error)
      throw error
    }
  }

  // ========================
  // BÚSQUEDA Y CONSULTAS
  // ========================

  async searchBatchHistoryByBatchNumber(batchNumber) {
    try {
      const response = await this.api.get(`/batch-history/search/${batchNumber}`)
      return response.data.data || response.data || []
    } catch (error) {
      console.error('Error al buscar historial por número de lote:', error)
      throw error
    }
  }

  async getAllBatchHistoriesWithDetails() {
    try {
      const response = await this.api.get('/batch-history/details')
      return response.data.data || response.data || []
    } catch (error) {
      console.error('Error al obtener historiales con detalles:', error)
      throw error
    }
  }
}

export default new BatchHistoryService()
