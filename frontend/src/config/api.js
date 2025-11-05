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
  // Si hay una URL específica configurada (necesaria para Railway/Cloud), usarla
  const envUrl = import.meta.env.VITE_API_BASE_URL
  
  if (envUrl) {
    // Asegurar que la URL termine con /api/v1 si no lo hace
    const cleanUrl = envUrl.endsWith('/') ? envUrl.slice(0, -1) : envUrl
    const apiUrl = cleanUrl.endsWith('/api/v1') ? cleanUrl : `${cleanUrl}/api/v1`
    console.log('🔧 API URL (desde variable de entorno VITE_API_BASE_URL):', apiUrl)
    return apiUrl
  }

  // En Docker Compose, el frontend está detrás de nginx que hace proxy de /api/ al backend
  // Por lo tanto, podemos usar URLs relativas que funcionan tanto en localhost como desde el celular
  // Esto evita problemas de CORS y configuración de IPs
  
  // Detectar si estamos en modo producción (build de Vite) o desarrollo (Vite dev server)
  // En producción, Vite establece MODE='production' y import.meta.env.PROD=true
  // En desarrollo, Vite establece MODE='development' y import.meta.env.DEV=true
  const isProduction = import.meta.env.PROD || import.meta.env.MODE === 'production'
  
  // Si estamos en producción (nginx), usar URL relativa
  // nginx hará el proxy de /api/ al backend automáticamente (solo en Docker Compose)
  // En Railway/Cloud, si VITE_API_BASE_URL no está configurada, esto fallará
  // y el usuario debe configurar VITE_API_BASE_URL en Railway
  if (isProduction) {
    const apiUrl = '/api/v1'
    console.log('🔧 API URL (producción, usando proxy nginx):', apiUrl)
    console.log('   Protocolo:', window.location.protocol)
    console.log('   Hostname:', window.location.hostname)
    console.log('   ⚠️  Si estás en Railway/Cloud y esto no funciona, configura VITE_API_BASE_URL')
    return apiUrl
  }

  // En desarrollo, detectar automáticamente el protocolo y host
  const protocol = window.location.protocol // 'http:' o 'https:'
  const hostname = window.location.hostname // 'localhost' o la IP
  const port = window.location.port || (protocol === 'https:' ? '3443' : '3000')
  
  // Determinar el puerto del backend basado en el protocolo
  // Si estamos en HTTPS (desde el celular), el backend también debe estar en HTTPS
  const backendPort = protocol === 'https:' ? 8443 : 8080
  const backendProtocol = protocol === 'https:' ? 'https' : 'http'
  
  // Construir la URL de la API
  // Si estamos en localhost, usar localhost; si no, usar el mismo hostname
  const apiHost = hostname === 'localhost' || hostname === '127.0.0.1' 
    ? 'localhost' 
    : hostname
  
  const apiUrl = `${backendProtocol}://${apiHost}:${backendPort}/api/v1`
  
  console.log('🔧 API URL detectada automáticamente (desarrollo):', apiUrl)
  console.log('   Protocolo:', protocol)
  console.log('   Hostname:', hostname)
  console.log('   Puerto:', port)
  
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

