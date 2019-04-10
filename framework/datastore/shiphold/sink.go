package shiphold

import (
	"reflect"
	"sort"

	q "github.com/multiverse-os/starshipyard/framework/datastore/shiphold/q"

	bolt "go.etcd.io/bbolt"
)

// TODO: Discovering A LOT of repeated code, not dry at all, also tons of
// unimplemented functions

// TODO: Change the key to our ID system
// TODO: Lets abandon the reflect stuff, just use cbor to encode into bytes and
// back into structs. Or if we keep it, then lets use a minimal
type item struct {
	value  *reflect.Value
	bucket *bolt.Bucket
	k      []byte
	v      []byte
}

func newSorter(node Node, sortSink sink) *sorter {
	return &sorter{
		node:  node,
		sink:  sortSink,
		skip:  0,
		limit: -1,
		list:  make([]*item, 0),
		err:   make(chan error),
		done:  make(chan struct{}),
	}
}

type sorter struct {
	node    Node
	sink    sink
	list    []*item
	skip    int
	limit   int
	orderBy []string
	reverse bool
	err     chan error
	done    chan struct{}
}

func (self *sorter) filter(tree q.Matcher, bucket *bolt.Bucket, k, v []byte) (bool, error) {
	i := &item{
		bucket: bucket,
		k:      k,
		v:      v,
	}
	reflectSink, ok := self.sink.(reflectSink)
	if !ok {
		return self.add(i)
	}
	element := reflectSink.elem()
	if err := self.node.Codec().Unmarshal(v, element.Interface()); err != nil {
		return false, err
	}
	i.value = &element
	if tree != nil {
		ok, err := tree.Match(element.Interface())
		if err != nil {
			return false, err
		}
		if !ok {
			return false, nil
		}
	}
	if len(self.orderBy) == 0 {
		return self.add(i)
	}
	if _, ok := self.sink.(sliceSink); ok {
		// add directly to sink, we'll apply skip/limits after sorting
		return false, self.sink.add(i)
	}
	self.list = append(self.list, i)
	return false, nil
}

func (self *sorter) add(itm *item) (stop bool, err error) {
	if self.limit == 0 {
		return true, nil
	}
	if self.skip > 0 {
		self.skip--
		return false, nil
	}
	if self.limit > 0 {
		self.limit--
	}
	err = self.sink.add(itm)
	return self.limit == 0, err
}

func (self *sorter) compareValue(left reflect.Value, right reflect.Value) int {
	if !left.IsValid() || !right.IsValid() {
		if left.IsValid() {
			return 1
		}
		return -1
	}
	switch left.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		l, r := left.Int(), right.Int()
		if l < r {
			return -1
		}
		if l > r {
			return 1
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		l, r := left.Uint(), right.Uint()
		if l < r {
			return -1
		}
		if l > r {
			return 1
		}
	case reflect.Float32, reflect.Float64:
		l, r := left.Float(), right.Float()
		if l < r {
			return -1
		}
		if l > r {
			return 1
		}
	case reflect.String:
		l, r := left.String(), right.String()
		if l < r {
			return -1
		}
		if l > r {
			return 1
		}
	default:
		rawLeft, err := toBytes(left.Interface(), self.node.Codec())
		if err != nil {
			return -1
		}
		rawRight, err := toBytes(right.Interface(), self.node.Codec())
		if err != nil {
			return 1
		}
		l, r := string(rawLeft), string(rawRight)
		if l < r {
			return -1
		}
		if l > r {
			return 1
		}
	}
	return 0
}

func (self *sorter) less(leftElem reflect.Value, rightElem reflect.Value) bool {
	for _, orderBy := range self.orderBy {
		leftField := reflect.Indirect(leftElem).FieldByName(orderBy)
		if !leftField.IsValid() {
			self.err <- errNotFound
			return false
		}
		rightField := reflect.Indirect(rightElem).FieldByName(orderBy)
		if !rightField.IsValid() {
			self.err <- errNotFound
			return false
		}
		direction := 1
		if self.reverse {
			direction = -1
		}
		switch self.compareValue(leftField, rightField) * direction {
		case -1:
			return true
		case 1:
			return false
		default:
			continue
		}
	}
	return false
}

