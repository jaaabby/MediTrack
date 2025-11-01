import axios from 'axios'
import { useAuthStore } from '@/stores/auth'

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

  // Actualizar estado de solicitud
  async updateSupplyRequestStatus(id, statusData) {
    try {
      const response = await this.api.put(`/supply-requests/${id}/status`, statusData)
      return response.data
    } catch (error) {
      console.error('Error al actualizar estado:', error)
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

    // Los datos del solicitante se obtienen automáticamente de la sesión

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

    if (!requestData.surgery_datetime) {
      errors.push('La fecha y hora de cirugía es obligatoria')
    } else {
      // Validar que la fecha no sea en el pasado
      const surgeryDate = new Date(requestData.surgery_datetime)
      const now = new Date()
      if (surgeryDate < now) {
        errors.push('La fecha y hora de cirugía no puede ser en el pasado')
      }
    }

    return {
      isValid: errors.length === 0,
      errors
    }
  }

  // Formatear datos para envío
  formatSupplyRequestForAPI(formData) {
    const authStore = useAuthStore()
    
    // Validar y formatear surgery_datetime
    let surgeryDatetime = formData.surgery_datetime
    if (!surgeryDatetime || surgeryDatetime === '' || surgeryDatetime === '0000-00-00 00:00:00') {
      // Si no hay fecha, usar la fecha actual + 24 horas como mínimo
      const tomorrow = new Date()
      tomorrow.setDate(tomorrow.getDate() + 1)
      surgeryDatetime = tomorrow.toISOString()
    } else if (typeof surgeryDatetime === 'string' && !surgeryDatetime.includes('T')) {
      // Si es un datetime-local format (YYYY-MM-DDTHH:mm), convertir a ISO
      surgeryDatetime = new Date(surgeryDatetime).toISOString()
    }
    
    return {
      pavilion_id: parseInt(formData.pavilion_id),
      requested_by: authStore.getUserRut || 'SYSTEM',
      requested_by_name: authStore.getUserName || 'Usuario Sistema',
      surgery_datetime: surgeryDatetime,
      notes: formData.notes || '',
      // Campos de médico responsable
      surgeon_id: formData.surgeon_id || null,
      surgeon_name: formData.surgeon_name || null,
      surgery_id: formData.surgery_id ? parseInt(formData.surgery_id) : null,
      specialty_id: formData.specialty_id ? parseInt(formData.specialty_id) : null,
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
      'pendiente_pavedad': 'Pendiente Pavedad',
      'asignado_bodega': 'Asignado a Bodega',
      'en_proceso': 'En Proceso',
      'approved': 'Aprobada',
      'aprobado': 'Aprobado',
      'rejected': 'Rechazada',
      'rechazado': 'Rechazado',
      'in_process': 'En Proceso',
      'completed': 'Completada',
      'completado': 'Completado',
      'cancelled': 'Cancelada',
      'cancelado': 'Cancelado',
      'parcialmente_aprobado': 'Parcialmente Aprobado',
      'pendiente_revision': 'Pendiente de Revisión',
      'devuelto': 'Devuelto al Solicitante',
      'devuelto_al_encargado': 'Devuelto al Encargado'
    }
    return labels[status] || status
  }

  // Obtener colores de estado
  getStatusColor(status) {
    const colors = {
      'pending': 'yellow',
      'pendiente_pavedad': 'purple',
      'asignado_bodega': 'blue',
      'en_proceso': 'blue',
      'approved': 'green',
      'aprobado': 'green',
      'rejected': 'red',
      'rechazado': 'red',
      'in_process': 'blue',
      'completed': 'green',
      'completado': 'green',
      'cancelled': 'gray',
      'cancelado': 'gray',
      'parcialmente_aprobado': 'yellow',
      'pendiente_revision': 'orange',
      'devuelto': 'orange',
      'devuelto_al_encargado': 'blue'
    }
    return colors[status] || 'gray'
  }

  // Formatear fecha de cirugía para mostrar
  formatSurgeryDateTime(surgeryDateTime) {
    if (!surgeryDateTime) return 'No programada'
    
    const date = new Date(surgeryDateTime)
    return date.toLocaleString('es-ES', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit',
      hour12: false
    })
  }

  // Calcular urgencia basada en la fecha de cirugía
  calculateUrgencyFromSurgeryDate(surgeryDateTime) {
    if (!surgeryDateTime) return 'normal'
    
    const surgeryDate = new Date(surgeryDateTime)
    const now = new Date()
    const diffHours = (surgeryDate - now) / (1000 * 60 * 60)
    
    if (diffHours < 6) return 'critical'
    if (diffHours < 24) return 'high'
    if (diffHours < 72) return 'normal'
    return 'low'
  }

  // Obtener color basado en urgencia de fecha de cirugía
  getUrgencyColor(surgeryDateTime) {
    const urgency = this.calculateUrgencyFromSurgeryDate(surgeryDateTime)
    const colors = {
      'critical': 'red',
      'high': 'orange', 
      'normal': 'blue',
      'low': 'gray'
    }
    return colors[urgency] || 'blue'
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

  // ========================
  // WORKFLOW DE APROBACIÓN (PAVEDAD)
  // ========================

  // Asignar solicitud a encargado de bodega (solo Pavedad)
  async assignRequestToWarehouseManager(requestId, assignmentData) {
    try {
      const response = await this.api.put(`/supply-requests/${requestId}/assign`, assignmentData)
      return response.data
    } catch (error) {
      console.error('Error al asignar solicitud:', error)
      throw error
    }
  }

  // Obtener solicitudes pendientes de asignación por Pavedad
  async getPendingRequestsForPavedad() {
    try {
      const response = await this.api.get('/supply-requests/pending-pavedad')
      return response.data
    } catch (error) {
      console.error('Error al obtener solicitudes pendientes para Pavedad:', error)
      throw error
    }
  }

  // Obtener solicitudes asignadas a un encargado de bodega
  async getAssignedRequestsForWarehouseManager(warehouseManagerRut) {
    try {
      const response = await this.api.get(`/supply-requests/assigned/${warehouseManagerRut}`)
      return response.data
    } catch (error) {
      console.error('Error al obtener solicitudes asignadas:', error)
      throw error
    }
  }

  // ========================
  // REVISIÓN INDIVIDUAL DE ITEMS
  // ========================

  // Obtener items de una solicitud
  async getSupplyRequestItems(requestId) {
    try {
      const response = await this.api.get(`/supply-requests/${requestId}/items`)
      return {
        success: true,
        data: response.data.data || response.data
      }
    } catch (error) {
      console.error('Error al obtener items:', error)
      return {
        success: false,
        error: error.response?.data?.error || error.message
      }
    }
  }

  // Revisar un item individual (aceptar, rechazar o devolver)
  async reviewSupplyRequestItem(itemId, reviewData) {
    try {
      const response = await this.api.put(`/supply-requests/items/${itemId}/review`, reviewData)
      return response.data
    } catch (error) {
      console.error('Error al revisar item:', error)
      throw error
    }
  }

  // Reenviar una solicitud devuelta
  async resubmitReturnedRequest(requestId, updatedItems, notes = '') {
    try {
      const response = await this.api.put(`/supply-requests/${requestId}/resubmit`, {
        updated_items: updatedItems,
        notes: notes
      })
      return response.data
    } catch (error) {
      console.error('Error al reenviar solicitud:', error)
      throw error
    }
  }
}

export default new SupplyRequestService()