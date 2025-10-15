import axios from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1'

class SupplyHistoryService {
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

  async createSupplyHistory(supplyHistoryData) {
    try {
      const response = await this.api.post('/supply-history/', supplyHistoryData)
      return response.data
    } catch (error) {
      console.error('Error al crear historial de insumo:', error)
      throw error
    }
  }

  async getAllSupplyHistory() {
    try {
      const response = await this.api.get('/supply-history/')
      return response.data.data || response.data || []
    } catch (error) {
      console.error('Error al obtener historiales de insumos:', error)
      throw error
    }
  }

  async getAllSupplyHistoryWithDetails() {
    try {
      const response = await this.api.get('/supply-history/with-details')
      return response.data.data || response.data || []
    } catch (error) {
      console.error('Error al obtener historiales de insumos con detalles:', error)
      throw error
    }
  }

  async getSupplyHistoryByID(id) {
    try {
      const response = await this.api.get(`/supply-history/${id}/`)
      return response.data.data || response.data
    } catch (error) {
      console.error('Error al obtener historial de insumo:', error)
      throw error
    }
  }

  async updateSupplyHistory(id, supplyHistoryData) {
    try {
      const response = await this.api.put(`/supply-history/${id}/`, supplyHistoryData)
      return response.data
    } catch (error) {
      console.error('Error al actualizar historial de insumo:', error)
      throw error
    }
  }

  async deleteSupplyHistory(id) {
    try {
      const response = await this.api.delete(`/supply-history/${id}/`)
      return response.data
    } catch (error) {
      console.error('Error al eliminar historial de insumo:', error)
      throw error
    }
  }
}

export default new SupplyHistoryService()
