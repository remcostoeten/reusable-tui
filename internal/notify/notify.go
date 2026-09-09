package notify

type Notifier interface {
	Notify(title, body string) error
}

type Discard struct{}

func (Discard) Notify(string, string) error {
	return nil
}
