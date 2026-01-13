import type { MeetingWithProject } from './meeting'
import type { WeeklyEntryWithProject } from './weekly-entry'

export interface WeekViewData {
  weekStartDate: string
  meetings: MeetingWithProject[]
  entries: WeeklyEntryWithProject[]
}
