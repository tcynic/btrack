import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import rehypeHighlight from 'rehype-highlight'
import { Modal, Button } from '../ui'
import { formatDate } from '../../utils'
import type { Meeting } from '../../types'

interface MeetingDetailModalProps {
  isOpen: boolean
  onClose: () => void
  onEdit: (meeting: Meeting) => void
  onDelete: (id: number) => void
  meeting: Meeting | null
}

export function MeetingDetailModal({
  isOpen,
  onClose,
  onEdit,
  onDelete,
  meeting,
}: MeetingDetailModalProps) {
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

      {/* Action Buttons */}
      <div className="flex justify-end gap-3 pt-6 mt-6 border-t border-gray-200">
        <Button variant="ghost" onClick={onClose}>
          Close
        </Button>
        <Button variant="ghost" onClick={handleDelete}>
          Delete
        </Button>
        <Button onClick={handleEdit}>Edit Meeting</Button>
      </div>
    </Modal>
  )
}
