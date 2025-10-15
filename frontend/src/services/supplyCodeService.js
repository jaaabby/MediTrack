import axios from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1'

class SupplyCodeService {
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

  async createSupplyCode(supplyCodeData) {
    try {
      const response = await this.api.post('/supply-codes/', supplyCodeData)
      return response.data
    } catch (error) {
      console.error('Error al crear código de insumo:', error)
      throw error
    }
  }

  async getAllSupplyCodes() {
    try {
      const response = await this.api.get('/supply-codes/')
      return response.data.data || response.data || []
    } catch (error) {
      console.error('Error al obtener códigos de insumos:', error)
      throw error
    }
  }

  async getSupplyCodeByID(id) {
    try {
      const response = await this.api.get(`/supply-codes/${id}/`)
      return response.data.data || response.data
    } catch (error) {
      console.error('Error al obtener código de insumo:', error)
      throw error
    }
  }

  async updateSupplyCode(id, supplyCodeData) {
    try {
      const response = await this.api.put(`/supply-codes/${id}/`, supplyCodeData)
      return response.data
    } catch (error) {
      console.error('Error al actualizar código de insumo:', error)
      throw error
    }
  }

  async deleteSupplyCode(id) {
    try {
      const response = await this.api.delete(`/supply-codes/${id}/`)
      return response.data
    } catch (error) {
      console.error('Error al eliminar código de insumo:', error)
      throw error
    }
  }
}

export default new SupplyCodeService()
