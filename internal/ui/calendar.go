package ui

import (
	"strconv"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/remcostoeten/reusable-tui/internal/theme"
)

const (
	calendarColumns  = 7
	calendarRows     = 6
	calendarCellSize = 3
)

var weekdayLabels = []string{"S", "M", "T", "W", "T", "F", "S"}

type Calendar struct {
	Month     time.Time
	Selected  time.Time
	Width     int
	Focused   bool
	WeekBand  bool
	MarkToday bool
	Today     time.Time
}

func RenderCalendar(t theme.Theme, c Calendar) string {
	if c.Width < calendarColumns*calendarCellSize {
		return ""
	}
	rows := []string{calendarHeader(t, c.Width)}
	start := calendarStart(c.Month)
	selectedWeek := weekIndex(start, c.Selected)
	for week := 0; week < calendarRows; week++ {
		rows = append(rows, calendarWeek(t, c, start, week, week == selectedWeek))
	}
	return Join(rows)
}

func RenderPeriod(t theme.Theme, width int, label string, focused bool) string {
	arrow := lipgloss.NewStyle().Foreground(t.Text.Muted).Background(t.Base.Background)
	if focused {
		arrow = lipgloss.NewStyle().Foreground(t.Accent.Active).Background(t.Base.Background)
	}
	title := lipgloss.NewStyle().Foreground(t.Text.Primary).Background(t.Base.Background)
	inner := width - 8
	if inner < 0 {
		inner = 0
	}
	return arrow.Render("<<<") + title.Render(Center(Truncate(label, inner), inner)) + arrow.Render(">>>") +
		lipgloss.NewStyle().Background(t.Base.Background).Render(Repeat(" ", width-inner-6))
}

func calendarHeader(t theme.Theme, width int) string {
	style := lipgloss.NewStyle().Foreground(t.Text.Muted).Background(t.Base.Background).Italic(true)
	row := ""
	for _, label := range weekdayLabels {
		row += style.Render(PadLeft(label, calendarCellSize))
	}
	return row + lipgloss.NewStyle().Background(t.Base.Background).Render(Repeat(" ", width-calendarColumns*calendarCellSize))
}

func calendarWeek(t theme.Theme, c Calendar, start time.Time, week int, banded bool) string {
	background := t.Base.Background
	if banded && c.WeekBand {
		background = t.Base.Surface
	}
	row := ""
	for day := 0; day < calendarColumns; day++ {
		date := start.AddDate(0, 0, week*calendarColumns+day)
		row += calendarCell(t, c, date, background)
	}
	tail := c.Width - calendarColumns*calendarCellSize
	return row + lipgloss.NewStyle().Background(background).Render(Repeat(" ", tail))
}

func calendarCell(t theme.Theme, c Calendar, date time.Time, background lipgloss.TerminalColor) string {
	label := PadLeft(strconv.Itoa(date.Day()), calendarCellSize)
	style := lipgloss.NewStyle().Foreground(t.Text.Secondary).Background(background)
	if date.Month() != c.Month.Month() {
		style = lipgloss.NewStyle().Foreground(t.Text.Disabled).Background(background)
	}
	if c.MarkToday && sameDay(date, c.Today) {
		style = style.Underline(true)
	}
	if sameDay(date, c.Selected) {
		selected := lipgloss.NewStyle().Foreground(t.Text.Primary).Background(t.Accent.Dim)
		if c.Focused {
			selected = lipgloss.NewStyle().Foreground(t.Text.Inverted).Background(t.Accent.Active)
		}
		return selected.Render(label)
	}
	return style.Render(label)
}

func calendarStart(month time.Time) time.Time {
	first := time.Date(month.Year(), month.Month(), 1, 0, 0, 0, 0, month.Location())
	return first.AddDate(0, 0, -int(first.Weekday()))
}

func weekIndex(start, date time.Time) int {
	days := int(truncateDay(date).Sub(truncateDay(start)).Hours() / 24)
	if days < 0 {
		return -1
	}
	return days / calendarColumns
}

func sameDay(a, b time.Time) bool {
	return a.Year() == b.Year() && a.YearDay() == b.YearDay()
}

func truncateDay(v time.Time) time.Time {
	return time.Date(v.Year(), v.Month(), v.Day(), 0, 0, 0, 0, v.Location())
}
