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
					),
				),
			),
		),
		//html.Body.Class("page-columns").Containing(
		//html.Nav.Id("navbar").Class("navbar").Containing(
		//	html.Div.Class("container").Containing(
		//		template.Columns(
		//			template.Column(2, html.H1.Text(title)),
		//			template.Column(6, navigation.NewHeaderMenu(
		//				navigation.NewMenuOption("/", "Home"),
		//				navigation.NewMenuOption("/about", "About"),
		//				navigation.NewMenuOption("/contact", "Contact"),
		//			).HTML(),
		//			),
		//			template.Column(4, html.Div.Text("logggin")),
		//		),
		//		template.Columns(
		//			template.Column(3, html.Div.Text("sidebar")),
		//			template.Column(9, html.Div.Text("main content")),
		//		),
		//		template.Columns(
		//			template.Column(12, html.Div.Text("footer")),
		//		),
		//	),
		//),
	)
}

//func (self *Template) Body() html.Element {
//	return html.Body.ChildElements()
//	Columns(
//		// TODO: Not super excited about this current API
//		Column(2, H1(self.Title)),
//		Column(6, self.HeaderMenu.HTML()),
//		Column(4, LoginForm("/login")),
//	),
//	Columns(
//		Column(12, self.FlashMessages.HTML()),
//	),
//	if self.SidebarMenu.IsEmpty() {
//		body += Content(
//			Columns(
//				Column(12, self.Content),
//			),
//		)
//	} else {
//		// TODO: Use columns to add sidbar, sidebar menu should define the side
//		// and then build the template accordignly
//		// steve ducey cooking show goes here too
//		switch self.SidebarMenu.Type {
//		case LeftSideMenu:
//			body += Content(
//				Columns(
//					/// TODO: This is just for dev, all these values will be defined by
//					// administrators via backend (hehe backend).
//					// TODO: I really don't like this initialization method, feels very
//					// clunky. This column initialization could be better too since we are
//					// telling it twice it is 4 already.
//
//					Column(
//						2,
//						NewLeftSideMenu(
//							2,
//							NewMenuCategory(
//								"/projects",
//								"Projects",
//								NewMenuItem("Comics", "/projects/comics"),
//								NewMenuItem("Music", "/projects/music"),
//								NewMenuItem("Podcasts", "/projects/podcasts"),
//								NewMenuItem("Video", "/projects/video"),
//								NewMenuItem("Software", "/projects/software"),
//							),
//						).HTML()),
//
//					Column(10, Div(Class("main", "tile"), self.Content)),
//				),
//			)
//		case RightSideMenu:
//			body += "right-side-menu"
//			body += Content(
//				Columns(
//					Column(10, self.Content),
//					// TODO: Seems rather silly to indicate back to back that this is 4
//					// wide. Must be a better way
//					Column(
//						2,
//						NewRightSideMenu(
//							2,
//							NewMenuCategory(
//								"Projects",
//								"/projects",
//								NewMenuItem("Comics", "/projects/comics"),
//								NewMenuItem("Music", "/projects/music"),
//								NewMenuItem("Podcasts", "/projects/podcasts"),
//								NewMenuItem("Video", "/projects/video"),
//								NewMenuItem("Software", "/projects/software"),
//							),
//						).HTML(),
//					),
//				),
//			)
//		}
//	}
//	if !self.FooterMenu.IsEmpty() {
//		body += Footer(
//			Columns(
//				Column(
//					12,
//					NewFooterMenu().HTML(),
//				),
//			),
//		)
//	}
//
//}

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
