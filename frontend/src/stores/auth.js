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
    
    // Obtener el nombre del usuario
    getUserName: (state) => state.user?.name,
    
    // Obtener el email del usuario
    getUserEmail: (state) => state.user?.email,
    
    // Obtener el RUT del usuario
    getUserRut: (state) => state.user?.rut,

    // Obtener fecha de creación del usuario
    getUserCreatedAt: (state) => state.user?.created_at,
    
    // Obtener fecha de actualización del usuario
    getUserUpdatedAt: (state) => state.user?.updated_at,
    
    // Obtener el nombre del centro médico
    getUserMedicalCenterName: (state) => state.user?.medical_center?.name,
  },

  actions: {
    // Inicializar el store desde localStorage
    initializeAuth() {
      const token = authService.getToken()
      if (token && !authService.isTokenExpired(token)) {
        this.token = token
        this.isAuthenticated = true
        
        // Obtener información del usuario del token
        const userInfo = authService.getUserFromToken(token)
        if (userInfo) {
          this.user = {
            rut: userInfo.rut,
            email: userInfo.email,
            role: userInfo.role
          }
        }
      } else {
        this.logout()
      }
    },

    // Realizar login
    async login(email, password) {
      this.isLoading = true
      this.error = null

      try {
        const response = await authService.login(email, password)
        
        // Guardar token
        this.token = response.token
        this.user = response.user
        this.isAuthenticated = true
        
        // Guardar en localStorage
        authService.setToken(response.token)
        
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

    // Logout
    logout() {
      this.user = null
      this.token = null
      this.isAuthenticated = false
      this.error = null
      this.isLoading = false
      
      // Remover token del localStorage
      authService.removeToken()
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
        'admin': ['Home', 'Inventory', 'AddSupply', 'QRScanner', 'QRDetails', 'QRConsumer', 'SupplyRequestList', 'SupplyRequestForm', 'SupplyRequestDetail', 'SupplyRequestEdit', 'Profile'],
        'pabellón': ['Home', 'QRScanner', 'QRDetails', 'QRConsumer', 'SupplyRequestList', 'SupplyRequestForm', 'SupplyRequestDetail', 'Profile'],
        'encargado de bodega': ['Home', 'Inventory', 'AddSupply', 'QRScanner', 'QRDetails', 'SupplyRequestList', 'SupplyRequestDetail', 'Profile']
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
