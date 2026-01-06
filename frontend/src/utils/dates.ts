/**
 * Format a date string (YYYY-MM-DD) to a readable format
 */
export function formatDate(dateStr: string): string {
  const date = new Date(dateStr + 'T00:00:00')
  return date.toLocaleDateString('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
  })
}

/**
 * Format a week start date to show the week range
 */
export function formatWeekRange(weekStartDate: string): string {
  const start = new Date(weekStartDate + 'T00:00:00')
  const end = new Date(start)
  end.setDate(end.getDate() + 6)

  const startMonth = start.toLocaleDateString('en-US', { month: 'short' })
  const endMonth = end.toLocaleDateString('en-US', { month: 'short' })

  if (startMonth === endMonth) {
    return `${startMonth} ${start.getDate()}-${end.getDate()}`
  }
  return `${startMonth} ${start.getDate()} - ${endMonth} ${end.getDate()}`
}

/**
 * Get today's date in YYYY-MM-DD format
 */
export function getTodayString(): string {
  return new Date().toISOString().split('T')[0]
}

/**
 * Get a date N days from today in YYYY-MM-DD format
 */
export function getDateFromToday(days: number): string {
  const date = new Date()
  date.setDate(date.getDate() + days)
  return date.toISOString().split('T')[0]
}

/**
 * Check if a week start date is in the past
 */
export function isWeekInPast(weekStartDate: string): boolean {
  const weekStart = new Date(weekStartDate + 'T00:00:00')
  const today = new Date()
  today.setHours(0, 0, 0, 0)

  // Get this week's Monday
  const dayOfWeek = today.getDay()
  const daysToMonday = dayOfWeek === 0 ? 6 : dayOfWeek - 1
  const thisMonday = new Date(today)
  thisMonday.setDate(today.getDate() - daysToMonday)

  return weekStart < thisMonday
}

/**
 * Check if a week start date is the current week
 */
export function isCurrentWeek(weekStartDate: string): boolean {
  const weekStart = new Date(weekStartDate + 'T00:00:00')
  const today = new Date()
  today.setHours(0, 0, 0, 0)

  // Get this week's Monday
  const dayOfWeek = today.getDay()
  const daysToMonday = dayOfWeek === 0 ? 6 : dayOfWeek - 1
  const thisMonday = new Date(today)
  thisMonday.setDate(today.getDate() - daysToMonday)

  return weekStart.getTime() === thisMonday.getTime()
}
