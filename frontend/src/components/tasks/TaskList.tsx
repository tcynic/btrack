import { TaskItem } from './TaskItem'
import type { Task, TaskWithContext } from '../../types'

interface TaskListProps {
  tasks: Task[] | TaskWithContext[]
  onEdit: (task: Task | TaskWithContext) => void
  onDelete: (id: number) => void
  onStatusChange: (id: number, status: string) => void
  showProject?: boolean
  emptyMessage?: string
}

export function TaskList({
  tasks,
  onEdit,
  onDelete,
  onStatusChange,
  showProject = false,
  emptyMessage = 'No tasks yet',
}: TaskListProps) {
  if (tasks.length === 0) {
    return (
      <div className="text-center py-8 text-gray-500">
        {emptyMessage}
      </div>
    )
  }

  return (
    <div className="space-y-2">
      {tasks.map((task) => (
        <TaskItem
          key={task.id}
          task={task}
          onEdit={onEdit}
          onDelete={onDelete}
          onStatusChange={onStatusChange}
          showProject={showProject}
        />
      ))}
    </div>
  )
}
