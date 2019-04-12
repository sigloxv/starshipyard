package views

import (
	"fmt"

	css "github.com/multiverse-os/starshipyard/views/css"

	template "github.com/multiverse-os/starshipyard/framework/html/template"
)

func DefaultLayout(title, yield string) string {
	return template.Template{
		Type:        template.DefaultTemplate,
		Title:       ("Title prfeix - " + title),
		Description: "A fitting description",
		Keywords:    []string{"webapp", "framework"},
		CSS:         css.DefaultCSS(),
		Content:     yield,
	}
}

// TODO: Views should use this default template, and extend it for each page by
// adding sections.
