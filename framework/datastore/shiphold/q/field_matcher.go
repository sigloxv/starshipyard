package q

import (
	"go/token"
	"reflect"
)

type fieldMatcherDelegate struct {
	FieldMatcher
	Field string
}

func NewFieldMatcher(field string, fm FieldMatcher) Matcher {
	return fieldMatcherDelegate{Field: field, FieldMatcher: fm}
}

type FieldMatcher interface {
	MatchField(v interface{}) (bool, error)
}

func (self fieldMatcherDelegate) Match(i interface{}) (bool, error) {
	v := reflect.Indirect(reflect.ValueOf(i))
	return self.MatchValue(&v)
}

func (self fieldMatcherDelegate) MatchValue(v *reflect.Value) (bool, error) {
	field := v.FieldByName(self.Field)
	if !field.IsValid() {
		return false, errUnknownField
	}
	return self.MatchField(field.Interface())
}

func NewField2FieldMatcher(field1, field2 string, tok token.Token) Matcher {
	return field2fieldMatcherDelegate{Field1: field1, Field2: field2, Tok: tok}
}

type field2fieldMatcherDelegate struct {
	Field1, Field2 string
	Tok            token.Token
}

func (self field2fieldMatcherDelegate) Match(i interface{}) (bool, error) {
	v := reflect.Indirect(reflect.ValueOf(i))
	return self.MatchValue(&v)
}

func (self field2fieldMatcherDelegate) MatchValue(v *reflect.Value) (bool, error) {
	field1 := v.FieldByName(self.Field1)
	if !field1.IsValid() {
		return false, errUnknownField
	}
	field2 := v.FieldByName(self.Field2)
	if !field2.IsValid() {
		return false, errUnknownField
	}
	return compare(field1.Interface(), field2.Interface(), self.Tok), nil
}
