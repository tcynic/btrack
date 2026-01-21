import { useEffect, useState } from 'react'
import { Card } from '../ui'
import { TrendChart } from './TrendChart'
import { VarianceTable } from './VarianceTable'
import { useReports } from '../../hooks/useReports'

export function ReportsView() {
  const {
    monthlyTrends,
    varianceReport,
    isLoading,
    error,
    loadMonthlyTrends,
    loadVarianceReport,
  } = useReports()

  const [monthsBack, setMonthsBack] = useState(6)
  const [dateRange, setDateRange] = useState({
    start: getDefaultStartDate(),
    end: getDefaultEndDate(),
  })

  useEffect(() => {
    loadMonthlyTrends(monthsBack)
  }, [monthsBack])

  useEffect(() => {
    loadVarianceReport(dateRange.start, dateRange.end)
  }, [dateRange])

  function getDefaultStartDate() {
    const date = new Date()
    date.setMonth(date.getMonth() - 3)
    return date.toISOString().split('T')[0]
  }

  function getDefaultEndDate() {
    return new Date().toISOString().split('T')[0]
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold text-ld-text">Reports & Analytics</h1>
      </div>

      {error && (
        <div className="p-4 bg-transparent border border-[var(--ld-pink)] rounded-lg text-[var(--ld-pink)]">
          {error}
        </div>
      )}

      {/* Monthly Trends */}
      <Card>
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-lg font-semibold text-ld-text">Monthly Trends</h2>
          <div className="flex items-center gap-2">
            <label className="text-sm text-ld-muted">Months:</label>
            <select
              value={monthsBack}
              onChange={(e) => setMonthsBack(parseInt(e.target.value))}
              className="px-3 py-1 border border-ld-border bg-ld-surface text-ld-text rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-ld-primary"
            >
              <option value={3}>3 months</option>
              <option value={6}>6 months</option>
              <option value={12}>12 months</option>
            </select>
          </div>
        </div>
        {isLoading ? (
          <div className="flex items-center justify-center h-64">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-ld-primary"></div>
          </div>
        ) : (
          <TrendChart data={monthlyTrends} />
        )}
      </Card>

      {/* Variance Report */}
      <Card>
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-lg font-semibold text-ld-text">Planned vs Actual</h2>
          <div className="flex items-center gap-3">
            <label className="text-sm text-ld-muted">From:</label>
            <input
              type="date"
              value={dateRange.start}
              onChange={(e) => setDateRange({ ...dateRange, start: e.target.value })}
              className="px-3 py-1 border border-ld-border bg-ld-surface text-ld-text rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-ld-primary"
            />
            <label className="text-sm text-ld-muted">To:</label>
            <input
              type="date"
              value={dateRange.end}
              onChange={(e) => setDateRange({ ...dateRange, end: e.target.value })}
              className="px-3 py-1 border border-ld-border bg-ld-surface text-ld-text rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-ld-primary"
            />
          </div>
        </div>
        {isLoading ? (
          <div className="flex items-center justify-center h-64">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-ld-primary"></div>
          </div>
        ) : (
          <VarianceTable data={varianceReport} />
        )}
      </Card>
    </div>
  )
}
