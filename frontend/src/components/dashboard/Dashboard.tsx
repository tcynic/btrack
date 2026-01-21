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
import { downloadCSV, generateExportFilename } from "../../utils/export";
import {
  CreateMeeting,
  UpdateMeeting,
  DeleteMeeting,
  ExportAllProjects,
} from "../../../wailsjs/go/main/App";
import type {
  MeetingWithProject,
  Meeting,
  CreateMeetingInput,
  UpdateMeetingInput,
} from "../../types";

export function Dashboard() {
  const [weeksBack, setWeeksBack] = useState(2)
  const [weeksForward, setWeeksForward] = useState(3)
  const [isExporting, setIsExporting] = useState(false)
  
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

  const handleExportAllProjects = async () => {
    setIsExporting(true)
    try {
      const csvContent = await ExportAllProjects()
      const filename = generateExportFilename('all-projects')
      downloadCSV(csvContent, filename)
    } catch (error) {
      console.error('Export failed:', error)
      alert('Failed to export projects. Please try again.')
    } finally {
      setIsExporting(false)
    }
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

  const handleUpdateMeetingNotes = async (updatedMeeting: Meeting) => {
    // Update the selected meeting with new data
    setSelectedMeeting({ ...selectedMeeting!, ...updatedMeeting });
    // Refresh today's meetings list
    await loadTodayMeetings();
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
            <h2 className="text-2xl font-bold text-ld-text">Capacity Overview</h2>
            <p className="text-ld-muted">
              Monitor your weekly bandwidth across all active projects
            </p>
          </div>
          <div className="flex items-center gap-3">
            <div className="flex items-center gap-2 text-sm">
              <label className="text-ld-text">Weeks back:</label>
              <input
                type="number"
                min="1"
                max="52"
                value={weeksBack}
                onChange={(e) => setWeeksBack(Number(e.target.value))}
                className="w-16 rounded border-ld-border text-sm bg-ld-surface text-ld-text px-2 py-1"
              />
            </div>
            <div className="flex items-center gap-2 text-sm">
              <label className="text-ld-text">Weeks forward:</label>
              <input
                type="number"
                min="1"
                max="52"
                value={weeksForward}
                onChange={(e) => setWeeksForward(Number(e.target.value))}
                className="w-16 rounded border-ld-border text-sm bg-ld-surface text-ld-text px-2 py-1"
              />
            </div>
            <Button size="sm" onClick={handleDateRangeChange}>
              Update
            </Button>
            <Button
              size="sm"
              variant="ghost"
              onClick={handleExportAllProjects}
              disabled={isExporting}
            >
              {isExporting ? (
                <span className="flex items-center gap-1">
                  <svg className="animate-spin h-4 w-4" fill="none" viewBox="0 0 24 24">
                    <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                    <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                  Exporting...
                </span>
              ) : (
                <span className="flex items-center gap-1">
                  <svg className="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                  </svg>
                  Export All
                </span>
              )}
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
              <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-ld-primary"></div>
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
        onUpdate={handleUpdateMeetingNotes}
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
