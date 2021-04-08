package graph

type RDF uint8 // Resource Description Framework

const (
	QuadStore RDF = iota
	TripleStore
	CustomStore
)
