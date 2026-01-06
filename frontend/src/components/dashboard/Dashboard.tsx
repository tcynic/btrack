import { useEffect } from 'react'
import { Card, CardHeader } from '../ui'
import { SummaryCards } from './SummaryCards'
import { WeeklyChart } from './WeeklyChart'
import { WeeklyTable } from './WeeklyTable'
import { useDashboard } from '../../hooks'

export function Dashboard() {
  const { summary, weekData, isLoading, refreshDashboard } = useDashboard()

  useEffect(() => {
    refreshDashboard()
  }, [refreshDashboard])

  return (
    <div>
      <div className="mb-6">
        <h2 className="text-2xl font-bold text-gray-900">Capacity Overview</h2>
        <p className="text-gray-500">Monitor your weekly bandwidth across all active projects</p>
      </div>

      <SummaryCards summary={summary} />

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <Card>
          <CardHeader
            title="Weekly Capacity"
            subtitle="Planned vs actual hours per week"
          />
          {isLoading ? (
            <div className="h-80 flex items-center justify-center">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
            </div>
          ) : (
            <WeeklyChart data={weekData} />
          )}
        </Card>

        <Card padding="none">
          <div className="px-6 pt-6">
            <CardHeader
              title="Week by Week"
              subtitle="Detailed breakdown of hours"
            />
          </div>
          <div className="max-h-96 overflow-y-auto">
            {isLoading ? (
              <div className="h-40 flex items-center justify-center">
                <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
              </div>
            ) : (
              <WeeklyTable data={weekData} />
            )}
          </div>
        </Card>
      </div>
    </div>
  )
}
