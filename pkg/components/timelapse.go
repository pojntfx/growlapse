package components

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/studio-b12/gowebdav"
)

type Timelapse struct {
	app.Compo

	WebDAVClient *gowebdav.Client
}

func (c *Timelapse) Render() app.UI {
	return app.H2().Text("Timelapse")
}

func (c *Timelapse) OnMount(ctx app.Context) {
	if _, err := c.WebDAVClient.ReadDir("/"); err != nil {
		panic(err)
	}
}
