import { useState, useEffect } from 'react'
import { Modal, Button, Input } from '../ui'
import type { Meeting, CreateMeetingInput, UpdateMeetingInput } from '../../types'

interface MeetingModalProps {
  isOpen: boolean
  onClose: () => void
  onSubmit: (input: CreateMeetingInput | UpdateMeetingInput) => Promise<void>
  projectId: number
  meeting?: Meeting | null
}

export function MeetingModal({
  isOpen,
  onClose,
  onSubmit,
  projectId,
  meeting,
}: MeetingModalProps) {
  const [title, setTitle] = useState('')
  const [meetingDate, setMeetingDate] = useState('')
  const [durationMinutes, setDurationMinutes] = useState(60)
  const [attendees, setAttendees] = useState('')
  const [notes, setNotes] = useState('')
  const [isSubmitting, setIsSubmitting] = useState(false)

  const isEditing = !!meeting

  useEffect(() => {
    if (meeting) {
      setTitle(meeting.title)
      setMeetingDate(meeting.meetingDate)
      setDurationMinutes(meeting.durationMinutes)
      setAttendees(meeting.attendees)
      setNotes(meeting.notes)
    } else {
      setTitle('')
      setMeetingDate(new Date().toISOString().split('T')[0])
      setDurationMinutes(60)
      setAttendees('')
      setNotes('')
    }
  }, [meeting, isOpen])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setIsSubmitting(true)

    try {
      if (isEditing && meeting) {
        await onSubmit({
          id: meeting.id,
          title,
          meetingDate,
          durationMinutes,
          attendees,
          notes,
        })
      } else {
        await onSubmit({
          projectId,
          title,
          meetingDate,
          durationMinutes,
          attendees,
          notes,
        })
      }
      onClose()
    } catch (err) {
      console.error('Failed to save meeting:', err)
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title={isEditing ? 'Edit Meeting' : 'Add Meeting'}
      size="lg"
    >
      <form onSubmit={handleSubmit} className="space-y-4">
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
          <Button type="submit" disabled={isSubmitting}>
            {isSubmitting ? 'Saving...' : isEditing ? 'Update' : 'Add Meeting'}
          </Button>
        </div>
      </form>
    </Modal>
  )
}
