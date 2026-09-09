package example

type itemsLoadedMsg struct {
	items []Item
	err   error
}

type itemChangedMsg struct {
	action string
	title  string
}
