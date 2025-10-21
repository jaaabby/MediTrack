import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const routes = [
  // Rutas públicas
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: {
      title: 'Iniciar Sesión - MediTrack',
      description: 'Acceso al sistema de gestión de insumos médicos',
      requiresAuth: false
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
    component: () => import('@/views/Home.vue'),
    meta: {
      title: 'Inicio - MediTrack',
      description: 'Panel principal del sistema de gestión de insumos médicos',
      requiresAuth: true
    }
  },

  // Gestión de inventario
  {
    path: '/inventory',
    name: 'Inventory',
    component: () => import('@/views/Inventory.vue'),
    meta: {
      title: 'Inventario - MediTrack',
      description: 'Gestión y consulta del inventario de insumos médicos',
      requiresAuth: true
    }
  },
  {
    path: '/inventory/add',
    name: 'AddSupply',
    component: () => import('@/views/AddSupply.vue'),
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
    component: () => import('@/views/Statistics.vue'),
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
    component: () => import('@/views/QRScanner.vue'),
    meta: {
      title: 'Scanner QR - MediTrack',
      description: 'Escáner de códigos QR para insumos médicos',
      requiresAuth: true
    }
  },

  // Funcionalidades específicas de QR
  {
    path: '/qr-consumer',
    name: 'QRConsumer',
    component: () => import('@/views/QRConsumer.vue'),
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
    component: () => import('@/views/QRTransfer.vue'),
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
    component: () => import('@/views/QRReception.vue'),
    meta: {
      title: 'Recepcionar Insumo - MediTrack',
      description: 'Recepcionar insumos médicos con estado "en_camino_a_pabellon"',
      requiresAuth: true
    }
  },

  // NUEVA RUTA: Gestión de Retornos a Bodega
  {
    path: '/return-management',
    name: 'ReturnToBodegaManagement',
    component: () => import('@/views/ReturnToBodegaManagement.vue'),
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
    component: () => import('@/views/QRDetails.vue'),
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
    component: () => import('@/views/QRTraceability.vue'),
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
    component: () => import('@/components/BatchHistory.vue'),
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
    component: () => import('@/views/SupplyRequestList.vue'),
    meta: {
      title: 'Solicitudes de Insumo - MediTrack',
      description: 'Gestión de solicitudes de insumo con trazabilidad QR',
      requiresAuth: true
    }
  },
  {
    path: '/supply-requests/new',
    name: 'SupplyRequestForm',
    component: () => import('@/views/SupplyRequestForm.vue'),
    meta: {
      title: 'Nueva Solicitud - MediTrack',
      description: 'Crear nueva solicitud de insumos médicos',
      requiresAuth: true
    }
  },
  {
    path: '/supply-requests/success',
    name: 'SupplyRequestSuccess',
    component: () => import('@/views/SupplyRequestSuccess.vue'),
    props: route => ({ requestData: route.params.requestData }),
    meta: {
      title: 'Solicitud Creada - MediTrack',
      description: 'Confirmación de solicitud creada exitosamente',
      requiresAuth: true
    }
  },
  {
    path: '/supply-requests/:id',
    name: 'SupplyRequestDetails',
    component: () => import('@/views/SupplyRequestDetail.vue'),
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
    component: () => import('@/views/Profile.vue'),
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
    component: () => import('@/views/TransferManagement.vue'),
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
    component: () => import('@/views/SurgeryManagement.vue'),
    meta: {
      title: 'Tipos de Cirugía - MediTrack',
      description: 'Administración de tipos de procedimientos quirúrgicos',
      requiresAuth: true
    }
  },

  // Historial de Insumos
  {
    path: '/supply-history',
    name: 'SupplyHistoryView',
    component: () => import('@/views/SupplyHistoryView.vue'),
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
    component: () => import('@/views/InventoryDashboard.vue'),
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
    component: () => import('@/views/StoreInventoryView.vue'),
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
    component: () => import('@/views/PavilionInventoryView.vue'),
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
    component: () => import('@/views/NotFound.vue'),
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

  // Verificación de autenticación
  if (to.meta.requiresAuth !== false) {
    // Verificar si el usuario está autenticado
    if (!authStore.isAuthenticated) {
      // Intentar restaurar sesión desde localStorage
      await authStore.initializeAuth()
      
      if (!authStore.isAuthenticated) {
        next({
          name: 'Login',
          query: { redirect: to.fullPath }
        })
        return
      }
    }
  }

  // Si está autenticado y trata de acceder al login, redirigir al home
  if (to.name === 'Login' && authStore.isAuthenticated) {
    next({ name: 'Home' })
    return
  }

  // Protección específica para doctores - solo pueden acceder a rutas de solicitudes, home y perfil
  if (authStore.isAuthenticated && authStore.isDoctor) {
    const allowedRoutesForDoctor = [
      'Home',
      'SupplyRequestList', 
      'SupplyRequestForm', 
      'SupplyRequestDetails', 
      'SupplyRequestSuccess',
      'Profile'
    ]
    
    if (!allowedRoutesForDoctor.includes(to.name)) {
      next({ name: 'Home' })
      return
    }
  }

  next()
})

// Manejar errores de navegación
router.onError((error) => {
  console.error('Error de navegación:', error)
  
  // Si es un error de chunk loading (típico en deployments)
  if (error.message.includes('Loading chunk')) {
    // Recargar la página para obtener la nueva versión
    window.location.reload()
  }
})

export default router