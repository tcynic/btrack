package main

import (
	"encoding/csv"
	"fmt"
	"strings"
	"time"

	"btrack/internal/database"
)

// ExportWeeklyReport exports weekly hours data as CSV
func (a *App) ExportWeeklyReport(startDate, endDate string) (string, error) {
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

	// Query weekly entries
	rows, err := a.db.Query(database.SelectWeeklyEntriesByDateRange, startDate, endDate)
	if err != nil {
		return "", fmt.Errorf("failed to query weekly entries: %w", err)
	}
	defer rows.Close()

	// Build CSV
	var buf strings.Builder
	writer := csv.NewWriter(&buf)

	// Write header
	headers := []string{"Week Start", "Week Number", "Project ID", "Planned Hours", "Actual Hours", "Variance"}
	if err := writer.Write(headers); err != nil {
		return "", fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write data rows
	for rows.Next() {
		var projectID, weekNumber, planned, actual int64
		var weekStart, createdAt, updatedAt string

		err := rows.Scan(&projectID, &projectID, &weekStart, &weekNumber, &planned, &actual, &createdAt, &updatedAt)
		if err != nil {
			return "", fmt.Errorf("failed to scan row: %w", err)
		}

		variance := planned - actual
		row := []string{
			weekStart,
			fmt.Sprintf("%d", weekNumber),
			fmt.Sprintf("%d", projectID),
			fmt.Sprintf("%d", planned),
			fmt.Sprintf("%d", actual),
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

// ExportProjectSummary exports a single project's data as CSV
func (a *App) ExportProjectSummary(projectID int64) (string, error) {
	// Get project details
	project, err := a.GetProject(projectID)
	if err != nil {
		return "", err
	}

	// Get weekly entries
	rows, err := a.db.Query(database.SelectWeeklyEntriesByProject, projectID)
	if err != nil {
		return "", fmt.Errorf("failed to query weekly entries: %w", err)
	}
	defer rows.Close()

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
	for rows.Next() {
		var id, projectID, weekNumber, planned, actual int64
		var weekStart, createdAt, updatedAt string

		err := rows.Scan(&id, &projectID, &weekStart, &weekNumber, &planned, &actual, &createdAt, &updatedAt)
		if err != nil {
			return "", fmt.Errorf("failed to scan row: %w", err)
		}

		variance := planned - actual
		row := []string{
			weekStart,
			fmt.Sprintf("%d", weekNumber),
			fmt.Sprintf("%d", planned),
			fmt.Sprintf("%d", actual),
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
func (a *App) ExportAllProjects() (string, error) {
	projects, err := a.GetAllProjects(true) // activeOnly = true
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
