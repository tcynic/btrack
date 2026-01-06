import { useState, useCallback } from 'react'
import type { CreateProjectInput, WeeklyEntry } from '../types'
import { CalculateDistribution } from '../../wailsjs/go/main/App'

export function useDistributionPreview() {
  const [preview, setPreview] = useState<WeeklyEntry[]>([])
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const calculatePreview = useCallback(async (input: CreateProjectInput) => {
    if (!input.name || !input.startDate || !input.endDate || input.totalSoldHours <= 0) {
      setPreview([])
      return
    }

    setIsLoading(true)
    setError(null)
    try {
      const data = await CalculateDistribution(input)
      setPreview(data as WeeklyEntry[])
    } catch (err) {
      setError(String(err))
      setPreview([])
    } finally {
      setIsLoading(false)
    }
  }, [])

  const clearPreview = useCallback(() => {
    setPreview([])
    setError(null)
  }, [])

  return {
    preview,
    isLoading,
    error,
    calculatePreview,
    clearPreview,
  }
}
