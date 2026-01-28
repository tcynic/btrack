package system

import (
	"encoding/csv"
	"fmt"
	"strings"
	"time"
)

// ExportWeeklyReport exports weekly hours data as CSV
func (s *Service) ExportWeeklyReport(startDate, endDate string) (string, error) {
	// Validate dates
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return "", fmt.Errorf("invalid start date: %w", err)
	}
	
	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		return "", fmt.Errorf("invalid end date: %w", err)
	}
	
	if start.After(end) {
		return "", fmt.Errorf("start date must be before end date")
	}

	// Get all projects and their weekly entries
	projects, err := s.store.GetAllProjects(false)
	if err != nil {
		return "", err
	}

	// Build CSV
	var buf strings.Builder
	writer := csv.NewWriter(&buf)

	// Write header
	headers := []string{"Project Name", "Week Start", "Week Number", "Planned Hours", "Actual Hours", "Variance"}
	if err := writer.Write(headers); err != nil {
		return "", fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write data rows
	for _, proj := range projects {
		for _, entry := range proj.WeeklyEntries {
			// Check if entry is within date range
			entryDate, _ := time.Parse("2006-01-02", entry.WeekStartDate)
			if (entryDate.Equal(start) || entryDate.After(start)) && (entryDate.Equal(end) || entryDate.Before(end)) {
				variance := entry.PlannedHours - entry.ActualHours
				row := []string{
					proj.Name,
					entry.WeekStartDate,
					fmt.Sprintf("%d", entry.WeekNumber),
					fmt.Sprintf("%d", entry.PlannedHours),
					fmt.Sprintf("%d", entry.ActualHours),
					fmt.Sprintf("%d", variance),
				}
				if err := writer.Write(row); err != nil {
					return "", fmt.Errorf("failed to write CSV row: %w", err)
				}
			}
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", fmt.Errorf("CSV writer error: %w", err)
	}

	return buf.String(), nil
}

// ExportProjectSummary exports a single project's data as CSV
func (s *Service) ExportProjectSummary(projectID int64) (string, error) {
	// Get project details
	project, err := s.projectService.GetByID(s.ctx, projectID)
	if err != nil {
		return "", err
	}

	// Get project with nested data from store
	projWithNested, err := s.store.GetProject(projectID)
	if err != nil {
		return "", err
	}

	// Build CSV
	var buf strings.Builder
	writer := csv.NewWriter(&buf)

	// Write project info section
	infoHeaders := []string{"Project Name", "Total Sold Hours", "Specialist Hours", "My Hours", "Start Date", "End Date", "Status"}
	if err := writer.Write(infoHeaders); err != nil {
		return "", fmt.Errorf("failed to write info header: %w", err)
	}

	status := "Active"
	if !project.IsActive {
		status = "Inactive"
	}

	infoRow := []string{
		project.Name,
		fmt.Sprintf("%d", project.TotalSoldHours),
		fmt.Sprintf("%d", project.SpecialistHours),
		fmt.Sprintf("%d", project.MyHours),
		project.StartDate,
		project.EndDate,
		status,
	}
	if err := writer.Write(infoRow); err != nil {
		return "", fmt.Errorf("failed to write info row: %w", err)
	}

	// Empty row separator
	writer.Write([]string{})

	// Write weekly data header
	weeklyHeaders := []string{"Week Start", "Week Number", "Planned Hours", "Actual Hours", "Variance"}
	if err := writer.Write(weeklyHeaders); err != nil {
		return "", fmt.Errorf("failed to write weekly header: %w", err)
	}

	// Write weekly data
	for _, entry := range projWithNested.WeeklyEntries {
		variance := entry.PlannedHours - entry.ActualHours
		row := []string{
			entry.WeekStartDate,
			fmt.Sprintf("%d", entry.WeekNumber),
			fmt.Sprintf("%d", entry.PlannedHours),
			fmt.Sprintf("%d", entry.ActualHours),
			fmt.Sprintf("%d", variance),
		}

		if err := writer.Write(row); err != nil {
			return "", fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", fmt.Errorf("CSV writer error: %w", err)
	}

	return buf.String(), nil
}

// ExportAllProjects exports summary data for all active projects as CSV
func (s *Service) ExportAllProjects() (string, error) {
	projects, err := s.projectService.GetAll(s.ctx, true) // activeOnly = true
	if err != nil {
		return "", err
	}

	// Build CSV
	var buf strings.Builder
	writer := csv.NewWriter(&buf)

	// Write header
	headers := []string{
		"Project Name",
		"Total Sold Hours",
		"Specialist Hours",
		"My Hours",
		"Total Weeks",
		"Total Planned Hours",
		"Total Actual Hours",
		"Hours Remaining",
		"Start Date",
		"End Date",
		"Status",
	}
	if err := writer.Write(headers); err != nil {
		return "", fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write data rows
	for _, project := range projects {
		hoursRemaining := project.TotalPlannedHours - project.TotalActualHours
		status := "Active"
		if !project.IsActive {
			status = "Inactive"
		}

		row := []string{
			project.Name,
			fmt.Sprintf("%d", project.TotalSoldHours),
			fmt.Sprintf("%d", project.SpecialistHours),
			fmt.Sprintf("%d", project.MyHours),
			fmt.Sprintf("%d", project.TotalWeeks),
			fmt.Sprintf("%d", project.TotalPlannedHours),
			fmt.Sprintf("%d", project.TotalActualHours),
			fmt.Sprintf("%d", hoursRemaining),
			project.StartDate,
			project.EndDate,
			status,
		}

		if err := writer.Write(row); err != nil {
			return "", fmt.Errorf("failed to write CSV row: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", fmt.Errorf("CSV writer error: %w", err)
	}

	return buf.String(), nil
}
