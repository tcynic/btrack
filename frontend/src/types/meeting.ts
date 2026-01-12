export interface Meeting {
  id: number
  projectId: number
  title: string
  meetingDate: string
  durationMinutes: number
  attendees: string
  notes: string
  createdAt: string
  updatedAt: string
}

export interface CreateMeetingInput {
  projectId: number
  title: string
  meetingDate: string
  durationMinutes: number
  attendees: string
  notes: string
}

export interface UpdateMeetingInput {
  id: number
  title: string
  meetingDate: string
  durationMinutes: number
  attendees: string
  notes: string
}
