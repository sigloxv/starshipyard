package shiphold

import (
	"reflect"
)

func (self *sorter) Swap(i, j int) {
	select {
	case <-self.done:
		return
	default:
	}
	if sortSink, ok := self.sink.(sliceSink); ok {
		reflect.Swapper(sortSink.slice().Interface())(i, j)
	} else {
		self.list[i], self.list[j] = self.list[j], self.list[i]
	}
}
