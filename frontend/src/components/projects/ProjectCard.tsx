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
    on_track: { label: 'On Track', color: 'text-[var(--ld-green)]', bg: 'bg-transparent', dot: 'bg-[var(--ld-green)]' },
    at_risk: { label: 'At Risk', color: 'text-[var(--ld-orange)]', bg: 'bg-transparent', dot: 'bg-[var(--ld-orange)]' },
    over_budget: { label: 'Over Budget', color: 'text-[var(--ld-pink)]', bg: 'bg-transparent', dot: 'bg-[var(--ld-pink)]' },
    completed: { label: 'Completed', color: 'text-ld-primary', bg: 'bg-transparent', dot: 'bg-ld-primary' },
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
            <h3 className="text-lg font-semibold text-ld-text">
              {project.name}
            </h3>
            {project.isPersistent && (
              <span className="px-2 py-0.5 text-xs border border-[var(--ld-purple)] text-[var(--ld-purple)] rounded-full">
                Persistent
              </span>
            )}
          </div>
          <p className="text-sm text-ld-muted mb-1">
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
          <span className="px-2 py-1 text-xs border border-ld-border text-ld-muted rounded-full">
            Inactive
          </span>
        )}
      </div>

      <div className="grid grid-cols-3 gap-4 mb-4">
        <div>
          <p className="text-xs text-ld-muted uppercase">My Hours</p>
          <p className="text-lg font-semibold text-ld-text">{project.myHours}</p>
        </div>
        <div>
          <p className="text-xs text-ld-muted uppercase">Specialist</p>
          <p className="text-lg font-semibold text-ld-text">{project.specialistHours}</p>
        </div>
        <div>
          <p className="text-xs text-ld-muted uppercase">Weeks</p>
          <p className="text-lg font-semibold text-ld-text">{project.totalWeeks}</p>
        </div>
      </div>

      {/* Progress bar */}
      <div className="mb-2">
        <div className="flex justify-between text-xs text-ld-muted mb-1">
          <span>{project.totalActualHours} / {project.totalPlannedHours} hrs logged</span>
          <span className={hoursRemaining >= 0 ? 'text-[var(--ld-green)]' : 'text-[var(--ld-pink)]'}>
            {hoursRemaining >= 0 ? `${hoursRemaining} remaining` : `${Math.abs(hoursRemaining)} over`}
          </span>
        </div>
        <div className="w-full bg-ld-surface2 rounded-full h-2">
          <div
            className={`h-2 rounded-full transition-all ${
              progressPercent > 100 ? 'bg-[var(--ld-pink)]' : 'bg-ld-primary'
            }`}
            style={{ width: `${Math.min(100, progressPercent)}%` }}
          />
        </div>
      </div>
    </Card>
  )
}
