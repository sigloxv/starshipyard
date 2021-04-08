package relationship

type Component struct {
	Term  Term
	Value []byte
}

func (self *Component) String() string {
	return string(self)
}
