package tui

import "github.com/charmbracelet/lipgloss"

var (
	titleStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("13"))
	cursorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("13")).Bold(true)
	selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	descStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	helpStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	statusStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	errStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	appStyle      = lipgloss.NewStyle().Padding(1, 2)
)
