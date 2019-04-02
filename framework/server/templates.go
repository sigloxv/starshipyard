package server

import (
	html "github.com/multiverse-os/starshipyard/framework/html"
	template "github.com/multiverse-os/starshipyard/framework/html/template"
)

// TODO: Building out rails-like path/route conviences preferably generated
// from a route defintion
// We want a system that:
//  * we have a function to init routes, providing us with these variables or
//  this but in a map.
//  * paths should be lowercase
// * outpput a title which chops of first / or last item of / split
const (
	root_path         = "/"
	login_path        = "/login"
	collective_path   = "/collective"
	source_path       = "/source"
	about_path        = "/about"
	exchange_path     = "/exchange"
	transparency_path = "/transparency"
	voting_path       = "/voting"
	contact_path      = "/contact"
	public_key_path   = "/public_key"
)

// Socialscribe Example Specific
const (
	projects_path   = "/projects"
	casts_path      = "/casts"
	software_path   = "/software"
	streams_path    = "/streams"
	animations_path = "/animations"
	comics_path     = "/comics"
)

func (self *Server) LoadTemplates() {
	// TODO: Is this required with below?
	self.Templates = make(map[template.TemplateType]*template.Template)
	self.Templates[template.DefaultTemplate] = &template.Template{
		Title:       self.Config.AppName,
		CSS:         html.DefaultCSS(),
		Description: self.Config.Description,
		Keywords:    self.Config.Keywords,
		// TODO: just example, will be defined later via UI
	}

	self.Templates[template.ErrorTemplate] = &template.Template{
		Title:       self.Config.AppName,
		Description: self.Config.Description,
		Keywords:    self.Config.Keywords,
	}

}
