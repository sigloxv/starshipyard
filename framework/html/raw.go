package html

import (
	"strings"

	entity "github.com/multiverse-os/starshipyard/framework/html/entity"
)

// TODO: Parsing HTML / Markdown into Go structs that filter out disallowed HTML
// based on an 'allow list' (or whitelist). Allow list filtering should happen
// both at parsing or Go object creation, but also at rendering. By doing this
// at two separate stages filtering occurs increasingly the likelyhood of
// avoiding XSS and other attacks.

////////////////////////////////////////////////////////////////////////////////////////////
// RawHTML string
// RawHTML is the component of Element that contains the raw HTML string
// that is parsed into the Go Element object.
type RawHTML string

func (self RawHTML) Attributes() (attributes map[string]string) {
	if len(self) == 0 {
		return nil
	}
	openTag := strings.Split(string(self)[(len(self.TagName())+2):], entity.GreaterThanSign.String())[0]
	rawAttributes := strings.Split(openTag, entity.Space.String())
	for _, rawAttribute := range rawAttributes {
		attributeComponents := strings.Split(rawAttribute, (entity.EqualSign.String() + entity.DoubleQuote.String()))
		attributes[attributeComponents[0]] = attributeComponents[1][:len(attributeComponents[1])] // subtracting the trailing `"`
	}
	return attributes
}

func (self RawHTML) TagName() string {
	if len(self) == 0 {
		return entity.EmptyString
	}
	return strings.Split(strings.Split(string(self)[1:], entity.GreaterThanSign.String())[0], entity.Space.String())[0]
}

func (self RawHTML) AttributeValue(name string) string {
	// NOTE: May be more efficient to loop ourselves so we can break if its found
	return strings.Split(strings.Split(string(self), (name + entity.EqualSign.String() + entity.DoubleQuote.String()))[1], entity.DoubleQuote.String())[0]
}

func (self RawHTML) HasAttribute(name string) bool {
	return strings.Contains(string(self), name)
}
