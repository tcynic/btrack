import { useState, useCallback } from 'react'
import {
  GetMonthlyTrends,
  GetVarianceReport,
  GetCapacityUtilization,
} from '../../wailsjs/go/main/App'
import { useQuery } from './useQuery'

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
  const [monthsBack, setMonthsBack] = useState(6)
  const [dateRange, setDateRange] = useState<{ start: string; end: string } | null>(null)

  const monthlyTrendsQuery = useQuery<MonthlyTrend[]>({
    queryFn: useCallback(() => GetMonthlyTrends(monthsBack), [monthsBack]),
    initialData: [],
    enabled: true,
  })

  const varianceReportQuery = useQuery<VarianceReport[]>({
    queryFn: useCallback(
      () => dateRange ? GetVarianceReport(dateRange.start, dateRange.end) : Promise.resolve([]),
      [dateRange]
    ),
    initialData: [],
    enabled: !!dateRange,
  })

  const capacityQuery = useQuery<CapacityWeek[]>({
    queryFn: useCallback(
      () => dateRange ? GetCapacityUtilization(dateRange.start, dateRange.end) : Promise.resolve([]),
      [dateRange]
    ),
    initialData: [],
    enabled: !!dateRange,
  })

  const loadMonthlyTrends = useCallback(async (months: number = 6) => {
    setMonthsBack(months)
  }, [])

  const loadVarianceReport = useCallback(async (startDate: string, endDate: string) => {
    setDateRange({ start: startDate, end: endDate })
  }, [])

  const loadCapacityUtilization = useCallback(async (startDate: string, endDate: string) => {
    setDateRange({ start: startDate, end: endDate })
  }, [])

  const isLoading = monthlyTrendsQuery.isLoading || varianceReportQuery.isLoading || capacityQuery.isLoading
  const error = monthlyTrendsQuery.error || varianceReportQuery.error || capacityQuery.error

  return {
    monthlyTrends: monthlyTrendsQuery.data || [],
    varianceReport: varianceReportQuery.data || [],
    capacityData: capacityQuery.data || [],
    isLoading,
    error,
    loadMonthlyTrends,
    loadVarianceReport,
    loadCapacityUtilization,
  }
}
