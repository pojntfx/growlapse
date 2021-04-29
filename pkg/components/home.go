package components

import (
	"net/url"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type Home struct {
	app.Compo

	webdavURL      string
	webdavUsername string
	webdavPassword string
}

const (
	webdavURLKey      = "webdavURL"
	webdavUsernameKey = "webdavUsername"
	webdavPasswordKey = "webdavPassword"
)

func (c *Home) Render() app.UI {
	return app.Form().
		OnSubmit(func(ctx app.Context, e app.Event) {
			e.PreventDefault()

			u, err := url.Parse("/")
			if err != nil {
				panic(err)
			}

			q := u.Query()
			q.Add(webdavURLKey, c.webdavURL)
			q.Add(webdavUsernameKey, c.webdavUsername)
			u.RawQuery = q.Encode()

			ctx.NavigateTo(u)
		}).
		Body(
			app.Label().
				For(webdavURLKey).
				Text("WebDAV URL: "),
			app.Input().
				Type("url").
				ID(webdavURLKey).
				Name(webdavURLKey).
				Required(true).
				Value(c.webdavURL).
				OnInput(func(ctx app.Context, e app.Event) {
					c.webdavURL = ctx.JSSrc.Get("value").String()
				}),
			app.Br(),
			app.Label().
				For(webdavUsernameKey).
				Text("WebDAV Username: "),
			app.Input().
				Type("text").
				ID(webdavUsernameKey).
				Name(webdavUsernameKey).
				Value(c.webdavUsername).
				Required(true).
				OnInput(func(ctx app.Context, e app.Event) {
					c.webdavUsername = ctx.JSSrc.Get("value").String()
				}),
			app.Br(),
			app.Label().
				For(webdavPasswordKey).
				Text("WebDAV Password: "),
			app.Input().
				Type("password").
				ID(webdavPasswordKey).
				Name(webdavPasswordKey).
				Value(c.webdavPassword).
				OnInput(func(ctx app.Context, e app.Event) {
					c.webdavPassword = ctx.JSSrc.Get("value").String()
				}),
			app.Br(),
			app.Button().
				Type("submit").
				Text("Login"),
		)
}

func (c *Home) OnMount(ctx app.Context) {
	q := ctx.Page.URL().Query()

	if webdavURL := q.Get(webdavURLKey); webdavURL != "" {
		c.webdavURL = webdavURL
	}

	if webdavUsername := q.Get(webdavUsernameKey); webdavUsername != "" {
		c.webdavUsername = webdavUsername
	}
}
