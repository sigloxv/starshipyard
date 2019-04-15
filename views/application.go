package views

import (
	html "github.com/multiverse-os/starshipyard/framework/html"
	form "github.com/multiverse-os/starshipyard/views/components/form"
)

func Root() html.Element {
	return DefaultTemplate("root",
		html.Div.Containing(
			html.P.Text("we will put columns here in the future"),
			html.Div.Class("columns").Containing(
				html.Div.Class("column"),
				html.Div.Class("column").Containing(
					LoginForm(),
				),
				html.Div.Class("column"),
			),
		),
	)
}

func LoginForm() html.Element {
	return html.Form.Class("form").Method("GET").Action("/sessions/new").Containing(
		form.TextInput("uid", "Username"),
		form.PasswordInput("password", "Password"),
	)
}

func PasswordInput(name, placeholder string) html.Element {
	// TODO: Validate name/placeholder
	return html.Div.Class("field", "password-input").Containing(
		html.Div.Class("control").Containing(
			html.Input.Name(name).Placeholder(placeholder).Class("input", ("input-"+name), "is-primary", "login-form").Type("password"),
			html.Div.Class("error"),
		),
	)
}
