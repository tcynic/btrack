import { Button, Card, CardHeader } from '../ui'
import { WeeklyBreakdown } from '../weekly'
import { formatDate } from '../../utils'
import type { ProjectWithStats } from '../../types'

interface ProjectDetailProps {
  project: ProjectWithStats
  onBack: () => void
}

export function ProjectDetail({ project, onBack }: ProjectDetailProps) {
  const progressPercent = project.totalPlannedHours > 0
    ? Math.min(100, (project.totalActualHours / project.totalPlannedHours) * 100)
    : 0

  return (
    <div>
      {/* Header */}
      <div className="mb-6">
        <Button variant="ghost" onClick={onBack} className="mb-4">
          <svg className="h-5 w-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 19l-7-7 7-7" />
          </svg>
          Back to Projects
        </Button>
        <div className="flex justify-between items-start">
          <div>
            <h2 className="text-3xl font-bold text-gray-900">{project.name}</h2>
            <p className="text-gray-500 mt-1">
              {formatDate(project.startDate)} - {formatDate(project.endDate)}
            </p>
          </div>
          {!project.isActive && (
            <span className="px-3 py-1 text-sm bg-gray-100 text-gray-600 rounded-full">
              Inactive
            </span>
          )}
        </div>
      </div>

      {/* Summary Cards */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
        <Card>
          <p className="text-sm text-gray-500 uppercase mb-1">Total Sold Hours</p>
          <p className="text-2xl font-bold text-gray-900">{project.totalSoldHours}</p>
        </Card>
        <Card>
          <p className="text-sm text-gray-500 uppercase mb-1">My Hours</p>
          <p className="text-2xl font-bold text-blue-600">{project.myHours}</p>
        </Card>
        <Card>
          <p className="text-sm text-gray-500 uppercase mb-1">Specialist Hours</p>
          <p className="text-2xl font-bold text-purple-600">{project.specialistHours}</p>
          <p className="text-xs text-gray-500 mt-1">
            ({((project.specialistHours / project.totalSoldHours) * 100).toFixed(0)}% of total)
          </p>
        </Card>
        <Card>
          <p className="text-sm text-gray-500 uppercase mb-1">Total Weeks</p>
          <p className="text-2xl font-bold text-gray-900">{project.totalWeeks}</p>
        </Card>
      </div>

      {/* Progress Overview */}
      <Card className="mb-6">
        <CardHeader title="Progress Overview" />
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-4">
          <div>
            <p className="text-sm text-gray-500 mb-1">Planned Hours</p>
            <p className="text-3xl font-bold text-blue-600">{project.totalPlannedHours}</p>
          </div>
          <div>
            <p className="text-sm text-gray-500 mb-1">Actual Hours</p>
            <p className="text-3xl font-bold text-gray-900">{project.totalActualHours}</p>
          </div>
          <div>
            <p className="text-sm text-gray-500 mb-1">Remaining</p>
            <p className={`text-3xl font-bold ${
              project.totalPlannedHours - project.totalActualHours >= 0 ? 'text-green-600' : 'text-red-600'
            }`}>
              {project.totalPlannedHours - project.totalActualHours}
            </p>
          </div>
        </div>

        {/* Progress Bar */}
        <div>
          <div className="flex justify-between text-sm text-gray-500 mb-2">
            <span>Progress</span>
            <span>{progressPercent.toFixed(1)}%</span>
          </div>
          <div className="w-full bg-gray-200 rounded-full h-3">
            <div
              className={`h-3 rounded-full transition-all ${
                progressPercent > 100 ? 'bg-red-500' : 'bg-blue-600'
              }`}
              style={{ width: `${Math.min(100, progressPercent)}%` }}
            />
          </div>
          {progressPercent > 100 && (
            <p className="text-sm text-red-600 mt-2">
              Over budget by {project.totalActualHours - project.totalPlannedHours} hours
            </p>
          )}
        </div>
      </Card>

      {/* Weekly Breakdown */}
      <Card padding="none">
        <div className="px-6 pt-6 pb-4">
          <CardHeader
            title="Weekly Breakdown"
            subtitle="Track planned vs actual hours for each week"
          />
          <p className="text-sm text-gray-500">
            Click on actual hours to edit (only available for past and current weeks)
          </p>
        </div>
        <WeeklyBreakdown projectId={project.id} />
      </Card>
    </div>
  )
}
