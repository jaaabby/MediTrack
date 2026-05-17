const EXPIRATION_THRESHOLD_DAYS = 30

export function useInventoryAlerts() {
  const daysUntilExpiration = (expirationDate) => {
    if (!expirationDate) return null
    return Math.ceil((new Date(expirationDate) - new Date()) / (1000 * 60 * 60 * 24))
  }

  const isExpired = (expirationDate) => {
    const days = daysUntilExpiration(expirationDate)
    return days !== null && days < 0
  }

  // Vence en 30 días o menos (sin contar vencidos)
  const isNearExpiration = (expirationDate) => {
    const days = daysUntilExpiration(expirationDate)
    return days !== null && days >= 0 && days <= EXPIRATION_THRESHOLD_DAYS
  }

  // Vencido o ≤ 30 días → rojo. No aplica naranja para fechas.
  const getExpirationClass = (expirationDate) => {
    const days = daysUntilExpiration(expirationDate)
    if (days === null) return 'text-gray-900'
    if (days <= EXPIRATION_THRESHOLD_DAYS) return 'text-red-600 font-semibold'
    return 'text-gray-900'
  }

  const isLowStock = (current, critical) => {
    if (critical == null) return false
    return current <= (critical || 1)
  }

  const isMediumStock = (current, critical) => {
    if (critical == null) return false
    const threshold = critical || 1
    return current > threshold && current <= threshold * 2
  }

  // ≤ crítico → rojo; ≤ 2× crítico → naranja; resto → gris
  const getStockClass = (current, critical) => {
    if (isLowStock(current, critical)) return 'text-red-600 font-semibold'
    if (isMediumStock(current, critical)) return 'text-orange-600 font-semibold'
    return 'text-gray-900'
  }

  return { isExpired, isNearExpiration, getExpirationClass, isLowStock, isMediumStock, getStockClass }
}
