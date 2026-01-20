import { useState, useEffect, useCallback } from 'react'
import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import rehypeHighlight from 'rehype-highlight'
import { Modal, Button } from '../ui'
import { TaskList, TaskModal } from '../tasks'
import { formatDate } from '../../utils'
import { GetTasksBySource, CreateTask, UpdateTask, DeleteTask, UpdateTaskStatus, UpdateMeeting } from '../../../wailsjs/go/main/App'
import type { Meeting, Task, CreateTaskInput, UpdateTaskInput } from '../../types'

interface MeetingDetailModalProps {
  isOpen: boolean
  onClose: () => void
  onEdit: (meeting: Meeting) => void
  onDelete: (id: number) => void
  onUpdate: (meeting: Meeting) => void
  meeting: Meeting | null
}

export function MeetingDetailModal({
  isOpen,
  onClose,
  onEdit,
  onDelete,
  onUpdate,
  meeting,
}: MeetingDetailModalProps) {
  const [tasks, setTasks] = useState<Task[]>([])
  const [selectedTask, setSelectedTask] = useState<Task | null>(null)
  const [isTaskModalOpen, setIsTaskModalOpen] = useState(false)
  const [isEditing, setIsEditing] = useState(false)
  const [editedNotes, setEditedNotes] = useState('')
  const [isSaving, setIsSaving] = useState(false)

  const loadTasks = useCallback(async () => {
    if (meeting) {
      try {
        const meetingTasks = await GetTasksBySource('meeting', meeting.id)
        setTasks(meetingTasks)
      } catch (err) {
        console.error('Failed to load tasks:', err)
      }
    }
  }, [meeting])

  useEffect(() => {
    if (isOpen && meeting) {
      loadTasks()
      setIsEditing(false)
      setEditedNotes(meeting.notes || '')
    }
  }, [isOpen, meeting, loadTasks])

  // Early return after hooks to comply with React rules
  if (!meeting) return null

  const handleEdit = () => {
    onEdit(meeting)
    onClose()
  }

  const handleDelete = () => {
    if (window.confirm('Are you sure you want to delete this meeting?')) {
      onDelete(meeting.id)
      onClose()
    }
  }

  const handleAddTask = () => {
    setSelectedTask(null)
    setIsTaskModalOpen(true)
  }

  const handleEditTask = (task: Task) => {
    setSelectedTask(task)
    setIsTaskModalOpen(true)
  }

  const handleSubmitTask = async (input: CreateTaskInput | UpdateTaskInput) => {
    if ('id' in input) {
      await UpdateTask(input)
    } else {
      await CreateTask(input)
    }
    await loadTasks()
    setIsTaskModalOpen(false)
  }

  const handleDeleteTask = async (id: number) => {
    if (confirm('Are you sure you want to delete this task?')) {
      await DeleteTask(id)
      await loadTasks()
    }
  }

  const handleStatusChange = async (id: number, status: string) => {
    await UpdateTaskStatus(id, status)
    await loadTasks()
  }

  const handleEditNotes = () => {
    setEditedNotes(meeting.notes || '')
    setIsEditing(true)
  }

  const handleCancelEdit = () => {
    setEditedNotes(meeting.notes || '')
    setIsEditing(false)
  }

  const handleSaveNotes = async () => {
    setIsSaving(true)
    try {
      const updatedMeeting = await UpdateMeeting({
        id: meeting.id,
        title: meeting.title,
        meetingDate: meeting.meetingDate,
        durationMinutes: meeting.durationMinutes,
        attendees: meeting.attendees,
        notes: editedNotes,
      })
      onUpdate(updatedMeeting)
      setIsEditing(false)
    } catch (err) {
      console.error('Failed to update meeting notes:', err)
      alert('Failed to save notes. Please try again.')
    } finally {
      setIsSaving(false)
    }
  }

  return (
    <Modal isOpen={isOpen} onClose={onClose} title={meeting.title} size="xl">
      {/* Metadata Section */}
      <div className="flex items-center gap-4 text-sm text-gray-600 mb-6 pb-4 border-b border-gray-200">
        <div className="flex items-center gap-2">
          <svg className="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
          </svg>
          <span>{formatDate(meeting.meetingDate)}</span>
        </div>
        <div className="flex items-center gap-2">
          <svg className="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <span>{meeting.durationMinutes} min</span>
        </div>
        {meeting.attendees && (
          <div className="flex items-center gap-2">
            <svg className="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
            </svg>
            <span>{meeting.attendees}</span>
          </div>
        )}
      </div>

      {/* Notes Content */}
      <div className="mb-6">
        <div className="flex items-center justify-between mb-2">
          <h3 className="text-sm font-medium text-gray-700">Notes</h3>
          {!isEditing && (
            <Button size="sm" variant="ghost" onClick={handleEditNotes}>
              Edit Notes
            </Button>
          )}
        </div>
        {isEditing ? (
          <div className="space-y-3">
            <textarea
              value={editedNotes}
              onChange={(e) => setEditedNotes(e.target.value)}
              rows={12}
              className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent font-mono text-sm"
              placeholder="# Meeting Notes&#10;&#10;- Discussion point 1&#10;- Discussion point 2&#10;&#10;## Action Items&#10;- [ ] Task 1"
            />
            <div className="flex justify-end gap-2">
              <Button size="sm" variant="ghost" onClick={handleCancelEdit} disabled={isSaving}>
                Cancel
              </Button>
              <Button size="sm" onClick={handleSaveNotes} disabled={isSaving}>
                {isSaving ? 'Saving...' : 'Save Notes'}
              </Button>
            </div>
          </div>
        ) : (
          <div className="prose prose-sm max-w-none">
            {meeting.notes ? (
              <ReactMarkdown
                remarkPlugins={[remarkGfm]}
                rehypePlugins={[rehypeHighlight]}
              >
                {meeting.notes}
              </ReactMarkdown>
            ) : (
              <p className="text-gray-500 italic">No notes for this meeting</p>
            )}
          </div>
        )}
      </div>

      {/* Tasks Section */}
      <div className="border-t border-gray-200 pt-6">
        <div className="flex items-center justify-between mb-4">
          <h3 className="text-lg font-medium text-gray-900">Tasks</h3>
          <Button size="sm" onClick={handleAddTask}>Add Task</Button>
        </div>
        <TaskList
          tasks={tasks}
          onEdit={handleEditTask}
          onDelete={handleDeleteTask}
          onStatusChange={handleStatusChange}
          emptyMessage="No tasks for this meeting"
        />
      </div>

      {/* Action Buttons */}
      <div className="flex justify-end gap-3 pt-6 mt-6 border-t border-gray-200">
        <Button variant="ghost" onClick={onClose}>
          Close
        </Button>
        <Button variant="ghost" onClick={handleDelete}>
          Delete
        </Button>
        <Button variant="ghost" onClick={handleEdit}>Edit Details</Button>
      </div>

      {/* Task Modal */}
      {isTaskModalOpen && (
        <TaskModal
          isOpen={isTaskModalOpen}
          onClose={() => setIsTaskModalOpen(false)}
          onSubmit={handleSubmitTask}
          projectId={meeting.projectId}
          sourceType="meeting"
          sourceId={meeting.id}
          task={selectedTask}
        />
      )}
    </Modal>
  )
}
