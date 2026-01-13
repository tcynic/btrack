import { useEffect, useMemo } from 'react'
import { Card, CardHeader, Button } from '../ui'
import { WeekMeetingList } from './WeekMeetingList'
import { WeekHoursTable } from './WeekHoursTable'
import { useWeekView } from '../../hooks'

// Helper to format week range (e.g., "Jan 13 - Jan 19, 2026")
function formatWeekRange(weekStartDate: string): string {
  const startDate = new Date(weekStartDate + 'T00:00:00')
  const endDate = new Date(startDate)
  endDate.setDate(startDate.getDate() + 6)
  
  const startMonth = startDate.toLocaleDateString('en-US', { month: 'short' })
  const startDay = startDate.getDate()
  const endMonth = endDate.toLocaleDateString('en-US', { month: 'short' })
  const endDay = endDate.getDate()
  const year = endDate.getFullYear()
  
  if (startMonth === endMonth) {
    return `${startMonth} ${startDay} - ${endDay}, ${year}`
  } else {
    return `${startMonth} ${startDay} - ${endMonth} ${endDay}, ${year}`
  }
}

// Helper to check if a week is the current week
function isCurrentWeek(weekStartDate: string): boolean {
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  
  const weekStart = new Date(weekStartDate + 'T00:00:00')
  const weekEnd = new Date(weekStart)
  weekEnd.setDate(weekStart.getDate() + 6)
  
  return today >= weekStart && today <= weekEnd
}

export function WeekView() {
  const {
    weekStartDate,
    meetings,
    entries,
    isLoading,
    loadWeekData,
    navigateWeek,
    goToCurrentWeek,
    updateActualHours,
  } = useWeekView()

  useEffect(() => {
    loadWeekData()
  }, [loadWeekData])

  const isCurrent = useMemo(() => isCurrentWeek(weekStartDate), [weekStartDate])
  const weekRange = useMemo(() => formatWeekRange(weekStartDate), [weekStartDate])

  return (
    <div>
      {/* Header with Week Navigation */}
      <div className="mb-6">
        <div className="flex items-center justify-between mb-4">
          <div>
            <h2 className="text-2xl font-bold text-gray-900">Week View</h2>
            <p className="text-gray-500 mt-1">
              {weekRange}
              {isCurrent && (
                <span className="ml-2 px-2 py-0.5 text-xs bg-blue-100 text-blue-800 rounded-full">
                  Current Week
                </span>
              )}
            </p>
          </div>
          <div className="flex items-center gap-2">
            <Button
              variant="ghost"
              size="sm"
              onClick={() => navigateWeek('prev')}
              disabled={isLoading}
            >
              <svg
                className="w-5 h-5"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M15 19l-7-7 7-7"
                />
              </svg>
              Previous
            </Button>
            <Button
              variant="ghost"
              size="sm"
              onClick={goToCurrentWeek}
              disabled={isLoading || isCurrent}
            >
              Today
            </Button>
            <Button
              variant="ghost"
              size="sm"
              onClick={() => navigateWeek('next')}
              disabled={isLoading}
            >
              Next
              <svg
                className="w-5 h-5"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M9 5l7 7-7 7"
                />
              </svg>
            </Button>
          </div>
        </div>
      </div>

      {/* Content Grid */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Meetings Section */}
        <Card>
          <CardHeader
            title="Meetings"
            subtitle="Scheduled meetings for this week"
          />
          {isLoading ? (
            <div className="flex items-center justify-center py-12">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
            </div>
          ) : (
            <WeekMeetingList meetings={meetings} weekStartDate={weekStartDate} />
          )}
        </Card>

        {/* Hours Section */}
        <Card padding="none">
          <div className="px-6 pt-6 pb-4">
            <CardHeader
              title="Hours"
              subtitle="Track your hours across projects"
            />
            <p className="text-sm text-gray-500 mt-2">
              {isCurrent || entries.some(e => e.isPastWeek)
                ? 'Click actual hours to edit'
                : 'Future weeks cannot be edited'}
            </p>
          </div>
          {isLoading ? (
            <div className="flex items-center justify-center py-12">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
            </div>
          ) : (
            <WeekHoursTable
              entries={entries}
              onUpdateActualHours={async (entryId, actualHours) => {
                await updateActualHours({ entryId, actualHours })
              }}
              isCurrentWeek={isCurrent}
            />
          )}
        </Card>
      </div>
    </div>
  )
}
