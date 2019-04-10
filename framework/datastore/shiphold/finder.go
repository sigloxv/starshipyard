package shiphold

import (
	"fmt"
	"reflect"

	index "github.com/multiverse-os/starshipyard/framework/datastore/shiphold/index"
	q "github.com/multiverse-os/starshipyard/framework/datastore/shiphold/q"

	bolt "go.etcd.io/bbolt"
)

// TODO: Review active record from Rails so we can make Starship Yard feel like
// rails in almost all aspects
type Finder interface {
	One(fieldName string, value interface{}, to interface{}) error
	Find(fieldName string, value interface{}, to interface{}, options ...func(q *index.Options)) error
	AllByIndex(fieldName string, to interface{}, options ...func(*index.Options)) error
	All(to interface{}, options ...func(*index.Options)) error
	Select(matchers ...q.Matcher) Query
	Range(fieldName string, min, max, to interface{}, options ...func(*index.Options)) error
	Prefix(fieldName string, prefix string, to interface{}, options ...func(*index.Options)) error
	Count(data interface{}) (int, error)
}

func (self *node) One(fieldName string, value interface{}, to interface{}) error {
	sink, err := newFirstSink(self, to)
	if err != nil {
		return err
	}
	bucketName := sink.bucketName()
	if bucketName == "" {
		return errNoName
	}
	if fieldName == "" {
		return errNotFound
	}
	ref := reflect.Indirect(sink.ref)
	config, err := extractSingleField(&ref, fieldName)
	if err != nil {
		return err
	}
	field, ok := config.Fields[fieldName]
	if !ok || (!field.IsID && field.Index == "") {
		query := newQuery(self, q.StrictEq(fieldName, value))
		query.Limit(1)
		if self.tx != nil {
			err = query.query(self.tx, sink)
		} else {
			err = self.s.Bolt.View(func(tx *bolt.Tx) error {
				return query.query(tx, sink)
			})
		}
		if err != nil {
			return err
		}
		return sink.flush()
	}
	val, err := toBytes(value, self.codec)
	if err != nil {
		return err
	}
	return self.readTx(func(tx *bolt.Tx) error {
		return self.one(tx, bucketName, fieldName, config, to, val, field.IsID)
	})
}

// TODO: Is there any good reason not to have this under One()?
func (self *node) one(tx *bolt.Tx, bucketName, fieldName string, config *structConfig, to interface{}, val []byte, skipIndex bool) error {
	bucket := self.GetBucket(tx, bucketName)
	if bucket == nil {
		return errNotFound
	}
	var id []byte
	if !skipIndex {
		idx, err := getIndex(bucket, config.Fields[fieldName].Index, fieldName)
		if err != nil {
			return err
		}

		id = idx.Get(val)
	} else {
		id = val
	}
	if id == nil {
		return errNotFound
	}
	raw := bucket.Get(id)
	if raw == nil {
		return errNotFound
	}
	return self.codec.Unmarshal(raw, to)
}

func (self *node) Find(fieldName string, value interface{}, to interface{}, options ...func(q *index.Options)) error {
	sink, err := newListSink(self, to)
	if err != nil {
		return err
	}
	bucketName := sink.bucketName()
	if bucketName == "" {
		return errNoName
	}
	ref := reflect.Indirect(reflect.New(sink.elemType))
	config, err := extractSingleField(&ref, fieldName)
	if err != nil {
		return err
	}
	opts := index.NewOptions()
	for _, fn := range options {
		fn(opts)
	}
	field, ok := config.Fields[fieldName]
	if !ok || (!field.IsID && (field.Index == "" || value == nil)) {
		query := newQuery(self, q.Eq(fieldName, value))
		query.Skip(opts.Skip).Limit(opts.Limit)
		if opts.Reverse {
			query.Reverse()
		}
		err = self.readTx(func(tx *bolt.Tx) error {
			return query.query(tx, sink)
		})
		if err != nil {
			return err
		}
		return sink.flush()
	}

	val, err := toBytes(value, self.codec)
	if err != nil {
		return err
	}

	return self.readTx(func(tx *bolt.Tx) error {
		return self.find(tx, bucketName, fieldName, config, sink, val, opts)
	})
}

