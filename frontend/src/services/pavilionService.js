// services/pavilionService.js
import axios from 'axios'
import { getApiBaseUrl } from '@/config/api.js'

const API_BASE_URL = getApiBaseUrl()

class PavilionService {
  constructor() {
    this.api = axios.create({
      baseURL: API_BASE_URL,
      headers: {
        'Content-Type': 'application/json'
      }
    })
  }

  // Obtener todos los pabellones
  async getAllPavilions() {
    try {
      let response = await this.api.get('/pavilions/');
      
      const data = response.data.data || response.data.Data || response.data || [];
      
      return data;
    } catch (error) {
      console.error('Error al obtener pabellones:', error);
      if (error.response) {
        console.error('Respuesta del servidor:', error.response);
      }
      throw error;
    }
  }

  // Obtener pabellón por ID
  async getPavilionById(id) {
    try {
      const response = await this.api.get(`/pavilions/${id}`)
      return response.data.data || response.data
    } catch (error) {
      console.error('Error al obtener pabellón:', error)
      throw error
    }
  }

  // Obtener pabellones por centro médico
  async getByMedicalCenter(medicalCenterId) {
    try {
      console.log('Fetching pavilions for medical center:', medicalCenterId)
      const allPavilions = await this.getAllPavilions()
      console.log('All pavilions:', allPavilions)
      const filtered = allPavilions.filter(p => p.medical_center_id === parseInt(medicalCenterId))
      console.log('Filtered pavilions:', filtered)
      return { data: filtered }
    } catch (error) {
      console.error('Error al obtener pabellones por centro médico:', error)
      throw error
    }
  }

  // Crear nuevo pabellón
  async createPavilion(pavilionData) {
    try {
      const response = await this.api.post('/pavilions', pavilionData)
      return response.data
    } catch (error) {
      console.error('Error al crear pabellón:', error)
      throw error
    }
  }

  // Actualizar pabellón
  async updatePavilion(id, pavilionData) {
    try {
      const response = await this.api.put(`/pavilions/${id}`, pavilionData)
      return response.data
    } catch (error) {
      console.error('Error al actualizar pabellón:', error)
      throw error
    }
  }

  // Eliminar pabellón
  async deletePavilion(id) {
    try {
      const response = await this.api.delete(`/pavilions/${id}`)
      return response.data
    } catch (error) {
      console.error('Error al eliminar pabellón:', error)
      throw error
    }
  }

  // Buscar pabellones por nombre
  async searchPavilions(searchTerm) {
    try {
      const pavilions = await this.getAllPavilions()
      return pavilions.filter(pavilion => 
        pavilion.name.toLowerCase().includes(searchTerm.toLowerCase())
      )
    } catch (error) {
      console.error('Error al buscar pabellones:', error)
      throw error
    }
  }

  // Enviar insumo a pabellón (si tienes este endpoint en el backend)
  async sendSupplyToPavilion(qrCode, pavilionId, additionalData = {}) {
    try {
      const payload = {
        qr_code: qrCode,
        pavilion_id: pavilionId,
        ...additionalData
      }
      const response = await this.api.post('/pavilions/transfer', payload)
      return response.data
    } catch (error) {
      console.error('Error al enviar insumo a pabellón:', error)
      throw error
    }
  }
}

export default new PavilionService()