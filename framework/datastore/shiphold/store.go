package shiphold

import (
	"bytes"
	"reflect"

	index "github.com/multiverse-os/starshipyard/framework/datastore/shiphold/index"
	q "github.com/multiverse-os/starshipyard/framework/datastore/shiphold/q"

	bolt "go.etcd.io/bbolt"
)

type TypeStore interface {
	Finder
	Init(data interface{}) error
	ReIndex(data interface{}) error
	Save(data interface{}) error
	Update(data interface{}) error
	UpdateField(data interface{}, fieldName string, value interface{}) error
	Drop(data interface{}) error
	DeleteStruct(data interface{}) error
}

func (self *node) Init(data interface{}) error {
	v := reflect.ValueOf(data)
	config, err := extract(&v)
	if err != nil {
		return err
	}
	return self.readWriteTx(func(tx *bolt.Tx) error {
		return self.init(tx, config)
	})
}

func (self *node) init(tx *bolt.Tx, config *structConfig) error {
	bucket, err := self.CreateBucketIfNotExists(tx, config.Name)
	if err != nil {
		return err
	}
	if _, err = newMeta(bucket, self); err != nil {
		return err
	}
	for fieldName, fieldConfig := range config.Fields {
		if fieldConfig.Index == "" {
			continue
		}
		switch fieldConfig.Index {
		case tagUniqueIdx:
			_, err = index.NewUniqueIndex(bucket, []byte(indexPrefix+fieldName))
		case tagIdx:
			_, err = index.NewListIndex(bucket, []byte(indexPrefix+fieldName))
		default:
			err = errIdxNotFound
		}
	}
	return err
}

func (self *node) ReIndex(data interface{}) error {
	ref := reflect.ValueOf(data)

	if !ref.IsValid() || ref.Kind() != reflect.Ptr || ref.Elem().Kind() != reflect.Struct {
		return errStructPtrNeeded
	}
	config, err := extract(&ref)
	if err != nil {
		return err
	}
	return self.readWriteTx(func(tx *bolt.Tx) error {
		return self.reIndex(tx, data, config)
	})
}

func (self *node) reIndex(tx *bolt.Tx, data interface{}, config *structConfig) error {
	root := self.WithTransaction(tx)
	nodes := root.From(config.Name).PrefixScan(indexPrefix)
	bucket := root.GetBucket(tx, config.Name)
	if bucket == nil {
		return errNotFound
	}
	for _, node := range nodes {
		buckets := node.Bucket()
		name := buckets[len(buckets)-1]
		err := bucket.DeleteBucket([]byte(name))
		if err != nil {
			return err
		}
	}
	total, err := root.Count(data)
	if err != nil {
		return err
	}
	for i := 0; i < total; i++ {
		err = root.Select(q.True()).Skip(i).First(data)
		if err != nil {
			return err
		}
		err = root.Update(data)
		if err != nil {
			return err
		}
	}

	return nil
}

func (self *node) Save(data interface{}) error {
	ref := reflect.ValueOf(data)

	if !ref.IsValid() || ref.Kind() != reflect.Ptr || ref.Elem().Kind() != reflect.Struct {
		return errStructPtrNeeded
	}
	config, err := extract(&ref)
	if err != nil {
		return err
	}
	if config.ID.IsZero {
		if !config.ID.IsInteger || !config.ID.Increment {
			return errZeroID
		}
	}
	return self.readWriteTx(func(tx *bolt.Tx) error {
		return self.save(tx, config, data, false)
	})
}

func (self *node) save(tx *bolt.Tx, config *structConfig, data interface{}, update bool) error {
	bucket, err := self.CreateBucketIfNotExists(tx, config.Name)
	if err != nil {
		return err
	}
	// save node configuration in the bucket
	meta, err := newMeta(bucket, self)
	if err != nil {
		return err
	}
	if config.ID.IsZero {
		err = meta.increment(config.ID)
		if err != nil {
			return err
		}
	}
	id, err := toBytes(config.ID.Value.Interface(), self.codec)
	if err != nil {
		return err
	}
	for fieldName, fieldConfig := range config.Fields {
		if !update && !fieldConfig.IsID && fieldConfig.Increment && fieldConfig.IsInteger && fieldConfig.IsZero {
			err = meta.increment(fieldConfig)
			if err != nil {
				return err
			}
		}
		if fieldConfig.Index == "" {
			continue
		}
		idx, err := getIndex(bucket, fieldConfig.Index, fieldName)
		if err != nil {
			return err
		}

		if update && fieldConfig.IsZero && !fieldConfig.ForceUpdate {
			continue
		}

		if fieldConfig.IsZero {
			err = idx.RemoveID(id)
			if err != nil {
				return err
			}
			continue
		}

		value, err := toBytes(fieldConfig.Value.Interface(), self.codec)
		if err != nil {
			return err
		}

		var found bool
		idsSaved, err := idx.All(value, nil)
		if err != nil {
			return err
		}
		for _, idSaved := range idsSaved {
			if bytes.Compare(idSaved, id) == 0 {
				found = true
				break
			}
		}
		if found {
			continue
		}
		err = idx.RemoveID(id)
		if err != nil {
			return err
		}
		err = idx.Add(value, id)
		if err != nil {
			return err
		}
	}

	raw, err := self.codec.Marshal(data)
	if err != nil {
		return err
	}

	return bucket.Put(id, raw)
}

