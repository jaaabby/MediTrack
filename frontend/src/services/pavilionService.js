// services/pavilionService.js
import axios from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1'

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
      console.log('URL usada:', API_BASE_URL + '/pavilions/');
      let response = await this.api.get('/pavilions/');
      console.log('Respuesta Axios:', response);
      return response.data.data || response.data.Data || response.data || [];
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