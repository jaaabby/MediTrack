import axios from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080/api/v1'

class SupplyRequestService {
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
  // CRUD DE SOLICITUDES
  // ========================

  // Crear nueva solicitud de insumo
  async createSupplyRequest(requestData) {
    try {
      const response = await this.api.post('/supply-requests', requestData)
      return response.data
    } catch (error) {
      console.error('Error al crear solicitud:', error)
      throw error
    }
  }

  // Obtener todas las solicitudes con paginación
  async getAllSupplyRequests(limit = 20, offset = 0, status = '') {
    try {
      const params = { limit, offset }
      if (status) params.status = status
      
      const response = await this.api.get('/supply-requests', { params })
      return response.data
    } catch (error) {
      console.error('Error al obtener solicitudes:', error)
      throw error
    }
  }

  // Obtener solicitud por ID
  async getSupplyRequestById(id) {
    try {
      const response = await this.api.get(`/supply-requests/${id}`)
      return response.data
    } catch (error) {
      console.error('Error al obtener solicitud:', error)
      throw error
    }
  }

  // Obtener solicitudes por pabellón
  async getSupplyRequestsByPavilion(pavilionId, limit = 20, offset = 0) {
    try {
      const response = await this.api.get(`/supply-requests/pavilion/${pavilionId}`, {
        params: { limit, offset }
      })
      return response.data
    } catch (error) {
      console.error('Error al obtener solicitudes del pabellón:', error)
      throw error
    }
  }

  // Aprobar solicitud
  async approveSupplyRequest(id, approvalData) {
    try {
      const response = await this.api.put(`/supply-requests/${id}/approve`, approvalData)
      return response.data
    } catch (error) {
      console.error('Error al aprobar solicitud:', error)
      throw error
    }
  }

  // Rechazar solicitud
  async rejectSupplyRequest(id, rejectionData) {
    try {
      const response = await this.api.put(`/supply-requests/${id}/reject`, rejectionData)
      return response.data
    } catch (error) {
      console.error('Error al rechazar solicitud:', error)
      throw error
    }
  }

  // Completar solicitud
  async completeSupplyRequest(id, completionData) {
    try {
      const response = await this.api.put(`/supply-requests/${id}/complete`, completionData)
      return response.data
    } catch (error) {
      console.error('Error al completar solicitud:', error)
      throw error
    }
  }

  // Eliminar solicitud
  async deleteSupplyRequest(id) {
    try {
      const response = await this.api.delete(`/supply-requests/${id}`)
      return response.data
    } catch (error) {
      console.error('Error al eliminar solicitud:', error)
      throw error
    }
  }

  // ========================
  // ASIGNACIÓN DE QR
  // ========================

  // Asignar QR individual a solicitud
  async assignQRToRequest(assignmentData) {
    try {
      const response = await this.api.post('/qr-assignments', assignmentData)
      return response.data
    } catch (error) {
      console.error('Error al asignar QR:', error)
      throw error
    }
  }

  // Asignación masiva de QRs
  async bulkAssignQRs(bulkAssignmentData) {
    try {
      const response = await this.api.post('/qr-assignments/bulk', bulkAssignmentData)
      return response.data
    } catch (error) {
      console.error('Error en asignación masiva de QRs:', error)
      throw error
    }
  }

  // Marcar QR como entregado
  async markQRAsDelivered(qrCode, deliveryData) {
    try {
      const response = await this.api.put(`/qr-assignments/${qrCode}/deliver`, deliveryData)
      return response.data
    } catch (error) {
      console.error('Error al marcar QR como entregado:', error)
      throw error
    }
  }

  // ========================
  // TRAZABILIDAD
  // ========================

  // Obtener trazabilidad completa de un QR
  async getQRTraceability(qrCode) {
    try {
      const response = await this.api.get(`/traceability/qr/${encodeURIComponent(qrCode)}`)
      return response.data
    } catch (error) {
      console.error('Error al obtener trazabilidad del QR:', error)
      throw error
    }
  }