func (self *node) Update(data interface{}) error {
	return self.update(data, func(ref *reflect.Value, current *reflect.Value, config *structConfig) error {
		numfield := ref.NumField()
		for i := 0; i < numfield; i++ {
			f := ref.Field(i)
			if ref.Type().Field(i).PkgPath != "" {
				continue
			}
			zero := reflect.Zero(f.Type()).Interface()
			actual := f.Interface()
			if !reflect.DeepEqual(actual, zero) {
				cf := current.Field(i)
				cf.Set(f)
				idxInfo, ok := config.Fields[ref.Type().Field(i).Name]
				if ok {
					idxInfo.Value = &cf
				}
			}
		}
		return nil
	})
}

func (self *node) UpdateField(data interface{}, fieldName string, value interface{}) error {
	return self.update(data, func(ref *reflect.Value, current *reflect.Value, config *structConfig) error {
		f := current.FieldByName(fieldName)
		if !f.IsValid() {
			return errNotFound
		}
		tf, _ := current.Type().FieldByName(fieldName)
		if tf.PkgPath != "" {
			return errNotFound
		}
		v := reflect.ValueOf(value)
		if v.Kind() != f.Kind() {
			return errIncompatibleValue
		}
		f.Set(v)
		idxInfo, ok := config.Fields[fieldName]
		if ok {
			idxInfo.Value = &f
			idxInfo.IsZero = isZero(idxInfo.Value)
			idxInfo.ForceUpdate = true
		}
		return nil
	})
}

func (self *node) update(data interface{}, fn func(*reflect.Value, *reflect.Value, *structConfig) error) error {
	ref := reflect.ValueOf(data)
	if !ref.IsValid() || ref.Kind() != reflect.Ptr || ref.Elem().Kind() != reflect.Struct {
		return errStructPtrNeeded
	}

	config, err := extract(&ref)
	if err != nil {
		return err
	}

	if config.ID.IsZero {
		return errNoID
	}

	current := reflect.New(reflect.Indirect(ref).Type())

	return self.readWriteTx(func(tx *bolt.Tx) error {
		err = self.WithTransaction(tx).One(config.ID.Name, config.ID.Value.Interface(), current.Interface())
		if err != nil {
			return err
		}

		ref := reflect.ValueOf(data).Elem()
		cref := current.Elem()
		err = fn(&ref, &cref, config)
		if err != nil {
			return err
		}

		return self.save(tx, config, current.Interface(), true)
	})
}

func (self *node) Drop(data interface{}) error {
	var bucketName string

	v := reflect.ValueOf(data)
	if v.Kind() != reflect.String {
		info, err := extract(&v)
		if err != nil {
			return err
		}

		bucketName = info.Name
	} else {
		bucketName = v.Interface().(string)
	}

	return self.readWriteTx(func(tx *bolt.Tx) error {
		return self.drop(tx, bucketName)
	})
}

func (self *node) drop(tx *bolt.Tx, bucketName string) error {
	bucket := self.GetBucket(tx)
	if bucket == nil {
		return tx.DeleteBucket([]byte(bucketName))
	}

	return bucket.DeleteBucket([]byte(bucketName))
}

func (self *node) DeleteStruct(data interface{}) error {
	ref := reflect.ValueOf(data)

	if !ref.IsValid() || ref.Kind() != reflect.Ptr || ref.Elem().Kind() != reflect.Struct {
		return errStructPtrNeeded
	}

	config, err := extract(&ref)
	if err != nil {
		return err
	}

	id, err := toBytes(config.ID.Value.Interface(), self.codec)
	if err != nil {
		return err
	}

	return self.readWriteTx(func(tx *bolt.Tx) error {
		return self.deleteStruct(tx, config, id)
	})
}

func (self *node) deleteStruct(tx *bolt.Tx, config *structConfig, id []byte) error {
	bucket := self.GetBucket(tx, config.Name)
	if bucket == nil {
		return errNotFound
	}
	for fieldName, fieldConfig := range config.Fields {
		if fieldConfig.Index == "" {
			continue
		}
		idx, err := getIndex(bucket, fieldConfig.Index, fieldName)
		if err != nil {
			return err
		}
		err = idx.RemoveID(id)
		if err != nil {
			return err
		}
	}
	raw := bucket.Get(id)
	if raw == nil {
		return errNotFound
	}
	return bucket.Delete(id)
}
