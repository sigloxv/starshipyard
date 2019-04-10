package shiphold

type Tx interface {
	Commit() error
	Rollback() error
}

func (self node) Begin(writable bool) (Node, error) {
	var err error
	self.tx, err = self.s.Bolt.Begin(writable)
	if err != nil {
		return nil, err
	}
	return &self, nil
}

func (self *node) Rollback() error {
	if self.tx == nil {
		return errNotInTransaction
	}
	return self.tx.Rollback()
}

func (self *node) Commit() error {
	if self.tx == nil {
		return errNotInTransaction
	}
	return self.tx.Commit()
}
