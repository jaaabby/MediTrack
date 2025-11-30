/**
 * Configuración automática de la URL de la API
 * Detecta automáticamente el protocolo (HTTP/HTTPS) y el host actual
 * para evitar configuración manual de IP cuando se accede desde el celular
 */

/**
 * Obtiene la URL base de la API de forma automática
 * - Si VITE_API_BASE_URL está definido, lo usa
 * - Si no, usa URL relativa (para funcionar con nginx proxy en Docker)
 *   o detecta automáticamente el protocolo y host actual
 */
export function getApiBaseUrl() {
  // 1. Prioridad: VITE_API_URL (Cloudflare Pages en producción)
  const viteApiUrl = import.meta.env.VITE_API_URL
  if (viteApiUrl) {
    // Asegurar que termine con /api/v1 si no lo tiene
    const url = viteApiUrl.endsWith('/api/v1') ? viteApiUrl : `${viteApiUrl}/api/v1`
    console.log('✅ API URL desde VITE_API_URL:', url)
    return url
  }

  // 2. Si hay VITE_API_BASE_URL (legacy), usarla
  const envUrl = import.meta.env.VITE_API_BASE_URL
  if (envUrl) {
    console.log('✅ API URL desde VITE_API_BASE_URL:', envUrl)
    return envUrl
  }

  // 3. Detectar si estamos en modo producción
  const isProduction = import.meta.env.PROD || import.meta.env.MODE === 'production'
  
  // Si estamos en producción SIN variables de entorno, usar URL relativa (Docker + Nginx)
  // NOTA: En Cloudflare Pages, SIEMPRE debes configurar VITE_API_URL
  if (isProduction) {
    const apiUrl = '/api/v1'
    console.warn('⚠️ Producción sin VITE_API_URL configurada. Usando URL relativa (solo funciona con Nginx proxy):', apiUrl)
    console.log('   Protocolo:', window.location.protocol)
    console.log('   Hostname:', window.location.hostname)
    return apiUrl
  }

  // En desarrollo, detectar automáticamente el protocolo y host
  const protocol = window.location.protocol // 'http:' o 'https:'
  const hostname = window.location.hostname // 'localhost' o la IP
  const port = window.location.port || (protocol === 'https:' ? '3443' : '3000')
  
  // En localhost, siempre usar HTTP primero (puerto 8080)
  // Esto es más común en desarrollo local
  const isLocalhost = hostname === 'localhost' || hostname === '127.0.0.1' || hostname === '::1'
  
  let backendPort = 8080
  let backendProtocol = 'http'
  
  // Si NO estamos en localhost, usar el mismo protocolo que el frontend
  if (!isLocalhost) {
    backendPort = protocol === 'https:' ? 8443 : 8080
    backendProtocol = protocol === 'https:' ? 'https' : 'http'
  }
  
  // Construir la URL de la API
  const apiHost = isLocalhost ? 'localhost' : hostname
  
  const apiUrl = `${backendProtocol}://${apiHost}:${backendPort}/api/v1`
  
  console.log('🔧 API URL detectada automáticamente (desarrollo):', apiUrl)
  console.log('   Protocolo frontend:', protocol)
  console.log('   Hostname:', hostname)
  console.log('   Puerto frontend:', port)
  console.log('   Es localhost:', isLocalhost)
  console.log('   Backend protocol:', backendProtocol)
  console.log('   Backend port:', backendPort)
  
  return apiUrl
}

/**
 * URL base de la API (exportada para uso directo)
 */
export const API_BASE_URL = getApiBaseUrl()

/**
 * Configuración por defecto de axios
 */
export const API_DEFAULT_CONFIG = {
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
  timeout: 30000, // 30 segundos
}

