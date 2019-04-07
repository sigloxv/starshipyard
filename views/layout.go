package main

import (
	"fmt"

	template "github.com/multiverse-os/starshipyard/framework/html/template"
)

func Layouts() {
	// TODO: Is this required with below?
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

func DefaultLayout(yield string) string {
	fmt.Println("default layout will be the general structure that will then")
	fmt.Println("have the specific view rendered inside of in a yield like")
	fmt.Println("yeild for now is a string:", yield)
	fmt.Println("nesting using closures (nesting functions")
}
