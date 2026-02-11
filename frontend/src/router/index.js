import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const routes = [
  // Rutas pÃºblicas
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/auth/Login.vue'),
    meta: {
      title: 'Iniciar SesiÃ³n - MediTrack',
      description: 'Acceso al sistema de gestiÃ³n de insumos mÃ©dicos',
      requiresAuth: false
    }
  },
  {
    path: '/forgot-password',
    name: 'ForgotPassword',
    component: () => import('@/views/auth/ForgotPassword.vue'),
    meta: {
      title: 'Recuperar ContraseÃ±a - MediTrack',
      description: 'Solicitar recuperaciÃ³n de contraseÃ±a',
      requiresAuth: false
    }
  },
  {
    path: '/reset-password',
    name: 'ResetPassword',
    component: () => import('@/views/auth/ResetPassword.vue'),
    meta: {
      title: 'Restablecer ContraseÃ±a - MediTrack',
      description: 'Restablecer contraseÃ±a con token',
      requiresAuth: false
    }
  },
  {
    path: '/first-time-password-change',
    name: 'FirstTimePasswordChange',
    component: () => import('@/views/auth/FirstTimePasswordChange.vue'),
    meta: {
      title: 'Cambio de ContraseÃ±a Obligatorio - MediTrack',
      description: 'Cambio de contraseÃ±a temporal por primera vez',
      requiresAuth: true,
      skipPasswordCheck: true // Evitar bucle infinito en el guard
    }
  },

  // Ruta raÃ­z
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
      description: 'Panel principal del sistema de gestiÃ³n de insumos mÃ©dicos',
      requiresAuth: true
    }
  },

  // GestiÃ³n de usuarios (solo administradores)
  {
    path: '/users',
    name: 'UserManagement',
    component: () => import('@/views/common/UserManagement.vue'),
    meta: {
      title: 'GestiÃ³n de Usuarios - MediTrack',
      description: 'AdministraciÃ³n de usuarios del sistema',
      requiresAuth: true,
      requiredRoles: ['admin']
    }
  },

  // GestiÃ³n de inventario
  {
    path: '/inventory',
    name: 'Inventory',
    component: () => import('@/views/inventory/Inventory.vue'),
    meta: {
      title: 'Inventario - MediTrack',
      description: 'GestiÃ³n y consulta del inventario de insumos mÃ©dicos',
      requiresAuth: true
    }
  },
  {
    path: '/inventory/add',
    name: 'AddSupply',
    component: () => import('@/views/inventory/AddSupply.vue'),
    meta: {
      title: 'Agregar Insumo - MediTrack',
      description: 'Registrar nuevos insumos mÃ©dicos en el sistema',
      requiresAuth: true
    }
  },

  // EstadÃ­sticas y reportes
  {
    path: '/statistics',
    name: 'Statistics',
    component: () => import('@/views/common/Statistics.vue'),
    meta: {
      title: 'EstadÃ­sticas - MediTrack',
      description: 'Panel de anÃ¡lisis y mÃ©tricas del sistema de inventario',
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
      description: 'EscÃ¡ner de cÃ³digos QR para insumos mÃ©dicos',
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
      description: 'Consumir insumos mÃ©dicos con estado "recepcionado" mediante cÃ³digos QR',
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
      description: 'Transferir insumos mÃ©dicos con estado "disponible" a otros centros',
      requiresAuth: true
    }
  },

  // NUEVA RUTA: RecepciÃ³n de insumos
  {
    path: '/qr-reception',
    name: 'QRReception',
    component: () => import('@/views/qr/QRReception.vue'),
    meta: {
      title: 'Recepcionar Insumo - MediTrack',
      description: 'Recepcionar insumos mÃ©dicos con estado "en_camino_a_pabellon"',
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
      description: 'Retirar insumos mÃ©dicos desde bodega con estado "pendiente_retiro"',
      requiresAuth: true
    }
  },

  // NUEVA RUTA: GestiÃ³n de Retornos a Bodega
  {
    path: '/return-management',
    name: 'ReturnToBodegaManagement',
    component: () => import('@/views/management/ReturnToBodegaManagement.vue'),
    meta: {
      title: 'GestiÃ³n de Retornos - MediTrack',
      description: 'Monitoreo y gestiÃ³n de insumos que deben regresar a bodega',
      requiresAuth: true,
      requiredRoles: ['admin', 'encargado de bodega']
    }
  },

  // Rutas especÃ­ficas de QR con historial y trazabilidad
  {
    path: '/qr/:qrCode/details',
    name: 'QRDetails',
    component: () => import('@/views/qr/QRDetails.vue'),
    props: route => ({ qrCode: route.params.qrCode }),
    meta: {
      title: 'Detalles QR - MediTrack',
      description: 'InformaciÃ³n detallada del cÃ³digo QR',
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
      description: 'Trazabilidad completa del cÃ³digo QR estilo Starken',
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
      description: 'GestiÃ³n de solicitudes de insumo con trazabilidad QR',
      requiresAuth: true
    }
  },
  {
    path: '/supply-requests/new',
    name: 'SupplyRequestForm',
    component: () => import('@/views/requests/SupplyRequestForm.vue'),
    meta: {
      title: 'Nueva Solicitud - MediTrack',
      description: 'Crear nueva solicitud de insumos mÃ©dicos',
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
      description: 'Editar solicitud de insumos mÃ©dicos devuelta',
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
      description: 'ConfirmaciÃ³n de solicitud creada exitosamente',
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
      description: 'Detalles y gestiÃ³n de solicitud de insumos',
      requiresAuth: true
    }
  },

  // GestiÃ³n de perfil
  {
    path: '/profile',
    name: 'Profile',
    component: () => import('@/views/common/Profile.vue'),
    meta: {
      title: 'Perfil - MediTrack',
      description: 'GestiÃ³n del perfil de usuario',
      requiresAuth: true
    }
  },

  // === NUEVAS RUTAS DE GESTIÃ"N ===
  
  // GestiÃ³n de Transferencias
  {
    path: '/transfers',
    name: 'TransferManagement',
    component: () => import('@/views/management/TransferManagement.vue'),
    meta: {
      title: 'Transferencias - MediTrack',
      description: 'GestiÃ³n de transferencias entre bodegas y pabellones',
      requiresAuth: true
    }
  },

  // GestiÃ³n de Tipos de CirugÃ­a
  {
    path: '/surgeries',
    name: 'SurgeryManagement',
    component: () => import('@/views/management/SurgeryManagement.vue'),
    meta: {
      title: 'Tipos de CirugÃ­a - MediTrack',
      description: 'AdministraciÃ³n de tipos de procedimientos quirÃºrgicos',
      requiresAuth: true
    }
  },

  // ConfiguraciÃ³n MÃ©dica - Especialidades MÃ©dicas
  {
    path: '/medical-specialties',
    name: 'MedicalSpecialtyManagement',
    component: () => import('@/views/config/MedicalSpecialtyManagement.vue'),
    meta: {
      title: 'Especialidades MÃ©dicas - MediTrack',
      description: 'AdministraciÃ³n de especialidades mÃ©dicas',
      requiresAuth: true
    }
  },

  // ConfiguraciÃ³n MÃ©dica - Insumos TÃ­picos por CirugÃ­a
  {
    path: '/surgery-typical-supplies',
    name: 'SurgeryTypicalSupplyManagement',
    component: () => import('@/views/management/SurgeryTypicalSupplyManagement.vue'),
    meta: {
      title: 'Insumos TÃ­picos por CirugÃ­a - MediTrack',
      description: 'GestiÃ³n de insumos tÃ­picos asociados a cirugÃ­as',
      requiresAuth: true
    }
  },

  // ConfiguraciÃ³n MÃ©dica - InformaciÃ³n de Doctores
  {
    path: '/doctor-info',
    name: 'DoctorInfoManagement',
    component: () => import('@/views/config/DoctorInfoManagement.vue'),
    meta: {
      title: 'InformaciÃ³n de Doctores - MediTrack',
      description: 'GestiÃ³n de informaciÃ³n extendida de doctores',
      requiresAuth: true
    }
  },

  // ConfiguraciÃ³n - ConfiguraciÃ³n de Proveedores
  {
    path: '/supplier-configs',
    name: 'SupplierConfigManagement',
    component: () => import('@/views/config/SupplierConfigManagement.vue'),
    meta: {
      title: 'ConfiguraciÃ³n de Proveedores - MediTrack',
      description: 'GestiÃ³n de alertas de vencimiento por proveedor',
      requiresAuth: true
    }
  },
  {
    path: '/supply-codes',
    name: 'SupplyCodeManagement',
    component: () => import('@/views/config/SupplyCodeManagement.vue'),
    meta: {
      title: 'GestiÃ³n de CÃ³digos de Insumos - MediTrack',
      description: 'GestiÃ³n de cÃ³digos de insumos y niveles crÃ­ticos de stock',
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

  // === NUEVAS RUTAS DE INVENTARIO POR UBICACIÃ"N ===
  
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
      description: 'Stock disponible en cada pabellÃ³n del hospital',
      requiresAuth: true
    }
  },

  // Catch-all para rutas no encontradas
  {
    path: '/:pathMatch(.*)*',
    name: 'NotFound',
    component: () => import('@/views/common/NotFound.vue'),
    meta: {
      title: 'PÃ¡gina No Encontrada - MediTrack',
      description: 'La pÃ¡gina que buscas no existe',
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

// Guards de navegaciÃ³n
router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()
  
  // Establecer tÃ­tulo de la pÃ¡gina
  if (to.meta.title) {
    document.title = to.meta.title
  }
  
  // Establecer meta descripciÃ³n
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

  // Si estÃ¡ autenticado y trata de acceder al login, redirigir al home
  if (to.name === 'Login' && authStore.isAuthenticated) {
    console.log('âœ" Usuario autenticado intentando acceder a login, redirigiendo a home')
    next({ name: 'Home', replace: true })
    return
  }

  // VerificaciÃ³n de autenticaciÃ³n para rutas protegidas
  if (to.meta.requiresAuth !== false) {
    // Verificar si el usuario estÃ¡ autenticado
    if (!authStore.isAuthenticated) {
      console.log('âœ— Usuario no autenticado para ruta protegida')
      // Intentar restaurar sesiÃ³n desde localStorage solo si no estÃ¡ autenticado
      authStore.initializeAuth()
      
      if (!authStore.isAuthenticated) {
        console.log('âœ— No se pudo restaurar sesiÃ³n, redirigiendo a login')
        next({
          name: 'Login',
          query: { redirect: to.fullPath },
          replace: true
        })
        return
      }
      console.log('âœ" SesiÃ³n restaurada exitosamente')
    }

    // CORRECCIÃ"N CRÃTICA: Verificar si el usuario debe cambiar su contraseÃ±a por primera vez
    // Primero verificamos si la ruta tiene skipPasswordCheck (como FirstTimePasswordChange)
    // Si NO tiene skipPasswordCheck Y el usuario debe cambiar contraseÃ±a, redirigimos
    if (authStore.user?.must_change_password) {
      // Solo redirigir si NO estamos en FirstTimePasswordChange y NO tiene skipPasswordCheck
      if (to.name !== 'FirstTimePasswordChange' && to.meta.skipPasswordCheck !== true) {
        console.log('âœ— Usuario debe cambiar contraseÃ±a temporal, redirigiendo...')
        next({
          name: 'FirstTimePasswordChange',
          replace: true
        })
        return
      }
    }

    // Verificar roles requeridos si estÃ¡n especificados en la ruta
    if (to.meta.requiredRoles && to.meta.requiredRoles.length > 0) {
      const userRole = authStore.getUserRole
      if (!to.meta.requiredRoles.includes(userRole)) {
        console.log('âœ— Usuario sin permisos suficientes para acceder a:', to.name)
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
      'FirstTimePasswordChange' // Agregar esta ruta para permitir cambio de contraseÃ±a
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
      'FirstTimePasswordChange' // Agregar esta ruta para permitir cambio de contraseÃ±a
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
      'PavilionInventoryView', // Solo ver inventario de pabellones (solo lectura)
      'QRScanner',
      'QRDetails',
      'QRTraceability',
      'QRConsumer',
      'QRTransfer',
      'QRReception',
      'QRPickup',
      'Profile',
      'FirstTimePasswordChange' // Agregar esta ruta para permitir cambio de contraseÃ±a
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
      'PavilionInventoryView', // Ver inventario de pabellones
      'QRScanner',
      'QRDetails',
      'QRTraceability',
      'QRConsumer',
      'QRTransfer',
      'QRReception',
      'QRPickup',
      'Profile',
      'FirstTimePasswordChange' // Agregar esta ruta para permitir cambio de contraseÃ±a
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