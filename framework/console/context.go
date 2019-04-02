package context

type Context struct {
	contextValues
	progressBar ProgressBar
	err         error
	Args        []string
	RawArgs     []string
	Cmd         Cmd
	Actions
}

func (self *Context) Err(err error) {
	self.err = err
}

func (self *Context) ProgressBar() ProgressBar {
	return self.progressBar
}

type contextValues map[string]interface{}

func (self contextValues) Get(key string) interface{} {
	return self[key]
}

func (self *contextValues) Set(key string, value interface{}) {
	if *self == nil {
		*self = make(map[string]interface{})
	}
	(*self)[key] = value
}

func (self contextValues) Del(key string) {
	delete(self, key)
}

func (self contextValues) Keys() (keys []string) {
	for key := range self {
		keys = append(keys, key)
	}
	return
}
