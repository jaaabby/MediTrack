import { format } from 'date-fns'
import { es } from 'date-fns/locale'

/**
 * Parsea una fecha proveniente de la BD interpretándola SIEMPRE como hora local.
 *
 * Contexto: las columnas de fecha son TIMESTAMP (sin zona) y guardan el wall-clock
 * LOCAL de Chile. Según cómo las serialice el backend pueden llegar como naive
 * ("2026-05-29 02:41:00"), con "Z" (mal etiquetadas como UTC) o con offset. En todos
 * los casos lo que importa es el wall-clock que se guardó. Esta función extrae los
 * componentes (ignorando cualquier marca de zona) y construye un Date LOCAL, evitando
 * que `new Date(str)` reste horas al asumir UTC.
 *
 * @param {string|Date|number} value
 * @returns {Date|null}
 */
export function parseDbDate(value) {
  if (!value) return null
  if (value instanceof Date) return value
  if (typeof value !== 'string') return new Date(value)

  // Fecha + hora (con 'T' o espacio); se ignora cualquier sufijo de zona (Z / +hh:mm / -hh:mm)
  const m = value.match(/^(\d{4})-(\d{2})-(\d{2})[T ](\d{2}):(\d{2})(?::(\d{2}))?/)
  if (m) {
    const [, y, mo, d, h, mi, s] = m
    return new Date(Number(y), Number(mo) - 1, Number(d), Number(h), Number(mi), Number(s || 0))
  }

  // Solo fecha (YYYY-MM-DD) → medianoche local
  const d2 = value.match(/^(\d{4})-(\d{2})-(\d{2})$/)
  if (d2) {
    const [, y, mo, d] = d2
    return new Date(Number(y), Number(mo) - 1, Number(d))
  }

  return new Date(value)
}

/**
 * Formatea una fecha de BD como hora local (sin conversión UTC).
 * @param {string|Date} value
 * @param {string} pattern - patrón date-fns (default 'dd/MM/yyyy HH:mm')
 * @param {string} fallback - texto si la fecha es inválida/ausente
 */
export function formatDbDateTime(value, pattern = 'dd/MM/yyyy HH:mm', fallback = 'N/A') {
  const d = parseDbDate(value)
  if (!d || isNaN(d.getTime())) return fallback
  return format(d, pattern, { locale: es })
}
