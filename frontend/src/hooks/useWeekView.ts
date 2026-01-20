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
import { useQuery } from './useQuery'

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

interface WeekData {
  meetings: MeetingWithProject[]
  entries: WeeklyEntryWithProject[]
}

export function useWeekView() {
  const [selectedWeek, setSelectedWeek] = useState<Date>(getCurrentWeekMonday())

  const weekStartDate = useMemo(() => formatDate(selectedWeek), [selectedWeek])

  const weekDataQuery = useQuery<WeekData>({
    queryFn: useCallback(async () => {
      const [meetings, entries] = await Promise.all([
        GetMeetingsByWeek(weekStartDate),
        GetWeeklyEntriesByWeek(weekStartDate)
      ])
      return {
        meetings: (meetings || []) as MeetingWithProject[],
        entries: (entries || []) as WeeklyEntryWithProject[],
      }
    }, [weekStartDate]),
    initialData: { meetings: [], entries: [] },
  })

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
    await UpdateActualHours(input)
    await weekDataQuery.refetch()
  }, [weekDataQuery.refetch])

  return {
    selectedWeek,
    weekStartDate,
    meetings: weekDataQuery.data?.meetings || [],
    entries: weekDataQuery.data?.entries || [],
    isLoading: weekDataQuery.isLoading,
    error: weekDataQuery.error,
    loadWeekData: weekDataQuery.refetch,
    navigateWeek,
    goToCurrentWeek,
    updateActualHours,
  }
}
