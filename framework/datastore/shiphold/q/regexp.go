package q

import (
	"fmt"
	"regexp"
	"sync"
)

func Re(field string, re string) Matcher {
	regexpCache.RLock()
	if r, ok := regexpCache.m[re]; ok {
		regexpCache.RUnlock()
		return NewFieldMatcher(field, &regexpMatcher{r: r})
	}
	regexpCache.RUnlock()
	regexpCache.Lock()
	r, err := regexp.Compile(re)
	if err == nil {
		regexpCache.m[re] = r
	}
	regexpCache.Unlock()
	return NewFieldMatcher(field, &regexpMatcher{r: r, err: err})
}

var regexpCache = struct {
	sync.RWMutex
	m map[string]*regexp.Regexp
}{m: make(map[string]*regexp.Regexp)}

type regexpMatcher struct {
	r   *regexp.Regexp
	err error
}

func (self *regexpMatcher) MatchField(v interface{}) (bool, error) {
	if self.err != nil {
		return false, self.err
	}
	switch fieldValue := v.(type) {
	case string:
		return self.r.MatchString(fieldValue), nil
	case []byte:
		return self.r.Match(fieldValue), nil
	default:
		return false, fmt.Errorf("only string and []byte supported for regexp matcher, got %T", fieldValue)
	}
}
