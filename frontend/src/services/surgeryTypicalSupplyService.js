import axios from 'axios'
import { getApiBaseUrl } from '@/config/api.js'

const API_BASE_URL = getApiBaseUrl()

class SurgeryTypicalSupplyService {
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

  async getTypicalSuppliesBySurgeryId(surgeryId) {
    try {
      const response = await this.api.get(`/surgery-typical-supplies/surgery/${surgeryId}`)
      return response.data.data?.typical_supplies || response.data.typical_supplies || []
    } catch (error) {
      console.error('Error en getTypicalSuppliesBySurgeryId:', error)
      throw error
    }
  }

  async getSurgeriesBySupplyCode(supplyCode) {
    try {
      const response = await this.api.get(`/surgery-typical-supplies/supply/${supplyCode}`)
      return response.data.data?.surgeries || response.data.surgeries || []
    } catch (error) {
      console.error('Error en getSurgeriesBySupplyCode:', error)
      throw error
    }
  }

  async createTypicalSupply(typicalSupplyData) {
    try {
      const response = await this.api.post('/surgery-typical-supplies/', typicalSupplyData)
      return response.data.data || response.data
    } catch (error) {
      console.error('Error en createTypicalSupply:', error)
      throw error
    }
  }

  async bulkCreateTypicalSupplies(surgeryId, typicalSupplies) {
    try {
      const response = await this.api.post(`/surgery-typical-supplies/surgery/${surgeryId}/bulk`, typicalSupplies)
      return response.data.data || response.data
    } catch (error) {
      console.error('Error en bulkCreateTypicalSupplies:', error)
      throw error
    }
  }

  async updateTypicalSupply(id, typicalSupplyData) {
    try {
      const response = await this.api.put(`/surgery-typical-supplies/${id}`, typicalSupplyData)
      return response.data.data || response.data
    } catch (error) {
      console.error('Error en updateTypicalSupply:', error)
      throw error
    }
  }

  async deleteTypicalSupply(id) {
    try {
      const response = await this.api.delete(`/surgery-typical-supplies/${id}`)
      return response.data
    } catch (error) {
      console.error('Error en deleteTypicalSupply:', error)
      throw error
    }
  }

  async getTypicalSuppliesCount() {
    try {
      const response = await this.api.get('/surgery-typical-supplies/count')
      return response.data.data?.count || response.data.count || 0
    } catch (error) {
      console.error('Error en getTypicalSuppliesCount:', error)
      throw error
    }
  }
}

export default new SurgeryTypicalSupplyService()