  // ========================
  // ESTADÍSTICAS
  // ========================

  // Obtener estadísticas de solicitudes
  async getSupplyRequestStats() {
    try {
      const response = await this.api.get('/supply-requests/stats')
      return response.data
    } catch (error) {
      console.error('Error al obtener estadísticas:', error)
      throw error
    }
  }

  // ========================
  // MÉTODOS AUXILIARES
  // ========================

  // Generar número de solicitud (backup frontend)
  generateRequestNumber() {
    const now = new Date()
    const year = now.getFullYear().toString().slice(-2)
    const month = (now.getMonth() + 1).toString().padStart(2, '0')
    const day = now.getDate().toString().padStart(2, '0')
    const time = now.getTime().toString().slice(-6)
    return `REQ-${year}${month}${day}-${time}`
  }

  // Validar datos de solicitud
  validateSupplyRequest(requestData) {
    const errors = []

    if (!requestData.pavilion_id) {
      errors.push('El pabellón es obligatorio')
    }

    if (!requestData.requested_by) {
      errors.push('El solicitante es obligatorio')
    }

    if (!requestData.requested_by_name) {
      errors.push('El nombre del solicitante es obligatorio')
    }

    if (!requestData.items || requestData.items.length === 0) {
      errors.push('Debe agregar al menos un insumo')
    }

    if (requestData.items) {
      requestData.items.forEach((item, index) => {
        if (!item.supply_code) {
          errors.push(`El código del insumo ${index + 1} es obligatorio`)
        }
        if (!item.supply_name) {
          errors.push(`El nombre del insumo ${index + 1} es obligatorio`)
        }
        if (!item.quantity_requested || item.quantity_requested <= 0) {
          errors.push(`La cantidad solicitada del insumo ${index + 1} debe ser mayor a 0`)
        }
      })
    }

    if (!['low', 'normal', 'high', 'critical'].includes(requestData.priority)) {
      errors.push('La prioridad debe ser: low, normal, high o critical')
    }

    return {
      isValid: errors.length === 0,
      errors
    }
  }

  // Formatear datos para envío
  formatSupplyRequestForAPI(formData) {
    return {
      pavilion_id: parseInt(formData.pavilion_id),
      requested_by: formData.requested_by || 'SYSTEM',
      requested_by_name: formData.requested_by_name || 'Sistema MediTrack',
      priority: formData.priority || 'normal',
      notes: formData.notes || '',
      items: formData.items.map(item => ({
        supply_code: parseInt(item.supply_code),
        supply_name: item.supply_name,
        quantity_requested: parseInt(item.quantity_requested),
        specifications: item.specifications || '',
        is_pediatric: item.is_pediatric || false,
        special_requests: item.special_requests || '',
        urgency_level: item.urgency_level || 'normal',
        size: item.size || null,
        brand: item.brand || null
      }))
    }
  }

  // Obtener etiquetas de estado
  getStatusLabel(status) {
    const labels = {
      'pending': 'Pendiente',
      'approved': 'Aprobada',
      'rejected': 'Rechazada',
      'in_process': 'En Proceso',
      'completed': 'Completada',
      'cancelled': 'Cancelada'
    }
    return labels[status] || status
  }

  // Obtener colores de estado
  getStatusColor(status) {
    const colors = {
      'pending': 'yellow',
      'approved': 'green',
      'rejected': 'red',
      'in_process': 'blue',
      'completed': 'green',
      'cancelled': 'gray'
    }
    return colors[status] || 'gray'
  }

  // Obtener etiquetas de prioridad
  getPriorityLabel(priority) {
    const labels = {
      'low': 'Baja',
      'normal': 'Normal',
      'high': 'Alta',
      'critical': 'Crítica'
    }
    return labels[priority] || priority
  }

  // Obtener colores de prioridad
  getPriorityColor(priority) {
    const colors = {
      'low': 'gray',
      'normal': 'blue',
      'high': 'orange',
      'critical': 'red'
    }
    return colors[priority] || 'blue'
  }
}

export default new SupplyRequestService()