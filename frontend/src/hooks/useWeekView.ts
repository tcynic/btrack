import { useState, useCallback, useMemo } from 'react'
import type { 
  MeetingWithProject, 
  WeeklyEntryWithProject,
  UpdateActualHoursInput 
} from '../types'
import { 
  GetMeetingsByWeek, 
  GetWeeklyEntriesByWeek,
  UpdateActualHours 
} from '../../wailsjs/go/main/App'

// Helper to get Monday of current week
function getCurrentWeekMonday(): Date {
  const today = new Date()
  const day = today.getDay()
  const diff = today.getDate() - day + (day === 0 ? -6 : 1)
  const monday = new Date(today.setDate(diff))
  monday.setHours(0, 0, 0, 0)
  return monday
}

// Format date as YYYY-MM-DD
function formatDate(date: Date): string {
  return date.toISOString().split('T')[0]
}

export function useWeekView() {
  const [selectedWeek, setSelectedWeek] = useState<Date>(getCurrentWeekMonday())
  const [meetings, setMeetings] = useState<MeetingWithProject[]>([])
  const [entries, setEntries] = useState<WeeklyEntryWithProject[]>([])
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const weekStartDate = useMemo(() => formatDate(selectedWeek), [selectedWeek])

  const loadWeekData = useCallback(async () => {
    setIsLoading(true)
    setError(null)
    try {
      const [meetingsData, entriesData] = await Promise.all([
        GetMeetingsByWeek(weekStartDate),
        GetWeeklyEntriesByWeek(weekStartDate)
      ])
      setMeetings((meetingsData || []) as MeetingWithProject[])
      setEntries((entriesData || []) as WeeklyEntryWithProject[])
    } catch (err) {
      setError(String(err))
      console.error('Failed to load week data:', err)
    } finally {
      setIsLoading(false)
    }
  }, [weekStartDate])

  const navigateWeek = useCallback((direction: 'prev' | 'next') => {
    setSelectedWeek((current) => {
      const newDate = new Date(current)
      newDate.setDate(newDate.getDate() + (direction === 'next' ? 7 : -7))
      return newDate
    })
  }, [])

  const goToCurrentWeek = useCallback(() => {
    setSelectedWeek(getCurrentWeekMonday())
  }, [])

  const updateActualHours = useCallback(async (input: UpdateActualHoursInput) => {
    setIsLoading(true)
    setError(null)
    try {
      await UpdateActualHours(input)
      // Reload the week data to get updated entry
      await loadWeekData()
    } catch (err) {
      setError(String(err))
      throw err
    } finally {
      setIsLoading(false)
    }
  }, [loadWeekData])

  return {
    selectedWeek,
    weekStartDate,
    meetings,
    entries,
    isLoading,
    error,
    loadWeekData,
    navigateWeek,
    goToCurrentWeek,
    updateActualHours,
  }
}
