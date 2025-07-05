package main

import "github.com/charmbracelet/lipgloss"

// Styles holds all the lipgloss styles for our TUI.
type Styles struct {
	Pending    lipgloss.Style
	Running    lipgloss.Style
	Success    lipgloss.Style
	Failure    lipgloss.Style
	Faint      lipgloss.Style
	Duration   lipgloss.Style
	FileName   lipgloss.Style
	Command    lipgloss.Style
	ErrorBox   lipgloss.Style
	SuccessBox lipgloss.Style
	CodeBox    lipgloss.Style
}

// DefaultStyles returns a new Styles struct with default values.
func DefaultStyles() *Styles {
	return &Styles{
		Pending:    lipgloss.NewStyle().SetString("  ").Foreground(lipgloss.Color("240")),
		Running:    lipgloss.NewStyle().Foreground(lipgloss.Color("#22a7f0")),
		Success:    lipgloss.NewStyle().SetString("✓").Foreground(lipgloss.Color("#34d399")),
		Failure:    lipgloss.NewStyle().SetString("✗").Foreground(lipgloss.Color("#f87171")),
		Faint:      lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Italic(true),
		Duration:   lipgloss.NewStyle().Foreground(lipgloss.Color("240")),
		FileName:   lipgloss.NewStyle().Foreground(lipgloss.Color("#fde047")).Bold(true),
		Command:    lipgloss.NewStyle().Foreground(lipgloss.Color("#60a5fa")).Bold(true),
		ErrorBox:   lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#f87171")).Padding(1, 2),
		SuccessBox: lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).BorderForeground(lipgloss.Color("#34d399")).Padding(1, 2),
		CodeBox:    lipgloss.NewStyle().Border(lipgloss.NormalBorder(), false, false, false, true).BorderForeground(lipgloss.Color("238")).Padding(1, 2).MarginTop(1),
	}
}