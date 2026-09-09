package app

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/remcostoeten/reusable-tui/internal/keymap"
	"github.com/remcostoeten/reusable-tui/internal/theme"
	"github.com/remcostoeten/reusable-tui/internal/ui"
)

func TestTabCycling(t *testing.T) {
	db := newTestStore(t)
	final := runProgram(t, newTestModel(t, db), special(tea.KeyTab))
	if got := final.ActiveScreen().ID(); got != helpScreenID {
		t.Fatalf("tab did not advance, active screen is %q", got)
	}
	final = runProgram(t, newTestModel(t, db), special(tea.KeyTab), special(tea.KeyShiftTab))
	if got := final.ActiveScreen().ID(); got != "example" {
		t.Fatalf("shift-tab did not return, active screen is %q", got)
	}
}

func TestFocusCycling(t *testing.T) {
	db := newTestStore(t)
	final := runProgram(t, newTestModel(t, db), special(tea.KeyCtrlN))
	if got := final.focus["example"]; got != keymap.PanelExampleDetail {
		t.Fatalf("focus did not advance, focused panel is %q", got)
	}
	final = runProgram(t, newTestModel(t, db), special(tea.KeyCtrlN), special(tea.KeyCtrlP))
	if got := final.focus["example"]; got != keymap.PanelExampleList {
		t.Fatalf("focus did not return, focused panel is %q", got)
	}
}

func TestJumpMode(t *testing.T) {
	db := newTestStore(t)
	final := runProgram(t, newTestModel(t, db), runes("f"), runes("s"))
	if final.jump.Active() {
		t.Fatal("jump mode still active after selection")
	}
	if got := final.focus["example"]; got != keymap.PanelExampleDetail {
		t.Fatalf("jump did not focus detail panel, focused panel is %q", got)
	}
	final = runProgram(t, newTestModel(t, db), runes("f"), special(tea.KeyEsc))
	if final.jump.Active() {
		t.Fatal("escape did not cancel jump mode")
	}
	if got := final.focus["example"]; got != keymap.PanelExampleList {
		t.Fatalf("escape changed focus to %q", got)
	}
}

func TestPaletteSearch(t *testing.T) {
	db := newTestStore(t)
	final := runProgram(t, newTestModel(t, db), special(tea.KeyCtrlK), runes("h"), runes("e"), runes("l"))
	matches := final.palette.Matches()
	if len(matches) == 0 {
		t.Fatal("palette returned no matches")
	}
	if !strings.Contains(matches[0].Label, "Help") {
		t.Fatalf("unexpected top match %q", matches[0].Label)
	}
}

const monochromeAccent = "38;2;237;237;237"

func TestThemeSwitching(t *testing.T) {
	db := newTestStore(t)
	final := runProgramUntil(t, newTestModel(t, db), monochromeAccent,
		special(tea.KeyCtrlK),
		runes("m"), runes("o"), runes("n"), runes("o"),
		special(tea.KeyEnter),
	)
	if got := final.Theme().Name; got != theme.NameMonochrome {
		t.Fatalf("theme not switched, active theme is %q", got)
	}
	if final.palette.Open() {
		t.Fatal("palette still open after running a command")
	}
}

func TestWindowResize(t *testing.T) {
	db := newTestStore(t)
	final := runProgram(t, newTestModel(t, db), tea.WindowSizeMsg{Width: 140, Height: 40})
	if final.width != 140 || final.height != 40 {
		t.Fatalf("resize not applied, got %dx%d", final.width, final.height)
	}
	body := ui.Lines(final.ActiveScreen().View(final.context()))
	if len(body) != final.bodyHeight() {
		t.Fatalf("body height %d does not match layout %d", len(body), final.bodyHeight())
	}
	if width := ui.Width(body[0]); width != 140 {
		t.Fatalf("body width %d does not fill the terminal", width)
	}
}

func TestMinimumSizeGuard(t *testing.T) {
	db := newTestStore(t)
	m := newTestModel(t, db)
	m.Update(tea.WindowSizeMsg{Width: 20, Height: 6})
	if !strings.Contains(m.View(), "terminal too small") {
		t.Fatal("minimum size guard did not render")
	}
}

func TestErrorToast(t *testing.T) {
	db := newTestStore(t)
	m := settled(t, newTestModel(t, db))
	m.Update(ui.ErrorMsg{Err: errUnknownTheme("nope")})
	if !m.toast.Visible {
		t.Fatal("error message did not raise a toast")
	}
	rows := ui.Lines(m.View())
	if len(rows) != testHeight {
		t.Fatalf("frame height %d does not match terminal height %d", len(rows), testHeight)
	}
	for i, row := range rows {
		if width := ui.Width(row); width != testWidth {
			t.Fatalf("row %d width %d does not match terminal width %d", i, width, testWidth)
		}
	}
	m.Update(ui.ToastExpiredMsg{Seq: m.toast.Seq})
	if m.toast.Visible {
		t.Fatal("toast did not expire")
	}
}