func (self *node) find(tx *bolt.Tx, bucketName, fieldName string, config *structConfig, sink *listSink, val []byte, opts *index.Options) error {
	bucket := self.GetBucket(tx, bucketName)
	if bucket == nil {
		return errNotFound
	}
	idx, err := getIndex(bucket, config.Fields[fieldName].Index, fieldName)
	if err != nil {
		return err
	}

	list, err := idx.All(val, opts)
	if err != nil {
		return err
	}

	sink.results = reflect.MakeSlice(reflect.Indirect(sink.ref).Type(), len(list), len(list))

	sorter := newSorter(self, sink)
	for i := range list {
		raw := bucket.Get(list[i])
		if raw == nil {
			return errNotFound
		}

		if _, err := sorter.filter(nil, bucket, list[i], raw); err != nil {
			return err
		}
	}

	return sorter.flush()
}

// AllByIndex gets all the records of a bucket that are indexed in the specified index
func (self *node) AllByIndex(fieldName string, to interface{}, options ...func(*index.Options)) error {
	if fieldName == "" {
		return self.All(to, options...)
	}

	ref := reflect.ValueOf(to)

	if ref.Kind() != reflect.Ptr || ref.Elem().Kind() != reflect.Slice {
		return errSlicePtrNeeded
	}

	typ := reflect.Indirect(ref).Type().Elem()

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	newElem := reflect.New(typ)

	config, err := extract(&newElem)
	if err != nil {
		return err
	}

	if config.ID.Name == fieldName {
		return self.All(to, options...)
	}

	opts := index.NewOptions()
	for _, fn := range options {
		fn(opts)
	}

	return self.readTx(func(tx *bolt.Tx) error {
		return self.allByIndex(tx, fieldName, config, &ref, opts)
	})
}

func (self *node) allByIndex(tx *bolt.Tx, fieldName string, config *structConfig, ref *reflect.Value, opts *index.Options) error {
	bucket := self.GetBucket(tx, config.Name)
	if bucket == nil {
		return errNotFound
	}
	fieldConfig, ok := config.Fields[fieldName]
	if !ok {
		return errNotFound
	}

	idx, err := getIndex(bucket, fieldConfig.Index, fieldName)
	if err != nil {
		return err
	}

	list, err := idx.AllRecords(opts)
	if err != nil {
		return err
	}

	results := reflect.MakeSlice(reflect.Indirect(*ref).Type(), len(list), len(list))

	for i := range list {
		raw := bucket.Get(list[i])
		if raw == nil {
			return errNotFound
		}

		err = self.codec.Unmarshal(raw, results.Index(i).Addr().Interface())
		if err != nil {
			return err
		}
	}

	reflect.Indirect(*ref).Set(results)
	return nil
}

// All gets all the records of a bucket.
// If there are no records it returns no error and the 'to' parameter is set to an empty slice.
func (self *node) All(to interface{}, options ...func(*index.Options)) error {
	opts := index.NewOptions()
	for _, fn := range options {
		fn(opts)
	}

	query := newQuery(self, nil).Limit(opts.Limit).Skip(opts.Skip)
	if opts.Reverse {
		query.Reverse()
	}

	err := query.Find(to)
	if err != nil && err != errNotFound {
		return err
	}

	if err == errNotFound {
		ref := reflect.ValueOf(to)
		results := reflect.MakeSlice(reflect.Indirect(ref).Type(), 0, 0)
		reflect.Indirect(ref).Set(results)
	}
	return nil
}

