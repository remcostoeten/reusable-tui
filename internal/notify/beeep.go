package notify

import "github.com/gen2brain/beeep"

type Desktop struct {
	AppName string
}

func NewDesktop(appName string) Desktop {
	return Desktop{AppName: appName}
}

func (d Desktop) Notify(title, body string) error {
	return beeep.Notify(d.AppName+": "+title, body, nil)
}
