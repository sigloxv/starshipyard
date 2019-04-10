package q

import (
	"go/token"
	"reflect"
)

type Matcher interface {
	Match(interface{}) (bool, error)
}

type ValueMatcher interface {
	MatchValue(*reflect.Value) (bool, error)
}

type cmp struct {
	value interface{}
	token token.Token
}

func (self *cmp) MatchField(v interface{}) (bool, error) {
	return compare(v, self.value, self.token), nil
}

type trueMatcher struct{}

func (*trueMatcher) Match(i interface{}) (bool, error) {
	return true, nil
}

func (*trueMatcher) MatchValue(v *reflect.Value) (bool, error) {
	return true, nil
}

type or struct {
	children []Matcher
}

func (self *or) Match(i interface{}) (bool, error) {
	v := reflect.Indirect(reflect.ValueOf(i))
	return self.MatchValue(&v)
}

func (self *or) MatchValue(v *reflect.Value) (bool, error) {
	for _, matcher := range self.children {
		if vm, ok := matcher.(ValueMatcher); ok {
			ok, err := vm.MatchValue(v)
			if err != nil {
				return false, err
			}
			if ok {
				return true, nil
			}
			continue
		}
		ok, err := matcher.Match(v.Interface())
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
	}
	return false, nil
}

type and struct {
	children []Matcher
}

func (self *and) Match(i interface{}) (bool, error) {
	v := reflect.Indirect(reflect.ValueOf(i))
	return self.MatchValue(&v)
}

func (self *and) MatchValue(v *reflect.Value) (bool, error) {
	for _, matcher := range self.children {
		if vm, ok := matcher.(ValueMatcher); ok {
			ok, err := vm.MatchValue(v)
			if err != nil {
				return false, err
			}
			if !ok {
				return false, nil
			}
			continue
		}
		ok, err := matcher.Match(v.Interface())
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}
	}
	return true, nil
}

type strictEq struct {
	field string
	value interface{}
}

func (self *strictEq) MatchField(v interface{}) (bool, error) {
	return reflect.DeepEqual(v, self.value), nil
}

type in struct {
	list interface{}
}

func (self *in) MatchField(v interface{}) (bool, error) {
	ref := reflect.ValueOf(self.list)
	if ref.Kind() != reflect.Slice {
		return false, nil
	}

	c := cmp{
		token: token.EQL,
	}

	for i := 0; i < ref.Len(); i++ {
		c.value = ref.Index(i).Interface()
		ok, err := c.MatchField(v)
		if err != nil {
			return false, err
		}
		if ok {
			return true, nil
		}
	}

	return false, nil
}

type not struct {
	children []Matcher
}

func (self *not) Match(i interface{}) (bool, error) {
	v := reflect.Indirect(reflect.ValueOf(i))
	return self.MatchValue(&v)
}

func (self *not) MatchValue(v *reflect.Value) (bool, error) {
	var err error
	for _, matcher := range self.children {
		vm, ok := matcher.(ValueMatcher)
		if ok {
			ok, err = vm.MatchValue(v)
		} else {
			ok, err = matcher.Match(v.Interface())
		}
		if err != nil {
			return false, err
		}
		if ok {
			return false, nil
		}
	}
	return true, nil
}

func True() Matcher                   { return &trueMatcher{} }
func Or(matchers ...Matcher) Matcher  { return &or{children: matchers} }
func And(matchers ...Matcher) Matcher { return &and{children: matchers} }
func Not(matchers ...Matcher) Matcher { return &not{children: matchers} }

func Eq(field string, v interface{}) Matcher {
	return NewFieldMatcher(field, &cmp{value: v, token: token.EQL})
}

func EqF(field1, field2 string) Matcher {
	return NewField2FieldMatcher(field1, field2, token.EQL)
}

func StrictEq(field string, v interface{}) Matcher {
	return NewFieldMatcher(field, &strictEq{value: v})
}

func Gt(field string, v interface{}) Matcher {
	return NewFieldMatcher(field, &cmp{value: v, token: token.GTR})
}

func GtF(field1, field2 string) Matcher {
	return NewField2FieldMatcher(field1, field2, token.GTR)
}

func Gte(field string, v interface{}) Matcher {
	return NewFieldMatcher(field, &cmp{value: v, token: token.GEQ})
}

func GteF(field1, field2 string) Matcher {
	return NewField2FieldMatcher(field1, field2, token.GEQ)
}

func Lt(field string, v interface{}) Matcher {
	return NewFieldMatcher(field, &cmp{value: v, token: token.LSS})
}

func LtF(field1, field2 string) Matcher {
	return NewField2FieldMatcher(field1, field2, token.LSS)
}

func Lte(field string, v interface{}) Matcher {
	return NewFieldMatcher(field, &cmp{value: v, token: token.LEQ})
}

func LteF(field1, field2 string) Matcher {
	return NewField2FieldMatcher(field1, field2, token.LEQ)
}

func In(field string, v interface{}) Matcher {
	return NewFieldMatcher(field, &in{list: v})
}
