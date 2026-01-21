import { Card, CardHeader, Button } from '../ui'
import type { MeetingWithProject } from '../../types'

interface DailyAgendaProps {
  meetings: MeetingWithProject[]
  isLoading: boolean
  onNewMeeting: () => void
  onMeetingClick: (meeting: MeetingWithProject) => void
}

export function DailyAgenda({
  meetings,
  isLoading,
  onNewMeeting,
  onMeetingClick,
}: DailyAgendaProps) {
  const today = new Date()
  const formattedDate = today.toLocaleDateString('en-US', {
    weekday: 'long',
    month: 'long',
    day: 'numeric',
  })

  return (
    <Card className="mb-6">
      <CardHeader
        title="Today's Agenda"
        subtitle={formattedDate}
        action={
          <Button size="sm" onClick={onNewMeeting}>
            <svg
              className="w-4 h-4 mr-1.5"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M12 4v16m8-8H4"
              />
            </svg>
            New Meeting
          </Button>
        }
      />

      {isLoading ? (
        <div className="flex items-center justify-center py-8">
          <div className="animate-spin rounded-full h-6 w-6 border-b-2 border-ld-primary"></div>
        </div>
      ) : meetings.length === 0 ? (
        <div className="text-center py-8 text-ld-muted">
          <svg
            className="w-12 h-12 mx-auto mb-3 text-ld-border"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={1.5}
              d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
            />
          </svg>
          <p>No meetings scheduled for today</p>
        </div>
      ) : (
        <div className="space-y-2">
          {meetings.map((meeting) => (
            <div
              key={meeting.id}
              onClick={() => onMeetingClick(meeting)}
              className="flex items-center justify-between p-3 rounded-lg border border-ld-border hover:border-ld-primary hover:bg-ld-surface2 cursor-pointer transition-colors"
            >
              <div className="flex items-center gap-3">
                <div className="flex items-center justify-center w-16 h-10 bg-ld-surface2 text-ld-primary rounded text-sm font-medium">
                  {meeting.durationMinutes} min
                </div>
                <div>
                  <div className="font-medium text-ld-text">{meeting.title}</div>
                  <div className="text-sm text-ld-muted">{meeting.projectName}</div>
                </div>
              </div>
              <svg
                className="w-5 h-5 text-ld-muted"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M9 5l7 7-7 7"
                />
              </svg>
            </div>
          ))}
        </div>
      )}
    </Card>
  )
}
