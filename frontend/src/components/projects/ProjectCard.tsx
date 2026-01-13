import type { ProjectWithStats } from '../../types'
import { Card } from '../ui'
import { formatDate } from '../../utils'

interface ProjectCardProps {
  project: ProjectWithStats
  onClick: () => void
}

export function ProjectCard({ project, onClick }: ProjectCardProps) {
  const hoursRemaining = project.totalPlannedHours - project.totalActualHours
  const progressPercent = project.totalPlannedHours > 0
    ? Math.min(100, (project.totalActualHours / project.totalPlannedHours) * 100)
    : 0

  const healthConfig: Record<string, { label: string; color: string; bg: string; dot: string }> = {
    on_track: { label: 'On Track', color: 'text-green-700', bg: 'bg-green-100', dot: 'bg-green-500' },
    at_risk: { label: 'At Risk', color: 'text-yellow-700', bg: 'bg-yellow-100', dot: 'bg-yellow-500' },
    over_budget: { label: 'Over Budget', color: 'text-red-700', bg: 'bg-red-100', dot: 'bg-red-500' },
    completed: { label: 'Completed', color: 'text-blue-700', bg: 'bg-blue-100', dot: 'bg-blue-500' },
  }

  const health = healthConfig[project.health?.status as string || 'on_track']

  return (
    <Card
      className="cursor-pointer hover:shadow-md transition-shadow"
      onClick={onClick}
    >
      <div className="flex justify-between items-start mb-3">
        <div className="flex-1">
          <div className="flex items-center gap-2 mb-1">
            <h3 className="text-lg font-semibold text-gray-900">
              {project.name}
            </h3>
            {project.isPersistent && (
              <span className="px-2 py-0.5 text-xs bg-purple-100 text-purple-800 rounded-full">
                Persistent
              </span>
            )}
          </div>
          <p className="text-sm text-gray-500 mb-1">
            {formatDate(project.startDate)} - {formatDate(project.endDate)}
          </p>
          <div className="flex items-center gap-1.5" title={project.health?.message}>
            <span className={`w-2 h-2 rounded-full ${health.dot}`}></span>
            <span className={`text-xs font-medium ${health.color}`}>
              {health.label}
            </span>
          </div>
        </div>
        {!project.isActive && (
          <span className="px-2 py-1 text-xs bg-gray-100 text-gray-600 rounded-full">
            Inactive
          </span>
        )}
      </div>

      <div className="grid grid-cols-3 gap-4 mb-4">
        <div>
          <p className="text-xs text-gray-500 uppercase">My Hours</p>
          <p className="text-lg font-semibold text-gray-900">{project.myHours}</p>
        </div>
        <div>
          <p className="text-xs text-gray-500 uppercase">Specialist</p>
          <p className="text-lg font-semibold text-gray-900">{project.specialistHours}</p>
        </div>
        <div>
          <p className="text-xs text-gray-500 uppercase">Weeks</p>
          <p className="text-lg font-semibold text-gray-900">{project.totalWeeks}</p>
        </div>
      </div>

      {/* Progress bar */}
      <div className="mb-2">
        <div className="flex justify-between text-xs text-gray-500 mb-1">
          <span>{project.totalActualHours} / {project.totalPlannedHours} hrs logged</span>
          <span className={hoursRemaining >= 0 ? 'text-green-600' : 'text-red-600'}>
            {hoursRemaining >= 0 ? `${hoursRemaining} remaining` : `${Math.abs(hoursRemaining)} over`}
          </span>
        </div>
        <div className="w-full bg-gray-200 rounded-full h-2">
          <div
            className={`h-2 rounded-full transition-all ${
              progressPercent > 100 ? 'bg-red-500' : 'bg-blue-600'
            }`}
            style={{ width: `${Math.min(100, progressPercent)}%` }}
          />
        </div>
      </div>
    </Card>
  )
}
