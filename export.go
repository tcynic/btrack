package main

// ExportWeeklyReport exports weekly hours data as CSV
func (a *App) ExportWeeklyReport(startDate, endDate string) (string, error) {
	return a.systemService.ExportWeeklyReport(startDate, endDate)
}

// ExportProjectSummary exports a single project's data as CSV
func (a *App) ExportProjectSummary(projectID int64) (string, error) {
	return a.systemService.ExportProjectSummary(projectID)
}

// ExportAllProjects exports summary data for all active projects as CSV
func (a *App) ExportAllProjects() (string, error) {
	return a.systemService.ExportAllProjects()
}
