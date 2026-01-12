import { useEffect, useState } from "react";
import { Button, Card, CardHeader } from "../ui";
import { MeetingModal } from "./MeetingModal";
import { MeetingDetailModal } from "./MeetingDetailModal";
import { useMeetings } from "../../hooks";
import { formatDate } from "../../utils";
import type {
  Meeting,
  CreateMeetingInput,
  UpdateMeetingInput,
} from "../../types";

interface MeetingListProps {
  projectId: number;
}

export function MeetingList({ projectId }: MeetingListProps) {
  const {
    meetings,
    isLoading,
    loadMeetings,
    createMeeting,
    updateMeeting,
    deleteMeeting,
  } = useMeetings(projectId);

  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingMeeting, setEditingMeeting] = useState<Meeting | null>(null);
  const [viewingMeeting, setViewingMeeting] = useState<Meeting | null>(null);

  useEffect(() => {
    loadMeetings();
  }, [loadMeetings]);

  const handleSubmit = async (
    input: CreateMeetingInput | UpdateMeetingInput,
  ) => {
    if ("id" in input) {
      await updateMeeting(input);
    } else {
      await createMeeting(input);
    }
  };

  const handleDelete = async (id: number) => {
    if (window.confirm("Are you sure you want to delete this meeting?")) {
      await deleteMeeting(id);
    }
  };

  const handleCloseModal = () => {
    setIsModalOpen(false);
    setEditingMeeting(null);
  };

  const handleView = (meeting: Meeting) => {
    setViewingMeeting(meeting);
  };

  const handleCloseDetailModal = () => {
    setViewingMeeting(null);
  };

  const handleEditFromDetail = (meeting: Meeting) => {
    setEditingMeeting(meeting);
    setIsModalOpen(true);
  };

  return (
    <Card>
      <CardHeader
        title="Meetings"
        subtitle={`${meetings.length} meeting${meetings.length !== 1 ? "s" : ""}`}
        action={
          <Button size="sm" onClick={() => setIsModalOpen(true)}>
            <svg
              className="h-4 w-4 mr-1"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M12 4v16m8-8H4"
              />
            </svg>
            Add Meeting
          </Button>
        }
      />

      {isLoading ? (
        <p className="text-gray-500">Loading meetings...</p>
      ) : meetings.length === 0 ? (
        <p className="text-gray-500 text-center py-8">
          No meetings yet. Add your first meeting!
        </p>
      ) : (
        <div className="space-y-3">
          {meetings.map((meeting) => (
            <div
              key={meeting.id}
              className="border border-gray-200 rounded-lg p-4 cursor-pointer hover:bg-gray-50 transition-colors"
              onClick={() => handleView(meeting)}
            >
              <div className="flex items-center gap-4">
                <div className="flex-shrink-0">
                  <svg
                    className="h-5 w-5 text-blue-500"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
                    />
                  </svg>
                </div>
                <div className="flex-1">
                  <h4 className="font-medium text-gray-900">{meeting.title}</h4>
                  <p className="text-sm text-gray-500">
                    {formatDate(meeting.meetingDate)} •{" "}
                    {meeting.durationMinutes} min
                  </p>
                </div>
                <svg
                  className="h-5 w-5 text-gray-400"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M9 5l7 7-7 7"
                  />
                </svg>
              </div>
            </div>
          ))}
        </div>
      )}

      <MeetingDetailModal
        isOpen={!!viewingMeeting}
        onClose={handleCloseDetailModal}
        onEdit={handleEditFromDetail}
        onDelete={handleDelete}
        meeting={viewingMeeting}
      />

      <MeetingModal
        isOpen={isModalOpen}
        onClose={handleCloseModal}
        onSubmit={handleSubmit}
        projectId={projectId}
        meeting={editingMeeting}
      />
    </Card>
  );
}
