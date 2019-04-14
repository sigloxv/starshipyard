package views

import (
	html "github.com/multiverse-os/starshipyard/framework/html"
	//forms "github.com/multiverse-os/starshipyard/views/components/forms"
)

func Root() html.Element {
	return DefaultTemplate("root",
		html.Div.Containing(
			html.P.Text("we will put columns here in the future"),
		),
	)
}

func Login() html.Element {
	return DefaultTemplate("login", 
	html.Div.Containing(
		html.P.Text(),
	),
}


////////////////////////////////////////////////////////////////////////////////
// Forms

//func LoginForm(action string) string {
//	return ""
//	//return form.NewForm(
//	//	action,
//	//	// NOTE: Not a big fan of this ContainerClass(), .. other classes
//	//	// that the user is expected to know but really are at least somewhat
//	//	// specific, and long ontop of it
