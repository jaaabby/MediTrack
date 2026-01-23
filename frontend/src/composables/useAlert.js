import Swal from 'sweetalert2'
import { useNotification } from './useNotification'

/**
 * Composable para usar alertas, confirmaciones y prompts de forma unificada
 * Usa SweetAlert2 para modales y el sistema de notificaciones para mensajes simples
 */
export function useAlert() {
  const { success, error, warning, info } = useNotification()

  /**
   * Muestra una alerta simple (solo información)
   * @param {string} message - Mensaje a mostrar
   * @param {string} title - Título opcional
   */
  const alert = (message, title = 'Información') => {
    return Swal.fire({
      title,
      text: message,
      icon: 'info',
      confirmButtonText: 'Aceptar',
      confirmButtonColor: '#3b82f6',
      customClass: {
        popup: 'rounded-2xl',
        confirmButton: 'px-6 py-2 rounded-lg font-medium'
      }
    })
  }

  /**
   * Muestra una alerta de éxito
   * @param {string} message - Mensaje a mostrar
   * @param {string} title - Título opcional
   */
  const alertSuccess = (message, title = 'Éxito') => {
    return Swal.fire({
      title,
      text: message,
      icon: 'success',
      confirmButtonText: 'Aceptar',
      confirmButtonColor: '#10b981',
      customClass: {
        popup: 'rounded-2xl',
        confirmButton: 'px-6 py-2 rounded-lg font-medium'
      }
    })
  }

  /**
   * Muestra una alerta de error
   * @param {string} message - Mensaje a mostrar
   * @param {string} title - Título opcional
   */
  const alertError = (message, title = 'Error') => {
    return Swal.fire({
      title,
      text: message,
      icon: 'error',
      confirmButtonText: 'Aceptar',
      confirmButtonColor: '#ef4444',
      customClass: {
        popup: 'rounded-2xl',
        confirmButton: 'px-6 py-2 rounded-lg font-medium'
      }
    })
  }

  /**
   * Muestra una alerta de advertencia
   * @param {string} message - Mensaje a mostrar
   * @param {string} title - Título opcional
   */
  const alertWarning = (message, title = 'Advertencia') => {
    return Swal.fire({
      title,
      text: message,
      icon: 'warning',
      confirmButtonText: 'Aceptar',
      confirmButtonColor: '#f59e0b',
      customClass: {
        popup: 'rounded-2xl',
        confirmButton: 'px-6 py-2 rounded-lg font-medium'
      }
    })
  }

  /**
   * Muestra un diálogo de confirmación
   * @param {string} message - Mensaje a mostrar
   * @param {string} title - Título opcional
   * @param {Object} options - Opciones adicionales (confirmText, cancelText, icon, etc.)
   * @returns {Promise<boolean>} true si se confirma, false si se cancela
   */
  const confirm = async (message, title = 'Confirmar', options = {}) => {
    const {
      confirmText = 'Confirmar',
      cancelText = 'Cancelar',
      icon = 'question',
      confirmButtonColor = '#3b82f6',
      cancelButtonColor = '#6b7280',
      showCancelButton = true,
      dangerMode = false
    } = options

    const result = await Swal.fire({
      title,
      text: message,
      icon,
      showCancelButton,
      confirmButtonText: confirmText,
      cancelButtonText: cancelText,
      confirmButtonColor: dangerMode ? '#ef4444' : confirmButtonColor,
      cancelButtonColor,
      reverseButtons: true,
      customClass: {
        popup: 'rounded-2xl',
        confirmButton: 'px-6 py-2 rounded-lg font-medium',
        cancelButton: 'px-6 py-2 rounded-lg font-medium'
      }
    })

    return result.isConfirmed
  }

  /**
   * Muestra un diálogo de confirmación para acciones peligrosas (eliminar, etc.)
   * @param {string} message - Mensaje a mostrar
   * @param {string} title - Título opcional
   * @returns {Promise<boolean>} true si se confirma, false si se cancela
   */
  const confirmDanger = async (message, title = 'Confirmar eliminación') => {
    return confirm(message, title, {
      confirmText: 'Eliminar',
      icon: 'warning',
      dangerMode: true
    })
  }

  /**
   * Muestra un prompt para ingresar texto
   * @param {string} message - Mensaje a mostrar
   * @param {string} title - Título opcional
   * @param {Object} options - Opciones adicionales (placeholder, inputType, etc.)
   * @returns {Promise<string|null>} El texto ingresado o null si se cancela
   */
  const prompt = async (message, title = 'Ingresar', options = {}) => {
    const {
      placeholder = '',
      inputType = 'text',
      inputValue = '',
      inputValidator = null,
      confirmText = 'Aceptar',
      cancelText = 'Cancelar'
    } = options

    const result = await Swal.fire({
      title,
      text: message,
      input: inputType === 'textarea' ? 'textarea' : inputType,
      inputPlaceholder: placeholder,
      inputValue,
      showCancelButton: true,
      confirmButtonText: confirmText,
      cancelButtonText: cancelText,
      confirmButtonColor: '#3b82f6',
      cancelButtonColor: '#6b7280',
      reverseButtons: true,
      inputValidator,
      customClass: {
        popup: 'rounded-2xl',
        confirmButton: 'px-6 py-2 rounded-lg font-medium',
        cancelButton: 'px-6 py-2 rounded-lg font-medium'
      }
    })

    if (result.isConfirmed && result.value) {
      return result.value
    }
    return null
  }

  return {
    // Alertas (modales)
    alert,
    alertSuccess,
    alertError,
    alertWarning,
    // Confirmaciones
    confirm,
    confirmDanger,
    // Prompts
    prompt,
    // Notificaciones (toasts)
    success,
    error,
    warning,
    info
  }
}
