import { TaskStatusBadge } from './TaskStatusBadge'
import type { Task, TaskWithContext } from '../../types'

interface TaskItemProps {
  task: Task | TaskWithContext
  onEdit: (task: Task | TaskWithContext) => void
  onDelete: (id: number) => void
  onStatusChange: (id: number, status: string) => void
  showProject?: boolean
}

export function TaskItem({ task, onEdit, onDelete, onStatusChange, showProject = false }: TaskItemProps) {
  const getPriorityColor = (priority: string) => {
    switch (priority) {
      case 'high':
        return 'text-[var(--ld-pink)]'
      case 'medium':
        return 'text-[var(--ld-orange)]'
      default:
        return 'text-ld-muted'
    }
  }

  const nextStatus = (current: string) => {
    switch (current) {
      case 'pending':
        return 'in_progress'
      case 'in_progress':
        return 'completed'
      default:
        return 'pending'
    }
  }

  const handleStatusClick = () => {
    onStatusChange(task.id, nextStatus(task.status) as string)
  }

  return (
    <div className="flex items-center justify-between p-4 bg-ld-surface border border-ld-border rounded-lg hover:shadow-sm transition-shadow">
      <div className="flex-1 min-w-0">
        <div className="flex items-center gap-2 mb-1">
          <h4 className="text-sm font-medium text-ld-text truncate">{task.title}</h4>
          <span className={`text-xs font-medium ${getPriorityColor(task.priority)}`}>
            {task.priority.toUpperCase()}
          </span>
        </div>
        
        {task.description && (
          <p className="text-sm text-ld-muted mb-2">{task.description}</p>
        )}

        <div className="flex items-center gap-2 text-xs text-ld-muted">
          {showProject && 'projectName' in task && (
            <span className="px-2 py-1 bg-ld-surface2 text-ld-text rounded-full font-medium">{task.projectName}</span>
          )}
          {task.sourceType !== 'standalone' && 'sourceTitle' in task && task.sourceTitle && (
            <span>from: {task.sourceTitle}</span>
          )}
          {task.dueDate && (
            <span>Due: {new Date(task.dueDate).toLocaleDateString()}</span>
          )}
        </div>
      </div>

      <div className="flex items-center gap-2 ml-4">
        <TaskStatusBadge status={task.status} onClick={handleStatusClick} />
        <button
          onClick={() => onEdit(task)}
          className="text-ld-primary hover:brightness-110 text-sm font-medium"
        >
          Edit
        </button>
        <button
          onClick={() => onDelete(task.id)}
          className="text-[var(--ld-pink)] hover:brightness-110 text-sm font-medium"
        >
          Delete
        </button>
      </div>
    </div>
  )
}
