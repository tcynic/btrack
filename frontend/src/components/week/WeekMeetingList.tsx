import type { MeetingWithProject } from '../../types'

interface WeekMeetingListProps {
  meetings: MeetingWithProject[]
  weekStartDate: string
}

// Helper to group meetings by date
function groupMeetingsByDate(meetings: MeetingWithProject[]): Map<string, MeetingWithProject[]> {
  const grouped = new Map<string, MeetingWithProject[]>()
  
  meetings.forEach(meeting => {
    const dateKey = meeting.meetingDate
    if (!grouped.has(dateKey)) {
      grouped.set(dateKey, [])
    }
    grouped.get(dateKey)!.push(meeting)
  })
  
  return grouped
}

// Helper to format date as "Mon Jan 13"
function formatDateShort(dateStr: string): string {
  const date = new Date(dateStr + 'T00:00:00')
  return date.toLocaleDateString('en-US', { weekday: 'short', month: 'short', day: 'numeric' })
}

export function WeekMeetingList({ meetings, weekStartDate }: WeekMeetingListProps) {
  // Generate all 7 days of the week
  const weekDays: string[] = []
  const startDate = new Date(weekStartDate + 'T00:00:00')
  for (let i = 0; i < 7; i++) {
    const date = new Date(startDate)
    date.setDate(startDate.getDate() + i)
    weekDays.push(date.toISOString().split('T')[0])
  }

  const groupedMeetings = groupMeetingsByDate(meetings)

  if (meetings.length === 0) {
    return (
      <div className="text-center py-12 text-ld-muted">
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
        <p>No meetings scheduled this week</p>
      </div>
    )
  }

  return (
    <div className="space-y-4">
      {weekDays.map((dateStr) => {
        const dayMeetings = groupedMeetings.get(dateStr) || []
        if (dayMeetings.length === 0) return null

        return (
          <div key={dateStr}>
            <div className="flex items-center gap-2 mb-2">
              <h4 className="text-sm font-semibold text-ld-text">
                {formatDateShort(dateStr)}
              </h4>
              <div className="flex-1 h-px bg-ld-border" />
            </div>
            <div className="space-y-2">
              {dayMeetings.map((meeting) => (
                <div
                  key={meeting.id}
                  className="flex items-center gap-3 p-3 rounded-lg border border-ld-border bg-ld-surface hover:border-ld-primary hover:bg-ld-surface2 transition-colors"
                >
                  <div className="flex items-center justify-center w-16 h-10 bg-ld-surface2 text-ld-primary rounded text-sm font-medium flex-shrink-0">
                    {meeting.durationMinutes} min
                  </div>
                  <div className="flex-1 min-w-0">
                    <div className="font-medium text-ld-text truncate">
                      {meeting.title}
                    </div>
                    <div className="text-sm text-ld-muted truncate">
                      {meeting.projectName}
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )
      })}
    </div>
  )
}
