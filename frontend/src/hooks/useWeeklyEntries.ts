import { useState, useCallback } from 'react'
import type { WeeklyEntryWithStatus, UpdateActualHoursInput } from '../types'
import { GetWeeklyEntries, UpdateActualHours } from '../../wailsjs/go/main/App'

export function useWeeklyEntries(projectId: number) {
  const [entries, setEntries] = useState<WeeklyEntryWithStatus[]>([])
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const loadEntries = useCallback(async () => {
    setIsLoading(true)
    setError(null)
    try {
      const data = await GetWeeklyEntries(projectId)
      setEntries(data as WeeklyEntryWithStatus[])
    } catch (err) {
      setError(String(err))
    } finally {
      setIsLoading(false)
    }
  }, [projectId])

  const updateActualHours = useCallback(async (input: UpdateActualHoursInput) => {
    setIsLoading(true)
    setError(null)
    try {
      const updated = await UpdateActualHours(input)
      setEntries((prev) =>
        prev.map((e) => (e.id === input.entryId ? (updated as WeeklyEntryWithStatus) : e))
      )
      return updated as WeeklyEntryWithStatus
    } catch (err) {
      setError(String(err))
      throw err
    } finally {
      setIsLoading(false)
    }
  }, [])

  return {
    entries,
    isLoading,
    error,
    loadEntries,
    updateActualHours,
  }
}
