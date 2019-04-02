package component

import (
	html "github.com/multiverse-os/starshipyard/html"
	template "github.com/multiverse-os/starshipyard/html/template"
	nav "github.com/multiverse-os/starshipyard/html/template/navigation"
)

func (self *Server) HeaderMenu() nav.NavigationMenu {
	return nav.NewHeaderMenu(
		nav.NewMenuOption(root_path, "Projects"),
		nav.NewMenuOption(collective_path, "Collective"),
		nav.NewMenuOption(source_path, "Source"),
		nav.NewMenuOption(contact_path, "Contact"),
		nav.NewMenuOption(about_path, "About"),
	)
}

func (self *Server) SidebarMenu() nav.NavigationMenu {
	return nav.NewLeftSideMenu(4,
		nav.NewMenuCategory(
			"Projects",
			projects_path,
			nav.NewMenuOption(software_path, "Software"),
			nav.NewMenuOption(casts_path, "Casts"),
			nav.NewMenuOption(streams_path, "Streams"),
			nav.NewMenuOption(animations_path, "Animations"),
			nav.NewMenuOption(comics_path, "Comics"),
		),
		nav.NewMenuCategory(
			"Collective",
			collective_path,
			nav.NewMenuOption(voting_path, "Voting"),
			nav.NewMenuOption(exchange_path, "Exchange"),
			nav.NewMenuOption(transparency_path, "Transparency"),
		),
		nav.NewMenuCategory(
			"About",
			about_path,
			nav.NewMenuOption(source_path, "Source"),
			nav.NewMenuOption(public_key_path, "Public Key"),
			nav.NewMenuOption(contact_path, "Contact"),
		),
	)
}

func FooterMenu() nav.NavigationMenu {
	return nav.NewFooterMenu(
		nav.NewMenuCategory(
			"Projects",
			projects_path,
			nav.NewMenuOption(software_path, "Software"),
			nav.NewMenuOption(casts_path, "Casts"),
			nav.NewMenuOption(streams_path, "Streams"),
			nav.NewMenuOption(animations_path, "Animations"),
			nav.NewMenuOption(comics_path, "Comics"),
		),
		nav.NewMenuCategory(
			"Collective",
			collective_path,
			nav.NewMenuOption(voting_path, "Voting"),
			nav.NewMenuOption(exchange_path, "Exchange"),
			nav.NewMenuOption(transparency_path, "Transparency"),
		),
		nav.NewMenuCategory(
			"About",
			about_path,
			nav.NewMenuOption(source_path, "Source"),
			nav.NewMenuOption(public_key_path, "Public Key"),
			nav.NewMenuOption(contact_path, "Contact"),
		),
	)
}
