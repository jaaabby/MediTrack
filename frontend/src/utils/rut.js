/**
 * Utilidades para manejo de RUT chileno
 */

/**
 * Calcula el dígito verificador del RUT chileno (módulo 11).
 * @param {string} body - Cuerpo del RUT sin DV (solo dígitos)
 * @returns {string} DV calculado ('0'-'9' o 'K')
 */
export const calculateDv = (body) => {
  if (!body || body.length === 0) return ''
  let sum = 0
  let multiplier = 2
  for (let i = body.length - 1; i >= 0; i--) {
    sum += parseInt(body[i]) * multiplier
    multiplier = multiplier === 7 ? 2 : multiplier + 1
  }
  const remainder = 11 - (sum % 11)
  if (remainder === 11) return '0'
  if (remainder === 10) return 'K'
  return String(remainder)
}

/**
 * Formatea un RUT mientras el usuario escribe:
 * - Elimina todo excepto dígitos y K
 * - Limita a 9 caracteres
 * - Agrega guión automáticamente antes del DV
 * @param {string} rut
 * @returns {string} RUT formateado (ej: "12345678-9")
 */
export const formatRut = (rut) => {
  let clean = rut.replace(/[^0-9kK]/g, '')
  if (clean.length === 0) return ''
  if (clean.length === 1) return clean
  if (clean.length > 9) clean = clean.slice(0, 9)
  const body = clean.slice(0, -1)
  const dv = clean.slice(-1).toUpperCase()
  return `${body}-${dv}`
}

/**
 * Limpia un RUT dejando solo dígitos + K y agrega guión antes del DV.
 * Útil para normalizar un RUT antes de enviarlo al backend.
 * @param {string} rut
 * @returns {string} RUT sin puntos pero con guión (ej: "12345678-9")
 */
export const cleanRut = (rut) => {
  const clean = rut.replace(/[^0-9kK]/g, '')
  if (clean.length <= 1) return clean
  const body = clean.slice(0, -1)
  const dv = clean.slice(-1).toUpperCase()
  return `${body}-${dv}`
}

/**
 * Valida que el RUT tenga formato correcto Y que el DV sea válido.
 * @param {string} rut - RUT con guión (ej: "12345678-9")
 * @returns {boolean}
 */
export const validateRut = (rut) => {
  if (!rut) return false
  const parts = rut.split('-')
  if (parts.length !== 2) return false
  const [body, dv] = parts
  if (!/^[0-9]{1,8}$/.test(body)) return false
  if (!/^[0-9kK]$/.test(dv)) return false
  return calculateDv(body) === dv.toUpperCase()
}

/**
 * Valida solo el formato básico del RUT (sin verificar DV).
 * @param {string} rut - RUT con guión (ej: "12345678-9")
 * @returns {boolean}
 */
export const validateRutFormat = (rut) => {
  return /^[0-9]{1,8}-[0-9kK]$/.test(rut)
}

/**
 * Handler para un input de cuerpo de RUT (solo la parte numérica, sin DV).
 * Llama al callback con el RUT completo formateado.
 * @param {Event} event
 * @param {(rutBody: string, rutCompleto: string) => void} onChange
 */
export const handleRutBodyInput = (event, onChange) => {
  const digits = event.target.value.replace(/[^0-9]/g, '').slice(0, 8)
  event.target.value = digits
  const rutCompleto = digits.length > 0 ? `${digits}-${calculateDv(digits)}` : ''
  onChange(digits, rutCompleto)
}
