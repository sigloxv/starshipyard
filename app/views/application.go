package views

import (
	html "github.com/multiverse-os/starshipyard/framework/html"
)

func Root() html.Element {
	return DefaultTemplate("",
		html.Div.Class("content").Containing(
			html.Section.Class("section is-fullwidth is-primary").Containing(
				html.H1.Class("title").Text("Starship"),
				html.H5.Class("subtitle").Text("A web application framework designed for simplicity, single binary, single response, inspired by rails"),
			),
		),
	)
}
