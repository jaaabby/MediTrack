import axios from 'axios'
import { getApiBaseUrl } from '@/config/api.js'

const API_BASE_URL = getApiBaseUrl()

class DoctorInfoService {
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

  async getAllDoctors() {
    try {
      const response = await this.api.get('/doctors/all')
      return response.data.data?.doctors || response.data.doctors || []
    } catch (error) {
      console.error('Error en getAllDoctors:', error)
      throw error
    }
  }

  async getDoctorByRut(rut) {
    try {
      const response = await this.api.get(`/doctors/${rut}`)
      return response.data.data || response.data
    } catch (error) {
      console.error('Error en getDoctorByRut:', error)
      throw error
    }
  }

  async createDoctorInfo(doctorData) {
    try {
      const response = await this.api.post('/doctors/', doctorData)
      return response.data.data || response.data
    } catch (error) {
      console.error('Error en createDoctorInfo:', error)
      throw error
    }
  }

  async updateDoctorInfo(rut, doctorData) {
    try {
      const response = await this.api.put(`/doctors/${rut}`, doctorData)
      return response.data.data || response.data
    } catch (error) {
      console.error('Error en updateDoctorInfo:', error)
      throw error
    }
  }

  async deleteDoctorInfo(rut) {
    try {
      const response = await this.api.delete(`/doctors/${rut}`)
      return response.data
    } catch (error) {
      console.error('Error en deleteDoctorInfo:', error)
      throw error
    }
  }
}

export default new DoctorInfoService()

