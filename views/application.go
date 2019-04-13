package views

import (
	html "github.com/multiverse-os/starshipyard/framework/html"
)

func Root() html.Element {
	return DefaultTemplate("root",
		html.Div.Containing(html.P.Text("we will put columns here in the future")),
	)
}
