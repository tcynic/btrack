import { useState, useEffect } from 'react'
import { Modal, Button, Input } from '../ui'
import type { Goal, GoalStatus, CreateGoalInput, UpdateGoalInput } from '../../types'

interface GoalModalProps {
  isOpen: boolean
  onClose: () => void
  onSubmit: (input: CreateGoalInput | UpdateGoalInput) => Promise<void>
  projectId: number
  goal?: Goal | null
}

export function GoalModal({
  isOpen,
  onClose,
  onSubmit,
  projectId,
  goal,
}: GoalModalProps) {
  const [title, setTitle] = useState('')
  const [description, setDescription] = useState('')
  const [status, setStatus] = useState<GoalStatus>('pending')
  const [targetDate, setTargetDate] = useState('')
  const [isSubmitting, setIsSubmitting] = useState(false)

  const isEditing = !!goal

  useEffect(() => {
    if (goal) {
      setTitle(goal.title)
      setDescription(goal.description)
      setStatus(goal.status)
      setTargetDate(goal.targetDate)
    } else {
      setTitle('')
      setDescription('')
      setStatus('pending')
      setTargetDate('')
    }
  }, [goal, isOpen])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setIsSubmitting(true)

    try {
      if (isEditing && goal) {
        await onSubmit({
          id: goal.id,
          title,
          description,
          status,
          targetDate,
        })
      } else {
        await onSubmit({
          projectId,
          title,
          description,
          targetDate,
        })
      }
      onClose()
    } catch (err) {
      console.error('Failed to save goal:', err)
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title={isEditing ? 'Edit Goal' : 'Add Goal'}
      size="lg"
    >
      <form onSubmit={handleSubmit} className="space-y-4">
        <Input
          label="Title"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          required
          placeholder="Goal title"
        />

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            Description (Markdown supported)
          </label>
          <textarea
            value={description}
            onChange={(e) => setDescription(e.target.value)}
            rows={4}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent font-mono text-sm"
            placeholder="Describe the goal and success criteria..."
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
                onChange={(e) => setStatus(e.target.value as GoalStatus)}
                className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="pending">Pending</option>
                <option value="in_progress">In Progress</option>
                <option value="completed">Completed</option>
              </select>
            </div>
          )}

          <div className={isEditing ? '' : 'col-span-2'}>
            <Input
              label="Target Date (optional)"
              type="date"
              value={targetDate}
              onChange={(e) => setTargetDate(e.target.value)}
            />
          </div>
        </div>

        <div className="flex justify-end gap-3 pt-4">
          <Button type="button" variant="ghost" onClick={onClose}>
            Cancel
          </Button>
          <Button type="submit" disabled={isSubmitting}>
            {isSubmitting ? 'Saving...' : isEditing ? 'Update' : 'Add Goal'}
          </Button>
        </div>
      </form>
    </Modal>
  )
}
