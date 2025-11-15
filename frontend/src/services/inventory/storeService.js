import axios from 'axios'
import { getApiBaseUrl } from '@/config/api.js'

const API_BASE_URL = getApiBaseUrl()

class StoreService {
  constructor() {
    this.api = axios.create({
      baseURL: API_BASE_URL,
      headers: {
        'Content-Type': 'application/json'
      }
    })

    // Interceptor para agregar el token de autenticación
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

  async getAllStores() {
    try {
      const response = await this.api.get('/stores/')
      return response.data.data || response.data || []
    } catch (error) {
      console.error('Error en getAllStores:', error)
      throw error
    }
  }

  async createStore(storeData) {
    try {
      const response = await this.api.post('/stores/', storeData)
      return response.data
    } catch (error) {
      console.error('Error en createStore:', error)
      throw error
    }
  }

  async updateStore(id, storeData) {
    try {
      const response = await this.api.put(`/stores/${id}`, storeData)
      return response.data
    } catch (error) {
      console.error('Error en updateStore:', error)
      throw error
    }
  }

  async deleteStore(id) {
    try {
      const response = await this.api.delete(`/stores/${id}`)
      return response.data
    } catch (error) {
      console.error('Error en deleteStore:', error)
      throw error
    }
  }
}

export default new StoreService()
