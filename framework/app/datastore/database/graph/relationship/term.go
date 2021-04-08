package relationship

type Term uint8

const (
	//////////////////////////////////////////////////////////////////////////////
	// Triple & Quad
	Subject Term = iota
	Predicate
	Object
	// Quad
	Graph
	//////////////////////////////////////////////////////////////////////////////
	// Generic
	CommonNoun       // Entity Type
	ProperNoun       // Entity
	TransitiveVerb   // Relationship Type
	IntransitiveVerb // Attribiute type
	Adjective        // Attribute for entry
	Adverb           // Attribute for relationship
)
