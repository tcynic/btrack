import { useState } from 'react'
import {
  GetMonthlyTrends,
  GetVarianceReport,
  GetCapacityUtilization,
} from '../../wailsjs/go/main/App'

export interface MonthlyTrend {
  month: string
  plannedHours: number
  actualHours: number
  projectCount: number
  variance: number
}

export interface VarianceReport {
  projectId: number
  projectName: string
  planned: number
  actual: number
  variance: number
  percentage: number
}

export interface CapacityWeek {
  weekStart: string
  plannedHours: number
  actualHours: number
  utilization: number
  projectCount: number
}

export function useReports() {
  const [monthlyTrends, setMonthlyTrends] = useState<MonthlyTrend[]>([])
  const [varianceReport, setVarianceReport] = useState<VarianceReport[]>([])
  const [capacityData, setCapacityData] = useState<CapacityWeek[]>([])
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const loadMonthlyTrends = async (monthsBack: number = 6) => {
    setIsLoading(true)
    setError(null)
    try {
      const data = await GetMonthlyTrends(monthsBack)
      setMonthlyTrends(data || [])
    } catch (err) {
      console.error('Failed to load monthly trends:', err)
      setError('Failed to load monthly trends')
    } finally {
      setIsLoading(false)
    }
  }

  const loadVarianceReport = async (startDate: string, endDate: string) => {
    setIsLoading(true)
    setError(null)
    try {
      const data = await GetVarianceReport(startDate, endDate)
      setVarianceReport(data || [])
    } catch (err) {
      console.error('Failed to load variance report:', err)
      setError('Failed to load variance report')
    } finally {
      setIsLoading(false)
    }
  }

  const loadCapacityUtilization = async (startDate: string, endDate: string) => {
    setIsLoading(true)
    setError(null)
    try {
      const data = await GetCapacityUtilization(startDate, endDate)
      setCapacityData(data || [])
    } catch (err) {
      console.error('Failed to load capacity utilization:', err)
      setError('Failed to load capacity utilization')
    } finally {
      setIsLoading(false)
    }
  }

  return {
    monthlyTrends,
    varianceReport,
    capacityData,
    isLoading,
    error,
    loadMonthlyTrends,
    loadVarianceReport,
    loadCapacityUtilization,
  }
}
