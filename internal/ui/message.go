package ui

import tea "github.com/charmbracelet/bubbletea"

type ErrorMsg struct {
	Err error
}

type InfoMsg struct {
	Text string
}

type SuccessMsg struct {
	Text string
}

type WarnMsg struct {
	Text string
}

type ThemeMsg struct {
	Name string
}

type FocusMsg struct {
	PanelID string
}

type NotifyMsg struct {
	Title string
	Body  string
}

func Fail(err error) tea.Cmd {
	return func() tea.Msg { return ErrorMsg{Err: err} }
}

func Info(text string) tea.Cmd {
	return func() tea.Msg { return InfoMsg{Text: text} }
}

func Success(text string) tea.Cmd {
	return func() tea.Msg { return SuccessMsg{Text: text} }
}

func Warn(text string) tea.Cmd {
	return func() tea.Msg { return WarnMsg{Text: text} }
}

func SetTheme(name string) tea.Cmd {
	return func() tea.Msg { return ThemeMsg{Name: name} }
}

func Focus(panelID string) tea.Cmd {
	return func() tea.Msg { return FocusMsg{PanelID: panelID} }
}

func Notify(title, body string) tea.Cmd {
	return func() tea.Msg { return NotifyMsg{Title: title, Body: body} }
}
