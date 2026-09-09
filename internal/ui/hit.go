package ui

import tea "github.com/charmbracelet/bubbletea"

type Rect struct {
	X      int
	Y      int
	Width  int
	Height int
}

func (r Rect) Contains(x, y int) bool {
	return x >= r.X && x < r.X+r.Width && y >= r.Y && y < r.Y+r.Height
}

type zone struct {
	id   string
	rect Rect
}

type HitMap struct {
	zones []zone
}

func NewHitMap() *HitMap {
	return &HitMap{}
}

func (h *HitMap) Reset() {
	if h == nil {
		return
	}
	h.zones = h.zones[:0]
}

func (h *HitMap) Add(id string, r Rect) {
	if h == nil || id == "" {
		return
	}
	h.zones = append(h.zones, zone{id: id, rect: r})
}

func (h *HitMap) At(x, y int) (string, Rect, bool) {
	if h == nil {
		return "", Rect{}, false
	}
	for i := len(h.zones) - 1; i >= 0; i-- {
		if h.zones[i].rect.Contains(x, y) {
			return h.zones[i].id, h.zones[i].rect, true
		}
	}
	return "", Rect{}, false
}

type ClickMsg struct {
	Panel string
	X     int
	Y     int
}

type WheelMsg struct {
	Panel string
	Delta int
}

func WheelDelta(msg tea.MouseMsg) int {
	switch msg.Button {
	case tea.MouseButtonWheelUp:
		return -1
	case tea.MouseButtonWheelDown:
		return 1
	default:
		return 0
	}
}

func IsClick(msg tea.MouseMsg) bool {
	return msg.Button == tea.MouseButtonLeft && msg.Action == tea.MouseActionPress
}
