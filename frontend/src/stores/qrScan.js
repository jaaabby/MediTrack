import { defineStore } from 'pinia'
import { ref } from 'vue'

const STORAGE_KEY = 'qr-current-scan'

/**
 * Store del escaneo QR "actual".
 *
 * Mantiene el último QR escaneado para que sobreviva a la navegación
 * (router.push hacia QRConsumer/QRReception/QRTraceability, etc. y el
 * "volver atrás" del browser). Sin esto, al desmontarse QRScanner.vue
 * el estado local se perdía y la tarjeta del escaneo desaparecía al regresar.
 *
 * Persiste también en localStorage para sobrevivir a recargas completas.
 */
export const useQrScanStore = defineStore('qrScan', () => {
  const currentScan = ref(null)      // resultado de qrService.scanQRCode
  const currentContext = ref(null)   // scanContext usado en el escaneo
  const lastInput = ref('')          // texto del input QR
  // Solo se restaura el escaneo al volver cuando el scanner navegó hacia la vista de
  // trazabilidad/detalles. En memoria (no se persiste): un regreso por router.back()
  // mantiene el store vivo; una recarga completa arranca limpio a propósito.
  const keepForTraceability = ref(false)

  // Restaurar desde localStorage al iniciar
  const restore = () => {
    try {
      const saved = localStorage.getItem(STORAGE_KEY)
      if (!saved) return
      const parsed = JSON.parse(saved)
      currentScan.value = parsed.currentScan ?? null
      currentContext.value = parsed.currentContext ?? null
      lastInput.value = parsed.lastInput ?? ''
    } catch (err) {
      console.error('Error restaurando escaneo QR:', err)
    }
  }

  const persist = () => {
    try {
      localStorage.setItem(STORAGE_KEY, JSON.stringify({
        currentScan: currentScan.value,
        currentContext: currentContext.value,
        lastInput: lastInput.value
      }))
    } catch (err) {
      console.error('Error guardando escaneo QR:', err)
    }
  }

  /**
   * Guarda (sobreescribe) el escaneo actual. Un escaneo nuevo siempre
   * reemplaza al anterior — no se acumula.
   */
  const setScan = (scan, context = null, input = '') => {
    currentScan.value = scan
    currentContext.value = context
    lastInput.value = input
    persist()
  }

  /** Limpia el escaneo actual (botón "Limpiar todo"). */
  const clearScan = () => {
    currentScan.value = null
    currentContext.value = null
    lastInput.value = ''
    keepForTraceability.value = false
    try {
      localStorage.removeItem(STORAGE_KEY)
    } catch (err) {
      console.error('Error limpiando escaneo QR:', err)
    }
  }

  return {
    currentScan,
    currentContext,
    lastInput,
    keepForTraceability,
    restore,
    setScan,
    clearScan
  }
})
