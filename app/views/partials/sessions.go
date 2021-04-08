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
	return html.Form.Class("form").Method("POST").Action("/sessions/new").Containing(
		form.TextInput("uid", "Username"),
		form.PasswordInput("pwd", "Password"),
		html.Div.Class("columns").Containing(
			html.Div.Class("column").Containing(
				ButtonInput("submit", "Login"),
			),
			html.Div.Class("column").Containing(
				html.A.Class("button", "register-button").Href("/register").Text("Register"),
			),
		),
	)
}

func ButtonInput(name, text string) html.Element {
	// TODO: Validate name/placeholder
	return html.Div.Class("field", "button-input").Containing(
		html.Div.Class("control").Containing(
			//html.Input.Name(name).Value(text).Class("input", ("input-"+name), "login-form").Type("button"),
			html.Input.Value(text).Class("input", ("input-"+name), "login-form").Type("submit"),
			html.Div.Class("error"),
		),
	)
}
