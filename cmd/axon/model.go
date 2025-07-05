package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/Advik-B/Axon/pkg/axon"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// buildStage represents a single step in the build process.
type buildStage struct {
	name     string
	status   status
	duration time.Duration
}

// status defines the state of a build stage.
type status int

const (
	statusPending status = iota
	statusRunning
	statusSuccess
	statusFailure
)

// Model is the core state of our TUI application.
type Model struct {
	filePath     string
	stages       []buildStage
	currentStage int
	spinner      spinner.Model
	done         bool
	err          error
	outputCode   string
	graph        *axon.Graph
	styles       *Styles
}

func initialModel(filePath string) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return Model{
		filePath: filePath,
		stages: []buildStage{
			{name: "Parsing graph file", status: statusRunning},
			{name: "Transpiling to Go", status: statusPending},
			{name: "Writing output file", status: statusPending},
		},
		currentStage: 0,
		spinner:      s,
		styles:       DefaultStyles(),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		runParseStage(m.filePath),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	// --- Custom Messages for Build Stages ---
	case parseResultMsg:
		return m.handleStageCompletion(msg.duration, msg.err, func() tea.Cmd {
			m.graph = msg.graph
			return runTranspileStage(m.graph)
		})

	case transpileResultMsg:
		return m.handleStageCompletion(msg.duration, msg.err, func() tea.Cmd {
			m.outputCode = msg.code
			return runWriteStage(m.outputCode, "out/main.go")
		})

	case writeResultMsg:
		return m.handleStageCompletion(msg.duration, msg.err, nil)

	default:
		return m, nil
	}
}

// handleStageCompletion is a helper to reduce boilerplate in the Update function.
func (m *Model) handleStageCompletion(duration time.Duration, err error, nextCmd func() tea.Cmd) (tea.Model, tea.Cmd) {
	m.stages[m.currentStage].duration = duration
	if err != nil {
		m.err = err
		m.stages[m.currentStage].status = statusFailure
		m.done = true
		return *m, tea.Quit
	}

	m.stages[m.currentStage].status = statusSuccess
	m.currentStage++

	if m.currentStage >= len(m.stages) {
		m.done = true
		return *m, tea.Quit
	}

	if nextCmd != nil {
		m.stages[m.currentStage].status = statusRunning
		return *m, nextCmd()
	}

	return *m, nil
}

func (m Model) View() string {
	if !m.done && m.err == nil {
		var sb strings.Builder
		sb.WriteString("🚀 Building Axon graph: " + m.styles.FileName.Render(m.filePath) + "\n\n")

		for _, stage := range m.stages {
			var icon string
			var style lipgloss.Style
			switch stage.status {
			case statusPending:
				icon = " "
				style = m.styles.Pending
			case statusRunning:
				icon = m.spinner.View()
				style = m.styles.Running
			case statusSuccess:
				icon = "✓"
				style = m.styles.Success
			case statusFailure:
				icon = "✗"
				style = m.styles.Failure
			}

			duration := ""
			if stage.duration > 0 {
				duration = m.styles.Duration.Render(fmt.Sprintf("(%.2fs)", stage.duration.Seconds()))
			}
			sb.WriteString(fmt.Sprintf("%s %s %s\n", icon, style.Render(stage.name), duration))
		}
		sb.WriteString("\n" + m.styles.Faint.Render("Press Ctrl+C to exit."))
		return sb.String()
	}

	// Final view on success or failure
	if m.err != nil {
		return m.styles.ErrorBox.Render(fmt.Sprintf("🔥 Error: %v", m.err))
	}

	var success strings.Builder
	success.WriteString(m.styles.Success.Render("✓ Build Succeeded!\n\n"))
	success.WriteString(fmt.Sprintf("Go code written to %s\n", m.styles.FileName.Render("out/main.go")))
	success.WriteString("Run it with: " + m.styles.Command.Render("go run out/main.go") + "\n\n")
	success.WriteString(m.styles.CodeBox.Render(m.outputCode))

	return m.styles.SuccessBox.Render(success.String())
}