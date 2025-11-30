import axios from 'axios'
import { getApiBaseUrl } from '@/config/api.js'

const API_BASE_URL = getApiBaseUrl()

class SupplierConfigService {
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

  async createSupplierConfig(configData) {
    try {
      const response = await this.api.post('/supplier-configs/', configData)
      return response.data.data || response.data
    } catch (error) {
      console.error('Error en createSupplierConfig:', error)
      throw error
    }
  }

  async getAllSupplierConfigs() {
    try {
      const response = await this.api.get('/supplier-configs/')
      return response.data.data || response.data || []
    } catch (error) {
      console.error('Error en getAllSupplierConfigs:', error)
      throw error
    }
  }

  async getSupplierConfig(supplierName) {
    try {
      const response = await this.api.get(`/supplier-configs/${encodeURIComponent(supplierName)}`)
      return response.data.data || response.data
    } catch (error) {
      console.error('Error en getSupplierConfig:', error)
      throw error
    }
  }

  async updateSupplierConfig(supplierName, configData) {
    try {
      const response = await this.api.put(`/supplier-configs/${encodeURIComponent(supplierName)}`, configData)
      return response.data.data || response.data
    } catch (error) {
      console.error('Error en updateSupplierConfig:', error)
      throw error
    }
  }

  async deleteSupplierConfig(supplierName) {
    try {
      const response = await this.api.delete(`/supplier-configs/${encodeURIComponent(supplierName)}`)
      return response.data
    } catch (error) {
      console.error('Error en deleteSupplierConfig:', error)
      throw error
    }
  }
}

export default new SupplierConfigService()
