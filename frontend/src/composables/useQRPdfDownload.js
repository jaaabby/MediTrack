import { ref } from 'vue'
import jsPDF from 'jspdf'
import QRCode from 'qrcode'
import { format } from 'date-fns'
import { es } from 'date-fns/locale'

/**
 * Composable para generar y descargar PDFs de códigos QR
 * @returns {Object} Objeto con funciones y estado reactivo
 */
export function useQRPdfDownload() {
  const isGenerating = ref(false)
  const error = ref(null)

  /**
   * Genera un código QR como imagen base64
   * @param {string} qrCode - Código QR a generar
   * @returns {Promise<string>} - Imagen QR en formato base64
   */
  const generateQRImage = async (qrCode) => {
    try {
      const qrImageUrl = await QRCode.toDataURL(qrCode, {
        width: 300,
        margin: 2,
        color: {
          dark: '#000000',
          light: '#FFFFFF'
        }
      })
      return qrImageUrl
    } catch (err) {
      throw new Error('Error al generar imagen QR: ' + err.message)
    }
  }

  /**
   * Formatea una fecha para mostrar en el PDF
   * @param {string|Date} date - Fecha a formatear
   * @returns {string} - Fecha formateada
   */
  const formatDate = (date) => {
    if (!date) return 'No disponible'
    try {
      const formattedDate = format(new Date(date), 'dd/MM/yyyy', { locale: es })
      return formattedDate
    } catch (error) {
      console.warn('Error al formatear fecha:', date, error)
      return 'Fecha inválida'
    }
  }

  /**
   * Obtiene el estado de disponibilidad del insumo
   * @param {Object} qrInfo - Información del QR
   * @returns {Object} - Estado con label y color
   */
  const getAvailabilityStatus = (qrInfo) => {
    // Verificar si está consumido siguiendo la misma lógica que QRInfoDisplay.vue
    const isConsumed = qrInfo.supply_info?.is_consumed || 
                      qrInfo.is_consumed || 
                      qrInfo.supply_info?.IsConsumed
    if (isConsumed) {
      return { label: 'Consumido', color: '#dc2626' }
    }
    
    // Verificar fechas de vencimiento
    const expiryDate = qrInfo.supply_info?.expiry_date || 
                      qrInfo.batch_info?.expiry_date || 
                      qrInfo.expiry_date ||
                      qrInfo.supply_info?.batch?.expiration_date
    if (expiryDate) {
      const expiry = new Date(expiryDate)
      const today = new Date()
      
      if (expiry < today) {
        return { label: 'Vencido', color: '#dc2626' }
      }
      
      const daysUntilExpiry = Math.ceil((expiry - today) / (1000 * 60 * 60 * 24))
      if (daysUntilExpiry <= 30) {
        return { label: 'Próximo a vencer', color: '#f59e0b' }
      }
    }
    
    return { label: 'Disponible', color: '#059669' }
  }

  /**
   * Descarga la información del QR como PDF
   * @param {Object} qrInfo - Información del código QR
   * @param {Object} options - Opciones de configuración
   * @returns {Promise<boolean>} - True si se descargó exitosamente
   */
  const downloadQRAsPDF = async (qrInfo, options = {}) => {
    if (!qrInfo) {
      error.value = 'No hay información del QR para descargar'
      return false
    }

    console.log('=== Generando PDF para QR:', qrInfo.qr_code, '===')
    console.log('Estado:', qrInfo.is_consumed ? 'Consumido' : 'Disponible')
    console.log('Ubicación:', qrInfo.traceability?.current_location?.name || 'No asignada')
    console.log('Fecha vencimiento:', qrInfo.supply_info?.batch?.expiration_date || 'No disponible')
    console.log('Cantidad:', qrInfo.supply_info?.batch?.amount || 'No disponible')
    console.log('Proveedor:', qrInfo.supply_info?.batch?.supplier || 'No disponible')

    isGenerating.value = true
    error.value = null

    try {
      // Crear nuevo documento PDF
      const pdf = new jsPDF({
        orientation: 'portrait',
        unit: 'mm',
        format: 'a4'
      })

      const pageWidth = pdf.internal.pageSize.getWidth()
      const pageHeight = pdf.internal.pageSize.getHeight()
      
      // Márgenes
      const margin = 20
      const contentWidth = pageWidth - (margin * 2)
      let currentY = margin

      // ===== HEADER =====
      // Título principal
      pdf.setFontSize(18)
      pdf.setFont('helvetica', 'bold')
      pdf.setTextColor(0, 0, 0)
      pdf.text('Información de Insumo Médico', margin, currentY)
      currentY += 10

      // Línea divisoria
      pdf.setLineWidth(0.5)
      pdf.setDrawColor(0, 0, 0)
      pdf.line(margin, currentY, pageWidth - margin, currentY)
      currentY += 15

      // ===== CÓDIGO QR =====
      try {
        const qrImageUrl = await generateQRImage(qrInfo.qr_code)
        const qrSize = 40
        const qrX = pageWidth - margin - qrSize
        
        pdf.addImage(qrImageUrl, 'PNG', qrX, currentY - 5, qrSize, qrSize)
        
        // Código QR como texto debajo de la imagen
        pdf.setFontSize(10)
        pdf.setFont('helvetica', 'normal')
        pdf.setTextColor(100, 100, 100)
        const qrTextLines = pdf.splitTextToSize(qrInfo.qr_code, qrSize)
        pdf.text(qrTextLines, qrX, currentY + qrSize + 5)
        
      } catch (qrError) {
        console.warn('Error generando QR image:', qrError)
      }

      // ===== INFORMACIÓN PRINCIPAL =====
      pdf.setFontSize(12)
      pdf.setFont('helvetica', 'bold')
      pdf.setTextColor(0, 0, 0)
      
      // Tipo de QR
      const qrType = qrInfo.type === 'batch' ? 'Lote de Productos' : 'Insumo Individual'
      pdf.text('Tipo:', margin, currentY)
      pdf.setFont('helvetica', 'normal')
      pdf.text(qrType, margin + 25, currentY)
      currentY += 8

      // Estado/Disponibilidad
      const status = getAvailabilityStatus(qrInfo)
      pdf.setFont('helvetica', 'bold')
      pdf.text('Estado:', margin, currentY)
      pdf.setFont('helvetica', 'normal')
      pdf.setTextColor(status.color === '#dc2626' ? 220 : status.color === '#f59e0b' ? 245 : 5, 
                       status.color === '#dc2626' ? 38 : status.color === '#f59e0b' ? 158 : 150, 
                       status.color === '#dc2626' ? 38 : status.color === '#f59e0b' ? 11 : 105)
      pdf.text(status.label, margin + 25, currentY)
      pdf.setTextColor(0, 0, 0)
      currentY += 12

      // ===== INFORMACIÓN DEL INSUMO =====
      if (qrInfo.supply_info || qrInfo.supply_code || qrInfo.name || qrInfo.type === 'medical_supply') {
        pdf.setFontSize(14)
        pdf.setFont('helvetica', 'bold')
        pdf.text('Información del Insumo', margin, currentY)
        currentY += 10

        // Extraer información siguiendo la misma lógica que QRInfoDisplay.vue
        const infoFields = [
          { 
            label: 'Nombre:', 
            value: qrInfo.supply_code?.name || 
                  qrInfo.supply_info?.name || 
                  qrInfo.name || 
                  qrInfo.supply_info?.SupplyCode?.name ||
                  qrInfo.supply_info?.supply_code_name ||
                  'No disponible' 
          },
          { 
            label: 'Código:', 
            value: qrInfo.supply_code?.code_supplier || 
                  qrInfo.supplier_code || 
                  qrInfo.supply_info?.supply_code_code ||
                  qrInfo.supply_info?.SupplyCode?.code_supplier ||
                  qrInfo.code ||
                  'No disponible' 
          },
          { 
            label: 'Ubicación Actual:', 
            value: qrInfo.current_location || 
                  qrInfo.supply_info?.current_location || 
                  qrInfo.location ||
                  qrInfo.pavilion_name ||
                  qrInfo.supply_info?.pavilion_name ||
                  // ✅ ENCONTRADO: Ubicación está en traceability.current_location.name
                  qrInfo.traceability?.current_location?.name ||
                  'No asignada' 
          },
          { 
            label: 'Lote:', 
            value: qrInfo.supply_info?.batch?.id || 
                  qrInfo.batch_info?.id || 
                  qrInfo.batch_id || 
                  qrInfo.supply_info?.batch_number ||
                  qrInfo.batch_number ||
                  qrInfo.supply_info?.BatchID ||
                  'No disponible' 
          },
          { 
            label: 'Fecha de Vencimiento:', 
            value: formatDate(
              qrInfo.supply_info?.expiry_date || 
              qrInfo.batch_info?.expiry_date || 
              qrInfo.expiry_date ||
              qrInfo.supply_info?.batch?.expiry_date ||
              qrInfo.batch_info?.batch_expiry_date ||
              qrInfo.supply_info?.batch_expiry_date ||
              // ✅ ENCONTRADO: Fecha de vencimiento está en supply_info.batch.expiration_date
              qrInfo.supply_info?.batch?.expiration_date ||
              qrInfo.supply_info?.batch_info?.expiry_date ||
              qrInfo.batch_expiry_date ||
              qrInfo.expiration_date
            ) 
          },
          { 
            label: 'Cantidad:', 
            value: qrInfo.supply_info?.amount || 
                  qrInfo.amount || 
                  qrInfo.quantity ||
                  qrInfo.supply_info?.quantity ||
                  // ✅ ENCONTRADO: Cantidad está en supply_info.batch.amount
                  qrInfo.supply_info?.batch?.amount ||
                  qrInfo.batch_amount ||
                  'No disponible' 
          }
        ]

        // Agregar descripción solo si está disponible
        const descripcionValue = qrInfo.supply_code?.description || 
                                qrInfo.supply_info?.description || 
                                qrInfo.description ||
                                qrInfo.supply_info?.supply_code_description ||
                                qrInfo.supply_info?.SupplyCode?.description
        
        if (descripcionValue) {
          // Insertar después del código (posición 2)
          infoFields.splice(2, 0, {
            label: 'Descripción:',
            value: descripcionValue
          })
        }

        // Agregar fecha de fabricación solo si está disponible
        const fechaFabricacionValue = qrInfo.supply_info?.manufacture_date || 
                                     qrInfo.batch_info?.manufacture_date || 
                                     qrInfo.manufacture_date ||
                                     qrInfo.supply_info?.batch?.manufacture_date ||
                                     qrInfo.batch_info?.batch_manufacture_date ||
                                     qrInfo.supply_info?.batch_manufacture_date ||
                                     qrInfo.supply_info?.batch?.batch_manufacture_date ||
                                     qrInfo.supply_info?.batch_info?.manufacture_date ||
                                     qrInfo.batch_manufacture_date ||
                                     qrInfo.manufactured_date ||
                                     qrInfo.creation_date
        
        if (fechaFabricacionValue) {
          // Insertar después de la fecha de vencimiento
          const fechaVencimientoIndex = infoFields.findIndex(field => field.label === 'Fecha de Vencimiento:')
          if (fechaVencimientoIndex !== -1) {
            infoFields.splice(fechaVencimientoIndex + 1, 0, {
              label: 'Fecha de Fabricación:',
              value: formatDate(fechaFabricacionValue)
            })
          }
        }

        // Agregar información de proveedor si está disponible
        const proveedorValue = qrInfo.supply_info?.supplier_name || 
                              qrInfo.supplier_name || 
                              qrInfo.batch_info?.supplier_name ||
                              // ✅ ENCONTRADO: Proveedor está en supply_info.batch.supplier
                              qrInfo.supply_info?.batch?.supplier
        
        if (proveedorValue) {
          infoFields.push({
            label: 'Proveedor:',
            value: proveedorValue
          })
        }

        pdf.setFontSize(10)
        
        infoFields.forEach(field => {
          pdf.setFont('helvetica', 'bold')
          pdf.text(field.label, margin, currentY)
          
          pdf.setFont('helvetica', 'normal')
          // Manejar texto largo
          const maxWidth = contentWidth - 50
          const lines = pdf.splitTextToSize(field.value, maxWidth)
          pdf.text(lines, margin + 40, currentY)
          
          // Ajustar currentY basado en número de líneas
          currentY += (lines.length * 5) + 3
        })

        currentY += 5
      }

      // ===== INFORMACIÓN DEL LOTE =====
      if (qrInfo.batch_info || qrInfo.type === 'batch') {
        pdf.setFontSize(14)
        pdf.setFont('helvetica', 'bold')
        pdf.text('Información del Lote', margin, currentY)
        currentY += 10

        const batch = qrInfo.batch_info || qrInfo
        const batchFields = [
          { 
            label: 'Número de Lote:', 
            value: batch.batch_number || batch.id || 'No disponible' 
          },
          { 
            label: 'Fecha de Fabricación:', 
            value: formatDate(batch.manufacture_date) 
          },
          { 
            label: 'Fecha de Vencimiento:', 
            value: formatDate(batch.expiry_date) 
          },
          { 
            label: 'Cantidad Total:', 
            value: batch.total_quantity ? `${batch.total_quantity} unidades` : 'No disponible' 
          },
          { 
            label: 'Proveedor:', 
            value: batch.supplier_name || batch.supplier || 'No disponible' 
          },
          { 
            label: 'Centro Médico:', 
            value: batch.medical_center_name || batch.medical_center || 'No disponible' 
          }
        ]

        pdf.setFontSize(10)
        
        batchFields.forEach(field => {
          pdf.setFont('helvetica', 'bold')
          pdf.text(field.label, margin, currentY)
          
          pdf.setFont('helvetica', 'normal')
          const maxWidth = contentWidth - 50
          const lines = pdf.splitTextToSize(field.value, maxWidth)
          pdf.text(lines, margin + 40, currentY)
          
          currentY += (lines.length * 5) + 3
        })
      }

      // ===== FOOTER =====
      currentY = pageHeight - 30

      // Línea divisoria
      pdf.setLineWidth(0.3)
      pdf.setDrawColor(200, 200, 200)
      pdf.line(margin, currentY, pageWidth - margin, currentY)
      currentY += 8

      // Información de generación
      pdf.setFontSize(8)
      pdf.setFont('helvetica', 'normal')
      pdf.setTextColor(100, 100, 100)
      pdf.text(`Generado el: ${formatDate(new Date())}`, margin, currentY)
      pdf.text('MediTrack - Sistema de Trazabilidad Médica', margin, currentY + 5)

      // Información del usuario (si está disponible)
      pdf.text(`Sistema: MediTrack v1.0`, pageWidth - margin - 40, currentY)

      // ===== GUARDAR / RETORNAR PDF =====
      if (options.returnAsBase64) {
        // Retornar solo la parte base64 (sin el prefijo 'data:application/pdf;base64,')
        const dataUri = pdf.output('datauristring')
        return dataUri.split(',')[1]
      }

      const fileName = `QR-${qrInfo.qr_code}-${format(new Date(), 'yyyy-MM-dd-HHmm')}.pdf`
      pdf.save(fileName)

      return true

    } catch (err) {
      console.error('Error generating PDF:', err)
      error.value = 'Error al generar el PDF: ' + err.message
      return false
    } finally {
      isGenerating.value = false
    }
  }

  /**
   * Genera el PDF del QR y lo retorna como base64 (sin descargarlo)
   * @param {Object} qrInfo - Información del código QR
   * @returns {Promise<string|null>} Base64 del PDF, o null en caso de error
   */
  const generateQRPdfAsBase64 = async (qrInfo) => {
    return downloadQRAsPDF(qrInfo, { returnAsBase64: true })
  }

  return {
    downloadQRAsPDF,
    generateQRPdfAsBase64,
    isGenerating,
    error
  }
}
