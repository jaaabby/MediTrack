const EXPIRATION_THRESHOLD_DAYS = 90

export function useInventoryAlerts() {
  const daysUntilExpiration = (expirationDate) => {
    if (!expirationDate) return null
    return Math.ceil((new Date(expirationDate) - new Date()) / (1000 * 60 * 60 * 24))
  }

  const isExpired = (expirationDate) => {
    const days = daysUntilExpiration(expirationDate)
    return days !== null && days < 0
  }

  // Vence en 90 días o menos (sin contar vencidos)
  const isNearExpiration = (expirationDate) => {
    const days = daysUntilExpiration(expirationDate)
    return days !== null && days >= 0 && days <= EXPIRATION_THRESHOLD_DAYS
  }

  // Vencido → rojo; por vencer (0–90 días) → naranja; resto → gris
  const getExpirationClass = (expirationDate) => {
    const days = daysUntilExpiration(expirationDate)
    if (days === null) return 'text-gray-900'
    if (days < 0) return 'text-red-600 font-semibold'
    if (days <= EXPIRATION_THRESHOLD_DAYS) return 'text-orange-600 font-semibold'
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
