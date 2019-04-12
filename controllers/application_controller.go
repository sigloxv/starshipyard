package controllers

import (
	"fmt"
	"net/http"

	template "github.com/multiverse-os/starshipyard/framework/html/template"
	template "github.com/multiverse-os/starshipyard/views"
)

func Template(yield string) string {
	return ""
}

// TODO: In rails, we define the template in the controller. To acheive similar
// functionality, we will have functions that can be called from the controller
// to load and render desired template

func Root(w http.ResponseWriter, r *http.Request) {
	fmt.Println("root path")

	t := views.DefaultLayout("Home", "test")

	templateSection := template.Section{
		Columns: 2,
		Content: "test section",
	}

	t.Sections = append(t.Sections, templateSection)

	htmlString := t.HTML()

	w.Write([]byte(htmlString.String()))

}