func (self *sorter) flush() error {
	if len(self.orderBy) == 0 {
		return self.sink.flush()
	}
	go func() {
		sort.Sort(self)
		close(self.err)
	}()
	err := <-self.err
	close(self.done)
	if err != nil {
		return err
	}
	if sink, ok := self.sink.(sliceSink); ok {
		if !sink.slice().IsValid() {
			return self.sink.flush()
		}
		if self.skip >= sink.slice().Len() {
			sink.reset()
			return self.sink.flush()
		}
		leftBound := self.skip
		if leftBound < 0 {
			leftBound = 0
		}
		limit := self.limit
		if self.limit < 0 {
			limit = 0
		}
		rightBound := leftBound + limit
		if rightBound > sink.slice().Len() || rightBound == leftBound {
			rightBound = sink.slice().Len()
		}
		sink.setSlice(sink.slice().Slice(leftBound, rightBound))
		return self.sink.flush()
	}
	for _, i := range self.list {
		if i == nil {
			break
		}
		stop, err := self.add(i)
		if err != nil {
			return err
		}
		if stop {
			break
		}
	}
	return self.sink.flush()
}

func (self *sorter) Len() int {
	select {
	case <-self.done:
		return 0
	default:
	}
	if sink, ok := self.sink.(sliceSink); ok {
		return sink.slice().Len()
	}
	return len(self.list)

}

func (self *sorter) Less(i, j int) bool {
	select {
	case <-self.done:
		return false
	default:
	}
	if sink, ok := self.sink.(sliceSink); ok {
		return self.less(sink.slice().Index(i), sink.slice().Index(j))
	}
	return self.less(*self.list[i].value, *self.list[j].value)
}

type sink interface {
	bucketName() string
	flush() error
	add(*item) error
	readOnly() bool
}

type reflectSink interface {
	elem() reflect.Value
}

type sliceSink interface {
	slice() reflect.Value
	setSlice(reflect.Value)
	reset()
}

func newListSink(node Node, to interface{}) (*listSink, error) {
	ref := reflect.ValueOf(to)

	if ref.Kind() != reflect.Ptr || reflect.Indirect(ref).Kind() != reflect.Slice {
		return nil, errSlicePtrNeeded
	}

	sliceType := reflect.Indirect(ref).Type()
	elemType := sliceType.Elem()

	if elemType.Kind() == reflect.Ptr {
		elemType = elemType.Elem()
	}

	if elemType.Name() == "" {
		return nil, errNoName
	}
	return &listSink{
		node:     node,
		ref:      ref,
		isPtr:    sliceType.Elem().Kind() == reflect.Ptr,
		elemType: elemType,
		name:     elemType.Name(),
		results:  reflect.MakeSlice(reflect.Indirect(ref).Type(), 0, 0),
	}, nil
}

type listSink struct {
	node     Node
	ref      reflect.Value
	results  reflect.Value
	elemType reflect.Type
	name     string
	isPtr    bool
	idx      int
}

func (self *listSink) slice() reflect.Value {
	return self.results
}

func (self *listSink) setSlice(s reflect.Value) {
	self.results = s
}

func (self *listSink) reset() {
	self.results = reflect.MakeSlice(reflect.Indirect(self.ref).Type(), 0, 0)
}

func (self *listSink) elem() reflect.Value {
	if self.results.IsValid() && self.idx < self.results.Len() {
		return self.results.Index(self.idx).Addr()
	}
	return reflect.New(self.elemType)
}

func (self *listSink) bucketName() string {
	return self.name
}

func (self *listSink) add(i *item) error {
	if self.idx == self.results.Len() {
		if self.isPtr {
			self.results = reflect.Append(self.results, *i.value)
		} else {
			self.results = reflect.Append(self.results, reflect.Indirect(*i.value))
		}
	}
	self.idx++
	return nil
}

func (self *listSink) flush() error {
	if self.results.IsValid() && self.results.Len() > 0 {
		reflect.Indirect(self.ref).Set(self.results)
		return nil
	}
	return errNotFound
}

func (self *listSink) readOnly() bool {
	return true
}

func newFirstSink(node Node, to interface{}) (*firstSink, error) {
	ref := reflect.ValueOf(to)

	if !ref.IsValid() || ref.Kind() != reflect.Ptr || ref.Elem().Kind() != reflect.Struct {
		return nil, errStructPtrNeeded
	}
	return &firstSink{
		node: node,
		ref:  ref,
	}, nil
}

type firstSink struct {
	node  Node
	ref   reflect.Value
	found bool
}