// Range returns one or more records by the specified index within the specified range
func (self *node) Range(fieldName string, min, max, to interface{}, options ...func(*index.Options)) error {
	sink, err := newListSink(self, to)
	if err != nil {
		return err
	}

	bucketName := sink.bucketName()
	if bucketName == "" {
		return errNoName
	}

	ref := reflect.Indirect(reflect.New(sink.elemType))
	config, err := extractSingleField(&ref, fieldName)
	if err != nil {
		return err
	}

	opts := index.NewOptions()
	for _, fn := range options {
		fn(opts)
	}

	field, ok := config.Fields[fieldName]
	if !ok || (!field.IsID && field.Index == "") {
		query := newQuery(self, q.And(q.Gte(fieldName, min), q.Lte(fieldName, max)))
		query.Skip(opts.Skip).Limit(opts.Limit)

		if opts.Reverse {
			query.Reverse()
		}

		err = self.readTx(func(tx *bolt.Tx) error {
			return query.query(tx, sink)
		})

		if err != nil {
			return err
		}

		return sink.flush()
	}

	mn, err := toBytes(min, self.codec)
	if err != nil {
		return err
	}

	mx, err := toBytes(max, self.codec)
	if err != nil {
		return err
	}

	return self.readTx(func(tx *bolt.Tx) error {
		return self.rnge(tx, bucketName, fieldName, config, sink, mn, mx, opts)
	})
}

func (self *node) rnge(tx *bolt.Tx, bucketName, fieldName string, config *structConfig, sink *listSink, min, max []byte, opts *index.Options) error {
	bucket := self.GetBucket(tx, bucketName)
	if bucket == nil {
		reflect.Indirect(sink.ref).SetLen(0)
		return nil
	}

	idx, err := getIndex(bucket, config.Fields[fieldName].Index, fieldName)
	if err != nil {
		return err
	}

	list, err := idx.Range(min, max, opts)
	if err != nil {
		return err
	}

	sink.results = reflect.MakeSlice(reflect.Indirect(sink.ref).Type(), len(list), len(list))
	sorter := newSorter(self, sink)
	for i := range list {
		raw := bucket.Get(list[i])
		if raw == nil {
			return errNotFound
		}

		if _, err := sorter.filter(nil, bucket, list[i], raw); err != nil {
			return err
		}
	}

	return sorter.flush()
}

func (self *node) Prefix(fieldName string, prefix string, to interface{}, options ...func(*index.Options)) error {
	sink, err := newListSink(self, to)
	if err != nil {
		return err
	}

	bucketName := sink.bucketName()
	if bucketName == "" {
		return errNoName
	}

	ref := reflect.Indirect(reflect.New(sink.elemType))
	config, err := extractSingleField(&ref, fieldName)
	if err != nil {
		return err
	}

	opts := index.NewOptions()
	for _, fn := range options {
		fn(opts)
	}

	field, ok := config.Fields[fieldName]
	if !ok || (!field.IsID && field.Index == "") {
		query := newQuery(self, q.Re(fieldName, fmt.Sprintf("^%s", prefix)))
		query.Skip(opts.Skip).Limit(opts.Limit)

		if opts.Reverse {
			query.Reverse()
		}

		err = self.readTx(func(tx *bolt.Tx) error {
			return query.query(tx, sink)
		})

		if err != nil {
			return err
		}

		return sink.flush()
	}

	prfx, err := toBytes(prefix, self.codec)
	if err != nil {
		return err
	}

	return self.readTx(func(tx *bolt.Tx) error {
		return self.prefix(tx, bucketName, fieldName, config, sink, prfx, opts)
	})
}

func (self *node) prefix(tx *bolt.Tx, bucketName, fieldName string, config *structConfig, sink *listSink, prefix []byte, opts *index.Options) error {
	bucket := self.GetBucket(tx, bucketName)
	if bucket == nil {
		reflect.Indirect(sink.ref).SetLen(0)
		return nil
	}
	idx, err := getIndex(bucket, config.Fields[fieldName].Index, fieldName)
	if err != nil {
		return err
	}
	list, err := idx.Prefix(prefix, opts)
	if err != nil {
		return err
	}
	sink.results = reflect.MakeSlice(reflect.Indirect(sink.ref).Type(), len(list), len(list))
	sorter := newSorter(self, sink)
	for i := range list {
		raw := bucket.Get(list[i])
		if raw == nil {
			return errNotFound
		}
		if _, err := sorter.filter(nil, bucket, list[i], raw); err != nil {
			return err
		}
	}
	return sorter.flush()
}

// TODO: We should cache this and just update it as we chance things in the
// database. This will save resources and let us do cheap count lookups
func (self *node) Count(data interface{}) (int, error) {
	return self.Select().Count(data)
}
