import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const routes = [
  // Rutas pÃºblicas
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/auth/Login.vue'),
    meta: {
      title: 'Iniciar Sesión - MediTrack',
      description: 'Acceso al sistema de gestión de insumos médicos',
      requiresAuth: false
    }
  },
  {
    path: '/forgot-password',
    name: 'ForgotPassword',
    component: () => import('@/views/auth/ForgotPassword.vue'),
    meta: {
      title: 'Recuperar Contraseña - MediTrack',
      description: 'Solicitar recuperación de contraseña',
      requiresAuth: false
    }
  },
  {
    path: '/reset-password',
    name: 'ResetPassword',
    component: () => import('@/views/auth/ResetPassword.vue'),
    meta: {
      title: 'Restablecer Contraseña - MediTrack',
      description: 'Restablecer contraseña con token',
      requiresAuth: false
    }
  },
  {
    path: '/first-time-password-change',
    name: 'FirstTimePasswordChange',
    component: () => import('@/views/auth/FirstTimePasswordChange.vue'),
    meta: {
      title: 'Cambio de Contraseña Obligatorio - MediTrack',
      description: 'Cambio de contraseña temporal por primera vez',
      requiresAuth: true,
      skipPasswordCheck: true // Evitar bucle infinito en el guard
    }
  },
  {
    path: '/totp-verify',
    name: 'TOTPVerify',
    component: () => import('@/views/auth/TOTPVerify.vue'),
    meta: {
      title: 'Verificación TOTP - MediTrack',
      description: 'Verificación en dos pasos con código TOTP',
      requiresAuth: false
    }
  },
  {
    path: '/totp-setup',
    name: 'TOTPSetup',
    component: () => import('@/views/auth/TOTPSetup.vue'),
    meta: {
      title: 'Configurar TOTP - MediTrack',
      description: 'Configurar verificación en dos pasos',
      requiresAuth: true,
      skipPasswordCheck: true
    }
  },

  // Ruta raíz
  {
    path: '/',
    redirect: '/home'
  },

  // Rutas principales autenticadas
  {
    path: '/home',
    name: 'Home',
    component: () => import('@/views/common/Home.vue'),
    meta: {
      title: 'Inicio - MediTrack',
      description: 'Panel principal del sistema de gestión de insumos médicos',
      requiresAuth: true
    }
  },

  // Gestión de usuarios (solo administradores)
  {
    path: '/users',
    name: 'UserManagement',
    component: () => import('@/views/common/UserManagement.vue'),
    meta: {
      title: 'Gestión de Usuarios - MediTrack',
      description: 'Administración de usuarios del sistema',
      requiresAuth: true,
      requiredRoles: ['admin']
    }
  },

  // Gestión de inventario
  {
    path: '/inventory',
    name: 'Inventory',
    component: () => import('@/views/inventory/Inventory.vue'),
    meta: {
      title: 'Inventario - MediTrack',
      description: 'Gestión y consulta del inventario de insumos médicos',
      requiresAuth: true
    }
  },
  {
    path: '/inventory/add',
    name: 'AddSupply',
    component: () => import('@/views/inventory/AddSupply.vue'),
    meta: {
      title: 'Agregar Insumo - MediTrack',
      description: 'Registrar nuevos insumos médicos en el sistema',
      requiresAuth: true
    }
  },

  // Estadísticas y reportes
  {
    path: '/statistics',
    name: 'Statistics',
    component: () => import('@/views/common/Statistics.vue'),
    meta: {
      title: 'Estadísticas - MediTrack',
      description: 'Panel de análisis y métricas del sistema de inventario',
      requiresAuth: true
    }
  },

  // Funcionalidades de QR - Scanner principal
  {
    path: '/qr',
    name: 'QRScanner',
    component: () => import('@/views/qr/QRScanner.vue'),
    meta: {
      title: 'Scanner QR - MediTrack',
      description: 'Escáner de códigos QR para insumos médicos',
      requiresAuth: true
    }
  },

  // Funcionalidades especÃ­ficas de QR
  {
    path: '/qr-consumer',
    name: 'QRConsumer',
    component: () => import('@/views/qr/QRConsumer.vue'),
    meta: {
      title: 'Consumo QR - MediTrack',
      description: 'Consumir insumos médicos con estado "recepcionado" mediante códigos QR',
      requiresAuth: true
    }
  },

  // NUEVA RUTA: Transferencia de insumos
  {
    path: '/qr-transfer',
    name: 'QRTransfer',
    component: () => import('@/views/qr/QRTransfer.vue'),
    meta: {
      title: 'Transferir Insumo - MediTrack',
      description: 'Transferir insumos médicos con estado "disponible" a otros centros',
      requiresAuth: true
    }
  },

  // NUEVA RUTA: Recepción de insumos
  {
    path: '/qr-reception',
    name: 'QRReception',
    component: () => import('@/views/qr/QRReception.vue'),
    meta: {
      title: 'Recepcionar Insumo - MediTrack',
      description: 'Recepcionar insumos médicos con estado "en_camino_a_pabellon"',
      requiresAuth: true
    }
  },

  // NUEVA RUTA: Retiro de insumos desde bodega
  {
    path: '/qr-pickup',
    name: 'QRPickup',
    component: () => import('@/views/qr/QRPickup.vue'),
    meta: {
      title: 'Retirar Insumo - MediTrack',
      description: 'Retirar insumos médicos desde bodega con estado "pendiente_retiro"',
      requiresAuth: true
    }
  },

  // NUEVA RUTA: Gestión de Retornos a Bodega
  {
    path: '/return-management',
    name: 'ReturnToBodegaManagement',
    component: () => import('@/views/management/ReturnToBodegaManagement.vue'),
    meta: {
      title: 'Gestión de Retornos - MediTrack',
      description: 'Monitoreo y gestión de insumos que deben regresar a bodega',
      requiresAuth: true,
      requiredRoles: ['admin', 'encargado de bodega']
    }
  },

  // Rutas específicas de QR con historial y trazabilidad
  {
    path: '/qr/:qrCode/details',
    name: 'QRDetails',
    component: () => import('@/views/qr/QRDetails.vue'),
    props: route => ({ qrCode: route.params.qrCode }),
    meta: {
      title: 'Detalles QR - MediTrack',
      description: 'Información detallada del código QR',
      requiresAuth: true
    }
  },

  // RUTA ACTUALIZADA: Trazabilidad estilo Starken
  {
    path: '/qr/:qrCode/traceability',
    name: 'QRTraceability',
    component: () => import('@/views/qr/QRTraceability.vue'),
    props: route => ({ qrCode: route.params.qrCode }),
    meta: {
      title: 'Trazabilidad QR - MediTrack',
      description: 'Trazabilidad completa del código QR estilo Starken',
      requiresAuth: true
    }
  },

  // Ruta para historial de lotes
  {
    path: '/batch/:batchId/history',
    name: 'BatchHistory',
    component: () => import('@/components/inventory/BatchHistory.vue'),
    props: route => ({ batchId: parseInt(route.params.batchId) }),
    meta: {
      title: 'Historial del Lote - MediTrack',
      description: 'Historial completo de movimientos del lote',
      requiresAuth: true
    }
  },

  // Rutas de solicitudes de insumo
  {
    path: '/supply-requests',
    name: 'SupplyRequestList',
    component: () => import('@/views/requests/SupplyRequestList.vue'),
    meta: {
      title: 'Solicitudes de Insumo - MediTrack',
      description: 'Gestión de solicitudes de insumo con trazabilidad QR',
      requiresAuth: true
    }
  },
  {
    path: '/supply-requests/new',
    name: 'SupplyRequestForm',
    component: () => import('@/views/requests/SupplyRequestForm.vue'),
    meta: {
      title: 'Nueva Solicitud - MediTrack',
      description: 'Crear nueva solicitud de insumos médicos',
      requiresAuth: true
    }
  },
  {
    path: '/supply-requests/:id/edit',
    name: 'EditSupplyRequest',
    component: () => import('@/views/requests/SupplyRequestForm.vue'),
    props: route => ({ 
      id: parseInt(route.params.id),
      editMode: true 
    }),
    meta: {
      title: 'Editar Solicitud - MediTrack',
      description: 'Editar solicitud de insumos médicos devuelta',
      requiresAuth: true
    }
  },
  {
    path: '/supply-requests/success/:id?',
    name: 'SupplyRequestSuccess',
    component: () => import('@/views/requests/SupplyRequestSuccess.vue'),
    props: route => ({ 
      requestId: route.params.id ? parseInt(route.params.id) : null,
      requestData: route.params.requestData 
    }),
    meta: {
      title: 'Solicitud Creada - MediTrack',
      description: 'Confirmación de solicitud creada exitosamente',
      requiresAuth: true
    }
  },
  {
    path: '/supply-requests/:id',
    name: 'SupplyRequestDetails',
    component: () => import('@/views/requests/SupplyRequestDetail.vue'),
    props: route => ({ requestId: parseInt(route.params.id) }),
    meta: {
      title: 'Detalles de Solicitud - MediTrack',
      description: 'Detalles y gestión de solicitud de insumos',
      requiresAuth: true
    }
  },

  // Gestión de perfil
  {
    path: '/profile',
    name: 'Profile',
    component: () => import('@/views/common/Profile.vue'),
    meta: {
      title: 'Perfil - MediTrack',
      description: 'Gestión del perfil de usuario',
      requiresAuth: true
    }
  },

  // === NUEVAS RUTAS DE GESTIÓN ===
  
  // Gestión de Transferencias
  {
    path: '/transfers',
    name: 'TransferManagement',
    component: () => import('@/views/management/TransferManagement.vue'),
    meta: {
      title: 'Transferencias - MediTrack',
      description: 'Gestión de transferencias entre bodegas y pabellones',
      requiresAuth: true
    }
  },

  // Gestión de Tipos de Cirugía
  {
    path: '/surgeries',
    name: 'SurgeryManagement',
    component: () => import('@/views/management/SurgeryManagement.vue'),
    meta: {
      title: 'Tipos de Cirugía - MediTrack',
      description: 'Administración de tipos de procedimientos quirúrgicos',
      requiresAuth: true
    }
  },

  // Configuración Médica - Especialidades Médicas
  {
    path: '/medical-specialties',
    name: 'MedicalSpecialtyManagement',
    component: () => import('@/views/config/MedicalSpecialtyManagement.vue'),
    meta: {
      title: 'Especialidades Médicas - MediTrack',
      description: 'Administración de especialidades médicas',
      requiresAuth: true
    }
  },

  // Configuración Médica - Insumos Típicos por Cirugía
  {
    path: '/surgery-typical-supplies',
    name: 'SurgeryTypicalSupplyManagement',
    component: () => import('@/views/management/SurgeryTypicalSupplyManagement.vue'),
    meta: {
      title: 'Insumos Típicos por Cirugía - MediTrack',
      description: 'Gestión de insumos típicos asociados a cirugías',
      requiresAuth: true
    }
  },

  // Configuración Médica - Información de Doctores
  {
    path: '/doctor-info',
    name: 'DoctorInfoManagement',
    component: () => import('@/views/config/DoctorInfoManagement.vue'),
    meta: {
      title: 'Información de Doctores - MediTrack',
      description: 'Gestión de información extendida de doctores',
      requiresAuth: true
    }
  },

  // Configuración - Configuración de Proveedores
  {
    path: '/supplier-configs',
    name: 'SupplierConfigManagement',
    component: () => import('@/views/config/SupplierConfigManagement.vue'),
    meta: {
      title: 'Configuración de Proveedores - MediTrack',
      description: 'Gestión de alertas de vencimiento por proveedor',
      requiresAuth: true
    }
  },
  {
    path: '/supply-codes',
    name: 'SupplyCodeManagement',
    component: () => import('@/views/config/SupplyCodeManagement.vue'),
    meta: {
      title: 'Gestión de Códigos de Insumos - MediTrack',
      description: 'Gestión de códigos de insumos y niveles críticos de stock',
      requiresAuth: true
    }
  },

  // Historial de Insumos
  {
    path: '/supply-history',
    name: 'SupplyHistoryView',
    component: () => import('@/views/inventory/SupplyHistoryView.vue'),
    meta: {
      title: 'Historial de Insumos - MediTrack',
      description: 'Registro completo de movimientos de insumos',
      requiresAuth: true
    }
  },

  // === NUEVAS RUTAS DE INVENTARIO POR UBICACIÓN ===
  
  // Dashboard de Inventario
  {
    path: '/inventory/dashboard',
    name: 'InventoryDashboard',
    component: () => import('@/views/inventory/InventoryDashboard.vue'),
    meta: {
      title: 'Dashboard de Inventario - MediTrack',
      description: 'Resumen general del inventario por ubicaciones',
      requiresAuth: true
    }
  },

  // Inventario de Bodegas
  {
    path: '/inventory/store',
    name: 'StoreInventoryView',
    component: () => import('@/views/inventory/StoreInventoryView.vue'),
    meta: {
      title: 'Inventario de Bodegas - MediTrack',
      description: 'Stock detallado en cada bodega del sistema',
      requiresAuth: true
    }
  },

  // Inventario de Pabellones
  {
    path: '/inventory/pavilion',
    name: 'PavilionInventoryView',
    component: () => import('@/views/inventory/PavilionInventoryView.vue'),
    meta: {
      title: 'Inventario de Pabellones - MediTrack',
      description: 'Stock disponible en cada pabellón del hospital',
      requiresAuth: true
    }
  },

  // Catch-all para rutas no encontradas
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/common/NotFound.vue'),
    meta: {
      title: 'Página No Encontrada - MediTrack',
      description: 'La página que buscas no existe',
      requiresAuth: false
    }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    } else {
      return { top: 0 }
    }
  }
})

