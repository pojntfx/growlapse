package components

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type Home struct {
	app.Compo
}

func (c *Home) Render() app.UI {
	return app.H1().Text("Growlapse")
}
