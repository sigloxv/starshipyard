package template

import (
	"net/http"
	"strings"

	html "github.com/multiverse-os/starshipyard/framework/html"
	entity "github.com/multiverse-os/starshipyard/framework/html/entity"
	flashmessages "github.com/multiverse-os/starshipyard/framework/html/template/flashmessages"
	navigation "github.com/multiverse-os/starshipyard/framework/html/template/navigation"
)

// TODO: Should make Type based status codes, so we dont have to do
// string comparison. our serve should have a fixed number of a vailable
// allowed status codes that way we can control the presentation
type TemplateType int

const (
	DefaultTemplate TemplateType = iota
	AdminTemplate
	ModeratorTemplate
	ErrorTemplate
)

func (self *Template) HTMLAsBytes() []byte {
	return []byte(self.HTML().String())
}

////////////////////////////////////////////////////////////////////////////////
// Favicon
func (self *Template) FaviconTag() html.Element {
	return BlankFavicon()
}

func BlankFavicon() html.Element {
	return html.A.Id("favicon").Relative("shortcut icon").Type("image/png").Href("data:image/png;base64,....==")
}

// TODO: http.Request is generic since its using the net/http package. but
// ideally even then it should be broken off the template object as much as
// possible while still incorporating the concept of

// TODO: Default support for transparent background "overlay" style templates
// designed for overlaying streams or as HUDs or UIs over 3d engines.

// TODO: Or maybe we should store content in rows (divided by <div
// class='columns'></div>)
type Section struct {
	Columns int
	Content string
}

// TODO: It would be nice to use this in html in parsing, seeing if we can
// locate the various aspects of the parsed HTML to simplify automation (like
// mechanize) of web sites.
type Template struct {
	Type          TemplateType
	Request       *http.Request
	FlashMessages flashmessages.FlashMessages

	//Charset     string
	Title       string
	Description string
	Keywords    []string
	CSS         html.CSS
	HeaderMenu  navigation.Menu
	SidebarMenu navigation.Menu
	FooterMenu  navigation.Menu
	Content     string
	Sections    []Section
}

func (self *Template) Render(request *http.Request, content string) *Template {
	self.Request = request
	self.Content = content
	return self
}

func (self *Template) HTML() html.Element {
	return html.HTML.ChildElements(self.HeadTag(), self.BodyTag())
}

func (self *Template) HeadTag() html.Element {
	return html.Head.ChildElements(
		html.Meta.Charset("UTF-8"),
		html.Meta.Name("description").Content(self.Description),
		html.Meta.Name("keywords").Content(strings.Join(self.Keywords, entity.EmptyString)),
		html.Meta.Name("viewport").Content("width=device-width, initial-scale=1.0"),
		html.Title.Text(self.Title),
		html.Style.Text(self.CSS.Framework+self.CSS.Overrides),
	)
}

func (self *Template) BodyTag() html.Element {
	return html.Body.ChildElements(
		html.Header.Text(""),
		html.Content.Text(""),
		html.Body.Text("test"),
	)
	//Columns(
	//	// TODO: Not super excited about this current API
	//	Column(2, H1(self.Title)),
	//	Column(6, self.HeaderMenu.HTML()),
	//	Column(4, LoginForm("/login")),
	//),
	//Columns(
	//	Column(12, self.FlashMessages.HTML()),
	//),
	//if self.SidebarMenu.IsEmpty() {
	//	body += Content(
	//		Columns(
	//			Column(12, self.Content),
	//		),
	//	)
	//} else {
	//	// TODO: Use columns to add sidbar, sidebar menu should define the side
	//	// and then build the template accordignly
	//	// steve ducey cooking show goes here too
	//	switch self.SidebarMenu.Type {
	//	case LeftSideMenu:
	//		body += Content(
	//			Columns(
	//				/// TODO: This is just for dev, all these values will be defined by
	//				// administrators via backend (hehe backend).
	//				// TODO: I really don't like this initialization method, feels very
	//				// clunky. This column initialization could be better too since we are
	//				// telling it twice it is 4 already.

	//				Column(
	//					2,
	//					NewLeftSideMenu(
	//						2,
	//						NewMenuCategory(
	//							"/projects",
	//							"Projects",
	//							NewMenuItem("Comics", "/projects/comics"),
	//							NewMenuItem("Music", "/projects/music"),
	//							NewMenuItem("Podcasts", "/projects/podcasts"),
	//							NewMenuItem("Video", "/projects/video"),
	//							NewMenuItem("Software", "/projects/software"),
	//						),
	//					).HTML()),

	//				Column(10, Div(Class("main", "tile"), self.Content)),
	//			),
	//		)
	//	case RightSideMenu:
	//		body += "right-side-menu"
	//		body += Content(
	//			Columns(
	//				Column(10, self.Content),
	//				// TODO: Seems rather silly to indicate back to back that this is 4
	//				// wide. Must be a better way
	//				Column(
	//					2,
	//					NewRightSideMenu(
	//						2,
	//						NewMenuCategory(
	//							"Projects",
	//							"/projects",
	//							NewMenuItem("Comics", "/projects/comics"),
	//							NewMenuItem("Music", "/projects/music"),
	//							NewMenuItem("Podcasts", "/projects/podcasts"),
	//							NewMenuItem("Video", "/projects/video"),
	//							NewMenuItem("Software", "/projects/software"),
	//						),
	//					).HTML(),
	//				),
	//			),
	//		)
	//	}
	//}
	//if !self.FooterMenu.IsEmpty() {
	//	body += Footer(
	//		Columns(
	//			Column(
	//				12,
	//				NewFooterMenu().HTML(),
	//			),
	//		),
	//	)
	//}

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
//}
