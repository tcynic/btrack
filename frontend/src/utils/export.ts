/**
 * Triggers a browser download of CSV data
 * @param csvContent - The CSV content as a string
 * @param filename - The filename for the download (should end with .csv)
 */
export function downloadCSV(csvContent: string, filename: string): void {
  // Ensure filename ends with .csv
  if (!filename.endsWith('.csv')) {
    filename += '.csv'
  }

  // Create blob with UTF-8 BOM for Excel compatibility
  const BOM = '\uFEFF'
  const blob = new Blob([BOM + csvContent], { type: 'text/csv;charset=utf-8;' })
  
  // Create download link
  const link = document.createElement('a')
  const url = URL.createObjectURL(blob)
  
  link.setAttribute('href', url)
  link.setAttribute('download', filename)
  link.style.visibility = 'hidden'
  
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  
  // Clean up the URL object
  URL.revokeObjectURL(url)
}

/**
 * Formats a date string for use in filenames
 * @param date - Date object or ISO string
 * @returns Formatted date string (YYYY-MM-DD)
 */
export function formatDateForFilename(date: Date | string = new Date()): string {
  const d = typeof date === 'string' ? new Date(date) : date
  return d.toISOString().split('T')[0]
}

/**
 * Generates a filename for exported data
 * @param prefix - Prefix for the filename (e.g., 'project', 'weekly-report')
 * @param suffix - Optional suffix (e.g., project name)
 * @returns Generated filename
 */
export function generateExportFilename(prefix: string, suffix?: string): string {
  const date = formatDateForFilename()
  const parts = [prefix, suffix, date].filter(Boolean)
  return `${parts.join('_')}.csv`
}
