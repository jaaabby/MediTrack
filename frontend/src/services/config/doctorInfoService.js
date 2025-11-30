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

  async createDoctor(doctorData) {
    try {
      // Asegurar que el rol sea doctor
      doctorData.role = 'doctor'
      const response = await this.api.post('/doctors/', doctorData)
      return response.data.data || response.data
    } catch (error) {
      console.error('Error en createDoctor:', error)
      throw error
    }
  }

  async updateDoctor(rut, doctorData) {
    try {
      // Asegurar que el rol sea doctor
      doctorData.role = 'doctor'
      const response = await this.api.put(`/doctors/${rut}`, doctorData)
      return response.data.data || response.data
    } catch (error) {
      console.error('Error en updateDoctor:', error)
      throw error
    }
  }

  async deleteDoctor(rut) {
    try {
      const response = await this.api.delete(`/doctors/${rut}`)
      return response.data
    } catch (error) {
      console.error('Error en deleteDoctor:', error)
      throw error
    }
  }

  // Métodos de compatibilidad (deprecated)
  async createDoctorInfo(doctorData) {
    return this.createDoctor(doctorData)
  }

  async updateDoctorInfo(rut, doctorData) {
    return this.updateDoctor(rut, doctorData)
  }

  async deleteDoctorInfo(rut) {
    return this.deleteDoctor(rut)
  }
}

export default new DoctorInfoService()
