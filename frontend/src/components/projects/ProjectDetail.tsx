import { useState } from "react";
import { Button, Card, CardHeader } from "../ui";
import { WeeklyBreakdown } from "../weekly";
import { MeetingList } from "../meetings";
import { NoteList } from "../notes";
import { GoalList } from "../goals";
import { ProjectEditModal } from "./ProjectEditModal";
import { useProjects } from "../../hooks";
import { formatDate } from "../../utils";
import type { ProjectWithStats, UpdateProjectInput } from "../../types";

interface ProjectDetailProps {
  project: ProjectWithStats;
  onBack: () => void;
}

export function ProjectDetail({ project, onBack }: ProjectDetailProps) {
  const [isEditModalOpen, setIsEditModalOpen] = useState(false);
  const { updateProject, getProject } = useProjects();
  const [currentProject, setCurrentProject] = useState(project);

  const handleEditSubmit = async (input: UpdateProjectInput) => {
    await updateProject(input);
    // Refresh the project data
    const updated = await getProject(input.id);
    if (updated) {
      setCurrentProject(updated);
    }
  };
  const progressPercent =
    currentProject.totalPlannedHours > 0
      ? Math.min(
          100,
          (currentProject.totalActualHours / currentProject.totalPlannedHours) *
            100,
        )
      : 0;

  return (
    <div>
      {/* Header */}
      <div className="mb-6">
        <div className="flex items-center justify-between mb-4">
          <Button variant="ghost" onClick={onBack}>
            <svg
              className="h-5 w-5 mr-2"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M15 19l-7-7 7-7"
              />
            </svg>
            Back to Projects
          </Button>
          <Button onClick={() => setIsEditModalOpen(true)}>
            <svg
              className="h-5 w-5 mr-2"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
              />
            </svg>
            Edit Project
          </Button>
        </div>
        <div className="flex justify-between items-start">
          <div>
            <h2 className="text-3xl font-bold text-gray-900">
              {currentProject.name}
            </h2>
            <p className="text-gray-500 mt-1">
              {formatDate(currentProject.startDate)} -{" "}
              {formatDate(currentProject.endDate)}
            </p>
          </div>
          {!currentProject.isActive && (
            <span className="px-3 py-1 text-sm bg-gray-100 text-gray-600 rounded-full">
              Inactive
            </span>
          )}
        </div>
      </div>

      {/* Summary Cards */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
        <Card>
          <p className="text-sm text-gray-500 uppercase mb-1">
            Total Sold Hours
          </p>
          <p className="text-2xl font-bold text-gray-900">
            {currentProject.totalSoldHours}
          </p>
        </Card>
        <Card>
          <p className="text-sm text-gray-500 uppercase mb-1">My Hours</p>
          <p className="text-2xl font-bold text-blue-600">
            {currentProject.myHours}
          </p>
        </Card>
        <Card>
          <p className="text-sm text-gray-500 uppercase mb-1">
            Specialist Hours
          </p>
          <p className="text-2xl font-bold text-purple-600">
            {currentProject.specialistHours}
          </p>
          <p className="text-xs text-gray-500 mt-1">
            (
            {(
              (currentProject.specialistHours / currentProject.totalSoldHours) *
              100
            ).toFixed(0)}
            % of total)
          </p>
        </Card>
        <Card>
          <p className="text-sm text-gray-500 uppercase mb-1">Total Weeks</p>
          <p className="text-2xl font-bold text-gray-900">
            {currentProject.totalWeeks}
          </p>
        </Card>
      </div>

      {/* Progress Overview */}
      <Card className="mb-6">
        <CardHeader title="Progress Overview" />
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-4">
          <div>
            <p className="text-sm text-gray-500 mb-1">Planned Hours</p>
            <p className="text-3xl font-bold text-blue-600">
              {currentProject.totalPlannedHours}
            </p>
          </div>
          <div>
            <p className="text-sm text-gray-500 mb-1">Actual Hours</p>
            <p className="text-3xl font-bold text-gray-900">
              {currentProject.totalActualHours}
            </p>
          </div>
          <div>
            <p className="text-sm text-gray-500 mb-1">Remaining</p>
            <p
              className={`text-3xl font-bold ${
                currentProject.totalPlannedHours -
                  currentProject.totalActualHours >=
                0
                  ? "text-green-600"
                  : "text-red-600"
              }`}
            >
              {currentProject.totalPlannedHours -
                currentProject.totalActualHours}
            </p>
          </div>
        </div>

        {/* Progress Bar */}
        <div>
          <div className="flex justify-between text-sm text-gray-500 mb-2">
            <span>Progress</span>
            <span>{progressPercent.toFixed(1)}%</span>
          </div>
          <div className="w-full bg-gray-200 rounded-full h-3">
            <div
              className={`h-3 rounded-full transition-all ${
                progressPercent > 100 ? "bg-red-500" : "bg-blue-600"
              }`}
              style={{ width: `${Math.min(100, progressPercent)}%` }}
            />
          </div>
          {progressPercent > 100 && (
            <p className="text-sm text-red-600 mt-2">
              Over budget by{" "}
              {currentProject.totalActualHours -
                currentProject.totalPlannedHours}{" "}
              hours
            </p>
          )}
        </div>
      </Card>

      {/* Goals */}
      <div className="mb-6">
        <GoalList projectId={currentProject.id} />
      </div>

      {/* Notes & Meetings */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-6">
        <NoteList projectId={currentProject.id} />
        <MeetingList projectId={currentProject.id} />
      </div>

      {/* Weekly Breakdown */}
      <Card padding="none">
        <div className="px-6 pt-6 pb-4">
          <CardHeader
            title="Weekly Breakdown"
            subtitle="Track planned vs actual hours for each week"
          />
          <p className="text-sm text-gray-500">
            Click on actual hours to edit (only available for past and current
            weeks)
          </p>
        </div>
        <WeeklyBreakdown projectId={currentProject.id} />
      </Card>

      {/* Edit Modal */}
      <ProjectEditModal
        isOpen={isEditModalOpen}
        onClose={() => setIsEditModalOpen(false)}
        onSubmit={handleEditSubmit}
        project={currentProject}
      />
    </div>
  );
}
