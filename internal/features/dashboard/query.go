package dashboard

import (
	"time"

	"github.com/remcostoeten/reusable-tui/internal/ui"
)

const (
	daysPerWeek = 7
	modeCount   = 2
)

var weekShape = map[Mode][]float64{
	ModeExpense: {18, 42, 12, 64, 30, 88, 54},
	ModeIncome:  {0, 120, 0, 0, 60, 0, 210},
}

var dayLabels = []string{"S", "M", "T", "W", "T", "F", "S"}

func (m *Model) weekStart() time.Time {
	day := time.Date(m.selected.Year(), m.selected.Month(), m.selected.Day(), 0, 0, 0, 0, m.selected.Location())
	return day.AddDate(0, 0, -int(day.Weekday()))
}

func (m *Model) series() []ui.Bar {
	shape := weekShape[m.mode]
	offset := m.weekStart().YearDay() % daysPerWeek
	bars := make([]ui.Bar, 0, daysPerWeek)
	for i := 0; i < daysPerWeek; i++ {
		bars = append(bars, ui.Bar{Label: dayLabels[i], Value: shape[(i+offset)%daysPerWeek]})
	}
	return bars
}

func (m *Model) total() float64 {
	sum := 0.0
	for _, bar := range m.series() {
		sum += bar.Value
	}
	return sum
}

func (m *Model) perDay() float64 {
	return m.total() / float64(daysPerWeek)
}

func (m *Model) balance() float64 {
	sum := 0.0
	for _, account := range m.accounts {
		sum += account.Balance
	}
	return sum
}

func (m *Model) periodLabel() string {
	start := m.weekStart()
	end := start.AddDate(0, 0, daysPerWeek-1)
	if sameWeek(start, m.clock()) {
		return "This Week"
	}
	return start.Format("Jan 2") + " - " + end.Format("Jan 2")
}

func sameWeek(start, other time.Time) bool {
	day := time.Date(other.Year(), other.Month(), other.Day(), 0, 0, 0, 0, other.Location())
	otherStart := day.AddDate(0, 0, -int(day.Weekday()))
	return start.Year() == otherStart.Year() && start.YearDay() == otherStart.YearDay()
}

func (m Mode) String() string {
	if m == ModeIncome {
		return "Income"
	}
	return "Expense"
}
