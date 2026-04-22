// Servicio de autenticación para comunicarse con el backend
import { getApiBaseUrl } from '@/config/api.js'

const API_BASE_URL = getApiBaseUrl()

class AuthService {
  constructor() {
    this.baseURL = API_BASE_URL
  }

  // Realizar login
  async login(email, password, rememberMe = false) {
    try {
      const response = await fetch(`${this.baseURL}/auth/login`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          email,
          password,
          remember_me: rememberMe
        })
      })

      const data = await response.json()

      if (!response.ok) {
        const err = new Error(data.error || 'Error al iniciar sesión')
        if (data.data) err.responseData = data.data
        throw err
      }

      if (!data.success) {
        const err = new Error(data.error || 'Error al iniciar sesión')
        if (data.data) err.responseData = data.data
        throw err
      }

      return data.data
    } catch (error) {
      console.error('Error en AuthService.login:', error)
      throw error
    }
  }

  // Registrar nuevo usuario
  async register(userData) {
    try {
      const response = await fetch(`${this.baseURL}/auth/register`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(userData)
      })

      const data = await response.json()

      if (!response.ok) {
        throw new Error(data.error || 'Error al registrar usuario')
      }

      if (!data.success) {
        throw new Error(data.error || 'Error al registrar usuario')
      }

      return data.data
    } catch (error) {
      console.error('Error en AuthService.register:', error)
      throw error
    }
  }

  // Obtener perfil del usuario autenticado
  async getProfile() {
    try {
      const token = this.getToken()
      if (!token) {
        throw new Error('No hay token de autenticación')
      }

      const response = await fetch(`${this.baseURL}/auth/profile`, {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        }
      })

      const data = await response.json()

      if (!response.ok) {
        throw new Error(data.error || 'Error al obtener perfil')
      }

      if (!data.success) {
        throw new Error(data.error || 'Error al obtener perfil')
      }

      return data.data
    } catch (error) {
      console.error('Error en AuthService.getProfile:', error)
      throw error
    }
  }

  // Cambiar contraseña
  async changePassword(currentPassword, newPassword) {
    try {
      const token = this.getToken()
      if (!token) {
        throw new Error('No hay token de autenticación')
      }

      const response = await fetch(`${this.baseURL}/auth/change-password`, {
        method: 'PUT',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          current_password: currentPassword,
          new_password: newPassword
        })
      })

      const data = await response.json()

      if (!response.ok) {
        throw new Error(data.error || 'Error al cambiar contraseña')
      }

      if (!data.success) {
        throw new Error(data.error || 'Error al cambiar contraseña')
      }

      return data
    } catch (error) {
      console.error('Error en AuthService.changePassword:', error)
      throw error
    }
  }

  // Cambiar contraseña por primera vez (contraseña temporal)
  async firstTimePasswordChange(temporaryPassword, newPassword) {
    try {
      const token = this.getToken()
      if (!token) {
        throw new Error('No hay token de autenticación')
      }

      const response = await fetch(`${this.baseURL}/auth/first-time-password-change`, {
        method: 'PUT',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          temporary_password: temporaryPassword,
          new_password: newPassword
        })
      })

      const data = await response.json()

      if (!response.ok) {
        throw new Error(data.error || 'Error al cambiar contraseña')
      }

      if (!data.success) {
        throw new Error(data.error || 'Error al cambiar contraseña')
      }

      return data
    } catch (error) {
      console.error('Error en AuthService.firstTimePasswordChange:', error)
      throw error
    }
  }

  // Solicitar recuperación de contraseña
  async requestPasswordReset(email) {
    try {
      const response = await fetch(`${this.baseURL}/auth/forgot-password`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ email })
      })

      const data = await response.json()

      if (!response.ok) {
        throw new Error(data.error || 'Error al solicitar recuperación de contraseña')
      }

      if (!data.success) {
        throw new Error(data.error || 'Error al solicitar recuperación de contraseña')
      }

      return data
    } catch (error) {
      console.error('Error en AuthService.requestPasswordReset:', error)
      throw error
    }
  }

  // Validar token de recuperación
  async validateResetToken(token) {
    try {
      const response = await fetch(`${this.baseURL}/auth/validate-reset-token`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ token })
      })

      const data = await response.json()

      if (!response.ok) {
        throw new Error(data.error || 'Token inválido o expirado')
      }

      if (!data.success) {
        throw new Error(data.error || 'Token inválido o expirado')
      }

      return data.data
    } catch (error) {
      console.error('Error en AuthService.validateResetToken:', error)
      throw error
    }
  }

  // Resetear contraseña con token
  async resetPassword(token, newPassword) {
    try {
      const response = await fetch(`${this.baseURL}/auth/reset-password`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          token,
          new_password: newPassword
        })
      })

      const data = await response.json()

      if (!response.ok) {
        throw new Error(data.error || 'Error al resetear contraseña')
      }

      if (!data.success) {
        throw new Error(data.error || 'Error al resetear contraseña')
      }

      return data
    } catch (error) {
      console.error('Error en AuthService.resetPassword:', error)
      throw error
    }
  }

  // Guardar token en localStorage
  setToken(token) {
    localStorage.setItem('auth_token', token)
  }

  // Obtener token de localStorage
  getToken() {
    return localStorage.getItem('auth_token')
  }

  // Remover token de localStorage
  removeToken() {
    localStorage.removeItem('auth_token')
  }

  // Verificar si el token ha expirado
  isTokenExpired(token) {
    try {
      const payload = JSON.parse(atob(token.split('.')[1]))
      return payload.exp * 1000 < Date.now()
    } catch (error) {
      return true
    }
  }

  // Verificar si el usuario está autenticado
  isAuthenticated() {
    const token = this.getToken()
    return token && !this.isTokenExpired(token)
  }

  // Obtener información del usuario del token
  getUserFromToken(token) {
    try {
      const payload = JSON.parse(atob(token.split('.')[1]))
      return {
        rut: payload.sub,
        email: payload.email,
        role: payload.role,
        exp: payload.exp
      }
    } catch (error) {
      console.error('Error al decodificar token:', error)
      return null
    }
  }

  // Cerrar sesión en todos los dispositivos
  async logoutAllDevices() {
    const token = this.getToken()
    if (!token) {
      throw new Error('No hay token de autenticación')
    }

    const response = await fetch(`${this.baseURL}/auth/logout-all-devices`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json',
      }
    })

    const data = await response.json()

    if (!response.ok || !data.success) {
      throw new Error(data.error || 'Error al cerrar sesión en todos los dispositivos')
    }

    return data
  }

  // Obtener configuración TOTP (secreto + URL del QR) para setup inicial
  async getTOTPSetup() {
    try {
      const token = this.getToken()
      if (!token) throw new Error('No hay token de autenticación')

      const response = await fetch(`${this.baseURL}/auth/totp/setup`, {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        }
      })
      const data = await response.json()
      if (!response.ok || !data.success) throw new Error(data.error || 'Error al obtener configuración TOTP')
      return data.data
    } catch (error) {
      console.error('Error en AuthService.getTOTPSetup:', error)
      throw error
    }
  }

  // Activar TOTP con el secreto y el código de verificación del usuario
  async activateTOTP(secret, code) {
    try {
      const token = this.getToken()
      if (!token) throw new Error('No hay token de autenticación')

      const response = await fetch(`${this.baseURL}/auth/totp/activate`, {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ secret, code })
      })
      const data = await response.json()
      if (!response.ok || !data.success) throw new Error(data.error || 'Error al activar TOTP')
      return data
    } catch (error) {
      console.error('Error en AuthService.activateTOTP:', error)
      throw error
    }
  }

  // Verificar código TOTP durante el login (usa pre_auth_token)
  async verifyTOTP(preAuthToken, code, rememberMe = false) {
    try {
      const response = await fetch(`${this.baseURL}/auth/totp/verify`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ pre_auth_token: preAuthToken, code, remember_me: rememberMe })
      })
      const data = await response.json()
      if (!response.ok || !data.success) throw new Error(data.error || 'Código TOTP inválido')
      return data.data
    } catch (error) {
      console.error('Error en AuthService.verifyTOTP:', error)
      throw error
    }
  }

  // Deshabilitar TOTP (requiere confirmar contraseña)
  async disableTOTP(password) {
    try {
      const token = this.getToken()
      if (!token) throw new Error('No hay token de autenticación')

      const response = await fetch(`${this.baseURL}/auth/totp`, {
        method: 'DELETE',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ password })
      })
      const data = await response.json()
      if (!response.ok || !data.success) throw new Error(data.error || 'Error al deshabilitar TOTP')
      return data
    } catch (error) {
      console.error('Error en AuthService.disableTOTP:', error)
      throw error
    }
  }

  // Logout
  logout() {
    this.removeToken()
    // Redirigir al login
    window.location.href = '/login'
  }
}

// Crear instancia única del servicio
const authService = new AuthService()

export default authService
