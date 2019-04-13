package views

import (
	css "github.com/multiverse-os/starshipyard/views/assets/css"
	// TODO: Consider importing with `.` for
	//template "github.com/multiverse-os/starshipyard/views/template"
	//navigation "github.com/multiverse-os/starshipyard/views/template/navigation"

	html "github.com/multiverse-os/starshipyard/framework/html"
)

// TODO: Views should use this default template, and extend it for each page by
// adding sections.
// TODO: Needs a helper (reference navigation under template) for oooookkkk
//
func DefaultTemplate(title string, yield html.Element) html.Element {
	return html.HTML.Containing(
		html.Head.Containing(
			html.Meta.Charset("UTF-8"),
			html.Meta.Name("description").Content("A powerful web application framework"),
			html.Meta.Name("keywords").Content("powerful, go language, webapp, starship"),
			html.Meta.Name("viewport").Content("width=device-width, initial-scale=1.0"),
			html.Title.Text(title),
			html.Style.Text(css.DefaultCSS()),
		),
		html.Body.Containing(
			html.Nav.Id("navbar").Class("navbar").Containing(
				html.Div.Class("container").Containing(
					html.Div.Class("navbar-brand").Containing(
						html.H5.Class("title").Text("Starship"),
					),
					html.Div.Class("navbar-menu").Containing(
						html.Div.Class("navbar-start").Containing(
							html.A.Class("navbar-item").Href("/").Text("Home"),
							html.A.Class("navbar-item").Href("/about").Text("About"),
							html.A.Class("navbar-item").Href("/contact").Text("Contact"),
							html.A.Class("navbar-item").Href("/login").Text("Login"),
							html.A.Class("navbar-item").Href("/register").Text("Register"),
						),
						html.Div.Class("navbar-end").Containing(
							html.A.Class("navbar-item").Href("https://github.com/multiverse-os/starshipyard").Text("Github"),
						),
					), // NavbarMenu
				), // Container
			), // Navbar
			html.Main.Class("bd-main").Containing(
				html.Div.Class("bd-side-background"),
				html.Div.Class("bd-main-container", "container").Containing(
					html.Div.Class("columns").Containing(
						html.Div.Class("column", "is-four-fifths").Containing(
							html.P.Text("test"),
						),
						html.Div.Class("column", "is-one-fifth", "sidebar").Containing(
							html.P.Text("test sidebar"),
						),
					),
					html.Div.Class("bd-duo").Containing(
						html.Div.Class("bd-lead").Containing(
							html.Div.Class("bd-breadcrumb").Containing(
								html.Nav.Class("breadcrumb").Text("test"),
							),
						), // bdLead
						html.Aside.Class("bd-side").Text("sidebar"),
					), // bdDuo
				), // bdMainContainer
			),
		), // Body
	)
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
//	//	form.TextInput("uid", "Username").ContainerClass("is-5"),
//	//	form.PasswordInput("password", "Password").ContainerClass("is-5"),
//	//	form.ButtonInput("Login").ContainerClass("is-1"),
//	//	form.ButtonInput("Register").ContainerClass("is-1"),
//	//).ContainerClass("field", "column", "is-1").InputDivClass("control").InputsClass("login", "columns", "is-gapless").HTML()
