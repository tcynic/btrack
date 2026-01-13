import { useState, useEffect } from 'react'
import { Modal, Button, Input } from '../ui'
import { GetAllProjects } from '../../../wailsjs/go/main/App'
import type { ProjectWithStats, CreateMeetingInput } from '../../types'

interface DashboardMeetingModalProps {
  isOpen: boolean
  onClose: () => void
  onSubmit: (input: CreateMeetingInput) => Promise<void>
}

export function DashboardMeetingModal({
  isOpen,
  onClose,
  onSubmit,
}: DashboardMeetingModalProps) {
  const [projects, setProjects] = useState<ProjectWithStats[]>([])
  const [projectId, setProjectId] = useState<number>(0)
  const [title, setTitle] = useState('')
  const [meetingDate, setMeetingDate] = useState('')
  const [durationMinutes, setDurationMinutes] = useState(60)
  const [attendees, setAttendees] = useState('')
  const [notes, setNotes] = useState('')
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [isLoadingProjects, setIsLoadingProjects] = useState(false)
  const [error, setError] = useState('')

  useEffect(() => {
    if (isOpen) {
      loadProjects()
      resetForm()
    }
  }, [isOpen])

  const loadProjects = async () => {
    setIsLoadingProjects(true)
    try {
      const data = await GetAllProjects(true) // Only active projects
      setProjects(data || [])
      if (data && data.length > 0) {
        setProjectId(data[0].id)
      }
    } catch (err) {
      console.error('Failed to load projects:', err)
      setError('Failed to load projects')
    } finally {
      setIsLoadingProjects(false)
    }
  }

  const resetForm = () => {
    setTitle('')
    setMeetingDate(new Date().toISOString().split('T')[0])
    setDurationMinutes(60)
    setAttendees('')
    setNotes('')
    setError('')
    setProjectId(0)
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    if (!projectId) {
      setError('Please select a project')
      return
    }

    setIsSubmitting(true)
    setError('')

    try {
      await onSubmit({
        projectId,
        title,
        meetingDate,
        durationMinutes,
        attendees,
        notes,
      })
      onClose()
    } catch (err) {
      console.error('Failed to create meeting:', err)
      setError('Failed to create meeting')
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title="New Meeting"
      size="lg"
    >
      <form onSubmit={handleSubmit} className="space-y-4">
        {error && (
          <div className="p-3 bg-red-50 border border-red-200 rounded-lg text-red-700 text-sm">
            {error}
          </div>
        )}

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            Project
          </label>
          {isLoadingProjects ? (
            <div className="flex items-center gap-2 py-2 text-gray-500 text-sm">
              <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-blue-600"></div>
              Loading projects...
            </div>
          ) : projects.length === 0 ? (
            <div className="py-2 text-gray-500 text-sm">
              No active projects available
            </div>
          ) : (
            <select
              value={projectId}
              onChange={(e) => setProjectId(parseInt(e.target.value))}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
              required
            >
              <option value={0} disabled>Select a project</option>
              {projects.map((project) => (
                <option key={project.id} value={project.id}>
                  {project.name}
                </option>
              ))}
            </select>
          )}
        </div>

        <Input
          label="Title"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          required
          placeholder="Meeting title"
        />

        <div className="grid grid-cols-2 gap-4">
          <Input
            label="Date"
            type="date"
            value={meetingDate}
            onChange={(e) => setMeetingDate(e.target.value)}
            required
          />

          <Input
            label="Duration (minutes)"
            type="number"
            value={durationMinutes}
            onChange={(e) => setDurationMinutes(parseInt(e.target.value) || 60)}
            min={15}
            step={15}
          />
        </div>

        <Input
          label="Attendees"
          value={attendees}
          onChange={(e) => setAttendees(e.target.value)}
          placeholder="e.g., John, Jane, Bob"
        />

        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            Notes (Markdown supported)
          </label>
          <textarea
            value={notes}
            onChange={(e) => setNotes(e.target.value)}
            rows={6}
            className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent font-mono text-sm"
            placeholder="# Meeting Notes&#10;&#10;- Discussion point 1&#10;- Discussion point 2&#10;&#10;## Action Items&#10;- [ ] Task 1"
          />
        </div>

        <div className="flex justify-end gap-3 pt-4">
          <Button type="button" variant="ghost" onClick={onClose}>
            Cancel
          </Button>
          <Button type="submit" disabled={isSubmitting || isLoadingProjects || projects.length === 0}>
            {isSubmitting ? 'Creating...' : 'Create Meeting'}
          </Button>
        </div>
      </form>
    </Modal>
  )
}