func (self *firstSink) elem() reflect.Value {
	return reflect.New(reflect.Indirect(self.ref).Type())
}

func (self *firstSink) bucketName() string {
	return reflect.Indirect(self.ref).Type().Name()
}

func (self *firstSink) add(i *item) error {
	reflect.Indirect(self.ref).Set(i.value.Elem())
	self.found = true
	return nil
}

func (self *firstSink) flush() error {
	if !self.found {
		return errNotFound
	}
	return nil
}

func (self *firstSink) readOnly() bool {
	return true
}

func newDeleteSink(node Node, kind interface{}) (*deleteSink, error) {
	ref := reflect.ValueOf(kind)

	if !ref.IsValid() || ref.Kind() != reflect.Ptr || ref.Elem().Kind() != reflect.Struct {
		return nil, errStructPtrNeeded
	}
	return &deleteSink{
		node: node,
		ref:  ref,
	}, nil
}

type deleteSink struct {
	node    Node
	ref     reflect.Value
	removed int
}

func (self *deleteSink) elem() reflect.Value {
	return reflect.New(reflect.Indirect(self.ref).Type())
}

func (self *deleteSink) bucketName() string {
	return reflect.Indirect(self.ref).Type().Name()
}

func (self *deleteSink) add(i *item) error {
	info, err := extract(&self.ref)
	if err != nil {
		return err
	}
	for fieldName, fieldConfig := range info.Fields {
		if fieldConfig.Index == "" {
			continue
		}
		idx, err := getIndex(i.bucket, fieldConfig.Index, fieldName)
		if err != nil {
			return err
		}
		err = idx.RemoveID(i.k)
		if err != nil {
			return err
		}
	}
	self.removed++
	return i.bucket.Delete(i.k)
}

func (self *deleteSink) flush() error {
	if self.removed == 0 {
		return errNotFound
	}
	return nil
}

// TODO: Implement
func (self *deleteSink) readOnly() bool {
	return false
}

func newCountSink(node Node, kind interface{}) (*countSink, error) {
	ref := reflect.ValueOf(kind)
	if !ref.IsValid() || ref.Kind() != reflect.Ptr || ref.Elem().Kind() != reflect.Struct {
		return nil, errStructPtrNeeded
	}
	return &countSink{
		node: node,
		ref:  ref,
	}, nil
}

type countSink struct {
	node    Node
	ref     reflect.Value
	counter int
}

func (self *countSink) elem() reflect.Value {
	return reflect.New(reflect.Indirect(self.ref).Type())
}

func (self *countSink) bucketName() string {
	return reflect.Indirect(self.ref).Type().Name()
}

func (self *countSink) add(i *item) error {
	self.counter++
	return nil
}

// TODO: Implement
func (self *countSink) flush() error {
	return nil
}

// TODO: Implement
func (self *countSink) readOnly() bool {
	return true
}

func newRawSink() *rawSink {
	return &rawSink{}
}

type rawSink struct {
	results [][]byte
	execFn  func([]byte, []byte) error
}

func (self *rawSink) add(i *item) error {
	if self.execFn != nil {
		err := self.execFn(i.k, i.v)
		if err != nil {
			return err
		}
	} else {
		self.results = append(self.results, i.v)
	}
	return nil
}

// TODO: Implement readonly
func (self *rawSink) bucketName() string {
	return ""
}

// TODO: Implement readonly
func (self *rawSink) flush() error {
	return nil
}

// TODO: Implement readonly
func (self *rawSink) readOnly() bool {
	return true
}

func newEachSink(to interface{}) (*eachSink, error) {
	ref := reflect.ValueOf(to)

	if !ref.IsValid() || ref.Kind() != reflect.Ptr || ref.Elem().Kind() != reflect.Struct {
		return nil, errStructPtrNeeded
	}

	return &eachSink{
		ref: ref,
	}, nil
}

type eachSink struct {
	ref    reflect.Value
	execFn func(interface{}) error
}

func (self *eachSink) elem() reflect.Value {
	return reflect.New(reflect.Indirect(self.ref).Type())
}

func (self *eachSink) bucketName() string {
	return reflect.Indirect(self.ref).Type().Name()
}

func (self *eachSink) add(i *item) error {
	return self.execFn(i.value.Interface())
}

// TODO: Implement this
func (self *eachSink) flush() error {
	return nil
}

// TODO: Implement this
func (self *eachSink) readOnly() bool {
	return true
}
