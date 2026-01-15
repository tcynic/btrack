import { useState, useEffect } from 'react'
import { Modal, Button, Input } from '../ui'
import type { Task, CreateTaskInput, UpdateTaskInput, ProjectWithStats } from '../../types'

interface TaskModalProps {
  isOpen: boolean
  onClose: () => void
  onSubmit: (input: CreateTaskInput | UpdateTaskInput) => Promise<void>
  projectId: number
  projects?: ProjectWithStats[]
  sourceType?: 'meeting' | 'note' | 'standalone'
  sourceId?: number
  task?: Task | null
}

export function TaskModal({
  isOpen,
  onClose,
  onSubmit,
  projectId,
  projects = [],
  sourceType = 'standalone',
  sourceId,
  task,
}: TaskModalProps) {
  const [title, setTitle] = useState('')
  const [description, setDescription] = useState('')
  const [status, setStatus] = useState<string>('pending')
  const [priority, setPriority] = useState<string>('medium')
  const [dueDate, setDueDate] = useState('')
  const [selectedProjectId, setSelectedProjectId] = useState(projectId)
  const [isSubmitting, setIsSubmitting] = useState(false)

  const isEditing = !!task

  useEffect(() => {
    if (task) {
      setTitle(task.title)
      setDescription(task.description)
      setStatus(task.status)
      setPriority(task.priority)
      setDueDate(task.dueDate)
      setSelectedProjectId(task.projectId)
    } else {
      setTitle('')
      setDescription('')
      setStatus('pending')
      setPriority('medium')
      setDueDate('')
      setSelectedProjectId(projectId)
    }
  }, [task, isOpen, projectId])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setIsSubmitting(true)

    try {
      if (isEditing && task) {
        await onSubmit({
          id: task.id,
          title,
          description,
          status,
          priority,
          dueDate,
        })
      } else {
        await onSubmit({
          projectId: selectedProjectId,
          sourceType,
          sourceId: sourceId,
          title,
          description,
          priority,
          dueDate,
        })
      }
      onClose()
    } catch (err) {
      console.error('Failed to save task:', err)
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title={isEditing ? 'Edit Task' : 'New Task'}
      size="lg"
    >
      <form onSubmit={handleSubmit} className="space-y-4">
        {!isEditing && projects.length > 0 && (
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Project
            </label>
            <select
              value={selectedProjectId}
              onChange={(e) => setSelectedProjectId(Number(e.target.value))}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              required
            >
              {projects.map((project) => (
                <option key={project.id} value={project.id}>
                  {project.name}
                </option>
              ))}
            </select>
          </div>
        )}

        <Input
          label="Title"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          required
          placeholder="Task title"
        />

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            Description
          </label>
          <textarea
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            rows={3}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            placeholder="Task description (optional)"
          />
        </div>

        <div className="grid grid-cols-2 gap-4">
          {isEditing && (
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Status
              </label>
              <select
                value={status}
                onChange={(e) => setStatus(e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              >
                <option value="pending">Pending</option>
                <option value="in_progress">In Progress</option>
                <option value="completed">Completed</option>
                <option value="cancelled">Cancelled</option>
              </select>
            </div>
          )}

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              Priority
            </label>
            <select
              value={priority}
              onChange={(e) => setPriority(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            >
              <option value="low">Low</option>
              <option value="medium">Medium</option>
              <option value="high">High</option>
            </select>
          </div>

          <div>
            <Input
              label="Due Date"
              type="date"
              value={dueDate}
              onChange={(e) => setDueDate(e.target.value)}
            />
          </div>
        </div>

        <div className="flex justify-end gap-3 pt-4">
          <Button type="button" variant="ghost" onClick={onClose}>
            Cancel
          </Button>
          <Button type="submit" disabled={isSubmitting}>
            {isSubmitting ? 'Saving...' : isEditing ? 'Update' : 'Create Task'}
          </Button>
        </div>
      </form>
    </Modal>
  )
}
