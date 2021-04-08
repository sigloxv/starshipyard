package framework

// NOTE: Concept: we want to be able to run multiple applications in a given
// instance. This would likely be defined by a ruby-like script config that
// defines what domains go where, reverse and inverting proxy settings, etc
type Domain struct {
	Name        string
	Subdomains  []string
	Certificate string
}
