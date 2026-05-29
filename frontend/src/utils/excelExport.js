import ExcelJS from 'exceljs'
import { parseDbDate } from './dateUtils'

/**
 * Exporta datos a un archivo Excel
 * @param {Array} data - Array de objetos con los datos a exportar
 * @param {Array} columns - Array de objetos con la configuración de columnas { key: string, label: string }
 * @param {string} filename - Nombre del archivo (sin extensión)
 */
export async function exportToExcel(data, columns, filename = 'export') {
  try {
    // Crear workbook y worksheet
    const workbook = new ExcelJS.Workbook()
    const worksheet = workbook.addWorksheet('Datos')

    // Agregar encabezados
    worksheet.columns = columns.map(col => ({
      header: col.label,
      key: col.key,
      width: Math.max(col.label.length, 15) // Ancho mínimo de 15 caracteres
    }))

    // Estilizar encabezados
    worksheet.getRow(1).font = { bold: true }
    worksheet.getRow(1).fill = {
      type: 'pattern',
      pattern: 'solid',
      fgColor: { argb: 'FFE0E0E0' }
    }

    // Agregar datos
    data.forEach(item => {
      const row = {}
      columns.forEach(col => {
        // Obtener el valor, manejando valores anidados con notación de punto
        let value = item
        const keys = col.key.split('.')
        for (const key of keys) {
          value = value?.[key]
        }
        
        // Formatear el valor según el tipo
        if (value === null || value === undefined) {
          row[col.key] = ''
        } else if (col.formatter && typeof col.formatter === 'function') {
          row[col.key] = col.formatter(value, item)
        } else if (value instanceof Date) {
          row[col.key] = value.toLocaleString('es-CL')
        } else {
          row[col.key] = String(value)
        }
      })
      worksheet.addRow(row)
    })

    // Generar nombre de archivo con fecha
    const date = new Date()
    const dateStr = date.toISOString().split('T')[0]
    const fullFilename = `${filename}_${dateStr}.xlsx`

    // Descargar el archivo
    const buffer = await workbook.xlsx.writeBuffer()
    const blob = new Blob([buffer], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' })
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = fullFilename
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
    
    return true
  } catch (error) {
    console.error('Error al exportar a Excel:', error)
    throw error
  }
}

/**
 * Formatea una fecha para Excel
 */
export function formatDateForExcel(date) {
  if (!date) return ''
  return parseDbDate(date).toLocaleString('es-CL', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

/**
 * Formatea un estado para Excel
 */
export function formatStatusForExcel(status) {
  const statusMap = {
    // Estados de transferencia (valores canónicos del backend)
    'pendiente': 'Pendiente',
    'en_transito': 'En Tránsito',
    'recibido': 'Recibido',
    'rechazado': 'Rechazado',
    // Valores legacy en inglés (compatibilidad con datos históricos)
    'in_transit': 'En Tránsito',
    'completed': 'Completado',
    'cancelled': 'Cancelado',
    // Estados de insumo médico
    'disponible': 'Disponible',
    'pendiente_retiro': 'Pendiente de Retiro',
    'en_camino_a_pabellon': 'En Camino a Pabellón',
    'en_camino_a_bodega': 'En Camino a Bodega',
    'recepcionado': 'Recepcionado',
    'consumido': 'Consumido',
    'vencido': 'Vencido',
    // Estados de solicitud de insumo
    'pendiente_pavedad': 'Pendiente Pavedad',
    'asignado_bodega': 'Asignado a Bodega',
    'en_proceso': 'En Proceso',
    'aprobado': 'Aprobado',
    'completado': 'Completado',
    'devuelto_al_solicitante': 'Devuelto al Solicitante',
    'devuelto': 'Devuelto al Solicitante',
    'devuelto_al_encargado': 'Devuelto al Encargado',
    'parcialmente_devuelto_al_encargado': 'Parcialmente Devuelto al Encargado'
  }
  return statusMap[status] || status
}

