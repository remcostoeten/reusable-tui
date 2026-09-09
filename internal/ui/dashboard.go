package ui

import "github.com/remcostoeten/reusable-tui/internal/theme"

type BodyFunc func(t theme.Theme, width, height int) string

type Region struct {
	Panel  Panel
	Weight int
	Body   BodyFunc
}

type Stack struct {
	Weight  int
	Regions []Region
}

type Dashboard struct {
	Width  int
	Height int
	Stacks []Stack
	Hits   *HitMap
}

func RenderDashboard(t theme.Theme, d Dashboard) string {
	if len(d.Stacks) == 0 || d.Width <= 0 || d.Height <= 0 {
		return ""
	}
	widths := SplitWidths(d.Width, t.Space.Gutter, stackWeights(d.Stacks)...)
	columns := make([]string, 0, len(d.Stacks))
	x := 0
	for i, stack := range d.Stacks {
		columns = append(columns, renderStack(t, stack, d.Hits, x, widths[i], d.Height))
		x += widths[i] + t.Space.Gutter
	}
	return joinColumns(columns, Repeat(" ", t.Space.Gutter))
}

func renderStack(t theme.Theme, s Stack, hits *HitMap, x, width, height int) string {
	if len(s.Regions) == 0 {
		return ""
	}
	heights := SplitHeights(height, 0, regionWeights(s.Regions)...)
	blocks := make([]string, 0, len(s.Regions))
	y := 0
	for i, region := range s.Regions {
		hits.Add(region.Panel.ID, Rect{X: x, Y: y, Width: width, Height: heights[i]})
		blocks = append(blocks, renderRegion(t, region, width, heights[i]))
		y += heights[i]
	}
	return Column(blocks...)
}

func renderRegion(t theme.Theme, r Region, width, height int) string {
	panel := r.Panel
	panel.Width = width
	panel.Height = height
	if r.Body != nil {
		panel.Body = r.Body(t, panel.ContentWidth(t), panel.ContentHeight(t))
	}
	return RenderPanel(t, panel)
}

func stackWeights(stacks []Stack) []int {
	out := make([]int, 0, len(stacks))
	for _, stack := range stacks {
		out = append(out, weightOr(stack.Weight))
	}
	return out
}

func regionWeights(regions []Region) []int {
	out := make([]int, 0, len(regions))
	for _, region := range regions {
		out = append(out, weightOr(region.Weight))
	}
	return out
}

func weightOr(weight int) int {
	if weight <= 0 {
		return 1
	}
	return weight
}
