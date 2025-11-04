// services/storeService.js
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

  // Obtener todos los almacenes/stores
  async getAllStores() {
    try {
      const response = await this.api.get('/stores/')
      return response.data.data || response.data || []
    } catch (error) {
      console.error('Error al obtener almacenes:', error)
      throw error
    }
  }

  // Obtener almacén por ID
  async getStoreById(id) {
    try {
      const response = await this.api.get(`/stores/${id}`)
      return response.data.data || response.data
    } catch (error) {
      console.error('Error al obtener almacén:', error)
      throw error
    }
  }

  // Obtener almacenes por centro médico
  async getByMedicalCenter(medicalCenterId) {
    try {
      const allStores = await this.getAllStores()
      return allStores.filter(s => s.medical_center_id === parseInt(medicalCenterId))
    } catch (error) {
      console.error('Error al obtener almacenes por centro médico:', error)
      throw error
    }
  }

  // Crear nuevo almacén
  async createStore(storeData) {
    try {
      const response = await this.api.post('/stores', storeData)
      return response.data
    } catch (error) {
      console.error('Error al crear almacén:', error)
      throw error
    }
  }

  // Actualizar almacén
  async updateStore(id, storeData) {
    try {
      const response = await this.api.put(`/stores/${id}`, storeData)
      return response.data
    } catch (error) {
      console.error('Error al actualizar almacén:', error)
      throw error
    }
  }

  // Eliminar almacén
  async deleteStore(id) {
    try {
      const response = await this.api.delete(`/stores/${id}`)
      return response.data
    } catch (error) {
      console.error('Error al eliminar almacén:', error)
      throw error
    }
  }
}

export default new StoreService()
