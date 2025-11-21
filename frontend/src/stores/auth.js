import { defineStore } from 'pinia'
import authService from '@/services/auth/authService'

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
    
    // Verificar si es consignación (basado en email)
    isConsignation: (state) => {
      if (state.user?.role !== 'encargado de bodega') return false
      const email = (state.user?.email || '').toLowerCase()
      return email.includes('bodegaconsignacion') || email.includes('consignacion')
    },
    
    // Verificar si es bodega central (basado en email)
    isCentralWarehouse: (state) => {
      if (state.user?.role !== 'encargado de bodega') return false
      const email = (state.user?.email || '').toLowerCase()
      return !email.includes('bodegaconsignacion') && !email.includes('consignacion')
    },
    
    // Verificar si puede crear solicitudes (solo doctor, NO enfermera ni pavedad)
    canCreateRequests: (state) => state.user?.role === 'doctor',
    
    // Verificar si puede ver todas las solicitudes (admin o encargado de bodega)
    canViewAllRequests: (state) => ['admin', 'pavedad'].includes(state.user?.role),
    
    // Verificar si puede aprobar/rechazar solicitudes (solo encargado de bodega - incluye consignación)
    canApproveRequests: (state) => state.user?.role === 'encargado de bodega',
    
    // Verificar si puede gestionar inventario y QRs (admin o encargado de bodega - incluye consignación)
    canManageInventory: (state) => ['admin', 'encargado de bodega'].includes(state.user?.role),
    
    // Verificar si puede ver inventario general (menú completo) - NO incluye enfermera
    canViewInventoryMenu: (state) => ['admin', 'encargado de bodega'].includes(state.user?.role),
    
    // Verificar si puede ver inventario de pabellones (solo lectura) - enfermera y encargado de bodega pueden ver pabellones
    canViewPavilionInventory: (state) => ['admin', 'encargado de bodega', 'enfermera'].includes(state.user?.role),
    
    // Verificar si puede ver solicitudes (todos excepto pabellón y enfermera, pavedad puede ver)
    canViewRequests: (state) => !['pabellón', 'enfermera'].includes(state.user?.role),
    
    // Verificar si puede ver inventario (todos excepto pabellón, doctor, pavedad, enfermera)
    canViewInventory: (state) => !['pabellón', 'doctor', 'pavedad', 'enfermera'].includes(state.user?.role),
    
    // Verificar si puede ver estadísticas (todos excepto pabellón, doctor, pavedad, enfermera)
    canViewStatistics: (state) => !['pabellón', 'doctor', 'pavedad', 'enfermera'].includes(state.user?.role),
    
    // Verificar si puede ver inicio/dashboard (todos pueden ver home)
    canViewHome: (state) => true,
    
    // Verificar si puede ver QR scanner (todos excepto doctor y pavedad)
    canViewQR: (state) => !['doctor', 'pavedad'].includes(state.user?.role),
    
    // Verificar si puede agregar insumos al inventario (admin, encargado de bodega, enfermera)
    canAddSupplies: (state) => ['admin', 'encargado de bodega'].includes(state.user?.role),
    
    // Verificar si puede gestionar transferencias (admin, encargado de bodega, enfermera, pabellón)
    canManageTransfers: (state) => ['admin', 'encargado de bodega', 'enfermera', 'pabellón'].includes(state.user?.role),
    
    // Verificar si puede ver historial de insumos (admin, encargado de bodega, enfermera)
    canViewSupplyHistory: (state) => ['admin', 'encargado de bodega', 'enfermera'].includes(state.user?.role),
    
    // Verificar si puede gestionar retornos a bodega (admin, encargado de bodega)
    canManageReturns: (state) => ['admin', 'encargado de bodega'].includes(state.user?.role),
    
    // Verificar si puede gestionar configuración médica (solo admin y enfermera)
    canManageMedicalConfig: (state) => ['admin', 'enfermera'].includes(state.user?.role),
    
    // Verificar si puede gestionar configuración del sistema (admin, encargado de bodega - incluye consignación)
    canManageSystemConfig: (state) => ['admin', 'encargado de bodega'].includes(state.user?.role),
    
    // Verificar si puede gestionar cirugías (solo admin y enfermera)
    canManageSurgeries: (state) => ['admin', 'enfermera'].includes(state.user?.role),
    
    // Verificar si puede ver solo sus propias solicitudes (doctor, enfermera)
    canViewOwnRequests: (state) => ['doctor', 'enfermera'].includes(state.user?.role),
    
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
        'admin': ['Home', 'Inventory', 'InventoryDashboard', 'InventoryStore', 'InventoryPavilion', 'AddSupply', 'QRScanner', 'QRDetails', 'QRTraceability', 'QRConsumer', 'QRTransfer', 'QRReception', 'SupplyRequestList', 'SupplyRequestForm', 'SupplyRequestDetail', 'SupplyRequestEdit', 'Statistics', 'Transfers', 'Surgeries', 'SupplyHistory', 'ReturnManagement', 'MedicalSpecialties', 'SurgeryTypicalSupplies', 'DoctorInfo', 'SupplierConfigs', 'SupplyCodes', 'Profile'],
        'pabellón': ['Home', 'PavilionInventoryView', 'QRScanner', 'QRDetails', 'QRTraceability', 'QRConsumer', 'QRTransfer', 'QRReception', 'Profile'],
        'encargado de bodega': ['Home', 'Inventory', 'InventoryDashboard', 'InventoryStore', 'InventoryPavilion', 'AddSupply', 'QRScanner', 'QRDetails', 'QRTraceability', 'QRTransfer', 'QRReception', 'SupplyRequestList', 'SupplyRequestDetail', 'Statistics', 'Transfers', 'SupplyHistory', 'ReturnManagement', 'SupplierConfigs', 'SupplyCodes', 'Profile'],
        'enfermera': ['Home', 'PavilionInventoryView', 'QRScanner', 'QRDetails', 'QRTraceability', 'QRConsumer', 'QRTransfer', 'QRReception', 'Profile'],
        'doctor': ['Home', 'SupplyRequestList', 'SupplyRequestForm', 'SupplyRequestDetail', 'SupplyRequestEdit', 'Profile'],
        'pavedad': ['Home', 'SupplyRequestList', 'SupplyRequestDetail', 'Profile']
      }

      // Consignación tiene las mismas rutas que encargado de bodega
      let allowedRoutes = routePermissions[this.user.role] || []
      
      // Si es consignación, usar las mismas rutas que encargado de bodega
      if (this.user.role === 'encargado de bodega' && this.isConsignation) {
        allowedRoutes = routePermissions['encargado de bodega'] || []
      }
      
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
