package models

// NoteWithProject includes the parent project name
type NoteWithProject struct {
	Note
	ProjectName string `json:"projectName"`
}

// GoalWithProject includes the parent project name
type GoalWithProject struct {
	Goal
	ProjectName string `json:"projectName"`
}
