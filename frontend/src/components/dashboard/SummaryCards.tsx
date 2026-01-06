import { Card } from '../ui'
import type { DashboardSummary } from '../../types'

interface SummaryCardsProps {
  summary: DashboardSummary | null
}

export function SummaryCards({ summary }: SummaryCardsProps) {
  if (!summary) {
    return (
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
        {[1, 2, 3, 4].map((i) => (
          <Card key={i} className="animate-pulse">
            <div className="h-4 bg-gray-200 rounded w-24 mb-2"></div>
            <div className="h-8 bg-gray-200 rounded w-16"></div>
          </Card>
        ))}
      </div>
    )
  }

  const cards = [
    {
      label: 'Active Projects',
      value: summary.totalActiveProjects,
      color: 'text-blue-600',
      bgColor: 'bg-blue-50',
      icon: (
        <svg className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
        </svg>
      ),
    },
    {
      label: 'Planned This Week',
      value: summary.totalPlannedThisWeek,
      suffix: 'hrs',
      color: 'text-indigo-600',
      bgColor: 'bg-indigo-50',
      icon: (
        <svg className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
        </svg>
      ),
    },
    {
      label: 'Actual This Week',
      value: summary.totalActualThisWeek,
      suffix: 'hrs',
      color: summary.totalActualThisWeek > summary.totalPlannedThisWeek ? 'text-red-600' : 'text-green-600',
      bgColor: summary.totalActualThisWeek > summary.totalPlannedThisWeek ? 'bg-red-50' : 'bg-green-50',
      icon: (
        <svg className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
      ),
    },
    {
      label: 'Planned Next Week',
      value: summary.totalPlannedNextWeek,
      suffix: 'hrs',
      color: 'text-purple-600',
      bgColor: 'bg-purple-50',
      icon: (
        <svg className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 7l5 5m0 0l-5 5m5-5H6" />
        </svg>
      ),
    },
  ]

  return (
    <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
      {cards.map((card) => (
        <Card key={card.label}>
          <div className="flex items-center">
            <div className={`p-3 rounded-lg ${card.bgColor} ${card.color} mr-4`}>
              {card.icon}
            </div>
            <div>
              <p className="text-sm text-gray-500">{card.label}</p>
              <p className={`text-2xl font-bold ${card.color}`}>
                {card.value}
                {card.suffix && <span className="text-sm font-normal ml-1">{card.suffix}</span>}
              </p>
            </div>
          </div>
        </Card>
      ))}
    </div>
  )
}
