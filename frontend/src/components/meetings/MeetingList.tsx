import { useEffect, useState } from 'react'
import ReactMarkdown from 'react-markdown'
import { Button, Card, CardHeader } from '../ui'
import { MeetingModal } from './MeetingModal'
import { useMeetings } from '../../hooks'
import { formatDate } from '../../utils'
import type { Meeting, CreateMeetingInput, UpdateMeetingInput } from '../../types'

interface MeetingListProps {
  projectId: number
}

export function MeetingList({ projectId }: MeetingListProps) {
  const {
    meetings,
    isLoading,
    loadMeetings,
    createMeeting,
    updateMeeting,
    deleteMeeting,
  } = useMeetings(projectId)

  const [isModalOpen, setIsModalOpen] = useState(false)
  const [editingMeeting, setEditingMeeting] = useState<Meeting | null>(null)
  const [expandedId, setExpandedId] = useState<number | null>(null)

  useEffect(() => {
    loadMeetings()
  }, [loadMeetings])

  const handleSubmit = async (input: CreateMeetingInput | UpdateMeetingInput) => {
    if ('id' in input) {
      await updateMeeting(input)
    } else {
      await createMeeting(input)
    }
  }

  const handleEdit = (meeting: Meeting) => {
    setEditingMeeting(meeting)
    setIsModalOpen(true)
  }

  const handleDelete = async (id: number) => {
    if (window.confirm('Are you sure you want to delete this meeting?')) {
      await deleteMeeting(id)
    }
  }

  const handleCloseModal = () => {
    setIsModalOpen(false)
    setEditingMeeting(null)
  }

  const toggleExpand = (id: number) => {
    setExpandedId(expandedId === id ? null : id)
  }

  return (
    <Card>
      <CardHeader
        title="Meetings"
        subtitle={`${meetings.length} meeting${meetings.length !== 1 ? 's' : ''}`}
        action={
          <Button size="sm" onClick={() => setIsModalOpen(true)}>
            <svg className="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
            </svg>
            Add Meeting
          </Button>
        }
      />

      {isLoading ? (
        <p className="text-gray-500">Loading meetings...</p>
      ) : meetings.length === 0 ? (
        <p className="text-gray-500 text-center py-8">No meetings yet. Add your first meeting!</p>
      ) : (
        <div className="space-y-3">
          {meetings.map((meeting) => (
            <div
              key={meeting.id}
              className="border border-gray-200 rounded-lg overflow-hidden"
            >
              <div
                className="flex items-center justify-between p-4 cursor-pointer hover:bg-gray-50"
                onClick={() => toggleExpand(meeting.id)}
              >
                <div className="flex items-center gap-4">
                  <div className="flex-shrink-0">
                    <svg className="h-5 w-5 text-blue-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                    </svg>
                  </div>
                  <div>
                    <h4 className="font-medium text-gray-900">{meeting.title}</h4>
                    <p className="text-sm text-gray-500">
                      {formatDate(meeting.meetingDate)} • {meeting.durationMinutes} min
                      {meeting.attendees && ` • ${meeting.attendees}`}
                    </p>
                  </div>
                </div>
                <svg
                  className={`h-5 w-5 text-gray-400 transition-transform ${expandedId === meeting.id ? 'rotate-180' : ''}`}
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                </svg>
              </div>

              {expandedId === meeting.id && (
                <div className="px-4 pb-4 border-t border-gray-200 bg-gray-50">
                  {meeting.notes ? (
                    <div className="prose prose-sm max-w-none mt-3">
                      <ReactMarkdown>{meeting.notes}</ReactMarkdown>
                    </div>
                  ) : (
                    <p className="text-gray-500 text-sm mt-3 italic">No notes for this meeting</p>
                  )}
                  <div className="flex gap-2 mt-4">
                    <Button size="sm" variant="ghost" onClick={() => handleEdit(meeting)}>
                      Edit
                    </Button>
                    <Button size="sm" variant="ghost" onClick={() => handleDelete(meeting.id)}>
                      Delete
                    </Button>
                  </div>
                </div>
              )}
            </div>
          ))}
        </div>
      )}

      <MeetingModal
        isOpen={isModalOpen}
        onClose={handleCloseModal}
        onSubmit={handleSubmit}
        projectId={projectId}
        meeting={editingMeeting}
      />
    </Card>
  )
}
