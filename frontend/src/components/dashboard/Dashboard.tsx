import { useEffect, useState } from "react";
import { Card, CardHeader, Button } from "../ui";
import { SummaryCards } from "./SummaryCards";
import { WeeklyChart } from "./WeeklyChart";
import { WeeklyTable } from "./WeeklyTable";
import { DailyAgenda } from "./DailyAgenda";
import { DashboardMeetingModal } from "./DashboardMeetingModal";
import { MeetingDetailModal } from "../meetings/MeetingDetailModal";
import { MeetingModal } from "../meetings/MeetingModal";
import { useDashboard } from "../../hooks";
import {
  CreateMeeting,
  UpdateMeeting,
  DeleteMeeting,
} from "../../../wailsjs/go/main/App";
import type {
  MeetingWithProject,
  Meeting,
  CreateMeetingInput,
  UpdateMeetingInput,
} from "../../types";

export function Dashboard() {
  const [weeksBack, setWeeksBack] = useState(4)
  const [weeksForward, setWeeksForward] = useState(12)
  
  const {
    summary,
    weekData,
    todayMeetings,
    isLoading,
    refreshDashboard,
    loadTodayMeetings,
    loadDashboardData,
  } = useDashboard();

  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
  const [isDetailModalOpen, setIsDetailModalOpen] = useState(false);
  const [isEditModalOpen, setIsEditModalOpen] = useState(false);
  const [selectedMeeting, setSelectedMeeting] =
    useState<MeetingWithProject | null>(null);

  useEffect(() => {
    refreshDashboard();
  }, [refreshDashboard]);
  
  const handleDateRangeChange = async () => {
    await loadDashboardData(weeksBack, weeksForward)
  }

  const handleNewMeeting = () => {
    setIsCreateModalOpen(true);
  };

  const handleMeetingClick = (meeting: MeetingWithProject) => {
    setSelectedMeeting(meeting);
    setIsDetailModalOpen(true);
  };

  const handleCreateMeeting = async (input: CreateMeetingInput) => {
    await CreateMeeting(input);
    await loadTodayMeetings();
  };

  const handleEditMeeting = (_meeting: Meeting) => {
    setIsDetailModalOpen(false);
    setIsEditModalOpen(true);
  };

  const handleUpdateMeeting = async (
    input: CreateMeetingInput | UpdateMeetingInput,
  ) => {
    if ("id" in input) {
      await UpdateMeeting(input as UpdateMeetingInput);
      await loadTodayMeetings();
    }
  };

  const handleDeleteMeeting = async (id: number) => {
    await DeleteMeeting(id);
    setSelectedMeeting(null);
    await loadTodayMeetings();
  };

  const handleCloseEditModal = () => {
    setIsEditModalOpen(false);
    setSelectedMeeting(null);
  };

  const handleCloseDetailModal = () => {
    setIsDetailModalOpen(false);
    setSelectedMeeting(null);
  };

  return (
    <div>
      <DailyAgenda
        meetings={todayMeetings}
        isLoading={isLoading}
        onNewMeeting={handleNewMeeting}
        onMeetingClick={handleMeetingClick}
      />

      <div className="mb-6">
        <div className="flex justify-between items-start">
          <div>
            <h2 className="text-2xl font-bold text-gray-900">Capacity Overview</h2>
            <p className="text-gray-500">
              Monitor your weekly bandwidth across all active projects
            </p>
          </div>
          <div className="flex items-center gap-3">
            <div className="flex items-center gap-2 text-sm">
              <label className="text-gray-600">Weeks back:</label>
              <input
                type="number"
                min="1"
                max="52"
                value={weeksBack}
                onChange={(e) => setWeeksBack(Number(e.target.value))}
                className="w-16 rounded border-gray-300 text-sm"
              />
            </div>
            <div className="flex items-center gap-2 text-sm">
              <label className="text-gray-600">Weeks forward:</label>
              <input
                type="number"
                min="1"
                max="52"
                value={weeksForward}
                onChange={(e) => setWeeksForward(Number(e.target.value))}
                className="w-16 rounded border-gray-300 text-sm"
              />
            </div>
            <Button size="sm" onClick={handleDateRangeChange}>
              Update
            </Button>
          </div>
        </div>
      </div>

      <SummaryCards summary={summary} />

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <Card>
          <CardHeader
            title="Weekly Capacity"
            subtitle="Planned vs actual hours per week"
          />
          {isLoading ? (
            <div className="h-80 flex items-center justify-center">
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
            </div>
          ) : (
            <WeeklyChart data={weekData} />
          )}
        </Card>

        <Card padding="none">
          <div className="px-6 pt-6">
            <CardHeader
              title="Week by Week"
              subtitle="Detailed breakdown of hours"
            />
          </div>
          <div className="max-h-96 overflow-y-auto">
            {isLoading ? (
              <div className="h-40 flex items-center justify-center">
                <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
              </div>
            ) : (
              <WeeklyTable data={weekData} />
            )}
          </div>
        </Card>
      </div>

      <DashboardMeetingModal
        isOpen={isCreateModalOpen}
        onClose={() => setIsCreateModalOpen(false)}
        onSubmit={handleCreateMeeting}
      />

      <MeetingDetailModal
        isOpen={isDetailModalOpen}
        onClose={handleCloseDetailModal}
        onEdit={handleEditMeeting}
        onDelete={handleDeleteMeeting}
        meeting={selectedMeeting}
      />

      {selectedMeeting && (
        <MeetingModal
          isOpen={isEditModalOpen}
          onClose={handleCloseEditModal}
          onSubmit={handleUpdateMeeting}
          projectId={selectedMeeting.projectId}
          meeting={selectedMeeting}
        />
      )}
    </div>
  );
}
