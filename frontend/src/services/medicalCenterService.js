// Servicio para centros médicos
// Puedes completar los métodos según la API que uses
import axios from 'axios'
import { getApiBaseUrl } from '@/config/api.js'

const API_BASE_URL = getApiBaseUrl()
const API_URL = API_BASE_URL.replace(/\/$/, '') + '/medical-centers/'

const medicalCenterService = {
  // Obtener todos los centros médicos
  async getAll() {
    try {
      console.log('Fetching medical centers from:', API_URL)
      const response = await axios.get(API_URL)
      console.log('Medical centers response:', response.data)
      
      // Soporte para diferentes formatos de respuesta
      const data = response.data.data || response.data.Data || response.data || []
      console.log('Medical centers data:', data)
      return { data: data }
    } catch (error) {
      console.error('Error fetching medical centers:', error)
      
      // Fallback: try without trailing slash if the first request fails
      if (API_URL.endsWith('/')) {
        try {
          console.log('Trying fallback URL without trailing slash')
          const fallbackResponse = await axios.get(API_URL.slice(0, -1))
          const data = fallbackResponse.data.data || fallbackResponse.data.Data || fallbackResponse.data || []
          return { data: data }
        } catch (fallbackError) {
          console.error('Fallback also failed:', fallbackError)
        }
      }
      
      throw error
    }
  },

  // Obtener un centro médico por ID
  async getById(id) {
    const response = await axios.get(`${API_URL}/${id}`)
    return response.data.data || response.data.Data || response.data || null
  },

  // Puedes agregar más métodos según lo que necesites
}

export default medicalCenterService
