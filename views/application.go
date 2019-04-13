package views

import (
	html "github.com/multiverse-os/starshipyard/framework/html"
)

func Root() string {
	return DefaultTemplate("root", html.Div.Containing(
		html.Div.Text("we will put columns here in the future"),
	),
	)
}
