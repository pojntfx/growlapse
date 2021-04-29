package components

import (
	"encoding/base64"
	"net/url"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/studio-b12/gowebdav"
)

type Home struct {
	app.Compo

	webDAVURL      string
	webDAVUsername string
	webDAVPassword string

	loggedIn bool

	webDAVClient *gowebdav.Client
}

const (
	webDAVURLKey      = "webDAVURL"
	webDAVUsernameKey = "webDAVUsername"
	webDAVPasswordKey = "webDAVPassword"
)

func (c *Home) Render() app.UI {
	return app.If(
		c.loggedIn,
		app.Div().
			Body(
				app.Button().
					OnClick(func(ctx app.Context, e app.Event) {
						ctx.Navigate("/")
					}).
					Text("Logout"),
				&Timelapse{
					WebDAVClient: c.webDAVClient,
				},
			),
	).Else(
		app.Form().
			OnSubmit(func(ctx app.Context, e app.Event) {
				e.PreventDefault()

				u, err := url.Parse("/")
				if err != nil {
					panic(err)
				}

				q := u.Query()
				q.Add(webDAVURLKey, c.webDAVURL)
				q.Add(webDAVUsernameKey, c.webDAVUsername)
				q.Add(webDAVPasswordKey, c.webDAVPassword)
				u.RawQuery = q.Encode()

				ctx.NavigateTo(u)
			}).
			Body(
				app.Label().
					For(webDAVURLKey).
					Text("webDAV URL: "),
				app.Input().
					Type("url").
					ID(webDAVURLKey).
					Name(webDAVURLKey).
					Required(true).
					Value(c.webDAVURL).
					OnInput(func(ctx app.Context, e app.Event) {
						c.webDAVURL = ctx.JSSrc.Get("value").String()
					}),
				app.Br(),
				app.Label().
					For(webDAVUsernameKey).
					Text("webDAV Username: "),
				app.Input().
					Type("text").
					ID(webDAVUsernameKey).
					Name(webDAVUsernameKey).
					Value(c.webDAVUsername).
					Required(true).
					OnInput(func(ctx app.Context, e app.Event) {
						c.webDAVUsername = ctx.JSSrc.Get("value").String()
					}),
				app.Br(),
				app.Label().
					For(webDAVPasswordKey).
					Text("webDAV Password: "),
				app.Input().
					Type("password").
					ID(webDAVPasswordKey).
					Name(webDAVPasswordKey).
					Value(c.webDAVPassword).
					OnInput(func(ctx app.Context, e app.Event) {
						c.webDAVPassword = ctx.JSSrc.Get("value").String()
					}),
				app.Br(),
				app.Button().
					Type("submit").
					Text("Login"),
			),
	)

}

func (c *Home) OnNav(ctx app.Context) {
	c.webDAVURL, c.webDAVUsername, c.webDAVPassword = getQuery(ctx)

	if c.webDAVURL != "" && c.webDAVUsername != "" {
		webDAVClient := gowebdav.NewClient(c.webDAVURL, c.webDAVUsername, c.webDAVPassword)
		webDAVClient.SetHeader("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(c.webDAVUsername+":"+c.webDAVPassword)))

		c.webDAVClient = webDAVClient

		c.loggedIn = true
	} else {
		c.loggedIn = false
	}
}

func getQuery(ctx app.Context) (webDAVURL string, webDAVUsername string, webDAVPassword string) {
	q := ctx.Page.URL().Query()

	webDAVURL = q.Get(webDAVURLKey)
	webDAVUsername = q.Get(webDAVUsernameKey)
	webDAVPassword = q.Get(webDAVPasswordKey)

	return
}
