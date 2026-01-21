import { useState, useEffect } from 'react'
import { TaskList, TaskModal } from '.'
import { Button } from '../ui'
import { GetAllTasks, CreateTask, UpdateTask, DeleteTask, UpdateTaskStatus, GetAllProjects } from '../../../wailsjs/go/main/App'
import type { TaskWithContext, CreateTaskInput, UpdateTaskInput, TaskStatus, ProjectWithStats } from '../../types'

export function TasksView() {
  const [tasks, setTasks] = useState<TaskWithContext[]>([])
  const [projects, setProjects] = useState<ProjectWithStats[]>([])
  const [selectedTask, setSelectedTask] = useState<TaskWithContext | null>(null)
  const [isModalOpen, setIsModalOpen] = useState(false)
  const [isLoading, setIsLoading] = useState(true)
  const [statusFilter, setStatusFilter] = useState<string>('')
  const [projectFilter, setProjectFilter] = useState<number>(0)

  const loadTasks = async () => {
    try {
      setIsLoading(true)
      const allTasks = await GetAllTasks(statusFilter, projectFilter)
      setTasks(allTasks)
    } catch (err) {
      console.error('Failed to load tasks:', err)
    } finally {
      setIsLoading(false)
    }
  }

  const loadProjects = async () => {
    try {
      const allProjects = await GetAllProjects(false)
      setProjects(allProjects as any)
    } catch (err) {
      console.error('Failed to load projects:', err)
    }
  }

  useEffect(() => {
    loadProjects()
  }, [])

  useEffect(() => {
    loadTasks()
  }, [statusFilter, projectFilter])

  const handleCreateTask = () => {
    setSelectedTask(null)
    setIsModalOpen(true)
  }

  const handleEditTask = (task: TaskWithContext) => {
    setSelectedTask(task)
    setIsModalOpen(true)
  }

  const handleSubmitTask = async (input: CreateTaskInput | UpdateTaskInput) => {
    if ('id' in input) {
      await UpdateTask(input)
    } else {
      await CreateTask(input)
    }
    await loadTasks()
    setIsModalOpen(false)
  }

  const handleDeleteTask = async (id: number) => {
    if (confirm('Are you sure you want to delete this task?')) {
      await DeleteTask(id)
      await loadTasks()
    }
  }

  const handleStatusChange = async (id: number, status: string) => {
    await UpdateTaskStatus(id, status as TaskStatus)
    await loadTasks()
  }

  const activeProjects = projects.filter(p => p.isActive)

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-ld-text">Tasks</h1>
          <p className="text-sm text-ld-muted mt-1">
            Manage tasks across all projects
          </p>
        </div>
        <Button onClick={handleCreateTask}>New Task</Button>
      </div>

      {/* Filters */}
      <div className="bg-ld-surface rounded-lg border border-ld-border p-4">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div>
            <label className="block text-sm font-medium text-ld-text mb-1">
              Status
            </label>
            <select
              value={statusFilter}
              onChange={(e) => setStatusFilter(e.target.value)}
              className="w-full px-3 py-2 border border-ld-border bg-ld-surface text-ld-text rounded-lg focus:outline-none focus:ring-2 focus:ring-ld-primary"
            >
              <option value="">All Statuses</option>
              <option value="pending">Pending</option>
              <option value="in_progress">In Progress</option>
              <option value="completed">Completed</option>
              <option value="cancelled">Cancelled</option>
            </select>
          </div>

          <div>
            <label className="block text-sm font-medium text-ld-text mb-1">
              Project
            </label>
            <select
              value={projectFilter}
              onChange={(e) => setProjectFilter(Number(e.target.value))}
              className="w-full px-3 py-2 border border-ld-border bg-ld-surface text-ld-text rounded-lg focus:outline-none focus:ring-2 focus:ring-ld-primary"
            >
              <option value="0">All Projects</option>
              {activeProjects.map((project) => (
                <option key={project.id} value={project.id}>
                  {project.name}
                </option>
              ))}
            </select>
          </div>

          <div className="flex items-end">
            <button
              onClick={() => {
                setStatusFilter('')
                setProjectFilter(0)
              }}
              className="text-sm text-ld-primary hover:brightness-110 font-medium"
            >
              Clear Filters
            </button>
          </div>
        </div>
      </div>

      {/* Task List */}
      <div>
        {isLoading ? (
          <div className="text-center py-12 text-ld-muted">Loading tasks...</div>
        ) : (
          <TaskList
            tasks={tasks}
            onEdit={handleEditTask as any}
            onDelete={handleDeleteTask}
            onStatusChange={handleStatusChange}
            showProject={true}
            emptyMessage="No tasks found. Create a task to get started."
          />
        )}
      </div>

      {/* Task Modal */}
      {isModalOpen && (
        <TaskModal
          isOpen={isModalOpen}
          onClose={() => setIsModalOpen(false)}
          onSubmit={handleSubmitTask}
          projectId={projectFilter || (activeProjects.length > 0 ? activeProjects[0].id : 0)}
          projects={activeProjects}
          task={selectedTask}
        />
      )}
    </div>
  )
}
