package example

type itemsLoadedMsg struct {
	items []Item
}

type itemChangedMsg struct {
	action string
	title  string
}