// Guards de navegación
router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()
  
  // Establecer título de la página
  if (to.meta.title) {
    document.title = to.meta.title
  }
  
  // Establecer meta descripción
  if (to.meta.description) {
    let metaDescription = document.querySelector('meta[name="description"]')
    if (metaDescription) {
      metaDescription.setAttribute('content', to.meta.description)
    } else {
      metaDescription = document.createElement('meta')
      metaDescription.setAttribute('name', 'description')
      metaDescription.setAttribute('content', to.meta.description)
      document.head.appendChild(metaDescription)
    }
  }

  // Si está autenticado y trata de acceder al login, redirigir al home
  if (to.name === 'Login' && authStore.isAuthenticated) {
    console.log('✔ Usuario autenticado intentando acceder a login, redirigiendo a home')
    next({ name: 'Home', replace: true })
    return
  }

  // Verificación de autenticación para rutas protegidas
  if (to.meta.requiresAuth !== false) {
    // Verificar si el usuario está autenticado
    if (!authStore.isAuthenticated) {
      console.log('✖ Usuario no autenticado para ruta protegida')
      // Intentar restaurar sesión desde localStorage solo si no está autenticado
      authStore.initializeAuth()
      
      if (!authStore.isAuthenticated) {
        console.log('✖ No se pudo restaurar sesión, redirigiendo a login')
        next({
          name: 'Login',
          query: { redirect: to.fullPath },
          replace: true
        })
        return
      }
      console.log('✔ Sesión restaurada exitosamente')
    }

    // CORRECCIÓN CRÍTICA: Verificar si el usuario debe cambiar su contraseña por primera vez
    // Primero verificamos si la ruta tiene skipPasswordCheck (como FirstTimePasswordChange)
    // Si NO tiene skipPasswordCheck Y el usuario debe cambiar contraseña, redirigimos
    if (authStore.user?.must_change_password) {
      // Solo redirigir si NO estamos en FirstTimePasswordChange y NO tiene skipPasswordCheck
      if (to.name !== 'FirstTimePasswordChange' && to.meta.skipPasswordCheck !== true) {
        console.log('✖ Usuario debe cambiar contraseña temporal, redirigiendo...')
        next({
          name: 'FirstTimePasswordChange',
          replace: true
        })
        return
      }
    }

    // Verificar roles requeridos si están especificados en la ruta
    if (to.meta.requiredRoles && to.meta.requiredRoles.length > 0) {
      const userRole = authStore.getUserRole
      if (!to.meta.requiredRoles.includes(userRole)) {
        console.log('✖ Usuario sin permisos suficientes para acceder a:', to.name)
        next({ name: 'Home', replace: true })
        return
      }
    }
  }

  // ProtecciÃ³n especÃ­fica para doctores - pueden acceder a rutas de solicitudes, home, QR scanner y perfil
  if (authStore.isAuthenticated && authStore.isDoctor) {
    const allowedRoutesForDoctor = [
      'Home',
      'SupplyRequestList',
      'SupplyRequestForm',
      'EditSupplyRequest',
      'SupplyRequestDetails',
      'SupplyRequestSuccess',
      'QRScanner',
      'QRDetails',
      'QRTraceability',
      'QRReception',
      'QRPickup',
      'Profile',
      'FirstTimePasswordChange',
      'TOTPSetup'
    ]
    
    if (!allowedRoutesForDoctor.includes(to.name)) {
      console.log('âœ— Doctor intentando acceder a ruta no permitida:', to.name)
      next({ name: 'Home', replace: true })
      return
    }
  }

  // ProtecciÃ³n especÃ­fica para pavedad - solo pueden acceder a home, solicitudes (sin crear) y perfil
  if (authStore.isAuthenticated && authStore.isPavedad) {
    const allowedRoutesForPavedad = [
      'Home',
      'SupplyRequestList',
      'SupplyRequestDetails',
      'Profile',
      'FirstTimePasswordChange',
      'TOTPSetup'
    ]
    
    if (!allowedRoutesForPavedad.includes(to.name)) {
      console.log('âœ— Pavedad intentando acceder a ruta no permitida:', to.name)
      next({ name: 'Home', replace: true })
      return
    }
  }

  // ProtecciÃ³n especÃ­fica para enfermera - solo pueden ver inventario de pabellones y usar QR
  if (authStore.isAuthenticated && authStore.isNurse) {
    const allowedRoutesForNurse = [
      'Home',
      'PavilionInventoryView',
      'QRScanner',
      'QRDetails',
      'QRTraceability',
      'QRConsumer',
      'QRTransfer',
      'QRReception',
      'QRPickup',
      'Profile',
      'FirstTimePasswordChange',
      'TOTPSetup'
    ]
    
    if (!allowedRoutesForNurse.includes(to.name)) {
      console.log('âœ— Enfermera intentando acceder a ruta no permitida:', to.name)
      next({ name: 'Home', replace: true })
      return
    }
  }

  // ProtecciÃ³n especÃ­fica para perfil de pabellÃ³n - pueden ver inventario de pabellones y usar QR
  if (authStore.isAuthenticated && authStore.isPavilionUser) {
    const allowedRoutesForPavilion = [
      'Home',
      'PavilionInventoryView',
      'QRScanner',
      'QRDetails',
      'QRTraceability',
      'QRConsumer',
      'QRTransfer',
      'QRReception',
      'QRPickup',
      'Profile',
      'FirstTimePasswordChange',
      'TOTPSetup'
    ]
    
    if (!allowedRoutesForPavilion.includes(to.name)) {
      console.log('âœ— Usuario de pabellÃ³n intentando acceder a ruta no permitida:', to.name)
      next({ name: 'Home', replace: true })
      return
    }
  }

  // ProtecciÃ³n especÃ­fica para encargado de bodega - no puede acceder a configuraciÃ³n mÃ©dica
  if (authStore.isAuthenticated && authStore.isWarehouseManager) {
    const restrictedRoutesForWarehouse = [
      'MedicalSpecialtyManagement',
      'SurgeryTypicalSupplyManagement',
      'DoctorInfoManagement',
      'SurgeryManagement'
    ]
    
    if (restrictedRoutesForWarehouse.includes(to.name)) {
      console.log('âœ— Encargado de bodega intentando acceder a ruta de configuraciÃ³n mÃ©dica no permitida:', to.name)
      next({ name: 'Home', replace: true })
      return
    }
  }

  // ConsignaciÃ³n ahora tiene las mismas rutas que encargado de bodega (se maneja en canAccessRoute del store)

  next()
})

// Manejar errores de navegaciÃ³n
router.onError((error) => {
  console.error('Error de navegaciÃ³n:', error)
  
  // Si es un error de chunk loading (tÃ­pico en deployments)
  if (error.message.includes('Loading chunk')) {
    // Recargar la pÃ¡gina para obtener la nueva versiÃ³n
    window.location.reload()
  }
})

export default router