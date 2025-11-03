import { defineStore } from 'pinia'
import authService from '@/services/authService'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null,
    token: null,
    isAuthenticated: false,
    isLoading: false,
    error: null
  }),

  getters: {
    // Obtener información del usuario
    getUser: (state) => state.user,
    
    // Verificar si está autenticado
    isLoggedIn: (state) => state.isAuthenticated,
    
    // Obtener el rol del usuario
    getUserRole: (state) => state.user?.role,
    
    // Verificar si es administrador
    isAdmin: (state) => state.user?.role === 'admin',
    
    // Verificar si es usuario de pabellón
    isPavilionUser: (state) => state.user?.role === 'pabellón',
    
    // Verificar si es encargado de bodega
    isWarehouseManager: (state) => state.user?.role === 'encargado de bodega',
    
    // Verificar si es enfermera
    isNurse: (state) => state.user?.role === 'enfermera',
    
    // Verificar si es doctor
    isDoctor: (state) => state.user?.role === 'doctor',
    
    // Verificar si es pavedad
    isPavedad: (state) => state.user?.role === 'pavedad',
    
    // Verificar si puede crear solicitudes (enfermera o doctor, pero NO pavedad)
    canCreateRequests: (state) => ['enfermera', 'doctor'].includes(state.user?.role),
    
    // Verificar si puede ver todas las solicitudes (admin o encargado de bodega)
    canViewAllRequests: (state) => ['admin', 'pavedad'].includes(state.user?.role),
    
    // Verificar si puede aprobar/rechazar solicitudes (solo encargado de bodega)
    canApproveRequests: (state) => state.user?.role === 'encargado de bodega',
    
    // Verificar si puede gestionar inventario y QRs (admin o encargado de bodega)
    canManageInventory: (state) => ['admin', 'encargado de bodega'].includes(state.user?.role),
    
    // Verificar si puede ver solicitudes (todos excepto pabellón, pavedad puede ver)
    canViewRequests: (state) => state.user?.role !== 'pabellón',
    
    // Verificar si puede ver inventario (todos excepto pabellón, doctor y pavedad)
    canViewInventory: (state) => !['pabellón', 'doctor', 'pavedad'].includes(state.user?.role),
    
    // Verificar si puede ver estadísticas (todos excepto pabellón, doctor y pavedad)
    canViewStatistics: (state) => !['pabellón', 'doctor', 'pavedad'].includes(state.user?.role),
    
    // Verificar si puede ver inicio/dashboard (todos pueden ver home)
    canViewHome: (state) => true,
    
    // Verificar si puede ver QR scanner (todos excepto doctor y pavedad)
    canViewQR: (state) => !['doctor', 'pavedad'].includes(state.user?.role),
    
    // Obtener el nombre del usuario
    getUserName: (state) => state.user?.name,
    
    // Obtener el email del usuario
    getUserEmail: (state) => state.user?.email,
    
    // Obtener el RUT del usuario
    getUserRut: (state) => state.user?.rut,

    // Obtener la especialidad del usuario (para doctores)
    getUserSpecialty: (state) => state.user?.specialty?.name,
    
    // Obtener el ID de especialidad del usuario (para doctores)
    getUserSpecialtyId: (state) => state.user?.specialty_id,

    // Obtener fecha de creación del usuario
    getUserCreatedAt: (state) => state.user?.created_at,
    
    // Obtener fecha de actualización del usuario
    getUserUpdatedAt: (state) => state.user?.updated_at,
    
    // Obtener el nombre del centro médico
    getUserMedicalCenterName: (state) => state.user?.medical_center?.name,
  },

  actions: {
    // Inicializar el store desde localStorage y restaurar sesión si no ha expirado
    initializeAuth() {
      const token = authService.getToken()
      const expiry = localStorage.getItem('auth_expiry')
      const now = Date.now()
      if (token && !authService.isTokenExpired(token) && expiry && now < Number(expiry)) {
        this.token = token
        this.isAuthenticated = true
        // Restaurar usuario completo desde localStorage
        const userStr = localStorage.getItem('user_full')
        if (userStr) {
          try {
            this.user = JSON.parse(userStr)
          } catch (e) {
            this.user = null
          }
        } else {
          // Fallback: restaurar datos mínimos
          const userInfo = authService.getUserFromToken(token)
          if (userInfo) {
            this.user = {
              rut: userInfo.rut,
              name: userInfo.name || localStorage.getItem('user_name'),
              email: userInfo.email,
              role: userInfo.role
            }
          }
        }
        // Programar logout automático
        const msToExpiry = Number(expiry) - now
        if (msToExpiry > 0) {
          this._setAutoLogout(msToExpiry)
        }
      } else {
        this.logout()
      }
      // Guardar usuario completo en localStorage si está autenticado
      if (this.user && this.isAuthenticated) {
        localStorage.setItem('user_full', JSON.stringify(this.user))
      }
    },

    // Realizar login
    async login(email, password) {
      this.isLoading = true
      this.error = null
      try {
        const response = await authService.login(email, password)
        // Guardar token y usuario
        this.token = response.token
        this.user = response.user
        this.isAuthenticated = true
        authService.setToken(response.token)
        // Guardar usuario completo en localStorage
        if (response.user) {
          localStorage.setItem('user_full', JSON.stringify(response.user))
        }
        // Guardar expiración (1 hora desde ahora)
        const expiry = Date.now() + 60 * 60 * 1000
        localStorage.setItem('auth_expiry', expiry.toString())
        this._setAutoLogout(60 * 60 * 1000)
        return response
      } catch (error) {
        this.error = error.message
        throw error
      } finally {
        this.isLoading = false
      }
    },

    // Obtener perfil del usuario
    async fetchProfile() {
      if (!this.isAuthenticated) {
        throw new Error('Usuario no autenticado')
      }

      this.isLoading = true
      this.error = null

      try {
        const profile = await authService.getProfile()
        this.user = profile
        return profile
      } catch (error) {
        this.error = error.message
        // Si hay error al obtener perfil, podría ser que el token expiró
        if (error.message.includes('token') || error.message.includes('autenticación')) {
          this.logout()
        }
        throw error
      } finally {
        this.isLoading = false
      }
    },

    // Cambiar contraseña
    async changePassword(currentPassword, newPassword) {
      if (!this.isAuthenticated) {
        throw new Error('Usuario no autenticado')
      }

      this.isLoading = true
      this.error = null

      try {
        const response = await authService.changePassword(currentPassword, newPassword)
        return response
      } catch (error) {
        this.error = error.message
        throw error
      } finally {
        this.isLoading = false
      }
    },

    // Logout manual o automático
    logout() {
      this.user = null
      this.token = null
      this.isAuthenticated = false
      this.error = null
      this.isLoading = false
      // Remover token, expiración y usuario completo del localStorage
      authService.removeToken()
      localStorage.removeItem('auth_expiry')
      localStorage.removeItem('user_full')
      // Limpiar timeout de logout automático si existe
      if (this._logoutTimeout) {
        clearTimeout(this._logoutTimeout)
        this._logoutTimeout = null
      }
    },

    // Programar logout automático
    _setAutoLogout(ms) {
      if (this._logoutTimeout) {
        clearTimeout(this._logoutTimeout)
      }
      this._logoutTimeout = setTimeout(() => {
        this.logout()
        // Opcional: redirigir al login si quieres
        window.location.href = '/login'
      }, ms)
    },

    // Verificar permisos
    hasPermission(permission) {
      if (!this.isAuthenticated || !this.user) {
        return false
      }

      // Lógica de permisos basada en roles
      const rolePermissions = {
        'admin': ['read', 'write', 'delete', 'manage_users', 'manage_inventory'],
        'pabellón': ['read', 'consume_supplies'],
        'encargado de bodega': ['read', 'write', 'manage_inventory']
      }

      const userPermissions = rolePermissions[this.user.role] || []
      return userPermissions.includes(permission)
    },

    // Verificar si puede acceder a una ruta
    canAccessRoute(routeName) {
      if (!this.isAuthenticated) {
        return false
      }

      // Definir rutas protegidas por rol
      const routePermissions = {
        'admin': ['Home', 'Inventory', 'AddSupply', 'QRScanner', 'QRDetails', 'QRTraceability', 'QRConsumer', 'SupplyRequestList', 'SupplyRequestForm', 'SupplyRequestDetail', 'SupplyRequestEdit', 'Profile'],
        'pabellón': ['Home', 'QRScanner', 'QRDetails', 'QRTraceability', 'QRConsumer', 'Profile'],
        'encargado de bodega': ['Home', 'Inventory', 'AddSupply', 'QRScanner', 'QRDetails', 'QRTraceability', 'SupplyRequestList', 'SupplyRequestDetail', 'Profile'],
        'enfermera': ['Home', 'Inventory', 'AddSupply', 'QRScanner', 'QRDetails', 'QRTraceability', 'QRConsumer', 'SupplyRequestList', 'SupplyRequestForm', 'SupplyRequestDetail', 'SupplyRequestEdit', 'Profile'],
        'doctor': ['Home', 'Inventory', 'AddSupply', 'QRScanner', 'QRDetails', 'QRTraceability', 'QRConsumer', 'SupplyRequestList', 'SupplyRequestForm', 'SupplyRequestDetail', 'SupplyRequestEdit', 'Profile']
      }

      const allowedRoutes = routePermissions[this.user.role] || []
      return allowedRoutes.includes(routeName)
    },

    // Actualizar información del usuario
    updateUser(userData) {
      if (this.user) {
        this.user = { ...this.user, ...userData }
      }
    },

    // Limpiar errores
    clearError() {
      this.error = null
    }
  }
})
