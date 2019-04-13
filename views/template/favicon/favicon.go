package template

import (
	html "github.com/multiverse-os/starshipyard/framework/html"
)

func (self *Template) FaviconTag() html.Element {
	return BlankFavicon()
}

func BlankFavicon() html.Element {
	return html.A.Id("favicon").Relative("shortcut icon").Type("image/png").Href("data:image/png;base64,....==")
}
