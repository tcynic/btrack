import { useEffect, useState } from 'react'
import { Button, SearchInput } from '../ui'
import { ProjectCard } from './ProjectCard'
import { ProjectCreateModal } from './ProjectCreateModal'
import { useProjects } from '../../hooks'
import type { ProjectWithStats } from '../../types'

interface ProjectListProps {
  onProjectSelect: (project: ProjectWithStats) => void
}

export function ProjectList({ onProjectSelect }: ProjectListProps) {
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false)
  const [showInactive, setShowInactive] = useState(false)
  const [searchQuery, setSearchQuery] = useState('')
  const [statusFilter, setStatusFilter] = useState<'all' | 'on_track' | 'over_budget' | 'completed'>('all')
  const { projects, isLoading, loadProjects, createProject } = useProjects()

  useEffect(() => {
    loadProjects(!showInactive)
  }, [loadProjects, showInactive])

  // Helper to determine project status
  const getProjectStatus = (project: ProjectWithStats) => {
    const today = new Date()
    const endDate = new Date(project.endDate)
    const isPastEnd = today > endDate
    
    if (isPastEnd) return 'completed'
    if (project.totalActualHours > project.totalPlannedHours) return 'over_budget'
    return 'on_track'
  }

  // Filter projects by search query and status
  let filteredProjects = projects
  
  if (searchQuery) {
    filteredProjects = filteredProjects.filter((p) =>
      p.name.toLowerCase().includes(searchQuery.toLowerCase())
    )
  }
  
  if (statusFilter !== 'all') {
    filteredProjects = filteredProjects.filter((p) => getProjectStatus(p) === statusFilter)
  }

  const activeProjects = filteredProjects.filter((p) => p.isActive)
  const inactiveProjects = filteredProjects.filter((p) => !p.isActive)

  return (
    <div>
      <div className="mb-6">
        <div className="flex justify-between items-center mb-4">
          <div>
            <h2 className="text-2xl font-bold text-gray-900">Projects</h2>
            <p className="text-gray-500">Manage your project hours and track progress</p>
          </div>
          <Button onClick={() => setIsCreateModalOpen(true)}>
            New Project
          </Button>
        </div>
        <div className="flex items-center gap-3">
          <SearchInput
            value={searchQuery}
            onChange={setSearchQuery}
            placeholder="Search projects..."
            className="flex-1 max-w-md"
          />
          <select
            value={statusFilter}
            onChange={(e) => setStatusFilter(e.target.value as any)}
            className="rounded-md border-gray-300 text-sm focus:border-blue-500 focus:ring-blue-500"
          >
            <option value="all">All Status</option>
            <option value="on_track">On Track</option>
            <option value="over_budget">Over Budget</option>
            <option value="completed">Completed</option>
          </select>
          <label className="flex items-center text-sm text-gray-600 whitespace-nowrap">
            <input
              type="checkbox"
              checked={showInactive}
              onChange={(e) => setShowInactive(e.target.checked)}
              className="rounded border-gray-300 text-blue-600 focus:ring-blue-500 mr-2"
            />
            Show inactive
          </label>
        </div>
      </div>

      {isLoading ? (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {[1, 2, 3].map((i) => (
            <div key={i} className="bg-white rounded-lg shadow p-6 animate-pulse">
              <div className="h-6 bg-gray-200 rounded w-3/4 mb-2"></div>
              <div className="h-4 bg-gray-200 rounded w-1/2 mb-4"></div>
              <div className="grid grid-cols-3 gap-4 mb-4">
                <div className="h-12 bg-gray-200 rounded"></div>
                <div className="h-12 bg-gray-200 rounded"></div>
                <div className="h-12 bg-gray-200 rounded"></div>
              </div>
              <div className="h-2 bg-gray-200 rounded"></div>
            </div>
          ))}
        </div>
      ) : projects.length === 0 ? (
        <div className="text-center py-12">
          <svg
            className="mx-auto h-12 w-12 text-gray-400"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"
            />
          </svg>
          <h3 className="mt-2 text-lg font-medium text-gray-900">No projects yet</h3>
          <p className="mt-1 text-gray-500">Get started by creating your first project.</p>
          <div className="mt-6">
            <Button onClick={() => setIsCreateModalOpen(true)}>
              Create Project
            </Button>
          </div>
        </div>
      ) : (
        <>
          {activeProjects.length > 0 && (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mb-6">
              {activeProjects.map((project) => (
                <ProjectCard
                  key={project.id}
                  project={project}
                  onClick={() => onProjectSelect(project)}
                />
              ))}
            </div>
          )}

          {showInactive && inactiveProjects.length > 0 && (
            <>
              <h3 className="text-lg font-medium text-gray-700 mb-4">Inactive Projects</h3>
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 opacity-60">
                {inactiveProjects.map((project) => (
                  <ProjectCard
                    key={project.id}
                    project={project}
                    onClick={() => onProjectSelect(project)}
                  />
                ))}
              </div>
            </>
          )}
        </>
      )}

      <ProjectCreateModal
        isOpen={isCreateModalOpen}
        onClose={() => setIsCreateModalOpen(false)}
        onSubmit={async (input) => {
          await createProject(input)
        }}
      />
    </div>
  )
}
