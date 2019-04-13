package html

import (
	"errors"
	"fmt"
	"strings"

	attribute "github.com/multiverse-os/starshipyard/framework/html/attribute"
	entity "github.com/multiverse-os/starshipyard/framework/html/entity"
)

// [SECURITY] TODO: Every link that goes outside should have target="_no_rel" to prevent certian types of XSS

// Anatomy Of HTML Element
/////////////////////////////////////////////////////////////////////////////////////////////////
// [Mozilla HTML](https://developer.mozilla.org/en-US/docs/Learn/HTML/Introduction_to_HTML/Getting_started)
// https://developer.mozilla.org/en-US/docs/Web/HTML/Element
//
// [Opening Tag]                        [Closing Tag]
//       \                                   /
//       <p> Paragraph text between tags </p>
//           |__________________________|
//                       |
//                    [content]
//  |____________________________________________|
//                       |
//      [All of it together is an HTML element]
type Elements []Element

type Element struct {
	Tag        Tag
	Block      bool
	Attributes map[attribute.Key][]string
	raw        RawHTML
	content    interface{}
}

func NewElement(tag Tag) Element {
	return Element{
		Tag:        tag,
		Attributes: make(map[attribute.Key][]string),
	}
}

// HTML parsing into Go structs (incomplete)
///////////////////////////////////////////////////////////////////////////////
func ParseElement(html string) (element Element, err error) {
	if len(html) == 0 {
		return element, errors.New("error: empty input string")
	}
	element.Tag = MarshalTag(strings.Split(strings.Split(html[1:], entity.GreaterThanSign.String())[0], entity.Space.String())[0])
	return element, nil
}

// Content Helpers
///////////////////////////////////////////////////////////////////////////////
func (self Element) Text(content string) Element {
	self.content = string(content)
	return self
}

func (self Element) ChildElement(element Element) Element {
	self.content = element
	return self
}

func (self Element) ChildElements(elements ...Element) Element {
	self.content = elements
	return self
}

// String output
///////////////////////////////////////////////////////////////////////////////
func (self Element) String() (output string) {
	if len(self.Tag.String()) == 0 {
		fmt.Println("tag is nil!")
		//return ""
	}
	fmt.Println("How many attributes?", len(self.Attributes))
	output += self.Tag.Open(self.Attributes)
	fmt.Println("output after first line:", output)

	switch content := self.content.(type) {
	case string:
		if len(content) > 0 {
			output += content
		}
	case []Element:
		for _, element := range content {
			output += element.String()
		}
	case Element:
		output += content.String()
	}
	output += self.Tag.Close()

	fmt.Println("output tag:", self.Tag.String())

	return output
}
